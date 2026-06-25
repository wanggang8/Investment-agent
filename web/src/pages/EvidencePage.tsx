import { useEffect, useState } from 'react'
import { EvidenceTable } from '../components/evidence/EvidenceTable'
import { StatusNotice } from '../components/status/StatusNotice'
import { refreshEvidence, getEvidenceVerification, listEvidence, rebuildEvidenceIndex } from '../services/evidence'
import type { PageErrorState } from '../shared/utils'
import { toPageErrorState } from '../shared/utils'
import type { EvidenceItem, SourceVerification } from '../types/evidence'

export function EvidencePage() {
  const [items, setItems] = useState<EvidenceItem[]>([])
  const [verification, setVerification] = useState<SourceVerification>()
  const [actionMessage, setActionMessage] = useState('')
  const [errorState, setErrorState] = useState<PageErrorState>()

  useEffect(() => {
    loadEvidence()
  }, [])

  function loadEvidence() {
    listEvidence()
      .then((res) => {
        setItems(res.data?.items ?? [])
        setErrorState(undefined)
      })
      .catch((error: unknown) => {
        setItems([])
        setErrorState(toPageErrorState(error))
      })
    getEvidenceVerification()
      .then((res) => setVerification(res.data))
      .catch(() => setVerification(undefined))
  }

  function handleRefreshEvidence() {
    refreshEvidence({ refresh_scope: 'all', include_background: true })
      .then((res) => {
        setActionMessage(`情报刷新完成；索引状态 ${res.data?.index_status ?? 'unknown'}。`)
        setErrorState(undefined)
        loadEvidence()
      })
      .catch((error: unknown) => {
        setActionMessage('')
        setErrorState(toPageErrorState(error))
      })
  }

  function handleRebuildIndex() {
    rebuildEvidenceIndex()
      .then((res) => {
        setActionMessage(`索引重建完成；已索引 ${res.data?.indexed_count ?? 0} 条，跳过 ${res.data?.skipped_count ?? 0} 条。`)
        setErrorState(undefined)
        loadEvidence()
      })
      .catch((error: unknown) => {
        setActionMessage('')
        setErrorState(toPageErrorState(error))
      })
  }

  return (
    <div>
      <h1 className="page-title">情报与证据</h1>
      {errorState && <StatusNotice state={errorState.state} safeMessage={errorState.message} code={errorState.code} />}
      <section className={`daily-hero daily-tone-${verification?.verification_status === 'satisfied' ? 'success' : 'warning'}`} aria-label="证据首屏状态">
        <div className="daily-hero-main">
          <div className="state-label">证据与规则快照</div>
          <h2>{verification?.verification_status === 'satisfied' ? '证据链满足当前核验要求' : '证据不足或尚未完成核验'}</h2>
          <p>高等级正式证据优先进入裁决链；背景材料只辅助理解，不替代规则裁决或人工复核。</p>
          <dl className="daily-hero-meta">
            <div>
              <dt>独立信源</dt>
              <dd>{verification?.independent_source_count ?? 0}</dd>
            </div>
            <div>
              <dt>高等级信源</dt>
              <dd>{verification?.high_grade_independent_source_count ?? 0}</dd>
            </div>
            <div>
              <dt>证据条数</dt>
              <dd>{items.length}</dd>
            </div>
            <div>
              <dt>最高等级</dt>
              <dd>{verification?.highest_source_level || '暂无'}</dd>
            </div>
          </dl>
        </div>
        <aside className="daily-hero-side" aria-label="证据下一步">
          <strong>下一步人工核查</strong>
          <p>先核对信源等级、发布时间和引用，再进入决策闭环或审计页面。</p>
          <strong>禁止动作</strong>
          <p>不把背景材料伪装为正式证据，不自动生成交易指令。</p>
        </aside>
      </section>
      <section className="cockpit-grid" aria-label="证据解释入口">
        <article className="cockpit-card">
          <div className="state-label">证据可信度</div>
          <h2>证据可信度</h2>
          <p>先看信源等级、独立信源和核验状态，再进入证据明细；背景级材料只辅助理解。</p>
          <p>独立信源数量：{verification?.independent_source_count ?? 0}</p>
          <p>最高信源等级：{verification?.highest_source_level || '暂无'}</p>
          <p>最新发布时间：{verification?.latest_published_at || '暂无'}</p>
        </article>
        <article className="cockpit-card">
          <div className="state-label">解释路径</div>
          <h2>决策解释入口</h2>
          <p>高可信等级可进入正式证据链；背景级材料只辅助理解，不替代规则裁决或人工复核。</p>
          <span className="reference-sr-only">S/A/B 级可作为正式证据；C 级和 background 只辅助理解，不替代规则裁决或人工复核。</span>
          <div className="link-row">
            <a href="/workbench">返回工作台</a>
            <a href="/decision-loop">查看决策闭环</a>
            <a href="/audit">查看审计</a>
          </div>
        </article>
      </section>
      <article className="cockpit-card form-card">
        <div className="state-label">证据维护</div>
        <h2>情报刷新与索引重建</h2>
        <p>只更新本地情报、检索索引和审计记录；刷新后仍需人工核验证据等级和引用。</p>
        <div className="action-row action-row-left">
          <button type="button" onClick={handleRefreshEvidence}>刷新情报</button>
          <button type="button" onClick={handleRebuildIndex}>重建索引</button>
        </div>
        {actionMessage && <p>{actionMessage}</p>}
      </article>
      {verification && (
        <article className="cockpit-card">
          <div className="state-label">多源验证</div>
          <p>独立信源数量：{verification.independent_source_count}</p>
          <p>高等级独立信源数量：{verification.high_grade_independent_source_count}</p>
          <p>最高信源等级：{verification.highest_source_level}</p>
          <p>最新发布时间：{verification.latest_published_at || '暂无'}</p>
          <p>证据引用：{verification.evidence_ids.join('、') || '暂无'}</p>
        </article>
      )}
      <div className="ledger-surface">
        <EvidenceTable items={items} />
      </div>
    </div>
  )
}
