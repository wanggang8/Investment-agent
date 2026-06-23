package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"investment-agent/internal/application/dto"
)

func TestUpdateSystemSettingsPersistsPagePreference(t *testing.T) {
	app, db := testApp(t)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/settings", bytes.NewBufferString(`{"notification_enabled":true,"page_preference":"compact","data_sources":["official"]}`))
	req.Header.Set("X-Request-ID", "req_settings")
	w := httptest.NewRecorder()

	app.UpdateSystemSettings(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var raw string
	if err := db.QueryRow(`SELECT notification_config_json FROM user_settings WHERE settings_id LIKE 'settings_%'`).Scan(&raw); err != nil {
		t.Fatalf("read settings: %v", err)
	}
	if raw == "{}" || !contains(raw, "compact") {
		t.Fatalf("expected page preference persisted, got %s", raw)
	}
	assertCount(t, db, "audit_events", 1)
}

func TestUpdateSystemSettingsRejectsRuleFields(t *testing.T) {
	app, db := testApp(t)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/settings", bytes.NewBufferString(`{"notification_enabled":true,"rule_thresholds":{"x":1}}`))
	req.Header.Set("X-Request-ID", "req_settings_rule")
	w := httptest.NewRecorder()

	app.UpdateSystemSettings(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
	assertCount(t, db, "user_settings", 0)
}

func TestUpdateAndGetCapabilitySettingsUseContractFields(t *testing.T) {
	app, _ := testApp(t)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/settings/capability", bytes.NewBufferString(`{"asset_types":["ETF"],"symbols":["510300"],"excluded_symbols":["159915"],"strategy_scope":["discipline_review"]}`))
	req.Header.Set("X-Request-ID", "req_capability")
	w := httptest.NewRecorder()

	app.UpdateCapabilitySettings(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/settings/capability", bytes.NewBuffer(nil))
	getReq.Header.Set("X-Request-ID", "req_capability_get")
	getW := httptest.NewRecorder()
	app.GetCapabilitySettings(getW, getReq)
	if getW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", getW.Code, getW.Body.String())
	}
	for _, key := range []string{"capability_id", "asset_types", "symbols", "excluded_symbols", "strategy_scope", "updated_at"} {
		if !contains(getW.Body.String(), key) {
			t.Fatalf("missing %s in %s", key, getW.Body.String())
		}
	}
}

func TestGetSystemSettingsReflectsPersistedDataSources(t *testing.T) {
	app, _ := testApp(t)
	putReq := httptest.NewRequest(http.MethodPut, "/api/v1/settings", bytes.NewBufferString(`{"notification_enabled":true,"page_preference":"compact","data_sources":["manual"]}`))
	putReq.Header.Set("X-Request-ID", "req_settings_source")
	putW := httptest.NewRecorder()
	app.UpdateSystemSettings(putW, putReq)
	if putW.Code != http.StatusOK {
		t.Fatalf("expected update 200, got %d body=%s", putW.Code, putW.Body.String())
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/settings/system", nil)
	req.Header.Set("X-Request-ID", "req_settings_system")
	w := httptest.NewRecorder()
	app.GetSystemSettings(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data dto.SystemStatusDTO `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Data.DataSources) != 1 || body.Data.DataSources[0] != "manual" || body.Data.SQLiteStatus != "ok" || body.Data.VecLiteStatus == "unavailable" {
		t.Fatalf("expected persisted settings and configured deps, got %+v", body.Data)
	}
}

func TestUpdateCapabilitySettingsRejectsOverlappedSymbols(t *testing.T) {
	app, db := testApp(t)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/settings/capability", bytes.NewBufferString(`{"asset_types":["ETF"],"symbols":["510300"],"excluded_symbols":["510300"]}`))
	req.Header.Set("X-Request-ID", "req_capability_overlap")
	w := httptest.NewRecorder()

	app.UpdateCapabilitySettings(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
	assertCount(t, db, "capability_configs", 0)
}

func contains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && bytes.Contains([]byte(s), []byte(sub)))
}
