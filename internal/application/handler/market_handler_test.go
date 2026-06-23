package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/application/workflow"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/pkg/apperr"
	"investment-agent/internal/pkg/idgen"
	"investment-agent/pkg/httputil"
)

func TestGetLatestMarketSnapshotReturnsFreshnessAndMetrics(t *testing.T) {
	app, db := testApp(t)
	_, err := db.Exec(`INSERT INTO market_snapshots (market_snapshot_id,symbol,trade_date,close_price,turnover_rate,pe_percentile,pb_percentile,volume_percentile,volatility_percentile,liquidity_state,sentiment_state,market_metrics_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`, "market_detail", "510300", "2026-01-02", 4.2, 0.3, 31, 41, 51, 21, "normal", "neutral", `{"close_price":1,"turnover_rate":9}`, "2026-01-02T01:00:00Z")
	if err != nil {
		t.Fatalf("seed market: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/market/latest", nil)
	req.Header.Set("X-Request-ID", "req_market_latest")
	w := httptest.NewRecorder()

	app.GetLatestMarketSnapshot(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data dto.MarketSnapshotDTO `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.TradeDate != "2026-01-02" || body.Data.DataStatus != "fresh" || !floatClose(body.Data.ClosePrice, 4.2) || !floatClose(body.Data.TurnoverRate, 0.3) || body.Data.MarketMetrics["close_price"] != float64(1) || body.Data.MarketMetrics["turnover_rate"] != float64(9) {
		t.Fatalf("expected freshness and metrics, got %+v", body.Data)
	}
}

func TestGetLatestMarketSnapshotFiltersBySymbol(t *testing.T) {
	app, db := testApp(t)
	if _, err := db.Exec(`INSERT INTO market_snapshots (market_snapshot_id,symbol,trade_date,close_price,turnover_rate,pe_percentile,pb_percentile,volume_percentile,volatility_percentile,liquidity_state,sentiment_state,market_metrics_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`, "market_old_symbol", "510300", "2026-01-02", 4.2, 0.3, 31, 41, 51, 21, "normal", "neutral", `{}`, "2026-01-02T01:00:00Z"); err != nil {
		t.Fatalf("seed first market: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO market_snapshots (market_snapshot_id,symbol,trade_date,close_price,turnover_rate,pe_percentile,pb_percentile,volume_percentile,volatility_percentile,liquidity_state,sentiment_state,market_metrics_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`, "market_target_symbol", "159915", "2026-01-01", 2.5, 0.2, 11, 12, 13, 14, "normal", "neutral", `{}`, "2026-01-01T01:00:00Z"); err != nil {
		t.Fatalf("seed target market: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/market/latest?symbol=159915", nil)
	req.Header.Set("X-Request-ID", "req_market_latest_symbol")
	w := httptest.NewRecorder()

	app.GetLatestMarketSnapshot(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data dto.MarketSnapshotDTO `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.MarketSnapshotID != "market_target_symbol" || body.Data.Symbol != "159915" {
		t.Fatalf("expected symbol-filtered latest snapshot, got %+v", body.Data)
	}
}

func TestGetMarketSourceHealthReturnsP34Categories(t *testing.T) {
	app, db := testApp(t)
	metrics := `{"source_name":"csindex","source_level":"A","source_type":"index_basic","trade_date":"2026-06-05","captured_at":"2026-06-06T01:00:00Z","metadata":{"p34_source_health":{"index_constituents":"fresh","index_weights":"no_data","index_valuation_files":"parse_error"},"p34_data_categories":["index_constituents","index_weights","index_valuation_files"]}}`
	_, err := db.Exec(`INSERT INTO market_snapshots (market_snapshot_id,symbol,trade_date,close_price,turnover_rate,pe_percentile,pb_percentile,volume_percentile,volatility_percentile,liquidity_state,sentiment_state,market_metrics_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`, "market_p34_health", "000300", "2026-06-05", 0, 0, 0, 0, 0, 0, "normal", "neutral", metrics, "2026-06-06T01:00:00Z")
	if err != nil {
		t.Fatalf("seed market: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/market/source-health?symbol=000300", nil)
	req.Header.Set("X-Request-ID", "req_market_source_health")
	w := httptest.NewRecorder()

	app.GetMarketSourceHealth(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data dto.SourceHealthResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Data.Sources) != 3 {
		t.Fatalf("expected three P34 source health rows, got %+v", body.Data)
	}
	if body.Data.Sources[0].SourceName != "csindex" || body.Data.Sources[0].DataCategory != "index_constituents" || body.Data.Sources[0].Freshness != "fresh" || body.Data.Sources[0].LastSuccessAt == "" {
		t.Fatalf("unexpected first health row: %+v", body.Data.Sources[0])
	}
	if body.Data.Sources[2].Freshness != "parse_error" || body.Data.Sources[2].FailureCategory != "parse_error" || len(body.Data.Sources[2].AffectedSymbols) != 1 || body.Data.Sources[2].AffectedSymbols[0] != "000300" {
		t.Fatalf("unexpected parse-error health row: %+v", body.Data.Sources[2])
	}
}

func TestGetMarketSourceHealthReturnsStructuredP34Health(t *testing.T) {
	app, db := testApp(t)
	metrics := `{"source_name":"csindex","source_level":"A","source_type":"index_basic","captured_at":"2026-06-06T01:00:00Z","metadata":{"p34_source_health":{"index_constituents":{"freshness":"fresh","last_success_at":"2026-06-06T01:00:00Z","data_date":"2026-06-05","affected_symbols":["000300"],"source_level":"A"},"index_weights":{"freshness":"source_unavailable","last_success_at":"2026-06-01T01:00:00Z","last_failure_at":"2026-06-06T01:00:00Z","failure_category":"source_unavailable","data_date":"2026-06-05","affected_symbols":["000300"],"source_level":"A"}},"p34_data_categories":["index_constituents","index_weights"]}}`
	_, err := db.Exec(`INSERT INTO market_snapshots (market_snapshot_id,symbol,trade_date,close_price,turnover_rate,pe_percentile,pb_percentile,volume_percentile,volatility_percentile,liquidity_state,sentiment_state,market_metrics_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`, "market_p34_structured_health", "000300", "2026-06-05", 0, 0, 0, 0, 0, 0, "normal", "neutral", metrics, "2026-06-06T01:00:00Z")
	if err != nil {
		t.Fatalf("seed market: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/market/source-health?symbol=000300", nil)
	w := httptest.NewRecorder()

	app.GetMarketSourceHealth(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data dto.SourceHealthResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Data.Sources) != 2 {
		t.Fatalf("expected two P34 health rows, got %+v", body.Data)
	}
	failed := body.Data.Sources[1]
	if failed.Freshness != "source_unavailable" || failed.LastSuccessAt != "2026-06-01T01:00:00Z" || failed.LastFailureAt != "2026-06-06T01:00:00Z" || failed.FailureCategory != "source_unavailable" || failed.DataDate != "2026-06-05" || failed.SourceLevel != "A" {
		t.Fatalf("structured health fields not preserved: %+v", failed)
	}
}

func TestRefreshMarketRejectsUnsupportedAsOfDate(t *testing.T) {
	app, _ := testApp(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/market/refresh", bytes.NewBufferString(`{"symbols":["510300"],"as_of_date":"2026-01-02"}`))
	req.Header.Set("X-Request-ID", "req_market_unsupported_date")
	w := httptest.NewRecorder()

	app.RefreshMarket(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestRefreshMarketRejectsInvalidAsOfDate(t *testing.T) {
	app, _ := testApp(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/market/refresh", bytes.NewBufferString(`{"symbols":["510300"],"as_of_date":"2026/01/02"}`))
	req.Header.Set("X-Request-ID", "req_market_bad_date")
	w := httptest.NewRecorder()

	app.RefreshMarket(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestRefreshMarketPartialSuccessReturnsFailedSymbols(t *testing.T) {
	app, db := testApp(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/market/refresh", bytes.NewBufferString(`{"symbols":["510300",""]}`))
	req.Header.Set("X-Request-ID", "req_market_partial")
	w := httptest.NewRecorder()

	app.RefreshMarket(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		httputil.Envelope
		Data dto.MarketRefreshResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.RefreshedCount != 1 || len(body.Data.LatestSnapshotIDs) != 1 || len(body.Data.FailedSymbols) != 1 || len(body.Data.AuditEventIDs) != 1 {
		t.Fatalf("unexpected market refresh response: %+v", body.Data)
	}
	var auditStatus string
	if err := db.QueryRow(`SELECT status FROM audit_events WHERE audit_event_id=?`, body.Data.AuditEventIDs[0]).Scan(&auditStatus); err != nil {
		t.Fatalf("read degraded audit: %v", err)
	}
	if auditStatus != "degraded" {
		t.Fatalf("expected degraded audit, got %s", auditStatus)
	}
	for _, auditID := range body.Data.AuditEventIDs {
		var exists int
		if err := db.QueryRow(`SELECT COUNT(*) FROM audit_events WHERE audit_event_id=?`, auditID).Scan(&exists); err != nil {
			t.Fatalf("check audit id: %v", err)
		}
		if exists != 1 {
			t.Fatalf("response audit id does not exist: %s", auditID)
		}
	}
	var notificationType, sourceType, sourceID string
	if err := db.QueryRow(`SELECT type,COALESCE(source_type,''),COALESCE(source_id,'') FROM notifications WHERE read_at IS NULL ORDER BY created_at DESC LIMIT 1`).Scan(&notificationType, &sourceType, &sourceID); err != nil {
		t.Fatalf("read partial market notification: %v", err)
	}
	if notificationType != "data_source_failure" || sourceType != "market_refresh" || sourceID != "partial_data_source_failure" {
		t.Fatalf("expected partial data source failure notification, got type=%s source=%s/%s", notificationType, sourceType, sourceID)
	}
}

func TestRefreshMarketInvalidJSONReturnsBadRequest(t *testing.T) {
	app, _ := testApp(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/market/refresh", bytes.NewBufferString(`{"symbols"`))
	req.Header.Set("X-Request-ID", "req_market_bad_json")
	w := httptest.NewRecorder()

	app.RefreshMarket(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestRefreshMarketWriteFailureReturnsMarketSnapshotWriteFailed(t *testing.T) {
	app, db := testApp(t)
	workflow.SetWorkflowIDGenerator(idgen.NewFixedGenerator(map[string][]string{"market": {"market_conflict"}, "audit": {"audit_conflict"}}))
	defer workflow.SetWorkflowIDGenerator(idgen.NewGenerator())
	err := app.Deps.MarketRepo.SaveMarketSnapshot(context.Background(), model.MarketSnapshot{MarketSnapshotID: "market_conflict", Symbol: "510300", LiquidityState: model.LiquidityNormal, SentimentState: model.SentimentNeutral}, "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed market snapshot: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/market/refresh", bytes.NewBufferString(`{"symbols":["510300"]}`))
	req.Header.Set("X-Request-ID", "req_market_write_failed")
	w := httptest.NewRecorder()

	app.RefreshMarket(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d body=%s", w.Code, w.Body.String())
	}
	var body httputil.Envelope
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Error == nil || body.Error.Code != string(apperr.CodeMarketSnapshotWriteFailed) {
		t.Fatalf("expected MARKET_SNAPSHOT_WRITE_FAILED, got %+v", body.Error)
	}
	var status, errorCode string
	if err := db.QueryRow(`SELECT status,error_code FROM audit_events WHERE request_id='req_market_write_failed' AND error_code=? ORDER BY created_at DESC LIMIT 1`, string(apperr.CodeMarketSnapshotWriteFailed)).Scan(&status, &errorCode); err != nil {
		t.Fatalf("read failure audit: %v", err)
	}
	if status != "failed" || errorCode != string(apperr.CodeMarketSnapshotWriteFailed) {
		t.Fatalf("unexpected failure audit status=%s code=%s", status, errorCode)
	}
	var auditCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM audit_events WHERE request_id='req_market_write_failed' AND error_code=? AND status='failed' AND action='refresh_market_data' AND node_name='MarketRefreshGraph'`, string(apperr.CodeMarketSnapshotWriteFailed)).Scan(&auditCount); err != nil {
		t.Fatalf("count failure audit: %v", err)
	}
	if auditCount != 1 {
		t.Fatalf("expected one graph failure audit, got %d", auditCount)
	}
	var successCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM audit_events WHERE request_id='req_market_write_failed' AND status='success'`).Scan(&successCount); err != nil {
		t.Fatalf("count success audit: %v", err)
	}
	if successCount != 0 {
		t.Fatalf("expected no success audit on write failure, got %d", successCount)
	}
}

func TestRefreshMarketAllFailedDedupesActiveNotification(t *testing.T) {
	app, db := testApp(t)
	for _, requestID := range []string{"req_market_failed_a", "req_market_failed_b"} {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/market/refresh", bytes.NewBufferString(`{"symbols":[""]}`))
		req.Header.Set("X-Request-ID", requestID)
		w := httptest.NewRecorder()
		app.RefreshMarket(w, req)
		if w.Code != http.StatusServiceUnavailable {
			t.Fatalf("expected 503, got %d body=%s", w.Code, w.Body.String())
		}
	}
	var unread int
	if err := db.QueryRow(`SELECT COUNT(*) FROM notifications WHERE type='data_source_failure' AND read_at IS NULL`).Scan(&unread); err != nil {
		t.Fatalf("read notification count: %v", err)
	}
	if unread != 1 {
		t.Fatalf("expected active data source failure notification to be deduped, got %d", unread)
	}
}

func TestRefreshMarketAllFailedReturnsDataSourceUnavailable(t *testing.T) {
	app, db := testApp(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/market/refresh", bytes.NewBufferString(`{"symbols":[""]}`))
	req.Header.Set("X-Request-ID", "req_market_failed")
	w := httptest.NewRecorder()

	app.RefreshMarket(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d body=%s", w.Code, w.Body.String())
	}
	var status, errorCode string
	if err := db.QueryRow(`SELECT status,error_code FROM audit_events WHERE request_id='req_market_failed' ORDER BY created_at DESC LIMIT 1`).Scan(&status, &errorCode); err != nil {
		t.Fatalf("read unavailable audit: %v", err)
	}
	if status != "failed" || errorCode != string(apperr.CodeDataSourceUnavailable) {
		t.Fatalf("unexpected unavailable audit status=%s code=%s", status, errorCode)
	}
	var unread int
	if err := db.QueryRow(`SELECT COUNT(*) FROM notifications WHERE type='data_source_failure' AND read_at IS NULL`).Scan(&unread); err != nil {
		t.Fatalf("read notification count: %v", err)
	}
	if unread != 1 {
		t.Fatalf("expected one unread data source failure notification, got %d", unread)
	}
}
