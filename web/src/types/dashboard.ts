import type { EvidenceSummary } from './evidence'

export type DashboardState =
  | 'first_use'
  | 'normal'
  | 'insufficient_data'
  | 'frozen_watch'
  | 'high_risk'
  | 'data_source_unavailable'
  | 'generic_failure'

export interface DashboardTodayResponse {
  dashboard_state: DashboardState
  discipline_status: string
  data_updated_at: string
  portfolio_summary: PortfolioSummary
  market_summary: MarketSummary
  triggered_rules: TriggeredRule[]
  decision_summary: DecisionSummary
  evidence_summary?: EvidenceSummary
}

export interface PortfolioSummary {
  total_assets: number
  cash_ratio: number
  high_risk_ratio: number
  position_count: number
}

export interface MarketSummary {
  pe_percentile?: number
  pb_percentile?: number
  sentiment_state: string
  liquidity_state: string
}

export interface TriggeredRule {
  rule_id: string
  rule_name: string
  severity: 'normal' | 'warning' | 'danger' | 'frozen_watch' | 'insufficient' | string
  description: string
}

export interface DecisionSummary {
  decision_id?: string
  verdict: string
  final_verdict_status: string
  prohibited_actions: string[] | null
  optional_actions: string[] | null
  action_required: boolean
  confirmation_status: string
}
