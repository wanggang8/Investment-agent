import { describe, expect, it, vi, afterEach } from 'vitest'
import { apiRequest, APIClientError, mapErrorCodeToDisplayState } from './client'

afterEach(() => {
  vi.unstubAllGlobals()
  vi.clearAllMocks()
})

describe('apiRequest', () => {
  it('normalizes success envelope request id', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => new Response(JSON.stringify({ request_id: 'rid_api', data: { ok: true } }), { status: 200, headers: { 'content-type': 'application/json' } })))

    const res = await apiRequest<{ ok: boolean }>('/api/example')

    expect(res.request_id).toBe('rid_api')
    expect(res.data?.ok).toBe(true)
  })

  it('maps error envelope to safe display state', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => new Response(JSON.stringify({ request_id: 'rid_error', error: { code: 'VECTOR_INDEX_UNAVAILABLE', message: 'raw', detail: 'internal path' } }), { status: 409, headers: { 'content-type': 'application/json' } })))

    await expect(apiRequest('/api/example')).rejects.toMatchObject({
      name: 'APIClientError',
      requestId: 'rid_error',
      code: 'VECTOR_INDEX_UNAVAILABLE',
      message: '索引暂不可用，请稍后重试或重建索引。',
    } satisfies Partial<APIClientError>)
  })

  it('maps known UI states', () => {
    expect(mapErrorCodeToDisplayState('SOURCE_VERIFICATION_FAILED')).toBe('frozen_watch')
    expect(mapErrorCodeToDisplayState('DATA_SOURCE_UNAVAILABLE')).toBe('data_source_unavailable')
    expect(mapErrorCodeToDisplayState('RULE_VERSION_MISSING')).toBe('high_risk')
    expect(mapErrorCodeToDisplayState('UNKNOWN', 503)).toBe('data_source_unavailable')
  })
})
