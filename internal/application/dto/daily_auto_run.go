package dto

// DailyAutoRunStatusResponse 是每日自动运行状态页的后端契约。
type DailyAutoRunStatusResponse struct {
	Enabled                bool   `json:"enabled"`
	RunTime                string `json:"run_time"`
	Timezone               string `json:"timezone"`
	Scope                  string `json:"scope"`
	Status                 string `json:"status"`
	RunID                  string `json:"run_id,omitempty"`
	IdempotencyKey         string `json:"idempotency_key,omitempty"`
	LocalDate              string `json:"local_date,omitempty"`
	LastRunAt              string `json:"last_run_at,omitempty"`
	NextRunAt              string `json:"next_run_at,omitempty"`
	FailureCode            string `json:"failure_code,omitempty"`
	FailureReason          string `json:"failure_reason,omitempty"`
	LatestDecisionLink     string `json:"latest_decision_link,omitempty"`
	LatestNotificationLink string `json:"latest_notification_link,omitempty"`
	LatestAuditLink        string `json:"latest_audit_link,omitempty"`
	MissingAction          string `json:"missing_action,omitempty"`
	SafetyNote             string `json:"safety_note"`
}
