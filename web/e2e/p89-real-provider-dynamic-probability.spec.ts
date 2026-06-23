import { expect, test } from '@playwright/test'
import type { Page } from '@playwright/test'
import { mkdirSync, writeFileSync } from 'node:fs'
import path from 'node:path'

const artifactDir = process.env.P89_ARTIFACT_DIR || path.resolve(process.cwd(), '../docs/release/ui-audit-assets/2026-06-22-p89-real-provider-dynamic-probability')
const shouldCapture = process.env.P89_CAPTURE_SCREENSHOTS === '1'

test('P89 dynamic probability and extreme-fear UI/API paths are read back from product state', async ({ page, request }) => {
  test.setTimeout(420_000)
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

  const baselineId = await runConsultation(page, '510300', 'P89：baseline 历史相似样本预期收益。')
  const baseline = await readDecision(request, baselineId)
  const baselineProbabilities = probabilities(baseline)
  expect(baselineProbabilities).toEqual([0.2, 0.6, 0.2])
  await expectExpectedReturnCard(page, baselineId, ['概率依据：historical_similar_sample_proportion', '乐观情景：8.00%~15.00%，概率 20.0%', '基准情景：0.00%~8.00%，概率 60.0%', '悲观情景：-12.00%~0.00%，概率 20.0%'])
  await capture(page, 'p89-baseline-expected-return.png')

  const dynamicId = await runConsultation(page, '159915', 'P89：估值、基本面和市场状态恶化后重新评估预期收益。')
  const dynamic = await readDecision(request, dynamicId)
  const dynamicProbabilities = probabilities(dynamic)
  expect(dynamicProbabilities[0]).toBeLessThan(baselineProbabilities[0])
  expect(dynamicProbabilities[1]).toBeLessThan(baselineProbabilities[1])
  expect(dynamicProbabilities[2]).toBeGreaterThan(baselineProbabilities[2])
  expect(dynamic.expected_return_scenarios.sell_evaluation.triggers).toEqual(expect.arrayContaining(['scenario_probability_downshift', 'two_month_assumption_downshift', 'one_month_pessimistic_path']))
  await expectExpectedReturnCard(page, dynamicId, ['触发因素：scenario_probability_downshift、two_month_assumption_downshift、one_month_pessimistic_path', '人工提示：估值、基本面或市场状态转弱', '盈利增速：预期 8.00%，实际 1.00%，低于预期 2 个月', '建议动作：复核情景概率、复核核心假设、手动调整情景概率'])
  await capture(page, 'p89-dynamic-downshift.png')

  const extremeId = await runConsultation(page, '600000', 'P89：极端恐惧时展示历史相似场景并暂停主动交易建议。')
  const extreme = await readDecision(request, extremeId)
  expect(extreme.final_verdict.prohibited_actions).toEqual(expect.arrayContaining(['主动交易建议']))
  expect(extreme.expected_return_scenarios.historical_contexts).toHaveLength(1)
  expect(extreme.expected_return_scenarios.sell_evaluation.triggers).toEqual(expect.arrayContaining(['extreme_fear_historical_context']))
  await expectExpectedReturnCard(page, extremeId, ['历史相似场景', '极端恐惧样本：2018Q4, 2020Q1, 2022Q4', '最大回撤 -18.0%', '建议动作：暂停主动交易建议'])
  await expect(page.getByText('禁止事项：主动交易建议')).toBeVisible()
  await capture(page, 'p89-extreme-fear-history.png')

  const providerReadback = await refreshStructuredProviderThroughUI(page, request)
  await capture(page, 'p89-runtime-provider-readback.png')

  expect(consoleErrors).toEqual([])
  expect(pageErrors).toEqual([])
  expect(failedApiResponses).toEqual([])
  await assertNoForbiddenActionButtons(page)

  writeResult({
    generated_at: new Date().toISOString(),
    status: 'passed',
    baseline: summarizeDecision(baseline),
    dynamic: summarizeDecision(dynamic),
    extreme_fear: summarizeDecision(extreme),
    provider_readback: providerReadback,
    ui_paths: ['/consultation', '/settings', `/decisions/${baselineId}`, `/decisions/${dynamicId}`, `/decisions/${extremeId}`],
    safety_boundaries: ['no broker UI', 'no automatic confirmation', 'no order placement', 'no external push', 'no automatic rule application'],
    console_errors: consoleErrors,
    page_errors: pageErrors,
    failed_api_responses: failedApiResponses,
  })
})

async function refreshStructuredProviderThroughUI(page: Page, request: any) {
  await page.goto('/settings')
  await expect(page.getByRole('heading', { name: '设置' })).toBeVisible()
  await page.getByRole('button', { name: '刷新市场数据' }).click()
  await expect(page.getByText('市场刷新完成；只更新本地行情事实和审计记录，不会执行交易。')).toBeVisible({ timeout: 180_000 })
  await expect(page.getByText('P89 结构化字段')).toBeVisible()
  await expect(page.getByText('融资融券：')).toBeVisible()
  await expect(page.getByText('成分财务：')).toBeVisible()
  await expect(page.getByText('资金流向：暂无真实 provider 读回')).toBeVisible()

  const response = await request.get('/api/v1/market/snapshots/latest?symbol=600000')
  await expect(response).toBeOK()
  const body = await response.json()
  const snapshot = body.data
  const structured = snapshot.market_metrics?.metadata?.p88_structured_fields ?? {}
  expect(structured.margin_financing?.date).toBeTruthy()
  expect(Number(structured.margin_financing?.margin_balance)).toBeGreaterThan(0)
  expect(structured.constituent_financial?.disclosure_date).toBeTruthy()
  expect(Number(structured.constituent_financial?.revenue)).toBeGreaterThan(0)
  expect(Number(structured.constituent_financial?.net_profit)).toBeGreaterThan(0)
  expect(structured.capital_flow).toBeUndefined()
  return {
    market_snapshot_id: snapshot.market_snapshot_id,
    symbol: snapshot.symbol,
    trade_date: snapshot.trade_date,
    structured_fields: structured,
  }
}

async function runConsultation(page: Page, symbol: string, question: string) {
  await page.goto('/consultation')
  await expect(page.getByRole('heading', { name: '主动咨询' })).toBeVisible()
  await page.getByLabel('标的代码').fill(symbol)
  await page.getByLabel('咨询场景').selectOption('hold_review')
  await page.getByLabel('咨询问题').fill(question)
  await page.getByRole('button', { name: '提交咨询' }).click()
  await expect(page.getByText('已生成本地决策材料。')).toBeVisible({ timeout: 240_000 })
  const href = await page.getByRole('link', { name: '打开生成的决策详情' }).getAttribute('href')
  const decisionId = href?.split('/').pop() ?? ''
  expect(decisionId).toMatch(/^decision_/)
  return decisionId
}

async function expectExpectedReturnCard(page: Page, decisionId: string, texts: string[]) {
  await page.goto(`/decisions/${decisionId}`)
  const card = page.locator('article.cockpit-card').filter({ has: page.getByText('预期收益情景', { exact: true }) })
  for (const text of texts) {
    await expect(card).toContainText(text)
  }
}

async function readDecision(request: any, decisionId: string) {
  const response = await request.get(`/api/v1/decisions/${decisionId}`)
  await expect(response).toBeOK()
  const body = await response.json()
  return body.data
}

function probabilities(decision: any) {
  return (decision.expected_return_scenarios?.scenarios ?? []).map((item: any) => Number(item.probability?.toFixed(4)))
}

function summarizeDecision(decision: any) {
  const expected = decision.expected_return_scenarios ?? {}
  return {
    decision_id: decision.decision_id,
    symbol: decision.symbol,
    final_verdict: decision.final_verdict?.status,
    prohibited_actions: decision.final_verdict?.prohibited_actions ?? [],
    probability_basis: expected.probability_basis,
    probabilities: probabilities(decision),
    triggers: expected.sell_evaluation?.triggers ?? [],
    assumption_checks: expected.assumption_checks ?? [],
    historical_contexts: expected.historical_contexts ?? [],
  }
}

async function assertNoForbiddenActionButtons(page: Page) {
  for (const name of ['自动交易', '一键交易', '代下单', '券商下单', '自动确认', '自动应用规则']) {
    await expect(page.getByRole('button', { name })).toHaveCount(0)
  }
}

async function capture(page: Page, name: string) {
  if (!shouldCapture) return
  await page.screenshot({ path: path.join(artifactDir, name), fullPage: true })
}

function writeResult(payload: unknown) {
  writeFileSync(path.join(artifactDir, 'browser-results.json'), JSON.stringify(payload, null, 2))
}
