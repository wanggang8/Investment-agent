package service

import (
	"investment-agent/internal/application/workflow"
	"investment-agent/internal/domain/model"
)

type RetrievalQualityFixture struct {
	Query               string
	Symbol              string
	ExpectedEvidenceIDs []string
	FormalOnly          bool
}

type RetrievalQualityEvaluation struct {
	Query                           string
	Symbol                          string
	Status                          string
	TopK                            int
	FallbackSource                  string
	IndexHealth                     string
	DegradedReason                  string
	MissingExpectedEvidenceIDs      []string
	UnexpectedBackgroundEvidenceIDs []string
}

func EvaluateRetrievalQuality(fixture RetrievalQualityFixture, result workflow.RetrievalResult) RetrievalQualityEvaluation {
	seen := make(map[string]model.Evidence, len(result.EvidenceSet.Items))
	for _, item := range result.EvidenceSet.Items {
		seen[item.EvidenceID] = item
	}
	missing := []string{}
	for _, expected := range fixture.ExpectedEvidenceIDs {
		if _, ok := seen[expected]; !ok {
			missing = append(missing, expected)
		}
	}
	background := []string{}
	if fixture.FormalOnly {
		for _, item := range result.EvidenceSet.Items {
			if item.Role != model.EvidenceFormal || item.SourceLevel == model.SourceLevelC {
				background = append(background, item.EvidenceID)
			}
		}
	}
	status := "hit"
	if len(missing) > 0 {
		status = "miss"
	}
	if fixture.FormalOnly && len(background) > 0 {
		status = "miss"
	}
	return RetrievalQualityEvaluation{Query: fixture.Query, Symbol: fixture.Symbol, Status: status, TopK: result.QualitySummary.TopK, FallbackSource: result.QualitySummary.FallbackSource, IndexHealth: result.QualitySummary.IndexHealth, DegradedReason: result.QualitySummary.DegradedReason, MissingExpectedEvidenceIDs: missing, UnexpectedBackgroundEvidenceIDs: background}
}
