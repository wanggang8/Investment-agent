package sqlite

import (
	"context"
	"database/sql"
	"strings"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

// RuleEffectRepository is the SQLite implementation for P36 validation and tracking facts.
type RuleEffectRepository struct{ db dbtx }

// NewRuleEffectRepository creates a rule effect repository.
func NewRuleEffectRepository(db *sql.DB) *RuleEffectRepository { return &RuleEffectRepository{db: db} }

func (r *RuleEffectRepository) SaveRuleEffectValidation(ctx context.Context, v repository.RuleEffectValidation) error {
	if err := validateRuleEffectValidation(v); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx, `INSERT INTO rule_effect_validations (validation_id,proposal_id,candidate_rule_version,validation_status,sample_count,sample_window,representativeness_status,overfit_risk,replay_result,guardrail_decision,source_explanation_json,metrics_json,risk_notes_json,related_error_cases_json,related_decision_ids_json,related_risk_alert_ids_json,related_audit_event_ids_json,safety_note,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, v.ValidationID, v.ProposalID, v.CandidateRuleVersion, string(v.ValidationStatus), v.SampleCount, v.SampleWindow, string(v.RepresentativenessStatus), string(v.OverfitRisk), string(v.ReplayResult), string(v.GuardrailDecision), nullString(v.SourceExplanationJSON), nullString(v.MetricsJSON), nullString(v.RiskNotesJSON), nullString(v.RelatedErrorCasesJSON), nullString(v.RelatedDecisionIDsJSON), nullString(v.RelatedRiskAlertIDsJSON), nullString(v.RelatedAuditEventIDsJSON), v.SafetyNote, v.CreatedAt, v.UpdatedAt)
	return apperr.FromRepositoryError(err)
}

func (r *RuleEffectRepository) GetRuleEffectValidation(ctx context.Context, validationID string) (repository.RuleEffectValidation, error) {
	row := r.db.QueryRowContext(ctx, ruleEffectValidationSelectSQL()+` WHERE validation_id=?`, validationID)
	return scanRuleEffectValidation(row)
}

func (r *RuleEffectRepository) ListRuleEffectValidations(ctx context.Context, filter repository.RuleEffectValidationFilter) ([]repository.RuleEffectValidation, error) {
	query := ruleEffectValidationSelectSQL()
	var clauses []string
	var args []any
	if strings.TrimSpace(filter.ProposalID) != "" {
		clauses = append(clauses, "proposal_id=?")
		args = append(args, strings.TrimSpace(filter.ProposalID))
	}
	if strings.TrimSpace(filter.CandidateRuleVersion) != "" {
		clauses = append(clauses, "candidate_rule_version=?")
		args = append(args, strings.TrimSpace(filter.CandidateRuleVersion))
	}
	if len(filter.Statuses) > 0 {
		placeholders := make([]string, 0, len(filter.Statuses))
		for _, status := range filter.Statuses {
			if !status.Valid() {
				return nil, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "invalid rule effect validation status")
			}
			placeholders = append(placeholders, "?")
			args = append(args, string(status))
		}
		clauses = append(clauses, "validation_status IN ("+strings.Join(placeholders, ",")+")")
	}
	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}
	query += " ORDER BY updated_at DESC, validation_id DESC"
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperr.FromRepositoryError(err)
	}
	defer rows.Close()
	var out []repository.RuleEffectValidation
	for rows.Next() {
		item, err := scanRuleEffectValidation(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, apperr.FromRepositoryError(rows.Err())
}

func (r *RuleEffectRepository) SaveRuleEffectTracking(ctx context.Context, t repository.RuleEffectTracking) error {
	if err := validateRuleEffectTracking(t); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx, `INSERT INTO rule_effect_tracking (tracking_id,applied_rule_version,proposal_id,period,hit_count,misjudgment_count,missing_evidence_count,degraded_count,risk_alert_count,trend_direction,metrics_json,related_proposal_ids_json,related_audit_event_ids_json,related_risk_alert_ids_json,safety_note,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, t.TrackingID, t.AppliedRuleVersion, nullString(t.ProposalID), t.Period, t.HitCount, t.MisjudgmentCount, t.MissingEvidenceCount, t.DegradedCount, t.RiskAlertCount, string(t.TrendDirection), nullString(t.MetricsJSON), nullString(t.RelatedProposalIDsJSON), nullString(t.RelatedAuditEventIDsJSON), nullString(t.RelatedRiskAlertIDsJSON), t.SafetyNote, t.CreatedAt, t.UpdatedAt)
	return apperr.FromRepositoryError(err)
}

func (r *RuleEffectRepository) GetRuleEffectTracking(ctx context.Context, trackingID string) (repository.RuleEffectTracking, error) {
	row := r.db.QueryRowContext(ctx, ruleEffectTrackingSelectSQL()+` WHERE tracking_id=?`, trackingID)
	return scanRuleEffectTracking(row)
}

func (r *RuleEffectRepository) ListRuleEffectTracking(ctx context.Context, filter repository.RuleEffectTrackingFilter) ([]repository.RuleEffectTracking, error) {
	query := ruleEffectTrackingSelectSQL()
	var clauses []string
	var args []any
	if strings.TrimSpace(filter.AppliedRuleVersion) != "" {
		clauses = append(clauses, "applied_rule_version=?")
		args = append(args, strings.TrimSpace(filter.AppliedRuleVersion))
	}
	if strings.TrimSpace(filter.ProposalID) != "" {
		clauses = append(clauses, "proposal_id=?")
		args = append(args, strings.TrimSpace(filter.ProposalID))
	}
	if strings.TrimSpace(filter.Period) != "" {
		clauses = append(clauses, "period=?")
		args = append(args, strings.TrimSpace(filter.Period))
	}
	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}
	query += " ORDER BY updated_at DESC, tracking_id DESC"
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperr.FromRepositoryError(err)
	}
	defer rows.Close()
	var out []repository.RuleEffectTracking
	for rows.Next() {
		item, err := scanRuleEffectTracking(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, apperr.FromRepositoryError(rows.Err())
}

func validateRuleEffectValidation(v repository.RuleEffectValidation) error {
	if strings.TrimSpace(v.ValidationID) == "" || strings.TrimSpace(v.ProposalID) == "" || strings.TrimSpace(v.CandidateRuleVersion) == "" || strings.TrimSpace(v.SampleWindow) == "" || strings.TrimSpace(v.SafetyNote) == "" {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "rule effect validation requires id, proposal, rule version, sample window and safety note")
	}
	if !v.ValidationStatus.Valid() || !v.RepresentativenessStatus.Valid() || !v.OverfitRisk.Valid() || !v.ReplayResult.Valid() || !v.GuardrailDecision.Valid() {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "invalid rule effect validation enum")
	}
	return nil
}

func validateRuleEffectTracking(t repository.RuleEffectTracking) error {
	if strings.TrimSpace(t.TrackingID) == "" || strings.TrimSpace(t.AppliedRuleVersion) == "" || strings.TrimSpace(t.Period) == "" || strings.TrimSpace(t.SafetyNote) == "" {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "rule effect tracking requires id, rule version, period and safety note")
	}
	if !t.TrendDirection.Valid() {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "invalid rule effect trend direction")
	}
	return nil
}

func ruleEffectValidationSelectSQL() string {
	return `SELECT validation_id,proposal_id,candidate_rule_version,validation_status,sample_count,sample_window,representativeness_status,overfit_risk,replay_result,guardrail_decision,COALESCE(source_explanation_json,''),COALESCE(metrics_json,''),COALESCE(risk_notes_json,''),COALESCE(related_error_cases_json,''),COALESCE(related_decision_ids_json,''),COALESCE(related_risk_alert_ids_json,''),COALESCE(related_audit_event_ids_json,''),safety_note,created_at,updated_at FROM rule_effect_validations`
}

func ruleEffectTrackingSelectSQL() string {
	return `SELECT tracking_id,applied_rule_version,COALESCE(proposal_id,''),period,hit_count,misjudgment_count,missing_evidence_count,degraded_count,risk_alert_count,trend_direction,COALESCE(metrics_json,''),COALESCE(related_proposal_ids_json,''),COALESCE(related_audit_event_ids_json,''),COALESCE(related_risk_alert_ids_json,''),safety_note,created_at,updated_at FROM rule_effect_tracking`
}

type ruleEffectScanner interface {
	Scan(dest ...any) error
}

func scanRuleEffectValidation(scanner ruleEffectScanner) (repository.RuleEffectValidation, error) {
	var out repository.RuleEffectValidation
	var status, representativeness, overfit, replay, guardrail string
	if err := scanner.Scan(&out.ValidationID, &out.ProposalID, &out.CandidateRuleVersion, &status, &out.SampleCount, &out.SampleWindow, &representativeness, &overfit, &replay, &guardrail, &out.SourceExplanationJSON, &out.MetricsJSON, &out.RiskNotesJSON, &out.RelatedErrorCasesJSON, &out.RelatedDecisionIDsJSON, &out.RelatedRiskAlertIDsJSON, &out.RelatedAuditEventIDsJSON, &out.SafetyNote, &out.CreatedAt, &out.UpdatedAt); err != nil {
		return repository.RuleEffectValidation{}, apperr.FromRepositoryError(err)
	}
	out.ValidationStatus = model.RuleEffectValidationStatus(status)
	out.RepresentativenessStatus = model.RuleEffectValidationStatus(representativeness)
	out.OverfitRisk = model.RuleEffectOverfitRisk(overfit)
	out.ReplayResult = model.RuleEffectReplayResult(replay)
	out.GuardrailDecision = model.RuleEffectGuardrailDecision(guardrail)
	return out, nil
}

func scanRuleEffectTracking(scanner ruleEffectScanner) (repository.RuleEffectTracking, error) {
	var out repository.RuleEffectTracking
	var trend string
	if err := scanner.Scan(&out.TrackingID, &out.AppliedRuleVersion, &out.ProposalID, &out.Period, &out.HitCount, &out.MisjudgmentCount, &out.MissingEvidenceCount, &out.DegradedCount, &out.RiskAlertCount, &trend, &out.MetricsJSON, &out.RelatedProposalIDsJSON, &out.RelatedAuditEventIDsJSON, &out.RelatedRiskAlertIDsJSON, &out.SafetyNote, &out.CreatedAt, &out.UpdatedAt); err != nil {
		return repository.RuleEffectTracking{}, apperr.FromRepositoryError(err)
	}
	out.TrendDirection = model.RuleEffectTrendDirection(trend)
	return out, nil
}
