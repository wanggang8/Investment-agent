package workflow

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/infrastructure/config"
)

func TestDailyAutoRunDisabledDoesNotWriteState(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	runner := NewDailyAutoRunner(config.DailyAutoRunConfig{Enabled: false}, deps)

	out, err := runner.RunOnce(context.Background(), time.Date(2026, 6, 7, 8, 30, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	if out.Status != "disabled" || out.IdempotencyKey != "" {
		t.Fatalf("unexpected disabled output: %+v", out)
	}
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM daily_auto_run_states`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("disabled auto-run must not write state, got %d", count)
	}
}

func TestDailyAutoRunWithHoldingsRunsRefreshAndDailyDiscipline(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	ctx := context.Background()
	if err := deps.PortfolioRepo.SavePortfolioSnapshot(ctx, repository.PortfolioSnapshot{SnapshotID: "portfolio_auto", SnapshotTime: "2026-06-07T00:00:00Z", Cash: 1000, TotalAssets: 10000, CashRatio: 0.1, PositionCount: 1, Source: "manual", CreatedAt: "2026-06-07T00:00:00Z"}, []repository.PositionSnapshot{{PositionSnapshotID: "ps_auto", SnapshotID: "portfolio_auto", Symbol: "510300", Name: "沪深300ETF", Quantity: 100, CostPrice: 4, CurrentPrice: 4.2, MarketValue: 420, UnrealizedProfitRatio: 0.05, PositionState: "normal", CreatedAt: "2026-06-07T00:00:00Z"}}); err != nil {
		t.Fatal(err)
	}
	if err := deps.PortfolioRepo.SavePosition(ctx, repository.Position{PositionID: "pos_auto", Symbol: "510300", Name: "沪深300ETF", Quantity: 100, CostPrice: 4, CurrentPrice: 4.2, MarketValue: 420, UnrealizedProfitRatio: 0.05, PositionState: "normal", UpdatedAt: "2026-06-07T00:00:00Z"}); err != nil {
		t.Fatal(err)
	}
	deps.MarketDataSource = testMarketDataSource{point: MarketDataPoint{ClosePrice: 4.21, SourceName: "stub", SourceLevel: model.SourceLevelB, SourceType: "daily_auto_run_test", TradeDate: "2026-06-07", PEPercentile: 50, PBPercentile: 45}}
	runner := NewDailyAutoRunner(config.DailyAutoRunConfig{Enabled: true, RunTime: "08:30", Timezone: "Asia/Shanghai", Scope: "holdings", Retry: 1, TimeoutSeconds: 900, MaxSymbols: 20}, deps)

	out, err := runner.RunOnce(ctx, time.Date(2026, 6, 7, 0, 30, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	if out.Status != "degraded" || out.FailureCode != "" {
		t.Fatalf("expected degraded auto-run for degraded daily discipline workflow, got %+v", out)
	}
	var stateStatus string
	if err := db.QueryRow(`SELECT status FROM daily_auto_run_states WHERE idempotency_key=?`, out.IdempotencyKey).Scan(&stateStatus); err != nil {
		t.Fatal(err)
	}
	if stateStatus != "degraded" {
		t.Fatalf("expected degraded state for degraded daily discipline workflow, got %s", stateStatus)
	}
	var decisionCount, marketCount, summaryCount, notificationCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM decision_records WHERE workflow_type='daily_discipline'`).Scan(&decisionCount); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM market_snapshots WHERE symbol='510300'`).Scan(&marketCount); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM intelligence_summary WHERE symbol='510300'`).Scan(&summaryCount); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM notifications WHERE type='daily_auto_run_success' AND source_id=?`, out.IdempotencyKey).Scan(&notificationCount); err != nil {
		t.Fatal(err)
	}
	if decisionCount != 1 || marketCount != 1 || summaryCount == 0 || notificationCount != 1 {
		t.Fatalf("expected decision, market snapshot, evidence summary and notification, got decision=%d market=%d summary=%d notification=%d", decisionCount, marketCount, summaryCount, notificationCount)
	}

	var report repository.DailyDisciplineReport
	if err := db.QueryRow(`SELECT report_id,local_date,scope,symbol_set_hash,source_type,source_id,COALESCE(decision_id,''),status,summary,COALESCE(failure_code,''),COALESCE(failure_reason,''),created_at,updated_at FROM daily_discipline_reports WHERE source_type='auto_run' AND source_id=?`, out.IdempotencyKey).Scan(&report.ReportID, &report.LocalDate, &report.Scope, &report.SymbolSetHash, &report.SourceType, &report.SourceID, &report.DecisionID, &report.Status, &report.Summary, &report.FailureCode, &report.FailureReason, &report.CreatedAt, &report.UpdatedAt); err != nil {
		t.Fatal(err)
	}
	if report.Status != "degraded" || report.SourceID != out.IdempotencyKey || report.DecisionID == "" || report.Summary != "今日纪律报告已生成" {
		t.Fatalf("unexpected degraded report row for degraded daily discipline workflow: %+v", report)
	}
	listed, err := deps.DailyDisciplineReportRepo.ListDailyDisciplineReports(ctx, repository.DailyDisciplineReportListFilter{Limit: 10})
	if err != nil {
		t.Fatalf("ListDailyDisciplineReports: %v", err)
	}
	if len(listed) != 1 || listed[0].SourceID != out.IdempotencyKey || listed[0].Status != "degraded" {
		t.Fatalf("expected degraded report visible in list, got %+v", listed)
	}
}

func TestDailyAutoRunMissingHoldingsIsIdempotent(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	runner := NewDailyAutoRunner(config.DailyAutoRunConfig{Enabled: true, RunTime: "08:30", Timezone: "Asia/Shanghai", Scope: "holdings", Retry: 1, TimeoutSeconds: 900, MaxSymbols: 20}, deps)
	runAt := time.Date(2026, 6, 7, 0, 30, 0, 0, time.UTC)

	first, err := runner.RunOnce(context.Background(), runAt)
	if err != nil {
		t.Fatalf("first RunOnce: %v", err)
	}
	second, err := runner.RunOnce(context.Background(), runAt)
	if err != nil {
		t.Fatalf("second RunOnce: %v", err)
	}
	if second.IdempotencyKey != first.IdempotencyKey {
		t.Fatalf("expected same idempotency key, got first=%s second=%s", first.IdempotencyKey, second.IdempotencyKey)
	}
	var stateCount, notificationCount, originalAuditCount, reuseAuditCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM daily_auto_run_states WHERE idempotency_key=?`, first.IdempotencyKey).Scan(&stateCount); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM notifications WHERE source_id=?`, first.IdempotencyKey).Scan(&notificationCount); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM audit_events WHERE workflow_type='daily_auto_run' AND input_ref=? AND node_action='daily_auto_run'`, first.IdempotencyKey).Scan(&originalAuditCount); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM audit_events WHERE workflow_type='daily_auto_run' AND input_ref=? AND node_action='daily_auto_run_reuse'`, first.IdempotencyKey).Scan(&reuseAuditCount); err != nil {
		t.Fatal(err)
	}
	if stateCount != 1 || notificationCount != 1 || originalAuditCount != 1 || reuseAuditCount != 1 {
		t.Fatalf("expected idempotent state/notification and reuse audit, got state=%d notification=%d original_audit=%d reuse_audit=%d", stateCount, notificationCount, originalAuditCount, reuseAuditCount)
	}
}

func TestDailyAutoRunMissingHoldingsWritesDiagnostics(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	runner := NewDailyAutoRunner(config.DailyAutoRunConfig{Enabled: true, RunTime: "08:30", Timezone: "Asia/Shanghai", Scope: "holdings", Retry: 1, TimeoutSeconds: 900, MaxSymbols: 20}, deps)

	out, err := runner.RunOnce(context.Background(), time.Date(2026, 6, 7, 0, 30, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	if out.Status != "failed" || out.FailureCode != "missing_prerequisites" {
		t.Fatalf("expected missing prerequisites failure, got %+v", out)
	}

	var stateStatus, failureCode string
	if err := db.QueryRow(`SELECT status, failure_code FROM daily_auto_run_states WHERE idempotency_key=?`, out.IdempotencyKey).Scan(&stateStatus, &failureCode); err != nil {
		t.Fatal(err)
	}
	if stateStatus != "failed" || failureCode != "missing_prerequisites" {
		t.Fatalf("unexpected persisted state status=%s failure=%s", stateStatus, failureCode)
	}
	var notificationCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM notifications WHERE type='daily_auto_run_failed' AND source_type='daily_auto_run' AND source_id=?`, out.IdempotencyKey).Scan(&notificationCount); err != nil {
		t.Fatal(err)
	}
	if notificationCount != 1 {
		t.Fatalf("expected one diagnostic notification, got %d", notificationCount)
	}
	var auditCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM audit_events WHERE workflow_type='daily_auto_run' AND action='run_local_task' AND status='failed' AND error_code='missing_prerequisites' AND output_ref LIKE '%step=prerequisites%'`).Scan(&auditCount); err != nil {
		t.Fatal(err)
	}
	if auditCount != 1 {
		t.Fatalf("expected failed diagnostic audit, got %d", auditCount)
	}

	var report repository.DailyDisciplineReport
	if err := db.QueryRow(`SELECT report_id,local_date,scope,symbol_set_hash,source_type,source_id,COALESCE(decision_id,''),status,summary,COALESCE(failure_code,''),COALESCE(failure_reason,''),created_at,updated_at FROM daily_discipline_reports WHERE source_type='auto_run' AND source_id=?`, out.IdempotencyKey).Scan(&report.ReportID, &report.LocalDate, &report.Scope, &report.SymbolSetHash, &report.SourceType, &report.SourceID, &report.DecisionID, &report.Status, &report.Summary, &report.FailureCode, &report.FailureReason, &report.CreatedAt, &report.UpdatedAt); err != nil {
		t.Fatal(err)
	}
	if report.Status != "insufficient_data" || report.FailureCode != "missing_prerequisites" || report.FailureReason == "" || report.Summary == "" {
		t.Fatalf("unexpected insufficient-data report row: %+v", report)
	}
}

func TestDailyAutoRunWithoutReportRepoKeepsLegacyPath(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	deps.DailyDisciplineReportRepo = nil
	runner := NewDailyAutoRunner(config.DailyAutoRunConfig{Enabled: true, RunTime: "08:30", Timezone: "Asia/Shanghai", Scope: "holdings", Retry: 1, TimeoutSeconds: 900, MaxSymbols: 20}, deps)

	out, err := runner.RunOnce(context.Background(), time.Date(2026, 6, 7, 0, 30, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("RunOnce with nil report repo: %v", err)
	}
	if out.Status != "failed" || out.FailureCode != "missing_prerequisites" {
		t.Fatalf("expected legacy missing prerequisites output, got %+v", out)
	}
}

func TestDailyAutoRunRunOncePersistsConfiguredNextRunAt(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	runner := NewDailyAutoRunner(config.DailyAutoRunConfig{Enabled: true, RunTime: "08:30", Timezone: "Asia/Shanghai", Scope: "holdings", Retry: 1, TimeoutSeconds: 900, MaxSymbols: 20}, deps)
	runner.SetNextRunAt("2026-06-08T08:30:00+08:00")

	out, err := runner.RunOnce(context.Background(), time.Date(2026, 6, 7, 0, 30, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	var nextRunAt string
	if err := db.QueryRow(`SELECT COALESCE(next_run_at,'') FROM daily_auto_run_states WHERE idempotency_key=?`, out.IdempotencyKey).Scan(&nextRunAt); err != nil {
		t.Fatal(err)
	}
	if nextRunAt != "2026-06-08T08:30:00+08:00" {
		t.Fatalf("expected configured next_run_at, got %q", nextRunAt)
	}
}

func TestDailyAutoRunRepeatedReuseWithSameNowDoesNotConflict(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	runner := NewDailyAutoRunner(config.DailyAutoRunConfig{Enabled: true, RunTime: "08:30", Timezone: "Asia/Shanghai", Scope: "holdings", Retry: 1, TimeoutSeconds: 900, MaxSymbols: 20}, deps)
	runAt := time.Date(2026, 6, 7, 0, 30, 0, 0, time.UTC)

	first, err := runner.RunOnce(context.Background(), runAt)
	if err != nil {
		t.Fatalf("first RunOnce: %v", err)
	}
	if _, err := runner.RunOnce(context.Background(), runAt); err != nil {
		t.Fatalf("second RunOnce: %v", err)
	}
	third, err := runner.RunOnce(context.Background(), runAt)
	if err != nil {
		t.Fatalf("third RunOnce with same now should reuse without audit conflict: %v", err)
	}
	if third.IdempotencyKey != first.IdempotencyKey || third.Status != first.Status {
		t.Fatalf("expected same reused output, got first=%+v third=%+v", first, third)
	}
}

func TestDailyAutoRunIdempotentReuseBackfillsMissingReportRow(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	runner := NewDailyAutoRunner(config.DailyAutoRunConfig{Enabled: true, RunTime: "08:30", Timezone: "Asia/Shanghai", Scope: "holdings", Retry: 1, TimeoutSeconds: 900, MaxSymbols: 20}, deps)
	ctx := context.Background()
	runAt := time.Date(2026, 6, 7, 0, 30, 0, 0, time.UTC)
	key := runner.idempotencyKey(runAt, nil)
	state := repository.DailyAutoRunState{RunID: stableAutoRunID(key), IdempotencyKey: key, LocalDate: "2026-06-07", Scope: "holdings", SymbolSetHash: symbolSetHash(nil), Status: "failed", LastRunAt: runAt.UTC().Format(time.RFC3339), FailureCode: "missing_prerequisites", FailureReason: "缺少本地持仓", CreatedAt: runAt.UTC().Format(time.RFC3339), UpdatedAt: runAt.UTC().Format(time.RFC3339)}
	if err := deps.DailyAutoRunRepo.UpsertDailyAutoRunState(ctx, state); err != nil {
		t.Fatal(err)
	}

	out, err := runner.RunOnce(ctx, runAt)
	if err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	if out.Status != "failed" || out.FailureCode != "missing_prerequisites" || out.IdempotencyKey != key {
		t.Fatalf("expected reused missing prerequisites output, got %+v", out)
	}
	var report repository.DailyDisciplineReport
	if err := db.QueryRow(`SELECT source_type,source_id,status,COALESCE(failure_code,''),COALESCE(failure_reason,'') FROM daily_discipline_reports WHERE source_type='auto_run' AND source_id=?`, key).Scan(&report.SourceType, &report.SourceID, &report.Status, &report.FailureCode, &report.FailureReason); err != nil {
		t.Fatal(err)
	}
	if report.Status != "insufficient_data" || report.FailureCode != "missing_prerequisites" || report.FailureReason != "缺少本地持仓" {
		t.Fatalf("unexpected backfilled report: %+v", report)
	}
}

func TestDailyAutoRunIdempotentReuseWritesReuseAuditOnly(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	runner := NewDailyAutoRunner(config.DailyAutoRunConfig{Enabled: true, RunTime: "08:30", Timezone: "Asia/Shanghai", Scope: "holdings", Retry: 1, TimeoutSeconds: 900, MaxSymbols: 20}, deps)
	runAt := time.Date(2026, 6, 7, 0, 30, 0, 0, time.UTC)

	first, err := runner.RunOnce(context.Background(), runAt)
	if err != nil {
		t.Fatalf("first RunOnce: %v", err)
	}
	second, err := runner.RunOnce(context.Background(), runAt)
	if err != nil {
		t.Fatalf("second RunOnce: %v", err)
	}
	if second.Status != "failed" || second.IdempotencyKey != first.IdempotencyKey {
		t.Fatalf("expected reused failed output, got first=%+v second=%+v", first, second)
	}
	var failedAuditCount, reuseAuditCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM audit_events WHERE workflow_type='daily_auto_run' AND input_ref=? AND node_action='daily_auto_run'`, first.IdempotencyKey).Scan(&failedAuditCount); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM audit_events WHERE workflow_type='daily_auto_run' AND input_ref=? AND node_action='daily_auto_run_reuse' AND status='success'`, first.IdempotencyKey).Scan(&reuseAuditCount); err != nil {
		t.Fatal(err)
	}
	if failedAuditCount != 1 || reuseAuditCount != 1 {
		t.Fatalf("expected one original audit and one reuse audit, got original=%d reuse=%d", failedAuditCount, reuseAuditCount)
	}
}

func TestDailyAutoRunWritesRunningStateBeforeSideEffects(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	ctx := context.Background()
	seedDailyAutoRunHolding(t, deps)
	source := &observingMarketSource{t: t, repo: deps.DailyAutoRunRepo}
	deps.MarketDataSource = source
	runner := NewDailyAutoRunner(config.DailyAutoRunConfig{Enabled: true, RunTime: "08:30", Timezone: "Asia/Shanghai", Scope: "holdings", Retry: 0, TimeoutSeconds: 900, MaxSymbols: 20}, deps)

	out, err := runner.RunOnce(ctx, time.Date(2026, 6, 7, 0, 30, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	if out.Status != "degraded" || !source.observedRunning {
		t.Fatalf("expected running state before side effects and degraded output, observed_running=%v out=%+v", source.observedRunning, out)
	}
}

func TestDailyAutoRunRunningStateReusesWithoutDuplicateSideEffects(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	ctx := context.Background()
	seedDailyAutoRunHolding(t, deps)
	now := time.Date(2026, 6, 7, 0, 30, 0, 0, time.UTC)
	runner := NewDailyAutoRunner(config.DailyAutoRunConfig{Enabled: true, RunTime: "08:30", Timezone: "Asia/Shanghai", Scope: "holdings", Retry: 0, TimeoutSeconds: 900, MaxSymbols: 20}, deps)
	key := runner.idempotencyKey(now, []repository.Position{{Symbol: "510300"}})
	if err := deps.DailyAutoRunRepo.UpsertDailyAutoRunState(ctx, repository.DailyAutoRunState{RunID: stableAutoRunID(key), IdempotencyKey: key, LocalDate: "2026-06-07", Scope: "holdings", SymbolSetHash: symbolSetHash([]repository.Position{{Symbol: "510300"}}), Status: "running", LastRunAt: now.UTC().Format(time.RFC3339), CreatedAt: now.UTC().Format(time.RFC3339), UpdatedAt: now.UTC().Format(time.RFC3339)}); err != nil {
		t.Fatal(err)
	}

	out, err := runner.RunOnce(ctx, now)
	if err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	if out.Status != "running" || out.IdempotencyKey != key {
		t.Fatalf("expected running output reuse, got %+v", out)
	}
	var decisionCount, marketCount, notificationCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM decision_records WHERE workflow_type='daily_discipline'`).Scan(&decisionCount); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM market_snapshots WHERE symbol='510300'`).Scan(&marketCount); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM notifications WHERE source_id=?`, key).Scan(&notificationCount); err != nil {
		t.Fatal(err)
	}
	if decisionCount != 0 || marketCount != 0 || notificationCount != 0 {
		t.Fatalf("running state must not duplicate side effects, got decisions=%d markets=%d notifications=%d", decisionCount, marketCount, notificationCount)
	}
}

func TestDailyAutoRunStaleRunningStateRerunsWorkflows(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	ctx := context.Background()
	seedDailyAutoRunHolding(t, deps)
	source := &countingSuccessMarketSource{}
	deps.MarketDataSource = source
	now := time.Date(2026, 6, 7, 0, 30, 0, 0, time.UTC)
	runner := NewDailyAutoRunner(config.DailyAutoRunConfig{Enabled: true, RunTime: "08:30", Timezone: "Asia/Shanghai", Scope: "holdings", Retry: 0, TimeoutSeconds: 900, MaxSymbols: 20}, deps)
	positions := []repository.Position{{Symbol: "510300"}}
	key := runner.idempotencyKey(now, positions)
	staleAt := now.Add(-2 * time.Hour).UTC().Format(time.RFC3339)
	if err := deps.DailyAutoRunRepo.UpsertDailyAutoRunState(ctx, repository.DailyAutoRunState{RunID: stableAutoRunID(key), IdempotencyKey: key, LocalDate: "2026-06-07", Scope: "holdings", SymbolSetHash: symbolSetHash(positions), Status: "running", LastRunAt: staleAt, CreatedAt: staleAt, UpdatedAt: staleAt}); err != nil {
		t.Fatal(err)
	}

	out, err := runner.RunOnce(ctx, now)
	if err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	if out.Status == "running" || source.calls == 0 {
		t.Fatalf("expected stale running state to rerun and finish, calls=%d out=%+v", source.calls, out)
	}
	var stateStatus string
	if err := db.QueryRow(`SELECT status FROM daily_auto_run_states WHERE idempotency_key=?`, key).Scan(&stateStatus); err != nil {
		t.Fatal(err)
	}
	if stateStatus == "running" {
		t.Fatalf("expected final persisted state not running")
	}
}

func TestDailyAutoRunRetriesMarketRefreshBeforeFailure(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	ctx := context.Background()
	seedDailyAutoRunHolding(t, deps)
	deps.MarketDataSource = &countingFailMarketSource{err: errors.New("temporary market failure")}
	runner := NewDailyAutoRunner(config.DailyAutoRunConfig{Enabled: true, RunTime: "08:30", Timezone: "Asia/Shanghai", Scope: "holdings", Retry: 2, TimeoutSeconds: 900, MaxSymbols: 20}, deps)

	out, err := runner.RunOnce(ctx, time.Date(2026, 6, 7, 0, 30, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	if out.Status != "failed" || out.FailureCode != "market_refresh_failed" {
		t.Fatalf("expected market refresh failure, got %+v", out)
	}
	source := deps.MarketDataSource.(*countingFailMarketSource)
	if source.calls != 3 {
		t.Fatalf("expected initial attempt plus 2 retries, got %d calls", source.calls)
	}
	var retryAuditCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM audit_events WHERE workflow_type='daily_auto_run' AND input_ref=? AND node_action='daily_auto_run_retry'`, out.IdempotencyKey).Scan(&retryAuditCount); err != nil {
		t.Fatal(err)
	}
	if retryAuditCount != 2 {
		t.Fatalf("expected 2 retry audits, got %d", retryAuditCount)
	}
}

func TestDailyAutoRunMarketRefreshFailureWritesFailedReport(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	ctx := context.Background()
	seedDailyAutoRunHolding(t, deps)
	deps.MarketDataSource = &countingFailMarketSource{err: errors.New("temporary market failure")}
	runner := NewDailyAutoRunner(config.DailyAutoRunConfig{Enabled: true, RunTime: "08:30", Timezone: "Asia/Shanghai", Scope: "holdings", Retry: 0, TimeoutSeconds: 900, MaxSymbols: 20}, deps)

	out, err := runner.RunOnce(ctx, time.Date(2026, 6, 7, 0, 30, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	if out.Status != "failed" || out.FailureCode != "market_refresh_failed" {
		t.Fatalf("expected market refresh failure, got %+v", out)
	}
	var report repository.DailyDisciplineReport
	if err := db.QueryRow(`SELECT status,COALESCE(failure_code,''),COALESCE(failure_reason,'') FROM daily_discipline_reports WHERE source_type='auto_run' AND source_id=?`, out.IdempotencyKey).Scan(&report.Status, &report.FailureCode, &report.FailureReason); err != nil {
		t.Fatal(err)
	}
	if report.Status != "failed" || report.FailureCode != "market_refresh_failed" || report.FailureReason == "" {
		t.Fatalf("unexpected market refresh failure report: %+v", report)
	}
}

func TestDailyAutoRunTimeoutWritesReadableDiagnostics(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	ctx := context.Background()
	seedDailyAutoRunHolding(t, deps)
	deps.MarketDataSource = blockingMarketSource{}
	runner := NewDailyAutoRunner(config.DailyAutoRunConfig{Enabled: true, RunTime: "08:30", Timezone: "Asia/Shanghai", Scope: "holdings", Retry: 0, TimeoutSeconds: 1, MaxSymbols: 20}, deps)

	out, err := runner.RunOnce(ctx, time.Date(2026, 6, 7, 0, 30, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	if out.Status != "failed" || out.FailureCode != "timeout" || !strings.Contains(out.FailureReason, "超时") {
		t.Fatalf("expected timeout failure, got %+v", out)
	}
	var auditOutput, notificationMessage string
	if err := db.QueryRow(`SELECT output_ref FROM audit_events WHERE workflow_type='daily_auto_run' AND input_ref=? AND node_action='daily_auto_run'`, out.IdempotencyKey).Scan(&auditOutput); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(`SELECT message FROM notifications WHERE source_id=?`, out.IdempotencyKey).Scan(&notificationMessage); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(auditOutput, "step=timeout") || !strings.Contains(notificationMessage, "超时") {
		t.Fatalf("expected timeout diagnostics, got audit=%q notification=%q", auditOutput, notificationMessage)
	}
}

func seedDailyAutoRunHolding(t *testing.T, deps WorkflowDependencies) {
	t.Helper()
	ctx := context.Background()
	if err := deps.PortfolioRepo.SavePortfolioSnapshot(ctx, repository.PortfolioSnapshot{SnapshotID: "portfolio_auto", SnapshotTime: "2026-06-07T00:00:00Z", Cash: 1000, TotalAssets: 10000, CashRatio: 0.1, PositionCount: 1, Source: "manual", CreatedAt: "2026-06-07T00:00:00Z"}, []repository.PositionSnapshot{{PositionSnapshotID: "ps_auto", SnapshotID: "portfolio_auto", Symbol: "510300", Name: "沪深300ETF", Quantity: 100, CostPrice: 4, CurrentPrice: 4.2, MarketValue: 420, UnrealizedProfitRatio: 0.05, PositionState: "normal", CreatedAt: "2026-06-07T00:00:00Z"}}); err != nil {
		t.Fatal(err)
	}
	if err := deps.PortfolioRepo.SavePosition(ctx, repository.Position{PositionID: "pos_auto", Symbol: "510300", Name: "沪深300ETF", Quantity: 100, CostPrice: 4, CurrentPrice: 4.2, MarketValue: 420, UnrealizedProfitRatio: 0.05, PositionState: "normal", UpdatedAt: "2026-06-07T00:00:00Z"}); err != nil {
		t.Fatal(err)
	}
}

type observingMarketSource struct {
	t               *testing.T
	repo            repository.DailyAutoRunRepository
	observedRunning bool
}

func (s *observingMarketSource) FetchMarketData(ctx context.Context, _ string) (MarketDataPoint, error) {
	state, err := s.repo.GetLatestDailyAutoRunState(ctx)
	if err != nil {
		s.t.Fatalf("expected running state before market refresh: %v", err)
	}
	s.observedRunning = state.Status == "running"
	return MarketDataPoint{ClosePrice: 4.21, SourceName: "stub", SourceLevel: model.SourceLevelB, SourceType: "daily_auto_run_test", TradeDate: "2026-06-07", PEPercentile: 50, PBPercentile: 45}, nil
}

type countingSuccessMarketSource struct {
	calls int
}

func (s *countingSuccessMarketSource) FetchMarketData(context.Context, string) (MarketDataPoint, error) {
	s.calls++
	return MarketDataPoint{ClosePrice: 4.21, SourceName: "stub", SourceLevel: model.SourceLevelB, SourceType: "daily_auto_run_test", TradeDate: "2026-06-07", PEPercentile: 50, PBPercentile: 45}, nil
}

type countingFailMarketSource struct {
	err   error
	calls int
}

func (s *countingFailMarketSource) FetchMarketData(context.Context, string) (MarketDataPoint, error) {
	s.calls++
	return MarketDataPoint{}, s.err
}

type blockingMarketSource struct{}

func (blockingMarketSource) FetchMarketData(ctx context.Context, _ string) (MarketDataPoint, error) {
	<-ctx.Done()
	return MarketDataPoint{}, ctx.Err()
}
