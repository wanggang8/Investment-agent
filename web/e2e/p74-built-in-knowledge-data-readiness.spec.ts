import { expect, test } from '@playwright/test'
import type { APIRequestContext, Page, Response } from '@playwright/test'
import { mkdirSync, writeFileSync } from 'node:fs'
import path from 'node:path'

type P74UIResult = {
  id: string
  route: string
  status: 'passed'
  evidence: string[]
}

test('P74 validates built-in knowledge, data readiness, LLM traceability, and mobile UI', async ({ page, request }) => {
  test.setTimeout(180_000)

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

  const results: P74UIResult[] = []
  results.push(await validateDataQualityReadinessPanel(page))
  results.push(await validateDecisionReadback(page))
  results.push(await validateBlockedReadinessAPI(request))
  results.push(await validateMobileReadiness(page))

  await assertNoForbiddenAffordance(page)
  expect(consoleErrors).toEqual([])
  expect(pageErrors).toEqual([])
  expect(unexpectedFailedApiResponses).toEqual([])
  writeResults({
    generated_at: new Date().toISOString(),
    status: 'passed',
    ui_tasks: results,
    console_errors: consoleErrors,
    page_errors: pageErrors,
    unexpected_failed_api_responses: unexpectedFailedApiResponses,
    accepted_boundaries: [
      'P74 UI acceptance uses a temporary local SQLite database prepared by the runner.',
      'Built-in knowledge is visible only as discipline/rule/LLM context and is not treated as formal market evidence.',
      'The UI does not expose full prompts, secrets, raw provider payloads, broker actions, auto trading, auto confirmation, or auto rule application.',
    ],
  })
})

async function validateDataQualityReadinessPanel(page: Page): Promise<P74UIResult> {
  await page.goto('/data-quality')
  await expect(page.getByRole('heading', { name: '数据质量可观测' })).toBeVisible()
  const panel = page.getByLabel('知识与数据准备度')
  await expect(panel.getByRole('heading', { name: '已准备' })).toBeVisible()
  await expect(panel.getByText('沪深300ETF · ETF · 跟踪 000300')).toBeVisible()
  await expect(panel.getByText('格雷厄姆：安全边际')).toBeVisible()
  await expect(panel.getByText('纪律：不凭单一信源决策')).toBeVisible()
  await expect(panel.getByText(/不作正式证据/).first()).toBeVisible()
  await expect(panel.getByText('LLM 上下文：已附加知识与数据准备度摘要')).toBeVisible()
  await expect(panel.getByText('生效规则 · 已准备')).toBeVisible()
  await expect(panel.getByText('正式证据 · 已准备')).toBeVisible()
  await capture(page, 'data-quality-readiness.png')
  return passed('data_quality_readiness_panel', '/data-quality', ['已准备', '格雷厄姆：安全边际', '生效规则 · 已准备', '正式证据 · 已准备'])
}

async function validateDecisionReadback(page: Page): Promise<P74UIResult> {
  await page.goto('/decisions/decision_smoke_p30')
  await expect(page.getByRole('heading', { name: '决策详情' })).toBeVisible()
  await expect(page.getByText('Agent 分析材料')).toBeVisible()
  await expect(page.getByText('LLM 已参考知识与数据准备度摘要')).toBeVisible()
  await expect(page.getByText(/prompt p74-knowledge-readiness-v1/)).toBeVisible()
  await expect(page.getByText(/仅展示脱敏摘要/)).toBeVisible()
  const body = await page.locator('body').innerText()
  expect(body).not.toMatch(/完整 prompt|principles=|data_readiness=|sk-|BEGIN RSA PRIVATE KEY/)
  await capture(page, 'decision-readiness-readback.png')
  return passed('decision_readiness_readback', '/decisions/decision_smoke_p30', ['LLM 已参考知识与数据准备度摘要', 'prompt p74-knowledge-readiness-v1', '仅展示脱敏摘要'])
}

async function validateBlockedReadinessAPI(request: APIRequestContext): Promise<P74UIResult> {
  const response = await request.get('/api/v1/knowledge-readiness?symbol=999999')
  await expect(response).toBeOK()
  const payload = await response.json()
  expect(payload.data.status).toBe('blocked')
  expect(payload.data.symbol_profile.known).toBe(false)
  expect(JSON.stringify(payload)).toContain('symbol_profile')
  expect(JSON.stringify(payload)).not.toMatch(/sk-|BEGIN RSA PRIVATE KEY|raw HTTP|prompt:/)
  return passed('blocked_symbol_api_readiness', '/api/v1/knowledge-readiness?symbol=999999', ['blocked', 'known=false', 'no sensitive leakage'])
}

async function validateMobileReadiness(page: Page): Promise<P74UIResult> {
  await page.setViewportSize({ width: 390, height: 844 })
  await page.goto('/data-quality')
  await assertPrimaryNavigationReachable(page)
  await expect(page.getByRole('heading', { name: '数据质量可观测' })).toBeVisible()
  await expect(page.getByLabel('知识与数据准备度').getByRole('heading', { name: '已准备' })).toBeVisible()
  await expectPageNotHorizontallyOverflowed(page)
  await capture(page, 'mobile-data-quality-readiness.png')
  await page.setViewportSize({ width: 1280, height: 900 })
  return passed('mobile_readiness_reflow', '/data-quality @390px', ['主导航可达', '知识与数据准备度', '无水平溢出'])
}

function passed(id: string, route: string, evidence: string[]): P74UIResult {
  return { id, route, status: 'passed', evidence }
}

async function capture(page: Page, name: string) {
  if (process.env.P74_CAPTURE_SCREENSHOTS !== '1') return
  const artifactDir = process.env.P74_ARTIFACT_DIR || path.resolve(process.cwd(), '..', 'docs', 'release', 'ui-audit-assets', '2026-06-19-p74')
  mkdirSync(artifactDir, { recursive: true })
  await page.screenshot({ path: path.join(artifactDir, name), fullPage: true })
}

function writeResults(payload: unknown) {
  const artifactDir = process.env.P74_ARTIFACT_DIR || path.resolve(process.cwd(), '..', 'docs', 'release', 'ui-audit-assets', '2026-06-19-p74')
  mkdirSync(artifactDir, { recursive: true })
  writeFileSync(path.join(artifactDir, 'browser-results.json'), `${JSON.stringify(payload, null, 2)}\n`)
}

function classifyFailedApiResponse(response: Response) {
  const method = response.request().method()
  const url = response.url()
  if (method === 'GET' && response.status() === 404 && url.includes('/api/v1/portfolio/current')) return 'expected'
  if (method === 'GET' && response.status() === 409 && url.includes('/api/v1/dashboard/today')) return 'expected'
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
  for (const pathName of ['/data-quality', '/decisions/decision_smoke_p30']) {
    await page.goto(pathName)
    await expect(page.getByRole('button', { name: forbidden })).toHaveCount(0)
    await expect(page.getByRole('link', { name: forbidden })).toHaveCount(0)
  }
}
