package rule

import (
	"testing"

	"investment-agent/internal/domain/model"
)

func TestP75ExecutableCriteriaVectors(t *testing.T) {
	t.Run("sentiment 90 or 10 percentile and abnormal 3-day inputs trigger cooldown state", func(t *testing.T) {
		cases := []SentimentSignal{
			{SentimentPercentile: 90},
			{SentimentPercentile: 10},
			{AbnormalFinancingDays: 3},
			{AbnormalVolumeDays: 3},
			{AbnormalMediaHeatDays: 3},
			{UserEmotionTags: []string{"恐慌"}},
			{UserEmotionTags: []string{"追热点"}},
		}
		for _, tc := range cases {
			if got := DetermineSentimentState(tc); got != model.SentimentExtreme {
				t.Fatalf("expected extreme sentiment for %+v, got %s", tc, got)
			}
		}
	})

	t.Run("liquidity 20-day 20x and same-day 5 percent thresholds block market-style actions", func(t *testing.T) {
		cases := []LiquiditySignal{
			{PlanAmount: 100, Average20DayAmount: 1999, SameDayAmount: 10000},
			{PlanAmount: 501, Average20DayAmount: 20000, SameDayAmount: 10000},
			{SpreadAbnormal: true, PlanAmount: 100, Average20DayAmount: 20000, SameDayAmount: 10000},
		}
		for _, tc := range cases {
			if got := DetermineLiquidityState(tc); got != model.LiquidityDanger {
				t.Fatalf("expected liquidity danger for %+v, got %s", tc, got)
			}
		}
		if got := DetermineLiquidityState(LiquiditySignal{PlanAmount: 500, Average20DayAmount: 10000, SameDayAmount: 10000}); got != model.LiquidityNormal {
			t.Fatalf("expected exact 20x and 5%% boundary to remain normal, got %s", got)
		}
	})
}

func TestP75RulePriorityAndRootRules(t *testing.T) {
	verifiedMajorBreak := []model.Evidence{{Role: model.EvidenceFormal, SourceLevel: model.SourceLevelA, EventType: model.EventBuyLogicBreak, HighGradeIndependentSourceCount: 2}}

	cases := []struct {
		name              string
		in                EvaluationInput
		wantStatus        model.FinalVerdictStatus
		wantProhibited    string
		forbidOptional    string
		wantTriggeredRule string
	}{
		{name: "R-1 multi-source verification freezes major event with one high-grade source", in: EvaluationInput{HasEvidence: true, Evidence: []model.Evidence{{Role: model.EvidenceFormal, SourceLevel: model.SourceLevelA, EventType: model.EventMajorNegative, HighGradeIndependentSourceCount: 1}}}, wantStatus: model.VerdictFrozenWatch, wantProhibited: "主动交易建议", wantTriggeredRule: "SOURCE"},
		{name: "R-2 extreme sentiment pauses active trading", in: EvaluationInput{HasEvidence: true, Evidence: formalEvidence(), SentimentState: model.SentimentExtreme}, wantStatus: model.VerdictHold, wantProhibited: "主动交易建议", wantTriggeredRule: "SENTIMENT"},
		{name: "R-3 satellite limit does not create buy guidance", in: EvaluationInput{HasEvidence: true, Evidence: formalEvidence(), SatelliteRatio: 0.40}, wantStatus: model.VerdictHold, forbidOptional: "新增买入", wantTriggeredRule: "ALLOCATION"},
		{name: "R-4 verified buy logic break enters sell only before valuation buy actions", in: EvaluationInput{HasEvidence: true, Evidence: verifiedMajorBreak, BuyLogicBroken: true, PEPercentile: 20, PBPercentile: 20}, wantStatus: model.VerdictSellOnly, wantProhibited: "新增买入", forbidOptional: "分批配置", wantTriggeredRule: "BUY_LOGIC"},
		{name: "R-5 low cash prohibits new buy", in: EvaluationInput{HasEvidence: true, Evidence: formalEvidence(), PEPercentile: 20, PBPercentile: 20, CashRatio: 0.049}, wantStatus: model.VerdictHold, wantProhibited: "新增买入", forbidOptional: "分批配置", wantTriggeredRule: "R-5"},
		{name: "R-6 liquidity danger blocks new buy and large market action", in: EvaluationInput{HasEvidence: true, Evidence: formalEvidence(), PEPercentile: 20, PBPercentile: 20, LiquidityState: model.LiquidityDanger}, wantStatus: model.VerdictHold, wantProhibited: "大额市价操作", forbidOptional: "分批配置", wantTriggeredRule: "LIQUIDITY"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Evaluate(model.WorkflowContext{}, tc.in)
			if got.Status != tc.wantStatus {
				t.Fatalf("status=%s want=%s verdict=%+v", got.Status, tc.wantStatus, got)
			}
			if tc.wantProhibited != "" && !containsString(got.ProhibitedActions, tc.wantProhibited) {
				t.Fatalf("expected prohibited action %q, got %+v", tc.wantProhibited, got.ProhibitedActions)
			}
			if tc.forbidOptional != "" && containsString(got.OptionalActions, tc.forbidOptional) {
				t.Fatalf("optional action %q should be absent, got %+v", tc.forbidOptional, got.OptionalActions)
			}
			if tc.wantTriggeredRule != "" && !containsTriggeredRule(got.TriggeredRules, tc.wantTriggeredRule) {
				t.Fatalf("expected triggered rule %q, got %+v", tc.wantTriggeredRule, got.TriggeredRules)
			}
		})
	}
}

func TestP75PositionStateAndCooldownBoundaries(t *testing.T) {
	stateCases := []struct {
		name string
		in   EvaluationInput
		want model.PositionState
	}{
		{name: "normal", in: EvaluationInput{HasEvidence: true, Evidence: formalEvidence()}, want: model.PositionNormal},
		{name: "sell only", in: EvaluationInput{HasEvidence: true, Evidence: []model.Evidence{{Role: model.EvidenceFormal, SourceLevel: model.SourceLevelA, EventType: model.EventBuyLogicBreak, HighGradeIndependentSourceCount: 2}}, BuyLogicBroken: true}, want: model.PositionSellOnly},
		{name: "frozen watch", in: EvaluationInput{HasEvidence: true, Evidence: []model.Evidence{{Role: model.EvidenceFormal, SourceLevel: model.SourceLevelA, EventType: model.EventMajorNegative, HighGradeIndependentSourceCount: 1}}}, want: model.PositionFrozenWatch},
	}
	for _, tc := range stateCases {
		t.Run(tc.name, func(t *testing.T) {
			verdict := Evaluate(model.WorkflowContext{}, tc.in)
			if got := PositionStateForVerdict(verdict.Status); got != tc.want {
				t.Fatalf("position state=%s want=%s verdict=%+v", got, tc.want, verdict)
			}
		})
	}

	if got := CooldownTradingDays(2); got != 1 {
		t.Fatalf("expected normal one-day cooldown before 3 consecutive triggers, got %d", got)
	}
	if got := CooldownTradingDays(3); got != 5 {
		t.Fatalf("expected 5-day cooldown at 3 consecutive triggers, got %d", got)
	}
}

func TestP75PortfolioAllocationAndTakeProfitReadback(t *testing.T) {
	v := Evaluate(model.WorkflowContext{}, EvaluationInput{HasEvidence: true, Evidence: formalEvidence(), CoreRatio: 0.49})
	if !containsString(v.OptionalActions, "提高核心资产占比") || !containsTriggeredRule(v.TriggeredRules, "ALLOCATION") {
		t.Fatalf("expected core underweight rebalance readback, got %+v", v)
	}

	v = Evaluate(model.WorkflowContext{}, EvaluationInput{HasEvidence: true, Evidence: formalEvidence(), SatelliteRatio: 0.41})
	if !containsString(v.OptionalActions, "再平衡卫星资产") || !containsTriggeredRule(v.TriggeredRules, "ALLOCATION") {
		t.Fatalf("expected satellite over-limit readback, got %+v", v)
	}

	v = Evaluate(model.WorkflowContext{}, EvaluationInput{HasEvidence: true, Evidence: formalEvidence(), UnrealizedProfitRatio: 0.21, SatelliteRatio: 0.41})
	if !containsString(v.OptionalActions, "止盈资金优先回归核心资产") {
		t.Fatalf("expected take-profit funds to return to core assets, got %+v", v)
	}
}
