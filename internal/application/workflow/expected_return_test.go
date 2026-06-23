package workflow

import (
	"context"
	"testing"

	"investment-agent/internal/domain/model"
)

func TestBuildExpectedReturnIncludesSampleContextForAllPrecisionStates(t *testing.T) {
	tests := []struct {
		name       string
		samples    int
		wantStatus model.PrecisionStatus
		wantCount  int
		wantProb   bool
		wantRanges bool
	}{
		{name: "available", samples: 20, wantStatus: model.PrecisionAvailable, wantCount: 3, wantProb: true, wantRanges: true},
		{name: "insufficient", samples: 5, wantStatus: model.PrecisionInsufficient, wantCount: 3, wantProb: false, wantRanges: true},
		{name: "unavailable", samples: 4, wantStatus: model.PrecisionUnavailable, wantCount: 0, wantProb: false, wantRanges: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := BuildExpectedReturn(tt.samples)
			if out.PrecisionStatus != tt.wantStatus || out.SampleCount != tt.samples {
				t.Fatalf("unexpected status/count: %+v", out)
			}
			if out.SampleWindow == "" || out.ScreeningCondition == "" {
				t.Fatalf("sample context missing: %+v", out)
			}
			if len(out.Scenarios) != tt.wantCount {
				t.Fatalf("scenario count=%d want=%d: %+v", len(out.Scenarios), tt.wantCount, out.Scenarios)
			}
			for _, scenario := range out.Scenarios {
				if scenario.Trigger == "" {
					t.Fatalf("scenario trigger missing: %+v", scenario)
				}
				if (scenario.Probability != nil) != tt.wantProb {
					t.Fatalf("probability presence mismatch: %+v", scenario)
				}
				if tt.wantRanges && scenario.ReturnRange == "" {
					t.Fatalf("return range missing: %+v", scenario)
				}
			}
		})
	}
}

func TestBuildExpectedReturnProducesAdvisorySellEvaluation(t *testing.T) {
	out := BuildExpectedReturnWithContext(ExpectedReturnInput{SampleCount: 20, CurrentPrice: 12, BasePrice: 10, PreviousBaseMidpoint: 0.06, TargetReturnRate: 0.15})

	if out.SellEvaluation.Status != "triggered" {
		t.Fatalf("expected triggered sell evaluation, got %+v", out.SellEvaluation)
	}
	if len(out.SellEvaluation.Triggers) == 0 || len(out.SellEvaluation.Actions) == 0 {
		t.Fatalf("expected advisory triggers and actions: %+v", out.SellEvaluation)
	}
	if out.SellEvaluation.NonTradingDisclaimer == "" {
		t.Fatalf("non-trading disclaimer missing: %+v", out.SellEvaluation)
	}
	if out.ReassessmentTrigger.Reason == "" || out.ReassessmentTrigger.CurrentValue == 0 {
		t.Fatalf("reassessment trigger missing: %+v", out.ReassessmentTrigger)
	}
}

func TestBuildExpectedReturnDoesNotTriggerTargetWithoutConfiguredTarget(t *testing.T) {
	out := BuildExpectedReturnWithContext(ExpectedReturnInput{SampleCount: 20, CurrentPrice: 12, BasePrice: 10})

	for _, trigger := range out.SellEvaluation.Triggers {
		if trigger == "target_return_reached" {
			t.Fatalf("target return must not trigger without explicit target: %+v", out.SellEvaluation)
		}
	}
}

func TestBuildExpectedReturnUsesScenarioBoundsForSellTriggers(t *testing.T) {
	out := BuildExpectedReturnWithContext(ExpectedReturnInput{SampleCount: 20, CurrentPrice: 10.4, BasePrice: 10})
	for _, trigger := range out.SellEvaluation.Triggers {
		if trigger == "base_upper_bound_exceeded" {
			t.Fatalf("base upper bound must not trigger below the base scenario upper bound: %+v", out.SellEvaluation)
		}
	}

	out = BuildExpectedReturnWithContext(ExpectedReturnInput{SampleCount: 20, CurrentPrice: 10.9, BasePrice: 10})
	found := false
	for _, trigger := range out.SellEvaluation.Triggers {
		if trigger == "base_upper_bound_exceeded" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected base upper bound trigger above the base scenario upper bound: %+v", out.SellEvaluation)
	}
}

func TestBuildExpectedReturnCoversAllSellEvaluationTriggers(t *testing.T) {
	out := BuildExpectedReturnWithContext(ExpectedReturnInput{SampleCount: 20, CurrentPrice: 10.8, BasePrice: 10, PreviousBaseMidpoint: 0.2, TargetReturnRate: 0.05})
	for _, want := range []string{"upside_lower_bound_reached", "base_upper_bound_exceeded", "base_midpoint_downshift", "target_return_reached"} {
		if !containsString(out.SellEvaluation.Triggers, want) {
			t.Fatalf("expected trigger %s, got %+v", want, out.SellEvaluation)
		}
	}
	if out.ReassessmentTrigger.Boundary != "base_midpoint_downshift" {
		t.Fatalf("expected base midpoint reassessment, got %+v", out.ReassessmentTrigger)
	}

	out = BuildExpectedReturnWithContext(ExpectedReturnInput{SampleCount: 20, CurrentPrice: 8.7, BasePrice: 10})
	if !containsString(out.SellEvaluation.Triggers, "downside_lower_bound_breached") {
		t.Fatalf("expected downside trigger, got %+v", out.SellEvaluation)
	}
}

func TestP88ExpectedReturnUsesHistoricalSampleProbabilitiesAndCoverage(t *testing.T) {
	out := BuildExpectedReturnWithContext(ExpectedReturnInput{
		SampleCount:      30,
		TargetName:       "沪深300ETF",
		TargetCode:       "510300",
		HoldingClass:     "broad_index_etf",
		HorizonLabel:     "未来 12 个月",
		SampleWindow:     "2021-2026 similar valuation samples",
		MarketState:      "neutral",
		FundamentalState: "in_line",
		HistoricalSamples: []ExpectedReturnHistoricalSample{
			{Scenario: "upside", Count: 6, ReturnRange: "12.00%~18.00%", LowerBound: 0.12, UpperBound: 0.18, ReturnRate: 0.15, Trigger: "业绩超预期，估值提升至历史高位"},
			{Scenario: "base", Count: 18, ReturnRange: "4.00%~9.00%", LowerBound: 0.04, UpperBound: 0.09, ReturnRate: 0.065, Trigger: "业绩符合预期，估值维持当前水平"},
			{Scenario: "downside", Count: 6, ReturnRange: "-10.00%~-2.00%", LowerBound: -0.10, UpperBound: -0.02, ReturnRate: -0.06, Trigger: "业绩低于预期，估值收缩"},
		},
		HoldingClassCoverage: []ExpectedReturnHoldingClassCoverage{
			{HoldingClass: "broad_index_etf", Symbol: "510300", Status: "covered"},
			{HoldingClass: "sector_growth_fund", Symbol: "159915", Status: "covered"},
			{HoldingClass: "equity_constituent_financial", Symbol: "600000", Status: "covered"},
		},
	})

	if out.TargetName != "沪深300ETF" || out.TargetCode != "510300" || out.HorizonLabel != "未来 12 个月" {
		t.Fatalf("target identity or horizon missing: %+v", out)
	}
	if out.PrecisionStatus != model.PrecisionAvailable || len(out.Scenarios) != 3 {
		t.Fatalf("expected available historical scenarios, got %+v", out)
	}
	if out.Scenarios[0].Name != "upside" || out.Scenarios[1].Name != "base" || out.Scenarios[2].Name != "downside" {
		t.Fatalf("scenario order must remain report-stable, got %+v", out.Scenarios)
	}
	if out.Scenarios[0].Probability == nil || *out.Scenarios[0].Probability != 0.2 || *out.Scenarios[1].Probability != 0.6 || *out.Scenarios[2].Probability != 0.2 {
		t.Fatalf("probabilities must come from sample proportions, got %+v", out.Scenarios)
	}
	if out.ProbabilityBasis == "" || !containsString(out.CoveredHoldingClasses(), "equity_constituent_financial") {
		t.Fatalf("expected probability basis and holding class coverage, got %+v", out)
	}
}

func TestP88ExpectedReturnDynamicMonitoringAndLowSampleSupplementData(t *testing.T) {
	downshift := BuildExpectedReturnWithContext(ExpectedReturnInput{
		SampleCount:      20,
		MarketState:      "stress",
		FundamentalState: "below_expectation",
		HistoricalSamples: []ExpectedReturnHistoricalSample{
			{Scenario: "upside", Count: 8, ReturnRange: "8.00%~15.00%", LowerBound: 0.08, UpperBound: 0.15, ReturnRate: 0.12, Trigger: "估值修复"},
			{Scenario: "base", Count: 9, ReturnRange: "0.00%~8.00%", LowerBound: 0, UpperBound: 0.08, ReturnRate: 0.04, Trigger: "维持当前"},
			{Scenario: "downside", Count: 3, ReturnRange: "-12.00%~0.00%", LowerBound: -0.12, UpperBound: 0, ReturnRate: -0.06, Trigger: "估值收缩"},
		},
		AssumptionChecks:      []ExpectedReturnAssumptionCheck{{Name: "盈利增速", Expected: 0.08, Actual: 0.01, MonthsBelow: 2}},
		PessimisticPathMonths: 1,
	})
	for _, want := range []string{"scenario_probability_downshift", "two_month_assumption_downshift", "one_month_pessimistic_path"} {
		if !containsString(downshift.SellEvaluation.Triggers, want) {
			t.Fatalf("expected monitoring trigger %s, got %+v", want, downshift.SellEvaluation)
		}
	}

	lowSample := BuildExpectedReturnWithContext(ExpectedReturnInput{SampleCount: 4, MissingCategories: []string{"capital_flow", "margin_financing", "constituent_financial"}})
	if lowSample.PrecisionStatus != model.PrecisionUnavailable || len(lowSample.Scenarios) != 0 {
		t.Fatalf("sample below five must not generate ranges, got %+v", lowSample)
	}
	for _, want := range []string{"capital_flow", "margin_financing", "constituent_financial"} {
		if !containsString(lowSample.SupplementData, want) {
			t.Fatalf("expected supplement data %s, got %+v", want, lowSample.SupplementData)
		}
	}
}

func TestP89ExpectedReturnInputReadsDynamicMonitoringMetadata(t *testing.T) {
	wf := &WorkflowContext{
		Symbol: "510300",
		MarketSnapshot: model.MarketSnapshot{
			Symbol:            "510300",
			MarketMetricsJSON: `{"metadata":{"expected_return_market_state":"stress","expected_return_fundamental_state":"below_expectation","expected_return_pessimistic_path_months":1,"expected_return_assumption_checks":[{"name":"盈利增速","expected":0.08,"actual":0.01,"months_below":2}]}}`,
		},
	}

	input := expectedReturnInputFromWorkflow(wf)

	if input.MarketState != "stress" || input.FundamentalState != "below_expectation" || input.PessimisticPathMonths != 1 {
		t.Fatalf("expected dynamic monitoring metadata in input, got %+v", input)
	}
	if len(input.AssumptionChecks) != 1 || input.AssumptionChecks[0].MonthsBelow != 2 {
		t.Fatalf("expected assumption checks from metadata, got %+v", input.AssumptionChecks)
	}
}

func TestP89ExpectedReturnLowersAffectedProbabilitiesWhenInputsDeteriorate(t *testing.T) {
	samples := []ExpectedReturnHistoricalSample{
		{Scenario: "upside", Count: 6, ReturnRange: "8.00%~15.00%", LowerBound: 0.08, UpperBound: 0.15, ReturnRate: 0.12, Trigger: "估值修复"},
		{Scenario: "base", Count: 18, ReturnRange: "0.00%~8.00%", LowerBound: 0, UpperBound: 0.08, ReturnRate: 0.04, Trigger: "维持当前"},
		{Scenario: "downside", Count: 6, ReturnRange: "-12.00%~0.00%", LowerBound: -0.12, UpperBound: 0, ReturnRate: -0.06, Trigger: "估值收缩"},
	}
	baseline := BuildExpectedReturnWithContext(ExpectedReturnInput{SampleCount: 30, MarketState: "neutral", FundamentalState: "in_line", HistoricalSamples: samples})
	downshift := BuildExpectedReturnWithContext(ExpectedReturnInput{SampleCount: 30, MarketState: "stress", FundamentalState: "below_expectation", HistoricalSamples: samples})

	if baseline.Scenarios[0].Probability == nil || downshift.Scenarios[0].Probability == nil || baseline.Scenarios[1].Probability == nil || downshift.Scenarios[1].Probability == nil {
		t.Fatalf("expected precise probabilities, baseline=%+v downshift=%+v", baseline.Scenarios, downshift.Scenarios)
	}
	if !(*downshift.Scenarios[0].Probability < *baseline.Scenarios[0].Probability) {
		t.Fatalf("upside probability should be lower after deterioration, baseline=%+v downshift=%+v", baseline.Scenarios, downshift.Scenarios)
	}
	if !(*downshift.Scenarios[1].Probability < *baseline.Scenarios[1].Probability) {
		t.Fatalf("base probability should be lower after deterioration, baseline=%+v downshift=%+v", baseline.Scenarios, downshift.Scenarios)
	}
	if downshift.Scenarios[2].Probability == nil || !(*downshift.Scenarios[2].Probability > *baseline.Scenarios[2].Probability) {
		t.Fatalf("downside probability should absorb the probability downshift, baseline=%+v downshift=%+v", baseline.Scenarios, downshift.Scenarios)
	}
}

func TestP89ExpectedReturnCarriesExtremeFearHistoricalSimilarContext(t *testing.T) {
	out := BuildExpectedReturnWithContext(ExpectedReturnInput{
		SampleCount:    30,
		SentimentState: "extreme",
		HistoricalSamples: []ExpectedReturnHistoricalSample{
			{Scenario: "upside", Count: 4, ReturnRange: "6.00%~12.00%", LowerBound: 0.06, UpperBound: 0.12, ReturnRate: 0.09, Trigger: "恐慌后估值修复"},
			{Scenario: "base", Count: 10, ReturnRange: "-2.00%~6.00%", LowerBound: -0.02, UpperBound: 0.06, ReturnRate: 0.02, Trigger: "情绪低位震荡"},
			{Scenario: "downside", Count: 6, ReturnRange: "-18.00%~-4.00%", LowerBound: -0.18, UpperBound: -0.04, ReturnRate: -0.10, Trigger: "恐慌继续扩散"},
		},
		HistoricalContexts: []ExpectedReturnHistoricalContext{
			{Label: "极端恐惧样本", Window: "2018Q4, 2020Q1, 2022Q4", SampleCount: 20, Outcome: "先暂停主动交易建议，再等待正式证据复核", MaxDrawdown: -0.18, Recovery: "多数样本在 3-9 个月内完成估值修复", Source: "local_public_history"},
		},
	})

	if len(out.HistoricalContexts) != 1 || out.HistoricalContexts[0].Label != "极端恐惧样本" {
		t.Fatalf("expected historical similar context, got %+v", out.HistoricalContexts)
	}
	if !containsString(out.SellEvaluation.Triggers, "extreme_fear_historical_context") {
		t.Fatalf("expected extreme fear context trigger, got %+v", out.SellEvaluation)
	}
	if !containsString(out.SellEvaluation.Actions, "暂停主动交易建议") {
		t.Fatalf("expected active trading lock guidance, got %+v", out.SellEvaluation)
	}
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func TestExpectedReturnNodeUsesWorkflowPricesForSellEvaluation(t *testing.T) {
	wf := &WorkflowContext{
		Symbol:                    "510300",
		ExpectedReturnSampleCount: 20,
		MarketSnapshot:            model.MarketSnapshot{ClosePrice: 12},
		PositionSnapshots:         []model.Position{{Symbol: "510300", CostPrice: 10, CurrentPrice: 12}},
	}

	result := RunExpectedReturnNode(context.Background(), wf, WorkflowDependencies{})

	if result.Status != StatusSuccess {
		t.Fatalf("expected success, got %+v", result)
	}
	if wf.ExpectedReturnSellEvaluation.Status != "triggered" {
		t.Fatalf("expected workflow price context to trigger sell evaluation, got %+v", wf.ExpectedReturnSellEvaluation)
	}
	if len(wf.ExpectedReturnSellEvaluation.Triggers) == 0 {
		t.Fatalf("expected advisory triggers from workflow context: %+v", wf.ExpectedReturnSellEvaluation)
	}
	if wf.ExpectedReturnReassessmentTrigger.Reason == "" {
		t.Fatalf("expected reassessment trigger from workflow context")
	}
}

func TestExpectedReturnNodeUsesMatchingSymbolPosition(t *testing.T) {
	wf := &WorkflowContext{
		Symbol:                    "159915",
		ExpectedReturnSampleCount: 20,
		MarketSnapshot:            model.MarketSnapshot{ClosePrice: 4},
		PositionSnapshots: []model.Position{
			{Symbol: "510300", CostPrice: 10, CurrentPrice: 12},
			{Symbol: "159915", CostPrice: 4, CurrentPrice: 4},
		},
	}

	result := RunExpectedReturnNode(context.Background(), wf, WorkflowDependencies{})

	if result.Status != StatusSuccess {
		t.Fatalf("expected success, got %+v", result)
	}
	if wf.ExpectedReturnSellEvaluation.Status != "not_triggered" {
		t.Fatalf("expected matching symbol position to avoid unrelated sell trigger, got %+v", wf.ExpectedReturnSellEvaluation)
	}
}

func TestExpectedReturnNodeUsesWorkflowDynamicSellInputs(t *testing.T) {
	wf := &WorkflowContext{
		Symbol:                             "510300",
		ExpectedReturnSampleCount:          20,
		ExpectedReturnPreviousBaseMidpoint: 0.2,
		ExpectedReturnTargetReturnRate:     0.15,
		MarketSnapshot:                     model.MarketSnapshot{ClosePrice: 12},
		PositionSnapshots:                  []model.Position{{Symbol: "510300", CostPrice: 10, CurrentPrice: 12}},
	}

	result := RunExpectedReturnNode(context.Background(), wf, WorkflowDependencies{})

	if result.Status != StatusSuccess {
		t.Fatalf("expected success, got %+v", result)
	}
	if wf.ExpectedReturnReassessmentTrigger.Boundary != "base_midpoint_downshift" {
		t.Fatalf("expected workflow previous base midpoint to trigger reassessment, got %+v", wf.ExpectedReturnReassessmentTrigger)
	}
	foundTarget := false
	for _, trigger := range wf.ExpectedReturnSellEvaluation.Triggers {
		if trigger == "target_return_reached" {
			foundTarget = true
		}
	}
	if !foundTarget {
		t.Fatalf("expected workflow target return to trigger advisory sell evaluation, got %+v", wf.ExpectedReturnSellEvaluation)
	}
}

func TestExpectedReturnNodeIncludesP34SupportingDataContext(t *testing.T) {
	wf := &WorkflowContext{
		Symbol:                    "000300",
		ExpectedReturnSampleCount: 20,
		MarketSnapshot:            model.MarketSnapshot{MarketSnapshotID: "market_p34", Symbol: "000300", ClosePrice: 4, MarketMetricsJSON: `{"metadata":{"p34_source_health":{"index_constituents":"fresh","sentiment_proxy":"fresh","index_valuation_files":"parse_error"},"p34_data_categories":["index_constituents","sentiment_proxy","index_valuation_files"]}}`},
		PositionSnapshots:         []model.Position{{Symbol: "000300", CostPrice: 4, CurrentPrice: 4}},
	}

	result := RunExpectedReturnNode(context.Background(), wf, WorkflowDependencies{})

	if result.Status != StatusSuccess {
		t.Fatalf("expected success, got %+v", result)
	}
	if wf.ExpectedReturnSupportingDataSummary == "" || wf.ExpectedReturnMissingCategories[0] != "index_valuation_files" {
		t.Fatalf("expected P34 freshness context in expected return, summary=%q missing=%+v", wf.ExpectedReturnSupportingDataSummary, wf.ExpectedReturnMissingCategories)
	}
	if wf.ExpectedReturnScreeningCondition == "" || wf.ExpectedReturnSampleWindow == "" {
		t.Fatalf("expected sample context to remain explicit")
	}
}

func TestExpectedReturnSampleCountFromWorkflowDataUsesMarketHistory(t *testing.T) {
	positions := []model.Position{{Symbol: "510300", CostPrice: 4.2}}
	market := model.MarketSnapshot{MarketSnapshotID: "market-1", Symbol: "510300", MarketMetricsJSON: `{"metadata":{"nav_history":[{"nav":4.1},{"nav":4.2},{"nav":4.3}]}}`}

	count := ExpectedReturnSampleCountFromWorkflowData(positions, market)

	if count != 5 {
		t.Fatalf("expected position + market snapshot + nav history count, got %d", count)
	}
}

func TestExpectedReturnSampleCountFromWorkflowDataDoesNotInventSamples(t *testing.T) {
	count := ExpectedReturnSampleCountFromWorkflowData(nil, model.MarketSnapshot{MarketSnapshotID: "market-1", Symbol: "510300", MarketMetricsJSON: `{}`})

	if count != 1 {
		t.Fatalf("expected only the real market snapshot sample, got %d", count)
	}
}

func TestBuildExpectedReturnExplainsMissingPriceContext(t *testing.T) {
	out := BuildExpectedReturnWithContext(ExpectedReturnInput{SampleCount: 20, CurrentPrice: 12})

	if out.SellEvaluation.Status != "not_applicable" {
		t.Fatalf("expected explicit not_applicable sell evaluation, got %+v", out.SellEvaluation)
	}
	if len(out.SellEvaluation.Prompts) == 0 || out.SellEvaluation.Prompts[0] == "" {
		t.Fatalf("expected missing price reason, got %+v", out.SellEvaluation)
	}
}
