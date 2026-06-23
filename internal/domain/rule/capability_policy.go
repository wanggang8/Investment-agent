package rule

import "investment-agent/internal/domain/model"

// evaluateCapability 执行能力圈规则；能力圈外直接拒绝交易类分析。
func evaluateCapability(in EvaluationInput) (model.RuleVerdict, bool) {
	if in.CapabilityStatus == "out_of_scope" {
		return model.RuleVerdict{
			Status:            model.VerdictRejected,
			Text:              "标的不在能力圈内，拒绝交易类分析",
			ProhibitedActions: []string{"交易类分析", "新增买入"},
			TriggeredRules:    []model.TriggeredRule{{RuleID: "CAPABILITY", RuleName: "能力圈约束", Severity: "danger", Description: "标的不在能力圈内"}},
		}, true
	}
	return model.RuleVerdict{}, false
}
