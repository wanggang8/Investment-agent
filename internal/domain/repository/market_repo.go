package repository

import (
	"context"

	"investment-agent/internal/domain/model"
)

// MarketRepository 定义市场快照的持久化边界。
type MarketRepository interface {
	SaveMarketSnapshot(ctx context.Context, snapshot model.MarketSnapshot, createdAt string) error
	GetMarketSnapshot(ctx context.Context, snapshotID string) (model.MarketSnapshot, error)
	GetLatestMarketSnapshot(ctx context.Context) (model.MarketSnapshot, error)
	GetLatestMarketSnapshotBySymbol(ctx context.Context, symbol string) (model.MarketSnapshot, error)
	MarketSnapshotExists(ctx context.Context, snapshotID string) (bool, error)
}
