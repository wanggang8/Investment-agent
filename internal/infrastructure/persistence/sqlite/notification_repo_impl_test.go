package sqlite

import (
	"context"
	"testing"

	"investment-agent/internal/domain/repository"
)

func TestNotificationRepositoryWriteListAndMarkRead(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	repo := NewNotificationRepository(db)

	notification := repository.Notification{NotificationID: "notif_1", Type: "data_source_failure", Severity: "warning", Title: "数据源失败", Message: "行情数据源不可用", SourceType: "workflow", SourceID: "req_1", CreatedAt: testTime}
	if err := repo.SaveNotification(ctx, notification); err != nil {
		t.Fatal(err)
	}
	duplicate := notification
	duplicate.NotificationID = "notif_2"
	duplicate.Message = "重复告警"
	if err := repo.SaveNotification(ctx, duplicate); err != nil {
		t.Fatal(err)
	}

	items, err := repo.ListNotifications(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].NotificationID != "notif_1" || items[0].ReadAt != "" || items[0].Message != "重复告警" {
		t.Fatalf("unexpected refreshed active notification: %+v", items)
	}
	if err := repo.MarkNotificationRead(ctx, "notif_1", "2026-06-01T00:00:00Z"); err != nil {
		t.Fatal(err)
	}
	items, err = repo.ListNotifications(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].ReadAt != "2026-06-01T00:00:00Z" {
		t.Fatalf("expected notification marked read: %+v", items)
	}
}

func TestNotificationRepositoryMarkAllRead(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	repo := NewNotificationRepository(db)
	items := []repository.Notification{
		{NotificationID: "notif_1", Type: "data_source_failure", Severity: "warning", Title: "数据源失败", Message: "行情数据源不可用", CreatedAt: testTime},
		{NotificationID: "notif_2", Type: "rule_proposal_pending", Severity: "info", Title: "规则提案待确认", Message: "有规则提案待确认", CreatedAt: testTime},
	}
	for _, item := range items {
		if err := repo.SaveNotification(ctx, item); err != nil {
			t.Fatal(err)
		}
	}
	if err := repo.MarkAllNotificationsRead(ctx, "2026-06-01T00:00:00Z"); err != nil {
		t.Fatal(err)
	}
	got, err := repo.ListNotifications(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].ReadAt == "" || got[1].ReadAt == "" {
		t.Fatalf("expected all notifications marked read: %+v", got)
	}
}
