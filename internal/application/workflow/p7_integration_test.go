package workflow

import (
	"context"
	"errors"
	"strings"
	"testing"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/pkg/apperr"
)

type testMarketDataSource struct {
	point MarketDataPoint
	err   error
}

func (s testMarketDataSource) FetchMarketData(context.Context, string) (MarketDataPoint, error) {
	return s.point, s.err
}

type testIntelligenceSource struct {
	items []IntelligenceSourceItem
	err   error
}

func (s testIntelligenceSource) FetchIntelligence(context.Context, string) ([]IntelligenceSourceItem, error) {
	return s.items, s.err
}

type testAnalystService struct {
	reports  map[string]string
	metadata map[string]string
	err      error
}

func (s testAnalystService) Analyze(context.Context, AnalystRequest) (AnalystResponse, error) {
	return AnalystResponse{Reports: s.reports, Metadata: s.metadata}, s.err
}

type categorizedAnalystError struct {
	category string
	metadata map[string]string
	err      error
}

func (e categorizedAnalystError) Error() string { return e.err.Error() }
func (e categorizedAnalystError) Unwrap() error { return e.err }
func (e categorizedAnalystError) Category() string {
	return e.category
}
func (e categorizedAnalystError) Metadata() map[string]string {
	return e.metadata
}

type testRetrievalService struct {
	result RetrievalResult
	err    error
}

func (s testRetrievalService) RetrieveEvidence(context.Context, RetrievalRequest) (RetrievalResult, error) {
	return s.result, s.err
}

func TestEvidenceRetrievalFallsBackToSQLiteSummary(t *testing.T) {
	wf := &WorkflowContext{RequestID: "req_retrieve", Symbol: "510300"}
	deps := WorkflowDependencies{RetrievalService: testRetrievalService{result: RetrievalResult{EvidenceSet: model.EvidenceSet{Items: []model.Evidence{{EvidenceID: "sum1", SourceLevel: model.SourceLevelA, Role: model.EvidenceFormal, EventType: model.EventNormal}}, VerificationStatus: model.VerificationSatisfied}, OutputRef: "sum1", DegradedReason: "veclite unavailable", QualitySummary: RetrievalQualitySummary{TopK: 1, IndexHealth: "missing", FallbackSource: "sqlite_summary", SourceConsistencyStatus: "checked", DegradedReason: "veclite unavailable"}}}}

	result := RunEvidenceRetrievalNode(context.Background(), wf, deps)
	if result.Status != StatusDegraded || result.ErrorCode != ErrCodeVectorIndexUnavailable || result.Audit.OutputRef != "sum1:topk=1:fallback=sqlite_summary:index=missing:consistency=checked:degraded=veclite unavailable" || wf.RetrievalDegradedReason != "veclite unavailable" {
		t.Fatalf("expected degraded sqlite fallback, result=%+v wf=%+v", result, wf)
	}
	if wf.RetrievalQualitySummary.FallbackSource != "sqlite_summary" || wf.RetrievalQualitySummary.TopK != 1 {
		t.Fatalf("expected retrieval quality summary preserved, got %+v", wf.RetrievalQualitySummary)
	}
}

func TestEvidenceRetrievalReturnsEvidenceNotFoundWhenFallbackEmpty(t *testing.T) {
	wf := &WorkflowContext{RequestID: "req_retrieve_empty", Symbol: "510300"}
	deps := WorkflowDependencies{RetrievalService: testRetrievalService{result: RetrievalResult{DegradedReason: "sqlite summary empty"}}}

	result := RunEvidenceRetrievalNode(context.Background(), wf, deps)
	if result.Status != StatusFailed || result.ErrorCode != ErrCodeEvidenceNotFound || wf.RetrievalDegradedReason != "sqlite summary empty" {
		t.Fatalf("expected evidence not found, result=%+v wf=%+v", result, wf)
	}
}

func TestMarketRefreshUsesConfiguredDataSource(t *testing.T) {
	deps := WorkflowDependencies{MarketDataSource: testMarketDataSource{point: MarketDataPoint{PEPercentile: 18, PBPercentile: 27, VolumePercentile: 45, VolatilityPercentile: 12, ClosePrice: 4.2, TurnoverRate: 0.3}}}

	out, err := NewMarketRefreshGraphWithDependencies(deps).Run(context.Background(), MarketRefreshInput{RequestID: "req_market_source", Symbol: "510300"})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if out.MarketSnapshot.PEPercentile != 18 || out.MarketSnapshot.PBPercentile != 27 || out.MarketSnapshot.VolumePercentile != 45 || out.MarketSnapshot.VolatilityPercentile != 12 || out.MarketSnapshot.ClosePrice != 4.2 || out.MarketSnapshot.TurnoverRate != 0.3 {
		t.Fatalf("market snapshot did not use data source point: %+v", out.MarketSnapshot)
	}
	if out.MarketSnapshot.MarketMetricsJSON != `{"close_price":4.2,"turnover_rate":0.3}` {
		t.Fatalf("market metrics not preserved: %s", out.MarketSnapshot.MarketMetricsJSON)
	}
}

func TestMarketRefreshWritesAuditForStaleSourcePoint(t *testing.T) {
	deps := WorkflowDependencies{MarketDataSource: testMarketDataSource{point: MarketDataPoint{Stale: true}}}
	graph := NewMarketRefreshGraphWithDependencies(deps)

	out, err := graph.Run(context.Background(), MarketRefreshInput{RequestID: "req_market_stale", Symbol: "510300"})
	if !apperr.IsCode(err, apperr.CodeDataStale) {
		t.Fatalf("expected DATA_STALE, got %v", err)
	}
	if len(out.AuditEvents) != 1 || out.AuditEvents[0].Status != model.AuditStatusFailed || out.AuditEvents[0].ErrorCode != string(apperr.CodeDataStale) {
		t.Fatalf("expected failed stale audit event, out=%+v", out)
	}
}

func TestMarketRefreshWritesAuditForSourceError(t *testing.T) {
	deps := WorkflowDependencies{MarketDataSource: testMarketDataSource{err: errors.New("market source unavailable")}}
	graph := NewMarketRefreshGraphWithDependencies(deps)

	out, err := graph.Run(context.Background(), MarketRefreshInput{RequestID: "req_market_source_error", Symbol: "510300"})
	if !apperr.IsCode(err, apperr.CodeDataSourceUnavailable) {
		t.Fatalf("expected DATA_SOURCE_UNAVAILABLE, got %v", err)
	}
	if len(out.AuditEvents) != 1 || out.AuditEvents[0].Status != model.AuditStatusFailed || out.AuditEvents[0].ErrorCode != string(apperr.CodeDataSourceUnavailable) {
		t.Fatalf("expected failed source audit event, out=%+v", out)
	}
}

func TestEvidenceVerificationIngestsIntelligenceSourceAndRestrictsCLevel(t *testing.T) {
	deps := WorkflowDependencies{IntelligenceSource: testIntelligenceSource{items: []IntelligenceSourceItem{{SourceName: "manual", SourceLevel: model.SourceLevelC, Title: "背景材料", Text: "C级信源只能作为背景", URL: "https://example.invalid/item"}}}}

	out, err := NewEvidenceVerificationGraphWithDependencies(deps).Run(context.Background(), EvidenceVerificationInput{RequestID: "req_intel_source", Symbol: "510300", Sources: []string{"manual", "official"}})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(out.IntelligenceItems) != 1 || out.IntelligenceSummary == "" || len(out.RAGChunks) != 1 {
		t.Fatalf("expected source-backed intelligence facts, got %+v", out)
	}
	if out.WorkflowContext.EvidenceSet.Items[0].SourceLevel != model.SourceLevelC || out.WorkflowContext.EvidenceSet.Items[0].Role != model.EvidenceBackground || out.WorkflowContext.EvidenceSet.Items[0].IndependentSourceCount != 1 {
		t.Fatalf("C level source must be background only with actual source count: %+v", out.WorkflowContext.EvidenceSet.Items)
	}
}

func TestEvidenceVerificationGeneratesStableContentHashes(t *testing.T) {
	items := []IntelligenceSourceItem{{SourceName: "official", SourceLevel: model.SourceLevelA, Title: "公告", Text: "同一正文", URL: "https://example.invalid/a", PublishedAt: "2026-01-01T00:00:00Z"}}
	deps := WorkflowDependencies{IntelligenceSource: testIntelligenceSource{items: items}}

	first, err := NewEvidenceVerificationGraphWithDependencies(deps).Run(context.Background(), EvidenceVerificationInput{RequestID: "req_hash_1", Symbol: "510300", Sources: []string{"official", "exchange"}})
	if err != nil {
		t.Fatalf("first Run: %v", err)
	}
	second, err := NewEvidenceVerificationGraphWithDependencies(deps).Run(context.Background(), EvidenceVerificationInput{RequestID: "req_hash_2", Symbol: "510300", Sources: []string{"official", "exchange"}})
	if err != nil {
		t.Fatalf("second Run: %v", err)
	}
	if first.WorkflowContext.EvidenceSet.Items[0].ContentHash == "" || first.WorkflowContext.EvidenceSet.Items[0].ContentHash != second.WorkflowContext.EvidenceSet.Items[0].ContentHash {
		t.Fatalf("expected stable content hash, first=%+v second=%+v", first.WorkflowContext.EvidenceSet.Items, second.WorkflowContext.EvidenceSet.Items)
	}
	changed, err := NewEvidenceVerificationGraphWithDependencies(WorkflowDependencies{IntelligenceSource: testIntelligenceSource{items: []IntelligenceSourceItem{{SourceName: "official", SourceLevel: model.SourceLevelA, Title: "公告", Text: "变更正文", URL: "https://example.invalid/a", PublishedAt: "2026-01-01T00:00:00Z"}}}}).Run(context.Background(), EvidenceVerificationInput{RequestID: "req_hash_3", Symbol: "510300", Sources: []string{"official", "exchange"}})
	if err != nil {
		t.Fatalf("changed Run: %v", err)
	}
	if first.WorkflowContext.EvidenceSet.Items[0].ContentHash == changed.WorkflowContext.EvidenceSet.Items[0].ContentHash || first.WorkflowContext.EvidenceSet.Items[0].ChunkHash == changed.WorkflowContext.EvidenceSet.Items[0].ChunkHash {
		t.Fatalf("expected hashes to change with content, first=%+v changed=%+v", first.WorkflowContext.EvidenceSet.Items, changed.WorkflowContext.EvidenceSet.Items)
	}
}

func TestAnalystServiceProvidesMaterialsWithoutChangingFinalVerdict(t *testing.T) {
	wf := &WorkflowContext{RequestID: "req_analyst", Symbol: "510300", AnalystReports: map[string]string{}, RuleVerdict: model.RuleVerdict{Status: model.VerdictHold}}
	deps := WorkflowDependencies{AnalystService: testAnalystService{reports: map[string]string{"value": "估值分析材料", "trend_risk": "趋势风险材料", "expected_return": "预期收益材料"}, metadata: map[string]string{"prompt_version": "p37-analyst-v1", "model": "gpt-5.4-mini", "parse_status": "parsed", "quality_status": "passed", "input_summary": "510300", "output_summary": "估值分析材料"}}}

	result := RunValueAnalystNode(context.Background(), wf, deps)
	if result.Status != StatusSuccess || wf.AnalystReports["value"] == "" {
		t.Fatalf("expected analyst material, result=%+v reports=%+v", result, wf.AnalystReports)
	}
	if wf.RuleVerdict.Status != model.VerdictHold {
		t.Fatalf("analyst service must not change final verdict: %+v", wf.RuleVerdict)
	}
	if wf.AnalystReportMetadata["value"]["prompt_version"] != "p37-analyst-v1" || wf.AnalystReportMetadata["value"]["model"] != "gpt-5.4-mini" {
		t.Fatalf("expected analyst metadata preserved: %+v", wf.AnalystReportMetadata)
	}
	record := buildDecisionRecord(*wf)
	if !strings.Contains(record.AnalystReportsJSON, `"prompt_version":"p37-analyst-v1"`) || !strings.Contains(record.AnalystReportsJSON, `"quality_status":"passed"`) {
		t.Fatalf("expected structured analyst metadata in decision record: %s", record.AnalystReportsJSON)
	}
}

func TestAnalystServiceUnavailableDegrades(t *testing.T) {
	wf := &WorkflowContext{RequestID: "req_analyst_down", Symbol: "510300", AnalystReports: map[string]string{}}
	deps := WorkflowDependencies{AnalystService: testAnalystService{err: errors.New("deepseek unavailable")}}

	result := RunValueAnalystNode(context.Background(), wf, deps)
	if result.Status != StatusDegraded || result.ErrorCode != ErrCodeAnalystUnavailable {
		t.Fatalf("expected degraded analyst result, got %+v", result)
	}
}

func TestAnalystServiceUnavailableIncludesStableCategory(t *testing.T) {
	wf := &WorkflowContext{RequestID: "req_analyst_timeout", Symbol: "510300", AnalystReports: map[string]string{}}
	deps := WorkflowDependencies{AnalystService: testAnalystService{err: categorizedAnalystError{category: "timeout", err: errors.New("deepseek timeout"), metadata: map[string]string{"prompt_version": "p37-analyst-v1", "model": "gpt-5.4-mini", "parse_status": "timeout", "quality_status": "not_evaluated"}}}}

	result := RunValueAnalystNode(context.Background(), wf, deps)

	if result.Status != StatusDegraded || result.ErrorCode != ErrCodeAnalystUnavailable {
		t.Fatalf("expected degraded analyst result, got %+v", result)
	}
	if result.Audit.OutputRef != "value:category=model_unavailable:prompt=p37-analyst-v1:model=gpt-5.4-mini:parse=timeout:quality=not_evaluated" {
		t.Fatalf("expected categorized analyst output ref, got %q", result.Audit.OutputRef)
	}
}

func TestAnalystAuditOutputSummaryUsesStableSeparators(t *testing.T) {
	got := summarizeAuditRef("结论: 不包含=最终裁决\n仅分析")
	if strings.Contains(got, ":") || strings.Contains(got, "=") || strings.Contains(got, "\n") {
		t.Fatalf("expected stable audit ref separators, got %q", got)
	}
}

func TestRetrievalAuditOutputRefSanitizesBase(t *testing.T) {
	got := retrievalAuditOutputRef("/private/tmp/sk-abc123:chunk=1", RetrievalQualitySummary{TopK: 1, FallbackSource: "sqlite_summary"})
	if strings.Contains(got, "/private/") || strings.Contains(got, "sk-abc123") || strings.Contains(got, ":chunk=") {
		t.Fatalf("expected sanitized retrieval audit base, got %q", got)
	}
	if !strings.Contains(got, "topk=1") || !strings.Contains(got, "fallback=sqlite_summary") {
		t.Fatalf("expected retrieval quality fields retained, got %q", got)
	}
}

func TestExpectedReturnNodeUsesAnalystService(t *testing.T) {
	wf := &WorkflowContext{RequestID: "req_expected_return", Symbol: "510300", AnalystReports: map[string]string{}, ExpectedReturnSampleCount: 20, RuleVerdict: model.RuleVerdict{Status: model.VerdictHold}}
	deps := WorkflowDependencies{AnalystService: testAnalystService{reports: map[string]string{"expected_return": "预期收益材料"}}}

	result := RunExpectedReturnNode(context.Background(), wf, deps)
	if result.Status != StatusSuccess || wf.AnalystReports["expected_return"] != "预期收益材料" || len(wf.ExpectedReturnScenarios) == 0 {
		t.Fatalf("expected return analyst material missing, result=%+v reports=%+v scenarios=%+v", result, wf.AnalystReports, wf.ExpectedReturnScenarios)
	}
	if wf.RuleVerdict.Status != model.VerdictHold {
		t.Fatalf("expected return analyst must not change final verdict: %+v", wf.RuleVerdict)
	}
}

func TestExpectedReturnNodeDegradesWhenAnalystUnavailable(t *testing.T) {
	wf := &WorkflowContext{RequestID: "req_expected_return_down", Symbol: "510300", AnalystReports: map[string]string{}, ExpectedReturnSampleCount: 20}
	deps := WorkflowDependencies{AnalystService: testAnalystService{err: errors.New("deepseek unavailable")}}

	result := RunExpectedReturnNode(context.Background(), wf, deps)
	if result.Status != StatusDegraded || result.ErrorCode != ErrCodeAnalystUnavailable || len(wf.ExpectedReturnScenarios) == 0 {
		t.Fatalf("expected degraded expected return with local scenarios, result=%+v scenarios=%+v", result, wf.ExpectedReturnScenarios)
	}
}

func TestExpectedReturnNodeUsesSafeLocalMaterialWhenAnalystQualityFails(t *testing.T) {
	wf := &WorkflowContext{RequestID: "req_expected_return_quality", Symbol: "510300", AnalystReports: map[string]string{}, ExpectedReturnSampleCount: 20}
	deps := WorkflowDependencies{AnalystService: testAnalystService{err: categorizedAnalystError{category: "quality_failure", err: errors.New("llm output failed safety quality gate"), metadata: map[string]string{"prompt_version": "p37-analyst-v1", "model": "gpt-5.4-mini", "parse_status": "parsed", "quality_status": "failed"}}}}

	result := RunExpectedReturnNode(context.Background(), wf, deps)
	if result.Status != StatusSuccess || result.ErrorCode != "" {
		t.Fatalf("expected safe local expected return material after quality failure, result=%+v", result)
	}
	if !strings.Contains(wf.AnalystReports["expected_return"], "本地预期收益情景") {
		t.Fatalf("expected local fallback analyst material, got %q", wf.AnalystReports["expected_return"])
	}
	metadata := wf.AnalystReportMetadata["expected_return"]
	if metadata["model"] != "deterministic-local" || metadata["quality_status"] != "passed" || metadata["fallback_reason"] != "llm_quality_failure" {
		t.Fatalf("expected deterministic fallback metadata, got %+v", metadata)
	}
}
