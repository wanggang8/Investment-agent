import type { NotificationListResponse } from '../types/notification'
import { apiRequest } from './client'

export function listNotifications() {
  return apiRequest<NotificationListResponse>('/api/v1/notifications')
}

export function markNotificationRead(notificationId: string) {
  return apiRequest<{ ok: boolean }>(`/api/v1/notifications/${notificationId}/read`, {
    method: 'POST',
  })
}

export function markAllNotificationsRead() {
  return apiRequest<{ ok: boolean }>('/api/v1/notifications/read-all', {
    method: 'POST',
  })
}
