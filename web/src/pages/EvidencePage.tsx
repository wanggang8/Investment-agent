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
      <section className="cockpit-grid" aria-label="证据解释入口">
        <article className="cockpit-card">
          <div className="state-label">证据可信度</div>
          <h2>证据可信度</h2>
          <p>先看信源等级、独立信源和核验状态，再进入证据明细；C/background 只作为背景材料。</p>
          <p>独立信源数量：{verification?.independent_source_count ?? 0}</p>
          <p>最高信源等级：{verification?.highest_source_level || '暂无'}</p>
          <p>最新发布时间：{verification?.latest_published_at || '暂无'}</p>
        </article>
        <article className="cockpit-card">
          <div className="state-label">解释路径</div>
          <h2>决策解释入口</h2>
          <p>S/A/B 级可作为正式证据；C 级和 background 只辅助理解，不替代规则裁决或人工复核。</p>
          <div className="link-row">
            <a href="/workbench">返回工作台</a>
            <a href="/decision-loop">查看决策闭环</a>
            <a href="/audit">查看审计</a>
          </div>
        </article>
      </section>
      <article className="cockpit-card">
        <div className="state-label">证据维护</div>
        <h2>情报刷新与索引重建</h2>
        <p>只更新本地情报、RAG 文本块、VecLite 索引和审计记录。</p>
        <button type="button" onClick={handleRefreshEvidence}>刷新情报</button>
        <button type="button" onClick={handleRebuildIndex}>重建索引</button>
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
