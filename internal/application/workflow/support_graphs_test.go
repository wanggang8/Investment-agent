package workflow

import (
	"context"
	"strings"
	"testing"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
)

func TestEvidenceVerificationGraphWritesFacts(t *testing.T) {
	out, err := NewEvidenceVerificationGraph(&MemoryAuditWriter{}).Run(context.Background(), EvidenceVerificationInput{RequestID: "req_ev", Symbol: "510300", Sources: []string{"official", "exchange"}})
	if err != nil {
		t.Fatalf("run evidence graph: %v", err)
	}
	if len(out.IntelligenceItems) == 0 || len(out.RAGChunks) == 0 || len(out.SourceVerifications) == 0 || out.IntelligenceSummary == "" {
		t.Fatalf("expected evidence facts: %+v", out)
	}
	if !out.VectorIndexRebuildable {
		t.Fatal("VecLite index must be rebuildable from rag chunks")
	}
}

func TestEvidenceVerificationGraphKeepsAllNormalizedEvidence(t *testing.T) {
	deps := WorkflowDependencies{IntelligenceSource: testIntelligenceSource{items: []IntelligenceSourceItem{
		{SourceName: "official", SourceLevel: model.SourceLevelA, Title: "公告一", Text: "正式证据一"},
		{SourceName: "exchange", SourceLevel: model.SourceLevelA, Title: "公告二", Text: "正式证据二"},
	}}}
	out, err := NewEvidenceVerificationGraphWithDependencies(deps).Run(context.Background(), EvidenceVerificationInput{RequestID: "req_ev_all", Symbol: "510300", Sources: []string{"official", "exchange"}})
	if err != nil {
		t.Fatalf("run evidence graph: %v", err)
	}
	if len(out.WorkflowContext.EvidenceSet.Items) != 2 {
		t.Fatalf("expected all evidence items in context, got %+v", out.WorkflowContext.EvidenceSet.Items)
	}
}

func TestEvidenceVerificationGraphPreservesQualityMetadata(t *testing.T) {
	deps := WorkflowDependencies{IntelligenceSource: testIntelligenceSource{items: []IntelligenceSourceItem{
		{SourceName: "official", SourceLevel: model.SourceLevelA, Title: "公告一", Text: "正式证据一", URL: "https://example.com/1", PublishedAt: "2026-06-01T01:00:00Z"},
		{SourceName: "exchange", SourceLevel: model.SourceLevelA, Title: "公告二", Text: "正式证据二", URL: "https://example.com/2", PublishedAt: "2026-06-01T02:00:00Z"},
	}}}
	out, err := NewEvidenceVerificationGraphWithDependencies(deps).Run(context.Background(), EvidenceVerificationInput{RequestID: "req_ev_quality", Symbol: "510300", Sources: []string{"official", "exchange"}})
	if err != nil {
		t.Fatalf("run evidence graph: %v", err)
	}
	if len(out.WorkflowContext.EvidenceSet.Items) != 2 || len(out.IntelligenceSummaries) != 2 {
		t.Fatalf("expected normalized evidence, got %+v", out)
	}
	item := out.WorkflowContext.EvidenceSet.Items[0]
	summary := out.IntelligenceSummaries[0]
	if item.TimeWeight == 0 || item.RelevanceScore == 0 || summary.TimeWeight == 0 || summary.RelevanceScore == 0 {
		t.Fatalf("expected quality scores preserved in evidence and summary: item=%+v summary=%+v", item, summary)
	}
	if item.IndependentSourceCount != 2 || item.HighGradeIndependentSourceCount != 2 || summary.IndependentSourceCount != 2 || summary.HighGradeIndependentSourceCount != 2 {
		t.Fatalf("expected source counts preserved: item=%+v summary=%+v", item, summary)
	}
	if item.SourceName != "official" || item.PublishedAt != "2026-06-01T01:00:00Z" || item.OriginalURL != "https://example.com/1" || item.ContentHash == "" {
		t.Fatalf("expected source metadata preserved: %+v", item)
	}
}

func TestEvidenceVerificationRequiresTwoHighGradeIndependentSources(t *testing.T) {
	deps := WorkflowDependencies{IntelligenceSource: testIntelligenceSource{items: []IntelligenceSourceItem{
		{SourceName: "official", SourceLevel: model.SourceLevelA, Title: "单源公告", Text: "只有单个高等级来源"},
	}}}
	out, err := NewEvidenceVerificationGraphWithDependencies(deps).Run(context.Background(), EvidenceVerificationInput{RequestID: "req_ev_single_high", Symbol: "510300", Sources: []string{"official", "exchange"}})
	if err != nil {
		t.Fatalf("run evidence graph: %v", err)
	}
	if out.WorkflowContext.EvidenceSet.VerificationStatus == model.VerificationSatisfied {
		t.Fatalf("single high grade source must not satisfy verification: %+v", out.WorkflowContext.EvidenceSet)
	}
}

func TestEvidenceVerificationGraphExposesAuditableSteps(t *testing.T) {
	graph := NewEvidenceVerificationGraph(&MemoryAuditWriter{})
	want := []string{"NewsFetchNode", "NewsClassifyNode", "EvidenceNormalizeNode", "EmbeddingNode", "VectorStoreNode", "SourceVerificationNode"}
	if !sameStrings(graph.NodeNames(), want) {
		t.Fatalf("node names=%+v want=%+v", graph.NodeNames(), want)
	}
	if !sameStrings(graph.RegisteredNodeNames(), want) {
		t.Fatalf("registered node names=%+v want=%+v", graph.RegisteredNodeNames(), want)
	}
}

func TestEvidenceVerificationGraphWritesAuditForEachStep(t *testing.T) {
	out, err := NewEvidenceVerificationGraph(&MemoryAuditWriter{}).Run(context.Background(), EvidenceVerificationInput{RequestID: "req_ev_steps", Symbol: "510300", Sources: []string{"official", "exchange"}})
	if err != nil {
		t.Fatalf("run evidence graph: %v", err)
	}
	for _, node := range []string{"NewsFetchNode", "NewsClassifyNode", "EvidenceNormalizeNode", "EmbeddingNode", "VectorStoreNode", "SourceVerificationNode"} {
		if !hasAuditNode(out.WorkflowContext.AuditEvents, node) {
			t.Fatalf("expected audit node %s in %+v", node, out.WorkflowContext.AuditEvents)
		}
	}
}

func TestMarketRefreshGraphWritesMarketSnapshot(t *testing.T) {
	out, err := NewMarketRefreshGraph(&MemoryAuditWriter{}).Run(context.Background(), MarketRefreshInput{RequestID: "req_mk", Symbol: "510300", PEPercentile: 45, PBPercentile: 40})
	if err != nil {
		t.Fatalf("run market graph: %v", err)
	}
	if out.MarketSnapshot.Symbol != "510300" || out.MarketSnapshot.MarketSnapshotID == "" {
		t.Fatalf("expected market snapshot: %+v", out.MarketSnapshot)
	}
	if !hasAudit(out.AuditEvents, model.AuditActionRefreshMarketData) {
		t.Fatalf("expected market audit: %+v", out.AuditEvents)
	}
}

func TestEvolutionProposalGraphDoesNotUpdateRuleVersion(t *testing.T) {
	out, err := NewEvolutionProposalGraph(&MemoryAuditWriter{}).Run(context.Background(), EvolutionProposalInput{RequestID: "req_evo", ErrorCaseID: "err_1", SampleCount: 5})
	if err != nil {
		t.Fatalf("run evolution graph: %v", err)
	}
	if out.RuleProposal.ProposalID == "" || out.UpdatedRuleVersion {
		t.Fatalf("proposal must not update formal rule: %+v", out)
	}
}

func TestGatekeeperAuditGraphWaitsForFinalConfirmation(t *testing.T) {
	out, err := NewGatekeeperAuditGraph(&MemoryAuditWriter{}).Run(context.Background(), GatekeeperAuditInput{RequestID: "req_gate", ProposalID: "prop_1", Approved: true})
	if err != nil {
		t.Fatalf("run gatekeeper graph: %v", err)
	}
	if out.ProposalStatus != model.ProposalPendingFinalConfirm || out.UpdatedRuleVersion {
		t.Fatalf("audit should wait final confirmation: %+v", out)
	}
	if len(out.GatekeeperAudits) != 1 {
		t.Fatalf("expected one gatekeeper audit: %+v", out.GatekeeperAudits)
	}
}

func TestGatekeeperAuditGraphExposesNodeLevelPlan(t *testing.T) {
	graph := NewGatekeeperAuditGraph(&MemoryAuditWriter{})
	want := gatekeeperAuditNodeNames()
	if !sameStrings(graph.NodeNames(), want) {
		t.Fatalf("node names=%+v want=%+v", graph.NodeNames(), want)
	}
	if !sameStrings(graph.RegisteredNodeNames(), want) {
		t.Fatalf("registered node names=%+v want=%+v", graph.RegisteredNodeNames(), want)
	}
}

func TestGatekeeperAuditGraphWritesAuditForEachNode(t *testing.T) {
	out, err := NewGatekeeperAuditGraph(&MemoryAuditWriter{}).Run(context.Background(), GatekeeperAuditInput{RequestID: "req_gate_nodes", ProposalID: "prop_nodes", Approved: true})
	if err != nil {
		t.Fatalf("run gatekeeper graph: %v", err)
	}
	want := gatekeeperAuditNodeNames()
	if len(out.AuditEvents) != len(want) {
		t.Fatalf("expected one audit event per node, got %+v", out.AuditEvents)
	}
	for _, node := range want {
		if !hasAuditNode(out.AuditEvents, node) {
			t.Fatalf("expected audit node %s in %+v", node, out.AuditEvents)
		}
	}
	for _, event := range out.AuditEvents {
		if event.WorkflowType != WorkflowGatekeeperAudit || event.NodeName == "" || event.NodeAction == "" || event.InputRef == "" || event.OutputRef == "" || event.Status == "" {
			t.Fatalf("missing gatekeeper node audit fields: %+v", event)
		}
	}
}

func TestGatekeeperAuditRecordsConflictAndBacktestMetrics(t *testing.T) {
	proposal := repository.RuleProposal{ProposalID: "prop_gate", BeforeRuleJSON: `{"rule":"hold"}`, AfterRuleJSON: `{"rule":"hold","auto_trade":true}`, SampleCount: 2}
	audit := gatekeeperAuditForProposal(proposal, model.AuditApproved, "v3.0")
	if !audit.ViolatesFundamentalRule || !audit.HasRuleConflict || audit.AllowApply {
		t.Fatalf("expected gatekeeper blocks: %+v", audit)
	}
	if !strings.Contains(audit.AuditReason, "FundamentalRuleCheck: failed") || !strings.Contains(audit.AuditReason, "ConflictCheck: failed") || !strings.Contains(audit.BacktestMetricsJSON, `"passed":false`) {
		t.Fatalf("expected detailed gatekeeper reason and metrics: %+v", audit)
	}
}

func hasAuditNode(events []model.AuditEvent, nodeName string) bool {
	for _, event := range events {
		if event.NodeName == nodeName {
			return true
		}
	}
	return false
}
