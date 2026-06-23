package handler

import "net/http"

func (a *App) ListNotifications(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	out, err := a.NotificationSvc.ListNotifications(r.Context())
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, out)
}

func (a *App) MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	if err := a.NotificationSvc.MarkNotificationRead(r.Context(), r.PathValue("notification_id")); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, map[string]bool{"ok": true})
}

func (a *App) MarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	if err := a.NotificationSvc.MarkAllNotificationsRead(r.Context()); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, map[string]bool{"ok": true})
}
