package repository

import "context"

// AuditEvent 是系统关键动作的审计记录，用于串联工作流、节点、规则和错误。
type AuditEvent struct {
	AuditEventID   string
	RequestID      string
	DecisionID     string
	WorkflowType   string
	NodeName       string
	Actor          string
	Action         string
	NodeAction     string
	ProposalID     string
	ConfirmationID string
	ErrorCaseID    string
	Status         string
	ErrorCode      string
	BeforeState    string
	AfterState     string
	RuleVersion    string
	SnapshotID     string
	InputRefType   string
	InputRef       string
	OutputRefType  string
	OutputRef      string
	CreatedAt      string
}

// AuditRepository 定义审计事件的追加写入与读取边界。
type AuditRepository interface {
	AppendAuditEvent(ctx context.Context, event AuditEvent) error
	GetAuditEvent(ctx context.Context, eventID string) (AuditEvent, error)
	ListAuditEvents(ctx context.Context) ([]AuditEvent, error)
}
