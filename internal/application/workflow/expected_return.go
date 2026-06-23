package workflow

import (
	"encoding/json"
	"strings"

	"investment-agent/internal/domain/model"
)

const nonTradingExpectedReturnDisclaimer = "预期收益仅为情景分析，不构成收益承诺；卖出评估仅提示人工复核，不会自动交易。"

func BuildExpectedReturn(sampleCount int) ExpectedReturnOutput {
	return BuildExpectedReturnWithContext(ExpectedReturnInput{SampleCount: sampleCount})
}

func ExpectedReturnSampleCountFromWorkflowData(positions []model.Position, market model.MarketSnapshot) int {
	count := 0
	for _, position := range positions {
		if position.Symbol == market.Symbol && position.CostPrice > 0 {
			count++
			break
		}
	}
	if market.MarketSnapshotID != "" {
		count++
	}
	var metadata map[string]any
	if err := json.Unmarshal([]byte(market.MarketMetricsJSON), &metadata); err != nil {
		return count
	}
	if nested, ok := metadata["metadata"].(map[string]any); ok {
		if history, ok := nested["nav_history"].([]any); ok {
			count += len(history)
		}
	}
	return count
}

func p34ExpectedReturnContext(market model.MarketSnapshot) (string, []string) {
	var metrics map[string]any
	if err := json.Unmarshal([]byte(market.MarketMetricsJSON), &metrics); err != nil {
		return "P34 扩展数据不可解析", []string{"p34_metrics"}
	}
	metadata, _ := metrics["metadata"].(map[string]any)
	health, _ := metadata["p34_source_health"].(map[string]any)
	if len(health) == 0 {
		return "P34 扩展数据缺失", []string{"p34_expanded_data"}
	}
	fresh := []string{}
	missing := []string{}
	for category, raw := range health {
		status := p34HealthFreshness(raw)
		if status == "fresh" || status == "stubbed" {
			fresh = append(fresh, category)
			continue
		}
		missing = append(missing, category)
	}
	if len(fresh) == 0 {
		return "P34 扩展数据均不可用", missing
	}
	return "P34 可用扩展数据：" + strings.Join(fresh, "、"), missing
}

func p34HealthFreshness(raw any) string {
	switch item := raw.(type) {
	case string:
		return item
	case map[string]any:
		status, _ := item["freshness"].(string)
		return status
	default:
		return "missing"
	}
}

func BuildExpectedReturnWithContext(input ExpectedReturnInput) ExpectedReturnOutput {
	out := ExpectedReturnOutput{
		SampleCount:           input.SampleCount,
		TargetName:            input.TargetName,
		TargetCode:            input.TargetCode,
		HoldingClass:          input.HoldingClass,
		HorizonLabel:          firstNonEmpty(input.HorizonLabel, "未来 12 个月"),
		SampleWindow:          firstNonEmpty(input.SampleWindow, sampleWindow(input.SampleCount)),
		ScreeningCondition:    firstNonEmpty(input.ScreeningCondition, "基于相似估值分位、市场状态和已验证本地样本筛选"),
		SupportingDataSummary: input.SupportingDataSummary,
		MissingCategories:     input.MissingCategories,
		SupplementData:        supplementData(input),
		AssumptionChecks:      input.AssumptionChecks,
		HistoricalContexts:    input.HistoricalContexts,
		HoldingClassCoverage:  input.HoldingClassCoverage,
		SellEvaluation: ExpectedReturnSellEvaluation{
			Status:               "not_triggered",
			NonTradingDisclaimer: nonTradingExpectedReturnDisclaimer,
		},
	}

	switch {
	case input.SampleCount >= 20 && len(input.HistoricalSamples) > 0:
		out.PrecisionStatus = model.PrecisionAvailable
		out.Scenarios = scenariosFromHistoricalSamples(input.HistoricalSamples)
		out.ProbabilityBasis = "historical_similar_sample_proportion"
	case input.SampleCount >= 20:
		up, base, down := 0.25, 0.50, 0.25
		out.PrecisionStatus = model.PrecisionAvailable
		out.Scenarios = []ExpectedReturnScenario{
			{Name: "upside", Probability: &up, ReturnRate: 0.08, ReturnRange: "8.00%~15.00%", LowerBound: 0.08, UpperBound: 0.15, Confidence: "medium", Trigger: "估值修复或情绪改善"},
			{Name: "base", Probability: &base, ReturnRate: 0.03, ReturnRange: "0.00%~8.00%", LowerBound: 0, UpperBound: 0.08, Confidence: "medium", Trigger: "估值维持当前区间"},
			{Name: "downside", Probability: &down, ReturnRate: -0.05, ReturnRange: "-12.00%~0.00%", LowerBound: -0.12, UpperBound: 0, Confidence: "medium", Trigger: "估值收缩或买入逻辑需复核"},
		}
		out.ProbabilityBasis = "deterministic_default_sample_distribution"
	case input.SampleCount >= 5:
		out.PrecisionStatus = model.PrecisionInsufficient
		out.Reason = "样本不足，不能返回精确概率"
		out.Scenarios = []ExpectedReturnScenario{
			{Name: "upside", ReturnRate: 0.08, ReturnRange: "8.00%~15.00%", LowerBound: 0.08, UpperBound: 0.15, Confidence: "low", Trigger: "样本有限，仅作上行情景参考"},
			{Name: "base", ReturnRate: 0.03, ReturnRange: "0.00%~8.00%", LowerBound: 0, UpperBound: 0.08, Confidence: "low", Trigger: "样本有限，仅作基准情景参考"},
			{Name: "downside", ReturnRate: -0.05, ReturnRange: "-12.00%~0.00%", LowerBound: -0.12, UpperBound: 0, Confidence: "low", Trigger: "样本有限，仅作下行情景参考"},
		}
	default:
		out.PrecisionStatus = model.PrecisionUnavailable
		out.Reason = "样本过少，仅能给出定性说明"
	}

	applySellEvaluation(&out, input)
	applyExpectedReturnMonitoring(&out, input)
	return out
}

func scenariosFromHistoricalSamples(samples []ExpectedReturnHistoricalSample) []ExpectedReturnScenario {
	total := 0
	for _, sample := range samples {
		if sample.Count > 0 {
			total += sample.Count
		}
	}
	byScenario := map[string]ExpectedReturnHistoricalSample{}
	for _, sample := range samples {
		byScenario[sample.Scenario] = sample
	}
	out := []ExpectedReturnScenario{}
	for _, name := range []string{"upside", "base", "downside"} {
		sample, ok := byScenario[name]
		if !ok {
			continue
		}
		var probability *float64
		if total > 0 {
			value := float64(sample.Count) / float64(total)
			probability = &value
		}
		out = append(out, ExpectedReturnScenario{Name: name, Probability: probability, ReturnRate: sample.ReturnRate, ReturnRange: sample.ReturnRange, LowerBound: sample.LowerBound, UpperBound: sample.UpperBound, Confidence: "historical", Trigger: sample.Trigger})
	}
	return out
}

func supplementData(input ExpectedReturnInput) []string {
	if len(input.MissingCategories) > 0 {
		return append([]string{}, input.MissingCategories...)
	}
	if input.SampleCount < 5 {
		return []string{"market_history", "valuation_percentiles", "fundamental_growth", "formal_evidence"}
	}
	return nil
}

func sampleWindow(sampleCount int) string {
	if sampleCount <= 0 {
		return "样本窗口不可用：样本数量为 0"
	}
	return "最近可比样本"
}

func applySellEvaluation(out *ExpectedReturnOutput, input ExpectedReturnInput) {
	if len(out.Scenarios) == 0 {
		out.SellEvaluation = ExpectedReturnSellEvaluation{Status: "not_applicable", Prompts: []string{"样本过少，无法生成可复核的情景边界"}, NonTradingDisclaimer: nonTradingExpectedReturnDisclaimer}
		return
	}
	if input.CurrentPrice <= 0 {
		out.SellEvaluation = ExpectedReturnSellEvaluation{Status: "not_applicable", Prompts: []string{"缺少当前价格，无法评估动态卖出边界"}, NonTradingDisclaimer: nonTradingExpectedReturnDisclaimer}
		return
	}
	if input.BasePrice <= 0 {
		out.SellEvaluation = ExpectedReturnSellEvaluation{Status: "not_applicable", Prompts: []string{"缺少持仓成本或可复现基准价格，无法评估动态卖出边界"}, NonTradingDisclaimer: nonTradingExpectedReturnDisclaimer}
		return
	}
	currentReturn := input.CurrentPrice/input.BasePrice - 1
	triggers := []string{}
	prompts := []string{}
	actions := []string{}

	if currentReturn >= out.Scenarios[0].LowerBound {
		triggers = append(triggers, "upside_lower_bound_reached")
		prompts = append(prompts, "当前价格已进入乐观情景下沿，请评估是否启动移动止盈")
		actions = append(actions, "评估移动止盈")
	}
	if currentReturn >= out.Scenarios[1].UpperBound {
		triggers = append(triggers, "base_upper_bound_exceeded")
		prompts = append(prompts, "当前价格已突破基准情景上沿，请评估是否分批止盈")
		actions = append(actions, "评估分批止盈")
	}
	if currentReturn <= out.Scenarios[2].LowerBound {
		triggers = append(triggers, "downside_lower_bound_breached")
		prompts = append(prompts, "当前价格跌破悲观情景下沿，请重新核验买入逻辑")
		actions = append(actions, "复核买入逻辑")
	}
	baseMidpoint := out.Scenarios[1].ReturnRate
	if input.PreviousBaseMidpoint > 0 && input.PreviousBaseMidpoint-baseMidpoint > 0.15 {
		triggers = append(triggers, "base_midpoint_downshift")
		prompts = append(prompts, "基准情景中枢下移超过 15%，请重新评估买入逻辑并考虑减仓")
		actions = append(actions, "复核基准情景")
		out.ReassessmentTrigger = ExpectedReturnReassessmentTrigger{Reason: "基准情景中枢下移超过 15%", Boundary: "base_midpoint_downshift", CurrentValue: baseMidpoint}
	}
	if input.TargetReturnRate > 0 && currentReturn >= input.TargetReturnRate {
		triggers = append(triggers, "target_return_reached")
		prompts = append(prompts, "当前收益达到用户目标，请查看或记录人工计划")
		actions = append(actions, "查看人工计划")
	}

	if len(triggers) == 0 {
		return
	}
	out.SellEvaluation = ExpectedReturnSellEvaluation{Status: "triggered", Triggers: triggers, Prompts: prompts, Actions: actions, NonTradingDisclaimer: nonTradingExpectedReturnDisclaimer}
	if out.ReassessmentTrigger.Reason == "" {
		out.ReassessmentTrigger = ExpectedReturnReassessmentTrigger{Reason: prompts[0], Boundary: triggers[0], CurrentValue: currentReturn}
	}
}

func applyExpectedReturnMonitoring(out *ExpectedReturnOutput, input ExpectedReturnInput) {
	triggers := append([]string{}, out.SellEvaluation.Triggers...)
	prompts := append([]string{}, out.SellEvaluation.Prompts...)
	actions := append([]string{}, out.SellEvaluation.Actions...)
	add := func(trigger, prompt, action string) {
		if containsExpectedReturnString(triggers, trigger) {
			return
		}
		triggers = append(triggers, trigger)
		prompts = append(prompts, prompt)
		actions = append(actions, action)
	}
	if input.MarketState == "stress" || input.FundamentalState == "below_expectation" {
		applyScenarioProbabilityDownshift(out)
		add("scenario_probability_downshift", "估值、基本面或市场状态转弱，需下调相关情景概率并人工复核", "复核情景概率")
		if out.ReassessmentTrigger.Reason == "" {
			out.ReassessmentTrigger = ExpectedReturnReassessmentTrigger{Reason: "情景概率下修", Boundary: "scenario_probability_downshift"}
		}
	}
	for _, check := range input.AssumptionChecks {
		if check.MonthsBelow >= 2 && check.Actual < check.Expected {
			add("two_month_assumption_downshift", check.Name+" 连续 2 个月低于预期，触发情景下修预警", "复核核心假设")
		}
	}
	if input.PessimisticPathMonths >= 1 {
		add("one_month_pessimistic_path", "实际走势连续 1 个月偏向悲观情景，建议用户手动调整情景概率", "手动调整情景概率")
	}
	if input.SentimentState == "extreme" && len(input.HistoricalContexts) > 0 {
		add("extreme_fear_historical_context", "极端恐惧状态已展示历史相似场景，暂停主动交易建议并等待正式证据复核", "暂停主动交易建议")
	}
	if len(triggers) == 0 {
		return
	}
	out.SellEvaluation = ExpectedReturnSellEvaluation{Status: "triggered", Triggers: triggers, Prompts: prompts, Actions: actions, NonTradingDisclaimer: nonTradingExpectedReturnDisclaimer}
}

func applyScenarioProbabilityDownshift(out *ExpectedReturnOutput) {
	delta := 0.0
	downsideIndex := -1
	for i := range out.Scenarios {
		scenario := &out.Scenarios[i]
		if scenario.Probability == nil {
			continue
		}
		switch scenario.Name {
		case "upside":
			reduction := *scenario.Probability * 0.25
			*scenario.Probability -= reduction
			delta += reduction
			scenario.Confidence = firstNonEmpty(scenario.Confidence, "adjusted")
		case "base":
			reduction := *scenario.Probability * 0.10
			*scenario.Probability -= reduction
			delta += reduction
			scenario.Confidence = firstNonEmpty(scenario.Confidence, "adjusted")
		case "downside":
			downsideIndex = i
		}
	}
	if delta > 0 && downsideIndex >= 0 && out.Scenarios[downsideIndex].Probability != nil {
		*out.Scenarios[downsideIndex].Probability += delta
		out.Scenarios[downsideIndex].Confidence = firstNonEmpty(out.Scenarios[downsideIndex].Confidence, "adjusted")
	}
}

func containsExpectedReturnString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
