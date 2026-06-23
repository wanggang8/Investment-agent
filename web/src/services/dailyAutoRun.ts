import type { DailyAutoRunStatus } from '../types/dailyAutoRun'
import { apiRequest } from './client'

export function getDailyAutoRunStatus() {
  return apiRequest<DailyAutoRunStatus>('/api/v1/daily-auto-run/status')
}
