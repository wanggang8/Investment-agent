import { afterEach, describe, expect, it, vi } from 'vitest'
import { getEvidenceVerification, listEvidence } from './evidence'

afterEach(() => {
  vi.unstubAllGlobals()
  vi.clearAllMocks()
})

describe('evidence service contract', () => {
  it('preserves high-grade source count from evidence list response', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => new Response(JSON.stringify({
      request_id: 'rid_ev_list',
      data: { items: [{ evidence_id: 'sum_1', source_name: '公告', source_level: 'A', summary: '摘要', high_grade_independent_source_count: 2 }], total: 1 },
    }), { status: 200, headers: { 'content-type': 'application/json' } })))

    const res = await listEvidence()

    expect(res.data?.items[0].high_grade_independent_source_count).toBe(2)
  })

  it('returns one verification DTO from verification endpoint', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => new Response(JSON.stringify({
      request_id: 'rid_ev_verification',
      data: { verification_id: 'ver_1', verification_status: 'satisfied', independent_source_count: 3, high_grade_independent_source_count: 2, highest_source_level: 'S', latest_published_at: '2026-01-02T00:00:00Z', evidence_ids: ['sum_1'] },
    }), { status: 200, headers: { 'content-type': 'application/json' } })))

    const res = await getEvidenceVerification()

    expect(res.data?.verification_id).toBe('ver_1')
    expect(res.data?.high_grade_independent_source_count).toBe(2)
    expect(res.data?.evidence_ids).toEqual(['sum_1'])
  })
})
