import { afterEach, describe, expect, it, vi } from 'vitest'
import { listNotifications, markAllNotificationsRead, markNotificationRead } from './notification'

afterEach(() => {
  vi.unstubAllGlobals()
  vi.clearAllMocks()
})

describe('notification service contract', () => {
  it('lists notifications with unread count', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => new Response(JSON.stringify({
      request_id: 'rid_notif_list',
      data: { unread_count: 1, items: [{ notification_id: 'notif_1', type: 'data_source_failure', severity: 'warning', title: '数据源失败', message: '行情数据源不可用', created_at: '2026-06-01T00:00:00Z' }] },
    }), { status: 200, headers: { 'content-type': 'application/json' } })))

    const res = await listNotifications()

    expect(res.data?.unread_count).toBe(1)
    expect(res.data?.items[0].notification_id).toBe('notif_1')
  })

  it('marks one notification and all notifications read', async () => {
    const fetchMock = vi.fn(async () => new Response(JSON.stringify({ request_id: 'rid_notif_read', data: { ok: true } }), { status: 200, headers: { 'content-type': 'application/json' } }))
    vi.stubGlobal('fetch', fetchMock)

    await markNotificationRead('notif_1')
    await markAllNotificationsRead()

    expect(fetchMock).toHaveBeenNthCalledWith(1, '/api/v1/notifications/notif_1/read', expect.objectContaining({ method: 'POST' }))
    expect(fetchMock).toHaveBeenNthCalledWith(2, '/api/v1/notifications/read-all', expect.objectContaining({ method: 'POST' }))
  })
})
