import type {
  BatchImportConfirmRequest,
  BatchImportValidationRequest,
  BatchImportValidationResponse,
  CorrectionRequest,
  HoldingEditRequest,
  HoldingRemoveRequest,
  LocalFactWriteResponse,
  OfflineTransactionRequest,
  PortfolioAdjustmentRequest,
  PortfolioCurrentResponse,
  PortfolioInitRequest,
  RebalanceReviewRequest,
  RebalanceReviewResponse,
  PortfolioWriteResponse,
} from '../types/portfolio'
import { apiRequest } from './client'

export function initPortfolio(body: PortfolioInitRequest) {
  return apiRequest<PortfolioWriteResponse>('/api/v1/portfolio/init', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export function getPortfolioCurrent() {
  return apiRequest<PortfolioCurrentResponse>('/api/v1/portfolio/current')
}

export function adjustPortfolio(body: PortfolioAdjustmentRequest) {
  return apiRequest<PortfolioWriteResponse>('/api/v1/portfolio/adjustments', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export function editHolding(body: HoldingEditRequest) {
  return apiRequest<LocalFactWriteResponse>('/api/v1/portfolio/holdings', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export function removeHolding(body: HoldingRemoveRequest) {
  return apiRequest<LocalFactWriteResponse>('/api/v1/portfolio/holdings/remove', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export function recordOfflineTransaction(body: OfflineTransactionRequest) {
  return apiRequest<LocalFactWriteResponse>('/api/v1/portfolio/offline-transactions', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export function validatePortfolioImport(body: BatchImportValidationRequest) {
  return apiRequest<BatchImportValidationResponse>('/api/v1/portfolio/imports/validate', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export function confirmPortfolioImport(body: BatchImportConfirmRequest) {
  return apiRequest<LocalFactWriteResponse>('/api/v1/portfolio/imports/confirm', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export function correctPortfolioFact(body: CorrectionRequest) {
  return apiRequest<LocalFactWriteResponse>('/api/v1/portfolio/corrections', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export function reviewQuarterlyRebalance(body: RebalanceReviewRequest) {
  return apiRequest<RebalanceReviewResponse>('/api/v1/portfolio/rebalance-review', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}
