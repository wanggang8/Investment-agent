package service

import (
	"encoding/json"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/domain/model"
)

// SourceHealthFromMarketSnapshot derives P34 source health DTOs from normalized market metadata.
func SourceHealthFromMarketSnapshot(market model.MarketSnapshot) []dto.SourceHealthItem {
	metrics := parseSourceHealthJSONMap(market.MarketMetricsJSON)
	metadata, _ := metrics["metadata"].(map[string]any)
	health, _ := metadata["p34_source_health"].(map[string]any)
	rawCategories, _ := metadata["p34_data_categories"].([]any)
	categories := make([]string, 0, len(rawCategories))
	for _, raw := range rawCategories {
		if category, ok := raw.(string); ok && category != "" {
			categories = append(categories, category)
		}
	}
	if len(categories) == 0 {
		for category := range health {
			categories = append(categories, category)
		}
	}
	items := make([]dto.SourceHealthItem, 0, len(categories))
	for _, category := range categories {
		item := dto.SourceHealthItem{
			SourceName:      sourceHealthStringFromMap(metrics, "source_name"),
			SourceLevel:     sourceHealthStringFromMap(metrics, "source_level"),
			SourceType:      sourceHealthStringFromMap(metrics, "source_type"),
			DataCategory:    category,
			Freshness:       "missing",
			DataDate:        normalizeSourceHealthDate(market.TradeDate),
			RequestID:       sourceHealthStringFromMap(metrics, "request_id"),
			AffectedSymbols: []string{market.Symbol},
		}
		switch raw := health[category].(type) {
		case map[string]any:
			item.Freshness = sourceHealthStringFromMap(raw, "freshness")
			if item.Freshness == "" {
				item.Freshness = "missing"
			}
			if value := sourceHealthStringFromMap(raw, "source_name"); value != "" {
				item.SourceName = value
			}
			if value := sourceHealthStringFromMap(raw, "source_level"); value != "" {
				item.SourceLevel = value
			}
			if value := sourceHealthStringFromMap(raw, "source_type"); value != "" {
				item.SourceType = value
			}
			if value := sourceHealthStringFromMap(raw, "data_date"); value != "" {
				item.DataDate = normalizeSourceHealthDate(value)
			}
			if value := sourceHealthStringFromMap(raw, "request_id"); value != "" {
				item.RequestID = value
			}
			item.LastSuccessAt = sourceHealthStringFromMap(raw, "last_success_at")
			item.LastFailureAt = sourceHealthStringFromMap(raw, "last_failure_at")
			item.FailureCategory = sourceHealthStringFromMap(raw, "failure_category")
			if symbols := sourceHealthStringSliceFromAny(raw["affected_symbols"]); len(symbols) > 0 {
				item.AffectedSymbols = symbols
			}
		case string:
			item.Freshness = raw
			if item.Freshness == "" {
				item.Freshness = "missing"
			}
			if item.Freshness == "fresh" {
				item.LastSuccessAt = sourceHealthStringFromMap(metrics, "captured_at")
			} else {
				item.LastFailureAt = sourceHealthStringFromMap(metrics, "captured_at")
				item.FailureCategory = item.Freshness
			}
		}
		if item.FailureCategory == "" && item.Freshness != "fresh" && item.Freshness != "stubbed" {
			item.FailureCategory = item.Freshness
		}
		items = append(items, item)
	}
	return items
}

func parseSourceHealthJSONMap(raw string) map[string]any {
	if raw == "" {
		return map[string]any{}
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return map[string]any{}
	}
	return out
}

func normalizeSourceHealthDate(value string) string {
	if len(value) >= 10 {
		return value[:10]
	}
	return value
}

func sourceHealthStringSliceFromAny(value any) []string {
	rawItems, _ := value.([]any)
	items := make([]string, 0, len(rawItems))
	for _, raw := range rawItems {
		if item, ok := raw.(string); ok && item != "" {
			items = append(items, item)
		}
	}
	return items
}

func sourceHealthStringFromMap(values map[string]any, key string) string {
	value, _ := values[key].(string)
	return value
}
