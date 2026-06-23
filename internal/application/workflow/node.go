package workflow

import (
	"context"
	"errors"
)

// NodeStatus 表示节点执行状态。
type NodeStatus string

const (
	// StatusSuccess 表示节点成功完成。
	StatusSuccess NodeStatus = "success"
	// StatusDegraded 表示节点降级完成，调用方可继续执行后续节点。
	StatusDegraded NodeStatus = "degraded"
	// StatusFailed 表示节点失败，通常需要写入错误码与审计事件。
	StatusFailed NodeStatus = "failed"
)

var (
	// ErrMissingAuditField 表示审计片段缺少必填字段。
	ErrMissingAuditField = errors.New("missing audit field")
	// ErrMissingErrorCode 表示失败或明确降级时缺少错误码。
	ErrMissingErrorCode = errors.New("missing error code")
)

// Node 定义工作流节点统一接口。
// 节点只读写 WorkflowContext，并通过 NodeResult 返回状态、错误码和审计片段。
type Node interface {
	Run(ctx context.Context, wf *WorkflowContext) NodeResult
}

// AuditFragment 是节点执行后产生的审计片段。
// AuditWriter 会把它转换为 audit_events 持久化记录。
type AuditFragment struct {
	Action        string
	NodeName      string
	NodeAction    string
	Status        NodeStatus
	InputRefType  string
	InputRef      string
	OutputRefType string
	OutputRef     string
	ErrorCode     string
}

// NodeResult 是节点统一返回值。
type NodeResult struct {
	Status    NodeStatus
	ErrorCode string
	Audit     AuditFragment
}

// Validate 校验节点返回值是否满足审计契约。
func (r NodeResult) Validate() error {
	if r.Audit.Action == "" || r.Audit.NodeName == "" || r.Audit.NodeAction == "" || r.Audit.Status == "" || r.Audit.InputRefType == "" || r.Audit.InputRef == "" {
		return ErrMissingAuditField
	}
	if r.Status == StatusFailed || (r.Status == StatusDegraded && r.ErrorCode != "") {
		if r.ErrorCode == "" && r.Audit.ErrorCode == "" {
			return ErrMissingErrorCode
		}
	}
	return nil
}
