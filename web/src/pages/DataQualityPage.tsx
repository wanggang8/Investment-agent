import { useEffect, useState } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { StatusNotice } from '../components/status/StatusNotice'
import { Button, SummaryCard, type UITone } from '../components/ui'
import { buildDataQualityExperienceModel } from '../features/dataQuality/dataQualityExperienceModel'
import { getEvidenceVerification, listEvidence } from '../services/evidence'
import { getLatestMarketSnapshot, getMarketSourceHealth } from '../services/market'
import { getReviewSummary } from '../services/review'
import { getSystemSettings } from '../services/settings'
import { createDataQualityGateResolution, getDataQualityGateResolution, getDataSourceQualityRegression, listDataQualityGateResolutions, retireDataQualityGateResolution } from '../services/dataSourceQuality'
import { getKnowledgeReadiness } from '../services/knowledgeReadiness'
import { opsStatusText, sourceCategoryText, sourceHealthStatusText, systemStatusText, textOrRaw, verificationStatusText } from '../shared/mappers/statusText'
import type { PageErrorState } from '../shared/utils'
import { toPageErrorState } from '../shared/utils'
import type { PageResult } from '../types/api'
import type { EvidenceItem, SourceVerification } from '../types/evidence'
import type { MarketSnapshot, SourceHealthItem } from '../types/market'
import type { ReviewSummary } from '../types/review'
import type { SystemStatus } from '../types/settings'
import type { DataQualityGateResolutionCheck, DataQualityGateResolutionListResponse, DataSourceQualityRegression } from '../types/dataSourceQuality'
import type { KnowledgeReadiness } from '../types/knowledgeReadiness'

type LoadState<T> = {
  data?: T
  error?: PageErrorState
}

type DataQualityState = {
  system: LoadState<SystemStatus>
  market: LoadState<MarketSnapshot>
  sourceHealth: LoadState<SourceHealthItem[]>
  evidence: LoadState<PageResult<EvidenceItem>>
  verification: LoadState<SourceVerification>
  review: LoadState<ReviewSummary>
  currentRegression: LoadState<DataSourceQualityRegression>
  gateResolution: LoadState<DataQualityGateResolutionCheck>
  resolutions: LoadState<DataQualityGateResolutionListResponse>
  knowledgeReadiness: LoadState<KnowledgeReadiness>
}

const initialState: DataQualityState = {
  system: {},
  market: {},
  sourceHealth: {},
  evidence: {},
  verification: {},
  review: {},
  currentRegression: {},
  gateResolution: {},
  resolutions: {},
  knowledgeReadiness: {},
}

let dataQualityLoadInFlight: { key: string; promise: Promise<DataQualityState> } | undefined

export function DataQualityPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const querySymbol = safeSymbol(searchParams.get('symbol') || '')
  const [state, setState] = useState<DataQualityState>(initialState)
  const [selectedSymbol, setSelectedSymbol] = useState(querySymbol)
  const [resolutionType, setResolutionType] = useState('scope_exclusion')
  const [scope, setScope] = useState('本次 release clean claim 排除 current local data health')
  const [reason, setReason] = useState('当前本地数据源存在降级，发布材料只声明有限范围')
  const [releaseImpact, setReleaseImpact] = useState('不得声明 current data healthy；不得把 resolution 描述为 policy passed')
  const [evidenceRef, setEvidenceRef] = useState('docs/release/acceptance')
  const [actionError, setActionError] = useState<PageErrorState | undefined>()
  const [actionBusy, setActionBusy] = useState(false)
  const [resolutionSymbol, setResolutionSymbol] = useState('')

  useEffect(() => {
    let mounted = true

    async function load() {
      const nextState = await loadDataQualityState(querySymbol)
      if (mounted) setState(nextState)
    }

    load()
    return () => {
      mounted = false
    }
  }, [querySymbol])

  useEffect(() => {
    setSelectedSymbol(querySymbol)
  }, [querySymbol])

  useEffect(() => {
    if (!resolutionSymbol && (state.gateResolution.data?.symbol || state.market.data?.symbol)) {
      setResolutionSymbol(state.gateResolution.data?.symbol || state.market.data?.symbol || '')
    }
  }, [resolutionSymbol, state.gateResolution.data?.symbol, state.market.data?.symbol])

  const sourceHealth = state.sourceHealth.data ?? []
  const evidenceItems = state.evidence.data?.items ?? []
  const review = state.review.data
  const degradedWorkflows = review?.degraded_workflows ?? []
  const degradedCount = review?.degraded_count ?? 0
  const missingEvidenceCount = review?.missing_evidence_count ?? 0
  const qualityGateText = `质量门禁：复盘降级 ${degradedCount} 条，缺证据 ${missingEvidenceCount} 条`
  const impactedDecisionId = degradedWorkflows[0]?.decision_id
  const gateResolution = state.gateResolution.data
  const currentPolicy = gateResolution?.policy ?? state.currentRegression.data?.policy
  const activeResolution = gateResolution?.active_resolution
  const activeResolutions = safeArray(state.resolutions.data?.items)
  const currentPolicyReasons = [
    ...safeArray(currentPolicy?.blocking_reasons),
    ...safeArray(currentPolicy?.waiver_reasons),
  ]
  const allowedClaims = safeArray(gateResolution?.allowed_claims)
  const prohibitedClaims = safeArray(gateResolution?.prohibited_claims)
  const qualityModel = buildDataQualityExperienceModel({
    system: state.system.data,
    market: state.market.data,
    sourceHealth,
    evidenceItems,
    evidenceTotal: state.evidence.data?.total,
    verification: state.verification.data,
    review,
    currentRegression: state.currentRegression.data,
    gateResolution: state.gateResolution.data,
    errors: [state.system.error, state.market.error, state.sourceHealth.error, state.evidence.error, state.verification.error, state.review.error, state.currentRegression.error, state.gateResolution.error, state.resolutions.error, state.knowledgeReadiness.error, actionError].filter((item): item is PageErrorState => Boolean(item)),
  })
  const knowledgeReadiness = state.knowledgeReadiness.data
  const readinessDeps = safeArray(knowledgeReadiness?.data_dependencies)
  const readinessReferences = safeArray(knowledgeReadiness?.knowledge_references)
  const readinessImpact = knowledgeReadiness?.feature_impacts?.[0]
  const activeSymbol = querySymbol || state.market.data?.symbol || knowledgeReadiness?.symbol || ''

  async function refreshResolutionState(symbol?: string) {
    const nextGate = await settle(getDataQualityGateResolution(symbol))
    const nextResolutions = await settle(listDataQualityGateResolutions(symbol, 'active'))
    if (nextGate.data?.policy.verdict === 'blocked') {
      setResolutionType('scope_exclusion')
    }
    setState((current) => ({ ...current, gateResolution: nextGate, resolutions: nextResolutions }))
  }

  async function handleCreateResolution() {
    const symbol = resolutionSymbol.trim() || gateResolution?.symbol || state.market.data?.symbol || '000300'
    const effectiveResolutionType = currentPolicy?.verdict === 'blocked' ? 'scope_exclusion' : resolutionType
    setActionBusy(true)
    setActionError(undefined)
    try {
      const response = await createDataQualityGateResolution({
        symbol,
        resolution_type: effectiveResolutionType,
        scope,
        reason,
        release_impact: releaseImpact,
        evidence_ref: evidenceRef,
      })
      const nextResolutions = await settle(listDataQualityGateResolutions(symbol, 'active'))
      setState((current) => ({ ...current, gateResolution: { data: response.data }, resolutions: nextResolutions }))
    } catch (error: unknown) {
      setActionError(toPageErrorState(error))
    } finally {
      setActionBusy(false)
    }
  }

  function handleApplySymbol() {
    const symbol = safeSymbol(selectedSymbol)
    if (symbol) {
      setSearchParams({ symbol })
    } else {
      setSearchParams({})
    }
  }

  async function handleRetireResolution(resolutionId: string) {
    setActionBusy(true)
    setActionError(undefined)
    try {
      const response = await retireDataQualityGateResolution(resolutionId)
      const symbol = response.data?.symbol || gateResolution?.symbol || state.market.data?.symbol
      const nextResolutions = await settle(listDataQualityGateResolutions(symbol, 'active'))
      setState((current) => ({ ...current, gateResolution: { data: response.data }, resolutions: nextResolutions }))
    } catch (error: unknown) {
      setActionError(toPageErrorState(error))
    } finally {
      setActionBusy(false)
    }
  }

  return (
    <div>
      <h1 className="page-title">数据质量可观测</h1>
      <p className="page-placeholder">只读聚合 source health、证据、检索、LLM 与诊断状态；这里只展示质量事实和导航，不发起后台变更或规则确认。</p>
      <section className="cockpit-card" aria-label="数据质量标的筛选">
        <div className="state-label">当前标的</div>
        <div className="quality-symbol-filter">
          <label>
            标的代码
            <input aria-label="数据质量标的" value={selectedSymbol} placeholder={activeSymbol || '默认最新标的'} onChange={(event) => setSelectedSymbol(event.target.value)} />
          </label>
          <Button onClick={handleApplySymbol}>切换标的</Button>
        </div>
        <p>当前查看：{activeSymbol || '默认最新本地市场事实'}</p>
      </section>

      <section className={`daily-hero daily-tone-${qualityModel.overallTone}`} aria-label="数据质量总览">
        <div className="daily-hero-main">
          <div className="state-label">数据质量总览</div>
          <h2>{qualityModel.overallLabel}</h2>
          <p>{qualityModel.safetyNotes[0]}</p>
          <div className="daily-signal-grid quality-signal-grid">
            {qualityModel.qualitySignals.map((signal) => (
              <SummaryCard key={signal.label} title={`${signal.label} 信号`} value={signal.value} detail={signal.detail} tone={signal.tone as UITone} />
            ))}
          </div>
        </div>
        <aside className="daily-hero-side" aria-label="数据质量下一步">
          <strong>下一步本地检查</strong>
          <ul>
            {qualityModel.nextActions.map((action) => (
              <li key={action.label}>
                <Link to={action.href}>{action.label}</Link>
                <span>{action.detail}</span>
              </li>
            ))}
          </ul>
        </aside>
      </section>

      <section className="cockpit-grid" aria-label="知识与数据准备度">
        <article className="cockpit-card">
          <div className="state-label">知识与数据准备度</div>
          <h2>{readinessStatusText(knowledgeReadiness?.status)}</h2>
          {state.knowledgeReadiness.error && <StatusNotice state={state.knowledgeReadiness.error.state} safeMessage={state.knowledgeReadiness.error.message} code={state.knowledgeReadiness.error.code} />}
          <p>{knowledgeReadiness?.symbol_profile?.known ? `${knowledgeReadiness.symbol_profile.name || knowledgeReadiness.symbol_profile.symbol} · ${knowledgeReadiness.symbol_profile.asset_type || '已记录'} · 跟踪 ${knowledgeReadiness.symbol_profile.tracked_index_symbol || '暂无'}` : '标的画像未准备，当前不能伪造成 ready。'}</p>
          <p>{safeArray(knowledgeReadiness?.safety_notes)[0] || '内置知识只作为纪律和分析上下文，不作为正式市场证据。'}</p>
          <p>LLM 上下文：{knowledgeReadiness?.llm_context_summary ? '已附加知识与数据准备度摘要' : '暂无可用摘要'}</p>
          {readinessImpact && <p>{readinessImpact.impact}</p>}
          <ul className="quality-list">
            {readinessReferences.map((entry) => (
              <li key={entry.knowledge_id}>
                <strong>{entry.title}</strong>
                <span>{knowledgeCategoryText(entry.category)} · {entry.llm_context_allowed ? '可作为 LLM 上下文' : '不进入 LLM 上下文'} · {entry.formal_evidence_allowed ? '可作正式证据' : '不作正式证据'}</span>
              </li>
            ))}
          </ul>
        </article>
        <article className="cockpit-card">
          <div className="state-label">数据依赖矩阵</div>
          <h2>缺口影响</h2>
          {readinessDeps.length ? (
            <ul className="quality-list">
              {readinessDeps.map((dep) => (
                <li key={dep.category}>
                  <strong>{readinessCategoryText(dep.category)} · {readinessStatusText(dep.status)}</strong>
                  <span>{dep.required ? '必需' : '可选'} · freshness：{safeStatusText(dep.freshness)} · 等级：{safeLevel(dep.source_level)}</span>
                  <span>来源：{safeReadinessSource(dep.source_name)} · 类型：{safeReadinessSource(dep.source_type)} · 日期：{safeDate(dep.data_date)} · request：{safeRequestID(dep.request_id)} · 标的：{safeAffectedSymbols(dep.affected_symbols)}</span>
                  <span>{dep.safe_degradation || '暂无降级说明'}</span>
                </li>
              ))}
            </ul>
          ) : (
            <p>暂无数据依赖矩阵。</p>
          )}
        </article>
      </section>

      <section className="cockpit-grid" aria-label="数据质量可观测区域">
        <article className="cockpit-card">
          <div className="state-label">数据源健康</div>
          <h2>source health</h2>
          {state.sourceHealth.error && <StatusNotice state={state.sourceHealth.error.state} safeMessage={state.sourceHealth.error.message} code={state.sourceHealth.error.code} />}
          {state.market.error && <StatusNotice state={state.market.error.state} safeMessage={state.market.error.message} code={state.market.error.code} />}
          <p>市场数据状态：{safeStatusText(state.market.data?.data_status)}</p>
          <p>影响标的：{state.market.data?.symbol ?? '暂无'}</p>
          {state.currentRegression.error && <StatusNotice state={state.currentRegression.error.state} safeMessage={state.currentRegression.error.message} code={state.currentRegression.error.code} />}
          {state.gateResolution.error && <StatusNotice state={state.gateResolution.error.state} safeMessage={state.gateResolution.error.message} code={state.gateResolution.error.code} />}
          {state.resolutions.error && <StatusNotice state={state.resolutions.error.state} safeMessage={state.resolutions.error.message} code={state.resolutions.error.code} />}
          {actionError && <StatusNotice state={actionError.state} safeMessage={actionError.message} code={actionError.code} />}
          {currentPolicy && (
            <div className="quality-list" aria-label="当前数据策略">
              <p>当前数据策略：{policyVerdictText(currentPolicy.verdict)}；release gate：{policyGateText(currentPolicy.release_gate)}</p>
              <p>{currentPolicy.safety_note}</p>
              {currentPolicyReasons.slice(0, 3).map((reason) => (
                <p key={reason}>{safePolicyReason(reason)}</p>
              ))}
            </div>
          )}
          {gateResolution && (
            <section className="quality-resolution-panel" aria-label="当前数据门禁处置">
              <div>
                <div className="state-label">发布声明边界</div>
                <h3>{releaseClaimStateText(gateResolution.release_claim_state)}</h3>
                <p>clean data claim：{gateResolution.clean_data_claim_allowed ? '允许' : '不允许'}</p>
                <p>{gateResolution.safety_note}</p>
              </div>
              <div className="quality-claim-grid">
                <div>
                  <strong>允许声明</strong>
                  <ul>
                    {allowedClaims.map((claim) => <li key={claim}>{claim}</li>)}
                  </ul>
                </div>
                <div>
                  <strong>禁止声明</strong>
                  <ul>
                    {prohibitedClaims.map((claim) => <li key={claim}>{claim}</li>)}
                  </ul>
                </div>
              </div>
              <div className="quality-resolution-form">
                <label>
                  标的
                  <input value={resolutionSymbol} onChange={(event) => setResolutionSymbol(event.target.value)} />
                </label>
              </div>
              {activeResolution ? (
                <div className="quality-active-resolution">
                  <strong>{resolutionTypeText(activeResolution.resolution_type)} · {activeResolution.status}</strong>
                  <span>{activeResolution.scope}</span>
                  <span>{activeResolution.release_impact}</span>
                  <button type="button" disabled={actionBusy} onClick={() => handleRetireResolution(activeResolution.resolution_id)}>退役处置</button>
                </div>
              ) : (
                <div className="quality-resolution-form">
                  <label>
                    类型
                    <select value={resolutionType} onChange={(event) => setResolutionType(event.target.value)} disabled={actionBusy || currentPolicy?.verdict === 'passed'}>
                      <option value="scope_exclusion">范围排除</option>
                      {currentPolicy?.verdict === 'waiver_required' && <option value="waiver">豁免记录</option>}
                    </select>
                  </label>
                  <label>
                    范围
                    <textarea value={scope} onChange={(event) => setScope(event.target.value)} />
                  </label>
                  <label>
                    原因
                    <textarea value={reason} onChange={(event) => setReason(event.target.value)} />
                  </label>
                  <label>
                    发布影响
                    <textarea value={releaseImpact} onChange={(event) => setReleaseImpact(event.target.value)} />
                  </label>
                  <label>
                    证据引用
                    <input value={evidenceRef} onChange={(event) => setEvidenceRef(event.target.value)} />
                  </label>
                  <button type="button" disabled={actionBusy || currentPolicy?.verdict === 'passed'} onClick={handleCreateResolution}>记录处置</button>
                </div>
              )}
              {activeResolutions.length > 0 && (
                <ul className="quality-list" aria-label="活跃处置记录">
                  {activeResolutions.slice(0, 3).map((item) => (
                    <li key={item.resolution_id}>
                      <strong>{resolutionTypeText(item.resolution_type)} · {item.symbol}</strong>
                      <span>{item.reason}</span>
                      <span>{item.evidence_ref || '暂无证据引用'}</span>
                    </li>
                  ))}
                </ul>
              )}
            </section>
          )}
          {sourceHealth.length ? (
            <ul className="quality-list">
              {sourceHealth.slice(0, 4).map((item) => (
                <li key={`${item.source_name}-${item.data_category}`}>
                  <strong>{item.source_name} · {textOrRaw(sourceCategoryText, item.data_category)} · {safeStatusText(item.freshness)}</strong>
                  <span>数据日：{item.data_date ?? '暂无'}；等级：{safeLevel(item.source_level)}</span>
                  <span>影响标的：{item.affected_symbols?.join('、') || '暂无'}</span>
                  <span>最近成功：{safeDate(item.last_success_at)}；最近失败：{safeDate(item.last_failure_at)}</span>
                  <span>失败类别：{safeStatusText(item.failure_category)}</span>
                </li>
              ))}
            </ul>
          ) : (
            <p>暂无数据源健康记录。</p>
          )}
          <div className="link-row">
            <Link to="/settings">查看设置</Link>
            <Link to="/workbench">返回工作台</Link>
            <button type="button" disabled={actionBusy} onClick={() => refreshResolutionState(resolutionSymbol.trim() || gateResolution?.symbol || state.market.data?.symbol)}>检查门禁处置</button>
          </div>
        </article>

        <article className="cockpit-card">
          <div className="state-label">证据与检索</div>
          <h2>Evidence / RAG</h2>
          {state.evidence.error && <StatusNotice state={state.evidence.error.state} safeMessage={state.evidence.error.message} code={state.evidence.error.code} />}
          {state.verification.error && <StatusNotice state={state.verification.error.state} safeMessage={state.verification.error.message} code={state.verification.error.code} />}
          <p>证据数量：{state.evidence.data?.total ?? evidenceItems.length}</p>
          <p>验证状态：{safeVerificationText(state.verification.data?.verification_status)}</p>
          <p>独立信源：{state.verification.data?.independent_source_count ?? 0}</p>
          <p>高等级独立信源：{state.verification.data?.high_grade_independent_source_count ?? 0}</p>
          <p>最高信源等级：{safeLevel(state.verification.data?.highest_source_level)}</p>
          <p>VecLite：{safeStatusText(state.system.data?.veclite_status)}</p>
          {evidenceItems.length ? (
            <ul className="quality-list">
              {evidenceItems.slice(0, 3).map((item) => (
                <li key={item.evidence_id}>
                  <strong>{item.source_name}</strong>
                  <span>{safeLevel(item.source_level)} · {safeVerificationText(item.verification_status)}</span>
                </li>
              ))}
            </ul>
          ) : (
            <p>暂无证据记录。</p>
          )}
          <div className="link-row">
            <Link to="/evidence">查看证据</Link>
            <Link to="/audit">查看审计</Link>
          </div>
        </article>

        <article className="cockpit-card">
          <div className="state-label">LLM 质量</div>
          <h2>分析材料质量</h2>
          {state.system.error && <StatusNotice state={state.system.error.state} safeMessage={state.system.error.message} code={state.system.error.code} />}
          {state.review.error && <StatusNotice state={state.review.error.state} safeMessage={state.review.error.message} code={state.review.error.code} />}
          <p>DeepSeek：{safeStatusText(state.system.data?.deepseek_status)}</p>
          <p>复盘状态：{safeStatusText(review?.ops_status?.review_status)}</p>
          <p>数据源状态：{safeStatusText(review?.ops_status?.data_source_status)}</p>
          <p>索引状态：{safeStatusText(review?.ops_status?.index_status)}</p>
          <p>{qualityGateText}</p>
          {review?.ops_status?.explanation && <p>{safeDiagnosticSummary()}</p>}
          <div className="link-row">
            <Link to="/review">查看质量复盘</Link>
          </div>
        </article>

        <article className="cockpit-card">
          <div className="state-label">影响范围与下一步</div>
          <h2>受影响工作流</h2>
          {degradedWorkflows.length ? (
            <ul className="quality-list">
              {degradedWorkflows.slice(0, 4).map((workflow) => (
                <li key={workflow.decision_id}>
                  <strong>{workflow.decision_id} · {workflow.symbol || '暂无标的'} · {safeStatusText(workflow.status)}</strong>
                </li>
              ))}
            </ul>
          ) : (
            <p>暂无受影响工作流记录。</p>
          )}
          <p>下一步只导航到权威页面检查，不发起后台变更、外部动作或规则生效。</p>
          <div className="link-row">
            <Link to="/risk-alerts">查看风险预警</Link>
            {impactedDecisionId && <Link to={`/decisions/${encodeURIComponent(impactedDecisionId)}`}>查看受影响决策</Link>}
            <Link to="/review">查看复盘</Link>
          </div>
        </article>
      </section>
    </div>
  )
}

function loadDataQualityState(symbol?: string) {
  const key = safeSymbol(symbol || '')
  if (!dataQualityLoadInFlight || dataQualityLoadInFlight.key !== key) {
    dataQualityLoadInFlight = { key, promise: loadDataQualityStateOnce(key || undefined).finally(() => {
      dataQualityLoadInFlight = undefined
    }) }
  }
  return dataQualityLoadInFlight.promise
}

async function loadDataQualityStateOnce(symbol?: string): Promise<DataQualityState> {
  const system = await settle(getSystemSettings())
  const market = await settle(getLatestMarketSnapshot(symbol))
  const targetSymbol = symbol || market.data?.symbol
  const sourceHealth = await settle(getMarketSourceHealth(targetSymbol))
  const evidence = await settle(listEvidence())
  const verification = await settle(getEvidenceVerification())
  const review = await settle(getReviewSummary())
  const currentRegression = await settle(getDataSourceQualityRegression('current', targetSymbol))
  const gateResolution = await settle(getDataQualityGateResolution(targetSymbol))
  const resolutions = await settle(listDataQualityGateResolutions(targetSymbol, 'active'))
  const knowledgeReadiness = await settle(getKnowledgeReadiness(targetSymbol))

  return {
    system,
    market,
    sourceHealth: sourceHealth.data ? { data: sourceHealth.data.sources } : { error: sourceHealth.error },
    evidence,
    verification,
    review,
    currentRegression,
    gateResolution,
    resolutions,
    knowledgeReadiness,
  }
}

function safeSymbol(value: string) {
  const symbol = value.trim()
  return /^[A-Za-z0-9._-]{1,24}$/.test(symbol) ? symbol : ''
}

async function settle<T>(promise: Promise<{ data?: T }>): Promise<LoadState<T>> {
  try {
    const response = await promise
    return { data: response.data }
  } catch (error: unknown) {
    return { error: toPageErrorState(error) }
  }
}

function safeStatusText(value?: string) {
  if (!value) return '暂无'
  return sourceHealthStatusText[value] || systemStatusText[value] || opsStatusText[value] || '已记录异常'
}

function safeArray<T>(value: T[] | null | undefined): T[] {
  return Array.isArray(value) ? value : []
}

function safeVerificationText(value?: string) {
  if (!value) return '暂无'
  return verificationStatusText[value] || '已记录异常'
}

function safeLevel(value?: string) {
  if (!value) return '暂无'
  return /^[A-Z]$/.test(value) ? value : '已记录'
}

function safeDate(value?: string) {
  if (!value) return '暂无'
  return value.includes('/') || value.includes('\\') ? '已记录' : value
}

function safeReadinessSource(value?: string) {
  if (!value) return '暂无'
  return value.replace(/[\\/][^\s]+/g, '已脱敏路径')
}

function safeRequestID(value?: string) {
  if (!value) return '暂无'
  return /^[A-Za-z0-9._:-]{1,96}$/.test(value) ? value : '已记录'
}

function safeAffectedSymbols(value: string[] | null | undefined) {
  const symbols = safeArray(value).filter((item) => /^[A-Za-z0-9._-]{1,24}$/.test(item))
  return symbols.length ? symbols.join('、') : '暂无'
}

function safeDiagnosticSummary() {
  return '存在已脱敏诊断摘要。'
}

function policyVerdictText(value?: string) {
  switch (value) {
    case 'passed': return '通过'
    case 'waiver_required': return '需豁免记录'
    case 'blocked': return '阻断'
    default: return '待检查'
  }
}

function policyGateText(value?: string) {
  switch (value) {
    case 'pass': return '通过'
    case 'waiver_required': return '需豁免记录'
    case 'block': return '阻断'
    default: return '待检查'
  }
}

function releaseClaimStateText(value?: string) {
  switch (value) {
    case 'pass': return '当前本地数据门禁通过'
    case 'resolved_with_waiver': return '已记录当前数据质量豁免'
    case 'resolved_with_scope_exclusion': return '已排除 current data clean claim'
    case 'requires_resolution': return '需要人工处置'
    default: return '待检查'
  }
}

function resolutionTypeText(value?: string) {
  switch (value) {
    case 'waiver': return '豁免记录'
    case 'scope_exclusion': return '范围排除'
    default: return '处置记录'
  }
}

function safePolicyReason(value: string) {
  return value
    .replace(/sk-[A-Za-z0-9_-]+/g, '[REDACTED_KEY]')
    .replace(/\/Users\/[A-Za-z0-9._-]+[^\s，。；;]*/g, '[REDACTED_PATH]')
    .replace(/\bSELECT\s+\*\s+FROM\b/gi, 'SELECT [REDACTED]')
    .replace(/prompt\s*:/gi, 'prompt [REDACTED]')
    .replace(/raw\s+(vendor|provider|HTTP)[^\s，。；;]*/gi, 'raw [REDACTED]')
}

function readinessStatusText(value?: string) {
  switch (value) {
    case 'ready': return '已准备'
    case 'degraded': return '降级'
    case 'blocked': return '阻断'
    default: return '待检查'
  }
}

function readinessCategoryText(value?: string) {
  switch (value) {
    case 'symbol_profile': return '标的画像'
    case 'fund_profile': return '基金画像'
    case 'tracked_index': return '跟踪指数'
    case 'market_price': return '市场价格'
    case 'valuation_percentiles': return '估值分位'
    case 'liquidity': return '流动性'
    case 'sentiment_proxy': return '情绪代理'
    case 'active_rule': return '生效规则'
    case 'formal_evidence': return '正式证据'
    case 'rag_index': return 'RAG 索引'
    case 'llm_context': return 'LLM 上下文'
    default: return value || '未分类'
  }
}

function knowledgeCategoryText(value?: string) {
  switch (value) {
    case 'master_principle': return '大师原则'
    case 'discipline_rule': return '纪律规则'
    case 'risk_sop': return '风险 SOP'
    case 'symbol_profile': return '标的画像'
    default: return value || '未分类'
  }
}
