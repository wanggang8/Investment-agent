package dto

import "investment-agent/internal/domain/model"

// ConsultDecisionRequest 是主动咨询请求体。
// scenario 用于前端表达咨询场景，最终裁决仍由规则工作流生成。
type ConsultDecisionRequest struct {
	Question                           string                `json:"question"`
	Symbol                             string                `json:"symbol"`
	Scenario                           model.ConsultScenario `json:"scenario,omitempty"`
	ExpectedReturnPreviousBaseMidpoint float64               `json:"expected_return_previous_base_midpoint,omitempty"`
	ExpectedReturnTargetReturnRate     float64               `json:"expected_return_target_return_rate,omitempty"`
}

// DecisionDetailResponse 是决策详情页使用的 data 字段。
type DecisionDetailResponse struct {
	DecisionID              string                   `json:"decision_id"`
	Question                string                   `json:"question,omitempty"`
	Symbol                  string                   `json:"symbol,omitempty"`
	GeneratedAt             string                   `json:"generated_at,omitempty"`
	CapabilityCheck         *CapabilityCheck         `json:"capability_check,omitempty"`
	WorkflowStatus          string                   `json:"workflow_status"`
	AccountSnapshot         *AccountSnapshot         `json:"account_snapshot,omitempty"`
	TriggeredRules          []TriggeredRuleDTO       `json:"triggered_rules"`
	EvidenceChain           []EvidenceDTO            `json:"evidence_chain"`
	AnalystReports          []AnalystReport          `json:"analyst_reports"`
	RetrievalQuality        *RetrievalQualitySummary `json:"retrieval_quality,omitempty"`
	MarketContext           *MarketContext           `json:"market_context,omitempty"`
	ExpectedReturnScenarios *ExpectedReturnScenarios `json:"expected_return_scenarios,omitempty"`
	ArbitrationChain        []ArbitrationStep        `json:"arbitration_chain"`
	FinalVerdict            FinalVerdict             `json:"final_verdict"`
	UserConfirmation        UserConfirmation         `json:"user_confirmation"`
}

type CapabilityCheck struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type AccountSnapshot struct {
	SnapshotID    string  `json:"snapshot_id"`
	Cash          float64 `json:"cash,omitempty"`
	TotalAssets   float64 `json:"total_assets,omitempty"`
	CashRatio     float64 `json:"cash_ratio"`
	HighRiskRatio float64 `json:"high_risk_ratio"`
}

type AnalystReport struct {
	AgentName     string   `json:"agent_name"`
	Conclusion    string   `json:"conclusion"`
	KeyReasons    []string `json:"key_reasons"`
	RiskWarnings  []string `json:"risk_warnings"`
	Confidence    string   `json:"confidence"`
	EvidenceIDs   []string `json:"evidence_ids"`
	PromptVersion string   `json:"prompt_version,omitempty"`
	Model         string   `json:"model,omitempty"`
	InputSummary  string   `json:"input_summary,omitempty"`
	OutputSummary string   `json:"output_summary,omitempty"`
	ParseStatus   string   `json:"parse_status,omitempty"`
	QualityStatus string   `json:"quality_status,omitempty"`
}

type RetrievalQualitySummary struct {
	QuerySummary            string `json:"query_summary,omitempty"`
	TopK                    int    `json:"top_k"`
	Status                  string `json:"status,omitempty"`
	IndexHealth             string `json:"index_health,omitempty"`
	IndexFreshness          string `json:"index_freshness,omitempty"`
	FallbackSource          string `json:"fallback_source,omitempty"`
	SourceConsistencyStatus string `json:"source_consistency_status,omitempty"`
	DegradedReason          string `json:"degraded_reason,omitempty"`
}

type MarketContext struct {
	Symbol       string  `json:"symbol,omitempty"`
	TradeDate    string  `json:"trade_date,omitempty"`
	CurrentPrice float64 `json:"current_price,omitempty"`
	PEPercentile float64 `json:"pe_percentile,omitempty"`
	PBPercentile float64 `json:"pb_percentile,omitempty"`
}

type ExpectedReturnScenarios struct {
	SampleCount           int                    `json:"sample_count"`
	TargetName            string                 `json:"target_name,omitempty"`
	TargetCode            string                 `json:"target_code,omitempty"`
	HoldingClass          string                 `json:"holding_class,omitempty"`
	HorizonLabel          string                 `json:"horizon_label,omitempty"`
	SampleWindow          string                 `json:"sample_window,omitempty"`
	ScreeningCondition    string                 `json:"screening_condition,omitempty"`
	PrecisionStatus       string                 `json:"precision_status"`
	ProbabilityBasis      string                 `json:"probability_basis,omitempty"`
	Scenarios             []ReturnScenario       `json:"scenarios"`
	Reason                string                 `json:"reason,omitempty"`
	SupportingDataSummary string                 `json:"supporting_data_summary,omitempty"`
	MissingCategories     []string               `json:"missing_categories,omitempty"`
	SupplementData        []string               `json:"supplement_data,omitempty"`
	AssumptionChecks      []AssumptionCheck      `json:"assumption_checks,omitempty"`
	HistoricalContexts    []HistoricalContext    `json:"historical_contexts,omitempty"`
	HoldingClassCoverage  []HoldingClassCoverage `json:"holding_class_coverage,omitempty"`
	Disclaimer            string                 `json:"disclaimer"`
	SellEvaluation        *SellEvaluation        `json:"sell_evaluation,omitempty"`
	ReassessmentTrigger   *ReassessmentTrigger   `json:"reassessment_trigger,omitempty"`
}

type AssumptionCheck struct {
	Name        string  `json:"name"`
	Expected    float64 `json:"expected"`
	Actual      float64 `json:"actual"`
	MonthsBelow int     `json:"months_below"`
}

type HoldingClassCoverage struct {
	HoldingClass string `json:"holding_class"`
	Symbol       string `json:"symbol"`
	Status       string `json:"status"`
}

type HistoricalContext struct {
	Label       string  `json:"label"`
	Window      string  `json:"window"`
	SampleCount int     `json:"sample_count"`
	Outcome     string  `json:"outcome"`
	MaxDrawdown  float64 `json:"max_drawdown"`
	Recovery    string  `json:"recovery"`
	Source      string  `json:"source"`
}

type SellEvaluation struct {
	Status               string   `json:"status"`
	Triggers             []string `json:"triggers,omitempty"`
	Prompts              []string `json:"prompts,omitempty"`
	Actions              []string `json:"actions,omitempty"`
	NonTradingDisclaimer string   `json:"non_trading_disclaimer,omitempty"`
}

type ReassessmentTrigger struct {
	Reason       string  `json:"reason"`
	Boundary     string  `json:"boundary,omitempty"`
	CurrentValue float64 `json:"current_value,omitempty"`
}

type ReturnScenario struct {
	Scenario    string   `json:"scenario"`
	ReturnRange string   `json:"return_range"`
	Probability *float64 `json:"probability,omitempty"`
	Trigger     string   `json:"trigger,omitempty"`
}

type ArbitrationStep struct {
	Priority int    `json:"priority"`
	RuleID   string `json:"rule_id"`
	Result   string `json:"result"`
}

type FinalVerdict struct {
	Status            string   `json:"status"`
	DisplayText       string   `json:"display_text"`
	ProhibitedActions []string `json:"prohibited_actions"`
	OptionalActions   []string `json:"optional_actions"`
}

type UserConfirmation struct {
	ConfirmationStatus string   `json:"confirmation_status"`
	AvailableActions   []string `json:"available_actions"`
}

type DecisionListItem struct {
	DecisionID         string   `json:"decision_id"`
	DisplayTitle       string   `json:"display_title"`
	Symbol             string   `json:"symbol"`
	FinalVerdict       string   `json:"final_verdict"`
	TriggeredRuleIDs   []string `json:"triggered_rule_ids"`
	ConfirmationStatus string   `json:"confirmation_status"`
	GeneratedAt        string   `json:"generated_at"`
}

type ConfirmationRequest struct {
	ConfirmationType string  `json:"confirmation_type"`
	OperationType    string  `json:"operation_type,omitempty"`
	Symbol           string  `json:"symbol,omitempty"`
	Quantity         float64 `json:"quantity,omitempty"`
	Price            float64 `json:"price,omitempty"`
	Fees             float64 `json:"fees,omitempty"`
	ExecutedAt       string  `json:"executed_at,omitempty"`
	ActualOutcome    string  `json:"actual_outcome,omitempty"`
	RootCauseTag     string  `json:"root_cause_tag,omitempty"`
	LessonLearned    string  `json:"lesson_learned,omitempty"`
	Note             string  `json:"note,omitempty"`
}

type ConfirmationResponse struct {
	ConfirmationID     string   `json:"confirmation_id"`
	DecisionID         string   `json:"decision_id"`
	ConfirmationStatus string   `json:"confirmation_status"`
	ErrorCaseID        string   `json:"error_case_id"`
	TransactionIDs     []string `json:"transaction_ids"`
	SnapshotID         string   `json:"snapshot_id"`
	AuditEventIDs      []string `json:"audit_event_ids"`
}
