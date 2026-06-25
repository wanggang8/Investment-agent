import { useEffect, useState } from 'react'
import { RuleProposalPanel } from '../components/rules/RuleProposalPanel'
import { Button, SummaryCard, type UITone } from '../components/ui'
import { APIClientError } from '../services/client'
import { confirmRuleProposal, createSOPAddendumProposal, finalConfirmRuleProposal, getCurrentRule, listRuleProposals } from '../services/rule'
import type { RuleProposal, RuleVersion } from '../types/rule'
import { buildRulesGovernanceModel } from '../features/governance'

export function RulesPage() {
  const [proposals, setProposals] = useState<RuleProposal[]>([])
  const [currentRule, setCurrentRule] = useState<RuleVersion>()
  const [message, setMessage] = useState('')

  function refresh(clearMessage = true) {
    listRuleProposals()
      .then((res) => {
        setProposals(res.data?.items ?? [])
        if (clearMessage) {
          setMessage('')
        }
      })
      .catch((error: unknown) => setMessage(error instanceof APIClientError ? error.message : '暂时无法读取规则提案。'))
  }

  useEffect(() => {
    refresh()
    getCurrentRule()
      .then((res) => setCurrentRule(res.data))
      .catch(() => setCurrentRule(undefined))
  }, [])

  function handleConfirm(proposalId: string, confirm: boolean) {
    confirmRuleProposal(proposalId, { confirm })
      .then(() => {
        setMessage(confirm ? '已确认送审。' : '已拒绝提案。')
        refresh(false)
      })
      .catch((error: unknown) => setMessage(error instanceof APIClientError ? error.message : '提案确认提交失败。'))
  }

  function handleFinalConfirm(proposalId: string, confirm: boolean) {
    finalConfirmRuleProposal(proposalId, { confirm })
      .then(() => {
        setMessage(confirm ? '已提交最终确认。' : '已拒绝应用。')
        refresh(false)
      })
      .catch((error: unknown) => setMessage(error instanceof APIClientError ? error.message : '最终确认提交失败。'))
  }

  function handleCreateSOPAddendum() {
    createSOPAddendumProposal({
      scenario_key: 'p88_uncovered_liquidity_gap',
      scenario_title: '连续流动性缺口未覆盖',
      occurrence_count: 4,
      sample_window: '2026-Q2',
    })
      .then(() => {
        setMessage('SOP 补充提案已生成，等待人工确认。')
        refresh(false)
      })
      .catch((error: unknown) => setMessage(error instanceof APIClientError ? error.message : 'SOP 补充提案生成失败。'))
  }

  const ruleObject = currentRule?.rules && typeof currentRule.rules === 'object' ? currentRule.rules as { priority?: string[], thresholds?: unknown } : undefined
  const governanceModel = buildRulesGovernanceModel({ currentRule, proposals })

  return (
    <div>
      <h1 className="page-title">规则与纪律</h1>
      {message && <div className="page-placeholder">{message}</div>}
      <section className={`daily-hero daily-tone-${governanceModel.overallTone}`} aria-label="规则治理总览">
        <div className="daily-hero-main">
          <div className="state-label">规则治理状态</div>
          <h2>{governanceModel.overallLabel}</h2>
          <p>{governanceModel.safetyNotes[0]}</p>
          <div className="daily-signal-grid quality-signal-grid">
            {governanceModel.metrics.map((metric) => (
              <SummaryCard key={metric.label} title={metric.label} value={metric.value} detail={metric.detail} tone={(metric.tone ?? 'unknown') as UITone} />
            ))}
          </div>
        </div>
        <aside className="daily-hero-side" aria-label="规则治理下一步">
          <strong>下一步人工治理</strong>
          <ul>
            {governanceModel.nextActions.map((action) => (
              <li key={action.label}>
                <a href={action.href} aria-label={`${action.label}入口`}>{action.label}</a>
                <span>{action.detail}</span>
              </li>
            ))}
          </ul>
        </aside>
      </section>
      {currentRule && (
        <article className="cockpit-card">
          <div className="state-label">当前规则库</div>
          <p>当前规则库：{currentRule.rule_version}</p>
          <p>状态：{currentRule.status}</p>
          <p>裁决优先级：{ruleObject?.priority?.join('、') || '暂无'}</p>
          <details className="raw-detail">
            <summary>查看规则阈值</summary>
            <pre aria-label="规则阈值">{JSON.stringify(ruleObject?.thresholds ?? currentRule.rules, null, 2)}</pre>
          </details>
        </article>
      )}
      <article className="cockpit-card">
        <div className="state-label">SOP 补充提案</div>
        <p>高频未覆盖场景只生成待确认提案；规则应用仍需后续人工确认与守门人审计。</p>
        <Button onClick={handleCreateSOPAddendum}>生成 SOP 补充提案</Button>
      </article>
      <RuleProposalPanel proposals={proposals} onConfirm={handleConfirm} onFinalConfirm={handleFinalConfirm} />
    </div>
  )
}
