package repository

import "context"

// RuleVersion 保存已生效或已归档的正式规则快照。
type RuleVersion struct {
	RuleVersion           string
	Status                string
	RulesJSON             string
	EffectiveAt           string
	CreatedFromProposalID string
	CreatedAt             string
}

// RuleProposal 保存规则演进提案，应用前必须经过用户确认和守门人审计。
type RuleProposal struct {
	ProposalID            string
	ProposalType          string
	Status                string
	SourceErrorCaseID     string
	Title                 string
	ProposalVersion       string
	BeforeRuleJSON        string
	AfterRuleJSON         string
	Reason                string
	ImpactScopeJSON       string
	RiskNotesJSON         string
	SampleCount           int
	FinalConfirmedAt      string
	FinalConfirmedNote    string
	AppliedRuleVersion    string
	RelatedErrorCasesJSON string
	CreatedAt             string
}

// GatekeeperAudit 保存守门人对规则提案的审计结论。
type GatekeeperAudit struct {
	GatekeeperAuditID       string
	ProposalID              string
	AuditResult             string
	AuditReason             string
	RequiredChanges         string
	ViolatesFundamentalRule bool
	HasRuleConflict         bool
	BacktestMetricsJSON     string
	AllowApply              bool
	AuditedRuleVersion      string
	CreatedAt               string
}

// RuleProposalWithAudit combines a proposal with the latest gatekeeper audit summary.
type RuleProposalWithAudit struct {
	RuleProposal
	AuditResult string
	AuditReason string
}

// RuleRepository 定义规则版本、规则提案和守门人审计的持久化边界。
type RuleRepository interface {
	SaveRuleVersion(ctx context.Context, version RuleVersion) error
	GetRuleVersion(ctx context.Context, ruleVersion string) (RuleVersion, error)
	GetActiveRuleVersion(ctx context.Context) (RuleVersion, error)
	SaveRuleProposal(ctx context.Context, proposal RuleProposal) error
	GetRuleProposal(ctx context.Context, proposalID string) (RuleProposal, error)
	ListRuleProposals(ctx context.Context) ([]RuleProposalWithAudit, error)
	UpdateRuleProposalStatus(ctx context.Context, proposalID string, status string) error
	ApplyRuleProposal(ctx context.Context, proposalID, status, finalConfirmedAt, finalConfirmedNote, appliedRuleVersion string) error
	ArchiveActiveRuleVersions(ctx context.Context) error
	SaveGatekeeperAudit(ctx context.Context, audit GatekeeperAudit) error
	GetGatekeeperAudit(ctx context.Context, auditID string) (GatekeeperAudit, error)
	GetLatestGatekeeperAuditByProposal(ctx context.Context, proposalID string) (GatekeeperAudit, error)
}
