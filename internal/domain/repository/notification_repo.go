package repository

import "context"

// Notification 是应用内通知中心展示和已读状态的持久化记录。
type Notification struct {
	NotificationID string
	Type           string
	Severity       string
	Title          string
	Message        string
	SourceType     string
	SourceID       string
	ReadAt         string
	CreatedAt      string
}

// NotificationRepository 定义应用内通知持久化边界。
type NotificationRepository interface {
	SaveNotification(ctx context.Context, notification Notification) error
	ListNotifications(ctx context.Context) ([]Notification, error)
	MarkNotificationRead(ctx context.Context, notificationID, readAt string) error
	MarkAllNotificationsRead(ctx context.Context, readAt string) error
}
