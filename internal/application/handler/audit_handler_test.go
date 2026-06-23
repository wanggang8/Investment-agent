package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"investment-agent/internal/application/dto"
)

func TestListAuditEventsReturnsContractFields(t *testing.T) {
	app, db := testApp(t)
	_, err := db.Exec(`INSERT INTO audit_events (audit_event_id,request_id,decision_id,workflow_type,node_name,actor,action,node_action,proposal_id,confirmation_id,error_case_id,status,error_code,before_state,after_state,rule_version,snapshot_id,input_ref_type,input_ref,output_ref_type,output_ref,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, "audit_full", "req_audit", "decision_1", "consultation", "RuleNode", "system", "generate_decision", "evaluate_rules", "proposal_1", "confirm_1", "error_1", "failed", "RULE_VERSION_MISSING", "pending", "failed", "v3.0", "snap_1", "symbol", "510300", "decision", "decision_1", "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed audit: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/audit/events", nil)
	req.Header.Set("X-Request-ID", "req_list_audit")
	w := httptest.NewRecorder()

	app.ListAuditEvents(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data struct {
			Items []dto.AuditEventDTO `json:"items"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Data.Items) != 1 {
		t.Fatalf("expected one audit event, got %+v", body.Data.Items)
	}
		got := body.Data.Items[0]
	if got.AuditEventID != "audit_full" || got.EventID != "audit_full" || got.RequestID != "req_audit" || got.DecisionID != "decision_1" || got.WorkflowType != "consultation" || got.NodeName != "RuleNode" || got.Actor != "system" || got.Action != "generate_decision" || got.NodeAction != "evaluate_rules" || got.ProposalID != "proposal_1" || got.ConfirmationID != "confirm_1" || got.ErrorCaseID != "error_1" || got.Status != "failed" || got.ErrorCode != "RULE_VERSION_MISSING" || got.BeforeState != "pending" || got.AfterState != "failed" || got.RuleVersion != "v3.0" || got.SnapshotID != "snap_1" || got.InputRefType != "symbol" || got.InputRef != "510300" || got.OutputRefType != "decision" || got.OutputRef != "decision_1" || got.CreatedAt == "" {
		t.Fatalf("audit contract fields not preserved: %+v", got)
	}
}
