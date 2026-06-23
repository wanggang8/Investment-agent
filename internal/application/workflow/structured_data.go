package workflow

import (
	"encoding/json"
	"strconv"
	"strings"

	"investment-agent/internal/domain/model"
)

type P88StructuredDataReadback struct {
	CapitalFlow          *P88CapitalFlow          `json:"capital_flow,omitempty"`
	MarginFinancing      *P88MarginFinancing      `json:"margin_financing,omitempty"`
	ConstituentFinancial *P88ConstituentFinancial `json:"constituent_financial,omitempty"`
}

type P88CapitalFlow struct {
	Date       string  `json:"date"`
	NetInflow  float64 `json:"net_inflow"`
	NetOutflow float64 `json:"net_outflow"`
	RawNetFlow float64 `json:"raw_net_flow,omitempty"`
}

type P88MarginFinancing struct {
	Date              string  `json:"date"`
	MarginBalance     float64 `json:"margin_balance"`
	BalanceChangeRate float64 `json:"balance_change_rate"`
}

type P88ConstituentFinancial struct {
	Revenue        float64 `json:"revenue"`
	NetProfit      float64 `json:"net_profit"`
	Growth         float64 `json:"growth"`
	DisclosureDate string  `json:"disclosure_date"`
}

func (r P88StructuredDataReadback) Empty() bool {
	return r.CapitalFlow == nil && r.MarginFinancing == nil && r.ConstituentFinancial == nil
}

func P88NormalizeStructuredDataMetadata(metadata map[string]any) P88StructuredDataReadback {
	return P88StructuredDataReadback{
		CapitalFlow:          p88CapitalFlowFromAny(metadata["capital_flow"]),
		MarginFinancing:      p88MarginFinancingFromAny(metadata["margin_financing"]),
		ConstituentFinancial: p88ConstituentFinancialFromAny(metadata["constituent_financial"]),
	}
}

func P88StructuredDataReadbackFromMarketSnapshot(snapshot model.MarketSnapshot) P88StructuredDataReadback {
	var metrics map[string]any
	if err := json.Unmarshal([]byte(snapshot.MarketMetricsJSON), &metrics); err != nil {
		return P88StructuredDataReadback{}
	}
	metadata, _ := metrics["metadata"].(map[string]any)
	raw := metadata["p88_structured_fields"]
	if raw == nil {
		raw = metrics["p88_structured_fields"]
	}
	data, err := json.Marshal(raw)
	if err != nil || string(data) == "null" {
		return P88StructuredDataReadback{}
	}
	var out P88StructuredDataReadback
	if err := json.Unmarshal(data, &out); err != nil {
		return P88StructuredDataReadback{}
	}
	if out.MarginFinancing != nil && out.MarginFinancing.MarginBalance == 0 && snapshot.MarginBalance != 0 {
		out.MarginFinancing.MarginBalance = snapshot.MarginBalance
	}
	if out.MarginFinancing != nil && out.MarginFinancing.BalanceChangeRate == 0 && snapshot.MarginBalanceChange != 0 {
		out.MarginFinancing.BalanceChangeRate = snapshot.MarginBalanceChange
	}
	return out
}

func p88CapitalFlowFromAny(raw any) *P88CapitalFlow {
	item, ok := rawMap(raw)
	if !ok {
		return nil
	}
	date := stringValue(item, "date", "trade_date", "data_date")
	inflow := floatValue(item, "net_inflow", "main_net_inflow")
	outflow := floatValue(item, "net_outflow", "main_net_outflow")
	rawNetFlow, hasRawNetFlow := floatValueWithPresence(item, "raw_net_flow", "net_flow")
	if date == "" || (inflow == 0 && outflow == 0 && !hasRawNetFlow) {
		return nil
	}
	return &P88CapitalFlow{Date: date, NetInflow: inflow, NetOutflow: outflow, RawNetFlow: rawNetFlow}
}

func p88MarginFinancingFromAny(raw any) *P88MarginFinancing {
	item, ok := rawMap(raw)
	if !ok {
		return nil
	}
	date := stringValue(item, "date", "trade_date", "data_date")
	balance := floatValue(item, "margin_balance", "balance")
	changeRate := floatValue(item, "balance_change_rate", "margin_balance_change", "change_rate")
	if date == "" || balance == 0 {
		return nil
	}
	return &P88MarginFinancing{Date: date, MarginBalance: balance, BalanceChangeRate: changeRate}
}

func p88ConstituentFinancialFromAny(raw any) *P88ConstituentFinancial {
	item, ok := rawMap(raw)
	if !ok {
		return nil
	}
	disclosureDate := stringValue(item, "disclosure_date", "date", "report_date")
	revenue := floatValue(item, "revenue", "operating_revenue")
	netProfit := floatValue(item, "net_profit", "profit")
	growth := floatValue(item, "growth", "net_profit_growth", "growth_rate")
	if disclosureDate == "" || revenue == 0 || netProfit == 0 {
		return nil
	}
	return &P88ConstituentFinancial{Revenue: revenue, NetProfit: netProfit, Growth: growth, DisclosureDate: disclosureDate}
}

func rawMap(raw any) (map[string]any, bool) {
	item, ok := raw.(map[string]any)
	if ok {
		return item, true
	}
	data, err := json.Marshal(raw)
	if err != nil || string(data) == "null" {
		return nil, false
	}
	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, false
	}
	return out, true
}

func stringValue(item map[string]any, keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(stringFromAny(item[key])); value != "" {
			return value
		}
	}
	return ""
}

func floatValue(item map[string]any, keys ...string) float64 {
	value, _ := floatValueWithPresence(item, keys...)
	return value
}

func floatValueWithPresence(item map[string]any, keys ...string) (float64, bool) {
	for _, key := range keys {
		switch value := item[key].(type) {
		case float64:
			return value, true
		case int:
			return float64(value), true
		case int64:
			return float64(value), true
		case json.Number:
			got, _ := value.Float64()
			return got, true
		case string:
			got, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
			if err == nil {
				return got, true
			}
		}
	}
	return 0, false
}

func stringFromAny(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case json.Number:
		return typed.String()
	default:
		return ""
	}
}
