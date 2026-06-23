package workflow

import "context"

// DailyDisciplineGraph 编排每日纪律工作流。
// 该流程不包含用户问题能力圈判断，主要读取事实快照、分析材料、预期收益并交给领域规则裁决。
type DailyDisciplineGraph struct {
	auditWriter AuditWriter
	deps        WorkflowDependencies
}

// NewDailyDisciplineGraph 创建每日纪律工作流。
func NewDailyDisciplineGraph(writer AuditWriter) *DailyDisciplineGraph {
	if writer == nil {
		writer = &MemoryAuditWriter{}
	}
	return &DailyDisciplineGraph{auditWriter: writer}
}

// NewDailyDisciplineGraphWithDependencies 创建带 SQLite 写入能力的每日纪律工作流。
func NewDailyDisciplineGraphWithDependencies(deps WorkflowDependencies) *DailyDisciplineGraph {
	return &DailyDisciplineGraph{auditWriter: NewRepositoryAuditWriter(deps.AuditRepo), deps: deps}
}

// Run 通过 Eino Graph 执行每日纪律节点，失败节点由 AuditWriter 记录审计片段。
func (g *DailyDisciplineGraph) Run(ctx context.Context, input WorkflowContext) (WorkflowContext, error) {
	input.WorkflowType = WorkflowDailyDiscipline
	einoGraph, err := BuildDailyEinoGraph(ctx, g.auditWriter, g.deps)
	if err != nil {
		return WorkflowContext{}, err
	}
	return einoGraph.Invoke(ctx, input)
}
