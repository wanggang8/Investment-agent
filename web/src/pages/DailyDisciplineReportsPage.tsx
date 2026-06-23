import { useEffect, useState } from 'react'
import { listDailyDisciplineReports } from '../services/dailyDisciplineReport'
import type { DailyDisciplineReport } from '../types/dailyDisciplineReport'
import { DailyDisciplineReportCard } from './DailyDisciplineReportComponents'
import { toPageErrorState } from '../shared/utils'
import { buildDailyOpsModel } from '../features/governance'

export function DailyDisciplineReportsPage() {
  const [reports, setReports] = useState<DailyDisciplineReport[]>([])
  const [loaded, setLoaded] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    listDailyDisciplineReports()
      .then((res) => {
        setReports(res.data?.reports ?? [])
        setLoaded(true)
        setError('')
      })
      .catch((err: unknown) => {
        const nextError = toPageErrorState(err, '每日纪律报告加载失败')
        setLoaded(true)
        setReports([])
        setError(nextError.message)
      })
  }, [])

  const dailyModel = buildDailyOpsModel({ reports })

  return (
    <section>
      <header className="page-header">
        <div>
          <h1>每日纪律报告历史</h1>
          <p>回看最近每日纪律报告，追踪数据缺口、证据覆盖和人工复核边界。</p>
        </div>
      </header>
      <section className={`daily-hero daily-tone-${dailyModel.overallTone}`} aria-label="每日纪律复盘总览">
        <div className="daily-hero-main">
          <div className="state-label">每日纪律复盘状态</div>
          <h2>{dailyModel.overallLabel}</h2>
          <p>{dailyModel.safetyNotes[0]}</p>
          <div className="daily-signal-grid quality-signal-grid">
            {dailyModel.metrics.map((metric) => (
              <article key={metric.label} className={`daily-signal daily-tone-${metric.tone ?? 'unknown'}`}>
                <h2>{metric.label}</h2>
                <strong>{metric.value}</strong>
                {metric.detail ? <p>{metric.detail}</p> : null}
              </article>
            ))}
          </div>
        </div>
        <aside className="daily-hero-side" aria-label="每日纪律下一步">
          <strong>下一步人工复盘</strong>
          <ul>
            {dailyModel.nextActions.map((action) => (
              <li key={action.label}>
                <a href={action.href} aria-label={`${action.label}入口`}>{action.label}</a>
                <span>{action.detail}</span>
              </li>
            ))}
          </ul>
        </aside>
      </section>
      {error ? <p role="alert">{error}</p> : null}
      {loaded && reports.length === 0 && !error ? <p>暂无每日纪律报告</p> : null}
      {reports.length > 0 ? (
        <div className="panel-list">
          {reports.map((report) => (
            <DailyDisciplineReportCard key={report.report_id} report={report} showDetailLink />
          ))}
        </div>
      ) : null}
    </section>
  )
}
