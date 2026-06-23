export interface NotificationItem {
  notification_id: string
  type: string
  severity: 'info' | 'warning' | 'critical'
  title: string
  message: string
  source_type?: string
  source_id?: string
  read_at?: string
  created_at: string
}

export interface NotificationListResponse {
  items: NotificationItem[]
  unread_count: number
}
