package service

import (
	"context"
	"testing"
	"time"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/infrastructure/persistence/sqlite"
	"investment-agent/internal/pkg/clock"
	"investment-agent/internal/pkg/idgen"
)

func TestRuleEffectValidationServiceEvaluatesInsufficientAndHighOverfit(t *testing.T) {
	store, err := sqlite.Open(":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	ctx := context.Background()
	if err := sqlite.Migrate(ctx, store.DB); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if err := sqlite.NewRuleRepository(store.DB).SaveRuleProposal(ctx, repository.RuleProposal{ProposalID: "prop_1", ProposalType: "threshold", Status: string(model.ProposalPendingUserConfirm), Title: "调高证据门槛", ProposalVersion: "v3.1-proposal", SampleCount: 2, CreatedAt: "2026-06-16T09:00:00Z"}); err != nil {
		t.Fatalf("seed proposal: %v", err)
	}

	svc := NewRuleEffectValidationService(sqlite.NewTransactor(store.DB))
	svc.clk = clock.FixedClock{Time: time.Date(2026, 6, 16, 9, 30, 0, 0, time.UTC)}
	svc.ids = idgen.NewFixedGenerator(map[string][]string{"rule_effect_validation": {"val_1"}, "audit": {"audit_1"}})

	out, err := svc.EvaluateProposal(ctx, RuleEffectEvaluationInput{RequestID: "req_1", ProposalID: "prop_1", CandidateRuleVersion: "v3.1-proposal", SampleWindow: "2026-Q2", SampleCount: 2, SourceCaseCount: 1, ConflictingOutcomeCount: 1, MissingEvidenceCount: 2, RelatedErrorCaseIDs: []string{"err_1"}, RelatedDecisionIDs: []string{"dec_1"}, RelatedRiskAlertIDs: []string{"risk_1"}})
	if err != nil {
		t.Fatalf("evaluate proposal: %v", err)
	}
	if out.ValidationStatus != model.RuleEffectValidationInsufficient || out.OverfitRisk != model.RuleEffectOverfitHigh || out.GuardrailDecision != model.RuleEffectGuardrailRejected {
		t.Fatalf("unexpected validation: %+v", out)
	}
	if out.SafetyNote == "" || out.SourceExplanationJSON == "" || out.RiskNotesJSON == "" {
		t.Fatalf("expected traceable safety output: %+v", out)
	}

	persisted, err := sqlite.NewRuleEffectRepository(store.DB).GetRuleEffectValidation(ctx, "val_1")
	if err != nil {
		t.Fatalf("get validation: %v", err)
	}
	if persisted.ProposalID != "prop_1" || persisted.RelatedRiskAlertIDsJSON == "" || persisted.RelatedAuditEventIDsJSON != `["audit_1"]` {
		t.Fatalf("unexpected persisted validation: %+v", persisted)
	}
}

func TestRuleEffectValidationServiceEvaluatesPassedAndTracksAppliedRule(t *testing.T) {
	store, err := sqlite.Open(":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	ctx := context.Background()
	if err := sqlite.Migrate(ctx, store.DB); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if err := sqlite.NewRuleRepository(store.DB).SaveRuleProposal(ctx, repository.RuleProposal{ProposalID: "prop_2", ProposalType: "risk_rule", Status: string(model.ProposalPendingUserConfirm), Title: "补充数据降级规则", ProposalVersion: "v3.2-proposal", SampleCount: 5, AppliedRuleVersion: "v3.2", CreatedAt: "2026-06-16T09:00:00Z"}); err != nil {
		t.Fatalf("seed proposal: %v", err)
	}

	svc := NewRuleEffectValidationService(sqlite.NewTransactor(store.DB))
	svc.clk = clock.FixedClock{Time: time.Date(2026, 6, 16, 10, 0, 0, 0, time.UTC)}
	svc.ids = idgen.NewFixedGenerator(map[string][]string{"rule_effect_validation": {"val_2"}, "rule_effect_tracking": {"track_1"}, "audit": {"audit_2", "audit_3"}})

	out, err := svc.EvaluateProposal(ctx, RuleEffectEvaluationInput{RequestID: "req_2", ProposalID: "prop_2", CandidateRuleVersion: "v3.2-proposal", SampleWindow: "2026-Q2", SampleCount: 5, SourceCaseCount: 4, HitCount: 8, MisjudgmentCount: 0, MissingEvidenceCount: 0, DegradedCount: 0, RiskAlertCount: 0, RelatedDecisionIDs: []string{"dec_1", "dec_2"}})
	if err != nil {
		t.Fatalf("evaluate passed proposal: %v", err)
	}
	if out.ValidationStatus != model.RuleEffectValidationPassed || out.OverfitRisk != model.RuleEffectOverfitLow || out.GuardrailDecision != model.RuleEffectGuardrailPassed {
		t.Fatalf("unexpected passed validation: %+v", out)
	}

	tracking, err := svc.TrackAppliedRule(ctx, RuleEffectTrackingInput{RequestID: "req_3", AppliedRuleVersion: "v3.2", ProposalID: "prop_2", Period: "2026-Q3", HitCount: 12, MisjudgmentCount: 0, MissingEvidenceCount: 0, DegradedCount: 0, RiskAlertCount: 0})
	if err != nil {
		t.Fatalf("track applied rule: %v", err)
	}
	if tracking.TrendDirection != model.RuleEffectTrendImproved || tracking.SafetyNote == "" {
		t.Fatalf("unexpected tracking: %+v", tracking)
	}

	for _, table := range []string{"operation_confirmations", "position_transactions"} {
		var count int
		if err := store.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+table).Scan(&count); err != nil {
			t.Fatalf("count %s: %v", table, err)
		}
		if count != 0 {
			t.Fatalf("expected no trading mutation in %s, got %d", table, count)
		}
	}
}
