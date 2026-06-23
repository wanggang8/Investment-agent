package workflow

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

// EinoWorkflowGraph 封装 CloudWeGo Eino Graph 的编译结果。
type EinoWorkflowGraph struct {
	runnable            compose.Runnable[WorkflowContext, WorkflowContext]
	nodeNames           []string
	registeredNodeNames []string
}

type einoWorkflowNode struct {
	name string
	step workflowStep
}

// BuildDailyEinoGraph 构建每日纪律 Eino Graph。
func BuildDailyEinoGraph(ctx context.Context, writer AuditWriter, deps WorkflowDependencies) (*EinoWorkflowGraph, error) {
	return buildEinoGraph(ctx, writer, deps, false)
}

// BuildConsultationEinoGraph 构建主动咨询 Eino Graph。
func BuildConsultationEinoGraph(ctx context.Context, writer AuditWriter, deps WorkflowDependencies) (*EinoWorkflowGraph, error) {
	return buildEinoGraph(ctx, writer, deps, true)
}

func buildEinoGraph(ctx context.Context, writer AuditWriter, deps WorkflowDependencies, includeCapability bool) (*EinoWorkflowGraph, error) {
	if writer == nil {
		writer = NewRepositoryAuditWriter(deps.AuditRepo)
	}
	nodes := workflowNodes(includeCapability)
	graph := compose.NewGraph[WorkflowContext, WorkflowContext]()
	for _, node := range nodes {
		node := node
		if err := graph.AddLambdaNode(node.name, compose.InvokableLambda(func(ctx context.Context, input WorkflowContext) (WorkflowContext, error) {
			return runEinoWorkflowNode(ctx, input, writer, deps, node)
		})); err != nil {
			return nil, err
		}
	}
	for i, node := range nodes {
		from := compose.START
		if i > 0 {
			from = nodes[i-1].name
		}
		if err := graph.AddEdge(from, node.name); err != nil {
			return nil, err
		}
	}
	if err := graph.AddEdge(nodes[len(nodes)-1].name, compose.END); err != nil {
		return nil, err
	}
	runnable, err := graph.Compile(ctx)
	if err != nil {
		return nil, err
	}
	nodeNames := workflowNodeNames(includeCapability)
	return &EinoWorkflowGraph{runnable: runnable, nodeNames: nodeNames, registeredNodeNames: nodeNames}, nil
}

func workflowNodes(includeCapability bool) []einoWorkflowNode {
	nodes := []einoWorkflowNode{{name: "StateSnapshotNode", step: stateSnapshotStep}}
	if includeCapability {
		nodes = append(nodes, einoWorkflowNode{name: "CapabilityCheckNode", step: capabilityCheckStep})
	}
	return append(nodes,
		einoWorkflowNode{name: "EvidenceRetrievalNode", step: evidenceRetrievalStep},
		einoWorkflowNode{name: "ValueAnalystNode", step: valueAnalystStep},
		einoWorkflowNode{name: "TrendRiskOfficerNode", step: trendRiskOfficerStep},
		einoWorkflowNode{name: "ExpectedReturnNode", step: expectedReturnStep},
		einoWorkflowNode{name: "RuleArbitrationNode", step: ruleArbitrationStep},
		einoWorkflowNode{name: "DecisionRecordNode", step: decisionRecordStep},
	)
}

func runEinoWorkflowNode(ctx context.Context, input WorkflowContext, writer AuditWriter, deps WorkflowDependencies, node einoWorkflowNode) (WorkflowContext, error) {
	if shouldSkipEinoNode(input, node.name) {
		return input, nil
	}
	result := node.step(ctx, &input, deps)
	if err := writeWorkflowAudit(ctx, writer, deps, &input, result); err != nil {
		return input, err
	}
	return input, nil
}

func shouldSkipEinoNode(wf WorkflowContext, nodeName string) bool {
	if containsWorkflowError(wf.Errors, ErrCodeDataRequired) || containsWorkflowError(wf.Errors, ErrCodeDataStale) || containsWorkflowError(wf.Errors, ErrCodeRuleVersionMissing) {
		return true
	}
	if containsWorkflowError(wf.Errors, ErrCodeEvidenceNotFound) || wf.CapabilityStatus == CapabilityOutOfScope {
		return nodeName != "RuleArbitrationNode" && nodeName != "DecisionRecordNode"
	}
	return false
}

func containsWorkflowError(errors []string, want string) bool {
	for _, item := range errors {
		if item == want {
			return true
		}
	}
	return false
}

func workflowNodeNames(includeCapability bool) []string {
	names := []string{"StateSnapshotNode"}
	if includeCapability {
		names = append(names, "CapabilityCheckNode")
	}
	return append(names, "EvidenceRetrievalNode", "ValueAnalystNode", "TrendRiskOfficerNode", "ExpectedReturnNode", "RuleArbitrationNode", "DecisionRecordNode")
}

// NodeNames 返回工作流的业务节点计划。
func (g *EinoWorkflowGraph) NodeNames() []string {
	return append([]string(nil), g.nodeNames...)
}

// RegisteredNodeNames 返回真实注册到 Eino Graph 的业务节点。
func (g *EinoWorkflowGraph) RegisteredNodeNames() []string {
	return append([]string(nil), g.registeredNodeNames...)
}

// Invoke 执行已编译的 Eino Graph。
func (g *EinoWorkflowGraph) Invoke(ctx context.Context, input WorkflowContext) (WorkflowContext, error) {
	return g.runnable.Invoke(ctx, input)
}
