package model

// MarketSnapshot 是领域规则读取的市场状态视图。
type MarketSnapshot struct {
	MarketSnapshotID     string
	Symbol               string
	TradeDate            string
	DataStatus           string
	MarketMetricsJSON    string
	ClosePrice           float64
	TurnoverRate         float64
	MarginBalance        float64
	MarginBalanceChange  float64
	PEPercentile         float64
	PBPercentile         float64
	LiquidityState       LiquidityState
	SentimentState       SentimentState
	VolumePercentile     float64
	VolatilityPercentile float64
}
