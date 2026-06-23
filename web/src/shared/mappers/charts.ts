import type { DashboardTodayResponse } from '../../types/dashboard'
import type { Position } from '../../types/portfolio'
import type { ReviewSummary } from '../../types/review'

export type ChartTone = 'safe' | 'warn' | 'danger' | 'info' | 'muted'

export interface ChartDatum {
  label: string
  value: number
  tone: ChartTone
}

export interface DashboardChartData {
  allocation: ChartDatum[]
  risk: ChartDatum[]
  evidence: ChartDatum[]
}

function percent(value?: number) {
  return Math.round((value ?? 0) * 100)
}

export function buildDashboardChartData(dashboard: Pick<DashboardTodayResponse, 'portfolio_summary' | 'market_summary' | 'evidence_summary'>): DashboardChartData {
  const cash = percent(dashboard.portfolio_summary.cash_ratio)
  const highRisk = percent(dashboard.portfolio_summary.high_risk_ratio)
  // 这里把 PE/PB 分位保留为 0-100 展示值，避免图表层再次理解后端市场含义。
  const pe = Math.round(dashboard.market_summary.pe_percentile ?? 0)
  const pb = Math.round(dashboard.market_summary.pb_percentile ?? 0)

  return {
    allocation: [
      { label: '现金', value: cash, tone: cash < 5 ? 'warn' : 'safe' },
      { label: '持仓', value: Math.max(0, 100 - cash), tone: 'info' },
    ],
    risk: [
      { label: '高风险仓位', value: highRisk, tone: highRisk > 30 ? 'danger' : 'safe' },
      { label: 'PE 分位', value: pe, tone: pe > 80 ? 'danger' : pe > 50 ? 'warn' : 'safe' },
      { label: 'PB 分位', value: pb, tone: pb > 80 ? 'danger' : pb > 50 ? 'warn' : 'safe' },
    ],
    evidence: [
      { label: '证据来源', value: dashboard.evidence_summary?.source_count ?? 0, tone: 'info' },
      { label: '持仓数量', value: dashboard.portfolio_summary.position_count, tone: 'muted' },
    ],
  }
}

export function buildPortfolioAllocationData(positions: Position[]): ChartDatum[] {
  const total = positions.reduce((sum, item) => sum + item.market_value, 0)
  if (total <= 0) return []
  const groups = new Map<string, { value: number; tone: ChartTone }>()
  for (const item of positions) {
    const label = item.name || item.symbol
    const current = groups.get(label) ?? { value: 0, tone: 'safe' as ChartTone }
    groups.set(label, {
      value: current.value + item.market_value,
      tone: riskierTone(current.tone, item.position_state === 'sell_only' ? 'danger' : item.position_state === 'frozen_watch' ? 'warn' : 'safe'),
    })
  }
  return Array.from(groups.entries()).map(([label, item]) => ({
    label,
    value: Math.round((item.value / total) * 100),
    // 资产状态来自 API DTO；frozen_watch/sell_only 用于提示风险，不在前端推导交易动作。
    tone: item.tone,
  }))
}

function riskierTone(current: ChartTone, next: ChartTone): ChartTone {
  const rank: Record<ChartTone, number> = { muted: 0, info: 1, safe: 2, warn: 3, danger: 4 }
  return rank[next] > rank[current] ? next : current
}

export function buildReviewActivityData(summary?: ReviewSummary): ChartDatum[] {
  return [
    { label: '确认动作', value: summary?.confirmation_count ?? 0, tone: 'safe' },
    { label: '错误案例', value: summary?.error_case_count ?? 0, tone: 'danger' },
    { label: '规则提案', value: summary?.rule_proposal_count ?? 0, tone: 'warn' },
    { label: '审计事件', value: summary?.audit_event_count ?? 0, tone: 'info' },
  ]
}
