package handler

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"investment-agent/internal/application/service"
	"investment-agent/internal/application/workflow"
	"investment-agent/internal/domain/repository"
	appsqlite "investment-agent/internal/infrastructure/persistence/sqlite"
)

func TestCreateConfirmationRejectsInvalidExecutedManuallyFields(t *testing.T) {
	for _, item := range []struct {
		name string
		body string
	}{
		{name: "invalid operation", body: `{"confirmation_type":"executed_manually","operation_type":"hold","symbol":"510300","quantity":1,"price":2,"executed_at":"2026-01-01T00:00:00Z"}`},
		{name: "future executed at", body: `{"confirmation_type":"executed_manually","operation_type":"buy","symbol":"510300","quantity":1,"price":2,"executed_at":"2999-01-01T00:00:00Z"}`},
	} {
		t.Run(item.name, func(t *testing.T) {
			app, db := testApp(t)
			seedDecision(t, db, "decision_invalid_exec", "pending", "formal_trade_advice")
			req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/decision_invalid_exec/confirmations", bytes.NewBufferString(item.body))
			req.SetPathValue("decision_id", "decision_invalid_exec")
			req.Header.Set("X-Request-ID", "req_invalid_exec")
			w := httptest.NewRecorder()

			app.CreateConfirmation(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
			}
			assertCount(t, db, "operation_confirmations", 0)
			assertCount(t, db, "position_transactions", 0)
			assertCount(t, db, "audit_events", 0)
		})
	}
}

func TestCreateConfirmationRejectsInvalidMarkedErrorFields(t *testing.T) {
	for _, item := range []struct {
		name string
		body string
	}{
		{name: "missing actual outcome", body: `{"confirmation_type":"marked_error","root_cause_tag":"analyst_error","lesson_learned":"复盘"}`},
		{name: "invalid root cause", body: `{"confirmation_type":"marked_error","actual_outcome":"missed","root_cause_tag":"bad_tag","lesson_learned":"复盘"}`},
		{name: "missing lesson", body: `{"confirmation_type":"marked_error","actual_outcome":"missed","root_cause_tag":"analyst_error"}`},
	} {
		t.Run(item.name, func(t *testing.T) {
			app, db := testApp(t)
			seedDecision(t, db, "decision_invalid_err", "pending", "formal_trade_advice")
			req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/decision_invalid_err/confirmations", bytes.NewBufferString(item.body))
			req.SetPathValue("decision_id", "decision_invalid_err")
			req.Header.Set("X-Request-ID", "req_invalid_err")
			w := httptest.NewRecorder()

			app.CreateConfirmation(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
			}
			assertCount(t, db, "operation_confirmations", 0)
			assertCount(t, db, "error_cases", 0)
			assertCount(t, db, "audit_events", 0)
		})
	}
}

func TestCreateConfirmationAcceptsContractRootCauseTags(t *testing.T) {
	for _, tag := range []string{"evidence_missed", "rule_threshold_issue", "analyst_error", "user_context_missing", "market_exception"} {
		t.Run(tag, func(t *testing.T) {
			app, db := testApp(t)
			seedDecision(t, db, "decision_tag_"+tag, "pending", "formal_trade_advice")
			req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/decision_tag/confirmations", bytes.NewBufferString(`{"confirmation_type":"marked_error","actual_outcome":"missed","root_cause_tag":"`+tag+`","lesson_learned":"复盘"}`))
			req.SetPathValue("decision_id", "decision_tag_"+tag)
			req.Header.Set("X-Request-ID", "req_tag_"+tag)
			w := httptest.NewRecorder()

			app.CreateConfirmation(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestCreateConfirmationMarkedErrorWritesErrorFacts(t *testing.T) {
	app, db := testApp(t)
	seedDecision(t, db, "decision_err", "pending", "formal_trade_advice")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/decision_err/confirmations", bytes.NewBufferString(`{"confirmation_type":"marked_error","actual_outcome":"missed","root_cause_tag":"analyst_error","lesson_learned":"复盘"}`))
	req.SetPathValue("decision_id", "decision_err")
	req.Header.Set("X-Request-ID", "req_err")
	w := httptest.NewRecorder()

	app.CreateConfirmation(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	assertCount(t, db, "operation_confirmations", 1)
	assertCount(t, db, "error_cases", 1)
	assertCount(t, db, "audit_events", 1)
	var actualOutcome, rootCauseTag, lessonLearned string
	if err := db.QueryRow(`SELECT actual_outcome,root_cause_tag,lesson_learned FROM error_cases WHERE decision_id='decision_err'`).Scan(&actualOutcome, &rootCauseTag, &lessonLearned); err != nil {
		t.Fatalf("read error case: %v", err)
	}
	if actualOutcome != "missed" || rootCauseTag != "analyst_error" || lessonLearned != "复盘" {
		t.Fatalf("unexpected error case fields: outcome=%q tag=%q lesson=%q", actualOutcome, rootCauseTag, lessonLearned)
	}
}

func TestCreateConfirmationRejectsNotRequiredDecision(t *testing.T) {
	app, db := testApp(t)
	seedDecision(t, db, "decision_nr", "not_required", "formal_trade_advice")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/decision_nr/confirmations", bytes.NewBufferString(`{"confirmation_type":"planned"}`))
	req.SetPathValue("decision_id", "decision_nr")
	req.Header.Set("X-Request-ID", "req_nr")
	w := httptest.NewRecorder()

	app.CreateConfirmation(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
	assertCount(t, db, "operation_confirmations", 0)
}

func TestCreateConfirmationRejectsBuyWhenCashIsInsufficient(t *testing.T) {
	app, db := testApp(t)
	seedDecision(t, db, "decision_cash", "pending", "formal_trade_advice")
	_, err := db.Exec(`INSERT INTO portfolio_snapshots (snapshot_id,snapshot_time,cash,total_assets,cash_ratio,high_risk_ratio,position_count,source,created_at) VALUES (?,?,?,?,?,?,?,?,?)`, "snap_cash", "2026-01-01T00:00:00Z", 10, 10, 1, 0, 0, "manual", "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed snapshot: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/decision_cash/confirmations", bytes.NewBufferString(`{"confirmation_type":"executed_manually","operation_type":"buy","symbol":"510300","quantity":10,"price":2,"executed_at":"2026-01-01T00:00:00Z"}`))
	req.SetPathValue("decision_id", "decision_cash")
	req.Header.Set("X-Request-ID", "req_cash")
	w := httptest.NewRecorder()

	app.CreateConfirmation(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
	assertCount(t, db, "operation_confirmations", 0)
	assertCount(t, db, "position_transactions", 0)
	assertCount(t, db, "positions", 0)
	assertCount(t, db, "audit_events", 0)
	assertCount(t, db, "portfolio_snapshots", 1)
}

func TestCreateConfirmationRejectsNonTradeRecord(t *testing.T) {
	app, db := testApp(t)
	seedDecision(t, db, "decision_non_trade", "pending", "non_trade_record")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/decision_non_trade/confirmations", bytes.NewBufferString(`{"confirmation_type":"planned"}`))
	req.SetPathValue("decision_id", "decision_non_trade")
	req.Header.Set("X-Request-ID", "req_non_trade")
	w := httptest.NewRecorder()

	app.CreateConfirmation(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
	assertCount(t, db, "operation_confirmations", 0)
}

func TestCreateConfirmationPlannedCanConvertToWatch(t *testing.T) {
	app, db := testApp(t)
	seedDecision(t, db, "decision_convert", "pending", "formal_trade_advice")

	for _, item := range []struct{ requestID, body string }{
		{"req_convert_plan", `{"confirmation_type":"planned"}`},
		{"req_convert_watch", `{"confirmation_type":"watch"}`},
	} {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/decision_convert/confirmations", bytes.NewBufferString(item.body))
		req.SetPathValue("decision_id", "decision_convert")
		req.Header.Set("X-Request-ID", item.requestID)
		w := httptest.NewRecorder()
		app.CreateConfirmation(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("%s expected 200, got %d body=%s", item.requestID, w.Code, w.Body.String())
		}
	}
	assertCount(t, db, "operation_confirmations", 2)
	assertCount(t, db, "position_transactions", 0)
}

func TestCreateConfirmationRejectsExecutedManuallyAfterWatch(t *testing.T) {
	app, db := testApp(t)
	seedDecision(t, db, "decision_watch_exec", "watch", "formal_trade_advice")
	_, err := db.Exec(`INSERT INTO portfolio_snapshots (snapshot_id,snapshot_time,cash,total_assets,cash_ratio,high_risk_ratio,position_count,source,created_at) VALUES (?,?,?,?,?,?,?,?,?)`, "snap_watch_exec", "2026-01-01T00:00:00Z", 100, 100, 1, 0, 0, "manual", "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed snapshot: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/decision_watch_exec/confirmations", bytes.NewBufferString(`{"confirmation_type":"executed_manually","operation_type":"buy","symbol":"510300","quantity":1,"price":2,"executed_at":"2026-01-01T00:00:00Z"}`))
	req.SetPathValue("decision_id", "decision_watch_exec")
	req.Header.Set("X-Request-ID", "req_watch_exec")
	w := httptest.NewRecorder()

	app.CreateConfirmation(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
	assertCount(t, db, "operation_confirmations", 0)
	assertCount(t, db, "position_transactions", 0)
}

func TestCreateConfirmationExecutedManuallyWritesRequiredFacts(t *testing.T) {
	app, db := testApp(t)
	seedDecision(t, db, "decision_exec", "pending", "formal_trade_advice")
	_, err := db.Exec(`INSERT INTO portfolio_snapshots (snapshot_id,snapshot_time,cash,total_assets,cash_ratio,high_risk_ratio,position_count,source,created_at) VALUES (?,?,?,?,?,?,?,?,?)`, "snap_exec", "2026-01-01T00:00:00Z", 100, 100, 1, 0, 0, "manual", "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed snapshot: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/decision_exec/confirmations", bytes.NewBufferString(`{"confirmation_type":"executed_manually","operation_type":"buy","symbol":"510300","quantity":10,"price":2.5,"fees":1,"executed_at":"2026-01-01T00:00:00Z"}`))
	req.SetPathValue("decision_id", "decision_exec")
	req.Header.Set("X-Request-ID", "req_exec")
	w := httptest.NewRecorder()

	app.CreateConfirmation(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	assertCount(t, db, "operation_confirmations", 1)
	assertCount(t, db, "position_transactions", 1)
	assertCount(t, db, "positions", 1)
	assertCount(t, db, "portfolio_snapshots", 2)
	assertCount(t, db, "position_snapshots", 1)
	assertCount(t, db, "audit_events", 1)
	var costPrice, cash, txFees, confirmationFees float64
	if err := db.QueryRow(`SELECT cost_price FROM positions WHERE symbol='510300'`).Scan(&costPrice); err != nil {
		t.Fatalf("read position cost: %v", err)
	}
	if err := db.QueryRow(`SELECT cash FROM portfolio_snapshots ORDER BY snapshot_time DESC LIMIT 1`).Scan(&cash); err != nil {
		t.Fatalf("read latest cash: %v", err)
	}
	if err := db.QueryRow(`SELECT fees FROM position_transactions WHERE confirmation_id IN (SELECT confirmation_id FROM operation_confirmations WHERE decision_id='decision_exec')`).Scan(&txFees); err != nil {
		t.Fatalf("read tx fees: %v", err)
	}
	if err := db.QueryRow(`SELECT fees FROM operation_confirmations WHERE decision_id='decision_exec'`).Scan(&confirmationFees); err != nil {
		t.Fatalf("read confirmation fees: %v", err)
	}
	if !floatClose(costPrice, 2.6) || !floatClose(cash, 74) || !floatClose(txFees, 1) || !floatClose(confirmationFees, 1) {
		t.Fatalf("expected fees reflected in cost, cash and tx, cost=%v cash=%v fees=%v", costPrice, cash, txFees)
	}
}

func TestCreateConfirmationExecutedSellSubtractsFeesFromCash(t *testing.T) {
	app, db := testApp(t)
	seedDecision(t, db, "decision_sell_fee", "pending", "formal_trade_advice")
	_, err := db.Exec(`INSERT INTO portfolio_snapshots (snapshot_id,snapshot_time,cash,total_assets,cash_ratio,high_risk_ratio,position_count,source,created_at) VALUES (?,?,?,?,?,?,?,?,?)`, "snap_sell_fee", "2026-01-01T00:00:00Z", 100, 130, 100.0/130.0, 0, 1, "manual", "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed snapshot: %v", err)
	}
	_, err = db.Exec(`INSERT INTO positions (position_id,symbol,name,quantity,cost_price,current_price,market_value,unrealized_profit_ratio,position_state,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?)`, "pos_sell_fee", "510300", "沪深300", 10, 2, 3, 30, 0.5, "normal", "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed position: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/decision_sell_fee/confirmations", bytes.NewBufferString(`{"confirmation_type":"executed_manually","operation_type":"sell","symbol":"510300","quantity":4,"price":3,"fees":1,"executed_at":"2026-01-01T00:00:00Z"}`))
	req.SetPathValue("decision_id", "decision_sell_fee")
	req.Header.Set("X-Request-ID", "req_sell_fee")
	w := httptest.NewRecorder()

	app.CreateConfirmation(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var cash, txFees float64
	if err := db.QueryRow(`SELECT cash FROM portfolio_snapshots ORDER BY snapshot_time DESC LIMIT 1`).Scan(&cash); err != nil {
		t.Fatalf("read latest cash: %v", err)
	}
	if err := db.QueryRow(`SELECT fees FROM position_transactions WHERE confirmation_id IN (SELECT confirmation_id FROM operation_confirmations WHERE decision_id='decision_sell_fee')`).Scan(&txFees); err != nil {
		t.Fatalf("read tx fees: %v", err)
	}
	if !floatClose(cash, 111) || !floatClose(txFees, 1) {
		t.Fatalf("expected sell fees subtracted from cash and saved, cash=%v fees=%v", cash, txFees)
	}
}

func TestCreateConfirmationExecutedClearRemovesCurrentPosition(t *testing.T) {
	app, db := testApp(t)
	seedDecision(t, db, "decision_clear", "pending", "formal_trade_advice")
	_, err := db.Exec(`INSERT INTO portfolio_snapshots (snapshot_id,snapshot_time,cash,total_assets,cash_ratio,high_risk_ratio,position_count,source,created_at) VALUES (?,?,?,?,?,?,?,?,?)`, "snap_before_clear", "2026-01-01T00:00:00Z", 100, 130, 100.0/130.0, 0, 1, "manual", "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed snapshot: %v", err)
	}
	_, err = db.Exec(`INSERT INTO positions (position_id,symbol,name,quantity,cost_price,current_price,market_value,unrealized_profit_ratio,position_state,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?)`, "pos_clear", "510300", "沪深300", 10, 2, 3, 30, 0.5, "normal", "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed position: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/decision_clear/confirmations", bytes.NewBufferString(`{"confirmation_type":"executed_manually","operation_type":"sell","symbol":"510300","quantity":10,"price":3,"executed_at":"2026-01-01T00:00:00Z"}`))
	req.SetPathValue("decision_id", "decision_clear")
	req.Header.Set("X-Request-ID", "req_clear")
	w := httptest.NewRecorder()

	app.CreateConfirmation(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	assertCount(t, db, "positions", 0)
	assertCount(t, db, "position_snapshots", 0)
	var cash, totalAssets float64
	var positionCount int
	if err := db.QueryRow(`SELECT cash,total_assets,position_count FROM portfolio_snapshots ORDER BY snapshot_time DESC LIMIT 1`).Scan(&cash, &totalAssets, &positionCount); err != nil {
		t.Fatalf("read latest snapshot: %v", err)
	}
	if cash != 130 || totalAssets != 130 || positionCount != 0 {
		t.Fatalf("expected cash-only snapshot after clear, cash=%v total=%v count=%d", cash, totalAssets, positionCount)
	}
}

func TestCreateConfirmationPlannedDoesNotWriteAccountFacts(t *testing.T) {
	app, db := testApp(t)
	seedDecision(t, db, "decision_plan", "pending", "formal_trade_advice")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/decision_plan/confirmations", bytes.NewBufferString(`{"confirmation_type":"planned","note":"明天观察"}`))
	req.SetPathValue("decision_id", "decision_plan")
	req.Header.Set("X-Request-ID", "req_plan")
	w := httptest.NewRecorder()

	app.CreateConfirmation(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	assertCount(t, db, "operation_confirmations", 1)
	assertCount(t, db, "position_transactions", 0)
	assertCount(t, db, "portfolio_snapshots", 0)
	assertCount(t, db, "audit_events", 1)
}

func TestCreateConfirmationExecutedRollsBackWhenPositionDeleteFails(t *testing.T) {
	app, db := testApp(t)
	app.ConfirmationSvc = service.NewConfirmationService(failingDeleteTransactor{inner: appsqlite.NewTransactor(db)})
	seedDecision(t, db, "decision_rollback_clear", "pending", "formal_trade_advice")
	_, err := db.Exec(`INSERT INTO portfolio_snapshots (snapshot_id,snapshot_time,cash,total_assets,cash_ratio,high_risk_ratio,position_count,source,created_at) VALUES (?,?,?,?,?,?,?,?,?)`, "snap_before_rollback", "2026-01-01T00:00:00Z", 100, 130, 100.0/130.0, 0, 1, "manual", "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed snapshot: %v", err)
	}
	_, err = db.Exec(`INSERT INTO positions (position_id,symbol,name,quantity,cost_price,current_price,market_value,unrealized_profit_ratio,position_state,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?)`, "pos_rollback", "510300", "沪深300", 10, 2, 3, 30, 0.5, "normal", "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed position: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/decision_rollback_clear/confirmations", bytes.NewBufferString(`{"confirmation_type":"executed_manually","operation_type":"sell","symbol":"510300","quantity":10,"price":3,"executed_at":"2026-01-01T00:00:00Z"}`))
	req.SetPathValue("decision_id", "decision_rollback_clear")
	req.Header.Set("X-Request-ID", "req_rollback_clear")
	w := httptest.NewRecorder()

	app.CreateConfirmation(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected internal error, got %d body=%s", w.Code, w.Body.String())
	}
	var status string
	if err := db.QueryRow(`SELECT confirmation_status FROM decision_records WHERE decision_id=?`, "decision_rollback_clear").Scan(&status); err != nil {
		t.Fatalf("read decision: %v", err)
	}
	if status != "pending" {
		t.Fatalf("expected decision rollback to pending, got %s", status)
	}
	assertCount(t, db, "operation_confirmations", 0)
	assertCount(t, db, "position_transactions", 0)
	assertCount(t, db, "audit_events", 0)
	var positionCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM positions WHERE position_id=?`, "pos_rollback").Scan(&positionCount); err != nil {
		t.Fatalf("count position: %v", err)
	}
	if positionCount != 1 {
		t.Fatalf("expected original position kept, got %d", positionCount)
	}
	var snapshotCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM portfolio_snapshots`).Scan(&snapshotCount); err != nil {
		t.Fatalf("count snapshots: %v", err)
	}
	if snapshotCount != 1 {
		t.Fatalf("expected manual snapshot rolled back, got %d", snapshotCount)
	}
}

func TestCreateConfirmationRejectsRepeatedTerminalState(t *testing.T) {
	app, db := testApp(t)
	seedDecision(t, db, "decision_done", "executed_manually", "formal_trade_advice")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/decision_done/confirmations", bytes.NewBufferString(`{"confirmation_type":"marked_error"}`))
	req.SetPathValue("decision_id", "decision_done")
	req.Header.Set("X-Request-ID", "req_done")
	w := httptest.NewRecorder()

	app.CreateConfirmation(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
	assertCount(t, db, "operation_confirmations", 0)
	assertCount(t, db, "error_cases", 0)
}

func floatClose(got, want float64) bool {
	if got > want {
		return got-want < 1e-9
	}
	return want-got < 1e-9
}

func testApp(t *testing.T) (*App, *sql.DB) {
	t.Helper()
	store, err := appsqlite.Open(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	if err := appsqlite.Migrate(context.Background(), store.DB); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	transactor := appsqlite.NewTransactor(store.DB)
	repos := repository.Repositories{
		DecisionRepo:                  appsqlite.NewDecisionRepository(store.DB),
		AuditRepo:                     appsqlite.NewAuditRepository(store.DB),
		RuleRepo:                      appsqlite.NewRuleRepository(store.DB),
		MarketRepo:                    appsqlite.NewMarketRepository(store.DB),
		SettingsRepo:                  appsqlite.NewSettingsRepository(store.DB),
		IntelligenceRepo:              appsqlite.NewIntelligenceRepository(store.DB),
		NotificationRepo:              appsqlite.NewNotificationRepository(store.DB),
		DailyAutoRunRepo:              appsqlite.NewDailyAutoRunRepository(store.DB),
		DailyDisciplineReportRepo:     appsqlite.NewDailyDisciplineReportRepository(store.DB),
		RiskAlertRepo:                 appsqlite.NewRiskAlertRepository(store.DB),
		RuleEffectRepo:                appsqlite.NewRuleEffectRepository(store.DB),
		PortfolioRepo:                 appsqlite.NewPortfolioRepository(store.DB),
		DataQualityGateResolutionRepo: appsqlite.NewDataQualityGateResolutionRepository(store.DB),
	}
	return NewApp(workflow.NewWorkflowDependencies(repos, transactor), repos, transactor), store.DB
}

func seedDecision(t *testing.T, db *sql.DB, decisionID, status, recordType string) {
	t.Helper()
	_, err := db.Exec(`INSERT INTO decision_records (decision_id,request_id,workflow_type,workflow_status,record_type,dashboard_state,final_verdict_status,final_verdict_text,confirmation_status,rule_version,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?)`, decisionID, "req_seed", "consultation", "completed", recordType, "normal", "hold", "持有", status, "v3.0", "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed decision: %v", err)
	}
}

func assertCount(t *testing.T, db *sql.DB, table string, want int) {
	t.Helper()
	var got int
	if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&got); err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	if got != want {
		t.Fatalf("%s count: want %d got %d", table, want, got)
	}
}

type failingDeleteTransactor struct {
	inner repository.Transactor
}

func (t failingDeleteTransactor) WithinTx(ctx context.Context, fn func(context.Context, repository.Repositories) error) error {
	return t.inner.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		repos.PortfolioRepo = failingDeletePortfolioRepo{PortfolioRepository: repos.PortfolioRepo}
		return fn(ctx, repos)
	})
}

type failingDeletePortfolioRepo struct {
	repository.PortfolioRepository
}

func (r failingDeletePortfolioRepo) DeletePosition(ctx context.Context, positionID string) error {
	return errors.New("delete position failed")
}

var _ repository.Transactor = failingDeleteTransactor{}
var _ repository.PortfolioRepository = failingDeletePortfolioRepo{}
