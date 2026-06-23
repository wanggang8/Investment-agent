package model

// RuleProposal 是领域层用于状态机判断的规则提案摘要。
type RuleProposal struct {
	ProposalID  string
	Status      RuleProposalStatus
	SampleCount int
	RuleVersion string
}

// ProposalTransitionInput 描述用户确认或守门人审计触发的提案流转输入。
type ProposalTransitionInput struct {
	Confirm     bool
	AuditResult AuditResult
}
