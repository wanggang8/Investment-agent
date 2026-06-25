import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { DailyDecisionHero } from '../components/dashboard/DailyDecisionHero'
import { ManualActionQueue } from '../components/dashboard/ManualActionQueue'
import { WorkbenchSignalGrid } from '../components/dashboard/WorkbenchSignalGrid'
import { EvidenceChecklist, ProgressTracker, SnapshotStrip } from '../components/reference'
import { StatusNotice } from '../components/status/StatusNotice'
import { buildDailyWorkbenchModel } from '../features/dashboard/dailyWorkbenchModel'
import { getTodayDailyDisciplineReport } from '../services/dailyDisciplineReport'
import { getDashboardToday } from '../services/dashboard'
import { getPortfolioCurrent } from '../services/portfolio'
import { getReviewSummary } from '../services/review'
import { listRiskAlerts } from '../services/riskAlert'
import { listRuleProposals } from '../services/rule'
import type { PageErrorState } from '../shared/utils'
import { formatCurrency, formatPercent, toPageErrorState } from '../shared/utils'
import type { DailyDisciplineReport } from '../types/dailyDisciplineReport'
import type { DashboardTodayResponse } from '../types/dashboard'
import type { PortfolioCurrentResponse } from '../types/portfolio'
import type { ReviewSummary } from '../types/review'
import type { RiskAlert } from '../types/riskAlert'
import type { RuleProposal } from '../types/rule'

type LoadState<T> = {
  data?: T
  error?: PageErrorState
}

type WorkbenchState = {
  dashboard: LoadState<DashboardTodayResponse>
  report: LoadState<DailyDisciplineReport>
  portfolio: LoadState<PortfolioCurrentResponse>
  risks: LoadState<RiskAlert[]>
  rules: LoadState<RuleProposal[]>
  review: LoadState<ReviewSummary>
}

const initialState: WorkbenchState = {
  dashboard: {},
  report: {},
  portfolio: {},
  risks: {},
  rules: {},
  review: {},
}

let workbenchLoadInFlight: Promise<WorkbenchState> | undefined

export function WorkbenchPage() {
  const [state, setState] = useState<WorkbenchState>(initialState)

  useEffect(() => {
    let mounted = true

    async function load() {
      const nextState = await loadWorkbenchState()

      if (!mounted) return
      setState(nextState)
    }

    load()
    return () => {
      mounted = false
    }
  }, [])

  const reportHref = state.report.data?.report_id
    ? `/daily-discipline/reports/${encodeURIComponent(state.report.data.report_id)}`
    : '/daily-discipline/reports'
  const pendingRuleCount = (state.rules.data ?? []).filter((proposal) => proposal.status !== 'applied' && proposal.status !== 'rejected').length
  const riskCount = state.risks.data?.length ?? 0
  const reviewCount = state.review.data?.decision_count ?? 0
  const snapshot = state.portfolio.data?.snapshot
  const dailyModel = buildDailyWorkbenchModel({
    dashboard: state.dashboard.data ?? fallbackWorkbenchDashboard,
    dashboardError: state.dashboard.error,
    report: state.report.data,
    reportError: state.report.error,
    portfolio: state.portfolio.data,
    portfolioError: state.portfolio.error,
    risks: state.risks.data,
    risksError: state.risks.error,
    rules: state.rules.data,
    rulesError: state.rules.error,
    review: state.review.data,
    reviewError: state.review.error,
  })

  return (
    <div>
      <h1 className="page-title">用户决策工作台</h1>
      <p className="page-intro">聚合今日纪律、本地事实、风险和复盘入口；这里只做只读提示和导航，最终动作仍由你线下决定。</p>

      <DailyDecisionHero model={dailyModel} />
      <section className="reference-command-grid" aria-label="工作台首屏">
        <ManualActionQueue actions={dailyModel.nextActions} />
        <div className="reference-side-stack">
          <WorkbenchSignalGrid signals={dailyModel.signals} />
          <SnapshotStrip
            title="持仓与资金快照"
            updatedAt={snapshot?.snapshot_time || dailyModel.updatedAtText}
            items={[
              { label: '总资产（估）', value: formatCurrency(snapshot?.total_assets ?? state.dashboard.data?.portfolio_summary.total_assets ?? 0) },
              { label: '可用资金', value: formatCurrency(snapshot?.cash ?? 0) },
              { label: '持仓数量', value: snapshot?.position_count ?? state.dashboard.data?.portfolio_summary.position_count ?? 0 },
              { label: '仓位水平', value: formatPercent(1 - (snapshot?.cash_ratio ?? state.dashboard.data?.portfolio_summary.cash_ratio ?? 0)), status: riskCount > 0 ? '需关注' : '正常' },
            ]}
          />
        </div>
      </section>
      <section className="reference-lower-grid" aria-label="工作台解释与证据">
        <ProgressTracker
          title="最近咨询 · 解释预览"
          actions={<Link to="/decision-loop">查看全部记录</Link>}
          steps={[
            { label: '输入假设', status: 'done' },
            { label: '信息核查', status: state.report.data ? 'done' : 'pending', detail: state.report.data ? `${state.report.data.evidence.evidence_count} 条证据` : '待检查' },
            { label: '分析材料', status: 'done', detail: '只作分析材料' },
            { label: '规则裁决', status: state.dashboard.data?.decision_summary?.final_verdict_status ? 'active' : 'pending' },
            { label: '最终建议', status: state.dashboard.data?.decision_summary?.verdict ? 'done' : 'pending' },
            { label: '等待人工确认', status: state.dashboard.data?.decision_summary?.action_required ? 'active' : 'pending' },
          ]}
        >
          <p>当前进展：系统只生成解释、证据和规则裁决记录；最终动作仍由你线下完成。</p>
        </ProgressTracker>
        <EvidenceChecklist
          title="证据与规则快照"
          items={[
            { label: '信息核查来源', value: state.report.data ? `${state.report.data.evidence.independent_source_count}/${state.report.data.evidence.evidence_count} 覆盖` : '待检查', status: state.report.data ? 'done' : 'pending' },
            { label: '分析材料', value: '只作材料', status: 'done' },
            { label: '关键规则通过率', value: pendingRuleCount > 0 ? `${pendingRuleCount} 项待确认` : '已通过', status: pendingRuleCount > 0 ? 'active' : 'done' },
            { label: '审计只读记录', value: `${reviewCount} 条复盘`, status: reviewCount > 0 ? 'done' : 'pending' },
          ]}
          action={{ label: '查看详情', href: '/evidence' }}
        />
      </section>

      {dailyModel.warnings.length > 0 ? (
        <section className="stacked-panel" aria-label="工作台状态提示">
          {dailyModel.warnings.map((warning) => (
            <StatusNotice key={warning} state="generic_failure" safeMessage={warning} />
          ))}
        </section>
      ) : null}

      <h2 className="section-title">任务展开</h2>
      <section className="cockpit-grid" aria-label="用户决策工作台区域">
        <article className="cockpit-card">
          <div className="state-label">今日复核</div>
          <h2>纪律报告与解释</h2>
          {state.dashboard.error && <StatusNotice state={state.dashboard.error.state} safeMessage={state.dashboard.error.message} code={state.dashboard.error.code} />}
          {state.report.error && <StatusNotice state={state.report.error.state} safeMessage={state.report.error.message} code={state.report.error.code} />}
          {state.report.data ? (
            <>
              <p>{state.report.data.summary || '今日纪律报告暂无摘要。'}</p>
              <p>状态：{statusText(state.report.data.status)}</p>
              {state.report.data.final_verdict && <p>最终裁决：{state.report.data.final_verdict}</p>}
              <p>证据：{state.report.data.evidence.evidence_count} 条，独立信源 {state.report.data.evidence.independent_source_count} 个</p>
            </>
          ) : (
            <p>暂无今日纪律报告，先查看历史报告或完成本地数据准备。</p>
          )}
          {state.dashboard.data?.triggered_rules?.length ? (
            <ul>
              {state.dashboard.data.triggered_rules.slice(0, 2).map((rule) => (
                <li key={rule.rule_id}>{rule.description}</li>
              ))}
            </ul>
          ) : null}
          <div className="link-row">
            <Link to={reportHref}>查看纪律报告</Link>
            <Link to="/">打开今日纪律</Link>
            <Link to="/data-quality">查看数据质量</Link>
          </div>
        </article>

        <article className="cockpit-card">
          <div className="state-label">组合与风险</div>
          <h2>持仓与 SOP</h2>
          {state.portfolio.error && <StatusNotice state={state.portfolio.error.state} safeMessage={state.portfolio.error.message} code={state.portfolio.error.code} />}
          {state.risks.error && <StatusNotice state={state.risks.error.state} safeMessage={state.risks.error.message} code={state.risks.error.code} />}
          {snapshot && snapshot.total_assets > 0 ? (
            <>
              <p>总资产：{formatCurrency(snapshot.total_assets)}</p>
              <p>现金占比：{formatPercent(snapshot.cash_ratio)}</p>
              <p>高风险占比：{formatPercent(snapshot.high_risk_ratio)}</p>
              <p>持仓数量：{snapshot.position_count}</p>
            </>
          ) : (
            <p>暂无持仓快照，先完成本地账户校准。</p>
          )}
          <p>活跃风险：{riskCount}</p>
          {state.risks.data?.length ? (
            <ul>
              {state.risks.data.slice(0, 2).map((alert) => (
                <li key={alert.alert_id}>{alert.symbol} · {alert.trigger_summary}</li>
              ))}
            </ul>
          ) : (
            <p>暂无活跃风险预警。</p>
          )}
          <div className="link-row">
            <Link to="/positions">查看持仓</Link>
            <Link to="/risk-alerts">查看风险预警</Link>
          </div>
        </article>

        <article className="cockpit-card">
          <div className="state-label">规则与复盘</div>
          <h2>规则提案与追踪</h2>
          {state.rules.error && <StatusNotice state={state.rules.error.state} safeMessage={state.rules.error.message} code={state.rules.error.code} />}
          {state.review.error && <StatusNotice state={state.review.error.state} safeMessage={state.review.error.message} code={state.review.error.code} />}
          <p>待确认规则：{pendingRuleCount}</p>
          <p>复盘决策：{reviewCount}</p>
          {pendingRuleCount > 0 ? (
            <ul>
              {(state.rules.data ?? []).slice(0, 2).map((proposal) => (
                <li key={proposal.proposal_id}>{proposal.title} · {ruleStatusText(proposal.status)}</li>
              ))}
            </ul>
          ) : (
            <p>暂无待确认规则提案。</p>
          )}
          {reviewCount > 0 ? (
            <p>错误样本：{state.review.data?.error_case_count ?? 0}，审计事件：{state.review.data?.audit_event_count ?? 0}</p>
          ) : (
            <p>暂无复盘活动数据。</p>
          )}
          <div className="link-row">
            <Link to="/rules">查看规则提案</Link>
            <Link to="/review">查看复盘摘要</Link>
            <Link to="/decision-loop">查看决策闭环</Link>
            <Link to="/audit">查看审计</Link>
          </div>
        </article>

        <article className="cockpit-card">
          <div className="state-label">主动咨询入口</div>
          <h2>带着问题进入咨询</h2>
          <p>主动咨询由你提交问题，系统只生成分析材料；最终动作仍由你线下决定。</p>
          <p>遇到数据缺失、数据源健康降级、分析模型或检索索引不可用时，页面只展示安全状态与现有规则事实。</p>
          <div className="link-row">
            <Link to="/consultation">发起主动咨询</Link>
          </div>
        </article>
      </section>
    </div>
  )
}

const fallbackWorkbenchDashboard: DashboardTodayResponse = {
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
    verdict: '等待本地数据补齐后生成今日纪律建议。',
    final_verdict_status: 'insufficient_data',
    prohibited_actions: ['暂停交易类建议'],
    optional_actions: ['补齐本地事实', '查看数据质量'],
    action_required: false,
    confirmation_status: 'not_required',
  },
}

function loadWorkbenchState() {
  if (!workbenchLoadInFlight) {
    workbenchLoadInFlight = loadWorkbenchStateOnce().finally(() => {
      workbenchLoadInFlight = undefined
    })
  }
  return workbenchLoadInFlight
}

async function loadWorkbenchStateOnce(): Promise<WorkbenchState> {
  const dashboard = await settle(getDashboardToday())
  const report = await settle(getTodayDailyDisciplineReport())
  const portfolio = await settle(getPortfolioCurrent())
  const risks = await settle(listRiskAlerts({ statuses: ['active', 'escalated'] }))
  const rules = await settle(listRuleProposals())
  const review = await settle(getReviewSummary())

  return {
    dashboard,
    report,
    portfolio,
    risks: risks.data ? { data: risks.data.items } : { error: risks.error },
    rules: rules.data ? { data: rules.data.items } : { error: rules.error },
    review,
  }
}

async function settle<T>(promise: Promise<{ data?: T }>): Promise<LoadState<T>> {
  try {
    const response = await promise
    return { data: response.data }
  } catch (error: unknown) {
    return { error: toPageErrorState(error) }
  }
}

function statusText(status: string) {
  const map: Record<string, string> = {
    not_started: '未开始',
    running: '运行中',
    success: '成功',
    degraded: '降级',
    failed: '失败',
    insufficient_data: '数据不足',
  }
  return map[status] ?? status
}

function ruleStatusText(status: string) {
  const map: Record<string, string> = {
    draft: '草稿',
    pending_user_confirm: '待用户确认',
    under_gatekeeper_audit: '守门人审计中',
    pending_final_confirm: '待最终确认',
    rejected: '已拒绝',
    applied: '已应用',
  }
  return map[status] ?? status
}
