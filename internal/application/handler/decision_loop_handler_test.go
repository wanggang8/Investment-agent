package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"investment-agent/internal/application/dto"
)

func TestListDecisionLoopsReadOnlyAndSafe(t *testing.T) {
	app, db := testApp(t)
	seedDecisionWithSymbol(t, db, "decision_loop_handler", "510300", "executed_manually", "formal_trade_advice")
	if _, err := db.Exec(`INSERT INTO operation_confirmations (confirmation_id,decision_id,confirmation_type,operation_type,symbol,quantity,price,fees,executed_at,payload_json,note,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`,
		"conf_loop_handler", "decision_loop_handler", "executed_manually", "buy", "510300", 10, 2.5, 1, "2026-01-01T09:30:00Z", `{"raw":"SELECT * FROM secret"}`, "Prompt: 完整 prompt raw HTTP GET /secret HTTP/1.1 /Users/private/key sk-proj-abcdef_123456 -----BEGIN OPENSSH PRIVATE KEY-----abc-----END OPENSSH PRIVATE KEY-----", "2026-01-01T09:20:00Z"); err != nil {
		t.Fatalf("seed confirmation: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO position_transactions (transaction_id,confirmation_id,symbol,operation_type,quantity,price,fees,occurred_at,created_at) VALUES (?,?,?,?,?,?,?,?,?)`,
		"tx_loop_handler", "conf_loop_handler", "510300", "buy", 10, 2.5, 1, "2026-01-01T09:30:00Z", "2026-01-01T09:31:00Z"); err != nil {
		t.Fatalf("seed transaction: %v", err)
	}
	before := handlerDecisionLoopCounts(t, db)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/decision-loops?symbol=510300&limit=5", nil)
	req.Header.Set("X-Request-ID", "req_decision_loop_list")
	w := httptest.NewRecorder()

	app.ListDecisionLoops(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var envelope struct {
		RequestID string                       `json:"request_id"`
		Data      dto.DecisionLoopListResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &envelope); err != nil {
		t.Fatal(err)
	}
	if envelope.RequestID != "req_decision_loop_list" || envelope.Data.Total != 1 {
		t.Fatalf("unexpected envelope: %#v", envelope)
	}
	if envelope.Data.Items[0].DecisionID != "decision_loop_handler" {
		t.Fatalf("unexpected decision loop: %#v", envelope.Data.Items)
	}
	body := w.Body.String()
	for _, forbidden := range []string{"SELECT * FROM", "/Users/private", "prompt:", "Prompt:", "完整 prompt", "raw HTTP", "GET /secret HTTP/1.1", "BEGIN OPENSSH PRIVATE KEY"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("handler response leaked %q: %s", forbidden, body)
		}
	}
	if regexp.MustCompile(`sk-[A-Za-z0-9]{12,}`).MatchString(body) {
		t.Fatalf("handler response leaked complete key-like fragment: %s", body)
	}
	if regexp.MustCompile(`sk-[A-Za-z0-9][A-Za-z0-9_-]{8,}`).MatchString(body) {
		t.Fatalf("handler response leaked hyphenated key-like fragment: %s", body)
	}
	after := handlerDecisionLoopCounts(t, db)
	if !sameHandlerDecisionLoopCounts(before, after) {
		t.Fatalf("handler wrote read-only tables: before=%v after=%v", before, after)
	}
}

func TestGetDecisionLoopReturnsNotFound(t *testing.T) {
	app, _ := testApp(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/decision-loops/missing_decision", nil)
	req.SetPathValue("decision_id", "missing_decision")
	req.Header.Set("X-Request-ID", "req_decision_loop_missing")
	w := httptest.NewRecorder()

	app.GetDecisionLoop(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"code":"NOT_FOUND"`) {
		t.Fatalf("expected NOT_FOUND envelope, got %s", w.Body.String())
	}
}

func TestListDecisionLoopsRejectsInvalidLimit(t *testing.T) {
	app, _ := testApp(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/decision-loops?limit=bad", nil)
	req.Header.Set("X-Request-ID", "req_decision_loop_bad_limit")
	w := httptest.NewRecorder()

	app.ListDecisionLoops(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
}

func seedDecisionWithSymbol(t *testing.T, db *sql.DB, decisionID, symbol, status, recordType string) {
	t.Helper()
	_, err := db.Exec(`INSERT INTO decision_records (decision_id,request_id,workflow_type,symbol,workflow_status,record_type,dashboard_state,source_verification_status,final_verdict_status,final_verdict_text,confirmation_status,rule_version,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`, decisionID, "req_seed", "consultation", symbol, "completed", recordType, "normal", "satisfied", "hold", "持有", status, "v3.0", "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed decision: %v", err)
	}
}

func handlerDecisionLoopCounts(t *testing.T, db *sql.DB) map[string]int {
	t.Helper()
	out := map[string]int{}
	for _, table := range []string{"decision_records", "operation_confirmations", "position_transactions", "error_cases", "risk_alerts", "audit_events", "notifications"} {
		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
			t.Fatalf("count %s: %v", table, err)
		}
		out[table] = count
	}
	return out
}

func sameHandlerDecisionLoopCounts(a, b map[string]int) bool {
	if len(a) != len(b) {
		return false
	}
	for key, av := range a {
		if b[key] != av {
			return false
		}
	}
	return true
}
