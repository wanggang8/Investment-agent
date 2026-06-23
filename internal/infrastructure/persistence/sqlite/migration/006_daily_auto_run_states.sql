CREATE TABLE IF NOT EXISTS daily_auto_run_states (
  run_id TEXT PRIMARY KEY,
  idempotency_key TEXT NOT NULL UNIQUE,
  local_date TEXT NOT NULL,
  scope TEXT NOT NULL CHECK (scope IN ('holdings')),
  symbol_set_hash TEXT NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('disabled', 'scheduled', 'running', 'success', 'degraded', 'failed')),
  last_run_at DATETIME,
  next_run_at DATETIME,
  failure_code TEXT,
  failure_reason TEXT,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  CHECK (status != 'failed' OR failure_code IS NOT NULL)
);

CREATE INDEX IF NOT EXISTS idx_daily_auto_run_states_date_scope ON daily_auto_run_states(local_date, scope);
