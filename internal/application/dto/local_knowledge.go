package dto

type LocalKnowledgeImportRow struct {
	Title     string   `json:"title"`
	Text      string   `json:"text"`
	Symbol    string   `json:"symbol,omitempty"`
	SourceURL string   `json:"source_url,omitempty"`
	AsOfDate  string   `json:"as_of_date,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

type LocalKnowledgeImportValidationRequest struct {
	SourceLabel   string                    `json:"source_label"`
	DefaultSymbol string                    `json:"default_symbol,omitempty"`
	Rows          []LocalKnowledgeImportRow `json:"rows"`
}

type LocalKnowledgeImportConfirmRequest struct {
	ImportBatchID string                    `json:"import_batch_id"`
	ConfirmReason string                    `json:"confirm_reason"`
	SourceLabel   string                    `json:"source_label"`
	DefaultSymbol string                    `json:"default_symbol,omitempty"`
	Rows          []LocalKnowledgeImportRow `json:"rows"`
}

type LocalKnowledgeImportRisk struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

type LocalKnowledgeImportRowResult struct {
	RowNumber      int                        `json:"row_number"`
	Status         string                     `json:"status"`
	Symbol         string                     `json:"symbol,omitempty"`
	TitlePreview   string                     `json:"title_preview"`
	TextPreview    string                     `json:"text_preview"`
	ContentHash    string                     `json:"content_hash"`
	EstimatedChunk int                        `json:"estimated_chunk_count"`
	Risks          []LocalKnowledgeImportRisk `json:"risks"`
}

type LocalKnowledgeImportIndexPlan struct {
	RAGChunkCount int    `json:"rag_chunk_count"`
	IndexStatus   string `json:"index_status"`
}

type LocalKnowledgeImportValidationSummary struct {
	TotalCount    int `json:"total_count"`
	ValidCount    int `json:"valid_count"`
	BlockingCount int `json:"blocking_count"`
	WarningCount  int `json:"warning_count"`
}

type LocalKnowledgeImportValidationResponse struct {
	ImportBatchID string                                `json:"import_batch_id"`
	Summary       LocalKnowledgeImportValidationSummary `json:"summary"`
	Rows          []LocalKnowledgeImportRowResult       `json:"rows"`
	IndexPlan     LocalKnowledgeImportIndexPlan         `json:"index_plan"`
	SafetyNote    string                                `json:"safety_note"`
}

type LocalKnowledgeImportConfirmResponse struct {
	ImportBatchID         string   `json:"import_batch_id"`
	IntelligenceItemCount int      `json:"intelligence_item_count"`
	SummaryCount          int      `json:"summary_count"`
	RAGChunkCount         int      `json:"rag_chunk_count"`
	VerificationCount     int      `json:"verification_count"`
	IndexStatus           string   `json:"index_status"`
	AuditEventIDs         []string `json:"audit_event_ids"`
	SafetyNote            string   `json:"safety_note"`
}
