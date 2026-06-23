import type { CapabilitySettings, SystemStatus } from '../../types/settings'
import type { MarketSnapshot, SourceHealthItem } from '../../types/market'
import { marketStateText, sourceCategoryText, sourceHealthStatusText, systemStatusText, textOrRaw } from '../../shared/mappers'

interface Props {
  capability?: CapabilitySettings
  system?: SystemStatus
  marketSnapshot?: MarketSnapshot
  sourceHealth?: SourceHealthItem[]
}

export function CapabilitySettingsPanel({ capability, system, marketSnapshot, sourceHealth = [] }: Props) {
  const readinessSummary = [
    `SQLite ${textOrRaw(systemStatusText, system?.sqlite_status)}`,
    `VecLite ${textOrRaw(systemStatusText, system?.veclite_status)}`,
    `DeepSeek ${textOrRaw(systemStatusText, system?.deepseek_status)}`,
  ].join('；')
  const structuredFields = readStructuredFields(marketSnapshot)
  const margin = structuredFields.margin_financing
  const financial = structuredFields.constituent_financial
  const capitalFlow = structuredFields.capital_flow

  return (
    <article className="cockpit-card">
      <div className="state-label">设置</div>
      <h2>能力圈配置</h2>
      <p>资产类型：{capability?.asset_types?.join('、') || '暂无'}</p>
      <p>纳入标的：{capability?.symbols?.join('、') || '暂无'}</p>
      <p>排除标的：{capability?.excluded_symbols?.join('、') || '暂无'}</p>
      <p>策略范围：{capability?.strategy_scope?.join('、') || '暂无'}</p>
      <h3>系统状态</h3>
      <p>SQLite：{textOrRaw(systemStatusText, system?.sqlite_status)}</p>
      <p>VecLite 索引状态：{textOrRaw(systemStatusText, system?.veclite_status)}</p>
      {/* DeepSeek 只展示配置状态，不展示完整密钥。 */}
      <p>DeepSeek：{textOrRaw(systemStatusText, system?.deepseek_status)}</p>
      <p>通知配置：当前后端未返回通知配置详情；只能通过保存设置接口更新。</p>
      <p>数据源：{system?.data_sources?.join('、') || '暂无'}</p>
      <h3>P40 本地运行就绪</h3>
      <p>就绪摘要：{readinessSummary}</p>
      <p>预检入口：go run ./cmd/agent --preflight --diagnostics ./tmp/preflight.json</p>
      <p>只展示本地诊断和人工处理提示；不发起资金动作、站外通知或规则生效。</p>
      <h3>市场快照状态</h3>
      <p>快照 ID：{marketSnapshot?.market_snapshot_id ?? '暂无'}</p>
      <p>交易日：{marketSnapshot?.trade_date ?? '暂无'}</p>
      <p>数据状态：{textOrRaw(marketStateText, marketSnapshot?.data_status)}</p>
      <p>PE/PB 分位：{marketSnapshot?.pe_percentile ?? '暂无'} / {marketSnapshot?.pb_percentile ?? '暂无'}</p>
      <p>情绪/流动性：{textOrRaw(marketStateText, marketSnapshot?.sentiment_state)} / {textOrRaw(marketStateText, marketSnapshot?.liquidity_state)}</p>
      <h3>结构化字段</h3>
      <p>融资融券：{margin ? `${margin.date ?? '暂无日期'}；余额 ${formatNumber(margin.margin_balance)}；变化率 ${formatPercent(margin.balance_change_rate)}` : '暂无'}</p>
      <p>成分财务：{financial ? `${financial.disclosure_date ?? '暂无日期'}；营收 ${formatNumber(financial.revenue)}；净利 ${formatNumber(financial.net_profit)}；增速 ${formatPercent(financial.growth)}` : '暂无'}</p>
      <p>资金流向：{capitalFlow ? `${capitalFlow.date ?? '暂无日期'}；净流入 ${formatNumber(capitalFlow.net_inflow)}；净流出 ${formatNumber(capitalFlow.net_outflow)}；净流向 ${formatNumber(capitalFlow.raw_net_flow)}` : '暂无真实 provider 读回'}</p>
      <h3>P40 数据源健康</h3>
      <p>仅展示公开只读数据状态；不会连接券商或发起交易。</p>
      {sourceHealth.length === 0 ? (
        <p>暂无 P40 数据源健康记录</p>
      ) : (
        <ul>
          {sourceHealth.map((item) => (
            <li key={`${item.source_name}-${item.data_category}`}>
              <strong>{item.source_name} · {textOrRaw(sourceCategoryText, item.data_category)} · {textOrRaw(sourceHealthStatusText, item.freshness)}</strong>
              <br />
              <span>数据日：{item.data_date ?? '暂无'}；等级：{item.source_level || '暂无'}；影响标的：{item.affected_symbols?.join('、') || '暂无'}</span>
              <br />
              <span>最近成功：{item.last_success_at || '暂无'}；最近失败：{item.last_failure_at || '暂无'}；失败类别：{textOrRaw(sourceHealthStatusText, item.failure_category) || '暂无'}</span>
            </li>
          ))}
        </ul>
      )}
    </article>
  )
}

type StructuredFields = {
  margin_financing?: Record<string, unknown>
  constituent_financial?: Record<string, unknown>
  capital_flow?: Record<string, unknown>
}

function readStructuredFields(marketSnapshot?: MarketSnapshot): StructuredFields {
  const metrics = marketSnapshot?.market_metrics
  const metadata = recordValue(metrics?.metadata)
  return recordValue(metadata?.p88_structured_fields) as StructuredFields
}

function recordValue(value: unknown): Record<string, unknown> {
  return value && typeof value === 'object' && !Array.isArray(value) ? value as Record<string, unknown> : {}
}

function formatNumber(value: unknown) {
  const numeric = Number(value)
  if (!Number.isFinite(numeric) || numeric === 0) return '暂无'
  return numeric.toLocaleString('zh-CN', { maximumFractionDigits: 2 })
}

function formatPercent(value: unknown) {
  const numeric = Number(value)
  if (!Number.isFinite(numeric)) return '暂无'
  return `${(numeric * 100).toFixed(2)}%`
}
