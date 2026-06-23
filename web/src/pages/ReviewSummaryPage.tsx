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
      <div className="link-row">
        <Link to="/decision-loop">查看决策闭环</Link>
      </div>
      <MetricBarChart title="复盘活动" data={buildReviewActivityData(summary)} emptyText="暂无复盘活动数据。" />
      <ReviewSummaryPanel summary={summary} />
    </div>
  )
}
