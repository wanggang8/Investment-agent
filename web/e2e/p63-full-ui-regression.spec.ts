import { expect, test } from '@playwright/test'
import type { Page, Response } from '@playwright/test'
import { mkdirSync, writeFileSync } from 'node:fs'
import path from 'node:path'

type RouteCheck = {
  path: string
  label: string
  heading?: string
  readyText?: string
  readyLabel?: string
  screenshot?: boolean
}

type RouteResult = {
  path: string
  label: string
  viewport: string
  bodyScrollWidth: number
  documentScrollWidth: number
  viewportWidth: number
  overflow: boolean
}

type FailedApiResponse = {
  url: string
  status: number
  method: string
  classification: 'expected_client_state' | 'unexpected'
}

const viewports = [
  { name: '390', width: 390, height: 844 },
  { name: '768', width: 768, height: 900 },
  { name: '1280', width: 1280, height: 900 },
]

const routeChecks: RouteCheck[] = [
  { path: '/', label: 'dashboard', heading: '今日纪律', readyText: '今日纪律报告', screenshot: true },
  { path: '/workbench', label: 'workbench', heading: '用户决策工作台', readyText: '数据可信度', screenshot: true },
  { path: '/consultation', label: 'consultation', heading: '主动咨询', readyLabel: '咨询问题', screenshot: true },
  { path: '/decisions/decision_smoke_p30', label: 'decision-detail', heading: '决策详情', readyText: '决策故事', screenshot: true },
  { path: '/evidence', label: 'evidence', heading: '情报与证据', readyText: '证据可信度', screenshot: true },
  { path: '/decision-loop', label: 'decision-loop', heading: '决策闭环解释', readyText: '只读决策生命周期', screenshot: true },
  { path: '/positions', label: 'positions', heading: '组合与持仓维护', readyText: '组合维护状态', screenshot: true },
  { path: '/data-quality', label: 'data-quality', heading: '数据质量可观测', readyText: '数据质量总览', screenshot: true },
  { path: '/risk-alerts', label: 'risk-alerts', heading: '风险预警中心', readyText: '风险处置队列', screenshot: true },
  { path: '/risk-alerts/risk_smoke_p39', label: 'risk-alert-detail', heading: '风险预警中心', readyText: 'P39 source health stale 触发数据降级风险', screenshot: true },
  { path: '/rules', label: 'rules', heading: '规则与纪律', readyText: '规则治理状态', screenshot: true },
  { path: '/audit', label: 'audit', heading: '复盘与审计', readyText: '审计检查状态', screenshot: true },
  { path: '/notifications', label: 'notifications', heading: '通知中心', readyText: '本地通知收件箱', screenshot: true },
  { path: '/daily-auto-run', label: 'daily-auto-run', heading: '每日自动运行', readyText: '每日自动运行健康', screenshot: true },
  { path: '/daily-discipline/reports', label: 'daily-reports', heading: '每日纪律报告历史', readyText: '每日纪律复盘状态', screenshot: true },
  { path: '/daily-discipline/reports/daily_report_smoke_p32', label: 'daily-report-detail', heading: '每日纪律报告详情', readyText: 'P32 smoke 今日纪律报告已生成', screenshot: true },
  { path: '/review', label: 'review', heading: '复盘摘要', readyText: '只读追踪', screenshot: true },
  { path: '/local-install', label: 'local-install', heading: '本地安装与诊断', readyText: '本地配置与诊断状态', screenshot: true },
  { path: '/local-knowledge', label: 'local-knowledge', heading: '本地知识导入', readyText: '脱敏预览', screenshot: true },
  { path: '/settings', label: 'settings', heading: '设置', readyText: '本地配置与诊断状态', screenshot: true },
]

test('P63 full UI regression covers all primary routes, reflow, safety, and real consultation classification', async ({ page, request }) => {
  test.setTimeout(240_000)

  const health = await request.get('/api/v1/health')
  await expect(health).toBeOK()

  const consoleErrors: string[] = []
  const pageErrors: string[] = []
  const failedApiResponses: FailedApiResponse[] = []
  const unexpectedFailedApiResponses: FailedApiResponse[] = []
  const failedResourceResponses: Array<{ url: string; status: number; method: string }> = []
  page.on('pageerror', (error) => pageErrors.push(error.message))
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
    const isApiResponse = response.url().includes('/api/v1/')
    if (response.status() >= 400 && isApiResponse) {
      const failedResponse = {
        url: redactResponseUrl(response.url()),
        status: response.status(),
        method: response.request().method(),
        classification: classifyFailedApiResponse(response),
      }
      failedApiResponses.push(failedResponse)
      if (failedResponse.classification === 'unexpected') {
        unexpectedFailedApiResponses.push(failedResponse)
      }
    } else if (response.status() >= 500 && !isApiResponse) {
      failedResourceResponses.push({
        url: redactResponseUrl(response.url()),
        status: response.status(),
        method: response.request().method(),
      })
    }
  })

  const routeResults: RouteResult[] = []
  for (const viewport of viewports) {
    await page.setViewportSize({ width: viewport.width, height: viewport.height })
    for (const target of routeChecks) {
      await page.goto(target.path)
      await assertRouteReady(page, target)
      const overflow = await page.evaluate(() => ({
        bodyScrollWidth: document.body.scrollWidth,
        documentScrollWidth: document.documentElement.scrollWidth,
        viewportWidth: window.innerWidth,
      }))
      const hasOverflow = overflow.bodyScrollWidth > overflow.viewportWidth || overflow.documentScrollWidth > overflow.viewportWidth
      routeResults.push({ path: target.path, label: target.label, viewport: viewport.name, ...overflow, overflow: hasOverflow })
      expect(hasOverflow, `${target.path} overflowed at ${viewport.name}px`).toBe(false)
      await assertNoForbiddenVisibleAffordance(page)
      if (shouldCaptureScreenshots() && target.screenshot) {
        await page.screenshot({ path: path.join(artifactDir(), `${target.label}-${viewport.name}.png`), fullPage: true })
      }
    }
  }

  const consultationResult = await runConsultationJourney(page)

  expect(consoleErrors).toEqual([])
  expect(pageErrors).toEqual([])
  expect(unexpectedFailedApiResponses).toEqual([])
  expect(failedResourceResponses).toEqual([])

  writeResults({
    generated_at: new Date().toISOString(),
    route_results: routeResults,
    failed_overflow_count: routeResults.filter((item) => item.overflow).length,
    consultation: consultationResult,
    console_errors: consoleErrors,
    page_errors: pageErrors,
    failed_api_responses: failedApiResponses,
    unexpected_failed_api_responses: unexpectedFailedApiResponses,
    failed_responses: unexpectedFailedApiResponses,
    failed_resource_responses: failedResourceResponses,
  })
})

async function assertRouteReady(page: Page, target: RouteCheck) {
  if (target.heading) {
    await expect(page.getByRole('heading', { name: target.heading, exact: true })).toBeVisible()
  }
  if (target.readyLabel) {
    await expect(page.getByLabel(target.readyLabel)).toBeVisible()
  }
  if (target.readyText) {
    await expect(page.getByText(target.readyText, { exact: false }).first()).toBeVisible()
  }
  await assertPrimaryNavigationReachable(page)
}

async function runConsultationJourney(page: Page) {
  await page.setViewportSize({ width: 1280, height: 900 })
  await ensurePortfolioCalibration(page)
  await page.goto('/consultation')
  await page.getByLabel('咨询问题').fill('P63 真实 UI 回归：510300 当前是否继续持有？')
  await page.getByLabel('标的代码').fill('510300')

  let consultResponse: Response | undefined
  const responsePromise = page.waitForResponse((response) => response.url().includes('/api/v1/decisions/consult') && response.request().method() === 'POST', { timeout: 150_000 })
  await page.getByRole('button', { name: '提交咨询' }).click()
  try {
    consultResponse = await responsePromise
  } catch {
    consultResponse = undefined
  }

  await page.getByText('决策故事', { exact: true }).waitFor({ timeout: 180_000 })
  await expect(page.getByRole('link', { name: '打开生成的决策详情' })).toBeVisible()
  const detailHref = await page.getByRole('link', { name: '打开生成的决策详情' }).getAttribute('href')
  await page.getByRole('link', { name: '打开生成的决策详情' }).click()
  await expect(page.getByRole('heading', { name: '决策详情' })).toBeVisible()
  await expect(page.getByText('以下内容仅作为分析材料，最终裁决仍以规则链为准。')).toBeVisible()
  await assertNoForbiddenVisibleAffordance(page)

  let body: any = null
  if (consultResponse) {
    try {
      body = await consultResponse.json()
    } catch {
      body = null
    }
  }
  const decision = body?.data ?? null
  const analystReports = Array.isArray(decision?.analyst_reports) ? decision.analyst_reports : []
  return {
    response_status: consultResponse?.status() ?? null,
    decision_id: decision?.decision_id ?? detailHref?.split('/').pop() ?? null,
    detail_href: detailHref,
    workflow_status: decision?.workflow_status ?? null,
    analyst_report_count: analystReports.length,
    parse_statuses: analystReports.map((item: any) => item?.parse_status ?? ''),
    quality_statuses: analystReports.map((item: any) => item?.quality_status ?? ''),
    llm_displayed: analystReports.length > 0,
  }
}

async function ensurePortfolioCalibration(page: Page) {
  await page.goto('/positions')
  await expect(page.getByRole('heading', { name: '组合与持仓维护', exact: true })).toBeVisible()
  await page.getByLabel('现金').fill('70')
  await page.getByLabel('总资产').fill('100')
  await page.getByLabel('标的代码').fill('510300')
  await page.getByLabel('标的名称').fill('沪深300ETF')
  await page.getByLabel('数量').fill('10')
  await page.getByLabel('成本价').fill('2')
  await page.getByLabel('现价').fill('3')
  await page.getByLabel('买入理由').fill('P63 真实 UI 回归前置校准')
  await page.getByLabel('资产标签').fill('长期配置')
  await page.getByRole('button', { name: '保存本地校准' }).click()
  await expect(page.getByText('账户校准已保存为本地事实；不会连接交易接口。')).toBeVisible()
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

async function assertNoForbiddenVisibleAffordance(page: Page) {
  const forbidden = /自动下单|一键交易|代下单|券商下单|券商接口|自动规则应用|自动应用规则|自动确认|自动修复|外部推送|短信|邮件|Webhook|第三方推送|收益承诺|完整密钥|API key|sk-|SELECT \* FROM|\/Users\/private|prompt:/
  await expect(page.getByRole('button', { name: forbidden })).toHaveCount(0)
  await expect(page.getByRole('link', { name: forbidden })).toHaveCount(0)
}

function shouldCaptureScreenshots() {
  return process.env.P63_CAPTURE_SCREENSHOTS === '1'
}

function artifactDir() {
  const output = process.env.P63_ARTIFACT_DIR ?? path.resolve(process.cwd(), '../tmp/p63-ui-regression')
  mkdirSync(output, { recursive: true })
  return output
}

function writeResults(payload: unknown) {
  const output = artifactDir()
  writeFileSync(path.join(output, 'browser-results.json'), `${JSON.stringify(payload, null, 2)}\n`)
}

function classifyFailedApiResponse(response: Response): FailedApiResponse['classification'] {
  if (response.status() === 404 || response.status() === 409) {
    return 'expected_client_state'
  }
  return 'unexpected'
}

function redactResponseUrl(url: string) {
  try {
    const parsed = new URL(url)
    return `${parsed.pathname}${parsed.search ? '?<redacted_query>' : ''}`
  } catch {
    return '<unparseable_url>'
  }
}
