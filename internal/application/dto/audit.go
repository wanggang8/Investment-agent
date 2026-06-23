package dto

// AuditEventDTO 是审计时间线展示结构。
// 字段与 audit_events 表保持对应，便于前端追踪请求、决策、规则和错误案例。
type AuditEventDTO struct {
	AuditEventID   string `json:"audit_event_id"`
	EventID        string `json:"event_id"`
	RequestID      string `json:"request_id"`
	DecisionID     string `json:"decision_id,omitempty"`
	WorkflowType   string `json:"workflow_type,omitempty"`
	NodeName       string `json:"node_name,omitempty"`
	Actor          string `json:"actor"`
	Action         string `json:"action"`
	NodeAction     string `json:"node_action,omitempty"`
	ProposalID     string `json:"proposal_id,omitempty"`
	ConfirmationID string `json:"confirmation_id,omitempty"`
	ErrorCaseID    string `json:"error_case_id,omitempty"`
	Status         string `json:"status"`
	ErrorCode      string `json:"error_code,omitempty"`
	BeforeState    string `json:"before_state,omitempty"`
	AfterState     string `json:"after_state,omitempty"`
	RuleVersion    string `json:"rule_version,omitempty"`
	SnapshotID     string `json:"snapshot_id,omitempty"`
	InputRefType   string `json:"input_ref_type,omitempty"`
	InputRef       string `json:"input_ref,omitempty"`
	OutputRefType  string `json:"output_ref_type,omitempty"`
	OutputRef      string `json:"output_ref,omitempty"`
	CreatedAt      string `json:"created_at"`
}
