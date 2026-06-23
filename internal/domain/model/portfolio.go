package model

// PortfolioSnapshot 是规则裁决读取的账户总览，不直接代表当前可交易状态。
type PortfolioSnapshot struct {
	SnapshotID    string
	Cash          float64
	TotalAssets   float64
	CashRatio     float64
	HighRiskRatio float64
	PositionCount int
}

// Position 是规则裁决使用的持仓视图，包含止盈、状态和买入逻辑破坏标记。
type Position struct {
	PositionID            string
	Symbol                string
	Name                  string
	Quantity              float64
	CostPrice             float64
	CurrentPrice          float64
	MarketValue           float64
	UnrealizedProfitRatio float64
	PositionState         PositionState
	AssetTag              string
	BuyLogicBroken        bool
	TakeProfitStarted     bool
	StageHighPrice        float64
	HandledProfit20       bool
	HandledProfit30       bool
}

// PortfolioAllocation 是核心-卫星仓位比例输入。
type PortfolioAllocation struct {
	CoreRatio      float64
	SatelliteRatio float64
}
