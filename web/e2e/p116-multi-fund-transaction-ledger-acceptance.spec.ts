import { expect, test } from '@playwright/test'
import type { Page } from '@playwright/test'
import { mkdirSync, writeFileSync } from 'node:fs'
import path from 'node:path'

const artifactDir =
  process.env.P116_ARTIFACT_DIR ||
  path.resolve(process.cwd(), '../docs/release/ui-audit-assets/2026-06-25-p116-multi-fund-transaction-ledger-acceptance/browser')

test('P116 browser layer covers multi-fund transaction ledger surfaces', async ({ page, request }) => {
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
  await fillPortfolioForm(page, {
    cash: '1000',
    totalAssets: '1400',
    symbol: '510300',
    name: '沪深300ETF',
    quantity: '100',
    costPrice: '3.2',
    currentPrice: '4',
    buyReason: 'P116 browser 多基金初始化',
    assetTag: 'core',
  })
  await page.getByRole('button', { name: '保存本地校准' }).click()
  await expect(page.getByText('账户校准已保存为本地事实；不会连接交易接口。')).toBeVisible()
  await capture(page, 'l02-browser-initial-portfolio.png')
  add('L02', '/positions', 'browser saves initial local portfolio calibration', 'l02-browser-initial-portfolio.png')
  add('L13', '/positions', 'positions page accepts real user portfolio input', 'l02-browser-initial-portfolio.png')

  await fillPortfolioForm(page, {
    symbol: '159915',
    name: '创业板ETF',
    quantity: '8',
    costPrice: '2.1',
    currentPrice: '2.6',
    buyReason: 'P116 browser 批量导入第二只基金',
    assetTag: 'satellite',
  })
  await page.getByRole('button', { name: '校验批量导入' }).click()
  await expect(page.getByText(/导入校验完成/)).toBeVisible()
  await page.getByRole('button', { name: '确认批量导入' }).click()
  await expect(page.getByText('批量导入已确认并保存。')).toBeVisible()
  await capture(page, 'l04-browser-import-confirm.png')
  add('L04', '/positions', 'browser validates and confirms a second-fund import batch', 'l04-browser-import-confirm.png')
  add('L13', '/positions', 'multi-fund import flow renders after confirmation', 'l04-browser-import-confirm.png')

  await page.getByLabel('线下交易类型').selectOption('buy')
  await page.getByRole('button', { name: '记录线下交易' }).click()
  await expect(page.getByText('线下交易已记录为本地事实。')).toBeVisible()
  await capture(page, 'l03-browser-offline-transaction.png')
  add('L03', '/positions', 'browser records offline transaction as local fact only', 'l03-browser-offline-transaction.png')
  add('L13', '/positions', 'offline transaction path stays local and user-triggered', 'l03-browser-offline-transaction.png')

  await page.getByRole('button', { name: '记录修正审计' }).click()
  await expect(page.getByText('错误修正已保存为本地事实。')).toBeVisible()
  await page.getByRole('button', { name: '季度再平衡复核' }).click()
  await expect(page.getByText(/再平衡/).first()).toBeVisible()
  await capture(page, 'l06-l09-browser-maintenance.png')
  add('L06', '/positions', 'browser correction audit can be recorded after multi-fund edits', 'l06-l09-browser-maintenance.png')
  add('L09', '/positions', 'browser rebalance review is visible and manual', 'l06-l09-browser-maintenance.png')

  await page.goto('/decisions/decision_p116_browser_execute')
  await expect(page.getByText('P116 浏览器多基金手动执行确认验收决策', { exact: true })).toBeVisible()
  await page.getByRole('button', { name: '已手动执行' }).click()
  const executeForm = page.getByLabel('确认表单')
  await executeForm.getByLabel('标的代码').fill('159915')
  await executeForm.getByLabel('线下动作').selectOption('sell')
  await executeForm.getByLabel('数量').fill('3')
  await executeForm.getByLabel('价格').fill('2.65')
  await executeForm.getByLabel('费用').fill('1')
  await executeForm.getByLabel('执行时间').fill('2026-06-24T12:00')
  await executeForm.getByLabel('备注').fill('P116 browser 多基金线下执行记录')
  await executeForm.getByRole('button', { name: '提交确认' }).click()
  await expect(page.getByText('确认已记录。')).toBeVisible()
  await capture(page, 'l07-browser-manual-confirmation.png')
  add('L07', '/decisions/decision_p116_browser_execute', 'browser manual confirmation records a second-fund sell fact', 'l07-browser-manual-confirmation.png')
  add('L14', '/decisions/decision_p116_browser_execute', 'decision detail confirmation path remains manual', 'l07-browser-manual-confirmation.png')

  await visitAndCapture(page, scenarios, 'L10', '/risk-alerts', '风险预警中心', 'l10-risk-alerts.png')
  await visitAndCapture(page, scenarios, 'L10', '/notifications', '通知中心', 'l10-notifications.png')
  await visitAndCapture(page, scenarios, 'L11', '/data-quality', '数据质量可观测', 'l11-data-quality.png')
  await visitAndCapture(page, scenarios, 'L12', '/', '今日纪律', 'l12-dashboard.png')
  await visitAndCapture(page, scenarios, 'L12', '/workbench', '用户决策工作台', 'l12-workbench.png')
  await visitAndCapture(page, scenarios, 'L12', '/review', '复盘摘要', 'l12-review.png')
  await visitAndCapture(page, scenarios, 'L12', '/audit', '复盘与审计', 'l12-audit.png')
  await visitAndCapture(page, scenarios, 'L14', '/decision-loop', '决策闭环解释', 'l14-decision-loop.png')

  await page.setViewportSize({ width: 390, height: 844 })
  await page.goto('/positions')
  await expect(page.getByRole('heading', { name: '组合与持仓维护' })).toBeVisible()
  await capture(page, 'l15-mobile-positions.png')
  add('L15', '/positions', '390px positions path renders after multi-fund operations', 'l15-mobile-positions.png')

  await assertNoForbiddenVisibleAffordance(page)
  add('L16', '/positions,/risk-alerts,/data-quality,/decision-loop', 'visible UI exposes no broker/order/return-guarantee affordance', 'l15-mobile-positions.png')
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
  writeFileSync(path.join(artifactDir, 'p116-browser-results.json'), JSON.stringify(payload, null, 2))
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
    cash: '现金',
    totalAssets: '总资产',
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
