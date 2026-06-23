package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"investment-agent/internal/application/service"
)

func TestGetKnowledgeReadinessReturnsSanitizedReadiness(t *testing.T) {
	app, db := testApp(t)
	seedKnowledgeReadinessHandlerFacts(t, db, "510300", readinessHandlerHealthJSON("fresh"))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/knowledge-readiness?symbol=510300", nil)
	w := httptest.NewRecorder()

	app.GetKnowledgeReadiness(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data service.KnowledgeReadinessResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.Status != "ready" || body.Data.SymbolProfile.TrackedIndexSymbol != "000300" {
		t.Fatalf("unexpected readiness response: %+v", body.Data)
	}
	if len(body.Data.KnowledgeReferences) < 5 || len(body.Data.DataDependencies) == 0 || !strings.Contains(body.Data.LLMContextSummary, "master.graham.margin_of_safety") {
		t.Fatalf("expected knowledge references, dependencies and LLM context summary, got %+v", body.Data)
	}
	for _, forbidden := range []string{"sk-", "raw HTTP", "BEGIN RSA PRIVATE KEY", "/Users/private", "prompt:"} {
		if strings.Contains(w.Body.String(), forbidden) {
			t.Fatalf("knowledge readiness API leaked %q: %s", forbidden, w.Body.String())
		}
	}
	assertHandlerTableCount(t, db, "audit_events", 0)
}

func TestGetKnowledgeReadinessBlocksUnknownSymbolWithoutFabricatingProfile(t *testing.T) {
	app, _ := testApp(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/knowledge-readiness?symbol=999999", nil)
	w := httptest.NewRecorder()

	app.GetKnowledgeReadiness(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data service.KnowledgeReadinessResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.Status != "blocked" || body.Data.SymbolProfile.Known {
		t.Fatalf("expected blocked unknown symbol readiness, got %+v", body.Data)
	}
}

func TestGetKnowledgeReadinessReturnsDegradedDependencyImpacts(t *testing.T) {
	app, db := testApp(t)
	seedKnowledgeReadinessHandlerFactsWithVerification(t, db, "510300", readinessHandlerHealthJSONForCategories(map[string]string{
		"symbol_profile":        "fresh",
		"fund_profile":          "fresh",
		"tracked_index":         "fresh",
		"market_price":          "fresh",
		"valuation_percentiles": "parse_error",
		"liquidity":             "missing",
		"sentiment_proxy":       "fresh",
	}), "satisfied", 1, 1)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/knowledge-readiness?symbol=510300", nil)
	w := httptest.NewRecorder()

	app.GetKnowledgeReadiness(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data service.KnowledgeReadinessResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.Status != "degraded" {
		t.Fatalf("expected degraded readiness, got %+v", body.Data)
	}
	for _, category := range []string{"valuation_percentiles", "liquidity", "formal_evidence"} {
		dep := handlerReadinessDependencyByCategory(body.Data.DataDependencies, category)
		if dep.Status != "degraded" {
			t.Fatalf("expected %s degraded, got %+v", category, dep)
		}
	}
	for _, want := range []string{"不得声明安全边际", "预期收益精度不足", "不得输出大额或市价式行动建议", "不生成交易确认", "不得输出交易确认"} {
		if !strings.Contains(w.Body.String(), want) {
			t.Fatalf("expected degraded API response to contain %q, got %s", want, w.Body.String())
		}
	}
	assertHandlerTableCount(t, db, "audit_events", 0)
}

func seedKnowledgeReadinessHandlerFacts(t *testing.T, db execDB, symbol string, metrics string) {
	t.Helper()
	seedKnowledgeReadinessHandlerFactsWithVerification(t, db, symbol, metrics, "satisfied", 3, 2)
}

func seedKnowledgeReadinessHandlerFactsWithVerification(t *testing.T, db execDB, symbol string, metrics string, verificationStatus string, independent int, highGrade int) {
	t.Helper()
	if _, err := db.Exec(`INSERT INTO market_snapshots (market_snapshot_id,symbol,trade_date,close_price,turnover_rate,pe_percentile,pb_percentile,volume_percentile,volatility_percentile,liquidity_state,sentiment_state,market_metrics_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`, "market_p74_handler_"+symbol, symbol, "2026-06-19", 4.75, 1.2, 28, 35, 40, 30, "normal", "neutral", metrics, "2026-06-19T08:00:00Z"); err != nil {
		t.Fatalf("seed market: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO source_verifications (verification_id,verification_group_id,event_id,symbol,event_type,evidence_role,verification_status,independent_source_count,high_grade_independent_source_count,highest_source_level,latest_published_at,evidence_ids_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`, "verify_p74_handler_"+symbol, "group_p74_handler", "event_p74_handler", symbol, "normal", "formal", verificationStatus, independent, highGrade, "A", "2026-06-19T08:00:00Z", `["sum_a","sum_b"]`, "2026-06-19T08:00:00Z"); err != nil {
		t.Fatalf("seed verification: %v", err)
	}
}

type execDB interface {
	Exec(query string, args ...any) (sql.Result, error)
}

func readinessHandlerHealthJSON(freshness string) string {
	return readinessHandlerHealthJSONForCategories(map[string]string{
		"symbol_profile":        freshness,
		"fund_profile":          freshness,
		"tracked_index":         freshness,
		"market_price":          freshness,
		"valuation_percentiles": freshness,
		"liquidity":             freshness,
		"sentiment_proxy":       freshness,
	})
}

func readinessHandlerHealthJSONForCategories(freshness map[string]string) string {
	return `{"source_name":"p74_handler_fixture","source_level":"A","source_type":"readiness_fixture","captured_at":"2026-06-19T08:00:00Z","metadata":{"p34_source_health":{"symbol_profile":{"freshness":"` + freshness["symbol_profile"] + `","data_date":"2026-06-19","affected_symbols":["510300"],"source_level":"A"},"fund_profile":{"freshness":"` + freshness["fund_profile"] + `","data_date":"2026-06-19","affected_symbols":["510300"],"source_level":"B"},"tracked_index":{"freshness":"` + freshness["tracked_index"] + `","data_date":"2026-06-19","affected_symbols":["000300"],"source_level":"A"},"market_price":{"freshness":"` + freshness["market_price"] + `","data_date":"2026-06-19","affected_symbols":["510300"],"source_level":"B"},"valuation_percentiles":{"freshness":"` + freshness["valuation_percentiles"] + `","data_date":"2026-06-19","affected_symbols":["000300"],"source_level":"A"},"liquidity":{"freshness":"` + freshness["liquidity"] + `","data_date":"2026-06-19","affected_symbols":["510300"],"source_level":"B"},"sentiment_proxy":{"freshness":"` + freshness["sentiment_proxy"] + `","data_date":"2026-06-19","affected_symbols":["510300"],"source_level":"C"}},"p34_data_categories":["symbol_profile","fund_profile","tracked_index","market_price","valuation_percentiles","liquidity","sentiment_proxy"]}}`
}

func handlerReadinessDependencyByCategory(items []service.KnowledgeDataDependency, category string) service.KnowledgeDataDependency {
	for _, item := range items {
		if item.Category == category {
			return item
		}
	}
	return service.KnowledgeDataDependency{}
}
