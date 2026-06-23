package dto

// ReviewSummaryResponse 是复盘页摘要，聚合错误案例、规则演进、审计和周期复盘概览。
type ReviewSummaryResponse struct {
	Period                string                    `json:"period"`
	DecisionCount         int                       `json:"decision_count"`
	ConfirmationCount     int                       `json:"confirmation_count"`
	ExecutedManuallyCount int                       `json:"executed_manually_count"`
	PlannedCount          int                       `json:"planned_count"`
	ErrorCaseCount        int                       `json:"error_case_count"`
	RuleProposalCount     int                       `json:"rule_proposal_count"`
	AuditEventCount       int                       `json:"audit_event_count"`
	RuleHitCount          int                       `json:"rule_hit_count"`
	MisjudgmentCount      int                       `json:"misjudgment_count"`
	MissingEvidenceCount  int                       `json:"missing_evidence_count"`
	DegradedCount         int                       `json:"degraded_count"`
	OpsStatus             ReviewOpsStatus           `json:"ops_status"`
	RuleSuggestions       []RuleSuggestionDTO       `json:"rule_suggestions"`
	TrackingLinks         []ReviewTrackingLink      `json:"tracking_links"`
	RecentDecisions       []DecisionListItem        `json:"recent_decisions"`
	AttributionSummaries  []ReviewAttribution       `json:"attribution_summaries"`
	RecurringErrorTags    []ReviewErrorTag          `json:"recurring_error_tags"`
	MissingEvidenceThemes []ReviewEvidenceTheme     `json:"missing_evidence_themes"`
	RuleProposalOutcomes  []ReviewRuleProposalState `json:"rule_proposal_outcomes"`
	RuleEffectTracking    []RuleEffectTrackingDTO   `json:"rule_effect_tracking"`
	DegradedWorkflows     []ReviewDegradedWorkflow  `json:"degraded_workflows"`
}

type ReviewOpsStatus struct {
	DataSourceStatus string `json:"data_source_status,omitempty"`
	IndexStatus      string `json:"index_status,omitempty"`
	ReviewStatus     string `json:"review_status,omitempty"`
	Explanation      string `json:"explanation,omitempty"`
}

// ReviewAttribution 是从本地决策、确认和证据状态归纳的可追溯归因项。
type ReviewAttribution struct {
	DecisionID         string `json:"decision_id"`
	Symbol             string `json:"symbol,omitempty"`
	Verdict            string `json:"verdict,omitempty"`
	ConfirmationStatus string `json:"confirmation_status,omitempty"`
	EvidenceStatus     string `json:"evidence_status,omitempty"`
	WorkflowStatus     string `json:"workflow_status,omitempty"`
	Outcome            string `json:"outcome"`
}

type ReviewErrorTag struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

type ReviewEvidenceTheme struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type ReviewRuleProposalState struct {
	ProposalID  string `json:"proposal_id"`
	Title       string `json:"title"`
	Status      string `json:"status"`
	AuditResult string `json:"audit_result,omitempty"`
}

type ReviewDegradedWorkflow struct {
	DecisionID string `json:"decision_id"`
	Symbol     string `json:"symbol,omitempty"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}

type RuleSuggestionDTO struct {
	ProposalID   string `json:"proposal_id"`
	Title        string `json:"title"`
	Status       string `json:"status"`
	Reason       string `json:"reason,omitempty"`
	CanAutoApply bool   `json:"can_auto_apply"`
}

// ReviewTrackingLink 指向复盘相关的审计、规则提案或错误案例记录。
type ReviewTrackingLink struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	Label string `json:"label"`
}
