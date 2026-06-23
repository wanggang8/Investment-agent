import type { RuleProposal, RuleVersion } from '../../types/rule'
import { ruleProposalStatusText } from '../../shared/mappers/statusText'
import type { OpsAction, OpsMetric, OpsTone } from './modelTypes'

export interface RuleProposalCard {
  proposalId: string
  title: string
  statusLabel: string
  reason: string
  sampleCount: number
  overfitRiskLabel: string
  guardrailLabel: string
  auditSummary: string
  safetyNote: string
}

export interface RulesGovernanceModel {
  overallLabel: string
  overallTone: OpsTone
  metrics: OpsMetric[]
  nextActions: OpsAction[]
  proposalCards: RuleProposalCard[]
  safetyNotes: string[]
}

export interface RulesGovernanceInput {
  currentRule?: RuleVersion
  proposals?: RuleProposal[]
}

const overfitRiskText: Record<string, string> = { low: '低', medium: '中', high: '高' }
const guardrailText: Record<string, string> = { passed: '通过', rejected: '拒绝', needs_user_review: '需要用户复核' }

export function buildRulesGovernanceModel(input: RulesGovernanceInput): RulesGovernanceModel {
  const proposals = input.proposals ?? []
  const pendingUser = proposals.filter((item) => item.status === 'pending_user_confirm').length
  const pendingFinal = proposals.filter((item) => item.status === 'pending_final_confirm').length
  const reviewNeeded = proposals.filter((item) => item.audit_result === 'needs_user_review' || item.effect_validation?.guardrail_decision === 'needs_user_review' || item.effect_validation?.validation_status === 'needs_user_review').length
  const rejectedOrFailed = proposals.filter((item) => item.audit_result === 'rejected' || item.effect_validation?.validation_status === 'failed' || item.effect_validation?.guardrail_decision === 'rejected').length
  const needsReview = pendingUser + pendingFinal + reviewNeeded + rejectedOrFailed > 0

  return {
    overallLabel: needsReview ? '规则治理需要人工复核' : proposals.length ? '规则治理暂无阻断' : '暂无规则提案',
    overallTone: rejectedOrFailed ? 'danger' : needsReview ? 'warning' : proposals.length ? 'success' : 'unknown',
    metrics: [
      { label: '当前规则', value: input.currentRule?.rule_version ?? '暂无', detail: input.currentRule?.status ?? 'unknown' },
      { label: '提案总数', value: String(proposals.length) },
      { label: '待用户确认', value: String(pendingUser), tone: pendingUser ? 'warning' : 'success' },
      { label: '待最终确认', value: String(pendingFinal), tone: pendingFinal ? 'warning' : 'success' },
      { label: '需复核门禁', value: String(reviewNeeded + rejectedOrFailed), tone: reviewNeeded + rejectedOrFailed ? 'warning' : 'success' },
    ],
    nextActions: buildRuleActions(pendingUser, pendingFinal, reviewNeeded + rejectedOrFailed),
    proposalCards: proposals.map(toProposalCard),
    safetyNotes: ['规则页只处理本地规则治理和人工确认边界，不连接券商、不触发交易、不绕过守门人。'],
  }
}

function buildRuleActions(pendingUser: number, pendingFinal: number, reviewNeeded: number): OpsAction[] {
  const actions: OpsAction[] = []
  if (pendingUser) actions.push({ label: '复核待确认提案', detail: '先确认是否送入守门人审计。', href: '/rules' })
  if (pendingFinal) actions.push({ label: '复核最终确认', detail: '守门人之后仍需人工最终确认。', href: '/rules' })
  if (reviewNeeded || pendingFinal) actions.push({ label: '检查守门人结果', detail: '查看验证、过拟合和审计摘要。', href: '/rules' })
  actions.push({ label: '查看审计记录', detail: '追踪规则治理相关审计事件。', href: '/audit' })
  return actions
}

function toProposalCard(proposal: RuleProposal): RuleProposalCard {
  return {
    proposalId: proposal.proposal_id,
    title: proposal.title,
    statusLabel: ruleProposalStatusText[proposal.status] ?? '未知状态',
    reason: proposal.reason ?? '暂无提案理由',
    sampleCount: proposal.effect_validation?.sample_count ?? proposal.sample_count,
    overfitRiskLabel: overfitRiskText[proposal.effect_validation?.overfit_risk ?? ''] ?? '未知',
    guardrailLabel: guardrailText[proposal.effect_validation?.guardrail_decision ?? ''] ?? '未知',
    auditSummary: proposal.audit_summary ?? '暂无审计摘要',
    safetyNote: proposal.effect_validation?.safety_note ?? '本地规则治理需人工复核，不会触发交易。',
  }
}
