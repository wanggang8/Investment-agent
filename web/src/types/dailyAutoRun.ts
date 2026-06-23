export interface DailyAutoRunStatus {
  enabled: boolean
  run_time?: string
  timezone?: string
  scope?: string
  status: 'disabled' | 'scheduled' | 'running' | 'success' | 'degraded' | 'failed'
  run_id?: string
  idempotency_key?: string
  local_date?: string
  last_run_at?: string
  next_run_at?: string
  failure_code?: string
  failure_reason?: string
  latest_decision_link?: string
  latest_notification_link?: string
  latest_audit_link?: string
  missing_action?: string
  safety_note: string
}
