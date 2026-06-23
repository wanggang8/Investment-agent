import { render, screen, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { cleanup } from '@testing-library/react'
import { APIClientError } from '../services/client'
import { AuditPage } from './AuditPage'

vi.mock('../services/audit', () => ({
  listAuditEvents: vi.fn(),
}))

import { listAuditEvents } from '../services/audit'

describe('AuditPage', () => {
  afterEach(() => {
    cleanup()
    vi.clearAllMocks()
  })

  it('shows API error state', async () => {
    vi.mocked(listAuditEvents).mockRejectedValue(new APIClientError({ requestId: 'rid', code: 'ANALYST_UNAVAILABLE', message: '分析服务暂不可用，页面仅展示规则与已有数据。', displayState: 'insufficient_data' }))

    render(<AuditPage />)

    await waitFor(() => expect(screen.getByText('分析服务暂不可用，页面仅展示规则与已有数据。')).toBeInTheDocument())
  })

  it('shows empty success state', async () => {
    vi.mocked(listAuditEvents).mockResolvedValue({ request_id: 'rid', data: { items: [], total: 0 } })

    render(<AuditPage />)

    await waitFor(() => expect(screen.getByText('审计检查状态')).toBeInTheDocument())
    expect(screen.getByText('暂无审计记录')).toBeInTheDocument()
    expect(screen.getByText('查看数据质量')).toBeInTheDocument()
    await waitFor(() => expect(screen.getByText('暂无匹配审计事件。')).toBeInTheDocument())
    expect(screen.queryByText('分析服务暂不可用，页面仅展示规则与已有数据。')).not.toBeInTheDocument()
  })

  it('shows audit summary and inspection actions before timeline', async () => {
    vi.mocked(listAuditEvents).mockResolvedValue({ request_id: 'rid', data: { items: [
      { request_id: 'req_1', audit_event_id: 'audit_1', actor: 'system', action: 'refresh_market_data', status: 'success', created_at: '2026-06-18T08:00:00Z' },
      { request_id: 'req_2', audit_event_id: 'audit_2', actor: 'gatekeeper', action: 'audit_rule_change', status: 'failed', proposal_id: 'prop_1', created_at: '2026-06-18T09:00:00Z' },
    ], total: 2 } })

    render(<AuditPage />)

    await waitFor(() => expect(screen.getByText('审计记录需要检查')).toBeInTheDocument())
    expect(screen.getByText('失败/降级')).toBeInTheDocument()
    expect(screen.getByText('查看失败事件')).toBeInTheDocument()
    expect(screen.getByText('检查规则治理')).toBeInTheDocument()
    expect(screen.getAllByText('审计规则变更').length).toBeGreaterThan(0)
  })
})
