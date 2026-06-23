import type { DashboardTodayResponse } from '../../types/dashboard'
import { dashboardStateText, textOrRaw } from '../../shared/mappers'

interface Props {
  dashboard: DashboardTodayResponse
  errorMessage?: string
}

const stateCopy: Record<string, string> = {
  first_use: '需要先录入账户和持仓，系统才能生成纪律报告。',
  normal: '今日未触发纪律红线。',
  insufficient_data: '当前信息不足，系统暂停生成交易类建议。',
  frozen_watch: '信息仍在核验中，进入冻结观察。',
  high_risk: '已触发高危纪律状态，禁止新增买入。',
}

export function DisciplineStatus({ dashboard, errorMessage }: Props) {
  return (
    <article className={`cockpit-card state-${dashboard.dashboard_state.replace('_', '-')}`}>
      <div className="state-label">纪律状态</div>
      <h2>{dashboard.discipline_status}</h2>
      <p>{errorMessage ?? stateCopy[dashboard.dashboard_state] ?? '等待系统更新状态。'}</p>
      <dl className="compact-list">
        <div>
          <dt>数据更新时间</dt>
          <dd>{dashboard.data_updated_at || '暂无'}</dd>
        </div>
        <div>
          <dt>当前状态</dt>
          <dd>{textOrRaw(dashboardStateText, dashboard.dashboard_state)}</dd>
        </div>
      </dl>
    </article>
  )
}
