package handler

import (
	"net/http"

	"investment-agent/internal/application/dto"
)

func (a *App) ValidateLocalKnowledgeImport(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	var req dto.LocalKnowledgeImportValidationRequest
	if err := decodeJSON(r, &req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	resp, err := a.LocalKnowledgeSvc.ValidateImport(r.Context(), req)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, resp)
}

func (a *App) ConfirmLocalKnowledgeImport(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	var req dto.LocalKnowledgeImportConfirmRequest
	if err := decodeJSON(r, &req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	resp, err := a.LocalKnowledgeSvc.ConfirmImport(r.Context(), requestID, req)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, resp)
}
