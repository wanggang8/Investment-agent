package handler

import (
	"net/http"
	"strconv"
	"strings"

	"investment-agent/internal/application/service"
	"investment-agent/internal/pkg/apperr"
)

// ListDecisionLoops returns read-only decision loop explanations.
func (a *App) ListDecisionLoops(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	filter := service.DecisionLoopListFilter{Symbol: strings.TrimSpace(r.URL.Query().Get("symbol"))}
	if rawLimit := strings.TrimSpace(r.URL.Query().Get("limit")); rawLimit != "" {
		limit, err := strconv.Atoi(rawLimit)
		if err != nil || limit <= 0 {
			WriteHandlerError(w, requestID, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "limit must be a positive integer"))
			return
		}
		filter.Limit = limit
	}
	resp, err := a.DecisionLoopSvc.ListDecisionLoops(r.Context(), filter)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, resp)
}

// GetDecisionLoop returns one read-only decision loop explanation.
func (a *App) GetDecisionLoop(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	decisionID := strings.TrimSpace(r.PathValue("decision_id"))
	item, err := a.DecisionLoopSvc.GetDecisionLoop(r.Context(), decisionID)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, item)
}
