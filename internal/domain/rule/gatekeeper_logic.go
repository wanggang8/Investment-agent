package rule

import (
	"errors"

	"investment-agent/internal/domain/model"
)

// ErrInvalidTransition 表示状态机不支持当前流转。
var ErrInvalidTransition = errors.New("invalid rule proposal transition")

// ErrInsufficientSamples 表示样本不足，不能进入审计或最终应用路径。
var ErrInsufficientSamples = errors.New("insufficient samples for rule application")

// ErrTerminalProposal 表示提案已处于终态，不能再次操作。
var ErrTerminalProposal = errors.New("terminal proposal cannot transition")

// AdvanceProposal 推进规则提案状态机。
// 返回值中的 bool 表示是否允许创建新的 active rule_version。
func AdvanceProposal(status model.RuleProposalStatus, sampleCount int, confirm bool, audit model.AuditResult) (model.RuleProposalStatus, bool, error) {
	switch status {
	case model.ProposalDraft:
		if confirm {
			return model.ProposalPendingUserConfirm, false, nil
		}
		return model.ProposalRejected, false, nil
	case model.ProposalPendingUserConfirm:
		if !confirm {
			return model.ProposalRejected, false, nil
		}
		if sampleCount < 3 {
			return model.ProposalPendingUserConfirm, false, ErrInsufficientSamples
		}
		return model.ProposalUnderGatekeeperAudit, false, nil
	case model.ProposalUnderGatekeeperAudit:
		switch audit {
		case model.AuditApproved:
			return model.ProposalPendingFinalConfirm, false, nil
		case model.AuditRejected:
			return model.ProposalRejected, false, nil
		case model.AuditNeedsUserReview:
			return model.ProposalPendingUserConfirm, false, nil
		default:
			return status, false, ErrInvalidTransition
		}
	case model.ProposalPendingFinalConfirm:
		if !confirm {
			return model.ProposalRejected, false, nil
		}
		if sampleCount < 3 {
			return status, false, ErrInsufficientSamples
		}
		return model.ProposalApplied, true, nil
	case model.ProposalRejected, model.ProposalApplied:
		return status, false, ErrTerminalProposal
	default:
		return status, false, ErrInvalidTransition
	}
}
