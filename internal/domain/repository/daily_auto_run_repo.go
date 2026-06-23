package repository

import "context"

// DailyAutoRunState 记录一个幂等运行范围的本地 scheduler 最新状态。
type DailyAutoRunState struct {
	RunID          string
	IdempotencyKey string
	LocalDate      string
	Scope          string
	SymbolSetHash  string
	Status         string
	LastRunAt      string
	NextRunAt      string
	FailureCode    string
	FailureReason  string
	CreatedAt      string
	UpdatedAt      string
}

// DailyAutoRunRepository 持久化本地每日自动运行状态。
type DailyAutoRunRepository interface {
	UpsertDailyAutoRunState(ctx context.Context, state DailyAutoRunState) error
	GetDailyAutoRunState(ctx context.Context, idempotencyKey string) (DailyAutoRunState, error)
	GetLatestDailyAutoRunState(ctx context.Context) (DailyAutoRunState, error)
}
