import { useMemo, useState } from 'react'
import type { AuditEvent } from '../../types/audit'
import { auditActionText, auditActorText, auditStatusText, textOrRaw } from '../../shared/mappers'

interface Props {
  events: AuditEvent[]
}

export function AuditEventTimeline({ events }: Props) {
  const [statusFilter, setStatusFilter] = useState('all')
  const [expandedId, setExpandedId] = useState<string>()
  const filteredEvents = useMemo(() => events.filter((event) => statusFilter === 'all' || event.status === statusFilter), [events, statusFilter])

  return (
    <article className="cockpit-card">
      <div className="state-label">审计时间线</div>
      <label className="filter-label">
        筛选审计状态
        <select value={statusFilter} onChange={(event) => setStatusFilter(event.target.value)}>
          <option value="all">全部</option>
          <option value="success">成功</option>
          <option value="degraded">降级</option>
          <option value="failed">失败</option>
        </select>
      </label>
      {filteredEvents.length === 0 ? (
        <p className="muted-text">暂无匹配审计事件。</p>
      ) : (
        <ol className="timeline">
          {filteredEvents.map((event) => {
            const eventId = event.audit_event_id ?? event.event_id ?? event.request_id
            const expanded = expandedId === eventId
            const detailId = `audit-event-detail-${eventId}`
            return (
            <li key={eventId}>
              <strong>{textOrRaw(auditActionText, event.action)}</strong>
              <p>
                节点：{event.node_name || '无'} / 动作：{textOrRaw(auditActionText, event.node_action, '未知动作')} / 执行方：{textOrRaw(auditActorText, event.actor, '未知执行方')} / 状态：{textOrRaw(auditStatusText, event.status)}
              </p>
              <p>时间：{event.created_at} / 工作流：{event.workflow_type || '无'} / 规则版本：{event.rule_version || '无'} / 快照：{event.snapshot_id || '无'}</p>
              <p className="muted-text">关联：{entityRefs(event)}</p>
              <p className="muted-text">引用：{compactRef(event.input_ref_type, event.input_ref)} → {compactRef(event.output_ref_type, event.output_ref)}</p>
              <p>状态变化：{event.before_state || '-'} → {event.after_state || '-'}</p>
              <p>错误码：{event.error_code || '无'}</p>
              <button className="link-button" type="button" aria-expanded={expanded} aria-controls={detailId} onClick={() => setExpandedId(expanded ? undefined : eventId)}>
                {expanded ? '收起' : '展开'}引用
              </button>
              {expanded && (
                <div id={detailId} role="region" aria-label="审计引用详情">
                  <p>input：{event.input_ref_type || '-'} {event.input_ref || ''}</p>
                  <p>output：{event.output_ref_type || '-'} {event.output_ref || ''}</p>
                  <p>request_id：{event.request_id}</p>
                </div>
              )}
            </li>
            )
          })}
        </ol>
      )}
    </article>
  )
}

function compactRef(type?: string, value?: string) {
  const refType = type || '-'
  const refValue = value || ''
  if (refValue.length <= 96) {
    return `${refType} ${refValue}`.trim()
  }
  return `${refType} ${refValue.slice(0, 96)}...`.trim()
}

function entityRefs(event: AuditEvent) {
  const refs = [
    event.decision_id ? `决策 ${event.decision_id}` : '',
    event.proposal_id ? `提案 ${event.proposal_id}` : '',
    event.confirmation_id ? `确认 ${event.confirmation_id}` : '',
    event.error_case_id ? `错误案例 ${event.error_case_id}` : '',
  ].filter(Boolean)
  return refs.length ? refs.join(' / ') : '无'
}
