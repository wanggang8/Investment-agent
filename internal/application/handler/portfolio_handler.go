package handler

import (
	"net/http"

	"investment-agent/internal/application/dto"
)

// InitPortfolio 首次写入本地账户事实。该接口只更新本地事实库，不连接券商账户。
func (a *App) InitPortfolio(w http.ResponseWriter, r *http.Request) {
	a.writePortfolioSnapshot(w, r, "manual")
}

// AdjustPortfolio 手动校准本地账户状态，不写 position_transactions。
func (a *App) AdjustPortfolio(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	var req dto.PortfolioAdjustmentRequest
	if err := decodeJSON(r, &req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	resp, err := a.PortfolioSvc.WriteAdjustment(r.Context(), requestID, req, "manual")
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, resp)
}

func (a *App) EditHolding(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	var req dto.HoldingEditRequest
	if err := decodeJSON(r, &req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	resp, err := a.PortfolioSvc.EditHolding(r.Context(), requestID, req)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, resp)
}

func (a *App) RemoveHolding(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	var req dto.HoldingRemoveRequest
	if err := decodeJSON(r, &req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	resp, err := a.PortfolioSvc.RemoveHolding(r.Context(), requestID, req)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, resp)
}

func (a *App) RecordOfflineTransaction(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	var req dto.OfflineTransactionRequest
	if err := decodeJSON(r, &req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	resp, err := a.PortfolioSvc.RecordOfflineTransaction(r.Context(), requestID, req)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, resp)
}

func (a *App) ValidatePortfolioImport(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	var req dto.BatchImportValidationRequest
	if err := decodeJSON(r, &req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	resp, err := a.PortfolioSvc.ValidateImport(r.Context(), requestID, req)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, resp)
}

func (a *App) ConfirmPortfolioImport(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	var req dto.BatchImportConfirmRequest
	if err := decodeJSON(r, &req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	resp, err := a.PortfolioSvc.ConfirmImport(r.Context(), requestID, req)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, resp)
}

func (a *App) CorrectPortfolioFact(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	var req dto.CorrectionRequest
	if err := decodeJSON(r, &req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	resp, err := a.PortfolioSvc.CorrectFact(r.Context(), requestID, req)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, resp)
}

func (a *App) ReviewQuarterlyRebalance(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	var req dto.RebalanceReviewRequest
	if err := decodeJSON(r, &req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	resp, err := a.PortfolioSvc.ReviewQuarterlyRebalance(r.Context(), requestID, req)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, resp)
}

func (a *App) writePortfolioSnapshot(w http.ResponseWriter, r *http.Request, source string) {
	requestID := RequestID(r)
	var req dto.PortfolioInitRequest
	if err := decodeJSON(r, &req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	resp, err := a.PortfolioSvc.WriteSnapshot(r.Context(), requestID, req, source)
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, resp)
}

// GetPortfolioCurrent 返回最新账户快照和当前持仓聚合态。
func (a *App) GetPortfolioCurrent(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	s, err := a.QuerySvc.LatestPortfolioSnapshot(r.Context())
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	positions, err := a.QuerySvc.ListPositions(r.Context())
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	items := make([]dto.PositionDTO, 0, len(positions))
	for _, p := range positions {
		items = append(items, positionDTO(p))
	}
	writeOK(w, requestID, dto.PortfolioCurrentResponse{Snapshot: portfolioDTO(s), Positions: items})
}
