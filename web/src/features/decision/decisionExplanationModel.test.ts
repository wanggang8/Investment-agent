import { describe, expect, it } from 'vitest'
import type { DecisionDetailResponse } from '../../types/decision'
import { buildDecisionExplanationModel } from './decisionExplanationModel'

const baseDecision: DecisionDetailResponse = {
  decision_id: 'decision_story_1',
  question: '510300 是否继续持有？',
  symbol: '510300',
  generated_at: '2026-06-18T08:00:00Z',
  capability_check: { status: 'in_scope', reason: '沪深300 ETF 在能力圈内' },
  workflow_status: 'completed',
  triggered_rules: [
    { rule_id: 'valuation_guard', rule_name: '估值纪律', severity: 'warning', description: '估值分位偏高，暂停新增买入' },
  ],
  evidence_chain: [
    {
      evidence_id: 'ev_story_1',
      source_name: '交易所公告',
      source_level: 'A',
      evidence_role: 'formal',
      verification_status: 'satisfied',
      summary: '公告显示成分股调入完成',
      high_grade_independent_source_count: 1,
    },
    {
      evidence_id: 'ev_story_2',
      source_name: '本地笔记',
      source_level: 'C',
      evidence_role: 'background',
      verification_status: 'background_only',
      summary: '用户长期配置背景',
    },
  ],
  analyst_reports: [
    {
      agent_name: '价值分析师',
      conclusion: '估值偏高，继续观察。',
      key_reasons: ['估值分位高于纪律阈值', '成交流动性正常'],
      risk_warnings: ['追高风险'],
      confidence: 'medium',
      evidence_ids: ['ev_story_1'],
      quality_status: 'passed',
      parse_status: 'parsed',
    },
  ],
  retrieval_quality: {
    query_summary: '510300 持有复核',
    top_k: 2,
    status: 'hit',
    index_health: 'healthy',
    index_freshness: 'fresh',
    fallback_source: 'veclite',
    source_consistency_status: 'checked',
  },
  expected_return_scenarios: {
    sample_count: 12,
    sample_window: '2025-01-01 至 2026-06-18',
    screening_condition: '同类 ETF 样本',
    precision_status: 'available',
    scenarios: [],
    disclaimer: '历史样本仅作参考，不构成收益承诺。',
  },
  arbitration_chain: [
    { priority: 1, rule_id: 'valuation_guard', result: '暂停新增买入' },
  ],
  audit_events: [
    { audit_event_id: 'audit_story_1', action: 'generate_decision', status: 'success', created_at: '2026-06-18T08:01:00Z', node_name: 'DecisionRecordNode' },
  ],
  final_verdict: {
    status: 'hold',
    display_text: '继续持有，等待人工复核',
    prohibited_actions: ['自动交易', '新增买入'],
    optional_actions: ['记录观察计划'],
  },
  user_confirmation: { confirmation_status: 'pending', available_actions: ['planned'] },
}

describe('buildDecisionExplanationModel', () => {
  it('turns a decision into a story-first explanation model', () => {
    const model = buildDecisionExplanationModel(baseDecision)

    expect(model.storyTitle).toBe('继续持有，等待人工复核')
    expect(model.decisionContext).toEqual([
      '决策 decision_story_1',
      '标的 510300',
      '问题：510300 是否继续持有？',
      '生成时间：2026-06-18T08:00:00Z',
    ])
    expect(model.keyReasons).toEqual([
      '估值分位高于纪律阈值',
      '成交流动性正常',
      '估值纪律：估值分位偏高，暂停新增买入',
      '交易所公告：公告显示成分股调入完成',
    ])
    expect(model.prohibitedActions).toEqual(['自动交易', '新增买入'])
    expect(model.optionalActions).toEqual(['记录观察计划'])
    expect(model.trustSummary).toEqual([
      '正式证据 1 条，背景材料 1 条',
      '最高信源等级 A',
      '检索状态 hit，召回 2 条',
      'LLM 材料 1 份；解析/质量通过 1 份',
    ])
    expect(model.explanationLinks).toEqual([
      { label: '查看证据', href: '/evidence' },
      { label: '查看决策闭环', href: '/decision-loop' },
      { label: '查看审计', href: '/audit' },
    ])
  })

  it('uses safe empty states for nullable real LLM-like fields', () => {
    const model = buildDecisionExplanationModel({
      ...baseDecision,
      decision_id: 'decision_nullable_story',
      evidence_chain: null,
      analyst_reports: [{
        agent_name: '趋势风控官',
        conclusion: '证据不足，暂停交易类建议。',
        key_reasons: null,
        risk_warnings: null,
        confidence: 'low',
        evidence_ids: null,
        quality_status: 'failed',
      }] as unknown as DecisionDetailResponse['analyst_reports'],
      triggered_rules: null,
      final_verdict: {
        status: 'insufficient_data',
        display_text: '证据不足，等待人工复核',
        prohibited_actions: null,
        optional_actions: null,
      },
    } as unknown as DecisionDetailResponse)

    expect(model.keyReasons).toEqual(['趋势风控官：证据不足，暂停交易类建议。'])
    expect(model.prohibitedActions).toEqual(['暂无；缺失字段不代表允许交易或自动执行'])
    expect(model.optionalActions).toEqual(['暂无；仅可人工复核'])
    expect(model.trustSummary).toContain('正式证据 0 条，背景材料 0 条')
    expect(model.missingDataWarnings).toContain('缺少可展示的正式证据，最终结论需要人工复核。')
    expect(model.safetyNotes).toContain('缺失、降级或 nullable 字段不会被解释为允许交易、自动确认或自动应用规则。')
  })

  it('keeps fallback analyst conclusions compact in the story summary', () => {
    const longConclusion = '第一段很长的分析材料。'.repeat(20)
    const model = buildDecisionExplanationModel({
      ...baseDecision,
      analyst_reports: [{
        agent_name: '趋势风控官',
        conclusion: longConclusion,
        key_reasons: [],
        risk_warnings: [],
        confidence: 'qualitative',
        evidence_ids: [],
        quality_status: 'passed',
        parse_status: 'parsed',
      }],
    })

    expect(model.keyReasons[0]).toMatch(/^趋势风控官：/)
    expect(model.keyReasons[0].length).toBeLessThanOrEqual(125)
    expect(model.keyReasons[0]).toContain('...')
  })
})
