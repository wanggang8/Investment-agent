import type { DecisionLoopItem, DecisionLoopListFilter, DecisionLoopListResponse } from '../types/decisionLoop'
import { apiRequest } from './client'

export function listDecisionLoops(filter: DecisionLoopListFilter = {}) {
  const params = new URLSearchParams()
  if (filter.symbol) {
    params.set('symbol', filter.symbol)
  }
  if (filter.limit) {
    params.set('limit', String(filter.limit))
  }
  const qs = params.toString()
  return apiRequest<DecisionLoopListResponse>(`/api/v1/decision-loops${qs ? `?${qs}` : ''}`)
}

export function getDecisionLoop(decisionId: string) {
  return apiRequest<DecisionLoopItem>(`/api/v1/decision-loops/${encodeURIComponent(decisionId)}`)
}
