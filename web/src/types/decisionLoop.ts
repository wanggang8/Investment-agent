export type DecisionLoopStatus = 'open' | 'planned' | 'recorded' | 'reviewed' | 'incomplete'

export type DecisionLoopStageStatus = 'complete' | 'pending' | 'not_required' | 'missing' | 'degraded'

export interface DecisionLoopStage {
  stage: 'recommendation' | 'confirmation' | 'manual_record' | 'risk_review' | 'review' | string
  status: DecisionLoopStageStatus
  label: string
  summary: string
  ref_type?: string
  ref_id?: string
  at?: string
}

export interface DecisionLoopManualAction {
  confirmation_id: string
  confirmation_type: string
  operation_type?: string
  symbol?: string
  quantity?: number
  price?: number
  fees?: number
  executed_at?: string
  transaction_ids: string[]
  note_preview?: string
}

export interface DecisionLoopLink {
  type: string
  id: string
  label: string
  href: string
  status?: string
}

export interface DecisionLoopItem {
  decision_id: string
  symbol?: string
  generated_at: string
  final_verdict_status: string
  final_verdict_text: string
  confirmation_status: string
  loop_status: DecisionLoopStatus
  stages: DecisionLoopStage[]
  manual_actions: DecisionLoopManualAction[]
  risk_links: DecisionLoopLink[]
  review_links: DecisionLoopLink[]
  audit_links: DecisionLoopLink[]
  missing_links: string[]
  safety_note: string
}

export interface DecisionLoopListResponse {
  items: DecisionLoopItem[]
  total: number
  safety_note: string
}

export interface DecisionLoopListFilter {
  symbol?: string
  limit?: number
}
