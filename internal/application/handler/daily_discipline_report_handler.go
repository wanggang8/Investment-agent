package handler

import (
	"net/http"
	"strconv"
	"time"

	"investment-agent/internal/pkg/apperr"
)

func (a *App) GetTodayDailyDisciplineReport(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	out, err := a.QuerySvc.TodayDailyDisciplineReport(r.Context(), time.Now())
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, out)
}

func (a *App) ListDailyDisciplineReports(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	limit := 0
	if raw := r.URL.Query().Get("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed <= 0 {
			WriteHandlerError(w, requestID, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "limit 必须为正整数"))
			return
		}
		limit = parsed
	}
	out, err := a.QuerySvc.ListDailyDisciplineReports(r.Context(), r.URL.Query().Get("status"), limit)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, out)
}

func (a *App) GetDailyDisciplineReport(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	out, err := a.QuerySvc.GetDailyDisciplineReport(r.Context(), r.PathValue("report_id"))
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, out)
}
