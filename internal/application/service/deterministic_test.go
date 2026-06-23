package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/application/workflow"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/infrastructure/config"
	"investment-agent/internal/pkg/apperr"
	"investment-agent/internal/pkg/clock"
	"investment-agent/internal/pkg/idgen"
)

type transactorStub struct{ repos repository.Repositories }

func (t transactorStub) WithinTx(ctx context.Context, fn func(context.Context, repository.Repositories) error) error {
	return fn(ctx, t.repos)
}

type portfolioRepoStub struct {
	snapshot          repository.PortfolioSnapshot
	position          repository.Position
	positions         []repository.Position
	positionSnapshots []repository.PositionSnapshot
	importBatch       repository.LocalAccountImportBatch
	importBatches     map[string]repository.LocalAccountImportBatch
	correction        repository.LocalAccountCorrection
	latestErr         error
}

func (r *portfolioRepoStub) SavePortfolioSnapshot(_ context.Context, s repository.PortfolioSnapshot, ps []repository.PositionSnapshot) error {
	r.snapshot = s
	r.positionSnapshots = ps
	return nil
}
func (r *portfolioRepoStub) GetPortfolioSnapshot(context.Context, string) (repository.PortfolioSnapshot, []repository.PositionSnapshot, error) {
	return repository.PortfolioSnapshot{}, nil, nil
}
func (r *portfolioRepoStub) GetLatestPortfolioSnapshot(context.Context) (repository.PortfolioSnapshot, error) {
	if r.latestErr != nil {
		return repository.PortfolioSnapshot{}, r.latestErr
	}
	if r.snapshot.SnapshotID == "" && r.snapshot.Cash == 0 && r.snapshot.TotalAssets == 0 {
		return repository.PortfolioSnapshot{SnapshotID: "snap_default", Cash: 100000, TotalAssets: 100000, CashRatio: 1}, nil
	}
	return r.snapshot, nil
}
func (r *portfolioRepoStub) SavePosition(_ context.Context, p repository.Position) error {
	r.position = p
	updated := false
	for i, item := range r.positions {
		if item.Symbol == p.Symbol {
			r.positions[i] = p
			updated = true
		}
	}
	if !updated {
		r.positions = append(r.positions, p)
	}
	return nil
}
func (r *portfolioRepoStub) ReplacePositions(_ context.Context, positions []repository.Position) error {
	r.positions = append([]repository.Position(nil), positions...)
	if len(positions) > 0 {
		r.position = positions[len(positions)-1]
	}
	return nil
}
func (r *portfolioRepoStub) DeletePosition(_ context.Context, positionID string) error {
	for i, item := range r.positions {
		if item.PositionID == positionID {
			r.positions = append(r.positions[:i], r.positions[i+1:]...)
			return nil
		}
	}
	return nil
}
func (r *portfolioRepoStub) GetPosition(context.Context, string) (repository.Position, error) {
	return repository.Position{}, nil
}
func (r *portfolioRepoStub) ListPositions(context.Context) ([]repository.Position, error) {
	return r.positions, nil
}
func (r *portfolioRepoStub) SaveLocalAccountImportBatch(_ context.Context, b repository.LocalAccountImportBatch) error {
	r.importBatch = b
	if r.importBatches == nil {
		r.importBatches = map[string]repository.LocalAccountImportBatch{}
	}
	r.importBatches[b.ImportBatchID] = b
	return nil
}
func (r *portfolioRepoStub) GetLocalAccountImportBatch(_ context.Context, importBatchID string) (repository.LocalAccountImportBatch, error) {
	if r.importBatches != nil {
		if batch, ok := r.importBatches[importBatchID]; ok {
			return batch, nil
		}
	}
	if r.importBatch.ImportBatchID == importBatchID {
		return r.importBatch, nil
	}
	return repository.LocalAccountImportBatch{}, apperr.New(apperr.CodeNotFound, apperr.CategoryNotFound, "import batch not found")
}
func (r *portfolioRepoStub) SaveLocalAccountCorrection(_ context.Context, c repository.LocalAccountCorrection) error {
	r.correction = c
	return nil
}

type auditRepoStub struct{ event repository.AuditEvent }

func (r *auditRepoStub) AppendAuditEvent(_ context.Context, e repository.AuditEvent) error {
	r.event = e
	return nil
}
func (r *auditRepoStub) GetAuditEvent(context.Context, string) (repository.AuditEvent, error) {
	return repository.AuditEvent{}, nil
}
func (r *auditRepoStub) ListAuditEvents(context.Context) ([]repository.AuditEvent, error) {
	return nil, nil
}

type decisionRepoStub struct {
	confirmation   repository.OperationConfirmation
	transaction    repository.PositionTransaction
	errorCase      repository.ErrorCase
	status         string
	staleUpdate    bool
	getRecordCalls int
}

func (r *decisionRepoStub) SaveDecisionRecord(context.Context, repository.DecisionRecord, []repository.EvidenceRef) error {
	return nil
}
func (r *decisionRepoStub) GetDecisionRecord(context.Context, string) (repository.DecisionRecord, []repository.EvidenceRef, error) {
	r.getRecordCalls++
	return repository.DecisionRecord{}, nil, nil
}
func (r *decisionRepoStub) ListDecisionRecords(context.Context) ([]repository.DecisionRecord, error) {
	return nil, nil
}
func (r *decisionRepoStub) ListErrorCases(context.Context) ([]repository.ErrorCase, error) {
	return nil, nil
}
func (r *decisionRepoStub) CountErrorCases(context.Context) (int, error) {
	return 0, nil
}
func (r *decisionRepoStub) GetDecisionConfirmationState(context.Context, string) (string, string, error) {
	return "formal_trade_advice", string(model.ConfirmationPending), nil
}
func (r *decisionRepoStub) ListOperationConfirmations(context.Context, string) ([]repository.OperationConfirmation, error) {
	return nil, nil
}
func (r *decisionRepoStub) SaveOperationConfirmation(_ context.Context, c repository.OperationConfirmation) error {
	r.confirmation = c
	return nil
}
func (r *decisionRepoStub) UpdateDecisionConfirmationStatus(_ context.Context, _ string, status string) error {
	r.status = status
	return nil
}
func (r *decisionRepoStub) UpdateDecisionConfirmationStatusIfCurrent(_ context.Context, _ string, expectedStatus, nextStatus string) (bool, error) {
	if r.staleUpdate {
		return false, nil
	}
	current := r.status
	if current == "" {
		current = string(model.ConfirmationPending)
	}
	if current != expectedStatus {
		return false, nil
	}
	r.status = nextStatus
	return true, nil
}
func (r *decisionRepoStub) ListPositionTransactionsByConfirmation(context.Context, string) ([]repository.PositionTransaction, error) {
	return nil, nil
}
func (r *decisionRepoStub) SavePositionTransaction(_ context.Context, tx repository.PositionTransaction) error {
	r.transaction = tx
	return nil
}
func (r *decisionRepoStub) SaveErrorCase(_ context.Context, e repository.ErrorCase) error {
	r.errorCase = e
	return nil
}
func (r *decisionRepoStub) GetOperationConfirmation(context.Context, string) (repository.OperationConfirmation, error) {
	return repository.OperationConfirmation{}, nil
}

type ruleRepoStub struct {
	proposalStatus string
	gatekeeper     repository.GatekeeperAudit
	ruleVersion    repository.RuleVersion
	appliedAt      string
	appliedVersion string
	archived       bool
}

func (r *ruleRepoStub) SaveRuleVersion(_ context.Context, v repository.RuleVersion) error {
	r.ruleVersion = v
	return nil
}
func (r *ruleRepoStub) GetRuleVersion(context.Context, string) (repository.RuleVersion, error) {
	return repository.RuleVersion{}, nil
}
func (r *ruleRepoStub) GetActiveRuleVersion(context.Context) (repository.RuleVersion, error) {
	return repository.RuleVersion{RuleVersion: "v_test", Status: "active"}, nil
}
func (r *ruleRepoStub) SaveRuleProposal(context.Context, repository.RuleProposal) error { return nil }
func (r *ruleRepoStub) GetRuleProposal(_ context.Context, proposalID string) (repository.RuleProposal, error) {
	return repository.RuleProposal{ProposalID: proposalID, Status: string(model.ProposalUnderGatekeeperAudit), SampleCount: 3, AfterRuleJSON: "{}"}, nil
}
func (r *ruleRepoStub) ListRuleProposals(context.Context) ([]repository.RuleProposalWithAudit, error) {
	return nil, nil
}
func (r *ruleRepoStub) UpdateRuleProposalStatus(_ context.Context, _ string, status string) error {
	r.proposalStatus = status
	return nil
}
func (r *ruleRepoStub) ApplyRuleProposal(_ context.Context, _ string, status, finalConfirmedAt, finalConfirmedNote, appliedRuleVersion string) error {
	r.proposalStatus = status
	r.appliedAt = finalConfirmedAt
	r.appliedVersion = appliedRuleVersion
	return nil
}
func (r *ruleRepoStub) ArchiveActiveRuleVersions(context.Context) error {
	r.archived = true
	return nil
}
func (r *ruleRepoStub) SaveGatekeeperAudit(_ context.Context, a repository.GatekeeperAudit) error {
	r.gatekeeper = a
	return nil
}
func (r *ruleRepoStub) GetGatekeeperAudit(context.Context, string) (repository.GatekeeperAudit, error) {
	return r.gatekeeper, nil
}
func (r *ruleRepoStub) GetLatestGatekeeperAuditByProposal(context.Context, string) (repository.GatekeeperAudit, error) {
	return r.gatekeeper, nil
}

func TestTodayDailyDisciplineReportUsesConfiguredDailyAutoRunTimezone(t *testing.T) {
	ctx := context.Background()
	autoRunRepo := &dailyAutoRunRepoStub{state: repository.DailyAutoRunState{RunID: "run_shanghai", IdempotencyKey: "key_shanghai", LocalDate: "2026-06-08", Scope: "holdings", SymbolSetHash: "hash_shanghai", Status: "failed", FailureCode: "missing_prerequisites", FailureReason: "缺少本地账户或持仓", UpdatedAt: "2026-06-08T00:30:00+08:00"}}
	reportRepo := &dailyDisciplineReportRepoStub{}
	svc := NewQueryServiceWithDailyAutoRunConfig(repository.Repositories{DailyAutoRunRepo: autoRunRepo, DailyDisciplineReportRepo: reportRepo}, config.DailyAutoRunConfig{Timezone: "Asia/Shanghai"})

	out, err := svc.TodayDailyDisciplineReport(ctx, time.Date(2026, 6, 7, 16, 30, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("TodayDailyDisciplineReport: %v", err)
	}
	if out.LocalDate != "2026-06-08" || out.Status != "insufficient_data" {
		t.Fatalf("expected Shanghai local date insufficient data, got local_date=%q status=%q", out.LocalDate, out.Status)
	}
	if out.FailureReason != "缺少本地账户或持仓" {
		t.Fatalf("expected auto-run failure reason, got %q", out.FailureReason)
	}
}

func TestTodayDailyDisciplineReportDoesNotFallbackToOldSameDayHashWhenLatestStateReportMissing(t *testing.T) {
	ctx := context.Background()
	autoRunRepo := &dailyAutoRunRepoStub{state: repository.DailyAutoRunState{RunID: "run_b", IdempotencyKey: "key_b", LocalDate: "2026-06-08", Scope: "holdings", SymbolSetHash: "hash_b", Status: "failed", FailureCode: "missing_prerequisites", FailureReason: "latest state missing report", UpdatedAt: "2026-06-08T02:00:00Z"}}
	reportRepo := &dailyDisciplineReportRepoStub{reports: []repository.DailyDisciplineReport{{ReportID: "report_a", LocalDate: "2026-06-08", Scope: "holdings", SymbolSetHash: "hash_a", SourceType: "auto_run", SourceID: "key_a", Status: "success", Summary: "old same-day report", UpdatedAt: "2026-06-08T01:00:00Z"}}}
	svc := NewQueryServiceWithDailyAutoRunConfig(repository.Repositories{DailyAutoRunRepo: autoRunRepo, DailyDisciplineReportRepo: reportRepo}, config.DailyAutoRunConfig{Timezone: "UTC"})

	out, err := svc.TodayDailyDisciplineReport(ctx, time.Date(2026, 6, 8, 3, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("TodayDailyDisciplineReport: %v", err)
	}
	if out.SourceID != "key_b" || out.Status != "insufficient_data" || out.Summary == "old same-day report" {
		t.Fatalf("expected synthesized latest state instead of old same-day report, got %+v", out)
	}
}

func TestTodayDailyDisciplineReportReportsMissingCategoriesForEmptyPrerequisites(t *testing.T) {
	ctx := context.Background()
	svc := NewQueryServiceWithDailyAutoRunConfig(repository.Repositories{DailyAutoRunRepo: &dailyAutoRunRepoStub{}, DailyDisciplineReportRepo: &dailyDisciplineReportRepoStub{}}, config.DailyAutoRunConfig{Timezone: "UTC"})

	out, err := svc.TodayDailyDisciplineReport(ctx, time.Date(2026, 6, 8, 3, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("TodayDailyDisciplineReport: %v", err)
	}
	if out.Status != "not_started" || !containsString(out.MissingCategories, "account") || !containsString(out.MissingCategories, "holdings") {
		t.Fatalf("expected empty prerequisites categories for not_started, got %+v", out)
	}
}

func TestListDailyDisciplineReportsCalculatesTrendOnce(t *testing.T) {
	ctx := context.Background()
	reportRepo := &dailyDisciplineReportRepoStub{reports: []repository.DailyDisciplineReport{
		{ReportID: "report_1", LocalDate: "2026-06-08", Scope: "holdings", Status: "success"},
		{ReportID: "report_2", LocalDate: "2026-06-07", Scope: "holdings", Status: "failed"},
		{ReportID: "report_3", LocalDate: "2026-06-06", Scope: "holdings", Status: "insufficient_data"},
	}}
	svc := NewQueryService(repository.Repositories{DailyDisciplineReportRepo: reportRepo})

	out, err := svc.ListDailyDisciplineReports(ctx, "", 30)
	if err != nil {
		t.Fatalf("ListDailyDisciplineReports: %v", err)
	}
	if reportRepo.listCalls != 2 {
		t.Fatalf("expected one list call for reports and one for shared trend, got %d", reportRepo.listCalls)
	}
	if len(out.Reports) != 3 {
		t.Fatalf("expected 3 reports, got %d", len(out.Reports))
	}
	for _, report := range out.Reports {
		if report.Trend.SuccessCount != 1 || report.Trend.FailedCount != 1 || report.Trend.InsufficientDataCount != 1 {
			t.Fatalf("expected shared trend on report %s, got %#v", report.ReportID, report.Trend)
		}
	}
}

func TestListDailyDisciplineReportsDoesNotLoadDecisionRecords(t *testing.T) {
	ctx := context.Background()
	reportRepo := &dailyDisciplineReportRepoStub{reports: []repository.DailyDisciplineReport{
		{ReportID: "report_1", LocalDate: "2026-06-08", Scope: "holdings", SourceType: "auto_run", SourceID: "auto_key_1", DecisionID: "decision_1", Status: "success"},
		{ReportID: "report_2", LocalDate: "2026-06-07", Scope: "holdings", SourceType: "auto_run", SourceID: "auto_key_2", DecisionID: "decision_2", Status: "degraded"},
	}}
	decisionRepo := &decisionRepoStub{}
	svc := NewQueryService(repository.Repositories{DailyDisciplineReportRepo: reportRepo, DecisionRepo: decisionRepo})

	out, err := svc.ListDailyDisciplineReports(ctx, "", 30)
	if err != nil {
		t.Fatalf("ListDailyDisciplineReports: %v", err)
	}
	if decisionRepo.getRecordCalls != 0 {
		t.Fatalf("list must not call DecisionRepo.GetDecisionRecord, got %d calls", decisionRepo.getRecordCalls)
	}
	if len(out.Reports) != 2 || out.Reports[0].DecisionID != "decision_1" || out.Reports[0].FinalVerdict != "" || out.Reports[0].Evidence.EvidenceCount != 0 {
		t.Fatalf("expected lightweight list items without verdict/evidence, got %+v", out.Reports)
	}
}

func TestDailyDisciplineReportLinksEscapeSourceID(t *testing.T) {
	ctx := context.Background()
	reportRepo := &dailyDisciplineReportRepoStub{reports: []repository.DailyDisciplineReport{
		{ReportID: "report_escape", LocalDate: "2026-06-08", Scope: "holdings", SourceType: "auto_run", SourceID: "manual&status=failed", Status: "failed", Summary: "failed"},
	}}
	svc := NewQueryService(repository.Repositories{DailyDisciplineReportRepo: reportRepo})

	out, err := svc.ListDailyDisciplineReports(ctx, "", 30)
	if err != nil {
		t.Fatalf("ListDailyDisciplineReports: %v", err)
	}
	if len(out.Reports) != 1 {
		t.Fatalf("expected one report, got %d", len(out.Reports))
	}
	if out.Reports[0].AuditLink != "/audit?input_ref=manual%26status%3Dfailed" || out.Reports[0].NotificationLink != "/notifications?source_id=manual%26status%3Dfailed" {
		t.Fatalf("expected escaped source links, got audit=%q notification=%q", out.Reports[0].AuditLink, out.Reports[0].NotificationLink)
	}
}

func containsString(items []string, want string) bool {
	for _, item := range items {
		if item == want {
			return true
		}
	}
	return false
}

type dailyAutoRunRepoStub struct {
	state repository.DailyAutoRunState
}

func (r *dailyAutoRunRepoStub) UpsertDailyAutoRunState(context.Context, repository.DailyAutoRunState) error {
	return nil
}
func (r *dailyAutoRunRepoStub) GetDailyAutoRunState(context.Context, string) (repository.DailyAutoRunState, error) {
	return repository.DailyAutoRunState{}, apperr.New(apperr.CodeNotFound, apperr.CategoryNotFound, "not found")
}
func (r *dailyAutoRunRepoStub) GetLatestDailyAutoRunState(context.Context) (repository.DailyAutoRunState, error) {
	if r.state.LocalDate == "" {
		return repository.DailyAutoRunState{}, apperr.New(apperr.CodeNotFound, apperr.CategoryNotFound, "not found")
	}
	return r.state, nil
}

type dailyDisciplineReportRepoStub struct {
	reports   []repository.DailyDisciplineReport
	listCalls int
}

func (r *dailyDisciplineReportRepoStub) UpsertDailyDisciplineReport(context.Context, repository.DailyDisciplineReport) error {
	return nil
}
func (r *dailyDisciplineReportRepoStub) GetDailyDisciplineReport(context.Context, string) (repository.DailyDisciplineReport, error) {
	return repository.DailyDisciplineReport{}, apperr.New(apperr.CodeNotFound, apperr.CategoryNotFound, "not found")
}
func (r *dailyDisciplineReportRepoStub) GetDailyDisciplineReportByKey(context.Context, string, string, string) (repository.DailyDisciplineReport, error) {
	return repository.DailyDisciplineReport{}, apperr.New(apperr.CodeNotFound, apperr.CategoryNotFound, "not found")
}
func (r *dailyDisciplineReportRepoStub) ListDailyDisciplineReports(context.Context, repository.DailyDisciplineReportListFilter) ([]repository.DailyDisciplineReport, error) {
	r.listCalls++
	return append([]repository.DailyDisciplineReport(nil), r.reports...), nil
}

func TestPortfolioServiceUsesInjectedClockAndIDs(t *testing.T) {
	portfolioRepo := &portfolioRepoStub{}
	auditRepo := &auditRepoStub{}
	svc := &PortfolioService{
		tx:  transactorStub{repos: repository.Repositories{PortfolioRepo: portfolioRepo, AuditRepo: auditRepo}},
		clk: clock.FixedClock{Time: time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)},
		ids: idgen.NewFixedGenerator(map[string][]string{"snap": {"snap_fixed"}, "audit": {"audit_fixed"}, "pos": {"pos_fixed"}, "ps": {"ps_fixed"}}),
	}
	out, err := svc.WriteSnapshot(context.Background(), "req_port", dto.PortfolioInitRequest{Cash: 98, TotalAssets: 100, Positions: []dto.PositionInput{{Symbol: "510300", Name: "ETF", Quantity: 1, CostPrice: 1, CurrentPrice: 2, BuyReason: "低估配置"}}}, "manual")
	if err != nil {
		t.Fatalf("WriteSnapshot: %v", err)
	}
	if out.SnapshotID != "snap_fixed" || out.AuditEventIDs[0] != "audit_fixed" {
		t.Fatalf("unexpected ids: %+v", out)
	}
	if portfolioRepo.snapshot.CreatedAt != "2026-05-29T04:00:00Z" || auditRepo.event.CreatedAt != "2026-05-29T04:00:00Z" {
		t.Fatalf("unexpected timestamps snapshot=%s audit=%s", portfolioRepo.snapshot.CreatedAt, auditRepo.event.CreatedAt)
	}
}

func TestPortfolioServiceRejectsInvalidInitializationWithoutWritingFacts(t *testing.T) {
	cases := []struct {
		name string
		req  dto.PortfolioInitRequest
	}{
		{name: "negative cash", req: dto.PortfolioInitRequest{Cash: -1, TotalAssets: 100}},
		{name: "missing symbol", req: dto.PortfolioInitRequest{Cash: 100, TotalAssets: 100, Positions: []dto.PositionInput{{Name: "ETF", Quantity: 1, CostPrice: 1, CurrentPrice: 1, BuyReason: "低估配置"}}}},
		{name: "missing name", req: dto.PortfolioInitRequest{Cash: 99, TotalAssets: 100, Positions: []dto.PositionInput{{Symbol: "510300", Quantity: 1, CostPrice: 1, CurrentPrice: 1, BuyReason: "低估配置"}}}},
		{name: "negative quantity", req: dto.PortfolioInitRequest{Cash: 101, TotalAssets: 100, Positions: []dto.PositionInput{{Symbol: "510300", Name: "ETF", Quantity: -1, CostPrice: 1, CurrentPrice: 1, BuyReason: "低估配置"}}}},
		{name: "missing cost basis", req: dto.PortfolioInitRequest{Cash: 100, TotalAssets: 100, Positions: []dto.PositionInput{{Symbol: "510300", Name: "ETF", Quantity: 1, CurrentPrice: 1, BuyReason: "低估配置"}}}},
		{name: "missing buy reason", req: dto.PortfolioInitRequest{Cash: 99, TotalAssets: 100, Positions: []dto.PositionInput{{Symbol: "510300", Name: "ETF", Quantity: 1, CostPrice: 1, CurrentPrice: 1}}}},
		{name: "inconsistent total assets", req: dto.PortfolioInitRequest{Cash: 10, TotalAssets: 100, Positions: []dto.PositionInput{{Symbol: "510300", Name: "ETF", Quantity: 1, CostPrice: 1, CurrentPrice: 1, BuyReason: "低估配置"}}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			portfolioRepo := &portfolioRepoStub{}
			auditRepo := &auditRepoStub{}
			svc := &PortfolioService{
				tx:  transactorStub{repos: repository.Repositories{PortfolioRepo: portfolioRepo, AuditRepo: auditRepo}},
				clk: clock.FixedClock{Time: time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)},
				ids: idgen.NewFixedGenerator(map[string][]string{"snap": {"snap_invalid"}, "audit": {"audit_invalid"}, "pos": {"pos_invalid"}, "ps": {"ps_invalid"}}),
			}
			if _, err := svc.WriteSnapshot(context.Background(), "req_invalid", tc.req, "manual"); err == nil {
				t.Fatal("expected validation error")
			}
			if portfolioRepo.snapshot.SnapshotID != "" || len(portfolioRepo.positions) != 0 || auditRepo.event.AuditEventID != "" {
				t.Fatalf("validation failure wrote facts: snapshot=%+v positions=%+v audit=%+v", portfolioRepo.snapshot, portfolioRepo.positions, auditRepo.event)
			}
		})
	}
}

func TestPortfolioServiceRecordOfflineBuyUpdatesCashAndWritesFacts(t *testing.T) {
	decisionRepo := &decisionRepoStub{}
	portfolioRepo := &portfolioRepoStub{snapshot: repository.PortfolioSnapshot{SnapshotID: "snap_before", Cash: 100, TotalAssets: 100, CashRatio: 1}}
	auditRepo := &auditRepoStub{}
	svc := &PortfolioService{
		tx:  transactorStub{repos: repository.Repositories{DecisionRepo: decisionRepo, PortfolioRepo: portfolioRepo, AuditRepo: auditRepo}},
		clk: clock.FixedClock{Time: time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)},
		ids: idgen.NewFixedGenerator(map[string][]string{"confirm": {"confirm_offline"}, "tx": {"tx_offline"}, "snap": {"snap_offline"}, "audit": {"audit_offline"}, "pos": {"pos_offline"}, "ps": {"ps_offline"}}),
	}

	out, err := svc.RecordOfflineTransaction(context.Background(), "req_offline", dto.OfflineTransactionRequest{OperationType: "buy", Symbol: "510300", Name: "沪深300ETF", Quantity: 10, Price: 3, Fees: 1, ExecutedAt: "2026-05-29T03:00:00Z", BuyReason: "低估配置", Note: "仅记录线下成交"})
	if err != nil {
		t.Fatalf("RecordOfflineTransaction: %v", err)
	}
	if out.TransactionID != "tx_offline" || out.SnapshotID != "snap_offline" || out.SafetyStatement == "" {
		t.Fatalf("unexpected response: %+v", out)
	}
	if decisionRepo.confirmation.DecisionID != "" || decisionRepo.transaction.TransactionID != "tx_offline" {
		t.Fatalf("expected local confirmation without synthetic decision and transaction facts, got confirmation=%+v tx=%+v", decisionRepo.confirmation, decisionRepo.transaction)
	}
	if portfolioRepo.snapshot.Cash != 69 || portfolioRepo.snapshot.TotalAssets != 99 || portfolioRepo.position.Quantity != 10 || portfolioRepo.position.BuyReason != "低估配置" {
		t.Fatalf("expected cash-aware portfolio facts, snapshot=%+v position=%+v", portfolioRepo.snapshot, portfolioRepo.position)
	}
	if auditRepo.event.AuditEventID != "audit_offline" {
		t.Fatalf("expected audit event, got %+v", auditRepo.event)
	}
}

func TestPortfolioServiceEditAndRemoveHoldingPreserveCashInSnapshots(t *testing.T) {
	portfolioRepo := &portfolioRepoStub{snapshot: repository.PortfolioSnapshot{SnapshotID: "snap_before", Cash: 70, TotalAssets: 100, CashRatio: 0.7}, positions: []repository.Position{{PositionID: "pos_a", Symbol: "510300", Name: "沪深300ETF", Quantity: 10, CostPrice: 2, CurrentPrice: 3, MarketValue: 30, PositionState: string(model.PositionNormal), BuyReason: "低估配置"}}}
	auditRepo := &auditRepoStub{}
	svc := &PortfolioService{
		tx:  transactorStub{repos: repository.Repositories{PortfolioRepo: portfolioRepo, AuditRepo: auditRepo}},
		clk: clock.FixedClock{Time: time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)},
		ids: idgen.NewFixedGenerator(map[string][]string{"snap": {"snap_edit", "snap_remove"}, "audit": {"audit_edit", "audit_remove"}, "ps": {"ps_edit"}}),
	}

	_, err := svc.EditHolding(context.Background(), "req_edit", dto.HoldingEditRequest{PositionID: "pos_a", Reason: "本地校准", Confirmation: "confirmed", Position: dto.PositionInput{Symbol: "510300", Name: "沪深300ETF", Quantity: 8, CostPrice: 2, CurrentPrice: 4, BuyReason: "低估配置"}})
	if err != nil {
		t.Fatalf("EditHolding: %v", err)
	}
	if portfolioRepo.snapshot.Cash != 70 || portfolioRepo.snapshot.TotalAssets != 102 || portfolioRepo.position.Quantity != 8 {
		t.Fatalf("expected edit snapshot to preserve cash, snapshot=%+v position=%+v", portfolioRepo.snapshot, portfolioRepo.position)
	}

	_, err = svc.RemoveHolding(context.Background(), "req_remove", dto.HoldingRemoveRequest{PositionID: "pos_a", Reason: "清仓后校准", Confirmation: "confirmed"})
	if err != nil {
		t.Fatalf("RemoveHolding: %v", err)
	}
	if portfolioRepo.snapshot.Cash != 70 || portfolioRepo.snapshot.TotalAssets != 70 || len(portfolioRepo.positions) != 0 {
		t.Fatalf("expected remove snapshot to preserve cash and remove current holding, snapshot=%+v positions=%+v", portfolioRepo.snapshot, portfolioRepo.positions)
	}
}

func TestPortfolioServiceBatchImportAndCorrection(t *testing.T) {
	portfolioRepo := &portfolioRepoStub{}
	auditRepo := &auditRepoStub{}
	svc := &PortfolioService{
		tx:  transactorStub{repos: repository.Repositories{PortfolioRepo: portfolioRepo, AuditRepo: auditRepo}},
		clk: clock.FixedClock{Time: time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)},
		ids: idgen.NewFixedGenerator(map[string][]string{"import": {"import_mixed", "import_validated"}, "snap": {"snap_import"}, "audit": {"audit_import", "audit_corr"}, "pos": {"pos_import"}, "ps": {"ps_import"}, "corr": {"corr_fixed"}}),
	}

	validation, err := svc.ValidateImport(context.Background(), "req_validate", dto.BatchImportValidationRequest{Rows: []dto.BatchImportRow{{RowNumber: 1, RowType: "holding", Symbol: "510300", Name: "沪深300ETF", Quantity: 10, CostPrice: 2, CurrentPrice: 3, BuyReason: "低估配置"}, {RowNumber: 2, RowType: "holding", Symbol: "", Name: "坏数据", Quantity: 1, CostPrice: 1, CurrentPrice: 1}}})
	if err != nil {
		t.Fatalf("ValidateImport: %v", err)
	}
	if validation.Summary.ValidCount != 1 || validation.Summary.InvalidCount != 1 || len(portfolioRepo.positions) != 0 {
		t.Fatalf("expected validation only, validation=%+v positions=%+v", validation, portfolioRepo.positions)
	}

	validValidation, err := svc.ValidateImport(context.Background(), "req_validate_valid", dto.BatchImportValidationRequest{Rows: []dto.BatchImportRow{{RowNumber: 1, RowType: "holding", Symbol: "510300", Name: "沪深300ETF", Quantity: 10, CostPrice: 2, CurrentPrice: 3, BuyReason: "低估配置"}}})
	if err != nil {
		t.Fatalf("ValidateImport valid: %v", err)
	}
	out, err := svc.ConfirmImport(context.Background(), "req_import", dto.BatchImportConfirmRequest{ImportBatchID: validValidation.ImportBatchID, ConfirmReason: "导入初始持仓", Rows: []dto.BatchImportRow{{RowNumber: 1, RowType: "holding", Symbol: "510300", Name: "沪深300ETF", Quantity: 10, CostPrice: 2, CurrentPrice: 3, BuyReason: "低估配置"}}})
	if err != nil {
		t.Fatalf("ConfirmImport: %v", err)
	}
	if out.ImportBatchID != validValidation.ImportBatchID || portfolioRepo.importBatch.Status != "committed" || len(portfolioRepo.positions) != 1 || auditRepo.event.AuditEventID != "audit_import" {
		t.Fatalf("expected import facts, out=%+v batch=%+v positions=%+v audit=%+v", out, portfolioRepo.importBatch, portfolioRepo.positions, auditRepo.event)
	}

	corr, err := svc.CorrectFact(context.Background(), "req_corr", dto.CorrectionRequest{TargetType: "position", TargetID: "pos_import", BeforeJSON: `{"quantity":10}`, AfterJSON: `{"quantity":8}`, CorrectionReason: "录入数量修正"})
	if err != nil {
		t.Fatalf("CorrectFact: %v", err)
	}
	if corr.CorrectionID != "corr_fixed" || portfolioRepo.correction.TargetID != "pos_import" || auditRepo.event.AuditEventID != "audit_corr" {
		t.Fatalf("expected correction facts, corr=%+v stored=%+v audit=%+v", corr, portfolioRepo.correction, auditRepo.event)
	}
}

func TestPortfolioServiceValidateImportPersistsBatchMetadataWithoutAccountFacts(t *testing.T) {
	portfolioRepo := &portfolioRepoStub{}
	svc := &PortfolioService{
		tx:  transactorStub{repos: repository.Repositories{PortfolioRepo: portfolioRepo}},
		clk: clock.FixedClock{Time: time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)},
		ids: idgen.NewFixedGenerator(map[string][]string{"import": {"import_validated"}}),
	}

	out, err := svc.ValidateImport(context.Background(), "req_validate", dto.BatchImportValidationRequest{Rows: []dto.BatchImportRow{{RowNumber: 1, RowType: "holding", Symbol: "510300", Name: "沪深300ETF", Quantity: 10, CostPrice: 2, CurrentPrice: 3, BuyReason: "低估配置"}}})
	if err != nil {
		t.Fatalf("ValidateImport: %v", err)
	}
	if out.ImportBatchID != "import_validated" || portfolioRepo.importBatch.Status != "validated" || portfolioRepo.importBatch.RowCount != 1 {
		t.Fatalf("expected validated import batch metadata, out=%+v batch=%+v", out, portfolioRepo.importBatch)
	}
	if len(portfolioRepo.positions) != 0 || portfolioRepo.snapshot.SnapshotID != "" {
		t.Fatalf("validation must not write account facts, snapshot=%+v positions=%+v", portfolioRepo.snapshot, portfolioRepo.positions)
	}
}

func TestPortfolioServiceConfirmImportWritesTransactionRows(t *testing.T) {
	decisionRepo := &decisionRepoStub{}
	rows := []dto.BatchImportRow{{RowNumber: 1, RowType: "transaction", OperationType: "buy", Symbol: "510300", Name: "沪深300ETF", Quantity: 10, Price: 3, Fees: 1, OccurredAt: "2026-05-29T03:00:00Z", BuyReason: "低估配置"}}
	portfolioRepo := &portfolioRepoStub{snapshot: repository.PortfolioSnapshot{SnapshotID: "snap_before", Cash: 100, TotalAssets: 100, CashRatio: 1}, importBatch: repository.LocalAccountImportBatch{ImportBatchID: "batch_tx", Status: "validated", RowCount: 1, ValidCount: 1, InvalidCount: 0, RowsHash: importRowsHash(rows)}}
	auditRepo := &auditRepoStub{}
	svc := &PortfolioService{
		tx:  transactorStub{repos: repository.Repositories{DecisionRepo: decisionRepo, PortfolioRepo: portfolioRepo, AuditRepo: auditRepo}},
		clk: clock.FixedClock{Time: time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)},
		ids: idgen.NewFixedGenerator(map[string][]string{"import": {"import_validated"}, "confirm": {"confirm_import_tx"}, "tx": {"tx_import"}, "snap": {"snap_import_tx"}, "audit": {"audit_import"}, "pos": {"pos_import_tx"}, "ps": {"ps_import_tx"}}),
	}

	out, err := svc.ConfirmImport(context.Background(), "req_import_tx", dto.BatchImportConfirmRequest{ImportBatchID: "batch_tx", ConfirmReason: "导入线下交易", Rows: rows})
	if err != nil {
		t.Fatalf("ConfirmImport transaction row: %v", err)
	}
	if out.ImportBatchID != "batch_tx" || decisionRepo.transaction.TransactionID != "tx_import" || portfolioRepo.position.Symbol != "510300" || portfolioRepo.snapshot.Cash != 69 {
		t.Fatalf("expected transaction import to write transaction and account facts, out=%+v tx=%+v snapshot=%+v position=%+v", out, decisionRepo.transaction, portfolioRepo.snapshot, portfolioRepo.position)
	}
}

func TestPortfolioServiceReturnsLatestSnapshotErrors(t *testing.T) {
	portfolioRepo := &portfolioRepoStub{latestErr: errors.New("read latest failed")}
	svc := &PortfolioService{
		tx:  transactorStub{repos: repository.Repositories{PortfolioRepo: portfolioRepo, AuditRepo: &auditRepoStub{}}},
		clk: clock.FixedClock{Time: time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)},
		ids: idgen.NewFixedGenerator(map[string][]string{"snap": {"snap_edit"}, "audit": {"audit_edit"}}),
	}

	_, err := svc.EditHolding(context.Background(), "req_edit", dto.HoldingEditRequest{Reason: "校准", Confirmation: "confirmed", Position: dto.PositionInput{Symbol: "510300", Name: "沪深300ETF", Quantity: 1, CostPrice: 2, CurrentPrice: 3, BuyReason: "低估配置"}})
	if err == nil || err.Error() != "read latest failed" {
		t.Fatalf("expected latest snapshot read error, got %v", err)
	}
}

func TestPortfolioServiceRejectsInvalidCorrectionTargetType(t *testing.T) {
	svc := &PortfolioService{
		tx:  transactorStub{repos: repository.Repositories{PortfolioRepo: &portfolioRepoStub{}, AuditRepo: &auditRepoStub{}}},
		clk: clock.FixedClock{Time: time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)},
		ids: idgen.NewFixedGenerator(map[string][]string{"corr": {"corr_invalid"}, "audit": {"audit_invalid"}}),
	}

	_, err := svc.CorrectFact(context.Background(), "req_corr", dto.CorrectionRequest{TargetType: "bad_target", TargetID: "pos_1", BeforeJSON: `{}`, AfterJSON: `{}`, CorrectionReason: "修正"})
	if err == nil {
		t.Fatal("expected invalid target type to be rejected before persistence")
	}
}

func TestPortfolioServiceConfirmImportRequiresValidatedBatch(t *testing.T) {
	portfolioRepo := &portfolioRepoStub{}
	svc := &PortfolioService{
		tx:  transactorStub{repos: repository.Repositories{PortfolioRepo: portfolioRepo, AuditRepo: &auditRepoStub{}}},
		clk: clock.FixedClock{Time: time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)},
		ids: idgen.NewFixedGenerator(map[string][]string{"import": {"import_revalidated"}, "snap": {"snap_import"}, "audit": {"audit_import"}}),
	}

	_, err := svc.ConfirmImport(context.Background(), "req_import", dto.BatchImportConfirmRequest{ImportBatchID: "missing_batch", ConfirmReason: "确认导入", Rows: []dto.BatchImportRow{{RowNumber: 1, RowType: "holding", Symbol: "510300", Name: "沪深300ETF", Quantity: 10, CostPrice: 2, CurrentPrice: 3, BuyReason: "低估配置"}}})
	if err == nil {
		t.Fatal("expected missing validated import batch to be rejected")
	}
	if portfolioRepo.snapshot.SnapshotID != "" || len(portfolioRepo.positions) != 0 {
		t.Fatalf("rejected import wrote account facts: snapshot=%+v positions=%+v", portfolioRepo.snapshot, portfolioRepo.positions)
	}
}

func TestPortfolioServiceConfirmImportRejectsInvalidValidatedBatch(t *testing.T) {
	portfolioRepo := &portfolioRepoStub{importBatch: repository.LocalAccountImportBatch{ImportBatchID: "batch_invalid", Status: "validated", RowCount: 1, ValidCount: 0, InvalidCount: 1}}
	svc := &PortfolioService{
		tx:  transactorStub{repos: repository.Repositories{PortfolioRepo: portfolioRepo, AuditRepo: &auditRepoStub{}}},
		clk: clock.FixedClock{Time: time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)},
		ids: idgen.NewFixedGenerator(map[string][]string{"import": {"import_revalidated"}, "snap": {"snap_import"}, "audit": {"audit_import"}}),
	}

	_, err := svc.ConfirmImport(context.Background(), "req_import", dto.BatchImportConfirmRequest{ImportBatchID: "batch_invalid", ConfirmReason: "确认导入", Rows: []dto.BatchImportRow{{RowNumber: 1, RowType: "holding", Symbol: "510300", Name: "沪深300ETF", Quantity: 10, CostPrice: 2, CurrentPrice: 3, BuyReason: "低估配置"}}})
	if err == nil {
		t.Fatal("expected invalid validated import batch to be rejected")
	}
}

func TestPortfolioServiceConfirmImportRejectsValidatedBatchWithoutRowsHash(t *testing.T) {
	portfolioRepo := &portfolioRepoStub{importBatch: repository.LocalAccountImportBatch{ImportBatchID: "batch_no_hash", Status: "validated", RowCount: 1, ValidCount: 1, InvalidCount: 0}}
	svc := &PortfolioService{
		tx:  transactorStub{repos: repository.Repositories{PortfolioRepo: portfolioRepo, AuditRepo: &auditRepoStub{}}},
		clk: clock.FixedClock{Time: time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)},
		ids: idgen.NewFixedGenerator(map[string][]string{"snap": {"snap_import"}, "audit": {"audit_import"}}),
	}

	_, err := svc.ConfirmImport(context.Background(), "req_import", dto.BatchImportConfirmRequest{ImportBatchID: "batch_no_hash", ConfirmReason: "确认导入", Rows: []dto.BatchImportRow{{RowNumber: 1, RowType: "holding", Symbol: "510300", Name: "沪深300ETF", Quantity: 10, CostPrice: 2, CurrentPrice: 3, BuyReason: "低估配置"}}})
	if err == nil {
		t.Fatal("expected validated import batch without rows hash to be rejected")
	}
	if portfolioRepo.snapshot.SnapshotID != "" || len(portfolioRepo.positions) != 0 {
		t.Fatalf("rejected import wrote account facts: snapshot=%+v positions=%+v", portfolioRepo.snapshot, portfolioRepo.positions)
	}
}

func TestPortfolioServiceConfirmImportRejectsRowsChangedAfterValidation(t *testing.T) {
	portfolioRepo := &portfolioRepoStub{}
	svc := &PortfolioService{
		tx:  transactorStub{repos: repository.Repositories{PortfolioRepo: portfolioRepo, AuditRepo: &auditRepoStub{}}},
		clk: clock.FixedClock{Time: time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)},
		ids: idgen.NewFixedGenerator(map[string][]string{"import": {"import_validated"}, "snap": {"snap_import"}, "audit": {"audit_import"}}),
	}

	validation, err := svc.ValidateImport(context.Background(), "req_validate", dto.BatchImportValidationRequest{Rows: []dto.BatchImportRow{{RowNumber: 1, RowType: "holding", Symbol: "510300", Name: "沪深300ETF", Quantity: 10, CostPrice: 2, CurrentPrice: 3, BuyReason: "低估配置"}}})
	if err != nil {
		t.Fatalf("ValidateImport: %v", err)
	}
	_, err = svc.ConfirmImport(context.Background(), "req_import", dto.BatchImportConfirmRequest{ImportBatchID: validation.ImportBatchID, ConfirmReason: "确认导入", Rows: []dto.BatchImportRow{{RowNumber: 1, RowType: "holding", Symbol: "159915", Name: "创业板ETF", Quantity: 10, CostPrice: 2, CurrentPrice: 3, BuyReason: "低估配置"}}})
	if err == nil {
		t.Fatal("expected changed rows to be rejected")
	}
	if portfolioRepo.snapshot.SnapshotID != "" || len(portfolioRepo.positions) != 0 {
		t.Fatalf("rejected changed rows wrote account facts: snapshot=%+v positions=%+v", portfolioRepo.snapshot, portfolioRepo.positions)
	}
}

func TestConfirmationRejectsStaleStatusBeforeWritingFacts(t *testing.T) {
	decisionRepo := &decisionRepoStub{staleUpdate: true}
	portfolioRepo := &portfolioRepoStub{}
	auditRepo := &auditRepoStub{}
	svc := &ConfirmationService{
		tx:  transactorStub{repos: repository.Repositories{DecisionRepo: decisionRepo, PortfolioRepo: portfolioRepo, AuditRepo: auditRepo}},
		clk: clock.FixedClock{Time: time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)},
		ids: idgen.NewFixedGenerator(map[string][]string{"confirm": {"confirm_stale"}, "audit": {"audit_stale"}}),
	}

	_, err := svc.Confirm(context.Background(), "req_stale", "decision_stale", dto.ConfirmationRequest{ConfirmationType: string(model.ConfirmationTypePlanned)})
	if err == nil {
		t.Fatal("expected stale confirmation status to be rejected")
	}
	if decisionRepo.confirmation.ConfirmationID != "" {
		t.Fatalf("confirmation fact was written before stale status rejection: %+v", decisionRepo.confirmation)
	}
	if auditRepo.event.AuditEventID != "" {
		t.Fatalf("audit was written for rejected stale confirmation: %+v", auditRepo.event)
	}
}

func TestManualExecutionUpdatesCashAndTotalAssets(t *testing.T) {
	decisionRepo := &decisionRepoStub{}
	portfolioRepo := &portfolioRepoStub{snapshot: repository.PortfolioSnapshot{SnapshotID: "snap_before", Cash: 100, TotalAssets: 130, CashRatio: 100.0 / 130.0}, positions: []repository.Position{{PositionID: "pos_a", Symbol: "510300", Name: "沪深300", Quantity: 10, CostPrice: 2, CurrentPrice: 3, MarketValue: 30, PositionState: string(model.PositionNormal)}}}
	auditRepo := &auditRepoStub{}
	svc := &ConfirmationService{
		tx:  transactorStub{repos: repository.Repositories{DecisionRepo: decisionRepo, PortfolioRepo: portfolioRepo, AuditRepo: auditRepo}},
		clk: clock.FixedClock{Time: time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)},
		ids: idgen.NewFixedGenerator(map[string][]string{"confirm": {"confirm_buy"}, "audit": {"audit_buy"}, "tx": {"tx_buy"}, "snap_confirm": {"snap_buy"}, "ps": {"ps_a"}}),
	}

	_, err := svc.confirmWithStatus(context.Background(), "req_buy", "decision_buy", "formal_trade_advice", string(model.ConfirmationPending), dto.ConfirmationRequest{ConfirmationType: string(model.ConfirmationTypeExecutedManually), OperationType: "buy", Symbol: "510300", Quantity: 5, Price: 4, ExecutedAt: "2026-05-29T03:00:00Z"})
	if err != nil {
		t.Fatalf("Confirm buy: %v", err)
	}
	if portfolioRepo.snapshot.Cash != 80 || portfolioRepo.snapshot.TotalAssets != 140 || portfolioRepo.snapshot.CashRatio != 80.0/140.0 {
		t.Fatalf("expected cash-aware snapshot, got %+v", portfolioRepo.snapshot)
	}
}

func TestManualExecutionRejectsBuyWhenCashIsInsufficient(t *testing.T) {
	decisionRepo := &decisionRepoStub{}
	portfolioRepo := &portfolioRepoStub{snapshot: repository.PortfolioSnapshot{SnapshotID: "snap_before", Cash: 10, TotalAssets: 10, CashRatio: 1}}
	auditRepo := &auditRepoStub{}
	svc := &ConfirmationService{
		tx:  transactorStub{repos: repository.Repositories{DecisionRepo: decisionRepo, PortfolioRepo: portfolioRepo, AuditRepo: auditRepo}},
		clk: clock.FixedClock{Time: time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)},
		ids: idgen.NewFixedGenerator(map[string][]string{"confirm": {"confirm_buy"}, "audit": {"audit_buy"}, "tx": {"tx_buy"}, "snap_confirm": {"snap_buy"}, "pos": {"pos_buy"}, "ps": {"ps_buy"}}),
	}

	_, err := svc.confirmWithStatus(context.Background(), "req_buy", "decision_buy", "formal_trade_advice", string(model.ConfirmationPending), dto.ConfirmationRequest{ConfirmationType: string(model.ConfirmationTypeExecutedManually), OperationType: "buy", Symbol: "510300", Quantity: 5, Price: 4, ExecutedAt: "2026-05-29T03:00:00Z"})
	if err == nil {
		t.Fatal("expected insufficient cash error")
	}
	if decisionRepo.transaction.TransactionID != "" || portfolioRepo.snapshot.SnapshotID != "snap_before" || auditRepo.event.AuditEventID != "" {
		t.Fatalf("should not persist transaction, snapshot, or audit: tx=%+v snapshot=%+v audit=%+v", decisionRepo.transaction, portfolioRepo.snapshot, auditRepo.event)
	}
}

func TestConfirmationServiceUsesInjectedClockIDsAndTransactor(t *testing.T) {
	decisionRepo := &decisionRepoStub{}
	portfolioRepo := &portfolioRepoStub{}
	auditRepo := &auditRepoStub{}
	svc := &ConfirmationService{
		tx:  transactorStub{repos: repository.Repositories{DecisionRepo: decisionRepo, PortfolioRepo: portfolioRepo, AuditRepo: auditRepo}},
		clk: clock.FixedClock{Time: time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)},
		ids: idgen.NewFixedGenerator(map[string][]string{"confirm": {"confirm_fixed"}, "audit": {"audit_fixed"}, "tx": {"tx_fixed"}, "snap_confirm": {"snap_fixed"}, "pos": {"pos_fixed"}, "ps": {"ps_fixed"}}),
	}
	out, err := svc.confirmWithStatus(context.Background(), "req_confirm", "decision_1", "formal_trade_advice", string(model.ConfirmationPending), dto.ConfirmationRequest{ConfirmationType: string(model.ConfirmationTypeExecutedManually), OperationType: "buy", Symbol: "510300", Quantity: 2, Price: 3, ExecutedAt: "2026-05-29T03:00:00Z"})
	if err != nil {
		t.Fatalf("Confirm: %v", err)
	}
	if out.ConfirmationID != "confirm_fixed" || out.TransactionIDs[0] != "tx_fixed" || out.SnapshotID != "snap_fixed" || out.AuditEventIDs[0] != "audit_fixed" {
		t.Fatalf("unexpected response ids: %+v", out)
	}
	if decisionRepo.confirmation.CreatedAt != "2026-05-29T04:00:00Z" || decisionRepo.transaction.TransactionID != "tx_fixed" || auditRepo.event.CreatedAt != "2026-05-29T04:00:00Z" {
		t.Fatalf("unexpected persisted facts: confirmation=%+v tx=%+v audit=%+v", decisionRepo.confirmation, decisionRepo.transaction, auditRepo.event)
	}
}

func TestConfirmationServiceExecutedSellKeepsFullPortfolioSnapshot(t *testing.T) {
	decisionRepo := &decisionRepoStub{}
	portfolioRepo := &portfolioRepoStub{positions: []repository.Position{
		{PositionID: "pos_a", Symbol: "510300", Name: "沪深300", Quantity: 10, CostPrice: 2, CurrentPrice: 3, MarketValue: 30, PositionState: string(model.PositionNormal)},
		{PositionID: "pos_b", Symbol: "159915", Name: "创业板", Quantity: 5, CostPrice: 4, CurrentPrice: 4, MarketValue: 20, PositionState: string(model.PositionNormal)},
	}}
	auditRepo := &auditRepoStub{}
	svc := &ConfirmationService{
		tx:  transactorStub{repos: repository.Repositories{DecisionRepo: decisionRepo, PortfolioRepo: portfolioRepo, AuditRepo: auditRepo}},
		clk: clock.FixedClock{Time: time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)},
		ids: idgen.NewFixedGenerator(map[string][]string{"confirm": {"confirm_sell"}, "audit": {"audit_sell"}, "tx": {"tx_sell"}, "snap_confirm": {"snap_sell"}, "ps": {"ps_a", "ps_b"}}),
	}

	_, err := svc.confirmWithStatus(context.Background(), "req_sell", "decision_sell", "formal_trade_advice", string(model.ConfirmationPending), dto.ConfirmationRequest{ConfirmationType: string(model.ConfirmationTypeExecutedManually), OperationType: "sell", Symbol: "510300", Quantity: 4, Price: 3, ExecutedAt: "2026-05-29T03:00:00Z"})
	if err != nil {
		t.Fatalf("Confirm sell: %v", err)
	}
	if portfolioRepo.position.Quantity != 6 || portfolioRepo.position.MarketValue != 18 {
		t.Fatalf("expected reduced current position, got %+v", portfolioRepo.position)
	}
	if portfolioRepo.snapshot.PositionCount != 2 || portfolioRepo.snapshot.TotalAssets != 100050 || portfolioRepo.snapshot.Cash != 100012 || len(portfolioRepo.positionSnapshots) != 2 {
		t.Fatalf("expected full portfolio snapshot, snapshot=%+v positions=%+v", portfolioRepo.snapshot, portfolioRepo.positionSnapshots)
	}
	if decisionRepo.transaction.BeforePositionJSON == "{}" || decisionRepo.transaction.AfterPositionJSON == "{}" {
		t.Fatalf("expected transaction before/after position JSON, got %+v", decisionRepo.transaction)
	}
}

func TestConfirmationServiceRejectsSellWithoutEnoughPosition(t *testing.T) {
	decisionRepo := &decisionRepoStub{}
	portfolioRepo := &portfolioRepoStub{positions: []repository.Position{{PositionID: "pos_a", Symbol: "510300", Quantity: 2, CostPrice: 2, CurrentPrice: 3, MarketValue: 6, PositionState: string(model.PositionNormal)}}}
	auditRepo := &auditRepoStub{}
	svc := &ConfirmationService{
		tx:  transactorStub{repos: repository.Repositories{DecisionRepo: decisionRepo, PortfolioRepo: portfolioRepo, AuditRepo: auditRepo}},
		clk: clock.FixedClock{Time: time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)},
		ids: idgen.NewFixedGenerator(map[string][]string{"confirm": {"confirm_sell"}, "audit": {"audit_sell"}, "tx": {"tx_sell"}, "snap_confirm": {"snap_sell"}}),
	}

	_, err := svc.confirmWithStatus(context.Background(), "req_sell", "decision_sell", "formal_trade_advice", string(model.ConfirmationPending), dto.ConfirmationRequest{ConfirmationType: string(model.ConfirmationTypeExecutedManually), OperationType: "sell", Symbol: "510300", Quantity: 4, Price: 3, ExecutedAt: "2026-05-29T03:00:00Z"})
	if err == nil {
		t.Fatal("expected insufficient position error")
	}
	if decisionRepo.transaction.TransactionID != "" || portfolioRepo.snapshot.SnapshotID != "" {
		t.Fatalf("should not persist transaction or snapshot: tx=%+v snapshot=%+v", decisionRepo.transaction, portfolioRepo.snapshot)
	}
}

func TestRuleProposalServiceRunsGatekeeperAuditPath(t *testing.T) {
	ruleRepo := &ruleRepoStub{}
	auditRepo := &auditRepoStub{}
	svc := &RuleProposalService{
		tx:   transactorStub{repos: repository.Repositories{RuleRepo: ruleRepo, AuditRepo: auditRepo}},
		deps: workflow.NewWorkflowDependencies(repository.Repositories{RuleRepo: ruleRepo, AuditRepo: auditRepo}, transactorStub{repos: repository.Repositories{RuleRepo: ruleRepo, AuditRepo: auditRepo}}),
		clk:  clock.FixedClock{Time: time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)},
		ids:  idgen.NewFixedGenerator(map[string][]string{"audit": {"user_audit_fixed", "gatekeeper_audit_event_fixed"}, "gatekeeper": {"gatekeeper_fixed"}}),
	}
	out, err := svc.submitForGatekeeper(context.Background(), "req_rule", loadedProposal{Status: string(model.ProposalPendingUserConfirm), SampleCount: 3, AfterRuleJSON: "{}"}, "proposal_1")
	if err != nil {
		t.Fatalf("submitForGatekeeper: %v", err)
	}
	if out.Status != string(model.ProposalPendingFinalConfirm) || out.GatekeeperAuditID == "" || len(out.AuditEventIDs) != 1 {
		t.Fatalf("expected gatekeeper audit path, got %+v", out)
	}
	if ruleRepo.gatekeeper.GatekeeperAuditID == "" || !ruleRepo.gatekeeper.AllowApply || ruleRepo.gatekeeper.AuditResult != string(model.AuditApproved) || ruleRepo.gatekeeper.AuditedRuleVersion != "v_test" {
		t.Fatalf("service should create approved gatekeeper audit: %+v", ruleRepo.gatekeeper)
	}
	if auditRepo.event.NodeName != "AuditRecordNode" || ruleRepo.proposalStatus != string(model.ProposalPendingFinalConfirm) {
		t.Fatalf("unexpected persisted state: audit=%+v status=%s", auditRepo.event, ruleRepo.proposalStatus)
	}
}
