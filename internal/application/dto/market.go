package dto

// MarketRefreshRequest 指定需要刷新的市场标的集合。
type MarketRefreshRequest struct {
	Symbols  []string `json:"symbols,omitempty"`
	AsOfDate string   `json:"as_of_date,omitempty"`
}

type MarketRefreshResponse struct {
	RefreshedCount    int                    `json:"refreshed_count"`
	FailedSymbols     []MarketRefreshFailure `json:"failed_symbols"`
	LatestSnapshotIDs []string               `json:"latest_snapshot_ids"`
	AuditEventIDs     []string               `json:"audit_event_ids"`
}

type MarketRefreshFailure struct {
	Symbol string `json:"symbol"`
	Reason string `json:"reason"`
}

type MarketSnapshotDTO struct {
	MarketSnapshotID     string         `json:"market_snapshot_id"`
	Symbol               string         `json:"symbol"`
	TradeDate            string         `json:"trade_date"`
	DataStatus           string         `json:"data_status"`
	ClosePrice           float64        `json:"close_price"`
	TurnoverRate         float64        `json:"turnover_rate"`
	PEPercentile         float64        `json:"pe_percentile"`
	PBPercentile         float64        `json:"pb_percentile"`
	VolumePercentile     float64        `json:"volume_percentile"`
	VolatilityPercentile float64        `json:"volatility_percentile"`
	LiquidityState       string         `json:"liquidity_state"`
	SentimentState       string         `json:"sentiment_state"`
	MarketMetrics        map[string]any `json:"market_metrics"`
}

type SourceHealthResponse struct {
	Sources []SourceHealthItem `json:"sources"`
}

type SourceHealthItem struct {
	SourceName      string   `json:"source_name"`
	SourceLevel     string   `json:"source_level"`
	SourceType      string   `json:"source_type"`
	DataCategory    string   `json:"data_category"`
	Freshness       string   `json:"freshness"`
	DataDate        string   `json:"data_date,omitempty"`
	RequestID       string   `json:"request_id,omitempty"`
	LastSuccessAt   string   `json:"last_success_at,omitempty"`
	LastFailureAt   string   `json:"last_failure_at,omitempty"`
	FailureCategory string   `json:"failure_category,omitempty"`
	AffectedSymbols []string `json:"affected_symbols,omitempty"`
}
