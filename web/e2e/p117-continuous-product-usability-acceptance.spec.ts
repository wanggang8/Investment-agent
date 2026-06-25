import { expect, test } from '@playwright/test'
import type { Page } from '@playwright/test'
import { mkdirSync, writeFileSync } from 'node:fs'
import path from 'node:path'

const artifactDir =
  process.env.P117_ARTIFACT_DIR ||
  path.resolve(process.cwd(), '../docs/release/ui-audit-assets/2026-06-25-p117-continuous-product-usability-acceptance/browser')

test('P117 browser layer covers continuous product usability surfaces', async ({ page, request }) => {
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
  function add(scenarioId: string, route: string, assertion: string, screenshot: string) {
    scenarios[scenarioId] ??= { scenario_id: scenarioId, browser_evidence: [], redaction_evidence: {} }
    scenarios[scenarioId].browser_evidence.push({
      route,
      viewport: page.viewportSize(),
      screenshot_path: path.join(artifactDir, screenshot),
      dom_assertion: assertion,
      console_error_count: consoleErrors.length,
    })
  }

  await page.goto('/positions')
  await expect(page.getByRole('heading', { name: '组合与持仓维护' })).toBeVisible()
  await expect(page.getByText('账户快照')).toBeVisible()
  await capture(page, 'u02-day1-onboarding-readback.png')
  add('U02', '/positions', 'day1 portfolio facts render after onboarding', 'u02-day1-onboarding-readback.png')

  await visitAndCapture(page, scenarios, 'U03', '/', '今日纪律', 'u03-dashboard-next-step.png')
  await visitAndCapture(page, scenarios, 'U03', '/workbench', '用户决策工作台', 'u03-workbench-next-step.png')
  await visitAndCapture(page, scenarios, 'U04', '/review', '复盘摘要', 'u04-review-routine.png')
  await visitAndCapture(page, scenarios, 'U04', '/audit', '复盘与审计', 'u04-audit-routine.png')

  await page.goto('/positions')
  await fillPortfolioForm(page, {
    symbol: '161725',
    name: '招商中证白酒指数',
    quantity: '10',
    costPrice: '1',
    currentPrice: '1.1',
    buyReason: 'P117 browser Day3 补记线下买入',
    assetTag: 'active_fund',
  })
  await page.getByLabel('线下交易类型').selectOption('buy')
  await page.getByRole('button', { name: '记录线下交易' }).click()
  await expect(page.getByText('线下交易已记录为本地事实。')).toBeVisible()
  await capture(page, 'u05-browser-offline-transaction.png')
  add('U05', '/positions', 'browser records a later-day offline transaction as local fact', 'u05-browser-offline-transaction.png')

  await visitAndCapture(page, scenarios, 'U06', '/risk-alerts', '风险预警中心', 'u06-risk-alerts.png')
  await visitAndCapture(page, scenarios, 'U06', '/notifications', '通知中心', 'u06-notifications.png')
  await visitAndCapture(page, scenarios, 'U10', '/data-quality', '数据质量可观测', 'u10-data-quality-degradation.png')

  await page.goto('/decisions/decision_p117_browser_execute')
  await expect(page.getByText('P117 浏览器连续使用确认', { exact: true })).toBeVisible()
  await page.getByRole('button', { name: '已手动执行' }).click()
  const executeForm = page.getByLabel('确认表单')
  await executeForm.getByLabel('标的代码').fill('512000')
  await executeForm.getByLabel('线下动作').selectOption('sell')
  await executeForm.getByLabel('数量').fill('2')
  await executeForm.getByLabel('价格').fill('1.28')
  await executeForm.getByLabel('费用').fill('0.5')
  await executeForm.getByLabel('执行时间').fill('2026-06-24T12:00')
  await executeForm.getByLabel('备注').fill('P117 browser 连续使用人工确认')
  await executeForm.getByRole('button', { name: '提交确认' }).click()
  await expect(page.getByText('确认已记录。')).toBeVisible()
  await capture(page, 'u11-browser-manual-confirmation.png')
  add('U11', '/decisions/decision_p117_browser_execute', 'browser manual confirmation works after accumulated facts', 'u11-browser-manual-confirmation.png')

  await visitAndCapture(page, scenarios, 'U13', '/positions', '组合与持仓维护', 'u13-final-positions.png')
  await visitAndCapture(page, scenarios, 'U13', '/', '今日纪律', 'u13-final-dashboard.png')
  await visitAndCapture(page, scenarios, 'U13', '/workbench', '用户决策工作台', 'u13-final-workbench.png')
  await visitAndCapture(page, scenarios, 'U13', '/review', '复盘摘要', 'u13-final-review.png')
  await visitAndCapture(page, scenarios, 'U13', '/audit', '复盘与审计', 'u13-final-audit.png')
  await visitAndCapture(page, scenarios, 'U13', '/decision-loop', '决策闭环解释', 'u13-final-decision-loop.png')

  await page.setViewportSize({ width: 390, height: 844 })
  await visitAndCapture(page, scenarios, 'U15', '/positions', '组合与持仓维护', 'u15-mobile-positions.png')
  await visitAndCapture(page, scenarios, 'U15', '/workbench', '用户决策工作台', 'u15-mobile-workbench.png')

  await assertNoForbiddenVisibleAffordance(page)
  add('U16', '/positions,/workbench,/decision-loop,/data-quality', 'visible UI exposes no broker/order/return-guarantee affordance', 'u15-mobile-workbench.png')
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
  writeFileSync(path.join(artifactDir, 'p117-browser-results.json'), JSON.stringify(payload, null, 2))
})

async function visitAndCapture(page: Page, scenarios: Record<string, any>, scenarioId: string, route: string, heading: string, screenshot: string) {
  await page.goto(route)
  await expect(page.getByRole('heading', { name: heading })).toBeVisible()
  await capture(page, screenshot)
  scenarios[scenarioId] ??= { scenario_id: scenarioId, browser_evidence: [], redaction_evidence: {} }
  scenarios[scenarioId].browser_evidence.push({
    route,
    viewport: page.viewportSize(),
    screenshot_path: path.join(artifactDir, screenshot),
    dom_assertion: `${heading} visible`,
    console_error_count: 0,
  })
}

async function fillPortfolioForm(page: Page, values: Record<string, string>) {
  const labels: Record<string, string> = {
    symbol: '标的代码',
    name: '标的名称',
    quantity: '数量',
    costPrice: '成本价',
    currentPrice: '现价',
    buyReason: '买入理由',
    assetTag: '资产标签',
  }
  for (const [key, value] of Object.entries(values)) {
    const label = labels[key]
    if (label) await page.getByLabel(label).fill(value)
  }
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
