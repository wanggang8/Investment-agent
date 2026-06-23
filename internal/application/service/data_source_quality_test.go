package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

func TestDataSourceQualityRegressionFixturePassesAndRedacts(t *testing.T) {
	svc := NewDataSourceQualityService(emptyQualityRepos())

	out, err := svc.Run(context.Background(), DataSourceQualityRegressionRequest{Mode: "fixture"})
	if err != nil {
		t.Fatalf("Run fixture: %v", err)
	}
	if out.Mode != "fixture" || out.Status != "passed" || len(out.Cases) != 6 || len(out.MissingCategories) != 0 {
		t.Fatalf("unexpected fixture regression response: %+v", out)
	}
	if out.Policy.Verdict != "passed" || out.Policy.ReleaseGate != "pass" {
		t.Fatalf("expected fixture policy pass, got %+v", out.Policy)
	}
	for _, forbidden := range []string{"自动修复", "自动确认", "自动应用规则", "一键交易", "代下单"} {
		if strings.Contains(out.Policy.SafetyNote, forbidden) {
			t.Fatalf("policy safety note should avoid forbidden affordance copy %q: %s", forbidden, out.Policy.SafetyNote)
		}
	}
	body := mustJSONText(t, out)
	for _, forbidden := range []string{"sk-123456789012", "sk-proj-abc_def-123456", "/Users/private", "select    *    from", "prompt:", "raw HTTP", "GET /secret HTTP/1.1", "HTTP/1.1 500", "BEGIN RSA PRIVATE KEY"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("fixture response leaked %q: %s", forbidden, body)
		}
	}
	foundRedactionCase := false
	for _, item := range out.Cases {
		if item.CaseID == "redaction" {
			foundRedactionCase = true
			if item.DiagnosticPreview == "" || strings.Contains(item.DiagnosticPreview, "secret") {
				t.Fatalf("expected sanitized redaction diagnostic, got %+v", item)
			}
		}
	}
	if !foundRedactionCase {
		t.Fatalf("expected redaction case, got %+v", out.Cases)
	}
}

func TestDataSourceQualityRegressionCurrentClassifiesSourceHealth(t *testing.T) {
	metrics := `{"source_name":"csindex","source_level":"A","source_type":"index_basic","captured_at":"2026-06-06T01:00:00Z","metadata":{"p34_source_health":{"index_constituents":{"freshness":"fresh","data_date":"2026-06-05","last_success_at":"2026-06-06T01:00:00Z","affected_symbols":["000300"],"source_level":"A"},"index_weights":{"freshness":"no_data","data_date":"2026-06-05","last_failure_at":"2026-06-06T01:00:00Z","failure_category":"no_data","affected_symbols":["000300"],"source_level":"A"},"index_valuation_files":{"freshness":"parse_error","data_date":"2026-06-05","last_failure_at":"2026-06-06T01:00:00Z","failure_category":"parse_error","affected_symbols":["000300"],"source_level":"A"},"capital_flow":{"freshness":"source_unavailable","data_date":"2026-06-05","last_failure_at":"2026-06-06T01:00:00Z","failure_category":"source_unavailable","affected_symbols":["000300"],"source_level":"B"},"sentiment_proxy":{"freshness":"stale","data_date":"2026-06-01","last_failure_at":"2026-06-06T01:00:00Z","failure_category":"stale","affected_symbols":["000300"],"source_level":"C"}},"p34_data_categories":["index_constituents","index_weights","index_valuation_files","capital_flow","sentiment_proxy"]}}`
	svc := NewDataSourceQualityService(qualityReposWithMarket(model.MarketSnapshot{MarketSnapshotID: "market_health", Symbol: "000300", TradeDate: "2026-06-05", MarketMetricsJSON: metrics}))

	out, err := svc.Run(context.Background(), DataSourceQualityRegressionRequest{Mode: "current", Symbol: "000300"})
	if err != nil {
		t.Fatalf("Run current: %v", err)
	}
	if out.Mode != "current" || out.Status != "degraded" || len(out.Cases) != 5 {
		t.Fatalf("unexpected current regression response: %+v", out)
	}
	got := map[string]string{}
	for _, item := range out.Cases {
		got[item.DataCategory] = item.ActualFreshness + ":" + item.Status
	}
	for category, want := range map[string]string{
		"index_constituents":    "fresh:passed",
		"index_weights":         "no_data:degraded",
		"index_valuation_files": "parse_error:degraded",
		"capital_flow":          "source_unavailable:degraded",
		"sentiment_proxy":       "stale:degraded",
	} {
		if got[category] != want {
			t.Fatalf("expected %s=%s, got map %+v", category, want, got)
		}
	}
	if len(out.MissingCategories) != 4 {
		t.Fatalf("expected degraded categories as missing categories, got %+v", out.MissingCategories)
	}
	if out.Policy.Verdict != "blocked" || out.Policy.ReleaseGate != "block" || out.Policy.BlockingCount == 0 {
		t.Fatalf("expected core degraded categories to block release policy, got %+v", out.Policy)
	}
}

func TestDataSourceQualityRegressionCurrentMissingHealthDegrades(t *testing.T) {
	svc := NewDataSourceQualityService(qualityReposWithMarket(model.MarketSnapshot{MarketSnapshotID: "market_empty", Symbol: "000300", TradeDate: "2026-06-05", MarketMetricsJSON: `{}`}))

	out, err := svc.Run(context.Background(), DataSourceQualityRegressionRequest{Mode: "current"})
	if err != nil {
		t.Fatalf("Run current missing health: %v", err)
	}
	if out.Status != "degraded" || len(out.Cases) != 1 || out.Cases[0].ActualFreshness != "missing" || out.MissingCategories[0] != "p34_source_health" {
		t.Fatalf("expected degraded missing source health, got %+v", out)
	}
	if out.Policy.Verdict != "blocked" || out.Policy.ReleaseGate != "block" {
		t.Fatalf("expected missing source health to block policy, got %+v", out.Policy)
	}
}

func TestDataSourceQualityRegressionCurrentUnknownFreshnessFails(t *testing.T) {
	metrics := `{"source_name":"csindex","source_level":"A","source_type":"index_basic","metadata":{"p34_source_health":{"index_weights":{"freshness":"mystery_status","data_date":"2026-06-05","affected_symbols":["000300"]}},"p34_data_categories":["index_weights"]}}`
	svc := NewDataSourceQualityService(qualityReposWithMarket(model.MarketSnapshot{MarketSnapshotID: "market_unknown", Symbol: "000300", TradeDate: "2026-06-05", MarketMetricsJSON: metrics}))

	out, err := svc.Run(context.Background(), DataSourceQualityRegressionRequest{Mode: "current"})
	if err != nil {
		t.Fatalf("Run current unknown: %v", err)
	}
	if out.Status != "failed" || out.Cases[0].Status != "failed" || !strings.Contains(out.Cases[0].DiagnosticPreview, "unrecognized freshness") {
		t.Fatalf("expected failed unknown freshness, got %+v", out)
	}
	if out.Policy.Verdict != "blocked" || out.Policy.ReleaseGate != "block" {
		t.Fatalf("expected unknown freshness to block policy, got %+v", out.Policy)
	}
}

func TestDataSourceQualityRegressionCurrentOptionalDegradedRequiresWaiver(t *testing.T) {
	metrics := `{"source_name":"sentiment_proxy_fixture","source_level":"C","source_type":"sentiment_proxy","metadata":{"p34_source_health":{"sentiment_proxy":{"freshness":"stale","data_date":"2026-06-01","last_failure_at":"2026-06-06T01:00:00Z","failure_category":"stale","affected_symbols":["510300"],"source_level":"C","source_type":"sentiment_proxy"}},"p34_data_categories":["sentiment_proxy"]}}`
	svc := NewDataSourceQualityService(qualityReposWithMarket(model.MarketSnapshot{MarketSnapshotID: "market_optional", Symbol: "510300", TradeDate: "2026-06-05", MarketMetricsJSON: metrics}))

	out, err := svc.Run(context.Background(), DataSourceQualityRegressionRequest{Mode: "current", Symbol: "510300"})
	if err != nil {
		t.Fatalf("Run current optional degraded: %v", err)
	}
	if out.Status != "degraded" || len(out.MissingCategories) != 1 || out.MissingCategories[0] != "sentiment_proxy" {
		t.Fatalf("expected legacy missing_categories to include optional degraded category, got %+v", out)
	}
	if out.Policy.Verdict != "waiver_required" || out.Policy.ReleaseGate != "waiver_required" || out.Policy.WaiverCount != 1 || out.Policy.BlockingCount != 0 {
		t.Fatalf("expected optional degraded category to require waiver, got %+v", out.Policy)
	}
}

func TestDataSourceQualityRegressionCurrentUnknownFailureCategoryBlocks(t *testing.T) {
	metrics := `{"source_name":"sentiment_proxy_fixture","source_level":"C","source_type":"sentiment_proxy","metadata":{"p34_source_health":{"sentiment_proxy":{"freshness":"stale","data_date":"2026-06-01","last_failure_at":"2026-06-06T01:00:00Z","failure_category":"vendor_mystery","affected_symbols":["510300"],"source_level":"C","source_type":"sentiment_proxy"}},"p34_data_categories":["sentiment_proxy"]}}`
	svc := NewDataSourceQualityService(qualityReposWithMarket(model.MarketSnapshot{MarketSnapshotID: "market_unknown_failure", Symbol: "510300", TradeDate: "2026-06-05", MarketMetricsJSON: metrics}))

	out, err := svc.Run(context.Background(), DataSourceQualityRegressionRequest{Mode: "current", Symbol: "510300"})
	if err != nil {
		t.Fatalf("Run current unknown failure category: %v", err)
	}
	if out.Policy.Verdict != "blocked" || out.Policy.ReleaseGate != "block" || out.Policy.BlockingCount == 0 {
		t.Fatalf("expected unknown failure category to block policy, got %+v", out.Policy)
	}
	if len(out.Policy.BlockingReasons) == 0 || !strings.Contains(out.Policy.BlockingReasons[0], "failure_category") {
		t.Fatalf("expected blocking reason to mention failure_category, got %+v", out.Policy)
	}
}

func TestDataQualityGateResolutionCheckRequiresResolutionWhenBlocked(t *testing.T) {
	metrics := `{"source_name":"csindex","source_level":"A","source_type":"index_basic","metadata":{"p34_source_health":{"index_valuation_files":{"freshness":"stale","data_date":"2026-06-05","failure_category":"stale","affected_symbols":["000300"],"source_level":"A","source_type":"index_basic"}},"p34_data_categories":["index_valuation_files"]}}`
	resolutions := newMemoryDataQualityGateResolutionRepo()
	svc := NewDataSourceQualityService(qualityReposWithMarketAndResolutions(model.MarketSnapshot{MarketSnapshotID: "market_block", Symbol: "000300", TradeDate: "2026-06-05", MarketMetricsJSON: metrics}, resolutions))

	out, err := svc.CheckGateResolution(context.Background(), DataQualityGateResolutionCheckRequest{Symbol: "000300"})
	if err != nil {
		t.Fatalf("CheckGateResolution blocked: %v", err)
	}
	if out.Policy.Verdict != "blocked" || out.ReleaseClaimState != "requires_resolution" || out.CleanDataClaimAllowed {
		t.Fatalf("expected blocked policy requiring resolution, got %+v", out)
	}
	if out.PolicyFingerprint == "" || len(out.ProhibitedClaims) == 0 {
		t.Fatalf("expected fingerprint and prohibited claims, got %+v", out)
	}
}

func TestDataQualityGateResolutionCreateRejectsBlockedWaiverAndAllowsScopeExclusion(t *testing.T) {
	metrics := `{"source_name":"csindex","source_level":"A","source_type":"index_basic","metadata":{"p34_source_health":{"index_valuation_files":{"freshness":"stale","data_date":"2026-06-05","failure_category":"stale","affected_symbols":["000300"],"source_level":"A","source_type":"index_basic"}},"p34_data_categories":["index_valuation_files"]}}`
	resolutions := newMemoryDataQualityGateResolutionRepo()
	svc := NewDataSourceQualityService(qualityReposWithMarketAndResolutions(model.MarketSnapshot{MarketSnapshotID: "market_block_scope", Symbol: "000300", TradeDate: "2026-06-05", MarketMetricsJSON: metrics}, resolutions))

	_, err := svc.CreateGateResolution(context.Background(), DataQualityGateResolutionCreateRequest{
		Symbol:         "000300",
		ResolutionType: "waiver",
		Scope:          "release scope",
		Reason:         "blocked source health",
		ReleaseImpact:  "do not claim current data clean",
	})
	if !apperr.IsCode(err, apperr.CodeBadRequest) {
		t.Fatalf("expected blocked waiver rejected, got %v", err)
	}

	created, err := svc.CreateGateResolution(context.Background(), DataQualityGateResolutionCreateRequest{
		Symbol:         "000300",
		ResolutionType: "scope_exclusion",
		Scope:          "release excludes current local data health",
		Reason:         "index valuation files stale",
		ReleaseImpact:  "current data health is not claimed clean",
		EvidenceRef:    "docs/release/acceptance/p66",
	})
	if err != nil {
		t.Fatalf("create scope exclusion: %v", err)
	}
	if created.ReleaseClaimState != "resolved_with_scope_exclusion" || created.ActiveResolution == nil || created.ActiveResolution.PolicyFingerprint == "" {
		t.Fatalf("expected scope exclusion resolution, got %+v", created)
	}
	if created.CleanDataClaimAllowed {
		t.Fatalf("scope exclusion must not allow clean data claim: %+v", created)
	}
}

func TestDataQualityGateResolutionWaiverRequiredPolicyAllowsWaiverAndRejectsConflictingActiveType(t *testing.T) {
	metrics := `{"source_name":"sentiment_proxy_fixture","source_level":"C","source_type":"sentiment_proxy","metadata":{"p34_source_health":{"sentiment_proxy":{"freshness":"stale","data_date":"2026-06-01","failure_category":"stale","affected_symbols":["510300"],"source_level":"C","source_type":"sentiment_proxy"}},"p34_data_categories":["sentiment_proxy"]}}`
	resolutions := newMemoryDataQualityGateResolutionRepo()
	svc := NewDataSourceQualityService(qualityReposWithMarketAndResolutions(model.MarketSnapshot{MarketSnapshotID: "market_waiver", Symbol: "510300", TradeDate: "2026-06-05", MarketMetricsJSON: metrics}, resolutions))

	created, err := svc.CreateGateResolution(context.Background(), DataQualityGateResolutionCreateRequest{
		Symbol:         "510300",
		ResolutionType: "waiver",
		Scope:          "optional sentiment proxy excluded from clean claim",
		Reason:         "optional source stale",
		ReleaseImpact:  "waiver recorded, not a clean pass",
	})
	if err != nil {
		t.Fatalf("create waiver: %v", err)
	}
	if created.ReleaseClaimState != "resolved_with_waiver" || created.ActiveResolution == nil {
		t.Fatalf("expected waiver resolution, got %+v", created)
	}

	reused, err := svc.CreateGateResolution(context.Background(), DataQualityGateResolutionCreateRequest{
		Symbol:         "510300",
		ResolutionType: "waiver",
		Scope:          "duplicate scope",
		Reason:         "duplicate reason",
		ReleaseImpact:  "duplicate impact",
	})
	if err != nil {
		t.Fatalf("expected duplicate waiver reuse, got %v", err)
	}
	if reused.ActiveResolution == nil || reused.ActiveResolution.ResolutionID != created.ActiveResolution.ResolutionID {
		t.Fatalf("expected duplicate to reuse active record, got created=%+v reused=%+v", created, reused)
	}

	_, err = svc.CreateGateResolution(context.Background(), DataQualityGateResolutionCreateRequest{
		Symbol:         "510300",
		ResolutionType: "scope_exclusion",
		Scope:          "conflicting scope",
		Reason:         "conflicting reason",
		ReleaseImpact:  "conflicting impact",
	})
	if !apperr.IsCode(err, apperr.CodeConflict) {
		t.Fatalf("expected conflicting active resolution rejected, got %v", err)
	}
}

func TestDataQualityGateResolutionRetireMakesPolicyRequireResolutionAgain(t *testing.T) {
	metrics := `{"source_name":"sentiment_proxy_fixture","source_level":"C","source_type":"sentiment_proxy","metadata":{"p34_source_health":{"sentiment_proxy":{"freshness":"stale","data_date":"2026-06-01","failure_category":"stale","affected_symbols":["510300"],"source_level":"C","source_type":"sentiment_proxy"}},"p34_data_categories":["sentiment_proxy"]}}`
	resolutions := newMemoryDataQualityGateResolutionRepo()
	svc := NewDataSourceQualityService(qualityReposWithMarketAndResolutions(model.MarketSnapshot{MarketSnapshotID: "market_retire", Symbol: "510300", TradeDate: "2026-06-05", MarketMetricsJSON: metrics}, resolutions))
	created, err := svc.CreateGateResolution(context.Background(), DataQualityGateResolutionCreateRequest{Symbol: "510300", ResolutionType: "waiver", Scope: "scope", Reason: "reason", ReleaseImpact: "impact"})
	if err != nil || created.ActiveResolution == nil {
		t.Fatalf("create waiver: out=%+v err=%v", created, err)
	}

	retired, err := svc.RetireGateResolution(context.Background(), created.ActiveResolution.ResolutionID)
	if err != nil {
		t.Fatalf("retire resolution: %v", err)
	}
	if retired.ReleaseClaimState != "requires_resolution" || retired.ActiveResolution != nil {
		t.Fatalf("expected requires resolution after retire, got %+v", retired)
	}
}

func TestDataQualityGateResolutionCreateRollsBackWhenAuditFails(t *testing.T) {
	metrics := `{"source_name":"csindex","source_level":"A","source_type":"index_basic","metadata":{"p34_source_health":{"index_valuation_files":{"freshness":"stale","data_date":"2026-06-05","failure_category":"stale","affected_symbols":["000300"],"source_level":"A","source_type":"index_basic"}},"p34_data_categories":["index_valuation_files"]}}`
	resolutions := newMemoryDataQualityGateResolutionRepo()
	market := qualityMarketRepo{snapshot: model.MarketSnapshot{MarketSnapshotID: "market_audit_fail_create", Symbol: "000300", TradeDate: "2026-06-05", MarketMetricsJSON: metrics}}
	tx := dataQualityGateResolutionMemoryTx{market: market, target: resolutions, audit: failingAuditRepo{}}
	svc := NewDataSourceQualityService(repository.Repositories{MarketRepo: market, DataQualityGateResolutionRepo: resolutions, AuditRepo: failingAuditRepo{}}, tx)

	_, err := svc.CreateGateResolution(context.Background(), DataQualityGateResolutionCreateRequest{
		Symbol:         "000300",
		ResolutionType: "scope_exclusion",
		Scope:          "release excludes current local data health",
		Reason:         "index valuation files stale",
		ReleaseImpact:  "current data health is not claimed clean",
	})
	if err == nil {
		t.Fatalf("expected audit failure")
	}
	items, listErr := resolutions.ListDataQualityGateResolutions(context.Background(), repository.DataQualityGateResolutionFilter{Symbol: "000300"})
	if listErr != nil {
		t.Fatalf("list resolutions: %v", listErr)
	}
	if len(items) != 0 {
		t.Fatalf("audit failure must roll back active resolution, got %+v", items)
	}
}

func TestDataQualityGateResolutionRetireRollsBackWhenAuditFails(t *testing.T) {
	metrics := `{"source_name":"sentiment_proxy_fixture","source_level":"C","source_type":"sentiment_proxy","metadata":{"p34_source_health":{"sentiment_proxy":{"freshness":"stale","data_date":"2026-06-01","failure_category":"stale","affected_symbols":["510300"],"source_level":"C","source_type":"sentiment_proxy"}},"p34_data_categories":["sentiment_proxy"]}}`
	resolutions := newMemoryDataQualityGateResolutionRepo()
	baseSvc := NewDataSourceQualityService(qualityReposWithMarketAndResolutions(model.MarketSnapshot{MarketSnapshotID: "market_audit_fail_retire", Symbol: "510300", TradeDate: "2026-06-05", MarketMetricsJSON: metrics}, resolutions))
	created, err := baseSvc.CreateGateResolution(context.Background(), DataQualityGateResolutionCreateRequest{Symbol: "510300", ResolutionType: "waiver", Scope: "scope", Reason: "reason", ReleaseImpact: "impact"})
	if err != nil || created.ActiveResolution == nil {
		t.Fatalf("create waiver: out=%+v err=%v", created, err)
	}

	market := qualityMarketRepo{snapshot: model.MarketSnapshot{MarketSnapshotID: "market_audit_fail_retire", Symbol: "510300", TradeDate: "2026-06-05", MarketMetricsJSON: metrics}}
	tx := dataQualityGateResolutionMemoryTx{market: market, target: resolutions, audit: failingAuditRepo{}}
	svc := NewDataSourceQualityService(repository.Repositories{MarketRepo: market, DataQualityGateResolutionRepo: resolutions, AuditRepo: failingAuditRepo{}}, tx)

	_, err = svc.RetireGateResolution(context.Background(), created.ActiveResolution.ResolutionID)
	if err == nil {
		t.Fatalf("expected audit failure")
	}
	check, checkErr := baseSvc.CheckGateResolution(context.Background(), DataQualityGateResolutionCheckRequest{Symbol: "510300"})
	if checkErr != nil {
		t.Fatalf("check after failed retire: %v", checkErr)
	}
	if check.ReleaseClaimState != "resolved_with_waiver" || check.ActiveResolution == nil {
		t.Fatalf("audit failure must preserve active resolution, got %+v", check)
	}
}

func TestDataSourceQualityRegressionRejectsUnsupportedMode(t *testing.T) {
	svc := NewDataSourceQualityService(emptyQualityRepos())

	_, err := svc.Run(context.Background(), DataSourceQualityRegressionRequest{Mode: "real"})
	if !apperr.IsCode(err, apperr.CodeBadRequest) {
		t.Fatalf("expected bad request for unsupported mode, got %v", err)
	}
}

func mustJSONText(t *testing.T, v any) string {
	t.Helper()
	buf, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal json: %v", err)
	}
	return string(buf)
}

func emptyQualityRepos() repository.Repositories {
	return repository.Repositories{}
}

func qualityReposWithMarket(snapshot model.MarketSnapshot) repository.Repositories {
	return repository.Repositories{MarketRepo: qualityMarketRepo{snapshot: snapshot}}
}

func qualityReposWithMarketAndResolutions(snapshot model.MarketSnapshot, resolutions repository.DataQualityGateResolutionRepository) repository.Repositories {
	return repository.Repositories{MarketRepo: qualityMarketRepo{snapshot: snapshot}, DataQualityGateResolutionRepo: resolutions}
}

type qualityMarketRepo struct {
	snapshot model.MarketSnapshot
}

func (r qualityMarketRepo) SaveMarketSnapshot(context.Context, model.MarketSnapshot, string) error {
	return nil
}

func (r qualityMarketRepo) GetMarketSnapshot(context.Context, string) (model.MarketSnapshot, error) {
	return r.snapshot, nil
}

func (r qualityMarketRepo) GetLatestMarketSnapshot(context.Context) (model.MarketSnapshot, error) {
	if r.snapshot.MarketSnapshotID == "" {
		return model.MarketSnapshot{}, apperr.New(apperr.CodeNotFound, apperr.CategoryNotFound, "record not found")
	}
	return r.snapshot, nil
}

func (r qualityMarketRepo) GetLatestMarketSnapshotBySymbol(_ context.Context, symbol string) (model.MarketSnapshot, error) {
	if r.snapshot.MarketSnapshotID == "" || (symbol != "" && r.snapshot.Symbol != symbol) {
		return model.MarketSnapshot{}, apperr.New(apperr.CodeNotFound, apperr.CategoryNotFound, "record not found")
	}
	return r.snapshot, nil
}

func (r qualityMarketRepo) MarketSnapshotExists(context.Context, string) (bool, error) {
	return r.snapshot.MarketSnapshotID != "", nil
}

type memoryDataQualityGateResolutionRepo struct {
	items []repository.DataQualityGateResolution
}

func newMemoryDataQualityGateResolutionRepo() *memoryDataQualityGateResolutionRepo {
	return &memoryDataQualityGateResolutionRepo{}
}

func (r *memoryDataQualityGateResolutionRepo) CreateDataQualityGateResolution(_ context.Context, resolution repository.DataQualityGateResolution) error {
	for _, item := range r.items {
		if item.Symbol == resolution.Symbol && item.PolicyFingerprint == resolution.PolicyFingerprint && item.Status == "active" {
			return apperr.New(apperr.CodeConflict, apperr.CategoryConflict, "active resolution already exists")
		}
	}
	r.items = append(r.items, resolution)
	return nil
}

func (r *memoryDataQualityGateResolutionRepo) GetDataQualityGateResolution(_ context.Context, resolutionID string) (repository.DataQualityGateResolution, error) {
	for _, item := range r.items {
		if item.ResolutionID == resolutionID {
			return item, nil
		}
	}
	return repository.DataQualityGateResolution{}, apperr.New(apperr.CodeNotFound, apperr.CategoryNotFound, "record not found")
}

func (r *memoryDataQualityGateResolutionRepo) GetActiveDataQualityGateResolution(_ context.Context, symbol, policyFingerprint string) (repository.DataQualityGateResolution, error) {
	for _, item := range r.items {
		if item.Symbol == symbol && item.PolicyFingerprint == policyFingerprint && item.Status == "active" {
			return item, nil
		}
	}
	return repository.DataQualityGateResolution{}, apperr.New(apperr.CodeNotFound, apperr.CategoryNotFound, "record not found")
}

func (r *memoryDataQualityGateResolutionRepo) ListDataQualityGateResolutions(_ context.Context, filter repository.DataQualityGateResolutionFilter) ([]repository.DataQualityGateResolution, error) {
	var out []repository.DataQualityGateResolution
	for _, item := range r.items {
		if filter.Symbol != "" && item.Symbol != filter.Symbol {
			continue
		}
		if filter.Status != "" && item.Status != filter.Status {
			continue
		}
		out = append(out, item)
	}
	return out, nil
}

func (r *memoryDataQualityGateResolutionRepo) RetireDataQualityGateResolution(_ context.Context, resolutionID, retiredBy, retiredAt string) error {
	for i, item := range r.items {
		if item.ResolutionID == resolutionID && item.Status == "active" {
			r.items[i].Status = "retired"
			r.items[i].RetiredBy = retiredBy
			r.items[i].RetiredAt = retiredAt
			return nil
		}
	}
	return apperr.New(apperr.CodeNotFound, apperr.CategoryNotFound, "record not found")
}

type failingAuditRepo struct{}

func (failingAuditRepo) AppendAuditEvent(context.Context, repository.AuditEvent) error {
	return errors.New("audit write failed")
}

func (failingAuditRepo) GetAuditEvent(context.Context, string) (repository.AuditEvent, error) {
	return repository.AuditEvent{}, apperr.New(apperr.CodeNotFound, apperr.CategoryNotFound, "record not found")
}

func (failingAuditRepo) ListAuditEvents(context.Context) ([]repository.AuditEvent, error) {
	return nil, nil
}

type dataQualityGateResolutionMemoryTx struct {
	market repository.MarketRepository
	target *memoryDataQualityGateResolutionRepo
	audit  repository.AuditRepository
}

func (tx dataQualityGateResolutionMemoryTx) WithinTx(ctx context.Context, fn func(context.Context, repository.Repositories) error) error {
	txRepo := &memoryDataQualityGateResolutionRepo{items: append([]repository.DataQualityGateResolution(nil), tx.target.items...)}
	err := fn(ctx, repository.Repositories{
		MarketRepo:                    tx.market,
		AuditRepo:                     tx.audit,
		DataQualityGateResolutionRepo: txRepo,
	})
	if err != nil {
		return err
	}
	tx.target.items = txRepo.items
	return nil
}
