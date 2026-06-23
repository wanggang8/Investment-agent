import { describe, expect, it } from 'vitest'
import { buildDailyWorkbenchModel } from './dailyWorkbenchModel'
import type { DashboardTodayResponse } from '../../types/dashboard'

const baseDashboard: DashboardTodayResponse = {
  dashboard_state: 'normal',
  discipline_status: '正常',
  data_updated_at: '2026-06-17T09:30:00Z',
  portfolio_summary: {
    total_assets: 100000,
    cash_ratio: 0.32,
    high_risk_ratio: 0.12,
    position_count: 2,
  },
  market_summary: {
    sentiment_state: 'neutral',
    liquidity_state: 'normal',
  },
  evidence_summary: {
    source_count: 3,
    highest_source_level: 'A',
    verification_status: 'satisfied',
  },
  triggered_rules: [],
  decision_summary: {
    decision_id: 'decision_today',
    verdict: '继续持有，等待人工复核',
    final_verdict_status: 'hold',
    prohibited_actions: ['追涨买入'],
    optional_actions: ['复核仓位'],
    action_required: true,
    confirmation_status: 'pending',
  },
}

describe('daily workbench view model', () => {
  it('builds an at-a-glance success model with manual actions', () => {
    const model = buildDailyWorkbenchModel({
      dashboard: baseDashboard,
      report: {
        report_id: 'report_today',
        local_date: '2026-06-17',
        scope: 'daily',
        status: 'success',
        summary: '今日纪律报告已生成',
        final_verdict: '继续持有，等待人工复核',
        evidence: { evidence_count: 4, independent_source_count: 3, high_grade_independent_source_count: 2 },
        trend: { success_count: 1, degraded_count: 0, failed_count: 0, insufficient_data_count: 0 },
        safety_note: '只读报告。',
      },
      portfolio: {
        snapshot: { snapshot_id: 'snap', snapshot_time: '2026-06-17T09:30:00Z', cash: 32000, total_assets: 100000, cash_ratio: 0.32, high_risk_ratio: 0.12, position_count: 2 },
        positions: [],
      },
      risks: [],
      rules: [],
      review: { decision_count: 8, error_case_count: 0, rule_proposal_count: 0 },
    })

    expect(model.statusLabel).toBe('成功')
    expect(model.statusTone).toBe('success')
    expect(model.verdictText).toBe('继续持有，等待人工复核')
    expect(model.trustSummary).toBe('证据 4 条 · 独立信源 3 个')
    expect(model.prohibitedActions).toEqual(['追涨买入'])
    expect(model.optionalActions).toEqual(['复核仓位'])
    expect(model.nextActions.map((action) => action.label)).toContain('查看决策详情')
    expect(model.nextActions.map((action) => action.label)).toContain('查看今日纪律报告')
    expect(model.nextActions.every((action) => action.href.startsWith('/'))).toBe(true)
  })

  it('prioritizes blocking steps when local facts or evidence are missing', () => {
    const model = buildDailyWorkbenchModel({
      dashboard: {
        ...baseDashboard,
        dashboard_state: 'insufficient_data',
        portfolio_summary: { total_assets: 0, cash_ratio: 0, high_risk_ratio: 0, position_count: 0 },
        evidence_summary: undefined,
        decision_summary: {
          ...baseDashboard.decision_summary,
          decision_id: undefined,
          verdict: '等待本地数据补齐后生成今日纪律建议。',
          final_verdict_status: 'insufficient_data',
          prohibited_actions: ['暂停交易类建议'],
          optional_actions: [],
        },
      },
      report: {
        report_id: 'report_gap',
        local_date: '2026-06-17',
        scope: 'daily',
        status: 'insufficient_data',
        summary: '缺少持仓快照。',
        missing_action: '请先录入本地账户和当前持仓。',
        evidence: { evidence_count: 0, independent_source_count: 0, high_grade_independent_source_count: 0 },
        trend: { success_count: 0, degraded_count: 0, failed_count: 0, insufficient_data_count: 1 },
        safety_note: '不会自动执行交易。',
      },
      risks: [],
      rules: [],
      review: { decision_count: 0, error_case_count: 0, rule_proposal_count: 0 },
    })

    expect(model.statusLabel).toBe('数据不足')
    expect(model.statusTone).toBe('warning')
    expect(model.trustSummary).toBe('证据不足或尚未完成核验')
    expect(model.nextActions[0]).toMatchObject({ label: '维护本地账户与持仓', href: '/positions', priority: 'blocking' })
    expect(model.nextActions.map((action) => action.label)).toContain('查看数据质量')
    expect(model.prohibitedActions).toEqual(['暂停交易类建议'])
  })

  it('does not style high risk or degraded states as success', () => {
    const highRiskModel = buildDailyWorkbenchModel({
      dashboard: { ...baseDashboard, dashboard_state: 'high_risk' },
      risks: [{
        alert_id: 'risk_1',
        risk_type: 'valuation_high',
        severity: 'warning',
        sop_status: 'active',
        symbol: '510300',
        trigger_summary: '估值分位偏高。',
        safety_note: '仅人工复核。',
        created_at: '2026-06-17T09:30:00Z',
        updated_at: '2026-06-17T09:30:00Z',
      }],
      rules: [],
      review: { decision_count: 0, error_case_count: 0, rule_proposal_count: 0 },
    })

    const degradedModel = buildDailyWorkbenchModel({
      dashboard: { ...baseDashboard, dashboard_state: 'data_source_unavailable' },
      report: {
        report_id: 'report_degraded',
        local_date: '2026-06-17',
        scope: 'daily',
        status: 'degraded',
        summary: '数据源降级。',
        evidence: { evidence_count: 1, independent_source_count: 0, high_grade_independent_source_count: 0 },
        trend: { success_count: 0, degraded_count: 1, failed_count: 0, insufficient_data_count: 0 },
        safety_note: '仅人工复核。',
      },
      risks: [],
      rules: [],
      review: { decision_count: 0, error_case_count: 0, rule_proposal_count: 0 },
    })

    expect(highRiskModel.statusTone).toBe('danger')
    expect(highRiskModel.nextActions.map((action) => action.label)).toContain('处理风险预警')
    expect(degradedModel.statusTone).toBe('degraded')
    expect(degradedModel.statusLabel).toBe('降级')
  })
})
