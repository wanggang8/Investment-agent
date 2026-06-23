import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { cleanup } from '@testing-library/react'
import { APIClientError } from '../services/client'
import { SettingsPage } from './SettingsPage'

vi.mock('../services/settings', () => ({
  getCapabilitySettings: vi.fn(),
  getSystemSettings: vi.fn(),
}))

vi.mock('../services/market', () => ({
  getLatestMarketSnapshot: vi.fn(),
  getMarketSourceHealth: vi.fn(),
  refreshMarket: vi.fn(),
}))

import { getLatestMarketSnapshot, getMarketSourceHealth, refreshMarket } from '../services/market'
import { getCapabilitySettings, getSystemSettings } from '../services/settings'

describe('SettingsPage', () => {
  afterEach(() => {
    cleanup()
    vi.clearAllMocks()
  })

  it('maps system and market statuses to safe Chinese text', async () => {
    vi.mocked(getCapabilitySettings).mockResolvedValue({ request_id: 'rid', data: { asset_types: ['ETF'], symbols: ['510300'], excluded_symbols: [], strategy_scope: ['定投'] } })
    vi.mocked(getSystemSettings).mockResolvedValue({ request_id: 'rid', data: { sqlite_status: 'rebuilding', veclite_status: 'degraded', deepseek_status: 'failed', data_sources: ['sqlite'], log_level: 'info' } })
    vi.mocked(getLatestMarketSnapshot).mockResolvedValue({ request_id: 'rid', data: { market_snapshot_id: 'market_1', symbol: '510300', trade_date: '2026-01-01', pe_percentile: 20, pb_percentile: 30, data_status: 'fresh', sentiment_state: 'neutral', liquidity_state: 'danger' } })
    vi.mocked(getMarketSourceHealth).mockResolvedValue({ request_id: 'rid', data: { sources: [] } })

    render(<SettingsPage />)

    await waitFor(() => expect(screen.getByText('本地配置与诊断状态')).toBeInTheDocument())
    expect(screen.getByText('本地配置与诊断需要检查')).toBeInTheDocument()
    expect(screen.getByText('查看数据源健康')).toBeInTheDocument()
    expect(screen.getByText('复验本地安装')).toBeInTheDocument()
    await waitFor(() => expect(screen.getByText('SQLite：重建中')).toBeInTheDocument())
    expect(screen.getByText('VecLite 索引状态：降级')).toBeInTheDocument()
    expect(screen.getByText('DeepSeek：失败')).toBeInTheDocument()
    expect(screen.getByText('P40 本地运行就绪')).toBeInTheDocument()
    expect(screen.getByText('就绪摘要：SQLite 重建中；VecLite 降级；DeepSeek 失败')).toBeInTheDocument()
    expect(screen.getByText('预检入口：go run ./cmd/agent --preflight --diagnostics ./tmp/preflight.json')).toBeInTheDocument()
    expect(screen.getByText('只展示本地诊断和人工处理提示；不发起资金动作、站外通知或规则生效。')).toBeInTheDocument()
    expect(screen.getByText('数据状态：新鲜')).toBeInTheDocument()
    expect(screen.getByText('情绪/流动性：中性 / 危险')).toBeInTheDocument()
    expect(screen.queryByText(/configured|fresh|neutral|danger/)).not.toBeInTheDocument()
  })

  it('shows P40 source health without trading actions', async () => {
    vi.mocked(getCapabilitySettings).mockResolvedValue({ request_id: 'rid', data: { asset_types: ['ETF'], symbols: ['000300'], excluded_symbols: [], strategy_scope: ['定投'] } })
    vi.mocked(getSystemSettings).mockResolvedValue({ request_id: 'rid', data: { sqlite_status: 'ok', veclite_status: 'ok', deepseek_status: 'configured', data_sources: ['csindex'], log_level: 'info' } })
    vi.mocked(getLatestMarketSnapshot).mockResolvedValue({ request_id: 'rid', data: { market_snapshot_id: 'market_1', symbol: '000300', trade_date: '2026-06-05', pe_percentile: 0, pb_percentile: 0, data_status: 'fresh', sentiment_state: 'neutral', liquidity_state: 'normal' } })
    vi.mocked(getMarketSourceHealth).mockResolvedValue({
      request_id: 'rid',
      data: {
        sources: [
          { source_name: 'csindex', source_level: 'A', source_type: 'index_basic', data_category: 'index_constituents', freshness: 'fresh', data_date: '2026-06-05', last_success_at: '2026-06-06T01:00:00Z', affected_symbols: ['000300'] },
          { source_name: 'csindex', source_level: 'A', source_type: 'index_basic', data_category: 'index_valuation_files', freshness: 'parse_error', data_date: '2026-06-05', failure_category: 'parse_error', last_failure_at: '2026-06-06T01:00:00Z', affected_symbols: ['000300'] },
          { source_name: 'eastmoney', source_level: 'B', source_type: 'fund_basic', data_category: 'capital_flow', freshness: 'stale', data_date: '2026-06-04', failure_category: 'failed', last_failure_at: '2026-06-06T02:00:00Z', affected_symbols: ['510300'] },
          { source_name: 'local-cache', source_level: 'C', source_type: 'fallback', data_category: 'sentiment_proxy', freshness: 'missing', data_date: '2026-06-03', failure_category: 'unknown', affected_symbols: ['510300'] },
          { source_name: 'snapshot-api', source_level: 'B', source_type: 'market_snapshot', data_category: 'margin_financing', freshness: 'failed', data_date: '2026-06-02', affected_symbols: ['510500'] },
          { source_name: 'manual-review', source_level: 'C', source_type: 'review', data_category: 'constituent_financials', freshness: 'unknown', data_date: '2026-06-01', affected_symbols: ['510500'] },
        ],
      },
    })

    render(<SettingsPage />)

    await waitFor(() => expect(screen.getByText('P40 数据源健康')).toBeInTheDocument())
    expect(screen.getByText('csindex · 指数样本 · 新鲜')).toBeInTheDocument()
    expect(screen.getByText('csindex · 指数估值文件 · 解析失败')).toBeInTheDocument()
    expect(screen.getByText('eastmoney · 资金流向 · 过期')).toBeInTheDocument()
    expect(screen.getByText('local-cache · 情绪替代指标 · 缺失')).toBeInTheDocument()
    expect(screen.getByText('snapshot-api · 融资融券 · 失败')).toBeInTheDocument()
    expect(screen.getByText('manual-review · 成分财务 · 未知状态')).toBeInTheDocument()
    expect(screen.getAllByText(/数据日：2026-06-05/).length).toBeGreaterThan(0)
    expect(screen.getByText(/最近成功：2026-06-06T01:00:00Z/)).toBeInTheDocument()
    expect(screen.getByText(/最近失败：2026-06-06T01:00:00Z/)).toBeInTheDocument()
    expect(screen.getByText(/失败类别：解析失败/)).toBeInTheDocument()
    expect(screen.getByText(/失败类别：失败/)).toBeInTheDocument()
    expect(screen.getByText(/失败类别：未知状态/)).toBeInTheDocument()
    expect(screen.getByText('仅展示公开只读数据状态；不会连接券商或发起交易。')).toBeInTheDocument()
    expect(screen.queryByText(/一键交易|自动下单/)).not.toBeInTheDocument()
  })

  it('shows P90 capital-flow structured fields with raw net flow', async () => {
    vi.mocked(getCapabilitySettings).mockResolvedValue({ request_id: 'rid', data: { asset_types: ['stock'], symbols: ['600000'], excluded_symbols: [], strategy_scope: ['hold_review'] } })
    vi.mocked(getSystemSettings).mockResolvedValue({ request_id: 'rid', data: { sqlite_status: 'ok', veclite_status: 'ok', deepseek_status: 'configured', data_sources: ['p89_structured_public'], log_level: 'info' } })
    vi.mocked(getLatestMarketSnapshot).mockResolvedValue({
      request_id: 'rid',
      data: {
        market_snapshot_id: 'market_p90',
        symbol: '600000',
        trade_date: '2026-06-22',
        pe_percentile: 20,
        pb_percentile: 30,
        data_status: 'fresh',
        sentiment_state: 'neutral',
        liquidity_state: 'normal',
        market_metrics: {
          metadata: {
            p88_structured_fields: {
              capital_flow: {
                date: '2026-06-22',
                net_inflow: 11895999,
                net_outflow: 0,
                raw_net_flow: 11895999,
              },
            },
          },
        },
      },
    })
    vi.mocked(getMarketSourceHealth).mockResolvedValue({ request_id: 'rid', data: { sources: [] } })

    render(<SettingsPage />)

    await waitFor(() => expect(screen.getByText('结构化字段')).toBeInTheDocument())
    expect(screen.getByText('资金流向：2026-06-22；净流入 11,895,999；净流出 暂无；净流向 11,895,999')).toBeInTheDocument()
  })

  it('shows API error notices instead of silently treating failed settings as empty', async () => {
    vi.mocked(getCapabilitySettings).mockRejectedValue(new APIClientError({ requestId: 'cap', code: 'INTERNAL_ERROR', message: '能力圈配置读取失败。', displayState: 'generic_failure' }))
    vi.mocked(getSystemSettings).mockRejectedValue(new APIClientError({ requestId: 'sys', code: 'INTERNAL_ERROR', message: '系统状态读取失败。', displayState: 'generic_failure' }))
    vi.mocked(getLatestMarketSnapshot).mockRejectedValue(new APIClientError({ requestId: 'market', code: 'DATA_STALE', message: '市场快照读取失败。', displayState: 'insufficient_data' }))
    vi.mocked(getMarketSourceHealth).mockRejectedValue(new APIClientError({ requestId: 'source', code: 'DATA_SOURCE_UNAVAILABLE', message: '数据源健康读取失败。', displayState: 'insufficient_data' }))

    render(<SettingsPage />)

    await waitFor(() => expect(screen.getByText('能力圈配置读取失败。')).toBeInTheDocument())
    expect(screen.getByText('系统状态读取失败。')).toBeInTheDocument()
    expect(screen.getByText('市场快照读取失败。')).toBeInTheDocument()
    expect(screen.getByText('数据源健康读取失败。')).toBeInTheDocument()
  })

  it('allows safe market refresh without trading copy', async () => {
    vi.mocked(getCapabilitySettings).mockResolvedValue({ request_id: 'rid', data: { asset_types: ['ETF'], symbols: ['510300'], excluded_symbols: [], strategy_scope: ['定投'] } })
    vi.mocked(getSystemSettings).mockResolvedValue({ request_id: 'rid', data: { sqlite_status: 'ok', veclite_status: 'ok', deepseek_status: 'configured', data_sources: ['sqlite'], log_level: 'info' } })
    vi.mocked(getLatestMarketSnapshot).mockResolvedValue({ request_id: 'rid', data: { market_snapshot_id: 'market_1', symbol: '510300', trade_date: '2026-01-01', pe_percentile: 20, pb_percentile: 30, data_status: 'fresh', sentiment_state: 'neutral', liquidity_state: 'normal' } })
    vi.mocked(getMarketSourceHealth).mockResolvedValue({ request_id: 'rid', data: { sources: [] } })
    vi.mocked(refreshMarket).mockResolvedValue({ request_id: 'refresh', data: { refreshed_count: 1, failed_symbols: [], latest_snapshot_ids: ['market_2'], audit_event_ids: ['audit_1'] } })

    render(<SettingsPage />)
    fireEvent.click(await screen.findByRole('button', { name: '刷新市场数据' }))

    expect(refreshMarket).toHaveBeenCalledWith({ symbols: ['510300'] })
    await waitFor(() => expect(screen.getByText('市场刷新完成；只更新本地行情事实和审计记录，不会执行交易。')).toBeInTheDocument())
    expect(screen.queryByText(/一键交易|自动下单/)).not.toBeInTheDocument()
  })
})
