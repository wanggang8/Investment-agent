package handler

import (
	"net/http"

	"investment-agent/internal/application/service"
)

// GetKnowledgeReadiness returns read-only built-in knowledge and data readiness for a symbol.
func (a *App) GetKnowledgeReadiness(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	out, err := a.KnowledgeReadinessSvc.Evaluate(r.Context(), service.KnowledgeReadinessRequest{Symbol: r.URL.Query().Get("symbol")})
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, out)
}
