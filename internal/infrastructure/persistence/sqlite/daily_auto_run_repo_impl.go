package sqlite

import (
	"context"
	"database/sql"

	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

// DailyAutoRunRepository 是本地 scheduler 状态的 SQLite 实现。
type DailyAutoRunRepository struct{ db dbtx }

// NewDailyAutoRunRepository 创建每日自动运行状态仓储。
func NewDailyAutoRunRepository(db *sql.DB) *DailyAutoRunRepository {
	return &DailyAutoRunRepository{db: db}
}

func (r *DailyAutoRunRepository) UpsertDailyAutoRunState(ctx context.Context, state repository.DailyAutoRunState) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO daily_auto_run_states (run_id,idempotency_key,local_date,scope,symbol_set_hash,status,last_run_at,next_run_at,failure_code,failure_reason,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?) ON CONFLICT(idempotency_key) DO UPDATE SET run_id=excluded.run_id,local_date=excluded.local_date,scope=excluded.scope,symbol_set_hash=excluded.symbol_set_hash,status=excluded.status,last_run_at=excluded.last_run_at,next_run_at=excluded.next_run_at,failure_code=excluded.failure_code,failure_reason=excluded.failure_reason,updated_at=excluded.updated_at`, state.RunID, state.IdempotencyKey, state.LocalDate, state.Scope, state.SymbolSetHash, state.Status, nullString(state.LastRunAt), nullString(state.NextRunAt), nullString(state.FailureCode), nullString(state.FailureReason), state.CreatedAt, state.UpdatedAt)
	return apperr.FromRepositoryError(err)
}

func (r *DailyAutoRunRepository) GetDailyAutoRunState(ctx context.Context, idempotencyKey string) (repository.DailyAutoRunState, error) {
	var state repository.DailyAutoRunState
	err := r.db.QueryRowContext(ctx, `SELECT run_id,idempotency_key,local_date,scope,symbol_set_hash,status,COALESCE(last_run_at,''),COALESCE(next_run_at,''),COALESCE(failure_code,''),COALESCE(failure_reason,''),created_at,updated_at FROM daily_auto_run_states WHERE idempotency_key=?`, idempotencyKey).Scan(&state.RunID, &state.IdempotencyKey, &state.LocalDate, &state.Scope, &state.SymbolSetHash, &state.Status, &state.LastRunAt, &state.NextRunAt, &state.FailureCode, &state.FailureReason, &state.CreatedAt, &state.UpdatedAt)
	return state, apperr.FromRepositoryError(err)
}

func (r *DailyAutoRunRepository) GetLatestDailyAutoRunState(ctx context.Context) (repository.DailyAutoRunState, error) {
	var state repository.DailyAutoRunState
	err := r.db.QueryRowContext(ctx, `SELECT run_id,idempotency_key,local_date,scope,symbol_set_hash,status,COALESCE(last_run_at,''),COALESCE(next_run_at,''),COALESCE(failure_code,''),COALESCE(failure_reason,''),created_at,updated_at FROM daily_auto_run_states ORDER BY updated_at DESC, created_at DESC LIMIT 1`).Scan(&state.RunID, &state.IdempotencyKey, &state.LocalDate, &state.Scope, &state.SymbolSetHash, &state.Status, &state.LastRunAt, &state.NextRunAt, &state.FailureCode, &state.FailureReason, &state.CreatedAt, &state.UpdatedAt)
	return state, apperr.FromRepositoryError(err)
}
