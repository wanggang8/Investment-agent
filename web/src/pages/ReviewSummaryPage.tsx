import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { MetricBarChart } from '../components/charts/MetricBarChart'
import { ReviewSummaryPanel } from '../components/review/ReviewSummaryPanel'
import { StatusNotice } from '../components/status/StatusNotice'
import { getReviewSummary } from '../services/review'
import { buildReviewActivityData } from '../shared/mappers/charts'
import type { PageErrorState } from '../shared/utils'
import { toPageErrorState } from '../shared/utils'
import type { ReviewSummary } from '../types/review'

export function ReviewSummaryPage() {
  const [summary, setSummary] = useState<ReviewSummary>()
  const [errorState, setErrorState] = useState<PageErrorState>()

  useEffect(() => {
    getReviewSummary()
      .then((res) => {
        setSummary(res.data)
        setErrorState(undefined)
      })
      .catch((error: unknown) => {
        setSummary(undefined)
        setErrorState(toPageErrorState(error))
      })
  }, [])

  return (
    <div>
      <h1 className="page-title">复盘摘要</h1>
      {errorState && <StatusNotice state={errorState.state} safeMessage={errorState.message} code={errorState.code} />}
      <section className="daily-hero daily-tone-readonly" aria-label="复盘摘要首屏状态">
        <div className="daily-hero-main">
          <div className="state-label">复盘追踪状态</div>
          <h2>{summary ? '复盘事实已汇总，规则变更仍需人工确认' : '等待本地复盘事实'}</h2>
          <p>本页只读展示月度/季度复盘、错误样本、降级工作流和规则提案线索。</p>
          <dl className="daily-hero-meta">
            <div>
              <dt>复盘决策</dt>
              <dd>{summary?.decision_count ?? 0}</dd>
            </div>
            <div>
              <dt>降级数量</dt>
              <dd>{summary?.degraded_count ?? 0}</dd>
            </div>
            <div>
              <dt>证据缺口</dt>
              <dd>{summary?.missing_evidence_count ?? 0}</dd>
            </div>
          </dl>
        </div>
        <aside className="daily-hero-side" aria-label="复盘摘要下一步">
          <strong>下一步人工复盘</strong>
          <p>进入决策闭环核对人工记录、风险追踪和审计引用。</p>
          <strong>禁止动作</strong>
          <p>规则变更不会自动应用，仍需守门人审计和用户最终确认。</p>
        </aside>
      </section>
      <div className="link-row">
        <Link to="/decision-loop">查看决策闭环</Link>
      </div>
      <MetricBarChart title="复盘活动" data={buildReviewActivityData(summary)} emptyText="暂无复盘活动数据。" />
      <ReviewSummaryPanel summary={summary} />
    </div>
  )
}
