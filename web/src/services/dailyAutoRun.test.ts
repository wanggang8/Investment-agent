import { afterEach, describe, expect, it, vi } from 'vitest'
import { getDailyAutoRunStatus } from './dailyAutoRun'

afterEach(() => {
  vi.unstubAllGlobals()
  vi.clearAllMocks()
})

describe('daily auto-run service contract', () => {
  it('queries the daily auto-run status endpoint', async () => {
    const fetchMock = vi.fn(async () => new Response(JSON.stringify({
      request_id: 'rid_auto_run',
      data: { enabled: true, status: 'failed', failure_code: 'missing_prerequisites', safety_note: '仅记录本地每日刷新和纪律评估结果，不会自动执行交易，需人工复核。' },
    }), { status: 200, headers: { 'content-type': 'application/json' } }))
    vi.stubGlobal('fetch', fetchMock)

    const res = await getDailyAutoRunStatus()

    expect(fetchMock).toHaveBeenCalledWith('/api/v1/daily-auto-run/status', expect.objectContaining({ headers: expect.any(Headers) }))
    expect(res.data?.status).toBe('failed')
    expect(res.data?.failure_code).toBe('missing_prerequisites')
  })
})
