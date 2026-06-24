import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import type { DecisionDetailResponse } from '../../types/decision'
import { DecisionTrace } from './DecisionTrace'

const decision: DecisionDetailResponse = {
  decision_id: 'decision_1',
  question: '可以买入吗',
  symbol: '510300',
  generated_at: '2026-06-01T00:00:00Z',
  capability_check: { status: 'new_capability_status', reason: '未配置能力圈' },
  workflow_status: 'completed',
  account_snapshot: { snapshot_id: 'snap_1', cash: 100, total_assets: 1000, cash_ratio: 0.1, high_risk_ratio: 0.2 },
  triggered_rules: [{ rule_id: 'rule_1', rule_name: '估值纪律', severity: 'new_severity', description: '估值偏高' }],
  evidence_chain: [],
  analyst_reports: [{ agent_name: 'analyst', conclusion: '观察', key_reasons: [], risk_warnings: [], confidence: 'new_confidence', evidence_ids: [] }],
  expected_return_scenarios: {
    sample_count: 3,
    precision_status: 'available',
    scenarios: [{ scenario: 'new_scenario', return_range: '-2% ~ 3%', probability: 0.5 }],
    disclaimer: '仅作参考',
  },
  arbitration_chain: [],
  final_verdict: { status: 'hold', display_text: '继续观察', prohibited_actions: [], optional_actions: [] },
  user_confirmation: { confirmation_status: 'not_required', available_actions: [] },
}

const positiveDecision: DecisionDetailResponse = {
  ...decision,
  capability_check: { status: 'in_scope', reason: '沪深300 ETF' },
  triggered_rules: [{ rule_id: 'rule_ok', rule_name: '纪律通过', severity: 'warning', description: '需要观察' }],
  analyst_reports: [{ agent_name: 'analyst', conclusion: '可以观察', key_reasons: ['估值合理'], risk_warnings: ['波动风险'], confidence: 'high', evidence_ids: ['ev_1'] }],
  expected_return_scenarios: {
    sample_count: 8,
    sample_window: '2024-01-01 至 2026-06-01',
    screening_condition: '同类 ETF 历史样本',
    precision_status: 'available',
    scenarios: [{ scenario: 'base', return_range: '0% ~ 5%', probability: 0.6, trigger: '价格回到基准区间' }],
    disclaimer: '历史样本仅作参考，不构成收益承诺。',
    sell_evaluation: {
      status: 'review_required',
      triggers: ['触及上行情景下沿'],
      prompts: ['人工复核止盈计划'],
      actions: ['记录人工计划'],
      non_trading_disclaimer: '卖出评估仅用于人工复核，不会自动执行交易。',
    },
    reassessment_trigger: {
      reason: '基准中位数下移超过 15%',
      boundary: 'base_midpoint_downshift_gt_15pct',
      current_value: 0.16,
    },
  },
}

describe('DecisionTrace', () => {
  afterEach(() => cleanup())

  it('shows a story-first verdict, safety and trust path before technical trace', () => {
    render(<DecisionTrace decision={positiveDecision} />)

    expect(screen.getByText('决策故事')).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: '继续观察' })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: '安全边界' })).toBeInTheDocument()
    expect(screen.getByText(/只展示本地分析和人工复核路径，不会自动交易、自动确认或自动应用规则/)).toBeInTheDocument()
    expect(screen.getAllByText(/估值合理/).length).toBeGreaterThanOrEqual(1)
    expect(screen.getByRole('heading', { name: '可信度' })).toBeInTheDocument()
    expect(screen.getByText(/正式证据 0 条，背景材料 0 条/)).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '查看证据' })).toHaveAttribute('href', '/evidence')
    expect(screen.getByRole('link', { name: '查看决策闭环' })).toHaveAttribute('href', '/decision-loop')
    expect(screen.getByRole('link', { name: '查看审计' })).toHaveAttribute('href', '/audit')
    expect(screen.getByRole('heading', { name: '技术追踪' })).toBeInTheDocument()
  })

  it('does not emit duplicate key warnings for repeated evidence summaries', () => {
    const consoleError = vi.spyOn(console, 'error').mockImplementation(() => undefined)
    render(<DecisionTrace decision={{
      ...positiveDecision,
      evidence_chain: [
        { evidence_id: 'ev_1', source_name: 'stub', source_level: 'A', evidence_role: 'formal', summary: '重复摘要' },
        { evidence_id: 'ev_2', source_name: 'stub', source_level: 'A', evidence_role: 'background', summary: '重复摘要' },
      ],
    }} />)

    expect(consoleError.mock.calls.some((call) => call.join(' ').includes('Encountered two children with the same key'))).toBe(false)
    consoleError.mockRestore()
  })

  it('maps decision enum fields to safe text without exposing unknown raw values', () => {
    render(<DecisionTrace decision={decision} />)
    fireEvent.click(screen.getByRole('button', { name: '展开 1 份分析材料' }))

    expect(screen.getAllByText(/未知状态/).length).toBeGreaterThanOrEqual(3)
    expect(screen.getByText(/未知情景/)).toBeInTheDocument()
    expect(screen.queryByText(/new_capability_status|new_severity|new_confidence|new_scenario/)).not.toBeInTheDocument()
  })

  it('shows mapped confidence scenario and disclaimer for normal values', () => {
    render(<DecisionTrace decision={positiveDecision} />)
    fireEvent.click(screen.getByRole('button', { name: '展开 1 份分析材料' }))

    expect(screen.getByText(/能力圈检查：能力圈内/)).toBeInTheDocument()
    expect(screen.getByText(/纪律通过：预警/)).toBeInTheDocument()
    expect(screen.getByText(/置信度：高/)).toBeInTheDocument()
    expect(screen.getByText(/基准情景：0% ~ 5%/)).toBeInTheDocument()
    expect(screen.getByText('历史样本仅作参考，不构成收益承诺。')).toBeInTheDocument()
  })

  it('reads back take-profit funds returning to core assets as a manual optional action', () => {
    render(<DecisionTrace decision={{
      ...positiveDecision,
      final_verdict: {
        ...positiveDecision.final_verdict,
        optional_actions: ['卖出 30%', '止盈资金优先回归核心资产'],
      },
    }} />)

    expect(screen.getByText(/可选动作：卖出 30%、止盈资金优先回归核心资产/)).toBeInTheDocument()
  })

  it('shows extreme-fear historical similar scenario context', () => {
    render(<DecisionTrace decision={{
      ...positiveDecision,
      expected_return_scenarios: {
        ...positiveDecision.expected_return_scenarios!,
        historical_contexts: [
          {
            label: '极端恐惧样本',
            window: '2018Q4, 2020Q1, 2022Q4',
            sample_count: 20,
            outcome: '暂停主动交易建议',
            max_drawdown: -0.18,
            recovery: '3-9 个月',
            source: 'local_public_history',
          },
        ],
        sell_evaluation: {
          status: 'triggered',
          triggers: ['extreme_fear_historical_context'],
          prompts: ['极端恐惧状态已展示历史相似场景'],
          actions: ['暂停主动交易建议'],
          non_trading_disclaimer: '卖出评估仅用于人工复核，不会自动执行交易。',
        },
      },
    }} />)

    expect(screen.getByRole('heading', { name: '历史相似场景' })).toBeInTheDocument()
    expect(screen.getByText(/极端恐惧样本：2018Q4, 2020Q1, 2022Q4/)).toBeInTheDocument()
    expect(screen.getByText(/最大回撤 -18.0%/)).toBeInTheDocument()
    expect(screen.getAllByText(/暂停主动交易建议/).length).toBeGreaterThanOrEqual(1)
  })

  it('shows knowledge readiness context metadata without exposing full prompts', () => {
    render(<DecisionTrace decision={{
      ...positiveDecision,
      analyst_reports: [{
        agent_name: 'value',
        conclusion: '估值证据不足，保持观察。',
        key_reasons: ['安全边际不足'],
        risk_warnings: ['估值分位降级'],
        confidence: 'medium',
        evidence_ids: ['ev_1'],
        input_summary: 'value 510300 principles=master.graham.margin_of_safety data_readiness=valuation_percentiles=degraded boundary=背景知识不能满足正式证据',
        prompt_version: 'p37-analyst-v1',
        quality_status: 'passed',
      }],
    }} />)
    fireEvent.click(screen.getByRole('button', { name: '展开 1 份分析材料' }))

    expect(screen.getByText('LLM 已参考知识与数据准备度摘要')).toBeInTheDocument()
    expect(screen.getByText(/prompt p37-analyst-v1/)).toBeInTheDocument()
    expect(document.body.textContent).not.toMatch(/完整 prompt|持仓上下文：|principles=|data_readiness=/)
  })

  it('shows P28 expected return sample context sell evaluation and reassessment trigger', () => {
    render(<DecisionTrace decision={{
      ...positiveDecision,
      market_context: {
        symbol: '510300',
        trade_date: '2026-06-20',
        current_price: 4.23,
        pe_percentile: 31,
        pb_percentile: 27,
      },
    }} />)

    expect(screen.getAllByText(/标的：510300/).length).toBeGreaterThanOrEqual(1)
    expect(screen.getByText(/当前日期：2026-06-20/)).toBeInTheDocument()
    expect(screen.getByText(/当前价格或净值：4.23/)).toBeInTheDocument()
    expect(screen.getByText(/PE\/PB 分位：31 \/ 27/)).toBeInTheDocument()
    expect(screen.getByText(/样本窗口：2024-01-01 至 2026-06-01/)).toBeInTheDocument()
    expect(screen.getByText(/筛选条件：同类 ETF 历史样本/)).toBeInTheDocument()
    expect(screen.getByText(/触发条件：价格回到基准区间/)).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: '动态卖出评估' })).toBeInTheDocument()
    expect(screen.getByText(/状态：需人工复核/)).toBeInTheDocument()
    expect(screen.getByText(/触发因素：触及上行情景下沿/)).toBeInTheDocument()
    expect(screen.getByText(/人工提示：人工复核止盈计划/)).toBeInTheDocument()
    expect(screen.getByText(/建议动作：记录人工计划/)).toBeInTheDocument()
    expect(screen.getByText('卖出评估仅用于人工复核，不会自动执行交易。')).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: '复核触发' })).toBeInTheDocument()
    expect(screen.getByText(/原因：基准中位数下移超过 15%/)).toBeInTheDocument()
    expect(screen.getByText(/边界：base_midpoint_downshift_gt_15pct/)).toBeInTheDocument()
    expect(screen.getByText(/当前值：16.00%/)).toBeInTheDocument()
  })

  it('shows degraded retrieval quality without exposing local index controls', () => {
    const degradedDecision = {
      ...positiveDecision,
      retrieval_quality: {
        query_summary: '510300 是否买入',
        top_k: 1,
        status: 'degraded',
        index_health: 'missing',
        index_freshness: 'unknown',
        fallback_source: 'sqlite_summary',
        source_consistency_status: 'checked',
        degraded_reason: 'veclite index not configured',
      },
    } as DecisionDetailResponse

    render(<DecisionTrace decision={degradedDecision} />)

    expect(screen.getByRole('heading', { name: '检索质量' })).toBeInTheDocument()
    expect(screen.getByText(/检索状态：降级/)).toBeInTheDocument()
    expect(screen.getByText(/召回数量：1/)).toBeInTheDocument()
    expect(screen.getByText(/降级原因：veclite index not configured/)).toBeInTheDocument()
    expect(screen.getByText(/可在证据页重建索引后再次复核/)).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: /重建/ })).not.toBeInTheDocument()
  })

  it('shows healthy retrieval quality as a read-only summary', () => {
    render(<DecisionTrace decision={{
      ...positiveDecision,
      retrieval_quality: {
        query_summary: '510300 公告',
        top_k: 3,
        status: 'hit',
        index_health: 'healthy',
        index_freshness: 'fresh',
        fallback_source: 'veclite',
        source_consistency_status: 'checked',
      },
    }} />)

    expect(screen.getByRole('heading', { name: '检索质量' })).toBeInTheDocument()
    expect(screen.getByText(/检索状态：命中/)).toBeInTheDocument()
    expect(screen.getByText(/索引健康：健康/)).toBeInTheDocument()
    expect(screen.getByText(/Fallback 来源：VecLite 索引/)).toBeInTheDocument()
    expect(screen.queryByText(/重建索引/)).not.toBeInTheDocument()
  })

  it('shows empty retrieval quality state without pretending evidence was found', () => {
    render(<DecisionTrace decision={{
      ...positiveDecision,
      retrieval_quality: {
        query_summary: '510300 缺失主题',
        top_k: 0,
        status: 'empty',
        index_health: 'fresh',
        index_freshness: 'fresh',
        fallback_source: 'none',
        source_consistency_status: 'not_checked',
      },
    }} />)

    expect(screen.getByText(/检索状态：无结果/)).toBeInTheDocument()
    expect(screen.getByText(/召回数量：0/)).toBeInTheDocument()
    expect(screen.getByText(/未召回可用于裁决的证据/)).toBeInTheDocument()
  })

  it('renders real LLM-like nullable verdict lists without crashing', () => {
    const realLikeDecision: DecisionDetailResponse = {
      ...positiveDecision,
      decision_id: 'decision_real_llm_nullable',
      final_verdict: {
        status: 'hold',
        display_text: '继续持有，等待人工复核',
        prohibited_actions: null,
        optional_actions: null,
      },
      analyst_reports: [{
        agent_name: 'value',
        conclusion: '估值证据不足，保持观察。',
        key_reasons: null,
        risk_warnings: null,
        confidence: 'medium',
        evidence_ids: null,
      } as unknown as DecisionDetailResponse['analyst_reports'][number]],
    }

    render(<DecisionTrace decision={realLikeDecision} />)
    fireEvent.click(screen.getByRole('button', { name: '展开 1 份分析材料' }))

    expect(screen.getByRole('heading', { name: '继续持有，等待人工复核' })).toBeInTheDocument()
    expect(screen.getByText(/建议编号：decision_real_llm_nullable/)).toBeInTheDocument()
    expect(screen.getByText(/禁止事项：暂无/)).toBeInTheDocument()
    expect(screen.getByText(/可选动作：暂无/)).toBeInTheDocument()
    expect(screen.getByText(/关键理由：暂无/)).toBeInTheDocument()
    expect(screen.getByText(/风险提示：暂无/)).toBeInTheDocument()
    expect(screen.getByText(/证据引用：暂无/)).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: /自动交易|一键交易|代下单/ })).not.toBeInTheDocument()
    expect(screen.queryByRole('link', { name: /自动交易|一键交易|代下单/ })).not.toBeInTheDocument()
  })

  it('renders missing verdict list fields as safe empty states', () => {
    const missingListDecision: DecisionDetailResponse = {
      ...positiveDecision,
      decision_id: 'decision_real_llm_missing_lists',
      final_verdict: {
        status: 'hold',
        display_text: '继续观察',
      } as DecisionDetailResponse['final_verdict'],
    }

    render(<DecisionTrace decision={missingListDecision} />)

    expect(screen.getByText(/建议编号：decision_real_llm_missing_lists/)).toBeInTheDocument()
    expect(screen.getByText(/禁止事项：暂无/)).toBeInTheDocument()
    expect(screen.getByText(/可选动作：暂无/)).toBeInTheDocument()
  })

  it('renders sparse real LLM-like decision substructures as safe empty states', () => {
    const sparseDecision = {
      ...positiveDecision,
      decision_id: 'decision_sparse_real_llm',
      final_verdict: null,
      triggered_rules: null,
      evidence_chain: null,
      analyst_reports: null,
      expected_return_scenarios: {
        sample_count: 0,
        precision_status: 'unavailable',
        scenarios: null,
        disclaimer: '',
      },
      arbitration_chain: null,
      user_confirmation: null,
    } as unknown as DecisionDetailResponse

    render(<DecisionTrace decision={sparseDecision} />)

    expect(screen.getByText(/建议编号：decision_sparse_real_llm/)).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: '暂无最终裁决' })).toBeInTheDocument()
    expect(screen.getByText(/禁止事项：暂无/)).toBeInTheDocument()
    expect(screen.getByText(/可选动作：暂无/)).toBeInTheDocument()
    expect(screen.getByText(/暂无触发规则。/)).toBeInTheDocument()
    expect(screen.getByText(/暂无证据链。/)).toBeInTheDocument()
    expect(screen.getByText(/暂无分析材料。/)).toBeInTheDocument()
    expect(screen.getByText(/暂无裁决链。/)).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: /自动交易|一键交易|代下单/ })).not.toBeInTheDocument()
  })
})
