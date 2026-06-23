package service

import (
	"context"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
	"investment-agent/internal/pkg/clock"
	"investment-agent/internal/pkg/idgen"
)

// NotificationService handles in-app notification listing and read state.
type NotificationService struct {
	tx  repository.Transactor
	clk clock.Clock
	ids idgen.Generator
}

func NewNotificationService(tx repository.Transactor) *NotificationService {
	return &NotificationService{tx: tx, clk: clock.SystemClock{}, ids: idgen.NewGenerator()}
}

func (s *NotificationService) AppendNotification(ctx context.Context, notification repository.Notification) error {
	if notification.NotificationID == "" {
		notification.NotificationID = s.ids.New("notif")
	}
	if notification.CreatedAt == "" {
		notification.CreatedAt = s.clk.NowRFC3339()
	}
	return s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		if repos.NotificationRepo == nil {
			return apperr.New(apperr.CodeInternalError, apperr.CategoryInternal, "notification repository not configured")
		}
		return repos.NotificationRepo.SaveNotification(ctx, notification)
	})
}

func (s *NotificationService) ListNotifications(ctx context.Context) (dto.NotificationListResponse, error) {
	var out dto.NotificationListResponse
	if err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		if repos.NotificationRepo == nil {
			return apperr.New(apperr.CodeInternalError, apperr.CategoryInternal, "notification repository not configured")
		}
		notifications, err := repos.NotificationRepo.ListNotifications(ctx)
		if err != nil {
			return err
		}
		out.Items = make([]dto.NotificationDTO, 0, len(notifications))
		for _, notification := range notifications {
			if notification.ReadAt == "" {
				out.UnreadCount++
			}
			out.Items = append(out.Items, notificationDTO(notification))
		}
		return nil
	}); err != nil {
		return dto.NotificationListResponse{}, err
	}
	return out, nil
}

func (s *NotificationService) MarkNotificationRead(ctx context.Context, notificationID string) error {
	return s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		if repos.NotificationRepo == nil {
			return apperr.New(apperr.CodeInternalError, apperr.CategoryInternal, "notification repository not configured")
		}
		return repos.NotificationRepo.MarkNotificationRead(ctx, notificationID, s.clk.NowRFC3339())
	})
}

func (s *NotificationService) MarkAllNotificationsRead(ctx context.Context) error {
	return s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		if repos.NotificationRepo == nil {
			return apperr.New(apperr.CodeInternalError, apperr.CategoryInternal, "notification repository not configured")
		}
		return repos.NotificationRepo.MarkAllNotificationsRead(ctx, s.clk.NowRFC3339())
	})
}

func notificationDTO(notification repository.Notification) dto.NotificationDTO {
	return dto.NotificationDTO{NotificationID: notification.NotificationID, Type: notification.Type, Severity: notification.Severity, Title: notification.Title, Message: notification.Message, SourceType: notification.SourceType, SourceID: notification.SourceID, ReadAt: notification.ReadAt, CreatedAt: notification.CreatedAt}
}
