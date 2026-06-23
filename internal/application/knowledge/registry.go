package knowledge

import (
	"sort"
	"strings"
)

const (
	ReadinessReady    = "ready"
	ReadinessDegraded = "degraded"
	ReadinessBlocked  = "blocked"
)

// Entry describes a deterministic built-in principle or rule reference.
type Entry struct {
	KnowledgeID           string   `json:"knowledge_id"`
	Title                 string   `json:"title"`
	Category              string   `json:"category"`
	Summary               string   `json:"summary"`
	AppliesTo             []string `json:"applies_to"`
	RuleMapping           []string `json:"rule_mapping"`
	LLMContextAllowed     bool     `json:"llm_context_allowed"`
	FormalEvidenceAllowed bool     `json:"formal_evidence_allowed"`
	SafetyBoundary        string   `json:"safety_boundary"`
}

// SymbolProfile describes built-in, locally accepted fund/ETF to index routing.
type SymbolProfile struct {
	Symbol             string
	Name               string
	AssetType          string
	TrackedIndexSymbol string
	TrackedIndexName   string
}

// Registry exposes built-in knowledge in a stable order.
type Registry struct {
	entries []Entry
}

// BuiltInRegistry returns the P74 deterministic built-in registry.
func BuiltInRegistry() Registry {
	entries := []Entry{
		{KnowledgeID: "master.graham.margin_of_safety", Title: "格雷厄姆：安全边际", Category: "master_principle", Summary: "估值分位越低，越需要用安全边际和证据质量约束分批行为。", AppliesTo: []string{"valuation_percentiles", "expected_return", "consultation"}, RuleMapping: []string{"valuation.low_zone", "valuation.high_risk"}, LLMContextAllowed: true, SafetyBoundary: "只能作为纪律原则和 LLM 分析背景，不能作为正式市场证据。"},
		{KnowledgeID: "master.buffett.circle_of_competence", Title: "巴菲特：能力圈", Category: "master_principle", Summary: "能力圈外标的应拒绝交易类分析，只允许提出研究前置问题。", AppliesTo: []string{"capability", "consultation"}, RuleMapping: []string{"capability.out_of_scope"}, LLMContextAllowed: true, SafetyBoundary: "能力圈原则不能绕过用户配置或规则裁决。"},
		{KnowledgeID: "master.livermore.trend_discipline", Title: "利弗莫尔：趋势与耐心", Category: "master_principle", Summary: "趋势和买入逻辑没有被正式证据破坏前，不因短期波动自动推翻计划。", AppliesTo: []string{"trend_risk", "risk_alerts"}, RuleMapping: []string{"buy_logic.broken", "take_profit.trailing"}, LLMContextAllowed: true, SafetyBoundary: "趋势原则不得输出确定涨跌预测或交易确认。"},
		{KnowledgeID: "master.dalio.risk_parity_cycle", Title: "达利欧：风险平价与周期", Category: "master_principle", Summary: "组合风险应看贡献而不只看金额，宏观周期只能作为风险复核问题。", AppliesTo: []string{"portfolio_risk", "macro_cycle", "review"}, RuleMapping: []string{"portfolio.risk_contribution", "cycle.review_required"}, LLMContextAllowed: true, SafetyBoundary: "周期判断不得输出确定宏观预测或自动调仓。"},
		{KnowledgeID: "master.marks.second_level_thinking", Title: "霍华德·马克斯：周期与第二层思维", Category: "master_principle", Summary: "狂热时要求反向风险检查，恐慌时要求证据和估值共同支持。", AppliesTo: []string{"sentiment_proxy", "risk_alerts"}, RuleMapping: []string{"sentiment.extreme", "evidence.insufficient"}, LLMContextAllowed: true, SafetyBoundary: "只能帮助组织分析角度，不替代 source verification。"},
		{KnowledgeID: "master.lynch.know_what_you_own", Title: "彼得·林奇：懂你持有什么", Category: "master_principle", Summary: "无法用简单语言说明标的、行业和基金跟踪逻辑时，先补研究而不是行动。", AppliesTo: []string{"capability", "fund_profile", "consultation"}, RuleMapping: []string{"capability.research_required", "fund_profile.required"}, LLMContextAllowed: true, SafetyBoundary: "生活常识只能生成研究问题，不能替代基金画像或正式证据。"},
		{KnowledgeID: "master.templeton.extreme_pessimism", Title: "邓普顿：极度悲观点", Category: "master_principle", Summary: "极端悲观只允许触发小仓位研究和风险复核，不自动形成买入建议。", AppliesTo: []string{"sentiment_proxy", "valuation_percentiles", "risk_alerts"}, RuleMapping: []string{"sentiment.extreme_fear", "position.small_probe_review"}, LLMContextAllowed: true, SafetyBoundary: "逆向原则不得绕过能力圈、正式证据或用户确认。"},
		{KnowledgeID: "discipline.no_single_source_decision", Title: "纪律：不凭单一信源决策", Category: "discipline_rule", Summary: "重大利好、重大利空或买入逻辑破坏必须满足多源验证。", AppliesTo: []string{"formal_evidence", "source_verification"}, RuleMapping: []string{"evidence.min_high_grade_sources"}, LLMContextAllowed: true, SafetyBoundary: "规则说明不是外部事实来源。"},
		{KnowledgeID: "discipline.no_extreme_emotion_trade", Title: "纪律：极端情绪冷静", Category: "discipline_rule", Summary: "情绪极端时暂停主动交易建议，只展示冷静提醒和既有计划。", AppliesTo: []string{"sentiment_proxy", "user_emotion_tags"}, RuleMapping: []string{"sentiment.cool_down"}, LLMContextAllowed: true, SafetyBoundary: "不得把情绪标签解释为确定市场方向。"},
		{KnowledgeID: "risk_sop.evidence_insufficient", Title: "SOP：证据不足", Category: "risk_sop", Summary: "正式证据不足时进入冻结观察、非交易记录或信息不足状态。", AppliesTo: []string{"formal_evidence", "rag_index", "consultation"}, RuleMapping: []string{"verdict.insufficient_data", "position.frozen_watch"}, LLMContextAllowed: true, SafetyBoundary: "不能用 C 级背景材料补足正式证据。"},
		{KnowledgeID: "risk_sop.valuation_high", Title: "SOP：估值高位", Category: "risk_sop", Summary: "核心估值高危时禁止新增买入，并优先检查止盈和仓位风险。", AppliesTo: []string{"valuation_percentiles", "risk_alerts"}, RuleMapping: []string{"valuation.high_risk", "take_profit"}, LLMContextAllowed: true, SafetyBoundary: "估值状态只约束纪律，不承诺未来涨跌。"},
		{KnowledgeID: "symbol_profile.510300", Title: "标的画像：510300", Category: "symbol_profile", Summary: "510300 作为沪深300 ETF 本地验收主路径，跟踪指数映射为 000300。", AppliesTo: []string{"510300", "000300"}, RuleMapping: []string{"symbol_profile.etf_index_mapping"}, LLMContextAllowed: true, SafetyBoundary: "标的画像用于路由数据依赖，不代表投资推荐。"},
		{KnowledgeID: "symbol_profile.159915", Title: "标的画像：159915", Category: "symbol_profile", Summary: "159915 作为创业板 ETF 动态标的验收路径，跟踪指数映射为 399006。", AppliesTo: []string{"159915", "399006"}, RuleMapping: []string{"symbol_profile.etf_index_mapping"}, LLMContextAllowed: true, SafetyBoundary: "标的画像用于路由数据依赖，不代表投资推荐。"},
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].KnowledgeID < entries[j].KnowledgeID })
	return Registry{entries: entries}
}

// BuiltInSymbolProfiles returns the known local symbol profiles accepted by
// current readiness and collector routing.
func BuiltInSymbolProfiles() []SymbolProfile {
	return []SymbolProfile{
		{Symbol: "510300", Name: "沪深300ETF", AssetType: "ETF", TrackedIndexSymbol: "000300", TrackedIndexName: "沪深300"},
		{Symbol: "159915", Name: "创业板ETF", AssetType: "ETF", TrackedIndexSymbol: "399006", TrackedIndexName: "创业板指"},
	}
}

// LookupSymbolProfile returns a known local symbol profile without fabricating
// a fallback for unsupported symbols.
func LookupSymbolProfile(symbol string) (SymbolProfile, bool) {
	symbol = strings.TrimSpace(symbol)
	for _, profile := range BuiltInSymbolProfiles() {
		if profile.Symbol == symbol {
			return profile, true
		}
	}
	return SymbolProfile{Symbol: symbol}, false
}

// DefaultTrackedIndexSymbol returns the built-in index routing for a known ETF.
func DefaultTrackedIndexSymbol(symbol string) string {
	profile, ok := LookupSymbolProfile(symbol)
	if !ok {
		return ""
	}
	return profile.TrackedIndexSymbol
}

// Entries returns a defensive copy of built-in entries.
func (r Registry) Entries() []Entry {
	out := append([]Entry(nil), r.entries...)
	return out
}

// EntriesForSymbol returns shared principles plus only the matching known
// symbol profile. Unknown symbols intentionally receive no symbol profile entry.
func (r Registry) EntriesForSymbol(symbol string) []Entry {
	symbol = strings.TrimSpace(symbol)
	out := make([]Entry, 0, len(r.entries))
	for _, entry := range r.entries {
		if entry.Category == "symbol_profile" && !entryAppliesToSymbol(entry, symbol) {
			continue
		}
		out = append(out, entry)
	}
	return out
}

func entryAppliesToSymbol(entry Entry, symbol string) bool {
	if symbol == "" {
		return false
	}
	for _, appliesTo := range entry.AppliesTo {
		if appliesTo == symbol {
			return true
		}
	}
	return false
}

// DataDependency is the minimal dependency shape needed to build LLM readiness context.
type DataDependency struct {
	Category string `json:"category"`
	Status   string `json:"status"`
}

// BuildLLMContextSummary returns the stable P74 readiness summary used by API/UI and workflow LLM prompts.
func BuildLLMContextSummary(entries []Entry, deps []DataDependency) string {
	ids := []string{}
	for _, entry := range entries {
		if entry.LLMContextAllowed {
			ids = append(ids, entry.KnowledgeID)
		}
	}
	depParts := []string{}
	for _, dep := range deps {
		depParts = append(depParts, dep.Category+"="+dep.Status)
	}
	return "principles=" + strings.Join(ids, ",") + "; data_readiness=" + strings.Join(depParts, ",") + "; boundary=LLM只生成分析材料，背景知识不能满足正式证据，最终裁决由规则引擎负责"
}
