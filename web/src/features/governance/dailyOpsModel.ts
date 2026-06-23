import type { DailyAutoRunStatus } from '../../types/dailyAutoRun'
import type { DailyDisciplineReport, DailyDisciplineReportStatus } from '../../types/dailyDisciplineReport'
import type { OpsAction, OpsMetric, OpsTone } from './modelTypes'

export interface DailyOpsModel {
  overallLabel: string
  overallTone: OpsTone
  metrics: OpsMetric[]
  nextActions: OpsAction[]
  latestReport?: DailyDisciplineReport
  safetyNotes: string[]
}

const reportStatusText: Record<DailyDisciplineReportStatus, string> = {
  not_started: '未开始',
  running: '运行中',
  success: '成功',
  degraded: '降级完成',
  failed: '失败',
  insufficient_data: '数据不足',
}

const autoRunStatusText: Record<string, string> = {
  disabled: '关闭',
  scheduled: '已计划',
  running: '运行中',
  success: '成功',
  degraded: '部分成功',
  failed: '失败',
  unknown: '未知状态',
}

export function buildDailyOpsModel({ reports, autoRun }: { reports?: DailyDisciplineReport[]; autoRun?: DailyAutoRunStatus }): DailyOpsModel {
  const sortedReports = [...(reports ?? [])].sort((a, b) => b.local_date.localeCompare(a.local_date))
  const latestReport = sortedReports[0]
  const reportNeedsCheck = latestReport ? ['failed', 'degraded', 'insufficient_data'].includes(latestReport.status) : true
  const autoRunNeedsCheck = autoRun ? ['failed', 'degraded'].includes(autoRun.status) || !autoRun.enabled : false

  return {
    overallLabel: reportNeedsCheck || autoRunNeedsCheck ? '每日纪律与自动运行需要检查' : '每日纪律与自动运行可用',
    overallTone: reportNeedsCheck || autoRunNeedsCheck ? 'warning' : 'success',
    metrics: [
      { label: '最新报告', value: latestReport ? reportStatusText[latestReport.status] : '暂无', detail: latestReport?.local_date, tone: reportNeedsCheck ? 'warning' : 'success' },
      { label: '证据覆盖', value: latestReport ? `${latestReport.evidence.evidence_count} 条` : '暂无' },
      { label: '自动运行', value: autoRun ? autoRunStatusText[autoRun.status] ?? '未知状态' : '暂无', tone: autoRunNeedsCheck ? 'warning' : 'success' },
      { label: '执行范围', value: autoRun?.scope === 'holdings' ? '本地当前持仓' : autoRun?.scope ?? latestReport?.scope ?? '暂无' },
    ],
    nextActions: buildDailyActions(latestReport, autoRun),
    latestReport,
    safetyNotes: ['每日纪律和自动运行只记录本地刷新、评估与审计结果；所有投资动作仍需人工复核。'],
  }
}

function buildDailyActions(report?: DailyDisciplineReport, autoRun?: DailyAutoRunStatus): OpsAction[] {
  const actions: OpsAction[] = []
  if (report?.status === 'insufficient_data' || report?.missing_categories?.some((item) => ['holdings', 'portfolio', 'account'].includes(item)) || autoRun?.failure_code === 'missing_prerequisites') {
    actions.push({ label: '补齐本地持仓', detail: '先完成账户、组合和持仓前提。', href: '/positions' })
  }
  if (autoRun?.latest_audit_link || report?.audit_link) {
    actions.push({ label: '查看审计记录', detail: '查看每日运行或报告生成的审计线索。', href: autoRun?.latest_audit_link ?? report?.audit_link ?? '/audit' })
  }
  if (report) {
    actions.push({ label: '查看报告详情', detail: '回看纪律报告里的证据、趋势和缺口。', href: `/daily-discipline/reports/${encodeURIComponent(report.report_id)}` })
  }
  actions.push({ label: '查看每日自动运行', detail: '确认本地自动运行是否启用和健康。', href: '/daily-auto-run' })
  return actions
}
