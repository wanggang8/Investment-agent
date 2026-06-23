package workflow

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/pkg/apperr"
)

func TestPublicEvidencePayloadNormalizesHashAndDedupes(t *testing.T) {
	capturedAt := time.Date(2026, 6, 5, 13, 0, 0, 0, time.UTC)
	items, err := NormalizePublicEvidenceItems([]PublicEvidencePayload{
		{SourceName: "cninfo", SourceLevel: model.SourceLevelA, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "ann-1", Title: "ETF 公告", Text: "公告正文", URL: "https://example.invalid/a", AttachmentURL: "https://example.invalid/a.pdf", PublishedAt: "2026-06-05T00:00:00+08:00", CapturedAt: capturedAt, Raw: map[string]any{"announcementId": "ann-1"}},
		{SourceName: "cninfo", SourceLevel: model.SourceLevelA, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "ann-1", Title: "ETF 公告重复", Text: "重复正文", URL: "https://example.invalid/a", PublishedAt: "2026-06-05T00:00:00+08:00", CapturedAt: capturedAt},
		{SourceName: "szse", SourceLevel: model.SourceLevelA, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", Title: "深交所公告", Text: "深交所正文", URL: "https://example.invalid/b", PublishedAt: "2026-06-05T01:00:00+08:00", CapturedAt: capturedAt},
	})

	if err != nil {
		t.Fatalf("NormalizePublicEvidenceItems: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected duplicate source record to be removed, got %+v", items)
	}
	if items[0].ContentHash == "" || !strings.HasPrefix(items[0].ContentHash, "sha256:") {
		t.Fatalf("expected sha256 content hash, got %+v", items[0])
	}
	if items[0].SourceRecordKey != "cninfo\x00ann-1" || items[0].AttachmentURL != "https://example.invalid/a.pdf" || items[0].Raw == nil {
		t.Fatalf("expected source identity and raw metadata preserved: %+v", items[0])
	}
	if items[1].SourceRecordKey == "" || items[1].ContentHash == items[0].ContentHash {
		t.Fatalf("expected fallback identity and distinct content hash: %+v", items)
	}
}

func TestPublicEvidencePayloadRejectsInvalidRequiredFields(t *testing.T) {
	_, err := NormalizePublicEvidenceItems([]PublicEvidencePayload{{SourceName: "cninfo", SourceLevel: model.SourceLevelA, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", Title: "缺正文"}})
	if err == nil {
		t.Fatal("expected invalid payload error")
	}
	var appErr *apperr.AppError
	if !errors.As(err, &appErr) || appErr.Code != apperr.CodeDataSourceUnavailable {
		t.Fatalf("unexpected error: %T %v", err, err)
	}
}

func TestPublicEvidencePayloadEnforcesSourceMetadataAndFormalBoundary(t *testing.T) {
	base := PublicEvidencePayload{SourceName: "cninfo", SourceLevel: model.SourceLevelA, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "ann-1", Title: "ETF 公告", Text: "公告正文", URL: "https://example.invalid/a", PublishedAt: "2026-06-05T00:00:00+08:00", CapturedAt: time.Date(2026, 6, 5, 13, 0, 0, 0, time.UTC)}

	missingLevel := base
	missingLevel.SourceLevel = ""
	if _, err := NormalizePublicEvidenceItems([]PublicEvidencePayload{missingLevel}); err == nil {
		t.Fatal("expected missing source level to be rejected")
	}

	missingRole := base
	missingRole.EvidenceRole = ""
	if _, err := NormalizePublicEvidenceItems([]PublicEvidencePayload{missingRole}); err == nil {
		t.Fatal("expected missing evidence role to be rejected")
	}

	cFormal := base
	cFormal.SourceLevel = model.SourceLevelC
	cFormal.EvidenceRole = "formal"
	items, err := NormalizePublicEvidenceItems([]PublicEvidencePayload{cFormal})
	if err != nil {
		t.Fatalf("C-level evidence should be retained as background, got %v", err)
	}
	if items[0].EvidenceRole != string(model.EvidenceBackground) {
		t.Fatalf("C-level evidence must not remain formal, got %+v", items[0])
	}
}

func TestPublicEvidencePayloadNormalizesEmotionalDescriptions(t *testing.T) {
	items, err := NormalizePublicEvidenceItems([]PublicEvidencePayload{{
		SourceName:     "financial_news",
		SourceLevel:    model.SourceLevelB,
		SourceType:     "public_disclosure",
		EvidenceRole:   "formal",
		Symbol:         "510300",
		SourceRecordID: "news-emotion",
		Title:          "市场评论",
		Text:           "市场爆雷，投资者恐慌踩踏，赶紧卖出。",
		URL:            "https://example.invalid/news",
		PublishedAt:    "2026-06-05T00:00:00+08:00",
		CapturedAt:     time.Date(2026, 6, 5, 13, 0, 0, 0, time.UTC),
	}})
	if err != nil {
		t.Fatalf("NormalizePublicEvidenceItems: %v", err)
	}
	if strings.Contains(items[0].Text, "爆雷") || strings.Contains(items[0].Text, "恐慌") || strings.Contains(items[0].Text, "赶紧卖出") {
		t.Fatalf("emotional wording should be converted before analysis, got %q", items[0].Text)
	}
	if !strings.Contains(items[0].Text, "客观化情绪描述") {
		t.Fatalf("expected objective conversion marker, got %q", items[0].Text)
	}
}

func TestPublicEvidenceCollectorFetchesFixtures(t *testing.T) {
	capturedAt := time.Date(2026, 6, 5, 13, 0, 0, 0, time.UTC)
	collector := FixturePublicEvidenceCollector{
		Fixtures: map[string][]PublicEvidencePayload{
			"510300": {
				{SourceName: "cninfo", SourceLevel: model.SourceLevelA, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "ann-1", Title: "ETF 公告", Text: "公告正文", URL: "https://example.invalid/a", PublishedAt: "2026-06-05T00:00:00+08:00", CapturedAt: capturedAt},
			},
		},
	}

	items, err := collector.FetchPublicEvidence(context.Background(), "510300", time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("FetchPublicEvidence: %v", err)
	}
	if len(items) != 1 || items[0].SourceName != "cninfo" || items[0].ContentHash == "" {
		t.Fatalf("expected fixture with hash, got %+v", items)
	}
}

func TestPublicEvidenceCollectorReturnsStableFailures(t *testing.T) {
	collector := FixturePublicEvidenceCollector{Fixtures: map[string][]PublicEvidencePayload{}}
	_, err := collector.FetchPublicEvidence(context.Background(), "missing", time.Time{}, time.Time{})

	if err == nil {
		t.Fatal("expected unavailable error")
	}
	var appErr *apperr.AppError
	if !errors.As(err, &appErr) || appErr.Code != apperr.CodeDataSourceUnavailable {
		t.Fatalf("unexpected error: %T %v", err, err)
	}
}

func TestConfiguredMarketDataSourceFetchesReadonlyFixture(t *testing.T) {
	source := ConfiguredMarketDataSource{
		Enabled: []string{"fixture"},
		MarketFixtures: map[string]MarketDataPoint{
			"510300": {PEPercentile: 61, PBPercentile: 52, ClosePrice: 4.21, TurnoverRate: 1.8, LiquidityState: model.LiquidityNormal, SentimentState: model.SentimentNeutral},
		},
	}

	point, err := source.FetchMarketData(context.Background(), "510300")
	if err != nil {
		t.Fatalf("FetchMarketData: %v", err)
	}
	if point.ClosePrice != 4.21 || point.PEPercentile != 61 || point.LiquidityState != model.LiquidityNormal {
		t.Fatalf("unexpected market point: %+v", point)
	}
}

func TestEastmoneyFundCollectorParsesNetWorthTrend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pingzhongdata/510300.js" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.RawQuery != "" {
			t.Fatalf("P27 collector must not append query parameters, got %s", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`var fS_name = "沪深300ETF"; var fS_code = "510300"; var Data_netWorthTrend = [{"x":1780588800000,"y":4.321,"equityReturn":0.12,"unitMoney":""}]; var Data_ACWorthTrend = [[1780588800000,5.678]]; var Data_assetAllocation = {"stock":"95.1","bond":"0.2","cash":"4.7","netAssets":"100000000"};`))
	}))
	defer server.Close()

	point, err := EastmoneyFundCollector{BaseURL: server.URL, HTTPClient: server.Client()}.FetchMarketData(context.Background(), "510300")
	if err != nil {
		t.Fatalf("FetchMarketData: %v", err)
	}
	if point.ClosePrice != 4.321 || point.SourceName != "eastmoney_fund" || point.SourceLevel != model.SourceLevelB || point.SourceType != "fund_nav" {
		t.Fatalf("unexpected point: %+v", point)
	}
	if point.TradeDate != "2026-06-05" || point.Metadata["fund_name"] != "沪深300ETF" || point.Metadata["fund_code"] != "510300" || point.Metadata["symbol"] != "510300" || point.Metadata["accumulated_nav"] != 5.678 {
		t.Fatalf("metadata not preserved: %+v", point)
	}
	if point.CapturedAt == "" || point.ContentHash == "" || point.Metadata["raw"] == nil {
		t.Fatalf("standard payload fields missing: %+v", point)
	}
	if point.PEPercentile != 0 || point.PBPercentile != 0 || point.VolumePercentile != 0 || point.VolatilityPercentile != 0 {
		t.Fatalf("B-level fund source must not fabricate percentiles: %+v", point)
	}
}

func TestEastmoneyFundCollectorParsesExtendedMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pingzhongdata/510300.js" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`var fS_name = "沪深300ETF"; var fS_code = "510300"; var Data_netWorthTrend = [{"x":1780502400000,"y":4.111,"equityReturn":0.1,"unitMoney":""},{"x":1780588800000,"y":4.321,"equityReturn":0.12,"unitMoney":""}]; var Data_ACWorthTrend = [[1780502400000,5.432],[1780588800000,5.678]]; var Data_assetAllocation = {"stock":"95.1","bond":"0.2","cash":"4.7","netAssets":"100000000"}; var Data_performanceEvaluation = {"avr":"优秀"}; var Data_currentFundManager = [{"name":"张三","workTime":"5年"}];`))
	}))
	defer server.Close()

	point, err := EastmoneyFundCollector{BaseURL: server.URL, HTTPClient: server.Client(), IncludeExtended: true}.FetchMarketData(context.Background(), "510300")
	if err != nil {
		t.Fatalf("FetchMarketData: %v", err)
	}
	history, ok := point.Metadata["nav_history"].([]map[string]any)
	if !ok || len(history) != 2 || history[1]["nav"] != 4.321 || history[1]["accumulated_nav"] != 5.678 {
		t.Fatalf("history not preserved: %+v", point.Metadata)
	}
	assetAllocation, ok := point.Metadata["asset_allocation"].(map[string]any)
	if !ok || assetAllocation["stock"] != "95.1" || assetAllocation["netAssets"] != "100000000" {
		t.Fatalf("asset allocation not preserved: %+v", point.Metadata)
	}
	profile, ok := point.Metadata["fund_profile"].(map[string]any)
	if !ok || profile["performance_evaluation"] == nil || profile["fund_managers"] == nil {
		t.Fatalf("fund profile not preserved: %+v", point.Metadata)
	}
}

func TestCsindexCollectorParsesCurrentPublicIndexBasics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/csindex-home/indexInfo/index-basic-info/000300" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"200","msg":"success","success":true,"data":{"indexCode":"000300","indexShortNameCn":"沪深300","indexShortNameEn":"CSI 300","indexFullNameCn":"沪深300指数","indexFullNameEn":"CSI 300 Index","basicDate":"2004-12-31","basicIndex":1000,"publishDate":"2005-04-08","publishChannelCn":"中证指数官网","currencyCn":"人民币","currencyEn":"CNY","consNumber":300,"adjFreqCn":"半年","indexType":"规模指数"}}`))
	}))
	defer server.Close()

	point, err := CsindexCollector{BaseURL: server.URL, HTTPClient: server.Client()}.FetchMarketData(context.Background(), "000300")
	if err != nil {
		t.Fatalf("FetchMarketData: %v", err)
	}
	if point.ClosePrice != 0 || point.SourceName != "csindex" || point.SourceLevel != model.SourceLevelA || point.SourceType != "index_basic" {
		t.Fatalf("unexpected point: %+v", point)
	}
	if point.Metadata["index_name"] != "沪深300" || point.Metadata["currency"] != "CNY" || point.Metadata["constituent_count"] != float64(300) {
		t.Fatalf("metadata not preserved: %+v", point.Metadata)
	}
	if point.TradeDate != "2005-04-08" || point.CapturedAt == "" || point.ContentHash == "" || point.Metadata["raw"] == nil {
		t.Fatalf("standard payload fields missing: %+v", point)
	}
}

func TestCsindexCollectorParsesIndexBasics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/csindex-home/indexInfo/index-basic-info/000300" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.RawQuery != "" {
			t.Fatalf("P27 collector must not append query parameters, got %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"indexCode":"000300","indexName":"沪深300","indexNameEn":"CSI 300","latestClose":3920.12,"currency":"CNY","publishDate":"2005-04-08"}}`))
	}))
	defer server.Close()

	point, err := CsindexCollector{BaseURL: server.URL, HTTPClient: server.Client()}.FetchMarketData(context.Background(), "000300")
	if err != nil {
		t.Fatalf("FetchMarketData: %v", err)
	}
	if point.ClosePrice != 3920.12 || point.SourceName != "csindex" || point.SourceLevel != model.SourceLevelA || point.SourceType != "index_basic" {
		t.Fatalf("unexpected point: %+v", point)
	}
	if point.Metadata["index_name"] != "沪深300" || point.Metadata["currency"] != "CNY" || point.Metadata["symbol"] != "000300" {
		t.Fatalf("metadata not preserved: %+v", point.Metadata)
	}
	if point.CapturedAt == "" || point.ContentHash == "" || point.Metadata["raw"] == nil {
		t.Fatalf("standard payload fields missing: %+v", point)
	}
}

func TestCsindexCollectorParsesExtendedMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/csindex-home/indexInfo/index-basic-info/000300":
			_, _ = w.Write([]byte(`{"data":{"indexCode":"000300","indexName":"沪深300","indexNameEn":"CSI 300","latestClose":3920.12,"currency":"CNY","publishDate":"2026-06-05"}}`))
		case "/csindex-home/indexInfo/index-consituent/000300":
			_, _ = w.Write([]byte(`{"data":[{"securityCode":"600000","securityName":"浦发银行","weight":0.8},{"securityCode":"000001","securityName":"平安银行","weight":0.7}]}`))
		case "/csindex-home/indexInfo/index-weight/000300":
			_, _ = w.Write([]byte(`{"data":[{"securityCode":"600000","weight":0.81,"tradeDate":"2026-06-05"}]}`))
		case "/csindex-home/indexInfo/index-valuation/000300":
			_, _ = w.Write([]byte(`{"data":[{"fileName":"估值表","fileUrl":"/files/000300.csv","publishDate":"2026-06-05"}]}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	point, err := CsindexCollector{BaseURL: server.URL, HTTPClient: server.Client(), IncludeExtended: true}.FetchMarketData(context.Background(), "000300")
	if err != nil {
		t.Fatalf("FetchMarketData: %v", err)
	}
	constituents, ok := point.Metadata["constituents"].([]map[string]any)
	if !ok || len(constituents) != 2 || constituents[0]["security_code"] != "600000" {
		t.Fatalf("constituents not preserved: %+v", point.Metadata)
	}
	weights, ok := point.Metadata["weights"].([]map[string]any)
	if !ok || len(weights) != 1 || weights[0]["weight"] != 0.81 {
		t.Fatalf("weights not preserved: %+v", point.Metadata)
	}
	valuationFiles, ok := point.Metadata["valuation_files"].([]map[string]any)
	if !ok || len(valuationFiles) != 1 || valuationFiles[0]["file_url"] != server.URL+"/files/000300.csv" {
		t.Fatalf("valuation files not preserved: %+v", point.Metadata)
	}
}

func TestCsindexCollectorRecordsP34ExtendedSourceHealth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/csindex-home/indexInfo/index-basic-info/000300":
			_, _ = w.Write([]byte(`{"data":{"indexCode":"000300","indexName":"沪深300","publishDate":"2026-06-05"}}`))
		case "/csindex-home/indexInfo/index-consituent/000300":
			_, _ = w.Write([]byte(`{"data":[{"securityCode":"600000","securityName":"浦发银行","weight":0.8}]}`))
		case "/csindex-home/indexInfo/index-weight/000300":
			_, _ = w.Write([]byte(`{"data":[]}`))
		case "/csindex-home/indexInfo/index-valuation/000300":
			_, _ = w.Write([]byte(`{"data":`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	point, err := CsindexCollector{BaseURL: server.URL, HTTPClient: server.Client(), IncludeExtended: true}.FetchMarketData(context.Background(), "000300")
	if err != nil {
		t.Fatalf("FetchMarketData: %v", err)
	}
	health, ok := point.Metadata["p34_source_health"].(map[string]any)
	if !ok {
		t.Fatalf("expected P34 source health metadata, got %+v", point.Metadata)
	}
	constituents, _ := health["index_constituents"].(map[string]any)
	weights, _ := health["index_weights"].(map[string]any)
	valuation, _ := health["index_valuation_files"].(map[string]any)
	if constituents["freshness"] != "fresh" || constituents["data_date"] != "2026-06-05" || weights["freshness"] != "no_data" || weights["data_date"] != "2026-06-05" || valuation["freshness"] != "parse_error" || valuation["failure_category"] != "parse_error" || valuation["data_date"] != "2026-06-05" {
		t.Fatalf("unexpected source health: %+v", health)
	}
	categories, ok := point.Metadata["p34_data_categories"].([]string)
	if !ok || len(categories) != 3 {
		t.Fatalf("expected P34 data categories, got %+v", point.Metadata["p34_data_categories"])
	}
}

func TestCsindexCollectorFallsBackToOfficialFilesWhenExtendedEndpointsAreUnavailable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/csindex-home/indexInfo/index-basic-info/000300":
			_, _ = w.Write([]byte(`{"data":{"indexCode":"000300","indexName":"沪深300","publishDate":"2026-06-05"}}`))
		case "/csindex-home/indexInfo/index-consituent/000300", "/csindex-home/indexInfo/index-weight/000300", "/csindex-home/indexInfo/index-valuation/000300":
			http.NotFound(w, r)
		case "/static/html/csindex/public/uploads/file/autofile/closeweight/000300closeweight.xls":
			w.Header().Set("Content-Type", "application/vnd.ms-excel")
			_, _ = w.Write([]byte("official close weight xls content"))
		case "/static/html/csindex/public/uploads/indices/detail/files/zh_CN/000300factsheet.pdf":
			w.Header().Set("Content-Type", "application/pdf")
			_, _ = w.Write([]byte("%PDF-1.7 official factsheet"))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	point, err := CsindexCollector{BaseURL: server.URL, HTTPClient: server.Client(), IncludeExtended: true}.FetchMarketData(context.Background(), "000300")
	if err != nil {
		t.Fatalf("FetchMarketData: %v", err)
	}
	health, ok := point.Metadata["p34_source_health"].(map[string]any)
	if !ok {
		t.Fatalf("expected P34 source health metadata, got %+v", point.Metadata)
	}
	for _, category := range []string{"index_constituents", "index_weights", "index_valuation_files"} {
		item, _ := health[category].(map[string]any)
		if item["freshness"] != "fresh" || item["failure_category"] != nil || item["last_success_at"] == "" {
			t.Fatalf("expected %s to be fresh from official file fallback, got %+v", category, item)
		}
	}
	if files, ok := point.Metadata["official_files"].([]map[string]any); !ok || len(files) != 2 {
		t.Fatalf("expected official file evidence metadata, got %+v", point.Metadata["official_files"])
	}
}

func TestP27CollectorsReturnStableFailures(t *testing.T) {
	tests := []struct {
		name      string
		collector MarketDataSource
		body      string
	}{
		{name: "eastmoney missing trend", collector: EastmoneyFundCollector{}, body: `var fS_name = "沪深300ETF";`},
		{name: "eastmoney invalid trend json", collector: EastmoneyFundCollector{}, body: `var Data_netWorthTrend = bad;`},
		{name: "eastmoney empty trend", collector: EastmoneyFundCollector{}, body: `var Data_netWorthTrend = [];`},
		{name: "csindex invalid json", collector: CsindexCollector{}, body: `{bad`},
		{name: "csindex missing code", collector: CsindexCollector{}, body: `{"data":{"publishDate":"2005-04-08"}}`},
		{name: "csindex missing data", collector: CsindexCollector{}, body: `{"code":"200","success":true}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()
			var err error
			switch tt.collector.(type) {
			case EastmoneyFundCollector:
				_, err = EastmoneyFundCollector{BaseURL: server.URL, HTTPClient: server.Client()}.FetchMarketData(context.Background(), "510300")
			case CsindexCollector:
				_, err = CsindexCollector{BaseURL: server.URL, HTTPClient: server.Client()}.FetchMarketData(context.Background(), "000300")
			}
			if !apperr.IsCode(err, apperr.CodeDataSourceUnavailable) {
				t.Fatalf("expected DATA_SOURCE_UNAVAILABLE, got %v", err)
			}
		})
	}
}

func TestCompositeMarketDataCollectorFallsBackAcrossSources(t *testing.T) {
	collector := CompositeMarketDataCollector{Collectors: []MarketDataSource{failingMarketSource{}, fixedMarketSource{point: MarketDataPoint{ClosePrice: 4.2, SourceName: "eastmoney_fund", SourceLevel: model.SourceLevelB, SourceType: "fund_nav"}}}}

	point, err := collector.FetchMarketData(context.Background(), "510300")
	if err != nil {
		t.Fatalf("FetchMarketData: %v", err)
	}
	if point.ClosePrice != 4.2 || point.SourceName != "eastmoney_fund" {
		t.Fatalf("expected fallback source point, got %+v", point)
	}
}

type failingMarketSource struct{}

func (failingMarketSource) FetchMarketData(context.Context, string) (MarketDataPoint, error) {
	return MarketDataPoint{}, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "failed")
}

type fixedMarketSource struct{ point MarketDataPoint }

func (s fixedMarketSource) FetchMarketData(context.Context, string) (MarketDataPoint, error) {
	return s.point, nil
}

func TestConfiguredMarketDataSourceReturnsStableFailures(t *testing.T) {
	tests := []struct {
		name   string
		source ConfiguredMarketDataSource
	}{
		{name: "unavailable", source: ConfiguredMarketDataSource{Enabled: []string{"fixture"}}},
		{name: "parse failure", source: ConfiguredMarketDataSource{Enabled: []string{"fixture"}, RawMarketFixtures: map[string]string{"510300": `{"close_price":"bad"}`}}},
		{name: "stale", source: ConfiguredMarketDataSource{Enabled: []string{"fixture"}, MarketFixtures: map[string]MarketDataPoint{"510300": {ClosePrice: 1, Stale: true}}}},
		{name: "timeout", source: ConfiguredMarketDataSource{Enabled: []string{"fixture"}, FetchDelay: 20 * time.Millisecond, Timeout: time.Millisecond, MarketFixtures: map[string]MarketDataPoint{"510300": {ClosePrice: 1}}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.source.FetchMarketData(context.Background(), "510300")
			if err == nil {
				t.Fatal("expected stable error")
			}
			var appErr *apperr.AppError
			if !errors.As(err, &appErr) {
				t.Fatalf("expected AppError, got %T %v", err, err)
			}
			if appErr.Code != apperr.CodeDataSourceUnavailable && appErr.Code != apperr.CodeDataStale {
				t.Fatalf("unexpected error code: %s", appErr.Code)
			}
		})
	}
}

func TestConfiguredIntelligenceSourceFetchesReadonlyFixture(t *testing.T) {
	source := ConfiguredIntelligenceSource{
		Enabled: []string{"fixture"},
		IntelligenceFixtures: map[string][]IntelligenceSourceItem{
			"510300": {{SourceName: "exchange", SourceLevel: model.SourceLevelA, Title: "公告", Text: "指数公告", URL: "https://example.invalid/news/1", PublishedAt: "2026-06-01T00:00:00Z"}},
		},
	}

	items, err := source.FetchIntelligence(context.Background(), "510300")
	if err != nil {
		t.Fatalf("FetchIntelligence: %v", err)
	}
	if len(items) != 1 || items[0].SourceName != "exchange" || items[0].URL == "" || items[0].SourceLevel != model.SourceLevelA {
		t.Fatalf("metadata not preserved: %+v", items)
	}
}

func TestConfiguredIntelligenceSourceReturnsStableFailures(t *testing.T) {
	tests := []struct {
		name   string
		source ConfiguredIntelligenceSource
	}{
		{name: "unavailable", source: ConfiguredIntelligenceSource{Enabled: []string{"fixture"}}},
		{name: "parse failure", source: ConfiguredIntelligenceSource{Enabled: []string{"fixture"}, RawIntelligenceFixtures: map[string]string{"510300": `{"bad":`}}},
		{name: "timeout", source: ConfiguredIntelligenceSource{Enabled: []string{"fixture"}, FetchDelay: 20 * time.Millisecond, Timeout: time.Millisecond, IntelligenceFixtures: map[string][]IntelligenceSourceItem{"510300": {{SourceName: "exchange", SourceLevel: model.SourceLevelA, Title: "公告", Text: "内容"}}}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.source.FetchIntelligence(context.Background(), "510300")
			if err == nil {
				t.Fatal("expected stable error")
			}
			var appErr *apperr.AppError
			if !errors.As(err, &appErr) || appErr.Code != apperr.CodeDataSourceUnavailable {
				t.Fatalf("unexpected error: %T %v", err, err)
			}
		})
	}
}

func TestConfiguredMarketDataSourceFetchesReadonlyHTTPProvider(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("symbol") != "510300" {
			t.Fatalf("symbol query not forwarded: %s", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"pe_percentile":62,"pb_percentile":51,"close_price":4.32,"turnover_rate":1.7}`))
	}))
	defer server.Close()

	source := ConfiguredMarketDataSource{Enabled: []string{"http"}, MarketEndpoint: server.URL, HTTPClient: server.Client()}
	point, err := source.FetchMarketData(context.Background(), "510300")
	if err != nil {
		t.Fatalf("FetchMarketData: %v", err)
	}
	if point.ClosePrice != 4.32 || point.PEPercentile != 62 || point.PBPercentile != 51 {
		t.Fatalf("unexpected point: %+v", point)
	}
}

func TestConfiguredMarketDataSourceHonorsEnabledSources(t *testing.T) {
	serverCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCalled = true
		_, _ = w.Write([]byte(`{"close_price":4.32}`))
	}))
	defer server.Close()

	source := ConfiguredMarketDataSource{Enabled: []string{"fixture"}, MarketEndpoint: server.URL, HTTPClient: server.Client()}
	_, err := source.FetchMarketData(context.Background(), "510300")

	if err == nil {
		t.Fatal("expected unavailable data source when http is not enabled")
	}
	if serverCalled {
		t.Fatal("http endpoint should not be called when not enabled")
	}
}
func TestConfiguredMarketDataSourceParsesPublicHTTPPayload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("symbol") != "510300" {
			t.Fatalf("symbol query not forwarded: %s", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"data":{"symbol":"510300","trade_date":"2026-06-01","close":4.56,"turnover_rate":1.9,"valuation":{"pe_percentile":63,"pb_percentile":52},"liquidity":"normal","sentiment":"neutral"}}`))
	}))
	defer server.Close()

	source := ConfiguredMarketDataSource{Enabled: []string{"public_http"}, MarketEndpoint: server.URL, HTTPClient: server.Client()}
	point, err := source.FetchMarketData(context.Background(), "510300")
	if err != nil {
		t.Fatalf("FetchMarketData: %v", err)
	}
	if point.ClosePrice != 4.56 || point.TurnoverRate != 1.9 || point.PEPercentile != 63 || point.PBPercentile != 52 || point.LiquidityState != model.LiquidityNormal || point.SentimentState != model.SentimentNeutral {
		t.Fatalf("unexpected public market point: %+v", point)
	}
}

func TestConfiguredMarketDataSourceUsesEndpointForConfiguredPublicSourceNames(t *testing.T) {
	serverCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCalled = true
		_, _ = w.Write([]byte(`{"close_price":4.88}`))
	}))
	defer server.Close()

	source := ConfiguredMarketDataSource{Enabled: []string{"official", "exchange"}, MarketEndpoint: server.URL, HTTPClient: server.Client(), Fallback: StubMarketDataSource{}}
	point, err := source.FetchMarketData(context.Background(), "510300")

	if err != nil {
		t.Fatalf("FetchMarketData: %v", err)
	}
	if !serverCalled || point.ClosePrice != 4.88 {
		t.Fatalf("expected configured public endpoint to be used, called=%v point=%+v", serverCalled, point)
	}
}

func TestConfiguredMarketDataSourceFallsBackWhenPublicHTTPFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unavailable", http.StatusInternalServerError)
	}))
	defer server.Close()

	source := ConfiguredMarketDataSource{Enabled: []string{"public_http"}, MarketEndpoint: server.URL, HTTPClient: server.Client(), Fallback: StubMarketDataSource{}}
	point, err := source.FetchMarketData(context.Background(), "510300")

	if err != nil {
		t.Fatalf("FetchMarketData: %v", err)
	}
	if point.ClosePrice <= 0 || point.LiquidityState == "" || point.SentimentState == "" {
		t.Fatalf("expected stub fallback market data, got %+v", point)
	}
}

func TestConfiguredMarketDataSourceParsesPublicHTTPDataArrayPayload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":[{"symbol":"510300","close":4.56,"valuation":{"pe_percentile":63,"pb_percentile":52}}]}`))
	}))
	defer server.Close()

	source := ConfiguredMarketDataSource{Enabled: []string{"public_http"}, MarketEndpoint: server.URL, HTTPClient: server.Client()}
	point, err := source.FetchMarketData(context.Background(), "510300")

	if err != nil {
		t.Fatalf("FetchMarketData: %v", err)
	}
	if point.ClosePrice != 4.56 || point.PEPercentile != 63 || point.PBPercentile != 52 {
		t.Fatalf("unexpected public market array point: %+v", point)
	}
}

func TestConfiguredIntelligenceSourceFetchesReadonlyHTTPProvider(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("symbol") != "510300" {
			t.Fatalf("symbol query not forwarded: %s", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`[{"source_name":"exchange","source_level":"A","title":"公告","text":"指数公告","url":"https://example.invalid/news/1","published_at":"2026-06-01T00:00:00Z"}]`))
	}))
	defer server.Close()

	source := ConfiguredIntelligenceSource{Enabled: []string{"http"}, IntelligenceEndpoint: server.URL, HTTPClient: server.Client()}
	items, err := source.FetchIntelligence(context.Background(), "510300")
	if err != nil {
		t.Fatalf("FetchIntelligence: %v", err)
	}
	if len(items) != 1 || items[0].SourceName != "exchange" || items[0].URL == "" {
		t.Fatalf("metadata not preserved: %+v", items)
	}
}

func TestConfiguredIntelligenceSourceParsesPublicHTTPTopLevelArrayPayload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"source":"exchange_disclosure","title":"ETF 公告","content":"公告正文","published_at":"2026-06-01T00:00:00Z"},{"source":"fund_company_disclosure","title":"基金公告","summary":"披露摘要","original_url":"https://example.invalid/fund","published_at":"2026-06-01T01:00:00Z"}]`))
	}))
	defer server.Close()

	source := ConfiguredIntelligenceSource{Enabled: []string{"public_http"}, IntelligenceEndpoint: server.URL, HTTPClient: server.Client()}
	items, err := source.FetchIntelligence(context.Background(), "510300")

	if err != nil {
		t.Fatalf("FetchIntelligence: %v", err)
	}
	if len(items) != 2 || items[0].SourceName != "exchange_disclosure" || items[0].Text != "公告正文" || items[0].SourceLevel != model.SourceLevelA || items[1].Text != "披露摘要" || items[1].URL != "https://example.invalid/fund" {
		t.Fatalf("expected public top-level array evidence payload, got %+v", items)
	}
}

func TestConfiguredIntelligenceSourceParsesPublicHTTPDataArrayPayload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":[{"source":"exchange_disclosure","title":"ETF 公告","content":"公告正文","published_at":"2026-06-01T00:00:00Z"},{"source":"fund_company_disclosure","title":"基金公告","content":"披露正文","published_at":"2026-06-01T01:00:00Z"}]}`))
	}))
	defer server.Close()

	source := ConfiguredIntelligenceSource{Enabled: []string{"public_http"}, IntelligenceEndpoint: server.URL, HTTPClient: server.Client()}
	items, err := source.FetchIntelligence(context.Background(), "510300")

	if err != nil {
		t.Fatalf("FetchIntelligence: %v", err)
	}
	if len(items) != 2 || items[0].SourceLevel != model.SourceLevelA || items[1].SourceLevel != model.SourceLevelA {
		t.Fatalf("expected data array public evidence payload, got %+v", items)
	}
}

func TestConfiguredIntelligenceSourceParsesPublicHTTPPayloadAndDedupes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("symbol") != "510300" {
			t.Fatalf("symbol query not forwarded: %s", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"items":[{"source":"exchange_disclosure","title":"ETF 公告","content":"510300 公告正文","url":"https://example.invalid/a","published_at":"2026-06-01T00:00:00Z"},{"source":"exchange_disclosure","title":"ETF 公告","content":"重复正文","url":"https://example.invalid/a","published_at":"2026-06-01T00:00:00Z"},{"source":"financial_news","title":"背景新闻","content":"基金背景材料","url":"https://example.invalid/b","published_at":"2026-06-01T01:00:00Z"}]}`))
	}))
	defer server.Close()

	source := ConfiguredIntelligenceSource{Enabled: []string{"public_http"}, IntelligenceEndpoint: server.URL, HTTPClient: server.Client()}
	items, err := source.FetchIntelligence(context.Background(), "510300")
	if err != nil {
		t.Fatalf("FetchIntelligence: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected duplicate URL to be removed, got %+v", items)
	}
	if items[0].SourceName != "exchange_disclosure" || items[0].SourceLevel != model.SourceLevelA || items[0].Text != "510300 公告正文" {
		t.Fatalf("expected exchange disclosure mapped to A-level evidence, got %+v", items[0])
	}
	if items[1].SourceName != "financial_news" || items[1].SourceLevel != model.SourceLevelB {
		t.Fatalf("expected financial news mapped to B-level background evidence, got %+v", items[1])
	}
}

func TestConfiguredIntelligenceSourceKeepsSameTitleFromDifferentSources(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"items":[{"source":"exchange_disclosure","title":"同题公告","content":"交易所公告","published_at":"2026-06-01T00:00:00Z"},{"source":"fund_company_disclosure","title":"同题公告","content":"基金公司公告","published_at":"2026-06-01T00:00:00Z"}]}`))
	}))
	defer server.Close()

	source := ConfiguredIntelligenceSource{Enabled: []string{"public_http"}, IntelligenceEndpoint: server.URL, HTTPClient: server.Client()}
	items, err := source.FetchIntelligence(context.Background(), "510300")

	if err != nil {
		t.Fatalf("FetchIntelligence: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected same title from different sources to remain independent, got %+v", items)
	}
}
func TestConfiguredIntelligenceSourceInfersChinesePublicSourceLevels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"items":[{"source":"巨潮资讯","title":"ETF 公告","content":"公告正文","published_at":"2026-06-01T00:00:00Z"},{"source":"基金公司公告","title":"基金公告","content":"公告正文","published_at":"2026-06-01T00:00:00Z"}]}`))
	}))
	defer server.Close()

	source := ConfiguredIntelligenceSource{Enabled: []string{"public_http"}, IntelligenceEndpoint: server.URL, HTTPClient: server.Client()}
	items, err := source.FetchIntelligence(context.Background(), "510300")

	if err != nil {
		t.Fatalf("FetchIntelligence: %v", err)
	}
	if len(items) != 2 || items[0].SourceLevel != model.SourceLevelA || items[1].SourceLevel != model.SourceLevelA {
		t.Fatalf("expected Chinese public disclosure sources to map to A level, got %+v", items)
	}
}
func TestStubSourcesRemainOfflineDeterministic(t *testing.T) {
	marketA, err := (StubMarketDataSource{}).FetchMarketData(context.Background(), "510300")
	if err != nil {
		t.Fatalf("stub market: %v", err)
	}
	marketB, err := (StubMarketDataSource{}).FetchMarketData(context.Background(), "510300")
	if err != nil {
		t.Fatalf("stub market second: %v", err)
	}
	if marketA.ClosePrice != marketB.ClosePrice || marketA.PEPercentile != marketB.PEPercentile || marketA.PBPercentile != marketB.PBPercentile || marketA.LiquidityState != marketB.LiquidityState || marketA.SentimentState != marketB.SentimentState {
		t.Fatalf("stub market must be deterministic: %+v %+v", marketA, marketB)
	}

	itemsA, err := (StubIntelligenceSource{}).FetchIntelligence(context.Background(), "510300")
	if err != nil {
		t.Fatalf("stub intelligence: %v", err)
	}
	itemsB, err := (StubIntelligenceSource{}).FetchIntelligence(context.Background(), "510300")
	if err != nil {
		t.Fatalf("stub intelligence second: %v", err)
	}
	if len(itemsA) != 1 || len(itemsB) != 1 || itemsA[0] != itemsB[0] {
		t.Fatalf("stub intelligence must be deterministic: %+v %+v", itemsA, itemsB)
	}
}
