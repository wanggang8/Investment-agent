package handler

import (
	"net/http"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/pkg/apperr"
)

// GetCurrentRule 返回当前 active 规则版本。
func (a *App) GetCurrentRule(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	rule, err := a.QuerySvc.ActiveRuleVersion(r.Context())
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, dto.RuleVersionDTO{RuleVersion: rule.RuleVersion, Status: rule.Status, Rules: parseJSONAny(rule.RulesJSON), EffectiveAt: rule.EffectiveAt, CreatedAt: rule.CreatedAt})
}

// ListRuleProposals 返回规则提案列表。
func (a *App) ListRuleProposals(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	proposals, err := a.QuerySvc.ListRuleProposals(r.Context())
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	items := make([]dto.RuleProposalDTO, 0, len(proposals))
	for _, proposal := range proposals {
		item := dto.RuleProposalDTO{ProposalID: proposal.ProposalID, ProposalType: proposal.ProposalType, Status: proposal.Status, SourceErrorCaseID: proposal.SourceErrorCaseID, Title: proposal.Title, ProposalVersion: proposal.ProposalVersion, BeforeRule: parseJSONAny(proposal.BeforeRuleJSON), AfterRule: parseJSONAny(proposal.AfterRuleJSON), Reason: proposal.Reason, ImpactScope: parseJSONAny(proposal.ImpactScopeJSON), RiskNotes: parseJSONAny(proposal.RiskNotesJSON), AuditResult: proposal.AuditResult, AuditSummary: proposal.AuditReason, SampleCount: proposal.SampleCount, CreatedAt: proposal.CreatedAt}
		validation, err := a.QuerySvc.LatestRuleEffectValidationByProposal(r.Context(), proposal.ProposalID)
		if err != nil && !apperr.IsCode(err, apperr.CodeNotFound) {
			WriteHandlerError(w, requestID, err)
			return
		}
		if err == nil {
			item.EffectValidation = &dto.RuleEffectValidationDTO{ValidationID: validation.ValidationID, ValidationStatus: string(validation.ValidationStatus), SampleCount: validation.SampleCount, SampleWindow: validation.SampleWindow, RepresentativenessStatus: string(validation.RepresentativenessStatus), OverfitRisk: string(validation.OverfitRisk), ReplayResult: string(validation.ReplayResult), GuardrailDecision: string(validation.GuardrailDecision), ValidationLink: "/rule-effect-validations/" + validation.ValidationID, SafetyNote: validation.SafetyNote}
		}
		items = append(items, item)
	}
	writeOK(w, requestID, dto.PageResult[dto.RuleProposalDTO]{Items: items, Total: len(items)})
}

func (a *App) CreateSOPAddendumProposal(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	var req dto.SOPAddendumProposalRequest
	if err := decodeJSON(r, &req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	resp, err := a.RuleSvc.GenerateSOPAddendumProposal(r.Context(), requestID, req)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, resp)
}

// ConfirmRuleProposal 将用户确认后的规则提案送入守门人审计，样本不足时拒绝。
func (a *App) ConfirmRuleProposal(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	proposalID := r.PathValue("proposal_id")
	var req dto.RuleProposalConfirmRequest
	if err := decodeOptionalJSON(r, &req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	resp, err := a.RuleSvc.ConfirmProposal(r.Context(), requestID, proposalID, req, false)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, resp)
}

// FinalConfirmRuleProposal 最终确认后才应用正式规则版本。
func (a *App) FinalConfirmRuleProposal(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	proposalID := r.PathValue("proposal_id")
	var req dto.RuleProposalConfirmRequest
	if err := decodeOptionalJSON(r, &req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	resp, err := a.RuleSvc.ConfirmProposal(r.Context(), requestID, proposalID, req, true)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, resp)
}
