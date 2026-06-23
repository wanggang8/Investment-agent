package workflow

import "investment-agent/internal/domain/model"

const (
	// WorkflowDailyDiscipline 表示每日纪律工作流。
	WorkflowDailyDiscipline = "daily_discipline"
	// WorkflowConsultation 表示用户主动咨询工作流。
	WorkflowConsultation = "consultation"
	// WorkflowEvidenceVerification 表示证据核查工作流。
	WorkflowEvidenceVerification = "evidence_verification"
	// WorkflowEvolutionProposal 表示规则提案生成工作流。
	WorkflowEvolutionProposal = "evolution_proposal"
	// WorkflowGatekeeperAudit 表示守门人审计工作流。
	WorkflowGatekeeperAudit = "gatekeeper_audit"
	// WorkflowMarketRefresh 表示市场刷新工作流；该流程只写 market_snapshots 与审计事件，不写决策记录。
	WorkflowMarketRefresh = "market_refresh"
)

const (
	// CapabilityInScope 表示标的在能力圈内。
	CapabilityInScope = "in_scope"
	// CapabilityOutOfScope 表示标的不在能力圈内。
	CapabilityOutOfScope = "out_of_scope"
	// CapabilityUnknown 表示能力圈未配置或无法判断。
	CapabilityUnknown = "unknown"
)

// WorkflowContext 是应用层工作流统一上下文。
// 它对齐 docs/workflow.md 第 4 节，用于节点间传递账户、证据、分析、裁决和审计信息。
type WorkflowContext struct {
	RequestID                           string
	WorkflowType                        string
	UserQuestion                        string
	Symbol                              string
	PortfolioSnapshot                   model.PortfolioSnapshot
	PositionSnapshots                   []model.Position
	MarketSnapshot                      model.MarketSnapshot
	RuleVersion                         string
	CapabilityStatus                    string
	CapabilityReason                    string
	SourceVerificationStatus            model.VerificationStatus
	MediaHeatSummary                    string
	UserEmotionTags                     []string
	EvidenceSet                         model.EvidenceSet
	AnalystReports                      map[string]string
	AnalystReportMetadata               map[string]map[string]string
	AnalystUnavailable                  bool
	RetrievalInput                      string
	RetrievalOutputRef                  string
	RetrievalDegradedReason             string
	RetrievalQualitySummary             RetrievalQualitySummary
	ExpectedReturnSampleCount           int
	ExpectedReturnPreviousBaseMidpoint  float64
	ExpectedReturnTargetReturnRate      float64
	ExpectedReturnTargetName            string
	ExpectedReturnTargetCode            string
	ExpectedReturnHoldingClass          string
	ExpectedReturnHorizonLabel          string
	ExpectedReturnSampleWindow          string
	ExpectedReturnScreeningCondition    string
	ExpectedReturnScenarios             []ExpectedReturnScenario
	ExpectedReturnSellEvaluation        ExpectedReturnSellEvaluation
	ExpectedReturnReassessmentTrigger   ExpectedReturnReassessmentTrigger
	ExpectedReturnPrecisionStatus       model.PrecisionStatus
	ExpectedReturnReason                string
	ExpectedReturnProbabilityBasis      string
	ExpectedReturnSupportingDataSummary string
	ExpectedReturnMissingCategories     []string
	ExpectedReturnSupplementData        []string
	ExpectedReturnAssumptionChecks      []ExpectedReturnAssumptionCheck
	ExpectedReturnHistoricalContexts    []ExpectedReturnHistoricalContext
	ExpectedReturnHoldingClassCoverage  []ExpectedReturnHoldingClassCoverage
	DecisionID                          string
	RuleVerdict                         model.RuleVerdict
	UserActionRequired                  bool
	AuditEvents                         []model.AuditEvent
	Errors                              []string
}

// ExpectedReturnScenario 是应用层预期收益 DTO。
// Probability 为指针，用于表达样本不足时不返回精确概率。
type ExpectedReturnScenario struct {
	Name        string
	Probability *float64
	ReturnRate  float64
	ReturnRange string
	LowerBound  float64
	UpperBound  float64
	Confidence  string
	Trigger     string
}

type ExpectedReturnHistoricalSample struct {
	Scenario    string  `json:"scenario"`
	Count       int     `json:"count"`
	ReturnRate  float64 `json:"return_rate"`
	ReturnRange string  `json:"return_range"`
	LowerBound  float64 `json:"lower_bound"`
	UpperBound  float64 `json:"upper_bound"`
	Trigger     string  `json:"trigger"`
}

type ExpectedReturnAssumptionCheck struct {
	Name        string  `json:"name"`
	Expected    float64 `json:"expected"`
	Actual      float64 `json:"actual"`
	MonthsBelow int     `json:"months_below"`
}

type ExpectedReturnHoldingClassCoverage struct {
	HoldingClass string `json:"holding_class"`
	Symbol       string `json:"symbol"`
	Status       string `json:"status"`
}

type ExpectedReturnSellEvaluation struct {
	Status               string   `json:"status"`
	Triggers             []string `json:"triggers"`
	Prompts              []string `json:"prompts"`
	Actions              []string `json:"actions"`
	NonTradingDisclaimer string   `json:"non_trading_disclaimer"`
}

type ExpectedReturnReassessmentTrigger struct {
	Reason       string  `json:"reason"`
	Boundary     string  `json:"boundary"`
	CurrentValue float64 `json:"current_value"`
}

type ExpectedReturnHistoricalContext struct {
	Label       string  `json:"label"`
	Window      string  `json:"window"`
	SampleCount int     `json:"sample_count"`
	Outcome     string  `json:"outcome"`
	MaxDrawdown  float64 `json:"max_drawdown"`
	Recovery    string  `json:"recovery"`
	Source      string  `json:"source"`
}

type ExpectedReturnInput struct {
	SampleCount           int
	TargetName            string
	TargetCode            string
	HoldingClass          string
	HorizonLabel          string
	CurrentPrice          float64
	BasePrice             float64
	PreviousBaseMidpoint  float64
	TargetReturnRate      float64
	SentimentState         string
	MarketState           string
	FundamentalState      string
	SampleWindow          string
	ScreeningCondition    string
	SupportingDataSummary string
	MissingCategories     []string
	HistoricalSamples     []ExpectedReturnHistoricalSample
	HistoricalContexts    []ExpectedReturnHistoricalContext
	AssumptionChecks      []ExpectedReturnAssumptionCheck
	PessimisticPathMonths int
	HoldingClassCoverage  []ExpectedReturnHoldingClassCoverage
}

// ExpectedReturnOutput 是 ExpectedReturnNode 的输出结构。
type ExpectedReturnOutput struct {
	PrecisionStatus       model.PrecisionStatus
	TargetName            string
	TargetCode            string
	HoldingClass          string
	HorizonLabel          string
	Scenarios             []ExpectedReturnScenario
	Reason                string
	SampleCount           int
	SampleWindow          string
	ScreeningCondition    string
	ProbabilityBasis      string
	SupportingDataSummary string
	MissingCategories     []string
	SupplementData        []string
	AssumptionChecks      []ExpectedReturnAssumptionCheck
	HistoricalContexts    []ExpectedReturnHistoricalContext
	HoldingClassCoverage  []ExpectedReturnHoldingClassCoverage
	SellEvaluation        ExpectedReturnSellEvaluation
	ReassessmentTrigger   ExpectedReturnReassessmentTrigger
}

func (o ExpectedReturnOutput) CoveredHoldingClasses() []string {
	out := []string{}
	for _, item := range o.HoldingClassCoverage {
		if item.Status == "covered" {
			out = append(out, item.HoldingClass)
		}
	}
	return out
}
