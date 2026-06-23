import type { AuditEvent } from '../../types/audit'
import { auditActionText, textOrRaw } from '../../shared/mappers'
import { countBy, type OpsAction, type OpsMetric, type OpsTone } from './modelTypes'

export interface AuditEventGroup {
  label: string
  count: number
  latestAt: string
}

export interface AuditOpsModel {
  overallLabel: string
  overallTone: OpsTone
  metrics: OpsMetric[]
  nextActions: OpsAction[]
  eventGroups: AuditEventGroup[]
  safetyNotes: string[]
}

export function buildAuditOpsModel({ events }: { events?: AuditEvent[] }): AuditOpsModel {
  const items = events ?? []
  const failedOrDegraded = items.filter((item) => item.status === 'failed' || item.status === 'degraded')
  const latest = [...items].sort((a, b) => b.created_at.localeCompare(a.created_at))[0]
  const byAction = countBy(items.map((item) => item.action))

  return {
    overallLabel: items.length === 0 ? '暂无审计记录' : failedOrDegraded.length ? '审计记录需要检查' : '审计记录正常',
    overallTone: items.length === 0 ? 'unknown' : failedOrDegraded.length ? 'warning' : 'success',
    metrics: [
      { label: '审计事件', value: String(items.length) },
      { label: '失败/降级', value: String(failedOrDegraded.length), tone: failedOrDegraded.length ? 'warning' : 'success' },
      { label: '最近事件', value: latest ? textOrRaw(auditActionText, latest.action) : '暂无', detail: latest?.created_at },
    ],
    nextActions: buildAuditActions(failedOrDegraded.length, items),
    eventGroups: Object.entries(byAction).map(([action, count]) => ({
      label: textOrRaw(auditActionText, action),
      count,
      latestAt: latestForAction(items, action),
    })),
    safetyNotes: ['审计页只展示本地审计事实和关联入口，不读取 raw log、私有路径或本地数据库文件。'],
  }
}

function buildAuditActions(failedCount: number, events: AuditEvent[]): OpsAction[] {
  const actions: OpsAction[] = []
  if (failedCount) actions.push({ label: '查看失败事件', detail: '优先定位 failed/degraded 审计事件。', href: '/audit' })
  if (events.some((item) => item.proposal_id || item.action === 'audit_rule_change')) actions.push({ label: '检查规则治理', detail: '回到规则页复核提案和守门人结果。', href: '/rules' })
  actions.push({ label: '查看数据质量', detail: '确认数据源、证据和索引是否影响审计结果。', href: '/data-quality' })
  return actions
}

function latestForAction(events: AuditEvent[], action: string) {
  return events.filter((item) => item.action === action).sort((a, b) => b.created_at.localeCompare(a.created_at))[0]?.created_at ?? '暂无'
}

