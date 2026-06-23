package repository

import "context"

// DataQualityGateResolution records a local manual handling decision for a current-data release gate.
type DataQualityGateResolution struct {
	ResolutionID        string
	Symbol              string
	PolicyFingerprint   string
	PolicyVerdict       string
	ReleaseGate         string
	PolicySummary       string
	ResolutionType      string
	Status              string
	Scope               string
	Reason              string
	ReleaseImpact       string
	EvidenceRef         string
	BlockingReasonsJSON string
	WaiverReasonsJSON   string
	CreatedBy           string
	RetiredBy           string
	CreatedAt           string
	RetiredAt           string
	SafetyNote          string
}

// DataQualityGateResolutionFilter filters local gate resolution records.
type DataQualityGateResolutionFilter struct {
	Symbol string
	Status string
}

// DataQualityGateResolutionRepository persists local current-data gate resolution records.
type DataQualityGateResolutionRepository interface {
	CreateDataQualityGateResolution(ctx context.Context, resolution DataQualityGateResolution) error
	GetDataQualityGateResolution(ctx context.Context, resolutionID string) (DataQualityGateResolution, error)
	GetActiveDataQualityGateResolution(ctx context.Context, symbol, policyFingerprint string) (DataQualityGateResolution, error)
	ListDataQualityGateResolutions(ctx context.Context, filter DataQualityGateResolutionFilter) ([]DataQualityGateResolution, error)
	RetireDataQualityGateResolution(ctx context.Context, resolutionID, retiredBy, retiredAt string) error
}
