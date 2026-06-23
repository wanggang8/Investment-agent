package sqlite

import (
	"context"
	"testing"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

func TestRuleEffectRepositoryValidationWriteReadAndList(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	repo := NewRuleEffectRepository(db)

	validation := repository.RuleEffectValidation{
		ValidationID:             "val_1",
		ProposalID:               "prop_1",
		CandidateRuleVersion:     "v3.1-proposal",
		ValidationStatus:         model.RuleEffectValidationInsufficient,
		SampleCount:              2,
		SampleWindow:             "2026-Q2",
		RepresentativenessStatus: model.RuleEffectValidationNeedsMoreSamples,
		OverfitRisk:              model.RuleEffectOverfitHigh,
		ReplayResult:             model.RuleEffectReplayFailed,
		GuardrailDecision:        model.RuleEffectGuardrailRejected,
		SourceExplanationJSON:    `{"error_cases":["err_1"]}`,
		MetricsJSON:              `{"misjudgment_rate_delta":0.12}`,
		RiskNotesJSON:            `["样本不足","过拟合风险高"]`,
		RelatedErrorCasesJSON:    `["err_1"]`,
		RelatedDecisionIDsJSON:   `["dec_1"]`,
		RelatedRiskAlertIDsJSON:  `["risk_1"]`,
		RelatedAuditEventIDsJSON: `["audit_1"]`,
		SafetyNote:               "规则效果验证只用于本地规则治理，不会自动应用规则或执行交易。",
		CreatedAt:                testTime,
		UpdatedAt:                testTime,
	}
	if err := repo.SaveRuleEffectValidation(ctx, validation); err != nil {
		t.Fatalf("save validation: %v", err)
	}

	got, err := repo.GetRuleEffectValidation(ctx, "val_1")
	if err != nil {
		t.Fatalf("get validation: %v", err)
	}
	if got.ValidationStatus != model.RuleEffectValidationInsufficient || got.OverfitRisk != model.RuleEffectOverfitHigh || got.GuardrailDecision != model.RuleEffectGuardrailRejected {
		t.Fatalf("unexpected validation: %+v", got)
	}

	items, err := repo.ListRuleEffectValidations(ctx, repository.RuleEffectValidationFilter{ProposalID: "prop_1"})
	if err != nil {
		t.Fatalf("list validations: %v", err)
	}
	if len(items) != 1 || items[0].ValidationID != "val_1" {
		t.Fatalf("unexpected validations: %+v", items)
	}
}

func TestRuleEffectRepositoryTrackingWriteReadAndNoTradingMutation(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	repo := NewRuleEffectRepository(db)

	tracking := repository.RuleEffectTracking{
		TrackingID:               "track_1",
		AppliedRuleVersion:       "v3.1",
		ProposalID:               "prop_1",
		Period:                   "2026-Q3",
		HitCount:                 8,
		MisjudgmentCount:         1,
		MissingEvidenceCount:     2,
		DegradedCount:            1,
		RiskAlertCount:           3,
		TrendDirection:           model.RuleEffectTrendWorsened,
		MetricsJSON:              `{"hit_rate":0.8}`,
		RelatedProposalIDsJSON:   `["prop_1"]`,
		RelatedAuditEventIDsJSON: `["audit_1"]`,
		RelatedRiskAlertIDsJSON:  `["risk_1"]`,
		SafetyNote:               "应用后追踪只读展示，不会自动回滚规则或执行交易。",
		CreatedAt:                testTime,
		UpdatedAt:                testTime,
	}
	if err := repo.SaveRuleEffectTracking(ctx, tracking); err != nil {
		t.Fatalf("save tracking: %v", err)
	}
	got, err := repo.GetRuleEffectTracking(ctx, "track_1")
	if err != nil {
		t.Fatalf("get tracking: %v", err)
	}
	if got.TrendDirection != model.RuleEffectTrendWorsened || got.RiskAlertCount != 3 {
		t.Fatalf("unexpected tracking: %+v", got)
	}
	items, err := repo.ListRuleEffectTracking(ctx, repository.RuleEffectTrackingFilter{AppliedRuleVersion: "v3.1"})
	if err != nil {
		t.Fatalf("list tracking: %v", err)
	}
	if len(items) != 1 || items[0].TrackingID != "track_1" {
		t.Fatalf("unexpected tracking list: %+v", items)
	}

	for _, table := range []string{"operation_confirmations", "position_transactions"} {
		var count int
		if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+table).Scan(&count); err != nil {
			t.Fatalf("count %s: %v", table, err)
		}
		if count != 0 {
			t.Fatalf("expected no trading mutation in %s, got %d", table, count)
		}
	}
}

func TestRuleEffectRepositoryClassifiesErrors(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	repo := NewRuleEffectRepository(db)

	if _, err := repo.GetRuleEffectValidation(ctx, "missing"); !apperr.IsCode(err, apperr.CodeNotFound) {
		t.Fatalf("expected validation not found, got %v", err)
	}
	bad := repository.RuleEffectValidation{ValidationID: "bad", ProposalID: "prop_1", CandidateRuleVersion: "v3.1-proposal", ValidationStatus: model.RuleEffectValidationStatus("invalid"), SampleWindow: "2026-Q2", RepresentativenessStatus: model.RuleEffectValidationPassed, OverfitRisk: model.RuleEffectOverfitLow, ReplayResult: model.RuleEffectReplayPassed, GuardrailDecision: model.RuleEffectGuardrailPassed, CreatedAt: testTime, UpdatedAt: testTime}
	if err := repo.SaveRuleEffectValidation(ctx, bad); !apperr.IsCode(err, apperr.CodeBadRequest) {
		t.Fatalf("expected validation bad request, got %v", err)
	}
}
