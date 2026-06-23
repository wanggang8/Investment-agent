package dto

// DecisionLoopListResponse is the read-only list response for decision loop explanations.
type DecisionLoopListResponse struct {
	Items      []DecisionLoopItem `json:"items"`
	Total      int                `json:"total"`
	SafetyNote string             `json:"safety_note"`
}

// DecisionLoopItem links a decision to local confirmations, manual records, risk, review, and audit traces.
type DecisionLoopItem struct {
	DecisionID         string                     `json:"decision_id"`
	Symbol             string                     `json:"symbol,omitempty"`
	GeneratedAt        string                     `json:"generated_at"`
	FinalVerdictStatus string                     `json:"final_verdict_status"`
	FinalVerdictText   string                     `json:"final_verdict_text"`
	ConfirmationStatus string                     `json:"confirmation_status"`
	LoopStatus         string                     `json:"loop_status"`
	Stages             []DecisionLoopStage        `json:"stages"`
	ManualActions      []DecisionLoopManualAction `json:"manual_actions"`
	RiskLinks          []DecisionLoopLink         `json:"risk_links"`
	ReviewLinks        []DecisionLoopLink         `json:"review_links"`
	AuditLinks         []DecisionLoopLink         `json:"audit_links"`
	MissingLinks       []string                   `json:"missing_links"`
	SafetyNote         string                     `json:"safety_note"`
}

// DecisionLoopStage describes one expected step in the local decision handling loop.
type DecisionLoopStage struct {
	Stage   string `json:"stage"`
	Status  string `json:"status"`
	Label   string `json:"label"`
	Summary string `json:"summary"`
	RefType string `json:"ref_type,omitempty"`
	RefID   string `json:"ref_id,omitempty"`
	At      string `json:"at,omitempty"`
}

// DecisionLoopManualAction summarizes one user-recorded local operation without exposing raw payloads.
type DecisionLoopManualAction struct {
	ConfirmationID   string   `json:"confirmation_id"`
	ConfirmationType string   `json:"confirmation_type"`
	OperationType    string   `json:"operation_type,omitempty"`
	Symbol           string   `json:"symbol,omitempty"`
	Quantity         float64  `json:"quantity,omitempty"`
	Price            float64  `json:"price,omitempty"`
	Fees             float64  `json:"fees,omitempty"`
	ExecutedAt       string   `json:"executed_at,omitempty"`
	TransactionIDs   []string `json:"transaction_ids"`
	NotePreview      string   `json:"note_preview,omitempty"`
}

// DecisionLoopLink is a safe local navigation link to related read-only facts.
type DecisionLoopLink struct {
	Type   string `json:"type"`
	ID     string `json:"id"`
	Label  string `json:"label"`
	Href   string `json:"href"`
	Status string `json:"status,omitempty"`
}
