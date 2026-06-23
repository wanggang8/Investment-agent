package service

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"investment-agent/internal/application/workflow"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

type intelligenceRepoForVectorTest struct {
	chunks    []repository.RAGChunk
	summaries []repository.IntelligenceSummary
}

func (r intelligenceRepoForVectorTest) SaveIntelligenceItem(context.Context, repository.IntelligenceItem) error {
	return nil
}
func (r intelligenceRepoForVectorTest) GetIntelligenceItem(context.Context, string) (repository.IntelligenceItem, error) {
	return repository.IntelligenceItem{}, nil
}
func (r intelligenceRepoForVectorTest) SaveIntelligenceSummary(context.Context, repository.IntelligenceSummary, []repository.RAGChunk) error {
	return nil
}
func (r intelligenceRepoForVectorTest) GetIntelligenceSummary(context.Context, string) (repository.IntelligenceSummary, []repository.RAGChunk, error) {
	return repository.IntelligenceSummary{}, nil, nil
}
func (r intelligenceRepoForVectorTest) ListEvidenceSummaries(context.Context) ([]repository.IntelligenceSummary, error) {
	return r.summaries, nil
}
func (r intelligenceRepoForVectorTest) SaveSourceVerification(context.Context, repository.SourceVerification) error {
	return nil
}
func (r intelligenceRepoForVectorTest) GetSourceVerification(context.Context, string) (repository.SourceVerification, error) {
	return repository.SourceVerification{}, nil
}
func (r intelligenceRepoForVectorTest) GetLatestSourceVerification(context.Context) (repository.SourceVerification, error) {
	return repository.SourceVerification{}, nil
}
func (r intelligenceRepoForVectorTest) GetLatestSourceVerificationByFilter(context.Context, string, string) (repository.SourceVerification, error) {
	return repository.SourceVerification{}, nil
}
func (r intelligenceRepoForVectorTest) UpdateRAGChunksIndexStatus(context.Context, []string, string) error {
	return nil
}
func (r intelligenceRepoForVectorTest) ListRAGChunks(context.Context) ([]repository.RAGChunk, error) {
	return r.chunks, nil
}
func (r intelligenceRepoForVectorTest) CountRAGChunks(context.Context) (int, error) {
	return len(r.chunks), nil
}

type recordingIntelligenceRepoForVectorTest struct {
	intelligenceRepoForVectorTest
	updatedIDs    []string
	updatedStatus string
}

func (r *recordingIntelligenceRepoForVectorTest) UpdateRAGChunksIndexStatus(_ context.Context, chunkIDs []string, status string) error {
	r.updatedIDs = append([]string{}, chunkIDs...)
	r.updatedStatus = status
	return nil
}

func TestFileVectorIndexPersistsChunksAcrossInstances(t *testing.T) {
	path := filepath.Join(t.TempDir(), "investment.vec.json")
	index := NewFileVectorIndex(path)
	chunk := repository.RAGChunk{ChunkID: "chunk1", SummaryID: "sum1", Symbol: "510300", ChunkText: "沪深300 证据", IndexStatus: "indexed"}

	if err := index.Upsert(context.Background(), chunk); err != nil {
		t.Fatalf("Upsert: %v", err)
	}

	reopened := NewFileVectorIndex(path)
	got, err := reopened.Search(context.Background(), "510300")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(got) != 1 || got[0].ChunkID != "chunk1" || got[0].SummaryID != "sum1" {
		t.Fatalf("expected persisted chunk, got %+v", got)
	}
}

func TestFileVectorIndexFiltersBySymbol(t *testing.T) {
	path := filepath.Join(t.TempDir(), "investment.vec.json")
	index := NewFileVectorIndex(path)
	if err := index.Upsert(context.Background(), repository.RAGChunk{ChunkID: "other", SummaryID: "sum_other", Symbol: "159915", ChunkText: "其他"}); err != nil {
		t.Fatalf("Upsert other: %v", err)
	}
	if err := index.Upsert(context.Background(), repository.RAGChunk{ChunkID: "target", SummaryID: "sum_target", Symbol: "510300", ChunkText: "目标"}); err != nil {
		t.Fatalf("Upsert target: %v", err)
	}

	got, err := index.Search(context.Background(), "510300")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(got) != 1 || got[0].ChunkID != "target" {
		t.Fatalf("expected target chunk only, got %+v", got)
	}
}

func TestFileVectorIndexReportsCorruptedFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "investment.vec.json")
	if err := os.WriteFile(path, []byte("not-json"), 0o600); err != nil {
		t.Fatalf("write corrupted file: %v", err)
	}

	_, err := NewFileVectorIndex(path).Search(context.Background(), "510300")
	if !apperr.IsCode(err, apperr.CodeVectorIndexUnavailable) {
		t.Fatalf("expected VECTOR_INDEX_UNAVAILABLE, got %v", err)
	}
}

func TestFileVectorIndexHealthStates(t *testing.T) {
	missing := NewFileVectorIndex(filepath.Join(t.TempDir(), "missing.json")).Health(context.Background())
	if missing.Status != VectorIndexHealthMissing || !missing.Rebuildable {
		t.Fatalf("expected missing rebuildable health, got %+v", missing)
	}

	corruptedPath := filepath.Join(t.TempDir(), "corrupted.json")
	if err := os.WriteFile(corruptedPath, []byte("not-json"), 0o600); err != nil {
		t.Fatalf("write corrupted file: %v", err)
	}
	corrupted := NewFileVectorIndex(corruptedPath).Health(context.Background())
	if corrupted.Status != VectorIndexHealthCorrupted || corrupted.DegradedReason == "" || !corrupted.Rebuildable {
		t.Fatalf("expected corrupted rebuildable health, got %+v", corrupted)
	}

	incompatiblePath := filepath.Join(t.TempDir(), "incompatible.json")
	if err := os.WriteFile(incompatiblePath, []byte(`{"version":999,"chunks":[]}`), 0o600); err != nil {
		t.Fatalf("write incompatible file: %v", err)
	}
	incompatible := NewFileVectorIndex(incompatiblePath).Health(context.Background())
	if incompatible.Status != VectorIndexHealthIncompatible || incompatible.Version != 999 || !incompatible.Rebuildable {
		t.Fatalf("expected incompatible rebuildable health, got %+v", incompatible)
	}

	healthyPath := filepath.Join(t.TempDir(), "healthy.json")
	healthyIndex := NewFileVectorIndex(healthyPath)
	if err := healthyIndex.Upsert(context.Background(), repository.RAGChunk{ChunkID: "chunk1", SummaryID: "sum1", Symbol: "510300", ChunkText: "证据"}); err != nil {
		t.Fatalf("Upsert healthy chunk: %v", err)
	}
	healthy := healthyIndex.Health(context.Background())
	if healthy.Status != VectorIndexHealthHealthy || healthy.ChunkCount != 1 || healthy.Version != CurrentVectorIndexVersion || !healthy.Rebuildable {
		t.Fatalf("expected healthy index, got %+v", healthy)
	}
}

func TestFileVectorIndexUsesVersionedEnvelope(t *testing.T) {
	path := filepath.Join(t.TempDir(), "investment.vec.json")
	index := NewFileVectorIndex(path)
	if err := index.Upsert(context.Background(), repository.RAGChunk{ChunkID: "chunk1", SummaryID: "sum1", ChunkText: "证据"}); err != nil {
		t.Fatalf("Upsert: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read index file: %v", err)
	}
	var envelope struct {
		Version int                   `json:"version"`
		Chunks  []repository.RAGChunk `json:"chunks"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		t.Fatalf("decode index envelope: %v", err)
	}
	if envelope.Version != CurrentVectorIndexVersion || len(envelope.Chunks) != 1 {
		t.Fatalf("expected versioned envelope, got %+v", envelope)
	}
}

func TestRebuildVectorIndexUsesSQLiteChunks(t *testing.T) {
	chunks := []repository.RAGChunk{{ChunkID: "chunk1", ChunkText: "证据文本", IndexStatus: "pending"}}
	svc := NewEvidenceService(transactorStub{repos: repository.Repositories{IntelligenceRepo: intelligenceRepoForVectorTest{chunks: chunks}}})
	index := &MemoryVectorIndex{}

	count, err := svc.RebuildVectorIndex(context.Background(), index)
	if err != nil {
		t.Fatalf("RebuildVectorIndex: %v", err)
	}
	if count != 1 || len(index.Chunks) != 1 || index.Chunks[0].ChunkID != "chunk1" {
		t.Fatalf("unexpected rebuild result count=%d chunks=%+v", count, index.Chunks)
	}
}

func TestRebuildVectorIndexReturnsStatistics(t *testing.T) {
	chunks := []repository.RAGChunk{
		{ChunkID: "chunk1", ChunkText: "证据文本", IndexStatus: "pending"},
		{ChunkID: "", ChunkText: "缺少 ID", IndexStatus: "pending"},
	}
	svc := NewEvidenceService(transactorStub{repos: repository.Repositories{IntelligenceRepo: intelligenceRepoForVectorTest{chunks: chunks}}})
	index := &MemoryVectorIndex{}

	stats, err := svc.RebuildVectorIndexWithStats(context.Background(), index)
	if err != nil {
		t.Fatalf("RebuildVectorIndexWithStats: %v", err)
	}
	if stats.IndexedCount != 1 || stats.SkippedCount != 1 || stats.Status != VectorIndexHealthHealthy || stats.LastRebuildAt == "" || stats.DegradedReason != "" {
		t.Fatalf("unexpected rebuild stats: %+v", stats)
	}
	if len(index.Chunks) != 1 || index.Chunks[0].ChunkID != "chunk1" || index.Chunks[0].IndexedAt == "" {
		t.Fatalf("expected only valid chunk indexed, got %+v", index.Chunks)
	}
}

func TestRebuildVectorIndexPersistsIndexedStatusToSQLite(t *testing.T) {
	repo := &recordingIntelligenceRepoForVectorTest{
		intelligenceRepoForVectorTest: intelligenceRepoForVectorTest{chunks: []repository.RAGChunk{
			{ChunkID: "chunk1", ChunkText: "证据文本", IndexStatus: "pending"},
			{ChunkID: "", ChunkText: "缺少 ID", IndexStatus: "pending"},
		}},
	}
	svc := NewEvidenceService(transactorStub{repos: repository.Repositories{IntelligenceRepo: repo}})
	index := &MemoryVectorIndex{}

	stats, err := svc.RebuildVectorIndexWithStats(context.Background(), index)
	if err != nil {
		t.Fatalf("RebuildVectorIndexWithStats: %v", err)
	}
	if stats.IndexedCount != 1 {
		t.Fatalf("expected one indexed chunk, got %+v", stats)
	}
	if repo.updatedStatus != "indexed" || len(repo.updatedIDs) != 1 || repo.updatedIDs[0] != "chunk1" {
		t.Fatalf("expected SQLite index status update for valid chunk, got ids=%+v status=%s", repo.updatedIDs, repo.updatedStatus)
	}
}

func TestRebuildVectorIndexReturnsDegradedStatisticsOnWriteFailure(t *testing.T) {
	chunks := []repository.RAGChunk{{ChunkID: "chunk1", ChunkText: "证据文本", IndexStatus: "pending"}}
	svc := NewEvidenceService(transactorStub{repos: repository.Repositories{IntelligenceRepo: intelligenceRepoForVectorTest{chunks: chunks}}})

	stats, err := svc.RebuildVectorIndexWithStats(context.Background(), failingUpsertVectorIndex{err: errors.New("disk full")})
	if !apperr.IsCode(err, apperr.CodeVectorIndexUnavailable) {
		t.Fatalf("expected VECTOR_INDEX_UNAVAILABLE, got %v", err)
	}
	if stats.Status != VectorIndexHealthDegraded || stats.DegradedReason == "" || stats.IndexedCount != 0 || stats.SkippedCount != 0 {
		t.Fatalf("expected degraded write stats, got %+v", stats)
	}
}

func TestRebuildVectorIndexUnavailable(t *testing.T) {
	svc := NewEvidenceService(transactorStub{repos: repository.Repositories{IntelligenceRepo: intelligenceRepoForVectorTest{}}})

	_, err := svc.RebuildVectorIndex(context.Background(), nil)
	if !apperr.IsCode(err, apperr.CodeVectorIndexUnavailable) {
		t.Fatalf("expected VECTOR_INDEX_UNAVAILABLE, got %v", err)
	}
}

type failingVectorIndex struct{ err error }

type failingUpsertVectorIndex struct{ err error }

func (f failingUpsertVectorIndex) Upsert(context.Context, repository.RAGChunk) error { return f.err }
func (f failingUpsertVectorIndex) Search(context.Context, string) ([]repository.RAGChunk, error) {
	return nil, nil
}

func (f failingVectorIndex) Upsert(context.Context, repository.RAGChunk) error { return nil }
func (f failingVectorIndex) Search(context.Context, string) ([]repository.RAGChunk, error) {
	return nil, f.err
}

func TestRetrievalAdapterFiltersSQLiteFallbackBySymbol(t *testing.T) {
	repo := intelligenceRepoForVectorTest{summaries: []repository.IntelligenceSummary{
		{SummaryID: "sum_other", Symbol: "159915", SourceLevel: "A", EvidenceRole: "formal", EventType: "normal", Summary: "其他标的摘要"},
		{SummaryID: "sum_target", Symbol: "510300", SourceLevel: "A", EvidenceRole: "formal", EventType: "normal", Summary: "目标标的摘要"},
	}}
	adapter := NewRetrievalAdapter(transactorStub{repos: repository.Repositories{IntelligenceRepo: repo}}, nil)

	out, err := adapter.RetrieveEvidence(context.Background(), workflow.RetrievalRequest{Symbol: "510300"})
	if err != nil {
		t.Fatalf("RetrieveEvidence: %v", err)
	}
	if out.OutputRef != "sum_target" || len(out.EvidenceSet.Items) != 1 || out.EvidenceSet.Items[0].EvidenceID != "sum_target" {
		t.Fatalf("expected target symbol summary only, got %+v", out)
	}
}

func TestRetrievalAdapterUsesSpecificReasonWhenIndexMissing(t *testing.T) {
	repo := intelligenceRepoForVectorTest{summaries: []repository.IntelligenceSummary{{SummaryID: "sum1", SourceLevel: "A", EvidenceRole: "formal", EventType: "normal", Summary: "摘要"}}}
	adapter := NewRetrievalAdapter(transactorStub{repos: repository.Repositories{IntelligenceRepo: repo}}, nil)

	out, err := adapter.RetrieveEvidence(context.Background(), workflow.RetrievalRequest{Symbol: "510300"})
	if err != nil {
		t.Fatalf("RetrieveEvidence: %v", err)
	}
	if out.DegradedReason != "veclite index not configured" {
		t.Fatalf("expected index missing reason, got %+v", out)
	}
}

func TestRetrievalAdapterUsesSpecificReasonWhenSearchFails(t *testing.T) {
	repo := intelligenceRepoForVectorTest{summaries: []repository.IntelligenceSummary{{SummaryID: "sum1", SourceLevel: "A", EvidenceRole: "formal", EventType: "normal", Summary: "摘要"}}}
	adapter := NewRetrievalAdapter(transactorStub{repos: repository.Repositories{IntelligenceRepo: repo}}, failingVectorIndex{err: errors.New("search down")})

	out, err := adapter.RetrieveEvidence(context.Background(), workflow.RetrievalRequest{Symbol: "510300"})
	if err != nil {
		t.Fatalf("RetrieveEvidence: %v", err)
	}
	if out.DegradedReason != "veclite search failed" {
		t.Fatalf("expected search failed reason, got %+v", out)
	}
}

func TestRetrievalAdapterUsesSpecificReasonWhenSearchEmpty(t *testing.T) {
	repo := intelligenceRepoForVectorTest{summaries: []repository.IntelligenceSummary{{SummaryID: "sum1", SourceLevel: "A", EvidenceRole: "formal", EventType: "normal", Summary: "摘要"}}}
	adapter := NewRetrievalAdapter(transactorStub{repos: repository.Repositories{IntelligenceRepo: repo}}, &MemoryVectorIndex{})

	out, err := adapter.RetrieveEvidence(context.Background(), workflow.RetrievalRequest{Symbol: "510300"})
	if err != nil {
		t.Fatalf("RetrieveEvidence: %v", err)
	}
	if out.DegradedReason != "veclite search empty" {
		t.Fatalf("expected search empty reason, got %+v", out)
	}
}

func TestRetrievalAdapterFallsBackToSQLiteSummary(t *testing.T) {
	repo := intelligenceRepoForVectorTest{summaries: []repository.IntelligenceSummary{{SummaryID: "sum1", SourceLevel: "A", EvidenceRole: "formal", EventType: "normal", VerificationStatus: "satisfied", VerificationEvidenceIDsJSON: `["sum1"]`, VerificationEvidenceRole: "formal", VerificationEventType: "normal", VerificationHighestSourceLevel: "A", Summary: "摘要"}}}
	adapter := NewRetrievalAdapter(transactorStub{repos: repository.Repositories{IntelligenceRepo: repo}}, nil)

	out, err := adapter.RetrieveEvidence(context.Background(), workflow.RetrievalRequest{Symbol: "510300"})
	if err != nil {
		t.Fatalf("RetrieveEvidence: %v", err)
	}
	if out.DegradedReason != "veclite index not configured" || out.OutputRef != "sum1" || len(out.EvidenceSet.Items) != 1 {
		t.Fatalf("unexpected fallback result: %+v", out)
	}
	if out.QualitySummary.TopK != 1 || out.QualitySummary.FallbackSource != "sqlite_summary" || out.QualitySummary.IndexHealth != VectorIndexHealthMissing || out.QualitySummary.DegradedReason != "veclite index not configured" {
		t.Fatalf("expected fallback quality summary, got %+v", out.QualitySummary)
	}
}

func TestEvaluateRetrievalQualityReportsHitAndBackgroundViolation(t *testing.T) {
	result := workflow.RetrievalResult{EvidenceSet: model.EvidenceSet{Items: []model.Evidence{
		{EvidenceID: "sum_formal", SourceLevel: model.SourceLevelA, Role: model.EvidenceFormal},
		{EvidenceID: "sum_background", SourceLevel: model.SourceLevelC, Role: model.EvidenceBackground},
	}}, QualitySummary: workflow.RetrievalQualitySummary{TopK: 2, FallbackSource: "veclite"}}
	fixture := RetrievalQualityFixture{Query: "监管处罚", Symbol: "510300", ExpectedEvidenceIDs: []string{"sum_formal"}, FormalOnly: true}

	evaluation := EvaluateRetrievalQuality(fixture, result)

	if evaluation.Status != "miss" || len(evaluation.MissingExpectedEvidenceIDs) != 0 {
		t.Fatalf("expected formal-only background violation to be a miss, got %+v", evaluation)
	}
	if len(evaluation.UnexpectedBackgroundEvidenceIDs) != 1 || evaluation.UnexpectedBackgroundEvidenceIDs[0] != "sum_background" {
		t.Fatalf("expected background-only diagnostic, got %+v", evaluation)
	}
	if evaluation.TopK != 2 || evaluation.FallbackSource != "veclite" {
		t.Fatalf("expected retrieval summary copied into evaluation, got %+v", evaluation)
	}
}

func TestRetrievalAdapterPreservesSummaryVerificationDetails(t *testing.T) {
	repo := intelligenceRepoForVectorTest{summaries: []repository.IntelligenceSummary{{SummaryID: "sum_verified", Symbol: "510300", SourceLevel: "A", EvidenceRole: "formal", EventType: "major_negative", VerificationStatus: "satisfied", VerificationEvidenceIDsJSON: `["sum_verified"]`, VerificationEvidenceRole: "formal", VerificationEventType: "major_negative", VerificationHighestSourceLevel: "A", IndependentSourceCount: 3, HighGradeIndependentSourceCount: 2, SourceName: "交易所公告", OriginalURL: "https://example.com/a", PublishedAt: "2026-05-29T01:00:00Z", CapturedAt: "2026-05-29T02:00:00Z", ContentHash: "hash_a", Summary: "多源确认风险"}}}
	adapter := NewRetrievalAdapter(transactorStub{repos: repository.Repositories{IntelligenceRepo: repo}}, nil)

	out, err := adapter.RetrieveEvidence(context.Background(), workflow.RetrievalRequest{Symbol: "510300"})
	if err != nil {
		t.Fatalf("RetrieveEvidence: %v", err)
	}
	item := out.EvidenceSet.Items[0]
	if item.IndependentSourceCount != 3 || item.HighGradeIndependentSourceCount != 2 || out.EvidenceSet.VerificationStatus != "satisfied" || item.EventType != "major_negative" {
		t.Fatalf("expected preserved verification details, got %+v", out.EvidenceSet)
	}
	if item.SourceName != "交易所公告" || item.OriginalURL != "https://example.com/a" || item.PublishedAt != "2026-05-29T01:00:00Z" || item.CapturedAt != "2026-05-29T02:00:00Z" || item.ContentHash != "hash_a" {
		t.Fatalf("expected preserved source metadata, got %+v", item)
	}
}

func TestRetrievalAdapterVecLiteHitDoesNotPromoteCLevelOrUnverifiedSummary(t *testing.T) {
	repo := intelligenceRepoForVectorTest{summaries: []repository.IntelligenceSummary{{SummaryID: "sum_c", Symbol: "510300", SourceLevel: "C", EvidenceRole: "formal", EventType: "major_negative", VerificationStatus: "background_only", IndependentSourceCount: 1, Summary: "背景材料"}}}
	index := &MemoryVectorIndex{Chunks: []repository.RAGChunk{{ChunkID: "chunk_c", SummaryID: "sum_c", Symbol: "510300", ChunkText: "背景材料", IndexStatus: "indexed"}}}
	adapter := NewRetrievalAdapter(transactorStub{repos: repository.Repositories{IntelligenceRepo: repo}}, index)

	out, err := adapter.RetrieveEvidence(context.Background(), workflow.RetrievalRequest{Symbol: "510300"})
	if err != nil {
		t.Fatalf("RetrieveEvidence: %v", err)
	}
	item := out.EvidenceSet.Items[0]
	if item.SourceLevel != "C" || item.Role != "background" || out.EvidenceSet.VerificationStatus != "background_only" {
		t.Fatalf("VecLite hit must preserve summary restrictions, got %+v", out.EvidenceSet)
	}
}

func TestRetrievalAdapterRanksVerifiedFormalEvidenceBeforeBackground(t *testing.T) {
	repo := intelligenceRepoForVectorTest{summaries: []repository.IntelligenceSummary{
		{SummaryID: "sum_background", Symbol: "510300", SourceLevel: "C", EvidenceRole: "background", EventType: "normal", VerificationStatus: "background_only", Summary: "背景材料", RelevanceScore: 0.99, TimeWeight: 1},
		{SummaryID: "sum_formal", Symbol: "510300", SourceLevel: "A", EvidenceRole: "formal", EventType: "normal", VerificationStatus: "satisfied", VerificationEvidenceIDsJSON: `["sum_formal"]`, VerificationEvidenceRole: "formal", VerificationEventType: "normal", VerificationHighestSourceLevel: "A", IndependentSourceCount: 2, HighGradeIndependentSourceCount: 2, Summary: "正式材料", RelevanceScore: 0.80, TimeWeight: 0.9},
	}}
	index := &MemoryVectorIndex{Chunks: []repository.RAGChunk{
		{ChunkID: "chunk_background", SummaryID: "sum_background", Symbol: "510300", ChunkText: "背景材料", IndexStatus: "indexed"},
		{ChunkID: "chunk_formal", SummaryID: "sum_formal", Symbol: "510300", ChunkText: "正式材料", IndexStatus: "indexed"},
	}}
	adapter := NewRetrievalAdapter(transactorStub{repos: repository.Repositories{IntelligenceRepo: repo}}, index)

	out, err := adapter.RetrieveEvidence(context.Background(), workflow.RetrievalRequest{Symbol: "510300"})
	if err != nil {
		t.Fatalf("RetrieveEvidence: %v", err)
	}
	if got := out.EvidenceSet.Items[0].EvidenceID; got != "sum_formal" {
		t.Fatalf("expected verified formal evidence first, got %s in %+v", got, out.EvidenceSet.Items)
	}
	if out.EvidenceSet.VerificationStatus != "satisfied" {
		t.Fatalf("verified formal evidence should keep satisfied status even with background context, got %+v", out.EvidenceSet)
	}
	if out.QualitySummary.SourceConsistencyStatus != "checked" {
		t.Fatalf("expected checked consistency summary, got %+v", out.QualitySummary)
	}
}

func TestRetrievalAdapterDegradesWhenChunkMetadataConflictsWithSummary(t *testing.T) {
	repo := intelligenceRepoForVectorTest{summaries: []repository.IntelligenceSummary{
		{SummaryID: "sum_formal", Symbol: "510300", SourceLevel: "A", EvidenceRole: "formal", EventType: "normal", VerificationStatus: "satisfied", VerificationEvidenceIDsJSON: `["sum_formal"]`, VerificationEvidenceRole: "formal", VerificationEventType: "normal", VerificationHighestSourceLevel: "A", Summary: "正式材料"},
	}}
	index := &MemoryVectorIndex{Chunks: []repository.RAGChunk{
		{ChunkID: "chunk_conflict", SummaryID: "sum_formal", Symbol: "510300", ChunkText: "正式材料", IndexStatus: "indexed", MetadataJSON: `{"source_level":"C","evidence_role":"background"}`},
	}}
	adapter := NewRetrievalAdapter(transactorStub{repos: repository.Repositories{IntelligenceRepo: repo}}, index)

	out, err := adapter.RetrieveEvidence(context.Background(), workflow.RetrievalRequest{Symbol: "510300"})
	if err != nil {
		t.Fatalf("RetrieveEvidence: %v", err)
	}
	if out.QualitySummary.FallbackSource != "sqlite_summary" || out.QualitySummary.DegradedReason != "veclite metadata inconsistent" || out.QualitySummary.SourceConsistencyStatus != "mismatch" {
		t.Fatalf("expected metadata mismatch degraded fallback, got %+v", out.QualitySummary)
	}
	if out.EvidenceSet.Items[0].Role != model.EvidenceFormal {
		t.Fatalf("expected SQLite summary to remain formal after fallback, got %+v", out.EvidenceSet.Items[0])
	}
}

func TestRetrievalAdapterDoesNotTrustVerificationGroupMissingSummaryID(t *testing.T) {
	for _, evidenceIDs := range []string{`["other_summary"]`, `[]`, ``} {
		t.Run("evidence_ids="+evidenceIDs, func(t *testing.T) {
			repo := intelligenceRepoForVectorTest{summaries: []repository.IntelligenceSummary{{
				SummaryID: "sum_formal", Symbol: "510300", SourceLevel: "A", EvidenceRole: "formal", EventType: "normal", VerificationStatus: "satisfied", VerificationEvidenceIDsJSON: evidenceIDs, VerificationEvidenceRole: "formal", VerificationEventType: "normal", VerificationHighestSourceLevel: "A", Summary: "正式材料",
			}}}
			adapter := NewRetrievalAdapter(transactorStub{repos: repository.Repositories{IntelligenceRepo: repo}}, nil)

			out, err := adapter.RetrieveEvidence(context.Background(), workflow.RetrievalRequest{Symbol: "510300"})
			if err != nil {
				t.Fatalf("RetrieveEvidence: %v", err)
			}
			if out.EvidenceSet.VerificationStatus == model.VerificationSatisfied || out.EvidenceSet.Items[0].Role == model.EvidenceFormal {
				t.Fatalf("verification evidence_ids mismatch must not satisfy formal evidence: %+v", out.EvidenceSet)
			}
		})
	}
}

func TestRetrievalAdapterReportsIndexFreshness(t *testing.T) {
	repo := intelligenceRepoForVectorTest{summaries: []repository.IntelligenceSummary{{SummaryID: "sum_formal", Symbol: "510300", SourceLevel: "A", EvidenceRole: "formal", EventType: "normal", VerificationStatus: "satisfied", VerificationEvidenceIDsJSON: `["sum_formal"]`, VerificationEvidenceRole: "formal", VerificationEventType: "normal", VerificationHighestSourceLevel: "A", Summary: "正式材料"}}}
	index := &MemoryVectorIndex{Chunks: []repository.RAGChunk{{ChunkID: "chunk_formal", SummaryID: "sum_formal", Symbol: "510300", ChunkText: "正式材料", IndexStatus: "indexed", IndexedAt: time.Now().UTC().Format(time.RFC3339)}}}
	adapter := NewRetrievalAdapter(transactorStub{repos: repository.Repositories{IntelligenceRepo: repo}}, index)

	out, err := adapter.RetrieveEvidence(context.Background(), workflow.RetrievalRequest{Symbol: "510300"})
	if err != nil {
		t.Fatalf("RetrieveEvidence: %v", err)
	}
	if out.QualitySummary.IndexHealth != VectorIndexHealthHealthy || out.QualitySummary.IndexFreshness != "fresh" {
		t.Fatalf("expected healthy fresh index summary, got %+v", out.QualitySummary)
	}
}

func TestRetrievalAdapterReportsStaleIndexFreshness(t *testing.T) {
	repo := intelligenceRepoForVectorTest{summaries: []repository.IntelligenceSummary{{SummaryID: "sum_formal", Symbol: "510300", SourceLevel: "A", EvidenceRole: "formal", EventType: "normal", VerificationStatus: "satisfied", VerificationEvidenceIDsJSON: `["sum_formal"]`, VerificationEvidenceRole: "formal", VerificationEventType: "normal", VerificationHighestSourceLevel: "A", Summary: "正式材料"}}}
	indexedAt := time.Now().UTC().AddDate(0, 0, -31).Format(time.RFC3339)
	index := &MemoryVectorIndex{Chunks: []repository.RAGChunk{{ChunkID: "chunk_formal", SummaryID: "sum_formal", Symbol: "510300", ChunkText: "正式材料", IndexStatus: "indexed", IndexedAt: indexedAt}}}
	adapter := NewRetrievalAdapter(transactorStub{repos: repository.Repositories{IntelligenceRepo: repo}}, index)

	out, err := adapter.RetrieveEvidence(context.Background(), workflow.RetrievalRequest{Symbol: "510300"})
	if err != nil {
		t.Fatalf("RetrieveEvidence: %v", err)
	}
	if out.QualitySummary.IndexFreshness != "stale" {
		t.Fatalf("expected stale index freshness, got %+v", out.QualitySummary)
	}
}

func TestRetrievalAdapterReportsUnknownFreshnessWhenIndexedAtMissing(t *testing.T) {
	repo := intelligenceRepoForVectorTest{summaries: []repository.IntelligenceSummary{{SummaryID: "sum_formal", Symbol: "510300", SourceLevel: "A", EvidenceRole: "formal", EventType: "normal", VerificationStatus: "satisfied", VerificationEvidenceIDsJSON: `["sum_formal"]`, VerificationEvidenceRole: "formal", VerificationEventType: "normal", VerificationHighestSourceLevel: "A", Summary: "正式材料"}}}
	index := &MemoryVectorIndex{Chunks: []repository.RAGChunk{{ChunkID: "chunk_formal", SummaryID: "sum_formal", Symbol: "510300", ChunkText: "正式材料", IndexStatus: "indexed"}}}
	adapter := NewRetrievalAdapter(transactorStub{repos: repository.Repositories{IntelligenceRepo: repo}}, index)

	out, err := adapter.RetrieveEvidence(context.Background(), workflow.RetrievalRequest{Symbol: "510300"})
	if err != nil {
		t.Fatalf("RetrieveEvidence: %v", err)
	}
	if out.QualitySummary.IndexFreshness != "unknown" {
		t.Fatalf("expected unknown freshness without indexed_at, got %+v", out.QualitySummary)
	}
}

func TestRetrievalAdapterReportsCorruptedIndexAsUnknownFreshness(t *testing.T) {
	repo := intelligenceRepoForVectorTest{summaries: []repository.IntelligenceSummary{{SummaryID: "sum_formal", Symbol: "510300", SourceLevel: "A", EvidenceRole: "formal", EventType: "normal", VerificationStatus: "satisfied", VerificationEvidenceIDsJSON: `["sum_formal"]`, VerificationEvidenceRole: "formal", VerificationEventType: "normal", VerificationHighestSourceLevel: "A", Summary: "正式材料"}}}
	path := filepath.Join(t.TempDir(), "corrupted.json")
	if err := os.WriteFile(path, []byte("not-json"), 0o600); err != nil {
		t.Fatalf("write corrupted index: %v", err)
	}
	adapter := NewRetrievalAdapter(transactorStub{repos: repository.Repositories{IntelligenceRepo: repo}}, NewFileVectorIndex(path))

	out, err := adapter.RetrieveEvidence(context.Background(), workflow.RetrievalRequest{Symbol: "510300"})
	if err != nil {
		t.Fatalf("RetrieveEvidence: %v", err)
	}
	if out.QualitySummary.IndexHealth != VectorIndexHealthCorrupted || out.QualitySummary.IndexFreshness != "unknown" || out.QualitySummary.FallbackSource != "sqlite_summary" {
		t.Fatalf("expected corrupted unknown freshness fallback summary, got %+v", out.QualitySummary)
	}
}

func TestRetrievalAdapterUnverifiedABSummaryIsNotSatisfied(t *testing.T) {
	repo := intelligenceRepoForVectorTest{summaries: []repository.IntelligenceSummary{{SummaryID: "sum_unverified", Symbol: "510300", SourceLevel: "A", EvidenceRole: "formal", EventType: "major_negative", Summary: "未验证摘要"}}}
	adapter := NewRetrievalAdapter(transactorStub{repos: repository.Repositories{IntelligenceRepo: repo}}, nil)

	out, err := adapter.RetrieveEvidence(context.Background(), workflow.RetrievalRequest{Symbol: "510300"})
	if err != nil {
		t.Fatalf("RetrieveEvidence: %v", err)
	}
	if out.EvidenceSet.VerificationStatus != "failed" || out.EvidenceSet.Items[0].Role != "background" {
		t.Fatalf("unverified A/B summary must not be formal satisfied evidence: %+v", out.EvidenceSet)
	}
}

func TestRetrievalAdapterRestrictsCLevelSummary(t *testing.T) {
	repo := intelligenceRepoForVectorTest{summaries: []repository.IntelligenceSummary{{SummaryID: "sum_c", SourceLevel: "C", EvidenceRole: "formal", EventType: "normal", Summary: "背景"}}}
	adapter := NewRetrievalAdapter(transactorStub{repos: repository.Repositories{IntelligenceRepo: repo}}, nil)

	out, err := adapter.RetrieveEvidence(context.Background(), workflow.RetrievalRequest{Symbol: "510300"})
	if err != nil {
		t.Fatalf("RetrieveEvidence: %v", err)
	}
	if out.EvidenceSet.Items[0].Role != "background" || out.EvidenceSet.VerificationStatus != "background_only" {
		t.Fatalf("C level summary must be background only: %+v", out.EvidenceSet)
	}
}
