package workflow

import (
	"context"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/idgen"
)

var auditIDGen idgen.Generator = idgen.NewGenerator()

// SetAuditIDGenerator allows tests to inject deterministic audit IDs.
func SetAuditIDGenerator(gen idgen.Generator) {
	if gen != nil {
		auditIDGen = gen
	}
}

// AuditWriter 定义节点审计片段的写入边界。
type AuditWriter interface {
	Write(ctx context.Context, wf *WorkflowContext, result NodeResult) error
}

// MemoryAuditWriter 是内存审计写入器，供工作流单元测试使用。
type MemoryAuditWriter struct{}

// Write 校验节点结果，并将完整审计片段追加到 WorkflowContext。
func (w *MemoryAuditWriter) Write(_ context.Context, wf *WorkflowContext, result NodeResult) error {
	if err := result.Validate(); err != nil {
		return err
	}
	event := buildDomainAuditEvent(wf, result)
	if event.Status == model.AuditStatusFailed || event.Status == model.AuditStatusDegraded {
		if result.ErrorCode != "" {
			wf.Errors = append(wf.Errors, result.ErrorCode)
		} else if result.Audit.ErrorCode != "" {
			wf.Errors = append(wf.Errors, result.Audit.ErrorCode)
		}
	}
	wf.AuditEvents = append(wf.AuditEvents, event)
	return nil
}

// RepositoryAuditWriter 把节点审计写入 SQLite audit_events，并同步追加到上下文。
type RepositoryAuditWriter struct{ repo repository.AuditRepository }

// NewRepositoryAuditWriter 创建持久化审计写入器。
func NewRepositoryAuditWriter(repo repository.AuditRepository) AuditWriter {
	return &RepositoryAuditWriter{repo: repo}
}

func writeAuditEvent(ctx context.Context, repo repository.AuditRepository, wf *WorkflowContext, result NodeResult) error {
	if repo == nil {
		return (&MemoryAuditWriter{}).Write(ctx, wf, result)
	}
	if err := result.Validate(); err != nil {
		return err
	}
	if err := repo.AppendAuditEvent(ctx, buildRepositoryAuditEvent(wf, result)); err != nil {
		return err
	}
	return (&MemoryAuditWriter{}).Write(ctx, wf, result)
}

// Write 校验节点审计片段，并写入 audit_events。
func (w *RepositoryAuditWriter) Write(ctx context.Context, wf *WorkflowContext, result NodeResult) error {
	return writeAuditEvent(ctx, w.repo, wf, result)
}

func buildDomainAuditEvent(wf *WorkflowContext, result NodeResult) model.AuditEvent {
	return model.AuditEvent{
		AuditEventID:  auditIDGen.New("audit"),
		RequestID:     wf.RequestID,
		WorkflowType:  wf.WorkflowType,
		NodeName:      result.Audit.NodeName,
		Actor:         model.AuditActorSystem,
		Action:        auditAction(result.Audit.Action),
		NodeAction:    result.Audit.NodeAction,
		Status:        auditStatus(result.Audit.Status),
		ErrorCode:     result.ErrorCode,
		RuleVersion:   wf.RuleVersion,
		InputRefType:  result.Audit.InputRefType,
		InputRef:      result.Audit.InputRef,
		OutputRefType: result.Audit.OutputRefType,
		OutputRef:     result.Audit.OutputRef,
	}
}

func buildRepositoryAuditEvent(wf *WorkflowContext, result NodeResult) repository.AuditEvent {
	code := result.Audit.ErrorCode
	if code == "" {
		code = result.ErrorCode
	}
		event := repository.AuditEvent{
		AuditEventID:  auditIDGen.New("audit"),
		RequestID:     wf.RequestID,
		DecisionID:    wf.DecisionID,
		WorkflowType:  wf.WorkflowType,
		NodeName:      result.Audit.NodeName,
		Actor:         string(model.AuditActorSystem),
		Action:        result.Audit.Action,
		NodeAction:    result.Audit.NodeAction,
		Status:        string(auditStatus(result.Audit.Status)),
		ErrorCode:     code,
		RuleVersion:   wf.RuleVersion,
		InputRefType:  result.Audit.InputRefType,
		InputRef:      result.Audit.InputRef,
		OutputRefType: result.Audit.OutputRefType,
		OutputRef:     result.Audit.OutputRef,
		CreatedAt:     workflowNowRFC3339(),
	}
	if result.Audit.InputRefType == "rule_proposal" {
		event.ProposalID = result.Audit.InputRef
	}
	return event
}
