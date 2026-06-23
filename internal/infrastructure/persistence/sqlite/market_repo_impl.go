package sqlite

import (
	"context"
	"database/sql"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/pkg/apperr"
)

// MarketRepository 是市场快照表的 SQLite 实现。
type MarketRepository struct{ db dbtx }

// NewMarketRepository 创建市场仓储实例。
func NewMarketRepository(db *sql.DB) *MarketRepository { return &MarketRepository{db: db} }

// SaveMarketSnapshot 保存行情、估值、流动性和情绪快照。
func (r *MarketRepository) SaveMarketSnapshot(ctx context.Context, s model.MarketSnapshot, createdAt string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO market_snapshots (market_snapshot_id,symbol,trade_date,close_price,turnover_rate,margin_balance,margin_balance_change,pe_percentile,pb_percentile,volume_percentile,volatility_percentile,liquidity_state,sentiment_state,market_metrics_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, s.MarketSnapshotID, s.Symbol, valueOr(s.TradeDate, createdAt[:10]), s.ClosePrice, s.TurnoverRate, s.MarginBalance, s.MarginBalanceChange, s.PEPercentile, s.PBPercentile, s.VolumePercentile, s.VolatilityPercentile, string(s.LiquidityState), string(s.SentimentState), nullString(s.MarketMetricsJSON), createdAt)
	return apperr.FromRepositoryError(err)
}

// GetMarketSnapshot 读取市场快照。
func (r *MarketRepository) GetMarketSnapshot(ctx context.Context, id string) (model.MarketSnapshot, error) {
	var s model.MarketSnapshot
	err := r.db.QueryRowContext(ctx, `SELECT market_snapshot_id,symbol,trade_date,COALESCE(close_price,0),COALESCE(turnover_rate,0),COALESCE(margin_balance,0),COALESCE(margin_balance_change,0),COALESCE(pe_percentile,0),COALESCE(pb_percentile,0),COALESCE(volume_percentile,0),COALESCE(volatility_percentile,0),COALESCE(liquidity_state,''),COALESCE(sentiment_state,''),COALESCE(market_metrics_json,'{}') FROM market_snapshots WHERE market_snapshot_id=?`, id).Scan(&s.MarketSnapshotID, &s.Symbol, &s.TradeDate, &s.ClosePrice, &s.TurnoverRate, &s.MarginBalance, &s.MarginBalanceChange, &s.PEPercentile, &s.PBPercentile, &s.VolumePercentile, &s.VolatilityPercentile, &s.LiquidityState, &s.SentimentState, &s.MarketMetricsJSON)
	return s, apperr.FromRepositoryError(err)
}

// GetLatestMarketSnapshot reads the newest market snapshot.
func (r *MarketRepository) GetLatestMarketSnapshot(ctx context.Context) (model.MarketSnapshot, error) {
	var s model.MarketSnapshot
	err := r.db.QueryRowContext(ctx, `SELECT market_snapshot_id,symbol,trade_date,COALESCE(close_price,0),COALESCE(turnover_rate,0),COALESCE(margin_balance,0),COALESCE(margin_balance_change,0),COALESCE(pe_percentile,0),COALESCE(pb_percentile,0),COALESCE(volume_percentile,0),COALESCE(volatility_percentile,0),COALESCE(liquidity_state,''),COALESCE(sentiment_state,''),COALESCE(market_metrics_json,'{}') FROM market_snapshots ORDER BY created_at DESC LIMIT 1`).Scan(&s.MarketSnapshotID, &s.Symbol, &s.TradeDate, &s.ClosePrice, &s.TurnoverRate, &s.MarginBalance, &s.MarginBalanceChange, &s.PEPercentile, &s.PBPercentile, &s.VolumePercentile, &s.VolatilityPercentile, &s.LiquidityState, &s.SentimentState, &s.MarketMetricsJSON)
	return s, apperr.FromRepositoryError(err)
}

// GetLatestMarketSnapshotBySymbol reads the newest market snapshot for a symbol.
func (r *MarketRepository) GetLatestMarketSnapshotBySymbol(ctx context.Context, symbol string) (model.MarketSnapshot, error) {
	var s model.MarketSnapshot
	err := r.db.QueryRowContext(ctx, `SELECT market_snapshot_id,symbol,trade_date,COALESCE(close_price,0),COALESCE(turnover_rate,0),COALESCE(margin_balance,0),COALESCE(margin_balance_change,0),COALESCE(pe_percentile,0),COALESCE(pb_percentile,0),COALESCE(liquidity_state,'normal'),COALESCE(sentiment_state,'neutral'),COALESCE(volume_percentile,0),COALESCE(volatility_percentile,0),COALESCE(market_metrics_json,'{}') FROM market_snapshots WHERE symbol=? ORDER BY created_at DESC LIMIT 1`, symbol).Scan(&s.MarketSnapshotID, &s.Symbol, &s.TradeDate, &s.ClosePrice, &s.TurnoverRate, &s.MarginBalance, &s.MarginBalanceChange, &s.PEPercentile, &s.PBPercentile, &s.LiquidityState, &s.SentimentState, &s.VolumePercentile, &s.VolatilityPercentile, &s.MarketMetricsJSON)
	return s, apperr.FromRepositoryError(err)
}

func valueOr(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

// MarketSnapshotExists reports whether a snapshot ID exists.
func (r *MarketRepository) MarketSnapshotExists(ctx context.Context, id string) (bool, error) {
	var exists int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM market_snapshots WHERE market_snapshot_id=?`, id).Scan(&exists)
	return exists > 0, apperr.FromRepositoryError(err)
}
