import { useEffect, useState } from 'react'
import { CapabilitySettingsPanel } from '../components/settings/CapabilitySettingsPanel'
import { StatusNotice } from '../components/status/StatusNotice'
import { Button, SummaryCard, type UITone } from '../components/ui'
import { buildLocalOpsModel, localOpsMetricTitle } from '../features/governance'
import { getLatestMarketSnapshot, getMarketSourceHealth, refreshMarket } from '../services/market'
import { getCapabilitySettings, getSystemSettings } from '../services/settings'
import type { PageErrorState } from '../shared/utils'
import { toPageErrorState } from '../shared/utils'
import type { MarketSnapshot, SourceHealthItem } from '../types/market'
import type { CapabilitySettings, SystemStatus } from '../types/settings'

export function SettingsPage() {
  const [capability, setCapability] = useState<CapabilitySettings>()
  const [system, setSystem] = useState<SystemStatus>()
  const [marketSnapshot, setMarketSnapshot] = useState<MarketSnapshot>()
  const [sourceHealth, setSourceHealth] = useState<SourceHealthItem[]>([])
  const [marketRefreshMessage, setMarketRefreshMessage] = useState('')
  const [errors, setErrors] = useState<PageErrorState[]>([])

  useEffect(() => {
    getCapabilitySettings()
      .then((res) => {
        setCapability(res.data)
        setErrors((items) => items.filter((item) => item.code !== 'capability'))
      })
      .catch((error: unknown) => {
        setCapability(undefined)
        setErrors((items) => [...items.filter((item) => item.code !== 'capability'), { ...toPageErrorState(error), code: 'capability' }])
      })
    getSystemSettings()
      .then((res) => {
        setSystem(res.data)
        setErrors((items) => items.filter((item) => item.code !== 'system'))
      })
      .catch((error: unknown) => {
        setSystem(undefined)
        setErrors((items) => [...items.filter((item) => item.code !== 'system'), { ...toPageErrorState(error), code: 'system' }])
      })
    getLatestMarketSnapshot()
      .then((res) => {
        setMarketSnapshot(res.data)
        setErrors((items) => items.filter((item) => item.code !== 'market'))
      })
      .catch((error: unknown) => {
        setMarketSnapshot(undefined)
        setErrors((items) => [...items.filter((item) => item.code !== 'market'), { ...toPageErrorState(error), code: 'market' }])
      })
    getMarketSourceHealth()
      .then((res) => {
        setSourceHealth(res.data?.sources ?? [])
        setErrors((items) => items.filter((item) => item.code !== 'source_health'))
      })
      .catch((error: unknown) => {
        setSourceHealth([])
        setErrors((items) => [...items.filter((item) => item.code !== 'source_health'), { ...toPageErrorState(error), code: 'source_health' }])
      })
  }, [])

  function handleMarketRefresh() {
    const symbols = capability?.symbols?.length ? capability.symbols : marketSnapshot?.symbol ? [marketSnapshot.symbol] : []
    refreshMarket(symbols.length ? { symbols } : {})
      .then(async () => {
        const [nextMarket, nextSourceHealth] = await Promise.all([getLatestMarketSnapshot(), getMarketSourceHealth()])
        setMarketSnapshot(nextMarket.data)
        setSourceHealth(nextSourceHealth.data?.sources ?? [])
        setMarketRefreshMessage('市场刷新完成；只更新本地行情事实和审计记录，不会执行交易。')
        setErrors((items) => items.filter((item) => item.code !== 'market_refresh'))
      })
      .catch((error: unknown) => {
        setMarketRefreshMessage('')
        setErrors((items) => [...items.filter((item) => item.code !== 'market_refresh'), { ...toPageErrorState(error), code: 'market_refresh' }])
      })
  }

  const localModel = buildLocalOpsModel({ system, capability, sourceHealth })

  return (
    <div>
      <h1 className="page-title">设置</h1>
      <section className={`daily-hero daily-tone-${localModel.overallTone}`} aria-label="本地配置与诊断总览">
        <div className="daily-hero-main">
          <div className="state-label">本地配置与诊断状态</div>
          <h2>{localModel.overallLabel}</h2>
          <p>{localModel.safetyNotes[0]}</p>
          <div className="daily-signal-grid quality-signal-grid">
            {localModel.metrics.map((metric) => (
              <SummaryCard key={metric.label} title={localOpsMetricTitle(metric.label)} value={metric.value} detail={metric.detail} tone={(metric.tone ?? 'unknown') as UITone} />
            ))}
          </div>
        </div>
        <aside className="daily-hero-side" aria-label="本地配置下一步">
          <strong>下一步本地复验</strong>
          <ul>
            {localModel.nextActions.map((action) => (
              <li key={action.label}>
                <a href={action.href} aria-label={`${action.label}入口`}>{action.label}</a>
                <span>{action.detail}</span>
              </li>
            ))}
          </ul>
        </aside>
      </section>
      {errors.length ? (
        <section className="status-notice-grid" aria-label="设置读取提示">
          {errors.map((error) => (
            <StatusNotice key={error.code ?? error.message} state={error.state} safeMessage={error.message} code={error.code} />
          ))}
        </section>
      ) : null}
      <CapabilitySettingsPanel
        capability={capability}
        system={system}
        marketSnapshot={marketSnapshot}
        sourceHealth={sourceHealth}
      />
      <article className="cockpit-card">
        <div className="state-label">市场刷新</div>
        <h2>本地市场数据</h2>
        <p>仅刷新行情事实与审计事件，不连接交易接口。</p>
        {(capability?.symbols?.length || marketSnapshot?.symbol) && <Button onClick={handleMarketRefresh}>刷新市场数据</Button>}
        {marketRefreshMessage && <p>{marketRefreshMessage}</p>}
      </article>
    </div>
  )
}
