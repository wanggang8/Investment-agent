import type { SourceHealthItem } from './market'
import type { RiskAlert } from './riskAlert'

export type DailyDisciplineReportStatus =
  | 'not_started'
  | 'running'
  | 'success'
  | 'degraded'
  | 'failed'
  | 'insufficient_data'

export interface DailyDisciplineReportEvidenceSummary {
  evidence_count: number
  independent_source_count: number
  high_grade_independent_source_count: number
}

export interface DailyDisciplineReportTrendSummary {
  success_count: number
  degraded_count: number
  failed_count: number
  insufficient_data_count: number
}

export interface DailyDisciplineReportP34Coverage {
  summary?: string
  missing_categories?: string[]
  source_health?: SourceHealthItem[]
}

export interface DailyDisciplineReport {
  report_id: string
  local_date: string
  scope: string
  status: DailyDisciplineReportStatus
  summary: string
  source_type?: string
  source_id?: string
  decision_id?: string
  decision_link?: string
  auto_run_link?: string
  audit_link?: string
  notification_link?: string
  failure_code?: string
  failure_reason?: string
  missing_action?: string
  missing_categories?: string[]
  final_verdict?: string
  verdict_status?: string
  evidence: DailyDisciplineReportEvidenceSummary
  p34_source_coverage?: DailyDisciplineReportP34Coverage
  risk_alerts?: RiskAlert[]
  trend: DailyDisciplineReportTrendSummary
  safety_note: string
  updated_at?: string
}

export interface DailyDisciplineReportListResponse {
  reports: DailyDisciplineReport[]
}
