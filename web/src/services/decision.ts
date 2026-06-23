import type { PageResult } from '../types/api'
import type {
  ConfirmationRequest,
  ConfirmationResponse,
  ConsultDecisionRequest,
  DecisionDetailResponse,
  DecisionListItem,
} from '../types/decision'
import { apiRequest } from './client'

export function consultDecision(body: ConsultDecisionRequest) {
  return apiRequest<DecisionDetailResponse>('/api/v1/decisions/consult', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export function getDecision(decisionId: string) {
  return apiRequest<DecisionDetailResponse>(`/api/v1/decisions/${decisionId}`)
}

export function listDecisions(query = '') {
  return apiRequest<PageResult<DecisionListItem>>(`/api/v1/decisions${query}`)
}

export function createConfirmation(decisionId: string, body: ConfirmationRequest) {
  return apiRequest<ConfirmationResponse>(`/api/v1/decisions/${decisionId}/confirmations`, {
    method: 'POST',
    body: JSON.stringify(body),
  })
}
