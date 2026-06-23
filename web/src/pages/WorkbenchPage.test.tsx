import { cleanup, render, screen, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { APIClientError } from '../services/client'
import { WorkbenchPage } from './WorkbenchPage'

vi.mock('../services/dashboard', () => ({
  getDashboardToday: vi.fn(),
}))

vi.mock('../services/dailyDisciplineReport', () => ({
  getTodayDailyDisciplineReport: vi.fn(),
}))

vi.mock('../services/portfolio', () => ({
  getPortfolioCurrent: vi.fn(),
}))

vi.mock('../services/riskAlert', () => ({
  listRiskAlerts: vi.fn(),
}))

vi.mock('../services/rule', () => ({
  listRuleProposals: vi.fn(),
}))

vi.mock('../services/review', () => ({
  getReviewSummary: vi.fn(),
}))

import { getDashboardToday } from '../services/dashboard'
import { getTodayDailyDisciplineReport } from '../services/dailyDisciplineReport'
import { getPortfolioCurrent } from '../services/portfolio'
import { getReviewSummary } from '../services/review'
import { listRiskAlerts } from '../services/riskAlert'
import { listRuleProposals } from '../services/rule'

describe('WorkbenchPage', () => {
  afterEach(() => {
    cleanup()
    vi.clearAllMocks()
  })

  it('aggregates today, portfolio, risk, rules, review, and consultation navigation', async () => {
    vi.mocked(getDashboardToday).mockResolvedValue({
      request_id: 'rid_dashboard',
      data: {
        dashboard_state: 'normal',
        discipline_status: 'ok',
        data_updated_at: '2026-06-16T09:30:00Z',
        portfolio_summary: { total_assets: 100000, cash_ratio: 0.32, high_risk_ratio: 0.12, position_count: 2 },
        market_summary: { sentiment_state: 'neutral', liquidity_state: 'normal' },
        triggered_rules: [],
        decision_summary: {
          decision_id: 'decision_today',
          verdict: 'hold',
          final_verdict_status: 'normal',
          prohibited_actions: ['追涨'],
          optional_actions: ['复核仓位'],
          action_required: true,
          confirmation_status: 'pending',
        },
      },
    })
    vi.mocked(getTodayDailyDisciplineReport).mockResolvedValue({
      request_id: 'rid_report',
      data: {
        report_id: 'report_today',
        local_date: '2026-06-16',
        scope: 'daily',
        status: 'success',
        summary: '今日纪律报告已生成',
        final_verdict: '继续持有，等待人工复核',
        verdict_status: 'normal',
        evidence: { evidence_count: 4, independent_source_count: 3, high_grade_independent_source_count: 2 },
        trend: { success_count: 3, degraded_count: 0, failed_count: 0, insufficient_data_count: 0 },
        safety_note: '仅展示本地纪律结果。',
      },
    })
    vi.mocked(getPortfolioCurrent).mockResolvedValue({
      request_id: 'rid_portfolio',
      data: {
        snapshot: {
          snapshot_id: 'snap_1',
          snapshot_time: '2026-06-16T09:30:00Z',
          cash: 32000,
          total_assets: 100000,
          cash_ratio: 0.32,
          high_risk_ratio: 0.12,
          position_count: 2,
        },
        positions: [],
      },
    })
    vi.mocked(listRiskAlerts).mockResolvedValue({
      request_id: 'rid_risk',
      data: {
        total: 1,
        items: [{
          alert_id: 'risk_1',
          risk_type: 'valuation_high',
          severity: 'warning',
          sop_status: 'active',
          symbol: '510300',
          trigger_summary: '估值分位偏高，暂停新增买入',
          safety_note: '仅用于人工复核。',
          created_at: '2026-06-16T09:30:00Z',
          updated_at: '2026-06-16T09:30:00Z',
        }],
      },
    })
    vi.mocked(listRuleProposals).mockResolvedValue({
      request_id: 'rid_rule',
      data: {
        total: 1,
        items: [{
          proposal_id: 'proposal_1',
          proposal_type: 'threshold',
          status: 'pending_final_confirm',
          title: '提高样本代表性阈值',
          proposal_version: 'p1',
          sample_count: 12,
          created_at: '2026-06-16T09:30:00Z',
        }],
      },
    })
    vi.mocked(getReviewSummary).mockResolvedValue({
      request_id: 'rid_review',
      data: {
        period: 'quarterly',
        decision_count: 8,
        confirmation_count: 5,
        executed_manually_count: 2,
        planned_count: 3,
        error_case_count: 1,
        rule_proposal_count: 1,
        audit_event_count: 6,
        recent_decisions: [],
      },
    })

    render(<MemoryRouter><WorkbenchPage /></MemoryRouter>)

    expect(await screen.findByRole('heading', { name: '用户决策工作台' })).toBeInTheDocument()
    expect(screen.getByRole('region', { name: '今日纪律状态' })).toHaveTextContent('继续持有，等待人工复核')
    expect(screen.getByRole('region', { name: '下一步人工动作' })).toHaveTextContent('查看决策详情')
    expect(screen.getByRole('region', { name: '下一步人工动作' })).toHaveTextContent('查看今日纪律报告')
    expect(screen.getByRole('region', { name: '今日信号摘要' })).toHaveTextContent('数据可信度')
    expect(screen.getByText('今日先看')).toBeInTheDocument()
    expect(screen.getByText('今日纪律报告已生成')).toBeInTheDocument()
    expect(screen.getByText('最终裁决：继续持有，等待人工复核')).toBeInTheDocument()
    expect(screen.getByText('组合与风险')).toBeInTheDocument()
    expect(screen.getByText('总资产：¥100,000.00')).toBeInTheDocument()
    expect(screen.getByText('活跃风险：1')).toBeInTheDocument()
    expect(screen.getAllByText('规则与复盘').length).toBeGreaterThanOrEqual(1)
    expect(screen.getByText('待确认规则：1')).toBeInTheDocument()
    expect(screen.getByText('复盘决策：8')).toBeInTheDocument()
    expect(screen.getByText('主动咨询入口')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '查看纪律报告' })).toHaveAttribute('href', '/daily-discipline/reports/report_today')
    expect(screen.getByRole('link', { name: '查看持仓' })).toHaveAttribute('href', '/positions')
    expect(screen.getByRole('link', { name: '查看风险预警' })).toHaveAttribute('href', '/risk-alerts')
    expect(screen.getByRole('link', { name: '查看规则提案' })).toHaveAttribute('href', '/rules')
    expect(screen.getByRole('link', { name: '查看复盘摘要' })).toHaveAttribute('href', '/review')
    expect(screen.getByRole('link', { name: '查看决策闭环' })).toHaveAttribute('href', '/decision-loop')
    expect(screen.getByRole('link', { name: '查看审计' })).toHaveAttribute('href', '/audit')
    expect(screen.getByRole('link', { name: '发起主动咨询' })).toHaveAttribute('href', '/consultation')
    expect(listRiskAlerts).toHaveBeenCalledWith({ statuses: ['active', 'escalated'] })
  })

  it('shows safe empty states without action buttons', async () => {
    vi.mocked(getDashboardToday).mockRejectedValue(new APIClientError({ requestId: 'rid_dashboard', code: 'DATA_REQUIRED', message: '需要先录入账户与持仓数据。', displayState: 'first_use' }))
    vi.mocked(getTodayDailyDisciplineReport).mockRejectedValue(new APIClientError({ requestId: 'rid_report', code: 'DATA_REQUIRED', message: '需要先录入账户与持仓数据。', displayState: 'first_use' }))
    vi.mocked(getPortfolioCurrent).mockResolvedValue({
      request_id: 'rid_portfolio_empty',
      data: {
        snapshot: { snapshot_id: 'empty', snapshot_time: '2026-06-16T09:30:00Z', cash: 0, total_assets: 0, cash_ratio: 0, high_risk_ratio: 0, position_count: 0 },
        positions: [],
      },
    })
    vi.mocked(listRiskAlerts).mockResolvedValue({ request_id: 'rid_risk_empty', data: { total: 0, items: [] } })
    vi.mocked(listRuleProposals).mockResolvedValue({ request_id: 'rid_rule_empty', data: { total: 0, items: [] } })
    vi.mocked(getReviewSummary).mockResolvedValue({
      request_id: 'rid_review_empty',
      data: { decision_count: 0, confirmation_count: 0, executed_manually_count: 0, planned_count: 0, error_case_count: 0, rule_proposal_count: 0, audit_event_count: 0, recent_decisions: [] },
    })

    render(<MemoryRouter><WorkbenchPage /></MemoryRouter>)

    await waitFor(() => expect(screen.getAllByText('需要先录入账户与持仓数据。').length).toBeGreaterThanOrEqual(1))
    expect(screen.getByRole('region', { name: '今日纪律状态' })).toHaveTextContent('等待本地数据补齐后生成今日纪律建议。')
    expect(screen.getByRole('region', { name: '下一步人工动作' })).toHaveTextContent('维护本地账户与持仓')
    expect(screen.getByRole('region', { name: '下一步人工动作' })).toHaveTextContent('查看数据质量')
    expect(screen.getByText('暂无持仓快照，先完成本地账户校准。')).toBeInTheDocument()
    expect(screen.getByText('暂无活跃风险预警。')).toBeInTheDocument()
    expect(screen.getAllByText('暂无待确认规则提案。').length).toBeGreaterThanOrEqual(1)
    expect(screen.getByText('暂无复盘活动数据。')).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: /交易|下单|确认|应用/ })).not.toBeInTheDocument()
  })

  it('surfaces degraded service states and keeps consultation manual', async () => {
    vi.mocked(getDashboardToday).mockResolvedValue({
      request_id: 'rid_dashboard_degraded',
      data: {
        dashboard_state: 'data_source_unavailable',
        discipline_status: 'degraded',
        data_updated_at: '2026-06-16T09:30:00Z',
        portfolio_summary: { total_assets: 0, cash_ratio: 0, high_risk_ratio: 0, position_count: 0 },
        market_summary: { sentiment_state: 'unknown', liquidity_state: 'degraded' },
        triggered_rules: [{ rule_id: 'rule_source', rule_name: '数据源降级', severity: 'warning', description: '行情源不可用' }],
        decision_summary: { verdict: 'pause', final_verdict_status: 'degraded', prohibited_actions: ['交易类建议'], optional_actions: [], action_required: true, confirmation_status: 'pending' },
      },
    })
    vi.mocked(getTodayDailyDisciplineReport).mockResolvedValue({
      request_id: 'rid_report_degraded',
      data: {
        report_id: 'report_degraded',
        local_date: '2026-06-16',
        scope: 'daily',
        status: 'degraded',
        summary: '数据源降级，仅展示已有本地事实',
        evidence: { evidence_count: 1, independent_source_count: 0, high_grade_independent_source_count: 0 },
        trend: { success_count: 1, degraded_count: 1, failed_count: 0, insufficient_data_count: 0 },
        safety_note: '降级时仅用于人工复核。',
      },
    })
    vi.mocked(getPortfolioCurrent).mockRejectedValue(new APIClientError({ requestId: 'rid_portfolio', code: 'DATA_STALE', message: '本地数据已过期，请刷新后再查看。', displayState: 'insufficient_data' }))
    vi.mocked(listRiskAlerts).mockRejectedValue(new APIClientError({ requestId: 'rid_risk', code: 'DATA_SOURCE_UNAVAILABLE', message: '数据源暂不可用，请检查数据源状态。', displayState: 'data_source_unavailable' }))
    vi.mocked(listRuleProposals).mockRejectedValue(new APIClientError({ requestId: 'rid_rule', code: 'VECTOR_INDEX_UNAVAILABLE', message: '索引暂不可用，请稍后重试或重建索引。', displayState: 'insufficient_data' }))
    vi.mocked(getReviewSummary).mockRejectedValue(new APIClientError({ requestId: 'rid_review', code: 'ANALYST_UNAVAILABLE', message: '分析服务暂不可用，页面仅展示规则与已有数据。', displayState: 'insufficient_data' }))

    render(<MemoryRouter><WorkbenchPage /></MemoryRouter>)

    await waitFor(() => expect(screen.getAllByText('数据源降级，仅展示已有本地事实').length).toBeGreaterThanOrEqual(1))
    expect(screen.getByRole('region', { name: '今日纪律状态' })).toHaveTextContent('数据源降级，仅展示已有本地事实')
    expect(screen.getByRole('region', { name: '今日信号摘要' })).toHaveTextContent('证据 1 条')
    expect(screen.getByText('行情源不可用')).toBeInTheDocument()
    expect(screen.getAllByText('本地数据已过期，请刷新后再查看。').length).toBeGreaterThanOrEqual(1)
    expect(screen.getAllByText('数据源暂不可用，请检查数据源状态。').length).toBeGreaterThanOrEqual(1)
    expect(screen.getAllByText('索引暂不可用，请稍后重试或重建索引。').length).toBeGreaterThanOrEqual(1)
    expect(screen.getAllByText('分析服务暂不可用，页面仅展示规则与已有数据。').length).toBeGreaterThanOrEqual(1)
    expect(screen.getByText('主动咨询由你提交问题，系统只生成分析材料；最终动作仍由你线下决定。')).toBeInTheDocument()
  })

  it('does not expose forbidden automatic action affordances', async () => {
    vi.mocked(getDashboardToday).mockResolvedValue({
      request_id: 'rid_dashboard_safe',
      data: {
        dashboard_state: 'normal',
        discipline_status: 'ok',
        data_updated_at: '2026-06-16T09:30:00Z',
        portfolio_summary: { total_assets: 0, cash_ratio: 0, high_risk_ratio: 0, position_count: 0 },
        market_summary: { sentiment_state: 'neutral', liquidity_state: 'normal' },
        triggered_rules: [],
        decision_summary: { verdict: 'hold', final_verdict_status: 'normal', prohibited_actions: [], optional_actions: [], action_required: false, confirmation_status: 'none' },
      },
    })
    vi.mocked(getTodayDailyDisciplineReport).mockResolvedValue({
      request_id: 'rid_report_safe',
      data: {
        report_id: 'report_safe',
        local_date: '2026-06-16',
        scope: 'daily',
        status: 'success',
        summary: '安全文案检查',
        evidence: { evidence_count: 0, independent_source_count: 0, high_grade_independent_source_count: 0 },
        trend: { success_count: 0, degraded_count: 0, failed_count: 0, insufficient_data_count: 0 },
        safety_note: '只读。',
      },
    })
    vi.mocked(getPortfolioCurrent).mockResolvedValue({
      request_id: 'rid_portfolio_safe',
      data: {
        snapshot: { snapshot_id: 'safe', snapshot_time: '2026-06-16T09:30:00Z', cash: 0, total_assets: 0, cash_ratio: 0, high_risk_ratio: 0, position_count: 0 },
        positions: [],
      },
    })
    vi.mocked(listRiskAlerts).mockResolvedValue({ request_id: 'rid_risk_safe', data: { total: 0, items: [] } })
    vi.mocked(listRuleProposals).mockResolvedValue({ request_id: 'rid_rule_safe', data: { total: 0, items: [] } })
    vi.mocked(getReviewSummary).mockResolvedValue({
      request_id: 'rid_review_safe',
      data: { decision_count: 0, confirmation_count: 0, executed_manually_count: 0, planned_count: 0, error_case_count: 0, rule_proposal_count: 0, audit_event_count: 0, recent_decisions: [] },
    })

    render(<MemoryRouter><WorkbenchPage /></MemoryRouter>)

    await waitFor(() => expect(screen.getAllByText('安全文案检查').length).toBeGreaterThanOrEqual(1))
    const forbidden = /自动交易|一键交易|代下单|券商下单|券商接口|自动规则应用|外部推送|收益承诺/
    expect(screen.queryByRole('button', { name: forbidden })).not.toBeInTheDocument()
    expect(screen.queryByRole('link', { name: forbidden })).not.toBeInTheDocument()
    expect(document.body.textContent).not.toMatch(/一键交易|代下单|券商下单|券商接口|自动规则应用|外部推送|收益承诺/)
  })
})
