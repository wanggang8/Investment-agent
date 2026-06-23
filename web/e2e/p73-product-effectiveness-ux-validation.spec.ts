import { expect, test } from '@playwright/test'
import type { Page, Response } from '@playwright/test'
import { mkdirSync, writeFileSync } from 'node:fs'
import path from 'node:path'

type UXTaskResult = {
  id: string
  goal: string
  route: string
  status: 'passed' | 'blocked'
  evidence: string[]
  severity?: 'critical' | 'major' | 'minor'
}

test('P73 validates product goal effectiveness and real UX task comprehension', async ({ page, request }) => {
  test.setTimeout(240_000)

  const health = await request.get('/api/v1/health')
  await expect(health).toBeOK()

  const consoleErrors: string[] = []
  const pageErrors: string[] = []
  const unexpectedFailedApiResponses: Array<{ url: string; status: number; method: string }> = []
  page.on('pageerror', (error) => pageErrors.push(error.stack || error.message))
  page.on('console', (message) => {
    if (
      message.type() === 'error'
      && !message.text().includes('status of 409 (Conflict)')
      && !message.text().includes('status of 404 (Not Found)')
      && !message.text().includes('Failed to load resource: the server responded with a status of')
    ) {
      consoleErrors.push(message.text())
    }
  })
  page.on('response', (response) => {
    if (response.status() >= 400 && response.url().includes('/api/v1/') && classifyFailedApiResponse(response) === 'unexpected') {
      unexpectedFailedApiResponses.push({ url: redactURL(response.url()), status: response.status(), method: response.request().method() })
    }
  })

  const results: UXTaskResult[] = []
  results.push(await validateDailyDisciplineGoal(page))
  results.push(await validatePortfolioGoal(page))
  results.push(await validateEvidenceAndDataQualityGoal(page))
  results.push(await validateDecisionTraceabilityGoal(page))
  results.push(await validateBackgroundOnlyBlockingGoal(page))
  results.push(await validateManualConfirmationGoal(page))
  results.push(await validateRiskReviewAndRuleEffectGoal(page))
  results.push(await validateMobileGoal(page))
  results.push(await validateUnsafeInputGoal(page))

  await assertNoForbiddenAffordance(page)
  expect(consoleErrors).toEqual([])
  expect(pageErrors).toEqual([])
  expect(unexpectedFailedApiResponses).toEqual([])
  expect(results.filter((item) => item.status !== 'passed')).toEqual([])
  writeResults({
    generated_at: new Date().toISOString(),
    status: 'passed',
    ux_tasks: results,
    console_errors: consoleErrors,
    page_errors: pageErrors,
    unexpected_failed_api_responses: unexpectedFailedApiResponses,
    critical_findings: [],
    accepted_gaps: [
      'P73 validates product-goal and UX effectiveness with deterministic local facts; it does not claim future investment returns.',
      'P73 relies on P71/P72 for live public-source and real LLM availability evidence.',
    ],
  })
})

async function validateDailyDisciplineGoal(page: Page): Promise<UXTaskResult> {
  await page.goto('/')
  await expect(page.getByRole('heading', { name: '今日纪律', exact: true })).toBeVisible()
  await expect(page.getByText('今日纪律报告', { exact: true })).toBeVisible()
  await expect(page.getByText('P32 smoke 今日纪律报告已生成')).toBeVisible()
  await expect(page.getByRole('navigation', { name: '主导航' }).getByRole('link', { name: '决策闭环' })).toHaveAttribute('href', '/decision-loop')
  await expect(page.getByRole('link', { name: '查看决策详情' })).toBeVisible()
  await capture(page, 'daily-discipline-goal.png')
  return passed('daily_discipline_goal', '用户能从首页识别今日纪律状态、决策详情入口和决策闭环导航。', '/', ['今日纪律报告', '查看决策详情', '决策闭环'])
}

async function validatePortfolioGoal(page: Page): Promise<UXTaskResult> {
  await page.goto('/positions')
  await expect(page.getByRole('heading', { name: '组合与持仓维护' })).toBeVisible()
  await expect(page.getByText('组合维护状态', { exact: true })).toBeVisible()
  await page.getByLabel('现金').fill('70')
  await page.getByLabel('总资产').fill('100')
  await page.getByLabel('标的代码').fill('510300')
  await page.getByLabel('标的名称').fill('沪深300ETF')
  await page.getByLabel('数量').fill('10')
  await page.getByLabel('成本价').fill('2')
  await page.getByLabel('现价').fill('3')
  await page.getByLabel('买入理由').fill('P73 UX 验收：本地纪律账户初始化')
  await page.getByLabel('资产标签').fill('长期配置')
  await page.getByRole('button', { name: '保存本地校准' }).click()
  await expect(page.getByText('账户校准已保存为本地事实；不会连接交易接口。')).toBeVisible()
  await expect(page.getByRole('cell', { name: '510300 沪深300ETF' })).toBeVisible()
  await capture(page, 'portfolio-goal.png')
  return passed('portfolio_goal', '用户能录入本地持仓事实，并看到系统不会连接交易接口。', '/positions', ['组合维护状态', '不会连接交易接口'])
}

async function validateEvidenceAndDataQualityGoal(page: Page): Promise<UXTaskResult> {
  await page.goto('/data-quality')
  await expect(page.getByRole('heading', { name: '数据质量可观测' })).toBeVisible()
  await expect(page.getByText('数据质量总览', { exact: true })).toBeVisible()
  await expect(page.getByText('证据与 RAG 信号')).toBeVisible()
  await expect(page.getByLabel('数据质量下一步').getByRole('link', { name: '查看证据' })).toHaveAttribute('href', '/evidence')
  await page.goto('/evidence')
  await expect(page.getByRole('heading', { name: '情报与证据' })).toBeVisible()
  await expect(page.getByText('P30 smoke 证据摘要')).toBeVisible()
  await expect(page.getByText('P73 背景材料')).toBeVisible()
  await capture(page, 'evidence-data-quality-goal.png')
  return passed('evidence_data_quality_goal', '用户能从数据质量跳到证据页，并识别正式证据与背景材料。', '/data-quality -> /evidence', ['证据与 RAG 信号', 'P73 背景材料'])
}

async function validateDecisionTraceabilityGoal(page: Page): Promise<UXTaskResult> {
  await page.goto('/decisions/decision_smoke_p30')
  await expect(page.getByRole('heading', { name: '决策详情' })).toBeVisible()
  await expect(page.getByText('决策故事', { exact: true })).toBeVisible()
  await expect(page.getByText('Agent 分析材料')).toBeVisible()
  await expect(page.getByText('最终裁决明细')).toBeVisible()
  await expect(page.getByText('裁决链', { exact: true })).toBeVisible()
  await expect(page.getByText('预期收益情景')).toBeVisible()
  await expect(page.getByText('卖出评估仅用于人工复核，不会自动交易。')).toBeVisible()
  await capture(page, 'decision-traceability-goal.png')
  return passed('decision_traceability_goal', '用户能区分 LLM 分析、证据、规则裁决和收益情景边界。', '/decisions/decision_smoke_p30', ['Agent 分析材料', '最终裁决明细', '裁决链', '不会自动交易'])
}

async function validateBackgroundOnlyBlockingGoal(page: Page): Promise<UXTaskResult> {
  await page.goto('/decisions/decision_smoke_p73_background_only')
  await expect(page.getByRole('heading', { name: '仅有背景材料，不能生成交易类建议' })).toBeVisible()
  await expect(page.getByText('工作流状态：已完成')).toBeVisible()
  await expect(page.getByText('最终裁决：仅有背景材料，不能生成交易类建议')).toBeVisible()
  await expect(page.getByText('状态：数据不足')).toBeVisible()
  await expect(page.getByText('当前确认状态：无需确认')).toBeVisible()
  await expect(page.getByText('C 级背景材料不得作为正式裁决依据').first()).toBeVisible()
  await expect(page.getByRole('button', { name: '已手动执行' })).toHaveCount(0)
  await capture(page, 'background-only-blocking-goal.png')
  return passed('background_only_blocking_goal', '用户能看懂 C 级背景材料不能生成交易类建议，且没有执行确认入口。', '/decisions/decision_smoke_p73_background_only', ['数据不足', '无需确认'])
}

async function validateManualConfirmationGoal(page: Page): Promise<UXTaskResult> {
  await page.goto('/decisions/decision_smoke_p30')
  await expect(page.getByRole('heading', { name: '决策详情' })).toBeVisible()
  await page.getByRole('button', { name: '记录计划' }).click()
  await page.getByLabel('确认表单').getByLabel('备注').fill('P73 UX 验收：只记录线下计划，不更新账户')
  await page.getByRole('button', { name: '提交确认' }).click()
  await expect(page.getByText('确认已记录。')).toBeVisible()
  await expect(page.getByText('当前确认状态：已记录计划')).toBeVisible()
  await expect(page.getByText('系统只记录你的线下动作，不会替你买入或卖出。')).toBeVisible()
  await capture(page, 'manual-confirmation-goal.png')
  return passed('manual_confirmation_goal', '用户能记录线下计划，并理解系统不替用户交易。', '/decisions/decision_smoke_p30', ['已记录计划', '不会替你买入或卖出'])
}

async function validateRiskReviewAndRuleEffectGoal(page: Page): Promise<UXTaskResult> {
  await page.goto('/risk-alerts')
  await expect(page.getByRole('heading', { name: '风险预警中心' })).toBeVisible()
  await expect(page.getByText('P39 source health stale 触发数据降级风险')).toBeVisible()
  await expect(page.getByText('禁止动作：自动交易、外部推送')).toBeVisible()
  await page.goto('/review')
  await expect(page.getByRole('heading', { name: '复盘摘要' })).toBeVisible()
  await expect(page.getByText('只读追踪')).toBeVisible()
  await expect(page.getByLabel('追踪关联风险预警')).toContainText('risk_smoke_p39')
  await page.goto('/rules')
  await expect(page.getByRole('heading', { name: '规则与纪律' })).toBeVisible()
  await expect(page.getByText('P39 E2E 规则提案')).toBeVisible()
  await expect(page.getByText('验证状态：已通过')).toBeVisible()
  await expect(page.getByText('守门人通过后仍需用户最终确认，正式规则不会自动生效。')).toBeVisible()
  await capture(page, 'risk-review-rule-effect-goal.png')
  return passed('risk_review_rule_effect_goal', '用户能从风险进入复盘和规则效果验证，并看到规则不会自动应用。', '/risk-alerts -> /review -> /rules', ['只读追踪', '验证状态：已通过'])
}

async function validateMobileGoal(page: Page): Promise<UXTaskResult> {
  await page.setViewportSize({ width: 390, height: 844 })
  for (const target of [
    { path: '/', readyText: '今日纪律报告' },
    { path: '/workbench', heading: '用户决策工作台', readyText: '主动咨询入口' },
    { path: '/data-quality', heading: '数据质量可观测', readyText: '数据质量总览' },
    { path: '/decisions/decision_smoke_p73_background_only', heading: '仅有背景材料，不能生成交易类建议', readyText: 'C 级背景材料不得作为正式裁决依据' },
  ]) {
    await page.goto(target.path)
    await assertPrimaryNavigationReachable(page)
    if (target.heading) await expect(page.getByRole('heading', { name: target.heading })).toBeVisible()
    await expect(page.getByText(target.readyText, { exact: false }).first()).toBeVisible()
    await expectPageNotHorizontallyOverflowed(page)
  }
  await capture(page, 'mobile-goal.png')
  await page.setViewportSize({ width: 1280, height: 900 })
  return passed('mobile_goal', '390px 视口下核心目标页面无水平溢出且导航可达。', 'mobile core routes', ['主导航', '无水平溢出'])
}

async function validateUnsafeInputGoal(page: Page): Promise<UXTaskResult> {
  await page.goto('/consultation')
  await expect(page.getByRole('heading', { name: '主动咨询' })).toBeVisible()
  await page.getByRole('button', { name: '提交咨询' }).click()
  await expect(page.getByText('请填写咨询问题和标的代码。')).toBeVisible()
  await capture(page, 'unsafe-input-goal.png')
  return passed('unsafe_input_goal', '用户提交不完整咨询时收到明确阻断，不生成伪建议。', '/consultation', ['请填写咨询问题和标的代码'])
}

function passed(id: string, goal: string, route: string, evidence: string[]): UXTaskResult {
  return { id, goal, route, status: 'passed', evidence }
}

async function capture(page: Page, name: string) {
  if (process.env.P73_CAPTURE_SCREENSHOTS !== '1') return
  const artifactDir = process.env.P73_ARTIFACT_DIR || path.resolve(process.cwd(), '..', 'docs', 'release', 'ui-audit-assets', '2026-06-19-p73')
  mkdirSync(artifactDir, { recursive: true })
  await page.screenshot({ path: path.join(artifactDir, name), fullPage: true })
}

function writeResults(payload: unknown) {
  const artifactDir = process.env.P73_ARTIFACT_DIR || path.resolve(process.cwd(), '..', 'docs', 'release', 'ui-audit-assets', '2026-06-19-p73')
  mkdirSync(artifactDir, { recursive: true })
  writeFileSync(path.join(artifactDir, 'browser-results.json'), `${JSON.stringify(payload, null, 2)}\n`)
}

function classifyFailedApiResponse(response: Response) {
  const method = response.request().method()
  const url = response.url()
  if (method === 'GET' && response.status() === 404 && url.includes('/api/v1/portfolio/current')) return 'expected'
  if (method === 'GET' && response.status() === 409 && url.includes('/api/v1/dashboard/today')) return 'expected'
  if (method === 'POST' && response.status() === 409 && url.includes('/api/v1/decisions/consult')) return 'expected'
  return 'unexpected'
}

function redactURL(value: string) {
  return value.replace(/[?&](api_key|token|key|authorization)=[^&]+/gi, '$1=REDACTED')
}

async function assertPrimaryNavigationReachable(page: Page) {
  const isMobile = await page.evaluate(() => window.innerWidth <= 760)
  if (!isMobile) {
    await expect(page.getByRole('navigation', { name: '主导航' })).toBeVisible()
    return
  }
  const toggle = page.getByRole('button', { name: '导航' })
  await expect(toggle).toBeVisible()
  await toggle.click()
  await expect(page.getByRole('navigation', { name: '主导航' })).toBeVisible()
  await toggle.click()
}

async function expectPageNotHorizontallyOverflowed(page: Page) {
  const overflow = await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)
  expect(overflow).toBeLessThanOrEqual(2)
}

async function assertNoForbiddenAffordance(page: Page) {
  const forbidden = /自动下单|一键交易|代下单|券商下单|券商接口|自动规则应用|自动应用规则|自动确认|自动修复|外部推送|短信|邮件|Webhook|第三方推送|收益承诺|完整密钥|API key|sk-|SELECT \* FROM|\/Users\/private|prompt:/
  for (const pathName of ['/', '/workbench', '/positions', '/data-quality', '/evidence', '/risk-alerts', '/review', '/rules', '/audit', '/notifications', '/decisions/decision_smoke_p73_background_only']) {
    await page.goto(pathName)
    await expect(page.getByRole('button', { name: forbidden })).toHaveCount(0)
    await expect(page.getByRole('link', { name: forbidden })).toHaveCount(0)
  }
}
