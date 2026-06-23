package dto

// RuleVersionDTO 是当前正式规则版本的 API 展示结构。
type RuleVersionDTO struct {
	RuleVersion string `json:"rule_version"`
	Status      string `json:"status"`
	Rules       any    `json:"rules,omitempty"`
	EffectiveAt string `json:"effective_at"`
	CreatedAt   string `json:"created_at"`
}

type RuleProposalDTO struct {
	ProposalID        string                   `json:"proposal_id"`
	ProposalType      string                   `json:"proposal_type"`
	Status            string                   `json:"status"`
	Title             string                   `json:"title"`
	ProposalVersion   string                   `json:"proposal_version"`
	Reason            string                   `json:"reason,omitempty"`
	SourceErrorCaseID string                   `json:"source_error_case_id,omitempty"`
	BeforeRule        any                      `json:"before_rule,omitempty"`
	AfterRule         any                      `json:"after_rule,omitempty"`
	ImpactScope       any                      `json:"impact_scope,omitempty"`
	RiskNotes         any                      `json:"risk_notes,omitempty"`
	AuditResult       string                   `json:"audit_result,omitempty"`
	AuditSummary      string                   `json:"audit_summary,omitempty"`
	SampleCount       int                      `json:"sample_count"`
	EffectValidation  *RuleEffectValidationDTO `json:"effect_validation,omitempty"`
	CreatedAt         string                   `json:"created_at"`
}

type RuleProposalConfirmRequest struct {
	Confirm *bool  `json:"confirm"`
	Note    string `json:"note,omitempty"`
}

type SOPAddendumProposalRequest struct {
	ScenarioKey     string `json:"scenario_key"`
	ScenarioTitle   string `json:"scenario_title"`
	OccurrenceCount int    `json:"occurrence_count"`
	SampleWindow    string `json:"sample_window"`
}

type SOPAddendumProposalResponse struct {
	ProposalID      string   `json:"proposal_id"`
	Status          string   `json:"status"`
	NotificationID  string   `json:"notification_id"`
	AuditEventIDs   []string `json:"audit_event_ids"`
	SafetyStatement string   `json:"safety_statement"`
}

type RuleEffectValidationDTO struct {
	ValidationID             string `json:"validation_id,omitempty"`
	ProposalID               string `json:"proposal_id,omitempty"`
	CandidateRuleVersion     string `json:"candidate_rule_version,omitempty"`
	ValidationStatus         string `json:"validation_status"`
	SampleCount              int    `json:"sample_count"`
	SampleWindow             string `json:"sample_window,omitempty"`
	RepresentativenessStatus string `json:"representativeness_status,omitempty"`
	OverfitRisk              string `json:"overfit_risk,omitempty"`
	ReplayResult             string `json:"replay_result,omitempty"`
	GuardrailDecision        string `json:"guardrail_decision,omitempty"`
	SourceExplanation        any    `json:"source_explanation,omitempty"`
	Metrics                  any    `json:"metrics,omitempty"`
	RiskNotes                any    `json:"risk_notes,omitempty"`
	RelatedErrorCases        any    `json:"related_error_cases,omitempty"`
	RelatedDecisionIDs       any    `json:"related_decision_ids,omitempty"`
	RelatedRiskAlertIDs      any    `json:"related_risk_alert_ids,omitempty"`
	RelatedAuditEventIDs     any    `json:"related_audit_event_ids,omitempty"`
	ValidationLink           string `json:"validation_link,omitempty"`
	SafetyNote               string `json:"safety_note,omitempty"`
	CreatedAt                string `json:"created_at,omitempty"`
	UpdatedAt                string `json:"updated_at,omitempty"`
}

type RuleEffectTrackingDTO struct {
	TrackingID           string `json:"tracking_id"`
	AppliedRuleVersion   string `json:"applied_rule_version"`
	ProposalID           string `json:"proposal_id,omitempty"`
	Period               string `json:"period"`
	HitCount             int    `json:"hit_count"`
	MisjudgmentCount     int    `json:"misjudgment_count"`
	MissingEvidenceCount int    `json:"missing_evidence_count"`
	DegradedCount        int    `json:"degraded_count"`
	RiskAlertCount       int    `json:"risk_alert_count"`
	TrendDirection       string `json:"trend_direction"`
	Metrics              any    `json:"metrics,omitempty"`
	RelatedProposalIDs   any    `json:"related_proposal_ids,omitempty"`
	RelatedAuditEventIDs any    `json:"related_audit_event_ids,omitempty"`
	RelatedRiskAlertIDs  any    `json:"related_risk_alert_ids,omitempty"`
	SafetyNote           string `json:"safety_note,omitempty"`
	CreatedAt            string `json:"created_at,omitempty"`
	UpdatedAt            string `json:"updated_at,omitempty"`
}

type RuleEffectValidationRefreshRequest struct {
	SampleWindow string `json:"sample_window"`
}

type RuleProposalConfirmResponse struct {
	ProposalID         string   `json:"proposal_id"`
	Status             string   `json:"status"`
	GatekeeperAuditID  string   `json:"gatekeeper_audit_id,omitempty"`
	AppliedRuleVersion string   `json:"applied_rule_version,omitempty"`
	CreatedRuleVersion string   `json:"created_rule_version,omitempty"`
	FinalConfirmedAt   string   `json:"final_confirmed_at,omitempty"`
	AuditEvents        []string `json:"audit_events,omitempty"`
	AuditEventIDs      []string `json:"audit_event_ids"`
}
