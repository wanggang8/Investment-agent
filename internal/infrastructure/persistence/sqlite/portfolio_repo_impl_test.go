package sqlite

import (
	"context"
	"testing"

	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

const testTime = "2026-05-28T00:00:00Z"

func TestPortfolioRepositoryWriteReadAndRollback(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	repo := NewPortfolioRepository(db)

	snap := repository.PortfolioSnapshot{SnapshotID: "snap1", SnapshotTime: testTime, Cash: 100, TotalAssets: 200, CashRatio: 0.5, HighRiskRatio: 0.1, PositionCount: 1, Source: "manual", CreatedAt: testTime}
	positions := []repository.PositionSnapshot{{PositionSnapshotID: "ps1", SnapshotID: "snap1", Symbol: "AAA", Name: "Alpha", Quantity: 1, CostPrice: 10, CurrentPrice: 11, MarketValue: 11, UnrealizedProfitRatio: 0.1, PositionState: "normal", CreatedAt: testTime}}
	if err := repo.SavePortfolioSnapshot(ctx, snap, positions); err != nil {
		t.Fatal(err)
	}
	got, gotPositions, err := repo.GetPortfolioSnapshot(ctx, "snap1")
	if err != nil {
		t.Fatal(err)
	}
	if got.SnapshotID != "snap1" || len(gotPositions) != 1 {
		t.Fatalf("unexpected read: %#v %#v", got, gotPositions)
	}

	bad := repository.PortfolioSnapshot{SnapshotID: "snap_bad", SnapshotTime: testTime, Cash: 1, TotalAssets: 1, Source: "manual", CreatedAt: testTime}
	badPositions := []repository.PositionSnapshot{{PositionSnapshotID: "badps", SnapshotID: "snap_bad", Symbol: "BBB", Name: "Beta", PositionState: "invalid", CreatedAt: testTime}}
	if err := repo.SavePortfolioSnapshot(ctx, bad, badPositions); err == nil {
		t.Fatal("expected rollback error")
	}
	if _, _, err := repo.GetPortfolioSnapshot(ctx, "snap_bad"); err == nil {
		t.Fatal("snapshot persisted after rollback")
	}
}

func TestPositionWriteRead(t *testing.T) {
	db := testDB(t)
	repo := NewPortfolioRepository(db)
	p := repository.Position{PositionID: "pos1", Symbol: "AAA", Name: "Alpha", Quantity: 2, CostPrice: 10, CurrentPrice: 12, MarketValue: 24, UnrealizedProfitRatio: 0.2, PositionState: "normal", UpdatedAt: testTime}
	if err := repo.SavePosition(context.Background(), p); err != nil {
		t.Fatal(err)
	}
	got, err := repo.GetPosition(context.Background(), "pos1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Symbol != "AAA" {
		t.Fatalf("got %#v", got)
	}
}

func TestPortfolioRepositoryClassifiesErrors(t *testing.T) {
	db := testDB(t)
	repo := NewPortfolioRepository(db)
	if _, _, err := repo.GetPortfolioSnapshot(context.Background(), "missing_snapshot"); !apperr.IsCode(err, apperr.CodeNotFound) {
		t.Fatalf("expected not found snapshot error, got %v", err)
	}
	if _, err := repo.GetPosition(context.Background(), "missing_position"); !apperr.IsCode(err, apperr.CodeNotFound) {
		t.Fatalf("expected not found position error, got %v", err)
	}
	bad := repository.Position{PositionID: "bad_position", Symbol: "AAA", Name: "Alpha", PositionState: "invalid", UpdatedAt: testTime}
	if err := repo.SavePosition(context.Background(), bad); !apperr.IsCode(err, apperr.CodeConflict) {
		t.Fatalf("expected conflict position error, got %v", err)
	}
}
