package dto

// DailyDisciplineReportResponse exposes the local daily discipline report for UI review.
type DailyDisciplineReportResponse struct {
	ReportID          string                           `json:"report_id"`
	LocalDate         string                           `json:"local_date"`
	Scope             string                           `json:"scope"`
	Status            string                           `json:"status"`
	Summary           string                           `json:"summary"`
	SourceType        string                           `json:"source_type"`
	SourceID          string                           `json:"source_id"`
	DecisionID        string                           `json:"decision_id"`
	DecisionLink      string                           `json:"decision_link"`
	AutoRunLink       string                           `json:"auto_run_link"`
	AuditLink         string                           `json:"audit_link"`
	NotificationLink  string                           `json:"notification_link"`
	FailureCode       string                           `json:"failure_code"`
	FailureReason     string                           `json:"failure_reason"`
	MissingAction     string                           `json:"missing_action"`
	MissingCategories []string                         `json:"missing_categories"`
	FinalVerdict      string                           `json:"final_verdict"`
	VerdictStatus     string                           `json:"verdict_status"`
	Evidence          DailyDisciplineReportEvidence    `json:"evidence"`
	P34SourceCoverage DailyDisciplineReportP34Coverage `json:"p34_source_coverage"`
	RiskAlerts        []DailyDisciplineReportRiskAlert `json:"risk_alerts"`
	Trend             DailyDisciplineReportTrend       `json:"trend"`
	SafetyNote        string                           `json:"safety_note"`
	UpdatedAt         string                           `json:"updated_at"`
}

// DailyDisciplineReportEvidence summarizes evidence coverage for the linked decision.
type DailyDisciplineReportEvidence struct {
	EvidenceCount                   int `json:"evidence_count"`
	IndependentSourceCount          int `json:"independent_source_count"`
	HighGradeIndependentSourceCount int `json:"high_grade_independent_source_count"`
}

// DailyDisciplineReportP34Coverage summarizes expanded public data coverage used by the report.
type DailyDisciplineReportP34Coverage struct {
	Summary           string             `json:"summary"`
	MissingCategories []string           `json:"missing_categories"`
	SourceHealth      []SourceHealthItem `json:"source_health"`
}

// DailyDisciplineReportTrend summarizes recent report statuses.
type DailyDisciplineReportTrend struct {
	SuccessCount          int `json:"success_count"`
	DegradedCount         int `json:"degraded_count"`
	FailedCount           int `json:"failed_count"`
	InsufficientDataCount int `json:"insufficient_data_count"`
}

// DailyDisciplineReportRiskAlert summarizes a related P35 risk alert for report surfaces.
type DailyDisciplineReportRiskAlert struct {
	AlertID           string   `json:"alert_id"`
	RiskType          string   `json:"risk_type"`
	Severity          string   `json:"severity"`
	SOPStatus         string   `json:"sop_status"`
	Symbol            string   `json:"symbol"`
	TriggerSummary    string   `json:"trigger_summary"`
	ProhibitedActions []string `json:"prohibited_actions"`
	SuggestedActions  []string `json:"suggested_actions"`
	Link              string   `json:"link"`
	SafetyNote        string   `json:"safety_note"`
}

// DailyDisciplineReportListResponse wraps daily discipline report list results.
type DailyDisciplineReportListResponse struct {
	Reports []DailyDisciplineReportResponse `json:"reports"`
}
