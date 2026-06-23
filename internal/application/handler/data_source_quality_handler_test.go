package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"investment-agent/internal/application/dto"
)

func TestGetDataSourceQualityRegressionFixtureReturnsSanitizedCases(t *testing.T) {
	app, db := testApp(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/data-source-quality/regression", nil)
	w := httptest.NewRecorder()

	app.GetDataSourceQualityRegression(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data dto.DataSourceQualityRegressionResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.Mode != "fixture" || body.Data.Status != "passed" || len(body.Data.Cases) != 6 {
		t.Fatalf("unexpected fixture response: %+v", body.Data)
	}
	if body.Data.Policy.Verdict != "passed" || body.Data.Policy.ReleaseGate != "pass" {
		t.Fatalf("expected fixture policy pass, got %+v", body.Data.Policy)
	}
	for _, forbidden := range []string{"sk-123456789012", "/Users/private", "select    *    from", "prompt:", "raw HTTP", "BEGIN RSA PRIVATE KEY"} {
		if strings.Contains(w.Body.String(), forbidden) {
			t.Fatalf("fixture API leaked %q: %s", forbidden, w.Body.String())
		}
	}
	assertHandlerTableCount(t, db, "audit_events", 0)
}

func assertHandlerTableCount(t *testing.T, db *sql.DB, table string, want int) {
	t.Helper()
	var got int
	if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&got); err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	if got != want {
		t.Fatalf("expected %s count %d, got %d", table, want, got)
	}
}

func TestGetDataSourceQualityRegressionCurrentReadsLocalSourceHealth(t *testing.T) {
	app, db := testApp(t)
	metrics := `{"source_name":"csindex","source_level":"A","source_type":"index_basic","captured_at":"2026-06-06T01:00:00Z","metadata":{"p34_source_health":{"index_constituents":"fresh","index_valuation_files":"parse_error"},"p34_data_categories":["index_constituents","index_valuation_files"]}}`
	_, err := db.Exec(`INSERT INTO market_snapshots (market_snapshot_id,symbol,trade_date,close_price,turnover_rate,pe_percentile,pb_percentile,volume_percentile,volatility_percentile,liquidity_state,sentiment_state,market_metrics_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`, "market_p48_current", "000300", "2026-06-05", 0, 0, 0, 0, 0, 0, "normal", "neutral", metrics, "2026-06-06T01:00:00Z")
	if err != nil {
		t.Fatalf("seed market: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/data-source-quality/regression?mode=current&symbol=000300", nil)
	w := httptest.NewRecorder()

	app.GetDataSourceQualityRegression(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data dto.DataSourceQualityRegressionResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.Mode != "current" || body.Data.Status != "degraded" || len(body.Data.Cases) != 2 || body.Data.MissingCategories[0] != "index_valuation_files" {
		t.Fatalf("unexpected current response: %+v", body.Data)
	}
	if body.Data.Policy.Verdict != "blocked" || body.Data.Policy.ReleaseGate != "block" {
		t.Fatalf("expected current core degradation to block policy, got %+v", body.Data.Policy)
	}
}

func TestGetDataSourceQualityRegressionCurrentMissingHealthDegrades(t *testing.T) {
	app, _ := testApp(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/data-source-quality/regression?mode=current", nil)
	w := httptest.NewRecorder()

	app.GetDataSourceQualityRegression(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data dto.DataSourceQualityRegressionResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.Status != "degraded" || body.Data.Cases[0].ActualFreshness != "missing" || body.Data.MissingCategories[0] != "p34_source_health" {
		t.Fatalf("expected degraded missing source health, got %+v", body.Data)
	}
	if body.Data.Policy.Verdict != "blocked" || body.Data.Policy.ReleaseGate != "block" {
		t.Fatalf("expected missing source health to block policy, got %+v", body.Data.Policy)
	}
}

func TestGetDataSourceQualityRegressionRejectsUnsupportedMode(t *testing.T) {
	app, _ := testApp(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/data-source-quality/regression?mode=real", nil)
	w := httptest.NewRecorder()

	app.GetDataSourceQualityRegression(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetDataQualityGateResolutionRequiresResolutionAndDoesNotAudit(t *testing.T) {
	app, db := testApp(t)
	seedBlockedCurrentSourceHealth(t, db)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/data-source-quality/gate-resolution?symbol=000300", nil)
	w := httptest.NewRecorder()

	app.GetDataQualityGateResolution(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data dto.DataQualityGateResolutionCheck `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.ReleaseClaimState != "requires_resolution" || body.Data.Policy.Verdict != "blocked" || body.Data.CleanDataClaimAllowed {
		t.Fatalf("unexpected resolution check: %+v", body.Data)
	}
	if body.Data.PolicyFingerprint == "" || body.Data.ActiveResolution != nil {
		t.Fatalf("expected fingerprint without active resolution, got %+v", body.Data)
	}
	assertHandlerTableCount(t, db, "audit_events", 0)
}

func TestCreateDataQualityGateResolutionScopeExclusionAndRetire(t *testing.T) {
	app, db := testApp(t)
	seedBlockedCurrentSourceHealth(t, db)
	body := `{"symbol":"000300","resolution_type":"scope_exclusion","scope":"本次 release clean claim 排除 current local data health","reason":"本地 current 源存在解析降级，发布材料只能声明 fixture 与既有功能验收","release_impact":"不得声明 current data healthy；只允许 limited release claim","evidence_ref":"docs/release/acceptance/p66"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/data-source-quality/resolutions", bytes.NewBufferString(body))
	req.Header.Set("X-Request-ID", "req_p67_create")
	w := httptest.NewRecorder()

	app.CreateDataQualityGateResolution(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var createResp struct {
		Data dto.DataQualityGateResolutionCheck `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if createResp.Data.ReleaseClaimState != "resolved_with_scope_exclusion" || createResp.Data.ActiveResolution == nil || createResp.Data.CleanDataClaimAllowed {
		t.Fatalf("unexpected create response: %+v", createResp.Data)
	}
	resolutionID := createResp.Data.ActiveResolution.ResolutionID
	assertHandlerTableCount(t, db, "data_quality_gate_resolutions", 1)
	assertHandlerTableCount(t, db, "audit_events", 1)

	duplicateReq := httptest.NewRequest(http.MethodPost, "/api/v1/data-source-quality/resolutions", bytes.NewBufferString(body))
	duplicateReq.Header.Set("X-Request-ID", "req_p67_duplicate")
	duplicateW := httptest.NewRecorder()

	app.CreateDataQualityGateResolution(duplicateW, duplicateReq)

	if duplicateW.Code != http.StatusOK {
		t.Fatalf("expected duplicate 200, got %d body=%s", duplicateW.Code, duplicateW.Body.String())
	}
	assertHandlerTableCount(t, db, "data_quality_gate_resolutions", 1)
	assertHandlerTableCount(t, db, "audit_events", 2)

	retireReq := httptest.NewRequest(http.MethodPost, "/api/v1/data-source-quality/resolutions/"+resolutionID+"/retire", nil)
	retireReq.SetPathValue("resolution_id", resolutionID)
	retireReq.Header.Set("X-Request-ID", "req_p67_retire")
	retireW := httptest.NewRecorder()

	app.RetireDataQualityGateResolution(retireW, retireReq)

	if retireW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", retireW.Code, retireW.Body.String())
	}
	var retireResp struct {
		Data dto.DataQualityGateResolutionCheck `json:"data"`
	}
	if err := json.Unmarshal(retireW.Body.Bytes(), &retireResp); err != nil {
		t.Fatalf("decode retire response: %v", err)
	}
	if retireResp.Data.ReleaseClaimState != "requires_resolution" || retireResp.Data.ActiveResolution != nil {
		t.Fatalf("unexpected retire response: %+v", retireResp.Data)
	}
	assertHandlerTableCount(t, db, "audit_events", 3)
}

func TestCreateDataQualityGateResolutionRejectsBlockedWaiver(t *testing.T) {
	app, db := testApp(t)
	seedBlockedCurrentSourceHealth(t, db)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/data-source-quality/resolutions", bytes.NewBufferString(`{"symbol":"000300","resolution_type":"waiver","scope":"current data","reason":"accept degraded current source","release_impact":"claim healthy"}`))
	w := httptest.NewRecorder()

	app.CreateDataQualityGateResolution(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
	assertHandlerTableCount(t, db, "data_quality_gate_resolutions", 0)
	assertHandlerTableCount(t, db, "audit_events", 0)
}

func seedBlockedCurrentSourceHealth(t *testing.T, db *sql.DB) {
	t.Helper()
	metrics := `{"source_name":"csindex","source_level":"A","source_type":"index_basic","captured_at":"2026-06-06T01:00:00Z","metadata":{"p34_source_health":{"index_constituents":{"freshness":"fresh","data_date":"2026-06-05"},"index_valuation_files":{"freshness":"parse_error","data_date":"2026-06-05","failure_category":"parse_error","diagnostic_preview":"valuation file parse failed"}},"p34_data_categories":["index_constituents","index_valuation_files"]}}`
	_, err := db.Exec(`INSERT INTO market_snapshots (market_snapshot_id,symbol,trade_date,close_price,turnover_rate,pe_percentile,pb_percentile,volume_percentile,volatility_percentile,liquidity_state,sentiment_state,market_metrics_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`, "market_p67_current", "000300", "2026-06-05", 0, 0, 0, 0, 0, 0, "normal", "neutral", metrics, "2026-06-06T01:00:00Z")
	if err != nil {
		t.Fatalf("seed market: %v", err)
	}
}
