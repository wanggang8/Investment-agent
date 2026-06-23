package handler

import (
	"net/http"
	"strings"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

func (a *App) ListRiskAlerts(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	out, err := a.QuerySvc.ListRiskAlerts(r.Context(), riskAlertFilterFromRequest(r))
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, out)
}

func (a *App) GetRiskAlert(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	out, err := a.QuerySvc.GetRiskAlert(r.Context(), r.PathValue("alert_id"))
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, out)
}

func (a *App) UpdateRiskAlertLifecycle(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	var req dto.RiskAlertLifecycleRequest
	if err := decodeJSON(r, &req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	status := model.RiskSOPStatus(strings.TrimSpace(req.Status))
	if !status.Valid() {
		WriteHandlerError(w, requestID, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "invalid risk alert status"))
		return
	}
	updated, err := a.RiskAlertSvc.UpdateRiskAlertLifecycle(r.Context(), r.PathValue("alert_id"), status, req.Reason)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	out, err := a.QuerySvc.GetRiskAlert(r.Context(), updated.AlertID)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, out)
}

func riskAlertFilterFromRequest(r *http.Request) repository.RiskAlertFilter {
	q := r.URL.Query()
	filter := repository.RiskAlertFilter{Symbol: strings.TrimSpace(q.Get("symbol"))}
	for _, raw := range q["status"] {
		for _, item := range strings.Split(raw, ",") {
			status := model.RiskSOPStatus(strings.TrimSpace(item))
			if status.Valid() {
				filter.SOPStatuses = append(filter.SOPStatuses, status)
			}
		}
	}
	return filter
}
