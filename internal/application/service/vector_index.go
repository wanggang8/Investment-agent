package service

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

// VectorIndex 定义 VecLite 索引写入和查询边界，便于本地 stub 和真实实现替换。
type VectorIndex interface {
	Upsert(ctx context.Context, chunk repository.RAGChunk) error
	Search(ctx context.Context, symbol string) ([]repository.RAGChunk, error)
}

const CurrentVectorIndexVersion = 1

const (
	VectorIndexHealthHealthy      = "healthy"
	VectorIndexHealthMissing      = "missing"
	VectorIndexHealthCorrupted    = "corrupted"
	VectorIndexHealthIncompatible = "incompatible"
	VectorIndexHealthDegraded     = "degraded"
)

type VectorIndexHealth struct {
	Status         string
	Path           string
	Version        int
	ChunkCount     int
	Rebuildable    bool
	DegradedReason string
}

type VectorIndexRebuildStats struct {
	Status         string
	IndexedCount   int
	SkippedCount   int
	LastRebuildAt  string
	DegradedReason string
}

type fileVectorIndexEnvelope struct {
	Version int                   `json:"version"`
	Chunks  []repository.RAGChunk `json:"chunks"`
}

// MemoryVectorIndex 是本地开发用索引实现，不依赖外部服务。
type MemoryVectorIndex struct {
	Path   string
	Chunks []repository.RAGChunk
}

func (m *MemoryVectorIndex) Upsert(_ context.Context, chunk repository.RAGChunk) error {
	m.Chunks = append(m.Chunks, chunk)
	return nil
}

func (m *MemoryVectorIndex) Search(_ context.Context, symbol string) ([]repository.RAGChunk, error) {
	if symbol == "" {
		return m.Chunks, nil
	}
	chunks := make([]repository.RAGChunk, 0, len(m.Chunks))
	for _, chunk := range m.Chunks {
		if chunk.Symbol == "" || chunk.Symbol == symbol {
			chunks = append(chunks, chunk)
		}
	}
	return chunks, nil
}

// FileVectorIndex 是本地 VecLite 文件索引适配器。
// 当前用 JSON 保存可重建 chunk 元数据，确保索引路径、持久化、损坏降级和 SQLite 重建流程可验证。
type FileVectorIndex struct {
	Path string
}

func NewFileVectorIndex(path string) *FileVectorIndex {
	return &FileVectorIndex{Path: path}
}

func (f *FileVectorIndex) Upsert(_ context.Context, chunk repository.RAGChunk) error {
	if f == nil || f.Path == "" {
		return apperr.New(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "VecLite 索引路径未配置")
	}
	chunks, err := f.readChunks()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	updated := false
	for i, existing := range chunks {
		if existing.ChunkID == chunk.ChunkID {
			chunks[i] = chunk
			updated = true
			break
		}
	}
	if !updated {
		chunks = append(chunks, chunk)
	}
	return f.writeChunks(chunks)
}

func (f *FileVectorIndex) Search(_ context.Context, symbol string) ([]repository.RAGChunk, error) {
	if f == nil || f.Path == "" {
		return nil, apperr.New(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "VecLite 索引路径未配置")
	}
	chunks, err := f.readChunks()
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if symbol == "" {
		return chunks, nil
	}
	out := make([]repository.RAGChunk, 0, len(chunks))
	for _, chunk := range chunks {
		if chunk.Symbol == "" || chunk.Symbol == symbol {
			out = append(out, chunk)
		}
	}
	return out, nil
}

func (f *FileVectorIndex) Health(_ context.Context) VectorIndexHealth {
	health := VectorIndexHealth{Status: VectorIndexHealthDegraded, Path: "", Rebuildable: true}
	if f == nil || f.Path == "" {
		health.DegradedReason = "VecLite 索引路径未配置"
		return health
	}
	health.Path = f.Path
	data, err := os.ReadFile(f.Path)
	if errors.Is(err, os.ErrNotExist) {
		health.Status = VectorIndexHealthMissing
		return health
	}
	if err != nil {
		health.DegradedReason = err.Error()
		return health
	}
	envelope, err := decodeFileVectorIndexEnvelope(data)
	if err != nil {
		health.Status = VectorIndexHealthCorrupted
		health.DegradedReason = err.Error()
		return health
	}
	health.Version = envelope.Version
	if envelope.Version != CurrentVectorIndexVersion {
		health.Status = VectorIndexHealthIncompatible
		health.DegradedReason = "VecLite 索引版本不兼容"
		return health
	}
	health.Status = VectorIndexHealthHealthy
	health.ChunkCount = len(envelope.Chunks)
	return health
}

func (f *FileVectorIndex) readChunks() ([]repository.RAGChunk, error) {
	data, err := os.ReadFile(f.Path)
	if err != nil {
		return nil, err
	}
	envelope, err := decodeFileVectorIndexEnvelope(data)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "VecLite 索引文件不可解析", err)
	}
	if envelope.Version != CurrentVectorIndexVersion {
		return nil, apperr.New(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "VecLite 索引版本不兼容")
	}
	return envelope.Chunks, nil
}

func (f *FileVectorIndex) writeChunks(chunks []repository.RAGChunk) error {
	if err := os.MkdirAll(filepath.Dir(f.Path), 0o700); err != nil {
		return apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "VecLite 索引目录不可写", err)
	}
	data, err := json.Marshal(fileVectorIndexEnvelope{Version: CurrentVectorIndexVersion, Chunks: chunks})
	if err != nil {
		return apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "VecLite 索引序列化失败", err)
	}
	if err := os.WriteFile(f.Path, data, 0o600); err != nil {
		return apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "VecLite 索引写入失败", err)
	}
	return nil
}

func decodeFileVectorIndexEnvelope(data []byte) (fileVectorIndexEnvelope, error) {
	var envelope fileVectorIndexEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return fileVectorIndexEnvelope{}, err
	}
	return envelope, nil
}

// RebuildVectorIndex 从 SQLite 文本块重建 VecLite 索引；索引失败时返回可识别降级错误。
func (s *EvidenceService) RebuildVectorIndex(ctx context.Context, index VectorIndex) (int, error) {
	stats, err := s.RebuildVectorIndexWithStats(ctx, index)
	return stats.IndexedCount, err
}

func (s *EvidenceService) RebuildVectorIndexWithStats(ctx context.Context, index VectorIndex) (VectorIndexRebuildStats, error) {
	stats := VectorIndexRebuildStats{Status: VectorIndexHealthDegraded}
	if index == nil {
		stats.DegradedReason = "VecLite 索引不可用"
		return stats, apperr.New(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, stats.DegradedReason)
	}
	var chunks []repository.RAGChunk
	if err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		got, err := repos.IntelligenceRepo.ListRAGChunks(ctx)
		if err != nil {
			return err
		}
		chunks = got
		return nil
	}); err != nil {
		stats.DegradedReason = err.Error()
		return stats, err
	}
	validChunks := make([]repository.RAGChunk, 0, len(chunks))
	validChunkIDs := make([]string, 0, len(chunks))
	rebuildAt := time.Now().UTC().Format(time.RFC3339)
	for _, chunk := range chunks {
		if chunk.ChunkID == "" {
			stats.SkippedCount++
			continue
		}
		chunk.IndexStatus = "indexed"
		chunk.IndexedAt = rebuildAt
		validChunks = append(validChunks, chunk)
		validChunkIDs = append(validChunkIDs, chunk.ChunkID)
	}
	if fileIndex, ok := index.(*FileVectorIndex); ok {
		if err := fileIndex.writeChunks(validChunks); err != nil {
			stats.DegradedReason = err.Error()
			return stats, apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "VecLite 索引写入失败", err)
		}
		if err := s.markRAGChunksIndexed(ctx, validChunkIDs); err != nil {
			stats.DegradedReason = err.Error()
			return stats, err
		}
		stats.IndexedCount = len(validChunks)
		stats.Status = VectorIndexHealthHealthy
		stats.LastRebuildAt = rebuildAt
		return stats, nil
	}
	for _, chunk := range validChunks {
		if err := index.Upsert(ctx, chunk); err != nil {
			stats.DegradedReason = err.Error()
			return stats, apperr.Wrap(apperr.CodeVectorIndexUnavailable, apperr.CategoryInternal, "VecLite 索引写入失败", err)
		}
		stats.IndexedCount++
	}
	if err := s.markRAGChunksIndexed(ctx, validChunkIDs); err != nil {
		stats.DegradedReason = err.Error()
		return stats, err
	}
	stats.Status = VectorIndexHealthHealthy
	stats.LastRebuildAt = rebuildAt
	return stats, nil
}

func (s *EvidenceService) markRAGChunksIndexed(ctx context.Context, chunkIDs []string) error {
	if len(chunkIDs) == 0 {
		return nil
	}
	return s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		return repos.IntelligenceRepo.UpdateRAGChunksIndexStatus(ctx, chunkIDs, "indexed")
	})
}
