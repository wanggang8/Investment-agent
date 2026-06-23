CREATE TABLE IF NOT EXISTS rule_effect_validations (
  validation_id TEXT PRIMARY KEY,
  proposal_id TEXT NOT NULL,
  candidate_rule_version TEXT NOT NULL,
  validation_status TEXT NOT NULL CHECK (validation_status IN ('not_evaluated','insufficient','passed','failed','needs_more_samples','needs_user_review')),
  sample_count INTEGER NOT NULL DEFAULT 0,
  sample_window TEXT NOT NULL,
  representativeness_status TEXT NOT NULL CHECK (representativeness_status IN ('not_evaluated','insufficient','passed','failed','needs_more_samples','needs_user_review')),
  overfit_risk TEXT NOT NULL CHECK (overfit_risk IN ('low','medium','high')),
  replay_result TEXT NOT NULL CHECK (replay_result IN ('passed','failed','mixed','unknown')),
  guardrail_decision TEXT NOT NULL CHECK (guardrail_decision IN ('passed','rejected','needs_user_review')),
  source_explanation_json TEXT,
  metrics_json TEXT,
  risk_notes_json TEXT,
  related_error_cases_json TEXT,
  related_decision_ids_json TEXT,
  related_risk_alert_ids_json TEXT,
  related_audit_event_ids_json TEXT,
  safety_note TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_rule_effect_validations_proposal ON rule_effect_validations(proposal_id, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_rule_effect_validations_rule_version ON rule_effect_validations(candidate_rule_version, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_rule_effect_validations_status ON rule_effect_validations(validation_status, updated_at DESC);

CREATE TABLE IF NOT EXISTS rule_effect_tracking (
  tracking_id TEXT PRIMARY KEY,
  applied_rule_version TEXT NOT NULL,
  proposal_id TEXT,
  period TEXT NOT NULL,
  hit_count INTEGER NOT NULL DEFAULT 0,
  misjudgment_count INTEGER NOT NULL DEFAULT 0,
  missing_evidence_count INTEGER NOT NULL DEFAULT 0,
  degraded_count INTEGER NOT NULL DEFAULT 0,
  risk_alert_count INTEGER NOT NULL DEFAULT 0,
  trend_direction TEXT NOT NULL CHECK (trend_direction IN ('improved','flat','worsened','unknown')),
  metrics_json TEXT,
  related_proposal_ids_json TEXT,
  related_audit_event_ids_json TEXT,
  related_risk_alert_ids_json TEXT,
  safety_note TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_rule_effect_tracking_rule_version ON rule_effect_tracking(applied_rule_version, period);
CREATE INDEX IF NOT EXISTS idx_rule_effect_tracking_proposal ON rule_effect_tracking(proposal_id, period);
CREATE INDEX IF NOT EXISTS idx_rule_effect_tracking_trend ON rule_effect_tracking(trend_direction, updated_at DESC);
