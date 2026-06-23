import { Link } from 'react-router-dom'
import { riskSeverityText, riskSOPStatusText, riskTypeText, textOrRaw } from '../shared/mappers'
import type { DailyDisciplineReport, DailyDisciplineReportStatus } from '../types/dailyDisciplineReport'
import type { RiskAlert } from '../types/riskAlert'

export const dailyDisciplineReportStatusLabels: Record<DailyDisciplineReportStatus, string> = {
  not_started: '未开始',
  running: '运行中',
  success: '成功',
  degraded: '降级完成',
  failed: '失败',
  insufficient_data: '数据不足',
}

interface ReportLinksProps {
  report: Pick<DailyDisciplineReport, 'decision_link' | 'audit_link' | 'notification_link' | 'auto_run_link'>
}

export function DailyDisciplineReportLinks({ report }: ReportLinksProps) {
  return (
    <div className="action-row">
      {report.decision_link ? <a href={report.decision_link}>查看决策详情</a> : null}
      {report.audit_link ? <a href={report.audit_link}>查看审计详情</a> : null}
      {report.notification_link ? <a href={report.notification_link}>查看通知</a> : null}
      {report.auto_run_link ? <a href={report.auto_run_link}>查看每日自动运行</a> : null}
    </div>
  )
}

export function safeDailyDisciplineFailureMessage(report: Pick<DailyDisciplineReport, 'failure_code' | 'failure_reason'>) {
  if (!report.failure_reason && !report.failure_code) return ''
  switch (report.failure_code) {
    case 'missing_prerequisites':
      return '缺少本地账户或持仓。'
    case 'timeout':
      return '每日纪律报告生成超时。'
    default:
      return '每日纪律报告暂时无法生成，请查看审计记录。'
  }
}

export function needsPortfolioOnboarding(report: Pick<DailyDisciplineReport, 'failure_code' | 'missing_categories' | 'status'>) {
  return report.failure_code === 'missing_prerequisites' || report.missing_categories?.some((category) => ['holdings', 'portfolio', 'account'].includes(category)) || report.status === 'insufficient_data'
}

export function PortfolioOnboardingLink({ report }: { report: Pick<DailyDisciplineReport, 'failure_code' | 'missing_categories' | 'status'> }) {
  if (!needsPortfolioOnboarding(report)) return null
  return <a href="/positions">前往账户初始化</a>
}

export function RiskAlertSummaryList({ alerts, title = '风险预警' }: { alerts?: RiskAlert[]; title?: string }) {
  if (!alerts?.length) return null
  return (
    <div className="risk-alert-summary">
      <strong>{title}</strong>
      <ul>
        {alerts.map((alert) => (
          <li key={alert.alert_id}>
            <span>{textOrRaw(riskTypeText, alert.risk_type)} · {textOrRaw(riskSeverityText, alert.severity)} · {textOrRaw(riskSOPStatusText, alert.sop_status)}</span>
            <p>{alert.trigger_summary}</p>
            {alert.prohibited_actions?.length ? <small>禁止动作：{alert.prohibited_actions.join('、')}</small> : null}
            {alert.suggested_actions?.length ? <small>建议人工动作：{alert.suggested_actions.join('、')}</small> : null}
            {alert.link ? <a href={alert.link}>查看风险预警</a> : null}
          </li>
        ))}
      </ul>
    </div>
  )
}

interface ReportCardProps {
  report: DailyDisciplineReport
  showDetailLink?: boolean
}

export function DailyDisciplineReportCard({ report, showDetailLink = false }: ReportCardProps) {
  return (
    <article className="panel-card">
      <div className="row-between">
        <strong>{report.local_date}</strong>
        <span>{dailyDisciplineReportStatusLabels[report.status]}</span>
      </div>
      <p>{report.summary}</p>
      {safeDailyDisciplineFailureMessage(report) ? <p>{safeDailyDisciplineFailureMessage(report)}</p> : null}
      {report.missing_action ? <p>{report.missing_action}</p> : null}
      <PortfolioOnboardingLink report={report} />
      {report.safety_note ? <small>{report.safety_note}</small> : null}
      <DailyDisciplineReportLinks report={report} />
      {showDetailLink ? (
        <Link to={`/daily-discipline/reports/${encodeURIComponent(report.report_id)}`}>查看报告</Link>
      ) : null}
    </article>
  )
}
