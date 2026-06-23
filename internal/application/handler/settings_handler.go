package handler

import (
	"encoding/json"
	"net/http"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/pkg/apperr"
)

// GetSystemSettings 返回本地依赖状态，不返回完整密钥。
func (a *App) GetSystemSettings(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	writeOK(w, requestID, a.QuerySvc.SystemStatus(r.Context(), a.VectorIndex != nil || a.Deps.VectorIndexWriter != nil, a.Deps.AnalystService != nil))
}

// UpdateSystemSettings 只保存通知、页面偏好和普通数据源；规则类变更必须走规则提案。
func (a *App) UpdateSystemSettings(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		WriteHandlerError(w, requestID, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "invalid json body"))
		return
	}
	if hasForbiddenRuleSetting(raw) {
		WriteHandlerError(w, requestID, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "规则类设置必须通过规则提案变更"))
		return
	}
	var req dto.SystemSettingsDTO
	body, _ := json.Marshal(raw)
	if err := json.Unmarshal(body, &req); err != nil {
		WriteHandlerError(w, requestID, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "invalid settings body"))
		return
	}
	if err := a.SettingsSvc.UpdateSystemSettings(r.Context(), requestID, req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, req)
}

func hasForbiddenRuleSetting(raw map[string]json.RawMessage) bool {
	for _, key := range []string{"rule_thresholds", "rule_priority", "arbitration_priority", "sop", "sop_config", "position_limits", "cash_min_ratio"} {
		if _, ok := raw[key]; ok {
			return true
		}
	}
	return false
}

// GetCapabilitySettings 返回能力圈配置。
func (a *App) GetCapabilitySettings(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	out, err := a.SettingsSvc.GetCapabilitySettings(r.Context())
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, out)
}

// UpdateCapabilitySettings 只保存能力圈；规则阈值和 SOP 不在此接口更新。
func (a *App) UpdateCapabilitySettings(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	var req dto.CapabilitySettingsDTO
	if err := decodeJSON(r, &req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	symbols := firstNonEmptySlice(req.Symbols, req.AllowedSymbols)
	if hasOverlap(symbols, req.ExcludedSymbols) {
		WriteHandlerError(w, requestID, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "symbols 与 excluded_symbols 不能重叠"))
		return
	}
	if err := a.SettingsSvc.UpdateCapabilitySettings(r.Context(), requestID, req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, req)
}

func firstNonEmptySlice(primary, fallback []string) []string {
	if len(primary) > 0 {
		return primary
	}
	return fallback
}

func hasOverlap(left, right []string) bool {
	seen := map[string]bool{}
	for _, item := range left {
		seen[item] = true
	}
	for _, item := range right {
		if seen[item] {
			return true
		}
	}
	return false
}
