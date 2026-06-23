package dto

// NotificationDTO 是应用内通知中心展示结构。
type NotificationDTO struct {
	NotificationID string `json:"notification_id"`
	Type           string `json:"type"`
	Severity       string `json:"severity"`
	Title          string `json:"title"`
	Message        string `json:"message"`
	SourceType     string `json:"source_type,omitempty"`
	SourceID       string `json:"source_id,omitempty"`
	ReadAt         string `json:"read_at,omitempty"`
	CreatedAt      string `json:"created_at"`
}

type NotificationListResponse struct {
	Items       []NotificationDTO `json:"items"`
	UnreadCount int               `json:"unread_count"`
}
