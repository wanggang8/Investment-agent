package workflow

import (
	"context"
	"testing"

	"investment-agent/internal/domain/model"
)

func TestP88StructuredDataNormalizesCapitalFlowReadback(t *testing.T) {
	readback := P88NormalizeStructuredDataMetadata(map[string]any{
		"capital_flow": map[string]any{
			"trade_date":  "2026-06-22",
			"net_inflow":  123456.78,
			"net_outflow": "23456.78",
		},
	})

	if readback.CapitalFlow == nil {
		t.Fatalf("expected capital-flow readback")
	}
	if readback.CapitalFlow.Date != "2026-06-22" || readback.CapitalFlow.NetInflow != 123456.78 || readback.CapitalFlow.NetOutflow != 23456.78 {
		t.Fatalf("unexpected capital-flow fields: %+v", readback.CapitalFlow)
	}
}

func TestP90StructuredDataNormalizesCapitalFlowRawNetFlow(t *testing.T) {
	readback := P88NormalizeStructuredDataMetadata(map[string]any{
		"capital_flow": map[string]any{
			"trade_date":   "2026-06-22",
			"net_inflow":   0,
			"net_outflow":  74801057.0,
			"raw_net_flow": -74801057.0,
		},
	})

	if readback.CapitalFlow == nil {
		t.Fatalf("expected capital-flow readback")
	}
	if readback.CapitalFlow.RawNetFlow != -74801057.0 {
		t.Fatalf("expected raw net-flow to survive readback, got %+v", readback.CapitalFlow)
	}
}

func TestP90StructuredDataKeepsZeroRawNetFlowDay(t *testing.T) {
	readback := P88NormalizeStructuredDataMetadata(map[string]any{
		"capital_flow": map[string]any{
			"trade_date":   "2026-06-22",
			"net_inflow":   0,
			"net_outflow":  0,
			"raw_net_flow": 0,
		},
	})

	if readback.CapitalFlow == nil {
		t.Fatalf("expected zero raw net-flow day to remain valid")
	}
	if readback.CapitalFlow.RawNetFlow != 0 || readback.CapitalFlow.NetInflow != 0 || readback.CapitalFlow.NetOutflow != 0 {
		t.Fatalf("unexpected zero raw net-flow readback: %+v", readback.CapitalFlow)
	}
}

func TestP88StructuredDataNormalizesMarginFinancingReadback(t *testing.T) {
	readback := P88NormalizeStructuredDataMetadata(map[string]any{
		"margin_financing": map[string]any{
			"data_date":           "2026-06-22",
			"margin_balance":      "987654321.12",
			"balance_change_rate": -0.0123,
		},
	})

	if readback.MarginFinancing == nil {
		t.Fatalf("expected margin-financing readback")
	}
	if readback.MarginFinancing.Date != "2026-06-22" || readback.MarginFinancing.MarginBalance != 987654321.12 || readback.MarginFinancing.BalanceChangeRate != -0.0123 {
		t.Fatalf("unexpected margin-financing fields: %+v", readback.MarginFinancing)
	}
}

func TestP88StructuredDataNormalizesConstituentFinancialReadback(t *testing.T) {
	readback := P88NormalizeStructuredDataMetadata(map[string]any{
		"constituent_financial": map[string]any{
			"operating_revenue": 4523000000.0,
			"net_profit":        812000000.0,
			"growth_rate":       0.137,
			"disclosure_date":   "2026-04-30",
		},
	})

	if readback.ConstituentFinancial == nil {
		t.Fatalf("expected constituent-financial readback")
	}
	if readback.ConstituentFinancial.Revenue != 4523000000.0 || readback.ConstituentFinancial.NetProfit != 812000000.0 || readback.ConstituentFinancial.Growth != 0.137 || readback.ConstituentFinancial.DisclosureDate != "2026-04-30" {
		t.Fatalf("unexpected constituent-financial fields: %+v", readback.ConstituentFinancial)
	}
}

func TestP88MarketRefreshPersistsStructuredDataReadback(t *testing.T) {
	deps := WorkflowDependencies{MarketDataSource: testMarketDataSource{point: MarketDataPoint{
		ClosePrice:  4.321,
		SourceName:  "p88_public_structured_source",
		SourceLevel: model.SourceLevelB,
		SourceType:  "public_structured_market",
		TradeDate:   "2026-06-22",
		Metadata: map[string]any{
			"capital_flow": map[string]any{
				"date":        "2026-06-22",
				"net_inflow":  1000.0,
				"net_outflow": 300.0,
			},
			"margin_financing": map[string]any{
				"date":                "2026-06-22",
				"margin_balance":      100000000.0,
				"balance_change_rate": 0.023,
			},
			"constituent_financial": map[string]any{
				"revenue":         4523000000.0,
				"net_profit":      812000000.0,
				"growth":          0.137,
				"disclosure_date": "2026-04-30",
			},
		},
	}}}

	out, err := NewMarketRefreshGraphWithDependencies(deps).Run(context.Background(), MarketRefreshInput{RequestID: "req_p88_structured_fields", Symbol: "600000"})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if out.MarketSnapshot.MarginBalance != 100000000.0 || out.MarketSnapshot.MarginBalanceChange != 0.023 {
		t.Fatalf("expected margin fields in SQLite snapshot model, got %+v", out.MarketSnapshot)
	}
	readback := P88StructuredDataReadbackFromMarketSnapshot(out.MarketSnapshot)
	if readback.CapitalFlow == nil || readback.MarginFinancing == nil || readback.ConstituentFinancial == nil {
		t.Fatalf("expected all P88 structured fields from market_metrics_json, got %+v", readback)
	}
	if readback.CapitalFlow.NetInflow != 1000.0 || readback.MarginFinancing.MarginBalance != 100000000.0 || readback.ConstituentFinancial.DisclosureDate != "2026-04-30" {
		t.Fatalf("unexpected P88 structured readback: %+v", readback)
	}
}
