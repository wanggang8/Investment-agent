import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'
import { AuditEventTimeline } from './AuditEventTimeline'

const events = [
  { audit_event_id: 'a1', request_id: 'r1', actor: 'system', action: 'refresh_market_data', status: 'success', created_at: '2026-05-31' },
  {
    audit_event_id: 'a2',
    request_id: 'r2',
    actor: 'gatekeeper',
    action: 'audit_rule_change',
    status: 'degraded',
    error_code: 'VECTOR_INDEX_UNAVAILABLE',
    input_ref_type: 'symbol',
    input_ref: '510300',
    created_at: '2026-05-31',
    workflow_type: 'decision',
    before_state: 'pending',
    after_state: 'degraded',
    rule_version: 'v1.0',
    snapshot_id: 'snap_1',
    decision_id: 'decision_a2',
  },
]

describe('AuditEventTimeline', () => {
  afterEach(() => cleanup())

  it('filters and expands audit event refs', () => {
    render(<AuditEventTimeline events={events} />)

    fireEvent.change(screen.getByLabelText('筛选审计状态'), { target: { value: 'degraded' } })
    expect(screen.queryByText('刷新市场数据')).not.toBeInTheDocument()
    expect(screen.getByText('审计规则变更')).toBeInTheDocument()

    const toggle = screen.getByRole('button', { name: '展开引用' })
    expect(toggle).toHaveAttribute('aria-expanded', 'false')
    fireEvent.click(toggle)
    expect(toggle).toHaveAttribute('aria-expanded', 'true')
    expect(screen.getByText(/request_id：r2/)).toBeInTheDocument()
    expect(screen.getByText(/工作流：decision/)).toBeInTheDocument()
    expect(screen.getByText(/规则版本：v1.0/)).toBeInTheDocument()
    expect(screen.getByText(/快照：snap_1/)).toBeInTheDocument()
    expect(screen.getByText(/关联：决策 decision_a2/)).toBeInTheDocument()
    expect(screen.getByText(/引用：symbol 510300/)).toBeInTheDocument()
    expect(screen.getByText(/状态变化：pending → degraded/)).toBeInTheDocument()
  })

  it('maps audit actor status and unknown node action safely', () => {
    render(<AuditEventTimeline events={[{ audit_event_id: 'a3', request_id: 'r3', actor: 'system', action: 'retrieve_evidence', node_action: 'new_backend_action', status: 'degraded', created_at: '2026-05-31' }]} />)

    expect(screen.getByText(/执行方：系统/)).toBeInTheDocument()
    expect(screen.getByText(/状态：降级/)).toBeInTheDocument()
    expect(screen.getByText(/动作：未知动作/)).toBeInTheDocument()
    expect(screen.queryByText(/new_backend_action/)).not.toBeInTheDocument()
  })

  it('maps backend audit action and actor enums', () => {
    render(<AuditEventTimeline events={[{ audit_event_id: 'a4', request_id: 'r4', actor: 'gatekeeper', action: 'generate_decision', status: 'success', created_at: '2026-05-31' }, { audit_event_id: 'a5', request_id: 'r5', actor: 'system', action: 'create_proposal', status: 'success', created_at: '2026-05-31' }, { audit_event_id: 'a6', request_id: 'r6', actor: 'system', action: 'update_settings', status: 'success', created_at: '2026-05-31' }, { audit_event_id: 'a7', request_id: 'r7', actor: 'system', action: 'update_capability', status: 'success', created_at: '2026-05-31' }, { audit_event_id: 'a8', request_id: 'r8', actor: 'system', action: 'rebuild_index', status: 'success', created_at: '2026-05-31' }]} />)

    expect(screen.getByText('生成决策')).toBeInTheDocument()
    expect(screen.getByText('创建规则提案')).toBeInTheDocument()
    expect(screen.getByText('更新系统设置')).toBeInTheDocument()
    expect(screen.getByText('更新能力圈')).toBeInTheDocument()
    expect(screen.getByText('重建索引')).toBeInTheDocument()
    expect(screen.getByText(/执行方：守门人/)).toBeInTheDocument()
  })

  it('supports backend event_id as audit event identifier', () => {
    render(<AuditEventTimeline events={[{ event_id: 'legacy_audit_1', request_id: 'r9', actor: 'system', action: 'refresh_market_data', status: 'success', created_at: '2026-05-31', input_ref: '510300' }, { event_id: 'legacy_audit_2', request_id: 'r10', actor: 'system', action: 'retrieve_evidence', status: 'success', created_at: '2026-05-31', input_ref: '159915' }]} />)

    const buttons = screen.getAllByRole('button', { name: '展开引用' })
    fireEvent.click(buttons[1])
    expect(screen.getByText(/request_id：r10/)).toBeInTheDocument()
    expect(screen.getByText(/input：- 159915/)).toBeInTheDocument()
    expect(screen.queryByText(/request_id：r9/)).not.toBeInTheDocument()
  })
})
