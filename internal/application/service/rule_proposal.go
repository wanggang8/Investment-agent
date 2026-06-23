package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/application/workflow"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/domain/rule"
	"investment-agent/internal/pkg/apperr"
	"investment-agent/internal/pkg/clock"
	"investment-agent/internal/pkg/idgen"
)

// RuleProposalService 处理规则提案确认和最终确认。
type RuleProposalService struct {
	tx   repository.Transactor
	deps workflow.WorkflowDependencies
	clk  clock.Clock
	ids  idgen.Generator
}

// NewRuleProposalService 创建规则提案服务。
func NewRuleProposalService(tx repository.Transactor, deps ...workflow.WorkflowDependencies) *RuleProposalService {
	var wfDeps workflow.WorkflowDependencies
	if len(deps) > 0 {
		wfDeps = deps[0]
	}
	return &RuleProposalService{tx: tx, deps: wfDeps, clk: clock.SystemClock{}, ids: idgen.NewGenerator()}
}

func (s *RuleProposalService) GenerateSOPAddendumProposal(ctx context.Context, requestID string, req dto.SOPAddendumProposalRequest) (dto.SOPAddendumProposalResponse, error) {
	if strings.TrimSpace(req.ScenarioKey) == "" || strings.TrimSpace(req.ScenarioTitle) == "" || strings.TrimSpace(req.SampleWindow) == "" {
		return dto.SOPAddendumProposalResponse{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "scenario_key、scenario_title 和 sample_window 不能为空")
	}
	if req.OccurrenceCount < 3 {
		return dto.SOPAddendumProposalResponse{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "高频未覆盖场景至少需要 3 次样本")
	}
	now := s.clk.NowRFC3339()
	proposalID := s.ids.New("prop")
	notificationID := s.ids.New("notif")
	auditID := s.ids.New("audit")
	beforeRule := map[string]any{"sop_addendum": "none", "scenario_key": req.ScenarioKey}
	afterRule := map[string]any{"proposal_subtype": "sop_addendum", "scenario_key": req.ScenarioKey, "scenario_title": req.ScenarioTitle, "manual_gate_required": true, "auto_apply": false}
	impact := map[string]any{"scope": "sop_addendum_proposal", "sample_window": req.SampleWindow, "occurrence_count": req.OccurrenceCount}
	riskNotes := []string{"SOP 补充只生成待确认提案", "需用户确认、守门人审计和最终确认后才可能应用", "不会自动应用规则或创建交易动作"}
	beforeJSON, _ := json.Marshal(beforeRule)
	afterJSON, _ := json.Marshal(afterRule)
	impactJSON, _ := json.Marshal(impact)
	riskJSON, _ := json.Marshal(riskNotes)
	if err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		if err := repos.RuleRepo.SaveRuleProposal(ctx, repository.RuleProposal{ProposalID: proposalID, ProposalType: "sop", Status: string(model.ProposalPendingUserConfirm), Title: "SOP 补充提案：" + req.ScenarioTitle, ProposalVersion: "p88-sop-addendum-" + proposalID, BeforeRuleJSON: string(beforeJSON), AfterRuleJSON: string(afterJSON), Reason: "高频未覆盖场景：" + req.ScenarioTitle, ImpactScopeJSON: string(impactJSON), RiskNotesJSON: string(riskJSON), SampleCount: req.OccurrenceCount, CreatedAt: now}); err != nil {
			return err
		}
		if err := repos.NotificationRepo.SaveNotification(ctx, repository.Notification{NotificationID: notificationID, Type: "rule_proposal_pending", Severity: "warning", Title: "SOP 补充提案待确认", Message: "高频未覆盖场景已生成待确认 SOP 提案：" + req.ScenarioTitle, SourceType: "rule_proposal", SourceID: proposalID, CreatedAt: now}); err != nil {
			return err
		}
		return repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: auditID, RequestID: requestID, Actor: string(model.AuditActorSystem), Action: string(model.AuditActionCreateProposal), Status: string(model.AuditStatusSuccess), ProposalID: proposalID, InputRefType: "uncovered_scenario", InputRef: req.ScenarioKey, OutputRefType: "rule_proposal", OutputRef: proposalID, CreatedAt: now})
	}); err != nil {
		return dto.SOPAddendumProposalResponse{}, err
	}
	return dto.SOPAddendumProposalResponse{ProposalID: proposalID, Status: string(model.ProposalPendingUserConfirm), NotificationID: notificationID, AuditEventIDs: []string{auditID}, SafetyStatement: "SOP 补充只生成待确认提案，不自动应用规则、不连接券商、不自动交易。"}, nil
}

// ConfirmProposal 处理用户确认或拒绝。
func (s *RuleProposalService) ConfirmProposal(ctx context.Context, requestID, proposalID string, req dto.RuleProposalConfirmRequest, final bool) (dto.RuleProposalConfirmResponse, error) {
	proposal, err := s.loadProposal(ctx, proposalID)
	if err != nil {
		return dto.RuleProposalConfirmResponse{}, err
	}
	if proposal.Status == string(model.ProposalRejected) || proposal.Status == string(model.ProposalApplied) {
		return dto.RuleProposalConfirmResponse{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "终态规则提案不能重复确认")
	}
	if final {
		return s.finalConfirm(ctx, requestID, proposal, proposalID, req)
	}
	if proposal.Status != string(model.ProposalPendingUserConfirm) {
		return dto.RuleProposalConfirmResponse{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "规则提案未处于用户确认状态")
	}
	if req.Confirm != nil && !*req.Confirm {
		return s.reject(ctx, requestID, proposal, proposalID, string(model.AuditActionAuditRuleChange))
	}
	if proposal.SampleCount < 3 {
		return dto.RuleProposalConfirmResponse{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "样本数不足，不能送审")
	}
	return s.submitForGatekeeper(ctx, requestID, proposal, proposalID)
}

type loadedProposal struct {
	Status          string
	SampleCount     int
	ProposalVersion string
	AfterRuleJSON   string
}

func (s *RuleProposalService) loadProposal(ctx context.Context, proposalID string) (loadedProposal, error) {
	var out loadedProposal
	if err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		proposal, err := repos.RuleRepo.GetRuleProposal(ctx, proposalID)
		if err != nil {
			return err
		}
		out = loadedProposal{Status: proposal.Status, SampleCount: proposal.SampleCount, ProposalVersion: proposal.ProposalVersion, AfterRuleJSON: proposal.AfterRuleJSON}
		return nil
	}); err != nil {
		return out, err
	}
	return out, nil
}

func (s *RuleProposalService) reject(ctx context.Context, requestID string, proposal loadedProposal, proposalID string, action string) (dto.RuleProposalConfirmResponse, error) {
	auditID := s.ids.New("audit")
	if err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		if err := repos.RuleRepo.UpdateRuleProposalStatus(ctx, proposalID, string(model.ProposalRejected)); err != nil {
			return err
		}
		return repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: auditID, RequestID: requestID, Actor: string(model.AuditActorUser), Action: action, Status: string(model.AuditStatusSuccess), BeforeState: proposal.Status, AfterState: string(model.ProposalRejected), CreatedAt: s.clk.NowRFC3339()})
	}); err != nil {
		return dto.RuleProposalConfirmResponse{}, err
	}
	return dto.RuleProposalConfirmResponse{ProposalID: proposalID, Status: string(model.ProposalRejected), AuditEvents: []string{auditID}, AuditEventIDs: []string{auditID}}, nil
}

func (s *RuleProposalService) submitForGatekeeper(ctx context.Context, requestID string, proposal loadedProposal, proposalID string) (dto.RuleProposalConfirmResponse, error) {
	now := s.clk.NowRFC3339()
	userAuditID := s.ids.New("audit")
	underAudit, _, err := rule.AdvanceProposal(model.RuleProposalStatus(proposal.Status), proposal.SampleCount, true, "")
	if err != nil {
		return dto.RuleProposalConfirmResponse{}, mapRuleTransitionError(err)
	}
	if err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		if err := repos.RuleRepo.UpdateRuleProposalStatus(ctx, proposalID, string(underAudit)); err != nil {
			return err
		}
		return repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: userAuditID, RequestID: requestID, Actor: string(model.AuditActorUser), Action: string(model.AuditActionAuditRuleChange), Status: string(model.AuditStatusSuccess), BeforeState: proposal.Status, AfterState: string(underAudit), ProposalID: proposalID, CreatedAt: now})
	}); err != nil {
		return dto.RuleProposalConfirmResponse{}, err
	}
	auditOut, err := workflow.NewGatekeeperAuditGraphWithDependencies(s.deps).Run(ctx, workflow.GatekeeperAuditInput{RequestID: requestID, ProposalID: proposalID, Approved: true})
	if err != nil {
		return dto.RuleProposalConfirmResponse{}, mapRuleTransitionError(err)
	}
	gatekeeperAuditID := ""
	if len(auditOut.GatekeeperAudits) > 0 {
		gatekeeperAuditID = auditOut.GatekeeperAudits[0]
	}
	return dto.RuleProposalConfirmResponse{ProposalID: proposalID, Status: string(auditOut.ProposalStatus), GatekeeperAuditID: gatekeeperAuditID, AuditEvents: []string{userAuditID}, AuditEventIDs: []string{userAuditID}}, nil
}

func mapRuleTransitionError(err error) error {
	if errors.Is(err, rule.ErrInsufficientSamples) {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "样本数不足，不能送审")
	}
	if errors.Is(err, rule.ErrTerminalProposal) || errors.Is(err, rule.ErrInvalidTransition) {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "规则提案状态不允许当前操作")
	}
	return err
}

func (s *RuleProposalService) approve(ctx context.Context, requestID string, proposal loadedProposal, proposalID string, note string) (dto.RuleProposalConfirmResponse, error) {
	version := "v_" + proposalID
	now := s.clk.NowRFC3339()
	auditID := s.ids.New("audit")
	if err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		proposal, err := repos.RuleRepo.GetRuleProposal(ctx, proposalID)
		if err != nil {
			return err
		}
		if proposal.Status != string(model.ProposalPendingFinalConfirm) || proposal.SampleCount < 3 {
			return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "规则提案未满足最终确认条件")
		}
		gatekeeperAudit, err := repos.RuleRepo.GetLatestGatekeeperAuditByProposal(ctx, proposalID)
		if err != nil {
			return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "缺少允许应用的守门人审计")
		}
		if gatekeeperAudit.AuditResult != string(model.AuditApproved) || !gatekeeperAudit.AllowApply {
			return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "缺少允许应用的守门人审计")
		}
		validations, err := repos.RuleEffectRepo.ListRuleEffectValidations(ctx, repository.RuleEffectValidationFilter{ProposalID: proposalID, CandidateRuleVersion: proposal.ProposalVersion})
		if err != nil {
			return err
		}
		if len(validations) == 0 || validations[0].ValidationStatus != model.RuleEffectValidationPassed || validations[0].GuardrailDecision != model.RuleEffectGuardrailPassed || validations[0].OverfitRisk == model.RuleEffectOverfitHigh || validations[0].ReplayResult == model.RuleEffectReplayFailed {
			return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "规则效果验证未通过，不能最终应用")
		}
		if err := repos.RuleRepo.ArchiveActiveRuleVersions(ctx); err != nil {
			return err
		}
		if err := repos.RuleRepo.SaveRuleVersion(ctx, repository.RuleVersion{RuleVersion: version, Status: "active", RulesJSON: proposal.AfterRuleJSON, EffectiveAt: now, CreatedFromProposalID: proposalID, CreatedAt: now}); err != nil {
			return err
		}
		if err := repos.RuleRepo.ApplyRuleProposal(ctx, proposalID, string(model.ProposalApplied), now, note, version); err != nil {
			return err
		}
		if err := repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: auditID, RequestID: requestID, Actor: string(model.AuditActorUser), Action: string(model.AuditActionUpdateRule), Status: string(model.AuditStatusSuccess), BeforeState: proposal.Status, AfterState: string(model.ProposalApplied), ProposalID: proposalID, OutputRefType: "rule_version", OutputRef: version, CreatedAt: now}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return dto.RuleProposalConfirmResponse{}, err
	}
	return dto.RuleProposalConfirmResponse{ProposalID: proposalID, Status: string(model.ProposalApplied), AppliedRuleVersion: version, CreatedRuleVersion: version, FinalConfirmedAt: now, AuditEvents: []string{auditID}, AuditEventIDs: []string{auditID}}, nil
}

func (s *RuleProposalService) finalConfirm(ctx context.Context, requestID string, proposal loadedProposal, proposalID string, req dto.RuleProposalConfirmRequest) (dto.RuleProposalConfirmResponse, error) {
	if proposal.SampleCount < 3 || proposal.Status != string(model.ProposalPendingFinalConfirm) {
		return dto.RuleProposalConfirmResponse{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "规则提案未满足最终确认条件")
	}
	if req.Confirm != nil && !*req.Confirm {
		return s.reject(ctx, requestID, proposal, proposalID, string(model.AuditActionUpdateRule))
	}
	return s.approve(ctx, requestID, proposal, proposalID, req.Note)
}
