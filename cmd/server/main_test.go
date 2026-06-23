package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"investment-agent/internal/application/workflow"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/infrastructure/config"
	appsqlite "investment-agent/internal/infrastructure/persistence/sqlite"
)

func TestStartDailyAutoRunSchedulerDisabledDoesNotWriteState(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	if err := appsqlite.Migrate(context.Background(), store.DB); err != nil {
		t.Fatal(err)
	}
	runner := workflow.NewDailyAutoRunner(config.DailyAutoRunConfig{Enabled: false}, workflow.WorkflowDependencies{})
	stop := startDailyAutoRunScheduler(context.Background(), runner, config.DailyAutoRunConfig{Enabled: false}, slog.Default())
	defer stop()
	time.Sleep(10 * time.Millisecond)
	var count int
	if err := store.DB.QueryRow(`SELECT COUNT(*) FROM daily_auto_run_states`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("disabled scheduler must not write state, got %d", count)
	}
}

func TestStartDailyAutoRunSchedulerEnabledWaitsForConfiguredRunTime(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	if err := appsqlite.Migrate(context.Background(), store.DB); err != nil {
		t.Fatal(err)
	}
	transactor := appsqlite.NewTransactor(store.DB)
	repos := repository.Repositories{
		AuditRepo:        appsqlite.NewAuditRepository(store.DB),
		NotificationRepo: appsqlite.NewNotificationRepository(store.DB),
		PortfolioRepo:    appsqlite.NewPortfolioRepository(store.DB),
		DailyAutoRunRepo: appsqlite.NewDailyAutoRunRepository(store.DB),
	}
	runner := workflow.NewDailyAutoRunner(config.DailyAutoRunConfig{Enabled: true, RunTime: "23:59", Timezone: "Asia/Shanghai", Scope: "holdings", Retry: 0, TimeoutSeconds: 1, MaxSymbols: 20}, workflow.NewWorkflowDependencies(repos, transactor))
	stop := startDailyAutoRunScheduler(context.Background(), runner, config.DailyAutoRunConfig{Enabled: true, RunTime: "23:59", Timezone: "Asia/Shanghai", Scope: "holdings", Retry: 0, TimeoutSeconds: 1, MaxSymbols: 20}, slog.Default())
	defer stop()
	time.Sleep(20 * time.Millisecond)
	var count int
	if err := store.DB.QueryRow(`SELECT COUNT(*) FROM daily_auto_run_states`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("enabled scheduler must wait for configured run time, got %d states", count)
	}
}

func TestRunValidatesConfigBeforeStartingServer(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(`server:
  host: "127.0.0.1"
  port: 0
sqlite:
  path: "`+filepath.Join(dir, "agent.db")+`"
veclite:
  path: "`+filepath.Join(dir, "veclite")+`"
data_sources:
  enabled:
    - "public-http"
  use_stub: false
  market_endpoint: ""
  intelligence_endpoint: ""
log:
  level: "info"
`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	t.Setenv("INVESTMENT_AGENT_CONFIG", configPath)

	err := run()

	if err == nil || !strings.Contains(err.Error(), "validate config") || !strings.Contains(err.Error(), "market_endpoint") {
		t.Fatalf("expected config validation error before server start, got %v", err)
	}
}
