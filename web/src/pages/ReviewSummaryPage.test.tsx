import { render, screen, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { cleanup } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { APIClientError } from '../services/client'
import { ReviewSummaryPage } from './ReviewSummaryPage'

vi.mock('../services/review', () => ({
  getReviewSummary: vi.fn(),
}))

import { getReviewSummary } from '../services/review'

describe('ReviewSummaryPage', () => {
  afterEach(() => {
    cleanup()
    vi.resetAllMocks()
  })

  it('shows API error message', async () => {
    vi.mocked(getReviewSummary).mockRejectedValue(new APIClientError({ requestId: 'rid', code: 'INTERNAL_ERROR', message: '系统暂时无法处理请求，请稍后重试。', displayState: 'generic_failure' }))

    render(<MemoryRouter><ReviewSummaryPage /></MemoryRouter>)

    await waitFor(() => expect(screen.getByText('系统暂时无法处理请求，请稍后重试。')).toBeInTheDocument())
  })

  it('shows empty success state with zero metrics', async () => {
    vi.mocked(getReviewSummary).mockResolvedValue({ request_id: 'rid', data: { decision_count: 0, confirmation_count: 0, executed_manually_count: 0, planned_count: 0, error_case_count: 0, rule_proposal_count: 0, audit_event_count: 0, recent_decisions: [] } })

    render(<MemoryRouter><ReviewSummaryPage /></MemoryRouter>)

    await waitFor(() => expect(screen.getByText('复盘活动')).toBeInTheDocument())
    expect(screen.getByText('建议数量')).toBeInTheDocument()
    expect(screen.getAllByText('0').length).toBeGreaterThanOrEqual(7)
    expect(screen.getByText('暂无规则建议。')).toBeInTheDocument()
    expect(screen.getByText('暂无追踪记录。')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '查看决策闭环' })).toHaveAttribute('href', '/decision-loop')
    expect(screen.queryByText('系统暂时无法处理请求，请稍后重试。')).not.toBeInTheDocument()
  })

  it('shows periodic summary, rule suggestions and tracking links', async () => {
    vi.mocked(getReviewSummary).mockResolvedValue({
      request_id: 'rid',
      data: {
        period: 'quarterly',
        decision_count: 2,
        confirmation_count: 2,
        executed_manually_count: 1,
        planned_count: 1,
        error_case_count: 1,
        rule_proposal_count: 1,
        audit_event_count: 2,
        rule_hit_count: 3,
        misjudgment_count: 1,
        missing_evidence_count: 1,
        degraded_count: 1,
        rule_suggestions: [{ proposal_id: 'prop_review', title: '季度阈值复盘', status: 'pending_user_confirm', reason: '误判样本触发', can_auto_apply: false }],
        attribution_summaries: [{ decision_id: 'decision_review_1', symbol: '510300', verdict: 'hold', confirmation_status: 'executed_manually', evidence_status: 'failed', workflow_status: 'degraded', outcome: 'missing_evidence' }],
        recurring_error_tags: [{ tag: 'rule_threshold_issue', count: 1 }],
        missing_evidence_themes: [{ status: 'failed', count: 1 }],
        rule_proposal_outcomes: [{ proposal_id: 'prop_review', title: '季度阈值复盘', status: 'pending_user_confirm' }],
        rule_effect_tracking: [{ tracking_id: 'track_review', applied_rule_version: 'v3.1', proposal_id: 'prop_review', period: 'quarterly', hit_count: 4, misjudgment_count: 2, missing_evidence_count: 1, degraded_count: 1, risk_alert_count: 2, trend_direction: 'worsened', metrics: { hit_count: 4, risk_alert_count: 2 }, related_proposal_ids: ['prop_review'], related_audit_event_ids: ['audit_review'], related_risk_alert_ids: ['risk_review'], safety_note: '应用后追踪只读展示，不会自动回滚规则或执行交易。' }],
        degraded_workflows: [{ decision_id: 'decision_review_1', symbol: '510300', status: 'degraded', created_at: '2026-05-01T00:00:00Z' }],
        tracking_links: [{ type: 'rule_proposal', id: 'prop_review', label: '规则提案 prop_review' }],
        recent_decisions: [],
      },
    })

    render(<MemoryRouter><ReviewSummaryPage /></MemoryRouter>)

    await waitFor(() => expect(screen.getByText('季度复盘')).toBeInTheDocument())
    expect(screen.getByText('规则命中')).toBeInTheDocument()
    expect(screen.getByText('缺证据')).toBeInTheDocument()
    expect(screen.getAllByText('季度阈值复盘').length).toBeGreaterThanOrEqual(1)
    expect(screen.getByText('归因摘要')).toBeInTheDocument()
    expect(screen.getByText('decision_review_1')).toBeInTheDocument()
    expect(screen.getByText('缺失证据')).toBeInTheDocument()
    expect(screen.getByText('高频错误标签')).toBeInTheDocument()
    expect(screen.getByText('rule_threshold_issue · 1')).toBeInTheDocument()
    expect(screen.getByText('缺证据主题')).toBeInTheDocument()
    expect(screen.getByText('failed · 1')).toBeInTheDocument()
    expect(screen.getByText('规则提案结果')).toBeInTheDocument()
    expect(screen.getByText('降级工作流')).toBeInTheDocument()
    expect(screen.getByText('decision_review_1 · 510300')).toBeInTheDocument()
    expect(screen.getByText('规则应用后效果追踪')).toBeInTheDocument()
    expect(screen.getByText('v3.1 · 趋势：变差')).toBeInTheDocument()
    expect(screen.getByText('命中 4 · 误判 2 · 缺证据 1 · 降级 1 · 风险预警 2')).toBeInTheDocument()
    expect(JSON.parse(screen.getByLabelText('追踪指标').textContent ?? '{}')).toEqual({ hit_count: 4, risk_alert_count: 2 })
    expect(JSON.parse(screen.getByLabelText('追踪关联提案').textContent ?? '[]')).toEqual(['prop_review'])
    expect(JSON.parse(screen.getByLabelText('追踪关联审计').textContent ?? '[]')).toEqual(['audit_review'])
    expect(JSON.parse(screen.getByLabelText('追踪关联风险预警').textContent ?? '[]')).toEqual(['risk_review'])
    expect(screen.getByText('应用后追踪只读展示，不会自动回滚规则或执行交易。')).toBeInTheDocument()
    expect(screen.getByText('规则变更仍需守门人审计和用户最终确认，不会自动应用。')).toBeInTheDocument()
    expect(screen.getAllByText('待用户确认').length).toBeGreaterThanOrEqual(1)
    const trackingLink = screen.getByRole('link', { name: /规则提案 prop_review/ })
    expect(trackingLink).toHaveAttribute('href', '#rule_proposal-prop_review')
  })

  it('shows ops status panel with degraded and unknown states safely', async () => {
    vi.mocked(getReviewSummary).mockResolvedValue({
      request_id: 'rid',
      data: {
        period: 'monthly',
        decision_count: 1,
        confirmation_count: 0,
        executed_manually_count: 0,
        planned_count: 0,
        error_case_count: 0,
        rule_proposal_count: 0,
        audit_event_count: 1,
        degraded_count: 2,
        ops_status: {
          data_source_status: 'degraded',
          index_status: 'missing',
          review_status: 'unexpected_status',
          explanation: 'VecLite 索引缺失，可通过本地索引任务重建。',
        },
        recent_decisions: [],
      },
    })

    render(<MemoryRouter><ReviewSummaryPage /></MemoryRouter>)

    await waitFor(() => expect(screen.getByText('运维状态')).toBeInTheDocument())
    expect(screen.getByText('数据源')).toBeInTheDocument()
    expect(screen.getAllByText('降级').length).toBeGreaterThanOrEqual(1)
    expect(screen.getByText('索引')).toBeInTheDocument()
    expect(screen.getByText('缺失')).toBeInTheDocument()
    expect(screen.getByText('复盘状态')).toBeInTheDocument()
    expect(screen.getByText('未知状态')).toBeInTheDocument()
    expect(screen.getByText('VecLite 索引缺失，可通过本地索引任务重建。')).toBeInTheDocument()
    expect(screen.getByText('仅展示状态与追踪入口，不执行交易，也不自动应用规则。')).toBeInTheDocument()
    expect(screen.queryByText(/自动交易|自动执行交易|一键交易|代下单|自动应用规则$/)).not.toBeInTheDocument()
  })

  it('shows ops status panel with success failed and explicit empty states', async () => {
    vi.mocked(getReviewSummary).mockResolvedValue({
      request_id: 'rid',
      data: {
        period: 'monthly',
        decision_count: 0,
        confirmation_count: 0,
        executed_manually_count: 0,
        planned_count: 0,
        error_case_count: 0,
        rule_proposal_count: 0,
        audit_event_count: 0,
        ops_status: {
          data_source_status: 'success',
          index_status: 'failed',
          review_status: 'empty',
        },
        recent_decisions: [],
      },
    })

    render(<MemoryRouter><ReviewSummaryPage /></MemoryRouter>)

    await waitFor(() => expect(screen.getByText('运维状态')).toBeInTheDocument())
    expect(screen.getByText('成功')).toBeInTheDocument()
    expect(screen.getByText('失败')).toBeInTheDocument()
    expect(screen.getByText('暂无数据')).toBeInTheDocument()
  })

  it('shows safe ops empty state', async () => {
    vi.mocked(getReviewSummary).mockResolvedValue({ request_id: 'rid', data: { decision_count: 0, confirmation_count: 0, executed_manually_count: 0, planned_count: 0, error_case_count: 0, rule_proposal_count: 0, audit_event_count: 0, recent_decisions: [] } })

    render(<MemoryRouter><ReviewSummaryPage /></MemoryRouter>)

    await waitFor(() => expect(screen.getByText('运维状态')).toBeInTheDocument())
    expect(screen.getByText('暂无运维状态数据。')).toBeInTheDocument()
  })
})
