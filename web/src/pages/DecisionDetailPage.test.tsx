import { fireEvent, render, screen, waitFor, cleanup } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { APIClientError } from '../services/client'
import { createConfirmation, getDecision, consultDecision } from '../services/decision'
import type { DecisionDetailResponse } from '../types/decision'
import { DecisionDetailPage } from './DecisionDetailPage'

vi.mock('../services/decision', () => ({
  getDecision: vi.fn(),
  createConfirmation: vi.fn(),
  consultDecision: vi.fn(),
}))

const firstDecision: DecisionDetailResponse = {
  decision_id: 'decision_1',
  question: '要卖出吗',
  symbol: '510300',
  generated_at: '2026-05-29T03:00:00Z',
  workflow_status: 'completed',
  triggered_rules: [],
  evidence_chain: [],
  analyst_reports: [],
  arbitration_chain: [],
  audit_events: [{ audit_event_id: 'audit_1', action: 'generate_decision', status: 'success', created_at: '2026-05-29T03:01:00Z', node_name: 'DecisionRecordNode' }],
  final_verdict: { status: 'hold', display_text: '继续持有', prohibited_actions: [], optional_actions: [] },
  user_confirmation: { confirmation_status: 'pending', available_actions: ['executed_manually'] },
}

const refreshedDecision: DecisionDetailResponse = {
  ...firstDecision,
  user_confirmation: { confirmation_status: 'executed_manually', available_actions: [] },
}

describe('DecisionDetailPage', () => {
  afterEach(() => {
    cleanup()
    vi.clearAllMocks()
  })

  it('submits confirmation and refreshes decision detail', async () => {
    vi.mocked(getDecision).mockResolvedValueOnce({ request_id: 'req_1', data: firstDecision }).mockResolvedValueOnce({ request_id: 'req_2', data: refreshedDecision })
    vi.mocked(createConfirmation).mockResolvedValue({ request_id: 'req_confirm', data: { confirmation_id: 'confirm_1', decision_id: 'decision_1', confirmation_status: 'executed_manually', audit_event_ids: ['audit_1'] } })

    render(
      <MemoryRouter initialEntries={['/decisions/decision_1']}>
        <Routes>
          <Route path="/decisions/:decisionId" element={<DecisionDetailPage />} />
        </Routes>
      </MemoryRouter>,
    )

    await waitFor(() => expect(screen.getByText('继续持有')).toBeInTheDocument())
    fireEvent.click(screen.getByRole('button', { name: '已手动执行' }))
    fireEvent.change(screen.getByLabelText('标的代码'), { target: { value: '510300' } })
    fireEvent.change(screen.getByLabelText('线下动作'), { target: { value: 'sell' } })
    fireEvent.change(screen.getByLabelText('数量'), { target: { value: '4' } })
    fireEvent.change(screen.getByLabelText('价格'), { target: { value: '3.5' } })
    fireEvent.change(screen.getByLabelText('执行时间'), { target: { value: '2026-05-29T03:00' } })
    fireEvent.click(screen.getByRole('button', { name: '提交确认' }))

    await waitFor(() => expect(createConfirmation).toHaveBeenCalledWith('decision_1', expect.objectContaining({ confirmation_type: 'executed_manually', operation_type: 'sell' })))
    await waitFor(() => expect(getDecision).toHaveBeenCalledTimes(2))
    expect(await screen.findByText(/确认已记录/)).toBeInTheDocument()
    expect(screen.getByText(/当前确认状态：已手动执行/)).toBeInTheDocument()
    expect(screen.getByText('审计时间线')).toBeInTheDocument()
    expect(screen.getByText(/DecisionRecordNode/)).toBeInTheDocument()
  })

  it('submits consultation from consultation route', async () => {
    vi.mocked(consultDecision).mockResolvedValue({ request_id: 'req_consult', data: firstDecision })

    render(
      <MemoryRouter initialEntries={['/consultation']}>
        <Routes>
          <Route path="/consultation" element={<DecisionDetailPage />} />
        </Routes>
      </MemoryRouter>,
    )

    expect([...screen.getByLabelText('咨询场景').querySelectorAll('option')].map((option) => option.value)).toEqual(['hold_review', 'buy_review', 'sell_review', 'rebalance_review'])
    expect(screen.getByRole('heading', { name: '主动咨询' })).toBeInTheDocument()
    expect(screen.getByText(/输入假设/)).toBeInTheDocument()
    expect(screen.getByText(/不会自动交易、自动确认或自动应用规则/)).toBeInTheDocument()
    fireEvent.change(screen.getByLabelText('咨询问题'), { target: { value: '是否继续持有' } })
    fireEvent.change(screen.getByLabelText('标的代码'), { target: { value: '510300' } })
    fireEvent.change(screen.getByLabelText('咨询场景'), { target: { value: 'hold_review' } })
    fireEvent.click(screen.getByRole('button', { name: '提交咨询' }))

    await waitFor(() => expect(consultDecision).toHaveBeenCalledWith({ question: '是否继续持有', symbol: '510300', scenario: 'hold_review' }))
    expect(await screen.findByText('继续持有')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '打开生成的决策详情' })).toHaveAttribute('href', '/decisions/decision_1')
    expect(screen.getAllByRole('link', { name: '查看证据' })[0]).toHaveAttribute('href', '/evidence')
    expect(screen.getAllByRole('link', { name: '查看决策闭环' })[0]).toHaveAttribute('href', '/decision-loop')
  })

  it('keeps previous decision state when confirmation fails', async () => {
    vi.mocked(getDecision).mockResolvedValueOnce({ request_id: 'req_1', data: firstDecision })
    vi.mocked(createConfirmation).mockRejectedValue(new APIClientError({ requestId: 'req_fail', code: 'INVALID_STATE', message: '当前状态不允许执行该操作。', displayState: 'frozen_watch' }))

    render(
      <MemoryRouter initialEntries={['/decisions/decision_1']}>
        <Routes>
          <Route path="/decisions/:decisionId" element={<DecisionDetailPage />} />
        </Routes>
      </MemoryRouter>,
    )

    await waitFor(() => expect(screen.getByText('继续持有')).toBeInTheDocument())
    fireEvent.click(screen.getByRole('button', { name: '已手动执行' }))
    fireEvent.change(screen.getByLabelText('标的代码'), { target: { value: '510300' } })
    fireEvent.change(screen.getByLabelText('线下动作'), { target: { value: 'sell' } })
    fireEvent.change(screen.getByLabelText('数量'), { target: { value: '4' } })
    fireEvent.change(screen.getByLabelText('价格'), { target: { value: '3.5' } })
    fireEvent.change(screen.getByLabelText('执行时间'), { target: { value: '2026-05-29T03:00' } })
    fireEvent.click(screen.getByRole('button', { name: '提交确认' }))

    await waitFor(() => expect(screen.getByText('当前状态不允许执行该操作。')).toBeInTheDocument())
    expect(screen.getByText(/当前确认状态：待确认/)).toBeInTheDocument()
    expect(screen.queryByText(/确认已记录/)).not.toBeInTheDocument()
    expect(getDecision).toHaveBeenCalledTimes(1)
  })
})
