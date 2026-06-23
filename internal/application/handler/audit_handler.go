package handler

import (
	"net/http"

	"investment-agent/internal/application/dto"
)

// ListAuditEvents 返回本地审计时间线。
func (a *App) ListAuditEvents(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	events, err := a.QuerySvc.ListAuditEvents(r.Context())
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	items := make([]dto.AuditEventDTO, 0, len(events))
	for _, event := range events {
		items = append(items, dto.AuditEventDTO{AuditEventID: event.AuditEventID, EventID: event.AuditEventID, RequestID: event.RequestID, DecisionID: event.DecisionID, WorkflowType: event.WorkflowType, NodeName: event.NodeName, Actor: event.Actor, Action: event.Action, NodeAction: event.NodeAction, ProposalID: event.ProposalID, ConfirmationID: event.ConfirmationID, ErrorCaseID: event.ErrorCaseID, Status: event.Status, ErrorCode: event.ErrorCode, BeforeState: event.BeforeState, AfterState: event.AfterState, RuleVersion: event.RuleVersion, SnapshotID: event.SnapshotID, InputRefType: event.InputRefType, InputRef: event.InputRef, OutputRefType: event.OutputRefType, OutputRef: event.OutputRef, CreatedAt: event.CreatedAt})
	}
	writeOK(w, requestID, dto.PageResult[dto.AuditEventDTO]{Items: items, Total: len(items)})
}
