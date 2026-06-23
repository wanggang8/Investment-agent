import { describe, expect, it } from 'vitest'
import { buildNotificationInboxModel } from './notificationInboxModel'
import type { NotificationItem } from '../../types/notification'

const notifications: NotificationItem[] = [
  { notification_id: 'n1', type: 'risk_alert', severity: 'critical', title: '风险预警', message: '组合风险升高', source_type: 'risk_alert', source_id: 'risk_1', created_at: '2026-06-18T08:00:00Z' },
  { notification_id: 'n2', type: 'data_source_failure', severity: 'warning', title: '数据源失败', message: '行情数据源不可用', read_at: '2026-06-18T09:00:00Z', created_at: '2026-06-18T07:00:00Z' },
]

describe('buildNotificationInboxModel', () => {
  it('summarizes local inbox state without external delivery promises', () => {
    const model = buildNotificationInboxModel({ notifications, unreadCount: 1 })

    expect(model.overallTone).toBe('danger')
    expect(model.overallLabel).toBe('本地通知需要处理')
    expect(model.metrics).toContainEqual(expect.objectContaining({ label: '未读', value: '1' }))
    expect(model.metrics).toContainEqual(expect.objectContaining({ label: '严重', value: '1' }))
    expect(model.sourceGroups.map((group) => group.label)).toContain('风险预警')
    expect(model.nextActions.map((action) => action.label)).toEqual(expect.arrayContaining(['查看未读通知', '查看风险预警']))
    expect(JSON.stringify(model)).not.toMatch(/短信|邮件|Webhook|第三方推送|外部推送|自动确认|自动交易/)
  })

  it('aggregates unknown notification sources under one stable display group', () => {
    const model = buildNotificationInboxModel({
      notifications: [
        { notification_id: 'n3', type: 'risk_alert', severity: 'warning', title: '风险', message: 'A', source_type: 'manual_daily_run', created_at: '2026-06-18T10:00:00Z' },
        { notification_id: 'n4', type: 'data_source_failure', severity: 'warning', title: '数据', message: 'B', source_type: 'p72_acceptance', created_at: '2026-06-18T11:00:00Z' },
      ],
    })

    expect(model.sourceGroups).toEqual([{ label: '其他来源', count: 2 }])
  })
})
