package workflow

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"investment-agent/internal/domain/model"
)

func TestEastmoneyFundEvidenceCollectorFetchesFundProfile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pingzhongdata/510300.js" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/javascript")
		w.Write([]byte(`var fS_name = "沪深300ETF华泰柏瑞";var fS_code = "510300";var Data_netWorthTrend = [{"x":1781712000000,"y":4.7533,"equityReturn":0.12,"unitMoney":""}];var Data_assetAllocation = {"series":[{"name":"股票占净比","data":[98.89]}]};`))
	}))
	defer server.Close()

	collector := EastmoneyFundEvidenceCollector{BaseURL: server.URL, HTTPClient: server.Client()}
	items, err := collector.FetchPublicEvidence(context.Background(), "510300", time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("FetchPublicEvidence: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %+v", items)
	}
	item := items[0]
	if item.SourceName != "eastmoney_fund" || item.SourceLevel != model.SourceLevelB || item.SourceType != "fund_profile" || item.EvidenceRole != "formal" {
		t.Fatalf("unexpected source metadata: %+v", item)
	}
	if item.Symbol != "510300" || item.Title != fundProfileEvidenceTitle("510300") || !strings.Contains(item.Text, "沪深300ETF华泰柏瑞") || !strings.Contains(item.Text, "4.7533") {
		t.Fatalf("unexpected item content: %+v", item)
	}
}

func TestEastmoneyFundEvidenceCollectorUsesRequestedNon510300Symbol(t *testing.T) {
	var seenPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		if r.URL.Path != "/pingzhongdata/159915.js" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/javascript")
		w.Write([]byte(`var fS_name = "创业板ETF";var fS_code = "159915";var Data_netWorthTrend = [{"x":1781712000000,"y":2.4130,"equityReturn":0.08,"unitMoney":""}];`))
	}))
	defer server.Close()

	collector := EastmoneyFundEvidenceCollector{BaseURL: server.URL, HTTPClient: server.Client()}
	items, err := collector.FetchPublicEvidence(context.Background(), "159915", time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("FetchPublicEvidence: %v", err)
	}
	if seenPath != "/pingzhongdata/159915.js" {
		t.Fatalf("expected request path for user symbol, got %s", seenPath)
	}
	if len(items) != 1 || items[0].Symbol != "159915" || !strings.Contains(items[0].Text, "创业板ETF") {
		t.Fatalf("expected non-510300 fund evidence bound to 159915, got %+v", items)
	}
}

func TestCsindexIndexEvidenceCollectorFetchesTrackedIndexProfile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/csindex-home/indexInfo/index-basic-info/000300" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code":"200","msg":"Success","data":{"indexFullNameCn":"沪深300指数","indexShortNameCn":"沪深300","indexCode":"000300","publishDate":"2005-04-08","indexCnDesc":"沪深300指数由沪深市场中规模大、流动性好的300只证券组成。"}}`))
	}))
	defer server.Close()

	collector := CsindexIndexEvidenceCollector{BaseURL: server.URL, HTTPClient: server.Client(), SymbolToIndex: map[string]string{"510300": "000300"}}
	items, err := collector.FetchPublicEvidence(context.Background(), "510300", time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("FetchPublicEvidence: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %+v", items)
	}
	item := items[0]
	if item.SourceName != "csindex_index" || item.SourceLevel != model.SourceLevelA || item.SourceType != "fund_profile" || item.EvidenceRole != "formal" {
		t.Fatalf("unexpected source metadata: %+v", item)
	}
	if item.Symbol != "510300" || item.Title != fundProfileEvidenceTitle("510300") || !strings.Contains(item.Text, "沪深300指数") || !strings.Contains(item.Text, "000300") {
		t.Fatalf("unexpected item content: %+v", item)
	}
}

func TestCsindexIndexEvidenceCollectorUsesKnownNon510300TrackedIndex(t *testing.T) {
	var seenPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		if r.URL.Path != "/csindex-home/indexInfo/index-basic-info/399006" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code":"200","msg":"Success","data":{"indexFullNameCn":"创业板指数","indexShortNameCn":"创业板指","indexCode":"399006","publishDate":"2010-06-01","indexCnDesc":"创业板指数反映创业板市场代表性证券表现。"}}`))
	}))
	defer server.Close()

	collector := CsindexIndexEvidenceCollector{BaseURL: server.URL, HTTPClient: server.Client()}
	items, err := collector.FetchPublicEvidence(context.Background(), "159915", time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("FetchPublicEvidence: %v", err)
	}
	if seenPath != "/csindex-home/indexInfo/index-basic-info/399006" {
		t.Fatalf("expected tracked index request for 399006, got %s", seenPath)
	}
	if len(items) != 1 || items[0].Symbol != "159915" || !strings.Contains(items[0].Text, "399006") {
		t.Fatalf("expected index evidence bound to 159915/399006, got %+v", items)
	}
}
