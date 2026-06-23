import type { TriggeredRule } from './dashboard'
import type { EvidenceItem } from './evidence'

export interface ConsultDecisionRequest {
  question: string
  symbol: string
  scenario?: string
  expected_return_previous_base_midpoint?: number
  expected_return_target_return_rate?: number
}

export interface DecisionDetailResponse {
  decision_id: string
  question?: string
  symbol?: string
  generated_at?: string
  capability_check?: CapabilityCheck
  workflow_status: string
  account_snapshot?: AccountSnapshot
  triggered_rules: TriggeredRule[]
  evidence_chain: EvidenceItem[]
  analyst_reports: AnalystReport[]
  retrieval_quality?: RetrievalQualitySummary
  market_context?: MarketContext
  expected_return_scenarios?: ExpectedReturnScenarios
  arbitration_chain: ArbitrationStep[]
  audit_events?: DecisionAuditEvent[]
  final_verdict: FinalVerdict
  user_confirmation: UserConfirmation
}

export interface CapabilityCheck {
  status: string
  reason: string
}

export interface AccountSnapshot {
  snapshot_id: string
  cash?: number
  total_assets?: number
  cash_ratio: number
  high_risk_ratio: number
}

export interface AnalystReport {
  agent_name: string
  conclusion: string
  key_reasons: string[]
  risk_warnings: string[]
  confidence: string
  evidence_ids: string[]
  prompt_version?: string
  model?: string
  input_summary?: string
  output_summary?: string
  parse_status?: string
  quality_status?: string
}

export interface RetrievalQualitySummary {
  query_summary?: string
  top_k: number
  status?: string
  index_health?: string
  index_freshness?: string
  fallback_source?: string
  source_consistency_status?: string
  degraded_reason?: string
}

export interface MarketContext {
  symbol?: string
  trade_date?: string
  current_price?: number
  pe_percentile?: number
  pb_percentile?: number
}

export type PrecisionStatus = 'available' | 'insufficient' | 'unavailable'
export type FinalVerdictStatus = 'buy_allowed' | 'hold' | 'reduce' | 'sell_only' | 'frozen_watch' | 'rejected' | 'insufficient_data' | 'high_risk'

export interface ExpectedReturnScenarios {
  sample_count: number
  target_name?: string
  target_code?: string
  holding_class?: string
  horizon_label?: string
  sample_window?: string
  screening_condition?: string
  precision_status: PrecisionStatus
  probability_basis?: string
  scenarios: ReturnScenario[]
  reason?: string
  supporting_data_summary?: string
  missing_categories?: string[]
  supplement_data?: string[]
  assumption_checks?: AssumptionCheck[]
  historical_contexts?: HistoricalContext[]
  holding_class_coverage?: HoldingClassCoverage[]
  disclaimer: string
  sell_evaluation?: SellEvaluation
  reassessment_trigger?: ReassessmentTrigger
}

export interface AssumptionCheck {
  name: string
  expected: number
  actual: number
  months_below: number
}

export interface HoldingClassCoverage {
  holding_class: string
  symbol: string
  status: string
}

export interface HistoricalContext {
  label: string
  window: string
  sample_count: number
  outcome: string
  max_drawdown: number
  recovery: string
  source: string
}

export interface SellEvaluation {
  status: string
  triggers?: string[]
  prompts?: string[]
  actions?: string[]
  non_trading_disclaimer?: string
}

export interface ReassessmentTrigger {
  reason: string
  boundary?: string
  current_value?: number
}

export interface ReturnScenario {
  scenario: string
  return_range: string
  probability?: number | null
  trigger?: string
}

export interface ArbitrationStep {
  priority: number
  rule_id: string
  result: string
}

export interface DecisionAuditEvent {
  audit_event_id: string
  action: string
  status: string
  created_at?: string
  node_name?: string
  error_code?: string
}

export interface FinalVerdict {
  status: FinalVerdictStatus
  display_text: string
  prohibited_actions?: string[] | null
  optional_actions?: string[] | null
}

export interface UserConfirmation {
  confirmation_status: string
  available_actions: ConfirmationType[]
}

export type ConfirmationType = 'planned' | 'executed_manually' | 'watch' | 'marked_error'

export interface DecisionListItem {
  decision_id: string
  display_title: string
  symbol: string
  final_verdict: string
  triggered_rule_ids: string[]
  confirmation_status: string
  generated_at: string
}

export interface ConfirmationRequest {
  confirmation_type: ConfirmationType
  operation_type?: 'buy' | 'sell' | 'reduce'
  symbol?: string
  quantity?: number
  price?: number
  fees?: number
  executed_at?: string
  actual_outcome?: string
  root_cause_tag?: string
  lesson_learned?: string
  note?: string
}

export interface ConfirmationResponse {
  confirmation_id: string
  decision_id: string
  confirmation_status: string
  error_case_id?: string
  transaction_ids?: string[]
  snapshot_id?: string
  audit_event_ids: string[]
}
