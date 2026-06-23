import type { PageResult } from '../types/api'
import type {
  RuleProposal,
  RuleProposalConfirmRequest,
  RuleProposalConfirmResponse,
  RuleEffectValidation,
  RuleEffectTracking,
  RuleEffectValidationRefreshRequest,
  RuleVersion,
  SOPAddendumProposalRequest,
  SOPAddendumProposalResponse,
} from '../types/rule'
import { apiRequest } from './client'

export function getCurrentRule() {
  return apiRequest<RuleVersion>('/api/v1/rules/current')
}

export function listRuleProposals() {
  return apiRequest<PageResult<RuleProposal>>('/api/v1/rule-proposals')
}

export function confirmRuleProposal(proposalId: string, body: RuleProposalConfirmRequest) {
  return apiRequest<RuleProposalConfirmResponse>(`/api/v1/rule-proposals/${proposalId}/confirm`, {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export function finalConfirmRuleProposal(proposalId: string, body: RuleProposalConfirmRequest) {
  return apiRequest<RuleProposalConfirmResponse>(`/api/v1/rule-proposals/${proposalId}/final-confirm`, {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export function createSOPAddendumProposal(body: SOPAddendumProposalRequest) {
  return apiRequest<SOPAddendumProposalResponse>('/api/v1/rule-proposals/sop-addendum', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export function getRuleEffectValidation(proposalId: string) {
  return apiRequest<RuleEffectValidation>(`/api/v1/rule-proposals/${proposalId}/effect-validation`)
}

export function refreshRuleEffectValidation(proposalId: string, body: RuleEffectValidationRefreshRequest) {
  return apiRequest<RuleEffectValidation>(`/api/v1/rule-proposals/${proposalId}/effect-validation`, {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export function listRuleEffectTracking(ruleVersion?: string) {
  const query = ruleVersion ? `?rule_version=${encodeURIComponent(ruleVersion)}` : ''
  return apiRequest<PageResult<RuleEffectTracking>>(`/api/v1/rule-effect-tracking${query}`)
}
