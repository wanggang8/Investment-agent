package rule

import (
	"errors"
	"strings"
	"testing"

	"investment-agent/internal/domain/model"
)

func formalEvidence() []model.Evidence {
	return []model.Evidence{{Role: model.EvidenceFormal, SourceLevel: model.SourceLevelA, EventType: model.EventNormal, IndependentSourceCount: 2, HighGradeIndependentSourceCount: 1}}
}

func TestEvaluatePriorityScenarios(t *testing.T) {
	cases := []struct {
		name string
		in   EvaluationInput
		want model.FinalVerdictStatus
	}{
		{name: "out of capability", in: EvaluationInput{CapabilityStatus: "out_of_scope", HasEvidence: true}, want: model.VerdictRejected},
		{name: "insufficient evidence", in: EvaluationInput{}, want: model.VerdictInsufficientData},
		{name: "major event lacks high grade sources", in: EvaluationInput{HasEvidence: true, Evidence: []model.Evidence{{Role: model.EvidenceFormal, SourceLevel: model.SourceLevelA, EventType: model.EventMajorNegative, HighGradeIndependentSourceCount: 1}}}, want: model.VerdictFrozenWatch},
		{name: "buy logic broken", in: EvaluationInput{HasEvidence: true, BuyLogicBroken: true, Evidence: formalEvidence()}, want: model.VerdictSellOnly},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Evaluate(model.WorkflowContext{}, tc.in)
			if got.Status != tc.want {
				t.Fatalf("status=%s want=%s", got.Status, tc.want)
			}
		})
	}
}

func TestP88SourceVerifiedBuyLogicTransitions(t *testing.T) {
	verifiedBreak := []model.Evidence{{Role: model.EvidenceFormal, SourceLevel: model.SourceLevelA, EventType: model.EventBuyLogicBreak, IndependentSourceCount: 2, HighGradeIndependentSourceCount: 2}}
	v := Evaluate(model.WorkflowContext{}, EvaluationInput{HasEvidence: true, Evidence: verifiedBreak})
	if v.Status != model.VerdictSellOnly {
		t.Fatalf("source-verified buy logic break must enter sell_only, got %+v", v)
	}
	for _, want := range []string{"新增买入", "加仓"} {
		if !containsString(v.ProhibitedActions, want) {
			t.Fatalf("expected prohibited action %q, got %+v", want, v.ProhibitedActions)
		}
	}
	if !containsTriggeredRule(v.TriggeredRules, "BUY_LOGIC") || !strings.Contains(v.Text, "A/S 独立信源=2") {
		t.Fatalf("expected source-count provenance in sell-only verdict, got %+v", v)
	}

	for _, eventType := range []model.EventType{model.EventBuyLogicBreak, model.EventMajorPositive, model.EventMajorNegative} {
		t.Run(string(eventType), func(t *testing.T) {
			frozen := Evaluate(model.WorkflowContext{}, EvaluationInput{HasEvidence: true, Evidence: []model.Evidence{{Role: model.EvidenceFormal, SourceLevel: model.SourceLevelA, EventType: eventType, IndependentSourceCount: 1, HighGradeIndependentSourceCount: 1}}})
			if frozen.Status != model.VerdictFrozenWatch {
				t.Fatalf("insufficient high-grade major evidence must enter frozen_watch, got %+v", frozen)
			}
			if !strings.Contains(frozen.Text, "A/S 独立信源=1") || !containsString(frozen.ProhibitedActions, "主动交易建议") {
				t.Fatalf("expected source-count provenance and pause guidance, got %+v", frozen)
			}
		})
	}
}

func TestEvaluateValuationAndCashRules(t *testing.T) {
	cases := []struct {
		name           string
		in             EvaluationInput
		wantStatus     model.FinalVerdictStatus
		wantRules      int
		wantProhibited string
		forbidOptional string
	}{
		{name: "high risk valuation", in: EvaluationInput{HasEvidence: true, Evidence: formalEvidence(), PEPercentile: 90, PBPercentile: 40, CashRatio: 0.04}, wantStatus: model.VerdictHold, wantRules: 2, wantProhibited: "新增买入"},
		{name: "observation valuation", in: EvaluationInput{HasEvidence: true, Evidence: formalEvidence(), PEPercentile: 63, PBPercentile: 40, CashRatio: 0.08}, wantStatus: model.VerdictHold, wantRules: 1},
		{name: "comfortable valuation", in: EvaluationInput{HasEvidence: true, Evidence: formalEvidence(), PEPercentile: 40, PBPercentile: 35, CashRatio: 0.08}, wantStatus: model.VerdictBuyAllowed, wantRules: 1},
		{name: "low valuation", in: EvaluationInput{HasEvidence: true, Evidence: formalEvidence(), PEPercentile: 20, PBPercentile: 20, CashRatio: 0.08}, wantStatus: model.VerdictBuyAllowed, wantRules: 1},
		{name: "normal cash does not prohibit", in: EvaluationInput{HasEvidence: true, Evidence: formalEvidence(), PEPercentile: 40, PBPercentile: 35, CashRatio: 0.08}, wantStatus: model.VerdictBuyAllowed, wantRules: 1},
		{name: "danger liquidity blocks new buy", in: EvaluationInput{HasEvidence: true, Evidence: formalEvidence(), PEPercentile: 20, PBPercentile: 20, CashRatio: 0.08, LiquidityState: model.LiquidityDanger}, wantStatus: model.VerdictHold, wantRules: 2, wantProhibited: "新增买入", forbidOptional: "分批配置"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Evaluate(model.WorkflowContext{}, tc.in)
			if got.Status != tc.wantStatus || len(got.TriggeredRules) < tc.wantRules {
				t.Fatalf("got %+v", got)
			}
			if tc.wantProhibited != "" && !containsString(got.ProhibitedActions, tc.wantProhibited) {
				t.Fatalf("expected prohibited action %s, got %+v", tc.wantProhibited, got.ProhibitedActions)
			}
			if tc.forbidOptional != "" && containsString(got.OptionalActions, tc.forbidOptional) {
				t.Fatalf("optional action %s should not appear when prohibited, got %+v", tc.forbidOptional, got.OptionalActions)
			}
		})
	}
}

func TestEvaluateValuationHighRiskBoundaryAtEightyPercent(t *testing.T) {
	v := Evaluate(model.WorkflowContext{}, EvaluationInput{HasEvidence: true, Evidence: formalEvidence(), PEPercentile: 80, PBPercentile: 40, CashRatio: 0.08})

	if v.Status != model.VerdictHold || !containsString(v.ProhibitedActions, "新增买入") {
		t.Fatalf("80%% valuation percentile must enter high-risk no-new-buy zone, got %+v", v)
	}
}

func TestEvaluateDoesNotTriggerAllocationWhenRatiosAreUnknown(t *testing.T) {
	v := Evaluate(model.WorkflowContext{}, EvaluationInput{HasEvidence: true, Evidence: formalEvidence(), PEPercentile: 40, PBPercentile: 35, CashRatio: 0.08})

	if containsString(v.OptionalActions, "再平衡卫星资产") || containsTriggeredRule(v.TriggeredRules, "ALLOCATION") {
		t.Fatalf("unknown allocation ratios must not trigger rebalance rules, got %+v", v)
	}
}

func TestEvaluateTakeProfitAndAllocationRules(t *testing.T) {
	v := Evaluate(model.WorkflowContext{}, EvaluationInput{HasEvidence: true, Evidence: formalEvidence(), UnrealizedProfitRatio: 0.21})
	if !containsString(v.OptionalActions, "卖出 30%") {
		t.Fatalf("expected 20%% take profit action, got %+v", v)
	}

	v = Evaluate(model.WorkflowContext{}, EvaluationInput{HasEvidence: true, Evidence: formalEvidence(), UnrealizedProfitRatio: 0.31, HandledProfit20: true})
	if !containsString(v.OptionalActions, "再卖出 30%") {
		t.Fatalf("expected 30%% take profit action, got %+v", v)
	}

	v = Evaluate(model.WorkflowContext{}, EvaluationInput{HasEvidence: true, Evidence: formalEvidence(), TakeProfitStarted: true, StageHighPrice: 100, CurrentPrice: 90})
	if v.Status != model.VerdictReduce {
		t.Fatalf("expected reduce, got %s", v.Status)
	}

	v = Evaluate(model.WorkflowContext{}, EvaluationInput{HasEvidence: true, Evidence: formalEvidence(), SatelliteRatio: 0.40})
	if len(v.TriggeredRules) == 0 {
		t.Fatal("expected satellite allocation rule")
	}

	v = Evaluate(model.WorkflowContext{}, EvaluationInput{HasEvidence: true, Evidence: formalEvidence(), CoreRatio: 0.50})
	if len(v.TriggeredRules) == 0 {
		t.Fatal("expected core allocation rule")
	}
}

func TestEvaluateSentimentAndSourceRules(t *testing.T) {
	v := Evaluate(model.WorkflowContext{}, EvaluationInput{HasEvidence: true, Evidence: formalEvidence(), SentimentState: model.SentimentExtreme})
	if !containsString(v.ProhibitedActions, "主动交易建议") {
		t.Fatal("expected sentiment prohibition")
	}

	v = Evaluate(model.WorkflowContext{}, EvaluationInput{HasEvidence: true, Evidence: []model.Evidence{{Role: model.EvidenceFormal, SourceLevel: model.SourceLevelC}}})
	if v.Status != model.VerdictInsufficientData {
		t.Fatalf("expected insufficient data, got %s", v.Status)
	}

	v = Evaluate(model.WorkflowContext{}, EvaluationInput{HasEvidence: true, SourceVerificationStatus: model.VerificationBackgroundOnly, Evidence: []model.Evidence{{Role: model.EvidenceBackground, SourceLevel: model.SourceLevelC, EventType: model.EventNormal, IndependentSourceCount: 1}}, PEPercentile: 20, PBPercentile: 20, CashRatio: 0.2})
	if v.Status == model.VerdictBuyAllowed {
		t.Fatalf("background-only evidence must not allow trading verdict: %+v", v)
	}
}

func TestExpectedReturnDoesNotOverrideVerdict(t *testing.T) {
	v := Evaluate(model.WorkflowContext{}, EvaluationInput{HasEvidence: true, Evidence: formalEvidence(), PEPercentile: 90, PBPercentile: 90})
	if v.Status != model.VerdictHold {
		t.Fatalf("expected hold, got %s", v.Status)
	}
	if len(v.ExpectedReturns) != 3 {
		t.Fatalf("expected scenarios")
	}
}

func TestAdvanceProposal(t *testing.T) {
	status, apply, err := AdvanceProposal(model.ProposalDraft, 3, true, "")
	if err != nil || apply || status != model.ProposalPendingUserConfirm {
		t.Fatalf("unexpected: %s %v %v", status, apply, err)
	}

	_, _, err = AdvanceProposal(model.ProposalPendingUserConfirm, 2, true, "")
	if !errors.Is(err, ErrInsufficientSamples) {
		t.Fatalf("expected insufficient samples, got %v", err)
	}

	status, _, err = AdvanceProposal(model.ProposalPendingUserConfirm, 3, true, "")
	if err != nil || status != model.ProposalUnderGatekeeperAudit {
		t.Fatalf("unexpected audit transition: %s %v", status, err)
	}

	status, _, err = AdvanceProposal(model.ProposalUnderGatekeeperAudit, 3, true, model.AuditApproved)
	if err != nil || status != model.ProposalPendingFinalConfirm {
		t.Fatalf("unexpected final confirm transition: %s %v", status, err)
	}

	status, apply, err = AdvanceProposal(model.ProposalPendingFinalConfirm, 3, true, "")
	if err != nil || !apply || status != model.ProposalApplied {
		t.Fatalf("unexpected applied transition: %s %v %v", status, apply, err)
	}

	status, apply, err = AdvanceProposal(model.ProposalPendingFinalConfirm, 3, false, "")
	if err != nil || apply || status != model.ProposalRejected {
		t.Fatalf("unexpected final reject transition: %s %v %v", status, apply, err)
	}

	_, _, err = AdvanceProposal(model.ProposalPendingFinalConfirm, 2, true, "")
	if !errors.Is(err, ErrInsufficientSamples) {
		t.Fatalf("expected final insufficient samples, got %v", err)
	}

	_, _, err = AdvanceProposal(model.ProposalApplied, 3, true, "")
	if !errors.Is(err, ErrTerminalProposal) {
		t.Fatalf("expected terminal error, got %v", err)
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

func containsTriggeredRule(values []model.TriggeredRule, want string) bool {
	for _, value := range values {
		if value.RuleID == want {
			return true
		}
	}
	return false
}
