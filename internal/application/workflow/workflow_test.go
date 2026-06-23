package workflow

import (
	"context"
	"testing"

	"investment-agent/internal/domain/model"
)

func TestNodeResultAuditValidation(t *testing.T) {
	result := NodeResult{
		Status: StatusSuccess,
		Audit: AuditFragment{
			Action:        "generate_decision",
			NodeName:      "StateSnapshotNode",
			NodeAction:    "load_state_snapshot",
			Status:        StatusSuccess,
			InputRefType:  "request",
			InputRef:      "req_1",
			OutputRefType: "snapshot",
			OutputRef:     "snap_1",
		},
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("expected valid audit fragment: %v", err)
	}

	failed := NodeResult{Status: StatusFailed, Audit: AuditFragment{Action: "generate_decision", NodeName: "StateSnapshotNode", NodeAction: "load_state_snapshot", Status: StatusFailed, InputRefType: "request", InputRef: "req_1"}}
	if err := failed.Validate(); err != ErrMissingErrorCode {
		t.Fatalf("expected missing error code, got %v", err)
	}
}

func TestAuditWriterAppendsFailedNodeAudit(t *testing.T) {
	writer := &MemoryAuditWriter{}
	ctx := WorkflowContext{RequestID: "req_1", WorkflowType: WorkflowDailyDiscipline, RuleVersion: "v3.0"}
	result := NodeResult{
		Status:    StatusFailed,
		ErrorCode: ErrCodeEvidenceNotFound,
		Audit: AuditFragment{
			Action:       "generate_decision",
			NodeName:     "EvidenceRetrievalNode",
			NodeAction:   "retrieve_evidence",
			Status:       StatusFailed,
			InputRefType: "symbol",
			InputRef:     "510300",
			ErrorCode:    ErrCodeEvidenceNotFound,
		},
	}
	if err := writer.Write(context.Background(), &ctx, result); err != nil {
		t.Fatalf("write audit: %v", err)
	}
	if len(ctx.AuditEvents) != 1 {
		t.Fatalf("expected one audit event, got %d", len(ctx.AuditEvents))
	}
	if ctx.AuditEvents[0].Status != model.AuditStatusFailed || ctx.AuditEvents[0].Action != model.AuditActionGenerateDecision {
		t.Fatalf("unexpected audit event: %+v", ctx.AuditEvents[0])
	}
}

func TestExpectedReturnPrecisionMapping(t *testing.T) {
	cases := []struct {
		name          string
		sampleCount   int
		want          model.PrecisionStatus
		wantScenarios bool
	}{
		{name: "available", sampleCount: 20, want: model.PrecisionAvailable, wantScenarios: true},
		{name: "insufficient", sampleCount: 5, want: model.PrecisionInsufficient, wantScenarios: true},
		{name: "unavailable", sampleCount: 4, want: model.PrecisionUnavailable, wantScenarios: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := BuildExpectedReturn(tc.sampleCount)
			if out.PrecisionStatus != tc.want {
				t.Fatalf("status=%s want=%s", out.PrecisionStatus, tc.want)
			}
			if (len(out.Scenarios) > 0) != tc.wantScenarios {
				t.Fatalf("unexpected scenarios: %+v", out.Scenarios)
			}
			if tc.want == model.PrecisionInsufficient && out.Scenarios[0].Probability != nil {
				t.Fatal("insufficient precision must not return exact probability")
			}
		})
	}
}
