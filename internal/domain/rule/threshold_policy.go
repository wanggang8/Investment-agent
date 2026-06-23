package rule

import (
	"strings"

	"investment-agent/internal/domain/model"
)

// SentimentSignal captures the non-LLM inputs allowed to trigger the cooldown rule.
type SentimentSignal struct {
	SentimentPercentile   float64
	AbnormalFinancingDays int
	AbnormalVolumeDays    int
	AbnormalMediaHeatDays int
	UserEmotionTags       []string
}

// LiquiditySignal captures sizing facts for market-style liquidity checks.
type LiquiditySignal struct {
	PlanAmount         float64
	Average20DayAmount float64
	SameDayAmount      float64
	SpreadAbnormal     bool
}

// DetermineSentimentState implements the executable sentiment thresholds from L1.
func DetermineSentimentState(signal SentimentSignal) model.SentimentState {
	if signal.SentimentPercentile >= 90 || (signal.SentimentPercentile > 0 && signal.SentimentPercentile <= 10) {
		return model.SentimentExtreme
	}
	if signal.AbnormalFinancingDays >= 3 || signal.AbnormalVolumeDays >= 3 || signal.AbnormalMediaHeatDays >= 3 {
		return model.SentimentExtreme
	}
	for _, tag := range signal.UserEmotionTags {
		switch strings.TrimSpace(tag) {
		case "恐慌", "狂喜", "清仓", "满仓", "追热点":
			return model.SentimentExtreme
		}
	}
	return model.SentimentNeutral
}

// DetermineLiquidityState implements the 20x 20-day amount and same-day 5% checks.
func DetermineLiquidityState(signal LiquiditySignal) model.LiquidityState {
	if signal.SpreadAbnormal {
		return model.LiquidityDanger
	}
	if signal.PlanAmount > 0 && signal.Average20DayAmount > 0 && signal.Average20DayAmount < signal.PlanAmount*20 {
		return model.LiquidityDanger
	}
	if signal.PlanAmount > 0 && signal.SameDayAmount > 0 && signal.PlanAmount > signal.SameDayAmount*0.05 {
		return model.LiquidityDanger
	}
	return model.LiquidityNormal
}

// CooldownTradingDays extends cooldown after three consecutive circuit-break triggers.
func CooldownTradingDays(consecutiveTriggers int) int {
	if consecutiveTriggers >= 3 {
		return 5
	}
	return 1
}

// PositionStateForVerdict maps a rule verdict to the position discipline state.
func PositionStateForVerdict(status model.FinalVerdictStatus) model.PositionState {
	switch status {
	case model.VerdictSellOnly:
		return model.PositionSellOnly
	case model.VerdictFrozenWatch:
		return model.PositionFrozenWatch
	default:
		return model.PositionNormal
	}
}
