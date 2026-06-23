package workflow

import (
	"context"
	"encoding/json"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
)

// EvolutionProposalInput 是错误案例驱动的规则提案输入。
type EvolutionProposalInput struct {
	RequestID               string
	ErrorCaseID             string
	ReviewPeriod            string
	ProposalType            string
	TargetRule              string
	SupportingDecisionIDs   []string
	SupportingAuditEventIDs []string
	SampleCount             int
}

// EvolutionProposalOutput 返回规则提案，不修改正式规则版本。
type EvolutionProposalOutput struct {
	RuleProposal       model.RuleProposal
	UpdatedRuleVersion bool
	AuditEvents        []model.AuditEvent
}

// EvolutionProposalGraph 从错误案例生成规则提案。
type EvolutionProposalGraph struct {
	auditWriter AuditWriter
	deps        WorkflowDependencies
}

// NewEvolutionProposalGraph 创建规则提案工作流。
func NewEvolutionProposalGraph(writer AuditWriter) *EvolutionProposalGraph {
	if writer == nil {
		writer = &MemoryAuditWriter{}
	}
	return &EvolutionProposalGraph{auditWriter: writer}
}

// NewEvolutionProposalGraphWithDependencies 创建带 SQLite 写入能力的规则提案工作流。
func NewEvolutionProposalGraphWithDependencies(deps WorkflowDependencies) *EvolutionProposalGraph {
	return &EvolutionProposalGraph{auditWriter: NewRepositoryAuditWriter(deps.AuditRepo), deps: deps}
}

// Run 只生成 rule_proposals，不写正式 rule_versions。
func (g *EvolutionProposalGraph) Run(ctx context.Context, in EvolutionProposalInput) (EvolutionProposalOutput, error) {
	wf := WorkflowContext{RequestID: in.RequestID, WorkflowType: WorkflowEvolutionProposal, RuleVersion: workflowRuleVersion(ctx, g.deps.RuleRepo)}
	proposalID := workflowID("proposal")
	status := model.ProposalPendingUserConfirm
	riskNotesJSON := "[]"
	if in.SampleCount < 3 || in.ErrorCaseID == "" && in.ReviewPeriod == "" {
		status = model.ProposalDraft
		riskNotesJSON = `[{"code":"INSUFFICIENT_SAMPLE","message":"样本数不足或来源事实缺失，提案只能作为草稿，不能送审或应用。"}]`
	}
	proposal := model.RuleProposal{ProposalID: proposalID, Status: status, SampleCount: in.SampleCount, RuleVersion: wf.RuleVersion}
	inputRefType := "error_case"
	inputRef := in.ErrorCaseID
	if in.ReviewPeriod != "" {
		inputRefType = "review_summary"
		inputRef = in.ReviewPeriod
	}
	if inputRef == "" {
		inputRefType = "missing_source"
		inputRef = firstNonEmpty(in.RequestID, proposalID)
	}
	repoProposal := buildEvolutionRuleProposal(proposalID, wf.RuleVersion, in, status, riskNotesJSON)
	result := NodeResult{Status: StatusSuccess, Audit: AuditFragment{Action: string(model.AuditActionCreateProposal), NodeName: "EvolutionProposalGraph", NodeAction: "create_rule_proposal", Status: StatusSuccess, InputRefType: inputRefType, InputRef: inputRef, OutputRefType: "rule_proposal", OutputRef: proposalID}}
	if g.deps.Transactor != nil && g.deps.RuleRepo != nil {
		err := g.deps.Transactor.WithinTx(ctx, func(txCtx context.Context, repos repository.Repositories) error {
			if err := repos.RuleRepo.SaveRuleProposal(txCtx, repoProposal); err != nil {
				return err
			}
			if status == model.ProposalPendingUserConfirm && repos.NotificationRepo != nil {
				if err := repos.NotificationRepo.SaveNotification(txCtx, repository.Notification{NotificationID: workflowID("notif"), Type: "rule_proposal_pending", Severity: "info", Title: "规则提案待确认", Message: repoProposal.Title, SourceType: "rule_proposal", SourceID: proposalID, CreatedAt: repoProposal.CreatedAt}); err != nil {
					return err
				}
			}
			return writeAuditEvent(txCtx, repos.AuditRepo, &wf, result)
		})
		if err != nil {
			return EvolutionProposalOutput{}, err
		}
		return EvolutionProposalOutput{RuleProposal: proposal, UpdatedRuleVersion: false, AuditEvents: wf.AuditEvents}, nil
	}
	if g.deps.RuleRepo != nil {
		if err := g.deps.RuleRepo.SaveRuleProposal(ctx, repoProposal); err != nil {
			return EvolutionProposalOutput{}, err
		}
	}
	if err := g.auditWriter.Write(ctx, &wf, result); err != nil {
		return EvolutionProposalOutput{}, err
	}
	return EvolutionProposalOutput{RuleProposal: proposal, UpdatedRuleVersion: false, AuditEvents: wf.AuditEvents}, nil
}

func buildEvolutionRuleProposal(proposalID, ruleVersion string, in EvolutionProposalInput, status model.RuleProposalStatus, riskNotesJSON string) repository.RuleProposal {
	proposalSubtype := normalizedEvolutionProposalSubtype(in.ProposalType)
	proposalType := evolutionProposalDBType(proposalSubtype)
	targetRule := firstNonEmpty(in.TargetRule, defaultEvolutionTargetRule(proposalSubtype))
	beforeRuleJSON := marshalEvolutionProposalPart(struct {
		TargetRule               string `json:"target_rule"`
		ProposalSubtype          string `json:"proposal_subtype"`
		RuleVersion              string `json:"rule_version"`
		RequiredHighGradeSources int    `json:"required_high_grade_sources"`
	}{TargetRule: targetRule, ProposalSubtype: proposalSubtype, RuleVersion: ruleVersion, RequiredHighGradeSources: 2})
	afterRuleJSON := marshalEvolutionProposalPart(struct {
		TargetRule               string `json:"target_rule"`
		ProposalSubtype          string `json:"proposal_subtype"`
		RuleVersion              string `json:"rule_version"`
		RequiredHighGradeSources int    `json:"required_high_grade_sources"`
		AppliesTo                string `json:"applies_to"`
	}{TargetRule: targetRule, ProposalSubtype: proposalSubtype, RuleVersion: ruleVersion, RequiredHighGradeSources: 2, AppliesTo: "a_share_etf_fund_evidence"})
	impactScopeJSON := marshalEvolutionProposalPart(struct {
		SourceVerification      bool     `json:"source_verification"`
		ProposalType            string   `json:"proposal_type"`
		ProposalSubtype         string   `json:"proposal_subtype"`
		SupportingDecisionIDs   []string `json:"supporting_decision_ids,omitempty"`
		SupportingAuditEventIDs []string `json:"supporting_audit_event_ids,omitempty"`
	}{SourceVerification: proposalType == "risk_rule", ProposalType: proposalType, ProposalSubtype: proposalSubtype, SupportingDecisionIDs: in.SupportingDecisionIDs, SupportingAuditEventIDs: in.SupportingAuditEventIDs})
	reason := proposalSubtype + " target_rule=" + targetRule + " required_high_grade_sources"
	if in.ErrorCaseID != "" {
		reason += " error_case=" + in.ErrorCaseID
	}
	if in.ReviewPeriod != "" {
		reason += " review_period=" + in.ReviewPeriod
	}
	return repository.RuleProposal{ProposalID: proposalID, ProposalType: proposalType, Status: string(status), SourceErrorCaseID: in.ErrorCaseID, Title: evolutionProposalTitle(proposalSubtype), ProposalVersion: "draft", BeforeRuleJSON: beforeRuleJSON, AfterRuleJSON: afterRuleJSON, Reason: reason, ImpactScopeJSON: impactScopeJSON, RiskNotesJSON: riskNotesJSON, SampleCount: in.SampleCount, RelatedErrorCasesJSON: evolutionRelatedSourceJSON(in), CreatedAt: workflowNowRFC3339()}
}

func normalizedEvolutionProposalSubtype(value string) string {
	switch value {
	case "threshold_adjustment", "sop_addition", "master_weight_adjustment", "behavior_pattern_alert", "risk_rule":
		return value
	default:
		return "risk_rule"
	}
}

func evolutionProposalDBType(proposalSubtype string) string {
	switch proposalSubtype {
	case "threshold_adjustment":
		return "threshold"
	case "sop_addition":
		return "sop"
	case "master_weight_adjustment":
		return "capability"
	default:
		return "risk_rule"
	}
}

func defaultEvolutionTargetRule(proposalSubtype string) string {
	switch proposalSubtype {
	case "threshold_adjustment":
		return "valuation_threshold"
	case "sop_addition":
		return "risk_sop.evidence_insufficient"
	case "master_weight_adjustment":
		return "master.graham.margin_of_safety"
	case "behavior_pattern_alert":
		return "emotion_bias"
	default:
		return "source_verification"
	}
}

func evolutionProposalTitle(proposalType string) string {
	switch proposalType {
	case "threshold_adjustment":
		return "阈值调整提案"
	case "sop_addition":
		return "SOP 增补提案"
	case "master_weight_adjustment":
		return "大师权重调整提案"
	case "behavior_pattern_alert":
		return "个人行为模式预警提案"
	default:
		return "证据源校验规则提案"
	}
}

func marshalEvolutionProposalPart(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func evolutionRelatedSourceJSON(in EvolutionProposalInput) string {
	related := struct {
		ReviewPeriod            string   `json:"review_period,omitempty"`
		ErrorCaseID             string   `json:"error_case_id,omitempty"`
		SupportingDecisionIDs   []string `json:"supporting_decision_ids,omitempty"`
		SupportingAuditEventIDs []string `json:"supporting_audit_event_ids,omitempty"`
	}{ReviewPeriod: in.ReviewPeriod, ErrorCaseID: in.ErrorCaseID, SupportingDecisionIDs: in.SupportingDecisionIDs, SupportingAuditEventIDs: in.SupportingAuditEventIDs}
	b, err := json.Marshal(related)
	if err != nil {
		return "{}"
	}
	return string(b)
}
