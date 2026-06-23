import { expect, test } from '@playwright/test'
import type { Page } from '@playwright/test'
import { mkdirSync, writeFileSync } from 'node:fs'
import path from 'node:path'

test('P75 continuous non-510300 browser journey binds UI, readiness, LLM, and data readback', async ({ page, request }) => {
  test.setTimeout(300_000)

  const health = await request.get('/api/v1/health')
  await expect(health).toBeOK()

  const capability = await request.put('/api/v1/settings/capability', {
    data: {
      asset_types: ['ETF'],
      symbols: ['159915', '510300'],
      excluded_symbols: [],
      strategy_scope: ['long_term_etf'],
      notes: 'P75 non-510300 accepted-local UI journey capability config',
    },
  })
  await expect(capability).toBeOK()

  const consoleErrors: string[] = []
  const pageErrors: string[] = []
  page.on('pageerror', (error) => pageErrors.push(error.stack || error.message))
  page.on('console', (message) => {
    if (message.type() === 'error' && !message.text().includes('Failed to load resource: the server responded with a status of')) {
      consoleErrors.push(message.text())
    }
  })

  await runPortfolioAddJourney(page, request)
  const readiness = await runDataQualityReadinessJourney(page, request)
  const decision = await runConsultationJourney(page)
  await runDerivedReadbackJourney(page, decision.decision_id)

  expect(consoleErrors).toEqual([])
  expect(pageErrors).toEqual([])

  writeResults({
    generated_at: new Date().toISOString(),
    status: 'passed',
    symbol: '159915',
    tracked_index_symbol: '399006',
    readiness,
    decision: sanitizeDecisionForArtifact(decision),
  })
})

async function runPortfolioAddJourney(page: Page, request: any) {
  await page.goto('/positions')
  await expect(page.getByRole('heading', { name: '组合与持仓维护', exact: true })).toBeVisible()
  await fillPortfolioForm(page, {
    cash: '99000',
    totalAssets: '99723.90',
    symbol: '159915',
    name: '创业板ETF',
    quantity: '300',
    costPrice: '2.20',
    currentPrice: '2.413',
    buyReason: 'P75 非510300连续真实UI场景：创业板ETF作为卫星仓位人工录入',
    assetTag: 'satellite',
    riskPreference: 'steady',
  })
  const adjustmentResponsePromise = page.waitForResponse((response) => response.url().includes('/api/v1/portfolio/adjustments') && response.request().method() === 'POST')
  await page.getByRole('button', { name: '保存本地校准' }).click()
  const adjustmentResponse = await adjustmentResponsePromise
  expect(adjustmentResponse.status()).toBe(200)
  await expect(page.getByText('账户校准已保存为本地事实；不会连接交易接口。')).toBeVisible()
  await expect(page.getByText('159915').first()).toBeVisible()
  await capture(page, 'p75-159915-portfolio.png')

  const portfolio = await request.get('/api/v1/portfolio/current')
  await expect(portfolio).toBeOK()
  const body = await portfolio.json()
  expect(body.data?.positions?.some((item: any) => item.symbol === '159915')).toBe(true)
}

async function runDataQualityReadinessJourney(page: Page, request: any) {
  const readinessAPI = await request.get('/api/v1/knowledge-readiness?symbol=159915')
  await expect(readinessAPI).toBeOK()
  const readinessBody = await readinessAPI.json()
  expect(readinessBody.data?.symbol_profile?.tracked_index_symbol).toBe('399006')
  expect(readinessBody.data?.data_dependencies?.some((item: any) => item.category === 'valuation_percentiles' && item.status === 'ready' && item.affected_symbols?.includes('399006'))).toBe(true)

  await page.goto('/data-quality?symbol=159915')
  await expect(page.getByRole('heading', { name: '数据质量可观测' })).toBeVisible()
  await expect(page.getByText('当前查看：159915')).toBeVisible()
  const panel = page.getByLabel('知识与数据准备度')
  await expect(panel.getByRole('heading', { name: '已准备' })).toBeVisible()
  await expect(panel.getByText('创业板ETF · ETF · 跟踪 399006')).toBeVisible()
  await expect(panel.getByText('跟踪指数 · 已准备')).toBeVisible()
  await expect(panel.getByText('估值分位 · 已准备')).toBeVisible()
  await expect(panel.getByText(/request：req_/).first()).toBeVisible()
  await expect(panel.getByText(/标的：399006/).first()).toBeVisible()
  await capture(page, 'p75-159915-data-quality.png')
  return readinessBody.data
}

async function runConsultationJourney(page: Page) {
  await page.goto('/consultation')
  await expect(page.getByRole('heading', { name: '主动咨询' })).toBeVisible()
  await page.getByLabel('标的代码').fill('159915')
  await page.getByLabel('咨询问题').fill('P75 非510300连续真实UI场景：159915 创业板ETF 当前数据准备度通过后，应如何人工复核持有风险？')
  const responsePromise = page.waitForResponse((response) => response.url().includes('/api/v1/decisions/consult') && response.request().method() === 'POST', { timeout: 180_000 })
  await page.getByRole('button', { name: '提交咨询' }).click()
  const response = await responsePromise
  expect(response.status()).toBe(200)
  const body = await response.json()
  const decision = body.data
  expect(decision.symbol).toBe('159915')
  expect(decision.workflow_status).toBe('completed')
  expect(decision.analyst_reports?.length ?? 0).toBeGreaterThanOrEqual(3)
  await expect(page.getByRole('link', { name: '打开生成的决策详情' })).toBeVisible()
  await capture(page, 'p75-159915-consultation.png')
  return decision
}

async function runDerivedReadbackJourney(page: Page, decisionID: string) {
  await page.goto(`/decisions/${decisionID}`)
  await expect(page.getByRole('heading', { name: '决策详情' })).toBeVisible()
  await expect(page.getByText('标的：159915').first()).toBeVisible()
  await expect(page.getByText('工作流状态：已完成')).toBeVisible()
  await expect(page.getByText('Agent 分析材料')).toBeVisible()
  await capture(page, 'p75-159915-decision-detail.png')

  await page.goto('/decision-loop')
  await expect(page.getByRole('heading', { name: '决策闭环解释' })).toBeVisible()
  await expect(page.getByText('159915').first()).toBeVisible()
  await capture(page, 'p75-159915-decision-loop.png')
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

async function capture(page: Page, fileName: string) {
  if (process.env.P75_CAPTURE_SCREENSHOTS !== '1') return
  mkdirSync(artifactDir(), { recursive: true })
  await page.screenshot({ path: path.join(artifactDir(), fileName), fullPage: true })
}

function artifactDir() {
  return process.env.P75_ARTIFACT_DIR ?? path.resolve(process.cwd(), '../tmp/p75-non-510300-real-ui')
}

function sanitizeDecisionForArtifact(decision: any) {
  const reports = Array.isArray(decision?.analyst_reports) ? decision.analyst_reports : []
  return {
    ...decision,
    analyst_reports: reports.map((report: any) => ({
      agent_name: report?.agent_name,
      model: report?.model,
      prompt_version: report?.prompt_version,
      parse_status: report?.parse_status,
      quality_status: report?.quality_status,
      confidence: report?.confidence,
      evidence_ids: report?.evidence_ids,
      input_summary: report?.input_summary,
      output_summary: report?.output_summary,
      conclusion_preview: preview(report?.conclusion),
    })),
  }
}

function preview(value: unknown) {
  const text = String(value ?? '').replace(/\s+/g, ' ').trim()
  return text.length > 180 ? `${text.slice(0, 180)}...` : text
}

function writeResults(payload: unknown) {
  mkdirSync(artifactDir(), { recursive: true })
  writeFileSync(path.join(artifactDir(), 'browser-results.json'), `${JSON.stringify(payload, null, 2)}\n`)
}
