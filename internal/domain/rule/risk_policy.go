package rule

import (
	"strconv"

	"investment-agent/internal/domain/model"
)

// evaluateBuyLogic 处理买入逻辑破坏后的只卖不买规则。
func evaluateBuyLogic(in EvaluationInput) (model.RuleVerdict, bool) {
	if verified, count := sourceVerifiedBuyLogicBreak(in.Evidence); verified {
		text := "买入逻辑破坏，只卖不买；A/S 独立信源=" + strconv.Itoa(count)
		return model.RuleVerdict{Status: model.VerdictSellOnly, Text: text, ProhibitedActions: []string{"新增买入", "加仓"}, TriggeredRules: []model.TriggeredRule{{RuleID: "BUY_LOGIC", RuleName: "买入逻辑破坏", Severity: "danger", Description: "禁止新增买入和加仓；A/S 独立信源=" + strconv.Itoa(count)}}}, true
	}
	if in.BuyLogicBroken {
		return model.RuleVerdict{Status: model.VerdictSellOnly, Text: "买入逻辑破坏，只卖不买", ProhibitedActions: []string{"新增买入", "加仓"}, TriggeredRules: []model.TriggeredRule{{RuleID: "BUY_LOGIC", RuleName: "买入逻辑破坏", Severity: "danger", Description: "禁止新增买入和加仓"}}}, true
	}
	return model.RuleVerdict{}, false
}

func sourceVerifiedBuyLogicBreak(evidence []model.Evidence) (bool, int) {
	for _, ev := range evidence {
		if ev.Role == model.EvidenceFormal && ev.EventType == model.EventBuyLogicBreak && ev.HighGradeIndependentSourceCount >= 2 {
			return true, ev.HighGradeIndependentSourceCount
		}
	}
	return false, 0
}

// applySentimentPolicy 在市场或用户情绪极端时暂停主动交易建议。
func applySentimentPolicy(in EvaluationInput, v *model.RuleVerdict) {
	if in.SentimentState == model.SentimentExtreme {
		addProhibited(v, "主动交易建议")
		addRule(v, "SENTIMENT", "极端情绪", "warning", "情绪极端时暂停主动交易建议")
	}
}

// applyLiquidityPolicy 在市场流动性不足时限制新增买入和大额操作。
func applyLiquidityPolicy(in EvaluationInput, v *model.RuleVerdict) {
	switch in.LiquidityState {
	case model.LiquidityDanger:
		v.Status = model.VerdictHold
		addProhibited(v, "新增买入")
		addProhibited(v, "大额市价操作")
		addRule(v, "LIQUIDITY", "流动性高危", "danger", "流动性不足时禁止新增买入和大额市价操作")
	case model.LiquidityWarning:
		addProhibited(v, "大额市价操作")
		addOptional(v, "分批或限价处理")
		addRule(v, "LIQUIDITY", "流动性观察", "warning", "流动性偏弱时仅允许分批或限价处理")
	}
}

// applyValuationPolicy 使用 PE/PB 分位映射估值区间。
func applyValuationPolicy(in EvaluationInput, v *model.RuleVerdict) {
	p := max(in.PEPercentile, in.PBPercentile)
	if p <= 0 {
		return
	}
	switch {
	case p >= 80:
		v.Status = model.VerdictHold
		addProhibited(v, "新增买入")
		addRule(v, "VALUATION", "估值高危", "danger", "PE/PB 分位超过 80%")
	case p >= 50:
		v.Status = model.VerdictHold
		addRule(v, "VALUATION", "估值观察", "warning", "PE/PB 分位位于 50%-80%")
	case p >= 30:
		if len(v.ProhibitedActions) == 0 {
			v.Status = model.VerdictBuyAllowed
		}
		addOptional(v, "按计划定投")
		addRule(v, "VALUATION", "估值舒适", "normal", "PE/PB 分位位于 30%-50%")
	default:
		if len(v.ProhibitedActions) == 0 {
			v.Status = model.VerdictBuyAllowed
		}
		addOptional(v, "分批配置")
		addRule(v, "VALUATION", "低估区", "normal", "PE/PB 分位低于 30%")
	}
}

// applyTakeProfitPolicy 处理浮盈 20%、30% 与移动止盈回撤规则。
func applyTakeProfitPolicy(in EvaluationInput, v *model.RuleVerdict) {
	switch {
	case in.TakeProfitStarted && in.StageHighPrice > 0 && in.CurrentPrice <= in.StageHighPrice*0.9:
		v.Status = model.VerdictReduce
		addOptional(v, "回撤 10% 后减仓或卖出评估")
		addRule(v, "TAKE_PROFIT", "移动止盈", "warning", "阶段高点回撤 10%")
	case in.UnrealizedProfitRatio >= 0.30 && in.HandledProfit20 && !in.HandledProfit30:
		addOptional(v, "再卖出 30%")
		addRule(v, "TAKE_PROFIT", "浮盈 30%", "normal", "分批止盈第二阶段")
	case in.UnrealizedProfitRatio >= 0.20 && !in.HandledProfit20:
		addOptional(v, "卖出 30%")
		addOptional(v, "启动移动止盈")
		addRule(v, "TAKE_PROFIT", "浮盈 20%", "normal", "分批止盈第一阶段")
	}
}

// applyCashPolicy 执行 R-5 现金冗余规则。
func applyCashPolicy(in EvaluationInput, v *model.RuleVerdict) {
	if in.CashRatio > 0 && in.CashRatio < 0.05 {
		v.Status = model.VerdictHold
		addProhibited(v, "新增买入")
		addRule(v, "R-5", "现金冗余", "danger", "现金比例低于 5%")
	}
}

// applyPortfolioPolicy 检查核心-卫星仓位是否偏离目标区间。
func applyPortfolioPolicy(in EvaluationInput, v *model.RuleVerdict) {
	if in.SatelliteRatio > 0 && (in.SatelliteRatio > 0.30 || abs(in.SatelliteRatio-0.25) > 0.15) {
		addOptional(v, "再平衡卫星资产")
		if in.UnrealizedProfitRatio >= 0.20 || in.TakeProfitStarted {
			addOptional(v, "止盈资金优先回归核心资产")
		}
		addRule(v, "ALLOCATION", "核心-卫星仓位", "warning", "卫星仓位超出目标区间")
	}
	if in.CoreRatio > 0 && in.CoreRatio < 0.60 {
		addOptional(v, "提高核心资产占比")
		addRule(v, "ALLOCATION", "核心资产不足", "warning", "核心资产低于目标 60%")
	}
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
