package workflow

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/cloudwego/eino/compose"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
)

// EvidenceVerificationInput 是证据核查工作流输入。
type EvidenceVerificationInput struct {
	RequestID string
	Symbol    string
	Sources   []string
}

// EvidenceVerificationOutput 汇总证据工作流允许写入的事实。
type EvidenceVerificationOutput struct {
	WorkflowContext         WorkflowContext
	IntelligenceItems       []string
	IntelligenceSummary     string
	IntelligenceSummaries   []repository.IntelligenceSummary
	RAGChunks               []string
	SourceVerifications     []model.VerificationStatus
	VectorIndexRebuildable  bool
	VectorIndexStatus       string
	VectorIndexFailedReason string
}

// EvidenceVerificationGraph 负责情报刷新、RAG 分块与多源验证事实写入。
type EvidenceVerificationGraph struct {
	auditWriter AuditWriter
	deps        WorkflowDependencies
	nodeNames   []string
}

type evidenceVerificationState struct {
	Input                   EvidenceVerificationInput
	WorkflowContext         WorkflowContext
	Status                  NodeStatus
	Code                    string
	Now                     string
	Items                   []IntelligenceSourceItem
	IntelligenceItems       []repository.IntelligenceItem
	Summaries               []repository.IntelligenceSummary
	Chunks                  []repository.RAGChunk
	EvidenceItems           []model.Evidence
	EvidenceIDs             []string
	Verification            repository.SourceVerification
	FirstSummaryID          string
	ResultOutputRef         string
	VectorIndexStatus       NodeStatus
	VectorIndexCode         string
	VectorIndexFailedReason string
}

// NewEvidenceVerificationGraph 创建证据核查工作流。
func NewEvidenceVerificationGraph(writer AuditWriter) *EvidenceVerificationGraph {
	if writer == nil {
		writer = &MemoryAuditWriter{}
	}
	return &EvidenceVerificationGraph{auditWriter: writer, nodeNames: evidenceVerificationNodeNames()}
}

// NewEvidenceVerificationGraphWithDependencies 创建带 SQLite 写入能力的证据核查工作流。
func NewEvidenceVerificationGraphWithDependencies(deps WorkflowDependencies) *EvidenceVerificationGraph {
	return &EvidenceVerificationGraph{auditWriter: NewRepositoryAuditWriter(deps.AuditRepo), deps: deps, nodeNames: evidenceVerificationNodeNames()}
}

func evidenceVerificationNodeNames() []string {
	return []string{"NewsFetchNode", "NewsClassifyNode", "EvidenceNormalizeNode", "EmbeddingNode", "VectorStoreNode", "SourceVerificationNode"}
}

func (g *EvidenceVerificationGraph) NodeNames() []string {
	return append([]string(nil), g.nodeNames...)
}

// RegisteredNodeNames 返回真实注册到 Eino Graph 的业务节点。
func (g *EvidenceVerificationGraph) RegisteredNodeNames() []string {
	return append([]string(nil), g.nodeNames...)
}

func writeEvidenceVectorIndex(ctx context.Context, index VectorIndexWriter, chunks []repository.RAGChunk) error {
	if index == nil {
		return nil
	}
	for _, chunk := range chunks {
		if err := index.Upsert(ctx, chunk); err != nil {
			return err
		}
	}
	return nil
}

// Run 生成并保存 intelligence_items、intelligence_summary、rag_chunks 和 source_verifications。
func (g *EvidenceVerificationGraph) Run(ctx context.Context, in EvidenceVerificationInput) (EvidenceVerificationOutput, error) {
	runnable, err := g.compile(ctx)
	if err != nil {
		return EvidenceVerificationOutput{}, err
	}
	state, err := runnable.Invoke(ctx, evidenceVerificationState{Input: in})
	if err != nil {
		return EvidenceVerificationOutput{}, err
	}
	return state.output(), nil
}

func (g *EvidenceVerificationGraph) compile(ctx context.Context) (compose.Runnable[evidenceVerificationState, evidenceVerificationState], error) {
	graph := compose.NewGraph[evidenceVerificationState, evidenceVerificationState]()
	nodes := []struct {
		name string
		fn   func(context.Context, evidenceVerificationState) (evidenceVerificationState, error)
	}{
		{name: "NewsFetchNode", fn: g.newsFetchNode},
		{name: "NewsClassifyNode", fn: g.newsClassifyNode},
		{name: "EvidenceNormalizeNode", fn: g.evidenceNormalizeNode},
		{name: "EmbeddingNode", fn: g.embeddingNode},
		{name: "VectorStoreNode", fn: g.vectorStoreNode},
		{name: "SourceVerificationNode", fn: g.sourceVerificationNode},
	}
	for _, node := range nodes {
		node := node
		if err := graph.AddLambdaNode(node.name, compose.InvokableLambda(node.fn)); err != nil {
			return nil, err
		}
	}
	for i, node := range nodes {
		from := compose.START
		if i > 0 {
			from = nodes[i-1].name
		}
		if err := graph.AddEdge(from, node.name); err != nil {
			return nil, err
		}
	}
	if err := graph.AddEdge(nodes[len(nodes)-1].name, compose.END); err != nil {
		return nil, err
	}
	return graph.Compile(ctx)
}

func (g *EvidenceVerificationGraph) newsFetchNode(ctx context.Context, state evidenceVerificationState) (evidenceVerificationState, error) {
	state.WorkflowContext = WorkflowContext{RequestID: state.Input.RequestID, WorkflowType: WorkflowEvidenceVerification, Symbol: state.Input.Symbol, RuleVersion: workflowRuleVersion(ctx, g.deps.RuleRepo)}
	state.Status, state.Code = StatusSuccess, ""
	state.VectorIndexStatus, state.VectorIndexCode = StatusSuccess, ""
	state.Now = workflowNowRFC3339()
	items, err := g.deps.intelligenceSource().FetchIntelligence(ctx, state.Input.Symbol)
	fetchStatus := StatusSuccess
	fetchCode := ""
	if err != nil {
		state.Status, state.Code = StatusFailed, ErrCodeSourceVerificationFailed
		fetchStatus, fetchCode = StatusFailed, ErrCodeSourceVerificationFailed
	}
	if err == nil && len(items) == 0 {
		state.Status, state.Code = StatusFailed, ErrCodeSourceVerificationFailed
		fetchStatus, fetchCode = StatusFailed, ErrCodeSourceVerificationFailed
	}
	if err == nil && len(items) > 0 && countIndependentSources(items) < 2 {
		state.Status, state.Code = StatusFailed, ErrCodeSourceVerificationFailed
	}
	state.Items = items
	return state, g.writeEvidenceAudit(ctx, &state.WorkflowContext, NodeResult{Status: fetchStatus, ErrorCode: fetchCode, Audit: AuditFragment{Action: string(model.AuditActionRebuildIndex), NodeName: "NewsFetchNode", NodeAction: "fetch_news", Status: fetchStatus, InputRefType: "symbol", InputRef: state.Input.Symbol, OutputRefType: "intelligence_items", OutputRef: state.Input.Symbol, ErrorCode: fetchCode}})
}

func (g *EvidenceVerificationGraph) newsClassifyNode(ctx context.Context, state evidenceVerificationState) (evidenceVerificationState, error) {
	return state, g.writeEvidenceAudit(ctx, &state.WorkflowContext, NodeResult{Status: StatusSuccess, Audit: AuditFragment{Action: string(model.AuditActionRebuildIndex), NodeName: "NewsClassifyNode", NodeAction: "classify_news", Status: StatusSuccess, InputRefType: "intelligence_items", InputRef: state.Input.Symbol, OutputRefType: "event_type", OutputRef: string(model.EventNormal)}})
}

func (g *EvidenceVerificationGraph) evidenceNormalizeNode(ctx context.Context, state evidenceVerificationState) (evidenceVerificationState, error) {
	verificationGroupID := workflowID("group")
	for i, src := range state.Items {
		summaryID := workflowID("summary")
		state.EvidenceIDs = append(state.EvidenceIDs, summaryID)
		level := src.SourceLevel
		if level == "" {
			level = model.SourceLevelC
		}
		role := model.EvidenceBackground
		if state.Status != StatusFailed && level.FormalAllowed() {
			role = model.EvidenceFormal
		}
		intelID := workflowID("intel")
		chunkID := workflowID("chunk")
		text := src.Text
		if text == "" {
			text = src.Title
		}
		contentHash := stableHash(src.SourceName, src.URL, src.PublishedAt, src.Title, text)
		chunkHash := stableHash(state.Input.Symbol, src.SourceName, src.URL, src.PublishedAt, src.Title, text)
		independentSourceCount := countIndependentSources(state.Items)
		highGradeCount := highGradeIndependentSourceCount(state.Items)
		timeWeight := evidenceTimeWeight(src.PublishedAt, state.Now)
		relevanceScore := evidenceRelevanceScore(state.Input.Symbol, src)
		item := repository.IntelligenceItem{IntelligenceID: intelID, SourceName: src.SourceName, SourceLevel: string(level), OriginalURL: src.URL, PublishedAt: src.PublishedAt, CapturedAt: state.Now, ContentHash: contentHash, RawTitle: src.Title, RawTextRef: text, CreatedAt: state.Now}
		summary := repository.IntelligenceSummary{SummaryID: summaryID, IntelligenceID: intelID, Symbol: state.Input.Symbol, EventType: string(model.EventNormal), ImpactDirection: "neutral", Summary: text, SourceLevel: string(level), EvidenceRole: string(role), TimeWeight: timeWeight, RelevanceScore: relevanceScore, VerificationGroupID: verificationGroupID, IndependentSourceCount: independentSourceCount, HighGradeIndependentSourceCount: highGradeCount, SourceName: src.SourceName, OriginalURL: src.URL, PublishedAt: src.PublishedAt, CapturedAt: state.Now, ContentHash: contentHash, CreatedAt: state.Now}
		chunk := repository.RAGChunk{ChunkID: chunkID, SummaryID: summaryID, Symbol: state.Input.Symbol, ChunkText: text, ChunkHash: chunkHash, IndexStatus: "pending", CreatedAt: state.Now}
		if i == 0 {
			state.FirstSummaryID = summaryID
			state.ResultOutputRef = chunkID
		}
		state.IntelligenceItems = append(state.IntelligenceItems, item)
		state.Summaries = append(state.Summaries, summary)
		state.Chunks = append(state.Chunks, chunk)
		state.EvidenceItems = append(state.EvidenceItems, model.Evidence{EvidenceID: summaryID, SummaryID: summaryID, SourceLevel: level, Role: role, EventType: model.EventNormal, IndependentSourceCount: independentSourceCount, HighGradeIndependentSourceCount: highGradeCount, SourceName: src.SourceName, PublishedAt: src.PublishedAt, CapturedAt: state.Now, OriginalURL: src.URL, Summary: text, ContentHash: contentHash, ChunkHash: chunkHash, TimeWeight: timeWeight, RelevanceScore: relevanceScore})
	}
	return state, g.writeEvidenceAudit(ctx, &state.WorkflowContext, NodeResult{Status: StatusSuccess, Audit: AuditFragment{Action: string(model.AuditActionRebuildIndex), NodeName: "EvidenceNormalizeNode", NodeAction: "normalize_evidence", Status: StatusSuccess, InputRefType: "symbol", InputRef: state.Input.Symbol, OutputRefType: "evidence", OutputRef: state.Input.Symbol}})
}

func (g *EvidenceVerificationGraph) embeddingNode(ctx context.Context, state evidenceVerificationState) (evidenceVerificationState, error) {
	inputRef := firstNonEmpty(state.ResultOutputRef, state.Input.Symbol)
	return state, g.writeEvidenceAudit(ctx, &state.WorkflowContext, NodeResult{Status: StatusSuccess, Audit: AuditFragment{Action: string(model.AuditActionRebuildIndex), NodeName: "EmbeddingNode", NodeAction: "embed_chunks", Status: StatusSuccess, InputRefType: "evidence", InputRef: inputRef, OutputRefType: "embedding", OutputRef: inputRef}})
}

func (g *EvidenceVerificationGraph) vectorStoreNode(ctx context.Context, state evidenceVerificationState) (evidenceVerificationState, error) {
	if err := g.writeIntelligenceFacts(ctx, state); err != nil {
		return state, err
	}
	indexStatus := "indexed"
	if err := writeEvidenceVectorIndex(ctx, g.deps.VectorIndexWriter, state.Chunks); err != nil {
		state.VectorIndexStatus, state.VectorIndexCode = StatusFailed, ErrCodeVectorIndexUnavailable
		state.VectorIndexFailedReason = err.Error()
		indexStatus = "failed"
	}
	if err := g.updateRAGChunksIndexStatus(ctx, state.Chunks, indexStatus); err != nil {
		return state, err
	}
	ref := firstNonEmpty(state.ResultOutputRef, state.Input.Symbol)
	return state, g.writeEvidenceAudit(ctx, &state.WorkflowContext, NodeResult{Status: state.VectorIndexStatus, ErrorCode: state.VectorIndexCode, Audit: AuditFragment{Action: string(model.AuditActionRebuildIndex), NodeName: "VectorStoreNode", NodeAction: "write_vector_index", Status: state.VectorIndexStatus, InputRefType: "embedding", InputRef: ref, OutputRefType: "rag_chunks", OutputRef: ref, ErrorCode: state.VectorIndexCode}})
}

func (g *EvidenceVerificationGraph) sourceVerificationNode(ctx context.Context, state evidenceVerificationState) (evidenceVerificationState, error) {
	verificationStatus := model.VerificationSatisfied
	verificationRole := model.EvidenceFormal
	if state.Status == StatusFailed {
		verificationStatus = model.VerificationFailed
		verificationRole = model.EvidenceBackground
	} else if highGradeIndependentSourceCount(state.Items) < 2 {
		verificationStatus = model.VerificationBackgroundOnly
		verificationRole = model.EvidenceBackground
	}
	evidenceIDsJSON, _ := json.Marshal(state.EvidenceIDs)
	state.Verification = repository.SourceVerification{VerificationID: workflowID("verify"), VerificationGroupID: firstNonEmptyVerificationGroup(state.Summaries), EventID: workflowID("event"), Symbol: state.Input.Symbol, EventType: string(model.EventNormal), EvidenceRole: string(verificationRole), VerificationStatus: string(verificationStatus), IndependentSourceCount: countIndependentSources(state.Items), HighGradeIndependentSourceCount: highGradeIndependentSourceCount(state.Items), HighestSourceLevel: string(highestSourceLevel(state.Items)), LatestPublishedAt: latestPublishedAt(state.Items), EvidenceIDsJSON: string(evidenceIDsJSON), CreatedAt: state.Now}
	if err := g.writeSourceVerification(ctx, state.Verification); err != nil {
		return state, err
	}
	state.WorkflowContext.EvidenceSet = model.EvidenceSet{Items: state.EvidenceItems, VerificationStatus: model.VerificationStatus(state.Verification.VerificationStatus)}
	return state, g.writeEvidenceAudit(ctx, &state.WorkflowContext, NodeResult{Status: state.Status, ErrorCode: state.Code, Audit: AuditFragment{Action: string(model.AuditActionRebuildIndex), NodeName: "SourceVerificationNode", NodeAction: "verify_sources", Status: state.Status, InputRefType: "symbol", InputRef: state.Input.Symbol, OutputRefType: "source_verification", OutputRef: state.Input.Symbol, ErrorCode: state.Code}})
}

func (g *EvidenceVerificationGraph) writeIntelligenceFacts(ctx context.Context, state evidenceVerificationState) error {
	write := func(txCtx context.Context, repos repository.Repositories) error {
		if repos.IntelligenceRepo == nil {
			return nil
		}
		for i, item := range state.IntelligenceItems {
			if err := repos.IntelligenceRepo.SaveIntelligenceItem(txCtx, item); err != nil {
				return err
			}
			if err := repos.IntelligenceRepo.SaveIntelligenceSummary(txCtx, state.Summaries[i], []repository.RAGChunk{state.Chunks[i]}); err != nil {
				return err
			}
		}
		return nil
	}
	if g.deps.Transactor != nil && g.deps.IntelligenceRepo != nil {
		return g.deps.Transactor.WithinTx(ctx, write)
	}
	return write(ctx, g.deps.repositories())
}

func (g *EvidenceVerificationGraph) updateRAGChunksIndexStatus(ctx context.Context, chunks []repository.RAGChunk, status string) error {
	chunkIDs := make([]string, 0, len(chunks))
	for _, chunk := range chunks {
		chunkIDs = append(chunkIDs, chunk.ChunkID)
	}
	write := func(txCtx context.Context, repos repository.Repositories) error {
		if repos.IntelligenceRepo == nil {
			return nil
		}
		return repos.IntelligenceRepo.UpdateRAGChunksIndexStatus(txCtx, chunkIDs, status)
	}
	if g.deps.Transactor != nil && g.deps.IntelligenceRepo != nil {
		return g.deps.Transactor.WithinTx(ctx, write)
	}
	return write(ctx, g.deps.repositories())
}

func (g *EvidenceVerificationGraph) writeSourceVerification(ctx context.Context, verification repository.SourceVerification) error {
	write := func(txCtx context.Context, repos repository.Repositories) error {
		if repos.IntelligenceRepo == nil {
			return nil
		}
		return repos.IntelligenceRepo.SaveSourceVerification(txCtx, verification)
	}
	if g.deps.Transactor != nil && g.deps.IntelligenceRepo != nil {
		return g.deps.Transactor.WithinTx(ctx, write)
	}
	return write(ctx, g.deps.repositories())
}

func (g *EvidenceVerificationGraph) writeEvidenceAudit(ctx context.Context, wf *WorkflowContext, result NodeResult) error {
	if g.deps.Transactor != nil && g.deps.AuditRepo != nil {
		return g.deps.Transactor.WithinTx(ctx, func(txCtx context.Context, repos repository.Repositories) error {
			return writeAuditEvent(txCtx, repos.AuditRepo, wf, result)
		})
	}
	writer := g.auditWriter
	if writer == nil {
		writer = &MemoryAuditWriter{}
	}
	return writer.Write(ctx, wf, result)
}

func (state evidenceVerificationState) output() EvidenceVerificationOutput {
	verifications := []model.VerificationStatus{model.VerificationStatus(state.Verification.VerificationStatus)}
	if state.Verification.VerificationStatus == "" {
		verifications = []model.VerificationStatus{model.VerificationFailed}
	}
	itemIDs := make([]string, 0, len(state.IntelligenceItems))
	for _, item := range state.IntelligenceItems {
		itemIDs = append(itemIDs, item.IntelligenceID)
	}
	chunkIDs := make([]string, 0, len(state.Chunks))
	for _, chunk := range state.Chunks {
		chunkIDs = append(chunkIDs, chunk.ChunkID)
	}
	return EvidenceVerificationOutput{WorkflowContext: state.WorkflowContext, IntelligenceItems: itemIDs, IntelligenceSummary: state.FirstSummaryID, IntelligenceSummaries: append([]repository.IntelligenceSummary(nil), state.Summaries...), RAGChunks: chunkIDs, SourceVerifications: verifications, VectorIndexRebuildable: true, VectorIndexStatus: string(state.VectorIndexStatus), VectorIndexFailedReason: state.VectorIndexFailedReason}
}

func evidenceTimeWeight(publishedAt, now string) float64 {
	if strings.TrimSpace(publishedAt) == "" {
		return 1
	}
	published, err := time.Parse(time.RFC3339, publishedAt)
	if err != nil {
		return 1
	}
	current, err := time.Parse(time.RFC3339, now)
	if err != nil || current.Before(published) {
		return 1
	}
	days := int(current.Sub(published).Hours() / 24)
	if days <= 7 {
		return 1
	}
	if days <= 30 {
		return 0.7
	}
	return 0.4
}

func evidenceRelevanceScore(symbol string, item IntelligenceSourceItem) float64 {
	text := strings.ToLower(strings.Join([]string{item.Title, item.Text, item.SourceName}, " "))
	if strings.TrimSpace(symbol) != "" && strings.Contains(text, strings.ToLower(symbol)) {
		return 1
	}
	return 0.8
}

func firstNonEmptyVerificationGroup(summaries []repository.IntelligenceSummary) string {
	for _, summary := range summaries {
		if summary.VerificationGroupID != "" {
			return summary.VerificationGroupID
		}
	}
	return workflowID("group")
}

func countIndependentSources(items []IntelligenceSourceItem) int {
	seen := map[string]bool{}
	for _, item := range items {
		name := item.SourceName
		if name == "" {
			name = item.URL
		}
		if name != "" {
			seen[name] = true
		}
	}
	if len(seen) == 0 && len(items) > 0 {
		return 1
	}
	return len(seen)
}

func highGradeIndependentSourceCount(items []IntelligenceSourceItem) int {
	seen := map[string]bool{}
	for _, item := range items {
		if item.SourceLevel != model.SourceLevelA && item.SourceLevel != model.SourceLevelS {
			continue
		}
		name := item.SourceName
		if name == "" {
			name = item.URL
		}
		if name != "" {
			seen[name] = true
		}
	}
	return len(seen)
}

func highestSourceLevel(items []IntelligenceSourceItem) model.SourceLevel {
	best := model.SourceLevelC
	for _, item := range items {
		switch item.SourceLevel {
		case model.SourceLevelS:
			return model.SourceLevelS
		case model.SourceLevelA:
			best = model.SourceLevelA
		case model.SourceLevelB:
			if best != model.SourceLevelA {
				best = model.SourceLevelB
			}
		}
	}
	return best
}

func latestPublishedAt(items []IntelligenceSourceItem) string {
	latest := ""
	for _, item := range items {
		if item.PublishedAt > latest {
			latest = item.PublishedAt
		}
	}
	return latest
}
