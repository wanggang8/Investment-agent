package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"investment-agent/internal/application/service"
	"investment-agent/internal/application/workflow"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	appsqlite "investment-agent/internal/infrastructure/persistence/sqlite"
)

func TestRunHelpShowsLocalTasksAndNoTrading(t *testing.T) {
	var stdout bytes.Buffer
	code := run(context.Background(), []string{"--help"}, &stdout, &stdout)

	if code != 0 {
		t.Fatalf("expected help exit code 0, got %d", code)
	}
	out := stdout.String()
	for _, want := range []string{"daily", "market-refresh", "evidence-index", "review", "public-evidence-refresh", "p34-expanded-refresh", "retrieval-quality-smoke", "data-source-quality-regression", "data-source-quality-resolution-check", "--source", "--start-date", "--end-date", "不会执行交易", "本地调度", "audit_events", "不会自动应用规则"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected help to contain %q, got %s", want, out)
		}
	}
}

func TestRunScheduleShowsSafeLocalSchedulerBoundary(t *testing.T) {
	var stdout bytes.Buffer
	code := run(context.Background(), []string{"--schedule"}, &stdout, &stdout)

	if code != 0 {
		t.Fatalf("expected schedule exit code 0, got %d", code)
	}
	out := stdout.String()
	for _, want := range []string{"默认不自动运行", "需要用户显式安装", "不会执行交易", "不会自动应用规则", "docs/ops-local-scheduler.md"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected schedule output to contain %q, got %s", want, out)
		}
	}
	for _, banned := range bannedAutomationTerms() {
		if strings.Contains(out, banned) {
			t.Fatalf("schedule output contains banned automation term %q: %s", banned, out)
		}
	}
}

func TestSchedulerExamplesAreSafeAndPlaceholderOnly(t *testing.T) {
	paths := []string{
		filepath.Join("..", "..", "docs", "ops-local-scheduler.md"),
		filepath.Join("..", "..", "examples", "scheduler", "launchd", "com.example.investment-agent.plist"),
		filepath.Join("..", "..", "examples", "scheduler", "cron", "investment-agent.cron"),
	}
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read scheduler artifact %s: %v", path, err)
		}
		text := string(content)
		for _, want := range []string{"/ABSOLUTE/PATH/TO/Investment-agent", "INVESTMENT_AGENT_CONFIG", "不会执行交易", "不会自动应用规则"} {
			if !strings.Contains(text, want) {
				t.Fatalf("expected %s to contain %q, got %s", path, want, text)
			}
		}
		for _, banned := range bannedAutomationTerms() {
			if strings.Contains(text, banned) {
				t.Fatalf("%s contains banned automation term %q", path, banned)
			}
		}
	}
}

func bannedAutomationTerms() []string {
	return []string{"broker", "order_id", "BUY ", "SELL ", "auto-confirm", "auto_apply", "place-order", "委托下单", "自动交易", "一键交易", "自动确认", "自动生效", "规则生效"}
}

func TestRunDataSourceQualityRegressionFixtureWritesSanitizedAuditWithoutTrading(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	configPath := writeTestConfig(t, dbPath)

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "data-source-quality-regression", "--source", "fixture", "--symbol", "000300"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected data source quality regression exit code 0, got %d stderr=%s", code, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{"data source quality regression completed", "status=passed", "cases=6", "不会执行交易"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected stdout to contain %q, got %s", want, out)
		}
	}
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	for _, table := range []string{"positions", "portfolio_snapshots", "operation_confirmations", "position_transactions", "local_account_import_batches", "local_account_corrections"} {
		assertTableCount(t, store.DB, table, 0)
	}
	assertAuditEvent(t, store.DB, "run_local_task", "data-source-quality-regression:source=fixture:symbol=000300", "data_source_quality:mode=fixture:status=passed:policy=passed:gate=pass:cases=6:degraded=0:failed=0:no_auto_trading")
	var auditOutput string
	if err := store.DB.QueryRow(`SELECT output_ref FROM audit_events WHERE input_ref='data-source-quality-regression:source=fixture:symbol=000300'`).Scan(&auditOutput); err != nil {
		t.Fatalf("read regression audit output: %v", err)
	}
	for _, forbidden := range []string{"sk-123456789012", "/Users/private", "select    *    from", "prompt:", "raw HTTP", "BEGIN RSA PRIVATE KEY"} {
		if strings.Contains(auditOutput, forbidden) || strings.Contains(out, forbidden) {
			t.Fatalf("regression output leaked %q stdout=%s audit=%s", forbidden, out, auditOutput)
		}
	}
}

func TestRunDataSourceQualityRegressionCurrentReadsExistingSnapshotOnly(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	configPath := writeTestConfig(t, dbPath)
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := appsqlite.Migrate(context.Background(), store.DB); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	metrics := `{"source_name":"csindex","source_level":"A","source_type":"index_basic","captured_at":"2026-06-06T01:00:00Z","metadata":{"p34_source_health":{"index_constituents":"fresh","index_valuation_files":"parse_error"},"p34_data_categories":["index_constituents","index_valuation_files"]}}`
	_, err = store.DB.Exec(`INSERT INTO market_snapshots (market_snapshot_id,symbol,trade_date,close_price,turnover_rate,pe_percentile,pb_percentile,volume_percentile,volatility_percentile,liquidity_state,sentiment_state,market_metrics_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`, "market_p48_cli", "000300", "2026-06-05", 0, 0, 0, 0, 0, 0, "normal", "neutral", metrics, "2026-06-06T01:00:00Z")
	if err != nil {
		t.Fatalf("seed market: %v", err)
	}
	store.Close()

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "data-source-quality-regression", "--source", "current", "--symbol", "000300"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected current regression exit code 0, got %d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "status=degraded") || !strings.Contains(stdout.String(), "cases=2") {
		t.Fatalf("expected degraded current summary, got %s", stdout.String())
	}
	store, err = appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite after regression: %v", err)
	}
	defer store.Close()
	assertTableCount(t, store.DB, "market_snapshots", 1)
	assertAuditEvent(t, store.DB, "run_local_task", "data-source-quality-regression:source=current:symbol=000300", "data_source_quality:mode=current:status=degraded:policy=blocked:gate=block:cases=2:degraded=1:failed=0:no_auto_trading")
}

func TestRunDataSourceQualityRegressionStrictGateFailsWhenPolicyBlocks(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	configPath := writeTestConfig(t, dbPath)
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := appsqlite.Migrate(context.Background(), store.DB); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	metrics := `{"source_name":"csindex","source_level":"A","source_type":"index_basic","captured_at":"2026-06-06T01:00:00Z","metadata":{"p34_source_health":{"index_constituents":"fresh","index_valuation_files":"parse_error"},"p34_data_categories":["index_constituents","index_valuation_files"]}}`
	_, err = store.DB.Exec(`INSERT INTO market_snapshots (market_snapshot_id,symbol,trade_date,close_price,turnover_rate,pe_percentile,pb_percentile,volume_percentile,volatility_percentile,liquidity_state,sentiment_state,market_metrics_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`, "market_p66_cli_gate", "000300", "2026-06-05", 0, 0, 0, 0, 0, 0, "normal", "neutral", metrics, "2026-06-06T01:00:00Z")
	if err != nil {
		t.Fatalf("seed market: %v", err)
	}
	store.Close()

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "data-source-quality-regression", "--source", "current", "--symbol", "000300", "--strict-quality-gate"}, &stdout, &stderr)

	if code == 0 {
		t.Fatalf("expected strict quality gate to fail, stdout=%s stderr=%s", stdout.String(), stderr.String())
	}
	for _, want := range []string{"policy=blocked", "gate=block"} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("expected stderr to contain %q, got %s", want, stderr.String())
		}
	}
	store, err = appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite after strict gate: %v", err)
	}
	defer store.Close()
	assertAuditEvent(t, store.DB, "run_local_task", "data-source-quality-regression:source=current:symbol=000300", "data_source_quality:mode=current:status=degraded:policy=blocked:gate=block:cases=2:degraded=1:failed=0:no_auto_trading")
}

func TestRunDataSourceQualityResolutionCheckFailsWhenResolutionIsRequired(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	configPath := writeTestConfig(t, dbPath)
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := appsqlite.Migrate(context.Background(), store.DB); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	seedAgentBlockedCurrentSourceHealth(t, store.DB, "market_p67_cli_required")
	store.Close()

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "data-source-quality-resolution-check", "--symbol", "000300"}, &stdout, &stderr)

	if code == 0 {
		t.Fatalf("expected resolution check to fail, stdout=%s stderr=%s", stdout.String(), stderr.String())
	}
	for _, want := range []string{"requires_resolution", "clean_data_claim=false", "no_auto_trading"} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("expected stderr to contain %q, got %s", want, stderr.String())
		}
	}
	store, err = appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite after check: %v", err)
	}
	defer store.Close()
	assertAuditEventContains(t, store.DB, "run_local_task", "data-source-quality-resolution-check:symbol=000300", []string{"data_quality_gate_resolution:claim_state=requires_resolution:policy=blocked:gate=block:fingerprint=", ":resolution=none:clean_data_claim=false:no_auto_trading"})
}

func TestRunDataSourceQualityResolutionCheckPassesWithScopeExclusion(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	configPath := writeTestConfig(t, dbPath)
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := appsqlite.Migrate(context.Background(), store.DB); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	seedAgentBlockedCurrentSourceHealth(t, store.DB, "market_p67_cli_scope")
	repos := repository.Repositories{
		MarketRepo:                    appsqlite.NewMarketRepository(store.DB),
		AuditRepo:                     appsqlite.NewAuditRepository(store.DB),
		DataQualityGateResolutionRepo: appsqlite.NewDataQualityGateResolutionRepository(store.DB),
	}
	if _, err := service.NewDataSourceQualityService(repos).CreateGateResolution(context.Background(), service.DataQualityGateResolutionCreateRequest{
		Symbol:         "000300",
		ResolutionType: service.DataQualityGateResolutionTypeScopeExclusion,
		Scope:          "本次 release clean claim 排除 current local data health",
		Reason:         "当前源解析降级，发布只声明有限范围",
		ReleaseImpact:  "不得声明 current data healthy",
		EvidenceRef:    "docs/release/acceptance/p66",
	}); err != nil {
		t.Fatalf("create resolution: %v", err)
	}
	store.Close()

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "data-source-quality-resolution-check", "--symbol", "000300"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected resolution check exit code 0, got %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	for _, want := range []string{"resolved_with_scope_exclusion", "clean_data_claim=false", "不会执行交易"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("expected stdout to contain %q, got %s", want, stdout.String())
		}
	}
}

func TestRunDataSourceQualityRegressionRejectsUnsupportedSource(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	configPath := writeTestConfig(t, dbPath)

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "data-source-quality-regression", "--source", "real"}, &stdout, &stderr)
	if code == 0 || !strings.Contains(stderr.String(), "unsupported data source quality regression mode") {
		t.Fatalf("expected unsupported regression source error, code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
}

func seedAgentBlockedCurrentSourceHealth(t *testing.T, db *sql.DB, snapshotID string) {
	t.Helper()
	metrics := `{"source_name":"csindex","source_level":"A","source_type":"index_basic","captured_at":"2026-06-06T01:00:00Z","metadata":{"p34_source_health":{"index_constituents":"fresh","index_valuation_files":"parse_error"},"p34_data_categories":["index_constituents","index_valuation_files"]}}`
	_, err := db.Exec(`INSERT INTO market_snapshots (market_snapshot_id,symbol,trade_date,close_price,turnover_rate,pe_percentile,pb_percentile,volume_percentile,volatility_percentile,liquidity_state,sentiment_state,market_metrics_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`, snapshotID, "000300", "2026-06-05", 0, 0, 0, 0, 0, 0, "normal", "neutral", metrics, "2026-06-06T01:00:00Z")
	if err != nil {
		t.Fatalf("seed market: %v", err)
	}
}

func TestRunDataSourceQualityRegressionSanitizesFailedAuditInput(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	configPath := writeTestConfig(t, dbPath)

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "data-source-quality-regression", "--source", "sk-123456789012", "--symbol", "/Users/private/secret.txt"}, &stdout, &stderr)
	if code == 0 {
		t.Fatalf("expected unsupported regression source failure")
	}
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	var inputRef, outputRef string
	if err := store.DB.QueryRow(`SELECT input_ref,output_ref FROM audit_events WHERE workflow_type='data-source-quality-regression' ORDER BY created_at DESC LIMIT 1`).Scan(&inputRef, &outputRef); err != nil {
		t.Fatalf("read failed regression audit: %v", err)
	}
	if inputRef != "data-source-quality-regression:source=unsupported:symbol=redacted" || outputRef != "task_failed" {
		t.Fatalf("expected sanitized failed audit, input=%s output=%s", inputRef, outputRef)
	}
	for _, forbidden := range []string{"sk-123456789012", "/Users/private"} {
		if strings.Contains(inputRef, forbidden) || strings.Contains(outputRef, forbidden) {
			t.Fatalf("failed audit leaked %q input=%s output=%s", forbidden, inputRef, outputRef)
		}
	}
}

func TestRunManualTaskWritesAuditWithoutTrading(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	configPath := writeTestConfig(t, dbPath)

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "review", "--period", "monthly"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected task exit code 0, got %d, stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "不会执行交易") {
		t.Fatalf("expected safety output, got %s", stdout.String())
	}
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	assertAuditEvent(t, store.DB, "run_local_task", "review:monthly", "no_auto_trading")
}

func TestRunMarketRefreshWritesSnapshotAndAudit(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	configPath := writeTestConfig(t, dbPath)

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "market-refresh"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected market task exit code 0, got %d, stderr=%s", code, stderr.String())
	}
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	assertTableCount(t, store.DB, "market_snapshots", 1)
	assertAuditEvent(t, store.DB, "refresh_market_data", "market-refresh", "no_auto_trading")
}

func TestRunP34ExpandedRefreshWritesSourceHealthAndAudit(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	configPath := writeTestConfig(t, dbPath)

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "p34-expanded-refresh", "--source", "sentiment_proxy_fixture", "--symbol", "000300", "--start-date", "2026-06-01", "--end-date", "2026-06-05"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected P34 refresh exit code 0, got %d, stderr=%s", code, stderr.String())
	}
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	var metrics string
	if err := store.DB.QueryRow(`SELECT market_metrics_json FROM market_snapshots WHERE symbol='000300' ORDER BY created_at DESC LIMIT 1`).Scan(&metrics); err != nil {
		t.Fatalf("read market metrics: %v", err)
	}
	if !strings.Contains(metrics, "sentiment_proxy") || !strings.Contains(metrics, "p34_source_health") {
		t.Fatalf("expected P34 sentiment source health metrics, got %s", metrics)
	}
	for _, table := range []string{"positions", "portfolio_snapshots", "operation_confirmations", "position_transactions", "local_account_import_batches", "local_account_corrections"} {
		assertTableCount(t, store.DB, table, 0)
	}
	assertAuditEvent(t, store.DB, "run_local_task", "p34-expanded-refresh:source=sentiment_proxy_fixture:symbol=000300:start=2026-06-01:end=2026-06-05", "no_auto_trading")
}

func TestRunP34ExpandedRefreshRejectsUnsupportedSource(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	configPath := writeTestConfig(t, dbPath)

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "p34-expanded-refresh", "--source", "unknown_source", "--symbol", "000300"}, &stdout, &stderr)
	if code == 0 || !strings.Contains(stderr.String(), "unsupported P34 source") {
		t.Fatalf("expected unsupported P34 source error, code=%d stderr=%s", code, stderr.String())
	}
}

func TestRunLLMSmokeCallsConfiguredModelAndWritesSanitizedAudit(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	var seenModel, authHeader, prompt string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Fatalf("unexpected LLM path: %s", r.URL.Path)
		}
		authHeader = r.Header.Get("Authorization")
		var body struct {
			Model    string `json:"model"`
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode LLM request: %v", err)
		}
		seenModel = body.Model
		if len(body.Messages) > 0 {
			prompt = body.Messages[len(body.Messages)-1].Content
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"choices":[{"message":{"content":"估值材料：样本仅用于 smoke，需结合规则引擎与人工确认。"}}]}`)
	}))
	defer server.Close()
	configPath := writeLLMSmokeTestConfig(t, dbPath, server.URL, "test-key", "gpt-5.4-mini")

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "llm-smoke", "--symbol", "510300"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected llm-smoke exit code 0, got %d, stderr=%s", code, stderr.String())
	}
	if seenModel != "gpt-5.4-mini" {
		t.Fatalf("expected configured model gpt-5.4-mini, got %q", seenModel)
	}
	if authHeader != "Bearer test-key" {
		t.Fatalf("expected bearer auth header, got %q", authHeader)
	}
	if strings.Contains(prompt, "test-key") {
		t.Fatalf("prompt must not contain API key: %s", prompt)
	}
	if !strings.Contains(stdout.String(), "不会执行交易") {
		t.Fatalf("expected safety output, got %s", stdout.String())
	}
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	for _, table := range []string{"positions", "portfolio_snapshots", "operation_confirmations", "position_transactions"} {
		assertTableCount(t, store.DB, table, 0)
	}
	assertAuditEvent(t, store.DB, "run_local_task", "llm-smoke:symbol=510300:model=gpt-5.4-mini", "llm_smoke:quality=passed:parse=parsed:no_auto_trading")
}

func TestRunLLMSmokeRequiresRealLLMConfig(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	configPath := writeTestConfig(t, dbPath)

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "llm-smoke", "--symbol", "510300"}, &stdout, &stderr)

	if code == 0 || !strings.Contains(stderr.String(), "deepseek.api_key") {
		t.Fatalf("expected missing deepseek.api_key error, code=%d stderr=%s", code, stderr.String())
	}
}

func TestRunEvidenceIndexWritesRAGChunkAndAudit(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	configPath := writeTestConfig(t, dbPath)

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "evidence-index"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected evidence task exit code 0, got %d, stderr=%s", code, stderr.String())
	}
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	assertTableCount(t, store.DB, "rag_chunks", 1)
	assertAuditEvent(t, store.DB, "rebuild_index", "evidence-index", "no_auto_trading")
}

func TestRunEvidenceIndexRebuildsCorruptedVecLiteFromSQLiteChunks(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	configPath := writeTestConfig(t, dbPath)
	store := seedRetrievalQualitySmokeFacts(t, dbPath)
	store.Close()
	indexPath := filepath.Join(filepath.Dir(dbPath), "veclite.json")
	if err := os.WriteFile(indexPath, []byte("not-json"), 0o600); err != nil {
		t.Fatalf("write corrupted index: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "evidence-index"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected evidence-index exit code 0, got %d, stderr=%s", code, stderr.String())
	}
	health := service.NewFileVectorIndex(indexPath).Health(context.Background())
	if health.Status != service.VectorIndexHealthHealthy || health.ChunkCount == 0 {
		t.Fatalf("expected evidence-index to rebuild a healthy VecLite index, got %+v", health)
	}
	stdout.Reset()
	stderr.Reset()
	code = run(context.Background(), []string{"--config", configPath, "--task", "retrieval-quality-smoke", "--symbol", "510300"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected retrieval-quality-smoke exit code 0, got %d, stderr=%s", code, stderr.String())
	}
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	assertAuditEvent(t, store.DB, "run_local_task", "retrieval-quality-smoke:symbol=510300", "retrieval_quality:status=hit:topk=1:fallback=veclite:index=healthy:consistency=checked:no_auto_trading")
}

func TestRunRetrievalQualitySmokeWritesQualityAuditWithoutTrading(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	configPath := writeTestConfig(t, dbPath)
	store := seedRetrievalQualitySmokeFacts(t, dbPath)
	store.Close()

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "retrieval-quality-smoke", "--symbol", "510300"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected retrieval quality smoke exit code 0, got %d, stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "不会执行交易") {
		t.Fatalf("expected safety output, got %s", stdout.String())
	}
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	for _, table := range []string{"positions", "portfolio_snapshots", "operation_confirmations", "position_transactions", "local_account_import_batches", "local_account_corrections"} {
		assertTableCount(t, store.DB, table, 0)
	}
	assertAuditEvent(t, store.DB, "run_local_task", "retrieval-quality-smoke:symbol=510300", "retrieval_quality:status=degraded:topk=1:fallback=sqlite_summary:index=missing:consistency=checked:no_auto_trading")
}

func TestRunPublicEvidenceRefreshWritesEvidenceAndAudit(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().UTC()
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/new/hisAnnouncement/query":
			_, _ = fmt.Fprintf(w, `{"totalAnnouncement":1,"totalRecordNum":1,"announcements":[{"secCode":"510300","secName":"300ETF","orgId":"org","announcementId":"ann-cninfo","announcementTitle":"ETF 公告","announcementTime":%d,"adjunctUrl":"/a.pdf","adjunctType":"PDF","announcementType":"基金公告"}],"hasMore":false,"totalpages":1}`, now.UnixMilli())
		case "/api/disc/announcement/searchQuery":
			_, _ = fmt.Fprintf(w, `{"recordCount":1,"data":[{"secCode":"510300","secName":"300ETF","announList":[{"id":"szse-id","title":"ETF 公告","attachPath":"/b.pdf","attachFormat":"PDF","attachSize":1,"annId":"ann-szse","bigCategoryName":"定期报告","publishTime":"%s"}]}]}`, now.Format("2006-01-02 15:04:05"))
		case "/searchList":
			_, _ = fmt.Fprintf(w, `{"data":{"page":1,"rows":50,"total":1,"results":[{"title":"基金监管规则","content":"监管正文","url":"/rule.shtml","publishedTime":"%s","channelCodeName":"rules","manuscriptId":"rule-1"}]}}`, now.Format(time.RFC3339))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	configPath := writePublicEvidenceTestConfig(t, dbPath, server.URL)

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "public-evidence-refresh", "--symbol", "510300"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected public evidence task exit code 0, got %d, stderr=%s", code, stderr.String())
	}
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	assertTableCount(t, store.DB, "intelligence_items", 3)
	assertTableCount(t, store.DB, "rag_chunks", 3)
	assertTableAtLeast(t, store.DB, "source_verifications", 2)
	assertAuditEvent(t, store.DB, "run_local_task", "public-evidence-refresh:symbol=510300", "no_auto_trading")
}
func TestRunPublicEvidenceRefreshUsesExplicitDateWindow(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	seenDate := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/new/hisAnnouncement/query" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		seenDate = r.Form.Get("seDate")
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"totalAnnouncement":1,"totalRecordNum":1,"announcements":[{"secCode":"510300","secName":"300ETF","orgId":"org","announcementId":"ann-cninfo","announcementTitle":"ETF 公告","announcementTime":1718841600000,"adjunctUrl":"/a.pdf","adjunctType":"PDF","announcementType":"基金公告"}],"hasMore":false,"totalpages":1}`)
	}))
	defer server.Close()
	configPath := writePublicEvidenceTestConfigWithSources(t, dbPath, server.URL, []string{"cninfo"})

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "public-evidence-refresh", "--symbol", "510300", "--start-date", "2024-06-01", "--end-date", "2024-06-30"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected public evidence task exit code 0, got %d, stderr=%s", code, stderr.String())
	}
	if seenDate != "2024-06-01~2024-06-30" {
		t.Fatalf("expected explicit seDate, got %s", seenDate)
	}
}

func TestRunPublicEvidenceRefreshRejectsInvalidDateWindow(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	configPath := writePublicEvidenceTestConfigWithSources(t, dbPath, "https://example.invalid", []string{"cninfo"})
	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "public-evidence-refresh", "--symbol", "510300", "--start-date", "2024/06/01", "--end-date", "2024-06-30"}, &stdout, &stderr)
	if code == 0 || !strings.Contains(stderr.String(), "YYYY-MM-DD") {
		t.Fatalf("expected invalid date error, code=%d stderr=%s", code, stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	code = run(context.Background(), []string{"--config", configPath, "--task", "public-evidence-refresh", "--symbol", "510300", "--start-date", "2024-07-01", "--end-date", "2024-06-30"}, &stdout, &stderr)
	if code == 0 || !strings.Contains(stderr.String(), "must not be after") {
		t.Fatalf("expected reversed date error, code=%d stderr=%s", code, stderr.String())
	}
}

func TestRunPublicEvidenceRefreshAuditsSymbolAndDateWindow(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"totalAnnouncement":1,"totalRecordNum":1,"announcements":[{"secCode":"159915","secName":"创业板ETF","orgId":"org","announcementId":"ann-window","announcementTitle":"创业板ETF 公告","announcementTime":1718841600000,"adjunctUrl":"/a.pdf","adjunctType":"PDF","announcementType":"基金公告"}],"hasMore":false,"totalpages":1}`)
	}))
	defer server.Close()
	configPath := writePublicEvidenceTestConfigWithSources(t, dbPath, server.URL, []string{"cninfo"})

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "public-evidence-refresh", "--symbol", "159915", "--start-date", "2024-06-01", "--end-date", "2024-06-30"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected public evidence task exit code 0, got %d, stderr=%s", code, stderr.String())
	}
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	var inputRef string
	if err := store.DB.QueryRow(`SELECT input_ref FROM audit_events WHERE node_name='cmd_agent' ORDER BY created_at DESC LIMIT 1`).Scan(&inputRef); err != nil {
		t.Fatalf("read cmd audit input ref: %v", err)
	}
	if inputRef != "public-evidence-refresh:symbol=159915:start=2024-06-01:end=2024-06-30" {
		t.Fatalf("expected symbol/date window audit input, got %s", inputRef)
	}
}

func TestRunNon510300DynamicAcceptanceBindsCollectorSourceHealthAuditAndReadiness(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	seen := struct {
		marketSymbol string
		cninfoStock  string
		cninfoDate   string
		szseKeyword  string
	}{}
	freshAnnouncementMillis := time.Date(2026, 6, 19, 0, 0, 0, 0, time.UTC).UnixMilli()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/market":
			seen.marketSymbol = r.URL.Query().Get("symbol")
			if seen.marketSymbol != "159915" {
				t.Fatalf("market refresh requested wrong symbol: %s", seen.marketSymbol)
			}
			_, _ = fmt.Fprint(w, `{"close_price":2.413,"turnover_rate":1.4,"pe_percentile":42,"pb_percentile":47,"volume_percentile":51,"volatility_percentile":33,"liquidity_state":"normal","sentiment_state":"neutral","source_name":"accepted_local_market","source_level":"B","source_type":"market_price","trade_date":"2026-06-19","captured_at":"2026-06-19T08:00:00Z","metadata":{"p34_source_health":{"symbol_profile":{"freshness":"fresh","data_date":"2026-06-19","affected_symbols":["159915"],"source_level":"A","source_type":"symbol_profile"},"fund_profile":{"freshness":"fresh","data_date":"2026-06-19","affected_symbols":["159915"],"source_level":"B","source_type":"fund_profile"},"tracked_index":{"freshness":"fresh","data_date":"2026-06-19","affected_symbols":["399006"],"source_level":"A","source_type":"index_profile"},"market_price":{"freshness":"fresh","data_date":"2026-06-19","affected_symbols":["159915"],"source_level":"B","source_type":"market_price"},"valuation_percentiles":{"freshness":"fresh","data_date":"2026-06-19","affected_symbols":["399006"],"source_level":"A","source_type":"index_valuation"},"liquidity":{"freshness":"fresh","data_date":"2026-06-19","affected_symbols":["159915"],"source_level":"B","source_type":"liquidity"},"sentiment_proxy":{"freshness":"fresh","data_date":"2026-06-19","affected_symbols":["159915"],"source_level":"C","source_type":"sentiment_proxy"},"rag_index":{"freshness":"fresh","data_date":"2026-06-19","affected_symbols":["159915"],"source_level":"local_index","source_name":"veclite","source_type":"rag_index"}},"p34_data_categories":["symbol_profile","fund_profile","tracked_index","market_price","valuation_percentiles","liquidity","sentiment_proxy","rag_index"]}}`)
		case "/new/hisAnnouncement/query":
			if err := r.ParseForm(); err != nil {
				t.Fatal(err)
			}
			seen.cninfoStock = r.Form.Get("stock")
			seen.cninfoDate = r.Form.Get("seDate")
			if seen.cninfoStock != "159915" {
				t.Fatalf("cninfo requested wrong stock: %s", seen.cninfoStock)
			}
			_, _ = fmt.Fprintf(w, `{"totalAnnouncement":1,"totalRecordNum":1,"announcements":[{"secCode":"159915","secName":"创业板ETF","orgId":"org","announcementId":"ann-159915","announcementTitle":"创业板ETF 公告","announcementTime":%d,"adjunctUrl":"/159915.pdf","adjunctType":"PDF","announcementType":"基金公告"}],"hasMore":false,"totalpages":1}`, freshAnnouncementMillis)
		case "/api/disc/announcement/searchQuery":
			seen.szseKeyword = r.URL.Query().Get("keyword")
			if seen.szseKeyword != "159915" {
				t.Fatalf("szse requested wrong keyword: %s", seen.szseKeyword)
			}
			_, _ = fmt.Fprint(w, `{"recordCount":1,"data":[{"secCode":"159915","secName":"创业板ETF","announList":[{"id":"szse-159915","title":"创业板ETF 公告","attachPath":"/159915-szse.pdf","attachFormat":"PDF","attachSize":1,"annId":"ann-159915-szse","bigCategoryName":"定期报告","publishTime":"2026-06-19 09:00:00"}]}]}`)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	configPath := writeDynamicAcceptanceTestConfig(t, dbPath, server.URL)

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "market-refresh", "--symbol", "159915"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected market refresh exit code 0, got %d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	stdout.Reset()
	stderr.Reset()
	code = run(context.Background(), []string{"--config", configPath, "--task", "public-evidence-refresh", "--symbol", "159915", "--start-date", "2026-06-01", "--end-date", "2026-06-30"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected public evidence refresh exit code 0, got %d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	if seen.marketSymbol != "159915" || seen.cninfoStock != "159915" || seen.cninfoDate != "2026-06-01~2026-06-30" || seen.szseKeyword != "159915" {
		t.Fatalf("collector request correlation mismatch: %+v", seen)
	}

	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	repos := repository.Repositories{
		AuditRepo:        appsqlite.NewAuditRepository(store.DB),
		RuleRepo:         appsqlite.NewRuleRepository(store.DB),
		MarketRepo:       appsqlite.NewMarketRepository(store.DB),
		IntelligenceRepo: appsqlite.NewIntelligenceRepository(store.DB),
	}
	var metrics string
	if err := store.DB.QueryRow(`SELECT market_metrics_json FROM market_snapshots WHERE symbol='159915' ORDER BY created_at DESC LIMIT 1`).Scan(&metrics); err != nil {
		t.Fatalf("read 159915 market metrics: %v", err)
	}
	for _, want := range []string{"p34_source_health", "159915", "399006", "valuation_percentiles", "2026-06-19"} {
		if !strings.Contains(metrics, want) {
			t.Fatalf("expected market source health to contain %q, got %s", want, metrics)
		}
	}
	var metricsPayload map[string]any
	if err := json.Unmarshal([]byte(metrics), &metricsPayload); err != nil {
		t.Fatalf("unmarshal market metrics: %v", err)
	}
	marketRequestID, _ := metricsPayload["request_id"].(string)
	if marketRequestID == "" {
		t.Fatalf("expected market metrics to carry request_id, got %s", metrics)
	}
	metadata, _ := metricsPayload["metadata"].(map[string]any)
	health, _ := metadata["p34_source_health"].(map[string]any)
	trackedHealth, _ := health["tracked_index"].(map[string]any)
	if trackedHealth["request_id"] != marketRequestID || trackedHealth["data_date"] != "2026-06-19" {
		t.Fatalf("expected source health request/data correlation, request=%s health=%+v", marketRequestID, trackedHealth)
	}
	var summaryCount int
	if err := store.DB.QueryRow(`SELECT COUNT(*) FROM intelligence_summary WHERE symbol='159915'`).Scan(&summaryCount); err != nil {
		t.Fatalf("count 159915 summaries: %v", err)
	}
	if summaryCount != 2 {
		t.Fatalf("expected 2 public evidence summaries for 159915, got %d", summaryCount)
	}
	var verificationStatus string
	if err := store.DB.QueryRow(`SELECT verification_status FROM source_verifications WHERE symbol='159915' ORDER BY created_at DESC LIMIT 1`).Scan(&verificationStatus); err != nil {
		t.Fatalf("read 159915 source verification: %v", err)
	}
	if verificationStatus != "satisfied" {
		t.Fatalf("expected 159915 multi-source formal evidence satisfied, got %s", verificationStatus)
	}
	var evidenceMetadata string
	if err := store.DB.QueryRow(`SELECT metadata_json FROM rag_chunks WHERE chunk_id IN (SELECT chunk_id FROM rag_chunks JOIN intelligence_summary USING(summary_id) WHERE intelligence_summary.symbol='159915') ORDER BY created_at DESC LIMIT 1`).Scan(&evidenceMetadata); err != nil {
		t.Fatalf("read 159915 evidence metadata: %v", err)
	}
	var evidencePayload map[string]any
	if err := json.Unmarshal([]byte(evidenceMetadata), &evidencePayload); err != nil {
		t.Fatalf("unmarshal evidence metadata: %v", err)
	}
	evidenceRequestID, _ := evidencePayload["request_id"].(string)
	if evidenceRequestID == "" || evidencePayload["source_name"] == "" {
		t.Fatalf("expected evidence metadata request/source correlation, got %s", evidenceMetadata)
	}
	var auditCount int
	if err := store.DB.QueryRow(`SELECT COUNT(*) FROM audit_events WHERE input_ref='public-evidence-refresh:symbol=159915:start=2026-06-01:end=2026-06-30'`).Scan(&auditCount); err != nil {
		t.Fatalf("count public evidence audit: %v", err)
	}
	if auditCount != 1 {
		t.Fatalf("expected one correlated public evidence audit event, got %d", auditCount)
	}
	var evidenceAuditRequestID string
	if err := store.DB.QueryRow(`SELECT request_id FROM audit_events WHERE input_ref='159915' AND output_ref='source=public_evidence count=2' ORDER BY created_at DESC LIMIT 1`).Scan(&evidenceAuditRequestID); err != nil {
		t.Fatalf("read public evidence ingestion audit request: %v", err)
	}
	if evidenceAuditRequestID == "" || evidenceAuditRequestID != evidenceRequestID {
		t.Fatalf("expected evidence request_id to match ingestion audit, metadata=%s audit=%s", evidenceRequestID, evidenceAuditRequestID)
	}
	readiness, err := service.NewKnowledgeReadinessService(repos).Evaluate(context.Background(), service.KnowledgeReadinessRequest{Symbol: "159915"})
	if err != nil {
		t.Fatalf("Evaluate readiness: %v", err)
	}
	if readiness.Status != "ready" || readiness.SymbolProfile.TrackedIndexSymbol != "399006" {
		t.Fatalf("expected ready readiness bound to 159915/399006, got %+v", readiness)
	}
	for _, category := range []string{"symbol_profile", "tracked_index", "market_price", "valuation_percentiles", "liquidity", "formal_evidence", "rag_index"} {
		dep := readinessDependencyByCategoryForAgentTest(readiness.DataDependencies, category)
		if dep.Status != "ready" {
			t.Fatalf("expected readiness category %s ready, got %+v", category, dep)
		}
		if category != "formal_evidence" && dep.RequestID != marketRequestID {
			t.Fatalf("expected readiness category %s to preserve request_id %s, got %+v", category, marketRequestID, dep)
		}
	}
}

func TestRunPublicEvidenceRefreshRequiresEnabledConfig(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	configPath := writeTestConfig(t, dbPath)
	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "public-evidence-refresh", "--symbol", "510300"}, &stdout, &stderr)
	if code == 0 || !strings.Contains(stderr.String(), "public evidence collector is disabled") {
		t.Fatalf("expected disabled public evidence error, code=%d stderr=%s", code, stderr.String())
	}
}

func TestRunDailyTaskWritesDecisionAuditAndReport(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	configPath := writeTestConfig(t, dbPath)
	store := seedDailyTaskFacts(t, dbPath)
	store.Close()

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "daily"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected daily task exit code 0, got %d, stderr=%s", code, stderr.String())
	}
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	assertTableCount(t, store.DB, "decision_records", 1)
	assertAuditEvent(t, store.DB, "generate_decision", "daily", "no_auto_trading")

	var sourceType, sourceID, decisionID, status, summary string
	if err := store.DB.QueryRow(`SELECT source_type,source_id,COALESCE(decision_id,''),status,summary FROM daily_discipline_reports WHERE source_type='manual'`).Scan(&sourceType, &sourceID, &decisionID, &status, &summary); err != nil {
		t.Fatalf("read manual daily discipline report: %v", err)
	}
	if sourceType != "manual" || sourceID == "" || decisionID == "" || status != "degraded" || summary != "今日纪律报告已生成" {
		t.Fatalf("unexpected manual daily discipline report source_type=%s source_id=%s decision_id=%s status=%s summary=%s", sourceType, sourceID, decisionID, status, summary)
	}
}

func TestRunDailyTaskWritesDecisionAuditReportAndRiskAlerts(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	configPath := writeTestConfig(t, dbPath)
	store := seedDailyTaskFacts(t, dbPath)
	if _, err := store.DB.Exec(`UPDATE market_snapshots SET pe_percentile=88 WHERE market_snapshot_id='ms_daily'`); err != nil {
		t.Fatalf("seed valuation risk: %v", err)
	}
	store.Close()

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "daily"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected daily task exit code 0, got %d, stderr=%s", code, stderr.String())
	}
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	var riskType, sopStatus, reportID string
	if err := store.DB.QueryRow(`SELECT risk_type,sop_status,COALESCE(related_report_id,'') FROM risk_alerts WHERE symbol='510300'`).Scan(&riskType, &sopStatus, &reportID); err != nil {
		t.Fatalf("read risk alert: %v", err)
	}
	if riskType != "valuation_high" || sopStatus != "active" || reportID == "" {
		t.Fatalf("unexpected risk alert risk_type=%s sop=%s report_id=%s", riskType, sopStatus, reportID)
	}
	assertTableCount(t, store.DB, "notifications", 1)
}

func TestManualDailyDisciplineReportLocalDateUsesConfiguredTimezone(t *testing.T) {
	utcBoundary := time.Date(2026, 6, 7, 23, 30, 0, 0, time.UTC)

	if got := manualDailyDisciplineReportLocalDate(utcBoundary, "Pacific/Kiritimati"); got != "2026-06-08" {
		t.Fatalf("expected Kiritimati local date 2026-06-08, got %s", got)
	}
	if got := manualDailyDisciplineReportLocalDate(utcBoundary, "Asia/Shanghai"); got != "2026-06-08" {
		t.Fatalf("expected Shanghai local date 2026-06-08, got %s", got)
	}
}

func TestUpsertManualDailyDisciplineReportWritesConfiguredTimezoneLocalDate(t *testing.T) {
	originalNow := manualDailyDisciplineReportNow
	manualDailyDisciplineReportNow = func() time.Time {
		return time.Date(2026, 6, 7, 23, 30, 0, 0, time.UTC)
	}
	t.Cleanup(func() { manualDailyDisciplineReportNow = originalNow })

	repo := &manualDailyDisciplineReportRepoSpy{}
	err := upsertManualDailyDisciplineReport(context.Background(), repository.Repositories{DailyDisciplineReportRepo: repo}, nil, "Pacific/Kiritimati", "req_manual", workflow.WorkflowContext{})

	if err != nil {
		t.Fatalf("upsert manual daily discipline report: %v", err)
	}
	if repo.report.LocalDate != "2026-06-08" {
		t.Fatalf("expected report.LocalDate from configured timezone, got %s", repo.report.LocalDate)
	}
}

func TestManualDailyDisciplineSymbolSetHashMatchesAutoRunEightByteHexScheme(t *testing.T) {
	got := manualDailyDisciplineSymbolSetHash([]model.Position{{Symbol: " 510500 "}, {Symbol: "510300"}})
	want := "9b805b68bb48e135"
	if got != want {
		t.Fatalf("expected auto-run compatible 8-byte symbol hash %s, got %s", want, got)
	}
}

func TestRunTaskFailureWritesFailedAuditWithErrorCode(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	configPath := writeTestConfig(t, dbPath)

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--task", "daily"}, &stdout, &stderr)

	if code == 0 {
		t.Fatalf("expected daily task to fail without required local facts")
	}
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	var status, errorCode, outputRef string
	if err := store.DB.QueryRow(`SELECT status,error_code,output_ref FROM audit_events ORDER BY created_at DESC LIMIT 1`).Scan(&status, &errorCode, &outputRef); err != nil {
		t.Fatalf("read failed audit: %v", err)
	}
	if status != "failed" || errorCode == "" {
		t.Fatalf("expected failed audit with error code, status=%s error=%s", status, errorCode)
	}
	if outputRef != "task_failed" {
		t.Fatalf("expected non data-source task failure output_ref=task_failed, got %s", outputRef)
	}
}

func TestRunValidateConfigReportsDiagnostics(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	configPath := writeTestConfig(t, dbPath)

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--validate-config"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected validate config exit code 0, got %d stderr=%s", code, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{"config validation passed", "sqlite", "veclite", "不会执行交易"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected validate output to contain %q, got %s", want, out)
		}
	}
}

func TestRunPreflightReportsChecksAndWritesSafeDiagnostics(t *testing.T) {
	dir := t.TempDir()
	missingDataDir := filepath.Join(dir, "missing-data")
	dbPath := filepath.Join(missingDataDir, "agent.db")
	configPath := writeLLMSmokeTestConfig(t, dbPath, "https://api.deepseek.com", "sk-test-secret", "test-model")
	diagnosticsPath := filepath.Join(dir, "diagnostics", "preflight.json")

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--preflight", "--diagnostics", diagnosticsPath}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected preflight exit code 0, got %d stderr=%s", code, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{"preflight generated_at", "config_validation:pass", "go_binary:", "playwright_browser:", "sqlite_path:warning", "veclite_path:warning", "不会执行交易", "不会自动应用规则"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected preflight output to contain %q, got %s", want, out)
		}
	}
	if _, err := os.Stat(missingDataDir); !os.IsNotExist(err) {
		t.Fatalf("preflight should not create configured data directory, stat err=%v", err)
	}
	if strings.Contains(out, "sk-test-secret") {
		t.Fatalf("preflight output leaked api key: %s", out)
	}
	data, err := os.ReadFile(diagnosticsPath)
	if err != nil {
		t.Fatalf("read diagnostics: %v", err)
	}
	if strings.Contains(string(data), "sk-test-secret") {
		t.Fatalf("diagnostics leaked api key: %s", string(data))
	}
	var report preflightReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("decode diagnostics: %v", err)
	}
	if len(report.Checks) == 0 || report.SafetyNote == "" {
		t.Fatalf("expected diagnostics checks and safety note, got %+v", report)
	}
}

func TestRunReleaseUpgradeCheckWritesSanitizedReadOnlyDiagnostics(t *testing.T) {
	dir := t.TempDir()
	missingDataDir := filepath.Join(dir, "missing-data")
	dbPath := filepath.Join(missingDataDir, "agent.db")
	configPath := writeLLMSmokeTestConfig(t, dbPath, "https://api.deepseek.com", "sk-test-secret", "test-model")
	diagnosticsPath := filepath.Join(dir, "diagnostics", "release-upgrade.json")

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--release-upgrade-check", "--target-version", "test-p49", "--diagnostics", diagnosticsPath}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected release upgrade check exit code 0, got %d stderr=%s", code, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{"release upgrade generated_at", "current=local-dev", "target=test-p49", "status=warning", "version_check:pass", "backup_reminder:warning", "migration_precheck:pass", "post-upgrade smoke commands", "不会执行升级", "不会运行迁移", "不会创建备份", "不会恢复或覆盖数据库", "不会执行交易", "不会外部推送", "不会自动确认", "不会自动应用规则", "不会自动修复"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected release upgrade output to contain %q, got %s", want, out)
		}
	}
	for _, leaked := range []string{"sk-test-secret", dbPath, configPath, missingDataDir} {
		if strings.Contains(out, leaked) {
			t.Fatalf("release upgrade stdout leaked %q: %s", leaked, out)
		}
	}
	if _, err := os.Stat(missingDataDir); !os.IsNotExist(err) {
		t.Fatalf("release upgrade check should not create configured data directory, stat err=%v", err)
	}
	data, err := os.ReadFile(diagnosticsPath)
	if err != nil {
		t.Fatalf("read diagnostics: %v", err)
	}
	for _, leaked := range []string{"sk-test-secret", dbPath, configPath, missingDataDir} {
		if strings.Contains(string(data), leaked) {
			t.Fatalf("release upgrade diagnostics leaked %q: %s", leaked, string(data))
		}
	}
	var report releaseUpgradeReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("decode release upgrade diagnostics: %v", err)
	}
	if report.TargetVersion != "test-p49" || report.Status != "warning" || len(report.Checks) == 0 || len(report.PostUpgradeSmokeCommands) == 0 {
		t.Fatalf("unexpected release upgrade report: %+v", report)
	}
}

func TestRunReleaseUpgradeCheckWarnsWithoutTargetVersionAndDoesNotWriteAudit(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "agent.db")
	configPath := writeTestConfig(t, dbPath)
	store := seedDailyTaskFacts(t, dbPath)
	store.Close()

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--release-upgrade-check"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected release upgrade check warning exit code 0, got %d stderr=%s", code, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{"target=<missing>", "status=warning", "version_check:warning", "--target-version", "backup_reminder:warning"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected missing target output to contain %q, got %s", want, out)
		}
	}
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	assertTableCount(t, store.DB, "audit_events", 0)
}

func TestRunReleaseUpgradeCheckRedactsUnsafeTargetVersion(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "agent.db")
	configPath := writeLLMSmokeTestConfig(t, dbPath, "https://api.deepseek.com", "sk-test-secret", "test-model")
	diagnosticsPath := filepath.Join(dir, "diagnostics", "release-upgrade.json")
	unsafeTarget := "/Users/private/sk-123456789012/raw HTTP prompt:secret"

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--release-upgrade-check", "--target-version", unsafeTarget, "--diagnostics", diagnosticsPath}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected release upgrade check exit code 0, got %d stderr=%s", code, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{"target=<redacted-target-version>", "version_check:warning", "target version redacted"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected unsafe target output to contain %q, got %s", want, out)
		}
	}
	data, err := os.ReadFile(diagnosticsPath)
	if err != nil {
		t.Fatalf("read diagnostics: %v", err)
	}
	for _, leaked := range []string{unsafeTarget, "/Users/private", "sk-123456789012", "raw HTTP", "prompt:secret", "sk-test-secret", dbPath, configPath} {
		if strings.Contains(out, leaked) {
			t.Fatalf("release upgrade stdout leaked %q: %s", leaked, out)
		}
		if strings.Contains(string(data), leaked) {
			t.Fatalf("release upgrade diagnostics leaked %q: %s", leaked, string(data))
		}
	}
	var report releaseUpgradeReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("decode release upgrade diagnostics: %v", err)
	}
	if report.TargetVersion != "<redacted-target-version>" || report.Status != "warning" {
		t.Fatalf("expected redacted warning report, got %+v", report)
	}
}

func TestRunReleaseUpgradeCheckRedactsSQLiteStatErrors(t *testing.T) {
	dir := t.TempDir()
	parentFile := filepath.Join(dir, "not-a-dir")
	if err := os.WriteFile(parentFile, []byte("not a directory"), 0o600); err != nil {
		t.Fatalf("write parent file: %v", err)
	}
	dbPath := filepath.Join(parentFile, "agent.db")
	configPath := writeTestConfig(t, dbPath)
	diagnosticsPath := filepath.Join(dir, "diagnostics", "release-upgrade.json")

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--release-upgrade-check", "--target-version", "test-p49", "--diagnostics", diagnosticsPath}, &stdout, &stderr)

	if code == 0 {
		t.Fatalf("expected release upgrade check to block on stat error")
	}
	out := stdout.String()
	if !strings.Contains(out, "backup_reminder:failed:sqlite path stat failed") {
		t.Fatalf("expected sanitized stat failure, got %s", out)
	}
	data, err := os.ReadFile(diagnosticsPath)
	if err != nil {
		t.Fatalf("read diagnostics: %v", err)
	}
	for _, leaked := range []string{dbPath, parentFile, dir} {
		if strings.Contains(out, leaked) {
			t.Fatalf("release upgrade stdout leaked %q: %s", leaked, out)
		}
		if strings.Contains(string(data), leaked) {
			t.Fatalf("release upgrade diagnostics leaked %q: %s", leaked, string(data))
		}
	}
}

func TestPlaywrightCheckReportsSkippedWarningAndExecutableBrowser(t *testing.T) {
	check := playwrightCheckWithPaths("", false, nil)
	if check.Status != "skipped" {
		t.Fatalf("expected missing cli to be skipped, got %+v", check)
	}

	dir := t.TempDir()
	check = playwrightCheckWithPaths("/tmp/playwright", true, []string{dir})
	if check.Status != "warning" {
		t.Fatalf("expected missing browser executable to warn, got %+v", check)
	}

	browserExecutable := filepath.Join(dir, "chromium-1234", "chrome-mac", "Chromium.app", "Contents", "MacOS", "Chromium")
	if err := os.MkdirAll(filepath.Dir(browserExecutable), 0o700); err != nil {
		t.Fatalf("create browser dir: %v", err)
	}
	if err := os.WriteFile(browserExecutable, []byte("#!/bin/sh\n"), 0o700); err != nil {
		t.Fatalf("write browser executable: %v", err)
	}
	check = playwrightCheckWithPaths("/tmp/playwright", true, []string{dir})
	if check.Status != "pass" || !strings.Contains(check.Detail, browserExecutable) {
		t.Fatalf("expected executable browser to pass, got %+v", check)
	}
}

func TestRunBackupRefusesMissingSQLiteSource(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "missing.db")
	configPath := writeTestConfig(t, dbPath)

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--config", configPath, "--backup", filepath.Join(dir, "backups")}, &stdout, &stderr)

	if code == 0 || !strings.Contains(stderr.String(), "sqlite.path does not exist") {
		t.Fatalf("expected missing sqlite.path backup failure, code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		t.Fatalf("backup should not create missing sqlite database, stat err=%v", err)
	}
}

func TestRunRecoverySmokeRestoresBackupAndVerifiesFacts(t *testing.T) {
	dir := t.TempDir()
	sourceDBPath := filepath.Join(dir, "source.db")
	sourceConfigPath := writeTestConfig(t, sourceDBPath)
	store := seedDailyTaskFacts(t, sourceDBPath)
	store.Close()
	backupDir := filepath.Join(dir, "backups")

	var backupOut, backupErr bytes.Buffer
	backupCode := run(context.Background(), []string{"--config", sourceConfigPath, "--backup", backupDir}, &backupOut, &backupErr)
	if backupCode != 0 {
		t.Fatalf("expected backup exit code 0, got %d stderr=%s", backupCode, backupErr.String())
	}
	backupFile := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(backupOut.String()), "backup created:"))

	restoreDBPath := filepath.Join(dir, "restore", "agent.db")
	restoreConfigPath := writeTestConfig(t, restoreDBPath)
	var smokeOut, smokeErr bytes.Buffer
	smokeCode := run(context.Background(), []string{"--config", restoreConfigPath, "--recovery-smoke", backupFile}, &smokeOut, &smokeErr)
	if smokeCode != 0 {
		t.Fatalf("expected recovery smoke exit code 0, got %d stderr=%s", smokeCode, smokeErr.String())
	}
	out := smokeOut.String()
	for _, want := range []string{"recovery smoke completed", "portfolio=1", "intelligence=1", "不会执行交易", "不会外部推送", "不会自动应用规则"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected recovery smoke output to contain %q, got %s", want, out)
		}
	}

	smokeOut.Reset()
	smokeErr.Reset()
	smokeCode = run(context.Background(), []string{"--config", restoreConfigPath, "--recovery-smoke", backupFile}, &smokeOut, &smokeErr)
	if smokeCode == 0 || !strings.Contains(smokeErr.String(), "refuse to overwrite existing sqlite") {
		t.Fatalf("expected recovery smoke to refuse existing db, code=%d stderr=%s", smokeCode, smokeErr.String())
	}
}

func TestRunBackupAndRestoreLocalFilesSafely(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "agent.db")
	configPath := writeTestConfig(t, dbPath)
	store := seedDailyTaskFacts(t, dbPath)
	store.Close()
	backupDir := filepath.Join(dir, "backups")

	var backupOut, backupErr bytes.Buffer
	backupCode := run(context.Background(), []string{"--config", configPath, "--backup", backupDir}, &backupOut, &backupErr)
	if backupCode != 0 {
		t.Fatalf("expected backup exit code 0, got %d stderr=%s", backupCode, backupErr.String())
	}
	backupPath := strings.TrimSpace(backupOut.String())
	if !strings.Contains(backupPath, "backup created:") {
		t.Fatalf("expected backup path output, got %s", backupOut.String())
	}
	backupFile := strings.TrimSpace(strings.TrimPrefix(backupPath, "backup created:"))
	if _, err := os.Stat(backupFile); err != nil {
		t.Fatalf("expected backup file: %v", err)
	}

	var restoreOut, restoreErr bytes.Buffer
	restoreCode := run(context.Background(), []string{"--config", configPath, "--restore", backupFile}, &restoreOut, &restoreErr)
	if restoreCode == 0 || !strings.Contains(restoreErr.String(), "--restore-confirm") {
		t.Fatalf("expected restore without confirm to refuse overwrite, code=%d stderr=%s", restoreCode, restoreErr.String())
	}

	if err := os.WriteFile(dbPath, []byte("existing"), 0o600); err != nil {
		t.Fatalf("write existing db: %v", err)
	}
	restoreOut.Reset()
	restoreErr.Reset()
	restoreCode = run(context.Background(), []string{"--config", configPath, "--restore", backupFile, "--restore-confirm"}, &restoreOut, &restoreErr)
	if restoreCode == 0 || !strings.Contains(restoreErr.String(), "refuse to overwrite existing sqlite") {
		t.Fatalf("expected confirmed restore to refuse existing db, code=%d stderr=%s", restoreCode, restoreErr.String())
	}

	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("remove db before restore: %v", err)
	}
	symlinkTarget := filepath.Join(dir, "unrelated.txt")
	if err := os.WriteFile(symlinkTarget, []byte("keep"), 0o600); err != nil {
		t.Fatalf("write symlink target: %v", err)
	}
	if err := os.Symlink(symlinkTarget, dbPath+".restore.tmp"); err != nil {
		t.Fatalf("create legacy restore symlink: %v", err)
	}
	restoreOut.Reset()
	restoreErr.Reset()
	restoreCode = run(context.Background(), []string{"--config", configPath, "--restore", backupFile, "--restore-confirm"}, &restoreOut, &restoreErr)
	if restoreCode != 0 {
		t.Fatalf("expected confirmed restore exit code 0, got %d stderr=%s", restoreCode, restoreErr.String())
	}
	contents, err := os.ReadFile(symlinkTarget)
	if err != nil {
		t.Fatalf("read symlink target: %v", err)
	}
	if string(contents) != "keep" {
		t.Fatalf("restore should not follow predictable tmp symlink, got %q", string(contents))
	}
	if !strings.Contains(restoreOut.String(), "restore completed") || !strings.Contains(restoreOut.String(), "不会执行交易") {
		t.Fatalf("expected safe restore output, got %s", restoreOut.String())
	}
}

type manualDailyDisciplineReportRepoSpy struct {
	report repository.DailyDisciplineReport
}

func (r *manualDailyDisciplineReportRepoSpy) UpsertDailyDisciplineReport(_ context.Context, report repository.DailyDisciplineReport) error {
	r.report = report
	return nil
}

func (r *manualDailyDisciplineReportRepoSpy) GetDailyDisciplineReport(_ context.Context, _ string) (repository.DailyDisciplineReport, error) {
	return repository.DailyDisciplineReport{}, fmt.Errorf("not implemented")
}

func (r *manualDailyDisciplineReportRepoSpy) GetDailyDisciplineReportByKey(_ context.Context, _, _, _ string) (repository.DailyDisciplineReport, error) {
	return repository.DailyDisciplineReport{}, fmt.Errorf("not implemented")
}

func (r *manualDailyDisciplineReportRepoSpy) ListDailyDisciplineReports(_ context.Context, _ repository.DailyDisciplineReportListFilter) ([]repository.DailyDisciplineReport, error) {
	return nil, fmt.Errorf("not implemented")
}

func seedDailyTaskFacts(t *testing.T, dbPath string) *appsqlite.Store {
	t.Helper()
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := appsqlite.Migrate(context.Background(), store.DB); err != nil {
		_ = store.Close()
		t.Fatalf("migrate sqlite: %v", err)
	}
	_, err = store.DB.Exec(`
INSERT INTO portfolio_snapshots (snapshot_id,snapshot_time,cash,total_assets,cash_ratio,high_risk_ratio,position_count,source,created_at) VALUES ('ps_daily','2026-01-01T00:00:00Z',20000,100000,0.2,0.1,1,'manual','2026-01-01T00:00:00Z');
INSERT INTO positions (position_id,symbol,name,quantity,cost_price,current_price,market_value,unrealized_profit_ratio,position_state,asset_tag,updated_at) VALUES ('pos_daily','510300','沪深300ETF',100,4,4.2,420,0.05,'normal','core','2026-01-01T00:00:00Z');
INSERT INTO position_snapshots (position_snapshot_id,snapshot_id,position_id,symbol,name,quantity,cost_price,current_price,market_value,unrealized_profit_ratio,position_state,asset_tag,created_at) VALUES ('psp_daily','ps_daily','pos_daily','510300','沪深300ETF',100,4,4.2,420,0.05,'normal','core','2026-01-01T00:00:00Z');
INSERT INTO market_snapshots (market_snapshot_id,symbol,trade_date,market_metrics_json,close_price,turnover_rate,pe_percentile,pb_percentile,liquidity_state,sentiment_state,volume_percentile,volatility_percentile,created_at) VALUES ('ms_daily','510300','2026-01-01','{}',4.2,1,50,50,'normal','neutral',50,50,'2026-01-01T00:00:00Z');
INSERT INTO source_verifications (verification_id,verification_group_id,event_id,symbol,event_type,evidence_role,verification_status,independent_source_count,high_grade_independent_source_count,highest_source_level,evidence_ids_json,created_at) VALUES ('sv_daily','vg_daily','ev_daily','510300','normal','formal','satisfied',2,1,'A','["sum_daily"]','2026-01-01T00:00:00Z');
INSERT INTO intelligence_items (intelligence_id,source_name,source_level,original_url,published_at,captured_at,content_hash,raw_title,raw_text_ref,created_at) VALUES ('intel_daily','official','A','https://example.com','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z','hash_daily','每日证据','raw_daily','2026-01-01T00:00:00Z');
INSERT INTO intelligence_summary (summary_id,intelligence_id,symbol,entity,event_type,impact_direction,summary,source_level,evidence_role,time_weight,relevance_score,verification_group_id,created_at) VALUES ('sum_daily','intel_daily','510300','510300','normal','neutral','正式证据摘要','A','formal',1,1,'vg_daily','2026-01-01T00:00:00Z');
INSERT INTO rag_chunks (chunk_id,summary_id,chunk_text,chunk_hash,index_status,metadata_json,created_at) VALUES ('chunk_daily','sum_daily','正式证据摘要','chunk_hash_daily','indexed','{}','2026-01-01T00:00:00Z')`)
	if err != nil {
		_ = store.Close()
		t.Fatalf("seed daily facts: %v", err)
	}
	return store
}

func seedRetrievalQualitySmokeFacts(t *testing.T, dbPath string) *appsqlite.Store {
	t.Helper()
	store, err := appsqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := appsqlite.Migrate(context.Background(), store.DB); err != nil {
		_ = store.Close()
		t.Fatalf("migrate sqlite: %v", err)
	}
	_, err = store.DB.Exec(`
INSERT INTO source_verifications (verification_id,verification_group_id,event_id,symbol,event_type,evidence_role,verification_status,independent_source_count,high_grade_independent_source_count,highest_source_level,evidence_ids_json,created_at) VALUES ('sv_retrieval_smoke','vg_retrieval_smoke','ev_retrieval_smoke','510300','announcement','formal','satisfied',2,1,'A','["sum_retrieval_smoke"]','2026-01-01T00:00:00Z');
INSERT INTO intelligence_items (intelligence_id,source_name,source_level,original_url,published_at,captured_at,content_hash,raw_title,raw_text_ref,created_at) VALUES ('intel_retrieval_smoke','official','A','https://example.invalid/retrieval-smoke','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z','hash_retrieval_smoke','检索质量 smoke 证据','raw_retrieval_smoke','2026-01-01T00:00:00Z');
INSERT INTO intelligence_summary (summary_id,intelligence_id,symbol,entity,event_type,impact_direction,summary,source_level,evidence_role,time_weight,relevance_score,verification_group_id,created_at) VALUES ('sum_retrieval_smoke','intel_retrieval_smoke','510300','510300','announcement','neutral','P38 retrieval quality smoke formal evidence','A','formal',1,1,'vg_retrieval_smoke','2026-01-01T00:00:00Z');
INSERT INTO rag_chunks (chunk_id,summary_id,chunk_text,chunk_hash,index_status,metadata_json,created_at) VALUES ('chunk_retrieval_smoke','sum_retrieval_smoke','P38 retrieval quality smoke formal evidence','chunk_hash_retrieval_smoke','indexed','{"source_level":"A","evidence_role":"formal"}','2026-01-01T00:00:00Z')`)
	if err != nil {
		_ = store.Close()
		t.Fatalf("seed retrieval quality smoke facts: %v", err)
	}
	return store
}

func writeTestConfig(t *testing.T, dbPath string) string {
	t.Helper()
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(configPath, []byte("server:\n  host: 127.0.0.1\n  port: 0\nsqlite:\n  path: "+dbPath+"\nveclite:\n  path: "+filepath.Join(filepath.Dir(dbPath), "veclite.json")+"\ndeepseek:\n  api_key: \"\"\n  base_url: https://api.deepseek.com\ndata_sources:\n  enabled:\n    - stub\n  use_stub: true\nlog:\n  level: error\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return configPath
}

func writeLLMSmokeTestConfig(t *testing.T, dbPath, baseURL, apiKey, llmModel string) string {
	t.Helper()
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	content := "server:\n  host: 127.0.0.1\n  port: 0\nsqlite:\n  path: " + dbPath + "\nveclite:\n  path: " + filepath.Join(filepath.Dir(dbPath), "veclite.json") + "\ndeepseek:\n  api_key: \"" + apiKey + "\"\n  base_url: " + baseURL + "\n  model: " + llmModel + "\ndata_sources:\n  enabled:\n    - stub\n  use_stub: true\nlog:\n  level: error\n"
	if err := os.WriteFile(configPath, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return configPath
}

func writePublicEvidenceTestConfig(t *testing.T, dbPath, baseURL string) string {
	t.Helper()
	return writePublicEvidenceTestConfigWithSources(t, dbPath, baseURL, []string{"cninfo", "szse", "csrc"})
}

func writePublicEvidenceTestConfigWithSources(t *testing.T, dbPath, baseURL string, sources []string) string {
	t.Helper()
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	enabled := make([]string, 0, len(sources))
	configuredSources := make([]string, 0, len(sources))
	for _, source := range sources {
		enabled = append(enabled, "    - "+source)
		configuredSources = append(configuredSources, "      - "+source)
	}
	content := "server:\n  host: 127.0.0.1\n  port: 0\nsqlite:\n  path: " + dbPath + "\nveclite:\n  path: " + filepath.Join(filepath.Dir(dbPath), "veclite.json") + "\ndeepseek:\n  api_key: \"\"\n  base_url: https://api.deepseek.com\ndata_sources:\n  enabled:\n" + strings.Join(enabled, "\n") + "\n  use_stub: true\n  public_evidence:\n    enabled: true\n    sources:\n" + strings.Join(configuredSources, "\n") + "\n    cninfo_base_url: " + baseURL + "\n    szse_base_url: " + baseURL + "\n    csrc_base_url: " + baseURL + "\nlog:\n  level: error\n"
	if err := os.WriteFile(configPath, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return configPath
}

func writeDynamicAcceptanceTestConfig(t *testing.T, dbPath, baseURL string) string {
	t.Helper()
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	content := "server:\n  host: 127.0.0.1\n  port: 0\nsqlite:\n  path: " + dbPath + "\nveclite:\n  path: " + filepath.Join(filepath.Dir(dbPath), "veclite.json") + "\ndeepseek:\n  api_key: \"\"\n  base_url: https://api.deepseek.com\ndata_sources:\n  enabled:\n    - configured\n  use_stub: false\n  market_endpoint: " + baseURL + "/market\n  public_evidence:\n    enabled: true\n    sources:\n      - cninfo\n      - szse\n    cninfo_base_url: " + baseURL + "\n    szse_base_url: " + baseURL + "\nlog:\n  level: error\n"
	if err := os.WriteFile(configPath, []byte(content), 0o600); err != nil {
		t.Fatalf("write dynamic acceptance config: %v", err)
	}
	return configPath
}

func assertTableCount(t *testing.T, db *sql.DB, table string, want int) {
	t.Helper()
	var got int
	if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&got); err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	if got != want {
		t.Fatalf("expected %s count %d, got %d", table, want, got)
	}
}

func readinessDependencyByCategoryForAgentTest(items []service.KnowledgeDataDependency, category string) service.KnowledgeDataDependency {
	for _, item := range items {
		if item.Category == category {
			return item
		}
	}
	return service.KnowledgeDataDependency{}
}

func assertTableAtLeast(t *testing.T, db *sql.DB, table string, want int) {
	t.Helper()
	var got int
	if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&got); err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	if got < want {
		t.Fatalf("expected %s count at least %d, got %d", table, want, got)
	}
}

func assertAuditEvent(t *testing.T, db *sql.DB, action, inputRef, outputRef string) {
	t.Helper()
	var gotAction, gotInputRef, gotOutputRef string
	if err := db.QueryRow(`SELECT action,input_ref,output_ref FROM audit_events ORDER BY created_at DESC LIMIT 1`).Scan(&gotAction, &gotInputRef, &gotOutputRef); err != nil {
		t.Fatalf("read audit event: %v", err)
	}
	if gotAction != action || gotInputRef != inputRef || gotOutputRef != outputRef {
		t.Fatalf("unexpected audit event action=%s input=%s output=%s", gotAction, gotInputRef, gotOutputRef)
	}
}

func assertAuditEventContains(t *testing.T, db *sql.DB, action, inputRef string, outputParts []string) {
	t.Helper()
	var gotAction, gotInputRef, gotOutputRef string
	if err := db.QueryRow(`SELECT action,input_ref,output_ref FROM audit_events ORDER BY created_at DESC LIMIT 1`).Scan(&gotAction, &gotInputRef, &gotOutputRef); err != nil {
		t.Fatalf("read audit event: %v", err)
	}
	if gotAction != action || gotInputRef != inputRef {
		t.Fatalf("unexpected audit event action=%s input=%s output=%s", gotAction, gotInputRef, gotOutputRef)
	}
	for _, part := range outputParts {
		if !strings.Contains(gotOutputRef, part) {
			t.Fatalf("expected audit output to contain %q, got %s", part, gotOutputRef)
		}
	}
}
