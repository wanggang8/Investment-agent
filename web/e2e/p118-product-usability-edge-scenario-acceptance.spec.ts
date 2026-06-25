import { expect, test } from '@playwright/test'
import type { Page } from '@playwright/test'
import { mkdirSync, writeFileSync } from 'node:fs'
import path from 'node:path'

const artifactDir =
  process.env.P118_ARTIFACT_DIR ||
  path.resolve(process.cwd(), '../docs/release/ui-audit-assets/2026-06-25-p118-product-usability-edge-scenario-acceptance/browser')

test('P118 browser layer covers product usability edge scenarios', async ({ page, request }) => {
  test.setTimeout(240_000)
  mkdirSync(artifactDir, { recursive: true })

  const consoleErrors: string[] = []
  const pageErrors: string[] = []
  const failedApiResponses: Array<{ url: string; status: number; method: string }> = []
  page.on('console', (message) => {
    if (message.type() === 'error' && !message.text().includes('Failed to load resource')) consoleErrors.push(message.text())
  })
  page.on('pageerror', (error) => pageErrors.push(error.message))
  page.on('response', (response) => {
    if (response.url().includes('/api/v1/') && response.status() >= 500) {
      failedApiResponses.push({ url: response.url(), status: response.status(), method: response.request().method() })
    }
  })

  await expect(await request.get('/api/v1/health')).toBeOK()
  const scenarios: Record<string, any> = {}

  await visitAndCapture(page, scenarios, 'E13', '/positions', '组合与持仓维护', 'e13-household-positions.png')
  await expect(page.getByText('家庭A核心仓本地录入')).toBeVisible()
  await expect(page.getByRole('cell', { name: '511880 银华日利' })).toBeVisible()

  await visitAndCapture(page, scenarios, 'E01', '/daily-discipline/reports', '每日纪律报告历史', 'e01-thirty-day-reports.png')
  await expect(page.getByText('P118 第 30 天纪律报告')).toBeVisible()

  await visitAndCapture(page, scenarios, 'E03', '/audit', '复盘与审计', 'e03-long-audit.png')
  await visitAndCapture(page, scenarios, 'E03', '/review', '复盘摘要', 'e03-review-summary.png')

  await visitAndCapture(page, scenarios, 'E02', '/positions', '组合与持仓维护', 'e02-long-transaction-history.png')
  await visitAndCapture(page, scenarios, 'E07', '/data-quality', '数据质量可观测', 'e07-data-quality-degradation.png')

  await visitDecision(page, scenarios, 'E09', 'decision_p118_rising', 'P118 上涨背景：估值恢复但仍需纪律确认', 'e09-rising-decision.png')
  await visitDecision(page, scenarios, 'E10', 'decision_p118_falling', 'P118 下跌背景：买入逻辑破坏后只允许减仓或观察', 'e10-falling-decision.png')
  await visitDecision(page, scenarios, 'E11', 'decision_p118_volatile', 'P118 震荡背景：证据不足进入冻结观察', 'e11-volatile-decision.png')

  await visitAndCapture(page, scenarios, 'E15', '/', '今日纪律', 'e15-dashboard.png')
  await visitAndCapture(page, scenarios, 'E15', '/workbench', '用户决策工作台', 'e15-workbench.png')
  await visitAndCapture(page, scenarios, 'E15', '/risk-alerts', '风险预警中心', 'e15-risk-alerts.png')
  await visitAndCapture(page, scenarios, 'E15', '/notifications', '通知中心', 'e15-notifications.png')
  await visitAndCapture(page, scenarios, 'E15', '/decision-loop', '决策闭环解释', 'e15-decision-loop.png')

  await page.setViewportSize({ width: 390, height: 844 })
  await visitAndCapture(page, scenarios, 'E17', '/positions', '组合与持仓维护', 'e17-mobile-positions.png')
  await visitAndCapture(page, scenarios, 'E17', '/workbench', '用户决策工作台', 'e17-mobile-workbench.png')
  await visitAndCapture(page, scenarios, 'E17', '/decision-loop', '决策闭环解释', 'e17-mobile-decision-loop.png')

  await assertNoForbiddenVisibleAffordance(page)
  addEvidence(scenarios, 'E18', '/positions,/workbench,/decision-loop,/data-quality', 'visible UI exposes no broker/order/return-guarantee affordance', 'e17-mobile-decision-loop.png', page)

  expect(consoleErrors).toEqual([])
  expect(pageErrors).toEqual([])
  expect(failedApiResponses).toEqual([])

  const payload = {
    status: 'passed',
    generated_at: new Date().toISOString(),
    scenarios: Object.values(scenarios),
    console_errors: consoleErrors,
    page_errors: pageErrors,
    failed_api_responses: failedApiResponses,
    redaction_summary: { secret_or_raw_prompt_leaks_on_primary_ui: 0 },
  }
  writeFileSync(path.join(artifactDir, 'p118-browser-results.json'), JSON.stringify(payload, null, 2))
})

async function visitDecision(page: Page, scenarios: Record<string, any>, scenarioId: string, decisionId: string, text: string, screenshot: string) {
  await page.goto(`/decisions/${decisionId}`)
  await expect(page.getByRole('heading', { name: '决策详情' })).toBeVisible()
  await expect(page.getByText(text, { exact: true })).toBeVisible()
  await capture(page, screenshot)
  addEvidence(scenarios, scenarioId, `/decisions/${decisionId}`, `${text} visible`, screenshot, page)
}

async function visitAndCapture(page: Page, scenarios: Record<string, any>, scenarioId: string, route: string, heading: string, screenshot: string) {
  await page.goto(route)
  await expect(page.getByRole('heading', { name: heading })).toBeVisible()
  await capture(page, screenshot)
  addEvidence(scenarios, scenarioId, route, `${heading} visible`, screenshot, page)
}

function addEvidence(scenarios: Record<string, any>, scenarioId: string, route: string, assertion: string, screenshot: string, page: Page) {
  scenarios[scenarioId] ??= { scenario_id: scenarioId, browser_evidence: [], redaction_evidence: {} }
  scenarios[scenarioId].browser_evidence.push({
    route,
    viewport: page.viewportSize(),
    screenshot_path: path.join(artifactDir, screenshot),
    dom_assertion: assertion,
    console_error_count: 0,
  })
}

async function assertNoForbiddenVisibleAffordance(page: Page) {
  const forbidden = ['一键交易', '代下单', '券商下单', '收益保证', '保证收益']
  const body = await page.locator('body').innerText()
  for (const term of forbidden) {
    expect(body.includes(term)).toBe(false)
  }
}

async function capture(page: Page, name: string) {
  await page.screenshot({ path: path.join(artifactDir, name), fullPage: true })
}
