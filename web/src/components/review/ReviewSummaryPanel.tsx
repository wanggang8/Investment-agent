import type { ReviewOpsStatus, ReviewSummary } from '../../types/review'
import { opsStatusText, ruleProposalStatusText, textOrRaw } from '../../shared/mappers/statusText'

interface Props {
  summary?: ReviewSummary
}

export function ReviewSummaryPanel({ summary }: Props) {
  const periodText = summary?.period === 'quarterly' ? '季度复盘' : '月度复盘'
  const suggestions = summary?.rule_suggestions ?? []
  const effectTracking = summary?.rule_effect_tracking ?? []
  const trackingLinks = summary?.tracking_links ?? []

  return (
    <article className="cockpit-card">
      <div className="state-label">复盘摘要</div>
      <h2>{periodText}</h2>
      <div className="metric-grid">
        <div><span>建议数量</span><strong>{summary?.decision_count ?? summary?.recent_decisions?.length ?? 0}</strong></div>
        <div><span>确认动作</span><strong>{summary?.confirmation_count ?? 0}</strong></div>
        <div><span>已手动执行</span><strong>{summary?.executed_manually_count ?? 0}</strong></div>
        <div><span>记录计划</span><strong>{summary?.planned_count ?? 0}</strong></div>
        <div><span>错误案例</span><strong>{summary?.error_case_count ?? 0}</strong></div>
        <div><span>规则提案</span><strong>{summary?.rule_proposal_count ?? 0}</strong></div>
        <div><span>审计事件</span><strong>{summary?.audit_event_count ?? 0}</strong></div>
        <div><span>规则命中</span><strong>{summary?.rule_hit_count ?? 0}</strong></div>
        <div><span>误判</span><strong>{summary?.misjudgment_count ?? 0}</strong></div>
        <div><span>缺证据</span><strong>{summary?.missing_evidence_count ?? 0}</strong></div>
        <div><span>降级</span><strong>{summary?.degraded_count ?? 0}</strong></div>
      </div>

      <ReviewOpsStatusPanel status={summary?.ops_status} />

      <section className="detail-section">
        <div className="state-label">归因摘要</div>
        {(summary?.attribution_summaries?.length ?? 0) === 0 ? <p>暂无归因摘要。</p> : summary?.attribution_summaries?.map((item) => (
          <div className="detail-row" key={item.decision_id}>
            <strong>{item.decision_id}</strong>
            <span>{item.symbol || '未知标的'}</span>
            <span>{textOrRaw(reviewOutcomeText, item.outcome)}</span>
            <span>{item.evidence_status || 'unknown'}</span>
          </div>
        ))}
      </section>

      <section className="detail-section">
        <div className="state-label">高频错误标签</div>
        {(summary?.recurring_error_tags?.length ?? 0) === 0 ? <p>暂无错误标签。</p> : summary?.recurring_error_tags?.map((item) => (
          <div className="detail-row" key={item.tag}>{item.tag} · {item.count}</div>
        ))}
      </section>

      <section className="detail-section">
        <div className="state-label">缺证据主题</div>
        {(summary?.missing_evidence_themes?.length ?? 0) === 0 ? <p>暂无缺证据主题。</p> : summary?.missing_evidence_themes?.map((item) => (
          <div className="detail-row" key={item.status}>{item.status} · {item.count}</div>
        ))}
      </section>

      <section className="detail-section">
        <div className="state-label">规则提案结果</div>
        {(summary?.rule_proposal_outcomes?.length ?? 0) === 0 ? <p>暂无规则提案结果。</p> : summary?.rule_proposal_outcomes?.map((item) => (
          <div className="detail-row" key={item.proposal_id}>
            <strong>{item.title}</strong>
            <span>{textOrRaw(ruleProposalStatusText, item.status)}</span>
            {item.audit_result && <span>{item.audit_result}</span>}
          </div>
        ))}
      </section>

      <section className="detail-section">
        <div className="state-label">降级工作流</div>
        {(summary?.degraded_workflows?.length ?? 0) === 0 ? <p>暂无降级工作流。</p> : summary?.degraded_workflows?.map((item) => (
          <div className="detail-row" key={item.decision_id}>
            <strong>{item.decision_id}{item.symbol ? ` · ${item.symbol}` : ''}</strong>
            <span>{textOrRaw(opsStatusText, item.status)}</span>
            <span>{item.created_at}</span>
          </div>
        ))}
      </section>

      <section className="detail-section">
        <div className="state-label">规则应用后效果追踪</div>
        {effectTracking.length === 0 ? <p>暂无规则效果追踪。</p> : effectTracking.map((item) => (
          <div className="detail-row" key={item.tracking_id}>
            <strong>{item.applied_rule_version} · 趋势：{trendText[item.trend_direction] ?? '未知'}</strong>
            <span>命中 {item.hit_count} · 误判 {item.misjudgment_count} · 缺证据 {item.missing_evidence_count} · 降级 {item.degraded_count} · 风险预警 {item.risk_alert_count}</span>
            {item.metrics !== undefined && <pre aria-label="追踪指标">{JSON.stringify(item.metrics, null, 2)}</pre>}
            {item.related_proposal_ids !== undefined && <pre aria-label="追踪关联提案">{JSON.stringify(item.related_proposal_ids, null, 2)}</pre>}
            {item.related_audit_event_ids !== undefined && <pre aria-label="追踪关联审计">{JSON.stringify(item.related_audit_event_ids, null, 2)}</pre>}
            {item.related_risk_alert_ids !== undefined && <pre aria-label="追踪关联风险预警">{JSON.stringify(item.related_risk_alert_ids, null, 2)}</pre>}
            {item.safety_note && <span>{item.safety_note}</span>}
          </div>
        ))}
      </section>

      <section className="detail-section">
        <div className="state-label">规则建议</div>
        <p>规则变更仍需守门人审计和用户最终确认，不会自动应用。</p>
        {suggestions.length === 0 ? <p>暂无规则建议。</p> : suggestions.map((item) => (
          <div className="detail-row" key={item.proposal_id}>
            <strong>{item.title}</strong>
            <span>{textOrRaw(ruleProposalStatusText, item.status)}</span>
            {item.reason && <span>{item.reason}</span>}
          </div>
        ))}
      </section>

      <section className="detail-section">
        <div className="state-label">追踪入口</div>
        {trackingLinks.length === 0 ? <p>暂无追踪记录。</p> : trackingLinks.map((item) => (
          <a className="detail-row" key={`${item.type}-${item.id}`} href={`#${item.type}-${item.id}`}>
            <strong>{item.label}</strong>
            <span>{item.id}</span>
          </a>
        ))}
      </section>
    </article>
  )
}

const reviewOutcomeText: Record<string, string> = {
  missing_evidence: '缺失证据',
  degraded: '降级',
  executed_manually: '已手动执行',
  planned: '记录计划',
  recorded: '已记录',
}

const trendText: Record<string, string> = {
  improved: '改善',
  flat: '持平',
  worsened: '变差',
  unknown: '未知',
}

function ReviewOpsStatusPanel({ status }: { status?: ReviewOpsStatus }) {
  return (
    <section className="detail-section">
      <div className="state-label">运维状态</div>
      {!status ? (
        <p>暂无运维状态数据。</p>
      ) : (
        <>
          <div className="metric-grid">
            <div><span>数据源</span><strong>{textOrRaw(opsStatusText, status.data_source_status)}</strong></div>
            <div><span>索引</span><strong>{textOrRaw(opsStatusText, status.index_status)}</strong></div>
            <div><span>复盘状态</span><strong>{textOrRaw(opsStatusText, status.review_status)}</strong></div>
          </div>
          {status.explanation && <p>{status.explanation}</p>}
        </>
      )}
      <p>仅展示状态与追踪入口，不执行交易，也不自动应用规则。</p>
    </section>
  )
}
