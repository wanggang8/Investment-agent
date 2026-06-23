package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	appsqlite "investment-agent/internal/infrastructure/persistence/sqlite"
	"investment-agent/internal/pkg/apperr"
)

func TestDecisionLoopServiceBuildsCompleteSafeLoop(t *testing.T) {
	ctx := context.Background()
	repos, db := decisionLoopRepos(t)
	seedDecisionLoopDecision(t, repos.DecisionRepo, repository.DecisionRecord{
		DecisionID:               "decision_loop_full",
		RequestID:                "req_loop_full",
		WorkflowType:             "consultation",
		Symbol:                   "510300",
		WorkflowStatus:           "completed",
		RecordType:               "formal_trade_advice",
		DashboardState:           "normal",
		SourceVerificationStatus: "satisfied",
		FinalVerdictStatus:       "hold",
		FinalVerdictText:         "继续持有，等待人工复核",
		ConfirmationStatus:       "executed_manually",
		RuleVersion:              "v3.0",
		CreatedAt:                "2026-06-16T09:00:00Z",
	})
	if err := repos.DecisionRepo.SaveOperationConfirmation(ctx, repository.OperationConfirmation{
		ConfirmationID:   "conf_loop_full",
		DecisionID:       "decision_loop_full",
		ConfirmationType: "executed_manually",
		OperationType:    "buy",
		Symbol:           "510300",
		Quantity:         10,
		Price:            2.5,
		Fees:             1,
		ExecutedAt:       "2026-06-16T09:30:00Z",
		PayloadJSON:      `{"raw":"SELECT * FROM secret"}`,
		Note:             "人工记录包含 prompt: 和 Prompt:、完整 prompt、/Users/private/local.txt、sk-123456789012、sk-proj-abc_def-123456、select    *    from secret、raw HTTP、GET /secret HTTP/1.1、HTTP/1.1 200 OK、-----BEGIN RSA PRIVATE KEY-----abc-----END RSA PRIVATE KEY-----",
		CreatedAt:        "2026-06-16T09:20:00Z",
	}); err != nil {
		t.Fatal(err)
	}
	if err := repos.DecisionRepo.SavePositionTransaction(ctx, repository.PositionTransaction{
		TransactionID:  "tx_loop_full",
		ConfirmationID: "conf_loop_full",
		Symbol:         "510300",
		OperationType:  "buy",
		Quantity:       10,
		Price:          2.5,
		Fees:           1,
		OccurredAt:     "2026-06-16T09:30:00Z",
		CreatedAt:      "2026-06-16T09:31:00Z",
	}); err != nil {
		t.Fatal(err)
	}
	if err := repos.DecisionRepo.SaveErrorCase(ctx, repository.ErrorCase{
		ErrorCaseID:    "err_loop_full",
		DecisionID:     "decision_loop_full",
		ConfirmationID: "conf_loop_full",
		ActualOutcome:  "missed",
		RootCauseTag:   "analyst_error",
		LessonLearned:  "后续需要复核样本。",
		CreatedAt:      "2026-06-16T10:00:00Z",
	}); err != nil {
		t.Fatal(err)
	}
	if err := repos.RiskAlertRepo.UpsertRiskAlert(ctx, repository.RiskAlert{
		AlertID:           "risk_loop_full",
		RiskType:          model.RiskTypeValuationHigh,
		Severity:          model.RiskSeverityWarning,
		SOPStatus:         model.RiskSOPActive,
		Symbol:            "510300",
		TriggerSummary:    "估值偏高，进入人工复核。",
		RelatedDecisionID: "decision_loop_full",
		CreatedAt:         "2026-06-16T09:40:00Z",
		UpdatedAt:         "2026-06-16T09:40:00Z",
	}); err != nil {
		t.Fatal(err)
	}
	if err := repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{
		AuditEventID:   "audit_loop_full",
		DecisionID:     "decision_loop_full",
		ConfirmationID: "conf_loop_full",
		ErrorCaseID:    "err_loop_full",
		Actor:          "user",
		Action:         "confirm_operation",
		Status:         "success",
		CreatedAt:      "2026-06-16T09:32:00Z",
	}); err != nil {
		t.Fatal(err)
	}
	beforeCounts := decisionLoopTableCounts(t, db)

	got, err := NewDecisionLoopService(repos).ListDecisionLoops(ctx, DecisionLoopListFilter{Symbol: "510300", Limit: 5})
	if err != nil {
		t.Fatal(err)
	}
	if got.Total != 1 || len(got.Items) != 1 {
		t.Fatalf("expected one loop, got %#v", got)
	}
	item := got.Items[0]
	if item.DecisionID != "decision_loop_full" || item.Symbol != "510300" || item.LoopStatus != "reviewed" {
		t.Fatalf("unexpected loop identity/status: %#v", item)
	}
	assertDecisionLoopStage(t, item.Stages, "recommendation", "complete")
	assertDecisionLoopStage(t, item.Stages, "confirmation", "complete")
	assertDecisionLoopStage(t, item.Stages, "manual_record", "complete")
	assertDecisionLoopStage(t, item.Stages, "risk_review", "complete")
	assertDecisionLoopStage(t, item.Stages, "review", "complete")
	if len(item.ManualActions) != 1 || len(item.ManualActions[0].TransactionIDs) != 1 || item.ManualActions[0].TransactionIDs[0] != "tx_loop_full" {
		t.Fatalf("manual action missing transaction ids: %#v", item.ManualActions)
	}
	if len(item.RiskLinks) != 1 || item.RiskLinks[0].ID != "risk_loop_full" {
		t.Fatalf("risk links not built: %#v", item.RiskLinks)
	}
	if len(item.ReviewLinks) == 0 || item.ReviewLinks[0].ID != "err_loop_full" {
		t.Fatalf("review links not built: %#v", item.ReviewLinks)
	}
	if len(item.AuditLinks) != 1 || item.AuditLinks[0].ID != "audit_loop_full" {
		t.Fatalf("audit links not built: %#v", item.AuditLinks)
	}
	if len(item.MissingLinks) != 0 {
		t.Fatalf("expected no missing links, got %#v", item.MissingLinks)
	}
	body, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"SELECT * FROM", "select    *    from", "/Users/private", "prompt:", "Prompt:", "完整 prompt", "raw HTTP", "GET /secret HTTP/1.1", "HTTP/1.1 200 OK", "BEGIN RSA PRIVATE KEY"} {
		if strings.Contains(string(body), forbidden) {
			t.Fatalf("response leaked forbidden fragment %q: %s", forbidden, string(body))
		}
	}
	for _, nullArray := range []string{`"manual_actions":null`, `"risk_links":null`, `"review_links":null`, `"audit_links":null`, `"missing_links":null`} {
		if strings.Contains(string(body), nullArray) {
			t.Fatalf("response returned null array %q: %s", nullArray, string(body))
		}
	}
	if regexp.MustCompile(`sk-[A-Za-z0-9]{12,}`).Match(body) {
		t.Fatalf("response leaked complete key-like fragment: %s", string(body))
	}
	if regexp.MustCompile(`sk-[A-Za-z0-9][A-Za-z0-9_-]{8,}`).Match(body) {
		t.Fatalf("response leaked hyphenated key-like fragment: %s", string(body))
	}
	if after := decisionLoopTableCounts(t, db); !sameDecisionLoopCounts(beforeCounts, after) {
		t.Fatalf("decision loop service wrote tables: before=%v after=%v", beforeCounts, after)
	}
}

func TestDecisionLoopServiceDetectsMissingManualRecord(t *testing.T) {
	ctx := context.Background()
	repos, _ := decisionLoopRepos(t)
	seedDecisionLoopDecision(t, repos.DecisionRepo, repository.DecisionRecord{
		DecisionID:         "decision_loop_missing",
		RequestID:          "req_loop_missing",
		WorkflowType:       "consultation",
		Symbol:             "510300",
		WorkflowStatus:     "completed",
		RecordType:         "formal_trade_advice",
		DashboardState:     "normal",
		FinalVerdictStatus: "hold",
		FinalVerdictText:   "持有但缺少线下记录",
		ConfirmationStatus: "executed_manually",
		RuleVersion:        "v3.0",
		CreatedAt:          "2026-06-16T11:00:00Z",
	})

	item, err := NewDecisionLoopService(repos).GetDecisionLoop(ctx, "decision_loop_missing")
	if err != nil {
		t.Fatal(err)
	}
	assertDecisionLoopStage(t, item.Stages, "confirmation", "missing")
	assertDecisionLoopStage(t, item.Stages, "manual_record", "missing")
	if len(item.MissingLinks) < 2 {
		t.Fatalf("expected confirmation and manual record gaps, got %#v", item.MissingLinks)
	}
	if item.LoopStatus != "incomplete" {
		t.Fatalf("expected incomplete loop, got %s", item.LoopStatus)
	}
}

func TestDecisionLoopServiceHandlesPlannedWatchAndNotRequired(t *testing.T) {
	ctx := context.Background()
	repos, _ := decisionLoopRepos(t)
	cases := []struct {
		id                 string
		confirmationStatus string
		confirmationType   string
		wantLoopStatus     string
		wantManualStatus   string
	}{
		{id: "decision_loop_planned", confirmationStatus: "planned", confirmationType: "planned", wantLoopStatus: "planned", wantManualStatus: "not_required"},
		{id: "decision_loop_watch", confirmationStatus: "watch", confirmationType: "watch", wantLoopStatus: "open", wantManualStatus: "not_required"},
		{id: "decision_loop_not_required", confirmationStatus: "not_required", wantLoopStatus: "reviewed", wantManualStatus: "not_required"},
	}
	for idx, tc := range cases {
		seedDecisionLoopDecision(t, repos.DecisionRepo, repository.DecisionRecord{
			DecisionID:         tc.id,
			RequestID:          "req_" + tc.id,
			WorkflowType:       "consultation",
			Symbol:             "510300",
			WorkflowStatus:     "completed",
			RecordType:         "formal_trade_advice",
			DashboardState:     "normal",
			FinalVerdictStatus: "hold",
			FinalVerdictText:   "持有",
			ConfirmationStatus: tc.confirmationStatus,
			RuleVersion:        "v3.0",
			CreatedAt:          "2026-06-16T12:0" + string(rune('0'+idx)) + ":00Z",
		})
		if tc.confirmationType != "" {
			if err := repos.DecisionRepo.SaveOperationConfirmation(ctx, repository.OperationConfirmation{
				ConfirmationID:   "conf_" + tc.id,
				DecisionID:       tc.id,
				ConfirmationType: tc.confirmationType,
				CreatedAt:        "2026-06-16T12:30:00Z",
			}); err != nil {
				t.Fatal(err)
			}
		}
		item, err := NewDecisionLoopService(repos).GetDecisionLoop(ctx, tc.id)
		if err != nil {
			t.Fatal(err)
		}
		if item.LoopStatus != tc.wantLoopStatus {
			t.Fatalf("%s loop status: want %s got %s", tc.id, tc.wantLoopStatus, item.LoopStatus)
		}
		assertDecisionLoopStage(t, item.Stages, "manual_record", tc.wantManualStatus)
	}
}

func TestDecisionLoopServiceLimitSymbolAndNotFound(t *testing.T) {
	ctx := context.Background()
	repos, _ := decisionLoopRepos(t)
	for _, item := range []struct {
		id      string
		symbol  string
		created string
	}{
		{"decision_loop_a", "510300", "2026-06-16T10:00:00Z"},
		{"decision_loop_b", "510300", "2026-06-16T11:00:00Z"},
		{"decision_loop_c", "510500", "2026-06-16T12:00:00Z"},
	} {
		seedDecisionLoopDecision(t, repos.DecisionRepo, repository.DecisionRecord{
			DecisionID:         item.id,
			RequestID:          "req_" + item.id,
			WorkflowType:       "consultation",
			Symbol:             item.symbol,
			WorkflowStatus:     "completed",
			RecordType:         "formal_trade_advice",
			DashboardState:     "normal",
			FinalVerdictStatus: "hold",
			FinalVerdictText:   "持有",
			ConfirmationStatus: "pending",
			RuleVersion:        "v3.0",
			CreatedAt:          item.created,
		})
	}

	got, err := NewDecisionLoopService(repos).ListDecisionLoops(ctx, DecisionLoopListFilter{Symbol: "510300", Limit: 1})
	if err != nil {
		t.Fatal(err)
	}
	if got.Total != 1 || got.Items[0].DecisionID != "decision_loop_b" {
		t.Fatalf("expected latest filtered item, got %#v", got)
	}
	if _, err := NewDecisionLoopService(repos).GetDecisionLoop(ctx, "missing_decision"); !apperr.IsCode(err, apperr.CodeNotFound) {
		t.Fatalf("expected NOT_FOUND, got %v", err)
	}
}

func decisionLoopRepos(t *testing.T) (repository.Repositories, *sql.DB) {
	t.Helper()
	store, err := appsqlite.Open(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	if err := appsqlite.Migrate(context.Background(), store.DB); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return repository.Repositories{
		DecisionRepo:  appsqlite.NewDecisionRepository(store.DB),
		AuditRepo:     appsqlite.NewAuditRepository(store.DB),
		RiskAlertRepo: appsqlite.NewRiskAlertRepository(store.DB),
	}, store.DB
}

func seedDecisionLoopDecision(t *testing.T, repo repository.DecisionRepository, d repository.DecisionRecord) {
	t.Helper()
	if d.CreatedAt == "" {
		d.CreatedAt = "2026-06-16T09:00:00Z"
	}
	if d.SourceVerificationStatus == "" {
		d.SourceVerificationStatus = "satisfied"
	}
	if err := repo.SaveDecisionRecord(context.Background(), d, nil); err != nil {
		t.Fatalf("seed decision %s: %v", d.DecisionID, err)
	}
}

func assertDecisionLoopStage(t *testing.T, stages []dto.DecisionLoopStage, stage string, wantStatus string) {
	t.Helper()
	for _, item := range stages {
		if item.Stage == stage {
			if item.Status != wantStatus {
				t.Fatalf("%s status: want %s got %s in %#v", stage, wantStatus, item.Status, stages)
			}
			return
		}
	}
	t.Fatalf("missing stage %s in %#v", stage, stages)
}

func decisionLoopTableCounts(t *testing.T, db *sql.DB) map[string]int {
	t.Helper()
	out := map[string]int{}
	for _, table := range []string{"decision_records", "operation_confirmations", "position_transactions", "error_cases", "risk_alerts", "audit_events"} {
		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
			t.Fatalf("count %s: %v", table, err)
		}
		out[table] = count
	}
	return out
}

func sameDecisionLoopCounts(a, b map[string]int) bool {
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
