import { useEffect, useState } from 'react'
import { Button, EmptyState, SummaryCard, type UITone } from '../components/ui'
import type { NotificationItem } from '../types/notification'
import { listNotifications, markAllNotificationsRead, markNotificationRead } from '../services/notification'
import { buildNotificationInboxModel } from '../features/governance'

const notificationPollIntervalMs = 30_000

export function NotificationPage({ pollIntervalMs = notificationPollIntervalMs }: { pollIntervalMs?: number } = {}) {
  const [items, setItems] = useState<NotificationItem[]>([])
  const [unreadCount, setUnreadCount] = useState(0)
  const [error, setError] = useState('')

  const load = async () => {
    try {
      const res = await listNotifications()
      setItems(res.data?.items ?? [])
      setUnreadCount(res.data?.unread_count ?? 0)
      setError('')
    } catch (err) {
      setError(err instanceof Error ? err.message : '通知加载失败')
    }
  }

  useEffect(() => {
    const refresh = () => {
      listNotifications()
        .then((res) => {
          setItems(res.data?.items ?? [])
          setUnreadCount(res.data?.unread_count ?? 0)
          setError('')
        })
        .catch((err: unknown) => {
          setError(err instanceof Error ? err.message : '通知加载失败')
        })
    }
    refresh()
    const timer = window.setInterval(refresh, pollIntervalMs)
    return () => window.clearInterval(timer)
  }, [pollIntervalMs])

  const markOneRead = async (notificationId: string) => {
    await markNotificationRead(notificationId)
    await load()
  }

  const markAllRead = async () => {
    await markAllNotificationsRead()
    await load()
  }

  const inboxModel = buildNotificationInboxModel({ notifications: items, unreadCount })

  return (
    <section className="reference-tight-page">
      <h1 className="page-title">通知中心</h1>
      <section className={`daily-hero daily-tone-${inboxModel.overallTone}`} aria-label="本地通知总览">
        <div className="daily-hero-main">
          <div className="state-label">本地通知收件箱</div>
          <h2>{inboxModel.overallLabel}</h2>
          <p>{inboxModel.safetyNotes[0]}</p>
          <div className="daily-signal-grid quality-signal-grid">
            {inboxModel.metrics.map((metric) => (
              <SummaryCard key={metric.label} title={metric.label} value={metric.value} tone={(metric.tone ?? 'unknown') as UITone} />
            ))}
          </div>
        </div>
        <aside className="daily-hero-side" aria-label="通知下一步">
          <strong>下一步本地处理</strong>
          <ul>
            {inboxModel.nextActions.map((action) => (
              <li key={action.label}>
                <a href={action.href} aria-label={`${action.label}入口`}>{action.label}</a>
                <span>{action.detail}</span>
              </li>
            ))}
          </ul>
          <Button onClick={markAllRead} disabled={unreadCount === 0}>全部标记已读</Button>
        </aside>
      </section>
      <p className="reference-page-note"><span>{`未读通知：${unreadCount}`}</span>。通知中心只处理本地应用内状态，不发送站外消息，也不代表任何交易许可。</p>
      {inboxModel.sourceGroups.length ? (
        <section className="cockpit-card" aria-label="通知来源分类">
          <div className="state-label">来源分类</div>
          <div className="metric-grid">
            {inboxModel.sourceGroups.map((group) => (
              <div key={group.label}>
                <span>{group.label}</span>
                <strong>{group.count}</strong>
              </div>
            ))}
          </div>
        </section>
      ) : null}
      {error ? <p role="alert">{error}</p> : null}
      {items.length === 0 ? <EmptyState title="暂无本地通知" description="当前没有需要处理的本地应用内通知。" /> : null}
      <div className="panel-list">
        {items.map((item) => (
          <article key={item.notification_id} className="panel-card">
            <div className="row-between">
              <strong>{item.title}</strong>
              <span>{item.severity}</span>
            </div>
            <p>{item.message}</p>
            <small>{item.created_at}{item.read_at ? ` · 已读 ${item.read_at}` : ' · 未读'}</small>
            {item.source_type === 'risk_alert' && item.source_id ? (
              <div><a href={`/risk-alerts/${encodeURIComponent(item.source_id)}`}>查看风险预警</a></div>
            ) : null}
            {!item.read_at ? (
              <div>
                <Button onClick={() => markOneRead(item.notification_id)}>标记已读</Button>
              </div>
            ) : null}
          </article>
        ))}
      </div>
    </section>
  )
}
