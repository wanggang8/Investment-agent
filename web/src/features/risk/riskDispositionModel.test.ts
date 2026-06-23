import { describe, expect, it } from 'vitest'
import { buildRiskDispositionModel } from './riskDispositionModel'
import type { RiskAlert } from '../../types/riskAlert'

function risk(overrides: Partial<RiskAlert>): RiskAlert {
  return {
    alert_id: overrides.alert_id ?? 'risk_1',
    risk_type: overrides.risk_type ?? 'valuation_high',
    severity: overrides.severity ?? 'warning',
    sop_status: overrides.sop_status ?? 'active',
    symbol: overrides.symbol ?? '510300',
    trigger_summary: overrides.trigger_summary ?? '估值分位高，需要人工复核',
    prohibited_actions: overrides.prohibited_actions ?? ['新增买入'],
    suggested_actions: overrides.suggested_actions ?? ['人工复核止盈计划'],
    safety_note: overrides.safety_note ?? '风险预警只用于本地人工复核，不会自动交易。',
    created_at: overrides.created_at ?? '2026-06-18T09:30:00Z',
    updated_at: overrides.updated_at ?? '2026-06-18T09:30:00Z',
    ...overrides,
  }
}

describe('buildRiskDispositionModel', () => {
  it('groups alerts into disposition queues by SOP status', () => {
    const model = buildRiskDispositionModel([
      risk({ alert_id: 'risk_pending', sop_status: 'triggered' }),
      risk({ alert_id: 'risk_active', sop_status: 'active' }),
      risk({ alert_id: 'risk_observing', sop_status: 'observing' }),
      risk({ alert_id: 'risk_escalated', sop_status: 'escalated', severity: 'critical', symbol: '159915' }),
      risk({ alert_id: 'risk_resolved', sop_status: 'resolved', severity: 'info' }),
      risk({ alert_id: 'risk_archived', sop_status: 'archived', severity: 'info' }),
    ])

    expect(model.summaryLabel).toBe('6 条风险事实，1 条需复盘')
    expect(model.highestSeverity).toBe('critical')
    expect(model.affectedSymbols).toEqual(['159915', '510300'])
    expect(model.queues.find((queue) => queue.id === 'pending_review')?.items.map((item) => item.alert_id)).toEqual(['risk_pending'])
    expect(model.queues.find((queue) => queue.id === 'in_progress')?.items.map((item) => item.alert_id)).toEqual(['risk_active', 'risk_observing'])
    expect(model.queues.find((queue) => queue.id === 'needs_review')?.items.map((item) => item.alert_id)).toEqual(['risk_escalated'])
    expect(model.queues.find((queue) => queue.id === 'recorded')?.items.map((item) => item.alert_id)).toEqual(['risk_resolved', 'risk_archived'])
  })

  it('shows safe empty state without execution affordance', () => {
    const model = buildRiskDispositionModel([])

    expect(model.summaryLabel).toBe('暂无风险预警')
    expect(model.highestSeverity).toBe('info')
    expect(model.nextActions[0]).toMatchObject({ label: '继续观察今日纪律' })
    expect(model.safetyNotes.join(' ')).toContain('不会自动交易')
  })

  it('prioritizes escalated risk review before routine observing', () => {
    const model = buildRiskDispositionModel([
      risk({ alert_id: 'risk_observing', sop_status: 'observing' }),
      risk({ alert_id: 'risk_escalated', sop_status: 'escalated', severity: 'critical' }),
    ])

    expect(model.nextActions[0]).toMatchObject({ label: '复核升级风险', priority: 'blocking' })
  })
})
