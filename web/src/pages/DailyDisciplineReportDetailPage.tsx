import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { getDailyDisciplineReport } from '../services/dailyDisciplineReport'
import type { DailyDisciplineReport } from '../types/dailyDisciplineReport'
import { DailyDisciplineReportLinks, PortfolioOnboardingLink, RiskAlertSummaryList, dailyDisciplineReportStatusLabels, safeDailyDisciplineFailureMessage } from './DailyDisciplineReportComponents'
import { toPageErrorState } from '../shared/utils'
import { sourceCategoryText, sourceHealthStatusText, textOrRaw } from '../shared/mappers'

export function DailyDisciplineReportDetailPage() {
  const { reportId } = useParams()
  const [report, setReport] = useState<DailyDisciplineReport>()
  const [error, setError] = useState('')

  useEffect(() => {
    if (!reportId) {
      setError('缺少报告 ID')
      return
    }
    getDailyDisciplineReport(reportId)
      .then((res) => {
        setReport(res.data)
        setError('')
      })
      .catch((err: unknown) => {
        const nextError = toPageErrorState(err, '每日纪律报告详情加载失败')
        setReport(undefined)
        setError(nextError.message)
      })
  }, [reportId])

  return (
    <section className="page-card">
      <header className="page-header">
        <div>
          <h1>每日纪律报告详情</h1>
          <p>查看单日报告摘要、证据覆盖、趋势和追踪链接。</p>
        </div>
        <strong>{report ? dailyDisciplineReportStatusLabels[report.status] : '加载中'}</strong>
      </header>
      {error ? <p role="alert">{error}</p> : null}
      {report ? (
        <div className="panel-list">
          <article className="panel-card">
            <div className="row-between">
              <strong>{report.local_date}</strong>
              <span>{dailyDisciplineReportStatusLabels[report.status]}</span>
            </div>
            <p>{report.summary}</p>
            {report.final_verdict ? <p>最终裁决：{report.final_verdict}</p> : null}
            {report.verdict_status ? <small>裁决状态：{report.verdict_status}</small> : null}
            {safeDailyDisciplineFailureMessage(report) ? <p>{safeDailyDisciplineFailureMessage(report)}</p> : null}
            {report.missing_action ? <p>{report.missing_action}</p> : null}
            <PortfolioOnboardingLink report={report} />
            {report.missing_categories?.length ? (
              <div>
                <strong>缺失数据类别</strong>
                <ul>
                  {report.missing_categories.map((category) => <li key={category}>{category}</li>)}
                </ul>
              </div>
            ) : null}
          </article>

          <article className="panel-card">
            <strong>证据覆盖</strong>
            <dl className="compact-list">
              <div><dt>证据数</dt><dd>{report.evidence.evidence_count}</dd></div>
              <div><dt>独立信源</dt><dd>{report.evidence.independent_source_count}</dd></div>
              <div><dt>高等级独立信源</dt><dd>{report.evidence.high_grade_independent_source_count}</dd></div>
            </dl>
          </article>

          {report.p34_source_coverage?.summary || report.p34_source_coverage?.missing_categories?.length || report.p34_source_coverage?.source_health?.length ? (
            <article className="panel-card">
              <strong>P34 扩展数据覆盖</strong>
              {report.p34_source_coverage.summary ? <p>{report.p34_source_coverage.summary}</p> : null}
              {report.p34_source_coverage.missing_categories?.length ? (
                <div>
                  <span>缺失或降级类别</span>
                  <ul>
                    {report.p34_source_coverage.missing_categories.map((category) => <li key={category}>{category}</li>)}
                  </ul>
                </div>
              ) : null}
              {report.p34_source_coverage.source_health?.length ? (
                <ul>
                  {report.p34_source_coverage.source_health.map((item) => (
                    <li key={`${item.source_name}-${item.data_category}`}>
                      {item.source_name} · {textOrRaw(sourceCategoryText, item.data_category)} · {textOrRaw(sourceHealthStatusText, item.freshness)}；数据日：{item.data_date || '暂无'}；等级：{item.source_level || '暂无'}；影响标的：{item.affected_symbols?.join('、') || '暂无'}
                    </li>
                  ))}
                </ul>
              ) : null}
            </article>
          ) : null}

          {report.risk_alerts?.length ? (
            <article className="panel-card">
              <RiskAlertSummaryList alerts={report.risk_alerts} />
            </article>
          ) : null}

          <article className="panel-card">
            <strong>最近趋势</strong>
            <dl className="compact-list">
              <div><dt>成功</dt><dd>{report.trend.success_count}</dd></div>
              <div><dt>降级</dt><dd>{report.trend.degraded_count}</dd></div>
              <div><dt>失败</dt><dd>{report.trend.failed_count}</dd></div>
              <div><dt>数据不足</dt><dd>{report.trend.insufficient_data_count}</dd></div>
            </dl>
          </article>

          <article className="panel-card">
            <strong>追踪入口</strong>
            <DailyDisciplineReportLinks report={report} />
          </article>

          <article className="panel-card">
            <strong>安全边界</strong>
            <p>{report.safety_note}</p>
          </article>
        </div>
      ) : null}
    </section>
  )
}
