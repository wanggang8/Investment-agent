import type { TriggeredRule } from '../../types/dashboard'

interface Props {
  rules: TriggeredRule[]
}

export function TriggeredRules({ rules }: Props) {
  return (
    <article className="cockpit-card">
      <div className="state-label">风险红线</div>
      {rules.length === 0 ? (
        <p>暂未触发纪律红线。</p>
      ) : (
        <ul className="rule-list">
          {rules.map((rule) => (
            <li key={rule.rule_id} className={`rule-item severity-${rule.severity}`}>
              <strong>{rule.rule_name}</strong>
              <span>{rule.description}</span>
            </li>
          ))}
        </ul>
      )}
    </article>
  )
}
