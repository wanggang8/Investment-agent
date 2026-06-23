CREATE TABLE IF NOT EXISTS data_quality_gate_resolutions (
  resolution_id TEXT PRIMARY KEY,
  symbol TEXT NOT NULL,
  policy_fingerprint TEXT NOT NULL,
  policy_verdict TEXT NOT NULL,
  release_gate TEXT NOT NULL,
  policy_summary TEXT NOT NULL,
  resolution_type TEXT NOT NULL,
  status TEXT NOT NULL,
  scope TEXT NOT NULL,
  reason TEXT NOT NULL,
  release_impact TEXT NOT NULL,
  evidence_ref TEXT,
  blocking_reasons_json TEXT,
  waiver_reasons_json TEXT,
  created_by TEXT NOT NULL,
  retired_by TEXT,
  created_at DATETIME NOT NULL,
  retired_at DATETIME,
  safety_note TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_data_quality_gate_resolutions_active_policy
ON data_quality_gate_resolutions(symbol, policy_fingerprint)
WHERE status = 'active';

CREATE INDEX IF NOT EXISTS idx_data_quality_gate_resolutions_symbol_created
ON data_quality_gate_resolutions(symbol, created_at DESC);
