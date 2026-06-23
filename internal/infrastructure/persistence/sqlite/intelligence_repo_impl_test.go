package sqlite

import (
	"context"
	"testing"

	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

func TestIntelligenceRepositoryWriteReadAndRollback(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	repo := NewIntelligenceRepository(db)

	summary := repository.IntelligenceSummary{
		SummaryID: "sum1", IntelligenceID: "intel1", Symbol: "AAA", Summary: "summary",
		SourceLevel: "A", EvidenceRole: "formal", CreatedAt: testTime,
	}
	chunks := []repository.RAGChunk{{
		ChunkID: "chunk1", SummaryID: "sum1", ChunkText: "chunk", ChunkHash: "hash",
		IndexStatus: "pending", CreatedAt: testTime,
	}}
	if err := repo.SaveIntelligenceSummary(ctx, summary, chunks); err != nil {
		t.Fatal(err)
	}
	got, gotChunks, err := repo.GetIntelligenceSummary(ctx, "sum1")
	if err != nil {
		t.Fatal(err)
	}
	if got.SummaryID != "sum1" || len(gotChunks) != 1 {
		t.Fatalf("unexpected read: %#v %#v", got, gotChunks)
	}

	bad := summary
	bad.SummaryID = "sum_bad"
	bad.SourceLevel = "C"
	bad.EvidenceRole = "formal"
	if err := repo.SaveIntelligenceSummary(ctx, bad, nil); err == nil {
		t.Fatal("expected rollback error")
	}
	if _, _, err := repo.GetIntelligenceSummary(ctx, "sum_bad"); err == nil {
		t.Fatal("summary persisted after rollback")
	}
}

func TestIntelligenceRepositoryListsSymbolsForFallbackAndRAGChunks(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	repo := NewIntelligenceRepository(db)

	items := []struct {
		summary repository.IntelligenceSummary
		chunk   repository.RAGChunk
	}{
		{summary: repository.IntelligenceSummary{SummaryID: "sum_510300", IntelligenceID: "intel_510300", Symbol: "510300", Summary: "目标摘要", SourceLevel: "A", EvidenceRole: "formal", CreatedAt: testTime}, chunk: repository.RAGChunk{ChunkID: "chunk_510300", SummaryID: "sum_510300", ChunkText: "目标文本", ChunkHash: "hash_510300", IndexStatus: "pending", CreatedAt: testTime}},
		{summary: repository.IntelligenceSummary{SummaryID: "sum_159915", IntelligenceID: "intel_159915", Symbol: "159915", Summary: "其他摘要", SourceLevel: "A", EvidenceRole: "formal", CreatedAt: testTime}, chunk: repository.RAGChunk{ChunkID: "chunk_159915", SummaryID: "sum_159915", ChunkText: "其他文本", ChunkHash: "hash_159915", IndexStatus: "pending", CreatedAt: testTime}},
	}
	for _, item := range items {
		if err := repo.SaveIntelligenceSummary(ctx, item.summary, []repository.RAGChunk{item.chunk}); err != nil {
			t.Fatal(err)
		}
	}

	summaries, err := repo.ListEvidenceSummaries(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 2 || summaries[0].Symbol == "" || summaries[1].Symbol == "" {
		t.Fatalf("expected symbols on evidence summaries: %+v", summaries)
	}
	chunks, err := repo.ListRAGChunks(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(chunks) != 2 || chunks[0].Symbol == "" || chunks[1].Symbol == "" {
		t.Fatalf("expected symbols on rag chunks: %+v", chunks)
	}
	_, gotChunks, err := repo.GetIntelligenceSummary(ctx, "sum_510300")
	if err != nil {
		t.Fatal(err)
	}
	if len(gotChunks) != 1 || gotChunks[0].Symbol != "510300" {
		t.Fatalf("expected summary chunks to include symbol: %+v", gotChunks)
	}
}

func TestIntelligenceRepositoryListEvidenceSummariesPreservesItemMetadata(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	repo := NewIntelligenceRepository(db)

	item := repository.IntelligenceItem{IntelligenceID: "intel_meta", SourceName: "交易所公告", SourceLevel: "A", OriginalURL: "https://example.com/a", PublishedAt: "2026-05-29T01:00:00Z", CapturedAt: "2026-05-29T02:00:00Z", ContentHash: "hash_meta", CreatedAt: testTime}
	summary := repository.IntelligenceSummary{SummaryID: "sum_meta", IntelligenceID: "intel_meta", Symbol: "510300", Summary: "真实摘要", SourceLevel: "A", EvidenceRole: "formal", CreatedAt: testTime}
	if err := repo.SaveIntelligenceItem(ctx, item); err != nil {
		t.Fatal(err)
	}
	if err := repo.SaveIntelligenceSummary(ctx, summary, nil); err != nil {
		t.Fatal(err)
	}

	summaries, err := repo.ListEvidenceSummaries(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 1 {
		t.Fatalf("expected one summary, got %+v", summaries)
	}
	got := summaries[0]
	if got.IntelligenceID != "intel_meta" || got.SourceName != "交易所公告" || got.OriginalURL != "https://example.com/a" || got.PublishedAt != "2026-05-29T01:00:00Z" || got.CapturedAt != "2026-05-29T02:00:00Z" || got.ContentHash != "hash_meta" {
		t.Fatalf("metadata not preserved: %+v", got)
	}
}

func TestIntelligenceRepositoryListEvidenceSummariesUsesLatestVerificationPerGroup(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	repo := NewIntelligenceRepository(db)
	item := repository.IntelligenceItem{IntelligenceID: "intel_group", SourceName: "cninfo", SourceLevel: "A", OriginalURL: "https://example.invalid/a", PublishedAt: testTime, CapturedAt: testTime, ContentHash: "hash_group", CreatedAt: testTime}
	if err := repo.SaveIntelligenceItem(ctx, item); err != nil {
		t.Fatal(err)
	}
	summary := repository.IntelligenceSummary{SummaryID: "sum_group", IntelligenceID: "intel_group", Symbol: "510300", Summary: "公告正文", SourceLevel: "A", EvidenceRole: "formal", VerificationGroupID: "group_same", CreatedAt: testTime}
	if err := repo.SaveIntelligenceSummary(ctx, summary, nil); err != nil {
		t.Fatal(err)
	}
	oldVerification := repository.SourceVerification{VerificationID: "ver_old", VerificationGroupID: "group_same", EventID: "event_same", Symbol: "510300", EventType: "public_disclosure", EvidenceRole: "formal", VerificationStatus: "failed", IndependentSourceCount: 1, HighGradeIndependentSourceCount: 1, HighestSourceLevel: "A", CreatedAt: "2026-06-05T00:00:00Z"}
	newVerification := repository.SourceVerification{VerificationID: "ver_new", VerificationGroupID: "group_same", EventID: "event_same", Symbol: "510300", EventType: "public_disclosure", EvidenceRole: "formal", VerificationStatus: "satisfied", IndependentSourceCount: 2, HighGradeIndependentSourceCount: 2, HighestSourceLevel: "A", CreatedAt: "2026-06-05T01:00:00Z"}
	if err := repo.SaveSourceVerification(ctx, oldVerification); err != nil {
		t.Fatal(err)
	}
	if err := repo.SaveSourceVerification(ctx, newVerification); err != nil {
		t.Fatal(err)
	}
	summaries, err := repo.ListEvidenceSummaries(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 1 {
		t.Fatalf("expected one summary row, got %+v", summaries)
	}
	if summaries[0].VerificationStatus != "satisfied" || summaries[0].IndependentSourceCount != 2 {
		t.Fatalf("expected latest verification metadata, got %+v", summaries[0])
	}
}

func TestSourceVerificationWriteRead(t *testing.T) {
	db := testDB(t)
	repo := NewIntelligenceRepository(db)
	verification := repository.SourceVerification{
		VerificationID: "ver1", VerificationGroupID: "group1", EventID: "event1",
		EvidenceRole: "formal", VerificationStatus: "satisfied", IndependentSourceCount: 2,
		HighGradeIndependentSourceCount: 1, HighestSourceLevel: "A", CreatedAt: testTime,
	}
	if err := repo.SaveSourceVerification(context.Background(), verification); err != nil {
		t.Fatal(err)
	}
	got, err := repo.GetSourceVerification(context.Background(), "ver1")
	if err != nil {
		t.Fatal(err)
	}
	if got.VerificationStatus != "satisfied" || got.HighGradeIndependentSourceCount != 1 {
		t.Fatalf("got %#v", got)
	}
	latest, err := repo.GetLatestSourceVerification(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if latest.VerificationID != "ver1" || latest.HighGradeIndependentSourceCount != 1 {
		t.Fatalf("latest %#v", latest)
	}
}

func TestIntelligenceRepositoryClassifiesErrors(t *testing.T) {
	db := testDB(t)
	repo := NewIntelligenceRepository(db)
	if _, err := repo.GetIntelligenceItem(context.Background(), "missing_item"); !apperr.IsCode(err, apperr.CodeNotFound) {
		t.Fatalf("expected not found intelligence item error, got %v", err)
	}
	bad := repository.SourceVerification{VerificationID: "bad_verification", VerificationGroupID: "group1", EventID: "event1", EvidenceRole: "formal", VerificationStatus: "invalid", IndependentSourceCount: 1, CreatedAt: testTime}
	if err := repo.SaveSourceVerification(context.Background(), bad); !apperr.IsCode(err, apperr.CodeConflict) {
		t.Fatalf("expected conflict source verification error, got %v", err)
	}
}
