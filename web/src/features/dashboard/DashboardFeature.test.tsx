import { render, screen, waitFor, cleanup } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { DashboardFeature } from './DashboardFeature'

vi.mock('../../services/dashboard', () => ({
  getDashboardToday: vi.fn(),
}))

vi.mock('../../services/dailyDisciplineReport', () => ({
  getTodayDailyDisciplineReport: vi.fn(),
}))

import { getDashboardToday } from '../../services/dashboard'
import { getTodayDailyDisciplineReport } from '../../services/dailyDisciplineReport'

describe('DashboardFeature', () => {
  afterEach(() => {
    cleanup()
    vi.clearAllMocks()
  })

  it('maps unknown market states to safe display text', async () => {
    vi.mocked(getTodayDailyDisciplineReport).mockRejectedValue(new Error('report unavailable'))
    vi.mocked(getDashboardToday).mockResolvedValue({
      request_id: 'req_1',
      data: {
        dashboard_state: 'normal',
        discipline_status: '正常',
        data_updated_at: '2026-06-01T00:00:00Z',
        portfolio_summary: {
          total_assets: 1000,
          cash_ratio: 0.2,
          high_risk_ratio: 0.1,
          position_count: 1,
        },
        market_summary: {
          pe_percentile: 0.3,
          pb_percentile: 0.4,
          sentiment_state: 'future_sentiment_state',
          liquidity_state: 'future_liquidity_state',
        },
        triggered_rules: [],
        decision_summary: {
          verdict: '维持观察',
          final_verdict_status: 'hold',
          prohibited_actions: [],
          optional_actions: ['继续观察'],
          action_required: false,
          confirmation_status: 'not_required',
        },
      },
    })

    render(<MemoryRouter><DashboardFeature /></MemoryRouter>)

    await waitFor(() => expect(screen.getByText('市场状态')).toBeInTheDocument())
    expect(screen.getByRole('heading', { name: '今日纪律' })).toBeInTheDocument()
    expect(screen.getByRole('region', { name: '纪律报告概览' })).toHaveClass('reference-hero')
    expect(screen.getByRole('region', { name: '下一步人工动作' })).toHaveClass('reference-action-queue')
    expect(screen.getByRole('region', { name: '状态总览' })).toHaveClass('reference-metric-grid')
    expect(screen.getByRole('region', { name: '持仓与资金快照' })).toHaveClass('reference-snapshot-strip')
    expect(screen.getByRole('region', { name: '最近咨询 · 解释预览' })).toHaveClass('reference-progress-tracker')
    expect(screen.getByRole('region', { name: '证据与规则快照' })).toHaveClass('reference-checklist')
    expect(screen.getByRole('region', { name: '今日纪律状态' })).toHaveTextContent('维持观察')
    expect(screen.getByRole('region', { name: '下一步人工动作' })).toHaveTextContent('查看数据质量')
    expect(screen.getByRole('region', { name: '今日信号摘要' })).toHaveTextContent('数据可信度')
    expect(screen.getAllByText('未知状态')).toHaveLength(2)
    expect(screen.queryByText('future_sentiment_state')).not.toBeInTheDocument()
    expect(screen.queryByText('future_liquidity_state')).not.toBeInTheDocument()
  })

  it('shows decision detail link instead of an inactive confirmation form', async () => {
    vi.mocked(getTodayDailyDisciplineReport).mockRejectedValue(new Error('report unavailable'))
    vi.mocked(getDashboardToday).mockResolvedValue({
      request_id: 'req_action',
      data: {
        dashboard_state: 'normal',
        discipline_status: '正常',
        data_updated_at: '2026-06-01T00:00:00Z',
        portfolio_summary: {
          total_assets: 1000,
          cash_ratio: 0.2,
          high_risk_ratio: 0.1,
          position_count: 1,
        },
        market_summary: {
          sentiment_state: 'neutral',
          liquidity_state: 'normal',
        },
        triggered_rules: [],
        decision_summary: {
          decision_id: 'decision_action',
          verdict: '允许按纪律处理',
          final_verdict_status: 'buy_allowed',
          prohibited_actions: [],
          optional_actions: ['继续观察'],
          action_required: true,
          confirmation_status: 'pending',
        },
      },
    })

    render(<MemoryRouter><DashboardFeature /></MemoryRouter>)

    await waitFor(() => expect(screen.getByRole('region', { name: '下一步人工动作' })).toHaveTextContent('查看决策详情'))
    expect(screen.getByRole('region', { name: '今日纪律状态' })).toHaveTextContent('允许按纪律处理')
    await waitFor(() => expect(screen.getByText('前往决策详情确认')).toBeInTheDocument())
    expect(screen.getByText('前往决策详情确认')).toHaveAttribute('href', '/decisions/decision_action')
    expect(screen.queryByRole('button', { name: '已手动执行' })).not.toBeInTheDocument()
  })

  it('shows today discipline report safety diagnostics without breaking the dashboard', async () => {
    vi.mocked(getDashboardToday).mockResolvedValue({
      request_id: 'req_report_dashboard',
      data: {
        dashboard_state: 'insufficient_data',
        discipline_status: '信息不足',
        data_updated_at: '2026-06-08T00:00:00Z',
        portfolio_summary: {
          total_assets: 0,
          cash_ratio: 0,
          high_risk_ratio: 0,
          position_count: 0,
        },
        market_summary: {
          sentiment_state: 'unknown',
          liquidity_state: 'unknown',
        },
        triggered_rules: [],
        decision_summary: {
          verdict: '等待数据补齐后生成今日纪律建议。',
          final_verdict_status: 'insufficient_data',
          prohibited_actions: ['暂停交易类建议'],
          optional_actions: ['刷新数据'],
          action_required: false,
          confirmation_status: 'not_required',
        },
      },
    })
    vi.mocked(getTodayDailyDisciplineReport).mockResolvedValue({ request_id: 'rid_today', data: {
      report_id: 'daily_report:2026-06-08:holdings:v1',
      local_date: '2026-06-08',
      scope: 'holdings',
      status: 'insufficient_data',
      summary: '缺少持仓快照，今日纪律报告仅展示数据缺口。',
      missing_action: '请先录入本地账户和当前持仓。',
      failure_code: 'missing_prerequisites',
      failure_reason: '本地持仓为空。',
      decision_link: '/decisions/decision_gap',
      audit_link: '/audit?input_ref=gap',
      risk_alerts: [
        {
          alert_id: 'risk_today_1',
          risk_type: 'valuation_high',
          severity: 'warning',
          sop_status: 'active',
          symbol: '510300',
          trigger_summary: '估值分位处于高位',
          prohibited_actions: ['新增买入'],
          suggested_actions: ['人工复核分批止盈'],
          link: '/risk-alerts/risk_today_1',
          safety_note: '风险预警只用于本地人工复核，不会自动交易。',
          created_at: '2026-06-08T02:00:00Z',
          updated_at: '2026-06-08T02:00:00Z',
        },
      ],
      evidence: { evidence_count: 0, independent_source_count: 0, high_grade_independent_source_count: 0 },
      trend: { success_count: 0, degraded_count: 0, failed_count: 1, insufficient_data_count: 2 },
      safety_note: '不会自动执行交易，需人工复核。',
    } })

    render(<MemoryRouter><DashboardFeature /></MemoryRouter>)

    await waitFor(() => expect(screen.getByText('今日纪律报告')).toBeInTheDocument())
    expect(screen.getByRole('region', { name: '今日纪律状态' })).toHaveTextContent('缺少持仓快照，今日纪律报告仅展示数据缺口。')
    expect(screen.getByRole('region', { name: '下一步人工动作' })).toHaveTextContent('维护本地账户与持仓')
    expect(screen.getByRole('region', { name: '今日信号摘要' })).toHaveTextContent('证据不足或尚未完成核验')
    expect(screen.getAllByText('缺少持仓快照，今日纪律报告仅展示数据缺口。').length).toBeGreaterThanOrEqual(1)
    expect(screen.getAllByText('数据不足').length).toBeGreaterThan(0)
    expect(screen.queryByText('本地持仓为空。')).not.toBeInTheDocument()
    expect(screen.getByText('缺少本地账户或持仓。')).toBeInTheDocument()
    expect(screen.getAllByText('请先录入本地账户和当前持仓。').length).toBeGreaterThanOrEqual(1)
    expect(screen.getByRole('link', { name: '前往账户初始化' })).toHaveAttribute('href', '/positions')
    expect(screen.getByText('不会自动执行交易，需人工复核。')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '查看决策详情' })).toHaveAttribute('href', '/decisions/decision_gap')
    expect(screen.getByRole('link', { name: '查看审计详情' })).toHaveAttribute('href', '/audit?input_ref=gap')
    expect(screen.getByText('今日风险预警')).toBeInTheDocument()
    expect(screen.getByText('估值分位处于高位')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '查看风险预警' })).toHaveAttribute('href', '/risk-alerts/risk_today_1')
  })

  it('keeps the main dashboard visible and shows a safe report error when today report loading fails', async () => {
    vi.mocked(getDashboardToday).mockResolvedValue({
      request_id: 'req_report_failed_dashboard',
      data: {
        dashboard_state: 'normal',
        discipline_status: '正常',
        data_updated_at: '2026-06-08T00:00:00Z',
        portfolio_summary: {
          total_assets: 1000,
          cash_ratio: 0.2,
          high_risk_ratio: 0.1,
          position_count: 1,
        },
        market_summary: {
          sentiment_state: 'neutral',
          liquidity_state: 'normal',
        },
        triggered_rules: [],
        decision_summary: {
          verdict: '维持观察',
          final_verdict_status: 'hold',
          prohibited_actions: [],
          optional_actions: ['继续观察'],
          action_required: false,
          confirmation_status: 'not_required',
        },
      },
    })
    vi.mocked(getTodayDailyDisciplineReport).mockRejectedValue(new Error('internal stack trace'))

    render(<MemoryRouter><DashboardFeature /></MemoryRouter>)

    await waitFor(() => expect(screen.getAllByText('维持观察').length).toBeGreaterThanOrEqual(1))
    expect(screen.getByRole('region', { name: '下一步人工动作' })).toHaveTextContent('查看数据质量')
    expect(screen.getByText('今日纪律报告')).toBeInTheDocument()
    expect(screen.getByRole('alert')).toHaveTextContent('每日纪律报告加载失败')
    expect(screen.queryByText(/internal stack trace/)).not.toBeInTheDocument()
  })
})
