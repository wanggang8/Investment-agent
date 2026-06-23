package model

// AuditEvent 是领域层使用的轻量审计事件摘要。
type AuditEvent struct {
	AuditEventID  string
	RequestID     string
	WorkflowType  string
	NodeName      string
	Actor         AuditActor
	Action        AuditAction
	NodeAction    string
	Status        AuditStatus
	ErrorCode     string
	RuleVersion   string
	InputRefType  string
	InputRef      string
	OutputRefType string
	OutputRef     string
}
