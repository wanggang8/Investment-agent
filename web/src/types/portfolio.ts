export interface PortfolioCurrentResponse {
  snapshot: PortfolioSnapshot
  positions: Position[]
}

export interface PortfolioSnapshot {
  snapshot_id: string
  snapshot_time: string
  cash: number
  total_assets: number
  cash_ratio: number
  high_risk_ratio: number
  position_count: number
}

export interface Position {
  position_id: string
  symbol: string
  name: string
  quantity: number
  cost_price: number
  current_price: number
  market_value: number
  unrealized_profit_ratio: number
  position_state: 'normal' | 'sell_only' | 'frozen_watch' | string
  buy_date?: string
  buy_reason?: string
  asset_tag?: string
}

export interface PositionInput {
  symbol: string
  name: string
  quantity: number
  cost_price: number
  current_price: number
  buy_date?: string
  position_state?: 'normal' | 'sell_only' | 'frozen_watch' | string
  buy_reason?: string
  asset_tag?: string
}

export interface PortfolioInitRequest {
  cash: number
  total_assets: number
  positions: PositionInput[]
}

export interface PortfolioAdjustmentRequest extends PortfolioInitRequest {
  adjust_reason: string
}

export interface PortfolioWriteResponse {
  snapshot_id: string
  position_count: number
  position_snapshot_count?: number
  audit_event_ids: string[]
}

export interface LocalFactWriteResponse {
  snapshot_id?: string
  position_id?: string
  transaction_id?: string
  import_batch_id?: string
  correction_id?: string
  audit_event_ids: string[]
  safety_statement: string
}

export interface HoldingEditRequest {
  position_id?: string
  reason: string
  confirmation: string
  position: PositionInput
}

export interface HoldingRemoveRequest {
  position_id: string
  reason: string
  confirmation: string
}

export interface OfflineTransactionRequest {
  operation_type: 'buy' | 'sell' | 'reduce' | string
  symbol: string
  name?: string
  quantity: number
  price: number
  fees?: number
  executed_at: string
  note?: string
  buy_reason?: string
  asset_tag?: string
}

export interface BatchImportRow {
  row_number?: number
  row_type: 'holding' | 'transaction' | string
  symbol: string
  name?: string
  quantity: number
  cost_price?: number
  current_price?: number
  operation_type?: 'buy' | 'sell' | 'reduce' | string
  price?: number
  fees?: number
  occurred_at?: string
  buy_date?: string
  buy_reason?: string
  position_state?: 'normal' | 'sell_only' | 'frozen_watch' | string
  asset_tag?: string
}

export interface BatchImportValidationRequest {
  rows: BatchImportRow[]
}

export interface BatchImportValidationSummary {
  row_count: number
  valid_count: number
  invalid_count: number
}

export interface BatchImportRowResult {
  row_number: number
  valid: boolean
  errors?: string[]
}

export interface BatchImportValidationResponse {
  import_batch_id: string
  summary: BatchImportValidationSummary
  rows: BatchImportRowResult[]
}

export interface BatchImportConfirmRequest extends BatchImportValidationRequest {
  import_batch_id: string
  confirm_reason: string
}

export interface CorrectionRequest {
  target_type: 'portfolio_snapshot' | 'position' | 'position_snapshot' | 'position_transaction' | 'import_batch' | string
  target_id: string
  before_json: string
  after_json: string
  correction_reason: string
}

export interface RebalanceReviewRequest {
  target_core_ratio: number
  target_satellite_ratio: number
  target_cash_ratio: number
  drift_threshold?: number
  review_date?: string
}

export interface RebalanceReviewItem {
  bucket: 'core' | 'satellite' | 'cash' | string
  target_ratio: number
  actual_ratio: number
  drift_ratio: number
  target_value: number
  actual_value: number
  recommendation: string
  manual_amount: number
}

export interface RebalanceReviewResponse {
  review_id: string
  review_date: string
  total_assets: number
  drift_threshold: number
  items: RebalanceReviewItem[]
  audit_event_ids: string[]
  safety_statement: string
}
