export interface AuditEvent {
  audit_event_id?: string
  event_id?: string
  request_id: string
  decision_id?: string
  workflow_type?: string
  node_name?: string
  actor: string
  action: string
  node_action?: string
  proposal_id?: string
  confirmation_id?: string
  error_case_id?: string
  status: string
  error_code?: string
  before_state?: string
  after_state?: string
  rule_version?: string
  snapshot_id?: string
  input_ref_type?: string
  input_ref?: string
  output_ref_type?: string
  output_ref?: string
  created_at: string
}
