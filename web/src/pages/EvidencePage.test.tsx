import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { APIClientError } from '../services/client'
import { EvidencePage } from './EvidencePage'

vi.mock('../services/evidence', () => ({
  getEvidenceVerification: vi.fn(),
  listEvidence: vi.fn(),
  rebuildEvidenceIndex: vi.fn(),
  refreshEvidence: vi.fn(),
}))

import { getEvidenceVerification, listEvidence, rebuildEvidenceIndex, refreshEvidence } from '../services/evidence'

describe('EvidencePage', () => {
  afterEach(() => {
    cleanup()
    vi.clearAllMocks()
  })

  it('shows vector unavailable API error', async () => {
    vi.mocked(getEvidenceVerification).mockResolvedValue({ request_id: 'rid', data: { verification_id: 'ver_1', verification_status: 'failed', independent_source_count: 0, high_grade_independent_source_count: 0, highest_source_level: 'C', latest_published_at: '', evidence_ids: [] } })
    vi.mocked(listEvidence).mockRejectedValue(new APIClientError({ requestId: 'rid', code: 'VECTOR_INDEX_UNAVAILABLE', message: '索引暂不可用，请稍后重试或重建索引。', displayState: 'insufficient_data' }))

    render(<EvidencePage />)

    await waitFor(() => expect(screen.getByText('索引暂不可用，请稍后重试或重建索引。')).toBeInTheDocument())
  })

  it('shows empty success state', async () => {
    vi.mocked(getEvidenceVerification).mockResolvedValue({ request_id: 'rid', data: { verification_id: 'ver_1', verification_status: 'failed', independent_source_count: 0, high_grade_independent_source_count: 0, highest_source_level: 'C', latest_published_at: '', evidence_ids: [] } })
    vi.mocked(listEvidence).mockResolvedValue({ request_id: 'rid', data: { items: [], total: 0 } })

    render(<EvidencePage />)

    await waitFor(() => expect(screen.getByText(/暂无匹配证据/)).toBeInTheDocument())
  })

  it('supports evidence refresh, index rebuild and source verification panel', async () => {
    vi.mocked(listEvidence)
      .mockResolvedValueOnce({ request_id: 'rid', data: { items: [{ evidence_id: 'ev_1', source_name: '交易所', source_level: 'A', evidence_role: 'formal', verification_status: 'satisfied', published_at: '2026-01-01T00:00:00Z', summary: '公告摘要', high_grade_independent_source_count: 1 }], total: 1 } })
      .mockResolvedValue({ request_id: 'rid2', data: { items: [{ evidence_id: 'ev_2', source_name: '交易所', source_level: 'A', evidence_role: 'formal', verification_status: 'satisfied', published_at: '2026-01-02T00:00:00Z', summary: '刷新后摘要', high_grade_independent_source_count: 1 }], total: 1 } })
    vi.mocked(getEvidenceVerification)
      .mockResolvedValueOnce({ request_id: 'rid', data: { verification_id: 'ver_1', verification_status: 'satisfied', independent_source_count: 2, high_grade_independent_source_count: 1, highest_source_level: 'A', latest_published_at: '2026-01-01T00:00:00Z', evidence_ids: ['ev_1'] } })
      .mockResolvedValue({ request_id: 'rid2', data: { verification_id: 'ver_2', verification_status: 'satisfied', independent_source_count: 3, high_grade_independent_source_count: 2, highest_source_level: 'S', latest_published_at: '2026-01-02T00:00:00Z', evidence_ids: ['ev_2'] } })
    vi.mocked(refreshEvidence).mockResolvedValue({ request_id: 'refresh', data: { intelligence_item_count: 1, summary_count: 1, rag_chunk_count: 1, verification_count: 1, index_status: 'indexed', audit_event_ids: ['audit_ev'] } })
    vi.mocked(rebuildEvidenceIndex).mockResolvedValue({ request_id: 'rebuild', data: { indexed_count: 1, skipped_count: 0, audit_event_ids: ['audit_idx'] } })

    render(<EvidencePage />)

    expect(await screen.findByRole('heading', { name: '证据可信度' })).toBeInTheDocument()
    expect(screen.getByText(/先看信源等级、独立信源和核验状态，再进入证据明细/)).toBeInTheDocument()
    await waitFor(() => expect(screen.getAllByText('独立信源数量：2').length).toBeGreaterThanOrEqual(1))
    expect(screen.getAllByText('最高信源等级：A').length).toBeGreaterThanOrEqual(1)
    expect(screen.getByText('证据引用：ev_1')).toBeInTheDocument()
    expect(screen.getByText(/S\/A\/B 级可作为正式证据/)).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '返回工作台' })).toHaveAttribute('href', '/workbench')
    expect(screen.getByRole('link', { name: '查看决策闭环' })).toHaveAttribute('href', '/decision-loop')
    expect(screen.getByRole('link', { name: '查看审计' })).toHaveAttribute('href', '/audit')

    fireEvent.click(screen.getByRole('button', { name: '刷新情报' }))
    await waitFor(() => expect(screen.getByText('情报刷新完成；索引状态 indexed。')).toBeInTheDocument())
    await waitFor(() => expect(screen.getByText('刷新后摘要')).toBeInTheDocument())
    expect(screen.getAllByText('独立信源数量：3').length).toBeGreaterThanOrEqual(1)
    expect(refreshEvidence).toHaveBeenCalledWith({ refresh_scope: 'all', include_background: true })

    fireEvent.click(screen.getByRole('button', { name: '重建索引' }))
    await waitFor(() => expect(screen.getByText('索引重建完成；已索引 1 条，跳过 0 条。')).toBeInTheDocument())
    expect(rebuildEvidenceIndex).toHaveBeenCalled()
    expect(listEvidence).toHaveBeenCalledTimes(3)
  })
})
