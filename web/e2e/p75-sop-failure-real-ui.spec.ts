import { expect, test } from '@playwright/test'
import type { Page } from '@playwright/test'
import { mkdirSync, writeFileSync } from 'node:fs'
import path from 'node:path'

const sopScenarios = [
  { id: 'SOP-A', summary: 'P75 SOP-A 持仓下跌超过5%', action: '记录继续观察', expected: '观察中' },
  { id: 'SOP-B', summary: 'P75 SOP-B 持仓上涨超过20%', action: '记录升级复核', expected: '已升级' },
  { id: 'SOP-C', summary: 'P75 SOP-C 热点追涨冲动', action: '记录继续观察', expected: '观察中' },
  { id: 'SOP-D', summary: 'P75 SOP-D 恐慌清仓语言', action: '记录升级复核', expected: '已升级' },
  { id: 'SOP-E', summary: 'P75 SOP-E 宏观灰犀牛', action: '记录本地解除预警', expected: '已解除' },
  { id: 'SOP-F', summary: 'P75 SOP-F 黑天鹅事件', action: '记录升级复核', expected: '已升级' },
]

test('P75 SOP A-F and failure states are verified through real browser UI operations', async ({ page, request }) => {
  test.setTimeout(180_000)

  const health = await request.get('/api/v1/health')
  await expect(health).toBeOK()

  const consoleErrors: string[] = []
  const pageErrors: string[] = []
  page.on('pageerror', (error) => pageErrors.push(error.stack || error.message))
  page.on('console', (message) => {
    if (message.type() === 'error' && !message.text().includes('Failed to load resource: the server responded with a status of')) {
      consoleErrors.push(message.text())
    }
  })

  await runSOPRiskJourney(page)
  await runFailureStateJourney(page)
  await runDecisionErrorMarkingJourney(page)
  await runRuleGatekeeperJourney(page)
  await runMobileLayoutJourney(page)

  expect(consoleErrors).toEqual([])
  expect(pageErrors).toEqual([])
  writeResults({
    generated_at: new Date().toISOString(),
    status: 'passed',
    sop_scenarios: sopScenarios.map(({ id, summary, expected }) => ({ id, summary, expected })),
    failure_states: ['unsupported_symbol', 'insufficient_data', 'stale_data', 'degraded_source', 'model_unavailable', 'validation_error', 'gatekeeper_deny', 'gatekeeper_user_review'],
  })
})

async function runSOPRiskJourney(page: Page) {
  await page.goto('/risk-alerts')
  await expect(page.getByRole('heading', { name: '风险预警中心' })).toBeVisible()
  await expect(page.getByLabel('风险处置队列')).toContainText('影响标的')

  for (const scenario of sopScenarios) {
    const card = page.locator('article.risk-alert-card').filter({ hasText: scenario.summary })
    await expect(card).toBeVisible()
    await expect(card).toContainText(`SOP：${scenario.id}`)
    await expect(card).toContainText('数据前提：')
    await expect(card).toContainText('LLM 角色：')
    await expect(card).toContainText('禁止动作')
    await expect(card).toContainText('建议人工动作')
    await expect(card).toContainText('不会自动交易')
    await card.getByRole('button', { name: scenario.action }).click()
    await expect(page.locator('article.risk-alert-card').filter({ hasText: scenario.summary })).toContainText(scenario.expected)
  }
  await capture(page, 'p75-sop-risk-alerts-after-ui-actions.png')

  await page.goto('/audit')
  await expect(page.getByRole('heading', { name: '复盘与审计' })).toBeVisible()
  for (const scenario of sopScenarios) {
    await expect(page.getByText(scenario.summary).first()).toBeVisible()
  }
  await capture(page, 'p75-sop-audit-readback.png')
}

async function runFailureStateJourney(page: Page) {
  await page.goto('/consultation')
  await expect(page.getByRole('heading', { name: '主动咨询' })).toBeVisible()
  await page.getByRole('button', { name: '提交咨询' }).click()
  await expect(page.getByText('请填写咨询问题和标的代码。')).toBeVisible()

  await page.goto('/data-quality?symbol=999999')
  await expect(page.getByRole('heading', { name: '数据质量可观测' })).toBeVisible()
  await expect(page.getByText(/999999/).first()).toBeVisible()
  await expect(page.getByText(/能力圈|不支持|暂无|降级|缺失/).first()).toBeVisible()

  await page.goto('/decisions/decision_p75_insufficient')
  await expect(page.getByRole('heading', { name: '决策详情' })).toBeVisible()
  await expect(page.getByText('最终裁决：证据不足，暂停交易类建议')).toBeVisible()
  await expect(page.getByText('状态：数据不足')).toBeVisible()
  await expect(page.getByRole('heading', { name: '证据不足，暂停交易类建议' })).toBeVisible()

  await page.goto('/decisions/decision_p75_model_unavailable')
  await expect(page.getByRole('heading', { name: '决策详情' })).toBeVisible()
  await expect(page.getByText('分析服务暂不可用')).toBeVisible()
  await expect(page.getByRole('heading', { name: 'LLM 降级，暂停交易类建议' })).toBeVisible()

  await page.goto('/settings')
  await expect(page.getByRole('heading', { name: '设置' })).toBeVisible()
  await expect(page.getByText('csindex_p75 · 指数估值文件 · 过期')).toBeVisible()
  await expect(page.getByText('eastmoney_p75 · 资金流向 · 解析失败')).toBeVisible()
  await capture(page, 'p75-failure-states.png')
}

async function runRuleGatekeeperJourney(page: Page) {
  await page.goto('/rules')
  await expect(page.getByRole('heading', { name: '规则与纪律' })).toBeVisible()
  const sendCard = page.locator('section.proposal-item').filter({ hasText: 'P75 UI 送审提案' })
  await expect(sendCard).toBeVisible()
  await sendCard.getByRole('button', { name: '确认送审' }).click()
  await expect(page.locator('section.proposal-item').filter({ hasText: 'P75 UI 送审提案' })).toContainText('守门人结果：审计通过')
  if (process.env.P75_FINAL_RULE_APPLY === '1') {
    const finalCard = page.locator('section.proposal-item').filter({ hasText: 'P75 UI 送审提案' })
    await expect(finalCard.getByRole('button', { name: '确认应用到正式规则' })).toBeVisible()
    await finalCard.getByRole('button', { name: '确认应用到正式规则' }).click()
    await expect(page.locator('section.proposal-item').filter({ hasText: 'P75 UI 送审提案' })).toContainText('状态：已应用')
    await page.reload()
    await expect(page.getByRole('heading', { name: '规则与纪律' })).toBeVisible()
    await expect(page.getByText(/当前规则库：v_prop_p75_ui_send_gatekeeper/)).toBeVisible()
  } else {
    await expect(page.locator('section.proposal-item').filter({ hasText: 'P75 UI 送审提案' })).toContainText('状态：待最终确认')
  }
  await expect(page.getByText('P75 守门人否决样例')).toBeVisible()
  await expect(page.getByText('守门人结果：审计否决')).toBeVisible()
  await expect(page.getByText('P75 守门人用户复核样例')).toBeVisible()
  await expect(page.getByText('守门人结果：需要用户复核')).toBeVisible()
  await expect(page.getByText('验证状态：样本不足')).toBeVisible()
  await expect(page.getByText('门禁结论：需要用户复核')).toBeVisible()
  await capture(page, 'p75-gatekeeper-deny-user-review.png')
}

async function runDecisionErrorMarkingJourney(page: Page) {
  await page.goto('/decisions/decision_p75_mark_error')
  await expect(page.getByRole('heading', { name: '决策详情' })).toBeVisible()
  await expect(page.getByText('当前确认状态：待确认')).toBeVisible()
  await page.getByRole('button', { name: '标记错误' }).click()
  const form = page.locator('[aria-label="确认表单"]')
  await form.getByLabel('实际结果').fill('P75 真实 UI 标记错误：后续结果与原判断偏离')
  await form.getByLabel('原因标签').selectOption('rule_threshold_issue')
  await form.getByLabel('复盘记录').fill('P75 真实 UI 标记错误后应写入 error_cases 并可在审计中读回。')
  await form.getByLabel('备注').fill('P75 critical mutation mark_error UI acceptance')
  await page.getByRole('button', { name: '提交确认' }).click()
  await expect(page.getByText('确认已记录。')).toBeVisible()
  await expect(page.getByText('当前确认状态：已标记错误')).toBeVisible()

  await page.goto('/review')
  await expect(page.getByRole('heading', { name: '复盘摘要' })).toBeVisible()
  await expect(page.getByText('rule_threshold_issue').first()).toBeVisible()

  await page.goto('/audit')
  await expect(page.getByRole('heading', { name: '复盘与审计' })).toBeVisible()
  await expect(page.getByText('decision_p75_mark_error').first()).toBeVisible()
  await capture(page, 'p75-mark-error-readback.png')
}

async function runMobileLayoutJourney(page: Page) {
  await page.setViewportSize({ width: 390, height: 844 })
  for (const target of [
    { path: '/risk-alerts', heading: '风险预警中心', ready: 'P75 SOP-F 黑天鹅事件' },
    { path: '/rules', heading: '规则与纪律', ready: 'P75 守门人用户复核样例' },
    { path: '/settings', heading: '设置', ready: 'P40 数据源健康' },
    { path: '/decisions/decision_p75_insufficient', heading: '决策详情', ready: '证据不足，暂停交易类建议' },
  ]) {
    await page.goto(target.path)
    await expect(page.getByRole('heading', { name: target.heading })).toBeVisible()
    await expect(page.getByText(target.ready).first()).toBeVisible()
    await assertNoHorizontalOverflow(page)
    await capture(page, `mobile-${target.path.replaceAll('/', '_') || 'home'}.png`)
  }
}

async function assertNoHorizontalOverflow(page: Page) {
  const overflow = await page.evaluate(() => {
    const body = document.body
    const doc = document.documentElement
    return Math.max(body.scrollWidth, doc.scrollWidth) - Math.max(body.clientWidth, doc.clientWidth)
  })
  expect(overflow).toBeLessThanOrEqual(2)
}

async function capture(page: Page, fileName: string) {
  if (process.env.P75_CAPTURE_SCREENSHOTS !== '1') return
  mkdirSync(artifactDir(), { recursive: true })
  await page.screenshot({ path: path.join(artifactDir(), fileName), fullPage: true })
}

function artifactDir() {
  return process.env.P75_ARTIFACT_DIR ?? path.resolve(process.cwd(), '../tmp/p75-sop-failure-real-ui')
}

function writeResults(payload: unknown) {
  mkdirSync(artifactDir(), { recursive: true })
  writeFileSync(path.join(artifactDir(), 'browser-results.json'), `${JSON.stringify(payload, null, 2)}\n`)
}
