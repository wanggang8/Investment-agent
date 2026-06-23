package repository

import "context"

// IntelligenceItem 是外部信源采集后的原始情报元信息。
type IntelligenceItem struct {
	IntelligenceID string
	SourceName     string
	SourceLevel    string
	OriginalURL    string
	PublishedAt    string
	CapturedAt     string
	ContentHash    string
	RawTitle       string
	RawTextRef     string
	CreatedAt      string
}

// IntelligenceSummary 是清洗后的结构化情报摘要，可进入证据链或 RAG 索引。
type IntelligenceSummary struct {
	SummaryID                       string
	IntelligenceID                  string
	Symbol                          string
	Entity                          string
	EventType                       string
	ImpactDirection                 string
	Summary                         string
	SourceLevel                     string
	EvidenceRole                    string
	TimeWeight                      float64
	RelevanceScore                  float64
	VerificationGroupID             string
	VerificationStatus              string
	VerificationEvidenceIDsJSON     string
	VerificationEvidenceRole        string
	VerificationEventType           string
	VerificationHighestSourceLevel  string
	IndependentSourceCount          int
	HighGradeIndependentSourceCount int
	SourceName                      string
	OriginalURL                     string
	PublishedAt                     string
	CapturedAt                      string
	ContentHash                     string
	CreatedAt                       string
}

// RAGChunk 是可重建 VecLite 索引的文本块元数据。
type RAGChunk struct {
	ChunkID          string
	SummaryID        string
	Symbol           string
	ChunkText        string
	ChunkHash        string
	VectorID         string
	VectorCollection string
	EmbeddingModel   string
	EmbeddingVersion string
	IndexVersion     string
	IndexStatus      string
	IndexedAt        string
	MetadataJSON     string
	CreatedAt        string
}

// SourceVerification 保存同一事件的多源验证结果。
type SourceVerification struct {
	VerificationID                  string
	VerificationGroupID             string
	EventID                         string
	Symbol                          string
	EventType                       string
	EvidenceRole                    string
	VerificationStatus              string
	IndependentSourceCount          int
	HighGradeIndependentSourceCount int
	HighestSourceLevel              string
	LatestPublishedAt               string
	EvidenceIDsJSON                 string
	CreatedAt                       string
}

// IntelligenceRepository 定义情报摘要、RAG 文本块与多源验证的持久化边界。
type IntelligenceRepository interface {
	SaveIntelligenceItem(ctx context.Context, item IntelligenceItem) error
	GetIntelligenceItem(ctx context.Context, intelligenceID string) (IntelligenceItem, error)
	SaveIntelligenceSummary(ctx context.Context, summary IntelligenceSummary, chunks []RAGChunk) error
	GetIntelligenceSummary(ctx context.Context, summaryID string) (IntelligenceSummary, []RAGChunk, error)
	ListEvidenceSummaries(ctx context.Context) ([]IntelligenceSummary, error)
	SaveSourceVerification(ctx context.Context, verification SourceVerification) error
	GetSourceVerification(ctx context.Context, verificationID string) (SourceVerification, error)
	GetLatestSourceVerification(ctx context.Context) (SourceVerification, error)
	GetLatestSourceVerificationByFilter(ctx context.Context, symbol, eventID string) (SourceVerification, error)
	UpdateRAGChunksIndexStatus(ctx context.Context, chunkIDs []string, status string) error
	ListRAGChunks(ctx context.Context) ([]RAGChunk, error)
	CountRAGChunks(ctx context.Context) (int, error)
}
