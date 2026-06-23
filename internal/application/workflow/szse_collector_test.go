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

func TestSzseCollectorFetchesAnnouncements(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/disc/announcement/searchQuery" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Query().Get("keyword") != "159915" {
			t.Fatalf("unexpected keyword: %s", r.URL.Query().Get("keyword"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"companyCount": 1,
			"announceCount": 1,
			"disclosureTip": "",
			"recordCount": 1,
			"data": [{
				"secCode": "159915",
				"secName": "易方达创业板ETF",
				"announList": [{
					"id": "12345678",
					"title": "易方达创业板交易型开放式指数证券投资基金2024年第1季度报告",
					"attachPath": "/finalpage/2024-04-22/159915_report.PDF",
					"attachFormat": "PDF",
					"attachSize": 345678,
					"annId": "ann-12345678",
					"bigCategoryId": "010301",
					"bigCategoryName": "定期报告",
					"publishTime": "2024-04-22 17:30:00",
					"importantRatio": "0"
				}]
			}]
		}`))
	}))
	defer server.Close()

	collector := SzseCollector{BaseURL: server.URL, HTTPClient: server.Client()}
	items, err := collector.FetchPublicEvidence(context.Background(), "159915", time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("FetchPublicEvidence: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	item := items[0]
	if item.SourceName != "szse" || item.SourceLevel != model.SourceLevelA || item.SourceType != "public_disclosure" {
		t.Fatalf("unexpected source metadata: %+v", item)
	}
	if item.SourceRecordID != "ann-12345678" || item.Title == "" || item.ContentHash == "" {
		t.Fatalf("expected normalized item with hash, got %+v", item)
	}
	if item.AttachmentURL != server.URL+"/finalpage/2024-04-22/159915_report.PDF" {
		t.Fatalf("unexpected attachment URL: %s", item.AttachmentURL)
	}
	if item.PublishedAt != "2024-04-22 17:30:00" {
		t.Fatalf("unexpected published time: %s", item.PublishedAt)
	}
}

func TestSzseCollectorUsesPaginationAndDateFilter(t *testing.T) {
	seenPages := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenPages = append(seenPages, r.URL.Query().Get("pageNum"))
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("pageNum") == "1" {
			w.Write([]byte(`{"recordCount":51,"data":[{"secCode":"159915","secName":"创业板ETF","announList":[{"id":"old","title":"范围外公告","attachPath":"/old.pdf","attachFormat":"PDF","attachSize":1,"annId":"old","bigCategoryName":"定期报告","publishTime":"2024-01-01 10:00:00"}]}]}`))
			return
		}
		w.Write([]byte(`{"recordCount":51,"data":[{"secCode":"159915","secName":"创业板ETF","announList":[{"id":"new","title":"范围内公告","attachPath":"/new.pdf","attachFormat":"PDF","attachSize":1,"annId":"new","bigCategoryName":"定期报告","publishTime":"2024-06-05 10:00:00"}]}]}`))
	}))
	defer server.Close()

	collector := SzseCollector{BaseURL: server.URL, HTTPClient: server.Client()}
	items, err := collector.FetchPublicEvidence(context.Background(), "159915", time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 6, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("FetchPublicEvidence: %v", err)
	}
	if len(items) != 1 || items[0].SourceRecordID != "new" || len(seenPages) != 2 {
		t.Fatalf("expected paginated in-range item only, pages=%+v items=%+v", seenPages, items)
	}
}
func TestSzseCollectorHandlesHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	collector := SzseCollector{BaseURL: server.URL, HTTPClient: server.Client()}
	_, err := collector.FetchPublicEvidence(context.Background(), "error", time.Time{}, time.Time{})
	if err == nil {
		t.Fatal("expected error for 500")
	}
}

func TestSzseCollectorHandlesEmptyData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"companyCount": 0, "announceCount": 0, "recordCount": 0, "data": []}`))
	}))
	defer server.Close()

	collector := SzseCollector{BaseURL: server.URL, HTTPClient: server.Client()}
	_, err := collector.FetchPublicEvidence(context.Background(), "empty", time.Time{}, time.Time{})
	if !apperr.IsCode(err, apperr.CodeNotFound) {
		t.Fatalf("expected not found for reachable empty data, got %v", err)
	}
	if sourceErr, ok := publicEvidenceErrorOf(err); !ok || sourceErr.ErrorCode != "no_data" || sourceErr.SourceName != "szse" {
		t.Fatalf("expected source-specific no_data error, got %+v", err)
	}
}
