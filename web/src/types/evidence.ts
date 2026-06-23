export interface EvidenceItem {
  evidence_id: string
  source_name: string
  source_level: string
  evidence_role?: string
  verification_status?: string
  published_at?: string
  captured_at?: string
  original_url?: string
  summary: string
  content_hash?: string
  time_weight?: number
  relevance_score?: number
  high_grade_independent_source_count?: number
}

export interface EvidenceSummary {
  source_count: number
  highest_source_level: string
  verification_status: string
}

export interface EvidenceRefreshRequest {
  symbol?: string
  refresh_scope: string
  include_background: boolean
}

export interface EvidenceRefreshResponse {
  intelligence_item_count: number
  summary_count: number
  rag_chunk_count: number
  verification_count: number
  index_status: string
  failed_reason?: string
  audit_event_ids: string[]
}

export interface SourceVerification {
  verification_id: string
  verification_status: string
  independent_source_count: number
  high_grade_independent_source_count: number
  highest_source_level: string
  latest_published_at: string
  evidence_ids: string[]
}

export interface RebuildIndexResponse {
  indexed_count: number
  skipped_count: number
  audit_event_ids: string[]
}
