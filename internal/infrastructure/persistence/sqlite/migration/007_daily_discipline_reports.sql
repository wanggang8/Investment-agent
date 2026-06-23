CREATE TABLE IF NOT EXISTS daily_discipline_reports (
  report_id TEXT PRIMARY KEY,
  local_date TEXT NOT NULL,
  scope TEXT NOT NULL CHECK (scope IN ('holdings')),
  symbol_set_hash TEXT NOT NULL,
  source_type TEXT NOT NULL CHECK (source_type IN ('auto_run','manual')),
  source_id TEXT,
  decision_id TEXT,
  status TEXT NOT NULL CHECK (status IN ('not_started','running','success','degraded','failed','insufficient_data')),
  summary TEXT,
  failure_code TEXT,
  failure_reason TEXT,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  UNIQUE(local_date, scope, symbol_set_hash)
);

CREATE INDEX IF NOT EXISTS idx_daily_discipline_reports_local_date ON daily_discipline_reports(local_date DESC);
CREATE INDEX IF NOT EXISTS idx_daily_discipline_reports_status_local_date ON daily_discipline_reports(status, local_date DESC);
