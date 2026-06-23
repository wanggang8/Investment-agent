package dto

// PortfolioInitRequest 是首次录入本地账户事实的请求。
type PortfolioInitRequest struct {
	Cash        float64         `json:"cash"`
	TotalAssets float64         `json:"total_assets"`
	Positions   []PositionInput `json:"positions"`
}

type PortfolioAdjustmentRequest struct {
	Cash         float64         `json:"cash"`
	TotalAssets  float64         `json:"total_assets"`
	AdjustReason string          `json:"adjust_reason"`
	Positions    []PositionInput `json:"positions"`
}

type PositionInput struct {
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	Quantity      float64 `json:"quantity"`
	CostPrice     float64 `json:"cost_price"`
	CurrentPrice  float64 `json:"current_price"`
	BuyDate       string  `json:"buy_date,omitempty"`
	PositionState string  `json:"position_state,omitempty"`
	BuyReason     string  `json:"buy_reason,omitempty"`
	AssetTag      string  `json:"asset_tag,omitempty"`
}

type HoldingEditRequest struct {
	PositionID   string        `json:"position_id,omitempty"`
	Reason       string        `json:"reason"`
	Confirmation string        `json:"confirmation"`
	Position     PositionInput `json:"position"`
}

type HoldingRemoveRequest struct {
	PositionID   string `json:"position_id"`
	Reason       string `json:"reason"`
	Confirmation string `json:"confirmation"`
}

type OfflineTransactionRequest struct {
	OperationType string  `json:"operation_type"`
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name,omitempty"`
	Quantity      float64 `json:"quantity"`
	Price         float64 `json:"price"`
	Fees          float64 `json:"fees,omitempty"`
	ExecutedAt    string  `json:"executed_at"`
	Note          string  `json:"note,omitempty"`
	BuyReason     string  `json:"buy_reason,omitempty"`
	AssetTag      string  `json:"asset_tag,omitempty"`
}

type BatchImportValidationRequest struct {
	Rows []BatchImportRow `json:"rows"`
}

type BatchImportRow struct {
	RowNumber     int     `json:"row_number"`
	RowType       string  `json:"row_type"`
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name,omitempty"`
	Quantity      float64 `json:"quantity"`
	CostPrice     float64 `json:"cost_price,omitempty"`
	CurrentPrice  float64 `json:"current_price,omitempty"`
	OperationType string  `json:"operation_type,omitempty"`
	Price         float64 `json:"price,omitempty"`
	Fees          float64 `json:"fees,omitempty"`
	OccurredAt    string  `json:"occurred_at,omitempty"`
	BuyDate       string  `json:"buy_date,omitempty"`
	BuyReason     string  `json:"buy_reason,omitempty"`
	PositionState string  `json:"position_state,omitempty"`
	AssetTag      string  `json:"asset_tag,omitempty"`
}

type BatchImportValidationResponse struct {
	ImportBatchID string                       `json:"import_batch_id"`
	Summary       BatchImportValidationSummary `json:"summary"`
	Rows          []BatchImportRowResult       `json:"rows"`
}

type BatchImportValidationSummary struct {
	RowCount     int `json:"row_count"`
	ValidCount   int `json:"valid_count"`
	InvalidCount int `json:"invalid_count"`
}

type BatchImportRowResult struct {
	RowNumber int      `json:"row_number"`
	Valid     bool     `json:"valid"`
	Errors    []string `json:"errors,omitempty"`
}

type BatchImportConfirmRequest struct {
	ImportBatchID string           `json:"import_batch_id"`
	Rows          []BatchImportRow `json:"rows"`
	ConfirmReason string           `json:"confirm_reason"`
}

type CorrectionRequest struct {
	TargetType       string `json:"target_type"`
	TargetID         string `json:"target_id"`
	BeforeJSON       string `json:"before_json"`
	AfterJSON        string `json:"after_json"`
	CorrectionReason string `json:"correction_reason"`
}

type RebalanceReviewRequest struct {
	TargetCoreRatio      float64 `json:"target_core_ratio"`
	TargetSatelliteRatio float64 `json:"target_satellite_ratio"`
	TargetCashRatio      float64 `json:"target_cash_ratio"`
	DriftThreshold       float64 `json:"drift_threshold,omitempty"`
	ReviewDate           string  `json:"review_date,omitempty"`
}

type RebalanceReviewResponse struct {
	ReviewID        string                `json:"review_id"`
	ReviewDate      string                `json:"review_date"`
	TotalAssets     float64               `json:"total_assets"`
	DriftThreshold  float64               `json:"drift_threshold"`
	Items           []RebalanceReviewItem `json:"items"`
	AuditEventIDs   []string              `json:"audit_event_ids"`
	SafetyStatement string                `json:"safety_statement"`
}

type RebalanceReviewItem struct {
	Bucket         string  `json:"bucket"`
	TargetRatio    float64 `json:"target_ratio"`
	ActualRatio    float64 `json:"actual_ratio"`
	DriftRatio     float64 `json:"drift_ratio"`
	TargetValue    float64 `json:"target_value"`
	ActualValue    float64 `json:"actual_value"`
	Recommendation string  `json:"recommendation"`
	ManualAmount   float64 `json:"manual_amount"`
}

type LocalFactWriteResponse struct {
	SnapshotID      string   `json:"snapshot_id,omitempty"`
	PositionID      string   `json:"position_id,omitempty"`
	TransactionID   string   `json:"transaction_id,omitempty"`
	ImportBatchID   string   `json:"import_batch_id,omitempty"`
	CorrectionID    string   `json:"correction_id,omitempty"`
	AuditEventIDs   []string `json:"audit_event_ids"`
	SafetyStatement string   `json:"safety_statement"`
}

type PortfolioWriteResponse struct {
	SnapshotID            string   `json:"snapshot_id"`
	PositionCount         int      `json:"position_count"`
	PositionSnapshotCount int      `json:"position_snapshot_count,omitempty"`
	AuditEventIDs         []string `json:"audit_event_ids"`
}

type PortfolioCurrentResponse struct {
	Snapshot  PortfolioSnapshotDTO `json:"snapshot"`
	Positions []PositionDTO        `json:"positions"`
}

type PortfolioSnapshotDTO struct {
	SnapshotID    string  `json:"snapshot_id"`
	SnapshotTime  string  `json:"snapshot_time"`
	Cash          float64 `json:"cash"`
	TotalAssets   float64 `json:"total_assets"`
	CashRatio     float64 `json:"cash_ratio"`
	HighRiskRatio float64 `json:"high_risk_ratio"`
	PositionCount int     `json:"position_count"`
}

type PositionDTO struct {
	PositionID            string  `json:"position_id"`
	Symbol                string  `json:"symbol"`
	Name                  string  `json:"name"`
	Quantity              float64 `json:"quantity"`
	CostPrice             float64 `json:"cost_price"`
	CurrentPrice          float64 `json:"current_price"`
	MarketValue           float64 `json:"market_value"`
	UnrealizedProfitRatio float64 `json:"unrealized_profit_ratio"`
	PositionState         string  `json:"position_state"`
	BuyDate               string  `json:"buy_date,omitempty"`
	BuyReason             string  `json:"buy_reason,omitempty"`
	AssetTag              string  `json:"asset_tag,omitempty"`
}
