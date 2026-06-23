import { expect, test } from '@playwright/test'
import type { Page } from '@playwright/test'
import { mkdirSync, writeFileSync } from 'node:fs'
import path from 'node:path'

const artifactDir = process.env.P88_ARTIFACT_DIR || path.resolve(process.cwd(), '../docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers')
const shouldCapture = process.env.P88_CAPTURE_SCREENSHOTS === '1'

test('P88 remaining blocker UI paths exercise expected return, rebalance, and SOP proposal', async ({ page, request }) => {
  test.setTimeout(420_000)
  mkdirSync(artifactDir, { recursive: true })

  const consoleErrors: string[] = []
  const pageErrors: string[] = []
  const failedApiResponses: Array<{ url: string; status: number; method: string }> = []
  page.on('console', (message) => {
    if (message.type() === 'error' && !message.text().includes('Failed to load resource')) {
      consoleErrors.push(message.text())
    }
  })
  page.on('pageerror', (error) => pageErrors.push(error.message))
  page.on('response', (response) => {
    if (response.url().includes('/api/v1/') && response.status() >= 500) {
      failedApiResponses.push({ url: response.url(), status: response.status(), method: response.request().method() })
    }
  })

  await expect(await request.get('/api/v1/health')).toBeOK()

  const decisionId = await runExpectedReturnConsultation(page)
  const decision = await readDecision(request, decisionId)
  expectExpectedReturnAPI(decision)
  const sourceTransitions = await runSourceVerifiedTransitions(page, request)

  const rebalance = await runRebalanceReview(page, request)
  const proposal = await runSOPProposal(page, request)

  expect(consoleErrors).toEqual([])
  expect(pageErrors).toEqual([])
  expect(failedApiResponses).toEqual([])
  await assertNoForbiddenActionButtons(page)

  writeResult({
    generated_at: new Date().toISOString(),
    status: 'passed',
    decision: summarizeDecision(decision),
    source_transitions: sourceTransitions,
    rebalance,
    proposal,
    ui_paths: ['/consultation', `/decisions/${decisionId}`, ...sourceTransitions.map((item) => `/decisions/${item.decision_id}`), '/positions', '/rules'],
    safety_boundaries: ['no broker UI', 'no automatic confirmation', 'no order placement', 'no external push', 'no automatic rule application'],
    console_errors: consoleErrors,
    page_errors: pageErrors,
    failed_api_responses: failedApiResponses,
  })
})

async function runExpectedReturnConsultation(page: Page) {
  await page.goto('/consultation')
  await expect(page.getByRole('heading', { name: '主动咨询' })).toBeVisible()
  await page.getByLabel('标的代码').fill('510300')
  await page.getByLabel('咨询场景').selectOption('hold_review')
  await page.getByLabel('咨询问题').fill('P88：请用历史相似样本分析未来 12 个月预期收益。')
  await page.getByRole('button', { name: '提交咨询' }).click()
  await expect(page.getByText('已生成本地决策材料。')).toBeVisible({ timeout: 240_000 })
  const href = await page.getByRole('link', { name: '打开生成的决策详情' }).getAttribute('href')
  const decisionId = href?.split('/').pop() ?? ''
  expect(decisionId).toMatch(/^decision_/)
  await page.goto(`/decisions/${decisionId}`)
  const card = expectedReturnCard(page)
  await expect(card).toContainText('标的：510300')
  await expect(card).toContainText('标的名称：沪深300ETF（510300）')
  await expect(card).toContainText('持仓类别：宽基指数 ETF')
  await expect(card).toContainText('分析期限：未来 12 个月')
  await expect(card).toContainText('概率依据：historical_similar_sample_proportion')
  await expect(card).toContainText('乐观情景：12.00%~18.00%，概率 20.0%')
  await expect(card).toContainText('基准情景：4.00%~9.00%，概率 60.0%')
  await expect(card).toContainText('悲观情景：-10.00%~-2.00%，概率 20.0%')
  await expect(card).toContainText('金融成分股路径：600000，covered')
  await expect(card).toContainText('预期收益仅为情景分析，不构成收益承诺。')
  await capture(page, 'p88-expected-return-detail.png')
  return decisionId
}

async function runSourceVerifiedTransitions(page: Page, request: any) {
  const sellOnlyDecisionId = await runConsultation(page, '159915', 'P88：请根据已核验的买入逻辑破坏证据判断是否进入只卖不买。')
  await page.goto(`/decisions/${sellOnlyDecisionId}`)
  await expect(page.getByRole('heading', { name: '决策详情' })).toBeVisible()
  await expect(page.getByRole('heading', { name: '买入逻辑破坏，只卖不买；A/S 独立信源=2' })).toBeVisible()
  await expect(page.getByText('禁止事项：新增买入、加仓')).toBeVisible()
  await expect(page.getByText(/买入逻辑破坏：危险/)).toBeVisible()
  await expect(page.getByText('样本数：2')).toBeVisible()
  await expect(page.getByText('样本过少，仅能给出定性说明')).toBeVisible()
  await expect(page.getByText('需补充数据：market_history、valuation_percentiles、fundamental_growth、formal_evidence')).toBeVisible()
  await expect(page.getByText('暂无预期收益情景。')).toBeVisible()
  await capture(page, 'p88-source-transition-sell-only.png')
  const sellOnly = await readDecision(request, sellOnlyDecisionId)
  expect(sellOnly.final_verdict.status).toBe('sell_only')
  expect(sellOnly.final_verdict.prohibited_actions).toEqual(expect.arrayContaining(['新增买入', '加仓']))

  const frozenDecisionId = await runConsultation(page, '600000', 'P88：请根据单一 A 级重大负面事件证据判断是否冻结观察。')
  await page.goto(`/decisions/${frozenDecisionId}`)
  await expect(page.getByRole('heading', { name: '决策详情' })).toBeVisible()
  await expect(page.getByRole('heading', { name: '重大事件缺少 2 个 A/S 独立信源；当前 A/S 独立信源=1' })).toBeVisible()
  await expect(page.getByText('禁止事项：主动交易建议')).toBeVisible()
  await expect(page.getByText(/多源验证：预警/)).toBeVisible()
  await capture(page, 'p88-source-transition-frozen-watch.png')
  const frozenWatch = await readDecision(request, frozenDecisionId)
  expect(frozenWatch.final_verdict.status).toBe('frozen_watch')
  expect(frozenWatch.final_verdict.prohibited_actions).toEqual(expect.arrayContaining(['主动交易建议']))

  return [
    summarizeTransitionDecision(sellOnly, 'source_verified_buy_logic_break_sell_only'),
    summarizeTransitionDecision(frozenWatch, 'single_high_grade_major_event_frozen_watch'),
  ]
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

async function runRebalanceReview(page: Page, request: any) {
  await page.goto('/positions')
  await expect(page.getByRole('heading', { name: '组合与持仓维护', exact: true })).toBeVisible()
  await page.getByRole('button', { name: '运行季度再平衡复核' }).click()
  await expect(page.getByText('季度再平衡复核已生成，仅作为人工计划。')).toBeVisible()
  await expect(page.getByText(/现金：目标 30.00%，实际/)).toBeVisible()
  await expect(page.getByText('季度再平衡仅生成人工计划金额，不连接券商、不自动交易、不创建订单。')).toBeVisible()
  await capture(page, 'p88-rebalance-review.png')

  const response = await request.get('/api/v1/audit-events')
  await expect(response).toBeOK()
  const body = await response.json()
  const items = body.data?.items ?? []
  const audit = items.find((item: any) => item.input_ref_type === 'rebalance_review')
  expect(audit).toBeTruthy()
  return { audit_event_id: audit.audit_event_id, action: audit.action, input_ref_type: audit.input_ref_type }
}

async function runSOPProposal(page: Page, request: any) {
  await page.goto('/rules')
  await expect(page.getByRole('heading', { name: '规则与纪律' })).toBeVisible()
  await page.getByRole('button', { name: '生成 SOP 补充提案' }).click()
  await expect(page.getByText('SOP 补充提案已生成，等待人工确认。')).toBeVisible()
  await expect(page.getByText('SOP 补充提案：连续流动性缺口未覆盖')).toBeVisible()
  await capture(page, 'p88-sop-proposal.png')

  const proposalsResponse = await request.get('/api/v1/rule-proposals')
  await expect(proposalsResponse).toBeOK()
  const proposalsBody = await proposalsResponse.json()
  const proposal = (proposalsBody.data?.items ?? []).find((item: any) => item.proposal_type === 'sop')
  expect(proposal).toBeTruthy()
  expect(proposal.status).toBe('pending_user_confirm')

  const notificationsResponse = await request.get('/api/v1/notifications')
  await expect(notificationsResponse).toBeOK()
  const notificationsBody = await notificationsResponse.json()
  const notification = (notificationsBody.data?.items ?? []).find((item: any) => item.source_id === proposal.proposal_id)
  expect(notification).toBeTruthy()
  return { proposal_id: proposal.proposal_id, status: proposal.status, notification_id: notification.notification_id }
}

async function readDecision(request: any, decisionId: string) {
  const response = await request.get(`/api/v1/decisions/${decisionId}`)
  await expect(response).toBeOK()
  const body = await response.json()
  return body.data
}

function expectExpectedReturnAPI(decision: any) {
  const expected = decision.expected_return_scenarios
  expect(expected.target_name).toBe('沪深300ETF')
  expect(expected.target_code).toBe('510300')
  expect(expected.horizon_label).toBe('未来 12 个月')
  expect(expected.probability_basis).toBe('historical_similar_sample_proportion')
  expect(expected.scenarios.map((item: any) => item.probability)).toEqual([0.2, 0.6, 0.2])
  expect(expected.holding_class_coverage.map((item: any) => item.holding_class)).toContain('equity_constituent_financial')
}

function expectedReturnCard(page: Page) {
  return page.locator('article.cockpit-card').filter({ has: page.getByText('预期收益情景', { exact: true }) })
}

function summarizeDecision(decision: any) {
  const expected = decision.expected_return_scenarios ?? {}
  return {
    decision_id: decision.decision_id,
    symbol: decision.symbol,
    precision_status: expected.precision_status,
    target_name: expected.target_name,
    target_code: expected.target_code,
    probability_basis: expected.probability_basis,
    probabilities: expected.scenarios?.map((item: any) => item.probability) ?? [],
    covered_holding_classes: expected.holding_class_coverage?.map((item: any) => item.holding_class) ?? [],
  }
}

function summarizeTransitionDecision(decision: any, scenario: string) {
  return {
    scenario,
    decision_id: decision.decision_id,
    symbol: decision.symbol,
    final_verdict: decision.final_verdict?.status,
    display_text: decision.final_verdict?.display_text,
    prohibited_actions: decision.final_verdict?.prohibited_actions ?? [],
    triggered_rules: decision.triggered_rules?.map((item: any) => item.rule_id) ?? [],
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
