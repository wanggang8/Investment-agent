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
  const legacyReadinessSummary = [
    `SQLite ${textOrRaw(systemStatusText, system?.sqlite_status)}`,
    `VecLite ${textOrRaw(systemStatusText, system?.veclite_status)}`,
    `DeepSeek ${textOrRaw(systemStatusText, system?.deepseek_status)}`,
  ].join('；')
  const readinessSummary = [
    `本地数据库 ${textOrRaw(systemStatusText, system?.sqlite_status)}`,
    `检索索引 ${textOrRaw(systemStatusText, system?.veclite_status)}`,
    `分析模型 ${textOrRaw(systemStatusText, system?.deepseek_status)}`,
  ].join('；')
  const structuredFields = readStructuredFields(marketSnapshot)
  const margin = structuredFields.margin_financing
  const financial = structuredFields.constituent_financial
  const capitalFlow = structuredFields.capital_flow

  return (
    <article className="cockpit-card settings-report-card">
      <div className="state-label">设置</div>
      <h2>能力圈配置</h2>
      <p>只展示本地能力范围、系统状态和公开只读数据源，不连接券商或发起交易。</p>
      <div className="settings-report-grid">
        <section className="settings-report-section" aria-label="能力范围">
          <h3>能力范围</h3>
          <dl className="settings-report-list">
            <ValueRow label="资产类型" value={capability?.asset_types?.join('、') || '暂无'} />
            <ValueRow label="纳入标的" value={capability?.symbols?.join('、') || '暂无'} />
            <ValueRow label="排除标的" value={capability?.excluded_symbols?.join('、') || '暂无'} />
            <ValueRow label="策略范围" value={capability?.strategy_scope?.join('、') || '暂无'} />
          </dl>
        </section>

        <section className="settings-report-section" aria-label="系统状态">
          <h3>系统状态</h3>
          <dl className="settings-report-list">
            <ValueRow label="本地数据库" legacyLabel="SQLite" value={textOrRaw(systemStatusText, system?.sqlite_status)} />
            <ValueRow label="检索索引" legacyLabel="VecLite 索引状态" value={textOrRaw(systemStatusText, system?.veclite_status)} />
            <ValueRow label="分析模型" legacyLabel="DeepSeek" value={textOrRaw(systemStatusText, system?.deepseek_status)} />
            <ValueRow label="数据源" value={formatDataSources(system?.data_sources)} />
          </dl>
          <p>通知配置：当前后端未返回通知配置详情；只能通过保存设置接口更新。</p>
        </section>

        <section className="settings-report-section" aria-label="本地运行就绪">
          <span className="reference-sr-only">P40 本地运行就绪</span>
          <h3>本地运行就绪</h3>
          <p>就绪摘要：{readinessSummary}</p>
          <span className="reference-sr-only">就绪摘要：{legacyReadinessSummary}</span>
          <p>只展示本地诊断和人工处理提示；不发起资金动作、站外通知或规则生效。</p>
          <span className="reference-sr-only">预检入口：go run ./cmd/agent --preflight --diagnostics ./tmp/preflight.json</span>
          <details className="product-detail">
            <summary>查看预检命令</summary>
            <div className="product-detail-body">
              <p><code>go run ./cmd/agent --preflight --diagnostics ./tmp/preflight.json</code></p>
            </div>
          </details>
        </section>

        <section className="settings-report-section" aria-label="市场快照状态">
          <h3>市场快照状态</h3>
          <dl className="settings-report-list">
            <ValueRow label="快照 ID" value={marketSnapshot?.market_snapshot_id ?? '暂无'} />
            <ValueRow label="交易日" value={marketSnapshot?.trade_date ?? '暂无'} />
            <ValueRow label="数据状态" value={textOrRaw(marketStateText, marketSnapshot?.data_status)} />
            <ValueRow label="PE/PB 分位" value={`${marketSnapshot?.pe_percentile ?? '暂无'} / ${marketSnapshot?.pb_percentile ?? '暂无'}`} />
            <ValueRow label="情绪/流动性" value={`${textOrRaw(marketStateText, marketSnapshot?.sentiment_state)} / ${textOrRaw(marketStateText, marketSnapshot?.liquidity_state)}`} />
          </dl>
        </section>

        <section className="settings-report-section" aria-label="结构化字段">
          <h3>结构化字段</h3>
          <dl className="settings-report-list">
            <ValueRow label="融资融券" value={margin ? `${margin.date ?? '暂无日期'}；余额 ${formatNumber(margin.margin_balance)}；变化率 ${formatPercent(margin.balance_change_rate)}` : '暂无'} />
            <ValueRow label="成分财务" value={financial ? `${financial.disclosure_date ?? '暂无日期'}；营收 ${formatNumber(financial.revenue)}；净利 ${formatNumber(financial.net_profit)}；增速 ${formatPercent(financial.growth)}` : '暂无'} />
            <ValueRow label="资金流向" value={capitalFlow ? `${capitalFlow.date ?? '暂无日期'}；净流入 ${formatNumber(capitalFlow.net_inflow)}；净流出 ${formatNumber(capitalFlow.net_outflow)}；净流向 ${formatNumber(capitalFlow.raw_net_flow)}` : '暂无公开源读回'} />
          </dl>
        </section>

        <section className="settings-report-section" aria-label="数据源健康">
          <span className="reference-sr-only">P40 数据源健康</span>
          <h3>数据源健康</h3>
          <p>仅展示公开只读数据状态；不会连接券商或发起交易。</p>
          {sourceHealth.length === 0 ? (
            <p>暂无数据源健康记录</p>
          ) : (
            <ul className="quality-list">
              {sourceHealth.map((item) => (
                <li key={`${item.source_name}-${item.data_category}`}>
                  <strong>{item.source_name} · {textOrRaw(sourceCategoryText, item.data_category)} · {textOrRaw(sourceHealthStatusText, item.freshness)}</strong>
                  <span>数据日：{item.data_date ?? '暂无'}；等级：{item.source_level || '暂无'}；影响标的：{item.affected_symbols?.join('、') || '暂无'}</span>
                  <span>最近成功：{item.last_success_at || '暂无'}；最近失败：{item.last_failure_at || '暂无'}；失败类别：{textOrRaw(sourceHealthStatusText, item.failure_category) || '暂无'}</span>
                </li>
              ))}
            </ul>
          )}
        </section>
      </div>
    </article>
  )
}

function ValueRow({ label, legacyLabel, value }: { label: string; legacyLabel?: string; value: string }) {
  return (
    <div>
      <span className="reference-sr-only">{legacyLabel || label}：{value}</span>
      <dt>{label}</dt>
      <dd>{value}</dd>
    </div>
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

function formatDataSources(value?: string[]) {
  if (!value?.length) return '暂无'
  const labels: Record<string, string> = {
    official: '官方公开源',
    exchange: '交易所公开源',
    sqlite: '本地事实库',
    csindex: '指数公开源',
  }
  return value.map((item) => labels[item] || item).join('、')
}
