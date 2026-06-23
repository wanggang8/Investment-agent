import type { DecisionSummary } from '../../types/dashboard'

interface Props {
  summary: DecisionSummary
}

export function FinalVerdictCard({ summary }: Props) {
  const prohibitedActions = Array.isArray(summary.prohibited_actions) ? summary.prohibited_actions : []
  const optionalActions = Array.isArray(summary.optional_actions) ? summary.optional_actions : []

  return (
    <article className="cockpit-card verdict-card">
      <div className="state-label">今日建议</div>
      <h2>{summary.verdict}</h2>
      <div className="verdict-section">
        <strong>禁止事项</strong>
        {prohibitedActions.length === 0 ? (
          <p>暂无新增禁止事项。</p>
        ) : (
          <ul>
            {prohibitedActions.map((action) => (
              <li key={action}>{action}</li>
            ))}
          </ul>
        )}
      </div>
      <div className="verdict-section">
        <strong>可选记录</strong>
        {optionalActions.length === 0 ? (
          <p>暂无可选记录。</p>
        ) : (
          <ul>
            {optionalActions.map((action) => (
              <li key={action}>{action}</li>
            ))}
          </ul>
        )}
      </div>
    </article>
  )
}
