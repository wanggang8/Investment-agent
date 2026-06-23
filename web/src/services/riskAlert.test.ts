import { afterEach, describe, expect, it, vi } from 'vitest'
import { getRiskAlert, listRiskAlerts, updateRiskAlertLifecycle } from './riskAlert'

afterEach(() => {
  vi.unstubAllGlobals()
  vi.clearAllMocks()
})

function jsonResponse(data: unknown) {
  return new Response(JSON.stringify({ request_id: 'rid_risk', data }), {
    status: 200,
    headers: { 'content-type': 'application/json' },
  })
}

describe('risk alert service contract', () => {
  it('lists risk alerts with optional status and symbol filters', async () => {
    const fetchMock = vi.fn(async () => jsonResponse({ items: [], total: 0 }))
    vi.stubGlobal('fetch', fetchMock)

    await listRiskAlerts()
    await listRiskAlerts({ statuses: ['active', 'escalated'], symbol: '510300' })

    expect(fetchMock).toHaveBeenNthCalledWith(1, '/api/v1/risk-alerts', expect.objectContaining({ headers: expect.any(Headers) }))
    expect(fetchMock).toHaveBeenNthCalledWith(2, '/api/v1/risk-alerts?status=active%2Cescalated&symbol=510300', expect.objectContaining({ headers: expect.any(Headers) }))
  })

  it('gets detail and updates lifecycle with encoded alert id', async () => {
    const fetchMock = vi.fn(async () => jsonResponse({ alert_id: 'risk/id 1', sop_status: 'resolved' }))
    vi.stubGlobal('fetch', fetchMock)

    await getRiskAlert('risk/id 1')
    await updateRiskAlertLifecycle('risk/id 1', { status: 'resolved', reason: '人工复核完成' })

    expect(fetchMock).toHaveBeenNthCalledWith(1, '/api/v1/risk-alerts/risk%2Fid%201', expect.objectContaining({ headers: expect.any(Headers) }))
    expect(fetchMock).toHaveBeenNthCalledWith(2, '/api/v1/risk-alerts/risk%2Fid%201/lifecycle', expect.objectContaining({ method: 'POST', body: JSON.stringify({ status: 'resolved', reason: '人工复核完成' }) }))
  })
})
