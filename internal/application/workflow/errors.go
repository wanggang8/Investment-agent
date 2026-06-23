package workflow

import (
	"investment-agent/internal/domain/model"
	"investment-agent/internal/pkg/apperr"
)

const (
	// ErrCodeDataRequired 表示账户或持仓等基础数据缺失。
	ErrCodeDataRequired = string(apperr.CodeDataRequired)
	// ErrCodeDataStale 表示行情或估值数据过期。
	ErrCodeDataStale = string(apperr.CodeDataStale)
	// ErrCodeRuleVersionMissing 表示缺少可用规则版本。
	ErrCodeRuleVersionMissing = string(apperr.CodeRuleVersionMissing)
	// ErrCodeEvidenceNotFound 表示没有可用正式证据。
	ErrCodeEvidenceNotFound = string(apperr.CodeEvidenceNotFound)
	// ErrCodeSourceVerificationFailed 表示多源验证未满足。
	ErrCodeSourceVerificationFailed = string(apperr.CodeSourceVerificationFailed)
	// ErrCodeVectorIndexUnavailable 表示 VecLite 索引不可用。
	ErrCodeVectorIndexUnavailable = string(apperr.CodeVectorIndexUnavailable)
	// ErrCodeAnalystUnavailable 表示 DeepSeek 分析节点不可用。
	ErrCodeAnalystUnavailable = string(apperr.CodeAnalystUnavailable)
	// ErrCodeDecisionRecordFailed 表示决策记录保存失败。
	ErrCodeDecisionRecordFailed = string(apperr.CodeDecisionRecordFailed)
)

// auditStatus 把节点状态转换为领域审计状态。
func auditStatus(status NodeStatus) model.AuditStatus {
	switch status {
	case StatusFailed:
		return model.AuditStatusFailed
	case StatusDegraded:
		return model.AuditStatusDegraded
	default:
		return model.AuditStatusSuccess
	}
}

// auditAction 把字符串动作转换为领域审计动作。
func auditAction(action string) model.AuditAction {
	v := model.AuditAction(action)
	if v.Valid() {
		return v
	}
	return model.AuditActionGenerateDecision
}
