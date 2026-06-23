import { fireEvent, render, screen, waitFor, cleanup } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { APIClientError } from '../services/client'
import { RulesPage } from './RulesPage'

vi.mock('../services/rule', () => ({
  listRuleProposals: vi.fn(),
  confirmRuleProposal: vi.fn(),
  createSOPAddendumProposal: vi.fn(),
  finalConfirmRuleProposal: vi.fn(),
  getCurrentRule: vi.fn(),
}))

import type { RuleProposal } from '../types/rule'
import { confirmRuleProposal, createSOPAddendumProposal, finalConfirmRuleProposal, getCurrentRule, listRuleProposals } from '../services/rule'

const proposal: RuleProposal = {
  proposal_id: 'prop_1',
  proposal_type: 'risk_rule',
  status: 'pending_final_confirm',
  title: '规则提案',
  proposal_version: 'draft',
  sample_count: 3,
  created_at: '2026-05-31T00:00:00Z',
  reason: '样本表现稳定',
  source_error_case_id: 'err_1',
  before_rule: { content: '旧规则文本', old: 1 },
  after_rule: { content: '新规则文本', next: 2 },
  impact_scope: { scope: 'portfolio' },
  risk_notes: { risk: 'low' },
  audit_result: 'approved',
  audit_summary: '通过',
  effect_validation: {
    validation_id: 'val_1',
    validation_status: 'passed',
    sample_count: 5,
    sample_window: '2026-Q2',
    representativeness_status: 'passed',
    overfit_risk: 'low',
    replay_result: 'passed',
    guardrail_decision: 'passed',
    source_explanation: { source_case_count: 3, related_error_case_ids: ['err_1'] },
    metrics: { hit_count: 4, misjudgment_count: 0 },
    risk_notes: ['本地样本暂未发现不利信号'],
    related_error_cases: ['err_1'],
    related_decision_ids: ['dec_1'],
    related_risk_alert_ids: ['risk_1'],
    related_audit_event_ids: ['audit_1'],
    validation_link: '/rule-effect-validations/val_1',
    safety_note: '规则效果验证只用于本地规则治理，不会自动应用规则或执行交易。',
  },
}

const pendingUserProposal: RuleProposal = {
  ...proposal,
  status: 'pending_user_confirm',
}

describe('RulesPage', () => {
  beforeEach(() => {
    vi.mocked(getCurrentRule).mockResolvedValue({ request_id: 'rule_current', data: { rule_version: 'v3.0', status: 'active', rules: { priority: ['safety_gate', 'risk_limit'], thresholds: { cash_ratio_min: 0.1 } }, effective_at: '2026-01-01T00:00:00Z', created_at: '2026-01-01T00:00:00Z' } })
  })

  afterEach(() => {
    cleanup()
    vi.clearAllMocks()
  })

  it('renders current rule version, priority and thresholds', async () => {
    vi.mocked(listRuleProposals).mockResolvedValue({ request_id: 'req_1', data: { items: [], total: 0 } })

    render(<RulesPage />)

    await waitFor(() => expect(screen.getByText('规则治理状态')).toBeInTheDocument())
    expect(screen.getByText('暂无规则提案')).toBeInTheDocument()
    expect(screen.getByText('查看审计记录')).toBeInTheDocument()
    await waitFor(() => expect(screen.getByText('当前规则库：v3.0')).toBeInTheDocument())
    expect(screen.getByText('裁决优先级：safety_gate、risk_limit')).toBeInTheDocument()
    expect(screen.getByText(/cash_ratio_min/)).toBeInTheDocument()
  })

  it('shows APIClientError message when rule proposals fail to load', async () => {
    vi.mocked(listRuleProposals).mockRejectedValue(new APIClientError({ requestId: 'rid', code: 'INTERNAL_ERROR', message: '系统暂时无法处理请求，请稍后重试。', displayState: 'generic_failure' }))

    render(<RulesPage />)

    await waitFor(() => expect(screen.getByText('系统暂时无法处理请求，请稍后重试。')).toBeInTheDocument())
  })

  it('shows APIClientError message when final confirmation fails', async () => {
    vi.mocked(listRuleProposals).mockResolvedValue({ request_id: 'req_1', data: { items: [proposal], total: 1 } })
    vi.mocked(finalConfirmRuleProposal).mockRejectedValue(new APIClientError({ requestId: 'rid', code: 'INVALID_STATE', message: '当前状态不允许执行该操作。', displayState: 'frozen_watch' }))

    render(<RulesPage />)

    await waitFor(() => expect(screen.getByText('规则提案')).toBeInTheDocument())
    fireEvent.click(screen.getByRole('button', { name: '确认应用到正式规则' }))

    await waitFor(() => expect(screen.getByText('当前状态不允许执行该操作。')).toBeInTheDocument())
  })

  it('supports user confirmation before gatekeeper audit', async () => {
    vi.mocked(listRuleProposals)
      .mockResolvedValueOnce({ request_id: 'req_1', data: { items: [pendingUserProposal], total: 1 } })
      .mockResolvedValueOnce({ request_id: 'req_2', data: { items: [proposal], total: 1 } })
    vi.mocked(confirmRuleProposal).mockResolvedValue({ request_id: 'req_confirm', data: { proposal_id: 'prop_1', status: 'pending_final_confirm' } })

    render(<RulesPage />)

    await waitFor(() => expect(screen.getByRole('button', { name: '确认送审' })).toBeInTheDocument())
    fireEvent.click(screen.getByRole('button', { name: '确认送审' }))

    await waitFor(() => expect(confirmRuleProposal).toHaveBeenCalledWith('prop_1', { confirm: true }))
    await waitFor(() => expect(screen.getByRole('button', { name: '确认应用到正式规则' })).toBeInTheDocument())
  })

  it('creates an SOP addendum proposal from the rules UI and refreshes the list', async () => {
    const sopProposal: RuleProposal = {
      ...pendingUserProposal,
      proposal_id: 'prop_sop',
      proposal_type: 'sop',
      title: 'SOP 补充提案：连续流动性缺口未覆盖',
      reason: '高频未覆盖场景：连续流动性缺口未覆盖',
    }
    vi.mocked(listRuleProposals)
      .mockResolvedValueOnce({ request_id: 'req_1', data: { items: [], total: 0 } })
      .mockResolvedValueOnce({ request_id: 'req_2', data: { items: [sopProposal], total: 1 } })
    vi.mocked(createSOPAddendumProposal).mockResolvedValue({ request_id: 'req_sop', data: { proposal_id: 'prop_sop', status: 'pending_user_confirm', notification_id: 'notif_sop', audit_event_ids: ['audit_sop'], safety_statement: 'SOP 补充只生成待确认提案，不自动应用规则、不连接券商、不自动交易。' } })

    render(<RulesPage />)

    const button = await screen.findByRole('button', { name: '生成 SOP 补充提案' })
    fireEvent.click(button)

    await waitFor(() => expect(createSOPAddendumProposal).toHaveBeenCalledWith(expect.objectContaining({ scenario_key: 'p88_uncovered_liquidity_gap', occurrence_count: 4, sample_window: '2026-Q2' })))
    await waitFor(() => expect(screen.getByText('SOP 补充提案：连续流动性缺口未覆盖')).toBeInTheDocument())
    expect(screen.getByText('SOP 补充提案已生成，等待人工确认。')).toBeInTheDocument()
  })

  it('renders rule proposal contract fields', async () => {
    vi.mocked(listRuleProposals).mockResolvedValue({ request_id: 'req_1', data: { items: [proposal], total: 1 } })

    render(<RulesPage />)

    await waitFor(() => expect(screen.getByText('规则治理需要人工复核')).toBeInTheDocument())
    expect(screen.getByText('复核最终确认')).toBeInTheDocument()
    expect(screen.getByText('检查守门人结果')).toBeInTheDocument()
    expect(screen.getByText('规则页只处理本地规则治理和人工确认边界，不连接券商、不触发交易、不绕过守门人。')).toBeInTheDocument()
    expect(screen.getByText('来源误判案例：err_1')).toBeInTheDocument()
    expect(screen.getByText('提案理由：样本表现稳定')).toBeInTheDocument()
    expect(screen.getByText('守门人结果：审计通过')).toBeInTheDocument()
    expect(screen.getByText('审计摘要：通过')).toBeInTheDocument()
    expect(JSON.parse(screen.getByLabelText('影响范围').textContent ?? '{}')).toEqual({ scope: 'portfolio' })
    expect(JSON.parse(screen.getByLabelText('风险提示').textContent ?? '{}')).toEqual({ risk: 'low' })
    expect(screen.getByLabelText('变更前规则')).toHaveTextContent('旧规则文本')
    expect(screen.getByLabelText('变更后规则')).toHaveTextContent('新规则文本')
    expect(screen.getByText('规则效果验证')).toBeInTheDocument()
    expect(screen.getByText('验证状态：已通过')).toBeInTheDocument()
    expect(screen.getByText('过拟合风险：低')).toBeInTheDocument()
    expect(screen.getByText('历史回放：通过')).toBeInTheDocument()
    expect(screen.getByText('门禁结论：通过')).toBeInTheDocument()
    expect(JSON.parse(screen.getByLabelText('验证来源').textContent ?? '{}')).toEqual({ source_case_count: 3, related_error_case_ids: ['err_1'] })
    expect(JSON.parse(screen.getByLabelText('验证指标').textContent ?? '{}')).toEqual({ hit_count: 4, misjudgment_count: 0 })
    expect(JSON.parse(screen.getByLabelText('验证风险提示').textContent ?? '[]')).toContain('本地样本暂未发现不利信号')
    expect(JSON.parse(screen.getByLabelText('关联误判案例').textContent ?? '[]')).toEqual(['err_1'])
    expect(JSON.parse(screen.getByLabelText('关联决策记录').textContent ?? '[]')).toEqual(['dec_1'])
    expect(JSON.parse(screen.getByLabelText('关联风险预警').textContent ?? '[]')).toEqual(['risk_1'])
    expect(JSON.parse(screen.getByLabelText('关联审计事件').textContent ?? '[]')).toEqual(['audit_1'])
    expect(screen.getByText('规则效果验证只用于本地规则治理；规则生效仍需用户手动确认。')).toBeInTheDocument()
    expect(document.body.textContent).not.toMatch(/自动规则应用|自动应用规则/)
  })

  it('clears old confirmation failure message after refresh succeeds', async () => {
    vi.mocked(listRuleProposals)
      .mockResolvedValueOnce({ request_id: 'req_1', data: { items: [proposal], total: 1 } })
      .mockResolvedValueOnce({ request_id: 'req_2', data: { items: [proposal], total: 1 } })
    vi.mocked(finalConfirmRuleProposal)
      .mockRejectedValueOnce(new APIClientError({ requestId: 'rid', code: 'INVALID_STATE', message: '当前状态不允许执行该操作。', displayState: 'frozen_watch' }))
      .mockResolvedValueOnce({ request_id: 'req_confirm', data: { proposal_id: 'prop_1', status: 'applied' } })

    render(<RulesPage />)

    await waitFor(() => expect(screen.getByText('规则提案')).toBeInTheDocument())
    fireEvent.click(screen.getByRole('button', { name: '确认应用到正式规则' }))
    await waitFor(() => expect(screen.getByText('当前状态不允许执行该操作。')).toBeInTheDocument())

    fireEvent.click(screen.getByRole('button', { name: '确认应用到正式规则' }))
    await waitFor(() => expect(screen.queryByText('当前状态不允许执行该操作。')).not.toBeInTheDocument())
  })
})
