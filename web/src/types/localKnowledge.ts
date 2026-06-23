export interface LocalKnowledgeImportRow {
  title?: string
  text: string
  symbol?: string
  as_of_date?: string
  tags?: string[]
}

export interface LocalKnowledgeImportValidationRequest {
  source_label: string
  default_symbol?: string
  rows: LocalKnowledgeImportRow[]
}

export interface LocalKnowledgeImportConfirmRequest extends LocalKnowledgeImportValidationRequest {
  import_batch_id: string
  confirm_reason: string
}

export interface LocalKnowledgeImportRisk {
  code: string
  severity: 'warning' | 'blocking'
  message: string
}

export interface LocalKnowledgeImportRowResult {
  row_number: number
  status: 'valid' | 'warning' | 'blocking'
  symbol: string
  title_preview: string
  text_preview: string
  content_hash: string
  risks: LocalKnowledgeImportRisk[]
}

export interface LocalKnowledgeImportIndexPlan {
  rag_chunk_count: number
  index_status: string
}

export interface LocalKnowledgeImportValidationSummary {
  total_count: number
  valid_count: number
  warning_count: number
  blocking_count: number
}

export interface LocalKnowledgeImportValidationResponse {
  import_batch_id: string
  summary: LocalKnowledgeImportValidationSummary
  rows: LocalKnowledgeImportRowResult[]
  index_plan: LocalKnowledgeImportIndexPlan
  safety_note: string
}

export interface LocalKnowledgeImportConfirmResponse {
  import_batch_id: string
  intelligence_item_count: number
  summary_count: number
  rag_chunk_count: number
  verification_count: number
  audit_event_ids: string[]
  index_status: string
  safety_note: string
}
