package workflow

import (
	"context"
	"strconv"
	"strings"

	"github.com/cloudwego/eino/compose"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/domain/rule"
)

// GatekeeperAuditInput 是守门人审计输入。
type GatekeeperAuditInput struct {
	RequestID   string
	ProposalID  string
	Approved    bool
	AuditResult model.AuditResult
}

// GatekeeperAuditOutput 表示审计产物和提案状态。
type GatekeeperAuditOutput struct {
	GatekeeperAudits   []string
	ProposalStatus     model.RuleProposalStatus
	UpdatedRuleVersion bool
	AuditEvents        []model.AuditEvent
}

// GatekeeperAuditGraph 只生成 gatekeeper_audits，不写正式规则。
type GatekeeperAuditGraph struct {
	auditWriter AuditWriter
	deps        WorkflowDependencies
	nodeNames   []string
}

type gatekeeperAuditState struct {
	Input          GatekeeperAuditInput
	Workflow       WorkflowContext
	Proposal       repository.RuleProposal
	AuditResult    model.AuditResult
	AuditID        string
	Audit          repository.GatekeeperAudit
	ProposalStatus model.RuleProposalStatus
}

// NewGatekeeperAuditGraph 创建守门人审计工作流。
func NewGatekeeperAuditGraph(writer AuditWriter) *GatekeeperAuditGraph {
	if writer == nil {
		writer = &MemoryAuditWriter{}
	}
	return &GatekeeperAuditGraph{auditWriter: writer, nodeNames: gatekeeperAuditNodeNames()}
}

// NewGatekeeperAuditGraphWithDependencies 创建带 SQLite 写入能力的守门人审计工作流。
func NewGatekeeperAuditGraphWithDependencies(deps WorkflowDependencies) *GatekeeperAuditGraph {
	return &GatekeeperAuditGraph{auditWriter: NewRepositoryAuditWriter(deps.AuditRepo), deps: deps, nodeNames: gatekeeperAuditNodeNames()}
}

func gatekeeperAuditNodeNames() []string {
	return []string{"ProposalLoadNode", "FundamentalRuleCheckNode", "ConflictCheckNode", "BacktestNode", "AuditDecisionNode", "AuditRecordNode"}
}

func (g *GatekeeperAuditGraph) NodeNames() []string {
	return append([]string(nil), g.nodeNames...)
}

func (g *GatekeeperAuditGraph) RegisteredNodeNames() []string {
	return append([]string(nil), g.nodeNames...)
}

func gatekeeperResult(in GatekeeperAuditInput) model.AuditResult {
	if in.AuditResult != "" {
		return in.AuditResult
	}
	if in.Approved {
		return model.AuditApproved
	}
	return model.AuditRejected
}

func gatekeeperAuditForProposal(proposal repository.RuleProposal, auditResult model.AuditResult, ruleVersion string) repository.GatekeeperAudit {
	violatesFundamental := containsUnsafeRuleChange(proposal.AfterRuleJSON)
	hasConflict := strings.TrimSpace(proposal.BeforeRuleJSON) == strings.TrimSpace(proposal.AfterRuleJSON) || strings.Contains(strings.ToLower(proposal.Reason), "conflict") || violatesFundamental
	backtestPassed := proposal.SampleCount >= 3
	allowApply := auditResult == model.AuditApproved && !violatesFundamental && !hasConflict && backtestPassed
	reasonParts := []string{
		"FundamentalRuleCheck: " + passFail(!violatesFundamental),
		"ConflictCheck: " + passFail(!hasConflict),
		"Backtest: sample_count " + passFail(backtestPassed),
		"AuditDecision: " + string(auditResult),
	}
	requiredChanges := ""
	if auditResult == model.AuditRejected {
		requiredChanges = "根据守门人意见调整规则后重新送审"
	}
	if auditResult == model.AuditNeedsUserReview {
		requiredChanges = "需要用户复核规则变化后重新送审"
	}
	if violatesFundamental {
		requiredChanges = appendRequiredChange(requiredChanges, "移除自动交易、主动荐股或收益承诺相关规则")
	}
	if hasConflict {
		requiredChanges = appendRequiredChange(requiredChanges, "解决与既有规则冲突或无效变更")
	}
	if !backtestPassed {
		requiredChanges = appendRequiredChange(requiredChanges, "补充足够回测样本后重新送审")
	}
	backtestMetrics := `{"sample_count":` + strconv.Itoa(proposal.SampleCount) + `,"min_sample_count":3,"passed":` + strconv.FormatBool(backtestPassed) + `}`
	return repository.GatekeeperAudit{ProposalID: proposal.ProposalID, AuditResult: string(auditResult), AuditReason: strings.Join(reasonParts, "; "), RequiredChanges: requiredChanges, ViolatesFundamentalRule: violatesFundamental, HasRuleConflict: hasConflict, BacktestMetricsJSON: backtestMetrics, AllowApply: allowApply, AuditedRuleVersion: ruleVersion}
}

func containsUnsafeRuleChange(text string) bool {
	lower := strings.ToLower(text)
	for _, keyword := range []string{"auto_trade", "broker", "guaranteed_return", "active_recommendation", "自动下单", "交易接口", "收益承诺", "主动荐股"} {
		if strings.Contains(lower, keyword) {
			return true
		}
	}
	return false
}

func passFail(ok bool) string {
	if ok {
		return "passed"
	}
	return "failed"
}

func appendRequiredChange(current, item string) string {
	if current == "" {
		return item
	}
	return current + "；" + item
}

func validateGatekeeperProposal(proposal repository.RuleProposal) error {
	if proposal.Status != string(model.ProposalUnderGatekeeperAudit) {
		return rule.ErrInvalidTransition
	}
	if proposal.SampleCount < 3 {
		return rule.ErrInsufficientSamples
	}
	return nil
}

func workflowRuleVersion(ctx context.Context, repo repository.RuleRepository) string {
	if repo == nil {
		return ""
	}
	version, err := repo.GetActiveRuleVersion(ctx)
	if err != nil {
		return ""
	}
	return version.RuleVersion
}

func activeRuleVersion(ctx context.Context, repo repository.RuleRepository) string {
	return workflowRuleVersion(ctx, repo)
}

// Run 审计通过后把提案推进到 pending_final_confirm，等待用户最终确认。
func (g *GatekeeperAuditGraph) Run(ctx context.Context, in GatekeeperAuditInput) (GatekeeperAuditOutput, error) {
	runnable, err := g.compile(ctx)
	if err != nil {
		return GatekeeperAuditOutput{}, err
	}
	state, err := runnable.Invoke(ctx, gatekeeperAuditState{Input: in})
	if err != nil {
		return GatekeeperAuditOutput{}, err
	}
	return GatekeeperAuditOutput{GatekeeperAudits: []string{state.AuditID}, ProposalStatus: state.ProposalStatus, UpdatedRuleVersion: false, AuditEvents: state.Workflow.AuditEvents}, nil
}

func (g *GatekeeperAuditGraph) compile(ctx context.Context) (compose.Runnable[gatekeeperAuditState, gatekeeperAuditState], error) {
	graph := compose.NewGraph[gatekeeperAuditState, gatekeeperAuditState]()
	nodes := []struct {
		name string
		fn   func(context.Context, gatekeeperAuditState) (gatekeeperAuditState, error)
	}{
		{name: "ProposalLoadNode", fn: g.proposalLoadNode},
		{name: "FundamentalRuleCheckNode", fn: g.fundamentalRuleCheckNode},
		{name: "ConflictCheckNode", fn: g.conflictCheckNode},
		{name: "BacktestNode", fn: g.backtestNode},
		{name: "AuditDecisionNode", fn: g.auditDecisionNode},
		{name: "AuditRecordNode", fn: g.auditRecordNode},
	}
	for _, node := range nodes {
		node := node
		if err := graph.AddLambdaNode(node.name, compose.InvokableLambda(node.fn)); err != nil {
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
	return graph.Compile(ctx)
}

func (g *GatekeeperAuditGraph) proposalLoadNode(ctx context.Context, state gatekeeperAuditState) (gatekeeperAuditState, error) {
	state.Workflow = WorkflowContext{RequestID: state.Input.RequestID, WorkflowType: WorkflowGatekeeperAudit, RuleVersion: activeRuleVersion(ctx, g.deps.RuleRepo)}
	state.AuditResult = gatekeeperResult(state.Input)
	state.AuditID = workflowID("gatekeeper")
	if g.deps.RuleRepo != nil {
		proposal, err := g.deps.RuleRepo.GetRuleProposal(ctx, state.Input.ProposalID)
		if err != nil {
			return state, err
		}
		if err := validateGatekeeperProposal(proposal); err != nil {
			return state, err
		}
		state.Proposal = proposal
	} else {
		state.Proposal = repository.RuleProposal{ProposalID: state.Input.ProposalID, Status: string(model.ProposalUnderGatekeeperAudit), SampleCount: 3, BeforeRuleJSON: `{"rule":"before"}`, AfterRuleJSON: `{"rule":"after"}`}
	}
	return state, g.writeGatekeeperAudit(ctx, &state.Workflow, gatekeeperNodeResult("ProposalLoadNode", "load_proposal", state.Input.ProposalID, state.Input.ProposalID, StatusSuccess, ""))
}

func (g *GatekeeperAuditGraph) fundamentalRuleCheckNode(ctx context.Context, state gatekeeperAuditState) (gatekeeperAuditState, error) {
	return state, g.writeGatekeeperAudit(ctx, &state.Workflow, gatekeeperNodeResult("FundamentalRuleCheckNode", "check_fundamental_rules", state.Input.ProposalID, state.Input.ProposalID, StatusSuccess, ""))
}

func (g *GatekeeperAuditGraph) conflictCheckNode(ctx context.Context, state gatekeeperAuditState) (gatekeeperAuditState, error) {
	return state, g.writeGatekeeperAudit(ctx, &state.Workflow, gatekeeperNodeResult("ConflictCheckNode", "check_rule_conflict", state.Input.ProposalID, state.Input.ProposalID, StatusSuccess, ""))
}

func (g *GatekeeperAuditGraph) backtestNode(ctx context.Context, state gatekeeperAuditState) (gatekeeperAuditState, error) {
	return state, g.writeGatekeeperAudit(ctx, &state.Workflow, gatekeeperNodeResult("BacktestNode", "check_backtest_samples", state.Input.ProposalID, state.Input.ProposalID, StatusSuccess, ""))
}

func (g *GatekeeperAuditGraph) auditDecisionNode(ctx context.Context, state gatekeeperAuditState) (gatekeeperAuditState, error) {
	state.Audit = gatekeeperAuditForProposal(state.Proposal, state.AuditResult, state.Workflow.RuleVersion)
	state.Audit.GatekeeperAuditID = state.AuditID
	state.Audit.CreatedAt = workflowNowRFC3339()
	status, _, err := rule.AdvanceProposal(model.RuleProposalStatus(state.Proposal.Status), state.Proposal.SampleCount, true, state.AuditResult)
	if err != nil {
		return state, err
	}
	if !state.Audit.AllowApply {
		status = model.ProposalRejected
	}
	state.ProposalStatus = status
	return state, g.writeGatekeeperAudit(ctx, &state.Workflow, gatekeeperNodeResult("AuditDecisionNode", "decide_gatekeeper_audit", state.Input.ProposalID, state.AuditID, StatusSuccess, ""))
}

func (g *GatekeeperAuditGraph) auditRecordNode(ctx context.Context, state gatekeeperAuditState) (gatekeeperAuditState, error) {
	if g.deps.Transactor != nil && g.deps.RuleRepo != nil {
		err := g.deps.Transactor.WithinTx(ctx, func(txCtx context.Context, repos repository.Repositories) error {
			if err := repos.RuleRepo.SaveGatekeeperAudit(txCtx, state.Audit); err != nil {
				return err
			}
			if err := repos.RuleRepo.UpdateRuleProposalStatus(txCtx, state.Input.ProposalID, string(state.ProposalStatus)); err != nil {
				return err
			}
			return writeAuditEvent(txCtx, repos.AuditRepo, &state.Workflow, gatekeeperNodeResult("AuditRecordNode", "record_gatekeeper_audit", state.Input.ProposalID, state.AuditID, StatusSuccess, ""))
		})
		return state, err
	}
	if g.deps.RuleRepo != nil {
		if err := g.deps.RuleRepo.SaveGatekeeperAudit(ctx, state.Audit); err != nil {
			return state, err
		}
		if err := g.deps.RuleRepo.UpdateRuleProposalStatus(ctx, state.Input.ProposalID, string(state.ProposalStatus)); err != nil {
			return state, err
		}
	}
	return state, g.writeGatekeeperAudit(ctx, &state.Workflow, gatekeeperNodeResult("AuditRecordNode", "record_gatekeeper_audit", state.Input.ProposalID, state.AuditID, StatusSuccess, ""))
}

func gatekeeperNodeResult(nodeName, nodeAction, inputRef, outputRef string, status NodeStatus, code string) NodeResult {
	return NodeResult{Status: status, ErrorCode: code, Audit: AuditFragment{Action: string(model.AuditActionAuditRuleChange), NodeName: nodeName, NodeAction: nodeAction, Status: status, InputRefType: "rule_proposal", InputRef: inputRef, OutputRefType: "gatekeeper_audit", OutputRef: outputRef, ErrorCode: code}}
}

func (g *GatekeeperAuditGraph) writeGatekeeperAudit(ctx context.Context, wf *WorkflowContext, result NodeResult) error {
	writer := g.auditWriter
	if writer == nil {
		writer = &MemoryAuditWriter{}
	}
	return writer.Write(ctx, wf, result)
}
