package repository

import "context"

// PortfolioSnapshot 是一次决策读取到的账户总览快照。
type PortfolioSnapshot struct {
	SnapshotID    string
	SnapshotTime  string
	Cash          float64
	TotalAssets   float64
	CashRatio     float64
	HighRiskRatio float64
	PositionCount int
	Source        string
	CreatedAt     string
}

// Position 是当前持仓聚合态，会随用户线下确认动作更新。
type Position struct {
	PositionID            string
	Symbol                string
	Name                  string
	Quantity              float64
	CostPrice             float64
	CurrentPrice          float64
	MarketValue           float64
	UnrealizedProfitRatio float64
	PositionState         string
	BuyDate               string
	BuyReason             string
	AssetTag              string
	UpdatedAt             string
}

// PositionSnapshot 是账户快照时点下的持仓明细，用于复现历史裁决上下文。
type PositionSnapshot struct {
	PositionSnapshotID    string
	SnapshotID            string
	PositionID            string
	Symbol                string
	Name                  string
	Quantity              float64
	CostPrice             float64
	CurrentPrice          float64
	MarketValue           float64
	UnrealizedProfitRatio float64
	PositionState         string
	BuyDate               string
	BuyReason             string
	AssetTag              string
	CreatedAt             string
}

type LocalAccountImportBatch struct {
	ImportBatchID         string
	RequestID             string
	Status                string
	RowCount              int
	ValidCount            int
	InvalidCount          int
	ValidationSummaryJSON string
	RowsHash              string
	CreatedAt             string
	CommittedAt           string
}

type LocalAccountCorrection struct {
	CorrectionID     string
	TargetType       string
	TargetID         string
	BeforeJSON       string
	AfterJSON        string
	CorrectionReason string
	SnapshotID       string
	AuditEventID     string
	CreatedAt        string
}

// PortfolioRepository 定义账户快照与持仓当前态的持久化边界。
type PortfolioRepository interface {
	SavePortfolioSnapshot(ctx context.Context, snapshot PortfolioSnapshot, positions []PositionSnapshot) error
	GetPortfolioSnapshot(ctx context.Context, snapshotID string) (PortfolioSnapshot, []PositionSnapshot, error)
	GetLatestPortfolioSnapshot(ctx context.Context) (PortfolioSnapshot, error)
	SavePosition(ctx context.Context, position Position) error
	ReplacePositions(ctx context.Context, positions []Position) error
	DeletePosition(ctx context.Context, positionID string) error
	GetPosition(ctx context.Context, positionID string) (Position, error)
	ListPositions(ctx context.Context) ([]Position, error)
	SaveLocalAccountImportBatch(ctx context.Context, batch LocalAccountImportBatch) error
	GetLocalAccountImportBatch(ctx context.Context, importBatchID string) (LocalAccountImportBatch, error)
	SaveLocalAccountCorrection(ctx context.Context, correction LocalAccountCorrection) error
}
