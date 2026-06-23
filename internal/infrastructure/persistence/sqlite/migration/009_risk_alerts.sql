CREATE TABLE IF NOT EXISTS risk_alerts (
  alert_id TEXT PRIMARY KEY,
  risk_type TEXT NOT NULL CHECK (risk_type IN ('valuation_high','buy_thesis_broken','liquidity_danger','sentiment_extreme','position_limit_breach','insufficient_evidence','data_degraded')),
  severity TEXT NOT NULL CHECK (severity IN ('info','warning','critical')),
  sop_status TEXT NOT NULL CHECK (sop_status IN ('triggered','active','observing','escalated','resolved','archived')),
  symbol TEXT NOT NULL,
  trigger_summary TEXT NOT NULL,
  trigger_context_json TEXT,
  prohibited_actions_json TEXT,
  suggested_actions_json TEXT,
  related_decision_id TEXT,
  related_report_id TEXT,
  related_notification_id TEXT,
  related_audit_event_id TEXT,
  last_triggered_at DATETIME,
  resolved_at DATETIME,
  resolution_reason TEXT,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_risk_alerts_active_identity ON risk_alerts(risk_type, symbol) WHERE sop_status IN ('triggered','active','observing','escalated');
CREATE INDEX IF NOT EXISTS idx_risk_alerts_status_updated_at ON risk_alerts(sop_status, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_risk_alerts_symbol_status ON risk_alerts(symbol, sop_status);
