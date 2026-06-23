package sqlite

import (
	"context"
	"testing"

	"investment-agent/internal/domain/repository"
)

func TestDailyAutoRunRepositoryUpsertAndGetState(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	repo := NewDailyAutoRunRepository(db)

	state := repository.DailyAutoRunState{
		RunID:          "run_1",
		IdempotencyKey: "2026-06-07:holdings:hash:v1",
		LocalDate:      "2026-06-07",
		Scope:          "holdings",
		SymbolSetHash:  "hash",
		Status:         "running",
		LastRunAt:      "2026-06-07T00:30:00Z",
		NextRunAt:      "2026-06-08T00:30:00Z",
		CreatedAt:      testTime,
		UpdatedAt:      testTime,
	}
	if err := repo.UpsertDailyAutoRunState(ctx, state); err != nil {
		t.Fatal(err)
	}

	state.Status = "success"
	state.FailureCode = ""
	state.FailureReason = ""
	state.UpdatedAt = "2026-06-07T00:31:00Z"
	if err := repo.UpsertDailyAutoRunState(ctx, state); err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetDailyAutoRunState(ctx, state.IdempotencyKey)
	if err != nil {
		t.Fatal(err)
	}
	if got.RunID != "run_1" || got.Status != "success" || got.UpdatedAt != "2026-06-07T00:31:00Z" {
		t.Fatalf("unexpected state: %+v", got)
	}
}

func TestDailyAutoRunRepositoryRejectsInvalidStatus(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	repo := NewDailyAutoRunRepository(db)

	err := repo.UpsertDailyAutoRunState(ctx, repository.DailyAutoRunState{
		RunID:          "run_1",
		IdempotencyKey: "2026-06-07:holdings:hash:v1",
		LocalDate:      "2026-06-07",
		Scope:          "holdings",
		SymbolSetHash:  "hash",
		Status:         "trading",
		CreatedAt:      testTime,
		UpdatedAt:      testTime,
	})
	if err == nil {
		t.Fatal("expected invalid status to fail")
	}
}
