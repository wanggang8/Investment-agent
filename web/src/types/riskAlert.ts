export type RiskSeverity = 'info' | 'warning' | 'critical'
export type RiskSOPStatus = 'triggered' | 'active' | 'observing' | 'escalated' | 'resolved' | 'archived'
export type RiskType =
  | 'valuation_high'
  | 'buy_thesis_broken'
  | 'liquidity_danger'
  | 'sentiment_extreme'
  | 'position_limit_breach'
  | 'insufficient_evidence'
  | 'data_degraded'

export interface RiskAlert {
  alert_id: string
  risk_type: RiskType | string
  severity: RiskSeverity
  sop_status: RiskSOPStatus
  symbol: string
  trigger_summary: string
  trigger_context?: unknown
  prohibited_actions?: string[]
  suggested_actions?: string[]
  related_decision_id?: string
  related_report_id?: string
  related_notification_id?: string
  related_audit_event_id?: string
  decision_link?: string
  report_link?: string
  notification_link?: string
  audit_link?: string
  link?: string
  last_triggered_at?: string
  resolved_at?: string
  resolution_reason?: string
  safety_note: string
  created_at: string
  updated_at: string
}

export interface RiskAlertListFilter {
  statuses?: RiskSOPStatus[]
  symbol?: string
}

export interface RiskAlertLifecycleRequest {
  status: RiskSOPStatus
  reason: string
}
