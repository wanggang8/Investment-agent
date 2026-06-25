import type { DailyDisciplineReport } from '../../types/dailyDisciplineReport'
import type { DashboardTodayResponse } from '../../types/dashboard'
import type { PortfolioCurrentResponse } from '../../types/portfolio'
import type { ReviewSummary } from '../../types/review'
import type { RiskAlert } from '../../types/riskAlert'
import type { RuleProposal } from '../../types/rule'
import type { PageErrorState } from '../../shared/utils'
import { formatCurrency, formatDateTime, formatPercent } from '../../shared/utils'

export type DailyTone = 'success' | 'warning' | 'danger' | 'degraded' | 'unknown'

export interface DailyAction {
  label: string
  href: string
  detail: string
  priority: 'blocking' | 'review' | 'follow_up'
}

export interface DailySignal {
  label: string
  value: string
  detail: string
  tone: DailyTone
  href?: string
}

export interface DailyWorkbenchModel {
  statusLabel: string
  statusTone: DailyTone
  verdictText: string
  trustSummary: string
  riskSummary: string
  updatedAtText: string
  prohibitedActions: string[]
  optionalActions: string[]
  nextActions: DailyAction[]
  explanationLinks: DailyAction[]
  signals: DailySignal[]
  warnings: string[]
}

export interface DailyWorkbenchInput {
  dashboard: DashboardTodayResponse
  dashboardError?: PageErrorState
  report?: DailyDisciplineReport
  reportError?: PageErrorState
  portfolio?: PortfolioCurrentResponse
  portfolioError?: PageErrorState
  risks?: RiskAlert[]
  risksError?: PageErrorState
  rules?: RuleProposal[]
  rulesError?: PageErrorState
  review?: ReviewSummary
  reviewError?: PageErrorState
}

export function buildDailyWorkbenchModel(input: DailyWorkbenchInput): DailyWorkbenchModel {
  const dashboard = input.dashboard
  const report = input.report
  const decision = dashboard.decision_summary
  const riskCount = input.risks?.length ?? report?.risk_alerts?.length ?? 0
  const pendingRuleCount = (input.rules ?? []).filter((proposal) => proposal.status !== 'applied' && proposal.status !== 'rejected').length
  const reviewCount = input.review?.decision_count ?? 0
  const snapshot = input.portfolio?.snapshot
  const evidenceCount = report?.evidence.evidence_count ?? dashboard.evidence_summary?.source_count ?? 0
  const independentSourceCount = report?.evidence.independent_source_count ?? dashboard.evidence_summary?.source_count ?? 0

  const statusTone = resolveTone(dashboard.dashboard_state, report?.status, riskCount)
  const statusLabel = report?.status ? dailyReportStatusText(report.status) : dashboardStateText(dashboard.dashboard_state)
  const verdictText = report?.final_verdict || report?.summary || decision.verdict || '等待本地数据补齐后生成今日纪律建议。'
  const prohibitedActions = safeList(decision.prohibited_actions)
  const optionalActions = safeList(decision.optional_actions)
  const updatedAtText = formatDateTime(report?.updated_at || dashboard.data_updated_at)
  const trustSummary = evidenceCount > 0
    ? `证据 ${evidenceCount} 条 · 独立信源 ${independentSourceCount} 个`
    : '证据不足或尚未完成核验'
  const riskSummary = riskCount > 0
    ? `${riskCount} 条活跃风险需要人工复核`
    : snapshot && snapshot.total_assets > 0
      ? `高风险占比 ${formatPercent(snapshot.high_risk_ratio)}`
      : `高风险占比 ${formatPercent(dashboard.portfolio_summary.high_risk_ratio)}`

  const warnings = uniqueStrings([
    input.dashboardError?.message,
    input.reportError?.message,
    input.portfolioError?.message,
    input.risksError?.message,
    input.rulesError?.message,
    input.reviewError?.message,
  ].filter(Boolean) as string[])

  const nextActions = buildNextActions({
    dashboard,
    report,
    riskCount,
    evidenceCount,
    hasPortfolio: Boolean(snapshot && snapshot.total_assets > 0),
    errors: warnings,
  })

  const explanationLinks: DailyAction[] = [
    { label: '查看证据', href: '/evidence', detail: '核对今日结论引用的本地证据。', priority: 'follow_up' },
    { label: '查看数据质量', href: '/data-quality', detail: '确认数据源、检索索引与分析质量状态。', priority: 'follow_up' },
    { label: '查看决策闭环', href: '/decision-loop', detail: '追踪建议、人工记录、风险和复盘链路。', priority: 'follow_up' },
    { label: '查看审计', href: '/audit', detail: '只读检查本地审计事件。', priority: 'follow_up' },
  ]

  const signals: DailySignal[] = [
    {
      label: '数据可信度',
      value: trustSummary,
      detail: report?.p34_source_coverage?.summary || '基于今日纪律报告和证据摘要。',
      tone: evidenceCount > 0 ? 'success' : 'warning',
      href: '/data-quality',
    },
    {
      label: '组合风险',
      value: snapshot && snapshot.total_assets > 0 ? formatCurrency(snapshot.total_assets) : formatCurrency(dashboard.portfolio_summary.total_assets),
      detail: riskSummary,
      tone: riskCount > 0 || dashboard.portfolio_summary.high_risk_ratio > 0.3 ? 'danger' : 'success',
      href: '/positions',
    },
    {
      label: '风险处置',
      value: riskCount > 0 ? `${riskCount} 条待复核` : '暂无活跃风险',
      detail: riskCount > 0 ? '先查看风险预警，再记录线下处理。' : '继续保持观察。',
      tone: riskCount > 0 ? 'warning' : 'success',
      href: '/risk-alerts',
    },
    {
      label: '规则与复盘',
      value: `待确认规则 ${pendingRuleCount} · 复盘 ${reviewCount}`,
      detail: pendingRuleCount > 0 ? '规则变更仍需守门人和用户最终确认。' : '暂无待确认规则提案。',
      tone: pendingRuleCount > 0 ? 'warning' : 'success',
      href: '/rules',
    },
  ]

  return {
    statusLabel,
    statusTone,
    verdictText,
    trustSummary,
    riskSummary,
    updatedAtText,
    prohibitedActions,
    optionalActions,
    nextActions,
    explanationLinks,
    signals,
    warnings,
  }
}

function buildNextActions(input: {
  dashboard: DashboardTodayResponse
  report?: DailyDisciplineReport
  riskCount: number
  evidenceCount: number
  hasPortfolio: boolean
  errors: string[]
}): DailyAction[] {
  const actions: DailyAction[] = []
  const decisionId = input.dashboard.decision_summary.decision_id || input.report?.decision_id

  if (!input.hasPortfolio || input.report?.missing_action || input.dashboard.dashboard_state === 'first_use') {
    actions.push({
      label: '维护本地账户与持仓',
      href: '/positions',
      detail: input.report?.missing_action || '补齐本地账户和持仓后再判断今日纪律。',
      priority: 'blocking',
    })
  }

  if (input.dashboard.dashboard_state === 'insufficient_data' || input.evidenceCount === 0 || input.errors.length > 0) {
    actions.push({
      label: '查看数据质量',
      href: '/data-quality',
      detail: '先确认数据源、证据或分析服务是否降级。',
      priority: 'blocking',
    })
  }

  if (input.riskCount > 0) {
    actions.push({
      label: '处理风险预警',
      href: '/risk-alerts',
      detail: '只读查看风险 SOP，并在线下人工处理后记录。',
      priority: 'review',
    })
  }

  if (decisionId) {
    actions.push({
      label: '查看决策详情',
      href: `/decisions/${encodeURIComponent(decisionId)}`,
      detail: '阅读规则、证据和分析材料后再决定线下动作。',
      priority: 'review',
    })
  }

  if (input.report?.report_id) {
    actions.push({
      label: '查看今日纪律报告',
      href: `/daily-discipline/reports/${encodeURIComponent(input.report.report_id)}`,
      detail: '复核今日纪律摘要、证据覆盖和风险线索。',
      priority: 'review',
    })
  }

  actions.push({
    label: '发起主动咨询',
    href: '/consultation',
    detail: '带着问题生成分析材料；最终动作仍由你线下决定。',
    priority: 'follow_up',
  })

  return dedupeActions(actions)
}

function dedupeActions(actions: DailyAction[]) {
  const seen = new Set<string>()
  return actions.filter((action) => {
    const key = `${action.label}:${action.href}`
    if (seen.has(key)) return false
    seen.add(key)
    return true
  })
}

function resolveTone(dashboardState: string, reportStatus?: string, riskCount = 0): DailyTone {
  if (dashboardState === 'high_risk' || riskCount > 0) return 'danger'
  if (dashboardState === 'data_source_unavailable' || reportStatus === 'degraded') return 'degraded'
  if (dashboardState === 'insufficient_data' || dashboardState === 'first_use' || reportStatus === 'insufficient_data') return 'warning'
  if (dashboardState === 'frozen_watch') return 'warning'
  if (dashboardState === 'normal' || reportStatus === 'success') return 'success'
  return 'unknown'
}

function safeList(value: unknown): string[] {
  return Array.isArray(value) ? value.filter((item): item is string => typeof item === 'string' && item.trim().length > 0) : []
}

function uniqueStrings(values: string[]) {
  return Array.from(new Set(values))
}

function dashboardStateText(state: string) {
  const map: Record<string, string> = {
    first_use: '首次使用',
    normal: '正常',
    insufficient_data: '数据不足',
    frozen_watch: '冻结观察',
    high_risk: '高风险',
    data_source_unavailable: '数据源不可用',
    generic_failure: '加载失败',
  }
  return map[state] ?? '未知状态'
}

function dailyReportStatusText(status: string) {
  const map: Record<string, string> = {
    not_started: '未开始',
    running: '运行中',
    success: '成功',
    degraded: '降级',
    failed: '失败',
    insufficient_data: '数据不足',
  }
  return map[status] ?? '未知状态'
}
