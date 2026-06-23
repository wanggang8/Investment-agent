import { expect, test } from '@playwright/test'
import type { APIRequestContext, Page } from '@playwright/test'
import { mkdirSync, writeFileSync } from 'node:fs'
import path from 'node:path'

test('P83 governance traceability is verified through real UI, API, and readback', async ({ page, request }) => {
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

  const monthly = await expectReviewSummary(request, 'monthly')
  const quarterly = await expectReviewSummary(request, 'quarterly')

  await runReviewReadback(page)
  await runRulesReadback(page)
  await runAuditAndNotificationReadback(page)
  await runOpsReadback(page)
  await runMobileGovernanceReadback(page)

  expect(consoleErrors).toEqual([])
  expect(pageErrors).toEqual([])
  writeResults({
    generated_at: new Date().toISOString(),
    status: 'passed',
    review_api: {
      monthly: summarizeReview(monthly),
      quarterly: summarizeReview(quarterly),
    },
    ui_routes: ['/review', '/rules', '/audit', '/notifications', '/daily-discipline/reports', '/local-install'],
    safety_boundaries: ['no broker UI', 'no external push action', 'no automatic rule application'],
  })
})

async function expectReviewSummary(request: APIRequestContext, period: 'monthly' | 'quarterly') {
  const res = await request.get(`/api/v1/review/summary?period=${period}`)
  await expect(res).toBeOK()
  const body = await res.json() as { data?: Record<string, any> }
  const data = body.data ?? {}
  expect(data.period).toBe(period)
  expect(data.decision_count).toBeGreaterThan(0)
  expect(data.confirmation_count).toBeGreaterThan(0)
  expect(data.error_case_count).toBeGreaterThan(0)
  expect(data.rule_proposal_count).toBeGreaterThan(0)
  expect(data.audit_event_count).toBeGreaterThan(0)
  expect(data.rule_hit_count).toBeGreaterThan(0)
  expect(data.rule_suggestions?.some((item: any) => item.can_auto_apply === false)).toBeTruthy()
  expect(data.attribution_summaries?.length).toBeGreaterThan(0)
  expect(data.recurring_error_tags?.length).toBeGreaterThan(0)
  expect(data.rule_proposal_outcomes?.length).toBeGreaterThan(0)
  expect(data.rule_effect_tracking?.length).toBeGreaterThan(0)
  expect(data.tracking_links?.length).toBeGreaterThan(0)
  return data
}

async function runReviewReadback(page: Page) {
  await page.goto('/review')
  await expect(page.getByRole('heading', { name: '复盘摘要' })).toBeVisible()
  await expect(page.getByText('归因摘要')).toBeVisible()
  await expect(page.getByText('高频错误标签')).toBeVisible()
  await expect(page.getByText('规则提案结果')).toBeVisible()
  await expect(page.getByText('规则应用后效果追踪')).toBeVisible()
  await expect(page.getByText('规则变更仍需守门人审计和用户最终确认，不会自动应用。')).toBeVisible()
  await expect(page.getByText('rule_threshold_issue').first()).toBeVisible()
  await expect(page.getByText(/P83|P39|阈值复盘/).first()).toBeVisible()
  await capture(page, 'p83-review-readback.png')
}

async function runRulesReadback(page: Page) {
  await page.goto('/rules')
  await expect(page.getByRole('heading', { name: '规则与纪律' })).toBeVisible()
  await expect(page.getByText('规则治理状态')).toBeVisible()
  const p39Proposal = page.locator('section.proposal-item').filter({ hasText: 'P39 E2E 规则提案' })
  await expect(p39Proposal).toBeVisible()
  await expect(p39Proposal).toContainText('守门人结果：审计通过')
  await expect(p39Proposal).toContainText('验证状态：已通过')
  const masterProposal = page.locator('section.proposal-item').filter({ hasText: '大师权重调整提案' })
  await expect(masterProposal).toBeVisible()
  await expect(masterProposal).toContainText('master.graham.margin_of_safety')
  await expect(masterProposal).toContainText('守门人通过后仍需用户最终确认，正式规则不会自动生效。')
  await capture(page, 'p83-rules-readback.png')
}

async function runAuditAndNotificationReadback(page: Page) {
  await page.goto('/audit')
  await expect(page.getByRole('heading', { name: '复盘与审计' })).toBeVisible()
  await expect(page.getByText('P83ReviewSeed').first()).toBeVisible()
  await expect(page.getByText('prop_p83_quarterly').first()).toBeVisible()
  await expect(page.getByText('review_summary monthly:2026-06').first()).toBeVisible()
  await expect(page.getByText('P83MasterSeed').first()).toBeVisible()
  await expect(page.getByText('prop_p83_master_weight').first()).toBeVisible()
  await capture(page, 'p83-audit-readback.png')

  await page.goto('/notifications')
  await expect(page.getByRole('heading', { name: '通知中心' })).toBeVisible()
  const notification = page.locator('article.panel-card').filter({ hasText: '复盘存在降级或缺失证据' }).first()
  await expect(notification).toBeVisible()
  await notification.getByRole('button', { name: '标记已读' }).click()
  await expect(page.locator('article.panel-card').filter({ hasText: '复盘存在降级或缺失证据' }).first()).toContainText('已读')
  await capture(page, 'p83-notification-readback.png')
}

async function runOpsReadback(page: Page) {
  await page.goto('/daily-discipline/reports')
  await expect(page.getByRole('heading', { name: '每日纪律报告历史' })).toBeVisible()
  await expect(page.getByText('P32 smoke 今日纪律报告已生成')).toBeVisible()
  await page.getByRole('link', { name: '查看报告', exact: true }).click()
  await expect(page.getByRole('heading', { name: '每日纪律报告详情' })).toBeVisible()
  await expect(page.getByText('不会自动执行交易')).toBeVisible()

  await page.goto('/local-install')
  await expect(page.getByRole('heading', { name: '本地安装与诊断' })).toBeVisible()
  await expect(page.getByText('本地配置与诊断状态')).toBeVisible()
  await expect(page.getByText('关键命令')).toBeVisible()
  await expect(page.getByText('本页仅展示本地诊断产物，不读取数据库路径、完整 key、SQL 或原始 HTTP 响应。')).toBeVisible()
  await page.getByLabel('sqlite 路径').fill('/private/tmp/p83-secret.db')
  await expect(page.getByLabel('启动配置草稿')).toContainText('<local-sqlite-path>')
  await capture(page, 'p83-local-install-readback.png')
}

async function runMobileGovernanceReadback(page: Page) {
  await page.setViewportSize({ width: 390, height: 844 })
  for (const target of [
    { route: '/review', heading: '复盘摘要', ready: '规则变更仍需守门人审计' },
    { route: '/rules', heading: '规则与纪律', ready: '大师权重调整提案' },
    { route: '/local-install', heading: '本地安装与诊断', ready: '本地配置与诊断状态' },
  ]) {
    await page.goto(target.route)
    await expect(page.getByRole('heading', { name: target.heading })).toBeVisible()
    await expect(page.getByText(target.ready, { exact: false }).first()).toBeVisible()
    await assertNoHorizontalOverflow(page)
    await capture(page, `p83-mobile-${target.route.replaceAll('/', '_') || 'home'}.png`)
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

function summarizeReview(data: Record<string, any>) {
  return {
    period: data.period,
    decision_count: data.decision_count,
    confirmation_count: data.confirmation_count,
    error_case_count: data.error_case_count,
    rule_proposal_count: data.rule_proposal_count,
    audit_event_count: data.audit_event_count,
    rule_hit_count: data.rule_hit_count,
    rule_effect_tracking_count: data.rule_effect_tracking?.length ?? 0,
    tracking_link_count: data.tracking_links?.length ?? 0,
  }
}

async function capture(page: Page, fileName: string) {
  if (process.env.P83_CAPTURE_SCREENSHOTS !== '1') return
  mkdirSync(artifactDir(), { recursive: true })
  await page.screenshot({ path: path.join(artifactDir(), fileName), fullPage: true })
}

function artifactDir() {
  return process.env.P83_ARTIFACT_DIR ?? path.resolve(process.cwd(), '../tmp/p83-governance-traceability')
}

function writeResults(payload: unknown) {
  mkdirSync(artifactDir(), { recursive: true })
  writeFileSync(path.join(artifactDir(), 'browser-results.json'), `${JSON.stringify(payload, null, 2)}\n`)
}
