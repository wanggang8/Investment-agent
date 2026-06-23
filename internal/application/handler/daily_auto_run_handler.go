package handler

import (
	"net/http"
	"net/url"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

const dailyAutoRunSafetyNote = "仅记录本地每日刷新和纪律评估结果，不会自动执行交易，需人工复核。"

func (a *App) GetDailyAutoRunStatus(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	out := dto.DailyAutoRunStatusResponse{
		Enabled:    a.Deps.DailyAutoRunConfig.Enabled,
		RunTime:    a.Deps.DailyAutoRunConfig.RunTime,
		Timezone:   a.Deps.DailyAutoRunConfig.Timezone,
		Scope:      a.Deps.DailyAutoRunConfig.Scope,
		Status:     "disabled",
		SafetyNote: dailyAutoRunSafetyNote,
	}
	state, err := a.QuerySvc.LatestDailyAutoRunState(r.Context())
	if err != nil {
		if apperr.IsCode(err, apperr.CodeNotFound) {
			if out.Enabled {
				out.Status = "scheduled"
			}
			writeOK(w, requestID, out)
			return
		}
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, dailyAutoRunStatusDTO(out, state))
}

func dailyAutoRunStatusDTO(base dto.DailyAutoRunStatusResponse, state repository.DailyAutoRunState) dto.DailyAutoRunStatusResponse {
	base.Status = state.Status
	base.RunID = state.RunID
	base.IdempotencyKey = state.IdempotencyKey
	base.LocalDate = state.LocalDate
	base.Scope = state.Scope
	base.LastRunAt = state.LastRunAt
	base.NextRunAt = state.NextRunAt
	base.FailureCode = state.FailureCode
	base.FailureReason = state.FailureReason
	if state.IdempotencyKey != "" {
		escapedKey := url.QueryEscape(state.IdempotencyKey)
		base.LatestNotificationLink = "/notifications?source_id=" + escapedKey
		base.LatestAuditLink = "/audit?input_ref=" + escapedKey
	}
	if state.RunID != "" {
		base.LatestDecisionLink = "/decisions?request_id=" + url.QueryEscape(state.RunID)
	}
	if state.FailureCode == "missing_prerequisites" {
		base.MissingAction = "请先录入本地账户、组合和当前持仓，再等待下一次自动运行或手动触发。"
	}
	return base
}
