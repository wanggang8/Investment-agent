package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/application/workflow"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

// ConsultDecision 同步执行主动咨询工作流，并返回可渲染的决策详情。
func (a *App) ConsultDecision(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	var req dto.ConsultDecisionRequest
	if err := decodeJSON(r, &req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	if req.Scenario != "" && !req.Scenario.Valid() {
		WriteHandlerError(w, requestID, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "scenario 不合法"))
		return
	}
	s, snapshotPositions, err := a.QuerySvc.PortfolioSnapshot(r.Context(), "")
	if err != nil || s.SnapshotID == "" {
		WriteHandlerError(w, requestID, apperr.New(apperr.CodeDataRequired, apperr.CategoryConflict, "需要先录入账户快照"))
		return
	}
	market, err := a.QuerySvc.LatestMarketSnapshotBySymbol(r.Context(), req.Symbol)
	if err != nil || market.MarketSnapshotID == "" {
		WriteHandlerError(w, requestID, apperr.New(apperr.CodeDataRequired, apperr.CategoryConflict, "需要先刷新市场快照"))
		return
	}
	activeRule, err := a.QuerySvc.ActiveRuleVersion(r.Context())
	if err != nil || activeRule.RuleVersion == "" {
		WriteHandlerError(w, requestID, apperr.New(apperr.CodeRuleVersionMissing, apperr.CategoryConflict, "规则版本缺失"))
		return
	}
	capabilityStatus, capabilityReason := a.capabilityStatusForSymbol(r.Context(), req.Symbol)
	positions := positionsFromSnapshots(snapshotPositions)
	wf := workflow.WorkflowContext{RequestID: requestID, UserQuestion: req.Question, Symbol: req.Symbol, PortfolioSnapshot: model.PortfolioSnapshot{SnapshotID: s.SnapshotID, Cash: s.Cash, TotalAssets: s.TotalAssets, CashRatio: s.CashRatio, HighRiskRatio: s.HighRiskRatio, PositionCount: s.PositionCount}, PositionSnapshots: positions, MarketSnapshot: market, RuleVersion: activeRule.RuleVersion, CapabilityStatus: capabilityStatus, CapabilityReason: capabilityReason, EvidenceSet: model.EvidenceSet{VerificationStatus: model.VerificationFailed}, ExpectedReturnSampleCount: workflow.ExpectedReturnSampleCountFromWorkflowData(positions, market), ExpectedReturnPreviousBaseMidpoint: req.ExpectedReturnPreviousBaseMidpoint, ExpectedReturnTargetReturnRate: req.ExpectedReturnTargetReturnRate}
	out, err := workflow.NewConsultationGraphWithDependencies(a.Deps).Run(r.Context(), wf)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, decisionDetailFromWorkflow(out))
}

// GetDecision 返回已持久化的决策详情。
func (a *App) GetDecision(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	record, refs, err := a.Deps.DecisionRepo.GetDecisionRecord(r.Context(), r.PathValue("decision_id"))
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	var snapshot *dto.AccountSnapshot
	if record.PortfolioSnapshotID != "" {
		if got, _, err := a.QuerySvc.PortfolioSnapshot(r.Context(), record.PortfolioSnapshotID); err == nil {
			snapshot = accountSnapshotDTO(model.PortfolioSnapshot{SnapshotID: got.SnapshotID, Cash: got.Cash, TotalAssets: got.TotalAssets, CashRatio: got.CashRatio, HighRiskRatio: got.HighRiskRatio})
		}
	}
	writeOK(w, requestID, decisionDetailFromRecord(record, refs, snapshot))
}

// ListDecisions 返回决策记录列表。
func (a *App) ListDecisions(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	from, to, err := parseDateRange(r)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	records, err := a.QuerySvc.ListDecisions(r.Context())
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	status := r.URL.Query().Get("status")
	items := make([]dto.DecisionListItem, 0, len(records))
	for _, record := range records {
		if status != "" && record.ConfirmationStatus != status {
			continue
		}
		if !includeDateRange(record.CreatedAt, from, to) {
			continue
		}
		items = append(items, dto.DecisionListItem{DecisionID: record.DecisionID, DisplayTitle: record.Symbol + " 决策", Symbol: record.Symbol, FinalVerdict: record.FinalVerdictStatus, TriggeredRuleIDs: splitJSONStrings(record.TriggeredRulesJSON), ConfirmationStatus: record.ConfirmationStatus, GeneratedAt: record.CreatedAt})
	}
	writeOK(w, requestID, dto.PageResult[dto.DecisionListItem]{Items: items, Total: len(items)})
}

func parseDateRange(r *http.Request) (time.Time, time.Time, error) {
	var from time.Time
	var to time.Time
	var err error
	if value := r.URL.Query().Get("from"); value != "" {
		from, err = parseDateParam(value)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}
	if value := r.URL.Query().Get("to"); value != "" {
		to, err = parseDateParam(value)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}
	if !from.IsZero() && !to.IsZero() && from.After(to) {
		return time.Time{}, time.Time{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "from 不能晚于 to")
	}
	return from, to, nil
}

func includeDateRange(createdAt string, from, to time.Time) bool {
	parsed, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		parsed, err = parseDateParam(createdAt)
		if err != nil {
			return true
		}
	}
	if !from.IsZero() && parsed.Before(from) {
		return false
	}
	if !to.IsZero() && parsed.After(to.AddDate(0, 0, 1).Add(-time.Nanosecond)) {
		return false
	}
	return true
}

func positionsFromSnapshots(positions []repository.PositionSnapshot) []model.Position {
	out := make([]model.Position, 0, len(positions))
	for _, p := range positions {
		out = append(out, model.Position{PositionID: p.PositionID, Symbol: p.Symbol, Name: p.Name, Quantity: p.Quantity, CostPrice: p.CostPrice, CurrentPrice: p.CurrentPrice, MarketValue: p.MarketValue, UnrealizedProfitRatio: p.UnrealizedProfitRatio, PositionState: model.PositionState(p.PositionState), AssetTag: p.AssetTag})
	}
	return out
}

func (a *App) capabilityStatusForSymbol(ctx context.Context, symbol string) (string, string) {
	cfg, err := a.QuerySvc.LatestCapabilityConfig(ctx)
	if err != nil {
		return workflow.CapabilityUnknown, "能力圈配置缺失"
	}
	var excluded []string
	_ = json.Unmarshal([]byte(cfg.ExcludedSymbolsJSON), &excluded)
	for _, item := range excluded {
		if item == symbol {
			return workflow.CapabilityOutOfScope, "标的不在能力圈"
		}
	}
	var allowed []string
	_ = json.Unmarshal([]byte(cfg.SymbolsJSON), &allowed)
	if len(allowed) == 0 {
		return workflow.CapabilityUnknown, "能力圈未配置具体标的"
	}
	for _, item := range allowed {
		if item == symbol {
			return workflow.CapabilityInScope, "标的在能力圈内"
		}
	}
	return workflow.CapabilityOutOfScope, "标的不在能力圈"
}

func decisionDetailFromWorkflow(wf workflow.WorkflowContext) dto.DecisionDetailResponse {
	confirmationStatus := confirmationStatusForDecision(wf.RuleVerdict.Status)
	return dto.DecisionDetailResponse{DecisionID: wf.DecisionID, Question: wf.UserQuestion, Symbol: wf.Symbol, GeneratedAt: nowRFC3339(), CapabilityCheck: capabilityCheckDTO(wf.CapabilityStatus, wf.CapabilityReason), WorkflowStatus: workflowStatusFromErrors(wf.Errors), AccountSnapshot: accountSnapshotDTO(wf.PortfolioSnapshot), TriggeredRules: triggeredRulesDTO(wf.RuleVerdict.TriggeredRules), EvidenceChain: evidenceDTOs(evidenceRefsFromWorkflow(wf)), AnalystReports: analystReportsDTOWithMetadata(wf.AnalystReports, wf.AnalystReportMetadata), RetrievalQuality: retrievalQualityDTO(wf.RetrievalQualitySummary), MarketContext: marketContextDTO(wf.MarketSnapshot, wf.Symbol), ExpectedReturnScenarios: expectedReturnDTO(expectedReturnDetail{
		Status:                wf.ExpectedReturnPrecisionStatus,
		Reason:                wf.ExpectedReturnReason,
		SampleCount:           wf.ExpectedReturnSampleCount,
		TargetName:            wf.ExpectedReturnTargetName,
		TargetCode:            wf.ExpectedReturnTargetCode,
		HoldingClass:          wf.ExpectedReturnHoldingClass,
		HorizonLabel:          wf.ExpectedReturnHorizonLabel,
		SampleWindow:          wf.ExpectedReturnSampleWindow,
		ScreeningCondition:    wf.ExpectedReturnScreeningCondition,
		ProbabilityBasis:      wf.ExpectedReturnProbabilityBasis,
		SupportingDataSummary: wf.ExpectedReturnSupportingDataSummary,
		MissingCategories:     wf.ExpectedReturnMissingCategories,
		SupplementData:        wf.ExpectedReturnSupplementData,
		AssumptionChecks:      wf.ExpectedReturnAssumptionChecks,
		HistoricalContexts:    wf.ExpectedReturnHistoricalContexts,
		HoldingClassCoverage:  wf.ExpectedReturnHoldingClassCoverage,
		Scenarios:             wf.ExpectedReturnScenarios,
		SellEvaluation:        wf.ExpectedReturnSellEvaluation,
		ReassessmentTrigger:   wf.ExpectedReturnReassessmentTrigger,
	}), ArbitrationChain: arbitrationChainDTO(wf.RuleVerdict.TriggeredRules), FinalVerdict: dto.FinalVerdict{Status: string(wf.RuleVerdict.Status), DisplayText: wf.RuleVerdict.Text, ProhibitedActions: wf.RuleVerdict.ProhibitedActions, OptionalActions: wf.RuleVerdict.OptionalActions}, UserConfirmation: dto.UserConfirmation{ConfirmationStatus: confirmationStatus, AvailableActions: availableConfirmationActions(confirmationStatus)}}
}

func workflowStatusFromErrors(errors []string) string {
	status := model.WorkflowCompleted
	for _, code := range errors {
		switch code {
		case workflow.ErrCodeDataRequired, workflow.ErrCodeDataStale, workflow.ErrCodeRuleVersionMissing, workflow.ErrCodeEvidenceNotFound, workflow.ErrCodeSourceVerificationFailed, workflow.ErrCodeDecisionRecordFailed:
			return string(model.WorkflowFailed)
		case workflow.ErrCodeVectorIndexUnavailable, workflow.ErrCodeAnalystUnavailable:
			status = model.WorkflowDegraded
		}
	}
	return string(status)
}

func confirmationStatusForDecision(status model.FinalVerdictStatus) string {
	switch status {
	case model.VerdictBuyAllowed, model.VerdictHold, model.VerdictReduce, model.VerdictSellOnly:
		return string(model.ConfirmationPending)
	default:
		return string(model.ConfirmationNotRequired)
	}
}

func availableConfirmationActions(status string) []string {
	if status != string(model.ConfirmationPending) {
		return []string{}
	}
	return []string{string(model.ConfirmationTypePlanned), string(model.ConfirmationTypeWatch), string(model.ConfirmationTypeExecutedManually), string(model.ConfirmationTypeMarkedError)}
}

func triggeredRulesDTO(rules []model.TriggeredRule) []dto.TriggeredRuleDTO {
	out := make([]dto.TriggeredRuleDTO, 0, len(rules))
	for _, rule := range rules {
		out = append(out, dto.TriggeredRuleDTO{RuleID: rule.RuleID, RuleName: rule.RuleName, Severity: rule.Severity, Description: rule.Description})
	}
	return out
}

func capabilityCheckDTO(status, reason string) *dto.CapabilityCheck {
	if status == "" && reason == "" {
		return nil
	}
	return &dto.CapabilityCheck{Status: status, Reason: reason}
}

func accountSnapshotDTO(snapshot model.PortfolioSnapshot) *dto.AccountSnapshot {
	if snapshot.SnapshotID == "" {
		return nil
	}
	return &dto.AccountSnapshot{SnapshotID: snapshot.SnapshotID, Cash: snapshot.Cash, TotalAssets: snapshot.TotalAssets, CashRatio: snapshot.CashRatio, HighRiskRatio: snapshot.HighRiskRatio}
}

func accountSnapshotFromRecord(record repository.DecisionRecord) *dto.AccountSnapshot {
	if record.PortfolioSnapshotID == "" {
		return nil
	}
	return &dto.AccountSnapshot{SnapshotID: record.PortfolioSnapshotID}
}

func analystReportsDTO(reports map[string]string) []dto.AnalystReport {
	return analystReportsDTOWithMetadata(reports, nil)
}

func analystReportsDTOWithMetadata(reports map[string]string, metadata map[string]map[string]string) []dto.AnalystReport {
	out := make([]dto.AnalystReport, 0, len(reports))
	for agent, conclusion := range reports {
		report := dto.AnalystReport{AgentName: agent, Conclusion: conclusion, KeyReasons: []string{}, RiskWarnings: []string{}, Confidence: "qualitative", EvidenceIDs: []string{}}
		if values := metadata[agent]; values != nil {
			report.PromptVersion = values["prompt_version"]
			report.Model = values["model"]
			report.InputSummary = values["input_summary"]
			report.OutputSummary = values["output_summary"]
			report.ParseStatus = values["parse_status"]
			report.QualityStatus = values["quality_status"]
		}
		out = append(out, report)
	}
	return out
}

func analystReportsFromJSON(raw string) []dto.AnalystReport {
	if raw == "" {
		return []dto.AnalystReport{}
	}
	var structured []dto.AnalystReport
	if err := json.Unmarshal([]byte(raw), &structured); err == nil && structured != nil {
		return normalizeAnalystReports(structured)
	}
	var mapped map[string]string
	if err := json.Unmarshal([]byte(raw), &mapped); err != nil {
		return []dto.AnalystReport{}
	}
	return analystReportsDTO(mapped)
}

func normalizeAnalystReports(reports []dto.AnalystReport) []dto.AnalystReport {
	out := make([]dto.AnalystReport, 0, len(reports))
	for _, report := range reports {
		if report.KeyReasons == nil {
			report.KeyReasons = []string{}
		}
		if report.RiskWarnings == nil {
			report.RiskWarnings = []string{}
		}
		if report.EvidenceIDs == nil {
			report.EvidenceIDs = []string{}
		}
		out = append(out, report)
	}
	return out
}

func retrievalQualityDTO(summary workflow.RetrievalQualitySummary) *dto.RetrievalQualitySummary {
	if summary.QuerySummary == "" && summary.TopK == 0 && summary.Status == "" && summary.IndexHealth == "" && summary.IndexFreshness == "" && summary.FallbackSource == "" && summary.SourceConsistencyStatus == "" && summary.DegradedReason == "" {
		return nil
	}
	return &dto.RetrievalQualitySummary{QuerySummary: summary.QuerySummary, TopK: summary.TopK, Status: summary.Status, IndexHealth: summary.IndexHealth, IndexFreshness: summary.IndexFreshness, FallbackSource: summary.FallbackSource, SourceConsistencyStatus: summary.SourceConsistencyStatus, DegradedReason: summary.DegradedReason}
}

type expectedReturnDetail struct {
	Status                model.PrecisionStatus
	Reason                string
	SampleCount           int
	TargetName            string
	TargetCode            string
	HoldingClass          string
	HorizonLabel          string
	SampleWindow          string
	ScreeningCondition    string
	ProbabilityBasis      string
	SupportingDataSummary string
	MissingCategories     []string
	SupplementData        []string
	AssumptionChecks      []workflow.ExpectedReturnAssumptionCheck
	HistoricalContexts    []workflow.ExpectedReturnHistoricalContext
	HoldingClassCoverage  []workflow.ExpectedReturnHoldingClassCoverage
	Scenarios             []workflow.ExpectedReturnScenario
	SellEvaluation        workflow.ExpectedReturnSellEvaluation
	ReassessmentTrigger   workflow.ExpectedReturnReassessmentTrigger
}

func expectedReturnDTO(input expectedReturnDetail) *dto.ExpectedReturnScenarios {
	if input.Status == "" && input.Reason == "" && input.SampleCount == 0 && input.TargetName == "" && input.TargetCode == "" && input.HorizonLabel == "" && input.SampleWindow == "" && input.ScreeningCondition == "" && input.ProbabilityBasis == "" && len(input.Scenarios) == 0 && input.SellEvaluation.Status == "" && input.ReassessmentTrigger.Reason == "" {
		return nil
	}
	out := &dto.ExpectedReturnScenarios{SampleCount: input.SampleCount, TargetName: input.TargetName, TargetCode: input.TargetCode, HoldingClass: input.HoldingClass, HorizonLabel: input.HorizonLabel, SampleWindow: input.SampleWindow, ScreeningCondition: input.ScreeningCondition, PrecisionStatus: string(input.Status), ProbabilityBasis: input.ProbabilityBasis, Scenarios: []dto.ReturnScenario{}, Reason: input.Reason, SupportingDataSummary: input.SupportingDataSummary, MissingCategories: input.MissingCategories, SupplementData: input.SupplementData, AssumptionChecks: assumptionCheckDTOs(input.AssumptionChecks), HistoricalContexts: historicalContextDTOs(input.HistoricalContexts), HoldingClassCoverage: holdingClassCoverageDTOs(input.HoldingClassCoverage), Disclaimer: "预期收益仅为情景分析，不构成收益承诺。"}
	for _, item := range input.Scenarios {
		out.Scenarios = append(out.Scenarios, dto.ReturnScenario{Scenario: item.Name, ReturnRange: firstNonEmptyString(item.ReturnRange, fmt.Sprintf("%.2f%%", item.ReturnRate*100)), Probability: item.Probability, Trigger: item.Trigger})
	}
	if input.SellEvaluation.Status != "" {
		out.SellEvaluation = &dto.SellEvaluation{Status: input.SellEvaluation.Status, Triggers: input.SellEvaluation.Triggers, Prompts: input.SellEvaluation.Prompts, Actions: input.SellEvaluation.Actions, NonTradingDisclaimer: input.SellEvaluation.NonTradingDisclaimer}
	}
	if input.ReassessmentTrigger.Reason != "" {
		out.ReassessmentTrigger = &dto.ReassessmentTrigger{Reason: input.ReassessmentTrigger.Reason, Boundary: input.ReassessmentTrigger.Boundary, CurrentValue: input.ReassessmentTrigger.CurrentValue}
	}
	return out
}

func assumptionCheckDTOs(items []workflow.ExpectedReturnAssumptionCheck) []dto.AssumptionCheck {
	if len(items) == 0 {
		return nil
	}
	out := make([]dto.AssumptionCheck, 0, len(items))
	for _, item := range items {
		out = append(out, dto.AssumptionCheck{Name: item.Name, Expected: item.Expected, Actual: item.Actual, MonthsBelow: item.MonthsBelow})
	}
	return out
}

func holdingClassCoverageDTOs(items []workflow.ExpectedReturnHoldingClassCoverage) []dto.HoldingClassCoverage {
	if len(items) == 0 {
		return nil
	}
	out := make([]dto.HoldingClassCoverage, 0, len(items))
	for _, item := range items {
		out = append(out, dto.HoldingClassCoverage{HoldingClass: item.HoldingClass, Symbol: item.Symbol, Status: item.Status})
	}
	return out
}

func historicalContextDTOs(items []workflow.ExpectedReturnHistoricalContext) []dto.HistoricalContext {
	if len(items) == 0 {
		return nil
	}
	out := make([]dto.HistoricalContext, 0, len(items))
	for _, item := range items {
		out = append(out, dto.HistoricalContext{Label: item.Label, Window: item.Window, SampleCount: item.SampleCount, Outcome: item.Outcome, MaxDrawdown: item.MaxDrawdown, Recovery: item.Recovery, Source: item.Source})
	}
	return out
}

func expectedReturnFromJSON(raw string) *dto.ExpectedReturnScenarios {
	if raw == "" {
		return nil
	}
	var stored struct {
		PrecisionStatus       string                                        `json:"precision_status"`
		Reason                string                                        `json:"reason"`
		SampleCount           int                                           `json:"sample_count"`
		TargetName            string                                        `json:"target_name"`
		TargetCode            string                                        `json:"target_code"`
		HoldingClass          string                                        `json:"holding_class"`
		HorizonLabel          string                                        `json:"horizon_label"`
		SampleWindow          string                                        `json:"sample_window"`
		ScreeningCondition    string                                        `json:"screening_condition"`
		ProbabilityBasis      string                                        `json:"probability_basis"`
		SupportingDataSummary string                                        `json:"supporting_data_summary"`
		MissingCategories     []string                                      `json:"missing_categories"`
		SupplementData        []string                                      `json:"supplement_data"`
		AssumptionChecks      []workflow.ExpectedReturnAssumptionCheck      `json:"assumption_checks"`
		HistoricalContexts    []workflow.ExpectedReturnHistoricalContext    `json:"historical_contexts"`
		HoldingClassCoverage  []workflow.ExpectedReturnHoldingClassCoverage `json:"holding_class_coverage"`
		Scenarios             []storedExpectedReturnScenario                `json:"scenarios"`
		SellEvaluation        workflow.ExpectedReturnSellEvaluation         `json:"sell_evaluation"`
		ReassessmentTrigger   workflow.ExpectedReturnReassessmentTrigger    `json:"reassessment_trigger"`
	}
	if err := json.Unmarshal([]byte(raw), &stored); err != nil {
		return nil
	}
	scenarios := make([]workflow.ExpectedReturnScenario, 0, len(stored.Scenarios))
	for _, item := range stored.Scenarios {
		probability := item.Probability
		if probability == nil {
			probability = item.WorkflowProbability
		}
		returnRate := item.ReturnRate
		if returnRate == 0 {
			returnRate = item.WorkflowReturnRate
		}
		scenarios = append(scenarios, workflow.ExpectedReturnScenario{Name: firstNonEmptyString(item.Name, item.Scenario, item.WorkflowName), ReturnRate: returnRate, ReturnRange: firstNonEmptyString(item.ReturnRange, item.WorkflowReturnRange), Probability: probability, Trigger: firstNonEmptyString(item.Trigger, item.WorkflowTrigger)})
	}
	out := expectedReturnDTO(expectedReturnDetail{Status: model.PrecisionStatus(stored.PrecisionStatus), Reason: stored.Reason, SampleCount: stored.SampleCount, TargetName: stored.TargetName, TargetCode: stored.TargetCode, HoldingClass: stored.HoldingClass, HorizonLabel: stored.HorizonLabel, SampleWindow: stored.SampleWindow, ScreeningCondition: stored.ScreeningCondition, ProbabilityBasis: stored.ProbabilityBasis, SupportingDataSummary: stored.SupportingDataSummary, MissingCategories: stored.MissingCategories, SupplementData: stored.SupplementData, AssumptionChecks: stored.AssumptionChecks, HistoricalContexts: stored.HistoricalContexts, HoldingClassCoverage: stored.HoldingClassCoverage, Scenarios: scenarios, SellEvaluation: stored.SellEvaluation, ReassessmentTrigger: stored.ReassessmentTrigger})
	return out
}

type storedExpectedReturnScenario struct {
	Name                string   `json:"name"`
	Scenario            string   `json:"scenario"`
	ReturnRate          float64  `json:"return_rate"`
	ReturnRange         string   `json:"return_range"`
	Probability         *float64 `json:"probability"`
	Trigger             string   `json:"trigger"`
	WorkflowName        string   `json:"Name"`
	WorkflowReturnRate  float64  `json:"ReturnRate"`
	WorkflowReturnRange string   `json:"ReturnRange"`
	WorkflowProbability *float64 `json:"Probability"`
	WorkflowTrigger     string   `json:"Trigger"`
}

func arbitrationChainDTO(rules []model.TriggeredRule) []dto.ArbitrationStep {
	out := make([]dto.ArbitrationStep, 0, len(rules))
	for i, rule := range rules {
		out = append(out, dto.ArbitrationStep{Priority: i + 1, RuleID: rule.RuleID, Result: rule.Description})
	}
	return out
}

func arbitrationChainFromJSON(raw, triggeredRulesRaw string) []dto.ArbitrationStep {
	if raw != "" {
		var out []dto.ArbitrationStep
		if err := json.Unmarshal([]byte(raw), &out); err == nil && out != nil {
			return out
		}
	}
	return arbitrationChainDTO(triggeredRulesModelFromJSON(triggeredRulesRaw))
}

func triggeredRulesFromJSON(raw string) []dto.TriggeredRuleDTO {
	return triggeredRulesDTO(triggeredRulesModelFromJSON(raw))
}

func triggeredRulesModelFromJSON(raw string) []model.TriggeredRule {
	if raw == "" {
		return nil
	}
	var rules []model.TriggeredRule
	if err := json.Unmarshal([]byte(raw), &rules); err != nil {
		return nil
	}
	return rules
}

func evidenceDTOs(refs []repository.EvidenceRef) []dto.EvidenceDTO {
	out := make([]dto.EvidenceDTO, 0, len(refs))
	for _, ref := range refs {
		out = append(out, evidenceDTO(ref))
	}
	return out
}

func evidenceRefsFromWorkflow(wf workflow.WorkflowContext) []repository.EvidenceRef {
	refs := make([]repository.EvidenceRef, 0, len(wf.EvidenceSet.Items))
	for _, item := range wf.EvidenceSet.Items {
		refs = append(refs, repository.EvidenceRef{EvidenceID: item.EvidenceID, DecisionID: wf.DecisionID, SummaryID: firstNonEmptyString(item.SummaryID, item.EvidenceID), SourceName: item.SourceName, SourceLevel: string(item.SourceLevel), EvidenceRole: string(item.Role), PublishedAt: item.PublishedAt, CapturedAt: item.CapturedAt, OriginalURL: item.OriginalURL, Summary: item.Summary, ContentHash: item.ContentHash, TimeWeight: item.TimeWeight, RelevanceScore: item.RelevanceScore, IndependentSourceCount: item.IndependentSourceCount, HighGradeIndependentSourceCount: item.HighGradeIndependentSourceCount})
	}
	return refs
}

func decisionDetailFromRecord(record repository.DecisionRecord, refs []repository.EvidenceRef, snapshot *dto.AccountSnapshot) dto.DecisionDetailResponse {
	evidence := evidenceDTOs(refs)
	if snapshot == nil {
		snapshot = accountSnapshotFromRecord(record)
	}
	return dto.DecisionDetailResponse{DecisionID: record.DecisionID, Question: record.Question, Symbol: record.Symbol, GeneratedAt: record.CreatedAt, CapabilityCheck: capabilityCheckDTO(record.CapabilityStatus, record.CapabilityReason), WorkflowStatus: record.WorkflowStatus, AccountSnapshot: snapshot, EvidenceChain: evidence, TriggeredRules: triggeredRulesFromJSON(record.TriggeredRulesJSON), AnalystReports: analystReportsFromJSON(record.AnalystReportsJSON), RetrievalQuality: retrievalQualityFromContextSnapshot(record.ContextSnapshotJSON), MarketContext: marketContextFromContextSnapshot(record.ContextSnapshotJSON, record.Symbol), ExpectedReturnScenarios: expectedReturnFromJSON(record.ExpectedReturnScenariosJSON), ArbitrationChain: arbitrationChainFromJSON(record.ArbitrationChainJSON, record.TriggeredRulesJSON), FinalVerdict: dto.FinalVerdict{Status: record.FinalVerdictStatus, DisplayText: record.FinalVerdictText, ProhibitedActions: splitJSONStrings(record.ProhibitedActionsJSON), OptionalActions: splitJSONStrings(record.OptionalActionsJSON)}, UserConfirmation: dto.UserConfirmation{ConfirmationStatus: record.ConfirmationStatus, AvailableActions: availableConfirmationActions(record.ConfirmationStatus)}}
}

func marketContextDTO(market model.MarketSnapshot, fallbackSymbol string) *dto.MarketContext {
	if market.Symbol == "" && fallbackSymbol == "" && market.TradeDate == "" && market.ClosePrice == 0 && market.PEPercentile == 0 && market.PBPercentile == 0 {
		return nil
	}
	symbol := firstNonEmptyString(market.Symbol, fallbackSymbol)
	return &dto.MarketContext{Symbol: symbol, TradeDate: market.TradeDate, CurrentPrice: market.ClosePrice, PEPercentile: market.PEPercentile, PBPercentile: market.PBPercentile}
}

func marketContextFromContextSnapshot(raw, fallbackSymbol string) *dto.MarketContext {
	if raw == "" {
		return nil
	}
	var stored struct {
		MarketSnapshot      model.MarketSnapshot `json:"MarketSnapshot"`
		MarketSnapshotLower model.MarketSnapshot `json:"market_snapshot"`
	}
	if err := json.Unmarshal([]byte(raw), &stored); err != nil {
		return nil
	}
	if got := marketContextDTO(stored.MarketSnapshot, fallbackSymbol); got != nil {
		return got
	}
	return marketContextDTO(stored.MarketSnapshotLower, fallbackSymbol)
}

func retrievalQualityFromContextSnapshot(raw string) *dto.RetrievalQualitySummary {
	if raw == "" {
		return nil
	}
	var stored struct {
		RetrievalQualitySummary workflow.RetrievalQualitySummary `json:"retrieval_quality_summary"`
		Legacy                  workflow.RetrievalQualitySummary `json:"RetrievalQualitySummary"`
	}
	if err := json.Unmarshal([]byte(raw), &stored); err != nil {
		return nil
	}
	if got := retrievalQualityDTO(stored.RetrievalQualitySummary); got != nil {
		return got
	}
	return retrievalQualityDTO(stored.Legacy)
}
