import type { PortfolioSummary as PortfolioSummaryData } from '../../types/dashboard'
import { formatCurrency, formatPercent } from '../../shared/utils'

interface Props {
  summary: PortfolioSummaryData
}

export function PortfolioSummary({ summary }: Props) {
  return (
    <article className="cockpit-card">
      <div className="state-label">账户摘要</div>
      <div className="metric-grid">
        <div>
          <span>总资产</span>
          <strong>{formatCurrency(summary.total_assets)}</strong>
        </div>
        <div>
          <span>现金比例</span>
          <strong>{formatPercent(summary.cash_ratio)}</strong>
        </div>
        <div>
          <span>高风险比例</span>
          <strong>{formatPercent(summary.high_risk_ratio)}</strong>
        </div>
        <div>
          <span>持仓数量</span>
          <strong>{summary.position_count}</strong>
        </div>
      </div>
    </article>
  )
}
