-- P9 兼容迁移：为既有 SQLite 库补充本地任务审计动作。
-- SQLite 无法直接修改 CHECK 约束，这里重建 audit_events，保留历史数据。
PRAGMA foreign_keys=off;

CREATE TABLE IF NOT EXISTS audit_events_p9 (
  audit_event_id TEXT PRIMARY KEY,
  request_id TEXT,
  decision_id TEXT,
  workflow_type TEXT,
  node_name TEXT,
  actor TEXT NOT NULL CHECK (actor IN ('system', 'user', 'gatekeeper')),
  action TEXT NOT NULL CHECK (action IN ('generate_decision', 'confirm_operation', 'mark_error', 'create_proposal', 'audit_rule_change', 'update_rule', 'refresh_market_data', 'update_settings', 'update_capability', 'rebuild_index', 'run_local_task', 'risk_alert')),
  node_action TEXT,
  proposal_id TEXT,
  confirmation_id TEXT,
  error_case_id TEXT,
  status TEXT NOT NULL CHECK (status IN ('success', 'degraded', 'failed')),
  error_code TEXT,
  before_state TEXT,
  after_state TEXT,
  rule_version TEXT,
  snapshot_id TEXT,
  input_ref_type TEXT,
  input_ref TEXT,
  output_ref_type TEXT,
  output_ref TEXT,
  created_at DATETIME NOT NULL,
  CHECK (status != 'failed' OR error_code IS NOT NULL)
);

INSERT OR IGNORE INTO audit_events_p9 (
  audit_event_id,request_id,decision_id,workflow_type,node_name,actor,action,node_action,proposal_id,confirmation_id,error_case_id,status,error_code,before_state,after_state,rule_version,snapshot_id,input_ref_type,input_ref,output_ref_type,output_ref,created_at
)
SELECT audit_event_id,request_id,decision_id,workflow_type,node_name,actor,action,node_action,proposal_id,confirmation_id,error_case_id,status,error_code,before_state,after_state,rule_version,snapshot_id,input_ref_type,input_ref,output_ref_type,output_ref,created_at
FROM audit_events;

DROP TABLE IF EXISTS audit_events;
ALTER TABLE audit_events_p9 RENAME TO audit_events;

CREATE INDEX IF NOT EXISTS idx_audit_events_request ON audit_events(request_id);
CREATE INDEX IF NOT EXISTS idx_audit_events_decision ON audit_events(decision_id);
CREATE INDEX IF NOT EXISTS idx_audit_events_created ON audit_events(created_at);

PRAGMA foreign_keys=on;
