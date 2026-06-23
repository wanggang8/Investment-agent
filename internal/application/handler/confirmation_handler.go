package handler

import (
	"net/http"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/pkg/apperr"
)

// CreateConfirmation 写入用户线下处理结果。
// planned/watch 只写确认和审计；executed_manually/marked_error 在一个 SQLite 事务内写入相关事实。
func (a *App) CreateConfirmation(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	decisionID := r.PathValue("decision_id")
	var req dto.ConfirmationRequest
	if err := decodeJSON(r, &req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	if req.ConfirmationType == "" {
		WriteHandlerError(w, requestID, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "confirmation_type 不能为空"))
		return
	}
	out, err := a.ConfirmationSvc.Confirm(r.Context(), requestID, decisionID, req)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, out)
}
