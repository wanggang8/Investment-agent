export interface MarketRefreshRequest {
  symbols?: string[]
}

export interface MarketRefreshResponse {
  refreshed_count: number
  failed_symbols: MarketRefreshFailure[]
  latest_snapshot_ids: string[]
  audit_event_ids: string[]
}

export interface MarketRefreshFailure {
  symbol: string
  reason: string
}

export interface MarketSnapshot {
  market_snapshot_id: string
  symbol: string
  trade_date?: string
  market_metrics?: Record<string, unknown>
  close_price?: number
  turnover_rate?: number
  pe_percentile: number
  pb_percentile: number
  volume_percentile?: number
  volatility_percentile?: number
  liquidity_state: string
  sentiment_state: string
  data_status?: string
}

export interface SourceHealthResponse {
  sources: SourceHealthItem[]
}

export interface SourceHealthItem {
  source_name: string
  source_level: string
  source_type: string
  data_category: string
  freshness: string
  data_date?: string
  request_id?: string
  last_success_at?: string
  last_failure_at?: string
  failure_category?: string
  affected_symbols?: string[]
}
