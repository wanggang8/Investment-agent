package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"investment-agent/internal/application/service"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/infrastructure/config"
)

const dailyDisciplineSafetyNote = "每日纪律报告只用于本地记录和人工复核，不会自动执行交易。"

func TestGetTodayDailyDisciplineReportUsesConfiguredAutoRunTimezone(t *testing.T) {
	app, db := testAppWithDailyAutoRunTimezone(t, "Asia/Shanghai")
	_, err := db.Exec(`INSERT INTO daily_auto_run_states (run_id,idempotency_key,local_date,scope,symbol_set_hash,status,last_run_at,failure_code,failure_reason,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
		"run_tz", "key_tz", "2026-06-08", "holdings", "hash_tz", "failed", "2026-06-07T16:30:00Z", "missing_prerequisites", "缺少本地账户或持仓", "2026-06-07T16:30:00Z", "2026-06-07T16:30:00Z")
	if err != nil {
		t.Fatalf("seed auto run: %v", err)
	}

	out, err := app.QuerySvc.TodayDailyDisciplineReport(testContext(t), time.Date(2026, 6, 7, 16, 30, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("TodayDailyDisciplineReport: %v", err)
	}
	if out.LocalDate != "2026-06-08" || out.Status != "insufficient_data" {
		t.Fatalf("expected configured timezone local date insufficient data, got local_date=%q status=%q", out.LocalDate, out.Status)
	}
}

func TestGetTodayDailyDisciplineReportSynthesizesInsufficientDataFromMissingPrerequisitesAutoRun(t *testing.T) {
	app, db := testAppWithDailyAutoRunTimezone(t, "UTC")
	today := time.Now().UTC().Format(time.DateOnly)
	_, err := db.Exec(`INSERT INTO daily_auto_run_states (run_id,idempotency_key,local_date,scope,symbol_set_hash,status,last_run_at,failure_code,failure_reason,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
		"run_missing", "key_missing", today, "holdings", "hash_missing", "failed", "2026-06-08T01:00:00Z", "missing_prerequisites", "缺少本地账户或持仓", "2026-06-08T01:00:00Z", "2026-06-08T01:00:00Z")
	if err != nil {
		t.Fatalf("seed auto run: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/daily-discipline/reports/today", nil)
	req.Header.Set("X-Request-ID", "req_today_missing")
	w := httptest.NewRecorder()

	app.GetTodayDailyDisciplineReport(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	data := decodeDailyDisciplineReportData(t, w)
	if data.Status != "insufficient_data" || data.FailureReason != "缺少本地账户或持仓" {
		t.Fatalf("unexpected status/failure: status=%q reason=%q", data.Status, data.FailureReason)
	}
	if data.MissingAction != "补齐本地账户与持仓后再生成每日纪律报告。" {
		t.Fatalf("unexpected missing action: %q", data.MissingAction)
	}
	if !containsStringLocal(data.MissingCategories, "account") || !containsStringLocal(data.MissingCategories, "holdings") {
		t.Fatalf("expected missing prerequisite categories, got %#v", data.MissingCategories)
	}
	if data.SafetyNote != dailyDisciplineSafetyNote {
		t.Fatalf("unexpected safety note: %q", data.SafetyNote)
	}
	if data.AutoRunLink != "/daily-auto-run" || data.AuditLink != "/audit?input_ref=key_missing" || data.NotificationLink != "/notifications?source_id=key_missing" {
		t.Fatalf("unexpected links: auto=%q audit=%q notification=%q", data.AutoRunLink, data.AuditLink, data.NotificationLink)
	}
	assertCount(t, db, "decision_records", 0)
}

func TestListDailyDisciplineReportsReturnsNewestFirstAndCapsLimit(t *testing.T) {
	app, db := testApp(t)
	seedDailyDisciplineReport(t, db, repository.DailyDisciplineReport{ReportID: "report_old", LocalDate: "2026-06-06", Scope: "holdings", SymbolSetHash: "hash_old", SourceType: "manual", Status: "success", Summary: "old", CreatedAt: "2026-06-06T00:00:00Z", UpdatedAt: "2026-06-06T00:00:00Z"})
	seedDailyDisciplineReport(t, db, repository.DailyDisciplineReport{ReportID: "report_new", LocalDate: "2026-06-08", Scope: "holdings", SymbolSetHash: "hash_new", SourceType: "manual", Status: "failed", Summary: "new", CreatedAt: "2026-06-08T00:00:00Z", UpdatedAt: "2026-06-08T00:00:00Z"})
	seedDailyDisciplineReport(t, db, repository.DailyDisciplineReport{ReportID: "report_mid", LocalDate: "2026-06-07", Scope: "holdings", SymbolSetHash: "hash_mid", SourceType: "manual", Status: "success", Summary: "mid", CreatedAt: "2026-06-07T00:00:00Z", UpdatedAt: "2026-06-07T00:00:00Z"})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/daily-discipline/reports?limit=999", nil)
	app.ListDailyDisciplineReports(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	list := decodeDailyDisciplineReportListData(t, w)
	if len(list.Reports) != 3 {
		t.Fatalf("expected all 3 reports under cap, got %d", len(list.Reports))
	}
	if list.Reports[0].ReportID != "report_new" || list.Reports[1].ReportID != "report_mid" || list.Reports[2].ReportID != "report_old" {
		t.Fatalf("reports not newest first: %#v", list.Reports)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/daily-discipline/reports?status=success", nil)
	app.ListDailyDisciplineReports(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	list = decodeDailyDisciplineReportListData(t, w)
	if len(list.Reports) != 2 || list.Reports[0].ReportID != "report_mid" || list.Reports[1].ReportID != "report_old" {
		t.Fatalf("status/default limit not honored: %#v", list.Reports)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/daily-discipline/reports?limit=bad", nil)
	app.ListDailyDisciplineReports(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid limit, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetDailyDisciplineReportUnknownReturns404(t *testing.T) {
	app, _ := testApp(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/daily-discipline/reports/report_missing", nil)
	req.SetPathValue("report_id", "report_missing")
	w := httptest.NewRecorder()

	app.GetDailyDisciplineReport(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetDailyDisciplineReportLinkedToDecisionIncludesVerdictEvidenceAndLinks(t *testing.T) {
	app, db := testApp(t)
	seedDecisionWithEvidence(t, db, "decision_report")
	seedDailyDisciplineReport(t, db, repository.DailyDisciplineReport{ReportID: "report_decision", LocalDate: "2026-06-08", Scope: "holdings", SymbolSetHash: "hash_decision", SourceType: "auto_run", SourceID: "key_decision", DecisionID: "decision_report", Status: "success", Summary: "已生成", CreatedAt: "2026-06-08T02:00:00Z", UpdatedAt: "2026-06-08T02:00:00Z"})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/daily-discipline/reports/report_decision", nil)
	req.SetPathValue("report_id", "report_decision")
	w := httptest.NewRecorder()

	app.GetDailyDisciplineReport(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	data := decodeDailyDisciplineReportData(t, w)
	if data.FinalVerdict != "持有" || data.VerdictStatus != "hold" {
		t.Fatalf("unexpected verdict: text=%q status=%q", data.FinalVerdict, data.VerdictStatus)
	}
	if data.Evidence.EvidenceCount != 2 || data.Evidence.IndependentSourceCount != 3 || data.Evidence.HighGradeIndependentSourceCount != 2 {
		t.Fatalf("unexpected evidence: %#v", data.Evidence)
	}
	if data.DecisionLink != "/decisions/decision_report" || data.AutoRunLink != "/daily-auto-run" || data.AuditLink != "/audit?input_ref=key_decision" || data.NotificationLink != "/notifications?source_id=key_decision" {
		t.Fatalf("unexpected links: %#v", data)
	}
	if data.Trend.SuccessCount != 1 {
		t.Fatalf("expected trend success count from reports, got %#v", data.Trend)
	}
}

func TestGetDailyDisciplineReportIncludesRiskAlertSummary(t *testing.T) {
	app, db := testApp(t)
	seedDecisionWithEvidence(t, db, "decision_risk_report")
	seedDailyDisciplineReport(t, db, repository.DailyDisciplineReport{ReportID: "report_risk", LocalDate: "2026-06-08", Scope: "holdings", SymbolSetHash: "hash_risk", SourceType: "auto_run", SourceID: "key_risk", DecisionID: "decision_risk_report", Status: "degraded", Summary: "已生成，存在风险预警", CreatedAt: "2026-06-08T02:00:00Z", UpdatedAt: "2026-06-08T02:00:00Z"})
	_, err := db.Exec(`INSERT INTO risk_alerts (alert_id,risk_type,severity,sop_status,symbol,trigger_summary,trigger_context_json,prohibited_actions_json,suggested_actions_json,related_decision_id,related_report_id,related_notification_id,related_audit_event_id,last_triggered_at,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, "risk_report_1", "valuation_high", "warning", "active", "510300", "PE 分位高于 80%", `{"pe_percentile":88}`, `["新增买入"]`, `["人工复核分批止盈"]`, "decision_risk_report", "report_risk", "notif_risk", "audit_risk", "2026-06-08T02:00:00Z", "2026-06-08T02:00:00Z", "2026-06-08T02:00:00Z")
	if err != nil {
		t.Fatalf("seed risk alert: %v", err)
	}
	_, err = db.Exec(`INSERT INTO risk_alerts (alert_id,risk_type,severity,sop_status,symbol,trigger_summary,related_decision_id,related_report_id,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?)`, "risk_report_archived", "data_degraded", "warning", "archived", "510300", "历史风险已归档", "decision_risk_report", "report_risk", "2026-06-08T01:00:00Z", "2026-06-08T01:00:00Z")
	if err != nil {
		t.Fatalf("seed risk alert: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/daily-discipline/reports/report_risk", nil)
	req.SetPathValue("report_id", "report_risk")
	w := httptest.NewRecorder()

	app.GetDailyDisciplineReport(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	data := decodeDailyDisciplineReportData(t, w)
	if len(data.RiskAlerts) != 1 {
		t.Fatalf("expected one risk alert summary, got %#v", data.RiskAlerts)
	}
	risk := data.RiskAlerts[0]
	if risk.AlertID != "risk_report_1" || risk.RiskType != "valuation_high" || risk.SOPStatus != "active" || risk.Link != "/risk-alerts/risk_report_1" || len(risk.ProhibitedActions) != 1 || risk.ProhibitedActions[0] != "新增买入" {
		t.Fatalf("unexpected risk summary: %#v", risk)
	}
	if risk.SafetyNote == "" || risk.SuggestedActions[0] != "人工复核分批止盈" {
		t.Fatalf("expected safe manual guidance, got %#v", risk)
	}
}

func TestGetDailyDisciplineReportIncludesP34SourceCoverage(t *testing.T) {
	app, db := testApp(t)
	seedDecisionWithEvidence(t, db, "decision_p34_report")
	_, err := db.Exec(`UPDATE decision_records SET expected_return_scenarios_json=? WHERE decision_id=?`, `{"supporting_data_summary":"P34 可用扩展数据：index_constituents、sentiment_proxy","missing_categories":["index_valuation_files"],"source_health":[{"source_name":"csindex","source_level":"A","source_type":"index_basic","data_category":"index_constituents","freshness":"fresh","data_date":"2026-06-05","last_success_at":"2026-06-06T01:00:00Z","affected_symbols":["000300"]},{"source_name":"csindex","source_level":"A","source_type":"index_basic","data_category":"index_valuation_files","freshness":"parse_error","failure_category":"parse_error","data_date":"2026-06-05","last_failure_at":"2026-06-06T01:00:00Z","affected_symbols":["000300"]}]}`, "decision_p34_report")
	if err != nil {
		t.Fatalf("seed expected return: %v", err)
	}
	seedDailyDisciplineReport(t, db, repository.DailyDisciplineReport{ReportID: "report_p34", LocalDate: "2026-06-08", Scope: "holdings", SymbolSetHash: "hash_p34", SourceType: "auto_run", SourceID: "key_p34", DecisionID: "decision_p34_report", Status: "degraded", Summary: "已生成，部分 P34 数据不足", CreatedAt: "2026-06-08T02:00:00Z", UpdatedAt: "2026-06-08T02:00:00Z"})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/daily-discipline/reports/report_p34", nil)
	req.SetPathValue("report_id", "report_p34")
	w := httptest.NewRecorder()

	app.GetDailyDisciplineReport(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	data := decodeDailyDisciplineReportData(t, w)
	if data.P34SourceCoverage.Summary == "" || !containsStringLocal(data.P34SourceCoverage.MissingCategories, "index_valuation_files") {
		t.Fatalf("expected P34 coverage from expected return context, got %#v", data.P34SourceCoverage)
	}
	if len(data.P34SourceCoverage.SourceHealth) != 2 || data.P34SourceCoverage.SourceHealth[1].Freshness != "parse_error" || data.P34SourceCoverage.SourceHealth[1].DataDate != "2026-06-05" || data.P34SourceCoverage.SourceHealth[1].SourceLevel != "A" {
		t.Fatalf("expected structured P34 source health in report, got %#v", data.P34SourceCoverage.SourceHealth)
	}
}

type dailyDisciplineReportEnvelope struct {
	Data dailyDisciplineReportPayload `json:"data"`
}

type dailyDisciplineReportPayload struct {
	ReportID          string   `json:"report_id"`
	LocalDate         string   `json:"local_date"`
	Scope             string   `json:"scope"`
	Status            string   `json:"status"`
	Summary           string   `json:"summary"`
	SourceType        string   `json:"source_type"`
	SourceID          string   `json:"source_id"`
	DecisionID        string   `json:"decision_id"`
	DecisionLink      string   `json:"decision_link"`
	AutoRunLink       string   `json:"auto_run_link"`
	AuditLink         string   `json:"audit_link"`
	NotificationLink  string   `json:"notification_link"`
	FailureCode       string   `json:"failure_code"`
	FailureReason     string   `json:"failure_reason"`
	MissingAction     string   `json:"missing_action"`
	MissingCategories []string `json:"missing_categories"`
	FinalVerdict      string   `json:"final_verdict"`
	VerdictStatus     string   `json:"verdict_status"`
	Evidence          struct {
		EvidenceCount                   int `json:"evidence_count"`
		IndependentSourceCount          int `json:"independent_source_count"`
		HighGradeIndependentSourceCount int `json:"high_grade_independent_source_count"`
	} `json:"evidence"`
	P34SourceCoverage struct {
		Summary           string   `json:"summary"`
		MissingCategories []string `json:"missing_categories"`
		SourceHealth      []struct {
			SourceName      string   `json:"source_name"`
			SourceLevel     string   `json:"source_level"`
			SourceType      string   `json:"source_type"`
			DataCategory    string   `json:"data_category"`
			Freshness       string   `json:"freshness"`
			DataDate        string   `json:"data_date"`
			AffectedSymbols []string `json:"affected_symbols"`
		} `json:"source_health"`
	} `json:"p34_source_coverage"`
	RiskAlerts []struct {
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
	} `json:"risk_alerts"`
	Trend struct {
		SuccessCount          int `json:"success_count"`
		DegradedCount         int `json:"degraded_count"`
		FailedCount           int `json:"failed_count"`
		InsufficientDataCount int `json:"insufficient_data_count"`
	} `json:"trend"`
	SafetyNote string `json:"safety_note"`
	UpdatedAt  string `json:"updated_at"`
}

type dailyDisciplineReportListEnvelope struct {
	Data struct {
		Reports []dailyDisciplineReportPayload `json:"reports"`
	} `json:"data"`
}

func containsStringLocal(items []string, want string) bool {
	for _, item := range items {
		if item == want {
			return true
		}
	}
	return false
}

func decodeDailyDisciplineReportData(t *testing.T, w *httptest.ResponseRecorder) dailyDisciplineReportPayload {
	t.Helper()
	var envelope dailyDisciplineReportEnvelope
	if err := json.Unmarshal(w.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode response: %v body=%s", err, w.Body.String())
	}
	return envelope.Data
}

func decodeDailyDisciplineReportListData(t *testing.T, w *httptest.ResponseRecorder) struct {
	Reports []dailyDisciplineReportPayload `json:"reports"`
} {
	t.Helper()
	var envelope dailyDisciplineReportListEnvelope
	if err := json.Unmarshal(w.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode response: %v body=%s", err, w.Body.String())
	}
	return envelope.Data
}

func testAppWithDailyAutoRunTimezone(t *testing.T, timezone string) (*App, *sql.DB) {
	t.Helper()
	app, db := testApp(t)
	app.Deps.DailyAutoRunConfig = config.DailyAutoRunConfig{Timezone: timezone}
	app.QuerySvc = service.NewQueryServiceWithDailyAutoRunConfig(repository.Repositories{
		DecisionRepo:              app.Deps.DecisionRepo,
		AuditRepo:                 app.Deps.AuditRepo,
		RuleRepo:                  app.Deps.RuleRepo,
		MarketRepo:                app.Deps.MarketRepo,
		SettingsRepo:              app.Deps.SettingsRepo,
		IntelligenceRepo:          app.Deps.IntelligenceRepo,
		NotificationRepo:          app.Deps.NotificationRepo,
		DailyAutoRunRepo:          app.Deps.DailyAutoRunRepo,
		DailyDisciplineReportRepo: app.Deps.DailyDisciplineReportRepo,
		RiskAlertRepo:             app.Deps.RiskAlertRepo,
		PortfolioRepo:             app.Deps.PortfolioRepo,
	}, app.Deps.DailyAutoRunConfig)
	return app, db
}

func testContext(t *testing.T) context.Context {
	t.Helper()
	return context.Background()
}

func seedDailyDisciplineReport(t *testing.T, db *sql.DB, report repository.DailyDisciplineReport) {
	t.Helper()
	_, err := db.Exec(`INSERT INTO daily_discipline_reports (report_id,local_date,scope,symbol_set_hash,source_type,source_id,decision_id,status,summary,failure_code,failure_reason,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`, report.ReportID, report.LocalDate, report.Scope, report.SymbolSetHash, report.SourceType, nullStringLocal(report.SourceID), nullStringLocal(report.DecisionID), report.Status, nullStringLocal(report.Summary), nullStringLocal(report.FailureCode), nullStringLocal(report.FailureReason), report.CreatedAt, report.UpdatedAt)
	if err != nil {
		t.Fatalf("seed report: %v", err)
	}
}

func seedDecisionWithEvidence(t *testing.T, db *sql.DB, decisionID string) {
	t.Helper()
	_, err := db.Exec(`INSERT INTO decision_records (decision_id,request_id,workflow_type,symbol,workflow_status,record_type,dashboard_state,final_verdict_status,final_verdict_text,confirmation_status,rule_version,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`, decisionID, "req_report", "consultation", "510300", "completed", "formal_trade_advice", "normal", "hold", "持有", "pending", "v3.0", "2026-06-08T01:00:00Z")
	if err != nil {
		t.Fatalf("seed decision: %v", err)
	}
	for _, item := range []struct {
		id        string
		ind       int
		highGrade int
	}{
		{id: "eref_1", ind: 1, highGrade: 1},
		{id: "eref_2", ind: 3, highGrade: 2},
	} {
		_, err = db.Exec(`INSERT INTO evidence_refs (evidence_ref_id,evidence_id,decision_id,summary_id,source_name,source_level,evidence_role,summary,time_weight,relevance_score,independent_source_count,high_grade_independent_source_count,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`, item.id, "evidence_"+item.id, decisionID, "summary_"+item.id, "source", "A", "formal", "摘要", 1.0, 1.0, item.ind, item.highGrade, "2026-06-08T01:00:00Z")
		if err != nil {
			t.Fatalf("seed evidence ref: %v", err)
		}
	}
}
