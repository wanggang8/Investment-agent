package repository

import (
	"context"

	"investment-agent/internal/domain/model"
)

// RuleEffectValidation 保存 P36 规则效果验证事实。
type RuleEffectValidation struct {
	ValidationID             string
	ProposalID               string
	CandidateRuleVersion     string
	ValidationStatus         model.RuleEffectValidationStatus
	SampleCount              int
	SampleWindow             string
	RepresentativenessStatus model.RuleEffectValidationStatus
	OverfitRisk              model.RuleEffectOverfitRisk
	ReplayResult             model.RuleEffectReplayResult
	GuardrailDecision        model.RuleEffectGuardrailDecision
	SourceExplanationJSON    string
	MetricsJSON              string
	RiskNotesJSON            string
	RelatedErrorCasesJSON    string
	RelatedDecisionIDsJSON   string
	RelatedRiskAlertIDsJSON  string
	RelatedAuditEventIDsJSON string
	SafetyNote               string
	CreatedAt                string
	UpdatedAt                string
}

// RuleEffectValidationFilter filters persisted validation facts.
type RuleEffectValidationFilter struct {
	ProposalID           string
	CandidateRuleVersion string
	Statuses             []model.RuleEffectValidationStatus
}

// RuleEffectTracking 保存已应用规则的后续效果追踪事实。
type RuleEffectTracking struct {
	TrackingID               string
	AppliedRuleVersion       string
	ProposalID               string
	Period                   string
	HitCount                 int
	MisjudgmentCount         int
	MissingEvidenceCount     int
	DegradedCount            int
	RiskAlertCount           int
	TrendDirection           model.RuleEffectTrendDirection
	MetricsJSON              string
	RelatedProposalIDsJSON   string
	RelatedAuditEventIDsJSON string
	RelatedRiskAlertIDsJSON  string
	SafetyNote               string
	CreatedAt                string
	UpdatedAt                string
}

// RuleEffectTrackingFilter filters applied rule tracking facts.
type RuleEffectTrackingFilter struct {
	AppliedRuleVersion string
	ProposalID         string
	Period             string
}

// RuleEffectRepository defines P36 rule effect validation storage.
type RuleEffectRepository interface {
	SaveRuleEffectValidation(ctx context.Context, validation RuleEffectValidation) error
	GetRuleEffectValidation(ctx context.Context, validationID string) (RuleEffectValidation, error)
	ListRuleEffectValidations(ctx context.Context, filter RuleEffectValidationFilter) ([]RuleEffectValidation, error)
	SaveRuleEffectTracking(ctx context.Context, tracking RuleEffectTracking) error
	GetRuleEffectTracking(ctx context.Context, trackingID string) (RuleEffectTracking, error)
	ListRuleEffectTracking(ctx context.Context, filter RuleEffectTrackingFilter) ([]RuleEffectTracking, error)
}
