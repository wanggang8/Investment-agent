package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"investment-agent/internal/application/workflow"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
	"investment-agent/pkg/httputil"
)

func TestDecisionDetailFromWorkflowReflectsFailedStatus(t *testing.T) {
	out := decisionDetailFromWorkflow(workflow.WorkflowContext{DecisionID: "dec_failed", Errors: []string{workflow.ErrCodeEvidenceNotFound}, RuleVerdict: model.RuleVerdict{Status: model.VerdictInsufficientData}})

	if out.WorkflowStatus != string(model.WorkflowFailed) {
		t.Fatalf("expected failed workflow status, got %s", out.WorkflowStatus)
	}
}

func TestConsultDecisionUsesActiveRuleVersion(t *testing.T) {
	app, db := testApp(t)
	seedPortfolioSnapshot(t, db)
	seedMarketSnapshot(t, db, "510300")
	seedCapabilityConfig(t, db, "510300")
	if _, err := db.Exec(`UPDATE rule_versions SET status='archived' WHERE status='active'`); err != nil {
		t.Fatalf("archive seeded rule: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO rule_versions (rule_version,status,rules_json,effective_at,created_at) VALUES (?,?,?,?,?)`, "v4.0", "active", "{}", "2026-01-01T00:00:00Z", "2026-01-01T00:00:00Z"); err != nil {
		t.Fatalf("seed active rule: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/consult", bytes.NewBufferString(`{"symbol":"510300","question":"是否买入","scenario":"hold_review"}`))
	req.Header.Set("X-Request-ID", "req_consult_active_rule")
	w := httptest.NewRecorder()
	app.ConsultDecision(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var ruleVersion string
	if err := db.QueryRow(`SELECT rule_version FROM decision_records WHERE request_id='req_consult_active_rule'`).Scan(&ruleVersion); err != nil {
		t.Fatalf("read decision rule version: %v", err)
	}
	if ruleVersion != "v4.0" {
		t.Fatalf("expected active rule v4.0, got %s", ruleVersion)
	}
}

func TestConsultDecisionReturnsRuleVersionMissingWhenActiveRuleAbsent(t *testing.T) {
	app, db := testApp(t)
	seedPortfolioSnapshot(t, db)
	seedMarketSnapshot(t, db, "510300")
	if _, err := db.Exec(`UPDATE rule_versions SET status='archived' WHERE status='active'`); err != nil {
		t.Fatalf("archive seeded rule: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/consult", bytes.NewBufferString(`{"symbol":"510300","question":"是否持有","scenario":"hold_review"}`))
	req.Header.Set("X-Request-ID", "req_consult_missing_rule")
	w := httptest.NewRecorder()
	app.ConsultDecision(w, req)
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d body=%s", w.Code, w.Body.String())
	}
	assertResponseCode(t, w.Body.Bytes(), string(apperr.CodeRuleVersionMissing))
	assertCount(t, db, "decision_records", 0)
}

func TestConsultDecisionRejectsInvalidScenario(t *testing.T) {
	app, db := testApp(t)
	seedPortfolioSnapshot(t, db)
	seedMarketSnapshot(t, db, "510300")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/consult", bytes.NewBufferString(`{"symbol":"510300","question":"是否持有","scenario":"hold"}`))
	req.Header.Set("X-Request-ID", "req_consult_invalid_scenario")
	w := httptest.NewRecorder()
	app.ConsultDecision(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
	assertResponseCode(t, w.Body.Bytes(), string(apperr.CodeBadRequest))
	assertCount(t, db, "decision_records", 0)
}

func TestConsultDecisionAcceptsContractScenarios(t *testing.T) {
	for _, scenario := range []string{"hold_review", "buy_review", "sell_review", "rebalance_review"} {
		t.Run(scenario, func(t *testing.T) {
			app, db := testApp(t)
			seedPortfolioSnapshot(t, db)
			seedMarketSnapshot(t, db, "510300")
			req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/consult", bytes.NewBufferString(`{"symbol":"510300","question":"是否持有","scenario":"`+scenario+`"}`))
			req.Header.Set("X-Request-ID", "req_consult_"+scenario)
			w := httptest.NewRecorder()
			app.ConsultDecision(w, req)
			if w.Code != http.StatusOK {
				t.Fatalf("expected 200 for scenario %s, got %d body=%s", scenario, w.Code, w.Body.String())
			}
			var persisted string
			if err := db.QueryRow(`SELECT context_snapshot_json FROM decision_records ORDER BY created_at DESC LIMIT 1`).Scan(&persisted); err != nil {
				t.Fatalf("read decision context: %v", err)
			}
			if !strings.Contains(persisted, `"WorkflowType":"consultation"`) {
				t.Fatalf("expected successful consultation context, got %s", persisted)
			}
		})
	}
}

func TestConsultDecisionPersistsConfirmableStatusConsistentWithDetail(t *testing.T) {
	app, db := testApp(t)
	seedPortfolioSnapshot(t, db)
	seedMarketSnapshot(t, db, "510300")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/consult", bytes.NewBufferString(`{"question":"要调仓吗","symbol":"510300"}`))
	req.Header.Set("X-Request-ID", "req_consult_confirmable")
	w := httptest.NewRecorder()

	app.ConsultDecision(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var decisionID, recordType, persistedStatus string
	if err := db.QueryRow(`SELECT decision_id,record_type,confirmation_status FROM decision_records ORDER BY created_at DESC LIMIT 1`).Scan(&decisionID, &recordType, &persistedStatus); err != nil {
		t.Fatalf("read decision: %v", err)
	}
	if recordType != "non_trade_record" || persistedStatus != "not_required" {
		t.Fatalf("expected persisted non-confirmable insufficient-data decision, record_type=%s status=%s", recordType, persistedStatus)
	}
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/decisions/"+decisionID, nil)
	getReq.SetPathValue("decision_id", decisionID)
	getReq.Header.Set("X-Request-ID", "req_get_decision")
	getW := httptest.NewRecorder()
	app.GetDecision(getW, getReq)
	if getW.Code != http.StatusOK {
		t.Fatalf("expected get 200, got %d body=%s", getW.Code, getW.Body.String())
	}
	var body struct {
		Data struct {
			UserConfirmation struct {
				ConfirmationStatus string   `json:"confirmation_status"`
				AvailableActions   []string `json:"available_actions"`
			} `json:"user_confirmation"`
		} `json:"data"`
	}
	if err := json.Unmarshal(getW.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode detail: %v", err)
	}
	if body.Data.UserConfirmation.ConfirmationStatus != persistedStatus || len(body.Data.UserConfirmation.AvailableActions) != 0 {
		t.Fatalf("detail confirmation inconsistent: %+v persisted=%s", body.Data.UserConfirmation, persistedStatus)
	}
}

func TestConsultDecisionRejectsMissingPortfolioInsteadOfFallbackContext(t *testing.T) {
	app, db := testApp(t)
	seedMarketSnapshot(t, db, "510300")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/consult", bytes.NewBufferString(`{"question":"能买吗","symbol":"510300"}`))
	req.Header.Set("X-Request-ID", "req_consult_missing_portfolio")
	w := httptest.NewRecorder()

	app.ConsultDecision(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d body=%s", w.Code, w.Body.String())
	}
	assertResponseCode(t, w.Body.Bytes(), string(apperr.CodeDataRequired))
	assertCount(t, db, "decision_records", 0)
}

func TestConsultDecisionRejectsMissingMarketInsteadOfFallbackMarket(t *testing.T) {
	app, db := testApp(t)
	seedPortfolioSnapshot(t, db)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/consult", bytes.NewBufferString(`{"question":"能买吗","symbol":"510300"}`))
	req.Header.Set("X-Request-ID", "req_consult_missing_market")
	w := httptest.NewRecorder()

	app.ConsultDecision(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d body=%s", w.Code, w.Body.String())
	}
	assertResponseCode(t, w.Body.Bytes(), string(apperr.CodeDataRequired))
	assertCount(t, db, "decision_records", 0)
}

func TestConsultDecisionRejectsExcludedSymbolFromCapabilityConfig(t *testing.T) {
	app, db := testApp(t)
	app.Deps.RetrievalService = staticRetrievalService{}
	seedPortfolioSnapshot(t, db)
	seedMarketSnapshot(t, db, "159915")
	_, err := db.Exec(`INSERT INTO capability_configs (capability_id,symbols_json,excluded_symbols_json,asset_types_json,strategy_scope_json,updated_at) VALUES (?,?,?,?,?,?)`, "cap_excluded", `["510300"]`, `["159915"]`, `[]`, `[]`, "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed capability: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/consult", bytes.NewBufferString(`{"question":"要买入吗","symbol":"159915"}`))
	req.Header.Set("X-Request-ID", "req_consult_excluded")
	w := httptest.NewRecorder()
	app.ConsultDecision(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var status, capabilityStatus, recordType, confirmationStatus string
	if err := db.QueryRow(`SELECT final_verdict_status,capability_status,record_type,confirmation_status FROM decision_records WHERE request_id='req_consult_excluded'`).Scan(&status, &capabilityStatus, &recordType, &confirmationStatus); err != nil {
		t.Fatalf("read decision: %v", err)
	}
	if status != string(model.VerdictRejected) || capabilityStatus != workflow.CapabilityOutOfScope || recordType != "rejection_record" || confirmationStatus != string(model.ConfirmationNotRequired) {
		t.Fatalf("expected excluded non-confirmable decision, status=%q capability=%q record=%q confirmation=%q", status, capabilityStatus, recordType, confirmationStatus)
	}
}

func TestConsultDecisionCanConfirmFormalTradeAdvice(t *testing.T) {
	app, db := testApp(t)
	app.Deps.RetrievalService = staticRetrievalService{}
	seedPortfolioSnapshot(t, db)
	seedMarketSnapshot(t, db, "510300")
	seedFormalEvidence(t, db, "510300")
	seedCapabilityConfig(t, db, "510300")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/consult", bytes.NewBufferString(`{"question":"要买入吗","symbol":"510300"}`))
	req.Header.Set("X-Request-ID", "req_consult_formal")
	w := httptest.NewRecorder()

	app.ConsultDecision(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected consult 200, got %d body=%s", w.Code, w.Body.String())
	}
	var immediate struct {
		Data struct {
			EvidenceChain           []any `json:"evidence_chain"`
			ExpectedReturnScenarios *struct {
				SampleCount int `json:"sample_count"`
			} `json:"expected_return_scenarios"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &immediate); err != nil {
		t.Fatalf("decode immediate detail: %v", err)
	}
	if len(immediate.Data.EvidenceChain) == 0 {
		t.Fatalf("expected immediate consult detail to include evidence chain")
	}
	if immediate.Data.ExpectedReturnScenarios == nil || immediate.Data.ExpectedReturnScenarios.SampleCount == 0 {
		t.Fatalf("expected consult detail to derive expected return sample count from local facts, got %+v", immediate.Data.ExpectedReturnScenarios)
	}
	var decisionID, recordType, persistedStatus string
	if err := db.QueryRow(`SELECT decision_id,record_type,confirmation_status FROM decision_records ORDER BY created_at DESC LIMIT 1`).Scan(&decisionID, &recordType, &persistedStatus); err != nil {
		t.Fatalf("read decision: %v", err)
	}
	if recordType != "formal_trade_advice" || persistedStatus != "pending" {
		t.Fatalf("expected persisted confirmable advice, record_type=%s status=%s", recordType, persistedStatus)
	}
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/decisions/"+decisionID, nil)
	getReq.SetPathValue("decision_id", decisionID)
	getReq.Header.Set("X-Request-ID", "req_get_formal")
	getW := httptest.NewRecorder()
	app.GetDecision(getW, getReq)
	if getW.Code != http.StatusOK {
		t.Fatalf("expected get 200, got %d body=%s", getW.Code, getW.Body.String())
	}
	var detail struct {
		Data struct {
			UserConfirmation struct {
				ConfirmationStatus string   `json:"confirmation_status"`
				AvailableActions   []string `json:"available_actions"`
			} `json:"user_confirmation"`
			TriggeredRules          []any `json:"triggered_rules"`
			EvidenceChain           []any `json:"evidence_chain"`
			AnalystReports          []any `json:"analyst_reports"`
			ArbitrationChain        []any `json:"arbitration_chain"`
			ExpectedReturnScenarios *struct {
				Disclaimer string `json:"disclaimer"`
			} `json:"expected_return_scenarios"`
			AccountSnapshot *struct {
				SnapshotID    string  `json:"snapshot_id"`
				Cash          float64 `json:"cash"`
				TotalAssets   float64 `json:"total_assets"`
				CashRatio     float64 `json:"cash_ratio"`
				HighRiskRatio float64 `json:"high_risk_ratio"`
			} `json:"account_snapshot"`
		} `json:"data"`
	}
	if err := json.Unmarshal(getW.Body.Bytes(), &detail); err != nil {
		t.Fatalf("decode detail: %v", err)
	}
	if detail.Data.UserConfirmation.ConfirmationStatus != "pending" || len(detail.Data.UserConfirmation.AvailableActions) == 0 || len(detail.Data.TriggeredRules) == 0 {
		t.Fatalf("expected confirmable detail with triggered rules, got %+v", detail.Data)
	}
	if len(detail.Data.EvidenceChain) == 0 || len(detail.Data.AnalystReports) == 0 || len(detail.Data.ArbitrationChain) == 0 || detail.Data.ExpectedReturnScenarios == nil || detail.Data.ExpectedReturnScenarios.Disclaimer == "" {
		t.Fatalf("expected replay detail fields, got %+v", detail.Data)
	}
	if detail.Data.AccountSnapshot == nil || detail.Data.AccountSnapshot.SnapshotID != "snap_seed" || detail.Data.AccountSnapshot.Cash != 100 || detail.Data.AccountSnapshot.TotalAssets != 1000 || detail.Data.AccountSnapshot.CashRatio != 0.1 || detail.Data.AccountSnapshot.HighRiskRatio != 0.2 {
		t.Fatalf("expected account snapshot values, got %+v", detail.Data.AccountSnapshot)
	}
	confirmReq := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/"+decisionID+"/confirmations", bytes.NewBufferString(`{"confirmation_type":"planned","note":"按计划观察"}`))
	confirmReq.SetPathValue("decision_id", decisionID)
	confirmReq.Header.Set("X-Request-ID", "req_confirm_formal")
	confirmW := httptest.NewRecorder()
	app.CreateConfirmation(confirmW, confirmReq)
	if confirmW.Code != http.StatusOK {
		t.Fatalf("expected confirm 200, got %d body=%s", confirmW.Code, confirmW.Body.String())
	}
	var finalStatus string
	if err := db.QueryRow(`SELECT confirmation_status FROM decision_records WHERE decision_id=?`, decisionID).Scan(&finalStatus); err != nil {
		t.Fatalf("read final status: %v", err)
	}
	if finalStatus != "planned" {
		t.Fatalf("expected planned confirmation, got %s", finalStatus)
	}
}

func TestListDecisionsFiltersByConfirmationStatusAndDateRange(t *testing.T) {
	app, db := testApp(t)
	if _, err := db.Exec(`INSERT INTO decision_records (decision_id,request_id,workflow_type,symbol,workflow_status,record_type,dashboard_state,final_verdict_status,final_verdict_text,confirmation_status,rule_version,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`, "decision_pending_old", "req_pending_old", "consultation", "510300", "completed", "formal_trade_advice", "normal", "hold", "持有", "pending", "v3.0", "2026-01-01T00:00:00Z"); err != nil {
		t.Fatalf("seed pending old decision: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO decision_records (decision_id,request_id,workflow_type,symbol,workflow_status,record_type,dashboard_state,final_verdict_status,final_verdict_text,confirmation_status,rule_version,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`, "decision_planned", "req_planned", "consultation", "159915", "completed", "formal_trade_advice", "normal", "hold", "持有", "planned", "v3.0", "2026-01-03T00:00:00Z"); err != nil {
		t.Fatalf("seed planned decision: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO decision_records (decision_id,request_id,workflow_type,symbol,workflow_status,record_type,dashboard_state,final_verdict_status,final_verdict_text,confirmation_status,rule_version,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`, "decision_pending_new", "req_pending_new", "consultation", "510500", "completed", "formal_trade_advice", "normal", "hold", "持有", "pending", "v3.0", "2026-01-04T00:00:00Z"); err != nil {
		t.Fatalf("seed pending new decision: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/decisions?status=pending&from=2026-01-02&to=2026-01-05", nil)
	req.Header.Set("X-Request-ID", "req_decision_filters")
	w := httptest.NewRecorder()
	app.ListDecisions(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data struct {
			Items []struct {
				DecisionID         string `json:"decision_id"`
				ConfirmationStatus string `json:"confirmation_status"`
			} `json:"items"`
			Total int `json:"total"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.Total != 1 || len(body.Data.Items) != 1 || body.Data.Items[0].DecisionID != "decision_pending_new" || body.Data.Items[0].ConfirmationStatus != "pending" {
		t.Fatalf("expected filtered pending decision, got %+v", body.Data)
	}
}

func TestListDecisionsRejectsInvalidDateRange(t *testing.T) {
	app, _ := testApp(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/decisions?from=2026-01-05&to=2026-01-02", nil)
	req.Header.Set("X-Request-ID", "req_decision_bad_range")
	w := httptest.NewRecorder()

	app.ListDecisions(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestDecisionDetailExpectedReturnUsesStoredSampleCount(t *testing.T) {
	record := repository.DecisionRecord{DecisionID: "decision_expected", RequestID: "req_expected", WorkflowType: "consultation", Symbol: "510300", WorkflowStatus: "completed", RecordType: "formal_trade_advice", DashboardState: "normal", FinalVerdictStatus: "hold", FinalVerdictText: "持有", ConfirmationStatus: "pending", ExpectedReturnScenariosJSON: `{"precision_status":"available","reason":"样本足够","sample_count":20,"sample_window":"2024-2026","screening_condition":"PE 30%-50%","sell_evaluation":{"status":"triggered","triggers":["base_upper_bound_exceeded"],"prompts":["评估分批止盈"],"actions":["评估分批止盈"],"non_trading_disclaimer":"仅人工评估"},"reassessment_trigger":{"reason":"基准下移","boundary":"base_midpoint_downshift","current_value":0.03},"scenarios":[{"name":"base","return_rate":0.03,"return_range":"3.00%","probability":0.5,"trigger":"估值维持"}]}`}

	detail := decisionDetailFromRecord(record, nil, nil)

	if detail.ExpectedReturnScenarios == nil || detail.ExpectedReturnScenarios.SampleCount != 20 || len(detail.ExpectedReturnScenarios.Scenarios) != 1 {
		t.Fatalf("expected stored sample_count independent of scenario count, got %+v", detail.ExpectedReturnScenarios)
	}
	if detail.ExpectedReturnScenarios.SampleWindow != "2024-2026" || detail.ExpectedReturnScenarios.ScreeningCondition != "PE 30%-50%" || detail.ExpectedReturnScenarios.Scenarios[0].Trigger != "估值维持" {
		t.Fatalf("expected stored sample context, got %+v", detail.ExpectedReturnScenarios)
	}
	if detail.ExpectedReturnScenarios.SellEvaluation == nil || detail.ExpectedReturnScenarios.SellEvaluation.Status != "triggered" || detail.ExpectedReturnScenarios.ReassessmentTrigger == nil {
		t.Fatalf("expected stored sell evaluation, got %+v", detail.ExpectedReturnScenarios)
	}
}

func TestDecisionDetailExpectedReturnReadsAPIShapedHistoricalJSON(t *testing.T) {
	record := repository.DecisionRecord{DecisionID: "decision_expected_api", RequestID: "req_expected_api", WorkflowType: "consultation", Symbol: "510300", WorkflowStatus: "completed", RecordType: "formal_trade_advice", DashboardState: "normal", FinalVerdictStatus: "hold", FinalVerdictText: "持有", ConfirmationStatus: "pending", ExpectedReturnScenariosJSON: `{"precision_status":"available","sample_count":20,"sample_window":"2024-2026","screening_condition":"PE 30%-50%","scenarios":[{"scenario":"base","return_range":"3.00%","probability":0.5,"trigger":"估值维持"}]}`}

	detail := decisionDetailFromRecord(record, nil, nil)

	if detail.ExpectedReturnScenarios == nil || len(detail.ExpectedReturnScenarios.Scenarios) != 1 {
		t.Fatalf("expected API-shaped expected return scenarios, got %+v", detail.ExpectedReturnScenarios)
	}
	got := detail.ExpectedReturnScenarios.Scenarios[0]
	if got.Scenario != "base" || got.ReturnRange != "3.00%" || got.Trigger != "估值维持" || got.Probability == nil || *got.Probability != 0.5 {
		t.Fatalf("expected API-shaped scenario fields, got %+v", got)
	}
}

func TestDecisionDetailExpectedReturnReadsWorkflowStoredJSONShape(t *testing.T) {
	record := repository.DecisionRecord{DecisionID: "decision_expected_workflow", RequestID: "req_expected_workflow", WorkflowType: "consultation", Symbol: "510300", WorkflowStatus: "completed", RecordType: "formal_trade_advice", DashboardState: "normal", FinalVerdictStatus: "hold", FinalVerdictText: "持有", ConfirmationStatus: "pending", ExpectedReturnScenariosJSON: `{"precision_status":"available","sample_count":20,"sample_window":"最近可比样本","screening_condition":"本地样本","scenarios":[{"Name":"base","ReturnRate":0.03,"ReturnRange":"0.00%~8.00%","Probability":0.5,"Trigger":"估值维持"}]}`}

	detail := decisionDetailFromRecord(record, nil, nil)

	if detail.ExpectedReturnScenarios == nil || len(detail.ExpectedReturnScenarios.Scenarios) != 1 {
		t.Fatalf("expected workflow-shaped expected return scenarios, got %+v", detail.ExpectedReturnScenarios)
	}
	got := detail.ExpectedReturnScenarios.Scenarios[0]
	if got.Scenario != "base" || got.ReturnRange != "0.00%~8.00%" || got.Trigger != "估值维持" || got.Probability == nil || *got.Probability != 0.5 {
		t.Fatalf("expected workflow-shaped scenario fields, got %+v", got)
	}
}

func TestDecisionDetailExpectedReturnKeepsSellEvaluationOnlyJSON(t *testing.T) {
	record := repository.DecisionRecord{DecisionID: "decision_expected_sell_only", RequestID: "req_expected_sell_only", WorkflowType: "consultation", Symbol: "510300", WorkflowStatus: "completed", RecordType: "formal_trade_advice", DashboardState: "normal", FinalVerdictStatus: "hold", FinalVerdictText: "持有", ConfirmationStatus: "pending", ExpectedReturnScenariosJSON: `{"sell_evaluation":{"status":"not_applicable","prompts":["缺少持仓成本"],"non_trading_disclaimer":"仅人工评估"}}`}

	detail := decisionDetailFromRecord(record, nil, nil)

	if detail.ExpectedReturnScenarios == nil || detail.ExpectedReturnScenarios.SellEvaluation == nil {
		t.Fatalf("expected sell evaluation to be retained, got %+v", detail.ExpectedReturnScenarios)
	}
	if detail.ExpectedReturnScenarios.SellEvaluation.Status != "not_applicable" || len(detail.ExpectedReturnScenarios.SellEvaluation.Prompts) != 1 {
		t.Fatalf("expected sell evaluation fields, got %+v", detail.ExpectedReturnScenarios.SellEvaluation)
	}
}

func TestDecisionDetailFromRecordIncludesRetrievalQualitySnapshot(t *testing.T) {
	record := repository.DecisionRecord{DecisionID: "decision_retrieval_quality", RequestID: "req_retrieval_quality", WorkflowType: "consultation", Symbol: "510300", WorkflowStatus: "completed", RecordType: "formal_trade_advice", DashboardState: "normal", FinalVerdictStatus: "hold", FinalVerdictText: "持有", ConfirmationStatus: "pending", ContextSnapshotJSON: `{"WorkflowType":"consultation","retrieval_quality_summary":{"query_summary":"510300","top_k":2,"status":"degraded","index_health":"missing","index_freshness":"unknown","fallback_source":"sqlite_summary","source_consistency_status":"checked","degraded_reason":"veclite index not configured"}}`}

	detail := decisionDetailFromRecord(record, nil, nil)

	if detail.RetrievalQuality == nil || detail.RetrievalQuality.TopK != 2 || detail.RetrievalQuality.FallbackSource != "sqlite_summary" || detail.RetrievalQuality.DegradedReason != "veclite index not configured" {
		t.Fatalf("expected retrieval quality restored from context snapshot, got %+v", detail.RetrievalQuality)
	}
}

func TestDecisionDetailFromWorkflowExpectedReturnUsesWorkflowSampleCount(t *testing.T) {
	probability := 0.5
	wf := workflow.WorkflowContext{DecisionID: "decision_wf_expected", RequestID: "req_wf_expected", Symbol: "510300", RuleVerdict: model.RuleVerdict{Status: model.VerdictHold, Text: "持有"}, MarketSnapshot: model.MarketSnapshot{Symbol: "510300", TradeDate: "2026-06-20", ClosePrice: 4.23, PEPercentile: 31, PBPercentile: 27}, ExpectedReturnSampleCount: 20, ExpectedReturnPrecisionStatus: model.PrecisionAvailable, ExpectedReturnReason: "样本足够", ExpectedReturnSampleWindow: "2024-2026", ExpectedReturnScreeningCondition: "PE 30%-50%", ExpectedReturnSellEvaluation: workflow.ExpectedReturnSellEvaluation{Status: "triggered", Triggers: []string{"base_upper_bound_exceeded"}, Prompts: []string{"评估分批止盈"}, Actions: []string{"评估分批止盈"}, NonTradingDisclaimer: "仅人工评估"}, ExpectedReturnReassessmentTrigger: workflow.ExpectedReturnReassessmentTrigger{Reason: "基准下移", Boundary: "base_midpoint_downshift", CurrentValue: 0.03}, ExpectedReturnScenarios: []workflow.ExpectedReturnScenario{{Name: "base", ReturnRate: 0.03, ReturnRange: "3.00%", Probability: &probability, Trigger: "估值维持"}}}

	detail := decisionDetailFromWorkflow(wf)

	if detail.ExpectedReturnScenarios == nil || detail.ExpectedReturnScenarios.SampleCount != 20 || len(detail.ExpectedReturnScenarios.Scenarios) != 1 {
		t.Fatalf("expected workflow sample count independent of scenario count, got %+v", detail.ExpectedReturnScenarios)
	}
	if detail.ExpectedReturnScenarios.SampleWindow != "2024-2026" || detail.ExpectedReturnScenarios.Scenarios[0].Trigger != "估值维持" || detail.ExpectedReturnScenarios.SellEvaluation == nil {
		t.Fatalf("expected workflow expected return context, got %+v", detail.ExpectedReturnScenarios)
	}
	if detail.MarketContext == nil || detail.MarketContext.CurrentPrice != 4.23 || detail.MarketContext.PEPercentile != 31 || detail.MarketContext.PBPercentile != 27 || detail.MarketContext.TradeDate != "2026-06-20" {
		t.Fatalf("expected market context for expected-return report fields, got %+v", detail.MarketContext)
	}
}

func TestDecisionDetailFromRecordRestoresMarketContextSnapshot(t *testing.T) {
	record := repository.DecisionRecord{DecisionID: "decision_market_context", RequestID: "req_market_context", WorkflowType: "consultation", Symbol: "510300", WorkflowStatus: "completed", RecordType: "formal_trade_advice", DashboardState: "normal", FinalVerdictStatus: "hold", FinalVerdictText: "持有", ConfirmationStatus: "pending", ContextSnapshotJSON: `{"MarketSnapshot":{"Symbol":"510300","TradeDate":"2026-06-20","ClosePrice":4.23,"PEPercentile":31,"PBPercentile":27}}`}

	detail := decisionDetailFromRecord(record, nil, nil)

	if detail.MarketContext == nil || detail.MarketContext.Symbol != "510300" || detail.MarketContext.TradeDate != "2026-06-20" || detail.MarketContext.CurrentPrice != 4.23 || detail.MarketContext.PEPercentile != 31 || detail.MarketContext.PBPercentile != 27 {
		t.Fatalf("expected market context restored from context snapshot, got %+v", detail.MarketContext)
	}
}

func TestDecisionDetailExpectedReturnSerializesEmptyScenariosAsArray(t *testing.T) {
	detail := decisionDetailFromWorkflow(workflow.WorkflowContext{DecisionID: "decision_empty_expected", RuleVerdict: model.RuleVerdict{Status: model.VerdictHold, Text: "持有"}, ExpectedReturnSampleCount: 3, ExpectedReturnPrecisionStatus: model.PrecisionUnavailable, ExpectedReturnReason: "样本过少"})

	body, err := json.Marshal(detail)
	if err != nil {
		t.Fatalf("marshal detail: %v", err)
	}
	if !strings.Contains(string(body), `"scenarios":[]`) {
		t.Fatalf("expected scenarios to serialize as empty array, got %s", body)
	}
}

func TestDecisionDetailAnalystReportsSerializeEmptyFieldsAsArrays(t *testing.T) {
	detail := decisionDetailFromWorkflow(workflow.WorkflowContext{DecisionID: "decision_reports", RuleVerdict: model.RuleVerdict{Status: model.VerdictHold, Text: "持有"}, AnalystReports: map[string]string{"risk": "谨慎持有"}})

	body, err := json.Marshal(detail)
	if err != nil {
		t.Fatalf("marshal detail: %v", err)
	}
	for _, want := range []string{`"key_reasons":[]`, `"risk_warnings":[]`, `"evidence_ids":[]`} {
		if !strings.Contains(string(body), want) {
			t.Fatalf("expected analyst report arrays to be non-null, missing %s in %s", want, body)
		}
	}
}

func TestDecisionDetailFromWorkflowIncludesAnalystMetadata(t *testing.T) {
	detail := decisionDetailFromWorkflow(workflow.WorkflowContext{
		DecisionID:     "decision_report_metadata",
		RuleVerdict:    model.RuleVerdict{Status: model.VerdictHold, Text: "持有"},
		AnalystReports: map[string]string{"value": "估值材料"},
		AnalystReportMetadata: map[string]map[string]string{
			"value": {
				"prompt_version": "p37-analyst-v1",
				"model":          "gpt-5.4-mini",
				"input_summary":  "value 510300",
				"output_summary": "估值材料",
				"parse_status":   "parsed",
				"quality_status": "passed",
			},
		},
	})

	if len(detail.AnalystReports) != 1 {
		t.Fatalf("expected one analyst report, got %+v", detail.AnalystReports)
	}
	report := detail.AnalystReports[0]
	if report.PromptVersion != "p37-analyst-v1" || report.Model != "gpt-5.4-mini" || report.ParseStatus != "parsed" || report.QualityStatus != "passed" {
		t.Fatalf("expected analyst metadata in immediate response, got %+v", report)
	}
}

func TestDecisionDetailFromWorkflowIncludesRetrievalQualitySummary(t *testing.T) {
	detail := decisionDetailFromWorkflow(workflow.WorkflowContext{
		DecisionID:              "decision_retrieval_quality",
		RuleVerdict:             model.RuleVerdict{Status: model.VerdictHold, Text: "持有"},
		RetrievalQualitySummary: workflow.RetrievalQualitySummary{QuerySummary: "510300", TopK: 2, Status: "degraded", IndexHealth: "missing", FallbackSource: "sqlite_summary", SourceConsistencyStatus: "checked", DegradedReason: "veclite index not configured"},
	})

	if detail.RetrievalQuality == nil {
		t.Fatalf("expected retrieval quality summary in immediate response")
	}
	if detail.RetrievalQuality.TopK != 2 || detail.RetrievalQuality.FallbackSource != "sqlite_summary" || detail.RetrievalQuality.DegradedReason != "veclite index not configured" {
		t.Fatalf("unexpected retrieval quality summary: %+v", detail.RetrievalQuality)
	}
}

func TestConsultDecisionAcceptsExpectedReturnDynamicInputs(t *testing.T) {
	app, db := testApp(t)
	app.Deps.RetrievalService = staticRetrievalService{}
	seedPortfolioSnapshot(t, db)
	seedMarketSnapshot(t, db, "510300")
	if _, err := db.Exec(`UPDATE market_snapshots SET market_metrics_json=? WHERE market_snapshot_id='market_seed'`, `{"metadata":{"nav_history":[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18]}}`); err != nil {
		t.Fatalf("seed market metrics: %v", err)
	}
	seedFormalEvidence(t, db, "510300")
	seedCapabilityConfig(t, db, "510300")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/decisions/consult", bytes.NewBufferString(`{"question":"要止盈吗","symbol":"510300","expected_return_previous_base_midpoint":0.20,"expected_return_target_return_rate":0.15}`))
	req.Header.Set("X-Request-ID", "req_consult_dynamic_expected")
	w := httptest.NewRecorder()

	app.ConsultDecision(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var record repository.DecisionRecord
	if err := db.QueryRow(`SELECT expected_return_scenarios_json FROM decision_records WHERE request_id='req_consult_dynamic_expected'`).Scan(&record.ExpectedReturnScenariosJSON); err != nil {
		t.Fatalf("read expected return json: %v", err)
	}
	if !strings.Contains(record.ExpectedReturnScenariosJSON, "target_return_reached") || !strings.Contains(record.ExpectedReturnScenariosJSON, "base_midpoint_downshift") {
		t.Fatalf("expected dynamic sell triggers in persisted expected return json, got %s", record.ExpectedReturnScenariosJSON)
	}
}

func TestP89ExpectedReturnFromJSONReadsHistoricalContexts(t *testing.T) {
	out := expectedReturnFromJSON(`{"precision_status":"available","sample_count":20,"scenarios":[],"historical_contexts":[{"label":"极端恐惧样本","window":"2018Q4, 2020Q1, 2022Q4","sample_count":20,"outcome":"暂停主动交易建议","max_drawdown":-0.18,"recovery":"3-9 个月","source":"local_public_history"}]}`)

	if out == nil || len(out.HistoricalContexts) != 1 {
		t.Fatalf("expected historical contexts readback, got %+v", out)
	}
	if out.HistoricalContexts[0].Label != "极端恐惧样本" || out.HistoricalContexts[0].MaxDrawdown != -0.18 {
		t.Fatalf("unexpected historical context readback: %+v", out.HistoricalContexts)
	}
}

func seedPortfolioSnapshot(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec(`INSERT INTO portfolio_snapshots (snapshot_id,snapshot_time,cash,total_assets,cash_ratio,high_risk_ratio,position_count,source,created_at) VALUES (?,?,?,?,?,?,?,?,?)`, "snap_seed", "2026-01-01T00:00:00Z", 100, 1000, 0.1, 0.2, 1, "manual", "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed portfolio: %v", err)
	}
	_, err = db.Exec(`INSERT INTO position_snapshots (position_snapshot_id,snapshot_id,position_id,symbol,name,quantity,cost_price,current_price,market_value,unrealized_profit_ratio,position_state,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`, "ps_seed", "snap_seed", "pos_seed", "510300", "沪深300", 10, 2, 3, 30, 0.5, "normal", "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed position snapshot: %v", err)
	}
}

func seedMarketSnapshot(t *testing.T, db *sql.DB, symbol string) {
	t.Helper()
	_, err := db.Exec(`INSERT INTO market_snapshots (market_snapshot_id,symbol,trade_date,pe_percentile,pb_percentile,liquidity_state,sentiment_state,volume_percentile,volatility_percentile,created_at) VALUES (?,?,?,?,?,?,?,?,?,?)`, "market_seed", symbol, "2026-01-01", 20, 20, "normal", "neutral", 50, 20, "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed market: %v", err)
	}
}

func seedFormalEvidence(t *testing.T, db *sql.DB, symbol string) {
	t.Helper()
	_, err := db.Exec(`INSERT INTO intelligence_items (intelligence_id,source_name,source_level,original_url,published_at,captured_at,content_hash,raw_title,raw_text_ref,created_at) VALUES (?,?,?,?,?,?,?,?,?,?)`, "intel_formal", "交易所公告", "A", "https://example.com/formal", "2026-01-01T00:00:00Z", "2026-01-01T01:00:00Z", "hash_formal", "标题", "raw", "2026-01-01T01:00:00Z")
	if err != nil {
		t.Fatalf("seed intelligence item: %v", err)
	}
	_, err = db.Exec(`INSERT INTO intelligence_summary (summary_id,intelligence_id,symbol,entity,source_level,evidence_role,event_type,summary,time_weight,relevance_score,verification_group_id,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`, "sum_formal", "intel_formal", symbol, symbol, "A", "formal", "normal", "正式证据", 1, 1, "vg_formal", "2026-01-01T01:00:00Z")
	if err != nil {
		t.Fatalf("seed summary: %v", err)
	}
	_, err = db.Exec(`INSERT INTO source_verifications (verification_id,verification_group_id,event_id,symbol,event_type,evidence_role,verification_status,independent_source_count,highest_source_level,latest_published_at,evidence_ids_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`, "ver_formal", "vg_formal", "event_formal", symbol, "normal", "formal", "satisfied", 2, "A", "2026-01-01T00:00:00Z", `["sum_formal"]`, "2026-01-01T01:00:00Z")
	if err != nil {
		t.Fatalf("seed source verification: %v", err)
	}
}

func seedCapabilityConfig(t *testing.T, db *sql.DB, symbol string) {
	t.Helper()
	_, err := db.Exec(`INSERT INTO capability_configs (capability_id,symbols_json,excluded_symbols_json,asset_types_json,strategy_scope_json,updated_at) VALUES (?,?,?,?,?,?)`, "cap_formal", `["`+symbol+`"]`, `[]`, `[]`, `[]`, "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed capability: %v", err)
	}
}

type staticRetrievalService struct{}

func (staticRetrievalService) RetrieveEvidence(ctx context.Context, req workflow.RetrievalRequest) (workflow.RetrievalResult, error) {
	return workflow.RetrievalResult{EvidenceSet: model.EvidenceSet{VerificationStatus: model.VerificationSatisfied, Items: []model.Evidence{{EvidenceID: "sum_formal", SummaryID: "sum_formal", SourceName: "交易所公告", SourceLevel: model.SourceLevelA, Role: model.EvidenceFormal, EventType: model.EventNormal, IndependentSourceCount: 2, Summary: "正式证据"}}}, OutputRef: "sum_formal"}, nil
}

func assertResponseCode(t *testing.T, body []byte, want string) {
	t.Helper()
	var envelope httputil.Envelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if envelope.Error == nil || envelope.Error.Code != want {
		t.Fatalf("expected error code %s, got %+v", want, envelope.Error)
	}
}
