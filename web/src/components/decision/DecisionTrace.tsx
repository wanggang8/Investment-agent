import { useState } from 'react'
import type { ConfirmationRequest, DecisionDetailResponse, ExpectedReturnScenarios } from '../../types/decision'
import { capabilityStatusText, confidenceText, evidenceRoleText, precisionStatusText, retrievalFallbackSourceText, retrievalQualityStatusText, scenarioText, sellEvaluationStatusText, severityText, sourceConsistencyStatusText, sourceHealthStatusText, textOrRaw, verdictStatusText, verificationStatusText, workflowStatusText } from '../../shared/mappers'
import { formatPercent } from '../../shared/utils'
import { buildDecisionExplanationModel } from '../../features/decision/decisionExplanationModel'
import { UserConfirmationPanel } from '../dashboard/UserConfirmationPanel'

interface Props {
  decision?: DecisionDetailResponse
  onConfirm?: (decisionId: string, payload: ConfirmationRequest) => void
}

function probabilityText(scenarios: ExpectedReturnScenarios, probability?: number | null) {
  if (scenarios.precision_status !== 'available') {
    return '样本不足，不展示精确概率'
  }
  return probability == null ? '无精确概率' : `${(probability * 100).toFixed(1)}%`
}

function percentText(value?: number | null) {
  return value == null ? '暂无' : `${(value * 100).toFixed(1)}%`
}

function textList(value: unknown) {
  return Array.isArray(value) && value.length ? value.filter(Boolean).join('、') : '暂无'
}

function productFallbackSourceText(value: string) {
  return value.replace(/VecLite 索引/g, '检索索引').replace(/RAG/g, '检索')
}

function productAnalysisText(value: string) {
  return value
    .replace(/暂无可展示的 LLM 分析材料/g, '暂无可展示的分析材料')
    .replace(/LLM 材料/g, '分析材料')
    .replace(/LLM 分析材料/g, '分析材料')
    .replace(/LLM 材料只作为分析输入/g, '分析材料只作为分析输入')
}

function holdingClassText(value?: string) {
  const mapped: Record<string, string> = {
    broad_index_etf: '宽基指数 ETF',
    sector_growth_fund: '行业/成长基金',
    equity_constituent_financial: '金融成分股路径',
  }
  return value ? (mapped[value] || value) : '暂无'
}

function safeArray<T>(value: T[] | null | undefined): T[] {
  return Array.isArray(value) ? value : []
}

function hasKnowledgeReadinessContext(value?: string | null) {
  if (!value) {
    return false
  }
  return value.includes('principles=') || value.includes('data_readiness=') || value.includes('master.')
}

function safePromptVersion(value?: string | null) {
  if (!value) {
    return '已记录'
  }
  return /^[A-Za-z0-9._-]+$/.test(value) ? value : '已记录'
}

export function DecisionTrace({ decision, onConfirm }: Props) {
  const [showEvidence, setShowEvidence] = useState(true)
  const [showAnalysis, setShowAnalysis] = useState(false)

  if (!decision) {
    return <div className="page-placeholder">选择一条建议后展示完整裁决链路。</div>
  }

  const explanation = buildDecisionExplanationModel(decision)
  const finalVerdict = decision.final_verdict ?? { display_text: '暂无明确裁决', status: 'insufficient_data' }
  const triggeredRules = safeArray(decision.triggered_rules)
  const evidenceChain = safeArray(decision.evidence_chain)
  const analystReports = safeArray(decision.analyst_reports)
  const expectedReturn = decision.expected_return_scenarios
  const expectedReturnScenarios = safeArray(expectedReturn?.scenarios)
  const marketContext = decision.market_context
  const arbitrationChain = safeArray(decision.arbitration_chain)
  const confirmation = decision.user_confirmation ?? { confirmation_status: 'not_required', available_actions: [] }
  const confirmationActions = safeArray(confirmation.available_actions)

  return (
    <div className="stacked-panel">
      <section className="decision-story" aria-label="决策解释故事">
        <article className="decision-story-hero">
          <div className="state-label">决策故事</div>
          <h2>{explanation.storyTitle}</h2>
          <dl className="daily-hero-meta">
            {explanation.decisionContext.map((item, index) => {
              const [label, ...rest] = item.split('：')
              const text = rest.length ? rest.join('：') : item
              return (
                <div key={`context-${index}-${item}`}>
                  <dt>{rest.length ? label : '上下文'}</dt>
                  <dd>{text}</dd>
                </div>
              )
            })}
          </dl>
          <div className="link-row">
            {explanation.explanationLinks.map((link) => (
              <a key={link.href} href={link.href}>{link.label}</a>
            ))}
          </div>
        </article>

        <section className="decision-story-grid" aria-label="决策解释摘要">
          <article className="cockpit-card state-frozen-watch">
            <div className="state-label">人工边界</div>
            <h2>安全边界</h2>
            <p>{explanation.safetyNotes[0]}</p>
            <section className="verdict-section">
              <h3>禁止动作</h3>
              <ul>
                {explanation.prohibitedActions.map((action, index) => <li key={`prohibited-${index}-${action}`}>{action}</li>)}
              </ul>
            </section>
            <section className="verdict-section">
              <h3>可选人工动作</h3>
              <ul>
                {explanation.optionalActions.map((action, index) => <li key={`optional-${index}-${action}`}>{action}</li>)}
              </ul>
            </section>
          </article>

          <article className="cockpit-card">
            <div className="state-label">为什么</div>
            <h2>关键原因</h2>
            <ul>
              {explanation.keyReasons.map((reason, index) => <li key={`reason-${index}-${reason}`}>{reason}</li>)}
            </ul>
            {explanation.missingDataWarnings.length ? (
              <ul className="quality-list">
                {explanation.missingDataWarnings.map((warning, index) => (
                  <li key={`warning-${index}-${warning}`}><strong>需复核</strong><span>{productAnalysisText(warning)}</span></li>
                ))}
              </ul>
            ) : null}
          </article>

          <article className="cockpit-card">
            <div className="state-label">可信度</div>
            <h2>可信度</h2>
            <ul>
              {explanation.trustSummary.map((item, index) => <li key={`trust-${index}-${item}`}>{productAnalysisText(item)}</li>)}
            </ul>
            <p className="muted-text">{productAnalysisText(explanation.safetyNotes[1])}</p>
          </article>
        </section>
      </section>

      <h2 className="section-title">技术追踪</h2>
      <article className="cockpit-card">
        <div className="state-label">最终裁决明细</div>
        <h2>裁决元数据</h2>
        <p>最终裁决：{finalVerdict.display_text || '暂无明确裁决'}</p>
        <p>状态：{textOrRaw(verdictStatusText, finalVerdict.status)}</p>
        <p>建议编号：{decision.decision_id}</p>
        <p>生成时间：{decision.generated_at || '暂无'}</p>
        <p>问题：{decision.question || '无'}</p>
        <p>标的：{decision.symbol || '无'}</p>
        <p>工作流状态：{textOrRaw(workflowStatusText, decision.workflow_status)}</p>
      </article>

      <article className="cockpit-card">
        <div className="state-label">摘要</div>
        <p>能力圈检查：{textOrRaw(capabilityStatusText, decision.capability_check?.status)} {decision.capability_check?.reason ?? ''}</p>
        <p>禁止事项：{textList(finalVerdict.prohibited_actions)}</p>
        <p>可选动作：{textList(finalVerdict.optional_actions)}</p>
      </article>

      <article className="cockpit-card">
        <div className="state-label">账户与规则</div>
        <p>账户快照：{decision.account_snapshot?.snapshot_id ?? '暂无'}</p>
        <p>现金：{decision.account_snapshot?.cash ?? '暂无'}</p>
        <p>总资产：{decision.account_snapshot?.total_assets ?? '暂无'}</p>
        <p>现金比例：{decision.account_snapshot?.cash_ratio ?? '暂无'}</p>
        <p>高风险比例：{decision.account_snapshot?.high_risk_ratio ?? '暂无'}</p>
        <ul>
          {triggeredRules.length ? triggeredRules.map((rule) => (
            <li key={rule.rule_id}>{rule.rule_name}：{textOrRaw(severityText, rule.severity)} / {rule.description}</li>
          )) : <li>暂无触发规则。</li>}
        </ul>
      </article>

      <article className="cockpit-card">
        <div className="state-label">证据链</div>
        <button className="link-button" type="button" onClick={() => setShowEvidence(!showEvidence)}>{showEvidence ? '收起' : '展开'}证据链</button>
        {showEvidence && (
          <ul>
            {evidenceChain.length ? evidenceChain.map((item) => (
              <li key={item.evidence_id}>
                {item.source_name} / {item.source_level} 级 / {textOrRaw(evidenceRoleText, item.evidence_role ?? 'formal')} / {textOrRaw(verificationStatusText, item.verification_status)}：{item.summary}
              </li>
            )) : <li>暂无证据链。</li>}
          </ul>
        )}
      </article>

      {decision.retrieval_quality && (
        <article className="cockpit-card">
          <div className="state-label">证据检索</div>
          <h3>检索质量</h3>
          <p>检索状态：{textOrRaw(retrievalQualityStatusText, decision.retrieval_quality.status)}</p>
          <p>查询摘要：{decision.retrieval_quality.query_summary || '暂无'}</p>
          <p>召回数量：{decision.retrieval_quality.top_k}</p>
          <p>索引健康：{textOrRaw(sourceHealthStatusText, decision.retrieval_quality.index_health)}</p>
          <p>索引新鲜度：{textOrRaw(sourceHealthStatusText, decision.retrieval_quality.index_freshness)}</p>
          <p>备用来源：{productFallbackSourceText(textOrRaw(retrievalFallbackSourceText, decision.retrieval_quality.fallback_source))}</p>
          <span className="reference-sr-only">Fallback 来源：{textOrRaw(retrievalFallbackSourceText, decision.retrieval_quality.fallback_source)}</span>
          <p>引用一致性：{textOrRaw(sourceConsistencyStatusText, decision.retrieval_quality.source_consistency_status)}</p>
          <p>降级原因：{decision.retrieval_quality.degraded_reason || '无'}</p>
          {decision.retrieval_quality.status === 'empty' && (
            <p>未召回可用于裁决的证据，最终判断仍按规则链与证据守门条件处理。</p>
          )}
          {decision.retrieval_quality.status === 'degraded' && (
            <p>可在证据页重建索引后再次复核；本页不会读取本地索引文件或自动修复。</p>
          )}
        </article>
      )}

      <article className="cockpit-card">
        <div className="state-label">Agent 分析材料</div>
        <p>以下 {analystReports.length} 份内容仅作为分析材料，最终裁决仍以规则链为准；默认收起以便优先复核裁决、安全边界和人工动作。</p>
        {analystReports.length ? (
          <button className="link-button" type="button" onClick={() => setShowAnalysis(!showAnalysis)}>{showAnalysis ? '收起' : `展开 ${analystReports.length} 份`}分析材料</button>
        ) : <p>暂无分析材料。</p>}
        {showAnalysis && (
          analystReports.length ? analystReports.map((report) => (
            <section key={report.agent_name} className="proposal-item">
              <h3>{report.agent_name}</h3>
              <p>{report.conclusion}</p>
              <p>关键理由：{textList(report.key_reasons)}</p>
              <p>风险提示：{textList(report.risk_warnings)}</p>
              <p>证据引用：{textList(report.evidence_ids)}</p>
              <p>置信度：{textOrRaw(confidenceText, report.confidence)}</p>
              {hasKnowledgeReadinessContext(report.input_summary) && (
                <>
                  <p>分析模型已参考知识与数据准备度摘要</p>
                  <span className="reference-sr-only">LLM 已参考知识与数据准备度摘要</span>
                  <p>prompt {safePromptVersion(report.prompt_version)}；仅展示脱敏摘要。</p>
                </>
              )}
            </section>
          )) : <p>暂无分析材料。</p>
        )}
      </article>

      {expectedReturn && (
        <article className="cockpit-card">
          <div className="state-label">预期收益情景</div>
          <p>标的：{marketContext?.symbol || decision.symbol || '暂无'}</p>
          <p>标的名称：{expectedReturn.target_name || '暂无'}{expectedReturn.target_code ? `（${expectedReturn.target_code}）` : ''}</p>
          <p>持仓类别：{holdingClassText(expectedReturn.holding_class)}</p>
          <p>分析期限：{expectedReturn.horizon_label || '暂无'}</p>
          <p>当前日期：{marketContext?.trade_date || decision.generated_at || '暂无'}</p>
          <p>当前价格或净值：{marketContext?.current_price ?? '暂无'}</p>
          <p>PE/PB 分位：{marketContext?.pe_percentile ?? '暂无'} / {marketContext?.pb_percentile ?? '暂无'}</p>
          <p>精度状态：{textOrRaw(precisionStatusText, expectedReturn.precision_status)}</p>
          <p>概率依据：{expectedReturn.probability_basis || '暂无'}</p>
          <p>样本数：{expectedReturn.sample_count ?? 0}</p>
          <p>样本窗口：{expectedReturn.sample_window || '暂无'}</p>
          <p>筛选条件：{expectedReturn.screening_condition || '暂无'}</p>
          <p>支持数据：{expectedReturn.supporting_data_summary || '暂无'}</p>
          <p>缺口数据：{textList(expectedReturn.missing_categories)}</p>
          <p>需补充数据：{textList(expectedReturn.supplement_data)}</p>
          <p>{expectedReturn.reason || '暂无说明'}</p>
          <p>{expectedReturn.disclaimer || '预期收益仅为情景分析，不构成收益承诺。'}</p>
          <ul>
            {expectedReturnScenarios.length ? expectedReturnScenarios.map((scenario, index) => (
              <li key={`${index}-${scenario.scenario || 'scenario'}`}>
                {scenarioText(scenario.scenario)}：{scenario.return_range}，概率 {probabilityText(expectedReturn, scenario.probability)}{scenario.trigger ? `，触发条件：${scenario.trigger}` : ''}
              </li>
            )) : <li>暂无预期收益情景。</li>}
          </ul>
          {safeArray(expectedReturn.holding_class_coverage).length ? (
            <section className="proposal-item">
              <h3>持仓类别覆盖</h3>
              <ul>
                {safeArray(expectedReturn.holding_class_coverage).map((item, index) => (
                  <li key={`${index}-${item.holding_class || 'holding'}-${item.symbol || 'symbol'}`}>{holdingClassText(item.holding_class)}：{item.symbol}，{item.status}</li>
                ))}
              </ul>
            </section>
          ) : null}
          {safeArray(expectedReturn.assumption_checks).length ? (
            <section className="proposal-item">
              <h3>假设监控</h3>
              <ul>
                {safeArray(expectedReturn.assumption_checks).map((item, index) => (
                  <li key={`${index}-${item.name || 'assumption'}-${item.months_below ?? 'months'}`}>{item.name}：预期 {formatPercent(item.expected)}，实际 {formatPercent(item.actual)}，低于预期 {item.months_below} 个月</li>
                ))}
              </ul>
            </section>
          ) : null}
          {safeArray(expectedReturn.historical_contexts).length ? (
            <section className="proposal-item">
              <h3>历史相似场景</h3>
              <ul>
                {safeArray(expectedReturn.historical_contexts).map((item, index) => (
                  <li key={`${index}-${item.label || 'history'}-${item.window || 'window'}`}>
                    {item.label}：{item.window}，样本 {item.sample_count}，结果：{item.outcome}，最大回撤 {percentText(item.max_drawdown)}，修复：{item.recovery}，来源：{item.source || '暂无'}
                  </li>
                ))}
              </ul>
            </section>
          ) : null}
          {expectedReturn.sell_evaluation && (
            <section className="proposal-item">
              <h3>动态卖出评估</h3>
              <p>状态：{textOrRaw(sellEvaluationStatusText, expectedReturn.sell_evaluation.status)}</p>
              <p>触发因素：{expectedReturn.sell_evaluation.triggers?.join('、') || '暂无'}</p>
              <p>人工提示：{expectedReturn.sell_evaluation.prompts?.join('、') || '暂无'}</p>
              <p>建议动作：{expectedReturn.sell_evaluation.actions?.join('、') || '暂无'}</p>
              <p>{expectedReturn.sell_evaluation.non_trading_disclaimer || '卖出评估仅用于人工复核，不会自动执行交易。'}</p>
            </section>
          )}
          {expectedReturn.reassessment_trigger && (
            <section className="proposal-item">
              <h3>复核触发</h3>
              <p>原因：{expectedReturn.reassessment_trigger.reason}</p>
              <p>边界：{expectedReturn.reassessment_trigger.boundary || '暂无'}</p>
              <p>当前值：{expectedReturn.reassessment_trigger.current_value == null ? '暂无' : formatPercent(expectedReturn.reassessment_trigger.current_value)}</p>
            </section>
          )}
        </article>
      )}

      <article className="cockpit-card">
        <div className="state-label">裁决链</div>
        <ol>
          {arbitrationChain.length ? arbitrationChain.map((step, index) => (
            <li key={`${index}-${step.priority ?? 'priority'}-${step.rule_id || 'rule'}`}>{step.priority}. {step.rule_id}：{step.result}</li>
          )) : <li>暂无裁决链。</li>}
        </ol>
      </article>

      <article className="cockpit-card">
        <div className="state-label">审计时间线</div>
        <p>关联决策 ID：{decision.decision_id}</p>
        {decision.audit_events?.length ? (
          <ol>
            {decision.audit_events.map((event) => (
              <li key={event.audit_event_id}>
                {event.created_at || '暂无时间'} / {event.node_name || event.action} / {event.status}{event.error_code ? ` / ${event.error_code}` : ''}
              </li>
            ))}
          </ol>
        ) : (
          <p>暂无内嵌审计事件，可在“复盘与审计”页按决策编号查询。</p>
        )}
      </article>

      <UserConfirmationPanel
        decisionId={decision.decision_id}
        availableActions={confirmationActions}
        confirmationStatus={confirmation.confirmation_status}
        onSubmit={onConfirm}
      />
    </div>
  )
}
