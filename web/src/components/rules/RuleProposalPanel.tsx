import type { RuleProposal } from '../../types/rule'

const statusText: Record<string, string> = {
  draft: '草稿',
  pending_user_confirm: '待用户确认',
  under_gatekeeper_audit: '守门人审计中',
  pending_final_confirm: '待最终确认',
  rejected: '已拒绝',
  applied: '已应用',
}

const auditResultText: Record<string, string> = {
  approved: '审计通过',
  rejected: '审计否决',
  needs_user_review: '需要用户复核',
}

const effectValidationStatusText: Record<string, string> = {
  not_evaluated: '未评估',
  insufficient: '样本不足',
  passed: '已通过',
  failed: '未通过',
  needs_more_samples: '需要更多样本',
  needs_user_review: '需要用户复核',
}

const overfitRiskText: Record<string, string> = {
  low: '低',
  medium: '中',
  high: '高',
}

const replayResultText: Record<string, string> = {
  passed: '通过',
  failed: '不利',
  mixed: '混合',
  unknown: '未知',
}

const guardrailDecisionText: Record<string, string> = {
  passed: '通过',
  rejected: '拒绝',
  needs_user_review: '需要用户复核',
}

interface Props {
  proposals: RuleProposal[]
  onConfirm?: (proposalId: string, confirm: boolean) => void
  onFinalConfirm?: (proposalId: string, confirm: boolean) => void
}

export function RuleProposalPanel({ proposals, onConfirm, onFinalConfirm }: Props) {
  return (
    <article className="cockpit-card">
      <div className="state-label">规则提案</div>
      {proposals.length === 0 ? (
        <p>暂无规则提案。</p>
      ) : (
        proposals.map((proposal) => (
          <section key={proposal.proposal_id} className="proposal-item">
            <h3>{proposal.title}</h3>
            <p>状态：{statusText[proposal.status] ?? '未知状态'}；样本数：{proposal.sample_count}</p>
            {proposal.source_error_case_id && <p>来源误判案例：{proposal.source_error_case_id}</p>}
            {proposal.reason && <p>提案理由：{proposal.reason}</p>}
            {proposal.audit_result && <p>守门人结果：{auditResultText[proposal.audit_result] ?? '未知状态'}</p>}
            {proposal.audit_summary && <p>审计摘要：{proposal.audit_summary}</p>}
            {proposal.impact_scope !== undefined && <DetailBlock label="影响范围" content={detailContent(proposal.impact_scope)} />}
            {proposal.risk_notes !== undefined && <DetailBlock label="风险提示" content={detailContent(proposal.risk_notes)} />}
            {proposal.effect_validation && (
              <section className="detail-section" aria-label="规则效果验证">
                <div className="state-label">规则效果验证</div>
                <p>验证状态：{effectValidationStatusText[proposal.effect_validation.validation_status] ?? proposal.effect_validation.validation_status}</p>
                <p>过拟合风险：{overfitRiskText[proposal.effect_validation.overfit_risk ?? ''] ?? (proposal.effect_validation.overfit_risk || '未知')}</p>
                <p>历史回放：{replayResultText[proposal.effect_validation.replay_result ?? ''] ?? (proposal.effect_validation.replay_result || '未知')}</p>
                <p>门禁结论：{guardrailDecisionText[proposal.effect_validation.guardrail_decision ?? ''] ?? (proposal.effect_validation.guardrail_decision || '未知')}</p>
                <p>样本窗口：{proposal.effect_validation.sample_window || '暂无'}；样本数：{proposal.effect_validation.sample_count}</p>
                {proposal.effect_validation.source_explanation !== undefined && <DetailBlock label="验证来源" content={detailContent(proposal.effect_validation.source_explanation)} />}
                {proposal.effect_validation.metrics !== undefined && <DetailBlock label="验证指标" content={detailContent(proposal.effect_validation.metrics)} />}
                {proposal.effect_validation.risk_notes !== undefined && <DetailBlock label="验证风险提示" content={detailContent(proposal.effect_validation.risk_notes)} />}
                {proposal.effect_validation.related_error_cases !== undefined && <DetailBlock label="关联误判案例" content={detailContent(proposal.effect_validation.related_error_cases)} />}
                {proposal.effect_validation.related_decision_ids !== undefined && <DetailBlock label="关联决策记录" content={detailContent(proposal.effect_validation.related_decision_ids)} />}
                {proposal.effect_validation.related_risk_alert_ids !== undefined && <DetailBlock label="关联风险预警" content={detailContent(proposal.effect_validation.related_risk_alert_ids)} />}
                {proposal.effect_validation.related_audit_event_ids !== undefined && <DetailBlock label="关联审计事件" content={detailContent(proposal.effect_validation.related_audit_event_ids)} />}
                {proposal.effect_validation.safety_note && <p>{safeRuleSafetyNote(proposal.effect_validation.safety_note)}</p>}
              </section>
            )}
            {proposal.before_rule !== undefined && <DetailBlock label="变更前规则" content={ruleContent(proposal.before_rule)} />}
            {proposal.after_rule !== undefined && <DetailBlock label="变更后规则" content={ruleContent(proposal.after_rule)} />}
            {proposal.status === 'pending_user_confirm' && (
              <div className="action-row">
                <button type="button" onClick={() => onConfirm?.(proposal.proposal_id, true)}>确认送审</button>
                <button type="button" onClick={() => onConfirm?.(proposal.proposal_id, false)}>拒绝提案</button>
                <small>确认后进入守门人审计，正式规则仍需最终确认。</small>
              </div>
            )}
            {proposal.status === 'pending_final_confirm' && (
              <div className="action-row">
                {/* 守门人通过后仍需用户最终确认，正式规则不会自动生效。 */}
                <button type="button" onClick={() => onFinalConfirm?.(proposal.proposal_id, true)}>确认应用到正式规则</button>
                <button type="button" onClick={() => onFinalConfirm?.(proposal.proposal_id, false)}>拒绝应用</button>
                <small>守门人通过后仍需用户最终确认，正式规则不会自动生效。</small>
              </div>
            )}
          </section>
        ))
      )}
    </article>
  )
}

function DetailBlock({ label, content }: { label: string; content: string }) {
  return (
    <details className="raw-detail">
      <summary>查看{label}</summary>
      <pre aria-label={label}>{content}</pre>
    </details>
  )
}

function detailContent(value: unknown) {
  return JSON.stringify(redactRuleDetail(value), null, 2)
}

function ruleContent(value: unknown) {
  if (value && typeof value === 'object' && 'content' in value && typeof value.content === 'string') {
    return value.content
  }
  return JSON.stringify(value, null, 2)
}

function safeRuleSafetyNote(value: string) {
  if (/自动应用规则|自动规则应用/.test(value)) {
    return '规则效果验证只用于本地规则治理；规则生效仍需用户手动确认。'
  }
  return value
}

function redactRuleDetail(value: unknown): unknown {
  if (typeof value === 'string') {
    return value
      .replace(/自动应用规则|自动规则应用/g, '规则生效需人工确认')
      .replace(/外部推送|第三方推送/g, '站外通知')
      .replace(/自动确认/g, '人工确认')
      .replace(/自动修复/g, '人工复验')
      .replace(/收益承诺/g, '收益边界说明')
  }
  if (Array.isArray(value)) {
    return value.map((item) => redactRuleDetail(item))
  }
  if (value && typeof value === 'object') {
    return Object.fromEntries(Object.entries(value).map(([key, item]) => [key, redactRuleDetail(item)]))
  }
  return value
}
