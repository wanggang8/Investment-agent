package workflow

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/pkg/apperr"
)

func TestCsrcCollectorFetchesRegulatoryInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/searchList" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Query().Get("keyword") != "基金" {
			t.Fatalf("unexpected keyword: %s", r.URL.Query().Get("keyword"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"data": {
				"page": 1,
				"rows": 50,
				"channelId": "all",
				"total": 1,
				"results": [{
					"title": "关于加强公开募集证券投资基金流动性风险管理的规定",
					"content": "为规范公开募集证券投资基金运作，保护投资者合法权益，根据《证券投资基金法》等法律法规，制定本规定。",
					"contentHtml": "<p>为规范公开募集证券投资基金运作...</p>",
					"memo": "基金流动性风险管理规定",
					"url": "/csrc/c100028/c1002866/content.shtml",
					"publishedTime": "2017-08-31T00:00:00+08:00",
					"publishedTimeStr": "2017年08月31日",
					"channelName": "部门规章",
					"channelCodeName": "rules",
					"manuscriptId": "1502866",
					"resList": [],
					"domainMetaList": []
				}]
			}
		}`))
	}))
	defer server.Close()

	collector := CsrcCollector{BaseURL: server.URL, HTTPClient: server.Client()}
	items, err := collector.FetchPublicEvidence(context.Background(), "基金", time.Date(2017, 8, 1, 0, 0, 0, 0, time.UTC), time.Date(2017, 9, 30, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("FetchPublicEvidence: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	item := items[0]
	if item.SourceName != "csrc" || item.SourceLevel != model.SourceLevelA || item.SourceType != "regulatory_disclosure" {
		t.Fatalf("unexpected source metadata: %+v", item)
	}
	if item.SourceRecordID != "1502866" || item.Title == "" || item.ContentHash == "" {
		t.Fatalf("expected normalized item with hash, got %+v", item)
	}
	if item.EvidenceRole != "formal" {
		t.Fatalf("expected formal evidence role for rules channel, got %s", item.EvidenceRole)
	}
	if item.URL != server.URL+"/csrc/c100028/c1002866/content.shtml" {
		t.Fatalf("unexpected URL: %s", item.URL)
	}
}

func TestCsrcCollectorHandlesBackgroundChannel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"data": {
				"page": 1,
				"rows": 50,
				"total": 1,
				"results": [{
					"title": "证监会新闻发布会",
					"content": "本周证监会召开例行新闻发布会...",
					"url": "/csrc/c100028/c1002900/content.shtml",
					"publishedTime": "2024-06-05T15:00:00+08:00",
					"publishedTimeStr": "2024年06月05日",
					"channelName": "新闻发布",
					"channelCodeName": "news",
					"manuscriptId": "1502900",
					"resList": [],
					"domainMetaList": []
				}]
			}
		}`))
	}))
	defer server.Close()

	collector := CsrcCollector{BaseURL: server.URL, HTTPClient: server.Client()}
	items, err := collector.FetchPublicEvidence(context.Background(), "基金", time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("FetchPublicEvidence: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	item := items[0]
	if item.EvidenceRole != "background" {
		t.Fatalf("expected background evidence role for news channel, got %s", item.EvidenceRole)
	}
}

func TestCsrcCollectorUsesPaginationDateFilterAndURLResolution(t *testing.T) {
	seenPages := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenPages = append(seenPages, r.URL.Query().Get("page"))
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("page") == "1" {
			w.Write([]byte(`{"data":{"page":1,"rows":50,"total":51,"results":[{"title":"范围外监管信息","content":"旧正文","url":"/old.shtml","publishedTime":"2024-01-01T00:00:00+08:00","channelCodeName":"rules","manuscriptId":"old"}]}}`))
			return
		}
		w.Write([]byte(`{"data":{"page":2,"rows":50,"total":51,"results":[{"title":"范围内监管信息","content":"新正文","url":"https://www.csrc.gov.cn/new.shtml","publishedTime":"2024-06-05T00:00:00+08:00","channelCodeName":"rules","manuscriptId":"new"}]}}`))
	}))
	defer server.Close()

	collector := CsrcCollector{BaseURL: server.URL, HTTPClient: server.Client()}
	items, err := collector.FetchPublicEvidence(context.Background(), "基金", time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 6, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("FetchPublicEvidence: %v", err)
	}
	if len(items) != 1 || items[0].SourceRecordID != "new" || len(seenPages) != 2 {
		t.Fatalf("expected paginated in-range item only, pages=%+v items=%+v", seenPages, items)
	}
	if items[0].URL != "https://www.csrc.gov.cn/new.shtml" {
		t.Fatalf("expected absolute URL to be preserved, got %s", items[0].URL)
	}
}
func TestCsrcCollectorSkipsMalformedRows(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"page":1,"rows":50,"total":2,"results":[{"title":"坏记录","content":"","url":"/bad.shtml","publishedTime":"not-a-date","channelCodeName":"rules","manuscriptId":"bad"},{"title":"好记录","content":"有效正文","url":"/good.shtml","publishedTime":"2024-06-05T00:00:00+08:00","channelCodeName":"rules","manuscriptId":"good"}]}}`))
	}))
	defer server.Close()

	collector := CsrcCollector{BaseURL: server.URL, HTTPClient: server.Client()}
	items, err := collector.FetchPublicEvidence(context.Background(), "基金", time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 6, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("FetchPublicEvidence: %v", err)
	}
	if len(items) != 1 || items[0].SourceRecordID != "good" {
		t.Fatalf("expected only valid row, got %+v", items)
	}
}

func TestCsrcCollectorHandlesEmptyResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data": {"page": 1, "rows": 50, "total": 0, "results": []}}`))
	}))
	defer server.Close()

	collector := CsrcCollector{BaseURL: server.URL, HTTPClient: server.Client()}
	_, err := collector.FetchPublicEvidence(context.Background(), "missing", time.Time{}, time.Time{})
	if !apperr.IsCode(err, apperr.CodeNotFound) {
		t.Fatalf("expected not found for reachable empty results, got %v", err)
	}
	if sourceErr, ok := publicEvidenceErrorOf(err); !ok || sourceErr.ErrorCode != "no_data" || sourceErr.SourceName != "csrc" {
		t.Fatalf("expected source-specific no_data error, got %+v", err)
	}
}
