package dto

// EvidenceDTO 是证据列表和决策证据链的统一展示结构。
type EvidenceDTO struct {
	EvidenceID                      string  `json:"evidence_id"`
	SourceName                      string  `json:"source_name"`
	SourceLevel                     string  `json:"source_level"`
	EvidenceRole                    string  `json:"evidence_role,omitempty"`
	PublishedAt                     string  `json:"published_at,omitempty"`
	CapturedAt                      string  `json:"captured_at,omitempty"`
	OriginalURL                     string  `json:"original_url,omitempty"`
	Summary                         string  `json:"summary"`
	ContentHash                     string  `json:"content_hash,omitempty"`
	TimeWeight                      float64 `json:"time_weight,omitempty"`
	RelevanceScore                  float64 `json:"relevance_score,omitempty"`
	IndependentSourceCount          int     `json:"independent_source_count,omitempty"`
	HighGradeIndependentSourceCount int     `json:"high_grade_independent_source_count,omitempty"`
}

// EvidenceRefreshRequest 描述本地证据刷新范围，不生成交易建议。
type EvidenceRefreshRequest struct {
	Symbol            string   `json:"symbol,omitempty"`
	RefreshScope      string   `json:"refresh_scope"`
	IncludeBackground bool     `json:"include_background"`
	Sources           []string `json:"sources,omitempty"`
}

type EvidenceRefreshResponse struct {
	IntelligenceItemCount int      `json:"intelligence_item_count"`
	SummaryCount          int      `json:"summary_count"`
	RAGChunkCount         int      `json:"rag_chunk_count"`
	VerificationCount     int      `json:"verification_count"`
	IndexStatus           string   `json:"index_status"`
	FailedReason          string   `json:"failed_reason,omitempty"`
	AuditEventIDs         []string `json:"audit_event_ids"`
}

type SourceVerificationDTO struct {
	VerificationID                  string   `json:"verification_id"`
	VerificationStatus              string   `json:"verification_status"`
	IndependentSourceCount          int      `json:"independent_source_count"`
	HighGradeIndependentSourceCount int      `json:"high_grade_independent_source_count"`
	HighestSourceLevel              string   `json:"highest_source_level"`
	LatestPublishedAt               string   `json:"latest_published_at"`
	EvidenceIDs                     []string `json:"evidence_ids"`
}

type RebuildIndexResponse struct {
	IndexedCount  int            `json:"indexed_count"`
	SkippedCount  int            `json:"skipped_count"`
	LastRebuildAt string         `json:"last_rebuild_at,omitempty"`
	IndexHealth   IndexHealthDTO `json:"index_health"`
	AuditEventIDs []string       `json:"audit_event_ids"`
}

type IndexHealthDTO struct {
	Status         string `json:"status"`
	Path           string `json:"path,omitempty"`
	Version        int    `json:"version,omitempty"`
	ChunkCount     int    `json:"chunk_count,omitempty"`
	Rebuildable    bool   `json:"rebuildable"`
	DegradedReason string `json:"degraded_reason,omitempty"`
}
