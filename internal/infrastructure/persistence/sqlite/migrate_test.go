package sqlite

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func testDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", t.TempDir()+"/test.db")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := Migrate(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestMigrateCreatesTablesAndSeeds(t *testing.T) {
	db := testDB(t)
	ctx := context.Background()
	if err := Migrate(ctx, db); err != nil {
		t.Fatalf("repeat migrate: %v", err)
	}
	for _, table := range []string{"portfolio_snapshots", "positions", "position_snapshots", "operation_confirmations", "position_transactions", "market_snapshots", "rule_versions", "decision_records", "intelligence_items", "intelligence_summary", "rag_chunks", "evidence_refs", "source_verifications", "capability_configs", "user_settings", "audit_events", "error_cases", "rule_proposals", "gatekeeper_audits", "notifications", "daily_auto_run_states", "daily_discipline_reports", "rule_effect_validations", "rule_effect_tracking", "data_quality_gate_resolutions"} {
		var name string
		if err := db.QueryRowContext(ctx, `SELECT name FROM sqlite_master WHERE type='table' AND name=?`, table).Scan(&name); err != nil {
			t.Fatalf("table %s missing: %v", table, err)
		}
	}
	var status string
	if err := db.QueryRowContext(ctx, `SELECT status FROM rule_versions WHERE rule_version='v3.0'`).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != "active" {
		t.Fatalf("rule v3.0 status=%s", status)
	}
}

func TestMigrateAllowsRunLocalTaskAuditOnExistingSchema(t *testing.T) {
	db, err := sql.Open("sqlite", t.TempDir()+"/legacy.db")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	_, err = db.Exec(`CREATE TABLE audit_events (
  audit_event_id TEXT PRIMARY KEY,
  request_id TEXT,
  decision_id TEXT,
  workflow_type TEXT,
  node_name TEXT,
  actor TEXT NOT NULL CHECK (actor IN ('system', 'user', 'gatekeeper')),
  action TEXT NOT NULL CHECK (action IN ('generate_decision', 'confirm_operation', 'mark_error', 'create_proposal', 'audit_rule_change', 'update_rule', 'refresh_market_data', 'update_settings', 'update_capability', 'rebuild_index')),
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
)`)
	if err != nil {
		t.Fatalf("create legacy audit table: %v", err)
	}
	if err := Migrate(context.Background(), db); err != nil {
		t.Fatalf("migrate legacy schema: %v", err)
	}
	_, err = db.Exec(`INSERT INTO audit_events (audit_event_id,actor,action,status,created_at) VALUES ('audit_local','user','run_local_task','success','2026-01-01T00:00:00Z')`)
	if err != nil {
		t.Fatalf("run_local_task should be allowed after migration: %v", err)
	}
}

func TestMigrateAddsEvidenceRefIndependentSourceCountOnExistingSchema(t *testing.T) {
	db, err := sql.Open("sqlite", t.TempDir()+"/legacy_evidence.db")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	_, err = db.Exec(`CREATE TABLE evidence_refs (
  evidence_ref_id TEXT PRIMARY KEY,
  evidence_id TEXT NOT NULL,
  decision_id TEXT NOT NULL,
  summary_id TEXT NOT NULL,
  source_name TEXT NOT NULL,
  source_level TEXT NOT NULL CHECK (source_level IN ('S', 'A', 'B', 'C')),
  evidence_role TEXT NOT NULL CHECK (evidence_role IN ('formal', 'background')),
  published_at DATETIME,
  captured_at DATETIME,
  original_url TEXT,
  summary TEXT NOT NULL,
  content_hash TEXT,
  time_weight REAL,
  relevance_score REAL,
  high_grade_independent_source_count INTEGER NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL,
  CHECK (source_level != 'C' OR evidence_role = 'background')
)`)
	if err != nil {
		t.Fatalf("create legacy evidence_refs: %v", err)
	}
	if err := Migrate(context.Background(), db); err != nil {
		t.Fatalf("migrate legacy evidence refs: %v", err)
	}
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('evidence_refs') WHERE name='independent_source_count'`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected independent_source_count column, got %d", count)
	}
}

func TestEnumCheckConstraint(t *testing.T) {
	db := testDB(t)
	_, err := db.Exec(`INSERT INTO positions (position_id,symbol,name,quantity,cost_price,current_price,market_value,unrealized_profit_ratio,position_state,updated_at) VALUES ('p','A','A',1,1,1,1,0,'invalid',CURRENT_TIMESTAMP)`)
	if err == nil {
		t.Fatal("expected CHECK constraint error")
	}
}
