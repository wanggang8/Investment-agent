import type { APIEnvelope } from '../types/api'
import type { DailyDisciplineReport, DailyDisciplineReportListResponse, DailyDisciplineReportStatus } from '../types/dailyDisciplineReport'
import { apiRequest } from './client'

export function getTodayDailyDisciplineReport(): Promise<APIEnvelope<DailyDisciplineReport>> {
  return apiRequest<DailyDisciplineReport>('/api/v1/daily-discipline/reports/today')
}

export function listDailyDisciplineReports(
  limit = 30,
  status?: DailyDisciplineReportStatus,
): Promise<APIEnvelope<DailyDisciplineReportListResponse>> {
  const params = new URLSearchParams({ limit: String(limit) })
  if (status) {
    params.set('status', status)
  }
  return apiRequest<DailyDisciplineReportListResponse>(`/api/v1/daily-discipline/reports?${params.toString()}`)
}

export function getDailyDisciplineReport(reportId: string): Promise<APIEnvelope<DailyDisciplineReport>> {
  return apiRequest<DailyDisciplineReport>(`/api/v1/daily-discipline/reports/${encodeURIComponent(reportId)}`)
}
