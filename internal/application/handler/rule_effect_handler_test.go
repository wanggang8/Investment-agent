package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRuleEffectValidationHandlersExposeValidationAndTracking(t *testing.T) {
	app, db := testApp(t)
	seedRuleProposal(t, db, "prop_effect", "pending_user_confirm", 5)
	seedPassedRuleEffectValidation(t, db, "prop_effect", "val_effect")
	if _, err := db.Exec(`INSERT INTO rule_effect_tracking (tracking_id,applied_rule_version,proposal_id,period,hit_count,misjudgment_count,missing_evidence_count,degraded_count,risk_alert_count,trend_direction,metrics_json,safety_note,created_at,updated_at) VALUES ('track_effect','v3.2','prop_effect','2026-Q3',9,1,0,0,0,'flat','{"hit_count":9}','只读追踪','2026-06-16T00:00:00Z','2026-06-16T00:00:00Z')`); err != nil {
		t.Fatalf("seed tracking: %v", err)
	}

	validationReq := httptest.NewRequest(http.MethodGet, "/api/v1/rule-proposals/prop_effect/effect-validation", nil)
	validationReq.SetPathValue("proposal_id", "prop_effect")
	validationW := httptest.NewRecorder()
	app.GetRuleEffectValidation(validationW, validationReq)
	if validationW.Code != http.StatusOK {
		t.Fatalf("expected validation 200, got %d body=%s", validationW.Code, validationW.Body.String())
	}
	var validationBody struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(validationW.Body.Bytes(), &validationBody); err != nil {
		t.Fatalf("decode validation body: %v", err)
	}
	if validationBody.Data["validation_status"] != "passed" || validationBody.Data["safety_note"] == "" || validationBody.Data["proposal_id"] != "prop_effect" || validationBody.Data["candidate_rule_version"] != "draft" || validationBody.Data["source_explanation"] == nil || validationBody.Data["metrics"] == nil || validationBody.Data["risk_notes"] == nil || validationBody.Data["created_at"] == "" {
		t.Fatalf("unexpected validation body: %+v", validationBody.Data)
	}

	trackingReq := httptest.NewRequest(http.MethodGet, "/api/v1/rule-effect-tracking?rule_version=v3.2", nil)
	trackingW := httptest.NewRecorder()
	app.ListRuleEffectTracking(trackingW, trackingReq)
	if trackingW.Code != http.StatusOK {
		t.Fatalf("expected tracking 200, got %d body=%s", trackingW.Code, trackingW.Body.String())
	}
	var trackingBody struct {
		Data struct {
			Items []map[string]any `json:"items"`
		} `json:"data"`
	}
	if err := json.Unmarshal(trackingW.Body.Bytes(), &trackingBody); err != nil {
		t.Fatalf("decode tracking body: %v", err)
	}
	if len(trackingBody.Data.Items) != 1 || trackingBody.Data.Items[0]["applied_rule_version"] != "v3.2" || trackingBody.Data.Items[0]["metrics"] == nil || trackingBody.Data.Items[0]["created_at"] == "" {
		t.Fatalf("unexpected tracking body: %+v", trackingBody.Data.Items)
	}
}

func TestRefreshRuleEffectValidationHandlerCreatesLocalOnlyValidation(t *testing.T) {
	app, db := testApp(t)
	seedRuleProposal(t, db, "prop_refresh", "pending_user_confirm", 2)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rule-proposals/prop_refresh/effect-validation", bytes.NewBufferString(`{"sample_window":"2026-Q2","sample_count":99,"source_case_count":99,"hit_count":99,"misjudgment_count":0,"missing_evidence_count":0,"degraded_count":0,"risk_alert_count":0}`))
	req.SetPathValue("proposal_id", "prop_refresh")
	req.Header.Set("X-Request-ID", "req_refresh_validation")
	w := httptest.NewRecorder()

	app.RefreshRuleEffectValidation(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected refresh 200, got %d body=%s", w.Code, w.Body.String())
	}
	var validationCount, confirmationCount int
	_ = db.QueryRow(`SELECT COUNT(*) FROM rule_effect_validations WHERE proposal_id='prop_refresh' AND validation_status='insufficient'`).Scan(&validationCount)
	_ = db.QueryRow(`SELECT COUNT(*) FROM operation_confirmations`).Scan(&confirmationCount)
	if validationCount != 1 || confirmationCount != 0 {
		t.Fatalf("expected local validation only, validations=%d confirmations=%d", validationCount, confirmationCount)
	}
}
