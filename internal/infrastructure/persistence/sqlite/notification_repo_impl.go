package sqlite

import (
	"context"
	"database/sql"

	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

// NotificationRepository 是应用内通知中心的 SQLite 实现。
type NotificationRepository struct{ db dbtx }

// NewNotificationRepository 创建通知仓储实例。
func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) SaveNotification(ctx context.Context, notification repository.Notification) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO notifications (notification_id,type,severity,title,message,source_type,source_id,read_at,created_at) VALUES (?,?,?,?,?,?,?,?,?) ON CONFLICT(type, source_type, source_id) WHERE read_at IS NULL AND source_type IS NOT NULL AND source_id IS NOT NULL DO UPDATE SET severity=excluded.severity,title=excluded.title,message=excluded.message,created_at=excluded.created_at`, notification.NotificationID, notification.Type, notification.Severity, notification.Title, notification.Message, nullString(notification.SourceType), nullString(notification.SourceID), nullString(notification.ReadAt), notification.CreatedAt)
	return apperr.FromRepositoryError(err)
}

func (r *NotificationRepository) ListNotifications(ctx context.Context) ([]repository.Notification, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT notification_id,type,severity,title,message,COALESCE(source_type,''),COALESCE(source_id,''),COALESCE(read_at,''),created_at FROM notifications ORDER BY created_at DESC, notification_id DESC`)
	if err != nil {
		return nil, apperr.FromRepositoryError(err)
	}
	defer rows.Close()
	var notifications []repository.Notification
	for rows.Next() {
		var notification repository.Notification
		if err := rows.Scan(&notification.NotificationID, &notification.Type, &notification.Severity, &notification.Title, &notification.Message, &notification.SourceType, &notification.SourceID, &notification.ReadAt, &notification.CreatedAt); err != nil {
			return nil, apperr.FromRepositoryError(err)
		}
		notifications = append(notifications, notification)
	}
	return notifications, apperr.FromRepositoryError(rows.Err())
}

func (r *NotificationRepository) MarkNotificationRead(ctx context.Context, notificationID, readAt string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE notifications SET read_at=? WHERE notification_id=? AND read_at IS NULL`, readAt, notificationID)
	return apperr.FromRepositoryError(err)
}

func (r *NotificationRepository) MarkAllNotificationsRead(ctx context.Context, readAt string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE notifications SET read_at=? WHERE read_at IS NULL`, readAt)
	return apperr.FromRepositoryError(err)
}
