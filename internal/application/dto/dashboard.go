package dto

// DashboardTodayResponse 对齐 `GET /api/v1/dashboard/today` 的 data 字段。
// 字段直接服务驾驶舱首屏展示，不暴露领域内部结构。
type DashboardTodayResponse struct {
	DashboardState   string             `json:"dashboard_state"`
	DisciplineStatus string             `json:"discipline_status"`
	DataUpdatedAt    string             `json:"data_updated_at"`
	PortfolioSummary PortfolioSummary   `json:"portfolio_summary"`
	MarketSummary    MarketSummary      `json:"market_summary"`
	TriggeredRules   []TriggeredRuleDTO `json:"triggered_rules"`
	DecisionSummary  DecisionSummary    `json:"decision_summary"`
	EvidenceSummary  *EvidenceSummary   `json:"evidence_summary,omitempty"`
}

type PortfolioSummary struct {
	TotalAssets   float64 `json:"total_assets"`
	CashRatio     float64 `json:"cash_ratio"`
	HighRiskRatio float64 `json:"high_risk_ratio"`
	PositionCount int     `json:"position_count"`
}

type MarketSummary struct {
	PEPercentile   float64 `json:"pe_percentile"`
	PBPercentile   float64 `json:"pb_percentile"`
	SentimentState string  `json:"sentiment_state"`
	LiquidityState string  `json:"liquidity_state"`
}

type TriggeredRuleDTO struct {
	RuleID      string `json:"rule_id"`
	RuleName    string `json:"rule_name"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
}

type DecisionSummary struct {
	DecisionID         string   `json:"decision_id"`
	Verdict            string   `json:"verdict"`
	FinalVerdictStatus string   `json:"final_verdict_status"`
	ProhibitedActions  []string `json:"prohibited_actions"`
	OptionalActions    []string `json:"optional_actions"`
	ActionRequired     bool     `json:"action_required"`
	ConfirmationStatus string   `json:"confirmation_status"`
}

type EvidenceSummary struct {
	SourceCount        int    `json:"source_count"`
	HighestSourceLevel string `json:"highest_source_level"`
	VerificationStatus string `json:"verification_status"`
}
