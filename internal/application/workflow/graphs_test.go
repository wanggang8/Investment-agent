package workflow

import (
	"context"
	"strings"
	"testing"

	appknowledge "investment-agent/internal/application/knowledge"
	"investment-agent/internal/domain/analyst"
	"investment-agent/internal/domain/model"
)

func TestDailyDisciplineGraphScenarios(t *testing.T) {
	cases := []struct {
		name      string
		input     WorkflowContext
		want      model.FinalVerdictStatus
		wantError string
	}{
		{name: "normal", input: sampleWorkflowContext(20), want: model.VerdictHold},
		{name: "insufficient evidence", input: func() WorkflowContext { c := sampleWorkflowContext(20); c.EvidenceSet.Items = nil; return c }(), want: model.VerdictInsufficientData},
		{name: "analyst unavailable", input: func() WorkflowContext { c := sampleWorkflowContext(20); c.AnalystUnavailable = true; return c }(), want: model.VerdictHold, wantError: ErrCodeAnalystUnavailable},
		{name: "sample too small", input: sampleWorkflowContext(4), want: model.VerdictHold},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := NewDailyDisciplineGraph(&MemoryAuditWriter{}).Run(context.Background(), tc.input)
			if err != nil {
				t.Fatalf("run daily graph: %v", err)
			}
			if out.RuleVerdict.Status != tc.want {
				t.Fatalf("status=%s want=%s", out.RuleVerdict.Status, tc.want)
			}
			if tc.wantError != "" && !hasString(out.Errors, tc.wantError) {
				t.Fatalf("expected error %s in %+v", tc.wantError, out.Errors)
			}
			if len(out.AuditEvents) == 0 {
				t.Fatal("expected audit events")
			}
		})
	}
}

func TestAnalystMaterialsCannotOverrideRuleVerdict(t *testing.T) {
	ctx := sampleWorkflowContext(20)
	ctx.MarketSnapshot.PEPercentile = 95
	graph, err := BuildDailyEinoGraph(context.Background(), &MemoryAuditWriter{}, WorkflowDependencies{AnalystService: analystReportStub{reports: map[string]string{
		"value":           "强烈买入",
		"trend_risk":      "强烈买入",
		"expected_return": "强烈买入",
	}}})
	if err != nil {
		t.Fatalf("BuildDailyEinoGraph: %v", err)
	}

	out, err := graph.Invoke(context.Background(), ctx)
	if err != nil {
		t.Fatalf("run daily graph: %v", err)
	}
	if out.RuleVerdict.Status == model.VerdictBuyAllowed {
		t.Fatalf("analyst reports must not override rule verdict: %+v", out.RuleVerdict)
	}
	if out.AnalystReports["value"] != "强烈买入" {
		t.Fatalf("expected analyst material preserved: %+v", out.AnalystReports)
	}
}

func TestAnalystRequestsIncludeKnowledgeReadinessContext(t *testing.T) {
	ctx := sampleWorkflowContext(20)
	recorder := &recordingAnalystService{}
	graph, err := BuildDailyEinoGraph(context.Background(), &MemoryAuditWriter{}, WorkflowDependencies{AnalystService: recorder})
	if err != nil {
		t.Fatalf("BuildDailyEinoGraph: %v", err)
	}

	out, err := graph.Invoke(context.Background(), ctx)
	if err != nil {
		t.Fatalf("run daily graph: %v", err)
	}
	if out.RuleVerdict.Status == model.VerdictBuyAllowed {
		t.Fatalf("LLM context must not override rule verdict: %+v", out.RuleVerdict)
	}
	if len(recorder.requests) < 3 {
		t.Fatalf("expected analyst requests to be recorded, got %d", len(recorder.requests))
	}
	knowledgeIDs := []string{}
	for _, entry := range appknowledge.BuiltInRegistry().Entries() {
		if entry.LLMContextAllowed && entry.Category != "symbol_profile" {
			knowledgeIDs = append(knowledgeIDs, entry.KnowledgeID)
		}
	}
	for _, req := range recorder.requests {
		for _, knowledgeID := range knowledgeIDs {
			if !strings.Contains(req.KnowledgeContextSummary, knowledgeID) {
				t.Fatalf("expected %s in knowledge readiness context on request %+v", knowledgeID, req)
			}
		}
		for _, readinessKey := range []string{"symbol_profile=", "fund_profile=", "tracked_index=", "active_rule=", "valuation_percentiles=", "formal_evidence=", "market_price=", "liquidity=", "sentiment_proxy=", "rag_index="} {
			if !strings.Contains(req.KnowledgeContextSummary, readinessKey) {
				t.Fatalf("expected %s in knowledge readiness context on request %+v", readinessKey, req)
			}
		}
		if !strings.Contains(req.KnowledgeContextSummary, "背景知识不能满足正式证据") {
			t.Fatalf("expected knowledge readiness context on request %+v", req)
		}
	}
}

func TestAnalystRequestsScopeSymbolProfileKnowledgeToWorkflowSymbol(t *testing.T) {
	ctx := sampleWorkflowContext(20)
	ctx.Symbol = "159915"
	ctx.MarketSnapshot.Symbol = "159915"
	recorder := &recordingAnalystService{}
	graph, err := BuildDailyEinoGraph(context.Background(), &MemoryAuditWriter{}, WorkflowDependencies{AnalystService: recorder})
	if err != nil {
		t.Fatalf("BuildDailyEinoGraph: %v", err)
	}

	if _, err := graph.Invoke(context.Background(), ctx); err != nil {
		t.Fatalf("run daily graph: %v", err)
	}
	if len(recorder.requests) == 0 {
		t.Fatalf("expected analyst requests")
	}
	for _, req := range recorder.requests {
		if !strings.Contains(req.KnowledgeContextSummary, "symbol_profile.159915") {
			t.Fatalf("expected 159915 profile in scoped knowledge context, got %q", req.KnowledgeContextSummary)
		}
		if strings.Contains(req.KnowledgeContextSummary, "symbol_profile.510300") {
			t.Fatalf("must not include unrelated 510300 profile in 159915 context: %q", req.KnowledgeContextSummary)
		}
	}
}

func TestAnalystRequestsPreferStructuredFinancialFacts(t *testing.T) {
	ctx := sampleWorkflowContext(20)
	ctx.MarketSnapshot.ClosePrice = 4.23
	ctx.MarketSnapshot.PEPercentile = 31
	ctx.MarketSnapshot.PBPercentile = 27
	ctx.MarketSnapshot.MarginBalance = 120000000
	ctx.MarketSnapshot.MarginBalanceChange = -0.08
	ctx.EvidenceSet.Items[0].Summary = "媒体评论称 PE 分位 99 且两融急升"
	recorder := &recordingAnalystService{}
	graph, err := BuildDailyEinoGraph(context.Background(), &MemoryAuditWriter{}, WorkflowDependencies{AnalystService: recorder})
	if err != nil {
		t.Fatalf("BuildDailyEinoGraph: %v", err)
	}

	if _, err := graph.Invoke(context.Background(), ctx); err != nil {
		t.Fatalf("run daily graph: %v", err)
	}
	if len(recorder.requests) == 0 {
		t.Fatalf("expected analyst requests")
	}
	for _, req := range recorder.requests {
		for _, want := range []string{"structured_financial_facts", "structured_facts_override_text_claims", "close_price=4.23", "pe_percentile=31", "pb_percentile=27", "margin_balance=120000000", "margin_balance_change=-0.08"} {
			if !strings.Contains(req.KnowledgeContextSummary, want) {
				t.Fatalf("expected %s in analyst knowledge context, got %q", want, req.KnowledgeContextSummary)
			}
		}
	}
}

func TestExpectedReturnMaterialIsExplanatoryOnly(t *testing.T) {
	ctx := sampleWorkflowContext(3)
	graph, err := BuildDailyEinoGraph(context.Background(), &MemoryAuditWriter{}, WorkflowDependencies{AnalystService: analystReportStub{reports: map[string]string{"expected_return": "收益翻倍"}}})
	if err != nil {
		t.Fatalf("BuildDailyEinoGraph: %v", err)
	}

	out, err := graph.Invoke(context.Background(), ctx)
	if err != nil {
		t.Fatalf("run daily graph: %v", err)
	}
	if out.ExpectedReturnPrecisionStatus != model.PrecisionUnavailable || out.RuleVerdict.Status != model.VerdictHold {
		t.Fatalf("expected return material must remain explanatory: precision=%s verdict=%+v", out.ExpectedReturnPrecisionStatus, out.RuleVerdict)
	}
	if out.AnalystReports["expected_return"] != "收益翻倍" {
		t.Fatalf("expected return analyst material not preserved: %+v", out.AnalystReports)
	}
}

func TestP88RuleArbitrationUsesMatchingSymbolForSourceVerifiedTransitions(t *testing.T) {
	wf := WorkflowContext{
		RequestID:         "req_p88_transition",
		Symbol:            "159915",
		PortfolioSnapshot: model.PortfolioSnapshot{SnapshotID: "snap_p88", CashRatio: 0.1, TotalAssets: 100000},
		MarketSnapshot:    model.MarketSnapshot{MarketSnapshotID: "market_p88", Symbol: "159915", PEPercentile: 20, PBPercentile: 20},
		PositionSnapshots: []model.Position{
			{Symbol: "510300", BuyLogicBroken: false, CostPrice: 3, CurrentPrice: 3},
			{Symbol: "159915", BuyLogicBroken: true, CostPrice: 2, CurrentPrice: 1.8},
		},
		EvidenceSet: model.EvidenceSet{VerificationStatus: model.VerificationSatisfied, Items: []model.Evidence{{EvidenceID: "ev_p88_break", SourceLevel: model.SourceLevelA, Role: model.EvidenceFormal, EventType: model.EventBuyLogicBreak, IndependentSourceCount: 2, HighGradeIndependentSourceCount: 2}}},
	}

	result := RunRuleArbitrationNode(context.Background(), &wf, WorkflowDependencies{})

	if result.Status != StatusSuccess {
		t.Fatalf("expected rule arbitration success, got %+v", result)
	}
	if wf.RuleVerdict.Status != model.VerdictSellOnly {
		t.Fatalf("expected matching symbol source-verified transition to sell_only, got %+v", wf.RuleVerdict)
	}
	if !strings.Contains(wf.RuleVerdict.Text, "A/S 独立信源=2") {
		t.Fatalf("expected source-count provenance, got %+v", wf.RuleVerdict)
	}
}

type analystReportStub struct{ reports map[string]string }

func (s analystReportStub) Analyze(context.Context, analyst.Request) (analyst.Response, error) {
	return analyst.Response{Reports: s.reports}, nil
}

type recordingAnalystService struct {
	requests []analyst.Request
}

func (s *recordingAnalystService) Analyze(_ context.Context, req analyst.Request) (analyst.Response, error) {
	s.requests = append(s.requests, req)
	return analyst.Response{Reports: map[string]string{req.AgentName: req.AgentName + " material"}}, nil
}

func TestConsultationGraphOutOfCapability(t *testing.T) {
	ctx := sampleWorkflowContext(20)
	ctx.CapabilityStatus = "out_of_scope"
	out, err := NewConsultationGraph(&MemoryAuditWriter{}).Run(context.Background(), ctx)
	if err != nil {
		t.Fatalf("run consultation graph: %v", err)
	}
	if out.RuleVerdict.Status != model.VerdictRejected {
		t.Fatalf("status=%s want=%s", out.RuleVerdict.Status, model.VerdictRejected)
	}
	if !hasAudit(out.AuditEvents, model.AuditActionGenerateDecision) {
		t.Fatalf("expected decision audit event: %+v", out.AuditEvents)
	}
}

func TestDailyEinoGraphExposesNodeLevelPlan(t *testing.T) {
	graph, err := BuildDailyEinoGraph(context.Background(), &MemoryAuditWriter{}, WorkflowDependencies{})
	if err != nil {
		t.Fatalf("BuildDailyEinoGraph: %v", err)
	}
	want := []string{"StateSnapshotNode", "EvidenceRetrievalNode", "ValueAnalystNode", "TrendRiskOfficerNode", "ExpectedReturnNode", "RuleArbitrationNode", "DecisionRecordNode"}
	if !sameStrings(graph.NodeNames(), want) {
		t.Fatalf("node names=%+v want=%+v", graph.NodeNames(), want)
	}
	if !sameStrings(graph.RegisteredNodeNames(), want) {
		t.Fatalf("registered node names=%+v want=%+v", graph.RegisteredNodeNames(), want)
	}
}

func TestConsultationEinoGraphExposesNodeLevelPlan(t *testing.T) {
	graph, err := BuildConsultationEinoGraph(context.Background(), &MemoryAuditWriter{}, WorkflowDependencies{})
	if err != nil {
		t.Fatalf("BuildConsultationEinoGraph: %v", err)
	}
	want := []string{"StateSnapshotNode", "CapabilityCheckNode", "EvidenceRetrievalNode", "ValueAnalystNode", "TrendRiskOfficerNode", "ExpectedReturnNode", "RuleArbitrationNode", "DecisionRecordNode"}
	if !sameStrings(graph.NodeNames(), want) {
		t.Fatalf("node names=%+v want=%+v", graph.NodeNames(), want)
	}
	if !sameStrings(graph.RegisteredNodeNames(), want) {
		t.Fatalf("registered node names=%+v want=%+v", graph.RegisteredNodeNames(), want)
	}
}

func sampleWorkflowContext(sampleCount int) WorkflowContext {
	return WorkflowContext{
		RequestID:                 "req_daily",
		WorkflowType:              WorkflowDailyDiscipline,
		Symbol:                    "510300",
		RuleVersion:               "v3.0",
		CapabilityStatus:          CapabilityInScope,
		PortfolioSnapshot:         model.PortfolioSnapshot{SnapshotID: "ps_1", CashRatio: 0.2, TotalAssets: 100000},
		MarketSnapshot:            model.MarketSnapshot{MarketSnapshotID: "ms_1", Symbol: "510300", PEPercentile: 50, PBPercentile: 50, SentimentState: model.SentimentNeutral},
		EvidenceSet:               model.EvidenceSet{VerificationStatus: model.VerificationSatisfied, Items: []model.Evidence{{EvidenceID: "ev_1", SourceLevel: model.SourceLevelA, Role: model.EvidenceFormal, EventType: model.EventNormal, IndependentSourceCount: 2, HighGradeIndependentSourceCount: 1}}},
		ExpectedReturnSampleCount: sampleCount,
	}
}

func hasString(values []string, want string) bool {
	for _, v := range values {
		if v == want {
			return true
		}
	}
	return false
}

func hasAudit(events []model.AuditEvent, want model.AuditAction) bool {
	for _, event := range events {
		if event.Action == want {
			return true
		}
	}
	return false
}

func sameStrings(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
