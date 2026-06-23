import type { DataQualityGateResolutionCheck, DataQualityGateResolutionCreateRequest, DataQualityGateResolutionListResponse, DataSourceQualityRegression } from '../types/dataSourceQuality'
import { apiRequest } from './client'

export function getDataSourceQualityRegression(mode = 'current', symbol?: string) {
  const params = new URLSearchParams({ mode })
  if (symbol) params.set('symbol', symbol)
  return apiRequest<DataSourceQualityRegression>(`/api/v1/data-source-quality/regression?${params.toString()}`)
}

export function getDataQualityGateResolution(symbol?: string) {
  const params = new URLSearchParams()
  if (symbol) params.set('symbol', symbol)
  const suffix = params.toString() ? `?${params.toString()}` : ''
  return apiRequest<DataQualityGateResolutionCheck>(`/api/v1/data-source-quality/gate-resolution${suffix}`)
}

export function listDataQualityGateResolutions(symbol?: string, status?: string) {
  const params = new URLSearchParams()
  if (symbol) params.set('symbol', symbol)
  if (status) params.set('status', status)
  const suffix = params.toString() ? `?${params.toString()}` : ''
  return apiRequest<DataQualityGateResolutionListResponse>(`/api/v1/data-source-quality/resolutions${suffix}`)
}

export function createDataQualityGateResolution(payload: DataQualityGateResolutionCreateRequest) {
  return apiRequest<DataQualityGateResolutionCheck>('/api/v1/data-source-quality/resolutions', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function retireDataQualityGateResolution(resolutionId: string) {
  return apiRequest<DataQualityGateResolutionCheck>(`/api/v1/data-source-quality/resolutions/${encodeURIComponent(resolutionId)}/retire`, {
    method: 'POST',
  })
}
