package handler

import (
	"context"
	"net/http"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/application/service"
	"investment-agent/internal/application/workflow"
	"investment-agent/internal/domain/repository"
)

// RefreshEvidence 同步执行证据核查工作流，写入 SQLite 事实并返回索引状态。
func (a *App) RefreshEvidence(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	var req dto.EvidenceRefreshRequest
	if err := decodeJSON(r, &req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	if req.Symbol == "" {
		req.Symbol = "market"
	}
	sources := req.Sources
	if len(sources) == 0 {
		sources = []string{"official", "exchange"}
	}
	out, err := workflow.NewEvidenceVerificationGraphWithDependencies(a.Deps).Run(r.Context(), workflow.EvidenceVerificationInput{RequestID: requestID, Symbol: req.Symbol, Sources: sources})
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	indexStatus := "indexed"
	failedReason := ""
	indexedCount := len(out.RAGChunks)
	if out.VectorIndexStatus == string(workflow.StatusFailed) {
		indexStatus = "failed"
		failedReason = firstNonEmptyString(out.VectorIndexFailedReason, "vector index unavailable")
		indexedCount = 0
		if err := a.NotificationSvc.AppendNotification(r.Context(), repository.Notification{Type: "vector_index_failure", Severity: "warning", Title: "证据索引失败", Message: failedReason, SourceType: "evidence_refresh", SourceID: req.Symbol}); err != nil {
			WriteHandlerError(w, requestID, err)
			return
		}
	}
	auditIDs := []string{}
	for _, audit := range out.WorkflowContext.AuditEvents {
		auditIDs = append(auditIDs, audit.AuditEventID)
	}
	writeOK(w, requestID, dto.EvidenceRefreshResponse{IntelligenceItemCount: len(out.IntelligenceItems), SummaryCount: len(out.IntelligenceSummaries), RAGChunkCount: indexedCount, VerificationCount: 1, IndexStatus: indexStatus, FailedReason: failedReason, AuditEventIDs: auditIDs})
}

// ListEvidence 返回当前已保存的证据摘要。P4 先提供本地查询接口，后续 P5 负责展示筛选。
func (a *App) ListEvidence(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	summaries, err := a.QuerySvc.ListEvidenceSummaries(r.Context())
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	items := make([]dto.EvidenceDTO, 0, len(summaries))
	for _, summary := range summaries {
		items = append(items, dto.EvidenceDTO{EvidenceID: summary.SummaryID, SourceName: firstNonEmptyString(summary.SourceName, summary.Entity), SourceLevel: summary.SourceLevel, EvidenceRole: summary.EvidenceRole, PublishedAt: summary.PublishedAt, CapturedAt: summary.CapturedAt, OriginalURL: summary.OriginalURL, Summary: summary.Summary, ContentHash: summary.ContentHash, TimeWeight: summary.TimeWeight, RelevanceScore: summary.RelevanceScore, IndependentSourceCount: summary.IndependentSourceCount, HighGradeIndependentSourceCount: summary.HighGradeIndependentSourceCount})
	}
	writeOK(w, requestID, dto.PageResult[dto.EvidenceDTO]{Items: items, Total: len(items)})
}

// GetEvidenceVerification 返回最近的多源验证结果。
func (a *App) GetEvidenceVerification(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	verification, err := a.QuerySvc.LatestSourceVerificationByFilter(r.Context(), r.URL.Query().Get("symbol"), r.URL.Query().Get("event_id"))
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	out := dto.SourceVerificationDTO{VerificationID: verification.VerificationID, VerificationStatus: verification.VerificationStatus, IndependentSourceCount: verification.IndependentSourceCount, HighGradeIndependentSourceCount: verification.HighGradeIndependentSourceCount, HighestSourceLevel: verification.HighestSourceLevel, LatestPublishedAt: verification.LatestPublishedAt, EvidenceIDs: splitJSONStrings(verification.EvidenceIDsJSON)}
	writeOK(w, requestID, out)
}

// RebuildEvidenceIndex 重建可由 SQLite 摘要恢复的本地索引，不改变事实数据。
func (a *App) RebuildEvidenceIndex(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	stats, err := a.EvidenceSvc.RebuildVectorIndexWithStats(r.Context(), a.VectorIndex)
	if err != nil {
		if notifyErr := a.NotificationSvc.AppendNotification(r.Context(), repository.Notification{Type: "vector_index_failure", Severity: "warning", Title: "证据索引重建失败", Message: err.Error(), SourceType: "evidence_rebuild", SourceID: "local_vector_index"}); notifyErr != nil {
			WriteHandlerError(w, requestID, notifyErr)
			return
		}
		WriteHandlerError(w, requestID, err)
		return
	}
	auditID, err := a.EvidenceSvc.AppendRebuildAudit(r.Context(), requestID)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	health := indexHealth(r.Context(), a.VectorIndex)
	if health.Status == "" {
		health.Status = stats.Status
	}
	writeOK(w, requestID, dto.RebuildIndexResponse{IndexedCount: stats.IndexedCount, SkippedCount: stats.SkippedCount, LastRebuildAt: stats.LastRebuildAt, IndexHealth: toIndexHealthDTO(health), AuditEventIDs: []string{auditID}})
}

func toIndexHealthDTO(health service.VectorIndexHealth) dto.IndexHealthDTO {
	return dto.IndexHealthDTO{Status: health.Status, Path: health.Path, Version: health.Version, ChunkCount: health.ChunkCount, Rebuildable: health.Rebuildable, DegradedReason: health.DegradedReason}
}

func indexHealth(ctx context.Context, index service.VectorIndex) service.VectorIndexHealth {
	if healthIndex, ok := index.(interface {
		Health(context.Context) service.VectorIndexHealth
	}); ok {
		return healthIndex.Health(ctx)
	}
	if index == nil {
		return service.VectorIndexHealth{Status: service.VectorIndexHealthMissing, Rebuildable: true, DegradedReason: "VecLite 索引不可用"}
	}
	return service.VectorIndexHealth{Status: service.VectorIndexHealthHealthy, Rebuildable: true}
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
