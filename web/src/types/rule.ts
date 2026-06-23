export interface RuleVersion {
  rule_version: string
  status: string
  rules?: unknown
  effective_at: string
  created_at: string
}

export type RuleProposalStatus = 'draft' | 'pending_user_confirm' | 'under_gatekeeper_audit' | 'pending_final_confirm' | 'rejected' | 'applied'

export interface RuleEffectValidation {
  validation_id?: string
  proposal_id?: string
  candidate_rule_version?: string
  validation_status: 'not_evaluated' | 'insufficient' | 'passed' | 'failed' | 'needs_more_samples' | 'needs_user_review' | string
  sample_count: number
  sample_window?: string
  representativeness_status?: string
  overfit_risk?: 'low' | 'medium' | 'high' | string
  replay_result?: 'passed' | 'failed' | 'mixed' | 'unknown' | string
  guardrail_decision?: 'passed' | 'rejected' | 'needs_user_review' | string
  source_explanation?: unknown
  metrics?: unknown
  risk_notes?: unknown
  related_error_cases?: unknown
  related_decision_ids?: unknown
  related_risk_alert_ids?: unknown
  related_audit_event_ids?: unknown
  validation_link?: string
  safety_note?: string
  created_at?: string
  updated_at?: string
}

export interface RuleEffectTracking {
  tracking_id: string
  applied_rule_version: string
  proposal_id?: string
  period: string
  hit_count: number
  misjudgment_count: number
  missing_evidence_count: number
  degraded_count: number
  risk_alert_count: number
  trend_direction: 'improved' | 'flat' | 'worsened' | 'unknown' | string
  metrics?: unknown
  related_proposal_ids?: unknown
  related_audit_event_ids?: unknown
  related_risk_alert_ids?: unknown
  safety_note?: string
  created_at?: string
  updated_at?: string
}

export interface RuleProposal {
  proposal_id: string
  proposal_type: string
  status: RuleProposalStatus
  title: string
  proposal_version: string
  reason?: string
  source_error_case_id?: string
  before_rule?: unknown
  after_rule?: unknown
  impact_scope?: unknown
  risk_notes?: unknown
  audit_result?: string
  audit_summary?: string
  effect_validation?: RuleEffectValidation
  sample_count: number
  created_at: string
}

export interface RuleProposalConfirmRequest {
  confirm: boolean
  note?: string
}

export interface SOPAddendumProposalRequest {
  scenario_key: string
  scenario_title: string
  occurrence_count: number
  sample_window: string
}

export interface SOPAddendumProposalResponse {
  proposal_id: string
  status: string
  notification_id: string
  audit_event_ids: string[]
  safety_statement: string
}

export interface RuleEffectValidationRefreshRequest {
  sample_window?: string
}

export interface RuleProposalConfirmResponse {
  proposal_id: string
  status: string
  gatekeeper_audit_id?: string
  applied_rule_version?: string
  created_rule_version?: string
  final_confirmed_at?: string
  audit_events?: string[]
  audit_event_ids?: string[]
}
