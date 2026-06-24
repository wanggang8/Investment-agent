package service

import (
	"context"
	"testing"

	"investment-agent/internal/application/workflow"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
)

func TestRetrievalQualityFixtureCoversRepresentativeEvidenceClasses(t *testing.T) {
	fixtures := []struct {
		name      string
		fixture   RetrievalQualityFixture
		result    workflow.RetrievalResult
		wantState string
	}{
		{
			name:    "announcement formal hit",
			fixture: RetrievalQualityFixture{Query: "510300 基金公告", Symbol: "510300", ExpectedEvidenceIDs: []string{"ev_announcement"}, FormalOnly: true},
			result: workflow.RetrievalResult{EvidenceSet: model.EvidenceSet{Items: []model.Evidence{
				{EvidenceID: "ev_announcement", SourceLevel: model.SourceLevelA, Role: model.EvidenceFormal},
			}}, QualitySummary: workflow.RetrievalQualitySummary{TopK: 1, FallbackSource: "veclite", IndexHealth: VectorIndexHealthHealthy}},
			wantState: "hit",
		},
		{
			name:    "regulatory formal hit",
			fixture: RetrievalQualityFixture{Query: "基金监管规则", Symbol: "510300", ExpectedEvidenceIDs: []string{"ev_regulatory"}, FormalOnly: true},
			result: workflow.RetrievalResult{EvidenceSet: model.EvidenceSet{Items: []model.Evidence{
				{EvidenceID: "ev_regulatory", SourceLevel: model.SourceLevelS, Role: model.EvidenceFormal},
			}}, QualitySummary: workflow.RetrievalQualitySummary{TopK: 1, FallbackSource: "sqlite_summary", IndexHealth: VectorIndexHealthMissing}},
			wantState: "hit",
		},
		{
			name:    "ETF info formal hit",
			fixture: RetrievalQualityFixture{Query: "510300 ETF 信息", Symbol: "510300", ExpectedEvidenceIDs: []string{"ev_etf_info"}, FormalOnly: true},
			result: workflow.RetrievalResult{EvidenceSet: model.EvidenceSet{Items: []model.Evidence{
				{EvidenceID: "ev_etf_info", SourceLevel: model.SourceLevelB, Role: model.EvidenceFormal},
			}}, QualitySummary: workflow.RetrievalQualitySummary{TopK: 1, FallbackSource: "veclite", IndexHealth: VectorIndexHealthHealthy}},
			wantState: "hit",
		},
		{
			name:    "market background is diagnostic only",
			fixture: RetrievalQualityFixture{Query: "市场情绪背景", Symbol: "000300", ExpectedEvidenceIDs: []string{"ev_market_background"}, FormalOnly: false},
			result: workflow.RetrievalResult{EvidenceSet: model.EvidenceSet{Items: []model.Evidence{
				{EvidenceID: "ev_market_background", SourceLevel: model.SourceLevelC, Role: model.EvidenceBackground},
			}}, QualitySummary: workflow.RetrievalQualitySummary{TopK: 1, FallbackSource: "sqlite_summary", IndexHealth: VectorIndexHealthMissing}},
			wantState: "hit",
		},
		{
			name:    "C level background cannot satisfy formal slot",
			fixture: RetrievalQualityFixture{Query: "C 级舆情作正式证据", Symbol: "000300", ExpectedEvidenceIDs: []string{"ev_c_background"}, FormalOnly: true},
			result: workflow.RetrievalResult{EvidenceSet: model.EvidenceSet{Items: []model.Evidence{
				{EvidenceID: "ev_c_background", SourceLevel: model.SourceLevelC, Role: model.EvidenceBackground},
			}}, QualitySummary: workflow.RetrievalQualitySummary{TopK: 1, FallbackSource: "sqlite_summary", IndexHealth: VectorIndexHealthMissing}},
			wantState: "miss",
		},
	}

	for _, tt := range fixtures {
		t.Run(tt.name, func(t *testing.T) {
			got := EvaluateRetrievalQuality(tt.fixture, tt.result)
			if got.Status != tt.wantState {
				t.Fatalf("expected %s, got %+v", tt.wantState, got)
			}
			if tt.fixture.FormalOnly && tt.result.EvidenceSet.Items[0].SourceLevel == model.SourceLevelC && len(got.UnexpectedBackgroundEvidenceIDs) != 1 {
				t.Fatalf("expected C-level background diagnostic, got %+v", got)
			}
		})
	}
}

func TestRetrievalQualityEvaluationAcceptsRerankedExpectedEvidence(t *testing.T) {
	repo := intelligenceRepoForVectorTest{summaries: []repository.IntelligenceSummary{
		{SummaryID: "sum_background", Symbol: "510300", SourceLevel: "C", EvidenceRole: "background", EventType: "normal", VerificationStatus: "background_only", Summary: "背景材料"},
		{SummaryID: "sum_valuation", Symbol: "510300", SourceLevel: "A", EvidenceRole: "formal", EventType: "normal", VerificationStatus: "satisfied", VerificationEvidenceIDsJSON: `["sum_valuation"]`, VerificationEvidenceRole: "formal", VerificationEventType: "normal", VerificationHighestSourceLevel: "A", Summary: "估值 分位 低估 买入纪律"},
	}}
	index := &recordingSemanticVectorIndex{chunks: []repository.RAGChunk{
		{ChunkID: "chunk_background", SummaryID: "sum_background", Symbol: "510300", ChunkText: "背景材料", IndexStatus: "indexed"},
		{ChunkID: "chunk_valuation", SummaryID: "sum_valuation", Symbol: "510300", ChunkText: "估值 分位 低估 买入纪律", IndexStatus: "indexed"},
	}}
	adapter := NewRetrievalAdapter(transactorStub{repos: repository.Repositories{IntelligenceRepo: repo}}, index)

	result, err := adapter.RetrieveEvidence(context.Background(), workflow.RetrievalRequest{Symbol: "510300", Query: "估值低估能买吗", TopK: 1})
	if err != nil {
		t.Fatalf("RetrieveEvidence: %v", err)
	}
	evaluation := EvaluateRetrievalQuality(RetrievalQualityFixture{Query: "估值低估能买吗", Symbol: "510300", ExpectedEvidenceIDs: []string{"sum_valuation"}, FormalOnly: true}, result)

	if evaluation.Status != "hit" || len(evaluation.MissingExpectedEvidenceIDs) != 0 || len(evaluation.UnexpectedBackgroundEvidenceIDs) != 0 {
		t.Fatalf("expected reranked formal valuation evidence to satisfy retrieval quality, got %+v result=%+v", evaluation, result.EvidenceSet.Items)
	}
}
