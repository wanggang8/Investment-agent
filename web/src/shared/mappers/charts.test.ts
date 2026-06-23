import { describe, expect, it } from 'vitest'
import { buildDashboardChartData, buildPortfolioAllocationData, buildReviewActivityData } from './charts'

const dashboard = {
  portfolio_summary: {
    total_assets: 100000,
    cash_ratio: 0.12,
    high_risk_ratio: 0.08,
    position_count: 4,
  },
  market_summary: {
    pe_percentile: 42,
    pb_percentile: 36,
    sentiment_state: 'neutral',
    liquidity_state: 'normal',
  },
  evidence_summary: {
    source_count: 3,
    highest_source_level: 'A',
    verification_status: 'satisfied',
  },
}

describe('chart mappers', () => {
  it('builds dashboard chart data from API DTO only', () => {
    const charts = buildDashboardChartData(dashboard)

    expect(charts.allocation).toContainEqual({ label: '现金', value: 12, tone: 'safe' })
    expect(charts.risk).toContainEqual({ label: '高风险仓位', value: 8, tone: 'safe' })
    expect(charts.evidence).toContainEqual({ label: '证据来源', value: 3, tone: 'info' })
  })

  it('builds portfolio allocation without local storage fields', () => {
    const chart = buildPortfolioAllocationData([
      { position_id: 'p1', symbol: '510300', name: '沪深300', quantity: 10, cost_price: 1, current_price: 1.2, market_value: 30000, unrealized_profit_ratio: 0.2, position_state: 'normal', asset_tag: 'core' },
      { position_id: 'p2', symbol: '159915', name: '创业板', quantity: 10, cost_price: 1, current_price: 1.1, market_value: 10000, unrealized_profit_ratio: 0.1, position_state: 'frozen_watch', asset_tag: 'satellite' },
    ])

    expect(chart).toEqual([
      { label: '沪深300', value: 75, tone: 'safe' },
      { label: '创业板', value: 25, tone: 'warn' },
    ])
  })

  it('aggregates duplicate portfolio labels for stable chart rows', () => {
    const chart = buildPortfolioAllocationData([
      { position_id: 'p1', symbol: '510300', name: '沪深300ETF', quantity: 100, cost_price: 3, current_price: 4, market_value: 400, unrealized_profit_ratio: 0.33, position_state: 'normal', asset_tag: 'core' },
      { position_id: 'p2', symbol: '510300', name: '沪深300ETF', quantity: 10, cost_price: 3.9, current_price: 4.1, market_value: 41, unrealized_profit_ratio: 0.05, position_state: 'normal', asset_tag: 'satellite' },
      { position_id: 'p3', symbol: '159915', name: '创业板', quantity: 10, cost_price: 2, current_price: 2, market_value: 59, unrealized_profit_ratio: 0, position_state: 'frozen_watch', asset_tag: 'satellite' },
    ])

    expect(chart).toEqual([
      { label: '沪深300ETF', value: 88, tone: 'safe' },
      { label: '创业板', value: 12, tone: 'warn' },
    ])
  })

  it('builds review activity summary', () => {
    expect(buildReviewActivityData({ error_case_count: 2, rule_proposal_count: 1, confirmation_count: 3, audit_event_count: 8 })).toEqual([
      { label: '确认动作', value: 3, tone: 'safe' },
      { label: '错误案例', value: 2, tone: 'danger' },
      { label: '规则提案', value: 1, tone: 'warn' },
      { label: '审计事件', value: 8, tone: 'info' },
    ])
  })
})
