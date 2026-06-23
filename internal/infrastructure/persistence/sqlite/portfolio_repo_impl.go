package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

// PortfolioRepository 是账户快照与持仓表的 SQLite 实现。
type PortfolioRepository struct{ db dbtx }

// NewPortfolioRepository 创建账户仓储实例。
func NewPortfolioRepository(db *sql.DB) *PortfolioRepository { return &PortfolioRepository{db: db} }

// SavePortfolioSnapshot 在同一事务中保存账户快照和对应的持仓时点快照。
func (r *PortfolioRepository) SavePortfolioSnapshot(ctx context.Context, s repository.PortfolioSnapshot, ps []repository.PositionSnapshot) error {
	err := withTx(ctx, r.db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO portfolio_snapshots (snapshot_id,snapshot_time,cash,total_assets,cash_ratio,high_risk_ratio,position_count,source,created_at) VALUES (?,?,?,?,?,?,?,?,?)`, s.SnapshotID, s.SnapshotTime, s.Cash, s.TotalAssets, s.CashRatio, s.HighRiskRatio, s.PositionCount, s.Source, s.CreatedAt)
		if err != nil {
			return err
		}
		for _, p := range ps {
			_, err = tx.ExecContext(ctx, `INSERT INTO position_snapshots (position_snapshot_id,snapshot_id,position_id,symbol,name,quantity,cost_price,current_price,market_value,unrealized_profit_ratio,position_state,buy_date,buy_reason,asset_tag,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, p.PositionSnapshotID, p.SnapshotID, nullString(p.PositionID), p.Symbol, p.Name, p.Quantity, p.CostPrice, p.CurrentPrice, p.MarketValue, p.UnrealizedProfitRatio, p.PositionState, nullString(p.BuyDate), nullString(p.BuyReason), nullString(p.AssetTag), p.CreatedAt)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return apperr.FromRepositoryError(err)
}

// GetPortfolioSnapshot 读取账户快照及其完整持仓明细。
func (r *PortfolioRepository) GetPortfolioSnapshot(ctx context.Context, id string) (repository.PortfolioSnapshot, []repository.PositionSnapshot, error) {
	var s repository.PortfolioSnapshot
	err := r.db.QueryRowContext(ctx, `SELECT snapshot_id,snapshot_time,cash,total_assets,cash_ratio,high_risk_ratio,position_count,source,created_at FROM portfolio_snapshots WHERE snapshot_id=?`, id).Scan(&s.SnapshotID, &s.SnapshotTime, &s.Cash, &s.TotalAssets, &s.CashRatio, &s.HighRiskRatio, &s.PositionCount, &s.Source, &s.CreatedAt)
	if err != nil {
		return s, nil, apperr.FromRepositoryError(err)
	}
	rows, err := r.db.QueryContext(ctx, `SELECT position_snapshot_id,snapshot_id,COALESCE(position_id,''),symbol,name,quantity,cost_price,current_price,market_value,unrealized_profit_ratio,position_state,COALESCE(buy_date,''),COALESCE(buy_reason,''),COALESCE(asset_tag,''),created_at FROM position_snapshots WHERE snapshot_id=? ORDER BY position_snapshot_id`, id)
	if err != nil {
		return s, nil, apperr.FromRepositoryError(err)
	}
	defer rows.Close()
	var out []repository.PositionSnapshot
	for rows.Next() {
		var p repository.PositionSnapshot
		if err := rows.Scan(&p.PositionSnapshotID, &p.SnapshotID, &p.PositionID, &p.Symbol, &p.Name, &p.Quantity, &p.CostPrice, &p.CurrentPrice, &p.MarketValue, &p.UnrealizedProfitRatio, &p.PositionState, &p.BuyDate, &p.BuyReason, &p.AssetTag, &p.CreatedAt); err != nil {
			return s, nil, apperr.FromRepositoryError(err)
		}
		out = append(out, p)
	}
	return s, out, apperr.FromRepositoryError(rows.Err())
}

// GetLatestPortfolioSnapshot reads the newest portfolio snapshot.
func (r *PortfolioRepository) GetLatestPortfolioSnapshot(ctx context.Context) (repository.PortfolioSnapshot, error) {
	var s repository.PortfolioSnapshot
	err := r.db.QueryRowContext(ctx, `SELECT snapshot_id,snapshot_time,cash,total_assets,cash_ratio,high_risk_ratio,position_count,source,created_at FROM portfolio_snapshots ORDER BY snapshot_time DESC LIMIT 1`).Scan(&s.SnapshotID, &s.SnapshotTime, &s.Cash, &s.TotalAssets, &s.CashRatio, &s.HighRiskRatio, &s.PositionCount, &s.Source, &s.CreatedAt)
	return s, apperr.FromRepositoryError(err)
}

// SavePosition 保存当前持仓聚合态。
func (r *PortfolioRepository) SavePosition(ctx context.Context, p repository.Position) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO positions (position_id,symbol,name,quantity,cost_price,current_price,market_value,unrealized_profit_ratio,position_state,buy_date,buy_reason,asset_tag,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?) ON CONFLICT(position_id) DO UPDATE SET symbol=excluded.symbol,name=excluded.name,quantity=excluded.quantity,cost_price=excluded.cost_price,current_price=excluded.current_price,market_value=excluded.market_value,unrealized_profit_ratio=excluded.unrealized_profit_ratio,position_state=excluded.position_state,buy_date=excluded.buy_date,buy_reason=excluded.buy_reason,asset_tag=excluded.asset_tag,updated_at=excluded.updated_at`, p.PositionID, p.Symbol, p.Name, p.Quantity, p.CostPrice, p.CurrentPrice, p.MarketValue, p.UnrealizedProfitRatio, p.PositionState, nullString(p.BuyDate), nullString(p.BuyReason), nullString(p.AssetTag), p.UpdatedAt)
	return apperr.FromRepositoryError(err)
}

// ReplacePositions 用一次新账户快照的持仓集合替换当前持仓聚合态。
func (r *PortfolioRepository) ReplacePositions(ctx context.Context, positions []repository.Position) error {
	err := withTx(ctx, r.db, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `DELETE FROM positions`); err != nil {
			return err
		}
		for _, p := range positions {
			_, err := tx.ExecContext(ctx, `INSERT INTO positions (position_id,symbol,name,quantity,cost_price,current_price,market_value,unrealized_profit_ratio,position_state,buy_date,buy_reason,asset_tag,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`, p.PositionID, p.Symbol, p.Name, p.Quantity, p.CostPrice, p.CurrentPrice, p.MarketValue, p.UnrealizedProfitRatio, p.PositionState, nullString(p.BuyDate), nullString(p.BuyReason), nullString(p.AssetTag), p.UpdatedAt)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return apperr.FromRepositoryError(err)
}

// DeletePosition 删除已清仓的当前持仓聚合态。
func (r *PortfolioRepository) DeletePosition(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM positions WHERE position_id=?`, id)
	return apperr.FromRepositoryError(err)
}

// GetPosition 读取当前持仓聚合态。
func (r *PortfolioRepository) GetPosition(ctx context.Context, id string) (repository.Position, error) {
	var p repository.Position
	err := r.db.QueryRowContext(ctx, `SELECT position_id,symbol,name,quantity,cost_price,current_price,market_value,unrealized_profit_ratio,position_state,COALESCE(buy_date,''),COALESCE(buy_reason,''),COALESCE(asset_tag,''),updated_at FROM positions WHERE position_id=?`, id).Scan(&p.PositionID, &p.Symbol, &p.Name, &p.Quantity, &p.CostPrice, &p.CurrentPrice, &p.MarketValue, &p.UnrealizedProfitRatio, &p.PositionState, &p.BuyDate, &p.BuyReason, &p.AssetTag, &p.UpdatedAt)
	return p, apperr.FromRepositoryError(err)
}

// ListPositions reads current positions ordered by symbol.
func (r *PortfolioRepository) ListPositions(ctx context.Context) ([]repository.Position, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT position_id,symbol,name,quantity,cost_price,current_price,market_value,unrealized_profit_ratio,position_state,COALESCE(buy_date,''),COALESCE(buy_reason,''),COALESCE(asset_tag,''),updated_at FROM positions ORDER BY symbol`)
	if err != nil {
		return nil, apperr.FromRepositoryError(err)
	}
	defer rows.Close()
	var out []repository.Position
	for rows.Next() {
		var p repository.Position
		if err := rows.Scan(&p.PositionID, &p.Symbol, &p.Name, &p.Quantity, &p.CostPrice, &p.CurrentPrice, &p.MarketValue, &p.UnrealizedProfitRatio, &p.PositionState, &p.BuyDate, &p.BuyReason, &p.AssetTag, &p.UpdatedAt); err != nil {
			return nil, apperr.FromRepositoryError(err)
		}
		out = append(out, p)
	}
	return out, apperr.FromRepositoryError(rows.Err())
}

func (r *PortfolioRepository) SaveLocalAccountImportBatch(ctx context.Context, b repository.LocalAccountImportBatch) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO local_account_import_batches (import_batch_id,request_id,status,row_count,valid_count,invalid_count,validation_summary_json,rows_hash,created_at,committed_at) VALUES (?,?,?,?,?,?,?,?,?,?) ON CONFLICT(import_batch_id) DO UPDATE SET status=excluded.status,row_count=excluded.row_count,valid_count=excluded.valid_count,invalid_count=excluded.invalid_count,validation_summary_json=excluded.validation_summary_json,rows_hash=excluded.rows_hash,committed_at=excluded.committed_at`, b.ImportBatchID, b.RequestID, b.Status, b.RowCount, b.ValidCount, b.InvalidCount, nullString(b.ValidationSummaryJSON), nullString(b.RowsHash), b.CreatedAt, nullString(b.CommittedAt))
	return apperr.FromRepositoryError(err)
}

func (r *PortfolioRepository) GetLocalAccountImportBatch(ctx context.Context, importBatchID string) (repository.LocalAccountImportBatch, error) {
	var b repository.LocalAccountImportBatch
	err := r.db.QueryRowContext(ctx, `SELECT import_batch_id,request_id,status,row_count,valid_count,invalid_count,COALESCE(validation_summary_json,''),COALESCE(rows_hash,''),created_at,COALESCE(committed_at,'') FROM local_account_import_batches WHERE import_batch_id=?`, importBatchID).Scan(&b.ImportBatchID, &b.RequestID, &b.Status, &b.RowCount, &b.ValidCount, &b.InvalidCount, &b.ValidationSummaryJSON, &b.RowsHash, &b.CreatedAt, &b.CommittedAt)
	return b, apperr.FromRepositoryError(err)
}

func (r *PortfolioRepository) SaveLocalAccountCorrection(ctx context.Context, c repository.LocalAccountCorrection) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO local_account_corrections (correction_id,target_type,target_id,before_json,after_json,correction_reason,snapshot_id,audit_event_id,created_at) VALUES (?,?,?,?,?,?,?,?,?)`, c.CorrectionID, c.TargetType, c.TargetID, c.BeforeJSON, c.AfterJSON, c.CorrectionReason, nullString(c.SnapshotID), nullString(c.AuditEventID), c.CreatedAt)
	return apperr.FromRepositoryError(err)
}

// withTx 封装事务提交和回滚，保证组合写入的原子性。
func withTx(ctx context.Context, db dbtx, fn func(*sql.Tx) error) error {
	if tx, ok := db.(txDB); ok {
		return fn(tx.Tx)
	}
	sqlDB, ok := db.(*sql.DB)
	if !ok {
		return errors.New("sqlite repository requires *sql.DB or txDB")
	}
	tx, err := sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		return errors.Join(err, tx.Rollback())
	}
	return tx.Commit()
}

// nullString 把空字符串写为 SQL NULL，读取时通过 COALESCE 还原为空字符串。
func nullString(v string) any {
	if v == "" {
		return nil
	}
	return v
}
