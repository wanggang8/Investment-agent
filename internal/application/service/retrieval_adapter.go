package service

import (
	"context"
	"encoding/json"
	"sort"
	"strings"
	"time"

	"investment-agent/internal/application/workflow"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

const maxIndexFreshnessAge = 30 * 24 * time.Hour

// RetrievalAdapter 优先查询 VecLite，失败后从 SQLite 摘要生成可裁决证据。
type RetrievalAdapter struct {
	tx    repository.Transactor
	index VectorIndex
}

func NewRetrievalAdapter(tx repository.Transactor, index VectorIndex) *RetrievalAdapter {
	return &RetrievalAdapter{tx: tx, index: index}
}

func (r *RetrievalAdapter) RetrieveEvidence(ctx context.Context, req workflow.RetrievalRequest) (workflow.RetrievalResult, error) {
	degradedReason := "veclite index not configured"
	indexHealth := r.indexHealth(ctx)
	summaries, summaryErr := r.listSummaries(ctx, req.Symbol)
	if r.index != nil {
		chunks, err := r.searchIndex(ctx, req)
		if err == nil && len(chunks) > 0 {
			matched, inconsistent := summariesForConsistentChunks(chunks, summaries)
			if inconsistent {
				degradedReason = "veclite metadata inconsistent"
			} else if len(matched) == 0 {
				degradedReason = "veclite metadata missing"
			} else {
				matched = rerankRetrievalSummaries(req, matched, chunks)
				return workflow.RetrievalResult{EvidenceSet: evidenceSetFromSummaries(matched), OutputRef: outputRefForRerankedEvidence(matched, chunks), QualitySummary: retrievalQualitySummary(querySummary(req), len(matched), indexHealth, r.indexFallbackSource(), "", chunks...)}, nil
			}
		} else if err != nil {
			degradedReason = "veclite search failed"
		} else {
			degradedReason = "veclite search empty"
		}
	}
	if summaryErr != nil {
		reason := "sqlite summary unavailable"
		return workflow.RetrievalResult{DegradedReason: reason, QualitySummary: retrievalQualitySummary(req.Symbol, 0, indexHealth, "none", reason)}, summaryErr
	}
	if len(summaries) == 0 {
		reason := "sqlite summary empty"
		return workflow.RetrievalResult{DegradedReason: reason, QualitySummary: retrievalQualitySummary(req.Symbol, 0, indexHealth, "none", reason)}, nil
	}
	// VecLite 不可用时只使用 SQLite 摘要降级，C 级信源仍只能转成 background。
	consistency := ""
	if degradedReason == "veclite metadata inconsistent" {
		consistency = "mismatch"
	}
	return workflow.RetrievalResult{EvidenceSet: evidenceSetFromSummaries(summaries), OutputRef: summaries[0].SummaryID, DegradedReason: degradedReason, QualitySummary: retrievalQualitySummaryWithConsistency(req.Symbol, len(summaries), indexHealth, "sqlite_summary", degradedReason, consistency)}, nil
}

func (r *RetrievalAdapter) searchIndex(ctx context.Context, req workflow.RetrievalRequest) ([]repository.RAGChunk, error) {
	if semantic, ok := r.index.(SemanticVectorIndex); ok {
		return semantic.SearchSimilar(ctx, VectorSearchQuery{Text: rewriteRetrievalQuery(req), Symbol: req.Symbol, TopK: semanticCandidateTopK(req.TopK)})
	}
	return r.index.Search(ctx, req.Symbol)
}

func (r *RetrievalAdapter) indexFallbackSource() string {
	if _, ok := r.index.(SemanticVectorIndex); ok {
		return "sqlite_vec"
	}
	return "veclite"
}

func querySummary(req workflow.RetrievalRequest) string {
	if strings.TrimSpace(req.Query) != "" {
		return req.Query
	}
	return req.Symbol
}

func rewriteRetrievalQuery(req workflow.RetrievalRequest) string {
	parts := []string{}
	if symbol := strings.TrimSpace(req.Symbol); symbol != "" {
		parts = append(parts, symbol)
	}
	if query := strings.TrimSpace(req.Query); query != "" {
		parts = append(parts, query)
	}
	raw := strings.ToLower(strings.Join(parts, " "))
	terms := []string{"公告", "风险", "估值", "source verification", "正式证据"}
	if strings.Contains(raw, "买") || strings.Contains(raw, "buy") || strings.Contains(raw, "能不能") || strings.Contains(raw, "能吗") {
		terms = append(terms, "买入纪律", "持仓状态", "现金比例")
	}
	if strings.Contains(raw, "卖") || strings.Contains(raw, "sell") || strings.Contains(raw, "止损") {
		terms = append(terms, "卖出纪律", "风险预警", "回撤")
	}
	if strings.Contains(raw, "资金") || strings.Contains(raw, "流") {
		terms = append(terms, "资金流", "两融", "净流入")
	}
	if strings.Contains(raw, "组合") || strings.Contains(raw, "仓") {
		terms = append(terms, "组合约束", "仓位纪律")
	}
	parts = append(parts, uniqueNonEmptyStrings(terms)...)
	return strings.Join(uniqueNonEmptyStrings(parts), " ")
}

func semanticCandidateTopK(requested int) int {
	if requested <= 0 {
		return 8
	}
	widened := requested * 4
	if widened < 8 {
		return 8
	}
	return widened
}

func (r *RetrievalAdapter) indexHealth(ctx context.Context) string {
	if r == nil || r.index == nil {
		return VectorIndexHealthMissing
	}
	if provider, ok := r.index.(interface {
		Health(context.Context) VectorIndexHealth
	}); ok {
		health := provider.Health(ctx)
		if health.Status != "" {
			return health.Status
		}
	}
	return VectorIndexHealthHealthy
}

func retrievalQualitySummary(query string, topK int, indexHealth string, fallbackSource string, degradedReason string, chunks ...repository.RAGChunk) workflow.RetrievalQualitySummary {
	return retrievalQualitySummaryWithConsistency(query, topK, indexHealth, fallbackSource, degradedReason, "", chunks...)
}

func retrievalQualitySummaryWithConsistency(query string, topK int, indexHealth string, fallbackSource string, degradedReason string, consistencyOverride string, chunks ...repository.RAGChunk) workflow.RetrievalQualitySummary {
	status := "hit"
	consistency := "checked"
	if topK == 0 {
		status = "miss"
		consistency = "not_checked"
	}
	if degradedReason != "" && topK > 0 {
		status = "degraded"
	}
	if consistencyOverride != "" {
		consistency = consistencyOverride
	}
	return workflow.RetrievalQualitySummary{QuerySummary: query, TopK: topK, Status: status, IndexHealth: indexHealth, IndexFreshness: retrievalIndexFreshness(indexHealth, chunks), FallbackSource: fallbackSource, SourceConsistencyStatus: consistency, DegradedReason: degradedReason}
}

func retrievalIndexFreshness(indexHealth string, chunks []repository.RAGChunk) string {
	if indexHealth != VectorIndexHealthHealthy {
		return "unknown"
	}
	validIndexedAt := false
	for _, chunk := range chunks {
		if chunk.IndexedAt == "" {
			continue
		}
		indexedAt, err := time.Parse(time.RFC3339, chunk.IndexedAt)
		if err != nil {
			return "unknown"
		}
		validIndexedAt = true
		if time.Since(indexedAt) > maxIndexFreshnessAge {
			return "stale"
		}
	}
	if !validIndexedAt {
		return "unknown"
	}
	return "fresh"
}

func (r *RetrievalAdapter) listSummaries(ctx context.Context, symbol string) ([]repository.IntelligenceSummary, error) {
	if r.tx == nil {
		return nil, apperr.New(apperr.CodeEvidenceNotFound, apperr.CategoryInternal, "检索依赖缺失")
	}
	var summaries []repository.IntelligenceSummary
	err := r.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		got, err := repos.IntelligenceRepo.ListEvidenceSummaries(ctx)
		if err != nil {
			return err
		}
		for _, summary := range got {
			if symbol == "" || summary.Symbol == "" || summary.Symbol == symbol {
				summaries = append(summaries, summary)
			}
		}
		return nil
	})
	return summaries, err
}

func (r *RetrievalAdapter) VectorIndex() VectorIndex {
	if r == nil {
		return nil
	}
	return r.index
}

func (r *RetrievalAdapter) VectorIndexPath() string {
	switch index := r.index.(type) {
	case *MemoryVectorIndex:
		return index.Path
	case *FileVectorIndex:
		return index.Path
	default:
		return ""
	}
}

func summariesForConsistentChunks(chunks []repository.RAGChunk, summaries []repository.IntelligenceSummary) ([]repository.IntelligenceSummary, bool) {
	byID := make(map[string]repository.IntelligenceSummary, len(summaries))
	for _, summary := range summaries {
		byID[summary.SummaryID] = summary
	}
	matched := make([]repository.IntelligenceSummary, 0, len(chunks))
	inconsistent := false
	for _, chunk := range chunks {
		if summary, ok := byID[chunk.SummaryID]; ok {
			if chunkMetadataConflicts(chunk, summary) {
				inconsistent = true
				continue
			}
			matched = append(matched, summary)
		} else {
			inconsistent = true
		}
	}
	return matched, inconsistent
}

func outputRefForRerankedEvidence(summaries []repository.IntelligenceSummary, chunks []repository.RAGChunk) string {
	if len(summaries) > 0 {
		for _, chunk := range chunks {
			if chunk.SummaryID == summaries[0].SummaryID && chunk.ChunkID != "" {
				return chunk.ChunkID
			}
		}
	}
	if len(chunks) > 0 {
		return chunks[0].ChunkID
	}
	return ""
}

type retrievalCandidate struct {
	summary repository.IntelligenceSummary
	chunk   repository.RAGChunk
	score   float64
	rank    int
}

func rerankRetrievalSummaries(req workflow.RetrievalRequest, summaries []repository.IntelligenceSummary, chunks []repository.RAGChunk) []repository.IntelligenceSummary {
	if len(summaries) == 0 {
		return nil
	}
	chunksBySummaryID := make(map[string][]repository.RAGChunk, len(chunks))
	for _, chunk := range chunks {
		chunksBySummaryID[chunk.SummaryID] = append(chunksBySummaryID[chunk.SummaryID], chunk)
	}
	best := make(map[string]retrievalCandidate, len(summaries))
	for rank, summary := range summaries {
		chunk := firstChunkForSummary(chunksBySummaryID[summary.SummaryID])
		candidate := retrievalCandidate{summary: summary, chunk: chunk, rank: rank}
		candidate.score = retrievalCandidateScore(req, candidate)
		if existing, ok := best[summary.SummaryID]; !ok || candidate.score > existing.score {
			best[summary.SummaryID] = candidate
		}
	}
	candidates := make([]retrievalCandidate, 0, len(best))
	for _, candidate := range best {
		candidates = append(candidates, candidate)
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].score == candidates[j].score {
			return candidates[i].rank < candidates[j].rank
		}
		return candidates[i].score > candidates[j].score
	})
	limit := req.TopK
	if limit <= 0 || limit > len(candidates) {
		limit = len(candidates)
	}
	selected := selectDiverseRetrievalCandidates(candidates, limit)
	out := make([]repository.IntelligenceSummary, 0, len(selected))
	for _, candidate := range selected {
		out = append(out, candidate.summary)
	}
	return out
}

func firstChunkForSummary(chunks []repository.RAGChunk) repository.RAGChunk {
	if len(chunks) == 0 {
		return repository.RAGChunk{}
	}
	return chunks[0]
}

func retrievalCandidateScore(req workflow.RetrievalRequest, candidate retrievalCandidate) float64 {
	summary := candidate.summary
	score := 100 - float64(candidate.rank)
	score += retrievalKeywordOverlapScore(rewriteRetrievalQuery(req), strings.Join([]string{summary.Summary, summary.EventType, summary.SourceName, summary.Entity, candidate.chunk.ChunkText}, " "))
	if strings.TrimSpace(req.Symbol) != "" && summary.Symbol == req.Symbol {
		score += 8
	} else if strings.TrimSpace(req.Symbol) != "" && summary.Symbol != "" && summary.Symbol != req.Symbol {
		score -= 40
	}
	switch model.SourceLevel(summary.SourceLevel) {
	case model.SourceLevelS:
		score += 12
	case model.SourceLevelA:
		score += 9
	case model.SourceLevelB:
		score += 5
	case model.SourceLevelC:
		score -= 8
	}
	if model.EvidenceRole(summary.EvidenceRole) == model.EvidenceFormal {
		score += 8
	} else {
		score -= 3
	}
	if model.VerificationStatus(summary.VerificationStatus) == model.VerificationSatisfied && summaryVerificationConsistent(summary) {
		score += 10
	} else if strings.TrimSpace(summary.VerificationStatus) == "" {
		score -= 4
	} else {
		score -= 6
	}
	if summary.HighGradeIndependentSourceCount >= 2 {
		score += 4
	} else if summary.IndependentSourceCount >= 2 {
		score += 2
	}
	if summary.RelevanceScore > 0 {
		score += summary.RelevanceScore * 3
	}
	if summary.TimeWeight > 0 {
		score += summary.TimeWeight * 2
	}
	if candidate.chunk.IndexedAt != "" {
		if indexedAt, err := time.Parse(time.RFC3339, candidate.chunk.IndexedAt); err == nil && time.Since(indexedAt) > maxIndexFreshnessAge {
			score -= 5
		}
	}
	return score
}

func retrievalKeywordOverlapScore(query string, text string) float64 {
	queryTokens := retrievalTokens(query)
	if len(queryTokens) == 0 {
		return 0
	}
	textLower := strings.ToLower(text)
	score := 0.0
	seen := map[string]bool{}
	for _, token := range queryTokens {
		if seen[token] {
			continue
		}
		seen[token] = true
		if strings.Contains(textLower, token) {
			score += 7
		}
	}
	return score
}

func retrievalTokens(text string) []string {
	normalized := strings.NewReplacer("：", " ", ":", " ", "，", " ", ",", " ", "。", " ", ".", " ", "？", " ", "?", " ", "\n", " ", "\t", " ").Replace(strings.ToLower(text))
	parts := strings.Fields(normalized)
	tokens := make([]string, 0, len(parts)+8)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if len([]rune(part)) >= 2 {
			tokens = append(tokens, part)
		}
	}
	for _, term := range []string{"买入", "买入纪律", "卖出", "卖出纪律", "估值", "低估", "分位", "风险", "公告", "资金流", "持仓", "仓位", "source", "verification"} {
		if strings.Contains(normalized, term) {
			tokens = append(tokens, term)
		}
	}
	return uniqueNonEmptyStrings(tokens)
}

func selectDiverseRetrievalCandidates(candidates []retrievalCandidate, limit int) []retrievalCandidate {
	if limit <= 0 || len(candidates) == 0 {
		return nil
	}
	selected := make([]retrievalCandidate, 0, limit)
	usedSummary := map[string]bool{}
	usedEvent := map[string]bool{}
	for _, candidate := range candidates {
		if len(selected) >= limit {
			break
		}
		eventType := strings.TrimSpace(candidate.summary.EventType)
		if usedSummary[candidate.summary.SummaryID] || (eventType != "" && usedEvent[eventType] && hasAlternativeEvent(candidates, usedSummary, usedEvent)) {
			continue
		}
		selected = append(selected, candidate)
		usedSummary[candidate.summary.SummaryID] = true
		if eventType != "" {
			usedEvent[eventType] = true
		}
	}
	for _, candidate := range candidates {
		if len(selected) >= limit {
			break
		}
		if usedSummary[candidate.summary.SummaryID] {
			continue
		}
		selected = append(selected, candidate)
		usedSummary[candidate.summary.SummaryID] = true
	}
	return selected
}

func hasAlternativeEvent(candidates []retrievalCandidate, usedSummary map[string]bool, usedEvent map[string]bool) bool {
	for _, candidate := range candidates {
		if usedSummary[candidate.summary.SummaryID] {
			continue
		}
		eventType := strings.TrimSpace(candidate.summary.EventType)
		if eventType == "" || !usedEvent[eventType] {
			return true
		}
	}
	return false
}

func chunkMetadataConflicts(chunk repository.RAGChunk, summary repository.IntelligenceSummary) bool {
	if chunk.MetadataJSON == "" {
		return false
	}
	var metadata struct {
		Symbol       string `json:"symbol"`
		SourceLevel  string `json:"source_level"`
		EvidenceRole string `json:"evidence_role"`
	}
	if err := json.Unmarshal([]byte(chunk.MetadataJSON), &metadata); err != nil {
		return true
	}
	if metadata.Symbol != "" && summary.Symbol != "" && metadata.Symbol != summary.Symbol {
		return true
	}
	if metadata.SourceLevel != "" && summary.SourceLevel != "" && metadata.SourceLevel != summary.SourceLevel {
		return true
	}
	if metadata.EvidenceRole != "" && summary.EvidenceRole != "" && metadata.EvidenceRole != summary.EvidenceRole {
		return true
	}
	return false
}

func evidenceSetFromSummaries(summaries []repository.IntelligenceSummary) model.EvidenceSet {
	items := make([]model.Evidence, 0, len(summaries))
	verification := model.VerificationFailed
	verificationSet := false
	for _, summary := range summaries {
		level := model.SourceLevel(summary.SourceLevel)
		role := model.EvidenceRole(summary.EvidenceRole)
		status := model.VerificationStatus(summary.VerificationStatus)
		if !status.Valid() {
			status = model.VerificationFailed
		}
		if status == model.VerificationSatisfied && !summaryVerificationConsistent(summary) {
			status = model.VerificationFailed
		}
		if !level.FormalAllowed() {
			role = model.EvidenceBackground
			status = model.VerificationBackgroundOnly
		} else if status != model.VerificationSatisfied {
			role = model.EvidenceBackground
		}
		if !verificationSet {
			verification = status
			verificationSet = true
		} else {
			verification = mergeVerificationStatus(verification, status)
		}
		count := summary.IndependentSourceCount
		if count == 0 {
			count = 1
		}
		highGradeCount := summary.HighGradeIndependentSourceCount
		if highGradeCount == 0 && (level == model.SourceLevelA || level == model.SourceLevelS) {
			highGradeCount = count
		}
		items = append(items, model.Evidence{EvidenceID: summary.SummaryID, SummaryID: summary.SummaryID, SourceLevel: level, Role: role, EventType: model.EventType(summary.EventType), IndependentSourceCount: count, HighGradeIndependentSourceCount: highGradeCount, SourceName: firstNonEmptyString(summary.SourceName, summary.Entity), PublishedAt: summary.PublishedAt, CapturedAt: summary.CapturedAt, OriginalURL: summary.OriginalURL, Summary: summary.Summary, ContentHash: summary.ContentHash, TimeWeight: summary.TimeWeight, RelevanceScore: summary.RelevanceScore})
	}
	sort.SliceStable(items, func(i, j int) bool {
		return evidenceQualityScore(items[i]) > evidenceQualityScore(items[j])
	})
	return model.EvidenceSet{Items: items, VerificationStatus: verification}
}

func summaryVerificationConsistent(summary repository.IntelligenceSummary) bool {
	if summary.VerificationEvidenceRole != "" && summary.EvidenceRole != "" && summary.VerificationEvidenceRole != summary.EvidenceRole {
		return false
	}
	if summary.VerificationEventType != "" && summary.EventType != "" && summary.VerificationEventType != summary.EventType {
		return false
	}
	if summary.VerificationHighestSourceLevel != "" && summary.SourceLevel != "" && sourceLevelRank(model.SourceLevel(summary.VerificationHighestSourceLevel)) < sourceLevelRank(model.SourceLevel(summary.SourceLevel)) {
		return false
	}
	ids := splitJSONStringList(summary.VerificationEvidenceIDsJSON)
	if len(ids) == 0 {
		return false
	}
	for _, id := range ids {
		if id == summary.SummaryID {
			return true
		}
	}
	return false
}

func splitJSONStringList(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	var values []string
	if err := json.Unmarshal([]byte(raw), &values); err == nil {
		out := make([]string, 0, len(values))
		for _, value := range values {
			value = strings.TrimSpace(value)
			if value != "" {
				out = append(out, value)
			}
		}
		return out
	}
	return nil
}

func sourceLevelRank(level model.SourceLevel) int {
	switch level {
	case model.SourceLevelS:
		return 4
	case model.SourceLevelA:
		return 3
	case model.SourceLevelB:
		return 2
	case model.SourceLevelC:
		return 1
	default:
		return 0
	}
}

func evidenceQualityScore(item model.Evidence) float64 {
	score := item.RelevanceScore + item.TimeWeight
	switch item.SourceLevel {
	case model.SourceLevelS:
		score += 4
	case model.SourceLevelA:
		score += 3
	case model.SourceLevelB:
		score += 2
	case model.SourceLevelC:
		score -= 2
	}
	if item.Role == model.EvidenceFormal {
		score += 3
	} else {
		score -= 1
	}
	if item.HighGradeIndependentSourceCount >= 2 {
		score += 2
	} else if item.IndependentSourceCount >= 2 {
		score += 1
	}
	return score
}

func mergeVerificationStatus(current, next model.VerificationStatus) model.VerificationStatus {
	if current == model.VerificationFailed || next == model.VerificationFailed {
		return model.VerificationFailed
	}
	if current == model.VerificationSatisfied || next == model.VerificationSatisfied {
		return model.VerificationSatisfied
	}
	if current == model.VerificationBackgroundOnly || next == model.VerificationBackgroundOnly {
		return model.VerificationBackgroundOnly
	}
	return model.VerificationSatisfied
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func uniqueNonEmptyStrings(values []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	return out
}
