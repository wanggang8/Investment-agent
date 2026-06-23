package service

import (
	"testing"

	"investment-agent/internal/application/workflow"
	"investment-agent/internal/domain/model"
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
