package handler

import (
	"net/http"

	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

// GetReviewSummary 返回复盘页聚合摘要，帮助用户查看错误案例、规则演进和审计追踪。
func (a *App) GetReviewSummary(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	period := r.URL.Query().Get("period")
	if period != "" && period != "monthly" && period != "quarterly" {
		WriteHandlerError(w, requestID, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "period 只支持 monthly 或 quarterly"))
		return
	}
	summary, err := a.QuerySvc.ReviewSummary(r.Context(), period)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	if summary.DegradedCount > 0 || summary.MissingEvidenceCount > 0 {
		if err := a.NotificationSvc.AppendNotification(r.Context(), repository.Notification{Type: "review_degraded", Severity: "warning", Title: "复盘存在降级或缺失证据", Message: summary.OpsStatus.Explanation, SourceType: "review_summary", SourceID: summary.Period}); err != nil {
			WriteHandlerError(w, requestID, err)
			return
		}
	}
	writeOK(w, requestID, summary)
}
