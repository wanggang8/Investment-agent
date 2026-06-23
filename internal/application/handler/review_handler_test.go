package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"investment-agent/internal/application/dto"
)

func TestGetReviewSummaryAggregatesPeriodicFacts(t *testing.T) {
	app, db := testApp(t)
	seedReviewFacts(t, db)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/review/summary?period=quarterly", nil)
	w := httptest.NewRecorder()

	app.GetReviewSummary(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data dto.ReviewSummaryResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Data.Period != "quarterly" || body.Data.DecisionCount != 2 || body.Data.ConfirmationCount != 2 || body.Data.ExecutedManuallyCount != 1 || body.Data.PlannedCount != 1 || body.Data.ErrorCaseCount != 1 || body.Data.RuleProposalCount != 1 || body.Data.AuditEventCount != 2 || body.Data.RuleHitCount != 3 || body.Data.MissingEvidenceCount != 1 || body.Data.DegradedCount != 1 || body.Data.MisjudgmentCount != 1 {
		t.Fatalf("unexpected review summary: %+v", body.Data)
	}
	if body.Data.OpsStatus.IndexStatus != "success" {
		t.Fatalf("expected review index status success, got %+v", body.Data.OpsStatus)
	}
	if len(body.Data.RuleSuggestions) != 1 || body.Data.RuleSuggestions[0].ProposalID != "prop_review" || body.Data.RuleSuggestions[0].CanAutoApply {
		t.Fatalf("expected gated rule suggestion, got %+v", body.Data.RuleSuggestions)
	}
	if len(body.Data.TrackingLinks) != 3 {
		t.Fatalf("expected tracking links, got %+v", body.Data.TrackingLinks)
	}
	if len(body.Data.AttributionSummaries) != 2 || body.Data.AttributionSummaries[0].DecisionID == "" || body.Data.AttributionSummaries[0].EvidenceStatus == "" || body.Data.AttributionSummaries[0].Outcome == "" {
		t.Fatalf("expected traceable attribution summaries, got %+v", body.Data.AttributionSummaries)
	}
	if len(body.Data.RecurringErrorTags) != 1 || body.Data.RecurringErrorTags[0].Tag != "rule_threshold_issue" || body.Data.RecurringErrorTags[0].Count != 1 {
		t.Fatalf("expected recurring error tags from error cases, got %+v", body.Data.RecurringErrorTags)
	}
	if len(body.Data.MissingEvidenceThemes) != 1 || body.Data.MissingEvidenceThemes[0].Status != "failed" || body.Data.MissingEvidenceThemes[0].Count != 1 {
		t.Fatalf("expected missing evidence themes from local facts, got %+v", body.Data.MissingEvidenceThemes)
	}
	if len(body.Data.RuleProposalOutcomes) != 1 || body.Data.RuleProposalOutcomes[0].ProposalID != "prop_review" || body.Data.RuleProposalOutcomes[0].Status != "pending_user_confirm" {
		t.Fatalf("expected rule proposal outcomes, got %+v", body.Data.RuleProposalOutcomes)
	}
	if len(body.Data.DegradedWorkflows) != 1 || body.Data.DegradedWorkflows[0].DecisionID != "decision_review_1" {
		t.Fatalf("expected degraded workflow trace, got %+v", body.Data.DegradedWorkflows)
	}
	if len(body.Data.RuleEffectTracking) != 1 || body.Data.RuleEffectTracking[0].AppliedRuleVersion != "v3.1" || body.Data.RuleEffectTracking[0].TrendDirection != "worsened" {
		t.Fatalf("expected rule effect tracking in review, got %+v", body.Data.RuleEffectTracking)
	}
	relatedRiskIDs, ok := body.Data.RuleEffectTracking[0].RelatedRiskAlertIDs.([]any)
	if !ok || len(relatedRiskIDs) != 1 || relatedRiskIDs[0] != "risk_review" {
		t.Fatalf("expected related risk alerts in review tracking, got %+v", body.Data.RuleEffectTracking[0].RelatedRiskAlertIDs)
	}
	var notificationType, sourceType, sourceID string
	if err := db.QueryRow(`SELECT type,COALESCE(source_type,''),COALESCE(source_id,'') FROM notifications WHERE read_at IS NULL ORDER BY created_at DESC LIMIT 1`).Scan(&notificationType, &sourceType, &sourceID); err != nil {
		t.Fatalf("read review notification: %v", err)
	}
	if notificationType != "review_degraded" || sourceType != "review_summary" || sourceID != "quarterly" {
		t.Fatalf("expected review degraded notification, got type=%s source=%s/%s", notificationType, sourceType, sourceID)
	}
}

func TestGetReviewSummaryUsesCurrentTimeWindow(t *testing.T) {
	app, db := testApp(t)
	seedReviewFacts(t, db)
	if _, err := db.ExecContext(context.Background(), `UPDATE decision_records SET created_at='2020-01-01T00:00:00Z'; UPDATE rule_proposals SET created_at='2020-01-01T00:00:00Z'; UPDATE audit_events SET created_at='2020-01-01T00:00:00Z'; UPDATE error_cases SET created_at='2020-01-01T00:00:00Z'`); err != nil {
		t.Fatalf("move review facts to old window: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/review/summary?period=monthly", nil)
	w := httptest.NewRecorder()

	app.GetReviewSummary(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data dto.ReviewSummaryResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Data.DecisionCount != 0 || body.Data.RuleProposalCount != 0 || body.Data.AuditEventCount != 0 || body.Data.ErrorCaseCount != 0 {
		t.Fatalf("expected old facts outside monthly window, got %+v", body.Data)
	}
	if body.Data.OpsStatus.ReviewStatus != "empty" || body.Data.OpsStatus.DataSourceStatus != "unknown" || body.Data.OpsStatus.IndexStatus != "unknown" {
		t.Fatalf("expected empty ops status for empty review window, got %+v", body.Data.OpsStatus)
	}
}

func TestGetReviewSummaryRejectsInvalidPeriod(t *testing.T) {
	app, _ := testApp(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/review/summary?period=weekly", nil)
	w := httptest.NewRecorder()

	app.GetReviewSummary(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
}

func seedReviewFacts(t *testing.T, db *sql.DB) {
	t.Helper()
	seedDecision(t, db, "decision_review_1", "executed_manually", "formal_trade_advice")
	seedDecision(t, db, "decision_review_2", "planned", "formal_trade_advice")
	if _, err := db.ExecContext(context.Background(), `UPDATE decision_records SET workflow_status='degraded', source_verification_status='failed', triggered_rules_json='["rule_a","rule_b"]', created_at='2026-05-01T00:00:00Z' WHERE decision_id='decision_review_1'`); err != nil {
		t.Fatalf("seed review decision fields: %v", err)
	}
	if _, err := db.ExecContext(context.Background(), `UPDATE decision_records SET triggered_rules_json='["rule_c"]', created_at='2026-05-01T00:00:00Z' WHERE decision_id='decision_review_2'`); err != nil {
		t.Fatalf("seed review rule hit: %v", err)
	}
	if _, err := db.ExecContext(context.Background(), `INSERT INTO operation_confirmations (confirmation_id,decision_id,confirmation_type,operation_type,symbol,quantity,price,fees,executed_at,payload_json,note,created_at) VALUES ('confirm_exec','decision_review_1','executed_manually','buy','510300',100,4,1,'2026-05-02T00:00:00Z','{}','线下执行','2026-05-02T00:00:00Z'),('confirm_plan','decision_review_2','planned',NULL,NULL,0,0,0,NULL,'{}','计划观察','2026-05-03T00:00:00Z')`); err != nil {
		t.Fatalf("seed confirmations: %v", err)
	}
	if _, err := db.ExecContext(context.Background(), `INSERT INTO error_cases (error_case_id,decision_id,confirmation_id,actual_outcome,root_cause_tag,lesson_learned,created_at) VALUES ('err_review','decision_review_1','confirm_exec','结果偏离','rule_threshold_issue','季度复盘样本','2026-05-04T00:00:00Z')`); err != nil {
		t.Fatalf("seed error case: %v", err)
	}
	if _, err := db.ExecContext(context.Background(), `INSERT INTO rule_proposals (proposal_id,proposal_type,status,title,proposal_version,before_rule_json,after_rule_json,reason,impact_scope_json,risk_notes_json,sample_count,related_error_cases_json,created_at) VALUES ('prop_review','threshold','pending_user_confirm','季度阈值复盘','v1','{}','{"content":"调整阈值"}','误判样本触发','[]','[]',3,'["err_review"]','2026-05-05T00:00:00Z')`); err != nil {
		t.Fatalf("seed proposal: %v", err)
	}
	if _, err := db.ExecContext(context.Background(), `INSERT INTO rule_effect_tracking (tracking_id,applied_rule_version,proposal_id,period,hit_count,misjudgment_count,missing_evidence_count,degraded_count,risk_alert_count,trend_direction,metrics_json,related_risk_alert_ids_json,safety_note,created_at,updated_at) VALUES ('track_review','v3.1','prop_review','quarterly',4,2,1,1,2,'worsened','{"hit_count":4}','["risk_review"]','只读追踪','2026-05-06T00:00:00Z','2026-05-06T00:00:00Z')`); err != nil {
		t.Fatalf("seed rule effect tracking: %v", err)
	}
	if _, err := db.ExecContext(context.Background(), `INSERT INTO audit_events (audit_event_id,request_id,actor,action,status,input_ref_type,input_ref,output_ref_type,output_ref,created_at) VALUES ('audit_review_1','req_review','user','confirm_operation','success','decision','decision_review_1','confirmation','confirm_exec','2026-05-02T00:00:00Z'),('audit_review_2','req_review','system','create_proposal','success','error_case','err_review','proposal','prop_review','2026-05-05T00:00:00Z')`); err != nil {
		t.Fatalf("seed audit: %v", err)
	}
}
