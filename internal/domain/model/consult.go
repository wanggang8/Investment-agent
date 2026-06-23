package model

// ConsultScenario 表示主动咨询前端场景。
type ConsultScenario string

const (
	ConsultScenarioHoldReview      ConsultScenario = "hold_review"
	ConsultScenarioBuyReview       ConsultScenario = "buy_review"
	ConsultScenarioSellReview      ConsultScenario = "sell_review"
	ConsultScenarioRebalanceReview ConsultScenario = "rebalance_review"
)

// Valid 判断主动咨询场景是否在契约枚举内。
func (s ConsultScenario) Valid() bool {
	switch s {
	case ConsultScenarioHoldReview, ConsultScenarioBuyReview, ConsultScenarioSellReview, ConsultScenarioRebalanceReview:
		return true
	default:
		return false
	}
}
