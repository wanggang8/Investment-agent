export interface DataSourceQualityRegression {
  mode: string
  status: string
  generated_at: string
  summary: string
  cases: DataSourceQualityCase[]
  missing_categories: string[]
  policy: DataSourceQualityPolicy
  safety_note: string
}

export interface DataSourceQualityCase {
  case_id: string
  source_name: string
  source_level: string
  source_type: string
  data_category: string
  expected_freshness: string
  actual_freshness: string
  status: string
  data_date?: string
  failure_category?: string
  affected_symbols?: string[]
  diagnostic_preview?: string
}

export interface DataSourceQualityPolicy {
  verdict: string
  release_gate: string
  degraded_count: number
  failed_count: number
  blocking_count: number
  waiver_count: number
  blocking_reasons: string[] | null
  waiver_reasons: string[] | null
  next_actions: string[] | null
  safety_note: string
}

export interface DataQualityGateResolutionRecord {
  resolution_id: string
  symbol: string
  policy_fingerprint: string
  policy_verdict: string
  release_gate: string
  policy_summary: string
  resolution_type: string
  status: string
  scope: string
  reason: string
  release_impact: string
  evidence_ref?: string
  created_by: string
  retired_by?: string
  created_at: string
  retired_at?: string
  safety_note: string
}

export interface DataQualityGateResolutionCheck {
  symbol: string
  policy_fingerprint: string
  policy_summary: string
  policy: DataSourceQualityPolicy
  release_claim_state: string
  clean_data_claim_allowed: boolean
  active_resolution?: DataQualityGateResolutionRecord
  allowed_claims: string[] | null
  prohibited_claims: string[] | null
  safety_note: string
}

export interface DataQualityGateResolutionCreateRequest {
  symbol: string
  resolution_type: string
  scope: string
  reason: string
  release_impact: string
  evidence_ref?: string
}

export interface DataQualityGateResolutionListResponse {
  items: DataQualityGateResolutionRecord[]
  total: number
}
