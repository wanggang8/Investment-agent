import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'
import { DisciplineStatus } from './DisciplineStatus'
import type { DashboardTodayResponse } from '../../types/dashboard'

const dashboard: DashboardTodayResponse = {
  dashboard_state: 'high_risk',
  discipline_status: '高危纪律状态',
  data_updated_at: '2026-05-31T00:00:00Z',
  portfolio_summary: { total_assets: 1000, cash_ratio: 0.1, high_risk_ratio: 0.7, position_count: 2 },
  market_summary: { sentiment_state: 'panic', liquidity_state: 'tight' },
  triggered_rules: [],
  decision_summary: {
    verdict: '禁止新增买入',
    final_verdict_status: 'high_risk',
    prohibited_actions: ['buy'],
    optional_actions: ['watch'],
    action_required: false,
    confirmation_status: 'not_required',
  },
}

describe('DisciplineStatus', () => {
  afterEach(() => cleanup())

  it('renders high risk discipline state', () => {
    render(<DisciplineStatus dashboard={dashboard} />)

    expect(screen.getByText('高危纪律状态')).toBeInTheDocument()
    expect(screen.getByText(/已触发高危纪律状态/)).toBeInTheDocument()
    expect(screen.getByText('高风险')).toBeInTheDocument()
  })

  it('renders normal state without missing-data warning', () => {
    render(<DisciplineStatus dashboard={{ ...dashboard, dashboard_state: 'normal', discipline_status: '正常' }} />)

    expect(screen.getByText('今日未触发纪律红线。')).toBeInTheDocument()
    expect(screen.queryByText(/缺少/)).not.toBeInTheDocument()
    expect(screen.getAllByText('正常').length).toBeGreaterThan(0)
  })
})
