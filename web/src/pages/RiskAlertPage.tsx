import { useCallback, useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { Button, EmptyState } from '../components/ui'
import { buildRiskDispositionModel } from '../features/risk/riskDispositionModel'
import { getRiskAlert, listRiskAlerts, updateRiskAlertLifecycle } from '../services/riskAlert'
import { riskSeverityText, riskSOPStatusText, riskTypeText, textOrRaw } from '../shared/mappers'
import { toPageErrorState } from '../shared/utils'
import type { RiskAlert, RiskSOPStatus } from '../types/riskAlert'

export function RiskAlertPage() {
  const { alertId } = useParams()
  const [alerts, setAlerts] = useState<RiskAlert[]>([])
  const [loaded, setLoaded] = useState(false)
  const [error, setError] = useState('')

  const load = useCallback(async () => {
    try {
      if (alertId) {
        const res = await getRiskAlert(alertId)
        setAlerts(res.data ? [res.data] : [])
      } else {
        const res = await listRiskAlerts()
        setAlerts(res.data?.items ?? [])
      }
      setLoaded(true)
      setError('')
    } catch (err) {
      const nextError = toPageErrorState(err, '风险预警加载失败')
      setAlerts([])
      setLoaded(true)
      setError(nextError.message)
    }
  }, [alertId])

  useEffect(() => {
    void load()
  }, [load])

  const updateStatus = async (alertId: string, status: RiskSOPStatus, label: string) => {
    await updateRiskAlertLifecycle(alertId, { status, reason: `前端人工 SOP 操作：${label}` })
    await load()
  }
  const disposition = buildRiskDispositionModel(alerts)

  return (
    <section className="page-card risk-alert-page reference-tight-page">
      <h1 className="page-title">风险预警中心</h1>
      <section className={`daily-hero daily-tone-${disposition.highestSeverity === 'critical' ? 'danger' : disposition.highestSeverity === 'warning' ? 'warning' : 'success'}`} aria-label="风险处置队列">
        <div className="daily-hero-main">
          <div className="state-label">风险处置队列</div>
          <h2>{disposition.summaryLabel}</h2>
          <p>最高严重程度：{textOrRaw(riskSeverityText, disposition.highestSeverity)}</p>
          <p>影响标的：{disposition.affectedSymbols.join('、') || '暂无'}</p>
        </div>
        <aside className="daily-hero-side" aria-label="风险下一步">
          <strong>{disposition.nextActions[0]?.label}</strong>
          <p>{disposition.nextActions[0]?.detail}</p>
          <small>{disposition.safetyNotes[0]}</small>
        </aside>
      </section>
      <p className="reference-page-note">风险预警中心集中查看风险类型、SOP 状态、触发依据和人工复核动作；只提供本地人工复核线索，不会自动交易，也不会调用券商接口。当前 {alerts.length} 条。</p>
      {error ? <p role="alert">{error}</p> : null}
      {loaded && alerts.length === 0 && !error ? (
        <EmptyState title="暂无需要处置的风险预警" description="当前没有需要处置的本地风险。" action={{ label: '返回工作台', href: '/workbench' }} />
      ) : null}
      <div className="risk-queue-grid" aria-label="风险队列">
        {disposition.queues.map((queue) => (
          <section key={queue.id} className="cockpit-card">
            <div className="row-between">
              <div>
                <div className="state-label">SOP 队列</div>
                <h2>{queue.label}队列</h2>
                <p className="muted-text">{queue.description}</p>
              </div>
              <strong>{queue.items.length}</strong>
            </div>
            {queue.items.length ? (
              <div className="panel-list">
                {queue.items.map((alert) => (
                  <RiskAlertCard key={alert.alert_id} alert={alert} onUpdate={updateStatus} />
                ))}
              </div>
            ) : (
              <p className="muted-text">暂无</p>
            )}
          </section>
        ))}
      </div>
    </section>
  )
}

function RiskAlertCard({ alert, onUpdate }: { alert: RiskAlert; onUpdate: (alertId: string, status: RiskSOPStatus, label: string) => void }) {
  return (
    <article className={`panel-card risk-alert-card risk-alert-card-${alert.severity}`}>
      <div className="row-between">
        <div>
          <strong>{textOrRaw(riskTypeText, alert.risk_type)}</strong>
          <p>{alert.symbol} · {textOrRaw(riskSeverityText, alert.severity)} · <span>{textOrRaw(riskSOPStatusText, alert.sop_status)}</span></p>
        </div>
        <small>{alert.updated_at}</small>
      </div>
      <p>{alert.trigger_summary}</p>
      <RiskContextSummary context={alert.trigger_context} />
      <p><span>禁止动作：{alert.prohibited_actions?.join('、') || '暂无'}</span></p>
      <p><span>建议人工动作：{alert.suggested_actions?.join('、') || '暂无'}</span></p>
      <RiskAlertLinks alert={alert} />
      <div className="action-row" aria-label="风险 SOP 操作">
        {!['resolved', 'archived'].includes(alert.sop_status) ? (
          <>
            <Button onClick={() => onUpdate(alert.alert_id, 'observing', '记录继续观察')}>记录继续观察</Button>
            <Button onClick={() => onUpdate(alert.alert_id, 'escalated', '记录升级复核')}>记录升级复核</Button>
            <Button variant="secondary" onClick={() => onUpdate(alert.alert_id, 'resolved', '记录本地解除预警')}>记录本地解除预警</Button>
          </>
        ) : null}
      </div>
      <small>{alert.safety_note}</small>
    </article>
  )
}

function RiskContextSummary({ context }: { context?: unknown }) {
  if (!context || typeof context !== 'object') return null
  const item = context as { sop?: string; data_prerequisites?: unknown; llm_role?: string }
  const prerequisites = Array.isArray(item.data_prerequisites) ? item.data_prerequisites.filter((value): value is string => typeof value === 'string') : []
  if (!item.sop && !prerequisites.length && !item.llm_role) return null
  return (
    <div className="muted-text" aria-label="SOP 上下文">
      {item.sop && <p>SOP：{item.sop}</p>}
      {prerequisites.length ? <p>数据前提：{prerequisites.join('、')}</p> : null}
      {item.llm_role && (
        <>
          <p>分析角色：{item.llm_role}</p>
          <span className="reference-sr-only">LLM 角色：{item.llm_role}</span>
        </>
      )}
    </div>
  )
}

function RiskAlertLinks({ alert }: { alert: RiskAlert }) {
  const links = [
    alert.report_link ? { to: alert.report_link, label: '关联报告' } : undefined,
    alert.decision_link ? { to: alert.decision_link, label: '关联决策' } : undefined,
    alert.notification_link ? { to: alert.notification_link, label: '关联通知' } : undefined,
    alert.audit_link ? { to: alert.audit_link, label: '关联审计' } : undefined,
  ].filter((item): item is { to: string; label: string } => Boolean(item))

  if (!links.length) {
    return null
  }
  return (
    <div className="action-row">
      {links.map((link) => (
        <Link key={link.label} to={link.to} className="link-button">{link.label}</Link>
      ))}
    </div>
  )
}
