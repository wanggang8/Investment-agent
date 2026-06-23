import { expect, test } from '@playwright/test'
import { mkdirSync, writeFileSync } from 'node:fs'
import path from 'node:path'

const artifactDir = process.env.P90_ARTIFACT_DIR || path.resolve(process.cwd(), '../docs/release/ui-audit-assets/2026-06-22-p90-capital-flow-provider')
const shouldCapture = process.env.P90_CAPTURE_SCREENSHOTS === '1'

test('P90 capital-flow provider is refreshed through Settings UI and read back through API', async ({ page, request }) => {
  test.setTimeout(240_000)
  mkdirSync(artifactDir, { recursive: true })

  const consoleErrors: string[] = []
  const pageErrors: string[] = []
  const failedApiResponses: Array<{ url: string; status: number; method: string }> = []
  page.on('console', (message) => {
    if (message.type() === 'error' && !message.text().includes('Failed to load resource')) consoleErrors.push(message.text())
  })
  page.on('pageerror', (error) => pageErrors.push(error.message))
  page.on('response', (response) => {
    if (response.url().includes('/api/v1/') && response.status() >= 500) {
      failedApiResponses.push({ url: response.url(), status: response.status(), method: response.request().method() })
    }
  })

  await expect(await request.get('/api/v1/health')).toBeOK()
  const before = await readLatestMarket(request)
  expect(before.market_snapshot_id).toBe('market_p90_seed')
  expect(before.market_metrics?.metadata?.p88_structured_fields?.capital_flow).toBeUndefined()

  await page.goto('/settings')
  await expect(page.getByRole('heading', { name: '设置' })).toBeVisible()
  await expect(page.getByText('资金流向：暂无真实 provider 读回')).toBeVisible()
  await page.getByRole('button', { name: '刷新市场数据' }).click()
  await expect(page.getByText('市场刷新完成；只更新本地行情事实和审计记录，不会执行交易。')).toBeVisible({ timeout: 180_000 })
  await expect(page.getByText(/资金流向：\d{4}-\d{2}-\d{2}；净流入/)).toBeVisible()
  await capture(page, 'p90-capital-flow-settings-readback.png')

  const after = await readLatestMarket(request)
  expect(after.market_snapshot_id).not.toBe('market_p90_seed')
  const capitalFlow = after.market_metrics?.metadata?.p88_structured_fields?.capital_flow
  expect(capitalFlow?.date).toMatch(/^\d{4}-\d{2}-\d{2}$/)
  expect(Number(capitalFlow?.raw_net_flow)).not.toBe(0)
  expect(Number(capitalFlow?.net_inflow) + Number(capitalFlow?.net_outflow)).toBeGreaterThan(0)
  expect([Number(capitalFlow?.net_inflow), Number(capitalFlow?.net_outflow)].filter((value) => value > 0).length).toBe(1)

  expect(consoleErrors).toEqual([])
  expect(pageErrors).toEqual([])
  expect(failedApiResponses).toEqual([])
  for (const name of ['自动交易', '一键交易', '代下单', '券商下单', '自动确认', '自动应用规则']) {
    await expect(page.getByRole('button', { name })).toHaveCount(0)
  }

  writeFileSync(path.join(artifactDir, 'browser-results.json'), JSON.stringify({
    generated_at: new Date().toISOString(),
    status: 'passed',
    pre_refresh: {
      market_snapshot_id: before.market_snapshot_id,
      capital_flow_absent: true,
    },
    provider_readback: {
      market_snapshot_id: after.market_snapshot_id,
      symbol: after.symbol,
      trade_date: after.trade_date,
      capital_flow: capitalFlow,
    },
    ui_paths: ['/settings', '/api/v1/market/snapshots/latest?symbol=600000'],
    safety_boundaries: ['no broker UI', 'no automatic confirmation', 'no order placement', 'no external push', 'no automatic rule application'],
    console_errors: consoleErrors,
    page_errors: pageErrors,
    failed_api_responses: failedApiResponses,
  }, null, 2))
})

async function readLatestMarket(request: any) {
  const response = await request.get('/api/v1/market/snapshots/latest?symbol=600000')
  await expect(response).toBeOK()
  const body = await response.json()
  return body.data
}

async function capture(page: any, name: string) {
  if (!shouldCapture) return
  await page.screenshot({ path: path.join(artifactDir, name), fullPage: true })
}
