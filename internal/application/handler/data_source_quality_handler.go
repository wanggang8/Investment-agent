package handler

import (
	"net/http"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/application/service"
)

// GetDataSourceQualityRegression returns a local, read-only data-source quality regression summary.
func (a *App) GetDataSourceQualityRegression(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	out, err := a.DataSourceQualitySvc.Run(r.Context(), service.DataSourceQualityRegressionRequest{
		Mode:   r.URL.Query().Get("mode"),
		Symbol: r.URL.Query().Get("symbol"),
	})
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, out)
}

// GetDataQualityGateResolution returns the local release-claim resolution state for current data.
func (a *App) GetDataQualityGateResolution(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	out, err := a.DataSourceQualitySvc.CheckGateResolution(r.Context(), service.DataQualityGateResolutionCheckRequest{
		Symbol: r.URL.Query().Get("symbol"),
	})
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, out)
}

// ListDataQualityGateResolutions returns sanitized local resolution records.
func (a *App) ListDataQualityGateResolutions(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	items, err := a.DataSourceQualitySvc.ListGateResolutions(r.Context(), service.DataQualityGateResolutionListRequest{
		Symbol: r.URL.Query().Get("symbol"),
		Status: r.URL.Query().Get("status"),
	})
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, dto.DataQualityGateResolutionListResponse{Items: items, Total: len(items)})
}

// CreateDataQualityGateResolution records a manual waiver or scope exclusion for the current policy fingerprint.
func (a *App) CreateDataQualityGateResolution(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	var body dto.DataQualityGateResolutionCreateRequest
	if err := decodeJSON(r, &body); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	out, err := a.DataSourceQualitySvc.CreateGateResolution(r.Context(), service.DataQualityGateResolutionCreateRequest{
		RequestID:      requestID,
		Symbol:         body.Symbol,
		ResolutionType: body.ResolutionType,
		Scope:          body.Scope,
		Reason:         body.Reason,
		ReleaseImpact:  body.ReleaseImpact,
		EvidenceRef:    body.EvidenceRef,
	})
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, out)
}

// RetireDataQualityGateResolution retires an active resolution without changing the source-health facts.
func (a *App) RetireDataQualityGateResolution(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	out, err := a.DataSourceQualitySvc.RetireGateResolution(r.Context(), r.PathValue("resolution_id"), requestID)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, out)
}
