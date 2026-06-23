import type { KnowledgeReadiness } from '../types/knowledgeReadiness'
import { apiRequest } from './client'

export function getKnowledgeReadiness(symbol?: string) {
  const params = new URLSearchParams()
  if (symbol) params.set('symbol', symbol)
  const suffix = params.toString() ? `?${params.toString()}` : ''
  return apiRequest<KnowledgeReadiness>(`/api/v1/knowledge-readiness${suffix}`)
}
