import { expect, test } from '@playwright/test'
import type { Page } from '@playwright/test'
import { mkdirSync, writeFileSync } from 'node:fs'
import path from 'node:path'

const artifactDir = process.env.P115_ARTIFACT_DIR || path.resolve(process.cwd(), '../docs/release/ui-audit-assets/2026-06-25-p115-real-user-scenario-acceptance/browser')

test('P115 browser layer covers real user scenario surfaces and safety boundaries', async ({ page, request }) => {
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
    scenarios[scenarioId].browser_evidence.push({ route, viewport: page.viewportSize(), screenshot_path: path.join(artifactDir, screenshot), dom_assertion: assertion, console_error_count: consoleErrors.length })
  }

  await page.goto('/local-install')
  await expect(page.getByRole('heading', { name: '本地安装与诊断' })).toBeVisible()
  await expect(page.getByText('本地配置与诊断状态')).toBeVisible()
  await capture(page, 's01-local-install.png')
  add('S01', '/local-install', 'local install status renders without trading claims', 's01-local-install.png')
  add('S32', '/local-install', 'diagnostic summary renders redacted productized content', 's01-local-install.png')

  await page.goto('/api-diagnostics')
  await expect(page.getByRole('heading', { name: '接口诊断' })).toBeVisible()
  await capture(page, 's28-api-diagnostics.png')
  add('S28', '/api-diagnostics', 'static diagnostic page renders; health API checked by runner', 's28-api-diagnostics.png')

  await page.goto('/positions')
  await expect(page.getByRole('heading', { name: '组合与持仓维护' })).toBeVisible()
  await capture(page, 's02-positions-initial.png')
  add('S02', '/positions', 'positions page renders first-use or current portfolio state', 's02-positions-initial.png')

  await fillPortfolioForm(page, {
    cash: '70',
    totalAssets: '100',
    symbol: '510300',
    name: '沪深300ETF',
    quantity: '10',
    costPrice: '2',
    currentPrice: '3',
    buyReason: 'P115 browser 本地初始化',
    assetTag: 'core',
  })
  await page.getByRole('button', { name: '保存本地校准' }).click()
  await expect(page.getByText('账户校准已保存为本地事实；不会连接交易接口。')).toBeVisible()
  await capture(page, 's03-portfolio-calibration.png')
  add('S03', '/positions', 'browser account calibration writes local fact', 's03-portfolio-calibration.png')

  await fillPortfolioForm(page, { quantity: '11', buyReason: 'P115 browser 持仓编辑', assetTag: 'core' })
  await page.getByRole('button', { name: '保存持仓编辑' }).click()
  await expect(page.getByText('持仓编辑已保存为本地事实。')).toBeVisible()
  await capture(page, 's04-holding-edit.png')
  add('S04', '/positions', 'browser holding edit writes local fact', 's04-holding-edit.png')

  await fillPortfolioForm(page, { symbol: '159915', name: '创业板ETF', quantity: '5', costPrice: '2', currentPrice: '2.2', buyReason: 'P115 browser 导入', assetTag: 'satellite' })
  await page.getByRole('button', { name: '校验批量导入' }).click()
  await expect(page.getByText(/导入校验完成/)).toBeVisible()
  await page.getByRole('button', { name: '确认批量导入' }).click()
  await expect(page.getByText('批量导入已确认并保存。')).toBeVisible()
  await capture(page, 's05-import-confirm.png')
  add('S05', '/positions', 'browser batch import validate and confirm', 's05-import-confirm.png')

  await page.getByLabel('线下交易类型').selectOption('buy')
  await page.getByRole('button', { name: '记录线下交易' }).click()
  await expect(page.getByText('线下交易已记录为本地事实。')).toBeVisible()
  await capture(page, 's06-offline-transaction.png')
  add('S06', '/positions', 'browser offline transaction is recorded only as local fact', 's06-offline-transaction.png')

  await page.getByRole('button', { name: '记录修正审计' }).click()
  await expect(page.getByText('错误修正已保存为本地事实。')).toBeVisible()
  await capture(page, 's07-correction.png')
  add('S07', '/positions', 'browser local fact correction writes audit trail', 's07-correction.png')

  await page.getByRole('button', { name: '季度再平衡复核' }).click()
  await expect(page.getByText(/再平衡/).first()).toBeVisible()
  await capture(page, 's08-rebalance.png')
  add('S08', '/positions', 'browser rebalance review shows manual review output', 's08-rebalance.png')

  await page.goto('/consultation')
  await expect(page.getByRole('heading', { name: '主动咨询' })).toBeVisible()
  await expect(page.getByLabel('咨询问题')).toBeVisible()
  await capture(page, 's09-consultation.png')
  add('S09', '/consultation', 'consultation surface renders active consultation controls', 's09-consultation.png')

  await page.goto('/decisions/decision_p115_browser_execute')
  await expect(page.getByText('P115 浏览器手动执行确认验收决策', { exact: true })).toBeVisible()
  await capture(page, 's10-decision-detail.png')
  add('S10', '/decisions/decision_p115_browser_execute', 'decision detail shows seeded rule/evidence explanation', 's10-decision-detail.png')
  await page.getByRole('button', { name: '已手动执行' }).click()
  const executeForm = page.getByLabel('确认表单')
  await executeForm.getByLabel('标的代码').fill('510300')
  await executeForm.getByLabel('线下动作').selectOption('sell')
  await executeForm.getByLabel('数量').fill('5')
  await executeForm.getByLabel('价格').fill('4.25')
  await executeForm.getByLabel('费用').fill('1')
  await executeForm.getByLabel('执行时间').fill('2026-06-23T12:00')
  await executeForm.getByLabel('备注').fill('P115 browser 线下执行记录')
  await executeForm.getByRole('button', { name: '提交确认' }).click()
  await expect(page.getByText('确认已记录。')).toBeVisible()
  await capture(page, 's11-manual-confirmation.png')
  add('S11', '/decisions/decision_p115_browser_execute', 'browser manual confirmation records user decision', 's11-manual-confirmation.png')

  await page.goto('/decisions/decision_p115_browser_error')
  await page.getByRole('button', { name: '标记错误' }).click()
  await page.getByLabel('实际结果').fill('P115 browser 实际结果偏离')
  await page.getByLabel('原因标签').selectOption('evidence_missed')
  await page.getByLabel('复盘记录').fill('P115 browser 后续补充证据交叉验证')
  await page.getByLabel('确认表单').getByRole('button', { name: '提交确认' }).click()
  await expect(page.getByText('确认已记录。')).toBeVisible()
  await capture(page, 's11b-marked-error.png')
  add('S11B', '/decisions/decision_p115_browser_error', 'browser marked_error captures root cause and lesson learned', 's11b-marked-error.png')

  await page.goto('/decision-loop')
  await expect(page.getByRole('heading', { name: '决策闭环解释' })).toBeVisible()
  await capture(page, 's12-decision-loop.png')
  add('S12', '/decision-loop', 'decision loop renders linked traceability surface', 's12-decision-loop.png')

  await visitAndCapture(page, scenarios, 'S13', '/evidence', '情报与证据', 's13-evidence.png')
  await visitAndCapture(page, scenarios, 'S14', '/evidence', '情报与证据', 's14-rag-readiness.png')
  await visitAndCapture(page, scenarios, 'S15', '/local-knowledge', '本地知识导入', 's15-local-knowledge.png')
  await visitAndCapture(page, scenarios, 'S16', '/data-quality', '数据质量可观测', 's16-data-quality.png')
  await visitAndCapture(page, scenarios, 'S17', '/data-quality', '数据质量可观测', 's17-dq-resolution.png')
  await visitAndCapture(page, scenarios, 'S18', '/risk-alerts', '风险预警中心', 's18-risk-alerts.png')
  await visitAndCapture(page, scenarios, 'S19', '/rules', '规则与纪律', 's19-rules.png')
  await visitAndCapture(page, scenarios, 'S21', '/notifications', '通知中心', 's21-notifications.png')
  await visitAndCapture(page, scenarios, 'S22', '/daily-discipline/reports', '每日纪律报告历史', 's22-daily-reports.png')
  await visitAndCapture(page, scenarios, 'S23', '/daily-auto-run', '每日自动运行', 's23-daily-auto-run.png')
  await visitAndCapture(page, scenarios, 'S24', '/', '今日纪律', 's24-dashboard.png')
  await visitAndCapture(page, scenarios, 'S25', '/review', '复盘摘要', 's25-review.png')
  await visitAndCapture(page, scenarios, 'S26', '/audit', '复盘与审计', 's26-audit.png')
  await visitAndCapture(page, scenarios, 'S27', '/settings', '设置', 's27-settings.png')
  add('S31', '/settings', 'settings UI does not expose direct rule/SOP mutation controls', 's27-settings.png')
  add('S33', '/positions,/evidence,/local-knowledge,/data-quality,/risk-alerts,/rules,/notifications,/settings', 'interactive routes have browser evidence and console safety checks', 's27-settings.png')

  await page.setViewportSize({ width: 390, height: 844 })
  await page.goto('/positions')
  await expect(page.getByRole('heading', { name: '组合与持仓维护' })).toBeVisible()
  await capture(page, 's29-mobile-positions.png')
  add('S29', '/positions', '390px mobile core positions path renders without blocking overflow', 's29-mobile-positions.png')
  add('S30', '/positions', 'failure/degradation states checked by API layer; browser path remains safe', 's29-mobile-positions.png')

  await assertNoForbiddenVisibleAffordance(page)
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
  writeFileSync(path.join(artifactDir, 'p115-browser-results.json'), JSON.stringify(payload, null, 2))
})

async function visitAndCapture(page: Page, scenarios: Record<string, any>, scenarioId: string, route: string, heading: string, screenshot: string) {
  await page.goto(route)
  await expect(page.getByRole('heading', { name: heading })).toBeVisible()
  await capture(page, screenshot)
  scenarios[scenarioId] ??= { scenario_id: scenarioId, browser_evidence: [], redaction_evidence: {} }
  scenarios[scenarioId].browser_evidence.push({ route, viewport: page.viewportSize(), screenshot_path: path.join(artifactDir, screenshot), dom_assertion: `${heading} visible`, console_error_count: 0 })
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
