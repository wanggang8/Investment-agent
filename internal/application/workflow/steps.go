package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"

	appknowledge "investment-agent/internal/application/knowledge"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	domainrule "investment-agent/internal/domain/rule"
)

// workflowStep 是 Eino 节点内部复用的应用层节点函数。
type workflowStep func(context.Context, *WorkflowContext, WorkflowDependencies) NodeResult

// RunStateSnapshotNode 执行状态快照节点。
func RunStateSnapshotNode(ctx context.Context, wf *WorkflowContext, deps WorkflowDependencies) NodeResult {
	return stateSnapshotStep(ctx, wf, deps)
}

// RunCapabilityCheckNode 执行能力圈检查节点。
func RunCapabilityCheckNode(ctx context.Context, wf *WorkflowContext, deps WorkflowDependencies) NodeResult {
	return capabilityCheckStep(ctx, wf, deps)
}

// RunEvidenceRetrievalNode 执行证据读取节点。
func RunEvidenceRetrievalNode(ctx context.Context, wf *WorkflowContext, deps WorkflowDependencies) NodeResult {
	return evidenceRetrievalStep(ctx, wf, deps)
}

// RunValueAnalystNode 执行价值分析节点。
func RunValueAnalystNode(ctx context.Context, wf *WorkflowContext, deps WorkflowDependencies) NodeResult {
	return valueAnalystStep(ctx, wf, deps)
}

// RunTrendRiskOfficerNode 执行趋势风险节点。
func RunTrendRiskOfficerNode(ctx context.Context, wf *WorkflowContext, deps WorkflowDependencies) NodeResult {
	return trendRiskOfficerStep(ctx, wf, deps)
}

// RunExpectedReturnNode 执行预期收益节点。
func RunExpectedReturnNode(ctx context.Context, wf *WorkflowContext, deps WorkflowDependencies) NodeResult {
	return expectedReturnStep(ctx, wf, deps)
}

// RunRuleArbitrationNode 执行规则裁决节点。
func RunRuleArbitrationNode(ctx context.Context, wf *WorkflowContext, deps WorkflowDependencies) NodeResult {
	return ruleArbitrationStep(ctx, wf, deps)
}

// RunDecisionRecordNode 执行决策记录节点。
func RunDecisionRecordNode(ctx context.Context, wf *WorkflowContext, deps WorkflowDependencies) NodeResult {
	return decisionRecordStep(ctx, wf, deps)
}

// stateSnapshotStep 校验账户、行情和规则版本快照。
func stateSnapshotStep(_ context.Context, wf *WorkflowContext, _ WorkflowDependencies) NodeResult {
	switch {
	case wf.PortfolioSnapshot.SnapshotID == "":
		return nodeResult("StateSnapshotNode", "load_state_snapshot", StatusFailed, ErrCodeDataRequired, "request", wf.RequestID, "", "")
	case wf.MarketSnapshot.MarketSnapshotID == "":
		return nodeResult("StateSnapshotNode", "load_state_snapshot", StatusFailed, ErrCodeDataStale, "request", wf.RequestID, "", "")
	case wf.RuleVersion == "":
		return nodeResult("StateSnapshotNode", "load_state_snapshot", StatusFailed, ErrCodeRuleVersionMissing, "request", wf.RequestID, "", "")
	default:
		return nodeResult("StateSnapshotNode", "load_state_snapshot", StatusSuccess, "", "request", wf.RequestID, "snapshot", wf.PortfolioSnapshot.SnapshotID)
	}
}

// capabilityCheckStep 识别主动咨询是否超出能力圈；未知配置不得伪造成能力圈内。
func capabilityCheckStep(_ context.Context, wf *WorkflowContext, _ WorkflowDependencies) NodeResult {
	if wf.CapabilityStatus == "" {
		wf.CapabilityStatus = CapabilityUnknown
	}
	return nodeResult("CapabilityCheckNode", "check_capability", StatusSuccess, "", "symbol", wf.Symbol, "capability", wf.CapabilityStatus)
}

// evidenceRetrievalStep 校验证据链是否可用于正式裁决。
func evidenceRetrievalStep(ctx context.Context, wf *WorkflowContext, deps WorkflowDependencies) NodeResult {
	wf.RetrievalInput = wf.Symbol
	if svc := deps.retrievalService(); svc != nil {
		result, err := svc.RetrieveEvidence(ctx, RetrievalRequest{Symbol: wf.Symbol, Query: wf.UserQuestion})
		if err != nil {
			wf.RetrievalDegradedReason = err.Error()
			return nodeResult("EvidenceRetrievalNode", "retrieve_evidence", StatusFailed, ErrCodeEvidenceNotFound, "symbol", wf.Symbol, "retrieval", wf.RetrievalDegradedReason)
		}
		wf.EvidenceSet = result.EvidenceSet
		wf.RetrievalOutputRef = result.OutputRef
		wf.RetrievalDegradedReason = result.DegradedReason
		wf.RetrievalQualitySummary = result.QualitySummary
		outputRef := retrievalAuditOutputRef(result.OutputRef, result.QualitySummary)
		if len(wf.EvidenceSet.Items) == 0 {
			return nodeResult("EvidenceRetrievalNode", "retrieve_evidence", StatusFailed, ErrCodeEvidenceNotFound, "symbol", wf.Symbol, "retrieval", result.DegradedReason)
		}
		if result.DegradedReason != "" {
			return nodeResult("EvidenceRetrievalNode", "retrieve_evidence", StatusDegraded, ErrCodeVectorIndexUnavailable, "symbol", wf.Symbol, "evidence_set", outputRef)
		}
		return nodeResult("EvidenceRetrievalNode", "retrieve_evidence", StatusSuccess, "", "symbol", wf.Symbol, "evidence_set", outputRef)
	}
	if len(wf.EvidenceSet.Items) == 0 {
		return nodeResult("EvidenceRetrievalNode", "retrieve_evidence", StatusFailed, ErrCodeEvidenceNotFound, "symbol", wf.Symbol, "", "")
	}
	wf.RetrievalOutputRef = "evidence_set"
	return nodeResult("EvidenceRetrievalNode", "retrieve_evidence", StatusSuccess, "", "symbol", wf.Symbol, "evidence_set", "evidence_set")
}

// valueAnalystStep 只生成价值分析材料，不写最终裁决。
func valueAnalystStep(ctx context.Context, wf *WorkflowContext, deps WorkflowDependencies) NodeResult {
	if wf.AnalystReports == nil {
		wf.AnalystReports = map[string]string{}
	}
	if wf.AnalystUnavailable {
		return nodeResult("ValueAnalystNode", "run_value_analyst", StatusDegraded, ErrCodeAnalystUnavailable, "symbol", wf.Symbol, "analysis_report", "value")
	}
	resp, err := deps.analystService().Analyze(ctx, AnalystRequest{AgentName: "value", Symbol: wf.Symbol, EvidenceSummary: evidenceSummary(wf), KnowledgeContextSummary: knowledgeContextSummary(wf), RuleBoundary: "DeepSeek 只生成分析材料，最终裁决由规则引擎负责"})
	if err != nil {
		return nodeResult("ValueAnalystNode", "run_value_analyst", StatusDegraded, ErrCodeAnalystUnavailable, "symbol", wf.Symbol, "analysis_report", analystOutputRef("value", err))
	}
	// LLM 输出只写 analyst_reports，不能改写 RuleVerdict。
	wf.AnalystReports["value"] = firstNonEmpty(resp.Reports["value"], "估值与基本面分析材料")
	recordAnalystMetadata(wf, "value", resp.Metadata)
	return nodeResult("ValueAnalystNode", "run_value_analyst", StatusSuccess, "", "symbol", wf.Symbol, "analysis_report", "value")
}

// trendRiskOfficerStep 只生成趋势和风险分析材料，不写最终裁决。
func trendRiskOfficerStep(ctx context.Context, wf *WorkflowContext, deps WorkflowDependencies) NodeResult {
	if wf.AnalystReports == nil {
		wf.AnalystReports = map[string]string{}
	}
	if wf.AnalystUnavailable {
		return nodeResult("TrendRiskOfficerNode", "run_trend_risk_officer", StatusDegraded, ErrCodeAnalystUnavailable, "symbol", wf.Symbol, "analysis_report", "trend_risk")
	}
	resp, err := deps.analystService().Analyze(ctx, AnalystRequest{AgentName: "trend_risk", Symbol: wf.Symbol, EvidenceSummary: evidenceSummary(wf), KnowledgeContextSummary: knowledgeContextSummary(wf), RuleBoundary: "DeepSeek 只生成分析材料，最终裁决由规则引擎负责"})
	if err != nil {
		return nodeResult("TrendRiskOfficerNode", "run_trend_risk_officer", StatusDegraded, ErrCodeAnalystUnavailable, "symbol", wf.Symbol, "analysis_report", analystOutputRef("trend_risk", err))
	}
	wf.AnalystReports["trend_risk"] = firstNonEmpty(resp.Reports["trend_risk"], "趋势与风险分析材料")
	recordAnalystMetadata(wf, "trend_risk", resp.Metadata)
	return nodeResult("TrendRiskOfficerNode", "run_trend_risk_officer", StatusSuccess, "", "symbol", wf.Symbol, "analysis_report", "trend_risk")
}

// expectedReturnStep 生成预期收益情景，并把 LLM 分析限制在材料层。
func expectedReturnStep(ctx context.Context, wf *WorkflowContext, deps WorkflowDependencies) NodeResult {
	out := BuildExpectedReturnWithContext(expectedReturnInputFromWorkflow(wf))
	wf.ExpectedReturnScenarios = out.Scenarios
	wf.ExpectedReturnPrecisionStatus = out.PrecisionStatus
	wf.ExpectedReturnReason = out.Reason
	wf.ExpectedReturnTargetName = out.TargetName
	wf.ExpectedReturnTargetCode = out.TargetCode
	wf.ExpectedReturnHoldingClass = out.HoldingClass
	wf.ExpectedReturnHorizonLabel = out.HorizonLabel
	wf.ExpectedReturnProbabilityBasis = out.ProbabilityBasis
	wf.ExpectedReturnSupportingDataSummary = out.SupportingDataSummary
	wf.ExpectedReturnMissingCategories = out.MissingCategories
	wf.ExpectedReturnSupplementData = out.SupplementData
	wf.ExpectedReturnAssumptionChecks = out.AssumptionChecks
	wf.ExpectedReturnHistoricalContexts = out.HistoricalContexts
	wf.ExpectedReturnHoldingClassCoverage = out.HoldingClassCoverage
	wf.ExpectedReturnSampleCount = out.SampleCount
	wf.ExpectedReturnSampleWindow = out.SampleWindow
	wf.ExpectedReturnScreeningCondition = out.ScreeningCondition
	wf.ExpectedReturnSellEvaluation = out.SellEvaluation
	wf.ExpectedReturnReassessmentTrigger = out.ReassessmentTrigger
	if wf.AnalystReports == nil {
		wf.AnalystReports = map[string]string{}
	}
	resp, err := deps.analystService().Analyze(ctx, AnalystRequest{AgentName: "expected_return", Symbol: wf.Symbol, EvidenceSummary: evidenceSummary(wf), KnowledgeContextSummary: knowledgeContextSummary(wf), RuleBoundary: "DeepSeek 只生成预期收益分析材料，最终裁决由规则引擎负责"})
	if err != nil {
		if isAnalystQualityFailure(err) {
			wf.AnalystReports["expected_return"] = expectedReturnLocalFallbackReport(out)
			recordAnalystMetadata(wf, "expected_return", map[string]string{
				"prompt_version":    "p37-analyst-v1",
				"model":             "deterministic-local",
				"parse_status":      "parsed",
				"quality_status":    "passed",
				"fallback_reason":   "llm_quality_failure",
				"input_summary":     wf.Symbol,
				"output_summary":    "local expected-return fallback",
				"decision_boundary": "analysis_material_only",
			})
			return nodeResult("ExpectedReturnNode", "estimate_expected_return", StatusSuccess, "", "symbol", wf.Symbol, "expected_return_scenarios", string(out.PrecisionStatus)+":fallback=deterministic_local:llm=quality_failure")
		}
		// 分析服务失败只让该节点降级；本地样本情景仍保留给规则流程参考。
		return nodeResult("ExpectedReturnNode", "estimate_expected_return", StatusDegraded, ErrCodeAnalystUnavailable, "symbol", wf.Symbol, "expected_return_scenarios", analystOutputRef(string(out.PrecisionStatus), err))
	}
	wf.AnalystReports["expected_return"] = firstNonEmpty(resp.Reports["expected_return"], "预期收益分析材料")
	recordAnalystMetadata(wf, "expected_return", resp.Metadata)
	return nodeResult("ExpectedReturnNode", "estimate_expected_return", StatusSuccess, "", "symbol", wf.Symbol, "expected_return_scenarios", string(out.PrecisionStatus))
}

// ruleArbitrationStep 调用领域规则生成最终裁决。
func ruleArbitrationStep(_ context.Context, wf *WorkflowContext, _ WorkflowDependencies) NodeResult {
	status := verificationStatus(*wf)
	wf.RuleVerdict = domainrule.Evaluate(wf.toDomainContext(), domainrule.EvaluationInput{
		CapabilityStatus:         wf.CapabilityStatus,
		HasEvidence:              len(wf.EvidenceSet.Items) > 0,
		Evidence:                 wf.EvidenceSet.Items,
		SourceVerificationStatus: status,
		BuyLogicBroken:           firstPosition(wf).BuyLogicBroken,
		SentimentState:           wf.MarketSnapshot.SentimentState,
		LiquidityState:           wf.MarketSnapshot.LiquidityState,
		PEPercentile:             wf.MarketSnapshot.PEPercentile,
		PBPercentile:             wf.MarketSnapshot.PBPercentile,
		CashRatio:                wf.PortfolioSnapshot.CashRatio,
		UnrealizedProfitRatio:    firstPosition(wf).UnrealizedProfitRatio,
		TakeProfitStarted:        firstPosition(wf).TakeProfitStarted,
		StageHighPrice:           firstPosition(wf).StageHighPrice,
		CurrentPrice:             firstPosition(wf).CurrentPrice,
		HandledProfit20:          firstPosition(wf).HandledProfit20,
		HandledProfit30:          firstPosition(wf).HandledProfit30,
	})
	return nodeResult("RuleArbitrationNode", "arbitrate_rule", StatusSuccess, "", "workflow_context", wf.RequestID, "rule_verdict", string(wf.RuleVerdict.Status))
}

// decisionRecordStep 写入决策记录和证据引用；有 Repository 时进入 SQLite 事实表。
func decisionRecordStep(_ context.Context, wf *WorkflowContext, _ WorkflowDependencies) NodeResult {
	if wf.DecisionID == "" {
		wf.DecisionID = workflowID("decision")
	}
	return nodeResult("DecisionRecordNode", "record_decision", StatusSuccess, "", "rule_verdict", string(wf.RuleVerdict.Status), "decision_record", wf.DecisionID)
}

func writeWorkflowAudit(ctx context.Context, writer AuditWriter, deps WorkflowDependencies, wf *WorkflowContext, result NodeResult) error {
	if deps.Transactor == nil || result.Audit.NodeName != "DecisionRecordNode" || deps.DecisionRepo == nil {
		return writer.Write(ctx, wf, result)
	}
	decision := buildDecisionRecord(*wf)
	refs := buildEvidenceRefs(*wf)
	return deps.Transactor.WithinTx(ctx, func(txCtx context.Context, repos repository.Repositories) error {
		if err := repos.DecisionRepo.SaveDecisionRecord(txCtx, decision, refs); err != nil {
			return err
		}
		return writeAuditEvent(txCtx, repos.AuditRepo, wf, result)
	})
}

func nodeResult(nodeName, nodeAction string, status NodeStatus, code, inputType, inputRef, outputType, outputRef string) NodeResult {
	return NodeResult{Status: status, ErrorCode: code, Audit: AuditFragment{Action: "generate_decision", NodeName: nodeName, NodeAction: nodeAction, Status: status, InputRefType: inputType, InputRef: inputRef, OutputRefType: outputType, OutputRef: outputRef, ErrorCode: code}}
}

func retrievalAuditOutputRef(base string, summary RetrievalQualitySummary) string {
	base = strings.TrimSpace(base)
	if base == "" {
		base = "retrieval"
	}
	base = summarizeAuditRef(base)
	if summary.TopK > 0 {
		base += ":topk=" + intString(summary.TopK)
	}
	if strings.TrimSpace(summary.FallbackSource) != "" {
		base += ":fallback=" + summarizeAuditRef(summary.FallbackSource)
	}
	if strings.TrimSpace(summary.IndexHealth) != "" {
		base += ":index=" + summarizeAuditRef(summary.IndexHealth)
	}
	if strings.TrimSpace(summary.SourceConsistencyStatus) != "" {
		base += ":consistency=" + summarizeAuditRef(summary.SourceConsistencyStatus)
	}
	if strings.TrimSpace(summary.DegradedReason) != "" {
		base += ":degraded=" + summarizeAuditRef(summary.DegradedReason)
	}
	return base
}

func intString(value int) string {
	data, _ := json.Marshal(value)
	return string(data)
}

type categorizedError interface {
	Category() string
}

type metadataError interface {
	Metadata() map[string]string
}

func analystOutputRef(base string, err error) string {
	base = strings.TrimSpace(base)
	if base == "" {
		base = "analyst"
	}
	var categorized categorizedError
	if err != nil && errors.As(err, &categorized) {
		category := strings.TrimSpace(categorized.Category())
		if category != "" {
			base += ":category=" + p52AnalystFailureCategory(category)
		}
	}
	var withMetadata metadataError
	if err != nil && errors.As(err, &withMetadata) {
		metadata := withMetadata.Metadata()
		if v := strings.TrimSpace(metadata["prompt_version"]); v != "" {
			base += ":prompt=" + v
		}
		if v := strings.TrimSpace(metadata["model"]); v != "" {
			base += ":model=" + v
		}
		if v := strings.TrimSpace(metadata["parse_status"]); v != "" {
			base += ":parse=" + v
		}
		if v := strings.TrimSpace(metadata["quality_status"]); v != "" {
			base += ":quality=" + v
		}
		if v := strings.TrimSpace(metadata["output_summary"]); v != "" {
			base += ":output=" + summarizeAuditRef(v)
		}
	}
	return base
}

func isAnalystQualityFailure(err error) bool {
	var categorized categorizedError
	if err == nil || !errors.As(err, &categorized) {
		return false
	}
	return p52AnalystFailureCategory(categorized.Category()) == "quality_failure"
}

func expectedReturnLocalFallbackReport(out ExpectedReturnOutput) string {
	parts := []string{
		"本地预期收益情景：LLM 预期收益材料未通过安全质量校验，已丢弃该输出，以下仅使用本地确定性情景。",
		"精度状态：" + string(out.PrecisionStatus) + "；样本数：" + strconv.Itoa(out.SampleCount) + "。",
	}
	if out.Reason != "" {
		parts = append(parts, "原因："+out.Reason+"。")
	}
	if len(out.Scenarios) == 0 {
		parts = append(parts, "当前样本不足，不能生成概率或收益区间。")
	} else {
		summaries := make([]string, 0, len(out.Scenarios))
		for _, scenario := range out.Scenarios {
			summaries = append(summaries, scenario.Name+"="+scenario.ReturnRange+"("+scenario.Confidence+")")
		}
		parts = append(parts, "情景区间："+strings.Join(summaries, "；")+"。")
	}
	if out.SellEvaluation.NonTradingDisclaimer != "" {
		parts = append(parts, out.SellEvaluation.NonTradingDisclaimer)
	} else {
		parts = append(parts, "仅提示人工复核，不构成交易指令，也不会自动交易。")
	}
	return strings.Join(parts, "\n")
}

func summarizeAuditRef(value string) string {
	value = strings.TrimSpace(strings.Join(strings.Fields(value), " "))
	value = regexp.MustCompile(`sk-[A-Za-z0-9_-]+`).ReplaceAllString(value, "sk_redacted")
	value = strings.NewReplacer(":", "_", "：", "_", "=", "_", "/", "_", "\\", "_", "\n", " ", "\r", " ", "\t", " ").Replace(value)
	runes := []rune(value)
	if len(runes) <= 40 {
		return value
	}
	return string(runes[:40])
}

func buildDecisionRecord(wf WorkflowContext) repository.DecisionRecord {
	now := workflowNowRFC3339()
	return repository.DecisionRecord{DecisionID: wf.DecisionID, RequestID: wf.RequestID, WorkflowType: wf.WorkflowType, Symbol: wf.Symbol, Question: wf.UserQuestion, WorkflowStatus: string(workflowStatus(wf)), RecordType: recordType(wf.RuleVerdict.Status), DashboardState: dashboardState(wf.RuleVerdict.Status), CapabilityStatus: wf.CapabilityStatus, CapabilityReason: wf.CapabilityReason, SourceVerificationStatus: string(verificationStatus(wf)), TriggeredRulesJSON: jsonString(wf.RuleVerdict.TriggeredRules), ErrorsJSON: jsonString(wf.Errors), FinalVerdictStatus: string(wf.RuleVerdict.Status), FinalVerdictText: wf.RuleVerdict.Text, ProhibitedActionsJSON: jsonString(wf.RuleVerdict.ProhibitedActions), OptionalActionsJSON: jsonString(wf.RuleVerdict.OptionalActions), ConfirmationStatus: confirmationStatusForVerdict(wf.RuleVerdict.Status), PortfolioSnapshotID: wf.PortfolioSnapshot.SnapshotID, MarketSnapshotID: wf.MarketSnapshot.MarketSnapshotID, RuleVersion: wf.RuleVersion, AnalystReportsJSON: jsonString(analystReportsForStorage(wf)), ExpectedReturnScenariosJSON: jsonString(map[string]any{"target_name": wf.ExpectedReturnTargetName, "target_code": wf.ExpectedReturnTargetCode, "holding_class": wf.ExpectedReturnHoldingClass, "horizon_label": wf.ExpectedReturnHorizonLabel, "precision_status": wf.ExpectedReturnPrecisionStatus, "reason": wf.ExpectedReturnReason, "sample_count": wf.ExpectedReturnSampleCount, "sample_window": wf.ExpectedReturnSampleWindow, "screening_condition": wf.ExpectedReturnScreeningCondition, "probability_basis": wf.ExpectedReturnProbabilityBasis, "supporting_data_summary": wf.ExpectedReturnSupportingDataSummary, "missing_categories": wf.ExpectedReturnMissingCategories, "supplement_data": wf.ExpectedReturnSupplementData, "assumption_checks": wf.ExpectedReturnAssumptionChecks, "historical_contexts": wf.ExpectedReturnHistoricalContexts, "holding_class_coverage": wf.ExpectedReturnHoldingClassCoverage, "source_health": p34SourceHealthForDecision(wf.MarketSnapshot), "scenarios": wf.ExpectedReturnScenarios, "sell_evaluation": wf.ExpectedReturnSellEvaluation, "reassessment_trigger": wf.ExpectedReturnReassessmentTrigger}), ContextSnapshotJSON: jsonString(wf.contextSnapshot()), CreatedAt: now}
}

func recordAnalystMetadata(wf *WorkflowContext, agent string, metadata map[string]string) {
	if len(metadata) == 0 {
		return
	}
	if wf.AnalystReportMetadata == nil {
		wf.AnalystReportMetadata = map[string]map[string]string{}
	}
	copied := make(map[string]string, len(metadata))
	for key, value := range metadata {
		copied[key] = value
	}
	wf.AnalystReportMetadata[agent] = copied
}

func analystReportsForStorage(wf WorkflowContext) any {
	if len(wf.AnalystReportMetadata) == 0 {
		return wf.AnalystReports
	}
	out := make([]map[string]any, 0, len(wf.AnalystReports))
	for agent, conclusion := range wf.AnalystReports {
		item := map[string]any{"agent_name": agent, "conclusion": conclusion}
		for key, value := range wf.AnalystReportMetadata[agent] {
			item[key] = value
		}
		out = append(out, item)
	}
	return out
}

func buildEvidenceRefs(wf WorkflowContext) []repository.EvidenceRef {
	now := workflowNowRFC3339()
	refs := make([]repository.EvidenceRef, 0, len(wf.EvidenceSet.Items))
	for _, e := range wf.EvidenceSet.Items {
		summaryID := firstNonEmpty(e.SummaryID, e.EvidenceID)
		refs = append(refs, repository.EvidenceRef{EvidenceRefID: workflowID("eref"), EvidenceID: e.EvidenceID, DecisionID: wf.DecisionID, SummaryID: summaryID, SourceName: e.SourceName, SourceLevel: string(e.SourceLevel), EvidenceRole: string(e.Role), PublishedAt: e.PublishedAt, CapturedAt: e.CapturedAt, OriginalURL: e.OriginalURL, Summary: e.Summary, ContentHash: e.ContentHash, TimeWeight: e.TimeWeight, RelevanceScore: e.RelevanceScore, IndependentSourceCount: e.IndependentSourceCount, HighGradeIndependentSourceCount: e.HighGradeIndependentSourceCount, CreatedAt: now})
	}
	return refs
}

func jsonString(v any) string { b, _ := json.Marshal(v); return string(b) }

func p34SourceHealthForDecision(market model.MarketSnapshot) []map[string]any {
	var metrics map[string]any
	if err := json.Unmarshal([]byte(market.MarketMetricsJSON), &metrics); err != nil {
		return nil
	}
	metadata, _ := metrics["metadata"].(map[string]any)
	health, _ := metadata["p34_source_health"].(map[string]any)
	categoriesRaw, _ := metadata["p34_data_categories"].([]any)
	categories := make([]string, 0, len(categoriesRaw))
	for _, raw := range categoriesRaw {
		if category, ok := raw.(string); ok && category != "" {
			categories = append(categories, category)
		}
	}
	if len(categories) == 0 {
		for category := range health {
			categories = append(categories, category)
		}
	}
	items := make([]map[string]any, 0, len(categories))
	for _, category := range categories {
		item := map[string]any{"source_name": metrics["source_name"], "source_level": metrics["source_level"], "source_type": metrics["source_type"], "data_category": category, "data_date": market.TradeDate, "affected_symbols": []string{market.Symbol}}
		switch raw := health[category].(type) {
		case map[string]any:
			for key, value := range raw {
				item[key] = value
			}
		case string:
			item["freshness"] = raw
		}
		items = append(items, item)
	}
	return items
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
func evidenceSummary(wf *WorkflowContext) string {
	if len(wf.EvidenceSet.Items) == 0 {
		return ""
	}
	return wf.EvidenceSet.Items[0].EvidenceID
}

func knowledgeContextSummary(wf *WorkflowContext) string {
	entries := appknowledge.BuiltInRegistry().EntriesForSymbol(wf.Symbol)
	deps := []appknowledge.DataDependency{
		{Category: "symbol_profile", Status: knowledgeStatus(profileKnown(wf.Symbol))},
		{Category: "active_rule", Status: knowledgeStatus(strings.TrimSpace(wf.RuleVersion) != "")},
	}
	for _, category := range []string{"fund_profile", "tracked_index"} {
		deps = append(deps, appknowledge.DataDependency{Category: category, Status: readinessStatusFromP34Health(wf.MarketSnapshot, category)})
	}
	valuationStatus := appknowledge.ReadinessDegraded
	if wf.MarketSnapshot.PEPercentile > 0 || wf.MarketSnapshot.PBPercentile > 0 {
		valuationStatus = appknowledge.ReadinessReady
	}
	marketPriceStatus := appknowledge.ReadinessDegraded
	if wf.MarketSnapshot.ClosePrice > 0 {
		marketPriceStatus = appknowledge.ReadinessReady
	}
	liquidityStatus := appknowledge.ReadinessDegraded
	if strings.TrimSpace(string(wf.MarketSnapshot.LiquidityState)) != "" {
		liquidityStatus = appknowledge.ReadinessReady
	}
	sentimentStatus := appknowledge.ReadinessDegraded
	if strings.TrimSpace(string(wf.MarketSnapshot.SentimentState)) != "" {
		sentimentStatus = appknowledge.ReadinessReady
	}
	evidenceStatus := string(verificationStatus(*wf))
	if evidenceStatus == "" {
		evidenceStatus = "failed"
	}
	ragStatus := appknowledge.ReadinessDegraded
	if len(wf.EvidenceSet.Items) > 0 && strings.TrimSpace(wf.RetrievalDegradedReason) == "" {
		ragStatus = appknowledge.ReadinessReady
	}
	deps = append(deps,
		appknowledge.DataDependency{Category: "market_price", Status: marketPriceStatus},
		appknowledge.DataDependency{Category: "valuation_percentiles", Status: valuationStatus},
		appknowledge.DataDependency{Category: "liquidity", Status: liquidityStatus},
		appknowledge.DataDependency{Category: "sentiment_proxy", Status: sentimentStatus},
		appknowledge.DataDependency{Category: "formal_evidence", Status: evidenceStatus},
		appknowledge.DataDependency{Category: "rag_index", Status: ragStatus},
		appknowledge.DataDependency{Category: "llm_context", Status: appknowledge.ReadinessReady},
	)
	summary := appknowledge.BuildLLMContextSummary(entries, deps)
	if facts := structuredFinancialFactsSummary(wf.MarketSnapshot); facts != "" {
		summary += "; " + facts
	}
	return summary
}

func structuredFinancialFactsSummary(market model.MarketSnapshot) string {
	facts := []string{}
	appendFact := func(name string, value float64) {
		if value == 0 {
			return
		}
		facts = append(facts, name+"="+strconv.FormatFloat(value, 'f', -1, 64))
	}
	appendFact("close_price", market.ClosePrice)
	appendFact("pe_percentile", market.PEPercentile)
	appendFact("pb_percentile", market.PBPercentile)
	appendFact("turnover_rate", market.TurnoverRate)
	appendFact("margin_balance", market.MarginBalance)
	appendFact("margin_balance_change", market.MarginBalanceChange)
	appendFact("volume_percentile", market.VolumePercentile)
	appendFact("volatility_percentile", market.VolatilityPercentile)
	if len(facts) == 0 {
		return ""
	}
	return "structured_financial_facts=" + strings.Join(facts, ",") + "; precedence=structured_facts_override_text_claims"
}

func profileKnown(symbol string) bool {
	_, ok := appknowledge.LookupSymbolProfile(symbol)
	return ok
}

func knowledgeStatus(ready bool) string {
	if ready {
		return appknowledge.ReadinessReady
	}
	return appknowledge.ReadinessDegraded
}

func readinessStatusFromP34Health(market model.MarketSnapshot, category string) string {
	for _, item := range p34SourceHealthForDecision(market) {
		if itemCategory, _ := item["data_category"].(string); itemCategory != category {
			continue
		}
		freshness, _ := item["freshness"].(string)
		if freshness == "fresh" || freshness == "stubbed" {
			return appknowledge.ReadinessReady
		}
		return appknowledge.ReadinessDegraded
	}
	return appknowledge.ReadinessDegraded
}
func expectedReturnInputFromWorkflow(wf *WorkflowContext) ExpectedReturnInput {
	position := positionForSymbol(wf)
	currentPrice := position.CurrentPrice
	if currentPrice == 0 {
		currentPrice = wf.MarketSnapshot.ClosePrice
	}
	p34Summary, p34Missing := p34ExpectedReturnContext(wf.MarketSnapshot)
	profileName := ""
	if profile, ok := appknowledge.LookupSymbolProfile(wf.Symbol); ok {
		profileName = profile.Name
	}
	marketState, fundamentalState, pessimisticPathMonths := expectedReturnDynamicMonitoringFromMarket(wf.MarketSnapshot)
	return ExpectedReturnInput{SampleCount: wf.ExpectedReturnSampleCount, TargetName: firstNonEmpty(position.Name, profileName), TargetCode: wf.Symbol, HoldingClass: expectedReturnHoldingClass(position, wf.Symbol), HorizonLabel: "未来 12 个月", CurrentPrice: currentPrice, BasePrice: position.CostPrice, PreviousBaseMidpoint: wf.ExpectedReturnPreviousBaseMidpoint, TargetReturnRate: wf.ExpectedReturnTargetReturnRate, SentimentState: string(wf.MarketSnapshot.SentimentState), MarketState: marketState, FundamentalState: fundamentalState, SampleWindow: "当前本地持仓、最新市场快照与可用公开净值历史", ScreeningCondition: "基于当前标的持仓成本、最新市场快照和已保存公开市场元数据；样本不足时仅作定性参考", SupportingDataSummary: p34Summary, MissingCategories: p34Missing, HistoricalSamples: expectedReturnHistoricalSamplesFromMarket(wf.MarketSnapshot), HistoricalContexts: expectedReturnHistoricalContextsFromMarket(wf.MarketSnapshot), AssumptionChecks: expectedReturnAssumptionChecksFromMarket(wf.MarketSnapshot), PessimisticPathMonths: pessimisticPathMonths, HoldingClassCoverage: expectedReturnHoldingCoverage(wf.PositionSnapshots)}
}

func expectedReturnHoldingClass(position model.Position, symbol string) string {
	if strings.TrimSpace(position.AssetTag) == "satellite" {
		return "sector_growth_fund"
	}
	if symbol == "600000" {
		return "equity_constituent_financial"
	}
	if _, ok := appknowledge.LookupSymbolProfile(symbol); ok {
		return "broad_index_etf"
	}
	if strings.TrimSpace(position.AssetTag) == "core" {
		return "broad_index_etf"
	}
	return "unknown"
}

func expectedReturnHistoricalSamplesFromMarket(market model.MarketSnapshot) []ExpectedReturnHistoricalSample {
	var metrics map[string]any
	if err := json.Unmarshal([]byte(market.MarketMetricsJSON), &metrics); err != nil {
		return nil
	}
	metadata, _ := metrics["metadata"].(map[string]any)
	raw := metadata["expected_return_historical_samples"]
	if raw == nil {
		raw = metrics["expected_return_historical_samples"]
	}
	data, err := json.Marshal(raw)
	if err != nil || string(data) == "null" {
		return nil
	}
	var samples []ExpectedReturnHistoricalSample
	if err := json.Unmarshal(data, &samples); err != nil {
		return nil
	}
	return samples
}

func expectedReturnAssumptionChecksFromMarket(market model.MarketSnapshot) []ExpectedReturnAssumptionCheck {
	var metrics map[string]any
	if err := json.Unmarshal([]byte(market.MarketMetricsJSON), &metrics); err != nil {
		return nil
	}
	metadata, _ := metrics["metadata"].(map[string]any)
	raw := metadata["expected_return_assumption_checks"]
	if raw == nil {
		raw = metrics["expected_return_assumption_checks"]
	}
	data, err := json.Marshal(raw)
	if err != nil || string(data) == "null" {
		return nil
	}
	var checks []ExpectedReturnAssumptionCheck
	if err := json.Unmarshal(data, &checks); err != nil {
		return nil
	}
	return checks
}

func expectedReturnHistoricalContextsFromMarket(market model.MarketSnapshot) []ExpectedReturnHistoricalContext {
	var metrics map[string]any
	if err := json.Unmarshal([]byte(market.MarketMetricsJSON), &metrics); err != nil {
		return nil
	}
	metadata, _ := metrics["metadata"].(map[string]any)
	raw := metadata["expected_return_historical_contexts"]
	if raw == nil {
		raw = metrics["expected_return_historical_contexts"]
	}
	data, err := json.Marshal(raw)
	if err != nil || string(data) == "null" {
		return nil
	}
	var contexts []ExpectedReturnHistoricalContext
	if err := json.Unmarshal(data, &contexts); err != nil {
		return nil
	}
	return contexts
}

func expectedReturnDynamicMonitoringFromMarket(market model.MarketSnapshot) (string, string, int) {
	var metrics map[string]any
	if err := json.Unmarshal([]byte(market.MarketMetricsJSON), &metrics); err != nil {
		return "", "", 0
	}
	metadata, _ := metrics["metadata"].(map[string]any)
	if metadata == nil {
		metadata = map[string]any{}
	}
	marketState := firstNonEmpty(stringFromAny(metadata["expected_return_market_state"]), stringFromAny(metrics["expected_return_market_state"]), stringFromAny(metadata["market_state"]), stringFromAny(metrics["market_state"]))
	fundamentalState := firstNonEmpty(stringFromAny(metadata["expected_return_fundamental_state"]), stringFromAny(metrics["expected_return_fundamental_state"]), stringFromAny(metadata["fundamental_state"]), stringFromAny(metrics["fundamental_state"]))
	pessimisticPathMonths := intFromAny(metadata["expected_return_pessimistic_path_months"])
	if pessimisticPathMonths == 0 {
		pessimisticPathMonths = intFromAny(metrics["expected_return_pessimistic_path_months"])
	}
	return marketState, fundamentalState, pessimisticPathMonths
}

func intFromAny(value any) int {
	switch typed := value.(type) {
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	case json.Number:
		got, _ := typed.Int64()
		return int(got)
	case string:
		got, err := strconv.Atoi(strings.TrimSpace(typed))
		if err == nil {
			return got
		}
	}
	return 0
}

func expectedReturnHoldingCoverage(positions []model.Position) []ExpectedReturnHoldingClassCoverage {
	seen := map[string]bool{}
	out := []ExpectedReturnHoldingClassCoverage{}
	for _, position := range positions {
		class := expectedReturnHoldingClass(position, position.Symbol)
		if class == "unknown" || seen[class] {
			continue
		}
		seen[class] = true
		out = append(out, ExpectedReturnHoldingClassCoverage{HoldingClass: class, Symbol: position.Symbol, Status: "covered"})
	}
	return out
}
func recordType(status model.FinalVerdictStatus) string {
	switch status {
	case model.VerdictRejected:
		return "rejection_record"
	case model.VerdictInsufficientData, model.VerdictFrozenWatch:
		return "non_trade_record"
	default:
		return "formal_trade_advice"
	}
}
func confirmationStatusForVerdict(status model.FinalVerdictStatus) string {
	switch status {
	case model.VerdictBuyAllowed, model.VerdictHold, model.VerdictReduce, model.VerdictSellOnly:
		return string(model.ConfirmationPending)
	default:
		return string(model.ConfirmationNotRequired)
	}
}
func verificationStatus(wf WorkflowContext) model.VerificationStatus {
	if wf.SourceVerificationStatus.Valid() {
		return wf.SourceVerificationStatus
	}
	return wf.EvidenceSet.VerificationStatus
}
func dashboardState(status model.FinalVerdictStatus) string {
	switch status {
	case model.VerdictInsufficientData:
		return string(model.DashboardInsufficientData)
	case model.VerdictFrozenWatch:
		return string(model.DashboardFrozenWatch)
	case model.VerdictSellOnly, model.VerdictRejected:
		return string(model.DashboardHighRisk)
	default:
		return string(model.DashboardNormal)
	}
}
func firstPosition(wf *WorkflowContext) model.Position {
	if len(wf.PositionSnapshots) == 0 {
		return model.Position{}
	}
	return wf.PositionSnapshots[0]
}
func positionForSymbol(wf *WorkflowContext) model.Position {
	for _, position := range wf.PositionSnapshots {
		if position.Symbol == wf.Symbol {
			return position
		}
	}
	return model.Position{}
}
func workflowStatus(wf WorkflowContext) model.WorkflowStatus {
	status := model.WorkflowCompleted
	for _, code := range wf.Errors {
		switch code {
		case ErrCodeDataRequired, ErrCodeDataStale, ErrCodeRuleVersionMissing, ErrCodeEvidenceNotFound, ErrCodeSourceVerificationFailed, ErrCodeDecisionRecordFailed:
			return model.WorkflowFailed
		case ErrCodeVectorIndexUnavailable, ErrCodeAnalystUnavailable:
			status = model.WorkflowDegraded
		}
	}
	return status
}

func (wf WorkflowContext) toDomainContext() model.WorkflowContext {
	return model.WorkflowContext{RequestID: wf.RequestID, WorkflowType: wf.WorkflowType, UserQuestion: wf.UserQuestion, Symbol: wf.Symbol, PortfolioSnapshot: wf.PortfolioSnapshot, PositionSnapshots: wf.PositionSnapshots, MarketSnapshot: wf.MarketSnapshot, RuleVersion: wf.RuleVersion, CapabilityStatus: wf.CapabilityStatus, CapabilityReason: wf.CapabilityReason, SourceVerificationStatus: wf.SourceVerificationStatus, MediaHeatSummary: wf.MediaHeatSummary, UserEmotionTags: wf.UserEmotionTags, EvidenceSet: wf.EvidenceSet, AuditEvents: wf.AuditEvents, Errors: wf.Errors}
}

func (wf WorkflowContext) contextSnapshot() any {
	return struct {
		model.WorkflowContext
		RetrievalQualitySummary RetrievalQualitySummary `json:"retrieval_quality_summary,omitempty"`
	}{WorkflowContext: wf.toDomainContext(), RetrievalQualitySummary: wf.RetrievalQualitySummary}
}
