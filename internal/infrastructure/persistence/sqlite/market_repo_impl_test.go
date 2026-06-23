package sqlite

import (
	"context"
	"testing"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/pkg/apperr"
)

func TestMarketRepositoryClassifiesErrors(t *testing.T) {
	db := testDB(t)
	repo := NewMarketRepository(db)
	if _, err := repo.GetMarketSnapshot(context.Background(), "missing_market"); !apperr.IsCode(err, apperr.CodeNotFound) {
		t.Fatalf("expected not found market error, got %v", err)
	}
	bad := model.MarketSnapshot{MarketSnapshotID: "bad_market", Symbol: "510300", LiquidityState: "invalid", SentimentState: model.SentimentNeutral}
	if err := repo.SaveMarketSnapshot(context.Background(), bad, testTime); !apperr.IsCode(err, apperr.CodeConflict) {
		t.Fatalf("expected conflict market error, got %v", err)
	}
}

func TestMarketRepositoryPreservesStructuredFinancialFields(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	repo := NewMarketRepository(db)
	snapshot := model.MarketSnapshot{
		MarketSnapshotID:    "market_financials",
		Symbol:              "510300",
		TradeDate:           "2026-06-20",
		ClosePrice:          4.23,
		TurnoverRate:        1.8,
		MarginBalance:       120000000,
		MarginBalanceChange: -0.08,
		PEPercentile:        31,
		PBPercentile:        27,
		LiquidityState:      model.LiquidityNormal,
		SentimentState:      model.SentimentNeutral,
	}

	if err := repo.SaveMarketSnapshot(ctx, snapshot, testTime); err != nil {
		t.Fatalf("SaveMarketSnapshot: %v", err)
	}
	got, err := repo.GetMarketSnapshot(ctx, "market_financials")
	if err != nil {
		t.Fatalf("GetMarketSnapshot: %v", err)
	}
	if got.MarginBalance != 120000000 || got.MarginBalanceChange != -0.08 || got.PEPercentile != 31 || got.PBPercentile != 27 {
		t.Fatalf("structured financial fields not preserved: %+v", got)
	}
	latest, err := repo.GetLatestMarketSnapshotBySymbol(ctx, "510300")
	if err != nil {
		t.Fatalf("GetLatestMarketSnapshotBySymbol: %v", err)
	}
	if latest.MarginBalance != 120000000 || latest.MarginBalanceChange != -0.08 {
		t.Fatalf("latest snapshot should preserve structured financial fields: %+v", latest)
	}
}
