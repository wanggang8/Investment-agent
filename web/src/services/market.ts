import type { MarketRefreshRequest, MarketRefreshResponse, MarketSnapshot, SourceHealthResponse } from '../types/market'
import { apiRequest } from './client'

export function refreshMarket(body: MarketRefreshRequest = {}) {
  return apiRequest<MarketRefreshResponse>('/api/v1/market/refresh', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export function getLatestMarketSnapshot(symbol?: string) {
  const query = symbol ? `?symbol=${encodeURIComponent(symbol)}` : ''
  return apiRequest<MarketSnapshot>(`/api/v1/market/snapshots/latest${query}`)
}

export function getMarketSourceHealth(symbol?: string) {
  const query = symbol ? `?symbol=${encodeURIComponent(symbol)}` : ''
  return apiRequest<SourceHealthResponse>(`/api/v1/market/source-health${query}`)
}
