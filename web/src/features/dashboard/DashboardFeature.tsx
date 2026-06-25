import { useEffect, useState } from 'react'
import { MetricBarChart } from '../../components/charts/MetricBarChart'
import { EvidenceSummary } from '../../components/dashboard/EvidenceSummary'
import { FinalVerdictCard } from '../../components/dashboard/FinalVerdictCard'
import { DisciplineStatus } from '../../components/dashboard/DisciplineStatus'
import { PortfolioSummary } from '../../components/dashboard/PortfolioSummary'
import { TriggeredRules } from '../../components/dashboard/TriggeredRules'
import { UserConfirmationPanel } from '../../components/dashboard/UserConfirmationPanel'
import { DailyDecisionHero } from '../../components/dashboard/DailyDecisionHero'
import { ManualActionQueue } from '../../components/dashboard/ManualActionQueue'
import { WorkbenchSignalGrid } from '../../components/dashboard/WorkbenchSignalGrid'
import { EvidenceChecklist, ProgressTracker, SnapshotStrip } from '../../components/reference'
import { StatusNotice } from '../../components/status/StatusNotice'
import { CockpitLayout } from '../../components/layout/CockpitLayout'
import { getDashboardToday } from '../../services/dashboard'
import { getTodayDailyDisciplineReport } from '../../services/dailyDisciplineReport'
import { DailyDisciplineReportLinks, PortfolioOnboardingLink, RiskAlertSummaryList, dailyDisciplineReportStatusLabels, safeDailyDisciplineFailureMessage } from '../../pages/DailyDisciplineReportComponents'
import { buildDailyWorkbenchModel } from './dailyWorkbenchModel'
import { buildDashboardChartData } from '../../shared/mappers/charts'
import { marketStateText, textOrRaw } from '../../shared/mappers'
import type { PageErrorState } from '../../shared/utils'
import { formatCurrency, formatPercent, toPageErrorState } from '../../shared/utils'
import type { DashboardTodayResponse } from '../../types/dashboard'
import type { DailyDisciplineReport } from '../../types/dailyDisciplineReport'

const fallbackDashboard: DashboardTodayResponse = {
  dashboard_state: 'insufficient_data',
  discipline_status: '信息不足',
  data_updated_at: '',
  portfolio_summary: {
    total_assets: 0,
    cash_ratio: 0,
    high_risk_ratio: 0,
    position_count: 0,
  },
  market_summary: {
    sentiment_state: 'unknown',
    liquidity_state: 'unknown',
  },
  triggered_rules: [],
  decision_summary: {
    verdict: '等待数据补齐后生成今日纪律建议。',
    final_verdict_status: 'insufficient_data',
    prohibited_actions: ['暂停交易类建议'],
    optional_actions: ['刷新数据', '查看缺失项'],
    action_required: false,
    confirmation_status: 'not_required',
  },
}

export function DashboardFeature() {
  const [dashboard, setDashboard] = useState<DashboardTodayResponse>(fallbackDashboard)
  const [todayReport, setTodayReport] = useState<DailyDisciplineReport>()
  const [reportErrorState, setReportErrorState] = useState<PageErrorState>()
  const [reportLoaded, setReportLoaded] = useState(false)
  const [errorState, setErrorState] = useState<PageErrorState>()

  useEffect(() => {
    getDashboardToday()
      .then((res) => {
        if (res.data) {
          setDashboard(res.data)
          setErrorState(undefined)
        }
      })
      .catch((error: unknown) => {
        const nextError = toPageErrorState(error)
        setDashboard({ ...fallbackDashboard, dashboard_state: nextError.state })
        setErrorState(nextError)
      })

    getTodayDailyDisciplineReport()
      .then((res) => {
        setTodayReport(res.data)
        setReportErrorState(undefined)
        setReportLoaded(true)
      })
      .catch((error: unknown) => {
        setTodayReport(undefined)
        setReportErrorState(toPageErrorState(error, '每日纪律报告加载失败'))
        setReportLoaded(true)
      })
  }, [])

  const charts = buildDashboardChartData(dashboard)
  const dailyModel = buildDailyWorkbenchModel({
    dashboard,
    dashboardError: errorState,
    report: todayReport,
    reportError: reportErrorState,
  })

  const systemPanel = (
    <>
      <DisciplineStatus dashboard={dashboard} errorMessage={errorState?.message} />
      <PortfolioSummary summary={dashboard.portfolio_summary} />
      <MetricBarChart title="仓位结构" data={charts.allocation} />
    </>
  )

  const decisionPanel = (
    <>
      <FinalVerdictCard summary={dashboard.decision_summary} />
      <TriggeredRules rules={dashboard.triggered_rules} />
      {dashboard.decision_summary.action_required ? (
        <article className="cockpit-card confirmation-card">
          <div className="state-label">用户确认</div>
          <p>这条纪律建议需要在决策详情页记录线下处理结果。</p>
          <a href={dashboard.decision_summary.decision_id ? `/decisions/${dashboard.decision_summary.decision_id}` : '/decisions'}>前往决策详情确认</a>
          <small>系统只记录你的线下动作，不会替你买入或卖出。</small>
        </article>
      ) : (
        <UserConfirmationPanel
          availableActions={[]}
          confirmationStatus={dashboard.decision_summary.confirmation_status}
        />
      )}
    </>
  )

  const evidencePanel = (
    <>
      <EvidenceSummary summary={dashboard.evidence_summary} dashboardState={dashboard.dashboard_state} />
      {errorState && <StatusNotice state={dashboard.dashboard_state} safeMessage={errorState.message} code={errorState.code} />}
      <MetricBarChart title="风险分位" data={charts.risk} />
      <MetricBarChart title="证据覆盖" data={charts.evidence} />
      <article className="cockpit-card">
        <div className="state-label">状态说明</div>
        <p>{dashboard.dashboard_state === 'frozen_watch' ? '等待条件：证据核验完成、市场快照可用、规则版本确认。' : '缺失项：账户快照、证据核验或市场快照；暂停原因：数据不足时不展示交易类动作。'}</p>
      </article>
      <article className="cockpit-card">
        <div className="state-label">市场状态</div>
        <dl className="compact-list">
          <div>
            <dt>PE 分位</dt>
            <dd>{dashboard.market_summary.pe_percentile ?? '暂无'}</dd>
          </div>
          <div>
            <dt>PB 分位</dt>
            <dd>{dashboard.market_summary.pb_percentile ?? '暂无'}</dd>
          </div>
          <div>
            <dt>情绪</dt>
            <dd>{textOrRaw(marketStateText, dashboard.market_summary.sentiment_state)}</dd>
          </div>
          <div>
            <dt>流动性</dt>
            <dd>{textOrRaw(marketStateText, dashboard.market_summary.liquidity_state)}</dd>
          </div>
        </dl>
      </article>
    </>
  )

  return (
    <div>
      <h1 className="page-title">今日纪律</h1>
      <DailyDecisionHero model={dailyModel} />
      <section className="reference-command-grid" aria-label="今日纪律首屏">
        <ManualActionQueue actions={dailyModel.nextActions} />
        <div className="reference-side-stack">
          <WorkbenchSignalGrid signals={dailyModel.signals} />
          <SnapshotStrip
            title="持仓与资金快照"
            updatedAt={dashboard.data_updated_at || dailyModel.updatedAtText}
            items={[
              { label: '总资产（估）', value: formatCurrency(dashboard.portfolio_summary.total_assets) },
              { label: '现金占比', value: formatPercent(dashboard.portfolio_summary.cash_ratio) },
              { label: '持仓数量', value: dashboard.portfolio_summary.position_count },
              { label: '高风险占比', value: formatPercent(dashboard.portfolio_summary.high_risk_ratio), status: dashboard.portfolio_summary.high_risk_ratio > 0.3 ? '需关注' : '正常' },
            ]}
          />
        </div>
      </section>
      <section className="reference-lower-grid" aria-label="今日解释与证据">
        <ProgressTracker
          title="最近咨询 · 解释预览"
          steps={[
            { label: '输入假设', status: 'done' },
            { label: '信息核查', status: todayReport ? 'done' : 'pending', detail: todayReport ? `${todayReport.evidence.evidence_count} 条证据` : '待检查' },
            { label: '分析材料', status: 'done', detail: '只作分析材料' },
            { label: '规则裁决', status: dashboard.decision_summary.final_verdict_status ? 'active' : 'pending' },
            { label: '最终建议', status: dashboard.decision_summary.verdict ? 'done' : 'pending' },
            { label: '等待人工确认', status: dashboard.decision_summary.action_required ? 'active' : 'pending' },
          ]}
        >
          <p>规则裁决进行中时，页面只展示解释与证据，不会发起任何交易或自动确认。</p>
        </ProgressTracker>
        <EvidenceChecklist
          title="证据与规则快照"
          items={[
            { label: '信息核查来源', value: todayReport ? `${todayReport.evidence.independent_source_count}/${todayReport.evidence.evidence_count} 覆盖` : '待检查', status: todayReport ? 'done' : 'pending' },
            { label: '分析材料', value: '只作材料', status: 'done' },
            { label: '关键规则通过率', value: dashboard.triggered_rules.length ? `${dashboard.triggered_rules.length} 条命中` : '已通过', status: dashboard.triggered_rules.length ? 'active' : 'done' },
            { label: '审计只读记录', value: todayReport?.audit_link ? '可查看' : '待记录', status: todayReport?.audit_link ? 'done' : 'pending' },
          ]}
          action={{ label: '查看详情', href: '/evidence' }}
        />
      </section>
      {dailyModel.warnings.length > 0 ? (
        <section className="stacked-panel" aria-label="今日状态提示">
          {dailyModel.warnings.map((warning) => (
            <StatusNotice key={warning} state="generic_failure" safeMessage={warning} />
          ))}
        </section>
      ) : null}
      {reportLoaded && (todayReport || reportErrorState) ? (
        <article className="cockpit-card" style={{ marginBottom: '1rem' }}>
          <div className="row-between">
            <div className="state-label">今日纪律报告</div>
            {todayReport ? <strong>{dailyDisciplineReportStatusLabels[todayReport.status]}</strong> : null}
          </div>
          {reportErrorState ? <p role="alert">{reportErrorState.message}</p> : null}
          {todayReport ? (
            <>
              <p>{todayReport.summary}</p>
              {safeDailyDisciplineFailureMessage(todayReport) ? <p>{safeDailyDisciplineFailureMessage(todayReport)}</p> : null}
              {todayReport.missing_action ? <p>{todayReport.missing_action}</p> : null}
              <PortfolioOnboardingLink report={todayReport} />
              <DailyDisciplineReportLinks report={todayReport} />
              <RiskAlertSummaryList alerts={todayReport.risk_alerts} title="今日风险预警" />
              <small>{todayReport.safety_note}</small>
            </>
          ) : null}
        </article>
      ) : null}
      <h2 className="section-title">详细信号</h2>
      <CockpitLayout systemPanel={systemPanel} decisionPanel={decisionPanel} evidencePanel={evidencePanel} />
    </div>
  )
}
