package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"testing"
	"time"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/infrastructure/persistence/sqlite"
	"investment-agent/internal/pkg/apperr"
)

const testTime = "2026-06-05T00:00:00Z"

func TestPublicEvidenceIngestionWritesRepositoriesAndDedupes(t *testing.T) {
	ctx := context.Background()
	store := testDB(t)
	intelligenceRepo := sqlite.NewIntelligenceRepository(store.DB)
	auditRepo := sqlite.NewAuditRepository(store.DB)

	capturedAt := time.Date(2026, 6, 5, 13, 0, 0, 0, time.UTC)
	collector := FixturePublicEvidenceCollector{
		Fixtures: map[string][]PublicEvidencePayload{
			"510300": {
				{SourceName: "cninfo", SourceLevel: model.SourceLevelA, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "ann-1", Title: "ETF 公告", Text: "公告正文", URL: "https://example.invalid/a", PublishedAt: "2026-06-05T00:00:00+08:00", CapturedAt: capturedAt},
			},
		},
	}

	service := PublicEvidenceIngestionService{
		Collector:        collector,
		IntelligenceRepo: intelligenceRepo,
		AuditRepo:        auditRepo,
		GenerateAuditID:  testIDGenerator(),
	}

	if err := service.IngestPublicEvidence(ctx, "510300", time.Time{}, time.Time{}); err != nil {
		t.Fatalf("first IngestPublicEvidence: %v", err)
	}
	if err := service.IngestPublicEvidence(ctx, "510300", time.Time{}, time.Time{}); err != nil {
		t.Fatalf("second IngestPublicEvidence: %v", err)
	}

	summaries, err := intelligenceRepo.ListEvidenceSummaries(ctx)
	if err != nil {
		t.Fatalf("ListEvidenceSummaries: %v", err)
	}
	if len(summaries) != 1 {
		t.Fatalf("expected one deduped summary, got %+v", summaries)
	}
	if summaries[0].SourceName != "cninfo" || summaries[0].SourceLevel != string(model.SourceLevelA) || summaries[0].EvidenceRole != "formal" || summaries[0].ContentHash == "" {
		t.Fatalf("expected summary with source metadata, got %+v", summaries[0])
	}

	chunks, err := intelligenceRepo.ListRAGChunks(ctx)
	if err != nil {
		t.Fatalf("ListRAGChunks: %v", err)
	}
	if len(chunks) != 1 || chunks[0].Symbol != "510300" || chunks[0].ChunkText != "公告正文" || chunks[0].ChunkHash == "" {
		t.Fatalf("expected one deduped RAG chunk, got %+v", chunks)
	}

	verification, err := intelligenceRepo.GetLatestSourceVerificationByFilter(ctx, "510300", publicEvidenceEventID(collector.Fixtures["510300"][0]))
	if err != nil {
		t.Fatalf("GetLatestSourceVerificationByFilter: %v", err)
	}
	if verification.VerificationStatus != "failed" || verification.IndependentSourceCount != 1 || verification.HighGradeIndependentSourceCount != 1 || verification.HighestSourceLevel != string(model.SourceLevelA) {
		t.Fatalf("expected single-source verification to remain failed/insufficient, got %+v", verification)
	}

	events, err := auditRepo.ListAuditEvents(ctx)
	if err != nil {
		t.Fatalf("ListAuditEvents: %v", err)
	}
	if len(events) != 2 || events[0].Action != "run_local_task" || events[0].Status != "success" {
		t.Fatalf("expected audit events for both refresh attempts, got %+v", events)
	}
}

func TestPublicEvidenceIngestionRepairsPartialWrites(t *testing.T) {
	ctx := context.Background()
	store := testDB(t)
	intelligenceRepo := sqlite.NewIntelligenceRepository(store.DB)
	auditRepo := sqlite.NewAuditRepository(store.DB)
	capturedAt := time.Date(2026, 6, 5, 13, 0, 0, 0, time.UTC)
	payload := PublicEvidencePayload{SourceName: "cninfo", SourceLevel: model.SourceLevelA, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "ann-1", Title: "ETF 公告", Text: "公告正文", URL: "https://example.invalid/a", PublishedAt: "2026-06-05T00:00:00+08:00", CapturedAt: capturedAt}
	items, err := NormalizePublicEvidenceItems([]PublicEvidencePayload{payload})
	if err != nil {
		t.Fatal(err)
	}
	payload = items[0]
	intelligenceID := deterministicEvidenceID("intel", payload)
	if err := intelligenceRepo.SaveIntelligenceItem(ctx, repository.IntelligenceItem{IntelligenceID: intelligenceID, SourceName: payload.SourceName, SourceLevel: string(payload.SourceLevel), OriginalURL: payload.URL, PublishedAt: payload.PublishedAt, CapturedAt: capturedAt.Format(time.RFC3339), ContentHash: payload.ContentHash, RawTitle: payload.Title, RawTextRef: payload.Text, CreatedAt: testTime}); err != nil {
		t.Fatal(err)
	}

	service := PublicEvidenceIngestionService{Collector: FixturePublicEvidenceCollector{Fixtures: map[string][]PublicEvidencePayload{"510300": {payload}}}, IntelligenceRepo: intelligenceRepo, AuditRepo: auditRepo, GenerateAuditID: testIDGenerator()}
	if err := service.IngestPublicEvidence(ctx, "510300", time.Time{}, time.Time{}); err != nil {
		t.Fatalf("IngestPublicEvidence should repair partial write: %v", err)
	}
	chunks, err := intelligenceRepo.ListRAGChunks(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(chunks) != 1 {
		t.Fatalf("expected repaired RAG chunk, got %+v", chunks)
	}
	if _, err := intelligenceRepo.GetLatestSourceVerificationByFilter(ctx, "510300", publicEvidenceEventID(payload)); err != nil {
		t.Fatalf("expected repaired source verification: %v", err)
	}
}

func TestPublicEvidenceIngestionUsesValidMetadataJSON(t *testing.T) {
	ctx := context.Background()
	store := testDB(t)
	intelligenceRepo := sqlite.NewIntelligenceRepository(store.DB)
	auditRepo := sqlite.NewAuditRepository(store.DB)
	collector := FixturePublicEvidenceCollector{Fixtures: map[string][]PublicEvidencePayload{"510300": {{SourceName: "cninfo", SourceLevel: model.SourceLevelA, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "ann-quote", Title: "ETF 公告", Text: "公告正文", URL: "https://example.invalid/a?x=\"quoted\"", AttachmentURL: "https://example.invalid/a\\file.pdf", PublishedAt: "2026-06-05T00:00:00+08:00", CapturedAt: time.Date(2026, 6, 5, 13, 0, 0, 0, time.UTC), Raw: map[string]any{"announcementId": "ann-quote"}}}}}
	service := PublicEvidenceIngestionService{Collector: collector, IntelligenceRepo: intelligenceRepo, AuditRepo: auditRepo, GenerateAuditID: testIDGenerator()}
	if err := service.IngestPublicEvidence(ctx, "510300", time.Time{}, time.Time{}); err != nil {
		t.Fatalf("IngestPublicEvidence: %v", err)
	}
	chunks, err := intelligenceRepo.ListRAGChunks(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var metadata map[string]any
	if err := json.Unmarshal([]byte(chunks[0].MetadataJSON), &metadata); err != nil {
		t.Fatalf("metadata_json should be valid JSON: %s", chunks[0].MetadataJSON)
	}
	if metadata["attachment_url"] == "" || metadata["raw"] == nil {
		t.Fatalf("expected metadata to preserve attachment and raw fields: %+v", metadata)
	}
}

func TestPublicEvidenceIngestionAggregatesIndependentSources(t *testing.T) {
	ctx := context.Background()
	store := testDB(t)
	intelligenceRepo := sqlite.NewIntelligenceRepository(store.DB)
	auditRepo := sqlite.NewAuditRepository(store.DB)
	capturedAt := time.Date(2026, 6, 5, 13, 0, 0, 0, time.UTC)
	collector := FixturePublicEvidenceCollector{Fixtures: map[string][]PublicEvidencePayload{"510300": {
		{SourceName: "cninfo", SourceLevel: model.SourceLevelA, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "cninfo-1", Title: "ETF 公告", Text: "公告正文 A", URL: "https://example.invalid/a", PublishedAt: "2026-06-05T00:00:00+08:00", CapturedAt: capturedAt},
		{SourceName: "szse", SourceLevel: model.SourceLevelA, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "szse-1", Title: "ETF 公告", Text: "公告正文 B", URL: "https://example.invalid/b", PublishedAt: "2026-06-05T00:00:00+08:00", CapturedAt: capturedAt},
	}}}
	service := PublicEvidenceIngestionService{Collector: collector, IntelligenceRepo: intelligenceRepo, AuditRepo: auditRepo, GenerateAuditID: testIDGenerator()}
	if err := service.IngestPublicEvidence(ctx, "510300", time.Time{}, time.Time{}); err != nil {
		t.Fatalf("IngestPublicEvidence: %v", err)
	}
	verification, err := intelligenceRepo.GetLatestSourceVerificationByFilter(ctx, "510300", publicEvidenceEventID(collector.Fixtures["510300"][0]))
	if err != nil {
		t.Fatal(err)
	}
	if verification.VerificationStatus != "satisfied" || verification.IndependentSourceCount != 2 || verification.HighGradeIndependentSourceCount != 2 {
		t.Fatalf("expected two A-level sources to satisfy verification, got %+v", verification)
	}
}

func TestPublicEvidenceIngestionSatisfiesNormalFormalEvidenceWithOneHighGradeAndTwoFormalSources(t *testing.T) {
	ctx := context.Background()
	store := testDB(t)
	intelligenceRepo := sqlite.NewIntelligenceRepository(store.DB)
	auditRepo := sqlite.NewAuditRepository(store.DB)
	capturedAt := time.Date(2026, 6, 18, 10, 0, 0, 0, time.UTC)
	collector := FixturePublicEvidenceCollector{Fixtures: map[string][]PublicEvidencePayload{"510300": {
		{SourceName: "csindex_index", SourceLevel: model.SourceLevelA, SourceType: "fund_profile", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "csindex-000300", Title: "510300 沪深300ETF 产品与跟踪指数事实", Text: "中证指数官方沪深300指数事实", URL: "https://www.csindex.com.cn/000300", PublishedAt: capturedAt.Format(time.RFC3339), CapturedAt: capturedAt},
		{SourceName: "eastmoney_fund", SourceLevel: model.SourceLevelB, SourceType: "fund_profile", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "eastmoney-510300", Title: "510300 沪深300ETF 产品与跟踪指数事实", Text: "东方财富基金公开净值与基金事实", URL: "https://fund.eastmoney.com/510300.html", PublishedAt: capturedAt.Format(time.RFC3339), CapturedAt: capturedAt},
	}}}
	service := PublicEvidenceIngestionService{Collector: collector, IntelligenceRepo: intelligenceRepo, AuditRepo: auditRepo, GenerateAuditID: testIDGenerator()}

	if err := service.IngestPublicEvidence(ctx, "510300", time.Time{}, time.Time{}); err != nil {
		t.Fatalf("IngestPublicEvidence: %v", err)
	}

	verification, err := intelligenceRepo.GetLatestSourceVerificationByFilter(ctx, "510300", "")
	if err != nil {
		t.Fatalf("GetLatestSourceVerificationByFilter: %v", err)
	}
	if verification.VerificationStatus != "satisfied" || verification.IndependentSourceCount != 2 || verification.HighGradeIndependentSourceCount != 1 {
		t.Fatalf("expected normal formal evidence to be satisfied with one A source and two formal independent sources, got %+v", verification)
	}
}

func TestPublicEvidenceIngestionMajorEventsRequireTwoHighGradeIndependentSources(t *testing.T) {
	ctx := context.Background()
	store := testDB(t)
	intelligenceRepo := sqlite.NewIntelligenceRepository(store.DB)
	auditRepo := sqlite.NewAuditRepository(store.DB)
	capturedAt := time.Date(2026, 6, 18, 10, 0, 0, 0, time.UTC)
	collector := FixturePublicEvidenceCollector{Fixtures: map[string][]PublicEvidencePayload{"510300": {
		{SourceName: "cninfo", SourceLevel: model.SourceLevelA, SourceType: string(model.EventMajorNegative), EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "cninfo-major", Title: "重大利空", Text: "重大利空正式公告", URL: "https://example.invalid/a", PublishedAt: capturedAt.Format(time.RFC3339), CapturedAt: capturedAt},
		{SourceName: "financial_news", SourceLevel: model.SourceLevelB, SourceType: string(model.EventMajorNegative), EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "news-major", Title: "重大利空", Text: "重大利空新闻跟进", URL: "https://example.invalid/b", PublishedAt: capturedAt.Format(time.RFC3339), CapturedAt: capturedAt},
	}}}
	service := PublicEvidenceIngestionService{Collector: collector, IntelligenceRepo: intelligenceRepo, AuditRepo: auditRepo, GenerateAuditID: testIDGenerator()}

	if err := service.IngestPublicEvidence(ctx, "510300", time.Time{}, time.Time{}); err != nil {
		t.Fatalf("IngestPublicEvidence: %v", err)
	}
	verification, err := intelligenceRepo.GetLatestSourceVerificationByFilter(ctx, "510300", publicEvidenceEventID(collector.Fixtures["510300"][0]))
	if err != nil {
		t.Fatalf("GetLatestSourceVerificationByFilter: %v", err)
	}
	if verification.VerificationStatus != "failed" || verification.IndependentSourceCount != 2 || verification.HighGradeIndependentSourceCount != 1 {
		t.Fatalf("major event must require two high-grade independent sources, got %+v", verification)
	}
}

func TestPublicEvidenceIngestionAppliesF4TimeDecayAndBackgroundBoundary(t *testing.T) {
	ctx := context.Background()
	store := testDB(t)
	intelligenceRepo := sqlite.NewIntelligenceRepository(store.DB)
	auditRepo := sqlite.NewAuditRepository(store.DB)
	now := time.Now().UTC()
	payloads := []PublicEvidencePayload{
		{SourceName: "cninfo", SourceLevel: model.SourceLevelA, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "f4-12h", Title: "F4 12h", Text: "F4 12h 正文", URL: "https://example.invalid/12h", PublishedAt: now.Add(-12 * time.Hour).Format(time.RFC3339), CapturedAt: now},
		{SourceName: "szse", SourceLevel: model.SourceLevelA, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "f4-3d", Title: "F4 3d", Text: "F4 3d 正文", URL: "https://example.invalid/3d", PublishedAt: now.Add(-72 * time.Hour).Format(time.RFC3339), CapturedAt: now},
		{SourceName: "cninfo", SourceLevel: model.SourceLevelA, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "f4-14d", Title: "F4 14d", Text: "F4 14d 正文", URL: "https://example.invalid/14d", PublishedAt: now.AddDate(0, 0, -14).Format(time.RFC3339), CapturedAt: now},
		{SourceName: "szse", SourceLevel: model.SourceLevelA, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "f4-45d", Title: "F4 45d", Text: "F4 45d 正文", URL: "https://example.invalid/45d", PublishedAt: now.AddDate(0, 0, -45).Format(time.RFC3339), CapturedAt: now},
	}
	collector := FixturePublicEvidenceCollector{Fixtures: map[string][]PublicEvidencePayload{"510300": payloads}}
	service := PublicEvidenceIngestionService{Collector: collector, IntelligenceRepo: intelligenceRepo, AuditRepo: auditRepo, GenerateAuditID: testIDGenerator()}

	if err := service.IngestPublicEvidence(ctx, "510300", time.Time{}, time.Time{}); err != nil {
		t.Fatalf("IngestPublicEvidence: %v", err)
	}

	summaries, err := intelligenceRepo.ListEvidenceSummaries(ctx)
	if err != nil {
		t.Fatalf("ListEvidenceSummaries: %v", err)
	}
	got := map[string]repository.IntelligenceSummary{}
	for _, summary := range summaries {
		got[summary.Summary] = summary
	}
	expectedWeights := map[string]float64{
		"F4 12h 正文": 1.0,
		"F4 3d 正文":  0.8,
		"F4 14d 正文": 0.5,
		"F4 45d 正文": 0.2,
	}
	for summaryText, want := range expectedWeights {
		summary, ok := got[summaryText]
		if !ok {
			t.Fatalf("expected summary %q in %+v", summaryText, summaries)
		}
		if math.Abs(summary.TimeWeight-want) > 0.0001 {
			t.Fatalf("expected %s time_weight %.1f, got %.3f", summaryText, want, summary.TimeWeight)
		}
	}
	if got["F4 45d 正文"].EvidenceRole != string(model.EvidenceBackground) {
		t.Fatalf(">30 day evidence must be background only, got %+v", got["F4 45d 正文"])
	}
}

func TestPublicEvidenceIngestionKeepsHighGradeCountDistinctBySource(t *testing.T) {
	ctx := context.Background()
	store := testDB(t)
	intelligenceRepo := sqlite.NewIntelligenceRepository(store.DB)
	auditRepo := sqlite.NewAuditRepository(store.DB)
	capturedAt := time.Date(2026, 6, 5, 13, 0, 0, 0, time.UTC)
	collector := FixturePublicEvidenceCollector{Fixtures: map[string][]PublicEvidencePayload{"510300": {
		{SourceName: "cninfo", SourceLevel: model.SourceLevelA, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "cninfo-1", Title: "ETF 公告", Text: "公告正文 A", URL: "https://example.invalid/a", PublishedAt: "2026-06-05T00:00:00+08:00", CapturedAt: capturedAt},
		{SourceName: "cninfo", SourceLevel: model.SourceLevelA, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "cninfo-2", Title: "ETF 公告", Text: "公告正文 A2", URL: "https://example.invalid/a2", PublishedAt: "2026-06-05T00:00:00+08:00", CapturedAt: capturedAt},
		{SourceName: "financial_news", SourceLevel: model.SourceLevelB, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "news-1", Title: "ETF 公告", Text: "新闻正文", URL: "https://example.invalid/news", PublishedAt: "2026-06-05T00:00:00+08:00", CapturedAt: capturedAt},
	}}}
	service := PublicEvidenceIngestionService{Collector: collector, IntelligenceRepo: intelligenceRepo, AuditRepo: auditRepo, GenerateAuditID: testIDGenerator()}
	if err := service.IngestPublicEvidence(ctx, "510300", time.Time{}, time.Time{}); err != nil {
		t.Fatalf("IngestPublicEvidence: %v", err)
	}
	verification, err := intelligenceRepo.GetLatestSourceVerificationByFilter(ctx, "510300", publicEvidenceEventID(collector.Fixtures["510300"][0]))
	if err != nil {
		t.Fatal(err)
	}
	if verification.VerificationStatus != "satisfied" || verification.IndependentSourceCount != 2 || verification.HighGradeIndependentSourceCount != 1 {
		t.Fatalf("expected satisfied normal evidence with only one high-grade independent source, got %+v", verification)
	}
}

func TestPublicEvidenceIngestionDoesNotDowngradeSatisfiedVerificationOnPartialRefresh(t *testing.T) {
	ctx := context.Background()
	store := testDB(t)
	intelligenceRepo := sqlite.NewIntelligenceRepository(store.DB)
	auditRepo := sqlite.NewAuditRepository(store.DB)
	capturedAt := time.Date(2026, 6, 5, 13, 0, 0, 0, time.UTC)
	cninfo := PublicEvidencePayload{SourceName: "cninfo", SourceLevel: model.SourceLevelA, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "cninfo-1", Title: "ETF 公告", Text: "公告正文 A", URL: "https://example.invalid/a", PublishedAt: "2026-06-05T00:00:00+08:00", CapturedAt: capturedAt}
	szse := PublicEvidencePayload{SourceName: "szse", SourceLevel: model.SourceLevelA, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "szse-1", Title: "ETF 公告", Text: "公告正文 B", URL: "https://example.invalid/b", PublishedAt: "2026-06-05T00:00:00+08:00", CapturedAt: capturedAt}
	collector := &CompositePublicEvidenceCollector{Collectors: []PublicEvidenceCollector{FixturePublicEvidenceCollector{Fixtures: map[string][]PublicEvidencePayload{"510300": {cninfo, szse}}}}}
	service := PublicEvidenceIngestionService{Collector: collector, IntelligenceRepo: intelligenceRepo, AuditRepo: auditRepo, GenerateAuditID: testIDGenerator()}
	if err := service.IngestPublicEvidence(ctx, "510300", time.Time{}, time.Time{}); err != nil {
		t.Fatalf("first ingest: %v", err)
	}
	verification, err := intelligenceRepo.GetLatestSourceVerificationByFilter(ctx, "510300", publicEvidenceEventID(cninfo))
	if err != nil {
		t.Fatal(err)
	}
	if verification.VerificationStatus != "satisfied" {
		t.Fatalf("expected initial satisfied verification, got %+v", verification)
	}
	collector.Collectors = []PublicEvidenceCollector{failingPublicEvidenceCollector{}, FixturePublicEvidenceCollector{Fixtures: map[string][]PublicEvidencePayload{"510300": {szse}}}}
	if err := service.IngestPublicEvidence(ctx, "510300", time.Time{}, time.Time{}); err != nil {
		t.Fatalf("partial refresh ingest: %v", err)
	}
	verification, err = intelligenceRepo.GetLatestSourceVerificationByFilter(ctx, "510300", publicEvidenceEventID(cninfo))
	if err != nil {
		t.Fatal(err)
	}
	if verification.VerificationStatus != "satisfied" || verification.IndependentSourceCount != 2 || verification.HighGradeIndependentSourceCount != 2 {
		t.Fatalf("partial refresh should preserve satisfied verification, got %+v", verification)
	}
}

func TestPublicEvidenceIngestionWritesLatestVerificationForIncrementalSources(t *testing.T) {
	ctx := context.Background()
	store := testDB(t)
	intelligenceRepo := sqlite.NewIntelligenceRepository(store.DB)
	auditRepo := sqlite.NewAuditRepository(store.DB)
	capturedAt := time.Date(2026, 6, 5, 13, 0, 0, 0, time.UTC)
	first := FixturePublicEvidenceCollector{Fixtures: map[string][]PublicEvidencePayload{"510300": {{SourceName: "cninfo", SourceLevel: model.SourceLevelA, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "cninfo-1", Title: "ETF 公告", Text: "公告正文 A", URL: "https://example.invalid/a", PublishedAt: "2026-06-05T00:00:00+08:00", CapturedAt: capturedAt}}}}
	service := PublicEvidenceIngestionService{Collector: first, IntelligenceRepo: intelligenceRepo, AuditRepo: auditRepo, GenerateAuditID: testIDGenerator()}
	if err := service.IngestPublicEvidence(ctx, "510300", time.Time{}, time.Time{}); err != nil {
		t.Fatalf("first ingest: %v", err)
	}
	second := FixturePublicEvidenceCollector{Fixtures: map[string][]PublicEvidencePayload{"510300": {
		first.Fixtures["510300"][0],
		{SourceName: "szse", SourceLevel: model.SourceLevelA, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "szse-1", Title: "ETF 公告", Text: "公告正文 B", URL: "https://example.invalid/b", PublishedAt: "2026-06-05T00:00:00+08:00", CapturedAt: capturedAt},
	}}}
	service.Collector = second
	if err := service.IngestPublicEvidence(ctx, "510300", time.Time{}, time.Time{}); err != nil {
		t.Fatalf("second ingest: %v", err)
	}
	verification, err := intelligenceRepo.GetLatestSourceVerificationByFilter(ctx, "510300", publicEvidenceEventID(first.Fixtures["510300"][0]))
	if err != nil {
		t.Fatal(err)
	}
	if verification.VerificationStatus != "satisfied" || verification.IndependentSourceCount != 2 || verification.HighGradeIndependentSourceCount != 2 {
		t.Fatalf("expected latest verification to be satisfied after second source, got %+v", verification)
	}
}
func TestPublicEvidenceIngestionValidatesDependencies(t *testing.T) {
	service := PublicEvidenceIngestionService{}
	if err := service.IngestPublicEvidence(context.Background(), "510300", time.Time{}, time.Time{}); !apperr.IsCode(err, apperr.CodeInternalError) {
		t.Fatalf("expected stable internal error for missing dependencies, got %v", err)
	}
}

func TestPublicEvidenceIngestionAuditsPartialSourceFailures(t *testing.T) {
	ctx := context.Background()
	store := testDB(t)
	intelligenceRepo := sqlite.NewIntelligenceRepository(store.DB)
	auditRepo := sqlite.NewAuditRepository(store.DB)
	collector := &CompositePublicEvidenceCollector{Collectors: []PublicEvidenceCollector{
		failingPublicEvidenceCollector{},
		FixturePublicEvidenceCollector{Fixtures: map[string][]PublicEvidencePayload{"510300": {{SourceName: "szse", SourceLevel: model.SourceLevelA, SourceType: "public_disclosure", EvidenceRole: "formal", Symbol: "510300", SourceRecordID: "szse-1", Title: "ETF 公告", Text: "公告正文", URL: "https://example.invalid/b", PublishedAt: "2026-06-05T00:00:00+08:00", CapturedAt: time.Date(2026, 6, 5, 13, 0, 0, 0, time.UTC)}}}},
	}}
	service := PublicEvidenceIngestionService{Collector: collector, IntelligenceRepo: intelligenceRepo, AuditRepo: auditRepo, GenerateAuditID: testIDGenerator()}
	if err := service.IngestPublicEvidence(ctx, "510300", time.Time{}, time.Time{}); err != nil {
		t.Fatalf("IngestPublicEvidence: %v", err)
	}
	events, err := auditRepo.ListAuditEvents(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) < 2 || events[1].Action != "run_local_task" || events[1].Status != "degraded" || events[1].ErrorCode != "cninfo:network" {
		t.Fatalf("expected degraded audit for partial source failure, got %+v", events)
	}
}

func TestPublicEvidenceIngestionTreatsAllSourceNoDataAsEmptySuccess(t *testing.T) {
	ctx := context.Background()
	store := testDB(t)
	intelligenceRepo := sqlite.NewIntelligenceRepository(store.DB)
	auditRepo := sqlite.NewAuditRepository(store.DB)
	collector := &CompositePublicEvidenceCollector{Collectors: []PublicEvidenceCollector{
		noDataPublicEvidenceCollector{source: "cninfo"},
		noDataPublicEvidenceCollector{source: "szse"},
	}}
	service := PublicEvidenceIngestionService{Collector: collector, IntelligenceRepo: intelligenceRepo, AuditRepo: auditRepo, GenerateAuditID: testIDGenerator()}

	if err := service.IngestPublicEvidence(ctx, "510300", time.Time{}, time.Time{}); err != nil {
		t.Fatalf("all-source no_data should be an empty successful refresh, got %v", err)
	}
	summaries, err := intelligenceRepo.ListEvidenceSummaries(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 0 {
		t.Fatalf("expected no summaries for empty no_data refresh, got %+v", summaries)
	}
	events, err := auditRepo.ListAuditEvents(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 3 {
		t.Fatalf("expected two degraded no_data audits and one success audit, got %+v", events)
	}
	statuses := map[string]int{}
	codes := map[string]bool{}
	for _, event := range events {
		statuses[event.Status]++
		codes[event.ErrorCode] = true
		if event.Status == "failed" {
			t.Fatalf("no_data refresh must not write failed audit: %+v", events)
		}
	}
	if statuses["success"] != 1 || statuses["degraded"] != 2 || !codes["cninfo:no_data"] || !codes["szse:no_data"] {
		t.Fatalf("expected degraded no_data diagnostics plus success audit, got %+v", events)
	}
	for _, event := range events {
		if event.Action != "run_local_task" {
			t.Fatalf("expected public evidence audits to use run_local_task, got %+v", events)
		}
	}
}

func TestPublicEvidenceIngestionAuditsSourceFailures(t *testing.T) {
	ctx := context.Background()
	store := testDB(t)
	auditRepo := sqlite.NewAuditRepository(store.DB)
	service := PublicEvidenceIngestionService{Collector: failingPublicEvidenceCollector{}, IntelligenceRepo: sqlite.NewIntelligenceRepository(store.DB), AuditRepo: auditRepo, GenerateAuditID: testIDGenerator()}
	if err := service.IngestPublicEvidence(ctx, "510300", time.Time{}, time.Time{}); err == nil {
		t.Fatal("expected source failure")
	}
	events, err := auditRepo.ListAuditEvents(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 || events[0].Action != "run_local_task" || events[0].Status != "failed" || events[0].ErrorCode != "cninfo:network" || events[0].OutputRef != "source=cninfo count=0" {
		t.Fatalf("expected source-specific audit failure, got %+v", events)
	}
}

type noDataPublicEvidenceCollector struct{ source string }

func (c noDataPublicEvidenceCollector) FetchPublicEvidence(context.Context, string, time.Time, time.Time) ([]PublicEvidencePayload, error) {
	return nil, PublicEvidenceError{SourceName: c.source, ErrorCode: "no_data", Count: 0, Err: apperr.New(apperr.CodeNotFound, apperr.CategoryNotFound, "窗口内无公告")}
}

type failingPublicEvidenceCollector struct{}

func (failingPublicEvidenceCollector) FetchPublicEvidence(context.Context, string, time.Time, time.Time) ([]PublicEvidencePayload, error) {
	return nil, PublicEvidenceError{SourceName: "cninfo", ErrorCode: "source_unavailable", Count: 0, Err: apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "源不可用")}
}
func testIDGenerator() func() string {
	counter := 0
	return func() string {
		counter++
		return fmt.Sprintf("test-id-%03d", counter)
	}
}

func testDB(t *testing.T) *sqlite.Store {
	store, err := sqlite.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	if err := sqlite.Migrate(context.Background(), store.DB); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}
