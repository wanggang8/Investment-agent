import { expect, test } from '@playwright/test'
import type { Page } from '@playwright/test'
import { mkdirSync, writeFileSync } from 'node:fs'
import path from 'node:path'

const artifactDir = process.env.P84_ARTIFACT_DIR || path.resolve(process.cwd(), '../docs/release/ui-audit-assets/2026-06-22-p84-portfolio-confirmation')
const shouldCapture = process.env.P84_CAPTURE_SCREENSHOTS === '1'

test('P84 portfolio and confirmation UI actions persist local data and downstream readbacks', async ({ page, request }) => {
  test.setTimeout(240_000)
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

  const portfolioResult = await runPortfolioActions(page, request)
  const confirmationResult = await runDecisionConfirmation(page, request)
  const downstream = await runDownstreamReadbacks(page)

  expect(consoleErrors).toEqual([])
  expect(pageErrors).toEqual([])
  expect(failedApiResponses).toEqual([])
  await assertNoForbiddenVisibleAffordance(page)

  writeResult({
    generated_at: new Date().toISOString(),
    status: 'passed',
    portfolio: portfolioResult,
    confirmation: confirmationResult,
    downstream,
    console_errors: consoleErrors,
    page_errors: pageErrors,
    failed_api_responses: failedApiResponses,
    safety_boundaries: ['no broker UI', 'no automatic confirmation', 'no order placement', 'no external push'],
  })
})

async function runPortfolioActions(page: Page, request: any) {
  await page.goto('/positions')
  await expect(page.getByRole('heading', { name: '组合与持仓维护', exact: true })).toBeVisible()

  await fillPortfolioForm(page, {
    cash: '12000',
    totalAssets: '72000',
    symbol: '510300',
    name: '沪深300ETF',
    quantity: '15000',
    costPrice: '3.6',
    currentPrice: '4',
    buyReason: 'P84 本地核心仓位校准',
    assetTag: 'core',
    riskPreference: 'steady',
  })
  await page.getByRole('button', { name: '保存本地校准' }).click()
  await expect(page.getByText('账户校准已保存为本地事实；不会连接交易接口。')).toBeVisible()
  await capture(page, 'p84-portfolio-calibration.png')

  await fillPortfolioForm(page, {
    quantity: '15250',
    costPrice: '3.62',
    currentPrice: '4.02',
    buyReason: 'P84 人工复核后更新核心仓位',
    assetTag: 'core',
  })
  await page.getByRole('button', { name: '保存持仓编辑' }).click()
  await expect(page.getByText('持仓编辑已保存为本地事实。')).toBeVisible()

  await fillPortfolioForm(page, {
    symbol: '159915',
    name: '创业板ETF',
    quantity: '9000',
    costPrice: '2.10',
    currentPrice: '2.30',
    buyReason: 'P84 卫星仓位批量导入',
    assetTag: 'satellite',
  })
  await page.getByRole('button', { name: '校验批量导入' }).click()
  await expect(page.getByText('导入校验完成：有效 1 行，无效 0 行。')).toBeVisible()
  await page.getByRole('button', { name: '确认批量导入' }).click()
  await expect(page.getByText('批量导入已确认并保存。')).toBeVisible()

  await fillPortfolioForm(page, {
    symbol: '511880',
    name: '货币基金',
    quantity: '3000',
    costPrice: '1',
    currentPrice: '1',
    buyReason: 'P84 现金管理线下买入记录',
    assetTag: 'cash',
  })
  await page.getByLabel('线下交易类型').selectOption('buy')
  await page.getByRole('button', { name: '记录线下交易' }).click()
  await expect(page.getByText('线下交易已记录为本地事实。')).toBeVisible()

  await fillPortfolioForm(page, {
    symbol: '511880',
    name: '货币基金',
    quantity: '3100',
    costPrice: '1',
    currentPrice: '1',
    buyReason: 'P84 本地事实修正审计',
    assetTag: 'cash',
  })
  await page.getByRole('button', { name: '记录修正审计' }).click()
  await expect(page.getByText('错误修正已保存为本地事实。')).toBeVisible()

  await page.reload()
  await expect(page.getByText('510300').first()).toBeVisible()
  await expect(page.getByText('159915').first()).toBeVisible()
  await expect(page.getByText('511880').first()).toBeVisible()
  await capture(page, 'p84-portfolio-after-actions.png')

  const response = await request.get('/api/v1/portfolio/current')
  await expect(response).toBeOK()
  const body = await response.json()
  expect(body.data?.positions?.some((item: any) => item.symbol === '510300')).toBe(true)
  expect(body.data?.positions?.some((item: any) => item.symbol === '159915')).toBe(true)
  expect(body.data?.positions?.some((item: any) => item.symbol === '511880')).toBe(true)

  return {
    snapshot: body.data?.snapshot,
    positions: body.data?.positions?.map((item: any) => ({
      symbol: item.symbol,
      quantity: item.quantity,
      current_price: item.current_price,
      market_value: item.market_value,
      asset_tag: item.asset_tag,
    })),
  }
}

async function runDecisionConfirmation(page: Page, request: any) {
  await page.goto('/decisions/decision_p84_pending')
  await expect(page.getByText('当前确认状态：待确认')).toBeVisible()
  await page.getByRole('button', { name: '已手动执行' }).click()
  const form = page.locator('[aria-label="确认表单"]')
  await form.getByLabel('标的代码').fill('510300')
  await form.getByLabel('线下动作').selectOption('sell')
  await form.getByLabel('数量').fill('100')
  await form.getByLabel('价格').fill('4.05')
  await form.getByLabel('费用').fill('1')
  await form.getByLabel('执行时间').fill(pastDateTimeLocal())
  await form.getByLabel('备注').fill('P84 用户线下卖出后回填确认')
  await form.getByRole('button', { name: '提交确认' }).click()
  await expect(page.getByText('确认已记录')).toBeVisible()
  await capture(page, 'p84-decision-confirmed.png')

  const detail = await request.get('/api/v1/decisions/decision_p84_pending')
  await expect(detail).toBeOK()
  const body = await detail.json()
  expect(body.data?.user_confirmation?.confirmation_status).toBe('executed_manually')
  return {
    decision_id: body.data?.decision_id,
    confirmation_status: body.data?.user_confirmation?.confirmation_status,
    final_verdict: body.data?.final_verdict?.status,
  }
}

async function runDownstreamReadbacks(page: Page) {
  const checks = [
    { path: '/decision-loop', text: 'decision_p84_pending', screenshot: 'p84-decision-loop-readback.png' },
    { path: '/review', text: 'decision_p84_pending', screenshot: 'p84-review-readback.png' },
    { path: '/audit', text: 'decision_p84_pending', screenshot: 'p84-audit-readback.png' },
    { path: '/workbench', text: '查看持仓', screenshot: 'p84-workbench-readback.png' },
  ]
  const seen: string[] = []
  for (const check of checks) {
    await page.goto(check.path)
    await expect(page.getByText(check.text).first()).toBeVisible()
    await capture(page, check.screenshot)
    seen.push(check.path)
  }
  return { routes: seen }
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
    buyReason: '买入理由',
    assetTag: '资产标签',
    riskPreference: '风险偏好',
  }
  for (const [key, value] of Object.entries(values)) {
    const label = labels[key]
    if (label) {
      await page.getByLabel(label).fill(value)
    }
  }
}

async function assertNoForbiddenVisibleAffordance(page: Page) {
  const forbidden = ['自动交易', '一键交易', '代下单', '券商下单', '外部推送', '自动确认', '自动应用规则']
  const body = await page.locator('body').innerText()
  for (const term of forbidden) {
    if (term === '自动交易' && body.includes('不自动交易')) continue
    if (term === '外部推送' && body.includes('不会外部推送')) continue
    if (term === '自动应用规则' && body.includes('不会自动应用规则')) continue
    expect(body.includes(term)).toBe(false)
  }
}

async function capture(page: Page, name: string) {
  if (!shouldCapture) return
  await page.screenshot({ path: path.join(artifactDir, name), fullPage: true })
}

function writeResult(payload: unknown) {
  writeFileSync(path.join(artifactDir, 'browser-results.json'), JSON.stringify(payload, null, 2))
}

function pastDateTimeLocal() {
  const date = new Date(Date.now() - 60 * 60 * 1000)
  const pad = (value: number) => String(value).padStart(2, '0')
  return `${date.getUTCFullYear()}-${pad(date.getUTCMonth() + 1)}-${pad(date.getUTCDate())}T${pad(date.getUTCHours())}:${pad(date.getUTCMinutes())}`
}
