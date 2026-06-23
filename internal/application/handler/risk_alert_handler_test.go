package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"investment-agent/internal/domain/repository"
)

func TestListRiskAlertsReturnsActiveItemsWithSafeLinks(t *testing.T) {
	app, db := testApp(t)
	seedRiskAlert(t, db, repository.RiskAlert{AlertID: "risk_1", RiskType: "valuation_high", Severity: "warning", SOPStatus: "active", Symbol: "510300", TriggerSummary: "PE 分位高于 80%", ProhibitedActionsJSON: `["新增买入"]`, SuggestedActionsJSON: `["人工复核分批止盈"]`, RelatedDecisionID: "dec_1", RelatedReportID: "report_1", RelatedNotificationID: "notif_1", RelatedAuditEventID: "audit_1", CreatedAt: "2026-06-15T09:30:00Z", UpdatedAt: "2026-06-15T09:30:00Z"})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/risk-alerts?status=active", nil)
	w := httptest.NewRecorder()
	app.ListRiskAlerts(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	out := decodeRiskAlertList(t, w)
	if len(out.Items) != 1 || out.Items[0].AlertID != "risk_1" || out.Items[0].Link != "/risk-alerts/risk_1" {
		t.Fatalf("unexpected risk list: %+v", out)
	}
	if out.Items[0].ProhibitedActions[0] != "新增买入" || out.Items[0].SafetyNote == "" {
		t.Fatalf("expected safe actions, got %+v", out.Items[0])
	}
}

func TestGetRiskAlertAndLifecycleAction(t *testing.T) {
	app, db := testApp(t)
	seedRiskAlert(t, db, repository.RiskAlert{AlertID: "risk_1", RiskType: "data_degraded", Severity: "warning", SOPStatus: "active", Symbol: "510300", TriggerSummary: "source health stale", CreatedAt: "2026-06-15T09:30:00Z", UpdatedAt: "2026-06-15T09:30:00Z"})

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/risk-alerts/risk_1", nil)
	getReq.SetPathValue("alert_id", "risk_1")
	getW := httptest.NewRecorder()
	app.GetRiskAlert(getW, getReq)
	if getW.Code != http.StatusOK {
		t.Fatalf("expected detail 200, got %d body=%s", getW.Code, getW.Body.String())
	}
	if detail := decodeRiskAlertDetail(t, getW); detail.AlertID != "risk_1" || detail.SOPStatus != "active" {
		t.Fatalf("unexpected detail: %+v", detail)
	}

	body := bytes.NewBufferString(`{"status":"resolved","reason":"数据恢复"}`)
	postReq := httptest.NewRequest(http.MethodPost, "/api/v1/risk-alerts/risk_1/lifecycle", body)
	postReq.SetPathValue("alert_id", "risk_1")
	postW := httptest.NewRecorder()
	app.UpdateRiskAlertLifecycle(postW, postReq)
	if postW.Code != http.StatusOK {
		t.Fatalf("expected lifecycle 200, got %d body=%s", postW.Code, postW.Body.String())
	}
	updated := decodeRiskAlertDetail(t, postW)
	if updated.SOPStatus != "resolved" || updated.ResolutionReason != "数据恢复" {
		t.Fatalf("unexpected updated risk: %+v", updated)
	}
	if confirmations := tableCount(t, db, "operation_confirmations"); confirmations != 0 {
		t.Fatalf("expected no trading confirmation, got %d", confirmations)
	}
}

type riskAlertEnvelope struct {
	Data riskAlertPayload `json:"data"`
}

type riskAlertListEnvelope struct {
	Data struct {
		Items []riskAlertPayload `json:"items"`
		Total int                `json:"total"`
	} `json:"data"`
}

type riskAlertPayload struct {
	AlertID           string   `json:"alert_id"`
	RiskType          string   `json:"risk_type"`
	Severity          string   `json:"severity"`
	SOPStatus         string   `json:"sop_status"`
	Symbol            string   `json:"symbol"`
	TriggerSummary    string   `json:"trigger_summary"`
	ProhibitedActions []string `json:"prohibited_actions"`
	SuggestedActions  []string `json:"suggested_actions"`
	Link              string   `json:"link"`
	SafetyNote        string   `json:"safety_note"`
	ResolutionReason  string   `json:"resolution_reason"`
}

func decodeRiskAlertList(t *testing.T, w *httptest.ResponseRecorder) struct {
	Items []riskAlertPayload `json:"items"`
	Total int                `json:"total"`
} {
	t.Helper()
	var envelope riskAlertListEnvelope
	if err := json.Unmarshal(w.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode risk alert list: %v body=%s", err, w.Body.String())
	}
	return envelope.Data
}

func decodeRiskAlertDetail(t *testing.T, w *httptest.ResponseRecorder) riskAlertPayload {
	t.Helper()
	var envelope riskAlertEnvelope
	if err := json.Unmarshal(w.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode risk alert detail: %v body=%s", err, w.Body.String())
	}
	return envelope.Data
}

func seedRiskAlert(t *testing.T, db *sql.DB, alert repository.RiskAlert) {
	t.Helper()
	_, err := db.Exec(`INSERT INTO risk_alerts (alert_id,risk_type,severity,sop_status,symbol,trigger_summary,trigger_context_json,prohibited_actions_json,suggested_actions_json,related_decision_id,related_report_id,related_notification_id,related_audit_event_id,last_triggered_at,resolved_at,resolution_reason,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, alert.AlertID, string(alert.RiskType), string(alert.Severity), string(alert.SOPStatus), alert.Symbol, alert.TriggerSummary, nullStringLocal(alert.TriggerContextJSON), nullStringLocal(alert.ProhibitedActionsJSON), nullStringLocal(alert.SuggestedActionsJSON), nullStringLocal(alert.RelatedDecisionID), nullStringLocal(alert.RelatedReportID), nullStringLocal(alert.RelatedNotificationID), nullStringLocal(alert.RelatedAuditEventID), nullStringLocal(alert.LastTriggeredAt), nullStringLocal(alert.ResolvedAt), nullStringLocal(alert.ResolutionReason), alert.CreatedAt, alert.UpdatedAt)
	if err != nil {
		t.Fatalf("seed risk alert: %v", err)
	}
}

func tableCount(t *testing.T, db *sql.DB, table string) int {
	t.Helper()
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	return count
}
