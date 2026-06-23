import { afterEach, describe, expect, it, vi } from 'vitest'
import {
  getDailyDisciplineReport,
  getTodayDailyDisciplineReport,
  listDailyDisciplineReports,
} from './dailyDisciplineReport'

afterEach(() => {
  vi.unstubAllGlobals()
  vi.clearAllMocks()
})

function jsonResponse(data: unknown) {
  return new Response(JSON.stringify({ request_id: 'rid_report', data }), {
    status: 200,
    headers: { 'content-type': 'application/json' },
  })
}

describe('daily discipline report service contract', () => {
  it('queries today report endpoint', async () => {
    const fetchMock = vi.fn(async () => jsonResponse({ report_id: 'daily_report:today', status: 'success' }))
    vi.stubGlobal('fetch', fetchMock)

    const res = await getTodayDailyDisciplineReport()

    expect(fetchMock).toHaveBeenCalledWith('/api/v1/daily-discipline/reports/today', expect.objectContaining({ headers: expect.any(Headers) }))
    expect(res.data?.report_id).toBe('daily_report:today')
  })

  it('lists reports with default limit and optional status', async () => {
    const fetchMock = vi.fn(async () => jsonResponse({ reports: [] }))
    vi.stubGlobal('fetch', fetchMock)

    await listDailyDisciplineReports()
    await listDailyDisciplineReports(12, 'insufficient_data')

    expect(fetchMock).toHaveBeenNthCalledWith(1, '/api/v1/daily-discipline/reports?limit=30', expect.objectContaining({ headers: expect.any(Headers) }))
    expect(fetchMock).toHaveBeenNthCalledWith(2, '/api/v1/daily-discipline/reports?limit=12&status=insufficient_data', expect.objectContaining({ headers: expect.any(Headers) }))
  })

  it('encodes detail report id as a single path segment', async () => {
    const fetchMock = vi.fn(async () => jsonResponse({ report_id: 'daily/report:2026-06-08:ETF 300', status: 'success' }))
    vi.stubGlobal('fetch', fetchMock)

    await getDailyDisciplineReport('daily/report:2026-06-08:ETF 300')

    expect(fetchMock).toHaveBeenCalledWith('/api/v1/daily-discipline/reports/daily%2Freport%3A2026-06-08%3AETF%20300', expect.objectContaining({ headers: expect.any(Headers) }))
  })
})
