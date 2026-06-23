import { describe, expect, it } from 'vitest'
import { buildRulesGovernanceModel } from './rulesGovernanceModel'
import type { RuleProposal, RuleVersion } from '../../types/rule'

const currentRule: RuleVersion = {
  rule_version: 'v3.0',
  status: 'active',
  effective_at: '2026-01-01T00:00:00Z',
  created_at: '2026-01-01T00:00:00Z',
  rules: { priority: ['safety_gate'], thresholds: { cash_ratio_min: 0.1 } },
}

const proposal: RuleProposal = {
  proposal_id: 'prop_1',
  proposal_type: 'risk_rule',
  status: 'pending_final_confirm',
  title: '降低高风险买入阈值',
  proposal_version: 'draft',
  sample_count: 12,
  created_at: '2026-06-18T08:00:00Z',
  reason: '历史样本显示高估值区间误判增加',
  audit_result: 'approved',
  audit_summary: '守门人通过，但仍需人工最终确认。',
  effect_validation: {
    validation_status: 'needs_user_review',
    sample_count: 12,
    sample_window: '2026-Q2',
    overfit_risk: 'medium',
    replay_result: 'mixed',
    guardrail_decision: 'needs_user_review',
    risk_notes: ['样本窗口偏短'],
    related_audit_event_ids: ['audit_1'],
    safety_note: '规则效果验证只用于本地规则治理，不会自动应用规则或执行交易。',
  },
}

describe('buildRulesGovernanceModel', () => {
  it('summarizes rule governance status and manual actions', () => {
    const model = buildRulesGovernanceModel({
      currentRule,
      proposals: [proposal, { ...proposal, proposal_id: 'prop_2', status: 'pending_user_confirm', audit_result: 'needs_user_review' }],
    })

    expect(model.overallTone).toBe('warning')
    expect(model.overallLabel).toBe('规则治理需要人工复核')
    expect(model.metrics).toContainEqual(expect.objectContaining({ label: '当前规则', value: 'v3.0' }))
    expect(model.metrics).toContainEqual(expect.objectContaining({ label: '待用户确认', value: '1' }))
    expect(model.metrics).toContainEqual(expect.objectContaining({ label: '待最终确认', value: '1' }))
    expect(model.nextActions.map((action) => action.label)).toEqual(expect.arrayContaining(['复核待确认提案', '检查守门人结果', '查看审计记录']))
    expect(model.proposalCards[0]).toMatchObject({
      proposalId: 'prop_1',
      title: '降低高风险买入阈值',
      statusLabel: '待最终确认',
      overfitRiskLabel: '中',
      guardrailLabel: '需要用户复核',
    })
  })

  it('keeps governance copy inside safety boundaries', () => {
    const model = buildRulesGovernanceModel({ currentRule, proposals: [proposal] })
    const text = JSON.stringify(model)

    expect(text).toContain('本地规则治理')
    expect(text).not.toMatch(/自动交易|一键交易|代下单|外部推送|自动确认|自动规则应用|自动修复|收益承诺/)
  })
})

