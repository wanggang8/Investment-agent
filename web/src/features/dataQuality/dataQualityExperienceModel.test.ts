import { describe, expect, it } from 'vitest'
import { buildDataQualityExperienceModel } from './dataQualityExperienceModel'
import type { EvidenceItem, SourceVerification } from '../../types/evidence'
import type { MarketSnapshot, SourceHealthItem } from '../../types/market'
import type { ReviewSummary } from '../../types/review'
import type { SystemStatus } from '../../types/settings'
import type { DataQualityGateResolutionCheck, DataSourceQualityRegression } from '../../types/dataSourceQuality'

const system: SystemStatus = {
  sqlite_status: 'ok',
  veclite_status: 'configured',
  deepseek_status: 'configured',
  data_sources: ['stub', 'csindex_extended'],
  log_level: 'error',
}

const market: MarketSnapshot = {
  market_snapshot_id: 'market_1',
  symbol: '510300',
  trade_date: '2026-06-18',
  pe_percentile: 0.8,
  pb_percentile: 0.7,
  liquidity_state: 'normal',
  sentiment_state: 'neutral',
  data_status: 'fresh',
}

const sourceHealth: SourceHealthItem[] = [{
  source_name: 'csindex_extended',
  source_level: 'A',
  source_type: 'public_file',
  data_category: 'index_valuation_files',
  freshness: 'stale',
  data_date: '2026-06-05',
  last_success_at: '2026-06-05T15:00:00Z',
  last_failure_at: '2026-06-18T08:00:00Z',
  failure_category: 'parse_error',
  affected_symbols: ['510300'],
}]

const verification: SourceVerification = {
  verification_id: 'verify_1',
  verification_status: 'satisfied',
  independent_source_count: 3,
  high_grade_independent_source_count: 2,
  highest_source_level: 'A',
  latest_published_at: '2026-06-18T09:00:00Z',
  evidence_ids: ['evidence_1'],
}

const evidence: EvidenceItem[] = [{
  evidence_id: 'evidence_1',
  source_name: 'CNInfo',
  source_level: 'A',
  evidence_role: 'formal',
  verification_status: 'satisfied',
  summary: '公告证据摘要',
}]

const review: ReviewSummary = {
  decision_count: 4,
  confirmation_count: 1,
  executed_manually_count: 0,
  planned_count: 1,
  error_case_count: 0,
  rule_proposal_count: 1,
  audit_event_count: 14,
  missing_evidence_count: 2,
  degraded_count: 1,
  ops_status: {
    data_source_status: 'degraded',
    index_status: 'success',
    review_status: 'degraded',
    explanation: 'sk-test-token SQL SELECT * FROM secrets /Users/private/raw vendor response',
  },
  degraded_workflows: [{ decision_id: 'decision_degraded', symbol: '510300', status: 'degraded', created_at: '2026-06-18T09:00:00Z' }],
  tracking_links: [],
  recent_decisions: [],
}

const blockedRegression: DataSourceQualityRegression = {
  mode: 'current',
  status: 'degraded',
  generated_at: '2026-06-18T10:00:00Z',
  summary: '数据源质量回归 mode=current status=degraded cases=1 degraded=1 failed=0',
  cases: [],
  missing_categories: ['index_valuation_files'],
  policy: {
    verdict: 'blocked',
    release_gate: 'block',
    degraded_count: 1,
    failed_count: 0,
    blocking_count: 1,
    waiver_count: 0,
    blocking_reasons: ['index_valuation_files core category degraded freshness=parse_error'],
    waiver_reasons: [],
    next_actions: ['不得把当前数据源质量声明为 clean'],
    safety_note: '当前数据质量策略只读取本地 source health。',
  },
  safety_note: '只读检查',
}

const scopeExclusionGate: DataQualityGateResolutionCheck = {
  symbol: '510300',
  policy_fingerprint: 'fp_test',
  policy_summary: 'data_source_quality:mode=current:status=degraded:policy=blocked:gate=block:no_auto_trading',
  policy: blockedRegression.policy,
  release_claim_state: 'resolved_with_scope_exclusion',
  clean_data_claim_allowed: false,
  active_resolution: {
    resolution_id: 'dqgr_1',
    symbol: '510300',
    policy_fingerprint: 'fp_test',
    policy_verdict: 'blocked',
    release_gate: 'block',
    policy_summary: 'data_source_quality:mode=current:status=degraded:policy=blocked:gate=block:no_auto_trading',
    resolution_type: 'scope_exclusion',
    status: 'active',
    scope: '本次 release clean claim 排除 current local data health',
    reason: '当前本地数据源存在降级',
    release_impact: '不得声明 current data healthy',
    created_by: 'local_user',
    created_at: '2026-06-18T10:00:00Z',
    safety_note: '当前数据门禁处置只记录本地人工声明边界。',
  },
  allowed_claims: ['可以声明当前本地数据健康已排除在 clean claim 外'],
  prohibited_claims: ['不得声明当前本地数据 clean', '不得声明 current data healthy'],
  safety_note: '当前数据门禁处置只记录本地人工声明边界。',
}

describe('buildDataQualityExperienceModel', () => {
  it('summarizes degraded source, evidence, LLM, and affected workflow signals', () => {
    const model = buildDataQualityExperienceModel({
      system,
      market,
      sourceHealth,
      evidenceItems: evidence,
      evidenceTotal: 1,
      verification,
      review,
    })

    expect(model.overallTone).toBe('warning')
    expect(model.overallLabel).toBe('数据质量需要检查')
    expect(model.qualitySignals.map((signal) => signal.label)).toEqual(['数据源健康', '证据与 RAG', 'LLM 分析', '影响范围'])
    expect(model.qualitySignals).toContainEqual(expect.objectContaining({ label: '数据源健康', value: '1 项需检查', tone: 'warning' }))
    expect(model.qualitySignals).toContainEqual(expect.objectContaining({ label: '影响范围', value: '1 个受影响工作流', tone: 'warning' }))
    expect(model.nextActions.map((action) => action.label)).toEqual(expect.arrayContaining(['查看数据源设置', '查看受影响决策', '查看质量复盘', '查看风险预警']))
  })

  it('promotes blocked current data policy to danger release-gate signal', () => {
    const model = buildDataQualityExperienceModel({
      system,
      market,
      sourceHealth,
      evidenceItems: evidence,
      evidenceTotal: 1,
      verification,
      review,
      currentRegression: blockedRegression,
    })

    expect(model.overallTone).toBe('danger')
    expect(model.overallLabel).toBe('当前数据质量阻断发布声明')
    expect(model.qualitySignals).toContainEqual(expect.objectContaining({ label: '当前数据策略', value: '阻断', tone: 'danger' }))
    expect(model.warnings).toContain('index_valuation_files core category degraded freshness=parse_error')
    expect(model.nextActions.map((action) => action.label)).toEqual(expect.arrayContaining(['查看当前数据策略']))
  })

  it('treats waiver_required policy as warning rather than clean pass', () => {
    const model = buildDataQualityExperienceModel({
      system,
      market,
      sourceHealth: [],
      evidenceItems: evidence,
      evidenceTotal: 1,
      verification,
      review: { ...review, degraded_count: 0, missing_evidence_count: 0, degraded_workflows: [] },
      currentRegression: {
        ...blockedRegression,
        status: 'degraded',
        missing_categories: ['sentiment_proxy'],
        policy: {
          ...blockedRegression.policy,
          verdict: 'waiver_required',
          release_gate: 'waiver_required',
          blocking_count: 0,
          waiver_count: 1,
          blocking_reasons: [],
          waiver_reasons: ['sentiment_proxy optional category degraded freshness=stale'],
          next_actions: ['发布材料中记录 waiver reason'],
        },
      },
    })

    expect(model.overallTone).toBe('warning')
    expect(model.overallLabel).toBe('数据质量需要检查')
    expect(model.qualitySignals).toContainEqual(expect.objectContaining({ label: '当前数据策略', value: '需豁免记录', tone: 'warning' }))
    expect(JSON.stringify(model)).not.toMatch(/自动修复|自动确认|自动应用规则|一键交易|代下单/)
  })

  it('does not expose raw diagnostics or private material in model text', () => {
    const model = buildDataQualityExperienceModel({
      system,
      market,
      sourceHealth,
      evidenceItems: evidence,
      evidenceTotal: 1,
      verification,
      review,
    })

    const text = JSON.stringify(model)
    expect(text).not.toMatch(/sk-test-token|SELECT \* FROM|\/Users\/private|raw vendor response/)
    expect(text).not.toMatch(/自动修复|自动确认|自动应用规则|自动规则应用|一键交易|代下单|收益承诺/)
    expect(text).toContain('存在已脱敏诊断摘要')
  })

  it('treats missing and unknown states as degraded instead of success', () => {
    const model = buildDataQualityExperienceModel({
      system: { ...system, veclite_status: 'missing', deepseek_status: 'unknown' },
      market: { ...market, data_status: 'missing' },
      sourceHealth: [],
      evidenceItems: [],
      evidenceTotal: 0,
      review: { ...review, degraded_count: 0, missing_evidence_count: 0, degraded_workflows: [] },
    })

    expect(model.overallTone).not.toBe('success')
    expect(model.qualitySignals.find((signal) => signal.label === '证据与 RAG')?.tone).toBe('warning')
    expect(model.qualitySignals.find((signal) => signal.label === 'LLM 分析')?.tone).toBe('warning')
  })

  it.each([
    { field: 'data_source_status', value: 'empty' },
    { field: 'index_status', value: 'failed' },
    { field: 'review_status', value: 'quality_failed' },
    { field: 'review_status', value: 'insufficient' },
    { field: 'review_status', value: 'insufficient_data' },
  ])('treats ops_status $field=$value as requiring inspection', ({ field, value }) => {
    const model = buildDataQualityExperienceModel({
      system,
      market,
      sourceHealth: [],
      evidenceItems: evidence,
      evidenceTotal: 1,
      verification,
      review: {
        ...review,
        missing_evidence_count: 0,
        degraded_count: 0,
        degraded_workflows: [],
        ops_status: {
          data_source_status: 'success',
          index_status: 'success',
          review_status: 'success',
          [field]: value,
        },
      },
    })

    expect(model.overallTone).not.toBe('success')
    expect(model.qualitySignals.find((signal) => signal.label === 'LLM 分析')?.tone).toBe('warning')
  })

  it('summarizes scope exclusion as limited release claim rather than clean pass', () => {
    const model = buildDataQualityExperienceModel({
      system,
      market,
      sourceHealth,
      evidenceItems: evidence,
      evidenceTotal: 1,
      verification,
      review,
      currentRegression: blockedRegression,
      gateResolution: scopeExclusionGate,
    })

    expect(model.overallTone).toBe('warning')
    expect(model.overallLabel).toBe('当前数据声明已限定范围')
    expect(model.qualitySignals).toContainEqual(expect.objectContaining({ label: '当前数据策略', value: '已排除 clean claim', tone: 'warning' }))
  })
})
