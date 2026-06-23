package model

// WorkflowContext 是应用层和工作流共享的领域上下文。
// 领域层只读取和追加状态，不依赖 HTTP、SQLite 或 Eino 实现细节。
type WorkflowContext struct {
	RequestID                string
	WorkflowType             string
	UserQuestion             string
	Symbol                   string
	PortfolioSnapshot        PortfolioSnapshot
	PositionSnapshots        []Position
	MarketSnapshot           MarketSnapshot
	RuleVersion              string
	CapabilityStatus         string
	CapabilityReason         string
	SourceVerificationStatus VerificationStatus
	MediaHeatSummary         string
	UserEmotionTags          []string
	EvidenceSet              EvidenceSet
	ExpectedReturnScenarios  []ExpectedReturnScenario
	RuleVerdict              RuleVerdict
	UserActionRequired       bool
	AuditEvents              []AuditEvent
	Errors                   []string
}

// RuleVerdict 是规则引擎输出的最终裁决。
// DeepSeek 等分析材料不得覆盖这里的 Status。
type RuleVerdict struct {
	Status            FinalVerdictStatus
	Text              string
	ProhibitedActions []string
	OptionalActions   []string
	TriggeredRules    []TriggeredRule
	ExpectedReturns   []ExpectedReturnScenario
}

// TriggeredRule 记录本次裁决命中的规则，用于前端展示和审计。
type TriggeredRule struct {
	RuleID      string
	RuleName    string
	Severity    string
	Description string
}

// ExpectedReturnScenario 是预期收益情景，只作为分析材料。
type ExpectedReturnScenario struct {
	Name        string
	Probability float64
	ReturnRate  float64
	Confidence  string
}
