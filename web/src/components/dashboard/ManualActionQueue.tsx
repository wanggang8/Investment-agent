import { Link } from 'react-router-dom'
import type { DailyAction } from '../../features/dashboard/dailyWorkbenchModel'

interface Props {
  actions: DailyAction[]
  title?: string
}

export function ManualActionQueue({ actions, title = '下一步人工动作' }: Props) {
  return (
    <section className="daily-action-queue" aria-label={title}>
      <div className="row-between">
        <div>
          <div className="state-label">只读导航</div>
          <h2>{title}</h2>
        </div>
        <small>系统只提示和记录，不会替你执行。</small>
      </div>
      <ol className="daily-action-list">
        {actions.map((action) => (
          <li key={`${action.label}:${action.href}`} className={`daily-action daily-action-${action.priority}`}>
            <div>
              <strong>{action.label}</strong>
              <p>{action.detail}</p>
            </div>
            <Link to={action.href}>进入</Link>
          </li>
        ))}
      </ol>
    </section>
  )
}
