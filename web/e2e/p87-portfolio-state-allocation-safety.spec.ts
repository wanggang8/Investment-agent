import { expect, test } from '@playwright/test'
import type { Page } from '@playwright/test'
import { mkdirSync, writeFileSync } from 'node:fs'
import path from 'node:path'

const artifactDir = process.env.P87_ARTIFACT_DIR || path.resolve(process.cwd(), '../docs/release/ui-audit-assets/2026-06-22-p87-portfolio-state-allocation-safety')
const shouldCapture = process.env.P87_CAPTURE_SCREENSHOTS === '1'

test('P87 portfolio state, allocation and safe-degradation paths are operated through real UI', async ({ page, request }) => {
  test.setTimeout(300_000)
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

  const portfolio = await runPortfolioStateAndAllocation(page, request)
  const decisions = await runDecisionSafetyReadback(page, request)

  expect(consoleErrors).toEqual([])
  expect(pageErrors).toEqual([])
  expect(failedApiResponses).toEqual([])
  await assertNoForbiddenActionButtons(page)

  writeResult({
    generated_at: new Date().toISOString(),
    status: 'passed',
    portfolio,
    decisions,
    ui_paths: ['/positions', '/decisions/decision_p87_sell_only', '/decisions/decision_p87_frozen_watch', '/decisions/decision_p87_insufficient'],
    safety_boundaries: ['no broker UI', 'no automatic confirmation', 'no order placement', 'no external push', 'frozen/insufficient decisions have no confirmation actions'],
    console_errors: consoleErrors,
    page_errors: pageErrors,
    failed_api_responses: failedApiResponses,
  })
})

async function runPortfolioStateAndAllocation(page: Page, request: any) {
  await page.goto('/positions')
  await expect(page.getByRole('heading', { name: '组合与持仓维护', exact: true })).toBeVisible()

  await fillPortfolioForm(page, {
    cash: '8000',
    totalAssets: '72000',
    symbol: '510300',
    name: '沪深300ETF',
    quantity: '16000',
    costPrice: '3.2',
    currentPrice: '4',
    buyDate: '2026-01-05',
    positionState: 'normal',
    buyReason: 'P87 核心仓位纪律验收',
    assetTag: 'core',
    riskPreference: 'steady',
  })
  await page.getByRole('button', { name: '保存本地校准' }).click()
  await expect(page.getByText('账户校准已保存为本地事实；不会连接交易接口。')).toBeVisible()
  await capture(page, 'p87-core-calibration.png')

  await importHolding(page, {
    symbol: '159915',
    name: '创业板ETF',
    quantity: '13500',
    costPrice: '1.6',
    currentPrice: '2',
    buyDate: '2026-01-06',
    positionState: 'sell_only',
    buyReason: 'P87 买入逻辑破坏后只卖不买',
    assetTag: 'satellite',
  })

  await importHolding(page, {
    symbol: '511880',
    name: '货币基金',
    quantity: '1000',
    costPrice: '1',
    currentPrice: '1',
    buyDate: '2026-01-07',
    positionState: 'frozen_watch',
    buyReason: 'P87 多源验证不足冻结观察',
    assetTag: 'cash',
  })

  await page.reload()
  await expect(page.getByText('510300').first()).toBeVisible()
  await expect(page.getByText('159915').first()).toBeVisible()
  await expect(page.getByText('511880').first()).toBeVisible()
  await expect(page.getByText('2026-01-05').first()).toBeVisible()
  await expect(page.getByText('2026-01-06').first()).toBeVisible()
  await expect(page.getByRole('row', { name: /510300 .*正常 2026-01-05/ })).toBeVisible()
  await expect(page.getByRole('row', { name: /159915 .*仅卖出 2026-01-06/ })).toBeVisible()
  await expect(page.getByRole('row', { name: /511880 .*冻结观察 2026-01-07/ })).toBeVisible()
  await expect(page.getByText('现金占比：8.00%')).toBeVisible()
  await capture(page, 'p87-portfolio-final-readback.png')

  const response = await request.get('/api/v1/portfolio/current')
  await expect(response).toBeOK()
  const body = await response.json()
  const positions = body.data?.positions ?? []
  expect(body.data?.snapshot?.cash).toBe(8000)
  expect(body.data?.snapshot?.total_assets).toBe(100000)
  expect(body.data?.snapshot?.cash_ratio).toBeCloseTo(0.08, 4)
  expect(positionBySymbol(positions, '510300')).toMatchObject({ asset_tag: 'core', position_state: 'normal', buy_date: '2026-01-05', market_value: 64000 })
  expect(positionBySymbol(positions, '159915')).toMatchObject({ asset_tag: 'satellite', position_state: 'sell_only', buy_date: '2026-01-06', market_value: 27000 })
  expect(positionBySymbol(positions, '511880')).toMatchObject({ asset_tag: 'cash', position_state: 'frozen_watch', buy_date: '2026-01-07', market_value: 1000 })

  const totalAssets = body.data.snapshot.total_assets
  const coreRatio = bucketValue(positions, 'core') / totalAssets
  const satelliteRatio = bucketValue(positions, 'satellite') / totalAssets
  const cashRatio = body.data.snapshot.cash_ratio
  expect(coreRatio).toBeGreaterThanOrEqual(0.60)
  expect(coreRatio).toBeLessThanOrEqual(0.70)
  expect(satelliteRatio).toBeGreaterThanOrEqual(0.20)
  expect(satelliteRatio).toBeLessThanOrEqual(0.30)
  expect(cashRatio).toBeGreaterThanOrEqual(0.05)
  expect(cashRatio).toBeLessThanOrEqual(0.10)

  return {
    snapshot: body.data.snapshot,
    positions: positions.map((position: any) => ({
      symbol: position.symbol,
      market_value: position.market_value,
      asset_tag: position.asset_tag,
      position_state: position.position_state,
      buy_date: position.buy_date,
    })),
    ratios: { core: coreRatio, satellite: satelliteRatio, cash: cashRatio },
  }
}

async function runDecisionSafetyReadback(page: Page, request: any) {
  const sellOnly = await readDecisionPage(page, request, 'decision_p87_sell_only', {
    statusText: '状态：仅卖出',
    confirmationText: '当前确认状态：待确认',
    expectedText: '禁止新增买入和加仓',
    screenshot: 'p87-sell-only-decision.png',
  })
  const frozenWatch = await readDecisionPage(page, request, 'decision_p87_frozen_watch', {
    statusText: '状态：冻结观察',
    confirmationText: '当前确认状态：无需确认',
    expectedText: '等待更多 A/S 级独立信源',
    screenshot: 'p87-frozen-watch-decision.png',
  })
  await expect(page.getByRole('button', { name: '已手动执行' })).toHaveCount(0)

  const insufficient = await readDecisionPage(page, request, 'decision_p87_insufficient', {
    statusText: '状态：数据不足',
    confirmationText: '当前确认状态：无需确认',
    expectedText: '不生成交易类建议',
    screenshot: 'p87-insufficient-decision.png',
  })
  await expect(page.getByRole('button', { name: '已手动执行' })).toHaveCount(0)

  return { sellOnly, frozenWatch, insufficient }
}

async function readDecisionPage(page: Page, request: any, decisionId: string, expected: { statusText: string; confirmationText: string; expectedText: string; screenshot: string }) {
  await page.goto(`/decisions/${decisionId}`)
  await expect(page.getByText(expected.statusText)).toBeVisible()
  await expect(page.getByText(expected.confirmationText)).toBeVisible()
  await expect(page.getByText(expected.expectedText).first()).toBeVisible()
  await capture(page, expected.screenshot)

  const response = await request.get(`/api/v1/decisions/${decisionId}`)
  await expect(response).toBeOK()
  const body = await response.json()
  return {
    decision_id: body.data?.decision_id,
    record_type: body.data?.record_type,
    final_verdict: body.data?.final_verdict?.status,
    confirmation_status: body.data?.user_confirmation?.confirmation_status,
    prohibited_actions: body.data?.final_verdict?.prohibited_actions ?? [],
  }
}

async function importHolding(page: Page, values: Record<string, string>) {
  await fillPortfolioForm(page, values)
  await page.getByRole('button', { name: '校验批量导入' }).click()
  await expect(page.getByText('导入校验完成：有效 1 行，无效 0 行。')).toBeVisible()
  await page.getByRole('button', { name: '确认批量导入' }).click()
  await expect(page.getByText('批量导入已确认并保存。')).toBeVisible()
}

async function fillPortfolioForm(page: Page, values: Record<string, string>) {
  const labels: Record<string, string> = {
    cash: '现金',
    totalAssets: '总资产',
    symbol: '标的代码',
    name: '标的名称',
    quantity: '数量',
    costPrice: '成本价',
    currentPrice: '现价',
    buyDate: '买入日期',
    buyReason: '买入理由',
    assetTag: '资产标签',
    riskPreference: '风险偏好',
  }
  for (const [key, value] of Object.entries(values)) {
    if (key === 'positionState') {
      await page.getByLabel('纪律状态').selectOption(value)
      continue
    }
    const label = labels[key]
    if (label) {
      await page.getByLabel(label).fill(value)
    }
  }
}

function positionBySymbol(positions: any[], symbol: string) {
  const position = positions.find((item) => item.symbol === symbol)
  expect(position).toBeTruthy()
  return position
}

function bucketValue(positions: any[], assetTag: string) {
  return positions.filter((item) => item.asset_tag === assetTag).reduce((sum, item) => sum + Number(item.market_value || 0), 0)
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
