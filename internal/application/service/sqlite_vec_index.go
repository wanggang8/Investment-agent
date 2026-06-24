package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	sqlitevec "github.com/asg017/sqlite-vec-go-bindings/cgo"
	_ "github.com/mattn/go-sqlite3"

	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

type EmbeddingProvider interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}

type VectorSearchQuery struct {
	Text   string
	Symbol string
	TopK   int
}

type SemanticVectorIndex interface {
	SearchSimilar(ctx context.Context, query VectorSearchQuery) ([]repository.RAGChunk, error)
}

type SQLiteVecVectorIndexConfig struct {
	Path       string
	Dimensions int
	TopK       int
	Embedder   EmbeddingProvider
}

type SQLiteVecVectorIndex struct {
	path       string
	dimensions int
	topK       int
	embedder   EmbeddingProvider
}

var sqliteVecAutoOnce sync.Once

func NewSQLiteVecVectorIndex(cfg SQLiteVecVectorIndexConfig) *SQLiteVecVectorIndex {
	topK := cfg.TopK
	if topK <= 0 {
		topK = 8
	}
	return &SQLiteVecVectorIndex{path: cfg.Path, dimensions: cfg.Dimensions, topK: topK, embedder: cfg.Embedder}
}

func (s *SQLiteVecVectorIndex) Upsert(ctx context.Context, chunk repository.RAGChunk) error {
	if s == nil || strings.TrimSpace(s.path) == "" {
		return apperr.New(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "sqlite-vec 索引路径未配置")
	}
	if s.embedder == nil {
		return apperr.New(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "embedding provider 未配置")
	}
	vec, err := s.embedder.Embed(ctx, chunk.ChunkText)
	if err != nil {
		return apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "embedding 生成失败", err)
	}
	if err := s.validateVector(vec); err != nil {
		return err
	}
	db, err := s.open(ctx)
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "sqlite-vec 事务启动失败", err)
	}
	rollback := true
	defer func() {
		if rollback {
			_ = tx.Rollback()
		}
	}()
	indexedAt := chunk.IndexedAt
	if indexedAt == "" {
		indexedAt = time.Now().UTC().Format(time.RFC3339)
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO vector_chunks (chunk_id,summary_id,symbol,chunk_text,chunk_hash,index_status,indexed_at,metadata_json,created_at) VALUES (?,?,?,?,?,?,?,?,?)
ON CONFLICT(chunk_id) DO UPDATE SET summary_id=excluded.summary_id,symbol=excluded.symbol,chunk_text=excluded.chunk_text,chunk_hash=excluded.chunk_hash,index_status=excluded.index_status,indexed_at=excluded.indexed_at,metadata_json=excluded.metadata_json,created_at=excluded.created_at`,
		chunk.ChunkID, chunk.SummaryID, chunk.Symbol, chunk.ChunkText, chunk.ChunkHash, "indexed", indexedAt, chunk.MetadataJSON, chunk.CreatedAt)
	if err != nil {
		return apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "sqlite-vec chunk metadata 写入失败", err)
	}
	var rowID int64
	if err := tx.QueryRowContext(ctx, `SELECT id FROM vector_chunks WHERE chunk_id=?`, chunk.ChunkID).Scan(&rowID); err != nil {
		return apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "sqlite-vec chunk rowid 查询失败", err)
	}
	serialized, err := sqlitevec.SerializeFloat32(vec)
	if err != nil {
		return apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "sqlite-vec embedding 序列化失败", err)
	}
	_, _ = tx.ExecContext(ctx, `DELETE FROM vector_embeddings WHERE rowid=?`, rowID)
	if _, err := tx.ExecContext(ctx, `INSERT INTO vector_embeddings(rowid, embedding) VALUES (?, ?)`, rowID, serialized); err != nil {
		return apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "sqlite-vec embedding 写入失败", err)
	}
	if err := tx.Commit(); err != nil {
		return apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "sqlite-vec 事务提交失败", err)
	}
	rollback = false
	return nil
}

func (s *SQLiteVecVectorIndex) Search(ctx context.Context, symbol string) ([]repository.RAGChunk, error) {
	db, err := s.open(ctx)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer db.Close()
	rows, err := db.QueryContext(ctx, `SELECT chunk_id,summary_id,symbol,chunk_text,chunk_hash,index_status,indexed_at,metadata_json,created_at FROM vector_chunks WHERE (?='' OR symbol='' OR symbol=?) ORDER BY id ASC`, symbol, symbol)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "sqlite-vec metadata 查询失败", err)
	}
	defer rows.Close()
	return scanSQLiteVecChunks(rows)
}

func (s *SQLiteVecVectorIndex) SearchSimilar(ctx context.Context, query VectorSearchQuery) ([]repository.RAGChunk, error) {
	if s == nil || s.embedder == nil {
		return nil, apperr.New(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "embedding provider 未配置")
	}
	text := firstNonEmptyString(query.Text, query.Symbol)
	vec, err := s.embedder.Embed(ctx, text)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "query embedding 生成失败", err)
	}
	if err := s.validateVector(vec); err != nil {
		return nil, err
	}
	db, err := s.open(ctx)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	serialized, err := sqlitevec.SerializeFloat32(vec)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "query embedding 序列化失败", err)
	}
	topK := query.TopK
	if topK <= 0 {
		topK = s.topK
	}
	rows, err := db.QueryContext(ctx, `SELECT c.chunk_id,c.summary_id,c.symbol,c.chunk_text,c.chunk_hash,c.index_status,c.indexed_at,c.metadata_json,c.created_at
FROM vector_embeddings v
JOIN vector_chunks c ON c.id=v.rowid
WHERE v.embedding MATCH ? AND k=? AND (?='' OR c.symbol='' OR c.symbol=?)
ORDER BY distance
LIMIT ?`, serialized, topK, query.Symbol, query.Symbol, topK)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "sqlite-vec topK 查询失败", err)
	}
	defer rows.Close()
	return scanSQLiteVecChunks(rows)
}

func (s *SQLiteVecVectorIndex) Health(ctx context.Context) VectorIndexHealth {
	health := VectorIndexHealth{Status: VectorIndexHealthDegraded, Path: "", Rebuildable: true}
	if s == nil || strings.TrimSpace(s.path) == "" {
		health.DegradedReason = "sqlite-vec 索引路径未配置"
		return health
	}
	health.Path = s.path
	if _, err := os.Stat(s.path); os.IsNotExist(err) {
		health.Status = VectorIndexHealthMissing
		return health
	}
	db, err := s.open(ctx)
	if err != nil {
		health.Status = VectorIndexHealthCorrupted
		health.DegradedReason = err.Error()
		return health
	}
	defer db.Close()
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM vector_chunks`).Scan(&health.ChunkCount); err != nil {
		health.Status = VectorIndexHealthCorrupted
		health.DegradedReason = err.Error()
		return health
	}
	health.Status = VectorIndexHealthHealthy
	health.Version = CurrentVectorIndexVersion
	return health
}

func (s *SQLiteVecVectorIndex) open(ctx context.Context) (*sql.DB, error) {
	if s == nil || strings.TrimSpace(s.path) == "" {
		return nil, apperr.New(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "sqlite-vec 索引路径未配置")
	}
	if s.dimensions <= 0 {
		return nil, apperr.New(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "embedding dimensions 未配置")
	}
	sqliteVecAutoOnce.Do(func() { sqlitevec.Auto() })
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return nil, apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "sqlite-vec 索引目录不可写", err)
	}
	db, err := sql.Open("sqlite3", s.path)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "sqlite-vec 打开失败", err)
	}
	if err := s.ensureSchema(ctx, db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func (s *SQLiteVecVectorIndex) ensureSchema(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS vector_chunks (
id INTEGER PRIMARY KEY AUTOINCREMENT,
chunk_id TEXT NOT NULL UNIQUE,
summary_id TEXT NOT NULL,
symbol TEXT,
chunk_text TEXT NOT NULL,
chunk_hash TEXT,
index_status TEXT NOT NULL,
indexed_at TEXT,
metadata_json TEXT,
created_at TEXT
)`); err != nil {
		return apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "sqlite-vec metadata schema 初始化失败", err)
	}
	stmt := fmt.Sprintf(`CREATE VIRTUAL TABLE IF NOT EXISTS vector_embeddings USING vec0(embedding float[%d])`, s.dimensions)
	if _, err := db.ExecContext(ctx, stmt); err != nil {
		return apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "sqlite-vec vector schema 初始化失败", err)
	}
	return nil
}

func (s *SQLiteVecVectorIndex) validateVector(vec []float32) error {
	if len(vec) == 0 {
		return apperr.New(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "embedding 为空")
	}
	if s.dimensions > 0 && len(vec) != s.dimensions {
		return apperr.New(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, fmt.Sprintf("embedding dimensions mismatch: got %d want %d", len(vec), s.dimensions))
	}
	return nil
}

func scanSQLiteVecChunks(rows *sql.Rows) ([]repository.RAGChunk, error) {
	var chunks []repository.RAGChunk
	for rows.Next() {
		var chunk repository.RAGChunk
		if err := rows.Scan(&chunk.ChunkID, &chunk.SummaryID, &chunk.Symbol, &chunk.ChunkText, &chunk.ChunkHash, &chunk.IndexStatus, &chunk.IndexedAt, &chunk.MetadataJSON, &chunk.CreatedAt); err != nil {
			return nil, apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "sqlite-vec chunk scan 失败", err)
		}
		chunks = append(chunks, chunk)
	}
	if err := rows.Err(); err != nil {
		return nil, apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "sqlite-vec chunk rows 失败", err)
	}
	return chunks, nil
}
