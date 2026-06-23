package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/domain/repository"
)

func (a *App) GetRuleEffectValidation(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	validation, err := a.QuerySvc.LatestRuleEffectValidationByProposal(r.Context(), r.PathValue("proposal_id"))
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, ruleEffectValidationDTO(validation))
}

func (a *App) RefreshRuleEffectValidation(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	proposalID := r.PathValue("proposal_id")
	var req dto.RuleEffectValidationRefreshRequest
	if err := decodeOptionalJSON(r, &req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	validation, err := a.RuleEffectSvc.EvaluateProposalFromLocalFacts(r.Context(), requestID, proposalID, defaultString(req.SampleWindow, "manual"))
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, ruleEffectValidationDTO(validation))
}

func (a *App) ListRuleEffectTracking(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	filter := repository.RuleEffectTrackingFilter{AppliedRuleVersion: strings.TrimSpace(r.URL.Query().Get("rule_version")), ProposalID: strings.TrimSpace(r.URL.Query().Get("proposal_id")), Period: strings.TrimSpace(r.URL.Query().Get("period"))}
	items, err := a.QuerySvc.ListRuleEffectTracking(r.Context(), filter)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	out := make([]dto.RuleEffectTrackingDTO, 0, len(items))
	for _, item := range items {
		out = append(out, ruleEffectTrackingDTO(item))
	}
	writeOK(w, requestID, dto.PageResult[dto.RuleEffectTrackingDTO]{Items: out, Total: len(out)})
}

func ruleEffectValidationDTO(item repository.RuleEffectValidation) dto.RuleEffectValidationDTO {
	return dto.RuleEffectValidationDTO{ValidationID: item.ValidationID, ProposalID: item.ProposalID, CandidateRuleVersion: item.CandidateRuleVersion, ValidationStatus: string(item.ValidationStatus), SampleCount: item.SampleCount, SampleWindow: item.SampleWindow, RepresentativenessStatus: string(item.RepresentativenessStatus), OverfitRisk: string(item.OverfitRisk), ReplayResult: string(item.ReplayResult), GuardrailDecision: string(item.GuardrailDecision), SourceExplanation: jsonValue(item.SourceExplanationJSON), Metrics: jsonValue(item.MetricsJSON), RiskNotes: jsonValue(item.RiskNotesJSON), RelatedErrorCases: jsonValue(item.RelatedErrorCasesJSON), RelatedDecisionIDs: jsonValue(item.RelatedDecisionIDsJSON), RelatedRiskAlertIDs: jsonValue(item.RelatedRiskAlertIDsJSON), RelatedAuditEventIDs: jsonValue(item.RelatedAuditEventIDsJSON), ValidationLink: "/rule-effect-validations/" + item.ValidationID, SafetyNote: item.SafetyNote, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt}
}

func ruleEffectTrackingDTO(item repository.RuleEffectTracking) dto.RuleEffectTrackingDTO {
	return dto.RuleEffectTrackingDTO{TrackingID: item.TrackingID, AppliedRuleVersion: item.AppliedRuleVersion, ProposalID: item.ProposalID, Period: item.Period, HitCount: item.HitCount, MisjudgmentCount: item.MisjudgmentCount, MissingEvidenceCount: item.MissingEvidenceCount, DegradedCount: item.DegradedCount, RiskAlertCount: item.RiskAlertCount, TrendDirection: string(item.TrendDirection), Metrics: jsonValue(item.MetricsJSON), RelatedProposalIDs: jsonValue(item.RelatedProposalIDsJSON), RelatedAuditEventIDs: jsonValue(item.RelatedAuditEventIDsJSON), RelatedRiskAlertIDs: jsonValue(item.RelatedRiskAlertIDsJSON), SafetyNote: item.SafetyNote, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt}
}

func jsonValue(raw string) any {
	if raw == "" {
		return nil
	}
	var out any
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return raw
	}
	return out
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}
