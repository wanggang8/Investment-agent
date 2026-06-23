package service

import (
	"context"
	"encoding/json"
	"fmt"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/clock"
	"investment-agent/internal/pkg/idgen"
)

const (
	ruleEffectValidationSafetyNote = "规则效果验证只用于本地规则治理，不会自动应用规则或执行交易。"
	ruleEffectTrackingSafetyNote   = "应用后追踪只读展示，不会自动回滚规则或执行交易。"
)

// RuleEffectValidationService evaluates rule proposals and tracks applied rule effects.
type RuleEffectValidationService struct {
	tx  repository.Transactor
	clk clock.Clock
	ids idgen.Generator
}

// RuleEffectEvaluationInput contains local facts used by P36 validation.
type RuleEffectEvaluationInput struct {
	RequestID               string
	ProposalID              string
	CandidateRuleVersion    string
	SampleWindow            string
	SampleCount             int
	SourceCaseCount         int
	HitCount                int
	MisjudgmentCount        int
	MissingEvidenceCount    int
	DegradedCount           int
	RiskAlertCount          int
	ConflictingOutcomeCount int
	RelatedErrorCaseIDs     []string
	RelatedDecisionIDs      []string
	RelatedRiskAlertIDs     []string
}

// RuleEffectTrackingInput contains post-application tracking facts.
type RuleEffectTrackingInput struct {
	RequestID            string
	AppliedRuleVersion   string
	ProposalID           string
	Period               string
	HitCount             int
	MisjudgmentCount     int
	MissingEvidenceCount int
	DegradedCount        int
	RiskAlertCount       int
}

// NewRuleEffectValidationService creates a P36 service.
func NewRuleEffectValidationService(tx repository.Transactor) *RuleEffectValidationService {
	return &RuleEffectValidationService{tx: tx, clk: clock.SystemClock{}, ids: idgen.NewGenerator()}
}

// EvaluateProposal persists validation and an audit event without applying rules.
func (s *RuleEffectValidationService) EvaluateProposal(ctx context.Context, input RuleEffectEvaluationInput) (repository.RuleEffectValidation, error) {
	validation := s.buildValidation(input)
	auditID := s.ids.New("audit")
	validation.RelatedAuditEventIDsJSON = mustJSON([]string{auditID})
	if err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		if _, err := repos.RuleRepo.GetRuleProposal(ctx, input.ProposalID); err != nil {
			return err
		}
		if err := repos.RuleEffectRepo.SaveRuleEffectValidation(ctx, validation); err != nil {
			return err
		}
		return repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: auditID, RequestID: input.RequestID, Actor: string(model.AuditActorSystem), Action: string(model.AuditActionRunLocalTask), Status: string(model.AuditStatusSuccess), ProposalID: input.ProposalID, OutputRefType: "rule_effect_validation", OutputRef: validation.ValidationID, CreatedAt: validation.CreatedAt})
	}); err != nil {
		return repository.RuleEffectValidation{}, err
	}
	return validation, nil
}

// EvaluateProposalFromLocalFacts builds validation only from local repositories.
func (s *RuleEffectValidationService) EvaluateProposalFromLocalFacts(ctx context.Context, requestID, proposalID, sampleWindow string) (repository.RuleEffectValidation, error) {
	var input RuleEffectEvaluationInput
	if err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		proposal, err := repos.RuleRepo.GetRuleProposal(ctx, proposalID)
		if err != nil {
			return err
		}
		decisions, err := repos.DecisionRepo.ListDecisionRecords(ctx)
		if err != nil {
			return err
		}
		errorCases, err := repos.DecisionRepo.ListErrorCases(ctx)
		if err != nil {
			return err
		}
		var riskAlerts []repository.RiskAlert
		if repos.RiskAlertRepo != nil {
			riskAlerts, err = repos.RiskAlertRepo.ListRiskAlerts(ctx, repository.RiskAlertFilter{SOPStatuses: []model.RiskSOPStatus{model.RiskSOPTriggered, model.RiskSOPActive, model.RiskSOPObserving, model.RiskSOPEscalated}})
			if err != nil {
				return err
			}
		}
		input = ruleEffectInputFromLocalFacts(requestID, proposal, sampleWindow, decisions, errorCases, riskAlerts)
		return nil
	}); err != nil {
		return repository.RuleEffectValidation{}, err
	}
	return s.EvaluateProposal(ctx, input)
}

// TrackAppliedRule persists read-only tracking and an audit event.
func (s *RuleEffectValidationService) TrackAppliedRule(ctx context.Context, input RuleEffectTrackingInput) (repository.RuleEffectTracking, error) {
	now := s.clk.NowRFC3339()
	auditID := s.ids.New("audit")
	tracking := repository.RuleEffectTracking{
		TrackingID:               s.ids.New("rule_effect_tracking"),
		AppliedRuleVersion:       input.AppliedRuleVersion,
		ProposalID:               input.ProposalID,
		Period:                   input.Period,
		HitCount:                 input.HitCount,
		MisjudgmentCount:         input.MisjudgmentCount,
		MissingEvidenceCount:     input.MissingEvidenceCount,
		DegradedCount:            input.DegradedCount,
		RiskAlertCount:           input.RiskAlertCount,
		TrendDirection:           trendFromMetrics(input.MisjudgmentCount, input.MissingEvidenceCount, input.DegradedCount, input.RiskAlertCount),
		MetricsJSON:              mustJSON(map[string]int{"hit_count": input.HitCount, "misjudgment_count": input.MisjudgmentCount, "missing_evidence_count": input.MissingEvidenceCount, "degraded_count": input.DegradedCount, "risk_alert_count": input.RiskAlertCount}),
		RelatedProposalIDsJSON:   mustJSON(nonEmptyStrings(input.ProposalID)),
		RelatedAuditEventIDsJSON: mustJSON([]string{auditID}),
		SafetyNote:               ruleEffectTrackingSafetyNote,
		CreatedAt:                now,
		UpdatedAt:                now,
	}
	if err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		if err := repos.RuleEffectRepo.SaveRuleEffectTracking(ctx, tracking); err != nil {
			return err
		}
		return repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: auditID, RequestID: input.RequestID, Actor: string(model.AuditActorSystem), Action: string(model.AuditActionRunLocalTask), Status: string(model.AuditStatusSuccess), ProposalID: input.ProposalID, OutputRefType: "rule_effect_tracking", OutputRef: tracking.TrackingID, CreatedAt: now})
	}); err != nil {
		return repository.RuleEffectTracking{}, err
	}
	return tracking, nil
}

func (s *RuleEffectValidationService) buildValidation(input RuleEffectEvaluationInput) repository.RuleEffectValidation {
	now := s.clk.NowRFC3339()
	status := model.RuleEffectValidationPassed
	representativeness := model.RuleEffectValidationPassed
	overfit := model.RuleEffectOverfitLow
	replay := model.RuleEffectReplayPassed
	guardrail := model.RuleEffectGuardrailPassed
	riskNotes := []string{}

	if input.SampleCount < 3 {
		status = model.RuleEffectValidationInsufficient
		representativeness = model.RuleEffectValidationNeedsMoreSamples
		guardrail = model.RuleEffectGuardrailRejected
		riskNotes = append(riskNotes, "样本数量不足")
	}
	if input.SourceCaseCount <= 1 || input.ConflictingOutcomeCount > 0 {
		overfit = model.RuleEffectOverfitHigh
		guardrail = model.RuleEffectGuardrailRejected
		riskNotes = append(riskNotes, "样本过窄或存在冲突结果")
	}
	if input.MisjudgmentCount > 0 || input.MissingEvidenceCount > 0 || input.DegradedCount > 0 || input.RiskAlertCount > 0 {
		replay = model.RuleEffectReplayFailed
		if status == model.RuleEffectValidationPassed {
			status = model.RuleEffectValidationFailed
		}
		guardrail = model.RuleEffectGuardrailRejected
		riskNotes = append(riskNotes, "历史回放指标不利")
	}
	if status == model.RuleEffectValidationPassed && guardrail == model.RuleEffectGuardrailPassed && len(riskNotes) == 0 {
		riskNotes = append(riskNotes, "本地样本暂未发现不利信号")
	}

	return repository.RuleEffectValidation{
		ValidationID:             s.ids.New("rule_effect_validation"),
		ProposalID:               input.ProposalID,
		CandidateRuleVersion:     input.CandidateRuleVersion,
		ValidationStatus:         status,
		SampleCount:              input.SampleCount,
		SampleWindow:             input.SampleWindow,
		RepresentativenessStatus: representativeness,
		OverfitRisk:              overfit,
		ReplayResult:             replay,
		GuardrailDecision:        guardrail,
		SourceExplanationJSON:    mustJSON(map[string]any{"source_case_count": input.SourceCaseCount, "related_error_case_ids": input.RelatedErrorCaseIDs, "related_decision_ids": input.RelatedDecisionIDs, "related_risk_alert_ids": input.RelatedRiskAlertIDs}),
		MetricsJSON:              mustJSON(map[string]int{"hit_count": input.HitCount, "misjudgment_count": input.MisjudgmentCount, "missing_evidence_count": input.MissingEvidenceCount, "degraded_count": input.DegradedCount, "risk_alert_count": input.RiskAlertCount}),
		RiskNotesJSON:            mustJSON(riskNotes),
		RelatedErrorCasesJSON:    mustJSON(input.RelatedErrorCaseIDs),
		RelatedDecisionIDsJSON:   mustJSON(input.RelatedDecisionIDs),
		RelatedRiskAlertIDsJSON:  mustJSON(input.RelatedRiskAlertIDs),
		SafetyNote:               ruleEffectValidationSafetyNote,
		CreatedAt:                now,
		UpdatedAt:                now,
	}
}

func ruleEffectInputFromLocalFacts(requestID string, proposal repository.RuleProposal, sampleWindow string, decisions []repository.DecisionRecord, errorCases []repository.ErrorCase, riskAlerts []repository.RiskAlert) RuleEffectEvaluationInput {
	decisionIDs := make([]string, 0, len(decisions))
	missingEvidence, degraded := 0, 0
	for _, decision := range decisions {
		if decision.DecisionID != "" {
			decisionIDs = append(decisionIDs, decision.DecisionID)
		}
		if decision.SourceVerificationStatus == string(model.VerificationFailed) || decision.SourceVerificationStatus == string(model.VerificationBackgroundOnly) || decision.FinalVerdictStatus == string(model.VerdictInsufficientData) {
			missingEvidence++
		}
		if decision.WorkflowStatus == string(model.WorkflowDegraded) {
			degraded++
		}
	}
	errorCaseIDs := make([]string, 0, len(errorCases))
	for _, item := range errorCases {
		if item.ErrorCaseID != "" {
			errorCaseIDs = append(errorCaseIDs, item.ErrorCaseID)
		}
	}
	riskAlertIDs := make([]string, 0, len(riskAlerts))
	for _, item := range riskAlerts {
		if item.AlertID != "" {
			riskAlertIDs = append(riskAlertIDs, item.AlertID)
		}
	}
	return RuleEffectEvaluationInput{
		RequestID:            requestID,
		ProposalID:           proposal.ProposalID,
		CandidateRuleVersion: proposal.ProposalVersion,
		SampleWindow:         defaultRuleEffectSampleWindow(sampleWindow),
		SampleCount:          proposal.SampleCount,
		SourceCaseCount:      len(errorCases),
		HitCount:             triggeredRulesInDecisions(decisions),
		MisjudgmentCount:     len(errorCases),
		MissingEvidenceCount: missingEvidence,
		DegradedCount:        degraded,
		RiskAlertCount:       len(riskAlerts),
		RelatedErrorCaseIDs:  errorCaseIDs,
		RelatedDecisionIDs:   decisionIDs,
		RelatedRiskAlertIDs:  riskAlertIDs,
	}
}

func defaultRuleEffectSampleWindow(value string) string {
	if value == "" {
		return "local_facts"
	}
	return value
}

func triggeredRulesInDecisions(decisions []repository.DecisionRecord) int {
	count := 0
	for _, decision := range decisions {
		var values []any
		if err := json.Unmarshal([]byte(decision.TriggeredRulesJSON), &values); err == nil {
			count += len(values)
		} else if decision.TriggeredRulesJSON != "" && decision.TriggeredRulesJSON != "[]" {
			count++
		}
	}
	return count
}

func trendFromMetrics(misjudgment, missingEvidence, degraded, riskAlerts int) model.RuleEffectTrendDirection {
	if misjudgment > 2 || missingEvidence > 2 || degraded > 1 || riskAlerts > 0 {
		return model.RuleEffectTrendWorsened
	}
	if misjudgment == 0 && missingEvidence <= 1 && degraded == 0 && riskAlerts == 0 {
		return model.RuleEffectTrendImproved
	}
	return model.RuleEffectTrendFlat
}

func nonEmptyStrings(values ...string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value != "" {
			out = append(out, value)
		}
	}
	return out
}

func mustJSON(value any) string {
	b, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf(`{"error":%q}`, err.Error())
	}
	return string(b)
}
