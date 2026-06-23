import { expect, test } from '@playwright/test'
import type { Page, Response } from '@playwright/test'
import { execFile } from 'node:child_process'
import { mkdirSync, writeFileSync } from 'node:fs'
import path from 'node:path'
import { promisify } from 'node:util'

const execFileAsync = promisify(execFile)

type FailedApiResponse = {
  url: string
  status: number
  method: string
  classification: 'expected_client_state' | 'unexpected'
}

test('P72 real user fund scenario validates UI operations, readbacks, LLM, RAG, daily discipline, and data impact', async ({ page, request }) => {
  test.setTimeout(600_000)

  const health = await request.get('/api/v1/health')
  await expect(health).toBeOK()

  const currentRegression = await request.get('/api/v1/data-source-quality/regression?mode=current&symbol=000300')
  await expect(currentRegression).toBeOK()
  const currentRegressionBody = await currentRegression.json()
  expect(currentRegressionBody.data?.policy?.verdict).toBe('passed')
  expect(currentRegressionBody.data?.policy?.release_gate).toBe('pass')
  expect(currentRegressionBody.data?.status).toBe('passed')

  const consoleErrors: string[] = []
  const pageErrors: string[] = []
  const failedApiResponses: FailedApiResponse[] = []
  const unexpectedFailedApiResponses: FailedApiResponse[] = []
  const failedResourceResponses: Array<{ url: string; status: number; method: string }> = []

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

  const invalidState = await runInvalidInputJourney(page)
  const portfolioResult = await runPortfolioJourney(page, request)
  const localKnowledgeResult = await runLocalKnowledgeJourney(page, request)
  const marketResult = await runMarketAndDataQualityJourney(page, request)
  const consultationResult = await runConsultationAndConfirmationJourney(page)
  const dailyResult = await runDailyDisciplineAndReadbackJourney(page)
  const reviewResult = await runGovernanceReadbackJourney(page)

  expect(consoleErrors).toEqual([])
  expect(pageErrors).toEqual([])
  expect(unexpectedFailedApiResponses).toEqual([])
  expect(failedResourceResponses).toEqual([])

  await assertNoForbiddenVisibleAffordance(page)

  writeResults({
    generated_at: new Date().toISOString(),
    current_data: currentRegressionBody.data,
    invalid_state: invalidState,
    portfolio: portfolioResult,
    local_knowledge: localKnowledgeResult,
    market: marketResult,
    consultation: consultationResult,
    daily: dailyResult,
    review: reviewResult,
    console_errors: consoleErrors,
    page_errors: pageErrors,
    failed_api_responses: failedApiResponses,
    unexpected_failed_api_responses: unexpectedFailedApiResponses,
    failed_resource_responses: failedResourceResponses,
  })
})

async function runInvalidInputJourney(page: Page) {
  await page.goto('/consultation')
  await expect(page.getByRole('heading', { name: '主动咨询' })).toBeVisible()
  await page.getByRole('button', { name: '提交咨询' }).click()
  await expect(page.getByText('请填写咨询问题和标的代码。')).toBeVisible()
  await assertNoForbiddenVisibleAffordance(page)
  return { route: '/consultation', blocked_message: '请填写咨询问题和标的代码。' }
}

async function runPortfolioJourney(page: Page, request: any) {
  await page.goto('/positions')
  await expect(page.getByRole('heading', { name: '组合与持仓维护', exact: true })).toBeVisible()
  await fillPortfolioForm(page, {
    cash: '96000',
    totalAssets: '100000',
    symbol: '510300',
    name: '沪深300ETF',
    quantity: '1000',
    costPrice: '3',
    currentPrice: '4',
    buyReason: 'P72 真实用户场景：长期核心仓位初始化',
    assetTag: 'core',
    riskPreference: 'steady',
  })
  await page.getByRole('button', { name: '保存本地校准' }).click()
  await expect(page.getByText('账户校准已保存为本地事实；不会连接交易接口。')).toBeVisible()
  await expect(page.getByText('510300').first()).toBeVisible()
  await capture(page, 'portfolio-calibration.png')

  await fillPortfolioForm(page, {
    quantity: '1200',
    costPrice: '3.2',
    currentPrice: '4.1',
    buyReason: 'P72 真实用户场景：人工复核后更新成本与现价',
    assetTag: 'core',
  })
  await page.getByRole('button', { name: '保存持仓编辑' }).click()
  await expect(page.getByText('持仓编辑已保存为本地事实。')).toBeVisible()

  await fillPortfolioForm(page, {
    quantity: '100',
    costPrice: '3.9',
    currentPrice: '4.1',
    buyReason: 'P72 真实用户场景：批量导入追加线下确认持仓',
    assetTag: 'satellite',
  })
  await page.getByRole('button', { name: '校验批量导入' }).click()
  await expect(page.getByText('导入校验完成：有效 1 行，无效 0 行。')).toBeVisible()
  await page.getByRole('button', { name: '确认批量导入' }).click()
  await expect(page.getByText('批量导入已确认并保存。')).toBeVisible()

  await page.getByRole('button', { name: '记录修正审计' }).click()
  await expect(page.getByText('错误修正已保存为本地事实。')).toBeVisible()

  await page.getByLabel('线下交易类型').selectOption('buy')
  await page.getByRole('button', { name: '记录线下交易' }).click()
  await expect(page.getByText('线下交易已记录为本地事实。')).toBeVisible()

  await page.reload()
  await expect(page.getByRole('heading', { name: '组合与持仓维护', exact: true })).toBeVisible()
  await expect(page.getByText('510300').first()).toBeVisible()
  await expect(page.getByText('总资产：¥101,330.00')).toBeVisible()
  await capture(page, 'portfolio-after-local-facts.png')

  const portfolio = await request.get('/api/v1/portfolio/current')
  await expect(portfolio).toBeOK()
  const body = await portfolio.json()
  expect(body.data?.snapshot?.total_assets).toBeCloseTo(101330, 2)
  expect(body.data?.snapshot?.cash).toBeCloseTo(95590, 2)
  expect(body.data?.snapshot?.position_count).toBe(2)
  expect(body.data?.positions?.length).toBe(2)

  await assertNoForbiddenVisibleAffordance(page)
  return {
    snapshot: body.data?.snapshot,
    positions: body.data?.positions?.map((item: any) => ({
      symbol: item.symbol,
      quantity: item.quantity,
      cost_price: item.cost_price,
      current_price: item.current_price,
      market_value: item.market_value,
      asset_tag: item.asset_tag,
    })),
  }
}

async function runLocalKnowledgeJourney(page: Page, request: any) {
  await page.goto('/local-knowledge')
  await expect(page.getByRole('heading', { name: '本地知识导入' })).toBeVisible()
  await page.getByLabel('来源标签').fill('p72_real_user_scenario_local_notes')
  await page.getByLabel('默认标的').fill('510300')
  await page.getByLabel('记录 JSON').fill(JSON.stringify([
    {
      title: 'P72 510300 真实用户场景背景',
      text: 'P72 真实用户验收导入的沪深300ETF持仓背景：用户将其作为核心仓位，关注数据质量、仓位纪律、估值风险和人工确认闭环。',
      symbol: '510300',
      tags: ['P72', '真实用户场景', '沪深300ETF'],
    },
  ], null, 2))
  await page.getByRole('button', { name: '校验预览' }).click()
  await expect(page.getByRole('cell', { name: '可写入' })).toBeVisible()
  await expect(page.getByText('预计片段：1')).toBeVisible()
  await page.getByLabel('确认理由').fill('P72 真实用户场景人工确认写入本地背景材料')
  await page.getByRole('button', { name: '写入本地事实' }).click()
  await expect(page.getByText('情报：1 条；摘要：1 条；片段：1 条。')).toBeVisible()

  const rebuild = await request.post('/api/v1/evidence/rebuild-index')
  await expect(rebuild).toBeOK()
  const rebuildBody = await rebuild.json()
  expect(rebuildBody.data?.indexed_count).toBeGreaterThan(0)
  expect(rebuildBody.data?.index_health?.status).toBe('healthy')
  expect(rebuildBody.data?.index_health?.chunk_count).toBeGreaterThan(0)
  await capture(page, 'local-knowledge-confirmed.png')
  await assertNoForbiddenVisibleAffordance(page)

  return {
    source_label: 'p72_real_user_scenario_local_notes',
    index_rebuild: rebuildBody.data,
  }
}

async function runMarketAndDataQualityJourney(page: Page, request: any) {
  const refresh = await request.post('/api/v1/market/refresh', {
    data: { symbols: ['510300'] },
  })
  const refreshBody = await refresh.json()
  expect([200, 503]).toContain(refresh.status())
  if (refresh.status() === 200) {
    expect(refreshBody.data?.failed_symbols ?? []).toEqual([])
    expect(refreshBody.data?.refreshed_count ?? 0).toBeGreaterThan(0)
    expect(refreshBody.data?.latest_snapshot_ids?.length ?? 0).toBeGreaterThan(0)
  } else {
    expect(refreshBody.error?.code).toBe('DATA_SOURCE_UNAVAILABLE')
    expect(refreshBody.error?.message).toContain('市场数据源不可用')
  }

  const sourceHealth = await request.get('/api/v1/market/source-health?symbol=000300')
  await expect(sourceHealth).toBeOK()
  const sourceHealthBody = await sourceHealth.json()

  await page.goto('/data-quality')
  await expect(page.getByRole('heading', { name: '数据质量可观测' })).toBeVisible()
  await expect(page.getByText('数据质量总览')).toBeVisible()
  await capture(page, 'data-quality-current-pass.png')
  await assertNoForbiddenVisibleAffordance(page)

  return {
    market_refresh_status: refresh.status(),
    market_refresh: refreshBody.data,
    market_refresh_error: refreshBody.error,
    source_health: sourceHealthBody.data,
  }
}

async function runConsultationAndConfirmationJourney(page: Page) {
  await page.goto('/consultation')
  await expect(page.getByRole('heading', { name: '主动咨询' })).toBeVisible()
  await page.getByLabel('咨询问题').fill('P72 真实用户场景：510300 作为核心仓位，在当前数据质量通过后是否继续持有并如何人工复核？')
  await page.getByLabel('标的代码').fill('510300')

  const responsePromise = page.waitForResponse((response) => response.url().includes('/api/v1/decisions/consult') && response.request().method() === 'POST', { timeout: 180_000 })
  await page.getByRole('button', { name: '提交咨询' }).click()
  const consultResponse = await responsePromise
  expect(consultResponse.status()).toBe(200)
  const consultBody = await consultResponse.json()
  const decision = consultBody.data
  expect(decision?.workflow_status).toBe('completed')
  expect(decision?.analyst_reports?.length).toBeGreaterThan(0)
  expect(decision?.analyst_reports.every((item: any) => item?.parse_status === 'parsed')).toBe(true)
  expect(decision?.analyst_reports.every((item: any) => item?.quality_status === 'passed')).toBe(true)
  expect(decision?.retrieval_quality?.status).toBe('hit')
  expect(decision?.retrieval_quality?.fallback_source).toBe('veclite')
  expect(decision?.retrieval_quality?.index_health).toBe('healthy')
  expect(decision?.retrieval_quality?.degraded_reason ?? '').toBe('')

  await page.getByText('决策故事', { exact: true }).waitFor({ timeout: 180_000 })
  await expect(page.getByRole('link', { name: '打开生成的决策详情' })).toBeVisible()
  const detailHref = await page.getByRole('link', { name: '打开生成的决策详情' }).getAttribute('href')
  await page.getByRole('link', { name: '打开生成的决策详情' }).click()
  await expect(page.getByRole('heading', { name: '决策详情' })).toBeVisible()
  await expect(page.getByText('Fallback 来源：sqlite_summary')).toHaveCount(0)
  await expect(page.getByText('索引健康：缺失')).toHaveCount(0)
  await expect(page.getByText('索引新鲜度：未知')).toHaveCount(0)
  await capture(page, 'decision-detail-before-confirmation.png')

  await page.getByRole('button', { name: '已手动执行' }).click()
  const confirmationForm = page.locator('[aria-label="确认表单"]')
  await confirmationForm.getByLabel('标的代码').fill('510300')
  await confirmationForm.getByLabel('线下动作').selectOption('sell')
  await confirmationForm.getByLabel('数量').fill('10')
  await confirmationForm.getByLabel('价格').fill('4.05')
  await confirmationForm.getByLabel('费用').fill('0')
  await confirmationForm.getByLabel('执行时间').fill(pastDateTimeLocal())
  await confirmationForm.getByLabel('备注').fill('P72 真实用户场景：用户在线下完成小额卖出后回填确认')
  await page.getByRole('button', { name: '提交确认' }).click()
  await expect(page.getByText('确认已记录。')).toBeVisible()
  await expect(page.getByText('当前确认状态：已手动执行')).toBeVisible()
  await capture(page, 'decision-detail-after-confirmation.png')
  await assertNoForbiddenVisibleAffordance(page)

  return {
    decision_id: decision?.decision_id,
    detail_href: detailHref,
    analyst_report_count: decision?.analyst_reports?.length ?? 0,
    retrieval_quality: decision?.retrieval_quality,
    confirmation_status: 'executed_manually',
  }
}

async function runDailyDisciplineAndReadbackJourney(page: Page) {
  const rootDir = path.resolve(process.cwd(), '..')
  const daily = await execFileAsync('go', ['run', './cmd/agent', '--task', 'daily'], {
    cwd: rootDir,
    env: process.env,
    timeout: 240_000,
    maxBuffer: 1024 * 1024,
  })
  expect(daily.stderr).toBe('')

  await page.goto('/')
  await expect(page.getByRole('heading', { name: '今日纪律', exact: true })).toBeVisible()
  await expect(page.getByText('今日纪律报告', { exact: true })).toBeVisible()
  await capture(page, 'dashboard-after-manual-daily.png')

  await page.goto('/daily-discipline/reports')
  await expect(page.getByRole('heading', { name: '每日纪律报告历史' })).toBeVisible()
  await expect(page.getByText('今日纪律报告已生成', { exact: true })).toBeVisible()
  await capture(page, 'daily-reports-after-manual-daily.png')

  await page.goto('/risk-alerts')
  await expect(page.getByRole('heading', { name: '风险预警中心' })).toBeVisible()
  await expect(page.getByText('风险处置队列')).toBeVisible()

  await page.goto('/notifications')
  await expect(page.getByRole('heading', { name: '通知中心' })).toBeVisible()
  await expect(page.getByText('本地通知收件箱')).toBeVisible()
  await assertNoForbiddenVisibleAffordance(page)

  return {
    stdout: redactLog(daily.stdout),
    stderr: redactLog(daily.stderr),
  }
}

async function runGovernanceReadbackJourney(page: Page) {
  const pages = [
    { path: '/decision-loop', heading: '决策闭环解释', ready: '只读决策生命周期', screenshot: 'decision-loop-readback.png' },
    { path: '/audit', heading: '复盘与审计', ready: '审计检查状态', screenshot: 'audit-readback.png' },
    { path: '/review', heading: '复盘摘要', ready: '只读追踪', screenshot: 'review-readback.png' },
    { path: '/rules', heading: '规则与纪律', ready: '规则治理状态', screenshot: 'rules-readback.png' },
    { path: '/workbench', heading: '用户决策工作台', ready: '数据可信度', screenshot: 'workbench-readback.png' },
  ]
  const out = []
  for (const target of pages) {
    await page.goto(target.path)
    await expect(page.getByRole('heading', { name: target.heading })).toBeVisible()
    await expect(page.getByText(target.ready, { exact: false }).first()).toBeVisible()
    await assertNoForbiddenVisibleAffordance(page)
    await capture(page, target.screenshot)
    out.push({ path: target.path, heading: target.heading })
  }
  return out
}

async function fillPortfolioForm(page: Page, values: {
  cash?: string
  totalAssets?: string
  symbol?: string
  name?: string
  quantity?: string
  costPrice?: string
  currentPrice?: string
  buyReason?: string
  assetTag?: string
  riskPreference?: string
}) {
  if (values.cash !== undefined) await page.getByLabel('现金').fill(values.cash)
  if (values.totalAssets !== undefined) await page.getByLabel('总资产').fill(values.totalAssets)
  if (values.symbol !== undefined) await page.getByLabel('标的代码').fill(values.symbol)
  if (values.name !== undefined) await page.getByLabel('标的名称').fill(values.name)
  if (values.quantity !== undefined) await page.getByLabel('数量').fill(values.quantity)
  if (values.costPrice !== undefined) await page.getByLabel('成本价').fill(values.costPrice)
  if (values.currentPrice !== undefined) await page.getByLabel('现价').fill(values.currentPrice)
  if (values.buyReason !== undefined) await page.getByLabel('买入理由').fill(values.buyReason)
  if (values.assetTag !== undefined) await page.getByLabel('资产标签').fill(values.assetTag)
  if (values.riskPreference !== undefined) await page.getByLabel('风险偏好').fill(values.riskPreference)
}

async function assertNoForbiddenVisibleAffordance(page: Page) {
  const forbidden = /自动下单|一键交易|代下单|券商下单|券商接口|自动规则应用|自动应用规则|自动确认|自动修复|外部推送|短信|邮件|Webhook|第三方推送|收益承诺|完整密钥|API key|sk-|SELECT \* FROM|\/Users\/private|prompt:/
  await expect(page.getByRole('button', { name: forbidden })).toHaveCount(0)
  await expect(page.getByRole('link', { name: forbidden })).toHaveCount(0)
}

async function capture(page: Page, fileName: string) {
  if (!shouldCaptureScreenshots()) return
  await page.screenshot({ path: path.join(artifactDir(), fileName), fullPage: true })
}

function shouldCaptureScreenshots() {
  return process.env.P72_CAPTURE_SCREENSHOTS === '1'
}

function artifactDir() {
  const output = process.env.P72_ARTIFACT_DIR ?? path.resolve(process.cwd(), '../tmp/p72-real-user-fund-scenario')
  mkdirSync(output, { recursive: true })
  return output
}

function writeResults(payload: unknown) {
  const output = artifactDir()
  writeFileSync(path.join(output, 'browser-results.json'), `${JSON.stringify(payload, null, 2)}\n`)
}

function pastDateTimeLocal() {
  const value = new Date(Date.now() - 60 * 60 * 1000)
  const year = value.getUTCFullYear()
  const month = `${value.getUTCMonth() + 1}`.padStart(2, '0')
  const day = `${value.getUTCDate()}`.padStart(2, '0')
  const hour = `${value.getUTCHours()}`.padStart(2, '0')
  const minute = `${value.getUTCMinutes()}`.padStart(2, '0')
  return `${year}-${month}-${day}T${hour}:${minute}`
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

function redactLog(value: string) {
  return value
    .replace(/[A-Za-z0-9_-]{24,}/g, '<redacted_token>')
    .replace(/\/Users\/[^\s"]+/g, '<redacted_path>')
}
