import type { RiskAlert, RiskSeverity, RiskSOPStatus } from '../../types/riskAlert'

export type RiskQueueID = 'pending_review' | 'in_progress' | 'needs_review' | 'recorded'
export type RiskActionPriority = 'blocking' | 'review' | 'follow_up'

export interface RiskDispositionAction {
  label: string
  detail: string
  priority: RiskActionPriority
}

export interface RiskDispositionQueue {
  id: RiskQueueID
  label: string
  description: string
  items: RiskAlert[]
}

export interface RiskDispositionModel {
  summaryLabel: string
  highestSeverity: RiskSeverity
  affectedSymbols: string[]
  queues: RiskDispositionQueue[]
  nextActions: RiskDispositionAction[]
  safetyNotes: string[]
}

const queueMeta: Record<RiskQueueID, Omit<RiskDispositionQueue, 'items'>> = {
  pending_review: { id: 'pending_review', label: '待看', description: '新触发风险，先确认禁止动作和触发依据。' },
  in_progress: { id: 'in_progress', label: '处理中', description: '仍在观察或处理中，需要继续人工跟踪。' },
  needs_review: { id: 'needs_review', label: '需复盘', description: '已升级风险，需要优先复核关联材料。' },
  recorded: { id: 'recorded', label: '已记录', description: '已解除或归档，仅作为本地审计线索。' },
}

const severityRank: Record<RiskSeverity, number> = { info: 0, warning: 1, critical: 2 }

export function buildRiskDispositionModel(alerts: RiskAlert[]): RiskDispositionModel {
  const queues = (Object.keys(queueMeta) as RiskQueueID[]).map((id) => ({ ...queueMeta[id], items: [] as RiskAlert[] }))

  alerts.forEach((alert) => {
    queues.find((queue) => queue.id === queueForStatus(alert.sop_status))?.items.push(alert)
  })

  const needsReviewCount = queues.find((queue) => queue.id === 'needs_review')?.items.length ?? 0
  const highestSeverity = alerts.reduce<RiskSeverity>((current, alert) => severityRank[alert.severity] > severityRank[current] ? alert.severity : current, 'info')
  const affectedSymbols = Array.from(new Set(alerts.map((alert) => alert.symbol).filter(Boolean))).sort()

  return {
    summaryLabel: alerts.length ? `${alerts.length} 条风险事实，${needsReviewCount} 条需复盘` : '暂无风险预警',
    highestSeverity,
    affectedSymbols,
    queues,
    nextActions: buildRiskActions(queues),
    safetyNotes: ['风险预警只记录本地 SOP 状态，不会自动交易、外部推送、自动确认或改变持仓。'],
  }
}

function queueForStatus(status: RiskSOPStatus): RiskQueueID {
  if (status === 'triggered') return 'pending_review'
  if (status === 'active' || status === 'observing') return 'in_progress'
  if (status === 'escalated') return 'needs_review'
  return 'recorded'
}

function buildRiskActions(queues: RiskDispositionQueue[]): RiskDispositionAction[] {
  if ((queues.find((queue) => queue.id === 'needs_review')?.items.length ?? 0) > 0) {
    return [{ label: '复核升级风险', detail: '先查看升级风险的触发依据、禁止动作和关联决策。', priority: 'blocking' }]
  }
  if ((queues.find((queue) => queue.id === 'pending_review')?.items.length ?? 0) > 0) {
    return [{ label: '查看新触发风险', detail: '确认风险是否需要继续观察、升级复核或解除。', priority: 'blocking' }]
  }
  if ((queues.find((queue) => queue.id === 'in_progress')?.items.length ?? 0) > 0) {
    return [{ label: '跟踪处理中风险', detail: '复核观察中的风险是否仍然有效。', priority: 'review' }]
  }
  return [{ label: '继续观察今日纪律', detail: '暂无风险预警时仍需保持每日纪律和数据质量检查。', priority: 'follow_up' }]
}
