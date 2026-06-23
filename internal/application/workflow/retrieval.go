package workflow

import (
	"context"

	"investment-agent/internal/domain/model"
)

// RetrievalRequest 是证据检索输入，当前按标的检索。
type RetrievalRequest struct {
	Symbol string
}

// RetrievalResult 保留命中证据和可审计的降级原因。
type RetrievalResult struct {
	EvidenceSet    model.EvidenceSet
	OutputRef      string
	DegradedReason string
	QualitySummary RetrievalQualitySummary
}

type RetrievalQualitySummary struct {
	QuerySummary            string `json:"query_summary,omitempty"`
	TopK                    int    `json:"top_k"`
	Status                  string `json:"status,omitempty"`
	IndexHealth             string `json:"index_health,omitempty"`
	IndexFreshness          string `json:"index_freshness,omitempty"`
	FallbackSource          string `json:"fallback_source,omitempty"`
	SourceConsistencyStatus string `json:"source_consistency_status,omitempty"`
	DegradedReason          string `json:"degraded_reason,omitempty"`
}

// RetrievalService 封装 VecLite/RAG 检索；实现可降级到 SQLite 摘要。
type RetrievalService interface {
	RetrieveEvidence(ctx context.Context, req RetrievalRequest) (RetrievalResult, error)
}
