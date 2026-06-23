import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { APIClientError } from '../services/client'
import { LocalKnowledgePage } from './LocalKnowledgePage'

vi.mock('../services/localKnowledge', () => ({
  confirmLocalKnowledgeImport: vi.fn(),
  validateLocalKnowledgeImport: vi.fn(),
}))

import {
  confirmLocalKnowledgeImport,
  validateLocalKnowledgeImport,
} from '../services/localKnowledge'

describe('LocalKnowledgePage', () => {
  afterEach(() => {
    cleanup()
    vi.clearAllMocks()
  })

  it('validates redacted preview and confirms with the bound batch id', async () => {
    vi.mocked(validateLocalKnowledgeImport).mockResolvedValue({
      request_id: 'rid_validate',
      data: {
        import_batch_id: 'lk_batch_123',
        summary: {
          total_count: 1,
          valid_count: 1,
          warning_count: 0,
          blocking_count: 0,
        },
        rows: [
          {
            row_number: 1,
            status: 'valid',
            symbol: '510300',
            title_preview: '510300 估值观察',
            text_preview: '本地研究记录：[REDACTED] 后续仅作为背景材料参与检索。',
            content_hash: 'hash_1',
            risks: [],
          },
        ],
        index_plan: {
          rag_chunk_count: 1,
          index_status: 'pending',
        },
        safety_note: 'server safety note',
      },
    })
    vi.mocked(confirmLocalKnowledgeImport).mockResolvedValue({
      request_id: 'rid_confirm',
      data: {
        import_batch_id: 'lk_batch_123',
        intelligence_item_count: 1,
        summary_count: 1,
        rag_chunk_count: 1,
        verification_count: 1,
        audit_event_ids: ['audit_1'],
        index_status: 'pending',
        safety_note: 'server safety note',
      },
    })

    render(<LocalKnowledgePage />)

    expect(screen.getByRole('heading', { name: '本地知识导入' })).toBeInTheDocument()
    expect(screen.getByText('本地配置与诊断状态')).toBeInTheDocument()
    expect(screen.getAllByText('脱敏预览').length).toBeGreaterThan(0)
    expect(screen.getAllByText('索引计划').length).toBeGreaterThan(0)

    fireEvent.click(screen.getByRole('button', { name: '校验预览' }))

    expect(await screen.findByText('lk_batch_123')).toBeInTheDocument()
    expect(screen.getByText('复核知识导入')).toBeInTheDocument()
    expect(screen.getByText('知识预览')).toBeInTheDocument()
    expect(screen.getByText('本地研究记录：[REDACTED] 后续仅作为背景材料参与检索。')).toBeInTheDocument()
    expect(screen.getByText('预计片段：1')).toBeInTheDocument()
    expect(screen.getByText('仅写入本地背景材料，后续需由用户在相关页面人工复核。')).toBeInTheDocument()

    fireEvent.click(screen.getByRole('button', { name: '写入本地事实' }))

    await waitFor(() => expect(confirmLocalKnowledgeImport).toHaveBeenCalledWith({
      import_batch_id: 'lk_batch_123',
      confirm_reason: '人工确认导入为本地背景材料',
      source_label: 'local_research_notes',
      default_symbol: '510300',
      rows: [
        {
          title: '510300 估值观察',
          text: '本地研究记录：指数估值处于偏高区间，后续仅作为背景材料参与检索。',
          symbol: '510300',
          as_of_date: '',
          tags: ['估值', '本地研究'],
        },
      ],
    }))
    expect(await screen.findByText('情报：1 条；摘要：1 条；片段：1 条。')).toBeInTheDocument()
    expect(screen.getByText('核验记录：1 条；审计事件：1 条。')).toBeInTheDocument()
  })

  it('keeps blocking rows out of confirm and displays safe errors', async () => {
    vi.mocked(validateLocalKnowledgeImport).mockResolvedValue({
      request_id: 'rid_blocking',
      data: {
        import_batch_id: 'lk_batch_blocking',
        summary: {
          total_count: 1,
          valid_count: 0,
          warning_count: 0,
          blocking_count: 1,
        },
        rows: [
          {
            row_number: 1,
            status: 'blocking',
            symbol: '510300',
            title_preview: '异常记录',
            text_preview: '[REDACTED]',
            content_hash: 'hash_blocking',
            risks: [{ code: 'unsafe_content', severity: 'blocking', message: '包含不适合入库的敏感片段。' }],
          },
        ],
        index_plan: {
          rag_chunk_count: 0,
          index_status: 'pending',
        },
        safety_note: 'server safety note',
      },
    })

    render(<LocalKnowledgePage />)
    fireEvent.click(screen.getByRole('button', { name: '校验预览' }))

    expect((await screen.findAllByText('阻断')).length).toBeGreaterThan(0)
    expect(screen.getByRole('button', { name: '写入本地事实' })).toBeDisabled()
    expect(confirmLocalKnowledgeImport).not.toHaveBeenCalled()

    vi.mocked(validateLocalKnowledgeImport).mockRejectedValue(new APIClientError({
      requestId: 'rid_error',
      code: 'INVALID_STATE',
      message: '当前状态不允许执行该操作。',
      displayState: 'frozen_watch',
    }))
    fireEvent.change(screen.getByLabelText('记录 JSON'), { target: { value: '[{"text":"新的本地记录"}]' } })
    fireEvent.click(screen.getByRole('button', { name: '校验预览' }))

    expect(await screen.findByText('当前状态不允许执行该操作。')).toBeInTheDocument()
  })
})
