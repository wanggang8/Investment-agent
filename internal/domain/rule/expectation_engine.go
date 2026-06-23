package rule

import "investment-agent/internal/domain/model"

// EstimateExpectedReturns 输出上行、基准、下行情景。
// 该结果只进入分析材料，不能覆盖 Evaluate 已给出的最终裁决。
func EstimateExpectedReturns(in EvaluationInput) []model.ExpectedReturnScenario {
	base := 0.03
	if in.PEPercentile < 30 && in.PBPercentile < 30 {
		base = 0.08
	}
	if in.PEPercentile >= 80 || in.PBPercentile >= 80 {
		base = -0.02
	}
	return []model.ExpectedReturnScenario{
		{Name: "upside", Probability: 0.25, ReturnRate: base + 0.08, Confidence: "medium"},
		{Name: "base", Probability: 0.50, ReturnRate: base, Confidence: "medium"},
		{Name: "downside", Probability: 0.25, ReturnRate: base - 0.08, Confidence: "medium"},
	}
}
