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

func TestCninfoCollectorFetchesAnnouncements(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/new/hisAnnouncement/query" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("stock") != "510300,9900000091" {
			t.Fatalf("unexpected stock: %s", r.Form.Get("stock"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"totalAnnouncement": 1,
			"totalRecordNum": 1,
			"announcements": [{
				"secCode": "510300",
				"secName": "300ETF",
				"orgId": "9900000091",
				"announcementId": "1218840123",
				"announcementTitle": "华泰柏瑞沪深300ETF基金产品资料概要更新",
				"announcementTime": 1717516800000,
				"adjunctUrl": "/finalpage/2024-06-05/1218840123.PDF",
				"adjunctSize": 245678,
				"adjunctType": "PDF",
				"announcementType": "基金公告"
			}],
			"hasMore": false,
			"totalpages": 1
		}`))
	}))
	defer server.Close()

	collector := CninfoCollector{BaseURL: server.URL, HTTPClient: server.Client()}
	items, err := collector.FetchPublicEvidence(context.Background(), "510300", time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("FetchPublicEvidence: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	item := items[0]
	if item.SourceName != "cninfo" || item.SourceLevel != model.SourceLevelA || item.SourceType != "public_disclosure" {
		t.Fatalf("unexpected source metadata: %+v", item)
	}
	if item.SourceRecordID != "1218840123" || item.Title == "" || item.ContentHash == "" {
		t.Fatalf("expected normalized item with hash, got %+v", item)
	}
	if item.AttachmentURL != server.URL+"/finalpage/2024-06-05/1218840123.PDF" {
		t.Fatalf("unexpected attachment URL: %s", item.AttachmentURL)
	}
}

func TestCninfoCollectorUsesDefaultDateRangeAndPagination(t *testing.T) {
	seenPages := []string{}
	seenDates := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		seenPages = append(seenPages, r.Form.Get("pageNum"))
		seenDates = append(seenDates, r.Form.Get("seDate"))
		w.Header().Set("Content-Type", "application/json")
		if r.Form.Get("pageNum") == "1" {
			w.Write([]byte(`{"totalAnnouncement":2,"totalRecordNum":2,"announcements":[{"secCode":"510300","secName":"300ETF","orgId":"org","announcementId":"ann-1","announcementTitle":"第一页公告","announcementTime":1717516800000,"adjunctUrl":"/a.pdf","adjunctType":"PDF","announcementType":"基金公告"}],"hasMore":true,"totalpages":2}`))
			return
		}
		w.Write([]byte(`{"totalAnnouncement":2,"totalRecordNum":2,"announcements":[{"secCode":"510300","secName":"300ETF","orgId":"org","announcementId":"ann-2","announcementTitle":"第二页公告","announcementTime":1717603200000,"adjunctUrl":"/b.pdf","adjunctType":"PDF","announcementType":"基金公告"}],"hasMore":false,"totalpages":2}`))
	}))
	defer server.Close()

	collector := CninfoCollector{BaseURL: server.URL, HTTPClient: server.Client()}
	items, err := collector.FetchPublicEvidence(context.Background(), "510300", time.Time{}, time.Date(2024, 6, 6, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("FetchPublicEvidence: %v", err)
	}
	if len(items) != 2 || seenPages[0] != "1" || seenPages[1] != "2" {
		t.Fatalf("expected two paged items, pages=%+v items=%+v", seenPages, items)
	}
	if seenDates[0] != "2024-03-08~2024-06-06" {
		t.Fatalf("expected default 90-day seDate, got %+v", seenDates)
	}
}
func TestCninfoCollectorPreservesExplicitOrgIDStockParam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("stock") != "159915,9900001234" {
			t.Fatalf("unexpected stock: %s", r.Form.Get("stock"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"totalAnnouncement":0,"totalRecordNum":0,"announcements":[],"hasMore":false,"totalpages":1}`))
	}))
	defer server.Close()

	collector := CninfoCollector{BaseURL: server.URL, HTTPClient: server.Client()}
	_, err := collector.FetchPublicEvidence(context.Background(), "159915,9900001234", time.Time{}, time.Time{})
	if !apperr.IsCode(err, apperr.CodeNotFound) {
		t.Fatalf("expected not found after preserving explicit stock param, got %v", err)
	}
}

func TestCninfoCollectorUsesConfiguredOrgIDMapping(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("stock") != "159915,9900001234" {
			t.Fatalf("unexpected stock: %s", r.Form.Get("stock"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"totalAnnouncement":0,"totalRecordNum":0,"announcements":[],"hasMore":false,"totalpages":1}`))
	}))
	defer server.Close()

	collector := CninfoCollector{BaseURL: server.URL, HTTPClient: server.Client(), OrgIDBySymbol: map[string]string{"159915": "9900001234"}}
	_, err := collector.FetchPublicEvidence(context.Background(), "159915", time.Time{}, time.Time{})
	if !apperr.IsCode(err, apperr.CodeNotFound) {
		t.Fatalf("expected not found after using configured orgId mapping, got %v", err)
	}
}

func TestCninfoCollectorStoresCleanSymbolForExplicitOrgIDInput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("stock") != "159915,9900001234" {
			t.Fatalf("unexpected stock: %s", r.Form.Get("stock"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"totalAnnouncement":1,"totalRecordNum":1,"announcements":[{"secCode":"159915","secName":"创业板ETF","orgId":"9900001234","announcementId":"ann-159915","announcementTitle":"创业板ETF 公告","announcementTime":1717516800000,"adjunctUrl":"/a.pdf","adjunctType":"PDF","announcementType":"基金公告"}],"hasMore":false,"totalpages":1}`))
	}))
	defer server.Close()

	collector := CninfoCollector{BaseURL: server.URL, HTTPClient: server.Client()}
	items, err := collector.FetchPublicEvidence(context.Background(), "159915,9900001234", time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("FetchPublicEvidence: %v", err)
	}
	if len(items) != 1 || items[0].Symbol != "159915" {
		t.Fatalf("expected clean symbol 159915, got %+v", items)
	}
}

func TestCninfoCollectorHandles404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}))
	defer server.Close()

	collector := CninfoCollector{BaseURL: server.URL, HTTPClient: server.Client()}
	_, err := collector.FetchPublicEvidence(context.Background(), "missing", time.Time{}, time.Time{})
	if err == nil {
		t.Fatal("expected error for 404")
	}
}

func TestCninfoCollectorHandlesEmptyAnnouncements(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"totalAnnouncement": 0, "totalRecordNum": 0, "announcements": [], "hasMore": false, "totalpages": 0}`))
	}))
	defer server.Close()

	collector := CninfoCollector{BaseURL: server.URL, HTTPClient: server.Client()}
	_, err := collector.FetchPublicEvidence(context.Background(), "empty", time.Time{}, time.Time{})
	if !apperr.IsCode(err, apperr.CodeNotFound) {
		t.Fatalf("expected not found for reachable empty announcements, got %v", err)
	}
	if sourceErr, ok := publicEvidenceErrorOf(err); !ok || sourceErr.ErrorCode != "no_data" || sourceErr.SourceName != "cninfo" {
		t.Fatalf("expected source-specific no_data error, got %+v", err)
	}
}
