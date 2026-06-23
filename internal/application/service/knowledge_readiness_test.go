package service

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	appsqlite "investment-agent/internal/infrastructure/persistence/sqlite"
)

func TestBuiltInKnowledgeRegistryHasStablePrinciplesAndSafetyBoundaries(t *testing.T) {
	registry := BuiltInKnowledgeRegistry()
	entries := registry.Entries()

	byID := map[string]KnowledgeEntry{}
	for _, entry := range entries {
		byID[entry.KnowledgeID] = entry
	}
	for _, id := range []string{
		"master.graham.margin_of_safety",
		"master.buffett.circle_of_competence",
		"master.livermore.trend_discipline",
		"master.dalio.risk_parity_cycle",
		"master.marks.second_level_thinking",
		"master.lynch.know_what_you_own",
		"master.templeton.extreme_pessimism",
		"discipline.no_single_source_decision",
		"risk_sop.evidence_insufficient",
	} {
		entry, ok := byID[id]
		if !ok {
			t.Fatalf("expected stable built-in knowledge id %q in registry, got ids=%v", id, knowledgeEntryIDs(entries))
		}
		if entry.Category == "" || entry.Summary == "" || len(entry.AppliesTo) == 0 || len(entry.RuleMapping) == 0 || entry.SafetyBoundary == "" {
			t.Fatalf("knowledge entry %s must be fully described, got %+v", id, entry)
		}
		if strings.HasPrefix(id, "master.") && entry.FormalEvidenceAllowed {
			t.Fatalf("master principle %s must not be formal market evidence: %+v", id, entry)
		}
	}
	if !byID["master.graham.margin_of_safety"].LLMContextAllowed {
		t.Fatalf("expected Graham principle to be allowed as summarized LLM context")
	}
	if byID["discipline.no_single_source_decision"].FormalEvidenceAllowed {
		t.Fatalf("discipline rules explain verdicts but must not become external formal evidence")
	}
}

func TestKnowledgeReadinessServiceReportsReadyFor510300WithFreshDependencies(t *testing.T) {
	ctx := context.Background()
	repos, db := knowledgeReadinessRepos(t)
	seedKnowledgeReadinessFacts(t, ctx, repos, model.MarketSnapshot{
		MarketSnapshotID:  "market_ready",
		Symbol:            "510300",
		TradeDate:         "2026-06-19",
		DataStatus:        "fresh",
		ClosePrice:        4.75,
		PEPercentile:      28,
		PBPercentile:      35,
		LiquidityState:    model.LiquidityNormal,
		SentimentState:    model.SentimentNeutral,
		MarketMetricsJSON: readinessHealthJSON("fresh", "fresh", "fresh", "fresh", "fresh", "fresh", "fresh"),
	}, repository.SourceVerification{
		VerificationID:                  "verify_ready",
		VerificationGroupID:             "group_ready",
		EventID:                         "event_ready",
		Symbol:                          "510300",
		EventType:                       "normal",
		EvidenceRole:                    "formal",
		VerificationStatus:              "satisfied",
		IndependentSourceCount:          3,
		HighGradeIndependentSourceCount: 2,
		HighestSourceLevel:              "A",
		CreatedAt:                       "2026-06-19T08:00:00Z",
	})
	svc := NewKnowledgeReadinessService(repos)

	out, err := svc.Evaluate(ctx, KnowledgeReadinessRequest{Symbol: "510300"})
	if err != nil {
		t.Fatalf("Evaluate readiness: %v", err)
	}
	if out.Status != "ready" || out.SymbolProfile.Symbol != "510300" || out.SymbolProfile.TrackedIndexSymbol != "000300" {
		t.Fatalf("expected ready 510300 profile, got %+v", out)
	}
	if len(out.KnowledgeReferences) < 5 || !strings.Contains(out.LLMContextSummary, "master.graham.margin_of_safety") {
		t.Fatalf("expected master/discipline references in LLM context, got %+v summary=%q", out.KnowledgeReferences, out.LLMContextSummary)
	}
	for _, dep := range out.DataDependencies {
		if dep.Required && dep.Status != "ready" {
			t.Fatalf("expected required dependency %s ready, got %+v", dep.Category, dep)
		}
	}
	if strings.Contains(mustJSONForReadiness(t, out), "sk-") || strings.Contains(mustJSONForReadiness(t, out), "raw HTTP") {
		t.Fatalf("readiness output leaked sensitive diagnostic: %s", mustJSONForReadiness(t, out))
	}
	assertKnowledgeReadinessNoAuditWrites(t, db)
}

func TestKnowledgeReadinessServiceReportsReadyForKnownNon510300Symbol(t *testing.T) {
	ctx := context.Background()
	repos, db := knowledgeReadinessRepos(t)
	seedKnowledgeReadinessFacts(t, ctx, repos, model.MarketSnapshot{
		MarketSnapshotID:  "market_159915_ready",
		Symbol:            "159915",
		TradeDate:         "2026-06-19",
		DataStatus:        "fresh",
		ClosePrice:        2.41,
		PEPercentile:      42,
		PBPercentile:      47,
		LiquidityState:    model.LiquidityNormal,
		SentimentState:    model.SentimentNeutral,
		MarketMetricsJSON: readinessHealthJSONForSymbolsWithRAG("159915", "399006", "fresh", "fresh", "fresh", "fresh", "fresh", "fresh", "fresh", "fresh", "req_159915_ready"),
	}, repository.SourceVerification{
		VerificationID:                  "verify_159915_ready",
		VerificationGroupID:             "group_159915_ready",
		EventID:                         "event_159915_ready",
		Symbol:                          "159915",
		EventType:                       "normal",
		EvidenceRole:                    "formal",
		VerificationStatus:              "satisfied",
		IndependentSourceCount:          3,
		HighGradeIndependentSourceCount: 2,
		HighestSourceLevel:              "A",
		CreatedAt:                       "2026-06-19T08:00:00Z",
	})
	svc := NewKnowledgeReadinessService(repos)

	out, err := svc.Evaluate(ctx, KnowledgeReadinessRequest{Symbol: "159915"})
	if err != nil {
		t.Fatalf("Evaluate readiness: %v", err)
	}
	if out.Status != "ready" {
		t.Fatalf("expected ready 159915 profile, got %+v", out)
	}
	if out.SymbolProfile.Symbol != "159915" || out.SymbolProfile.TrackedIndexSymbol != "399006" || !out.SymbolProfile.Known {
		t.Fatalf("expected known 159915 -> 399006 profile, got %+v", out.SymbolProfile)
	}
	if !strings.Contains(out.LLMContextSummary, "symbol_profile.159915") {
		t.Fatalf("LLM context summary must include the concrete known symbol profile, got %q", out.LLMContextSummary)
	}
	trackedIndex := readinessDependencyByCategory(out.DataDependencies, "tracked_index")
	if trackedIndex.Status != "ready" || trackedIndex.SourceLevel != "A" || trackedIndex.RequestID != "req_159915_ready" || trackedIndex.DataDate != "2026-06-19" {
		t.Fatalf("expected tracked index dependency ready for 399006, got %+v", trackedIndex)
	}
	ragIndex := readinessDependencyByCategory(out.DataDependencies, "rag_index")
	if ragIndex.Status != "ready" || ragIndex.RequestID != "req_159915_ready" || !knowledgeReadinessContainsString(ragIndex.AffectedSymbols, "159915") {
		t.Fatalf("expected dynamic rag_index readiness bound to 159915, got %+v", ragIndex)
	}
	assertKnowledgeReadinessNoAuditWrites(t, db)
}

func TestKnowledgeReadinessServiceDegradesMissingValuationWithoutFabricatingPass(t *testing.T) {
	ctx := context.Background()
	repos, _ := knowledgeReadinessRepos(t)
	seedKnowledgeReadinessFacts(t, ctx, repos, model.MarketSnapshot{
		MarketSnapshotID:  "market_degraded",
		Symbol:            "510300",
		TradeDate:         "2026-06-19",
		DataStatus:        "fresh",
		ClosePrice:        4.75,
		MarketMetricsJSON: readinessHealthJSON("fresh", "fresh", "fresh", "fresh", "parse_error", "fresh", "fresh"),
	}, repository.SourceVerification{
		VerificationID:                  "verify_degraded",
		VerificationGroupID:             "group_degraded",
		EventID:                         "event_degraded",
		Symbol:                          "510300",
		EventType:                       "normal",
		EvidenceRole:                    "formal",
		VerificationStatus:              "satisfied",
		IndependentSourceCount:          2,
		HighGradeIndependentSourceCount: 2,
		HighestSourceLevel:              "A",
		CreatedAt:                       "2026-06-19T08:00:00Z",
	})
	if err := repos.RuleRepo.ArchiveActiveRuleVersions(ctx); err != nil {
		t.Fatalf("archive active rules: %v", err)
	}
	svc := NewKnowledgeReadinessService(repos)

	out, err := svc.Evaluate(ctx, KnowledgeReadinessRequest{Symbol: "510300"})
	if err != nil {
		t.Fatalf("Evaluate readiness: %v", err)
	}
	if out.Status != "degraded" {
		t.Fatalf("expected degraded readiness, got %+v", out)
	}
	valuation := readinessDependencyByCategory(out.DataDependencies, "valuation_percentiles")
	if valuation.Status != "degraded" || !strings.Contains(strings.Join(valuation.AffectedFeatures, ","), "margin_of_safety") {
		t.Fatalf("expected valuation dependency degraded with feature impact, got %+v", valuation)
	}
	if !strings.Contains(out.LLMContextSummary, "valuation_percentiles=degraded") {
		t.Fatalf("LLM context summary must include degraded valuation, got %q", out.LLMContextSummary)
	}
}

func TestKnowledgeReadinessServicePropagatesCriticalDataGapsToFeatureImpacts(t *testing.T) {
	ctx := context.Background()
	repos, _ := knowledgeReadinessRepos(t)
	seedKnowledgeReadinessFacts(t, ctx, repos, model.MarketSnapshot{
		MarketSnapshotID:  "market_critical_gaps",
		Symbol:            "510300",
		TradeDate:         "2026-06-19",
		DataStatus:        "fresh",
		ClosePrice:        4.75,
		MarketMetricsJSON: readinessHealthJSON("fresh", "fresh", "fresh", "fresh", "parse_error", "missing", "fresh"),
	}, repository.SourceVerification{
		VerificationID:                  "verify_single_source",
		VerificationGroupID:             "group_single_source",
		EventID:                         "event_single_source",
		Symbol:                          "510300",
		EventType:                       "normal",
		EvidenceRole:                    "formal",
		VerificationStatus:              "satisfied",
		IndependentSourceCount:          1,
		HighGradeIndependentSourceCount: 1,
		HighestSourceLevel:              "A",
		CreatedAt:                       "2026-06-19T08:00:00Z",
	})
	svc := NewKnowledgeReadinessService(repos)

	out, err := svc.Evaluate(ctx, KnowledgeReadinessRequest{Symbol: "510300"})
	if err != nil {
		t.Fatalf("Evaluate readiness: %v", err)
	}
	if out.Status != "degraded" {
		t.Fatalf("expected degraded readiness, got %+v", out)
	}
	for _, category := range []string{"valuation_percentiles", "liquidity", "formal_evidence"} {
		dep := readinessDependencyByCategory(out.DataDependencies, category)
		if dep.Status != "degraded" {
			t.Fatalf("expected %s degraded, got %+v", category, dep)
		}
	}
	assertFeatureImpactForReadiness(t, out.FeatureImpacts, "valuation_percentiles", "margin_of_safety", "不得声明安全边际")
	assertFeatureImpactForReadiness(t, out.FeatureImpacts, "valuation_percentiles", "expected_return", "预期收益精度不足")
	assertFeatureImpactForReadiness(t, out.FeatureImpacts, "liquidity", "risk_alerts", "不得输出大额或市价式行动建议")
	assertFeatureImpactForReadiness(t, out.FeatureImpacts, "formal_evidence", "consultation", "不生成交易确认")
	for _, impact := range out.FeatureImpacts {
		if impact.Category == "valuation_percentiles" || impact.Category == "liquidity" || impact.Category == "formal_evidence" {
			if !knowledgeReadinessContainsString(impact.Claims, "不得伪造成 ready") || !knowledgeReadinessContainsString(impact.Claims, "不得输出交易确认") {
				t.Fatalf("impact must carry safety claims, got %+v", impact)
			}
		}
	}
}

func TestKnowledgeReadinessServiceDoesNotSubstituteStubBackgroundOrLLMForRequiredData(t *testing.T) {
	ctx := context.Background()
	repos, _ := knowledgeReadinessRepos(t)
	seedKnowledgeReadinessFacts(t, ctx, repos, model.MarketSnapshot{
		MarketSnapshotID:  "market_stubbed_required",
		Symbol:            "510300",
		TradeDate:         "2026-06-19",
		DataStatus:        "fresh",
		ClosePrice:        4.75,
		MarketMetricsJSON: readinessHealthJSON("fresh", "fresh", "fresh", "stubbed", "stubbed", "stubbed", "stubbed"),
	}, repository.SourceVerification{
		VerificationID:                  "verify_background_only",
		VerificationGroupID:             "group_background_only",
		EventID:                         "event_background_only",
		Symbol:                          "510300",
		EventType:                       "normal",
		EvidenceRole:                    "background",
		VerificationStatus:              "background_only",
		IndependentSourceCount:          1,
		HighGradeIndependentSourceCount: 0,
		HighestSourceLevel:              "C",
		CreatedAt:                       "2026-06-19T08:00:00Z",
	})
	svc := NewKnowledgeReadinessService(repos)

	out, err := svc.Evaluate(ctx, KnowledgeReadinessRequest{Symbol: "510300"})
	if err != nil {
		t.Fatalf("Evaluate readiness: %v", err)
	}
	if out.Status != "degraded" {
		t.Fatalf("stub/background/LLM context must not substitute required data, got %+v", out)
	}
	for _, category := range []string{"market_price", "valuation_percentiles", "liquidity", "formal_evidence"} {
		dep := readinessDependencyByCategory(out.DataDependencies, category)
		if dep.Status != "degraded" {
			t.Fatalf("%s must remain degraded, got %+v", category, dep)
		}
	}
	if readinessDependencyByCategory(out.DataDependencies, "llm_context").Status != "ready" {
		t.Fatalf("LLM context should remain available as context only: %+v", out.DataDependencies)
	}
	for _, token := range []string{"market_price=degraded", "valuation_percentiles=degraded", "liquidity=degraded", "formal_evidence=degraded", "llm_context=ready"} {
		if !strings.Contains(out.LLMContextSummary, token) {
			t.Fatalf("expected %s in readiness LLM summary, got %q", token, out.LLMContextSummary)
		}
	}
	assertFeatureImpactForReadiness(t, out.FeatureImpacts, "formal_evidence", "consultation", "不生成交易确认")
}

func TestKnowledgeReadinessServiceDegradesWhenActiveRuleIsMissing(t *testing.T) {
	ctx := context.Background()
	repos, db := knowledgeReadinessRepos(t)
	seedKnowledgeReadinessFacts(t, ctx, repos, model.MarketSnapshot{
		MarketSnapshotID:  "market_missing_rule",
		Symbol:            "510300",
		TradeDate:         "2026-06-19",
		DataStatus:        "fresh",
		ClosePrice:        4.75,
		PEPercentile:      28,
		PBPercentile:      35,
		LiquidityState:    model.LiquidityNormal,
		SentimentState:    model.SentimentNeutral,
		MarketMetricsJSON: readinessHealthJSON("fresh", "fresh", "fresh", "fresh", "fresh", "fresh", "fresh"),
	}, repository.SourceVerification{
		VerificationID:                  "verify_missing_rule",
		VerificationGroupID:             "group_missing_rule",
		EventID:                         "event_missing_rule",
		Symbol:                          "510300",
		EventType:                       "normal",
		EvidenceRole:                    "formal",
		VerificationStatus:              "satisfied",
		IndependentSourceCount:          3,
		HighGradeIndependentSourceCount: 2,
		HighestSourceLevel:              "A",
		CreatedAt:                       "2026-06-19T08:00:00Z",
	})
	if _, err := db.Exec(`DELETE FROM rule_versions WHERE status='active'`); err != nil {
		t.Fatalf("delete active rules: %v", err)
	}
	svc := NewKnowledgeReadinessService(repos)

	out, err := svc.Evaluate(ctx, KnowledgeReadinessRequest{Symbol: "510300"})
	if err != nil {
		t.Fatalf("Evaluate readiness: %v", err)
	}
	if out.Status != "degraded" {
		t.Fatalf("missing active rule must degrade readiness, got %+v", out)
	}
	activeRule := readinessDependencyByCategory(out.DataDependencies, "active_rule")
	if activeRule.Status != "degraded" || activeRule.Freshness != "missing" || !strings.Contains(activeRule.SafeDegradation, "规则裁决边界") {
		t.Fatalf("expected active rule dependency to explain missing rule boundary, got %+v", activeRule)
	}
	if !strings.Contains(out.LLMContextSummary, "active_rule=degraded") {
		t.Fatalf("LLM context summary must include active rule readiness, got %q", out.LLMContextSummary)
	}
}

func TestKnowledgeReadinessServiceBlocksUnknownSymbolProfile(t *testing.T) {
	ctx := context.Background()
	repos, _ := knowledgeReadinessRepos(t)
	svc := NewKnowledgeReadinessService(repos)

	out, err := svc.Evaluate(ctx, KnowledgeReadinessRequest{Symbol: "999999"})
	if err != nil {
		t.Fatalf("Evaluate readiness: %v", err)
	}
	if out.Status != "blocked" {
		t.Fatalf("unknown symbol profile must block readiness, got %+v", out)
	}
	profile := readinessDependencyByCategory(out.DataDependencies, "symbol_profile")
	if profile.Status != "blocked" || !strings.Contains(profile.SafeDegradation, "不生成正式交易类建议") {
		t.Fatalf("expected blocked symbol profile safe degradation, got %+v", profile)
	}
	if out.SymbolProfile.Symbol != "999999" || out.SymbolProfile.Known {
		t.Fatalf("expected explicit unknown profile for requested symbol, got %+v", out.SymbolProfile)
	}
}

func knowledgeReadinessRepos(t *testing.T) (repository.Repositories, *sql.DB) {
	t.Helper()
	store, err := appsqlite.Open(":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	if err := appsqlite.Migrate(context.Background(), store.DB); err != nil {
		t.Fatalf("migrate sqlite: %v", err)
	}
	repos := repository.Repositories{
		AuditRepo:        appsqlite.NewAuditRepository(store.DB),
		RuleRepo:         appsqlite.NewRuleRepository(store.DB),
		MarketRepo:       appsqlite.NewMarketRepository(store.DB),
		IntelligenceRepo: appsqlite.NewIntelligenceRepository(store.DB),
	}
	return repos, store.DB
}

func seedKnowledgeReadinessFacts(t *testing.T, ctx context.Context, repos repository.Repositories, market model.MarketSnapshot, verification repository.SourceVerification) {
	t.Helper()
	if market.LiquidityState == "" {
		market.LiquidityState = model.LiquidityNormal
	}
	if market.SentimentState == "" {
		market.SentimentState = model.SentimentNeutral
	}
	if err := repos.MarketRepo.SaveMarketSnapshot(ctx, market, "2026-06-19T08:00:00Z"); err != nil {
		t.Fatalf("seed market: %v", err)
	}
	if err := repos.IntelligenceRepo.SaveSourceVerification(ctx, verification); err != nil {
		t.Fatalf("seed source verification: %v", err)
	}
}

func readinessHealthJSON(symbolProfile, fundProfile, trackedIndex, marketPrice, valuation, liquidity, sentiment string) string {
	return readinessHealthJSONForSymbols("510300", "000300", symbolProfile, fundProfile, trackedIndex, marketPrice, valuation, liquidity, sentiment)
}

func readinessHealthJSONForSymbols(symbol, trackedIndexSymbol, symbolProfile, fundProfile, trackedIndex, marketPrice, valuation, liquidity, sentiment string) string {
	return readinessHealthJSONForSymbolsWithRAG(symbol, trackedIndexSymbol, symbolProfile, fundProfile, trackedIndex, marketPrice, valuation, liquidity, sentiment, "", "")
}

func readinessHealthJSONForSymbolsWithRAG(symbol, trackedIndexSymbol, symbolProfile, fundProfile, trackedIndex, marketPrice, valuation, liquidity, sentiment, ragIndex, requestID string) string {
	ragJSON := ""
	ragCategory := ""
	if strings.TrimSpace(ragIndex) != "" {
		ragJSON = `,"rag_index":{"freshness":"` + ragIndex + `","data_date":"2026-06-19","affected_symbols":["` + symbol + `"],"source_level":"local_index","source_name":"veclite","source_type":"rag_index","request_id":"` + requestID + `"}`
		ragCategory = `,"rag_index"`
	}
	return `{"source_name":"p74_fixture","source_level":"A","source_type":"readiness_fixture","captured_at":"2026-06-19T08:00:00Z","metadata":{"p34_source_health":{"symbol_profile":{"freshness":"` + symbolProfile + `","data_date":"2026-06-19","affected_symbols":["` + symbol + `"],"source_level":"A","request_id":"` + requestID + `"},"fund_profile":{"freshness":"` + fundProfile + `","data_date":"2026-06-19","affected_symbols":["` + symbol + `"],"source_level":"B","request_id":"` + requestID + `"},"tracked_index":{"freshness":"` + trackedIndex + `","data_date":"2026-06-19","affected_symbols":["` + trackedIndexSymbol + `"],"source_level":"A","request_id":"` + requestID + `"},"market_price":{"freshness":"` + marketPrice + `","data_date":"2026-06-19","affected_symbols":["` + symbol + `"],"source_level":"B","request_id":"` + requestID + `"},"valuation_percentiles":{"freshness":"` + valuation + `","data_date":"2026-06-19","affected_symbols":["` + trackedIndexSymbol + `"],"source_level":"A","failure_category":"` + valuation + `","request_id":"` + requestID + `"},"liquidity":{"freshness":"` + liquidity + `","data_date":"2026-06-19","affected_symbols":["` + symbol + `"],"source_level":"B","request_id":"` + requestID + `"},"sentiment_proxy":{"freshness":"` + sentiment + `","data_date":"2026-06-19","affected_symbols":["` + symbol + `"],"source_level":"C","request_id":"` + requestID + `"}` + ragJSON + `},"p34_data_categories":["symbol_profile","fund_profile","tracked_index","market_price","valuation_percentiles","liquidity","sentiment_proxy"` + ragCategory + `]}}`
}

func knowledgeEntryIDs(entries []KnowledgeEntry) []string {
	ids := make([]string, 0, len(entries))
	for _, entry := range entries {
		ids = append(ids, entry.KnowledgeID)
	}
	return ids
}

func readinessDependencyByCategory(items []KnowledgeDataDependency, category string) KnowledgeDataDependency {
	for _, item := range items {
		if item.Category == category {
			return item
		}
	}
	return KnowledgeDataDependency{}
}

func mustJSONForReadiness(t *testing.T, value any) string {
	t.Helper()
	return mustJSONText(t, value)
}

func assertKnowledgeReadinessNoAuditWrites(t *testing.T, db *sql.DB) {
	t.Helper()
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM audit_events`).Scan(&count); err != nil {
		t.Fatalf("count audit_events: %v", err)
	}
	if count != 0 {
		t.Fatalf("readiness evaluation must be read-only, audit_events count=%d", count)
	}
}

func assertFeatureImpactForReadiness(t *testing.T, impacts []KnowledgeFeatureImpact, category string, feature string, impactText string) {
	t.Helper()
	for _, impact := range impacts {
		if impact.Category == category && impact.Feature == feature && strings.Contains(impact.Impact, impactText) {
			return
		}
	}
	t.Fatalf("missing feature impact category=%s feature=%s text=%q in %+v", category, feature, impactText, impacts)
}
