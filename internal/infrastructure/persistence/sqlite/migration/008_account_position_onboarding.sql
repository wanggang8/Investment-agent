-- P33 local account onboarding support tables.
-- These tables preserve batch validation/correction traceability without replacing
-- portfolio_snapshots, positions, position_snapshots, position_transactions, or audit_events.

CREATE TABLE IF NOT EXISTS local_account_import_batches (
  import_batch_id TEXT PRIMARY KEY,
  request_id TEXT NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('validated', 'committed', 'rejected')),
  row_count INTEGER NOT NULL,
  valid_count INTEGER NOT NULL,
  invalid_count INTEGER NOT NULL,
  validation_summary_json TEXT,
  rows_hash TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  committed_at DATETIME
);
CREATE INDEX IF NOT EXISTS idx_local_account_import_batches_status ON local_account_import_batches(status);
CREATE INDEX IF NOT EXISTS idx_local_account_import_batches_created_at ON local_account_import_batches(created_at DESC);

CREATE TABLE IF NOT EXISTS local_account_corrections (
  correction_id TEXT PRIMARY KEY,
  target_type TEXT NOT NULL CHECK (target_type IN ('portfolio_snapshot', 'position', 'position_snapshot', 'position_transaction', 'import_batch')),
  target_id TEXT NOT NULL,
  before_json TEXT NOT NULL,
  after_json TEXT NOT NULL,
  correction_reason TEXT NOT NULL,
  snapshot_id TEXT,
  audit_event_id TEXT,
  created_at DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_local_account_corrections_target ON local_account_corrections(target_type, target_id);
CREATE INDEX IF NOT EXISTS idx_local_account_corrections_created_at ON local_account_corrections(created_at DESC);
