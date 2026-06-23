import { expect, test } from '@playwright/test'
import type { Page } from '@playwright/test'
import { mkdirSync, writeFileSync } from 'node:fs'
import path from 'node:path'

const artifactDir = process.env.P85_ARTIFACT_DIR || path.resolve(process.cwd(), '../docs/release/ui-audit-assets/2026-06-22-p85-expected-return-analysis')
const shouldCapture = process.env.P85_CAPTURE_SCREENSHOTS === '1'

test('P85 expected-return analysis is operated through real UI and read back from API', async ({ page, request }) => {
  test.setTimeout(600_000)
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

  const available = await submitConsultation(page, {
    symbol: '510300',
    question: 'P85：当前涨幅达到目标后是否需要人工止盈评估？',
    previousBaseMidpoint: '20',
    targetReturn: '15',
    screenshot: 'p85-available-consultation.png',
  })
  await expectAvailableExpectedReturn(page, available.decisionId)
  const availableApi = await readDecision(request, available.decisionId)
  expectAvailableAPI(availableApi)

  const downside = await submitConsultation(page, {
    symbol: '159915',
    question: 'P85：当前跌破悲观情景下沿后是否需要核验买入逻辑？',
    screenshot: 'p85-downside-consultation.png',
  })
  await expectDownsideExpectedReturn(page, downside.decisionId)
  const downsideApi = await readDecision(request, downside.decisionId)
  expect(downsideApi.expected_return_scenarios?.sell_evaluation?.triggers ?? []).toContain('downside_lower_bound_breached')

  const unavailable = await submitConsultation(page, {
    symbol: '512000',
    question: 'P85：样本不足时是否仍会生成收益区间？',
    screenshot: 'p85-unavailable-consultation.png',
  })
  await expectUnavailableExpectedReturn(page, unavailable.decisionId)
  const unavailableApi = await readDecision(request, unavailable.decisionId)
  expect(unavailableApi.expected_return_scenarios?.precision_status).toBe('unavailable')
  expect(unavailableApi.expected_return_scenarios?.scenarios ?? []).toHaveLength(0)

  expect(consoleErrors).toEqual([])
  expect(pageErrors).toEqual([])
  expect(failedApiResponses).toEqual([])
  await assertNoForbiddenActionButtons(page)

  writeResult({
    generated_at: new Date().toISOString(),
    status: 'passed',
    decisions: {
      available: summarizeDecision(availableApi),
      downside: summarizeDecision(downsideApi),
      unavailable: summarizeDecision(unavailableApi),
    },
    ui_paths: ['/consultation', `/decisions/${available.decisionId}`, `/decisions/${downside.decisionId}`, `/decisions/${unavailable.decisionId}`],
    safety_boundaries: ['no automatic confirmation', 'no trade execution', 'no broker order', 'no external push'],
    console_errors: consoleErrors,
    page_errors: pageErrors,
    failed_api_responses: failedApiResponses,
  })
})

async function submitConsultation(page: Page, input: { symbol: string; question: string; previousBaseMidpoint?: string; targetReturn?: string; screenshot: string }) {
  await page.goto('/consultation')
  await expect(page.getByRole('heading', { name: '主动咨询' })).toBeVisible()
  await page.getByLabel('标的代码').fill(input.symbol)
  await page.getByLabel('咨询场景').selectOption('hold_review')
  await page.getByLabel('上一轮基准情景中枢（%）').fill(input.previousBaseMidpoint ?? '')
  await page.getByLabel('目标收益率（%）').fill(input.targetReturn ?? '')
  await page.getByLabel('咨询问题').fill(input.question)
  await page.getByRole('button', { name: '提交咨询' }).click()
  await expect(page.getByText('已生成本地决策材料。')).toBeVisible({ timeout: 240_000 })
  const link = page.getByRole('link', { name: '打开生成的决策详情' })
  const href = await link.getAttribute('href')
  expect(href).toMatch(/^\/decisions\/decision_/)
  const decisionId = href?.split('/').pop() ?? ''
  await link.click()
  await expect(page).toHaveURL(new RegExp(`/decisions/${decisionId}$`))
  await capture(page, input.screenshot)
  return { decisionId }
}

async function expectAvailableExpectedReturn(page: Page, decisionId: string) {
  await page.goto(`/decisions/${decisionId}`)
  const card = expectedReturnCard(page)
  await expect(card).toBeVisible()
  await expect(card).toContainText('标的：510300')
  await expect(card).toContainText('当前日期：2026-06-22')
  await expect(card).toContainText('当前价格或净值：3')
  await expect(card).toContainText('PE/PB 分位：20 / 18')
  await expect(card).toContainText('精度状态：可展示区间')
  await expect(card).toContainText('样本数：20')
  await expect(card).toContainText('样本窗口：当前本地持仓、最新市场快照与可用公开净值历史')
  await expect(card).toContainText('筛选条件：基于当前标的持仓成本、最新市场快照和已保存公开市场元数据')
  await expect(card).toContainText('乐观情景：8.00%~15.00%，概率 25.0%')
  await expect(card).toContainText('基准情景：0.00%~8.00%，概率 50.0%')
  await expect(card).toContainText('悲观情景：-12.00%~0.00%，概率 25.0%')
  await expect(card).toContainText('触发因素：upside_lower_bound_reached、base_upper_bound_exceeded、base_midpoint_downshift、target_return_reached')
  await expect(card).toContainText('预期收益仅为情景分析，不构成收益承诺。')
  await expect(card).toContainText('卖出评估仅提示人工复核，不会自动交易。')
  await expect(page.getByText('最终裁决仍以规则链为准')).toBeVisible()
  await capture(page, 'p85-available-detail.png')
}

async function expectDownsideExpectedReturn(page: Page, decisionId: string) {
  await page.goto(`/decisions/${decisionId}`)
  const card = expectedReturnCard(page)
  await expect(card).toContainText('标的：159915')
  await expect(card).toContainText('精度状态：可展示区间')
  await expect(card).toContainText('样本数：20')
  await expect(card).toContainText('触发因素：downside_lower_bound_breached')
  await expect(card).toContainText('当前价格跌破悲观情景下沿，请重新核验买入逻辑')
  await capture(page, 'p85-downside-detail.png')
}

async function expectUnavailableExpectedReturn(page: Page, decisionId: string) {
  await page.goto(`/decisions/${decisionId}`)
  const card = expectedReturnCard(page)
  await expect(card).toContainText('标的：512000')
  await expect(card).toContainText('精度状态：不可用')
  await expect(card).toContainText('样本数：1')
  await expect(card).toContainText('样本过少，仅能给出定性说明')
  await expect(card).toContainText('暂无预期收益情景。')
  await expect(card).toContainText('状态：不适用')
  await expect(card).toContainText('样本过少，无法生成可复核的情景边界')
  await capture(page, 'p85-unavailable-detail.png')
}

function expectedReturnCard(page: Page) {
  return page.locator('article.cockpit-card').filter({ has: page.getByText('预期收益情景', { exact: true }) })
}

async function readDecision(request: any, decisionId: string) {
  const response = await request.get(`/api/v1/decisions/${decisionId}`)
  await expect(response).toBeOK()
  const body = await response.json()
  return body.data
}

function expectAvailableAPI(decision: any) {
  const expected = decision.expected_return_scenarios
  expect(expected.precision_status).toBe('available')
  expect(expected.sample_count).toBe(20)
  expect(expected.sample_window).toBe('当前本地持仓、最新市场快照与可用公开净值历史')
  expect(expected.scenarios).toHaveLength(3)
  expect(expected.scenarios.map((item: any) => item.return_range)).toEqual(['8.00%~15.00%', '0.00%~8.00%', '-12.00%~0.00%'])
  expect(expected.scenarios.map((item: any) => item.probability)).toEqual([0.25, 0.5, 0.25])
  expect(expected.sell_evaluation.triggers).toEqual(['upside_lower_bound_reached', 'base_upper_bound_exceeded', 'base_midpoint_downshift', 'target_return_reached'])
  expect(expected.reassessment_trigger.boundary).toBe('base_midpoint_downshift')
  expect(decision.final_verdict?.status).toBeTruthy()
  expect(decision.user_confirmation?.confirmation_status).toBe('pending')
}

function summarizeDecision(decision: any) {
  const expected = decision.expected_return_scenarios ?? {}
  return {
    decision_id: decision.decision_id,
    symbol: decision.symbol,
    workflow_status: decision.workflow_status,
    final_verdict: decision.final_verdict?.status,
    confirmation_status: decision.user_confirmation?.confirmation_status,
    precision_status: expected.precision_status,
    sample_count: expected.sample_count,
    scenario_count: expected.scenarios?.length ?? 0,
    sell_evaluation_status: expected.sell_evaluation?.status,
    sell_triggers: expected.sell_evaluation?.triggers ?? [],
    retrieval_quality: decision.retrieval_quality,
    analyst_count: decision.analyst_reports?.length ?? 0,
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
