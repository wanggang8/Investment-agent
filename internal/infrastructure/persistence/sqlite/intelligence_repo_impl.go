package sqlite

import (
	"context"
	"database/sql"
	"strings"

	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

// IntelligenceRepository 是情报摘要、RAG 文本块和多源验证表的 SQLite 实现。
type IntelligenceRepository struct{ db dbtx }

// NewIntelligenceRepository 创建情报仓储实例。
func NewIntelligenceRepository(db *sql.DB) *IntelligenceRepository {
	return &IntelligenceRepository{db: db}
}

// SaveIntelligenceItem 保存原始情报元信息。
func (r *IntelligenceRepository) SaveIntelligenceItem(ctx context.Context, item repository.IntelligenceItem) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO intelligence_items (intelligence_id,source_name,source_level,original_url,published_at,captured_at,content_hash,raw_title,raw_text_ref,created_at) VALUES (?,?,?,?,?,?,?,?,?,?)`, item.IntelligenceID, item.SourceName, item.SourceLevel, nullString(item.OriginalURL), nullString(item.PublishedAt), item.CapturedAt, item.ContentHash, nullString(item.RawTitle), nullString(item.RawTextRef), item.CreatedAt)
	return apperr.FromRepositoryError(err)
}

// GetIntelligenceItem 读取原始情报元信息。
func (r *IntelligenceRepository) GetIntelligenceItem(ctx context.Context, id string) (repository.IntelligenceItem, error) {
	var item repository.IntelligenceItem
	err := r.db.QueryRowContext(ctx, `SELECT intelligence_id,source_name,source_level,COALESCE(original_url,''),COALESCE(published_at,''),captured_at,content_hash,COALESCE(raw_title,''),COALESCE(raw_text_ref,''),created_at FROM intelligence_items WHERE intelligence_id=?`, id).Scan(&item.IntelligenceID, &item.SourceName, &item.SourceLevel, &item.OriginalURL, &item.PublishedAt, &item.CapturedAt, &item.ContentHash, &item.RawTitle, &item.RawTextRef, &item.CreatedAt)
	return item, apperr.FromRepositoryError(err)
}

// SaveIntelligenceSummary 在同一事务中保存情报摘要和对应 RAG 文本块。
func (r *IntelligenceRepository) SaveIntelligenceSummary(ctx context.Context, s repository.IntelligenceSummary, chunks []repository.RAGChunk) error {
	err := withTx(ctx, r.db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO intelligence_summary (summary_id,intelligence_id,symbol,entity,event_type,impact_direction,summary,source_level,evidence_role,time_weight,relevance_score,verification_group_id,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`, s.SummaryID, s.IntelligenceID, nullString(s.Symbol), nullString(s.Entity), nullString(s.EventType), nullString(s.ImpactDirection), s.Summary, s.SourceLevel, s.EvidenceRole, s.TimeWeight, s.RelevanceScore, nullString(s.VerificationGroupID), s.CreatedAt)
		if err != nil {
			return err
		}
		for _, c := range chunks {
			chunk := c
			if chunk.Symbol == "" {
				chunk.Symbol = s.Symbol
			}
			_, err = tx.ExecContext(ctx, `INSERT INTO rag_chunks (chunk_id,summary_id,chunk_text,chunk_hash,vector_id,vector_collection,embedding_model,embedding_version,index_version,index_status,indexed_at,metadata_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`, chunk.ChunkID, chunk.SummaryID, chunk.ChunkText, chunk.ChunkHash, nullString(chunk.VectorID), nullString(chunk.VectorCollection), nullString(chunk.EmbeddingModel), nullString(chunk.EmbeddingVersion), nullString(chunk.IndexVersion), chunk.IndexStatus, nullString(chunk.IndexedAt), nullString(chunk.MetadataJSON), chunk.CreatedAt)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return apperr.FromRepositoryError(err)
}

// GetIntelligenceSummary 读取情报摘要及其 RAG 文本块。
func (r *IntelligenceRepository) GetIntelligenceSummary(ctx context.Context, id string) (repository.IntelligenceSummary, []repository.RAGChunk, error) {
	var s repository.IntelligenceSummary
	err := r.db.QueryRowContext(ctx, `SELECT summary_id,intelligence_id,COALESCE(symbol,''),COALESCE(entity,''),COALESCE(event_type,''),COALESCE(impact_direction,''),summary,source_level,evidence_role,COALESCE(time_weight,0),COALESCE(relevance_score,0),COALESCE(verification_group_id,''),created_at FROM intelligence_summary WHERE summary_id=?`, id).Scan(&s.SummaryID, &s.IntelligenceID, &s.Symbol, &s.Entity, &s.EventType, &s.ImpactDirection, &s.Summary, &s.SourceLevel, &s.EvidenceRole, &s.TimeWeight, &s.RelevanceScore, &s.VerificationGroupID, &s.CreatedAt)
	if err != nil {
		return s, nil, apperr.FromRepositoryError(err)
	}
	rows, err := r.db.QueryContext(ctx, `SELECT r.chunk_id,r.summary_id,COALESCE(s.symbol,''),r.chunk_text,r.chunk_hash,COALESCE(r.vector_id,''),COALESCE(r.vector_collection,''),COALESCE(r.embedding_model,''),COALESCE(r.embedding_version,''),COALESCE(r.index_version,''),r.index_status,COALESCE(r.indexed_at,''),COALESCE(r.metadata_json,''),r.created_at FROM rag_chunks r JOIN intelligence_summary s ON s.summary_id=r.summary_id WHERE r.summary_id=? ORDER BY r.chunk_id`, id)
	if err != nil {
		return s, nil, apperr.FromRepositoryError(err)
	}
	defer rows.Close()
	var chunks []repository.RAGChunk
	for rows.Next() {
		var c repository.RAGChunk
		if err := rows.Scan(&c.ChunkID, &c.SummaryID, &c.Symbol, &c.ChunkText, &c.ChunkHash, &c.VectorID, &c.VectorCollection, &c.EmbeddingModel, &c.EmbeddingVersion, &c.IndexVersion, &c.IndexStatus, &c.IndexedAt, &c.MetadataJSON, &c.CreatedAt); err != nil {
			return s, nil, apperr.FromRepositoryError(err)
		}
		chunks = append(chunks, c)
	}
	return s, chunks, apperr.FromRepositoryError(rows.Err())
}

// ListEvidenceSummaries reads evidence summaries for list pages.
func (r *IntelligenceRepository) ListEvidenceSummaries(ctx context.Context) ([]repository.IntelligenceSummary, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT s.summary_id,s.intelligence_id,COALESCE(s.symbol,''),COALESCE(s.entity,''),s.source_level,s.evidence_role,COALESCE(s.event_type,''),s.summary,COALESCE(s.time_weight,0),COALESCE(s.relevance_score,0),COALESCE(v.verification_status,''),COALESCE(v.evidence_ids_json,''),COALESCE(v.evidence_role,''),COALESCE(v.event_type,''),COALESCE(v.highest_source_level,''),COALESCE(v.independent_source_count,0),COALESCE(v.high_grade_independent_source_count,0),COALESCE(i.source_name,''),COALESCE(i.original_url,''),COALESCE(i.published_at,''),COALESCE(i.captured_at,''),COALESCE(i.content_hash,''),s.created_at FROM intelligence_summary s LEFT JOIN source_verifications v ON v.verification_id=(SELECT sv.verification_id FROM source_verifications sv WHERE sv.verification_group_id=s.verification_group_id ORDER BY sv.created_at DESC, sv.verification_id DESC LIMIT 1) LEFT JOIN intelligence_items i ON i.intelligence_id=s.intelligence_id ORDER BY s.created_at DESC`)
	if err != nil {
		return nil, apperr.FromRepositoryError(err)
	}
	defer rows.Close()
	var out []repository.IntelligenceSummary
	for rows.Next() {
		var s repository.IntelligenceSummary
		if err := rows.Scan(&s.SummaryID, &s.IntelligenceID, &s.Symbol, &s.Entity, &s.SourceLevel, &s.EvidenceRole, &s.EventType, &s.Summary, &s.TimeWeight, &s.RelevanceScore, &s.VerificationStatus, &s.VerificationEvidenceIDsJSON, &s.VerificationEvidenceRole, &s.VerificationEventType, &s.VerificationHighestSourceLevel, &s.IndependentSourceCount, &s.HighGradeIndependentSourceCount, &s.SourceName, &s.OriginalURL, &s.PublishedAt, &s.CapturedAt, &s.ContentHash, &s.CreatedAt); err != nil {
			return nil, apperr.FromRepositoryError(err)
		}
		out = append(out, s)
	}
	return out, apperr.FromRepositoryError(rows.Err())
}

// SaveSourceVerification 保存多源验证结果。
func (r *IntelligenceRepository) SaveSourceVerification(ctx context.Context, v repository.SourceVerification) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO source_verifications (verification_id,verification_group_id,event_id,symbol,event_type,evidence_role,verification_status,independent_source_count,high_grade_independent_source_count,highest_source_level,latest_published_at,evidence_ids_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?) ON CONFLICT(verification_id) DO UPDATE SET verification_group_id=excluded.verification_group_id,event_id=excluded.event_id,symbol=excluded.symbol,event_type=excluded.event_type,evidence_role=excluded.evidence_role,verification_status=excluded.verification_status,independent_source_count=excluded.independent_source_count,high_grade_independent_source_count=excluded.high_grade_independent_source_count,highest_source_level=excluded.highest_source_level,latest_published_at=excluded.latest_published_at,evidence_ids_json=excluded.evidence_ids_json,created_at=excluded.created_at`, v.VerificationID, v.VerificationGroupID, v.EventID, nullString(v.Symbol), nullString(v.EventType), v.EvidenceRole, v.VerificationStatus, v.IndependentSourceCount, v.HighGradeIndependentSourceCount, nullString(v.HighestSourceLevel), nullString(v.LatestPublishedAt), nullString(v.EvidenceIDsJSON), v.CreatedAt)
	return apperr.FromRepositoryError(err)
}

// GetSourceVerification 读取多源验证结果。
func (r *IntelligenceRepository) GetSourceVerification(ctx context.Context, id string) (repository.SourceVerification, error) {
	var v repository.SourceVerification
	err := r.db.QueryRowContext(ctx, `SELECT verification_id,verification_status,independent_source_count,COALESCE(high_grade_independent_source_count,0),COALESCE(highest_source_level,''),COALESCE(latest_published_at,''),COALESCE(evidence_ids_json,''),created_at FROM source_verifications WHERE verification_id=?`, id).Scan(&v.VerificationID, &v.VerificationStatus, &v.IndependentSourceCount, &v.HighGradeIndependentSourceCount, &v.HighestSourceLevel, &v.LatestPublishedAt, &v.EvidenceIDsJSON, &v.CreatedAt)
	return v, apperr.FromRepositoryError(err)
}

// GetLatestSourceVerification reads the latest source verification.
func (r *IntelligenceRepository) GetLatestSourceVerification(ctx context.Context) (repository.SourceVerification, error) {
	var v repository.SourceVerification
	err := r.db.QueryRowContext(ctx, `SELECT verification_id,verification_status,independent_source_count,COALESCE(high_grade_independent_source_count,0),COALESCE(highest_source_level,''),COALESCE(latest_published_at,''),COALESCE(evidence_ids_json,'[]'),created_at FROM source_verifications ORDER BY created_at DESC, verification_id DESC LIMIT 1`).Scan(&v.VerificationID, &v.VerificationStatus, &v.IndependentSourceCount, &v.HighGradeIndependentSourceCount, &v.HighestSourceLevel, &v.LatestPublishedAt, &v.EvidenceIDsJSON, &v.CreatedAt)
	return v, apperr.FromRepositoryError(err)
}

func (r *IntelligenceRepository) GetLatestSourceVerificationByFilter(ctx context.Context, symbol, eventID string) (repository.SourceVerification, error) {
	query := `SELECT verification_id,verification_group_id,COALESCE(event_id,''),COALESCE(symbol,''),COALESCE(event_type,''),evidence_role,verification_status,independent_source_count,COALESCE(high_grade_independent_source_count,0),COALESCE(highest_source_level,''),COALESCE(latest_published_at,''),COALESCE(evidence_ids_json,'[]'),created_at FROM source_verifications`
	args := []any{}
	where := []string{}
	if symbol != "" {
		where = append(where, "symbol=?")
		args = append(args, symbol)
	}
	if eventID != "" {
		where = append(where, "event_id=?")
		args = append(args, eventID)
	}
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += " ORDER BY created_at DESC, verification_id DESC LIMIT 1"
	var v repository.SourceVerification
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&v.VerificationID, &v.VerificationGroupID, &v.EventID, &v.Symbol, &v.EventType, &v.EvidenceRole, &v.VerificationStatus, &v.IndependentSourceCount, &v.HighGradeIndependentSourceCount, &v.HighestSourceLevel, &v.LatestPublishedAt, &v.EvidenceIDsJSON, &v.CreatedAt)
	return v, apperr.FromRepositoryError(err)
}

// UpdateRAGChunksIndexStatus updates index status for selected chunks.
func (r *IntelligenceRepository) UpdateRAGChunksIndexStatus(ctx context.Context, chunkIDs []string, status string) error {
	if len(chunkIDs) == 0 {
		return nil
	}
	args := make([]any, 0, len(chunkIDs)+1)
	args = append(args, status)
	placeholders := "?"
	for i, chunkID := range chunkIDs {
		if i > 0 {
			placeholders += ",?"
		}
		args = append(args, chunkID)
	}
	_, err := r.db.ExecContext(ctx, `UPDATE rag_chunks SET index_status=? WHERE chunk_id IN (`+placeholders+`)`, args...)
	return apperr.FromRepositoryError(err)
}

// ListRAGChunks 读取所有可用于重建 VecLite 索引的文本块。
func (r *IntelligenceRepository) ListRAGChunks(ctx context.Context) ([]repository.RAGChunk, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT r.chunk_id,r.summary_id,COALESCE(s.symbol,''),r.chunk_text,r.chunk_hash,COALESCE(r.vector_id,''),COALESCE(r.vector_collection,''),COALESCE(r.embedding_model,''),COALESCE(r.embedding_version,''),COALESCE(r.index_version,''),r.index_status,COALESCE(r.indexed_at,''),COALESCE(r.metadata_json,''),r.created_at FROM rag_chunks r JOIN intelligence_summary s ON s.summary_id=r.summary_id ORDER BY r.created_at ASC, r.chunk_id ASC`)
	if err != nil {
		return nil, apperr.FromRepositoryError(err)
	}
	defer rows.Close()
	var chunks []repository.RAGChunk
	for rows.Next() {
		var c repository.RAGChunk
		if err := rows.Scan(&c.ChunkID, &c.SummaryID, &c.Symbol, &c.ChunkText, &c.ChunkHash, &c.VectorID, &c.VectorCollection, &c.EmbeddingModel, &c.EmbeddingVersion, &c.IndexVersion, &c.IndexStatus, &c.IndexedAt, &c.MetadataJSON, &c.CreatedAt); err != nil {
			return nil, apperr.FromRepositoryError(err)
		}
		chunks = append(chunks, c)
	}
	return chunks, apperr.FromRepositoryError(rows.Err())
}

// CountRAGChunks counts rebuildable RAG chunks.
func (r *IntelligenceRepository) CountRAGChunks(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM rag_chunks`).Scan(&count)
	return count, apperr.FromRepositoryError(err)
}
