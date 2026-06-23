import type { DashboardTodayResponse } from '../types/dashboard'
import { apiRequest } from './client'

export function getDashboardToday() {
  return apiRequest<DashboardTodayResponse>('/api/v1/dashboard/today')
}
