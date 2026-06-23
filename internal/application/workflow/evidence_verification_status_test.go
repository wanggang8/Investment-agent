package workflow

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"investment-agent/internal/domain/model"
)

type evidenceVerificationSourceStub struct {
	items []IntelligenceSourceItem
}

func (s evidenceVerificationSourceStub) FetchIntelligence(context.Context, string) ([]IntelligenceSourceItem, error) {
	return s.items, nil
}

func TestEvidenceVerificationUsesActualFetchedSources(t *testing.T) {
	graph := NewEvidenceVerificationGraphWithDependencies(WorkflowDependencies{IntelligenceSource: evidenceVerificationSourceStub{items: []IntelligenceSourceItem{
		{SourceName: "交易所", SourceLevel: model.SourceLevelA, Title: "公告", Text: "510300 公告", URL: "https://example.com/a", PublishedAt: "2026-01-01T00:00:00Z"},
		{SourceName: "基金公司", SourceLevel: model.SourceLevelA, Title: "披露", Text: "510300 披露", URL: "https://example.com/b", PublishedAt: "2026-01-01T01:00:00Z"},
	}}})

	out, err := graph.Run(context.Background(), EvidenceVerificationInput{RequestID: "req_actual_sources", Symbol: "510300", Sources: []string{"official"}})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if out.SourceVerifications[0] != model.VerificationSatisfied {
		t.Fatalf("expected actual high-grade independent sources to satisfy verification, got %+v", out.SourceVerifications)
	}
	if len(out.IntelligenceSummaries) != 2 || out.IntelligenceSummaries[0].IndependentSourceCount != 2 || out.IntelligenceSummaries[0].HighGradeIndependentSourceCount != 2 {
		t.Fatalf("expected actual source counts on summaries, got %+v", out.IntelligenceSummaries)
	}
}

func TestEvidenceVerificationUsesPublicHTTPPayloadSources(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("symbol") != "510300" {
			t.Fatalf("symbol query not forwarded: %s", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"items":[{"source":"exchange_disclosure","title":"ETF 公告","content":"510300 公告正文","url":"https://example.invalid/a","published_at":"2026-06-01T00:00:00Z"},{"source":"exchange_disclosure","title":"ETF 公告","content":"重复正文","url":"https://example.invalid/a","published_at":"2026-06-01T00:00:00Z"},{"source":"fund_company_disclosure","title":"基金披露","content":"510300 基金披露","url":"https://example.invalid/b","published_at":"2026-06-01T01:00:00Z"}]}`))
	}))
	defer server.Close()

	graph := NewEvidenceVerificationGraphWithDependencies(WorkflowDependencies{IntelligenceSource: ConfiguredIntelligenceSource{Enabled: []string{"public_http"}, IntelligenceEndpoint: server.URL, HTTPClient: server.Client()}})
	out, err := graph.Run(context.Background(), EvidenceVerificationInput{RequestID: "req_public_sources", Symbol: "510300", Sources: []string{"public_http"}})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if out.SourceVerifications[0] != model.VerificationSatisfied {
		t.Fatalf("expected public high-grade independent sources to satisfy verification, got %+v", out.SourceVerifications)
	}
	if len(out.IntelligenceSummaries) != 2 || out.IntelligenceSummaries[0].IndependentSourceCount != 2 || out.IntelligenceSummaries[0].HighGradeIndependentSourceCount != 2 {
		t.Fatalf("expected deduped public source counts on summaries, got %+v", out.IntelligenceSummaries)
	}
}

