import type { NotificationItem } from '../../types/notification'
import { type OpsAction, type OpsMetric, type OpsTone } from './modelTypes'

export interface NotificationSourceGroup {
  label: string
  count: number
}

export interface NotificationInboxModel {
  overallLabel: string
  overallTone: OpsTone
  metrics: OpsMetric[]
  nextActions: OpsAction[]
  sourceGroups: NotificationSourceGroup[]
  safetyNotes: string[]
}

const sourceText: Record<string, string> = {
  risk_alert: '风险预警',
  data_source_failure: '数据源',
  daily_auto_run: '每日自动运行',
  rule_proposal: '规则治理',
}

export function buildNotificationInboxModel({ notifications, unreadCount }: { notifications?: NotificationItem[]; unreadCount?: number }): NotificationInboxModel {
  const items = notifications ?? []
  const unread = unreadCount ?? items.filter((item) => !item.read_at).length
  const critical = items.filter((item) => item.severity === 'critical').length
  const warning = items.filter((item) => item.severity === 'warning').length
  const bySourceLabel = new Map<string, number>()
  for (const item of items) {
    const source = item.source_type || item.type || 'unknown'
    const label = sourceText[source] ?? '其他来源'
    bySourceLabel.set(label, (bySourceLabel.get(label) ?? 0) + 1)
  }

  return {
    overallLabel: unread || critical ? '本地通知需要处理' : items.length ? '本地通知已处理' : '暂无本地通知',
    overallTone: critical ? 'danger' : unread || warning ? 'warning' : items.length ? 'success' : 'unknown',
    metrics: [
      { label: '未读', value: String(unread), tone: unread ? 'warning' : 'success' },
      { label: '严重', value: String(critical), tone: critical ? 'danger' : 'success' },
      { label: '预警', value: String(warning), tone: warning ? 'warning' : 'success' },
      { label: '总数', value: String(items.length) },
    ],
    nextActions: buildInboxActions(unread, items),
    sourceGroups: Array.from(bySourceLabel.entries()).map(([label, count]) => ({ label, count })),
    safetyNotes: ['通知中心只处理本地应用内状态，不发送站外消息，也不代表任何交易许可。'],
  }
}

function buildInboxActions(unread: number, items: NotificationItem[]): OpsAction[] {
  const actions: OpsAction[] = []
  if (unread) actions.push({ label: '查看未读通知', detail: '先处理未读和严重通知。', href: '/notifications' })
  if (items.some((item) => item.source_type === 'risk_alert')) actions.push({ label: '查看风险预警', detail: '回到风险处置队列人工复核。', href: '/risk-alerts' })
  actions.push({ label: '查看审计记录', detail: '追踪通知来源和处理状态。', href: '/audit' })
  return actions
}
