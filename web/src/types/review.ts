import type { DecisionListItem } from './decision'
import type { RuleEffectTracking } from './rule'

export interface ReviewSummary {
  period?: 'monthly' | 'quarterly'
  decision_count?: number
  confirmation_count?: number
  executed_manually_count?: number
  planned_count?: number
  error_case_count: number
  rule_proposal_count: number
  audit_event_count?: number
  rule_hit_count?: number
  misjudgment_count?: number
  missing_evidence_count?: number
  degraded_count?: number
  ops_status?: ReviewOpsStatus
  rule_suggestions?: RuleSuggestion[]
  attribution_summaries?: ReviewAttribution[]
  recurring_error_tags?: ReviewErrorTag[]
  missing_evidence_themes?: ReviewEvidenceTheme[]
  rule_proposal_outcomes?: ReviewRuleProposalOutcome[]
  rule_effect_tracking?: RuleEffectTracking[]
  degraded_workflows?: ReviewDegradedWorkflow[]
  tracking_links?: ReviewTrackingLink[]
  recent_decisions?: DecisionListItem[]
}

export interface ReviewOpsStatus {
  data_source_status?: string
  index_status?: string
  review_status?: string
  explanation?: string
}

export interface RuleSuggestion {
  proposal_id: string
  title: string
  status: string
  reason?: string
  can_auto_apply: boolean
}

export interface ReviewAttribution {
  decision_id: string
  symbol?: string
  verdict?: string
  confirmation_status?: string
  evidence_status?: string
  workflow_status?: string
  outcome: string
}

export interface ReviewErrorTag {
  tag: string
  count: number
}

export interface ReviewEvidenceTheme {
  status: string
  count: number
}

export interface ReviewRuleProposalOutcome {
  proposal_id: string
  title: string
  status: string
  audit_result?: string
}

export interface ReviewDegradedWorkflow {
  decision_id: string
  symbol?: string
  status: string
  created_at: string
}
export interface ReviewTrackingLink {
  type: 'audit_event' | 'rule_proposal' | 'error_case' | string
  id: string
  label: string
}
