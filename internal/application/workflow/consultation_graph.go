package workflow

import "context"

// ConsultationGraph 编排用户主动咨询工作流。
// 该流程在状态快照后加入能力圈检查，能力圈外时最终裁决必须由领域规则给出 rejected。
type ConsultationGraph struct {
	auditWriter AuditWriter
	deps        WorkflowDependencies
}

// NewConsultationGraph 创建主动咨询工作流。
func NewConsultationGraph(writer AuditWriter) *ConsultationGraph {
	if writer == nil {
		writer = &MemoryAuditWriter{}
	}
	return &ConsultationGraph{auditWriter: writer}
}

// NewConsultationGraphWithDependencies 创建带 SQLite 写入能力的主动咨询工作流。
func NewConsultationGraphWithDependencies(deps WorkflowDependencies) *ConsultationGraph {
	return &ConsultationGraph{auditWriter: NewRepositoryAuditWriter(deps.AuditRepo), deps: deps}
}

// Run 通过 Eino Graph 执行主动咨询节点，并保留所有节点审计事件。
func (g *ConsultationGraph) Run(ctx context.Context, input WorkflowContext) (WorkflowContext, error) {
	input.WorkflowType = WorkflowConsultation
	einoGraph, err := BuildConsultationEinoGraph(ctx, g.auditWriter, g.deps)
	if err != nil {
		return WorkflowContext{}, err
	}
	return einoGraph.Invoke(ctx, input)
}
