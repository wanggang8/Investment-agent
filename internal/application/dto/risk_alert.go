package dto

// RiskAlertDTO exposes a local P35 risk alert and SOP status.
type RiskAlertDTO struct {
	AlertID               string   `json:"alert_id"`
	RiskType              string   `json:"risk_type"`
	Severity              string   `json:"severity"`
	SOPStatus             string   `json:"sop_status"`
	Symbol                string   `json:"symbol"`
	TriggerSummary        string   `json:"trigger_summary"`
	TriggerContext        any      `json:"trigger_context,omitempty"`
	ProhibitedActions     []string `json:"prohibited_actions"`
	SuggestedActions      []string `json:"suggested_actions"`
	RelatedDecisionID     string   `json:"related_decision_id,omitempty"`
	RelatedReportID       string   `json:"related_report_id,omitempty"`
	RelatedNotificationID string   `json:"related_notification_id,omitempty"`
	RelatedAuditEventID   string   `json:"related_audit_event_id,omitempty"`
	DecisionLink          string   `json:"decision_link,omitempty"`
	ReportLink            string   `json:"report_link,omitempty"`
	NotificationLink      string   `json:"notification_link,omitempty"`
	AuditLink             string   `json:"audit_link,omitempty"`
	Link                  string   `json:"link"`
	LastTriggeredAt       string   `json:"last_triggered_at,omitempty"`
	ResolvedAt            string   `json:"resolved_at,omitempty"`
	ResolutionReason      string   `json:"resolution_reason,omitempty"`
	SafetyNote            string   `json:"safety_note"`
	CreatedAt             string   `json:"created_at"`
	UpdatedAt             string   `json:"updated_at"`
}

// RiskAlertLifecycleRequest updates a local risk alert SOP status.
type RiskAlertLifecycleRequest struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}
