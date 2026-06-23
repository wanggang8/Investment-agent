import { describe, expect, it } from 'vitest'
import { buildAuditOpsModel } from './auditOpsModel'
import type { AuditEvent } from '../../types/audit'

const events: AuditEvent[] = [
  { request_id: 'req_1', audit_event_id: 'audit_1', actor: 'system', action: 'refresh_market_data', status: 'success', created_at: '2026-06-18T08:00:00Z', workflow_type: 'market' },
  { request_id: 'req_2', audit_event_id: 'audit_2', actor: 'gatekeeper', action: 'audit_rule_change', status: 'failed', error_code: 'RULE_REJECTED', created_at: '2026-06-18T09:00:00Z', proposal_id: 'prop_1' },
  { request_id: 'req_3', audit_event_id: 'audit_3', actor: 'system', action: 'run_local_task', status: 'degraded', error_code: 'DATA_STALE', created_at: '2026-06-18T10:00:00Z' },
]

describe('buildAuditOpsModel', () => {
  it('builds audit summary, categories, and inspection actions', () => {
    const model = buildAuditOpsModel({ events })

    expect(model.overallTone).toBe('warning')
    expect(model.overallLabel).toBe('审计记录需要检查')
    expect(model.metrics).toContainEqual(expect.objectContaining({ label: '审计事件', value: '3' }))
    expect(model.metrics).toContainEqual(expect.objectContaining({ label: '失败/降级', value: '2' }))
    expect(model.eventGroups.map((group) => group.label)).toEqual(expect.arrayContaining(['刷新市场数据', '审计规则变更', '运行本地任务']))
    expect(model.nextActions.map((action) => action.label)).toEqual(expect.arrayContaining(['查看失败事件', '检查规则治理', '查看数据质量']))
  })

  it('returns a safe empty state', () => {
    const model = buildAuditOpsModel({ events: [] })

    expect(model.overallTone).toBe('unknown')
    expect(model.overallLabel).toBe('暂无审计记录')
    expect(JSON.stringify(model)).not.toMatch(/raw stack|SELECT \* FROM|sk-/)
  })
})

