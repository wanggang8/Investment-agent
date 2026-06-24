import { cleanup, render, screen, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { APIClientError } from '../services/client'
import { DecisionLoopPage } from './DecisionLoopPage'

vi.mock('../services/decisionLoop', () => ({
  getDecisionLoop: vi.fn(),
  listDecisionLoops: vi.fn(),
}))

import { getDecisionLoop, listDecisionLoops } from '../services/decisionLoop'

describe('DecisionLoopPage', () => {
  afterEach(() => {
    cleanup()
    vi.clearAllMocks()
  })

  it('renders loop stages, manual records, gaps and trace links without write actions', async () => {
    vi.mocked(listDecisionLoops).mockResolvedValue({
      request_id: 'rid_loop',
      data: {
        total: 1,
        safety_note: '只读解释链，仅展示本地事实和导航，不改变事实状态。',
        items: [{
          decision_id: 'decision_loop_1',
          symbol: '510300',
          generated_at: '2026-06-16T09:00:00Z',
          final_verdict_status: 'hold',
          final_verdict_text: '继续持有，等待人工复核',
          confirmation_status: 'executed_manually',
          loop_status: 'reviewed',
          safety_note: '只读解释链，仅展示本地事实和导航，不改变事实状态。',
          stages: [
            { stage: 'recommendation', status: 'complete', label: '建议生成', summary: '继续持有', ref_type: 'decision', ref_id: 'decision_loop_1', at: '2026-06-16T09:00:00Z' },
            { stage: 'confirmation', status: 'complete', label: '用户记录', summary: '已记录线下处理', ref_type: 'confirmation', ref_id: 'conf_loop_1', at: '2026-06-16T09:20:00Z' },
            { stage: 'manual_record', status: 'complete', label: '线下记录', summary: '1 条本地流水', ref_type: 'transaction', ref_id: 'tx_loop_1', at: '2026-06-16T09:30:00Z' },
            { stage: 'risk_review', status: 'complete', label: '风险线索', summary: '1 条风险线索', ref_type: 'risk_alert', ref_id: 'risk_loop_1' },
            { stage: 'review', status: 'complete', label: '复盘线索', summary: '已有复盘或审计线索', ref_type: 'error_case', ref_id: 'err_loop_1' },
          ],
          manual_actions: [{
            confirmation_id: 'conf_loop_1',
            confirmation_type: 'executed_manually',
            operation_type: 'buy',
            symbol: '510300',
            quantity: 10,
            price: 2.5,
            fees: 1,
            executed_at: '2026-06-16T09:30:00Z',
            transaction_ids: ['tx_loop_1'],
            note_preview: '人工记录已脱敏',
          }],
          risk_links: [{ type: 'risk_alert', id: 'risk_loop_1', label: '估值风险', href: '/risk-alerts/risk_loop_1', status: 'active' }],
          review_links: [{ type: 'error_case', id: 'err_loop_1', label: '错误案例', href: '/review#error_case-err_loop_1', status: 'reviewed' }],
          audit_links: [{ type: 'audit_event', id: 'audit_loop_1', label: '审计事件', href: '/audit#audit_loop_1', status: 'success' }],
          missing_links: [],
        }],
      },
    })

    render(<MemoryRouter><DecisionLoopPage /></MemoryRouter>)

    expect(await screen.findByRole('heading', { name: '决策闭环解释' })).toBeInTheDocument()
    expect(screen.getByText('只读解释链，仅展示本地事实和导航，不改变事实状态。')).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: '只读决策生命周期' })).toBeInTheDocument()
    expect(screen.getByText(/建议生成 -> 用户确认 -> 线下记录 -> 风险\/复盘 -> 审计/)).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '查看决策详情' })).toHaveAttribute('href', '/decisions/decision_loop_1')
    expect(screen.getByText('decision_loop_1 · 510300')).toBeInTheDocument()
    expect(screen.getByText('继续持有，等待人工复核')).toBeInTheDocument()
    expect(screen.getByText('建议生成')).toBeInTheDocument()
    expect(screen.getByText('用户记录')).toBeInTheDocument()
    expect(screen.getByText('线下记录')).toBeInTheDocument()
    expect(screen.getByText('人工记录已脱敏')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '估值风险' })).toHaveAttribute('href', '/risk-alerts/risk_loop_1')
    expect(screen.getByRole('link', { name: '错误案例' })).toHaveAttribute('href', '/review#error_case-err_loop_1')
    expect(screen.getByRole('link', { name: '审计事件' })).toHaveAttribute('href', '/audit#audit_loop_1')
    expect(screen.queryByRole('button', { name: /交易|下单|确认|应用|推送/ })).not.toBeInTheDocument()
    expect(document.body.textContent).not.toMatch(/SELECT \* FROM|\/Users\/private|sk-|prompt:/)
  })

  it('shows empty and gap states safely', async () => {
    vi.mocked(listDecisionLoops).mockResolvedValue({
      request_id: 'rid_empty',
      data: {
        total: 1,
        safety_note: '只读解释链，仅展示本地事实和导航，不改变事实状态。',
        items: [{
          decision_id: 'decision_loop_gap',
          generated_at: '2026-06-16T10:00:00Z',
          final_verdict_status: 'hold',
          final_verdict_text: '持有但缺少线下记录',
          confirmation_status: 'executed_manually',
          loop_status: 'incomplete',
          safety_note: '只读解释链，仅展示本地事实和导航，不改变事实状态。',
          stages: [{ stage: 'manual_record', status: 'missing', label: '线下记录', summary: '缺少本地流水记录' }],
          manual_actions: [],
          risk_links: [],
          review_links: [],
          audit_links: [],
          missing_links: ['缺少用户确认记录', '缺少线下记录'],
        }],
      },
    })

    render(<MemoryRouter><DecisionLoopPage /></MemoryRouter>)

    expect(await screen.findByText('decision_loop_gap')).toBeInTheDocument()
    expect(screen.getByText('缺少用户确认记录')).toBeInTheDocument()
    expect(screen.getByText('缺少线下记录')).toBeInTheDocument()
    expect(screen.getByText('暂无人工处理记录。')).toBeInTheDocument()
    expect(screen.getByText('暂无风险、复盘或审计链接。')).toBeInTheDocument()
  })

  it('shows safe API errors', async () => {
    vi.mocked(listDecisionLoops).mockRejectedValue(new APIClientError({
      requestId: 'rid_loop_error',
      code: 'INTERNAL_ERROR',
      message: '系统暂时无法处理请求，请稍后重试。',
      displayState: 'generic_failure',
    }))

    render(<MemoryRouter><DecisionLoopPage /></MemoryRouter>)

    await waitFor(() => expect(screen.getByText('系统暂时无法处理请求，请稍后重试。')).toBeInTheDocument())
  })

  it('focuses a decision loop from decision_id query without rendering the full list', async () => {
    vi.mocked(getDecisionLoop).mockResolvedValue({
      request_id: 'rid_focus',
      data: {
        decision_id: 'decision_focus',
        symbol: '510300',
        generated_at: '2026-06-24T02:29:48Z',
        final_verdict_status: 'hold',
        final_verdict_text: '按纪律观察',
        confirmation_status: 'planned',
        loop_status: 'reviewed',
        safety_note: '只读解释链，仅展示本地事实和导航，不改变事实状态。',
        stages: [{ stage: 'recommendation', status: 'complete', label: '建议生成', summary: '按纪律观察', ref_type: 'decision', ref_id: 'decision_focus' }],
        manual_actions: [],
        risk_links: [],
        review_links: [],
        audit_links: [{ type: 'audit_event', id: 'audit_focus', label: '审计事件', href: '/audit#audit_focus', status: 'success' }],
        missing_links: [],
      },
    })

    render(<MemoryRouter initialEntries={['/decision-loop?decision_id=decision_focus']}><DecisionLoopPage /></MemoryRouter>)

    await waitFor(() => expect(getDecisionLoop).toHaveBeenCalledWith('decision_focus'))
    expect(listDecisionLoops).not.toHaveBeenCalled()
    expect(screen.getByText('当前聚焦：decision_focus · 510300')).toBeInTheDocument()
    expect(screen.getByText('decision_focus · 510300')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '审计事件' })).toHaveAttribute('href', '/audit#audit_focus')
  })
})
