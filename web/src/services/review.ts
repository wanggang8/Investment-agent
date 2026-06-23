import type { ReviewSummary } from '../types/review'
import { apiRequest } from './client'

export function getReviewSummary() {
  return apiRequest<ReviewSummary>('/api/v1/review/summary')
}
