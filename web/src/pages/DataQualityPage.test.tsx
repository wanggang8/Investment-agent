import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { APIClientError } from '../services/client'
import { DataQualityPage } from './DataQualityPage'

vi.mock('../services/settings', () => ({
  getSystemSettings: vi.fn(),
}))

vi.mock('../services/market', () => ({
  getLatestMarketSnapshot: vi.fn(),
  getMarketSourceHealth: vi.fn(),
}))

vi.mock('../services/evidence', () => ({
  getEvidenceVerification: vi.fn(),
  listEvidence: vi.fn(),
}))

vi.mock('../services/review', () => ({
  getReviewSummary: vi.fn(),
}))

vi.mock('../services/dataSourceQuality', () => ({
  createDataQualityGateResolution: vi.fn(),
  getDataQualityGateResolution: vi.fn(),
  getDataSourceQualityRegression: vi.fn(),
  listDataQualityGateResolutions: vi.fn(),
  retireDataQualityGateResolution: vi.fn(),
}))

vi.mock('../services/knowledgeReadiness', () => ({
  getKnowledgeReadiness: vi.fn(),
}))

import { getEvidenceVerification, listEvidence } from '../services/evidence'
import { getLatestMarketSnapshot, getMarketSourceHealth } from '../services/market'
import { getReviewSummary } from '../services/review'
import { getSystemSettings } from '../services/settings'
import { createDataQualityGateResolution, getDataQualityGateResolution, getDataSourceQualityRegression, listDataQualityGateResolutions, retireDataQualityGateResolution } from '../services/dataSourceQuality'
import { getKnowledgeReadiness } from '../services/knowledgeReadiness'

function knowledgeReadinessFixture(overrides: Partial<Awaited<ReturnType<typeof getKnowledgeReadiness>>['data']> = {}) {
  return {
    symbol: '510300',
    status: 'degraded',
    symbol_profile: {
      symbol: '510300',
      name: '沪深300ETF',
      asset_type: 'ETF',
      tracked_index_symbol: '000300',
      tracked_index_name: '沪深300',
      known: true,
    },
    knowledge_references: [
      {
        knowledge_id: 'master.graham.margin_of_safety',
        title: '格雷厄姆：安全边际',
        category: 'master_principle',
        summary: '估值分位越低，越需要安全边际。',
        applies_to: ['valuation_percentiles'],
        rule_mapping: ['valuation.low_zone'],
        llm_context_allowed: true,
        formal_evidence_allowed: false,
        safety_boundary: '只能作为纪律原则和 LLM 分析背景，不能作为正式市场证据。',
      },
      {
        knowledge_id: 'discipline.no_single_source_decision',
        title: '纪律：不凭单一信源决策',
        category: 'discipline_rule',
        summary: '重大事件必须满足多源验证。',
        applies_to: ['formal_evidence'],
        rule_mapping: ['evidence.min_high_grade_sources'],
        llm_context_allowed: true,
        formal_evidence_allowed: false,
        safety_boundary: '规则说明不是外部事实来源。',
      },
    ],
    data_dependencies: [
      {
        category: 'symbol_profile',
        status: 'ready',
        required: true,
        freshness: 'fresh',
        affected_features: ['consultation'],
        safe_degradation: '画像缺失时不生成正式交易类建议。',
      },
      {
        category: 'valuation_percentiles',
        status: 'degraded',
        required: true,
        source_level: 'A',
        freshness: 'parse_error',
        affected_features: ['margin_of_safety', 'expected_return'],
        safe_degradation: '估值分位缺失时不得声明安全边际或估值高低。',
      },
    ],
    feature_impacts: [
      {
        feature: 'margin_of_safety',
        category: 'valuation_percentiles',
        impact: '估值分位缺失时不得声明安全边际或估值高低。',
        claims: ['不得伪造成 ready'],
      },
    ],
    llm_context_summary: 'principles=master.graham.margin_of_safety; data_readiness=valuation_percentiles=degraded; boundary=背景知识不能满足正式证据',
    safety_notes: ['内置知识只作为纪律、规则映射和 LLM 分析上下文，不作为正式市场证据。'],
    ...overrides,
  }
}

function gateResolutionFixture(overrides: Partial<Awaited<ReturnType<typeof getDataQualityGateResolution>>['data']> = {}) {
  return {
    symbol: '510300',
    policy_fingerprint: 'fp_test',
    policy_summary: 'data_source_quality:mode=current:status=degraded:policy=blocked:gate=block:no_auto_trading',
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
    release_claim_state: 'requires_resolution',
    clean_data_claim_allowed: false,
    allowed_claims: ['可以声明需要人工处置当前数据门禁'],
    prohibited_claims: ['不得声明当前本地数据 clean', '不得声明 current data healthy', '不得把 resolution 描述为 policy passed'],
    safety_note: '当前数据门禁处置只记录本地人工声明边界，不改变数据质量事实、不刷新数据、不触发交易。',
    ...overrides,
  }
}

function mockMinimalDataQualityDependencies(symbol = '510300') {
  vi.mocked(getSystemSettings).mockResolvedValue({ request_id: 'rid_system_minimal', data: { sqlite_status: 'ok', veclite_status: 'configured', deepseek_status: 'configured', data_sources: ['csindex'], log_level: 'error' } })
  vi.mocked(getLatestMarketSnapshot).mockResolvedValue({ request_id: 'rid_market_minimal', data: { market_snapshot_id: 'market_minimal', symbol, trade_date: '2026-06-19', pe_percentile: 0.5, pb_percentile: 0.5, liquidity_state: 'normal', sentiment_state: 'neutral', data_status: 'fresh' } })
  vi.mocked(getMarketSourceHealth).mockResolvedValue({ request_id: 'rid_source_minimal', data: { sources: [] } })
  vi.mocked(getEvidenceVerification).mockResolvedValue({ request_id: 'rid_verification_minimal', data: { verification_id: 'verification_minimal', verification_status: 'satisfied', independent_source_count: 2, high_grade_independent_source_count: 2, highest_source_level: 'A', latest_published_at: '2026-06-19T08:00:00Z', evidence_ids: [] } })
  vi.mocked(listEvidence).mockResolvedValue({ request_id: 'rid_evidence_minimal', data: { total: 0, items: [] } })
  vi.mocked(getReviewSummary).mockResolvedValue({ request_id: 'rid_review_minimal', data: { decision_count: 0, confirmation_count: 0, executed_manually_count: 0, planned_count: 0, error_case_count: 0, rule_proposal_count: 0, audit_event_count: 0, ops_status: { data_source_status: 'success', index_status: 'success', review_status: 'success' }, recent_decisions: [] } })
  vi.mocked(getDataSourceQualityRegression).mockResolvedValue({ request_id: 'rid_policy_minimal', data: { mode: 'current', status: 'passed', generated_at: '2026-06-19T08:00:00Z', summary: 'passed', cases: [], missing_categories: [], policy: { verdict: 'passed', release_gate: 'pass', degraded_count: 0, failed_count: 0, blocking_count: 0, waiver_count: 0, blocking_reasons: [], waiver_reasons: [], next_actions: [], safety_note: '只读检查' }, safety_note: '只读检查' } })
  vi.mocked(getDataQualityGateResolution).mockResolvedValue({ request_id: 'rid_gate_minimal', data: gateResolutionFixture({ policy: { verdict: 'passed', release_gate: 'pass', degraded_count: 0, failed_count: 0, blocking_count: 0, waiver_count: 0, blocking_reasons: [], waiver_reasons: [], next_actions: [], safety_note: '只读检查' }, release_claim_state: 'pass', clean_data_claim_allowed: true }) })
  vi.mocked(listDataQualityGateResolutions).mockResolvedValue({ request_id: 'rid_resolutions_minimal', data: { items: [], total: 0 } })
}

describe('DataQualityPage', () => {
  beforeEach(() => {
    vi.mocked(getKnowledgeReadiness).mockResolvedValue({
      request_id: 'rid_knowledge_readiness',
      data: knowledgeReadinessFixture(),
    })
  })

  afterEach(() => {
    cleanup()
    vi.clearAllMocks()
  })

  it('aggregates source health, evidence, retrieval, LLM, and impact navigation', async () => {
    vi.mocked(getSystemSettings).mockResolvedValue({
      request_id: 'rid_system',
      data: {
        sqlite_status: 'ok',
        veclite_status: 'configured',
        deepseek_status: 'configured',
        data_sources: ['stub', 'csindex_extended'],
        log_level: 'error',
      },
    })
    vi.mocked(getLatestMarketSnapshot).mockResolvedValue({
      request_id: 'rid_market',
      data: {
        market_snapshot_id: 'market_1',
        symbol: '510300',
        trade_date: '2026-06-16',
        pe_percentile: 0.82,
        pb_percentile: 0.76,
        liquidity_state: 'warning',
        sentiment_state: 'hot',
        data_status: 'fresh',
      },
    })
    vi.mocked(getMarketSourceHealth).mockResolvedValue({
      request_id: 'rid_source',
      data: {
        sources: [
          {
            source_name: 'csindex_extended',
            source_level: 'A',
            source_type: 'public_file',
            data_category: 'index_valuation_files',
            freshness: 'stale',
            data_date: '2026-06-05',
            last_success_at: '2026-06-05T15:00:00Z',
            last_failure_at: '2026-06-16T08:00:00Z',
            failure_category: 'stale',
            affected_symbols: ['510300'],
          },
          {
            source_name: 'parse_source',
            source_level: 'B',
            source_type: 'public_file',
            data_category: 'constituent_financials',
            freshness: 'parse_error',
            data_date: '2026-06-16',
            failure_category: 'parse_error',
            affected_symbols: ['510500'],
          },
        ],
      },
    })
    vi.mocked(getEvidenceVerification).mockResolvedValue({
      request_id: 'rid_verification',
      data: {
        verification_id: 'verify_1',
        verification_status: 'satisfied',
        independent_source_count: 3,
        high_grade_independent_source_count: 2,
        highest_source_level: 'A',
        latest_published_at: '2026-06-15T09:00:00Z',
        evidence_ids: ['evidence_1', 'evidence_2'],
      },
    })
    vi.mocked(listEvidence).mockResolvedValue({
      request_id: 'rid_evidence',
      data: {
        total: 2,
        items: [
          {
            evidence_id: 'evidence_1',
            source_name: 'CNInfo',
            source_level: 'A',
            evidence_role: 'formal',
            verification_status: 'satisfied',
            published_at: '2026-06-15T09:00:00Z',
            captured_at: '2026-06-15T09:05:00Z',
            summary: '公告证据摘要',
            time_weight: 0.9,
          },
          {
            evidence_id: 'evidence_2',
            source_name: 'CSRC',
            source_level: 'A',
            evidence_role: 'formal',
            verification_status: 'satisfied',
            published_at: '2026-06-14T09:00:00Z',
            captured_at: '2026-06-14T09:05:00Z',
            summary: '监管证据摘要',
            time_weight: 0.8,
          },
        ],
      },
    })
    vi.mocked(getReviewSummary).mockResolvedValue({
      request_id: 'rid_review',
      data: {
        period: 'monthly',
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
          explanation: '复盘窗口内存在降级记录。',
        },
        degraded_workflows: [{ decision_id: 'decision_degraded', symbol: '510300', status: 'degraded', created_at: '2026-06-16T09:00:00Z' }],
        tracking_links: [{ type: 'audit_event', id: 'audit_1', label: '审计事件 audit_1' }],
        recent_decisions: [],
      },
    })
    vi.mocked(getDataSourceQualityRegression).mockResolvedValue({
      request_id: 'rid_current_policy',
      data: {
        mode: 'current',
        status: 'degraded',
        generated_at: '2026-06-18T10:00:00Z',
        summary: '数据源质量回归 mode=current status=degraded cases=2 degraded=1 failed=0',
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
      },
    })
    vi.mocked(getDataQualityGateResolution).mockResolvedValue({
      request_id: 'rid_gate',
      data: gateResolutionFixture(),
    })
    vi.mocked(listDataQualityGateResolutions).mockResolvedValue({
      request_id: 'rid_resolutions',
      data: { items: [], total: 0 },
    })

    render(<MemoryRouter><DataQualityPage /></MemoryRouter>)

    expect(await screen.findByRole('heading', { name: '数据质量可观测' })).toBeInTheDocument()
    expect(screen.getByText('数据质量总览')).toBeInTheDocument()
    expect(screen.getByText('当前数据质量阻断发布声明')).toBeInTheDocument()
    expect(screen.getByText('阻断')).toBeInTheDocument()
    expect(screen.getByText('需要人工处置')).toBeInTheDocument()
    expect(screen.getByText('clean data claim：不允许')).toBeInTheDocument()
    expect(screen.getByText('不得声明 current data healthy')).toBeInTheDocument()
    expect(screen.getByText('index_valuation_files core category degraded freshness=parse_error')).toBeInTheDocument()
    expect(screen.getByRole('article', { name: '数据源健康 信号' })).toHaveClass('ui-summary-card')
    expect(screen.getByText('2 项需检查')).toBeInTheDocument()
    expect(screen.getByText('1 个受影响工作流')).toBeInTheDocument()
    const summaryActions = within(screen.getByLabelText('数据质量下一步'))
    expect(summaryActions.getByRole('link', { name: '查看数据源设置' })).toHaveAttribute('href', '/settings')
    expect(summaryActions.getByRole('link', { name: '查看当前数据策略' })).toHaveAttribute('href', '/data-quality')
    expect(summaryActions.getByRole('link', { name: '查看受影响决策' })).toHaveAttribute('href', '/decisions/decision_degraded')
    expect(summaryActions.getByRole('link', { name: '查看证据' })).toHaveAttribute('href', '/evidence')
    expect(summaryActions.getByRole('link', { name: '查看风险预警' })).toHaveAttribute('href', '/risk-alerts')
    expect(screen.getByText('数据源健康')).toBeInTheDocument()
    expect(screen.getByText('csindex_extended · 指数估值文件 · 过期')).toBeInTheDocument()
    expect(screen.getByText('parse_source · 成分财务 · 解析失败')).toBeInTheDocument()
    expect(screen.getByText('失败类别：过期')).toBeInTheDocument()
    expect(screen.getByText('证据与检索')).toBeInTheDocument()
    expect(screen.getByText('证据数量：2')).toBeInTheDocument()
    expect(screen.getByText('独立信源：3')).toBeInTheDocument()
    expect(screen.getByText('VecLite：已配置')).toBeInTheDocument()
    expect(screen.getByText('LLM 质量')).toBeInTheDocument()
    expect(screen.getByText('DeepSeek：已配置')).toBeInTheDocument()
    expect(screen.getByText('质量门禁：复盘降级 1 条，缺证据 2 条')).toBeInTheDocument()
    expect(screen.getByText('影响范围与下一步')).toBeInTheDocument()
    expect(screen.getByText('知识与数据准备度')).toBeInTheDocument()
    expect(screen.getByText('沪深300ETF · ETF · 跟踪 000300')).toBeInTheDocument()
    expect(screen.getByText('格雷厄姆：安全边际')).toBeInTheDocument()
    expect(screen.getByText('估值分位 · 降级')).toBeInTheDocument()
    expect(screen.getAllByText('估值分位缺失时不得声明安全边际或估值高低。').length).toBeGreaterThanOrEqual(1)
    expect(screen.getByText('LLM 上下文：已附加知识与数据准备度摘要')).toBeInTheDocument()
    expect(screen.getByText('decision_degraded · 510300 · 降级')).toBeInTheDocument()
    const details = within(screen.getByLabelText('数据质量可观测区域'))
    expect(details.getByRole('link', { name: '查看设置' })).toHaveAttribute('href', '/settings')
    expect(details.getByRole('link', { name: '查看证据' })).toHaveAttribute('href', '/evidence')
    expect(details.getByRole('link', { name: '查看受影响决策' })).toHaveAttribute('href', '/decisions/decision_degraded')
    expect(details.getByRole('link', { name: '查看复盘' })).toHaveAttribute('href', '/review')
    expect(details.getByRole('link', { name: '查看审计' })).toHaveAttribute('href', '/audit')
    expect(details.getByRole('link', { name: '查看风险预警' })).toHaveAttribute('href', '/risk-alerts')
    expect(details.getByRole('link', { name: '返回工作台' })).toHaveAttribute('href', '/workbench')
  })

  it('shows ready knowledge and data readiness state', async () => {
    mockMinimalDataQualityDependencies()
    vi.mocked(getKnowledgeReadiness).mockResolvedValue({
      request_id: 'rid_knowledge_ready',
      data: knowledgeReadinessFixture({
        status: 'ready',
        data_dependencies: [
          { category: 'symbol_profile', status: 'ready', required: true, freshness: 'fresh', affected_features: ['consultation'], safe_degradation: '画像缺失时不生成正式交易类建议。' },
          { category: 'valuation_percentiles', status: 'ready', required: true, source_level: 'A', freshness: 'fresh', affected_features: ['margin_of_safety'], safe_degradation: '估值分位缺失时不得声明安全边际或估值高低。' },
        ],
        feature_impacts: [],
        llm_context_summary: 'principles=master.graham.margin_of_safety; data_readiness=valuation_percentiles=ready',
      }),
    })

    render(<MemoryRouter><DataQualityPage /></MemoryRouter>)

    const panel = await screen.findByLabelText('知识与数据准备度')
    expect(await within(panel).findByRole('heading', { name: '已准备' })).toBeInTheDocument()
    expect(within(panel).getByText('估值分位 · 已准备')).toBeInTheDocument()
    expect(within(panel).getByText('LLM 上下文：已附加知识与数据准备度摘要')).toBeInTheDocument()
  })

  it('loads data quality dependencies for the symbol in the URL query', async () => {
    mockMinimalDataQualityDependencies('159915')
    vi.mocked(getKnowledgeReadiness).mockResolvedValue({
      request_id: 'rid_knowledge_159915',
      data: knowledgeReadinessFixture({
        symbol: '159915',
        status: 'ready',
        symbol_profile: {
          symbol: '159915',
          name: '创业板ETF',
          asset_type: 'ETF',
          tracked_index_symbol: '399006',
          tracked_index_name: '创业板指',
          known: true,
        },
        data_dependencies: [
          { category: 'tracked_index', status: 'ready', required: true, freshness: 'fresh', request_id: 'req-readiness-159915', affected_symbols: ['399006'], affected_features: ['consultation'], safe_degradation: '跟踪指数缺失时 ETF 分析降级为信息不足。' },
          { category: 'valuation_percentiles', status: 'ready', required: true, source_level: 'A', freshness: 'fresh', request_id: 'req-readiness-159915', affected_symbols: ['399006'], affected_features: ['margin_of_safety'], safe_degradation: '估值分位缺失时不得声明安全边际或估值高低。' },
        ],
        feature_impacts: [],
        llm_context_summary: 'symbol_profile.159915;tracked_index.399006;valuation_percentiles=ready',
      }),
    })

    render(<MemoryRouter initialEntries={['/data-quality?symbol=159915']}><DataQualityPage /></MemoryRouter>)

    expect(await screen.findByText('当前查看：159915')).toBeInTheDocument()
    expect(screen.getByText('创业板ETF · ETF · 跟踪 399006')).toBeInTheDocument()
    expect(screen.getByText('跟踪指数 · 已准备')).toBeInTheDocument()
    expect(screen.getByText('估值分位 · 已准备')).toBeInTheDocument()
    await waitFor(() => expect(getLatestMarketSnapshot).toHaveBeenCalledWith('159915'))
    expect(getMarketSourceHealth).toHaveBeenCalledWith('159915')
    expect(getKnowledgeReadiness).toHaveBeenCalledWith('159915')
    expect(getDataSourceQualityRegression).toHaveBeenCalledWith('current', '159915')
    expect(getDataQualityGateResolution).toHaveBeenCalledWith('159915')
    expect(listDataQualityGateResolutions).toHaveBeenCalledWith('159915', 'active')
  })

  it('shows blocked knowledge readiness without fabricating symbol profile', async () => {
    mockMinimalDataQualityDependencies('999999')
    vi.mocked(getKnowledgeReadiness).mockResolvedValue({
      request_id: 'rid_knowledge_blocked',
      data: knowledgeReadinessFixture({
        symbol: '999999',
        status: 'blocked',
        symbol_profile: { symbol: '999999', known: false },
        data_dependencies: [
          { category: 'symbol_profile', status: 'blocked', required: true, freshness: 'missing', affected_features: ['consultation'], safe_degradation: '标的画像未知时不生成正式交易类建议。' },
        ],
        feature_impacts: [
          { feature: 'consultation', category: 'symbol_profile', impact: '标的画像未知时不生成正式交易类建议。', claims: ['不得伪造成 ready'] },
        ],
      }),
    })

    render(<MemoryRouter><DataQualityPage /></MemoryRouter>)

    const panel = await screen.findByLabelText('知识与数据准备度')
    expect(await within(panel).findByRole('heading', { name: '阻断' })).toBeInTheDocument()
    expect(within(panel).getByText('标的画像未准备，当前不能伪造成 ready。')).toBeInTheDocument()
    expect(within(panel).getByText('标的画像 · 阻断')).toBeInTheDocument()
    expect(within(panel).getAllByText(/不生成正式交易类建议/).length).toBeGreaterThanOrEqual(1)
  })

  it('shows degraded readiness impacts for valuation liquidity and formal evidence gaps', async () => {
    mockMinimalDataQualityDependencies()
    vi.mocked(getKnowledgeReadiness).mockResolvedValue({
      request_id: 'rid_knowledge_critical_gaps',
      data: knowledgeReadinessFixture({
        status: 'degraded',
        data_dependencies: [
          { category: 'valuation_percentiles', status: 'degraded', required: true, source_name: 'csindex_extended', source_level: 'A', source_type: 'index_valuation', freshness: 'parse_error', data_date: '2026-06-19', request_id: 'req-readiness-159915', affected_symbols: ['399006'], affected_features: ['margin_of_safety', 'expected_return', 'risk_alerts'], safe_degradation: '估值分位缺失时不得声明安全边际或估值高低，只能标记预期收益精度不足。' },
          { category: 'liquidity', status: 'degraded', required: true, source_level: 'B', freshness: 'missing', affected_features: ['risk_alerts', 'consultation'], safe_degradation: '流动性缺失时不得输出大额或市价式行动建议。' },
          { category: 'formal_evidence', status: 'degraded', required: true, source_level: 'A', freshness: 'insufficient', affected_features: ['consultation', 'decision_detail', 'risk_alerts'], safe_degradation: '正式证据不足时进入冻结观察或信息不足，不生成交易确认。' },
        ],
        feature_impacts: [
          { feature: 'margin_of_safety', category: 'valuation_percentiles', impact: '估值分位缺失时不得声明安全边际或估值高低，只能标记预期收益精度不足。', claims: ['不得伪造成 ready', '不得输出交易确认'] },
          { feature: 'risk_alerts', category: 'liquidity', impact: '流动性缺失时不得输出大额或市价式行动建议。', claims: ['不得伪造成 ready', '不得输出交易确认'] },
          { feature: 'consultation', category: 'formal_evidence', impact: '正式证据不足时进入冻结观察或信息不足，不生成交易确认。', claims: ['不得伪造成 ready', '不得输出交易确认'] },
        ],
      }),
    })

    render(<MemoryRouter><DataQualityPage /></MemoryRouter>)

    const readinessPanel = await screen.findByLabelText('知识与数据准备度')
    expect(await within(readinessPanel).findByRole('heading', { name: '降级' })).toBeInTheDocument()
    const dependencyPanel = screen.getByLabelText('知识与数据准备度').parentElement as HTMLElement
    expect(within(dependencyPanel).getByText('估值分位 · 降级')).toBeInTheDocument()
    expect(within(dependencyPanel).getByText('来源：csindex_extended · 类型：index_valuation · 日期：2026-06-19 · request：req-readiness-159915 · 标的：399006')).toBeInTheDocument()
    expect(within(dependencyPanel).getByText('流动性 · 降级')).toBeInTheDocument()
    expect(within(dependencyPanel).getByText('正式证据 · 降级')).toBeInTheDocument()
    expect(within(dependencyPanel).getAllByText(/不得声明安全边际/).length).toBeGreaterThanOrEqual(1)
    expect(within(dependencyPanel).getByText(/不得输出大额或市价式行动建议/)).toBeInTheDocument()
    expect(within(dependencyPanel).getAllByText(/不生成交易确认/).length).toBeGreaterThanOrEqual(1)
  })

  it('records and retires a current data gate scope exclusion from the UI', async () => {
    const activeResolution = {
      resolution_id: 'dqgr_1',
      symbol: '510300',
      policy_fingerprint: 'fp_test',
      policy_verdict: 'blocked',
      release_gate: 'block',
      policy_summary: 'data_source_quality:mode=current:status=degraded:policy=blocked:gate=block:no_auto_trading',
      resolution_type: 'scope_exclusion',
      status: 'active',
      scope: '本次 release clean claim 排除 current local data health',
      reason: '当前本地数据源存在降级，发布材料只声明有限范围',
      release_impact: '不得声明 current data healthy',
      evidence_ref: 'docs/release/acceptance',
      created_by: 'local_user',
      created_at: '2026-06-18T10:00:00Z',
      safety_note: '当前数据门禁处置只记录本地人工声明边界。',
    }
    vi.mocked(getSystemSettings).mockResolvedValue({ request_id: 'rid_system', data: { sqlite_status: 'ok', veclite_status: 'configured', deepseek_status: 'configured', data_sources: [], log_level: 'error' } })
    vi.mocked(getLatestMarketSnapshot).mockResolvedValue({ request_id: 'rid_market', data: { market_snapshot_id: 'market_1', symbol: '510300', pe_percentile: 0.5, pb_percentile: 0.5, liquidity_state: 'normal', sentiment_state: 'neutral', data_status: 'fresh' } })
    vi.mocked(getMarketSourceHealth).mockResolvedValue({ request_id: 'rid_source', data: { sources: [] } })
    vi.mocked(getEvidenceVerification).mockResolvedValue({ request_id: 'rid_verification', data: { verification_id: 'verify_1', verification_status: 'satisfied', independent_source_count: 1, high_grade_independent_source_count: 1, highest_source_level: 'A', latest_published_at: '2026-06-18T10:00:00Z', evidence_ids: [] } })
    vi.mocked(listEvidence).mockResolvedValue({ request_id: 'rid_evidence', data: { total: 0, items: [] } })
    vi.mocked(getReviewSummary).mockResolvedValue({ request_id: 'rid_review', data: { decision_count: 0, confirmation_count: 0, executed_manually_count: 0, planned_count: 0, error_case_count: 0, rule_proposal_count: 0, audit_event_count: 0, missing_evidence_count: 0, degraded_count: 0, ops_status: { data_source_status: 'degraded', index_status: 'success', review_status: 'success' }, recent_decisions: [] } })
    vi.mocked(getDataSourceQualityRegression).mockResolvedValue({ request_id: 'rid_policy', data: { mode: 'current', status: 'degraded', generated_at: '2026-06-18T10:00:00Z', summary: 'degraded', cases: [], missing_categories: [], policy: gateResolutionFixture().policy, safety_note: '只读检查' } })
    vi.mocked(getDataQualityGateResolution).mockResolvedValue({ request_id: 'rid_gate', data: gateResolutionFixture() })
    vi.mocked(listDataQualityGateResolutions)
      .mockResolvedValueOnce({ request_id: 'rid_resolutions_initial', data: { items: [], total: 0 } })
      .mockResolvedValueOnce({ request_id: 'rid_resolutions_active', data: { items: [activeResolution], total: 1 } })
      .mockResolvedValueOnce({ request_id: 'rid_resolutions_retired', data: { items: [], total: 0 } })
    vi.mocked(createDataQualityGateResolution).mockResolvedValue({
      request_id: 'rid_create',
      data: gateResolutionFixture({
        release_claim_state: 'resolved_with_scope_exclusion',
        active_resolution: activeResolution,
        allowed_claims: ['可以声明当前本地数据健康已排除在 clean claim 外'],
      }),
    })
    vi.mocked(retireDataQualityGateResolution).mockResolvedValue({ request_id: 'rid_retire', data: gateResolutionFixture() })

    render(<MemoryRouter><DataQualityPage /></MemoryRouter>)

    fireEvent.click(await screen.findByRole('button', { name: '记录处置' }))

    await waitFor(() => expect(createDataQualityGateResolution).toHaveBeenCalledWith(expect.objectContaining({
      symbol: '510300',
      resolution_type: 'scope_exclusion',
    })))
    expect(await screen.findByText('已排除 current data clean claim')).toBeInTheDocument()
    expect(screen.getByText('范围排除 · active')).toBeInTheDocument()

    fireEvent.click(screen.getByRole('button', { name: '退役处置' }))

    await waitFor(() => expect(retireDataQualityGateResolution).toHaveBeenCalledWith('dqgr_1'))
    expect(await screen.findByText('需要人工处置')).toBeInTheDocument()
  })

  it('resets resolution type when symbol check moves from waiver-required to blocked policy', async () => {
    const waiverPolicy = {
      verdict: 'waiver_required',
      release_gate: 'waiver_required',
      degraded_count: 1,
      failed_count: 0,
      blocking_count: 0,
      waiver_count: 1,
      blocking_reasons: [],
      waiver_reasons: ['sentiment_proxy optional category degraded freshness=stale'],
      next_actions: ['在发布材料中记录 waiver reason 和影响范围'],
      safety_note: '当前数据质量策略只读取本地 source health。',
    }
    vi.mocked(getSystemSettings).mockResolvedValue({ request_id: 'rid_system', data: { sqlite_status: 'ok', veclite_status: 'configured', deepseek_status: 'configured', data_sources: [], log_level: 'error' } })
    vi.mocked(getLatestMarketSnapshot).mockResolvedValue({ request_id: 'rid_market', data: { market_snapshot_id: 'market_waiver', symbol: '510300', pe_percentile: 0.5, pb_percentile: 0.5, liquidity_state: 'normal', sentiment_state: 'neutral', data_status: 'fresh' } })
    vi.mocked(getMarketSourceHealth).mockResolvedValue({ request_id: 'rid_source', data: { sources: [] } })
    vi.mocked(getEvidenceVerification).mockResolvedValue({ request_id: 'rid_verification', data: { verification_id: 'verify_1', verification_status: 'satisfied', independent_source_count: 1, high_grade_independent_source_count: 1, highest_source_level: 'A', latest_published_at: '2026-06-18T10:00:00Z', evidence_ids: [] } })
    vi.mocked(listEvidence).mockResolvedValue({ request_id: 'rid_evidence', data: { total: 0, items: [] } })
    vi.mocked(getReviewSummary).mockResolvedValue({ request_id: 'rid_review', data: { decision_count: 0, confirmation_count: 0, executed_manually_count: 0, planned_count: 0, error_case_count: 0, rule_proposal_count: 0, audit_event_count: 0, missing_evidence_count: 0, degraded_count: 0, ops_status: { data_source_status: 'degraded', index_status: 'success', review_status: 'success' }, recent_decisions: [] } })
    vi.mocked(getDataSourceQualityRegression).mockResolvedValue({ request_id: 'rid_policy', data: { mode: 'current', status: 'degraded', generated_at: '2026-06-18T10:00:00Z', summary: 'degraded', cases: [], missing_categories: [], policy: waiverPolicy, safety_note: '只读检查' } })
    vi.mocked(getDataQualityGateResolution)
      .mockResolvedValueOnce({ request_id: 'rid_gate_waiver', data: gateResolutionFixture({ policy: waiverPolicy, release_claim_state: 'requires_resolution' }) })
      .mockResolvedValueOnce({ request_id: 'rid_gate_blocked', data: gateResolutionFixture({ symbol: '000300' }) })
    vi.mocked(listDataQualityGateResolutions).mockResolvedValue({ request_id: 'rid_resolutions', data: { items: [], total: 0 } })
    vi.mocked(createDataQualityGateResolution).mockResolvedValue({
      request_id: 'rid_create',
      data: gateResolutionFixture({ symbol: '000300', release_claim_state: 'resolved_with_scope_exclusion' }),
    })

    render(<MemoryRouter><DataQualityPage /></MemoryRouter>)

    const typeSelect = await screen.findByLabelText('类型')
    expect(screen.getByRole('option', { name: '豁免记录' })).toBeInTheDocument()
    fireEvent.change(typeSelect, { target: { value: 'waiver' } })
    fireEvent.change(screen.getByLabelText('标的'), { target: { value: '000300' } })
    fireEvent.click(screen.getByRole('button', { name: '检查门禁处置' }))

    await waitFor(() => expect(getDataQualityGateResolution).toHaveBeenCalledWith('000300'))
    expect(screen.queryByRole('option', { name: '豁免记录' })).not.toBeInTheDocument()

    fireEvent.click(screen.getByRole('button', { name: '记录处置' }))

    await waitFor(() => expect(createDataQualityGateResolution).toHaveBeenCalledWith(expect.objectContaining({
      symbol: '000300',
      resolution_type: 'scope_exclusion',
    })))
  })

  it('shows degraded and unknown states safely without execution controls', async () => {
    vi.mocked(getSystemSettings).mockResolvedValue({
      request_id: 'rid_system_unknown',
      data: {
        sqlite_status: 'ok',
        veclite_status: 'missing',
        deepseek_status: 'unknown',
        data_sources: [],
        log_level: 'error',
      },
    })
    vi.mocked(getLatestMarketSnapshot).mockResolvedValue({
      request_id: 'rid_market_unknown',
      data: {
        market_snapshot_id: 'market_missing',
        symbol: '510300',
        pe_percentile: 0,
        pb_percentile: 0,
        liquidity_state: 'unknown',
        sentiment_state: 'unknown',
        data_status: 'missing',
      },
    })
    vi.mocked(getMarketSourceHealth).mockRejectedValue(new APIClientError({ requestId: 'rid_source', code: 'DATA_SOURCE_UNAVAILABLE', message: '数据源暂不可用，请检查数据源状态。', displayState: 'data_source_unavailable' }))
    vi.mocked(getEvidenceVerification).mockRejectedValue(new APIClientError({ requestId: 'rid_index', code: 'VECTOR_INDEX_UNAVAILABLE', message: '索引暂不可用，请稍后重试或重建索引。', displayState: 'insufficient_data' }))
    vi.mocked(listEvidence).mockResolvedValue({ request_id: 'rid_empty', data: { total: 0, items: [] } })
    vi.mocked(getReviewSummary).mockRejectedValue(new APIClientError({ requestId: 'rid_review', code: 'ANALYST_UNAVAILABLE', message: '分析服务暂不可用，页面仅展示规则与已有数据。', displayState: 'insufficient_data' }))
    vi.mocked(getDataSourceQualityRegression).mockRejectedValue(new APIClientError({ requestId: 'rid_policy', code: 'DATA_SOURCE_UNAVAILABLE', message: '当前数据策略暂不可用。', displayState: 'data_source_unavailable' }))
    vi.mocked(getDataQualityGateResolution).mockRejectedValue(new APIClientError({ requestId: 'rid_gate', code: 'DATA_SOURCE_UNAVAILABLE', message: '当前数据门禁处置暂不可用。', displayState: 'data_source_unavailable' }))
    vi.mocked(listDataQualityGateResolutions).mockResolvedValue({ request_id: 'rid_resolutions_empty', data: { items: [], total: 0 } })

    render(<MemoryRouter><DataQualityPage /></MemoryRouter>)

    expect(await screen.findByText('数据源暂不可用，请检查数据源状态。')).toBeInTheDocument()
    expect(screen.getByText('索引暂不可用，请稍后重试或重建索引。')).toBeInTheDocument()
    expect(screen.getByText('分析服务暂不可用，页面仅展示规则与已有数据。')).toBeInTheDocument()
    expect(screen.getByText('当前数据策略暂不可用。')).toBeInTheDocument()
    expect(screen.getByText('市场数据状态：缺失')).toBeInTheDocument()
    expect(screen.getByText('VecLite：缺失')).toBeInTheDocument()
    expect(screen.getByText('DeepSeek：未知状态')).toBeInTheDocument()
    expect(screen.getByText('暂无证据记录。')).toBeInTheDocument()
    expect(screen.getByText('暂无受影响工作流记录。')).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: /刷新|修复|重建|交易|下单|确认|应用|推送/ })).not.toBeInTheDocument()
  })

  it('does not expose secrets, private paths, SQL, or forbidden automatic action affordances', async () => {
    vi.mocked(getSystemSettings).mockResolvedValue({
      request_id: 'rid_system_safe',
      data: {
        sqlite_status: 'ok',
        sqlite_path: '/Users/private/investment-agent.db',
        veclite_status: 'configured',
        veclite_path: '/Users/private/veclite',
        deepseek_status: 'configured',
        data_sources: ['stub'],
        log_level: 'error',
      },
    })
    vi.mocked(getLatestMarketSnapshot).mockResolvedValue({
      request_id: 'rid_market_safe',
      data: {
        market_snapshot_id: 'market_safe',
        symbol: '510300',
        pe_percentile: 0,
        pb_percentile: 0,
        liquidity_state: 'normal',
        sentiment_state: 'neutral',
        data_status: 'fresh',
      },
    })
    vi.mocked(getMarketSourceHealth).mockResolvedValue({
      request_id: 'rid_source_safe',
      data: {
        sources: [{
          source_name: 'safe_source',
          source_level: 'A',
          source_type: 'public',
          data_category: 'index_valuation_files',
          freshness: 'fresh',
          failure_category: 'SQL failed near SELECT * FROM secrets',
          affected_symbols: ['510300'],
        }],
      },
    })
    vi.mocked(getEvidenceVerification).mockResolvedValue({
      request_id: 'rid_verification_safe',
      data: {
        verification_id: 'verify_safe',
        verification_status: 'satisfied',
        independent_source_count: 1,
        high_grade_independent_source_count: 1,
        highest_source_level: 'A',
        latest_published_at: '',
        evidence_ids: [],
      },
    })
    vi.mocked(listEvidence).mockResolvedValue({
      request_id: 'rid_evidence_safe',
      data: {
        total: 1,
        items: [{
          evidence_id: 'evidence_safe',
          source_name: 'CNInfo',
          source_level: 'A',
          evidence_role: 'formal',
          verification_status: 'satisfied',
          summary: 'prompt: sk-0277 secret SQL /Users/private/path should not render',
        }],
      },
    })
    vi.mocked(getReviewSummary).mockResolvedValue({
      request_id: 'rid_review_safe',
      data: {
        decision_count: 0,
        confirmation_count: 0,
        executed_manually_count: 0,
        planned_count: 0,
        error_case_count: 0,
        rule_proposal_count: 0,
        audit_event_count: 0,
        ops_status: {
          data_source_status: 'degraded',
          index_status: 'failed',
          review_status: 'degraded',
          explanation: 'DeepSeek 401 invalid_api_key; SQL syntax error near secrets; /home/vick/private; vendor raw response',
        },
        recent_decisions: [],
      },
    })
    vi.mocked(getDataSourceQualityRegression).mockResolvedValue({
      request_id: 'rid_policy_safe',
      data: {
        mode: 'current',
        status: 'passed',
        generated_at: '2026-06-18T10:00:00Z',
        summary: '数据源质量回归 mode=current status=passed cases=1 degraded=0 failed=0',
        cases: [],
        missing_categories: [],
        policy: {
          verdict: 'passed',
          release_gate: 'pass',
          degraded_count: 0,
          failed_count: 0,
          blocking_count: 0,
          waiver_count: 0,
          blocking_reasons: [],
          waiver_reasons: [],
          next_actions: ['保留当前只读数据质量证据'],
          safety_note: '当前数据质量策略只读取本地 source health。',
        },
        safety_note: '只读检查',
      },
    })
    vi.mocked(getDataQualityGateResolution).mockResolvedValue({
      request_id: 'rid_gate_safe',
      data: gateResolutionFixture({
        policy: {
          verdict: 'passed',
          release_gate: 'pass',
          degraded_count: 0,
          failed_count: 0,
          blocking_count: 0,
          waiver_count: 0,
          blocking_reasons: [],
          waiver_reasons: [],
          next_actions: ['保留当前只读数据质量证据'],
          safety_note: '当前数据质量策略只读取本地 source health。',
        },
        release_claim_state: 'pass',
        clean_data_claim_allowed: true,
        allowed_claims: ['可以声明当前本地数据门禁通过'],
        prohibited_claims: [],
      }),
    })
    vi.mocked(listDataQualityGateResolutions).mockResolvedValue({ request_id: 'rid_resolutions_safe', data: { items: [], total: 0 } })

    render(<MemoryRouter><DataQualityPage /></MemoryRouter>)

    await screen.findByText('safe_source · 指数估值文件 · 新鲜')
    const text = document.body.textContent ?? ''
    expect(text).not.toMatch(/sk-0277|secret|SELECT \* FROM|\/Users\/private|prompt:|invalid_api_key|SQL syntax error|\/home\/vick|vendor raw response/)
    expect(screen.getByText('存在已脱敏诊断摘要。')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '检查门禁处置' })).toBeInTheDocument()
    expect(screen.queryByRole('link', { name: /自动交易|一键交易|代下单|券商接口|外部推送|自动确认|自动应用规则|收益承诺/ })).not.toBeInTheDocument()
  })

  it('renders nullable gate claim arrays without crashing', async () => {
    vi.mocked(getSystemSettings).mockResolvedValue({ request_id: 'rid_system_nullable', data: { sqlite_status: 'ok', veclite_status: 'configured', deepseek_status: 'configured', data_sources: ['csindex'], log_level: 'error' } })
    vi.mocked(getLatestMarketSnapshot).mockResolvedValue({ request_id: 'rid_market_nullable', data: { market_snapshot_id: 'market_nullable', symbol: '000300', trade_date: '2026-06-18', pe_percentile: 0.5, pb_percentile: 0.5, liquidity_state: 'normal', sentiment_state: 'neutral', data_status: 'fresh' } })
    vi.mocked(getMarketSourceHealth).mockResolvedValue({ request_id: 'rid_source_nullable', data: { sources: [] } })
    vi.mocked(listEvidence).mockResolvedValue({ request_id: 'rid_evidence_nullable', data: { items: [], total: 0 } })
    vi.mocked(getEvidenceVerification).mockResolvedValue({ request_id: 'rid_verification_nullable', data: { verification_id: 'verification_nullable', verification_status: 'background_only', independent_source_count: 0, high_grade_independent_source_count: 0, highest_source_level: 'C', latest_published_at: '', evidence_ids: [] } })
    vi.mocked(getReviewSummary).mockResolvedValue({ request_id: 'rid_review_nullable', data: { decision_count: 0, confirmation_count: 0, executed_manually_count: 0, planned_count: 0, error_case_count: 0, rule_proposal_count: 0, audit_event_count: 0, ops_status: { data_source_status: 'success', index_status: 'success', review_status: 'success', explanation: 'ok' }, recent_decisions: [] } })
    vi.mocked(getDataSourceQualityRegression).mockResolvedValue({ request_id: 'rid_regression_nullable', data: { mode: 'current', status: 'passed', generated_at: '2026-06-18T00:00:00Z', summary: 'passed', cases: [], missing_categories: [], policy: { verdict: 'passed', release_gate: 'pass', degraded_count: 0, failed_count: 0, blocking_count: 0, waiver_count: 0, blocking_reasons: null, waiver_reasons: null, next_actions: null, safety_note: '只读。' }, safety_note: '只读。' } })
    vi.mocked(getDataQualityGateResolution).mockResolvedValue({ request_id: 'rid_gate_nullable', data: gateResolutionFixture({ policy: { verdict: 'passed', release_gate: 'pass', degraded_count: 0, failed_count: 0, blocking_count: 0, waiver_count: 0, blocking_reasons: null, waiver_reasons: null, next_actions: null, safety_note: '只读。' }, release_claim_state: 'pass', clean_data_claim_allowed: true, allowed_claims: null, prohibited_claims: null }) })
    vi.mocked(listDataQualityGateResolutions).mockResolvedValue({ request_id: 'rid_resolutions_nullable', data: { items: [], total: 0 } })

    render(<MemoryRouter><DataQualityPage /></MemoryRouter>)

    await screen.findByText('当前数据策略：通过；release gate：通过')
    expect(screen.getByText('clean data claim：允许')).toBeInTheDocument()
  })
})
