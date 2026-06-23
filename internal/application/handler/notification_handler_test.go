package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNotificationHandlersListAndMarkRead(t *testing.T) {
	app, db := testApp(t)
	_, err := db.Exec(`INSERT INTO notifications (notification_id,type,severity,title,message,source_type,source_id,created_at) VALUES ('notif_1','data_source_failure','warning','数据源失败','行情数据源不可用','workflow','req_1','2026-06-01T00:00:00Z')`)
	if err != nil {
		t.Fatal(err)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/notifications", nil)
	listW := httptest.NewRecorder()
	app.ListNotifications(listW, listReq)
	if listW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", listW.Code, listW.Body.String())
	}
	if !strings.Contains(listW.Body.String(), "notif_1") || !strings.Contains(listW.Body.String(), `"unread_count":1`) {
		t.Fatalf("unexpected list body: %s", listW.Body.String())
	}

	readReq := httptest.NewRequest(http.MethodPost, "/api/v1/notifications/notif_1/read", nil)
	readReq.SetPathValue("notification_id", "notif_1")
	readW := httptest.NewRecorder()
	app.MarkNotificationRead(readW, readReq)
	if readW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", readW.Code, readW.Body.String())
	}
	var readAt string
	if err := db.QueryRow(`SELECT COALESCE(read_at,'') FROM notifications WHERE notification_id='notif_1'`).Scan(&readAt); err != nil {
		t.Fatal(err)
	}
	if readAt == "" {
		t.Fatal("expected notification read_at to be set")
	}
}

func TestNotificationHandlersMarkAllRead(t *testing.T) {
	app, db := testApp(t)
	_, err := db.Exec(`INSERT INTO notifications (notification_id,type,severity,title,message,created_at) VALUES ('notif_1','data_source_failure','warning','数据源失败','行情数据源不可用','2026-06-01T00:00:00Z'),('notif_2','rule_proposal_pending','info','规则提案待确认','有规则提案待确认','2026-06-01T00:00:01Z')`)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/notifications/read-all", nil)
	w := httptest.NewRecorder()
	app.MarkAllNotificationsRead(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var unread int
	if err := db.QueryRow(`SELECT COUNT(*) FROM notifications WHERE read_at IS NULL`).Scan(&unread); err != nil {
		t.Fatal(err)
	}
	if unread != 0 {
		t.Fatalf("expected no unread notifications, got %d", unread)
	}
}
