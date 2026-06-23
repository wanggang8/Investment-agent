package workflow

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"investment-agent/internal/domain/analyst"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/pkg/apperr"
)

// MarketDataPoint 是行情数据源返回的标准化市场点位。
type MarketDataPoint struct {
	PEPercentile         float64              `json:"pe_percentile"`
	PBPercentile         float64              `json:"pb_percentile"`
	VolumePercentile     float64              `json:"volume_percentile"`
	VolatilityPercentile float64              `json:"volatility_percentile"`
	ClosePrice           float64              `json:"close_price"`
	TurnoverRate         float64              `json:"turnover_rate"`
	LiquidityState       model.LiquidityState `json:"liquidity_state"`
	SentimentState       model.SentimentState `json:"sentiment_state"`
	Stale                bool                 `json:"stale"`
	SourceName           string               `json:"source_name"`
	SourceLevel          model.SourceLevel    `json:"source_level"`
	SourceType           string               `json:"source_type"`
	TradeDate            string               `json:"trade_date"`
	CapturedAt           string               `json:"captured_at"`
	ContentHash          string               `json:"content_hash"`
	Metadata             map[string]any       `json:"metadata"`
}

// MarketDataSource 隔离真实行情供应商和本地 stub，业务层只依赖标准化结果。
type MarketDataSource interface {
	FetchMarketData(ctx context.Context, symbol string) (MarketDataPoint, error)
}

// SentimentProxyPoint 是 P34 情绪替代指标的 fixture/公开源标准化输入。
type SentimentProxyPoint struct {
	SourceName     string               `json:"source_name"`
	SourceLevel    model.SourceLevel    `json:"source_level"`
	DataDate       string               `json:"data_date"`
	HeatScore      float64              `json:"heat_score"`
	SentimentState model.SentimentState `json:"sentiment_state"`
	Raw            map[string]any       `json:"raw"`
}

// FixtureSentimentProxyCollector 提供 P34 情绪替代指标的确定性 fixture，不依赖公网。
type FixtureSentimentProxyCollector struct {
	Fixtures map[string]SentimentProxyPoint
}

func (c FixtureSentimentProxyCollector) FetchMarketData(ctx context.Context, symbol string) (MarketDataPoint, error) {
	_ = ctx
	symbol = strings.TrimSpace(symbol)
	item, ok := c.Fixtures[symbol]
	if !ok || strings.TrimSpace(item.DataDate) == "" {
		return MarketDataPoint{}, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "P34 情绪替代指标无可用数据")
	}
	sourceName := strings.TrimSpace(item.SourceName)
	if sourceName == "" {
		sourceName = "sentiment_proxy_fixture"
	}
	sourceLevel := item.SourceLevel
	if sourceLevel == "" {
		sourceLevel = model.SourceLevelC
	}
	sentiment := item.SentimentState
	if sentiment == "" {
		sentiment = model.SentimentNeutral
	}
	metadata := map[string]any{
		"data_category":       "sentiment_proxy",
		"heat_score":          item.HeatScore,
		"p34_data_categories": []string{"sentiment_proxy"},
		"p34_source_health": map[string]any{"sentiment_proxy": map[string]any{
			"freshness":        "stubbed",
			"last_failure_at":  workflowNowRFC3339(),
			"failure_category": "stubbed",
			"data_date":        item.DataDate,
			"affected_symbols": []string{symbol},
			"source_level":     string(sourceLevel),
			"source_type":      "sentiment_proxy",
		}},
	}
	if item.Raw != nil {
		metadata["raw"] = item.Raw
	}
	return MarketDataPoint{SourceName: sourceName, SourceLevel: sourceLevel, SourceType: "sentiment_proxy", TradeDate: item.DataDate, CapturedAt: workflowNowRFC3339(), ContentHash: "sha256:" + stableHash(sourceName, symbol, item.DataDate, strconv.FormatFloat(item.HeatScore, 'f', -1, 64)), SentimentState: sentiment, Metadata: metadata}, nil
}

// IntelligenceSourceItem 是情报数据源返回的标准化原始材料。
type IntelligenceSourceItem struct {
	SourceName  string            `json:"source_name"`
	SourceLevel model.SourceLevel `json:"source_level"`
	Title       string            `json:"title"`
	Text        string            `json:"text"`
	URL         string            `json:"url"`
	PublishedAt string            `json:"published_at"`
}

// IntelligenceSource 隔离新闻、公告、手工导入等情报来源。
type IntelligenceSource interface {
	FetchIntelligence(ctx context.Context, symbol string) ([]IntelligenceSourceItem, error)
}

// PublicEvidencePayload 是 P26 公告与证据源 collector 返回的标准化 payload。
type PublicEvidencePayload struct {
	SourceName      string            `json:"source_name"`
	SourceLevel     model.SourceLevel `json:"source_level"`
	SourceType      string            `json:"source_type"`
	EvidenceRole    string            `json:"evidence_role"`
	Symbol          string            `json:"symbol"`
	SourceRecordID  string            `json:"source_record_id"`
	Title           string            `json:"title"`
	Text            string            `json:"text"`
	URL             string            `json:"url"`
	AttachmentURL   string            `json:"attachment_url"`
	PublishedAt     string            `json:"published_at"`
	CapturedAt      time.Time         `json:"captured_at"`
	TimeWeight      float64           `json:"time_weight,omitempty"`
	ContentHash     string            `json:"content_hash"`
	SourceRecordKey string            `json:"source_record_key"`
	Raw             map[string]any    `json:"raw"`
}

// PublicEvidenceCollector 定义证据源 collector 接口（P26）。
type PublicEvidenceCollector interface {
	FetchPublicEvidence(ctx context.Context, symbol string, start, end time.Time) ([]PublicEvidencePayload, error)
}

// CompositePublicEvidenceCollector 按源顺序调用多个 collector，单源失败时降级并继续其他源。
type CompositePublicEvidenceCollector struct {
	Collectors []PublicEvidenceCollector
	Failures   []PublicEvidenceError
}

func (c *CompositePublicEvidenceCollector) FetchPublicEvidence(ctx context.Context, symbol string, start, end time.Time) ([]PublicEvidencePayload, error) {
	var all []PublicEvidencePayload
	var lastErr error
	allErrorsAreNoData := len(c.Collectors) > 0
	c.Failures = nil
	for _, collector := range c.Collectors {
		items, err := collector.FetchPublicEvidence(ctx, symbol, start, end)
		if err != nil {
			lastErr = err
			if sourceErr, ok := publicEvidenceErrorOf(err); ok {
				c.Failures = append(c.Failures, sourceErr)
				if sourceErr.ErrorCode != "no_data" {
					allErrorsAreNoData = false
				}
			} else {
				allErrorsAreNoData = false
			}
			continue
		}
		allErrorsAreNoData = false
		all = append(all, items...)
	}
	if len(all) == 0 {
		if lastErr != nil {
			if allErrorsAreNoData {
				return []PublicEvidencePayload{}, nil
			}
			return nil, lastErr
		}
		return nil, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "证据数据源无可用材料")
	}
	return NormalizePublicEvidenceItems(all)
}

func (c *CompositePublicEvidenceCollector) PublicEvidenceFailures() []PublicEvidenceError {
	return c.Failures
}

// FixturePublicEvidenceCollector 提供测试和示例用的固定 payload。
type FixturePublicEvidenceCollector struct {
	Fixtures map[string][]PublicEvidencePayload
}

func (f FixturePublicEvidenceCollector) FetchPublicEvidence(ctx context.Context, symbol string, start, end time.Time) ([]PublicEvidencePayload, error) {
	items, ok := f.Fixtures[symbol]
	if !ok || len(items) == 0 {
		return nil, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "证据数据源无可用材料")
	}
	return NormalizePublicEvidenceItems(items)
}

// NormalizePublicEvidenceItems 对 P26 collector 输出去重、hash 和校验必填字段。
func NormalizePublicEvidenceItems(payloads []PublicEvidencePayload) ([]PublicEvidencePayload, error) {
	seen := make(map[string]bool, len(payloads))
	result := make([]PublicEvidencePayload, 0, len(payloads))

	for _, p := range payloads {
		p.SourceName = strings.TrimSpace(p.SourceName)
		p.SourceType = strings.TrimSpace(p.SourceType)
		p.Symbol = strings.TrimSpace(p.Symbol)
		p.SourceRecordID = strings.TrimSpace(p.SourceRecordID)
		p.Title = strings.TrimSpace(p.Title)
		p.Text = normalizeObjectiveEvidenceText(strings.TrimSpace(p.Text))
		p.URL = strings.TrimSpace(p.URL)
		p.AttachmentURL = strings.TrimSpace(p.AttachmentURL)
		p.PublishedAt = strings.TrimSpace(p.PublishedAt)

		// 验证必填字段
		if p.SourceName == "" || p.Title == "" || p.Text == "" {
			return nil, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "证据数据缺少来源、标题或正文")
		}
		if !p.SourceLevel.Valid() {
			return nil, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "证据数据缺少有效信源等级")
		}
		role := model.EvidenceRole(strings.TrimSpace(p.EvidenceRole))
		if !role.Valid() {
			return nil, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "证据数据缺少有效证据角色")
		}
		if p.SourceLevel == model.SourceLevelC {
			role = model.EvidenceBackground
		}
		p.EvidenceRole = string(role)

		// 计算 source record key：优先用 source_record_id，回退到 URL+PublishedAt
		sourceRecordKey := ""
		if p.SourceRecordID != "" {
			sourceRecordKey = p.SourceName + "\x00" + p.SourceRecordID
		} else if p.URL != "" {
			sourceRecordKey = p.URL + "\x00" + p.PublishedAt
		} else {
			sourceRecordKey = p.SourceName + "\x00" + p.Title + "\x00" + p.PublishedAt
		}

		// 去重：相同 source record key 只保留第一条
		if seen[sourceRecordKey] {
			continue
		}
		seen[sourceRecordKey] = true

		// 计算 content hash
		contentHash := stableHash(p.SourceName, p.Title, p.Text, p.URL, p.PublishedAt)

		p.SourceRecordKey = sourceRecordKey
		p.ContentHash = "sha256:" + contentHash
		result = append(result, p)
	}

	if len(result) == 0 {
		return nil, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "证据数据源未返回有效材料")
	}

	return result, nil
}

func normalizeObjectiveEvidenceText(text string) string {
	replacements := []struct {
		from string
		to   string
	}{
		{from: "恐慌踩踏", to: "客观化情绪描述: 投资者情绪显著偏弱"},
		{from: "赶紧卖出", to: "客观化情绪描述: 出现卖出倾向表达"},
		{from: "赶紧卖", to: "客观化情绪描述: 出现卖出倾向表达"},
		{from: "爆雷", to: "客观化情绪描述: 负面风险事件被市场讨论"},
		{from: "恐慌", to: "客观化情绪描述: 投资者情绪偏弱"},
		{from: "崩盘", to: "客观化情绪描述: 价格大幅下跌表述"},
		{from: "暴跌", to: "客观化情绪描述: 价格快速下跌表述"},
		{from: "暴涨", to: "客观化情绪描述: 价格快速上涨表述"},
		{from: "大利空", to: "客观化情绪描述: 重大负面信息表述"},
	}
	for _, replacement := range replacements {
		text = strings.ReplaceAll(text, replacement.from, replacement.to)
	}
	return text
}

type AnalystRequest = analyst.Request

type AnalystResponse = analyst.Response

type AnalystService = analyst.Service

type ConfiguredMarketDataSource struct {
	Enabled           []string
	MarketEndpoint    string
	HTTPClient        *http.Client
	MarketFixtures    map[string]MarketDataPoint
	RawMarketFixtures map[string]string
	Fallback          MarketDataSource
	Timeout           time.Duration
	FetchDelay        time.Duration
}

func publicHTTPEndpointEnabled(enabled []string) bool {
	if len(enabled) == 0 {
		return true
	}
	for _, item := range enabled {
		source := strings.TrimSpace(item)
		if source != "fixture" && source != "stub" {
			return true
		}
	}
	return false
}

func (s ConfiguredMarketDataSource) FetchMarketData(ctx context.Context, symbol string) (MarketDataPoint, error) {
	symbol = strings.TrimSpace(symbol)
	if symbol == "" {
		return MarketDataPoint{}, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "市场数据源不可用")
	}
	if err := waitConfiguredSource(ctx, s.FetchDelay, s.Timeout); err != nil {
		return MarketDataPoint{}, err
	}
	if raw, ok := s.RawMarketFixtures[symbol]; ok {
		var point MarketDataPoint
		if err := json.Unmarshal([]byte(raw), &point); err != nil {
			return MarketDataPoint{}, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "市场数据源解析失败", err)
		}
		return validateMarketPoint(point)
	}
	if point, ok := s.MarketFixtures[symbol]; ok {
		return validateMarketPoint(point)
	}
	if publicHTTPEndpointEnabled(s.Enabled) && strings.TrimSpace(s.MarketEndpoint) != "" {
		raw, err := fetchReadonlyBytes(ctx, s.HTTPClient, s.MarketEndpoint, symbol)
		if err != nil {
			return s.fetchMarketFallback(ctx, symbol, err)
		}
		point, err := parseMarketDataPoint(raw)
		if err != nil {
			return s.fetchMarketFallback(ctx, symbol, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "数据源解析失败", err))
		}
		point, err = validateMarketPoint(point)
		if err != nil {
			return s.fetchMarketFallback(ctx, symbol, err)
		}
		return point, nil
	}
	return s.fetchMarketFallback(ctx, symbol, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "真实市场数据源未返回可用数据"))
}

func (s ConfiguredMarketDataSource) fetchMarketFallback(ctx context.Context, symbol string, primaryErr error) (MarketDataPoint, error) {
	if s.Fallback == nil {
		return MarketDataPoint{}, primaryErr
	}
	point, err := s.Fallback.FetchMarketData(ctx, symbol)
	if err != nil {
		return MarketDataPoint{}, primaryErr
	}
	return validateMarketPoint(point)
}

type ConfiguredIntelligenceSource struct {
	Enabled                 []string
	IntelligenceEndpoint    string
	HTTPClient              *http.Client
	IntelligenceFixtures    map[string][]IntelligenceSourceItem
	RawIntelligenceFixtures map[string]string
	Timeout                 time.Duration
	FetchDelay              time.Duration
}

func (s ConfiguredIntelligenceSource) FetchIntelligence(ctx context.Context, symbol string) ([]IntelligenceSourceItem, error) {
	symbol = strings.TrimSpace(symbol)
	if symbol == "" {
		return nil, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "情报数据源不可用")
	}
	if err := waitConfiguredSource(ctx, s.FetchDelay, s.Timeout); err != nil {
		return nil, err
	}
	if raw, ok := s.RawIntelligenceFixtures[symbol]; ok {
		var items []IntelligenceSourceItem
		if err := json.Unmarshal([]byte(raw), &items); err != nil {
			return nil, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "情报数据源解析失败", err)
		}
		return validateIntelligenceItems(items)
	}
	if items, ok := s.IntelligenceFixtures[symbol]; ok {
		return validateIntelligenceItems(items)
	}
	if publicHTTPEndpointEnabled(s.Enabled) && strings.TrimSpace(s.IntelligenceEndpoint) != "" {
		raw, err := fetchReadonlyBytes(ctx, s.HTTPClient, s.IntelligenceEndpoint, symbol)
		if err != nil {
			return nil, err
		}
		items, err := parseIntelligenceItems(raw)
		if err != nil {
			return nil, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "数据源解析失败", err)
		}
		return validateIntelligenceItems(items)
	}
	return nil, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "真实情报数据源未返回可用材料")
}

type publicMarketPayload struct {
	Data      publicMarketDataList `json:"data"`
	Market    publicMarketDataList `json:"market"`
	Quote     publicMarketDataList `json:"quote"`
	Valuation publicValuation      `json:"valuation"`
}

type publicMarketDataList []publicMarketData

func (l *publicMarketDataList) UnmarshalJSON(raw []byte) error {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	var items []publicMarketData
	if err := json.Unmarshal(raw, &items); err == nil {
		*l = items
		return nil
	}
	var item publicMarketData
	if err := json.Unmarshal(raw, &item); err != nil {
		return err
	}
	*l = []publicMarketData{item}
	return nil
}

type publicIntelligencePayload struct {
	Items []publicIntelligenceItem   `json:"items"`
	Data  publicIntelligenceItemList `json:"data"`
}

type publicIntelligenceItemList []publicIntelligenceItem

func (l *publicIntelligenceItemList) UnmarshalJSON(raw []byte) error {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	var items []publicIntelligenceItem
	if err := json.Unmarshal(raw, &items); err == nil {
		*l = items
		return nil
	}
	var item publicIntelligenceItem
	if err := json.Unmarshal(raw, &item); err != nil {
		return err
	}
	*l = []publicIntelligenceItem{item}
	return nil
}

type publicIntelligenceItem struct {
	SourceName  string            `json:"source_name"`
	Source      string            `json:"source"`
	SourceLevel model.SourceLevel `json:"source_level"`
	Title       string            `json:"title"`
	Text        string            `json:"text"`
	Content     string            `json:"content"`
	Summary     string            `json:"summary"`
	URL         string            `json:"url"`
	OriginalURL string            `json:"original_url"`
	PublishedAt string            `json:"published_at"`
}

func parseIntelligenceItems(raw []byte) ([]IntelligenceSourceItem, error) {
	var publicItems []publicIntelligenceItem
	if err := json.Unmarshal(raw, &publicItems); err == nil && len(publicItems) > 0 {
		items := make([]IntelligenceSourceItem, 0, len(publicItems))
		for _, item := range publicItems {
			items = append(items, publicIntelligenceSourceItem(item))
		}
		return dedupeIntelligenceItems(items), nil
	}
	var items []IntelligenceSourceItem
	if err := json.Unmarshal(raw, &items); err == nil && len(items) > 0 {
		return dedupeIntelligenceItems(items), nil
	}
	var payload publicIntelligencePayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}
	payloadItems := payload.Items
	if len(payloadItems) == 0 {
		payloadItems = payload.Data
	}
	items = make([]IntelligenceSourceItem, 0, len(payloadItems))
	for _, item := range payloadItems {
		items = append(items, publicIntelligenceSourceItem(item))
	}
	return dedupeIntelligenceItems(items), nil
}

func publicIntelligenceSourceItem(item publicIntelligenceItem) IntelligenceSourceItem {
	sourceName := item.SourceName
	if sourceName == "" {
		sourceName = item.Source
	}
	text := item.Text
	if text == "" {
		text = item.Content
	}
	if text == "" {
		text = item.Summary
	}
	itemURL := item.URL
	if itemURL == "" {
		itemURL = item.OriginalURL
	}
	return IntelligenceSourceItem{SourceName: sourceName, SourceLevel: publicSourceLevel(sourceName, item.SourceLevel), Title: item.Title, Text: text, URL: itemURL, PublishedAt: item.PublishedAt}
}

func publicSourceLevel(sourceName string, level model.SourceLevel) model.SourceLevel {
	if level != "" {
		return level
	}
	normalized := strings.ToLower(strings.TrimSpace(sourceName))
	switch normalized {
	case "exchange_disclosure", "fund_company_disclosure", "official_fund_data":
		return model.SourceLevelA
	case "financial_news":
		return model.SourceLevelB
	}
	if strings.Contains(sourceName, "交易所") || strings.Contains(sourceName, "基金公司公告") || strings.Contains(sourceName, "巨潮资讯") || strings.Contains(sourceName, "上交所") || strings.Contains(sourceName, "深交所") {
		return model.SourceLevelA
	}
	return ""
}

func dedupeIntelligenceItems(items []IntelligenceSourceItem) []IntelligenceSourceItem {
	seen := make(map[string]bool, len(items))
	result := make([]IntelligenceSourceItem, 0, len(items))
	for _, item := range items {
		key := strings.TrimSpace(item.URL)
		if key == "" {
			key = strings.TrimSpace(item.SourceName) + "\x00" + strings.TrimSpace(item.Title) + "\x00" + strings.TrimSpace(item.PublishedAt)
		}
		if key != "" && seen[key] {
			continue
		}
		if key != "" {
			seen[key] = true
		}
		result = append(result, item)
	}
	return result
}

type publicMarketData struct {
	Symbol               string               `json:"symbol"`
	TradeDate            string               `json:"trade_date"`
	Close                float64              `json:"close"`
	ClosePrice           float64              `json:"close_price"`
	NAV                  float64              `json:"nav"`
	NetValue             float64              `json:"net_value"`
	TurnoverRate         float64              `json:"turnover_rate"`
	PEPercentile         float64              `json:"pe_percentile"`
	PBPercentile         float64              `json:"pb_percentile"`
	VolumePercentile     float64              `json:"volume_percentile"`
	VolatilityPercentile float64              `json:"volatility_percentile"`
	Liquidity            model.LiquidityState `json:"liquidity"`
	LiquidityState       model.LiquidityState `json:"liquidity_state"`
	Sentiment            model.SentimentState `json:"sentiment"`
	SentimentState       model.SentimentState `json:"sentiment_state"`
	Stale                bool                 `json:"stale"`
	Valuation            publicValuation      `json:"valuation"`
}

type publicValuation struct {
	PEPercentile float64 `json:"pe_percentile"`
	PBPercentile float64 `json:"pb_percentile"`
}

func parseMarketDataPoint(raw []byte) (MarketDataPoint, error) {
	var point MarketDataPoint
	if err := json.Unmarshal(raw, &point); err != nil {
		return MarketDataPoint{}, err
	}
	if point.ClosePrice > 0 || point.Stale {
		return point, nil
	}
	var payload publicMarketPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return MarketDataPoint{}, err
	}
	data := firstPublicMarketData(payload.Data, payload.Market, payload.Quote)
	valuation := data.Valuation
	if valuation.PEPercentile == 0 && valuation.PBPercentile == 0 {
		valuation = payload.Valuation
	}
	return publicMarketDataPoint(data, valuation), nil
}

func firstPublicMarketData(lists ...publicMarketDataList) publicMarketData {
	for _, list := range lists {
		for _, item := range list {
			if item.Close > 0 || item.ClosePrice > 0 || item.NAV > 0 || item.NetValue > 0 || item.Stale {
				return item
			}
		}
	}
	return publicMarketData{}
}

func publicMarketDataPoint(data publicMarketData, valuation publicValuation) MarketDataPoint {
	closePrice := data.ClosePrice
	if closePrice == 0 {
		closePrice = data.Close
	}
	if closePrice == 0 {
		closePrice = data.NAV
	}
	if closePrice == 0 {
		closePrice = data.NetValue
	}
	liquidity := data.LiquidityState
	if liquidity == "" {
		liquidity = data.Liquidity
	}
	sentiment := data.SentimentState
	if sentiment == "" {
		sentiment = data.Sentiment
	}
	pePercentile := data.PEPercentile
	if pePercentile == 0 {
		pePercentile = valuation.PEPercentile
	}
	pbPercentile := data.PBPercentile
	if pbPercentile == 0 {
		pbPercentile = valuation.PBPercentile
	}
	return MarketDataPoint{PEPercentile: pePercentile, PBPercentile: pbPercentile, VolumePercentile: data.VolumePercentile, VolatilityPercentile: data.VolatilityPercentile, ClosePrice: closePrice, TurnoverRate: data.TurnoverRate, LiquidityState: liquidity, SentimentState: sentiment, Stale: data.Stale}
}

func fetchReadonlyBytes(ctx context.Context, client *http.Client, endpoint, symbol string) ([]byte, error) {
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "数据源地址不可用", err)
	}
	query := parsed.Query()
	query.Set("symbol", symbol)
	parsed.RawQuery = query.Encode()
	return fetchReadonlyURL(ctx, client, parsed.String())
}

func fetchReadonlyURL(ctx context.Context, client *http.Client, endpoint string) ([]byte, error) {
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "数据源地址不可用", err)
	}
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsed.String(), nil)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "数据源请求不可用", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "数据源请求失败", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, fmt.Sprintf("数据源返回状态 %d", resp.StatusCode))
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "数据源读取失败", err)
	}
	return raw, nil
}

func waitConfiguredSource(ctx context.Context, delay, timeout time.Duration) error {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if delay <= 0 {
		return nil
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "数据源请求超时", ctx.Err())
	case <-timer.C:
		return nil
	}
}

func validateMarketPoint(point MarketDataPoint) (MarketDataPoint, error) {
	if point.Stale {
		return MarketDataPoint{}, apperr.New(apperr.CodeDataStale, apperr.CategoryInvalidState, "市场数据已过期")
	}
	if point.ClosePrice <= 0 && point.SourceType != "structured_public_fields" {
		return MarketDataPoint{}, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "市场数据缺少有效价格")
	}
	if point.LiquidityState == "" {
		point.LiquidityState = model.LiquidityNormal
	}
	if point.SentimentState == "" {
		point.SentimentState = model.SentimentNeutral
	}
	return point, nil
}

func validateIntelligenceItems(items []IntelligenceSourceItem) ([]IntelligenceSourceItem, error) {
	if len(items) == 0 {
		return nil, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "情报数据源未返回材料")
	}
	for i := range items {
		items[i].SourceName = strings.TrimSpace(items[i].SourceName)
		items[i].Title = strings.TrimSpace(items[i].Title)
		items[i].Text = strings.TrimSpace(items[i].Text)
		if items[i].SourceName == "" || items[i].Title == "" || items[i].Text == "" {
			return nil, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "情报数据缺少来源、标题或正文")
		}
		if items[i].SourceLevel == "" {
			items[i].SourceLevel = model.SourceLevelB
		}
	}
	return items, nil
}

// StubMarketDataSource 为本地开发提供可预测行情，避免测试依赖外部网络。
type StubMarketDataSource struct{}

func (StubMarketDataSource) FetchMarketData(_ context.Context, symbol string) (MarketDataPoint, error) {
	if strings.TrimSpace(symbol) == "" {
		return MarketDataPoint{}, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "市场数据源不可用")
	}
	return MarketDataPoint{PEPercentile: 50, PBPercentile: 50, VolumePercentile: 50, VolatilityPercentile: 20, ClosePrice: 1, TurnoverRate: 1, LiquidityState: model.LiquidityNormal, SentimentState: model.SentimentNeutral}, nil
}

func stableHash(parts ...string) string {
	h := sha256.New()
	for _, part := range parts {
		_, _ = h.Write([]byte(part))
		_, _ = h.Write([]byte("\x00"))
	}
	return hex.EncodeToString(h.Sum(nil))
}

// StubIntelligenceSource 为本地开发提供可索引情报材料，不写入真实来源凭证。
type StubIntelligenceSource struct{}

func (StubIntelligenceSource) FetchIntelligence(_ context.Context, symbol string) ([]IntelligenceSourceItem, error) {
	if strings.TrimSpace(symbol) == "" {
		return nil, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "情报数据源不可用")
	}
	return []IntelligenceSourceItem{{SourceName: "stub", SourceLevel: model.SourceLevelA, Title: "本地情报", Text: "本地 stub 情报摘要", URL: "stub://local", PublishedAt: "2026-01-01T00:00:00Z"}}, nil
}

// StaticAnalystService 提供本地分析材料；真实 DeepSeek 客户端实现同一接口。
type StaticAnalystService struct{}

func (StaticAnalystService) Analyze(_ context.Context, req AnalystRequest) (AnalystResponse, error) {
	if strings.TrimSpace(req.Symbol) == "" {
		return AnalystResponse{}, apperr.New(apperr.CodeAnalystUnavailable, apperr.CategoryInternal, "分析服务不可用")
	}
	return AnalystResponse{Reports: map[string]string{
		"value":           "估值与基本面分析材料",
		"trend_risk":      "趋势与风险分析材料",
		"expected_return": "预期收益分析材料",
	}}, nil
}
