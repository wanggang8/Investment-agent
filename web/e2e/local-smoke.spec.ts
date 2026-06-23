import { expect, test } from '@playwright/test'
import type { Page } from '@playwright/test'

test('local server and critical read-only pages render smoke data', async ({ page, request }) => {
  const health = await request.get('/api/v1/health')
  await expect(health).toBeOK()
  await expect(await health.json()).toMatchObject({ status: 'ok' })

  const errors: string[] = []
  page.on('pageerror', (error) => errors.push(error.message))
  page.on('console', (message) => {
    if (message.type() === 'error' && !message.text().includes('status of 409 (Conflict)') && !message.text().includes('status of 404 (Not Found)')) {
      errors.push(message.text())
    }
  })

  await page.goto('/decisions/decision_smoke_p30')
  await expect(page.getByRole('heading', { name: '决策详情' })).toBeVisible()
  await expect(page.getByText('P30 本地 E2E smoke 决策', { exact: true }).first()).toBeVisible()
  await expect(page.getByText('预期收益情景')).toBeVisible()
  await expect(page.getByText('动态卖出评估')).toBeVisible()
  await expect(page.getByText('卖出评估仅用于人工复核，不会自动交易。')).toBeVisible()
  await expect(page.getByText('样本不足，不展示精确概率')).toBeVisible()

  await page.goto('/evidence')
  await expect(page.getByRole('heading', { name: '情报与证据' })).toBeVisible()
  const p30EvidenceRow = page.getByRole('row').filter({ hasText: 'P30SmokeSource' })
  await expect(p30EvidenceRow.getByText('A 级')).toBeVisible()
  await expect(p30EvidenceRow.getByText('正式证据')).toBeVisible()
  await expect(p30EvidenceRow.getByText('P30 smoke 证据摘要')).toBeVisible()

  await page.goto('/audit')
  await expect(page.getByRole('heading', { name: '复盘与审计' })).toBeVisible()
  await expect(page.getByText('审计检查状态')).toBeVisible()
  await expect(page.getByLabel('审计下一步')).toBeVisible()
  const p30AuditItem = page.getByRole('listitem').filter({ hasText: 'P30SmokeSeed' })
  await expect(p30AuditItem.getByText('运行本地任务')).toBeVisible()
  await p30AuditItem.getByRole('button', { name: '展开引用' }).click()
  await expect(p30AuditItem.getByText('p30-real-e2e-smoke').first()).toBeVisible()

  await page.goto('/')
  await expect(page.getByRole('heading', { name: '今日纪律', exact: true })).toBeVisible()
  await expect(page.getByText('今日纪律报告', { exact: true })).toBeVisible()
  await expect(page.getByText('P32 smoke 今日纪律报告已生成')).toBeVisible()

  await page.goto('/workbench')
  await expect(page.getByRole('heading', { name: '用户决策工作台' })).toBeVisible()
  await expect(page.getByText('今日先看')).toBeVisible()
  const workbenchRegion = page.getByRole('region', { name: '用户决策工作台区域' })
  await expect(workbenchRegion.getByText('组合与风险')).toBeVisible()
  await expect(workbenchRegion.getByText('规则与复盘')).toBeVisible()
  await expect(workbenchRegion.getByText('主动咨询入口')).toBeVisible()
  await expect(page.getByRole('link', { name: '查看决策闭环' })).toHaveAttribute('href', '/decision-loop')
  await expect(page.getByRole('link', { name: '发起主动咨询' })).toHaveAttribute('href', '/consultation')

  await page.goto('/decision-loop')
  await expect(page.getByRole('heading', { name: '决策闭环解释' })).toBeVisible()
  await expect(page.getByText('只读解释链，仅展示本地事实和导航，不改变事实状态。')).toBeVisible()
  await expect(page.getByText('闭环概览')).toBeVisible()

  await page.goto('/data-quality')
  await expect(page.getByRole('heading', { name: '数据质量可观测' })).toBeVisible()
  await expect(page.getByText('数据质量总览', { exact: true })).toBeVisible()
  await expect(page.getByText('数据源健康 信号')).toBeVisible()
  await expect(page.getByText('证据与 RAG 信号')).toBeVisible()
  await expect(page.getByText('数据源健康', { exact: true })).toBeVisible()
  await expect(page.getByText('证据与检索', { exact: true })).toBeVisible()
  await expect(page.getByText('LLM 质量', { exact: true })).toBeVisible()
  await expect(page.getByText('影响范围与下一步', { exact: true })).toBeVisible()
  await expect(page.getByLabel('数据质量下一步').getByRole('link', { name: '查看证据' })).toHaveAttribute('href', '/evidence')

  await page.goto('/daily-discipline/reports')
  await expect(page.getByRole('heading', { name: '每日纪律报告历史' })).toBeVisible()
  await expect(page.getByText('每日纪律复盘状态')).toBeVisible()
  await expect(page.getByText('P32 smoke 今日纪律报告已生成')).toBeVisible()
  await page.getByRole('link', { name: '查看报告', exact: true }).click()

  await expect(page.getByRole('heading', { name: '每日纪律报告详情' })).toBeVisible()
  await expect(page.getByText('P32 smoke 今日纪律报告已生成')).toBeVisible()
  await expect(page.getByText('不会自动执行交易')).toBeVisible()

  await page.goto('/daily-auto-run')
  await expect(page.getByRole('heading', { name: '每日自动运行' })).toBeVisible()
  await expect(page.getByText('每日自动运行健康')).toBeVisible()
  await expect(page.locator('header').getByText('失败', { exact: true })).toBeVisible()
  await expect(page.getByText('缺少本地账户或持仓。')).toBeVisible()
  await expect(page.getByText('仅记录本地每日刷新和纪律评估结果，不会自动执行交易，需人工复核。')).toBeVisible()

  await page.goto('/positions')
  await expect(page.getByRole('heading', { name: '组合与持仓维护' })).toBeVisible()
  await expect(page.getByText('组合维护状态', { exact: true })).toBeVisible()
  await expect(page.getByText('首次使用引导')).toBeVisible()
  await page.goto('/local-install')
  await expect(page.getByRole('heading', { name: '本地安装与诊断' })).toBeVisible()
  await expect(page.getByText('本地配置与诊断状态')).toBeVisible()
  await expect(page.getByText('关键命令')).toBeVisible()
  await expect(page.getByText('启动草稿')).toBeVisible()
  await expect(page.getByText(/该页用于本地安装/)).toBeVisible()

  await page.goto('/local-knowledge')
  await expect(page.getByRole('heading', { name: '本地知识导入' })).toBeVisible()
  await expect(page.getByText('本地配置与诊断状态')).toBeVisible()
  const localKnowledgeRegion = page.getByLabel('本地知识导入区域')
  await expect(localKnowledgeRegion.getByText('脱敏预览', { exact: true }).first()).toBeVisible()
  await expect(localKnowledgeRegion.getByText('索引计划', { exact: true }).first()).toBeVisible()

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
  await page.getByLabel('买入理由').fill('P33 smoke 本地初始化')
  await page.getByLabel('资产标签').fill('长期配置')
  await page.getByRole('button', { name: '保存本地校准' }).click()
  await expect(page.getByText('账户校准已保存为本地事实；不会连接交易接口。')).toBeVisible()
  await expect(page.getByText('P33 smoke 本地初始化')).toBeVisible()
  await expect(page.getByRole('cell', { name: '510300 沪深300ETF' })).toBeVisible()

  expect(errors).toEqual([])
})

test('P39 full local user journey covers review, rules, risk, settings, and mobile safety', async ({ page, request }) => {
  test.setTimeout(180_000)

  const health = await request.get('/api/v1/health')
  await expect(health).toBeOK()

  const errors = captureUnexpectedErrors(page)

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
  await page.getByLabel('买入理由').fill('P39 E2E 本地初始化')
  await page.getByLabel('资产标签').fill('长期配置')
  await page.getByRole('button', { name: '保存本地校准' }).click()
  await expect(page.getByText('账户校准已保存为本地事实；不会连接交易接口。')).toBeVisible()

  await page.goto('/consultation')
  await page.getByLabel('咨询问题').fill('P39 E2E 主动咨询：510300 是否继续持有？')
  await page.getByLabel('标的代码').fill('510300')
  await page.getByRole('button', { name: '提交咨询' }).click()
  await page.getByText('决策故事', { exact: true }).waitFor({ timeout: 150_000 })
  await expect(page.getByRole('link', { name: '打开生成的决策详情' })).toBeVisible()
  await expect(page.getByText('以下内容仅作为分析材料，最终裁决仍以规则链为准。')).toBeVisible()
  await page.getByRole('button', { name: '记录计划' }).click()
  await page.getByLabel('确认表单').getByLabel('备注').fill('P39 E2E 仅记录线下计划')
  await page.getByRole('button', { name: '提交确认' }).click()
  await expect(page.getByText('确认已记录。')).toBeVisible()
  await expect(page.getByText('当前确认状态：已记录计划')).toBeVisible()
  await expect(page.getByText('系统只记录你的线下动作，不会替你买入或卖出。')).toBeVisible()

  await page.goto('/decisions/decision_smoke_p30')
  await expect(page.getByRole('heading', { name: '决策详情' })).toBeVisible()
  await expect(page.getByRole('heading', { name: '检索质量' })).toBeVisible()
  await expect(page.getByText('分析服务暂不可用', { exact: false })).toHaveCount(0)

  await page.goto('/decisions/decision_smoke_p39_out_of_scope')
  await expect(page.getByRole('heading', { name: '能力圈外，拒绝交易类建议' })).toBeVisible()
  await expect(page.getByText('能力圈检查：能力圈外 P39 fixture excluded symbol')).toBeVisible()
  await expect(page.getByText('能力圈外，不生成收益判断。')).toBeVisible()
  await expect(page.getByText('检索状态：无结果')).toBeVisible()

  await page.goto('/decisions/decision_smoke_p39_llm_degraded')
  await expect(page.getByRole('heading', { name: 'LLM 降级，暂停交易类建议' })).toBeVisible()
  await expect(page.getByText('工作流状态：降级')).toBeVisible()
  await expect(page.getByText('分析服务暂不可用，页面仅展示规则与已有数据。')).toBeVisible()
  await expect(page.getByText('LLM 降级时不展示精确收益概率。')).toBeVisible()

  const missingMarket = await request.post('/api/v1/decisions/consult', {
    data: {
      question: 'P39 E2E 缺市场降级检查',
      symbol: 'NO_MARKET',
      scenario: 'hold_review',
    },
  })
  expect(missingMarket.status()).toBe(409)
  const missingMarketBody = await missingMarket.json() as { error?: { code?: string } }
  expect(missingMarketBody.error?.code).toBe('DATA_REQUIRED')

  await page.goto('/review')
  await expect(page.getByRole('heading', { name: '复盘摘要' })).toBeVisible()
  await expect(page.getByText('季度阈值复盘')).toBeVisible()
  await expect(page.getByText('规则变更仍需守门人审计和用户最终确认，不会自动应用。')).toBeVisible()
  await expect(page.getByText('只读追踪')).toBeVisible()
  await expect(page.getByLabel('追踪关联风险预警')).toContainText('risk_smoke_p39')

  await page.goto('/rules')
  await expect(page.getByRole('heading', { name: '规则与纪律' })).toBeVisible()
  await expect(page.getByText('规则治理状态')).toBeVisible()
  await expect(page.getByText('P39 E2E 规则提案')).toBeVisible()
  await expect(page.getByText('守门人结果：审计通过')).toBeVisible()
  await expect(page.getByText('验证状态：已通过')).toBeVisible()
  await expect(page.getByText('守门人通过后仍需用户最终确认，正式规则不会自动生效。')).toBeVisible()

  await page.goto('/risk-alerts')
  await expect(page.getByRole('heading', { name: '风险预警中心' })).toBeVisible()
  await expect(page.getByText('风险处置队列', { exact: true })).toBeVisible()
  await expect(page.getByText('待看队列')).toBeVisible()
  await expect(page.getByText('处理中队列')).toBeVisible()
  await expect(page.getByText('需复盘队列')).toBeVisible()
  await expect(page.getByText('P39 source health stale 触发数据降级风险')).toBeVisible()
  await expect(page.getByText('禁止动作：自动交易、外部推送')).toBeVisible()
  await expect(page.getByRole('link', { name: '关联决策' })).toBeVisible()

  await page.goto('/settings')
  await expect(page.getByRole('heading', { name: '设置' })).toBeVisible()
  await expect(page.getByText('本地配置与诊断状态')).toBeVisible()
  await expect(page.getByText('csindex_extended · 指数估值文件 · 过期')).toBeVisible()
  await expect(page.getByText('仅刷新行情事实与审计事件，不连接交易接口。')).toBeVisible()

  await page.setViewportSize({ width: 390, height: 844 })
  await page.goto('/workbench')
  await assertPrimaryNavigationReachable(page)
  await expect(page.getByRole('heading', { name: '用户决策工作台' })).toBeVisible()
  const consultationEntry = page.getByRole('link', { name: '发起主动咨询' })
  await consultationEntry.scrollIntoViewIfNeeded()
  await expect(consultationEntry).toBeInViewport()
  await expectPageNotHorizontallyOverflowed(page)

  await page.goto('/data-quality')
  await assertPrimaryNavigationReachable(page)
  await expect(page.getByRole('heading', { name: '数据质量可观测' })).toBeVisible()
  await expect(page.getByText('数据质量总览', { exact: true })).toBeVisible()
  const evidenceEntry = page.getByLabel('数据质量可观测区域').getByRole('link', { name: '查看证据' })
  await evidenceEntry.scrollIntoViewIfNeeded()
  await expect(evidenceEntry).toBeInViewport()
  await expectPageNotHorizontallyOverflowed(page)

  for (const target of [
    { path: '/consultation', heading: '主动咨询' },
    { path: '/decisions/decision_smoke_p30', heading: '决策详情', readyText: '决策故事' },
    { path: '/evidence', heading: '情报与证据', readyText: '证据可信度' },
    { path: '/decision-loop', heading: '决策闭环解释', readyText: '只读决策生命周期' },
    { path: '/rules', heading: '规则与纪律', readyText: '规则治理状态' },
    { path: '/audit', heading: '复盘与审计', readyText: '审计检查状态' },
    { path: '/notifications', heading: '通知中心', readyText: '本地通知收件箱' },
    { path: '/daily-discipline/reports', heading: '每日纪律报告历史', readyText: '每日纪律复盘状态' },
    { path: '/daily-auto-run', heading: '每日自动运行', readyText: '每日自动运行健康' },
    { path: '/local-install', heading: '本地安装与诊断', readyText: '本地配置与诊断状态' },
    { path: '/local-knowledge', heading: '本地知识导入', readyText: '本地配置与诊断状态' },
    { path: '/settings', heading: '设置', readyText: '本地配置与诊断状态' },
  ]) {
    await page.goto(target.path)
    await assertPrimaryNavigationReachable(page)
    await expect(page.getByRole('heading', { name: target.heading })).toBeVisible()
    if (target.readyText) {
      await expect(page.getByText(target.readyText, { exact: false }).first()).toBeVisible()
    }
    await expectPageNotHorizontallyOverflowed(page)
  }

  await page.goto('/positions')
  await assertPrimaryNavigationReachable(page)
  await expect(page.getByRole('heading', { name: '组合与持仓维护' })).toBeVisible()
  await expect(page.getByText('组合维护状态', { exact: true })).toBeVisible()
  await expectPageNotHorizontallyOverflowed(page)

  await page.goto('/risk-alerts')
  await assertPrimaryNavigationReachable(page)
  await expect(page.getByRole('button', { name: '记录继续观察' })).toBeVisible()
  const pageTitle = page.getByRole('heading', { name: '风险预警中心' })
  await expect(page.getByText('风险处置队列', { exact: true })).toBeVisible()
  const observeButton = page.getByRole('button', { name: '记录继续观察' })
  await expect(pageTitle).toBeInViewport()
  await observeButton.scrollIntoViewIfNeeded()
  await expect(observeButton).toBeInViewport()

  await assertNoForbiddenTradingAffordanceAcrossKeyPages(page, [
    { path: '/', readyText: '今日纪律报告' },
    { path: '/workbench', heading: '用户决策工作台' },
    { path: '/decision-loop', heading: '决策闭环解释', readyText: '闭环概览', scanBody: true },
    { path: '/data-quality', heading: '数据质量可观测', readyText: '数据质量总览', scanBody: true },
    { path: '/positions', heading: '组合与持仓维护', readyText: '组合维护状态' },
    { path: '/consultation', readyLabel: '咨询问题' },
    { path: '/decisions/decision_smoke_p30', readyText: 'P30 本地 E2E smoke 决策' },
    { path: '/decisions/decision_smoke_p39_out_of_scope', readyText: '能力圈外，拒绝交易类建议' },
    { path: '/decisions/decision_smoke_p39_llm_degraded', readyText: 'LLM 降级，暂停交易类建议' },
    { path: '/evidence', readyText: 'P30 smoke 证据摘要' },
    { path: '/daily-discipline/reports', readyText: 'P32 smoke 今日纪律报告已生成' },
    { path: '/daily-auto-run', heading: '每日自动运行', readyText: '每日自动运行健康', scanBody: true },
    { path: '/review', readyText: '只读追踪' },
    { path: '/rules', readyText: 'P39 E2E 规则提案', scanBody: true },
    { path: '/audit', heading: '复盘与审计', readyText: '审计检查状态', scanBody: true },
    { path: '/notifications', heading: '通知中心', readyText: '本地通知收件箱', scanBody: true },
    { path: '/risk-alerts', readyText: 'P39 source health stale 触发数据降级风险' },
    { path: '/settings', readyText: 'csindex_extended · 指数估值文件 · 过期', scanBody: true },
    { path: '/local-install', heading: '本地安装与诊断', readyText: '启动草稿', scanBody: true },
    { path: '/local-knowledge', heading: '本地知识导入', readyText: '脱敏预览', scanBody: true },
  ])
  expect(errors).toEqual([])
})

test('P62 design system gates cover keyboard paths and three viewport reflow', async ({ page, request }) => {
  test.setTimeout(180_000)

  const health = await request.get('/api/v1/health')
  await expect(health).toBeOK()

  const errors = captureUnexpectedErrors(page)

  for (const viewport of [
    { width: 390, height: 844 },
    { width: 768, height: 900 },
    { width: 1280, height: 900 },
  ]) {
    await page.setViewportSize(viewport)
    for (const target of [
      { path: '/', readyText: '今日纪律报告' },
      { path: '/workbench', heading: '用户决策工作台', readyText: '数据可信度' },
      { path: '/positions', heading: '组合与持仓维护', readyText: '组合维护状态' },
      { path: '/data-quality', heading: '数据质量可观测', readyText: '数据质量总览' },
      { path: '/risk-alerts', heading: '风险预警中心', readyText: '风险处置队列' },
      { path: '/rules', heading: '规则与纪律', readyText: '规则治理状态' },
      { path: '/audit', heading: '复盘与审计', readyText: '审计检查状态' },
      { path: '/notifications', heading: '通知中心', readyText: '本地通知收件箱' },
      { path: '/local-install', heading: '本地安装与诊断', readyText: '本地配置与诊断状态' },
      { path: '/local-knowledge', heading: '本地知识导入', readyText: '本地配置与诊断状态' },
      { path: '/settings', heading: '设置', readyText: '本地配置与诊断状态' },
    ]) {
      await page.goto(target.path)
      if (target.heading) {
        await expect(page.getByRole('heading', { name: target.heading })).toBeVisible()
      }
      if (target.readyText) {
        await expect(page.getByText(target.readyText, { exact: false }).first()).toBeVisible()
      }
      await expectPageNotHorizontallyOverflowed(page)
    }
  }

  await page.setViewportSize({ width: 390, height: 844 })
  await page.goto('/')
  await assertMobileNavigationKeyboardPath(page)

  await page.setViewportSize({ width: 1280, height: 900 })
  await page.goto('/')
  await assertDesktopNavigationKeyboardPath(page)

  await page.setViewportSize({ width: 768, height: 900 })
  await page.goto('/positions')
  await expect(page.getByLabel('现金')).toBeVisible()
  await page.getByLabel('现金').focus()
  await page.keyboard.insertText('88')
  await expect(page.getByLabel('现金')).toHaveValue('88')
  await page.keyboard.press('Tab')
  await expect(page.getByLabel('总资产')).toBeFocused()

  await page.goto('/audit')
  const p30AuditItem = page.getByRole('listitem').filter({ hasText: 'P30SmokeSeed' })
  const auditToggle = p30AuditItem.getByRole('button', { name: '展开引用' })
  await auditToggle.focus()
  await auditToggle.press('Enter')
  await expect(p30AuditItem.getByRole('button', { name: '收起引用' })).toHaveAttribute('aria-expanded', 'true')
  await expect(p30AuditItem.getByText('p30-real-e2e-smoke').first()).toBeVisible()

  await page.goto('/local-knowledge')
  const validateButton = page.getByRole('button', { name: '校验预览' })
  await validateButton.focus()
  await page.keyboard.press('Enter')
  await expect(page.getByLabel('本地知识导入区域').getByText('批次', { exact: true })).toBeVisible()
  await expect(page.getByLabel('本地知识导入区域').getByText('索引计划', { exact: true })).toBeVisible()

  expect(errors).toEqual([])
})

function captureUnexpectedErrors(page: Page) {
  const errors: string[] = []
  page.on('pageerror', (error) => errors.push(error.message))
  page.on('console', (message) => {
    if (message.type() === 'error' && !message.text().includes('status of 409 (Conflict)') && !message.text().includes('status of 404 (Not Found)')) {
      errors.push(message.text())
    }
  })
  return errors
}

type SafetyScanPage = {
  path: string
  heading?: string
  readyLabel?: string
  readyText?: string
  scanBody?: boolean
}

async function assertNoForbiddenTradingAffordanceAcrossKeyPages(page: Page, paths: SafetyScanPage[]) {
  const forbidden = /自动下单|一键交易|代下单|券商下单|券商接口|自动规则应用|自动应用规则|自动确认|自动修复|外部推送|短信|邮件|Webhook|第三方推送|收益承诺|完整密钥|API key|sk-|SELECT \* FROM|\/Users\/private|prompt:/
  for (const target of paths) {
    await page.goto(target.path)
    await assertPrimaryNavigationReachable(page)
    if (target.heading) {
      await expect(page.getByRole('heading', { name: target.heading })).toBeVisible()
    }
    if (target.readyLabel) {
      await expect(page.getByLabel(target.readyLabel)).toBeVisible()
    }
    if (target.readyText) {
      await expect(page.getByText(target.readyText, { exact: false }).first()).toBeVisible()
    }
    await page.waitForLoadState('networkidle')
    await expect(page.getByRole('button', { name: forbidden })).toHaveCount(0)
    await expect(page.getByRole('link', { name: forbidden })).toHaveCount(0)
    if (target.scanBody) {
      await expect(page.locator('body')).not.toContainText(forbidden)
    }
  }
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

async function assertMobileNavigationKeyboardPath(page: Page) {
  const toggle = page.getByRole('button', { name: '导航' })
  await expect(toggle).toBeVisible()
  await toggle.focus()
  await page.keyboard.press('Enter')
  await expect(toggle).toHaveAttribute('aria-expanded', 'true')
  const workbenchLink = page.getByRole('link', { name: '决策工作台' })
  await workbenchLink.focus()
  await expect(workbenchLink).toBeFocused()
  await page.keyboard.press('Enter')
  await expect(page.getByRole('heading', { name: '用户决策工作台' })).toBeVisible()
}

async function assertDesktopNavigationKeyboardPath(page: Page) {
  await expect(page.getByRole('navigation', { name: '主导航' })).toBeVisible()
  const workbenchLink = page.getByRole('link', { name: '决策工作台' })
  await workbenchLink.focus()
  await expect(workbenchLink).toBeFocused()
  await workbenchLink.press('Enter')
  await expect(page.getByRole('heading', { name: '用户决策工作台' })).toBeVisible()
}

async function expectPageNotHorizontallyOverflowed(page: Page) {
  await expect.poll(async () => page.evaluate(() => {
    const viewport = window.innerWidth
    return document.body.scrollWidth <= viewport && document.documentElement.scrollWidth <= viewport
  })).toBe(true)
}
