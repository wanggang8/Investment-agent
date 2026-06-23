package workflow

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestP89StructuredPublicCollectorFetchesCapitalFlowMarginAndFinancialFields(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/qt/stock/fflow/daykline/get", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"rc":0,"data":{"code":"600000","market":1,"name":"浦发银行","klines":["2026-06-22,11895999.0,1280999.0,-13176999.0,3760767.0,8135232.0,1.75,0.19,-1.94,0.55,1.20,9.16,0.77"]}}`))
	})
	mux.HandleFunc("/api/data/v1/get", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"success":true,"result":{"data":[{"SECURITY_CODE":"600000","TOTAL_OPERATE_INCOME":46573000000,"PARENT_NETPROFIT":17861000000,"SJLTZ":1.49,"NOTICE_DATE":"2026-04-30 00:00:00"}]}}`))
	})
	mux.HandleFunc("/marketdata/tradedata/queryMargin.do", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Referer") == "" {
			t.Fatalf("expected SSE referer header")
		}
		_, _ = w.Write([]byte(`{"pageHelp":{"data":[{"opDate":"20260618","rzye":1495185799372},{"opDate":"20260617","rzye":1488211916696}]}}`))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	point, err := P89StructuredPublicCollector{
		EastmoneyPush2BaseURL:      server.URL,
		EastmoneyH5BaseURL:         server.URL,
		EastmoneyDatacenterBaseURL: server.URL,
		SSEQueryBaseURL:            server.URL,
		HTTPClient:                 server.Client(),
	}.FetchMarketData(context.Background(), "600000")
	if err != nil {
		t.Fatalf("fetch structured public data: %v", err)
	}

	readback := P88NormalizeStructuredDataMetadata(point.Metadata)
	if readback.CapitalFlow == nil || readback.CapitalFlow.Date != "2026-06-22" || readback.CapitalFlow.NetInflow != 11895999 || readback.CapitalFlow.NetOutflow != 13176999 {
		t.Fatalf("unexpected capital flow readback: %+v", readback.CapitalFlow)
	}
	if readback.MarginFinancing == nil || readback.MarginFinancing.Date != "2026-06-18" || readback.MarginFinancing.MarginBalance != 1495185799372 {
		t.Fatalf("unexpected margin readback: %+v", readback.MarginFinancing)
	}
	if readback.MarginFinancing.BalanceChangeRate <= 0 {
		t.Fatalf("expected computed margin balance change rate, got %+v", readback.MarginFinancing)
	}
	if readback.ConstituentFinancial == nil || readback.ConstituentFinancial.Revenue != 46573000000 || readback.ConstituentFinancial.NetProfit != 17861000000 || readback.ConstituentFinancial.DisclosureDate != "2026-04-30" {
		t.Fatalf("unexpected financial readback: %+v", readback.ConstituentFinancial)
	}
	if point.SourceName != "p89_public_structured_sources" || point.SourceType != "structured_public_fields" {
		t.Fatalf("unexpected provenance: %+v", point)
	}
}

func TestP89StructuredPublicCollectorKeepsVerifiedFieldsWhenCapitalFlowBlocked(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/qt/stock/fflow/daykline/get", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "capital flow blocked", http.StatusBadGateway)
	})
	mux.HandleFunc("/api/data/v1/get", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"success":true,"result":{"data":[{"SECURITY_CODE":"600000","TOTAL_OPERATE_INCOME":46573000000,"PARENT_NETPROFIT":17861000000,"SJLTZ":1.49,"NOTICE_DATE":"2026-04-30 00:00:00"}]}}`))
	})
	mux.HandleFunc("/marketdata/tradedata/queryMargin.do", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"pageHelp":{"data":[{"opDate":"20260618","rzye":1495185799372},{"opDate":"20260617","rzye":1488211916696}]}}`))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	point, err := P89StructuredPublicCollector{
		EastmoneyPush2BaseURL:      server.URL,
		EastmoneyH5BaseURL:         server.URL,
		EastmoneyDatacenterBaseURL: server.URL,
		SSEQueryBaseURL:            server.URL,
		HTTPClient:                 server.Client(),
	}.FetchMarketData(context.Background(), "600000")
	if err != nil {
		t.Fatalf("capital-flow blocker must not discard verified margin/financial fields: %v", err)
	}

	readback := P88NormalizeStructuredDataMetadata(point.Metadata)
	if readback.CapitalFlow != nil {
		t.Fatalf("capital flow must remain absent when provider is blocked: %+v", readback.CapitalFlow)
	}
	if readback.MarginFinancing == nil || readback.MarginFinancing.Date != "2026-06-18" {
		t.Fatalf("expected margin financing readback despite capital-flow blocker: %+v", readback.MarginFinancing)
	}
	if readback.ConstituentFinancial == nil || readback.ConstituentFinancial.DisclosureDate != "2026-04-30" {
		t.Fatalf("expected financial readback despite capital-flow blocker: %+v", readback.ConstituentFinancial)
	}
	health, _ := point.Metadata["p34_source_health"].(map[string]any)
	if health["capital_flow"] == nil || health["margin_financing"] == nil || health["constituent_financials"] == nil {
		t.Fatalf("expected source-health metadata for passed and blocked categories: %+v", health)
	}
}

func TestP90StructuredPublicCollectorUsesH5CapitalFlowWhenPush2Blocked(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/qt/stock/fflow/daykline/get", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "push2 blocked", http.StatusBadGateway)
	})
	mux.HandleFunc("/dc/ZJLX/getDBHistoryData", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("secid") != "1.600000" {
			t.Fatalf("expected secid for 600000, got %q", r.URL.Query().Get("secid"))
		}
		_, _ = w.Write([]byte(`{"rc":0,"data":{"code":"600000","market":1,"name":"浦发银行","klines":["2026-06-18,-74801057.0,9.09,-1.62","2026-06-22,11895999.0,9.16,0.77"]}}`))
	})
	mux.HandleFunc("/api/data/v1/get", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"success":true,"result":{"data":[{"SECURITY_CODE":"600000","TOTAL_OPERATE_INCOME":46573000000,"PARENT_NETPROFIT":17861000000,"SJLTZ":1.49,"NOTICE_DATE":"2026-04-30 00:00:00"}]}}`))
	})
	mux.HandleFunc("/marketdata/tradedata/queryMargin.do", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"pageHelp":{"data":[{"opDate":"20260618","rzye":1495185799372},{"opDate":"20260617","rzye":1488211916696}]}}`))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	point, err := P89StructuredPublicCollector{
		EastmoneyPush2BaseURL:      server.URL,
		EastmoneyH5BaseURL:         server.URL,
		EastmoneyDatacenterBaseURL: server.URL,
		SSEQueryBaseURL:            server.URL,
		HTTPClient:                 server.Client(),
	}.FetchMarketData(context.Background(), "600000")
	if err != nil {
		t.Fatalf("fetch H5 structured capital flow: %v", err)
	}

	readback := P88NormalizeStructuredDataMetadata(point.Metadata)
	if readback.CapitalFlow == nil {
		t.Fatalf("expected H5 capital-flow readback")
	}
	if readback.CapitalFlow.Date != "2026-06-22" || readback.CapitalFlow.NetInflow != 11895999 || readback.CapitalFlow.NetOutflow != 0 || readback.CapitalFlow.RawNetFlow != 11895999 {
		t.Fatalf("unexpected H5 capital-flow fields: %+v", readback.CapitalFlow)
	}
}
