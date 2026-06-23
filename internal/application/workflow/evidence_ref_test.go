package workflow

import (
	"testing"

	"investment-agent/internal/domain/model"
)

func TestBuildEvidenceRefsPreservesRetrievedMetadata(t *testing.T) {
	wf := WorkflowContext{
		DecisionID: "decision_1",
		EvidenceSet: model.EvidenceSet{Items: []model.Evidence{{
			EvidenceID:      "ev_1",
			SummaryID:       "sum_1",
			SourceName:      "交易所公告",
			SourceLevel:     model.SourceLevelA,
			Role:            model.EvidenceFormal,
			PublishedAt:     "2026-05-29T01:00:00Z",
			CapturedAt:      "2026-05-29T02:00:00Z",
			OriginalURL:     "https://example.com/a",
			Summary:         "真实摘要",
			ContentHash:     "hash_1",
			TimeWeight:                      0.8,
			RelevanceScore:                  0.9,
			IndependentSourceCount:          3,
			HighGradeIndependentSourceCount: 2,
			EventType:                       model.EventMajorNegative,
		}}, VerificationStatus: model.VerificationSatisfied},
	}

	refs := buildEvidenceRefs(wf)

	if len(refs) != 1 {
		t.Fatalf("expected one ref, got %+v", refs)
	}
	ref := refs[0]
	if ref.SummaryID != "sum_1" || ref.SourceName != "交易所公告" || ref.Summary != "真实摘要" || ref.OriginalURL != "https://example.com/a" || ref.ContentHash != "hash_1" {
		t.Fatalf("metadata not preserved: %+v", ref)
	}
	if ref.TimeWeight != 0.8 || ref.RelevanceScore != 0.9 || ref.EvidenceRole != "formal" || ref.SourceLevel != "A" || ref.IndependentSourceCount != 3 || ref.HighGradeIndependentSourceCount != 2 {
		t.Fatalf("score/source fields not preserved: %+v", ref)
	}
}
