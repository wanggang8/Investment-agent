package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListRuleProposalsReturnsFrontendFields(t *testing.T) {
	app, db := testApp(t)
	_, err := db.Exec(`INSERT INTO rule_proposals (proposal_id,proposal_type,status,source_error_case_id,title,proposal_version,before_rule_json,after_rule_json,reason,impact_scope_json,risk_notes_json,sample_count,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`, "prop_list", "risk_rule", "pending_final_confirm", "err_1", "测试提案", "draft", `{"old":1}`, `{"new":2}`, "原因", `{"scope":"portfolio"}`, `{"risk":"low"}`, 3, "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed proposal: %v", err)
	}
	_, err = db.Exec(`INSERT INTO gatekeeper_audits (gatekeeper_audit_id,proposal_id,audit_result,audit_reason,violates_fundamental_rule,has_rule_conflict,allow_apply,audited_rule_version,created_at) VALUES (?,?,?,?,?,?,?,?,?)`, "gk_list", "prop_list", "approved", "通过", 0, 0, 1, "draft", "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed audit: %v", err)
	}
	_, err = db.Exec(`INSERT INTO rule_effect_validations (validation_id,proposal_id,candidate_rule_version,validation_status,sample_count,sample_window,representativeness_status,overfit_risk,replay_result,guardrail_decision,metrics_json,risk_notes_json,safety_note,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, "val_list", "prop_list", "draft", "passed", 5, "2026-Q2", "passed", "low", "passed", "passed", `{"hit_count":5}`, `[]`, "只读验证", "2026-01-01T00:00:00Z", "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed validation: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/rule-proposals", bytes.NewBuffer(nil))
	req.Header.Set("X-Request-ID", "req_rule_list")
	w := httptest.NewRecorder()

	app.ListRuleProposals(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data struct {
			Items []map[string]any `json:"items"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Data.Items) != 1 {
		t.Fatalf("expected one proposal, got %+v", body.Data.Items)
	}
	item := body.Data.Items[0]
	auditResult, okResult := item["audit_result"].(string)
	auditSummary, okSummary := item["audit_summary"].(string)
	for _, key := range []string{"source_error_case_id", "before_rule", "after_rule", "impact_scope", "risk_notes"} {
		if _, ok := item[key]; !ok {
			t.Fatalf("missing %s in %+v", key, item)
		}
	}
	if _, ok := item["impact_scope"].(map[string]any); !ok {
		t.Fatalf("impact_scope should be structured, got %+v", item["impact_scope"])
	}
	if _, ok := item["risk_notes"].(map[string]any); !ok {
		t.Fatalf("risk_notes should be structured, got %+v", item["risk_notes"])
	}
	if !okResult || !okSummary || auditResult == "" || auditSummary == "" {
		t.Fatalf("missing audit fields in %+v", item)
	}
	validation, ok := item["effect_validation"].(map[string]any)
	if !ok || validation["validation_status"] != "passed" || validation["overfit_risk"] != "low" || validation["guardrail_decision"] != "passed" || validation["validation_link"] == "" {
		t.Fatalf("missing effect validation summary in %+v", item)
	}
}

func TestP88SOPAddendumProposalCreatesPendingProposalNotificationAndAudit(t *testing.T) {
	app, db := testApp(t)
	var activeBefore string
	if err := db.QueryRow(`SELECT rule_version FROM rule_versions WHERE status='active' LIMIT 1`).Scan(&activeBefore); err != nil {
		t.Fatalf("read active rule: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/rule-proposals/sop-addendum", bytes.NewBufferString(`{"scenario_key":"p88_uncovered_liquidity_gap","scenario_title":"连续流动性缺口未覆盖","occurrence_count":4,"sample_window":"2026-Q2"}`))
	req.Header.Set("X-Request-ID", "req_p88_sop_addendum")
	w := httptest.NewRecorder()

	app.CreateSOPAddendumProposal(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data struct {
			ProposalID     string   `json:"proposal_id"`
			Status         string   `json:"status"`
			NotificationID string   `json:"notification_id"`
			AuditEventIDs  []string `json:"audit_event_ids"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.ProposalID == "" || body.Data.Status != "pending_user_confirm" || body.Data.NotificationID == "" || len(body.Data.AuditEventIDs) != 1 {
		t.Fatalf("unexpected response: %+v", body.Data)
	}
	var proposalType, status, title, afterRule, riskNotes string
	var sampleCount int
	if err := db.QueryRow(`SELECT proposal_type,status,title,after_rule_json,risk_notes_json,sample_count FROM rule_proposals WHERE proposal_id=?`, body.Data.ProposalID).Scan(&proposalType, &status, &title, &afterRule, &riskNotes, &sampleCount); err != nil {
		t.Fatalf("read proposal: %v", err)
	}
	if proposalType != "sop" || status != "pending_user_confirm" || sampleCount != 4 || !bytes.Contains([]byte(title), []byte("SOP 补充提案")) || !bytes.Contains([]byte(afterRule), []byte(`"auto_apply":false`)) || !bytes.Contains([]byte(riskNotes), []byte("不会自动应用规则")) {
		t.Fatalf("unexpected stored proposal type=%s status=%s title=%q sample=%d after=%s risk=%s", proposalType, status, title, sampleCount, afterRule, riskNotes)
	}
	var notificationCount, auditCount, activeAfterCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM notifications WHERE notification_id=? AND source_type='rule_proposal' AND source_id=?`, body.Data.NotificationID, body.Data.ProposalID).Scan(&notificationCount); err != nil {
		t.Fatalf("count notification: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM audit_events WHERE request_id='req_p88_sop_addendum' AND action='create_proposal' AND proposal_id=?`, body.Data.ProposalID).Scan(&auditCount); err != nil {
		t.Fatalf("count audit: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM rule_versions WHERE status='active' AND rule_version=?`, activeBefore).Scan(&activeAfterCount); err != nil {
		t.Fatalf("count active rule: %v", err)
	}
	if notificationCount != 1 || auditCount != 1 || activeAfterCount != 1 {
		t.Fatalf("expected notification/audit and unchanged active rule, notification=%d audit=%d active=%d", notificationCount, auditCount, activeAfterCount)
	}
}

func TestRuleProposalConfirmRunsGatekeeperAuditAndAllowsFinalConfirm(t *testing.T) {
	app, db := testApp(t)
	seedRuleProposal(t, db, "prop_confirm", "pending_user_confirm", 3)

	_, err := db.Exec(`UPDATE rule_versions SET status='archived' WHERE status='active'`)
	if err != nil {
		t.Fatalf("archive seeded rule: %v", err)
	}
	_, err = db.Exec(`INSERT INTO rule_versions (rule_version,status,rules_json,effective_at,created_at) VALUES (?,?,?,?,?)`, "v4.0", "active", "{}", "2026-01-01T00:00:00Z", "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed active rule: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rule-proposals/prop_confirm/confirm", bytes.NewBufferString(`{"confirm":true}`))
	req.SetPathValue("proposal_id", "prop_confirm")
	req.Header.Set("X-Request-ID", "req_rule_confirm")
	w := httptest.NewRecorder()

	app.ConfirmRuleProposal(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var status string
	_ = db.QueryRow(`SELECT status FROM rule_proposals WHERE proposal_id='prop_confirm'`).Scan(&status)
	if status != "pending_final_confirm" {
		t.Fatalf("expected pending_final_confirm, got %s", status)
	}
	var auditResult, auditReason, backtestMetrics, auditedRuleVersion string
	if err := db.QueryRow(`SELECT audit_result,audit_reason,COALESCE(backtest_metrics_json,''),audited_rule_version FROM gatekeeper_audits WHERE proposal_id='prop_confirm'`).Scan(&auditResult, &auditReason, &backtestMetrics, &auditedRuleVersion); err != nil {
		t.Fatalf("read gatekeeper audit: %v", err)
	}
	if auditResult != "approved" || !bytes.Contains([]byte(auditReason), []byte("FundamentalRuleCheck")) || backtestMetrics == "" || auditedRuleVersion != "v4.0" {
		t.Fatalf("expected explicit gatekeeper checks with active rule, result=%s reason=%q metrics=%q rule=%q", auditResult, auditReason, backtestMetrics, auditedRuleVersion)
	}
	assertCount(t, db, "gatekeeper_audits", 1)
	var gatekeeperNodeEvents int
	if err := db.QueryRow(`SELECT COUNT(*) FROM audit_events WHERE proposal_id='prop_confirm' AND node_name IN ('ProposalLoadNode','FundamentalRuleCheckNode','ConflictCheckNode','BacktestNode','AuditDecisionNode','AuditRecordNode')`).Scan(&gatekeeperNodeEvents); err != nil {
		t.Fatalf("count gatekeeper node events: %v", err)
	}
	if gatekeeperNodeEvents != 6 {
		t.Fatalf("expected gatekeeper node audit events, got %d", gatekeeperNodeEvents)
	}
	assertCount(t, db, "audit_events", 7)

	if _, err := db.Exec(`INSERT INTO rule_effect_validations (validation_id,proposal_id,candidate_rule_version,validation_status,sample_count,sample_window,representativeness_status,overfit_risk,replay_result,guardrail_decision,metrics_json,risk_notes_json,safety_note,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, "val_confirm", "prop_confirm", "draft", "passed", 3, "2026-Q1", "passed", "low", "passed", "passed", `{"hit_count":3}`, `[]`, "只读验证", "2026-01-01T00:00:00Z", "2026-01-01T00:00:00Z"); err != nil {
		t.Fatalf("seed passed validation: %v", err)
	}

	finalReq := httptest.NewRequest(http.MethodPost, "/api/v1/rule-proposals/prop_confirm/final-confirm", bytes.NewBufferString(`{"confirm":true}`))
	finalReq.SetPathValue("proposal_id", "prop_confirm")
	finalReq.Header.Set("X-Request-ID", "req_rule_final_confirm")
	finalW := httptest.NewRecorder()
	app.FinalConfirmRuleProposal(finalW, finalReq)
	if finalW.Code != http.StatusOK {
		t.Fatalf("expected final confirm 200, got %d body=%s", finalW.Code, finalW.Body.String())
	}
	var finalStatus string
	_ = db.QueryRow(`SELECT status FROM rule_proposals WHERE proposal_id='prop_confirm'`).Scan(&finalStatus)
	if finalStatus != "applied" {
		t.Fatalf("expected applied after final confirm, got %s", finalStatus)
	}
}

func TestRuleProposalConfirmFalseRejectsProposalAndWritesAudit(t *testing.T) {
	app, db := testApp(t)
	seedRuleProposal(t, db, "prop_reject", "pending_user_confirm", 3)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rule-proposals/prop_reject/confirm", bytes.NewBufferString(`{"confirm":false,"note":"暂不采用"}`))
	req.SetPathValue("proposal_id", "prop_reject")
	req.Header.Set("X-Request-ID", "req_rule_reject")
	w := httptest.NewRecorder()

	app.ConfirmRuleProposal(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	assertCount(t, db, "gatekeeper_audits", 0)
	assertCount(t, db, "audit_events", 1)
	var status string
	_ = db.QueryRow(`SELECT status FROM rule_proposals WHERE proposal_id='prop_reject'`).Scan(&status)
	if status != "rejected" {
		t.Fatalf("expected rejected, got %s", status)
	}
}

func TestRuleProposalConfirmFalseRejectsInvalidState(t *testing.T) {
	app, db := testApp(t)
	seedRuleProposal(t, db, "prop_draft_reject", "draft", 3)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rule-proposals/prop_draft_reject/confirm", bytes.NewBufferString(`{"confirm":false,"note":"不通过"}`))
	req.SetPathValue("proposal_id", "prop_draft_reject")
	req.Header.Set("X-Request-ID", "req_rule_draft_reject")
	w := httptest.NewRecorder()

	app.ConfirmRuleProposal(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
	assertCount(t, db, "gatekeeper_audits", 0)
	var status string
	_ = db.QueryRow(`SELECT status FROM rule_proposals WHERE proposal_id='prop_draft_reject'`).Scan(&status)
	if status != "draft" {
		t.Fatalf("expected draft, got %s", status)
	}
}

func TestRuleProposalConfirmRejectsTerminalStates(t *testing.T) {
	for _, item := range []struct {
		proposalID string
		status     string
	}{
		{proposalID: "prop_rejected", status: "rejected"},
		{proposalID: "prop_applied", status: "applied"},
	} {
		t.Run(item.status, func(t *testing.T) {
			app, db := testApp(t)
			seedRuleProposal(t, db, item.proposalID, item.status, 3)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/rule-proposals/"+item.proposalID+"/confirm", bytes.NewBufferString(`{"confirm":true}`))
			req.SetPathValue("proposal_id", item.proposalID)
			req.Header.Set("X-Request-ID", "req_"+item.proposalID)
			w := httptest.NewRecorder()

			app.ConfirmRuleProposal(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
			}
			assertCount(t, db, "gatekeeper_audits", 0)
		})
	}
}

func TestRuleProposalConfirmRejectsInvalidState(t *testing.T) {
	app, db := testApp(t)
	seedRuleProposal(t, db, "prop_draft", "draft", 3)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/rule-proposals/prop_draft/confirm", bytes.NewBufferString(`{"confirm":true}`))
	req.SetPathValue("proposal_id", "prop_draft")
	req.Header.Set("X-Request-ID", "req_prop_draft")
	w := httptest.NewRecorder()

	app.ConfirmRuleProposal(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestRuleProposalConfirmRejectsSmallSampleWithoutAudit(t *testing.T) {
	app, db := testApp(t)
	seedRuleProposal(t, db, "prop_small", "pending_user_confirm", 2)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rule-proposals/prop_small/confirm", bytes.NewBuffer(nil))
	req.SetPathValue("proposal_id", "prop_small")
	req.Header.Set("X-Request-ID", "req_rule_small")
	w := httptest.NewRecorder()

	app.ConfirmRuleProposal(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
	assertCount(t, db, "gatekeeper_audits", 0)
}

func TestFinalConfirmFalseRejectsWithoutCreatingRuleVersion(t *testing.T) {
	app, db := testApp(t)
	seedRuleProposal(t, db, "prop_final_reject", "pending_final_confirm", 3)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rule-proposals/prop_final_reject/final-confirm", bytes.NewBufferString(`{"confirm":false,"note":"不同意"}`))
	req.SetPathValue("proposal_id", "prop_final_reject")
	req.Header.Set("X-Request-ID", "req_rule_final_reject")
	w := httptest.NewRecorder()

	app.FinalConfirmRuleProposal(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var active int
	_ = db.QueryRow(`SELECT COUNT(*) FROM rule_versions WHERE status='active'`).Scan(&active)
	if active != 1 {
		t.Fatalf("expected only original active rule, got %d", active)
	}
	var status string
	_ = db.QueryRow(`SELECT status FROM rule_proposals WHERE proposal_id='prop_final_reject'`).Scan(&status)
	if status != "rejected" {
		t.Fatalf("expected rejected, got %s", status)
	}
	assertCount(t, db, "audit_events", 1)
}

func TestFinalConfirmRequiresApprovedGatekeeperAudit(t *testing.T) {
	app, db := testApp(t)
	seedRuleProposal(t, db, "prop_final_no_audit", "pending_final_confirm", 3)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rule-proposals/prop_final_no_audit/final-confirm", bytes.NewBufferString(`{"confirm":true}`))
	req.SetPathValue("proposal_id", "prop_final_no_audit")
	req.Header.Set("X-Request-ID", "req_rule_final_no_audit")
	w := httptest.NewRecorder()

	app.FinalConfirmRuleProposal(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
	var active int
	_ = db.QueryRow(`SELECT COUNT(*) FROM rule_versions WHERE status='active'`).Scan(&active)
	if active != 1 {
		t.Fatalf("should not create new active rule without audit, active=%d", active)
	}
}

func TestFinalConfirmArchivesOldActiveAndCreatesNewActive(t *testing.T) {
	app, db := testApp(t)
	seedRuleProposal(t, db, "prop_final", "pending_final_confirm", 3)
	_, err := db.Exec(`INSERT INTO gatekeeper_audits (gatekeeper_audit_id,proposal_id,audit_result,audit_reason,violates_fundamental_rule,has_rule_conflict,allow_apply,audited_rule_version,created_at) VALUES (?,?,?,?,?,?,?,?,?)`, "gk_final", "prop_final", "approved", "通过", 0, 0, 1, "draft", "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed audit: %v", err)
	}
	seedPassedRuleEffectValidation(t, db, "prop_final", "val_final")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rule-proposals/prop_final/final-confirm", bytes.NewBufferString(`{"confirm":true,"note":"确认应用"}`))
	req.SetPathValue("proposal_id", "prop_final")
	req.Header.Set("X-Request-ID", "req_rule_final")
	w := httptest.NewRecorder()

	app.FinalConfirmRuleProposal(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var active, archived int
	_ = db.QueryRow(`SELECT COUNT(*) FROM rule_versions WHERE status='active'`).Scan(&active)
	_ = db.QueryRow(`SELECT COUNT(*) FROM rule_versions WHERE status='archived'`).Scan(&archived)
	if active != 1 || archived != 1 {
		t.Fatalf("unexpected rule version state: active=%d archived=%d", active, archived)
	}
	var finalConfirmedAt, finalConfirmedNote, appliedRuleVersion string
	_ = db.QueryRow(`SELECT COALESCE(final_confirmed_at,''),COALESCE(final_confirmed_note,''),COALESCE(applied_rule_version,'') FROM rule_proposals WHERE proposal_id='prop_final'`).Scan(&finalConfirmedAt, &finalConfirmedNote, &appliedRuleVersion)
	if finalConfirmedAt == "" || finalConfirmedNote != "确认应用" || appliedRuleVersion != "v_prop_final" {
		t.Fatalf("proposal metadata not persisted: final_confirmed_at=%q note=%q applied=%q", finalConfirmedAt, finalConfirmedNote, appliedRuleVersion)
	}
	assertCount(t, db, "audit_events", 1)
}

func TestFinalConfirmRollsBackWhenRuleVersionSaveFails(t *testing.T) {
	app, db := testApp(t)
	seedRuleProposal(t, db, "prop_final_conflict", "pending_final_confirm", 3)
	_, err := db.Exec(`INSERT INTO gatekeeper_audits (gatekeeper_audit_id,proposal_id,audit_result,audit_reason,violates_fundamental_rule,has_rule_conflict,allow_apply,audited_rule_version,created_at) VALUES (?,?,?,?,?,?,?,?,?)`, "gk_conflict", "prop_final_conflict", "approved", "通过", 0, 0, 1, "draft", "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed audit: %v", err)
	}
	_, err = db.Exec(`INSERT INTO rule_versions (rule_version,status,rules_json,effective_at,created_at) VALUES (?,?,?,?,?)`, "v_prop_final_conflict", "archived", "{}", "2026-01-01T00:00:00Z", "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed conflicting rule version: %v", err)
	}
	seedPassedRuleEffectValidation(t, db, "prop_final_conflict", "val_final_conflict")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rule-proposals/prop_final_conflict/final-confirm", bytes.NewBufferString(`{"confirm":true}`))
	req.SetPathValue("proposal_id", "prop_final_conflict")
	req.Header.Set("X-Request-ID", "req_rule_final_conflict")
	w := httptest.NewRecorder()

	app.FinalConfirmRuleProposal(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d body=%s", w.Code, w.Body.String())
	}
	var active int
	_ = db.QueryRow(`SELECT COUNT(*) FROM rule_versions WHERE status='active'`).Scan(&active)
	if active != 1 {
		t.Fatalf("expected original active rule to remain active, got %d", active)
	}
	var status string
	_ = db.QueryRow(`SELECT status FROM rule_proposals WHERE proposal_id='prop_final_conflict'`).Scan(&status)
	if status != "pending_final_confirm" {
		t.Fatalf("expected proposal status rollback, got %s", status)
	}
	assertCount(t, db, "audit_events", 0)
}

func TestFinalConfirmRejectsTerminalStates(t *testing.T) {
	for _, item := range []struct {
		proposalID string
		status     string
	}{
		{proposalID: "prop_final_rejected", status: "rejected"},
		{proposalID: "prop_final_applied", status: "applied"},
	} {
		t.Run(item.status, func(t *testing.T) {
			app, db := testApp(t)
			seedRuleProposal(t, db, item.proposalID, item.status, 3)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/rule-proposals/"+item.proposalID+"/final-confirm", bytes.NewBufferString(`{"confirm":true}`))
			req.SetPathValue("proposal_id", item.proposalID)
			req.Header.Set("X-Request-ID", "req_"+item.proposalID)
			w := httptest.NewRecorder()

			app.FinalConfirmRuleProposal(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
			}
			assertCount(t, db, "audit_events", 0)
		})
	}
}

func TestFinalConfirmRuleProposalRejectsFailedEffectValidation(t *testing.T) {
	app, db := testApp(t)
	seedRuleProposal(t, db, "prop_failed_validation", "pending_final_confirm", 5)
	if _, err := db.Exec(`INSERT INTO gatekeeper_audits (gatekeeper_audit_id,proposal_id,audit_result,audit_reason,violates_fundamental_rule,has_rule_conflict,allow_apply,audited_rule_version,created_at) VALUES (?,?,?,?,?,?,?,?,?)`, "gk_failed_validation", "prop_failed_validation", "approved", "通过", 0, 0, 1, "draft", "2026-01-01T00:00:00Z"); err != nil {
		t.Fatalf("seed gatekeeper audit: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO rule_effect_validations (validation_id,proposal_id,candidate_rule_version,validation_status,sample_count,sample_window,representativeness_status,overfit_risk,replay_result,guardrail_decision,metrics_json,risk_notes_json,safety_note,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, "val_failed", "prop_failed_validation", "draft", "failed", 5, "2026-Q2", "passed", "high", "failed", "rejected", `{"misjudgment_count":2}`, `["回放不利"]`, "只读验证", "2026-01-01T00:00:00Z", "2026-01-01T00:00:00Z"); err != nil {
		t.Fatalf("seed failed validation: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rule-proposals/prop_failed_validation/final-confirm", bytes.NewBufferString(`{"confirm":true}`))
	req.SetPathValue("proposal_id", "prop_failed_validation")
	req.Header.Set("X-Request-ID", "req_failed_validation")
	w := httptest.NewRecorder()

	app.FinalConfirmRuleProposal(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
	var status string
	_ = db.QueryRow(`SELECT status FROM rule_proposals WHERE proposal_id='prop_failed_validation'`).Scan(&status)
	if status != "pending_final_confirm" {
		t.Fatalf("validation failure must not apply rule, got status %s", status)
	}
}

func TestFinalConfirmRejectsEffectValidationForOldProposalVersion(t *testing.T) {
	app, db := testApp(t)
	seedRuleProposal(t, db, "prop_old_validation", "pending_final_confirm", 5)
	if _, err := db.Exec(`UPDATE rule_proposals SET proposal_version='draft_new' WHERE proposal_id='prop_old_validation'`); err != nil {
		t.Fatalf("update proposal version: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO gatekeeper_audits (gatekeeper_audit_id,proposal_id,audit_result,audit_reason,violates_fundamental_rule,has_rule_conflict,allow_apply,audited_rule_version,created_at) VALUES (?,?,?,?,?,?,?,?,?)`, "gk_old_validation", "prop_old_validation", "approved", "通过", 0, 0, 1, "draft_new", "2026-01-01T00:00:00Z"); err != nil {
		t.Fatalf("seed gatekeeper audit: %v", err)
	}
	seedPassedRuleEffectValidation(t, db, "prop_old_validation", "val_old_validation")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rule-proposals/prop_old_validation/final-confirm", bytes.NewBufferString(`{"confirm":true}`))
	req.SetPathValue("proposal_id", "prop_old_validation")
	w := httptest.NewRecorder()

	app.FinalConfirmRuleProposal(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
}

func seedRuleProposal(t *testing.T, db execer, proposalID, status string, sampleCount int) {
	t.Helper()
	_, err := db.Exec(`INSERT INTO rule_proposals (proposal_id,proposal_type,status,title,proposal_version,before_rule_json,after_rule_json,sample_count,created_at) VALUES (?,?,?,?,?,?,?,?,?)`, proposalID, "risk_rule", status, "测试提案", "draft", "{}", "{\"version\":1}", sampleCount, "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed proposal: %v", err)
	}
}

func seedPassedRuleEffectValidation(t *testing.T, db execer, proposalID, validationID string) {
	t.Helper()
	_, err := db.Exec(`INSERT INTO rule_effect_validations (validation_id,proposal_id,candidate_rule_version,validation_status,sample_count,sample_window,representativeness_status,overfit_risk,replay_result,guardrail_decision,source_explanation_json,metrics_json,risk_notes_json,related_error_cases_json,related_decision_ids_json,related_risk_alert_ids_json,related_audit_event_ids_json,safety_note,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, validationID, proposalID, "draft", "passed", 5, "2026-Q1", "passed", "low", "passed", "passed", `{"source_case_count":3}`, `{"hit_count":5}`, `["本地样本暂未发现不利信号"]`, `["err_1"]`, `["dec_1"]`, `["risk_1"]`, `["audit_1"]`, "只读验证", "2026-01-01T00:00:00Z", "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed passed rule effect validation: %v", err)
	}
}

type execer interface {
	Exec(query string, args ...any) (sql.Result, error)
}
