package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"investment-agent/internal/application/service"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/infrastructure/config"
	appsqlite "investment-agent/internal/infrastructure/persistence/sqlite"
)

func TestRunSeedsP31AndP32WhenP30DecisionAlreadyExists(t *testing.T) {
	dbPath := configureSmokeSeedTest(t, "UTC")

	if err := run(); err != nil {
		t.Fatalf("initial run failed: %v", err)
	}
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer store.Close()
	ctx := context.Background()
	if _, err := store.DB.ExecContext(ctx, `DELETE FROM daily_discipline_reports`); err != nil {
		t.Fatalf("delete reports: %v", err)
	}
	if _, err := store.DB.ExecContext(ctx, `DELETE FROM daily_auto_run_states`); err != nil {
		t.Fatalf("delete auto-run states: %v", err)
	}
	if _, err := store.DB.ExecContext(ctx, `DELETE FROM notifications WHERE notification_id='notif_smoke_p31_auto_run'`); err != nil {
		t.Fatalf("delete notification: %v", err)
	}
	if _, err := store.DB.ExecContext(ctx, `DELETE FROM audit_events WHERE audit_event_id='audit_smoke_p31_auto_run'`); err != nil {
		t.Fatalf("delete audit event: %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}

	if err := run(); err != nil {
		t.Fatalf("rerun failed: %v", err)
	}

	store, err = appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("reopen db: %v", err)
	}
	defer store.Close()
	var autoRunCount, reportCount int
	if err := store.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM daily_auto_run_states`).Scan(&autoRunCount); err != nil {
		t.Fatalf("count auto-run states: %v", err)
	}
	if err := store.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM daily_discipline_reports`).Scan(&reportCount); err != nil {
		t.Fatalf("count reports: %v", err)
	}
	if autoRunCount != 1 || reportCount != 1 {
		t.Fatalf("expected P31/P32 rows to be backfilled, got auto-run=%d reports=%d", autoRunCount, reportCount)
	}
}

func TestRunUsesConfiguredTimezoneLocalDateForP31AndP32(t *testing.T) {
	dbPath := configureSmokeSeedTest(t, "Pacific/Kiritimati")
	loc, err := time.LoadLocation("Pacific/Kiritimati")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}
	expectedLocalDate := time.Now().In(loc).Format(time.DateOnly)
	expectedSourceIDPrefix := expectedLocalDate + ":holdings:"

	if err := run(); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if err := run(); err != nil {
		t.Fatalf("idempotent rerun failed: %v", err)
	}

	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer store.Close()
	ctx := context.Background()
	var reportLocalDate, reportSourceID, stateLocalDate, stateKey string
	if err := store.DB.QueryRowContext(ctx, `SELECT local_date, source_id FROM daily_discipline_reports WHERE local_date=? AND scope='holdings' AND symbol_set_hash='p32smokereport'`, expectedLocalDate).Scan(&reportLocalDate, &reportSourceID); err != nil {
		t.Fatalf("read report: %v", err)
	}
	if err := store.DB.QueryRowContext(ctx, `SELECT local_date, idempotency_key FROM daily_auto_run_states WHERE run_id='auto_run_smoke_p31'`).Scan(&stateLocalDate, &stateKey); err != nil {
		t.Fatalf("read auto-run state: %v", err)
	}
	if reportLocalDate != expectedLocalDate || stateLocalDate != expectedLocalDate {
		t.Fatalf("expected local date %s, got report=%s state=%s", expectedLocalDate, reportLocalDate, stateLocalDate)
	}
	if !strings.HasPrefix(stateKey, expectedSourceIDPrefix) || reportSourceID != "auto_run_smoke_p32_success" || stateKey == reportSourceID {
		t.Fatalf("expected state key to use %q prefix and P32 report to use separated source id, got source_id=%q state_key=%q", expectedSourceIDPrefix, reportSourceID, stateKey)
	}
	var notificationCount, auditCount int
	if err := store.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM notifications WHERE source_id=?`, stateKey).Scan(&notificationCount); err != nil {
		t.Fatalf("count notifications: %v", err)
	}
	if err := store.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM audit_events WHERE audit_event_id='audit_smoke_p31_auto_run'`).Scan(&auditCount); err != nil {
		t.Fatalf("count audit events: %v", err)
	}
	if notificationCount != 1 || auditCount != 1 {
		t.Fatalf("expected idempotent notification/audit rows, got notifications=%d audit=%d", notificationCount, auditCount)
	}
}

func TestRunSeedsP31FailedStateAndP32SuccessReportWithDifferentKeys(t *testing.T) {
	dbPath := configureSmokeSeedTest(t, "UTC")
	expectedLocalDate := time.Now().UTC().Format(time.DateOnly)

	if err := run(); err != nil {
		t.Fatalf("run failed: %v", err)
	}

	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer store.Close()
	ctx := context.Background()
	var stateKey, stateHash, reportSourceID, reportHash string
	if err := store.DB.QueryRowContext(ctx, `SELECT idempotency_key,symbol_set_hash FROM daily_auto_run_states WHERE run_id='auto_run_smoke_p31'`).Scan(&stateKey, &stateHash); err != nil {
		t.Fatalf("read auto-run state: %v", err)
	}
	if err := store.DB.QueryRowContext(ctx, `SELECT source_id,symbol_set_hash FROM daily_discipline_reports WHERE report_id='daily_report_smoke_p32'`).Scan(&reportSourceID, &reportHash); err != nil {
		t.Fatalf("read P32 report: %v", err)
	}
	if stateKey == reportSourceID || stateHash == reportHash {
		t.Fatalf("expected P31 failed state and P32 success report to use different keys, state=%q/%q report=%q/%q", stateKey, stateHash, reportSourceID, reportHash)
	}

	svc := service.NewQueryServiceWithDailyAutoRunConfig(repository.Repositories{
		DailyAutoRunRepo:          appsqlite.NewDailyAutoRunRepository(store.DB),
		DailyDisciplineReportRepo: appsqlite.NewDailyDisciplineReportRepository(store.DB),
		DecisionRepo:              appsqlite.NewDecisionRepository(store.DB),
	}, config.DailyAutoRunConfig{Timezone: "UTC"})
	out, err := svc.TodayDailyDisciplineReport(ctx, time.Now().UTC())
	if err != nil {
		t.Fatalf("TodayDailyDisciplineReport: %v", err)
	}
	if out.ReportID != "daily_report_smoke_p32" || out.Status != "success" || out.SourceID != reportSourceID || out.LocalDate != expectedLocalDate {
		t.Fatalf("expected P32 success report not shadowed by failed auto-run state, got %+v", out)
	}
}

func configureSmokeSeedTest(t *testing.T, timezone string) string {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "smoke.db")
	vecPath := filepath.Join(dir, "veclite")
	configPath := filepath.Join(dir, "config.yaml")
	configYAML := fmt.Sprintf(`server:
  host: "127.0.0.1"
  port: 0
sqlite:
  path: %q
veclite:
  path: %q
deepseek:
  api_key: ""
  base_url: "https://api.deepseek.com"
data_sources:
  enabled:
    - "stub"
  use_stub: true
  market_endpoint: ""
  intelligence_endpoint: ""
  public_evidence:
    enabled: false
    sources:
      - "cninfo"
    cninfo_base_url: "https://www.cninfo.com.cn"
    cninfo_org_ids:
      "510300": "9900000091"
    szse_base_url: "https://www.szse.cn"
    csrc_base_url: "https://www.csrc.gov.cn"
  market_collectors:
    enabled: false
    sources:
      - "csindex"
    csindex_base_url: "https://www.csindex.com.cn"
    eastmoney_fund_base_url: "https://fund.eastmoney.com"
daily_auto_run:
  enabled: false
  run_time: "08:30"
  timezone: %q
  scope: "holdings"
  retry: 1
  timeout_seconds: 900
  max_symbols: 20
log:
  level: "error"
`, dbPath, vecPath, timezone)
	if err := os.WriteFile(configPath, []byte(configYAML), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	t.Setenv("INVESTMENT_AGENT_CONFIG", configPath)
	return dbPath
}
