package rule

import (
	"strconv"

	"investment-agent/internal/domain/model"
)

// evaluateEvidence 处理缺少有效证据的情况。
func evaluateEvidence(in EvaluationInput) (model.RuleVerdict, bool) {
	if !in.HasEvidence {
		return model.RuleVerdict{Status: model.VerdictInsufficientData, Text: "证据不足，暂停交易类建议", ProhibitedActions: []string{"交易类建议"}, TriggeredRules: []model.TriggeredRule{{RuleID: "EVIDENCE", RuleName: "证据不足", Severity: "warning", Description: "缺少有效证据"}}}, true
	}
	return model.RuleVerdict{}, false
}

// evaluateSource 处理信源等级和重大事件多源验证规则。
func evaluateSource(in EvaluationInput) (model.RuleVerdict, bool) {
	if in.SourceVerificationStatus == model.VerificationFailed {
		return frozenWatch("多源验证未满足"), true
	}
	if in.SourceVerificationStatus == model.VerificationBackgroundOnly {
		return model.RuleVerdict{Status: model.VerdictInsufficientData, Text: "证据仅可作为背景材料，暂停交易类建议", ProhibitedActions: []string{"交易类建议"}, TriggeredRules: []model.TriggeredRule{{RuleID: "SOURCE", RuleName: "证据角色", Severity: "warning", Description: "背景材料不能作为正式裁决依据"}}}, true
	}
	for _, ev := range in.Evidence {
		if ev.Role == model.EvidenceFormal && !ev.SourceLevel.FormalAllowed() {
			return model.RuleVerdict{Status: model.VerdictInsufficientData, Text: "C 级信源只能作为背景材料", ProhibitedActions: []string{"使用 C 级信源作正式裁决"}}, true
		}
		if isMajorEvent(ev.EventType) && ev.HighGradeIndependentSourceCount < 2 {
			return frozenWatch("重大事件缺少 2 个 A/S 独立信源；当前 A/S 独立信源=" + strconv.Itoa(ev.HighGradeIndependentSourceCount)), true
		}
	}
	return model.RuleVerdict{}, false
}

func isMajorEvent(t model.EventType) bool {
	return t == model.EventMajorPositive || t == model.EventMajorNegative || t == model.EventBuyLogicBreak
}

func frozenWatch(text string) model.RuleVerdict {
	return model.RuleVerdict{Status: model.VerdictFrozenWatch, Text: text, ProhibitedActions: []string{"主动交易建议"}, TriggeredRules: []model.TriggeredRule{{RuleID: "SOURCE", RuleName: "多源验证", Severity: "warning", Description: text}}}
}
