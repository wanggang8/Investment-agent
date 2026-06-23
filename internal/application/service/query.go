package service

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/infrastructure/config"
	"investment-agent/internal/pkg/apperr"
)

// QueryService handles read-only application queries for HTTP handlers.
type QueryService struct {
	repos              repository.Repositories
	dailyAutoRunConfig config.DailyAutoRunConfig
}

// NewQueryService creates a query service from repository interfaces.
func NewQueryService(repos repository.Repositories) *QueryService {
	return &QueryService{repos: repos}
}

// NewQueryServiceWithDailyAutoRunConfig creates a query service that uses the scheduler timezone for daily auto-run queries.
func NewQueryServiceWithDailyAutoRunConfig(repos repository.Repositories, cfg config.DailyAutoRunConfig) *QueryService {
	return &QueryService{repos: repos, dailyAutoRunConfig: cfg}
}

func (s *QueryService) LatestPortfolioSnapshot(ctx context.Context) (repository.PortfolioSnapshot, error) {
	return s.repos.PortfolioRepo.GetLatestPortfolioSnapshot(ctx)
}

func (s *QueryService) PortfolioSnapshot(ctx context.Context, snapshotID string) (repository.PortfolioSnapshot, []repository.PositionSnapshot, error) {
	if snapshotID != "" {
		return s.repos.PortfolioRepo.GetPortfolioSnapshot(ctx, snapshotID)
	}
	snapshot, err := s.repos.PortfolioRepo.GetLatestPortfolioSnapshot(ctx)
	if err != nil {
		return snapshot, nil, err
	}
	fullSnapshot, positions, err := s.repos.PortfolioRepo.GetPortfolioSnapshot(ctx, snapshot.SnapshotID)
	if err != nil {
		return snapshot, nil, err
	}
	return fullSnapshot, positions, nil
}

func (s *QueryService) ListPositions(ctx context.Context) ([]repository.Position, error) {
	return s.repos.PortfolioRepo.ListPositions(ctx)
}

func (s *QueryService) LatestMarketSnapshot(ctx context.Context) (model.MarketSnapshot, error) {
	return s.repos.MarketRepo.GetLatestMarketSnapshot(ctx)
}

func (s *QueryService) LatestMarketSnapshotBySymbol(ctx context.Context, symbol string) (model.MarketSnapshot, error) {
	return s.repos.MarketRepo.GetLatestMarketSnapshotBySymbol(ctx, symbol)
}

func (s *QueryService) ListDecisions(ctx context.Context) ([]repository.DecisionRecord, error) {
	return s.repos.DecisionRepo.ListDecisionRecords(ctx)
}

func (s *QueryService) ListEvidenceSummaries(ctx context.Context) ([]repository.IntelligenceSummary, error) {
	return s.repos.IntelligenceRepo.ListEvidenceSummaries(ctx)
}

func (s *QueryService) LatestSourceVerification(ctx context.Context) (repository.SourceVerification, error) {
	return s.repos.IntelligenceRepo.GetLatestSourceVerification(ctx)
}

func (s *QueryService) LatestSourceVerificationByFilter(ctx context.Context, symbol, eventID string) (repository.SourceVerification, error) {
	if symbol == "" && eventID == "" {
		return s.repos.IntelligenceRepo.GetLatestSourceVerification(ctx)
	}
	return s.repos.IntelligenceRepo.GetLatestSourceVerificationByFilter(ctx, symbol, eventID)
}

func (s *QueryService) SystemStatus(ctx context.Context, vectorConfigured bool, deepSeekConfigured bool) dto.SystemStatusDTO {
	status := dto.SystemStatusDTO{SQLiteStatus: "ok", SQLitePath: "local", VecLiteStatus: "missing", DeepSeekStatus: "missing", DataSources: []string{"official", "exchange"}, LogLevel: "info"}
	if vectorConfigured {
		status.VecLiteStatus = "configured"
	}
	if deepSeekConfigured {
		status.DeepSeekStatus = "configured"
	}
	settings, err := s.repos.SettingsRepo.GetLatestSystemSettings(ctx)
	if err == nil && strings.TrimSpace(settings.DataSourcesJSON) != "" {
		var sources []string
		if json.Unmarshal([]byte(settings.DataSourcesJSON), &sources) == nil && len(sources) > 0 {
			status.DataSources = sources
		}
	}
	return status
}

func (s *QueryService) ActiveRuleVersion(ctx context.Context) (repository.RuleVersion, error) {
	return s.repos.RuleRepo.GetActiveRuleVersion(ctx)
}

func (s *QueryService) ListRuleProposals(ctx context.Context) ([]repository.RuleProposalWithAudit, error) {
	return s.repos.RuleRepo.ListRuleProposals(ctx)
}

func (s *QueryService) RuleProposal(ctx context.Context, proposalID string) (repository.RuleProposal, error) {
	return s.repos.RuleRepo.GetRuleProposal(ctx, proposalID)
}

func (s *QueryService) LatestRuleEffectValidationByProposal(ctx context.Context, proposalID string) (repository.RuleEffectValidation, error) {
	if s.repos.RuleEffectRepo == nil {
		return repository.RuleEffectValidation{}, apperr.New(apperr.CodeNotFound, apperr.CategoryNotFound, "rule effect validation not found")
	}
	items, err := s.repos.RuleEffectRepo.ListRuleEffectValidations(ctx, repository.RuleEffectValidationFilter{ProposalID: proposalID})
	if err != nil {
		return repository.RuleEffectValidation{}, err
	}
	if len(items) == 0 {
		return repository.RuleEffectValidation{}, apperr.New(apperr.CodeNotFound, apperr.CategoryNotFound, "rule effect validation not found")
	}
	return items[0], nil
}

func (s *QueryService) ListRuleEffectTracking(ctx context.Context, filter repository.RuleEffectTrackingFilter) ([]repository.RuleEffectTracking, error) {
	if s.repos.RuleEffectRepo == nil {
		return nil, nil
	}
	return s.repos.RuleEffectRepo.ListRuleEffectTracking(ctx, filter)
}

func (s *QueryService) ListAuditEvents(ctx context.Context) ([]repository.AuditEvent, error) {
	return s.repos.AuditRepo.ListAuditEvents(ctx)
}

func (s *QueryService) LatestDailyAutoRunState(ctx context.Context) (repository.DailyAutoRunState, error) {
	return s.repos.DailyAutoRunRepo.GetLatestDailyAutoRunState(ctx)
}

func (s *QueryService) LatestCapabilityConfig(ctx context.Context) (repository.CapabilityConfig, error) {
	return s.repos.SettingsRepo.GetLatestCapabilityConfig(ctx)
}

func (s *QueryService) ReviewCounts(ctx context.Context) (int, int, error) {
	summary, err := s.ReviewSummary(ctx, "monthly")
	if err != nil {
		return 0, 0, err
	}
	return summary.ErrorCaseCount, summary.RuleProposalCount, nil
}

// ReviewSummary 聚合本地复盘事实，所有统计均来自现有事实表，避免形成第二套复盘来源。
func (s *QueryService) ReviewSummary(ctx context.Context, period string) (dto.ReviewSummaryResponse, error) {
	if period != "quarterly" {
		period = "monthly"
	}
	decisions, err := s.repos.DecisionRepo.ListDecisionRecords(ctx)
	if err != nil {
		return dto.ReviewSummaryResponse{}, err
	}
	proposals, err := s.repos.RuleRepo.ListRuleProposals(ctx)
	if err != nil {
		return dto.ReviewSummaryResponse{}, err
	}
	var ruleEffectTracking []repository.RuleEffectTracking
	if s.repos.RuleEffectRepo != nil {
		ruleEffectTracking, err = s.repos.RuleEffectRepo.ListRuleEffectTracking(ctx, repository.RuleEffectTrackingFilter{Period: period})
		if err != nil {
			return dto.ReviewSummaryResponse{}, err
		}
	}
	audits, err := s.repos.AuditRepo.ListAuditEvents(ctx)
	if err != nil {
		return dto.ReviewSummaryResponse{}, err
	}
	errorCases, err := s.repos.DecisionRepo.ListErrorCases(ctx)
	if err != nil {
		return dto.ReviewSummaryResponse{}, err
	}
	from := reviewPeriodStart(period, time.Now().UTC())
	decisions = filterDecisionsByCreatedAt(decisions, from)
	proposals = filterProposalsByCreatedAt(proposals, from)
	audits = filterAuditsByCreatedAt(audits, from)
	errorCases = filterErrorCasesByCreatedAt(errorCases, from)

	summary := dto.ReviewSummaryResponse{
		Period:                period,
		DecisionCount:         len(decisions),
		ErrorCaseCount:        len(errorCases),
		MisjudgmentCount:      len(errorCases),
		RuleProposalCount:     len(proposals),
		AuditEventCount:       len(audits),
		RecentDecisions:       reviewDecisionItems(decisions),
		AttributionSummaries:  reviewAttributions(decisions),
		RecurringErrorTags:    reviewErrorTags(errorCases),
		MissingEvidenceThemes: reviewEvidenceThemes(decisions),
		RuleProposalOutcomes:  reviewProposalOutcomes(proposals),
		DegradedWorkflows:     reviewDegradedWorkflows(decisions),
		OpsStatus:             reviewOpsStatus(decisions, audits),
	}
	for _, item := range decisions {
		summary.RuleHitCount += triggeredRuleCount(item.TriggeredRulesJSON)
		if item.SourceVerificationStatus == string(model.VerificationFailed) || item.SourceVerificationStatus == string(model.VerificationBackgroundOnly) || item.FinalVerdictStatus == string(model.VerdictInsufficientData) {
			summary.MissingEvidenceCount++
		}
		if item.WorkflowStatus == string(model.WorkflowDegraded) {
			summary.DegradedCount++
			summary.OpsStatus.DataSourceStatus = "degraded"
			summary.OpsStatus.ReviewStatus = "degraded"
			summary.OpsStatus.Explanation = "复盘窗口内存在降级记录，请查看追踪入口。"
		}
		switch item.ConfirmationStatus {
		case string(model.ConfirmationExecutedManually):
			summary.ConfirmationCount++
			summary.ExecutedManuallyCount++
		case string(model.ConfirmationPlanned):
			summary.ConfirmationCount++
			summary.PlannedCount++
		case string(model.ConfirmationWatch), string(model.ConfirmationMarkedError):
			summary.ConfirmationCount++
		}
	}
	for _, audit := range audits {
		if audit.AuditEventID != "" && len(summary.TrackingLinks) < 3 {
			summary.TrackingLinks = append(summary.TrackingLinks, dto.ReviewTrackingLink{Type: "audit_event", ID: audit.AuditEventID, Label: "审计事件 " + audit.AuditEventID})
		}
	}
	for _, proposal := range proposals {
		summary.RuleSuggestions = append(summary.RuleSuggestions, dto.RuleSuggestionDTO{ProposalID: proposal.ProposalID, Title: proposal.Title, Status: proposal.Status, Reason: proposal.Reason, CanAutoApply: false})
		if len(summary.TrackingLinks) < 3 {
			summary.TrackingLinks = append(summary.TrackingLinks, dto.ReviewTrackingLink{Type: "rule_proposal", ID: proposal.ProposalID, Label: "规则提案 " + proposal.ProposalID})
		}
	}
	if summary.ErrorCaseCount > 0 && len(summary.TrackingLinks) < 3 {
		summary.TrackingLinks = append(summary.TrackingLinks, dto.ReviewTrackingLink{Type: "error_case", ID: "error_cases", Label: "错误案例"})
	}
	for _, tracking := range ruleEffectTracking {
		summary.RuleEffectTracking = append(summary.RuleEffectTracking, dto.RuleEffectTrackingDTO{TrackingID: tracking.TrackingID, AppliedRuleVersion: tracking.AppliedRuleVersion, ProposalID: tracking.ProposalID, Period: tracking.Period, HitCount: tracking.HitCount, MisjudgmentCount: tracking.MisjudgmentCount, MissingEvidenceCount: tracking.MissingEvidenceCount, DegradedCount: tracking.DegradedCount, RiskAlertCount: tracking.RiskAlertCount, TrendDirection: string(tracking.TrendDirection), Metrics: parseJSONValue(tracking.MetricsJSON), RelatedProposalIDs: parseJSONValue(tracking.RelatedProposalIDsJSON), RelatedAuditEventIDs: parseJSONValue(tracking.RelatedAuditEventIDsJSON), RelatedRiskAlertIDs: parseJSONValue(tracking.RelatedRiskAlertIDsJSON), SafetyNote: tracking.SafetyNote})
		if len(summary.TrackingLinks) < 3 {
			summary.TrackingLinks = append(summary.TrackingLinks, dto.ReviewTrackingLink{Type: "rule_effect_tracking", ID: tracking.TrackingID, Label: "规则效果追踪 " + tracking.AppliedRuleVersion})
		}
	}
	return summary, nil
}

func reviewPeriodStart(period string, now time.Time) time.Time {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	if period == "quarterly" {
		return now.AddDate(0, -3, 0)
	}
	return now.AddDate(0, -1, 0)
}

func filterDecisionsByCreatedAt(items []repository.DecisionRecord, from time.Time) []repository.DecisionRecord {
	out := make([]repository.DecisionRecord, 0, len(items))
	for _, item := range items {
		if includeReviewCreatedAt(item.CreatedAt, from) {
			out = append(out, item)
		}
	}
	return out
}

func filterProposalsByCreatedAt(items []repository.RuleProposalWithAudit, from time.Time) []repository.RuleProposalWithAudit {
	out := make([]repository.RuleProposalWithAudit, 0, len(items))
	for _, item := range items {
		if includeReviewCreatedAt(item.CreatedAt, from) {
			out = append(out, item)
		}
	}
	return out
}

func filterAuditsByCreatedAt(items []repository.AuditEvent, from time.Time) []repository.AuditEvent {
	out := make([]repository.AuditEvent, 0, len(items))
	for _, item := range items {
		if includeReviewCreatedAt(item.CreatedAt, from) {
			out = append(out, item)
		}
	}
	return out
}

func filterErrorCasesByCreatedAt(items []repository.ErrorCase, from time.Time) []repository.ErrorCase {
	out := make([]repository.ErrorCase, 0, len(items))
	for _, item := range items {
		if includeReviewCreatedAt(item.CreatedAt, from) {
			out = append(out, item)
		}
	}
	return out
}

func includeReviewCreatedAt(createdAt string, from time.Time) bool {
	parsed, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return true
	}
	return !parsed.Before(from)
}

func triggeredRuleCount(raw string) int {
	if strings.TrimSpace(raw) == "" || raw == "[]" {
		return 0
	}
	var values []any
	if err := json.Unmarshal([]byte(raw), &values); err == nil {
		return len(values)
	}
	return 1
}

func reviewDecisionItems(decisions []repository.DecisionRecord) []dto.DecisionListItem {
	out := make([]dto.DecisionListItem, 0, len(decisions))
	for _, item := range decisions {
		out = append(out, dto.DecisionListItem{DecisionID: item.DecisionID, Symbol: item.Symbol, FinalVerdict: item.FinalVerdictStatus, ConfirmationStatus: item.ConfirmationStatus, GeneratedAt: item.CreatedAt})
	}
	return out
}

func reviewOpsStatus(decisions []repository.DecisionRecord, audits []repository.AuditEvent) dto.ReviewOpsStatus {
	if len(decisions) == 0 && len(audits) == 0 {
		return dto.ReviewOpsStatus{DataSourceStatus: "unknown", IndexStatus: "unknown", ReviewStatus: "empty", Explanation: "复盘窗口内暂无本地事实，无法推断运行状态。"}
	}
	status := dto.ReviewOpsStatus{DataSourceStatus: "success", IndexStatus: "success", ReviewStatus: "success", Explanation: "复盘摘要来自本地事实表与 API DTO。"}
	for _, audit := range audits {
		if audit.Status == string(model.AuditStatusFailed) {
			status.ReviewStatus = "degraded"
			status.Explanation = "复盘窗口内存在失败审计事件，请查看追踪入口。"
		}
		if audit.Action == string(model.AuditActionRefreshMarketData) && audit.Status == string(model.AuditStatusFailed) {
			status.DataSourceStatus = "degraded"
		}
		if audit.Action == string(model.AuditActionRebuildIndex) && audit.Status == string(model.AuditStatusFailed) {
			status.IndexStatus = "degraded"
		}
	}
	return status
}

func reviewAttributions(decisions []repository.DecisionRecord) []dto.ReviewAttribution {
	out := make([]dto.ReviewAttribution, 0, len(decisions))
	for _, item := range decisions {
		evidenceStatus := item.SourceVerificationStatus
		if evidenceStatus == "" {
			evidenceStatus = "unknown"
		}
		out = append(out, dto.ReviewAttribution{DecisionID: item.DecisionID, Symbol: item.Symbol, Verdict: item.FinalVerdictStatus, ConfirmationStatus: item.ConfirmationStatus, EvidenceStatus: evidenceStatus, WorkflowStatus: item.WorkflowStatus, Outcome: reviewAttributionOutcome(item)})
	}
	return out
}

func reviewAttributionOutcome(item repository.DecisionRecord) string {
	if item.FinalVerdictStatus == string(model.VerdictInsufficientData) || item.SourceVerificationStatus == string(model.VerificationFailed) || item.SourceVerificationStatus == string(model.VerificationBackgroundOnly) {
		return "missing_evidence"
	}
	if item.WorkflowStatus == string(model.WorkflowDegraded) {
		return "degraded"
	}
	if item.ConfirmationStatus != "" {
		return item.ConfirmationStatus
	}
	return "recorded"
}

func reviewErrorTags(errorCases []repository.ErrorCase) []dto.ReviewErrorTag {
	counts := make(map[string]int)
	for _, item := range errorCases {
		if strings.TrimSpace(item.RootCauseTag) != "" {
			counts[item.RootCauseTag]++
		}
	}
	out := make([]dto.ReviewErrorTag, 0, len(counts))
	for tag, count := range counts {
		out = append(out, dto.ReviewErrorTag{Tag: tag, Count: count})
	}
	return out
}

func reviewEvidenceThemes(decisions []repository.DecisionRecord) []dto.ReviewEvidenceTheme {
	counts := make(map[string]int)
	for _, item := range decisions {
		if item.SourceVerificationStatus == string(model.VerificationFailed) || item.SourceVerificationStatus == string(model.VerificationBackgroundOnly) || item.FinalVerdictStatus == string(model.VerdictInsufficientData) {
			status := item.SourceVerificationStatus
			if status == "" {
				status = "unknown"
			}
			counts[status]++
		}
	}
	out := make([]dto.ReviewEvidenceTheme, 0, len(counts))
	for status, count := range counts {
		out = append(out, dto.ReviewEvidenceTheme{Status: status, Count: count})
	}
	return out
}

func reviewProposalOutcomes(proposals []repository.RuleProposalWithAudit) []dto.ReviewRuleProposalState {
	out := make([]dto.ReviewRuleProposalState, 0, len(proposals))
	for _, item := range proposals {
		out = append(out, dto.ReviewRuleProposalState{ProposalID: item.ProposalID, Title: item.Title, Status: item.Status, AuditResult: item.AuditResult})
	}
	return out
}

func reviewDegradedWorkflows(decisions []repository.DecisionRecord) []dto.ReviewDegradedWorkflow {
	out := make([]dto.ReviewDegradedWorkflow, 0)
	for _, item := range decisions {
		if item.WorkflowStatus == string(model.WorkflowDegraded) {
			out = append(out, dto.ReviewDegradedWorkflow{DecisionID: item.DecisionID, Symbol: item.Symbol, Status: item.WorkflowStatus, CreatedAt: item.CreatedAt})
		}
	}
	return out
}

func (s *QueryService) ListRiskAlerts(ctx context.Context, filter repository.RiskAlertFilter) (dto.PageResult[dto.RiskAlertDTO], error) {
	if s.repos.RiskAlertRepo == nil {
		return dto.PageResult[dto.RiskAlertDTO]{Items: []dto.RiskAlertDTO{}}, nil
	}
	alerts, err := s.repos.RiskAlertRepo.ListRiskAlerts(ctx, filter)
	if err != nil {
		return dto.PageResult[dto.RiskAlertDTO]{}, err
	}
	items := make([]dto.RiskAlertDTO, 0, len(alerts))
	for _, alert := range alerts {
		items = append(items, riskAlertDTO(alert))
	}
	return dto.PageResult[dto.RiskAlertDTO]{Items: items, Total: len(items)}, nil
}

func (s *QueryService) GetRiskAlert(ctx context.Context, alertID string) (dto.RiskAlertDTO, error) {
	if s.repos.RiskAlertRepo == nil {
		return dto.RiskAlertDTO{}, apperr.New(apperr.CodeInternalError, apperr.CategoryInternal, "risk alert repository not configured")
	}
	alert, err := s.repos.RiskAlertRepo.GetRiskAlert(ctx, alertID)
	if err != nil {
		return dto.RiskAlertDTO{}, err
	}
	return riskAlertDTO(alert), nil
}

func riskAlertDTO(alert repository.RiskAlert) dto.RiskAlertDTO {
	out := dto.RiskAlertDTO{AlertID: alert.AlertID, RiskType: string(alert.RiskType), Severity: string(alert.Severity), SOPStatus: string(alert.SOPStatus), Symbol: alert.Symbol, TriggerSummary: alert.TriggerSummary, TriggerContext: parseJSONValue(alert.TriggerContextJSON), ProhibitedActions: stringListFromJSON(alert.ProhibitedActionsJSON), SuggestedActions: stringListFromJSON(alert.SuggestedActionsJSON), RelatedDecisionID: alert.RelatedDecisionID, RelatedReportID: alert.RelatedReportID, RelatedNotificationID: alert.RelatedNotificationID, RelatedAuditEventID: alert.RelatedAuditEventID, Link: "/risk-alerts/" + url.PathEscape(alert.AlertID), LastTriggeredAt: alert.LastTriggeredAt, ResolvedAt: alert.ResolvedAt, ResolutionReason: alert.ResolutionReason, SafetyNote: "风险预警只用于本地人工复核，不会自动交易，也不会调用券商接口。", CreatedAt: alert.CreatedAt, UpdatedAt: alert.UpdatedAt}
	if alert.RelatedDecisionID != "" {
		out.DecisionLink = "/decisions/" + url.PathEscape(alert.RelatedDecisionID)
	}
	if alert.RelatedReportID != "" {
		out.ReportLink = "/daily-discipline/reports/" + url.PathEscape(alert.RelatedReportID)
	}
	if alert.RelatedNotificationID != "" {
		out.NotificationLink = "/notifications?source_id=" + url.QueryEscape(alert.RelatedNotificationID)
	}
	if alert.RelatedAuditEventID != "" {
		out.AuditLink = "/audit?event_id=" + url.QueryEscape(alert.RelatedAuditEventID)
	}
	return out
}

func parseJSONValue(raw string) any {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var out any
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	return out
}

func (s *QueryService) TodayDailyDisciplineReport(ctx context.Context, now time.Time) (dto.DailyDisciplineReportResponse, error) {
	if now.IsZero() {
		now = time.Now()
	}
	localDate := dailyDisciplineLocalDate(now, s.dailyAutoRunConfig.Timezone)
	state, stateErr := s.repos.DailyAutoRunRepo.GetLatestDailyAutoRunState(ctx)
	if stateErr == nil && state.LocalDate == localDate && state.SymbolSetHash != "" {
		report, err := s.repos.DailyDisciplineReportRepo.GetDailyDisciplineReportByKey(ctx, localDate, "holdings", state.SymbolSetHash)
		if err == nil {
			return s.dailyDisciplineReportDTO(ctx, report, nil)
		}
		if !apperr.IsCode(err, apperr.CodeNotFound) {
			return dto.DailyDisciplineReportResponse{}, err
		}
		return s.latestSameDayReportOrState(ctx, localDate, state)
	}

	reports, err := s.repos.DailyDisciplineReportRepo.ListDailyDisciplineReports(ctx, repository.DailyDisciplineReportListFilter{Limit: 30})
	if err != nil {
		return dto.DailyDisciplineReportResponse{}, err
	}
	for _, report := range reports {
		if report.LocalDate == localDate {
			return s.dailyDisciplineReportDTO(ctx, report, nil)
		}
	}
	if stateErr == nil && state.LocalDate == localDate {
		return s.dailyDisciplineReportFromAutoRunState(ctx, state)
	}
	if stateErr != nil && !apperr.IsCode(stateErr, apperr.CodeNotFound) {
		return dto.DailyDisciplineReportResponse{}, stateErr
	}
	out := dto.DailyDisciplineReportResponse{LocalDate: localDate, Scope: "holdings", Status: "not_started", Summary: "今日尚未生成每日纪律报告。", MissingAction: "补齐本地账户与持仓，并等待每日自动运行或手动执行本地纪律检查。", MissingCategories: missingPrerequisiteCategories(), SafetyNote: dailyDisciplineReportSafetyNote}
	out.Trend, err = s.dailyDisciplineReportTrend(ctx)
	if err != nil {
		return dto.DailyDisciplineReportResponse{}, err
	}
	return out, nil
}

func (s *QueryService) latestSameDayReportOrState(ctx context.Context, localDate string, state repository.DailyAutoRunState) (dto.DailyDisciplineReportResponse, error) {
	reports, err := s.repos.DailyDisciplineReportRepo.ListDailyDisciplineReports(ctx, repository.DailyDisciplineReportListFilter{Limit: 30})
	if err != nil {
		return dto.DailyDisciplineReportResponse{}, err
	}
	for _, report := range reports {
		if report.LocalDate == localDate && report.UpdatedAt >= state.UpdatedAt {
			return s.dailyDisciplineReportDTO(ctx, report, nil)
		}
	}
	return s.dailyDisciplineReportFromAutoRunState(ctx, state)
}

func dailyDisciplineLocalDate(now time.Time, timezone string) string {
	loc := time.UTC
	if strings.TrimSpace(timezone) != "" {
		loaded, err := time.LoadLocation(timezone)
		if err == nil {
			loc = loaded
		}
	}
	return now.In(loc).Format("2006-01-02")
}

func (s *QueryService) ListDailyDisciplineReports(ctx context.Context, status string, limit int) (dto.DailyDisciplineReportListResponse, error) {
	limit = normalizeDailyDisciplineReportLimit(limit)
	reports, err := s.repos.DailyDisciplineReportRepo.ListDailyDisciplineReports(ctx, repository.DailyDisciplineReportListFilter{Status: status, Limit: limit})
	if err != nil {
		return dto.DailyDisciplineReportListResponse{}, err
	}
	out := dto.DailyDisciplineReportListResponse{Reports: make([]dto.DailyDisciplineReportResponse, 0, len(reports))}
	trend, err := s.dailyDisciplineReportTrend(ctx)
	if err != nil {
		return dto.DailyDisciplineReportListResponse{}, err
	}
	for _, report := range reports {
		out.Reports = append(out.Reports, dailyDisciplineReportListItemDTO(report, trend))
	}
	return out, nil
}

func (s *QueryService) GetDailyDisciplineReport(ctx context.Context, reportID string) (dto.DailyDisciplineReportResponse, error) {
	report, err := s.repos.DailyDisciplineReportRepo.GetDailyDisciplineReport(ctx, reportID)
	if err != nil {
		return dto.DailyDisciplineReportResponse{}, err
	}
	return s.dailyDisciplineReportDTO(ctx, report, nil)
}

const dailyDisciplineReportSafetyNote = "每日纪律报告只用于本地记录和人工复核，不会自动执行交易。"

func normalizeDailyDisciplineReportLimit(limit int) int {
	if limit <= 0 {
		return 30
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func dailyDisciplineReportListItemDTO(report repository.DailyDisciplineReport, trend dto.DailyDisciplineReportTrend) dto.DailyDisciplineReportResponse {
	out := dto.DailyDisciplineReportResponse{ReportID: report.ReportID, LocalDate: report.LocalDate, Scope: report.Scope, Status: report.Status, Summary: report.Summary, SourceType: report.SourceType, SourceID: report.SourceID, DecisionID: report.DecisionID, FailureCode: report.FailureCode, FailureReason: report.FailureReason, SafetyNote: dailyDisciplineReportSafetyNote, UpdatedAt: report.UpdatedAt, Trend: trend}
	applyDailyDisciplineReportLinks(&out)
	return out
}

func (s *QueryService) dailyDisciplineReportDTO(ctx context.Context, report repository.DailyDisciplineReport, trend *dto.DailyDisciplineReportTrend) (dto.DailyDisciplineReportResponse, error) {
	out := dto.DailyDisciplineReportResponse{ReportID: report.ReportID, LocalDate: report.LocalDate, Scope: report.Scope, Status: report.Status, Summary: report.Summary, SourceType: report.SourceType, SourceID: report.SourceID, DecisionID: report.DecisionID, FailureCode: report.FailureCode, FailureReason: report.FailureReason, SafetyNote: dailyDisciplineReportSafetyNote, UpdatedAt: report.UpdatedAt}
	applyDailyDisciplineReportLinks(&out)
	if report.DecisionID != "" {
		decision, refs, err := s.repos.DecisionRepo.GetDecisionRecord(ctx, report.DecisionID)
		if err != nil {
			return dto.DailyDisciplineReportResponse{}, err
		}
		out.FinalVerdict = decision.FinalVerdictText
		out.VerdictStatus = decision.FinalVerdictStatus
		out.Evidence.EvidenceCount = len(refs)
		for _, ref := range refs {
			if ref.IndependentSourceCount > out.Evidence.IndependentSourceCount {
				out.Evidence.IndependentSourceCount = ref.IndependentSourceCount
			}
			if ref.HighGradeIndependentSourceCount > out.Evidence.HighGradeIndependentSourceCount {
				out.Evidence.HighGradeIndependentSourceCount = ref.HighGradeIndependentSourceCount
			}
		}
		out.P34SourceCoverage = dailyDisciplineP34Coverage(decision.ExpectedReturnScenariosJSON)
		out.RiskAlerts, err = s.dailyDisciplineRiskAlerts(ctx, report)
		if err != nil {
			return dto.DailyDisciplineReportResponse{}, err
		}
	}
	if trend == nil {
		calculated, err := s.dailyDisciplineReportTrend(ctx)
		if err != nil {
			return dto.DailyDisciplineReportResponse{}, err
		}
		trend = &calculated
	}
	out.Trend = *trend
	return out, nil
}

func dailyDisciplineP34Coverage(raw string) dto.DailyDisciplineReportP34Coverage {
	if strings.TrimSpace(raw) == "" {
		return dto.DailyDisciplineReportP34Coverage{}
	}
	var payload struct {
		SupportingDataSummary string                 `json:"supporting_data_summary"`
		MissingCategories     []string               `json:"missing_categories"`
		SourceHealth          []dto.SourceHealthItem `json:"source_health"`
	}
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return dto.DailyDisciplineReportP34Coverage{Summary: "P34 扩展数据状态不可解析", MissingCategories: []string{"p34_expanded_data"}}
	}
	return dto.DailyDisciplineReportP34Coverage{Summary: payload.SupportingDataSummary, MissingCategories: payload.MissingCategories, SourceHealth: payload.SourceHealth}
}

func (s *QueryService) dailyDisciplineRiskAlerts(ctx context.Context, report repository.DailyDisciplineReport) ([]dto.DailyDisciplineReportRiskAlert, error) {
	if s.repos.RiskAlertRepo == nil {
		return nil, nil
	}
	alerts, err := s.repos.RiskAlertRepo.ListRiskAlerts(ctx, repository.RiskAlertFilter{SOPStatuses: []model.RiskSOPStatus{model.RiskSOPTriggered, model.RiskSOPActive, model.RiskSOPObserving, model.RiskSOPEscalated}})
	if err != nil {
		return nil, err
	}
	out := make([]dto.DailyDisciplineReportRiskAlert, 0, len(alerts))
	for _, alert := range alerts {
		if alert.RelatedReportID != report.ReportID && alert.RelatedDecisionID != report.DecisionID {
			continue
		}
		out = append(out, dailyDisciplineRiskAlertDTO(alert))
	}
	return out, nil
}

func dailyDisciplineRiskAlertDTO(alert repository.RiskAlert) dto.DailyDisciplineReportRiskAlert {
	return dto.DailyDisciplineReportRiskAlert{AlertID: alert.AlertID, RiskType: string(alert.RiskType), Severity: string(alert.Severity), SOPStatus: string(alert.SOPStatus), Symbol: alert.Symbol, TriggerSummary: alert.TriggerSummary, ProhibitedActions: stringListFromJSON(alert.ProhibitedActionsJSON), SuggestedActions: stringListFromJSON(alert.SuggestedActionsJSON), Link: "/risk-alerts/" + url.PathEscape(alert.AlertID), SafetyNote: "风险预警只用于本地人工复核，不会自动执行交易。"}
}

func stringListFromJSON(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var items []string
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return nil
	}
	return items
}

func (s *QueryService) dailyDisciplineReportFromAutoRunState(ctx context.Context, state repository.DailyAutoRunState) (dto.DailyDisciplineReportResponse, error) {
	status := state.Status
	if state.FailureCode == "missing_prerequisites" {
		status = "insufficient_data"
	}
	out := dto.DailyDisciplineReportResponse{LocalDate: state.LocalDate, Scope: state.Scope, Status: status, Summary: state.FailureReason, SourceType: "auto_run", SourceID: state.IdempotencyKey, FailureCode: state.FailureCode, FailureReason: state.FailureReason, SafetyNote: dailyDisciplineReportSafetyNote, UpdatedAt: state.UpdatedAt}
	if out.Summary == "" {
		out.Summary = state.Status
	}
	if state.FailureCode == "missing_prerequisites" {
		out.MissingAction = "补齐本地账户与持仓后再生成每日纪律报告。"
		out.MissingCategories = missingPrerequisiteCategories()
	}
	applyDailyDisciplineReportLinks(&out)
	trend, err := s.dailyDisciplineReportTrend(ctx)
	if err != nil {
		return dto.DailyDisciplineReportResponse{}, err
	}
	out.Trend = trend
	return out, nil
}

func missingPrerequisiteCategories() []string {
	return []string{"account", "holdings"}
}

func applyDailyDisciplineReportLinks(out *dto.DailyDisciplineReportResponse) {
	if out.DecisionID != "" {
		out.DecisionLink = "/decisions/" + out.DecisionID
	}
	if out.SourceType == "auto_run" {
		out.AutoRunLink = "/daily-auto-run"
	}
	if out.SourceID != "" {
		escapedSourceID := url.QueryEscape(out.SourceID)
		out.AuditLink = "/audit?input_ref=" + escapedSourceID
		out.NotificationLink = "/notifications?source_id=" + escapedSourceID
	}
}

func (s *QueryService) dailyDisciplineReportTrend(ctx context.Context) (dto.DailyDisciplineReportTrend, error) {
	reports, err := s.repos.DailyDisciplineReportRepo.ListDailyDisciplineReports(ctx, repository.DailyDisciplineReportListFilter{Limit: 30})
	if err != nil {
		return dto.DailyDisciplineReportTrend{}, err
	}
	var trend dto.DailyDisciplineReportTrend
	for _, report := range reports {
		switch report.Status {
		case "success":
			trend.SuccessCount++
		case "degraded":
			trend.DegradedCount++
		case "failed":
			trend.FailedCount++
		case "insufficient_data":
			trend.InsufficientDataCount++
		}
	}
	return trend, nil
}

func (s *QueryService) DataRequiredError(message string) error {
	return apperr.New(apperr.CodeDataRequired, apperr.CategoryConflict, message)
}
