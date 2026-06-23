package workflow

import (
	"context"
	"testing"
	"time"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/pkg/clock"
	"investment-agent/internal/pkg/idgen"
)

func TestAuditEventIDsAreDeterministicWithInjectedGenerator(t *testing.T) {
	oldAuditGen := auditIDGen
	oldWorkflowClock := workflowClock
	defer func() {
		auditIDGen = oldAuditGen
		workflowClock = oldWorkflowClock
	}()
	auditIDGen = idgen.NewFixedGenerator(map[string][]string{"audit": {"audit_one", "audit_two", "audit_repo"}})
	workflowClock = clock.FixedClock{Time: time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)}

	wf := WorkflowContext{RequestID: "req_1", WorkflowType: WorkflowConsultation}
	result := NodeResult{Audit: AuditFragment{Action: "generate_decision", NodeName: "TestNode", NodeAction: "write", Status: StatusSuccess}}
	event1 := buildDomainAuditEvent(&wf, result)
	event2 := buildDomainAuditEvent(&wf, result)
	if event1.AuditEventID != "audit_one" || event2.AuditEventID != "audit_two" {
		t.Fatalf("unexpected audit ids: %s %s", event1.AuditEventID, event2.AuditEventID)
	}
	repoEvent := buildRepositoryAuditEvent(&wf, result)
	if repoEvent.AuditEventID != "audit_repo" || repoEvent.CreatedAt != "2026-05-29T04:00:00Z" {
		t.Fatalf("unexpected repository audit event: %+v", repoEvent)
	}
}

func TestWorkflowDecisionAndEvidenceIDsRespectInjectedGeneratorAndClock(t *testing.T) {
	oldWorkflowIDGen := workflowIDGen
	oldWorkflowClock := workflowClock
	defer func() {
		workflowIDGen = oldWorkflowIDGen
		workflowClock = oldWorkflowClock
	}()
	workflowIDGen = idgen.NewFixedGenerator(map[string][]string{
		"decision": {"decision_fixed"},
		"eref":     {"eref_one", "eref_two"},
	})
	workflowClock = clock.FixedClock{Time: time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)}

	wf := WorkflowContext{
		RequestID:         "req_fixed",
		WorkflowType:      WorkflowConsultation,
		Symbol:            "510300",
		UserQuestion:      "test",
		PortfolioSnapshot: model.PortfolioSnapshot{SnapshotID: "snap_1"},
		EvidenceSet:       model.EvidenceSet{Items: []model.Evidence{{EvidenceID: "e1", SummaryID: "summary_one", SourceLevel: model.SourceLevelA, Role: model.EvidenceFormal}, {EvidenceID: "e2", SummaryID: "summary_two", SourceLevel: model.SourceLevelB, Role: model.EvidenceFormal}}},
		RuleVerdict:       model.RuleVerdict{Status: model.VerdictHold},
	}
	decisionRecordStep(context.Background(), &wf, WorkflowDependencies{})
	if wf.DecisionID != "decision_fixed" {
		t.Fatalf("decision id = %s", wf.DecisionID)
	}
	refs := buildEvidenceRefs(wf)
	if len(refs) != 2 || refs[0].EvidenceRefID != "eref_one" || refs[1].SummaryID != "summary_two" {
		t.Fatalf("unexpected refs: %+v", refs)
	}
	if got := buildDecisionRecord(wf).CreatedAt; got != "2026-05-29T04:00:00Z" {
		t.Fatalf("created_at = %s", got)
	}
}
