package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetDailyAutoRunStatusReturnsConfigAndLatestState(t *testing.T) {
	app, db := testApp(t)
	app.Deps.DailyAutoRunConfig.Enabled = true
	app.Deps.DailyAutoRunConfig.RunTime = "08:30"
	app.Deps.DailyAutoRunConfig.Timezone = "Asia/Shanghai"
	_, err := db.Exec(`INSERT INTO daily_auto_run_states (run_id,idempotency_key,local_date,scope,symbol_set_hash,status,last_run_at,next_run_at,failure_code,failure_reason,created_at,updated_at) VALUES ('auto_run_1','2026-06-07:holdings:abc:v1','2026-06-07','holdings','abc','failed','2026-06-07T00:30:00Z','2026-06-08T00:30:00Z','missing_prerequisites','缺少本地持仓','2026-06-07T00:30:00Z','2026-06-07T00:30:00Z')`)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/daily-auto-run/status", nil)
	w := httptest.NewRecorder()
	app.GetDailyAutoRunStatus(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	for _, want := range []string{`"enabled":true`, `"status":"failed"`, `"last_run_at":"2026-06-07T00:30:00Z"`, `"next_run_at":"2026-06-08T00:30:00Z"`, `"failure_code":"missing_prerequisites"`, `"failure_reason":"缺少本地持仓"`, `"latest_notification_link":"/notifications?source_id=2026-06-07%3Aholdings%3Aabc%3Av1"`, `"latest_audit_link":"/audit?input_ref=2026-06-07%3Aholdings%3Aabc%3Av1"`, `"safety_note":"仅记录本地每日刷新和纪律评估结果，不会自动执行交易，需人工复核。"`} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected body to contain %s, got %s", want, body)
		}
	}
}

func TestGetDailyAutoRunStatusEscapesLatestLinks(t *testing.T) {
	app, db := testApp(t)
	_, err := db.Exec(`INSERT INTO daily_auto_run_states (run_id,idempotency_key,local_date,scope,symbol_set_hash,status,last_run_at,next_run_at,failure_code,failure_reason,created_at,updated_at) VALUES ('run&status=failed','manual&status=failed','2026-06-07','holdings','abc','failed','2026-06-07T00:30:00Z','','daily_discipline_failed','x','2026-06-07T00:30:00Z','2026-06-07T00:30:00Z')`)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/daily-auto-run/status", nil)
	w := httptest.NewRecorder()
	app.GetDailyAutoRunStatus(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	for _, want := range []string{`"latest_notification_link":"/notifications?source_id=manual%26status%3Dfailed"`, `"latest_audit_link":"/audit?input_ref=manual%26status%3Dfailed"`, `"latest_decision_link":"/decisions?request_id=run%26status%3Dfailed"`} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected body to contain %s, got %s", want, body)
		}
	}
}

func TestGetDailyAutoRunStatusDisabledWithoutState(t *testing.T) {
	app, _ := testApp(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/daily-auto-run/status", nil)
	w := httptest.NewRecorder()
	app.GetDailyAutoRunStatus(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	for _, want := range []string{`"enabled":false`, `"status":"disabled"`, `"safety_note":"仅记录本地每日刷新和纪律评估结果，不会自动执行交易，需人工复核。"`} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected body to contain %s, got %s", want, body)
		}
	}
}
