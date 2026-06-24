import { useEffect, useState } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { StatusNotice } from '../components/status/StatusNotice'
import { getDecisionLoop, listDecisionLoops } from '../services/decisionLoop'
import type { PageErrorState } from '../shared/utils'
import { toPageErrorState } from '../shared/utils'
import type { DecisionLoopItem, DecisionLoopLink } from '../types/decisionLoop'

export function DecisionLoopPage() {
  const [searchParams] = useSearchParams()
  const focusedDecisionId = searchParams.get('decision_id')?.trim() ?? ''
  const [items, setItems] = useState<DecisionLoopItem[]>([])
  const [safetyNote, setSafetyNote] = useState('只读解释链，仅展示本地事实和导航，不改变事实状态。')
  const [errorState, setErrorState] = useState<PageErrorState>()

  useEffect(() => {
    let mounted = true
    const request = focusedDecisionId
      ? getDecisionLoop(focusedDecisionId).then((res) => ({
        request_id: res.request_id,
        data: {
          items: res.data ? [res.data] : [],
          total: res.data ? 1 : 0,
          safety_note: res.data?.safety_note || '只读解释链，仅展示本地事实和导航，不改变事实状态。',
        },
      }))
      : listDecisionLoops({ limit: 20 })
    request
      .then((res) => {
        if (!mounted) return
        setItems(res.data?.items ?? [])
        setSafetyNote(res.data?.safety_note || '只读解释链，仅展示本地事实和导航，不改变事实状态。')
        setErrorState(undefined)
      })
      .catch((error: unknown) => {
        if (!mounted) return
        setItems([])
        setErrorState(toPageErrorState(error))
      })
    return () => {
      mounted = false
    }
  }, [focusedDecisionId])

  const safeItems = items.map(normalizeLoopItem)
  const incompleteCount = safeItems.filter((item) => item.loop_status === 'incomplete').length
  const latest = safeItems[0]

  return (
    <div>
      <h1 className="page-title">决策闭环解释</h1>
      <p className="page-placeholder">{safetyNote}</p>
      {errorState && <StatusNotice state={errorState.state} safeMessage={errorState.message} code={errorState.code} />}

      <section className="cockpit-grid" aria-label="决策闭环概览">
        <article className="cockpit-card">
          <div className="state-label">生命周期</div>
          <h2>只读决策生命周期</h2>
          <p>建议生成 -&gt; 用户确认 -&gt; 线下记录 -&gt; 风险/复盘 -&gt; 审计。</p>
          <p>本页只展示已有本地事实和缺口，不创建确认、不记录交易、不改变风险状态、不应用规则。</p>
        </article>

        <article className="cockpit-card">
          <div className="state-label">闭环概览</div>
          <h2>本地事实链</h2>
          <p>闭环条数：{safeItems.length}</p>
          <p>未闭合：{incompleteCount}</p>
          <p>{focusedDecisionId ? '当前聚焦' : '最近决策'}：{latest ? `${latest.decision_id}${latest.symbol ? ` · ${latest.symbol}` : ''}` : (focusedDecisionId || '暂无')}</p>
        </article>

        <article className="cockpit-card">
          <div className="state-label">查看入口</div>
          <h2>相关页面</h2>
          <div className="link-row">
            <Link to="/workbench">返回工作台</Link>
            {latest ? <Link to={`/decisions/${latest.decision_id}`}>查看决策详情</Link> : null}
            <Link to="/review">查看复盘摘要</Link>
            <Link to="/audit">查看审计</Link>
          </div>
        </article>
      </section>

      {safeItems.length === 0 && !errorState ? (
        <article className="cockpit-card">
          <div className="state-label">空态</div>
          <h2>{focusedDecisionId ? '未找到目标决策闭环' : '暂无决策闭环记录'}</h2>
          <p>{focusedDecisionId ? '请回到决策详情或闭环列表确认该决策是否已经生成本地记录。' : '完成本地决策记录后，这里会展示建议、用户记录、线下事实与追踪线索。'}</p>
        </article>
      ) : null}

      <section className="stacked-panel" aria-label="决策闭环列表">
        {safeItems.map((item) => (
          <article className="cockpit-card ledger-surface" key={item.decision_id}>
            <div className="state-label">{loopStatusText(item.loop_status)}</div>
            <h2>{item.symbol ? `${item.decision_id} · ${item.symbol}` : item.decision_id}</h2>
            <p>{item.final_verdict_text}</p>
            <p>裁决状态：{item.final_verdict_status}；确认状态：{item.confirmation_status}</p>
            <p>生成时间：{item.generated_at || '暂无'}</p>

            <section aria-label={`${item.decision_id} 阶段链路`}>
              <h3>阶段链路</h3>
              <ul className="trace-list">
                {item.stages.map((stage) => (
                  <li key={`${item.decision_id}-${stage.stage}`}>
                    <strong>{stage.label}</strong>
                    <span> · {stageStatusText(stage.status)}</span>
                    <br />
                    <span>{stage.summary}</span>
                    {stage.ref_id && <small> · {stage.ref_type || 'ref'}:{stage.ref_id}</small>}
                  </li>
                ))}
              </ul>
            </section>

            <section aria-label={`${item.decision_id} 缺口`}>
              <h3>缺口</h3>
              {item.missing_links.length ? (
                <ul>
                  {item.missing_links.map((gap) => <li key={`${item.decision_id}-${gap}`}>{gap}</li>)}
                </ul>
              ) : (
                <p>暂无缺口。</p>
              )}
            </section>

            <section aria-label={`${item.decision_id} 人工记录`}>
              <h3>人工记录</h3>
              {item.manual_actions.length ? (
                <ul>
                  {item.manual_actions.map((action) => (
                    <li key={action.confirmation_id}>
                      <strong>{action.confirmation_id}</strong>
                      <br />
                      <span>{action.confirmation_type}{action.operation_type ? ` · ${action.operation_type}` : ''}{action.symbol ? ` · ${action.symbol}` : ''}</span>
                      <br />
                      <span>数量：{valueOrDash(action.quantity)}；价格：{valueOrDash(action.price)}；费用：{valueOrDash(action.fees)}</span>
                      <br />
                      <span>流水：{action.transaction_ids.length ? action.transaction_ids.join('、') : '暂无'}</span>
                      {action.note_preview && (
                        <>
                          <br />
                          <span>{action.note_preview}</span>
                        </>
                      )}
                    </li>
                  ))}
                </ul>
              ) : (
                <p>暂无人工处理记录。</p>
              )}
            </section>

            <section aria-label={`${item.decision_id} 追踪链接`}>
              <h3>追踪链接</h3>
              {item.risk_links.length || item.review_links.length || item.audit_links.length ? (
                <div className="link-row">
                  {renderLinks([...item.risk_links, ...item.review_links, ...item.audit_links])}
                </div>
              ) : (
                <p>暂无风险、复盘或审计链接。</p>
              )}
            </section>
          </article>
        ))}
      </section>
    </div>
  )
}

function renderLinks(links: DecisionLoopLink[]) {
  return links.map((link) => (
    <Link key={`${link.type}-${link.id}`} to={link.href}>
      {link.label}
    </Link>
  ))
}

function normalizeLoopItem(item: DecisionLoopItem): DecisionLoopItem {
  return {
    ...item,
    stages: item.stages ?? [],
    manual_actions: item.manual_actions ?? [],
    risk_links: item.risk_links ?? [],
    review_links: item.review_links ?? [],
    audit_links: item.audit_links ?? [],
    missing_links: item.missing_links ?? [],
  }
}

function loopStatusText(status: string) {
  const map: Record<string, string> = {
    open: '观察中',
    planned: '已记录计划',
    recorded: '已有本地记录',
    reviewed: '已有追踪线索',
    incomplete: '存在缺口',
  }
  return map[status] ?? '已记录'
}

function stageStatusText(status: string) {
  const map: Record<string, string> = {
    complete: '完成',
    pending: '待补充',
    not_required: '无需补充',
    missing: '缺失',
    degraded: '降级',
  }
  return map[status] ?? status
}

function valueOrDash(value?: number) {
  if (value === undefined || value === null || Number.isNaN(value)) {
    return '暂无'
  }
  return value
}
