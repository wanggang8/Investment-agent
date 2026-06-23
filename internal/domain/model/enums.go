package model

// DashboardState 是前端驾驶舱主状态。
type DashboardState string

// WorkflowStatus 是一次工作流执行结果。
type WorkflowStatus string

// PositionState 是当前持仓的交易约束状态。
type PositionState string

// VerificationStatus 是多源验证或证据验证状态。
type VerificationStatus string

// ConfirmationStatus 是决策记录上的用户处理状态。
type ConfirmationStatus string

// ConfirmationType 是用户本次提交的确认动作类型。
type ConfirmationType string

// FinalVerdictStatus 是领域规则给出的最终裁决状态。
type FinalVerdictStatus string

// OperationType 是线下手工交易动作类型。
type OperationType string

// AuditResult 是守门人审计结论。
type AuditResult string

// RuleProposalStatus 是规则提案状态机状态。
type RuleProposalStatus string

// AuditActor 是审计事件操作者类型。
type AuditActor string

// AuditAction 是审计事件业务动作。
type AuditAction string

// AuditStatus 是审计事件执行状态。
type AuditStatus string

// LiquidityState 是市场流动性状态。
type LiquidityState string

// SentimentState 是市场或用户情绪状态。
type SentimentState string

// PrecisionStatus 是预期收益精度状态。
type PrecisionStatus string

// RootCauseTag 是错误案例归因标签。
type RootCauseTag string

// RiskType 是 P35 风险预警类型。
type RiskType string

// RiskSeverity 是风险预警严重程度。
type RiskSeverity string

// RiskSOPStatus 是风险预警 SOP 生命周期状态。
type RiskSOPStatus string

// RuleEffectValidationStatus 是规则效果验证状态。
type RuleEffectValidationStatus string

// RuleEffectOverfitRisk 是规则效果验证中的过拟合风险等级。
type RuleEffectOverfitRisk string

// RuleEffectReplayResult 是候选规则历史回放结果。
type RuleEffectReplayResult string

// RuleEffectGuardrailDecision 是效果验证对门禁的建议结论。
type RuleEffectGuardrailDecision string

// RuleEffectTrendDirection 是应用后追踪趋势。
type RuleEffectTrendDirection string

const (
	DashboardFirstUse         DashboardState = "first_use"
	DashboardNormal           DashboardState = "normal"
	DashboardInsufficientData DashboardState = "insufficient_data"
	DashboardFrozenWatch      DashboardState = "frozen_watch"
	DashboardHighRisk         DashboardState = "high_risk"

	WorkflowCompleted WorkflowStatus = "completed"
	WorkflowDegraded  WorkflowStatus = "degraded"
	WorkflowFailed    WorkflowStatus = "failed"

	PositionNormal      PositionState = "normal"
	PositionSellOnly    PositionState = "sell_only"
	PositionFrozenWatch PositionState = "frozen_watch"

	VerificationSatisfied      VerificationStatus = "satisfied"
	VerificationFailed         VerificationStatus = "failed"
	VerificationBackgroundOnly VerificationStatus = "background_only"

	ConfirmationNotRequired      ConfirmationStatus = "not_required"
	ConfirmationPending          ConfirmationStatus = "pending"
	ConfirmationPlanned          ConfirmationStatus = "planned"
	ConfirmationExecutedManually ConfirmationStatus = "executed_manually"
	ConfirmationWatch            ConfirmationStatus = "watch"
	ConfirmationMarkedError      ConfirmationStatus = "marked_error"

	ConfirmationTypePlanned          ConfirmationType = "planned"
	ConfirmationTypeExecutedManually ConfirmationType = "executed_manually"
	ConfirmationTypeWatch            ConfirmationType = "watch"
	ConfirmationTypeMarkedError      ConfirmationType = "marked_error"

	VerdictBuyAllowed       FinalVerdictStatus = "buy_allowed"
	VerdictHold             FinalVerdictStatus = "hold"
	VerdictReduce           FinalVerdictStatus = "reduce"
	VerdictSellOnly         FinalVerdictStatus = "sell_only"
	VerdictFrozenWatch      FinalVerdictStatus = "frozen_watch"
	VerdictRejected         FinalVerdictStatus = "rejected"
	VerdictInsufficientData FinalVerdictStatus = "insufficient_data"

	OperationBuy    OperationType = "buy"
	OperationSell   OperationType = "sell"
	OperationReduce OperationType = "reduce"

	AuditApproved        AuditResult = "approved"
	AuditRejected        AuditResult = "rejected"
	AuditNeedsUserReview AuditResult = "needs_user_review"

	ProposalDraft                RuleProposalStatus = "draft"
	ProposalPendingUserConfirm   RuleProposalStatus = "pending_user_confirm"
	ProposalUnderGatekeeperAudit RuleProposalStatus = "under_gatekeeper_audit"
	ProposalPendingFinalConfirm  RuleProposalStatus = "pending_final_confirm"
	ProposalRejected             RuleProposalStatus = "rejected"
	ProposalApplied              RuleProposalStatus = "applied"

	AuditActorSystem     AuditActor = "system"
	AuditActorUser       AuditActor = "user"
	AuditActorGatekeeper AuditActor = "gatekeeper"

	AuditActionGenerateDecision  AuditAction = "generate_decision"
	AuditActionConfirmOperation  AuditAction = "confirm_operation"
	AuditActionMarkError         AuditAction = "mark_error"
	AuditActionCreateProposal    AuditAction = "create_proposal"
	AuditActionAuditRuleChange   AuditAction = "audit_rule_change"
	AuditActionUpdateRule        AuditAction = "update_rule"
	AuditActionRefreshMarketData AuditAction = "refresh_market_data"
	AuditActionUpdateSettings    AuditAction = "update_settings"
	AuditActionUpdateCapability  AuditAction = "update_capability"
	AuditActionRebuildIndex      AuditAction = "rebuild_index"
	AuditActionRunLocalTask      AuditAction = "run_local_task"
	AuditActionRiskAlert         AuditAction = "risk_alert"

	AuditStatusSuccess  AuditStatus = "success"
	AuditStatusDegraded AuditStatus = "degraded"
	AuditStatusFailed   AuditStatus = "failed"

	LiquidityNormal  LiquidityState = "normal"
	LiquidityWarning LiquidityState = "warning"
	LiquidityDanger  LiquidityState = "danger"

	SentimentCold    SentimentState = "cold"
	SentimentNeutral SentimentState = "neutral"
	SentimentHot     SentimentState = "hot"
	SentimentExtreme SentimentState = "extreme"

	PrecisionAvailable    PrecisionStatus = "available"
	PrecisionInsufficient PrecisionStatus = "insufficient"
	PrecisionUnavailable  PrecisionStatus = "unavailable"

	RootCauseEvidenceMissed     RootCauseTag = "evidence_missed"
	RootCauseRuleThresholdIssue RootCauseTag = "rule_threshold_issue"
	RootCauseAnalystError       RootCauseTag = "analyst_error"
	RootCauseUserContextMissing RootCauseTag = "user_context_missing"
	RootCauseMarketException    RootCauseTag = "market_exception"

	RiskTypeValuationHigh        RiskType = "valuation_high"
	RiskTypeBuyThesisBroken      RiskType = "buy_thesis_broken"
	RiskTypeLiquidityDanger      RiskType = "liquidity_danger"
	RiskTypeSentimentExtreme     RiskType = "sentiment_extreme"
	RiskTypePositionLimitBreach  RiskType = "position_limit_breach"
	RiskTypeInsufficientEvidence RiskType = "insufficient_evidence"
	RiskTypeDataDegraded         RiskType = "data_degraded"

	RiskSeverityInfo     RiskSeverity = "info"
	RiskSeverityWarning  RiskSeverity = "warning"
	RiskSeverityCritical RiskSeverity = "critical"

	RiskSOPTriggered RiskSOPStatus = "triggered"
	RiskSOPActive    RiskSOPStatus = "active"
	RiskSOPObserving RiskSOPStatus = "observing"
	RiskSOPEscalated RiskSOPStatus = "escalated"
	RiskSOPResolved  RiskSOPStatus = "resolved"
	RiskSOPArchived  RiskSOPStatus = "archived"

	RuleEffectValidationNotEvaluated     RuleEffectValidationStatus = "not_evaluated"
	RuleEffectValidationInsufficient     RuleEffectValidationStatus = "insufficient"
	RuleEffectValidationPassed           RuleEffectValidationStatus = "passed"
	RuleEffectValidationFailed           RuleEffectValidationStatus = "failed"
	RuleEffectValidationNeedsMoreSamples RuleEffectValidationStatus = "needs_more_samples"
	RuleEffectValidationNeedsUserReview  RuleEffectValidationStatus = "needs_user_review"

	RuleEffectOverfitLow    RuleEffectOverfitRisk = "low"
	RuleEffectOverfitMedium RuleEffectOverfitRisk = "medium"
	RuleEffectOverfitHigh   RuleEffectOverfitRisk = "high"

	RuleEffectReplayPassed  RuleEffectReplayResult = "passed"
	RuleEffectReplayFailed  RuleEffectReplayResult = "failed"
	RuleEffectReplayMixed   RuleEffectReplayResult = "mixed"
	RuleEffectReplayUnknown RuleEffectReplayResult = "unknown"

	RuleEffectGuardrailPassed          RuleEffectGuardrailDecision = "passed"
	RuleEffectGuardrailRejected        RuleEffectGuardrailDecision = "rejected"
	RuleEffectGuardrailNeedsUserReview RuleEffectGuardrailDecision = "needs_user_review"

	RuleEffectTrendImproved RuleEffectTrendDirection = "improved"
	RuleEffectTrendFlat     RuleEffectTrendDirection = "flat"
	RuleEffectTrendWorsened RuleEffectTrendDirection = "worsened"
	RuleEffectTrendUnknown  RuleEffectTrendDirection = "unknown"
)

// valid 判断枚举值是否在契约允许范围内。
func valid[T ~string](v T, allowed ...T) bool {
	for _, item := range allowed {
		if v == item {
			return true
		}
	}
	return false
}

func (v DashboardState) Valid() bool {
	return valid(v, DashboardFirstUse, DashboardNormal, DashboardInsufficientData, DashboardFrozenWatch, DashboardHighRisk)
}
func (v WorkflowStatus) Valid() bool {
	return valid(v, WorkflowCompleted, WorkflowDegraded, WorkflowFailed)
}
func (v PositionState) Valid() bool {
	return valid(v, PositionNormal, PositionSellOnly, PositionFrozenWatch)
}
func (v VerificationStatus) Valid() bool {
	return valid(v, VerificationSatisfied, VerificationFailed, VerificationBackgroundOnly)
}
func (v ConfirmationStatus) Valid() bool {
	return valid(v, ConfirmationNotRequired, ConfirmationPending, ConfirmationPlanned, ConfirmationExecutedManually, ConfirmationWatch, ConfirmationMarkedError)
}
func (v ConfirmationType) Valid() bool {
	return valid(v, ConfirmationTypePlanned, ConfirmationTypeExecutedManually, ConfirmationTypeWatch, ConfirmationTypeMarkedError)
}
func (v FinalVerdictStatus) Valid() bool {
	return valid(v, VerdictBuyAllowed, VerdictHold, VerdictReduce, VerdictSellOnly, VerdictFrozenWatch, VerdictRejected, VerdictInsufficientData)
}
func (v OperationType) Valid() bool { return valid(v, OperationBuy, OperationSell, OperationReduce) }
func (v AuditResult) Valid() bool {
	return valid(v, AuditApproved, AuditRejected, AuditNeedsUserReview)
}
func (v RuleProposalStatus) Valid() bool {
	return valid(v, ProposalDraft, ProposalPendingUserConfirm, ProposalUnderGatekeeperAudit, ProposalPendingFinalConfirm, ProposalRejected, ProposalApplied)
}
func (v AuditActor) Valid() bool {
	return valid(v, AuditActorSystem, AuditActorUser, AuditActorGatekeeper)
}
func (v AuditAction) Valid() bool {
	return valid(v, AuditActionGenerateDecision, AuditActionConfirmOperation, AuditActionMarkError, AuditActionCreateProposal, AuditActionAuditRuleChange, AuditActionUpdateRule, AuditActionRefreshMarketData, AuditActionUpdateSettings, AuditActionUpdateCapability, AuditActionRebuildIndex, AuditActionRunLocalTask, AuditActionRiskAlert)
}
func (v AuditStatus) Valid() bool {
	return valid(v, AuditStatusSuccess, AuditStatusDegraded, AuditStatusFailed)
}
func (v LiquidityState) Valid() bool {
	return valid(v, LiquidityNormal, LiquidityWarning, LiquidityDanger)
}
func (v SentimentState) Valid() bool {
	return valid(v, SentimentCold, SentimentNeutral, SentimentHot, SentimentExtreme)
}
func (v PrecisionStatus) Valid() bool {
	return valid(v, PrecisionAvailable, PrecisionInsufficient, PrecisionUnavailable)
}
func (v RootCauseTag) Valid() bool {
	return valid(v, RootCauseEvidenceMissed, RootCauseRuleThresholdIssue, RootCauseAnalystError, RootCauseUserContextMissing, RootCauseMarketException)
}
func (v RiskType) Valid() bool {
	return valid(v, RiskTypeValuationHigh, RiskTypeBuyThesisBroken, RiskTypeLiquidityDanger, RiskTypeSentimentExtreme, RiskTypePositionLimitBreach, RiskTypeInsufficientEvidence, RiskTypeDataDegraded)
}
func (v RiskSeverity) Valid() bool {
	return valid(v, RiskSeverityInfo, RiskSeverityWarning, RiskSeverityCritical)
}
func (v RiskSOPStatus) Valid() bool {
	return valid(v, RiskSOPTriggered, RiskSOPActive, RiskSOPObserving, RiskSOPEscalated, RiskSOPResolved, RiskSOPArchived)
}
func (v RiskSOPStatus) IsTerminal() bool {
	return v == RiskSOPResolved || v == RiskSOPArchived
}
func (v RuleEffectValidationStatus) Valid() bool {
	return valid(v, RuleEffectValidationNotEvaluated, RuleEffectValidationInsufficient, RuleEffectValidationPassed, RuleEffectValidationFailed, RuleEffectValidationNeedsMoreSamples, RuleEffectValidationNeedsUserReview)
}
func (v RuleEffectOverfitRisk) Valid() bool {
	return valid(v, RuleEffectOverfitLow, RuleEffectOverfitMedium, RuleEffectOverfitHigh)
}
func (v RuleEffectReplayResult) Valid() bool {
	return valid(v, RuleEffectReplayPassed, RuleEffectReplayFailed, RuleEffectReplayMixed, RuleEffectReplayUnknown)
}
func (v RuleEffectGuardrailDecision) Valid() bool {
	return valid(v, RuleEffectGuardrailPassed, RuleEffectGuardrailRejected, RuleEffectGuardrailNeedsUserReview)
}
func (v RuleEffectTrendDirection) Valid() bool {
	return valid(v, RuleEffectTrendImproved, RuleEffectTrendFlat, RuleEffectTrendWorsened, RuleEffectTrendUnknown)
}
