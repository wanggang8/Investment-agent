import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { RiskAlertPage } from './RiskAlertPage'

vi.mock('../services/riskAlert', () => ({
  getRiskAlert: vi.fn(),
  listRiskAlerts: vi.fn(),
  updateRiskAlertLifecycle: vi.fn(),
}))

import { getRiskAlert, listRiskAlerts, updateRiskAlertLifecycle } from '../services/riskAlert'

describe('RiskAlertPage', () => {
  afterEach(() => {
    cleanup()
    vi.clearAllMocks()
  })

  it('loads alert detail when alert id exists in the route', async () => {
    vi.mocked(getRiskAlert).mockResolvedValue({ request_id: 'rid_detail', data: {
      alert_id: 'risk_detail',
      risk_type: 'liquidity_danger',
      severity: 'critical',
      sop_status: 'active',
      symbol: '510300',
      trigger_summary: '流动性处于 danger，禁止市价式大额操作',
      prohibited_actions: ['市价式大额操作'],
      suggested_actions: ['人工复核流动性'],
      safety_note: '风险预警只用于本地人工复核，不会自动交易。',
      created_at: '2026-06-15T10:30:00Z',
      updated_at: '2026-06-15T10:30:00Z',
    } })

    render(
      <MemoryRouter initialEntries={['/risk-alerts/risk_detail']}>
        <Routes><Route path="/risk-alerts/:alertId" element={<RiskAlertPage />} /></Routes>
      </MemoryRouter>,
    )

    await waitFor(() => expect(getRiskAlert).toHaveBeenCalledWith('risk_detail'))
    expect(screen.getByText('流动性危险')).toBeInTheDocument()
    expect(screen.getByText('流动性处于 danger，禁止市价式大额操作')).toBeInTheDocument()
    expect(listRiskAlerts).not.toHaveBeenCalled()
  })

  it('shows active and escalated alerts with safe SOP actions', async () => {
    vi.mocked(listRiskAlerts).mockResolvedValue({ request_id: 'rid_risk', data: { total: 2, items: [
      {
        alert_id: 'risk_1',
        risk_type: 'valuation_high',
        severity: 'warning',
        sop_status: 'active',
        symbol: '510300',
        trigger_summary: '估值分位处于高位，暂停新增买入并复核止盈计划',
        trigger_context: { sop: 'SOP-B', data_prerequisites: ['profit_ratio', 'valuation_percentile'], llm_role: 'explain_only' },
        prohibited_actions: ['新增买入'],
        suggested_actions: ['人工复核分批止盈'],
        related_report_id: 'report_1',
        report_link: '/daily-discipline/reports/report_1',
        related_decision_id: 'dec_1',
        decision_link: '/decisions/dec_1',
        safety_note: '风险预警只用于本地人工复核，不会自动交易。',
        created_at: '2026-06-15T09:30:00Z',
        updated_at: '2026-06-15T09:30:00Z',
      },
      {
        alert_id: 'risk_2',
        risk_type: 'data_degraded',
        severity: 'critical',
        sop_status: 'escalated',
        symbol: '159915',
        trigger_summary: '数据源新鲜度降级',
        prohibited_actions: ['交易类建议'],
        suggested_actions: ['检查数据源'],
        safety_note: '风险预警只用于本地人工复核，不会自动交易。',
        created_at: '2026-06-15T10:30:00Z',
        updated_at: '2026-06-15T10:30:00Z',
      },
    ] } })
    vi.mocked(updateRiskAlertLifecycle).mockResolvedValue({ request_id: 'rid_update', data: { alert_id: 'risk_1', sop_status: 'observing' } as never })

    render(<MemoryRouter><RiskAlertPage /></MemoryRouter>)

    expect(await screen.findByRole('heading', { name: '风险预警中心' })).toBeInTheDocument()
    expect(screen.getByText('风险处置队列')).toBeInTheDocument()
    expect(screen.getByText('2 条风险事实，1 条需复盘')).toBeInTheDocument()
    expect(screen.getByText('最高严重程度：严重')).toBeInTheDocument()
    expect(screen.getByText('处理中队列')).toBeInTheDocument()
    expect(screen.getByText('需复盘队列')).toBeInTheDocument()
    expect(screen.getByText('估值高位')).toBeInTheDocument()
    expect(screen.getByText('SOP：SOP-B')).toBeInTheDocument()
    expect(screen.getByText('数据前提：profit_ratio、valuation_percentile')).toBeInTheDocument()
    expect(screen.getByText('LLM 角色：explain_only')).toBeInTheDocument()
    expect(screen.getByText('已升级')).toBeInTheDocument()
    expect(screen.getByText('禁止动作：新增买入')).toBeInTheDocument()
    expect(screen.getByText('建议人工动作：人工复核分批止盈')).toBeInTheDocument()
    expect(screen.getAllByText(/不会自动交易/)[0]).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '关联报告' })).toHaveAttribute('href', '/daily-discipline/reports/report_1')

    fireEvent.click(screen.getAllByRole('button', { name: '记录继续观察' })[0])
    await waitFor(() => expect(updateRiskAlertLifecycle).toHaveBeenCalledWith('risk_1', { status: 'observing', reason: '前端人工 SOP 操作：记录继续观察' }))
  })

  it('renders resolved and archived alerts without lifecycle controls', async () => {
    vi.mocked(listRiskAlerts).mockResolvedValue({ request_id: 'rid_done', data: { total: 2, items: [
      { alert_id: 'risk_resolved', risk_type: 'insufficient_evidence', severity: 'warning', sop_status: 'resolved', symbol: '510300', trigger_summary: '证据恢复', safety_note: '风险预警只用于本地人工复核，不会自动交易。', created_at: '2026-06-15T09:30:00Z', updated_at: '2026-06-15T09:30:00Z' },
      { alert_id: 'risk_archived', risk_type: 'position_limit_breach', severity: 'info', sop_status: 'archived', symbol: '159915', trigger_summary: '历史风险归档', safety_note: '风险预警只用于本地人工复核，不会自动交易。', created_at: '2026-06-15T09:30:00Z', updated_at: '2026-06-15T09:30:00Z' },
    ] } })

    render(<MemoryRouter><RiskAlertPage /></MemoryRouter>)

    expect(await screen.findByText('已解除')).toBeInTheDocument()
    expect(screen.getByText('已归档')).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: '记录继续观察' })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: '记录升级复核' })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: '记录本地解除预警' })).not.toBeInTheDocument()
  })

  it('shows empty state without trading action controls', async () => {
    vi.mocked(listRiskAlerts).mockResolvedValue({ request_id: 'rid_empty', data: { total: 0, items: [] } })

    render(<MemoryRouter><RiskAlertPage /></MemoryRouter>)

    expect(await screen.findByRole('region', { name: '暂无需要处置的风险预警' })).toHaveTextContent('当前没有需要处置的本地风险。')
    expect(await screen.findByText('暂无风险预警')).toBeInTheDocument()
    expect(screen.getByText('继续观察今日纪律')).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: /买入|卖出|交易/ })).not.toBeInTheDocument()
  })
})
