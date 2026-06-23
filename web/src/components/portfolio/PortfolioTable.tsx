import type { Position } from '../../types/portfolio'
import { formatCurrency, formatPercent } from '../../shared/utils'
import { positionStateText, textOrRaw } from '../../shared/mappers'
import { ResponsiveTable } from '../ui'

interface Props {
  positions: Position[]
}

export function PortfolioTable({ positions }: Props) {
  return (
    <article className="cockpit-card">
      <div className="state-label">当前持仓</div>
      <ResponsiveTable
        caption="当前持仓明细"
        rows={positions}
        getRowKey={(position) => position.position_id}
        emptyText="暂无持仓记录。"
        columns={[
          { key: 'symbol', header: '标的', render: (position) => `${position.symbol} ${position.name}` },
          { key: 'quantity', header: '数量', render: (position) => position.quantity },
          { key: 'cost', header: '成本', render: (position) => formatCurrency(position.cost_price) },
          { key: 'current', header: '现价', render: (position) => formatCurrency(position.current_price) },
          { key: 'market', header: '市值', render: (position) => formatCurrency(position.market_value) },
          { key: 'profit', header: '浮盈', render: (position) => formatPercent(position.unrealized_profit_ratio) },
          { key: 'state', header: '状态', render: (position) => textOrRaw(positionStateText, position.position_state) },
          { key: 'buy_date', header: '买入日期', render: (position) => position.buy_date || '暂无买入日期' },
          { key: 'reason', header: '买入理由', render: (position) => position.buy_reason || '暂无买入理由' },
        ]}
      />
    </article>
  )
}
