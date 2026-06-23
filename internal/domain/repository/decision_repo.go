package repository

import "context"

// DecisionRecord 保存一次正式建议或非交易记录的最终结果。
type DecisionRecord struct {
	DecisionID                  string
	RequestID                   string
	WorkflowType                string
	Symbol                      string
	Question                    string
	WorkflowStatus              string
	RecordType                  string
	DashboardState              string
	CapabilityStatus            string
	CapabilityReason            string
	SourceVerificationStatus    string
	RiskReasonCode              string
	MediaHeatSummaryJSON        string
	UserEmotionTagsJSON         string
	TriggeredRulesJSON          string
	ErrorsJSON                  string
	FinalVerdictStatus          string
	FinalVerdictText            string
	ProhibitedActionsJSON       string
	OptionalActionsJSON         string
	ConfirmationStatus          string
	PortfolioSnapshotID         string
	MarketSnapshotID            string
	RuleVersion                 string
	AnalystReportsJSON          string
	ExpectedReturnScenariosJSON string
	ArbitrationChainJSON        string
	ContextSnapshotJSON         string
	CreatedAt                   string
}

// EvidenceRef 保存决策引用过的证据快照，避免历史详情受后续情报清洗影响。
type EvidenceRef struct {
	EvidenceRefID                   string
	EvidenceID                      string
	DecisionID                      string
	SummaryID                       string
	SourceName                      string
	SourceLevel                     string
	EvidenceRole                    string
	PublishedAt                     string
	CapturedAt                      string
	OriginalURL                     string
	Summary                         string
	ContentHash                     string
	TimeWeight                      float64
	RelevanceScore                  float64
	IndependentSourceCount          int
	HighGradeIndependentSourceCount int
	CreatedAt                       string
}

// PositionTransaction records an executed manual operation.
type PositionTransaction struct {
	TransactionID      string
	ConfirmationID     string
	Symbol             string
	OperationType      string
	Quantity           float64
	Price              float64
	Fees               float64
	OccurredAt         string
	BeforePositionJSON string
	AfterPositionJSON  string
	CreatedAt          string
}

// ErrorCase records a user-marked decision error.
type ErrorCase struct {
	ErrorCaseID    string
	DecisionID     string
	ConfirmationID string
	ActualOutcome  string
	RootCauseTag   string
	LessonLearned  string
	CreatedAt      string
}

// OperationConfirmation 只记录用户线下处理结果，不触发自动交易。
type OperationConfirmation struct {
	ConfirmationID   string
	DecisionID       string
	ConfirmationType string
	OperationType    string
	Symbol           string
	Quantity         float64
	Price            float64
	Fees             float64
	ExecutedAt       string
	ErrorCaseID      string
	PayloadJSON      string
	Note             string
	CreatedAt        string
}

// DecisionRepository 定义决策记录、证据引用和用户确认的持久化边界。
type DecisionRepository interface {
	SaveDecisionRecord(ctx context.Context, record DecisionRecord, evidenceRefs []EvidenceRef) error
	GetDecisionRecord(ctx context.Context, decisionID string) (DecisionRecord, []EvidenceRef, error)
	ListDecisionRecords(ctx context.Context) ([]DecisionRecord, error)
	ListErrorCases(ctx context.Context) ([]ErrorCase, error)
	CountErrorCases(ctx context.Context) (int, error)
	GetDecisionConfirmationState(ctx context.Context, decisionID string) (recordType, confirmationStatus string, err error)
	ListOperationConfirmations(ctx context.Context, decisionID string) ([]OperationConfirmation, error)
	SaveOperationConfirmation(ctx context.Context, confirmation OperationConfirmation) error
	UpdateDecisionConfirmationStatus(ctx context.Context, decisionID, status string) error
	UpdateDecisionConfirmationStatusIfCurrent(ctx context.Context, decisionID, expectedStatus, nextStatus string) (bool, error)
	ListPositionTransactionsByConfirmation(ctx context.Context, confirmationID string) ([]PositionTransaction, error)
	SavePositionTransaction(ctx context.Context, transaction PositionTransaction) error
	SaveErrorCase(ctx context.Context, errorCase ErrorCase) error
	GetOperationConfirmation(ctx context.Context, confirmationID string) (OperationConfirmation, error)
}
