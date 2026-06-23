package dto

// DataSourceQualityRegressionResponse describes a local source-quality regression run.
type DataSourceQualityRegressionResponse struct {
	Mode              string                  `json:"mode"`
	Status            string                  `json:"status"`
	GeneratedAt       string                  `json:"generated_at"`
	Summary           string                  `json:"summary"`
	Cases             []DataSourceQualityCase `json:"cases"`
	MissingCategories []string                `json:"missing_categories"`
	Policy            DataSourceQualityPolicy `json:"policy"`
	SafetyNote        string                  `json:"safety_note"`
}

// DataSourceQualityPolicy describes release impact for current source-health quality.
type DataSourceQualityPolicy struct {
	Verdict         string   `json:"verdict"`
	ReleaseGate     string   `json:"release_gate"`
	DegradedCount   int      `json:"degraded_count"`
	FailedCount     int      `json:"failed_count"`
	BlockingCount   int      `json:"blocking_count"`
	WaiverCount     int      `json:"waiver_count"`
	BlockingReasons []string `json:"blocking_reasons"`
	WaiverReasons   []string `json:"waiver_reasons"`
	NextActions     []string `json:"next_actions"`
	SafetyNote      string   `json:"safety_note"`
}

// DataQualityGateResolutionRecord describes a local manual resolution record for release claims.
type DataQualityGateResolutionRecord struct {
	ResolutionID      string `json:"resolution_id"`
	Symbol            string `json:"symbol"`
	PolicyFingerprint string `json:"policy_fingerprint"`
	PolicyVerdict     string `json:"policy_verdict"`
	ReleaseGate       string `json:"release_gate"`
	PolicySummary     string `json:"policy_summary"`
	ResolutionType    string `json:"resolution_type"`
	Status            string `json:"status"`
	Scope             string `json:"scope"`
	Reason            string `json:"reason"`
	ReleaseImpact     string `json:"release_impact"`
	EvidenceRef       string `json:"evidence_ref,omitempty"`
	CreatedBy         string `json:"created_by"`
	RetiredBy         string `json:"retired_by,omitempty"`
	CreatedAt         string `json:"created_at"`
	RetiredAt         string `json:"retired_at,omitempty"`
	SafetyNote        string `json:"safety_note"`
}

// DataQualityGateResolutionCheck combines current policy with local resolution state.
type DataQualityGateResolutionCheck struct {
	Symbol                string                           `json:"symbol"`
	PolicyFingerprint     string                           `json:"policy_fingerprint"`
	PolicySummary         string                           `json:"policy_summary"`
	Policy                DataSourceQualityPolicy          `json:"policy"`
	ReleaseClaimState     string                           `json:"release_claim_state"`
	CleanDataClaimAllowed bool                             `json:"clean_data_claim_allowed"`
	ActiveResolution      *DataQualityGateResolutionRecord `json:"active_resolution,omitempty"`
	AllowedClaims         []string                         `json:"allowed_claims"`
	ProhibitedClaims      []string                         `json:"prohibited_claims"`
	SafetyNote            string                           `json:"safety_note"`
}

// DataQualityGateResolutionCreateRequest is the API payload for a manual resolution.
type DataQualityGateResolutionCreateRequest struct {
	Symbol         string `json:"symbol"`
	ResolutionType string `json:"resolution_type"`
	Scope          string `json:"scope"`
	Reason         string `json:"reason"`
	ReleaseImpact  string `json:"release_impact"`
	EvidenceRef    string `json:"evidence_ref"`
}

// DataQualityGateResolutionListResponse lists local resolution records.
type DataQualityGateResolutionListResponse struct {
	Items []DataQualityGateResolutionRecord `json:"items"`
	Total int                               `json:"total"`
}

// DataSourceQualityCase describes one source-health regression assertion.
type DataSourceQualityCase struct {
	CaseID            string   `json:"case_id"`
	SourceName        string   `json:"source_name"`
	SourceLevel       string   `json:"source_level"`
	SourceType        string   `json:"source_type"`
	DataCategory      string   `json:"data_category"`
	ExpectedFreshness string   `json:"expected_freshness"`
	ActualFreshness   string   `json:"actual_freshness"`
	Status            string   `json:"status"`
	DataDate          string   `json:"data_date,omitempty"`
	FailureCategory   string   `json:"failure_category,omitempty"`
	AffectedSymbols   []string `json:"affected_symbols,omitempty"`
	DiagnosticPreview string   `json:"diagnostic_preview,omitempty"`
}
