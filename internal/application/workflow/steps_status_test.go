package workflow

import (
	"testing"

	"investment-agent/internal/domain/model"
)

func TestBuildDecisionRecordWorkflowStatus(t *testing.T) {
	cases := []struct {
		name   string
		errors []string
		want   model.WorkflowStatus
	}{
		{name: "completed without errors", errors: nil, want: model.WorkflowCompleted},
		{name: "degraded for analyst unavailable", errors: []string{ErrCodeAnalystUnavailable}, want: model.WorkflowDegraded},
		{name: "degraded for vector index unavailable", errors: []string{ErrCodeVectorIndexUnavailable}, want: model.WorkflowDegraded},
		{name: "failed for evidence not found", errors: []string{ErrCodeEvidenceNotFound}, want: model.WorkflowFailed},
		{name: "failed wins over degraded", errors: []string{ErrCodeAnalystUnavailable, ErrCodeDataRequired}, want: model.WorkflowFailed},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			record := buildDecisionRecord(WorkflowContext{DecisionID: "decision_status", RequestID: "req_status", WorkflowType: WorkflowDailyDiscipline, Errors: tc.errors, RuleVerdict: model.RuleVerdict{Status: model.VerdictHold, Text: "持有"}})
			if record.WorkflowStatus != string(tc.want) {
				t.Fatalf("workflow_status=%s, want %s", record.WorkflowStatus, tc.want)
			}
		})
	}
}
