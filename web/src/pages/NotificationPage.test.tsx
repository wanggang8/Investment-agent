import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { NotificationPage } from './NotificationPage'

vi.mock('../services/notification', () => ({
  listNotifications: vi.fn(),
  markAllNotificationsRead: vi.fn(),
  markNotificationRead: vi.fn(),
}))

import { listNotifications, markAllNotificationsRead, markNotificationRead } from '../services/notification'

describe('NotificationPage', () => {
  afterEach(() => {
    cleanup()
    vi.clearAllMocks()
  })

  it('polls notifications while the page is open', async () => {
    vi.mocked(listNotifications)
      .mockResolvedValueOnce({ request_id: 'rid_1', data: { unread_count: 0, items: [] } })
      .mockResolvedValue({ request_id: 'rid_2', data: { unread_count: 1, items: [{ notification_id: 'notif_poll', type: 'risk_alert', severity: 'warning', title: '新通知', message: '行情数据源不可用', source_type: 'risk_alert', source_id: 'risk_poll', created_at: '2026-06-01T00:00:00Z' }] } })

    render(<NotificationPage pollIntervalMs={10} />)

    await waitFor(() => expect(listNotifications).toHaveBeenCalledTimes(2))
    expect(screen.getByText('本地通知收件箱')).toBeInTheDocument()
    expect(screen.getByText('本地通知需要处理')).toBeInTheDocument()
    expect(screen.getByText('查看未读通知')).toBeInTheDocument()
    expect(screen.getAllByText('查看风险预警').length).toBeGreaterThan(0)
    expect(screen.getByText('未读通知：1')).toBeInTheDocument()
    expect(screen.getByText('新通知')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '查看风险预警' })).toHaveAttribute('href', '/risk-alerts/risk_poll')
  })

  it('shows unread notifications and supports marking read', async () => {
    vi.mocked(listNotifications)
      .mockResolvedValueOnce({ request_id: 'rid_1', data: { unread_count: 1, items: [{ notification_id: 'notif_1', type: 'data_source_failure', severity: 'warning', title: '数据源失败', message: '行情数据源不可用', created_at: '2026-06-01T00:00:00Z' }] } })
      .mockResolvedValue({ request_id: 'rid_2', data: { unread_count: 0, items: [{ notification_id: 'notif_1', type: 'data_source_failure', severity: 'warning', title: '数据源失败', message: '行情数据源不可用', read_at: '2026-06-01T01:00:00Z', created_at: '2026-06-01T00:00:00Z' }] } })
    vi.mocked(markNotificationRead).mockResolvedValue({ request_id: 'rid_read', data: { ok: true } })

    render(<NotificationPage />)

    expect(await screen.findByText('未读通知：1')).toBeInTheDocument()
    expect(screen.getByText('数据源失败')).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: '标记已读' }))

    await waitFor(() => expect(markNotificationRead).toHaveBeenCalledWith('notif_1'))
    await waitFor(() => expect(screen.getByText('未读通知：0')).toBeInTheDocument())
  })

  it('marks all notifications read', async () => {
    vi.mocked(listNotifications)
      .mockResolvedValueOnce({ request_id: 'rid_1', data: { unread_count: 2, items: [{ notification_id: 'notif_1', type: 'data_source_failure', severity: 'warning', title: '数据源失败', message: '行情数据源不可用', created_at: '2026-06-01T00:00:00Z' }] } })
      .mockResolvedValue({ request_id: 'rid_2', data: { unread_count: 0, items: [] } })
    vi.mocked(markAllNotificationsRead).mockResolvedValue({ request_id: 'rid_all', data: { ok: true } })

    render(<NotificationPage />)

    fireEvent.click(await screen.findByRole('button', { name: '全部标记已读' }))
    await waitFor(() => expect(markAllNotificationsRead).toHaveBeenCalled())
    await waitFor(() => expect(screen.getByText('未读通知：0')).toBeInTheDocument())
  })

  it('keeps notification copy local-only', async () => {
    vi.mocked(listNotifications).mockResolvedValue({ request_id: 'rid_1', data: { unread_count: 0, items: [] } })

    render(<NotificationPage />)

    await waitFor(() => expect(screen.getAllByText('暂无本地通知').length).toBeGreaterThan(0))
    expect(screen.getByText('通知中心只处理本地应用内状态，不发送站外消息，也不代表任何交易许可。')).toBeInTheDocument()
    expect(document.body.textContent).not.toMatch(/短信|邮件|Webhook|第三方推送|外部推送|自动确认|自动交易/)
  })
})
