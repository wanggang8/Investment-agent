package service

import (
	"context"
	"strings"

	appknowledge "investment-agent/internal/application/knowledge"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

const (
	KnowledgeReadinessReady    = appknowledge.ReadinessReady
	KnowledgeReadinessDegraded = appknowledge.ReadinessDegraded
	KnowledgeReadinessBlocked  = appknowledge.ReadinessBlocked
)

type KnowledgeEntry = appknowledge.Entry

type KnowledgeRegistry = appknowledge.Registry

// BuiltInKnowledgeRegistry returns the P74 deterministic built-in registry.
func BuiltInKnowledgeRegistry() KnowledgeRegistry {
	return appknowledge.BuiltInRegistry()
}

type KnowledgeReadinessRequest struct {
	Symbol string
}

type KnowledgeReadinessResponse struct {
	Symbol              string                    `json:"symbol"`
	Status              string                    `json:"status"`
	SymbolProfile       KnowledgeSymbolProfile    `json:"symbol_profile"`
	KnowledgeReferences []KnowledgeEntry          `json:"knowledge_references"`
	DataDependencies    []KnowledgeDataDependency `json:"data_dependencies"`
	FeatureImpacts      []KnowledgeFeatureImpact  `json:"feature_impacts"`
	LLMContextSummary   string                    `json:"llm_context_summary"`
	SafetyNotes         []string                  `json:"safety_notes"`
}

type KnowledgeSymbolProfile struct {
	Symbol             string `json:"symbol"`
	Name               string `json:"name,omitempty"`
	AssetType          string `json:"asset_type,omitempty"`
	TrackedIndexSymbol string `json:"tracked_index_symbol,omitempty"`
	TrackedIndexName   string `json:"tracked_index_name,omitempty"`
	Known              bool   `json:"known"`
}

type KnowledgeDataDependency struct {
	Category         string   `json:"category"`
	Status           string   `json:"status"`
	Required         bool     `json:"required"`
	SourceName       string   `json:"source_name,omitempty"`
	SourceLevel      string   `json:"source_level,omitempty"`
	SourceType       string   `json:"source_type,omitempty"`
	Freshness        string   `json:"freshness,omitempty"`
	DataDate         string   `json:"data_date,omitempty"`
	RequestID        string   `json:"request_id,omitempty"`
	AffectedSymbols  []string `json:"affected_symbols,omitempty"`
	AffectedFeatures []string `json:"affected_features,omitempty"`
	SafeDegradation  string   `json:"safe_degradation,omitempty"`
}

type KnowledgeFeatureImpact struct {
	Feature  string   `json:"feature"`
	Category string   `json:"category"`
	Impact   string   `json:"impact"`
	Claims   []string `json:"claims,omitempty"`
}

type KnowledgeReadinessService struct {
	repos    repository.Repositories
	registry KnowledgeRegistry
}

func NewKnowledgeReadinessService(repos repository.Repositories) *KnowledgeReadinessService {
	return &KnowledgeReadinessService{repos: repos, registry: BuiltInKnowledgeRegistry()}
}

func (s *KnowledgeReadinessService) Evaluate(ctx context.Context, req KnowledgeReadinessRequest) (KnowledgeReadinessResponse, error) {
	symbol := strings.TrimSpace(req.Symbol)
	if symbol == "" {
		symbol = "510300"
	}
	profile := builtInSymbolProfile(symbol)
	out := KnowledgeReadinessResponse{
		Symbol:              symbol,
		Status:              KnowledgeReadinessReady,
		SymbolProfile:       profile,
		KnowledgeReferences: s.matchedKnowledge(profile),
		SafetyNotes: []string{
			"内置知识只作为纪律、规则映射和 LLM 分析上下文，不作为正式市场证据。",
			"准备度检查只读取本地事实，不刷新数据、不修改规则、不改变账户或确认记录。",
		},
	}

	health := map[string]sourceHealthReadiness{}
	if s.repos.MarketRepo != nil {
		market, err := s.repos.MarketRepo.GetLatestMarketSnapshotBySymbol(ctx, symbol)
		if err != nil && !apperr.IsCode(err, apperr.CodeNotFound) {
			return KnowledgeReadinessResponse{}, err
		}
		if err == nil {
			for _, item := range SourceHealthFromMarketSnapshot(market) {
				health[item.DataCategory] = sourceHealthReadiness{
					freshness:       item.Freshness,
					sourceName:      item.SourceName,
					sourceLevel:     item.SourceLevel,
					sourceType:      item.SourceType,
					dataDate:        item.DataDate,
					requestID:       item.RequestID,
					affectedSymbols: item.AffectedSymbols,
				}
			}
		}
	}

	out.DataDependencies = append(out.DataDependencies, symbolProfileDependency(profile, health["symbol_profile"]))
	for _, category := range []string{"fund_profile", "tracked_index", "market_price", "valuation_percentiles", "liquidity", "sentiment_proxy"} {
		out.DataDependencies = append(out.DataDependencies, dataDependencyFromHealth(category, health[category]))
	}
	out.DataDependencies = append(out.DataDependencies, s.activeRuleDependency(ctx))
	out.DataDependencies = append(out.DataDependencies, s.formalEvidenceDependency(ctx, symbol))
	out.DataDependencies = append(out.DataDependencies, ragIndexDependencyFromHealth(health["rag_index"]))
	out.DataDependencies = append(out.DataDependencies, KnowledgeDataDependency{Category: "llm_context", Status: KnowledgeReadinessReady, Required: false, SafeDegradation: "LLM 只生成分析材料，不能覆盖规则最终裁决。", AffectedFeatures: []string{"consultation", "decision_detail"}})

	out.Status = aggregateKnowledgeReadinessStatus(out.DataDependencies)
	out.FeatureImpacts = featureImpactsFromDependencies(out.DataDependencies)
	out.LLMContextSummary = BuildKnowledgeLLMContextSummary(out.KnowledgeReferences, out.DataDependencies)
	return out, nil
}

func (s *KnowledgeReadinessService) matchedKnowledge(profile KnowledgeSymbolProfile) []KnowledgeEntry {
	return s.registry.EntriesForSymbol(profile.Symbol)
}

func (s *KnowledgeReadinessService) activeRuleDependency(ctx context.Context) KnowledgeDataDependency {
	dep := KnowledgeDataDependency{Category: "active_rule", Status: KnowledgeReadinessDegraded, Required: true, Freshness: "missing", SafeDegradation: "当前 active rule 缺失时规则裁决边界不可确认，只能展示数据准备度降级，不生成交易确认。", AffectedFeatures: []string{"consultation", "decision_detail", "rules"}}
	if s.repos.RuleRepo == nil {
		return dep
	}
	ruleVersion, err := s.repos.RuleRepo.GetActiveRuleVersion(ctx)
	if err != nil {
		return dep
	}
	if strings.TrimSpace(ruleVersion.RuleVersion) == "" || ruleVersion.Status != "active" {
		return dep
	}
	dep.Status = KnowledgeReadinessReady
	dep.Freshness = "fresh"
	dep.SourceLevel = "local_rule"
	dep.SafeDegradation = "规则版本只作为本地裁决边界，不代表市场事实或未来收益。"
	return dep
}

func (s *KnowledgeReadinessService) formalEvidenceDependency(ctx context.Context, symbol string) KnowledgeDataDependency {
	dep := KnowledgeDataDependency{Category: "formal_evidence", Status: KnowledgeReadinessDegraded, Required: true, SafeDegradation: "正式证据不足时进入冻结观察或信息不足，不生成交易确认。", AffectedFeatures: []string{"consultation", "decision_detail", "risk_alerts"}}
	if s.repos.IntelligenceRepo == nil {
		dep.Freshness = "missing"
		return dep
	}
	verification, err := s.repos.IntelligenceRepo.GetLatestSourceVerificationByFilter(ctx, symbol, "")
	if err != nil {
		dep.Freshness = "missing"
		return dep
	}
	dep.SourceLevel = verification.HighestSourceLevel
	switch {
	case verification.VerificationStatus == string(model.VerificationSatisfied) && verification.HighGradeIndependentSourceCount >= 2:
		dep.Status = KnowledgeReadinessReady
		dep.Freshness = "fresh"
	case verification.VerificationStatus == string(model.VerificationBackgroundOnly):
		dep.Status = KnowledgeReadinessDegraded
		dep.Freshness = "background_only"
	default:
		dep.Status = KnowledgeReadinessDegraded
		dep.Freshness = "insufficient"
	}
	return dep
}

type sourceHealthReadiness struct {
	freshness       string
	sourceName      string
	sourceLevel     string
	sourceType      string
	dataDate        string
	requestID       string
	affectedSymbols []string
}

func builtInSymbolProfile(symbol string) KnowledgeSymbolProfile {
	profile, ok := appknowledge.LookupSymbolProfile(symbol)
	if !ok {
		return KnowledgeSymbolProfile{Symbol: strings.TrimSpace(symbol), Known: false}
	}
	return KnowledgeSymbolProfile{
		Symbol:             profile.Symbol,
		Name:               profile.Name,
		AssetType:          profile.AssetType,
		TrackedIndexSymbol: profile.TrackedIndexSymbol,
		TrackedIndexName:   profile.TrackedIndexName,
		Known:              true,
	}
}

func symbolProfileDependency(profile KnowledgeSymbolProfile, health sourceHealthReadiness) KnowledgeDataDependency {
	if !profile.Known {
		return KnowledgeDataDependency{Category: "symbol_profile", Status: KnowledgeReadinessBlocked, Required: true, Freshness: "missing", SourceName: health.sourceName, SourceLevel: health.sourceLevel, SourceType: health.sourceType, DataDate: health.dataDate, RequestID: health.requestID, AffectedSymbols: health.affectedSymbols, AffectedFeatures: []string{"consultation", "daily_discipline", "data_quality"}, SafeDegradation: "标的画像未知时不生成正式交易类建议，只提示补齐标的画像或能力圈。"}
	}
	return KnowledgeDataDependency{Category: "symbol_profile", Status: KnowledgeReadinessReady, Required: true, Freshness: knowledgeReadinessFirstNonEmpty(health.freshness, "fresh"), SourceName: health.sourceName, SourceLevel: health.sourceLevel, SourceType: health.sourceType, DataDate: health.dataDate, RequestID: health.requestID, AffectedSymbols: health.affectedSymbols, AffectedFeatures: []string{"consultation", "daily_discipline", "data_quality"}}
}

func dataDependencyFromHealth(category string, health sourceHealthReadiness) KnowledgeDataDependency {
	dep := KnowledgeDataDependency{Category: category, Required: requiredReadinessCategory(category), Freshness: knowledgeReadinessFirstNonEmpty(health.freshness, "missing"), SourceName: health.sourceName, SourceLevel: health.sourceLevel, SourceType: health.sourceType, DataDate: health.dataDate, RequestID: health.requestID, AffectedSymbols: health.affectedSymbols, AffectedFeatures: affectedFeaturesForReadinessCategory(category), SafeDegradation: safeDegradationForReadinessCategory(category)}
	switch dep.Freshness {
	case "fresh":
		dep.Status = KnowledgeReadinessReady
	case "missing", "":
		dep.Status = KnowledgeReadinessDegraded
	default:
		dep.Status = KnowledgeReadinessDegraded
	}
	return dep
}

func ragIndexDependencyFromHealth(health sourceHealthReadiness) KnowledgeDataDependency {
	dep := dataDependencyFromHealth("rag_index", health)
	dep.Required = false
	dep.AffectedFeatures = []string{"evidence", "decision_detail"}
	dep.SafeDegradation = "VecLite 不可用时只能回退 SQLite 摘要并标记检索降级。"
	return dep
}

func requiredReadinessCategory(category string) bool {
	switch category {
	case "symbol_profile", "tracked_index", "market_price", "valuation_percentiles", "liquidity":
		return true
	default:
		return false
	}
}

func affectedFeaturesForReadinessCategory(category string) []string {
	switch category {
	case "fund_profile":
		return []string{"consultation", "evidence"}
	case "tracked_index":
		return []string{"consultation", "expected_return", "data_quality"}
	case "market_price":
		return []string{"positions", "daily_discipline", "risk_alerts"}
	case "valuation_percentiles":
		return []string{"margin_of_safety", "expected_return", "consultation", "risk_alerts"}
	case "liquidity":
		return []string{"risk_alerts", "consultation"}
	case "sentiment_proxy":
		return []string{"cool_down", "risk_alerts"}
	default:
		return []string{"data_quality"}
	}
}

func safeDegradationForReadinessCategory(category string) string {
	switch category {
	case "valuation_percentiles":
		return "估值分位缺失时不得声明安全边际或估值高低，只能标记预期收益精度不足。"
	case "market_price":
		return "行情缺失或过期时暂停交易类建议。"
	case "liquidity":
		return "流动性缺失时不得输出大额或市价式行动建议。"
	case "tracked_index":
		return "跟踪指数缺失时 ETF 分析降级为信息不足。"
	case "sentiment_proxy":
		return "情绪代理缺失时不得声明情绪风险已通过。"
	default:
		return "缺失时在数据质量页面显示降级，不伪造成已满足。"
	}
}

func aggregateKnowledgeReadinessStatus(deps []KnowledgeDataDependency) string {
	status := KnowledgeReadinessReady
	for _, dep := range deps {
		if dep.Status == KnowledgeReadinessBlocked {
			return KnowledgeReadinessBlocked
		}
		if dep.Required && dep.Status != KnowledgeReadinessReady {
			status = KnowledgeReadinessDegraded
		}
	}
	return status
}

func featureImpactsFromDependencies(deps []KnowledgeDataDependency) []KnowledgeFeatureImpact {
	out := []KnowledgeFeatureImpact{}
	for _, dep := range deps {
		if dep.Status == KnowledgeReadinessReady {
			continue
		}
		for _, feature := range dep.AffectedFeatures {
			out = append(out, KnowledgeFeatureImpact{Feature: feature, Category: dep.Category, Impact: dep.SafeDegradation, Claims: []string{"不得伪造成 ready", "不得输出交易确认"}})
		}
	}
	return out
}

// BuildKnowledgeLLMContextSummary returns the stable P74 readiness summary used by API/UI and workflow LLM prompts.
func BuildKnowledgeLLMContextSummary(entries []KnowledgeEntry, deps []KnowledgeDataDependency) string {
	parts := make([]appknowledge.DataDependency, 0, len(deps))
	for _, dep := range deps {
		parts = append(parts, appknowledge.DataDependency{Category: dep.Category, Status: dep.Status})
	}
	return appknowledge.BuildLLMContextSummary(entries, parts)
}

func knowledgeReadinessContainsString(items []string, want string) bool {
	for _, item := range items {
		if item == want {
			return true
		}
	}
	return false
}

func knowledgeReadinessFirstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
