import type { PageResult } from '../types/api'
import type { RiskAlert, RiskAlertLifecycleRequest, RiskAlertListFilter } from '../types/riskAlert'
import { apiRequest } from './client'

export function listRiskAlerts(filter: RiskAlertListFilter = {}) {
  const params = new URLSearchParams()
  if (filter.statuses?.length) {
    params.set('status', filter.statuses.join(','))
  }
  if (filter.symbol) {
    params.set('symbol', filter.symbol)
  }
  const qs = params.toString()
  return apiRequest<PageResult<RiskAlert>>(`/api/v1/risk-alerts${qs ? `?${qs}` : ''}`)
}

export function getRiskAlert(alertId: string) {
  return apiRequest<RiskAlert>(`/api/v1/risk-alerts/${encodeURIComponent(alertId)}`)
}

export function updateRiskAlertLifecycle(alertId: string, input: RiskAlertLifecycleRequest) {
  return apiRequest<RiskAlert>(`/api/v1/risk-alerts/${encodeURIComponent(alertId)}/lifecycle`, {
    method: 'POST',
    body: JSON.stringify(input),
  })
}
