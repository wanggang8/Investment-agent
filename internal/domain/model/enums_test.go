package model

import "testing"

func TestEnumValidation(t *testing.T) {
	cases := []struct {
		name  string
		valid bool
	}{
		{"dashboard", DashboardNormal.Valid()},
		{"workflow", WorkflowCompleted.Valid()},
		{"position", PositionSellOnly.Valid()},
		{"verification", VerificationSatisfied.Valid()},
		{"confirmation status", ConfirmationExecutedManually.Valid()},
		{"confirmation type", ConfirmationTypeWatch.Valid()},
		{"verdict", VerdictFrozenWatch.Valid()},
		{"proposal", ProposalPendingFinalConfirm.Valid()},
		{"audit result", AuditNeedsUserReview.Valid()},
		{"audit action", AuditActionGenerateDecision.Valid()},
		{"audit status", AuditStatusFailed.Valid()},
		{"liquidity", LiquidityDanger.Valid()},
		{"sentiment", SentimentExtreme.Valid()},
		{"precision", PrecisionUnavailable.Valid()},
		{"source", SourceLevelC.Valid()},
		{"evidence role", EvidenceBackground.Valid()},
		{"root cause tag", RootCauseEvidenceMissed.Valid()},
	}
	for _, tc := range cases {
		if !tc.valid {
			t.Fatalf("%s should be valid", tc.name)
		}
	}

	if DashboardState("bad").Valid() || FinalVerdictStatus("bad").Valid() || SourceLevel("X").Valid() || RootCauseTag("bad").Valid() {
		t.Fatal("invalid enum accepted")
	}
}

func TestSourceLevelPolicy(t *testing.T) {
	if !SourceLevelA.FormalAllowed() || SourceLevelC.FormalAllowed() {
		t.Fatal("source formal policy mismatch")
	}
	if !SourceLevelS.HighGrade() || SourceLevelB.HighGrade() {
		t.Fatal("source high grade policy mismatch")
	}
}
