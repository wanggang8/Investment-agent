import { useEffect, useState } from 'react'
import { AuditEventTimeline } from '../components/audit/AuditEventTimeline'
import { StatusNotice } from '../components/status/StatusNotice'
import { SummaryCard, type UITone } from '../components/ui'
import { buildAuditOpsModel } from '../features/governance'
import { listAuditEvents } from '../services/audit'
import type { PageErrorState } from '../shared/utils'
import { toPageErrorState } from '../shared/utils'
import type { AuditEvent } from '../types/audit'

export function AuditPage() {
  const [events, setEvents] = useState<AuditEvent[]>([])
  const [errorState, setErrorState] = useState<PageErrorState>()

  useEffect(() => {
    listAuditEvents()
      .then((res) => {
        setEvents(res.data?.items ?? [])
        setErrorState(undefined)
      })
      .catch((error: unknown) => {
        setEvents([])
        setErrorState(toPageErrorState(error))
      })
  }, [])

  const auditModel = buildAuditOpsModel({ events })

  return (
    <div>
      <h1 className="page-title">复盘与审计</h1>
      {errorState && <StatusNotice state={errorState.state} safeMessage={errorState.message} code={errorState.code} />}
      <section className={`daily-hero daily-tone-${auditModel.overallTone}`} aria-label="审计检查总览">
        <div className="daily-hero-main">
          <div className="state-label">审计检查状态</div>
          <h2>{auditModel.overallLabel}</h2>
          <p>{auditModel.safetyNotes[0]}</p>
          <div className="daily-signal-grid quality-signal-grid">
            {auditModel.metrics.map((metric) => (
              <SummaryCard key={metric.label} title={metric.label} value={metric.value} detail={metric.detail} tone={(metric.tone ?? 'unknown') as UITone} />
            ))}
          </div>
        </div>
        <aside className="daily-hero-side" aria-label="审计下一步">
          <strong>下一步本地检查</strong>
          <ul>
            {auditModel.nextActions.map((action) => (
              <li key={action.label}>
                <a href={action.href} aria-label={`${action.label}入口`}>{action.label}</a>
                <span>{action.detail}</span>
              </li>
            ))}
          </ul>
        </aside>
      </section>
      {auditModel.eventGroups.length ? (
        <section className="cockpit-card" aria-label="审计分类摘要">
          <div className="state-label">审计分类</div>
          <div className="metric-grid">
            {auditModel.eventGroups.map((group) => (
              <div key={group.label}>
                <span>{group.label}</span>
                <strong>{group.count}</strong>
                <small>最近：{group.latestAt}</small>
              </div>
            ))}
          </div>
        </section>
      ) : null}
      <AuditEventTimeline events={events} />
    </div>
  )
}
