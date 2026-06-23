package workflow

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/pkg/apperr"
)

type CompositeMarketDataCollector struct {
	Collectors []MarketDataSource
}

func (c CompositeMarketDataCollector) FetchMarketData(ctx context.Context, symbol string) (MarketDataPoint, error) {
	var lastErr error
	for _, collector := range c.Collectors {
		point, err := collector.FetchMarketData(ctx, symbol)
		if err != nil {
			lastErr = err
			continue
		}
		return validateMarketPoint(point)
	}
	if lastErr != nil {
		return MarketDataPoint{}, lastErr
	}
	return MarketDataPoint{}, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "市场数据源未返回可用数据")
}

type EastmoneyFundCollector struct {
	BaseURL         string
	HTTPClient      *http.Client
	IncludeExtended bool
}

func (c EastmoneyFundCollector) FetchMarketData(ctx context.Context, symbol string) (MarketDataPoint, error) {
	base := strings.TrimRight(c.BaseURL, "/")
	if base == "" {
		base = "https://fund.eastmoney.com"
	}
	raw, err := fetchReadonlyURL(ctx, c.HTTPClient, base+"/pingzhongdata/"+strings.TrimSpace(symbol)+".js")
	if err != nil {
		return MarketDataPoint{}, err
	}
	return parseEastmoneyFundPoint(raw, symbol, c.IncludeExtended)
}

type CsindexCollector struct {
	BaseURL         string
	HTTPClient      *http.Client
	IncludeExtended bool
}

func (c CsindexCollector) FetchMarketData(ctx context.Context, symbol string) (MarketDataPoint, error) {
	base := strings.TrimRight(c.BaseURL, "/")
	if base == "" {
		base = "https://www.csindex.com.cn"
	}
	raw, err := fetchReadonlyURL(ctx, c.HTTPClient, base+"/csindex-home/indexInfo/index-basic-info/"+strings.TrimSpace(symbol))
	if err != nil {
		return MarketDataPoint{}, err
	}
	point, err := parseCsindexPoint(raw)
	if err != nil || !c.IncludeExtended {
		return point, err
	}
	c.fetchCsindexExtendedMetadata(ctx, base, strings.TrimSpace(symbol), point.Metadata)
	return point, nil
}

func parseEastmoneyFundPoint(raw []byte, symbol string, includeExtended bool) (MarketDataPoint, error) {
	name := eastmoneyStringVar(raw, "fS_name")
	code := eastmoneyStringVar(raw, "fS_code")
	if code == "" {
		code = strings.TrimSpace(symbol)
	}
	trend, err := eastmoneyNetWorthTrend(raw)
	if err != nil {
		return MarketDataPoint{}, err
	}
	accumulated := eastmoneyAccumulatedNAV(raw, trend.Timestamp)
	chinaTime, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return MarketDataPoint{}, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "东方财富基金交易日解析失败", err)
	}
	tradeDate := time.UnixMilli(trend.Timestamp).In(chinaTime).Format("2006-01-02")
	metadata := map[string]any{"fund_name": name, "fund_code": code, "symbol": code, "equity_return": trend.EquityReturn, "unit_money": trend.UnitMoney, "raw": map[string]any{"net_worth_timestamp": trend.Timestamp}}
	if accumulated > 0 {
		metadata["accumulated_nav"] = accumulated
	}
	if includeExtended {
		metadata["nav_history"] = eastmoneyNAVHistory(raw)
		if assetAllocation := eastmoneyJSONVar(raw, "Data_assetAllocation"); len(assetAllocation) > 0 {
			metadata["asset_allocation"] = assetAllocation
		}
		profile := map[string]any{}
		if performance := eastmoneyJSONVar(raw, "Data_performanceEvaluation"); len(performance) > 0 {
			profile["performance_evaluation"] = performance
		}
		if managers := eastmoneyJSONArrayVar(raw, "Data_currentFundManager"); len(managers) > 0 {
			profile["fund_managers"] = managers
		}
		if len(profile) > 0 {
			metadata["fund_profile"] = profile
		}
	}
	return MarketDataPoint{ClosePrice: trend.NAV, SourceName: "eastmoney_fund", SourceLevel: model.SourceLevelB, SourceType: "fund_nav", TradeDate: tradeDate, CapturedAt: workflowNowRFC3339(), ContentHash: "sha256:" + stableHash("eastmoney_fund", code, tradeDate, strconv.FormatFloat(trend.NAV, 'f', -1, 64)), Metadata: metadata}, nil
}

func eastmoneyNAVHistory(raw []byte) []map[string]any {
	trend, err := eastmoneyNetWorthTrendItems(raw)
	if err != nil {
		return nil
	}
	accumulated := eastmoneyAccumulatedNAVByTimestamp(raw)
	chinaTime, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil
	}
	items := make([]map[string]any, 0, len(trend))
	for _, item := range trend {
		row := map[string]any{"trade_date": time.UnixMilli(item.Timestamp).In(chinaTime).Format("2006-01-02"), "nav": item.NAV, "equity_return": item.EquityReturn, "unit_money": item.UnitMoney}
		if value := accumulated[item.Timestamp]; value > 0 {
			row["accumulated_nav"] = value
		}
		items = append(items, row)
	}
	return items
}

func eastmoneyNetWorthTrendItems(raw []byte) ([]eastmoneyTrendItem, error) {
	re := regexp.MustCompile(`var\s+Data_netWorthTrend\s*=\s*(\[[\s\S]*?\]);`)
	match := re.FindSubmatch(raw)
	if len(match) != 2 {
		return nil, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "东方财富基金净值缺失")
	}
	var items []eastmoneyTrendItem
	if err := json.Unmarshal(match[1], &items); err != nil {
		return nil, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "东方财富基金净值解析失败", err)
	}
	if len(items) == 0 {
		return nil, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "东方财富基金净值无有效记录")
	}
	return items, nil
}

func eastmoneyAccumulatedNAVByTimestamp(raw []byte) map[int64]float64 {
	values := map[int64]float64{}
	re := regexp.MustCompile(`var\s+Data_ACWorthTrend\s*=\s*(\[[\s\S]*?\]);`)
	match := re.FindSubmatch(raw)
	if len(match) != 2 {
		return values
	}
	var rows [][]float64
	if err := json.Unmarshal(match[1], &rows); err != nil {
		return values
	}
	for _, row := range rows {
		if len(row) >= 2 {
			values[int64(row[0])] = row[1]
		}
	}
	return values
}

func eastmoneyJSONVar(raw []byte, name string) map[string]any {
	re := regexp.MustCompile(`var\s+` + regexp.QuoteMeta(name) + `\s*=\s*(\{[\s\S]*?\});`)
	match := re.FindSubmatch(raw)
	if len(match) != 2 {
		return nil
	}
	var value map[string]any
	if err := json.Unmarshal(match[1], &value); err != nil {
		return nil
	}
	return value
}

func eastmoneyJSONArrayVar(raw []byte, name string) []map[string]any {
	re := regexp.MustCompile(`var\s+` + regexp.QuoteMeta(name) + `\s*=\s*(\[[\s\S]*?\]);`)
	match := re.FindSubmatch(raw)
	if len(match) != 2 {
		return nil
	}
	var value []map[string]any
	if err := json.Unmarshal(match[1], &value); err != nil {
		return nil
	}
	return value
}

type eastmoneyTrendItem struct {
	Timestamp    int64   `json:"x"`
	NAV          float64 `json:"y"`
	EquityReturn float64 `json:"equityReturn"`
	UnitMoney    string  `json:"unitMoney"`
}

func eastmoneyStringVar(raw []byte, name string) string {
	re := regexp.MustCompile(`var\s+` + regexp.QuoteMeta(name) + `\s*=\s*"([^"]*)"`)
	match := re.FindSubmatch(raw)
	if len(match) != 2 {
		return ""
	}
	return string(match[1])
}

func eastmoneyNetWorthTrend(raw []byte) (eastmoneyTrendItem, error) {
	items, err := eastmoneyNetWorthTrendItems(raw)
	if err != nil {
		return eastmoneyTrendItem{}, err
	}
	last := items[len(items)-1]
	if last.NAV <= 0 || last.Timestamp <= 0 {
		return eastmoneyTrendItem{}, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "东方财富基金净值无有效记录")
	}
	return last, nil
}

func eastmoneyAccumulatedNAV(raw []byte, timestamp int64) float64 {
	re := regexp.MustCompile(`var\s+Data_ACWorthTrend\s*=\s*(\[[\s\S]*?\]);`)
	match := re.FindSubmatch(raw)
	if len(match) != 2 {
		return 0
	}
	var rows [][]float64
	if err := json.Unmarshal(match[1], &rows); err != nil {
		return 0
	}
	for i := len(rows) - 1; i >= 0; i-- {
		if len(rows[i]) >= 2 && int64(rows[i][0]) == timestamp {
			return rows[i][1]
		}
	}
	return 0
}

func (c CsindexCollector) fetchCsindexExtendedMetadata(ctx context.Context, base, symbol string, metadata map[string]any) {
	health := map[string]any{}
	categories := []string{"index_constituents", "index_weights", "index_valuation_files"}
	capturedAt := workflowNowRFC3339()
	officialFiles := make([]map[string]any, 0, 2)

	if items, status := fetchCsindexP34Rows(ctx, c.HTTPClient, base+"/csindex-home/indexInfo/index-consituent/"+symbol); len(items) > 0 {
		metadata["constituents"] = normalizeCsindexRows(items, base)
		health["index_constituents"] = p34SourceHealthMetadata("csindex", string(model.SourceLevelA), "index_basic", "index_constituents", status, capturedAt, "", symbol, pointTradeDate(metadata))
	} else {
		health["index_constituents"] = p34SourceHealthMetadata("csindex", string(model.SourceLevelA), "index_basic", "index_constituents", status, "", capturedAt, symbol, pointTradeDate(metadata))
	}
	if items, status := fetchCsindexP34Rows(ctx, c.HTTPClient, base+"/csindex-home/indexInfo/index-weight/"+symbol); len(items) > 0 {
		metadata["weights"] = normalizeCsindexRows(items, base)
		health["index_weights"] = p34SourceHealthMetadata("csindex", string(model.SourceLevelA), "index_basic", "index_weights", status, capturedAt, "", symbol, pointTradeDate(metadata))
	} else {
		health["index_weights"] = p34SourceHealthMetadata("csindex", string(model.SourceLevelA), "index_basic", "index_weights", status, "", capturedAt, symbol, pointTradeDate(metadata))
	}
	if items, status := fetchCsindexP34Rows(ctx, c.HTTPClient, base+"/csindex-home/indexInfo/index-valuation/"+symbol); len(items) > 0 {
		metadata["valuation_files"] = normalizeCsindexRows(items, base)
		health["index_valuation_files"] = p34SourceHealthMetadata("csindex", string(model.SourceLevelA), "index_basic", "index_valuation_files", status, capturedAt, "", symbol, pointTradeDate(metadata))
	} else {
		health["index_valuation_files"] = p34SourceHealthMetadata("csindex", string(model.SourceLevelA), "index_basic", "index_valuation_files", status, "", capturedAt, symbol, pointTradeDate(metadata))
	}
	constituentsNeedsFallback := p34CollectorHealthFreshness(health, "index_constituents") == "source_unavailable"
	weightsNeedsFallback := p34CollectorHealthFreshness(health, "index_weights") == "source_unavailable"
	if constituentsNeedsFallback || weightsNeedsFallback {
		if file, ok := c.fetchCsindexOfficialFile(ctx, base, symbol, "closeweight", "/static/html/csindex/public/uploads/file/autofile/closeweight/"+symbol+"closeweight.xls"); ok {
			officialFiles = append(officialFiles, file)
			metadata["constituents_file"] = file
			metadata["weights_file"] = file
			if constituentsNeedsFallback {
				health["index_constituents"] = p34SourceHealthMetadata("csindex", string(model.SourceLevelA), "index_basic", "index_constituents", "fresh", capturedAt, "", symbol, pointTradeDate(metadata))
			}
			if weightsNeedsFallback {
				health["index_weights"] = p34SourceHealthMetadata("csindex", string(model.SourceLevelA), "index_basic", "index_weights", "fresh", capturedAt, "", symbol, pointTradeDate(metadata))
			}
		}
	}
	if p34CollectorHealthFreshness(health, "index_valuation_files") == "source_unavailable" {
		if file, ok := c.fetchCsindexOfficialFile(ctx, base, symbol, "factsheet", "/static/html/csindex/public/uploads/indices/detail/files/zh_CN/"+symbol+"factsheet.pdf"); ok {
			officialFiles = append(officialFiles, file)
			metadata["valuation_file"] = file
			health["index_valuation_files"] = p34SourceHealthMetadata("csindex", string(model.SourceLevelA), "index_basic", "index_valuation_files", "fresh", capturedAt, "", symbol, pointTradeDate(metadata))
		}
	}

	metadata["p34_data_categories"] = categories
	metadata["p34_source_health"] = health
	if len(officialFiles) > 0 {
		metadata["official_files"] = officialFiles
	}
}

func (c CsindexCollector) fetchCsindexOfficialFile(ctx context.Context, base, symbol, category, path string) (map[string]any, bool) {
	url := strings.TrimRight(base, "/") + path
	raw, err := fetchReadonlyURL(ctx, c.HTTPClient, url)
	if err != nil || len(raw) == 0 {
		return nil, false
	}
	return map[string]any{
		"category":    category,
		"symbol":      symbol,
		"url":         url,
		"size_bytes":  len(raw),
		"captured_at": workflowNowRFC3339(),
	}, true
}

func p34CollectorHealthFreshness(health map[string]any, category string) string {
	item, _ := health[category].(map[string]any)
	freshness, _ := item["freshness"].(string)
	return strings.TrimSpace(freshness)
}

func p34SourceHealthMetadata(sourceName, sourceLevel, sourceType, category, freshness, lastSuccessAt, lastFailureAt, symbol, dataDate string) map[string]any {
	item := map[string]any{"freshness": freshness, "source_name": sourceName, "source_level": sourceLevel, "source_type": sourceType, "data_category": category, "affected_symbols": []string{symbol}}
	if dataDate != "" {
		item["data_date"] = dataDate
	}
	if lastSuccessAt != "" {
		item["last_success_at"] = lastSuccessAt
	}
	if lastFailureAt != "" {
		item["last_failure_at"] = lastFailureAt
		item["failure_category"] = freshness
	}
	return item
}

func pointTradeDate(metadata map[string]any) string {
	tradeDate, _ := metadata["trade_date"].(string)
	if tradeDate != "" {
		return tradeDate
	}
	publishDate, _ := metadata["publish_date"].(string)
	return publishDate
}

func fetchCsindexP34Rows(ctx context.Context, client *http.Client, endpoint string) ([]map[string]any, string) {
	items, err := fetchCsindexRows(ctx, client, endpoint)
	if err != nil {
		if strings.Contains(err.Error(), "解析失败") {
			return nil, "parse_error"
		}
		return nil, "source_unavailable"
	}
	if len(items) == 0 {
		return nil, "no_data"
	}
	return items, "fresh"
}

type csindexRowsResponse struct {
	Data []map[string]any `json:"data"`
}

func fetchCsindexRows(ctx context.Context, client *http.Client, endpoint string) ([]map[string]any, error) {
	raw, err := fetchReadonlyURL(ctx, client, endpoint)
	if err != nil {
		return nil, err
	}
	var resp csindexRowsResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "中证指数扩展数据解析失败", err)
	}
	return resp.Data, nil
}

func normalizeCsindexRows(rows []map[string]any, base string) []map[string]any {
	out := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		item := map[string]any{}
		for key, value := range row {
			item[toSnakeCase(key)] = value
		}
		if fileURL, ok := item["file_url"].(string); ok && strings.HasPrefix(fileURL, "/") {
			item["file_url"] = strings.TrimRight(base, "/") + fileURL
		}
		out = append(out, item)
	}
	return out
}

func toSnakeCase(value string) string {
	var b strings.Builder
	for i, r := range value {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('_')
		}
		b.WriteRune(r)
	}
	return strings.ToLower(b.String())
}

type csindexResponse struct {
	Data csindexData `json:"data"`
}

type csindexData struct {
	IndexCode        string  `json:"indexCode"`
	IndexName        string  `json:"indexName"`
	IndexNameEn      string  `json:"indexNameEn"`
	IndexShortNameCn string  `json:"indexShortNameCn"`
	IndexShortNameEn string  `json:"indexShortNameEn"`
	IndexFullNameCn  string  `json:"indexFullNameCn"`
	IndexFullNameEn  string  `json:"indexFullNameEn"`
	LatestClose      float64 `json:"latestClose"`
	Currency         string  `json:"currency"`
	CurrencyCn       string  `json:"currencyCn"`
	CurrencyEn       string  `json:"currencyEn"`
	PublishDate      string  `json:"publishDate"`
	BasicDate        string  `json:"basicDate"`
	BasicIndex       float64 `json:"basicIndex"`
	ConsNumber       float64 `json:"consNumber"`
	PublishChannelCn string  `json:"publishChannelCn"`
	IndexType        string  `json:"indexType"`
	AdjFreqCn        string  `json:"adjFreqCn"`
}

func parseCsindexPoint(raw []byte) (MarketDataPoint, error) {
	var resp csindexResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return MarketDataPoint{}, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "中证指数数据解析失败", err)
	}
	if strings.TrimSpace(resp.Data.IndexCode) == "" {
		return MarketDataPoint{}, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "中证指数数据缺少指数代码")
	}
	indexName := firstNonEmptyMarketString(resp.Data.IndexName, resp.Data.IndexShortNameCn, resp.Data.IndexFullNameCn)
	indexNameEn := firstNonEmptyMarketString(resp.Data.IndexNameEn, resp.Data.IndexShortNameEn, resp.Data.IndexFullNameEn)
	currency := firstNonEmptyMarketString(resp.Data.Currency, resp.Data.CurrencyEn, resp.Data.CurrencyCn)
	metadata := map[string]any{"index_code": resp.Data.IndexCode, "symbol": resp.Data.IndexCode, "index_name": indexName, "index_name_en": indexNameEn, "currency": currency, "publish_date": resp.Data.PublishDate, "raw": map[string]any{"data": resp.Data}}
	if resp.Data.IndexFullNameCn != "" {
		metadata["index_full_name"] = resp.Data.IndexFullNameCn
	}
	if resp.Data.BasicDate != "" {
		metadata["basic_date"] = resp.Data.BasicDate
	}
	if resp.Data.BasicIndex > 0 {
		metadata["basic_index"] = resp.Data.BasicIndex
	}
	if resp.Data.ConsNumber > 0 {
		metadata["constituent_count"] = resp.Data.ConsNumber
	}
	if resp.Data.PublishChannelCn != "" {
		metadata["publish_channel"] = resp.Data.PublishChannelCn
	}
	if resp.Data.IndexType != "" {
		metadata["index_type"] = resp.Data.IndexType
	}
	if resp.Data.AdjFreqCn != "" {
		metadata["adjustment_frequency"] = resp.Data.AdjFreqCn
	}
	return MarketDataPoint{ClosePrice: resp.Data.LatestClose, SourceName: "csindex", SourceLevel: model.SourceLevelA, SourceType: "index_basic", TradeDate: resp.Data.PublishDate, CapturedAt: workflowNowRFC3339(), ContentHash: "sha256:" + stableHash("csindex", resp.Data.IndexCode, resp.Data.PublishDate, strconv.FormatFloat(resp.Data.LatestClose, 'f', -1, 64), indexName), Metadata: metadata}, nil
}

func firstNonEmptyMarketString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
