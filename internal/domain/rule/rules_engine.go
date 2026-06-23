package rule

import "investment-agent/internal/domain/model"

// EvaluationInput 汇总规则裁决需要的关键指标。
// 输入来自 WorkflowContext、账户快照、市场快照和证据验证结果。
type EvaluationInput struct {
	CapabilityStatus         string
	HasEvidence              bool
	Evidence                 []model.Evidence
	SourceVerificationStatus model.VerificationStatus
	BuyLogicBroken           bool
	SentimentState           model.SentimentState
	LiquidityState           model.LiquidityState
	PEPercentile             float64
	PBPercentile             float64
	CashRatio                float64
	CoreRatio                float64
	SatelliteRatio           float64
	UnrealizedProfitRatio    float64
	TakeProfitStarted        bool
	StageHighPrice           float64
	CurrentPrice             float64
	HandledProfit20          bool
	HandledProfit30          bool
}

// Evaluate 是领域规则裁决入口。
// 安全类规则先返回终态，风险类规则只追加禁止事项或可选动作。
func Evaluate(_ model.WorkflowContext, in EvaluationInput) model.RuleVerdict {
	verdict := baseVerdict()
	if rejected, ok := evaluateCapability(in); ok {
		return rejected
	}
	if insufficient, ok := evaluateEvidence(in); ok {
		return insufficient
	}
	if frozen, ok := evaluateSource(in); ok {
		return frozen
	}
	if sellOnly, ok := evaluateBuyLogic(in); ok {
		return sellOnly
	}

	applySentimentPolicy(in, &verdict)
	applyLiquidityPolicy(in, &verdict)
	applyValuationPolicy(in, &verdict)
	applyTakeProfitPolicy(in, &verdict)
	applyCashPolicy(in, &verdict)
	applyPortfolioPolicy(in, &verdict)
	verdict.ExpectedReturns = EstimateExpectedReturns(in)
	return verdict
}

func baseVerdict() model.RuleVerdict {
	return model.RuleVerdict{Status: model.VerdictHold, Text: "按纪律观察"}
}

func addProhibited(v *model.RuleVerdict, action string) {
	v.ProhibitedActions = append(v.ProhibitedActions, action)
	if action == "新增买入" {
		v.OptionalActions = removeBuyLikeOptionalActions(v.OptionalActions)
	}
}
func addOptional(v *model.RuleVerdict, action string) {
	if actionProhibited(v, action) {
		return
	}
	v.OptionalActions = append(v.OptionalActions, action)
}
func addRule(v *model.RuleVerdict, id, name, severity, desc string) {
	v.TriggeredRules = append(v.TriggeredRules, model.TriggeredRule{RuleID: id, RuleName: name, Severity: severity, Description: desc})
}

func actionProhibited(v *model.RuleVerdict, action string) bool {
	for _, prohibited := range v.ProhibitedActions {
		if prohibited == "新增买入" && (action == "分批配置" || action == "按计划定投") {
			return true
		}
	}
	return false
}

func removeBuyLikeOptionalActions(values []string) []string {
	out := values[:0]
	for _, value := range values {
		if value == "分批配置" || value == "按计划定投" {
			continue
		}
		out = append(out, value)
	}
	return out
}
