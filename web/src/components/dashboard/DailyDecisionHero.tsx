import type { DailyWorkbenchModel } from '../../features/dashboard/dailyWorkbenchModel'

interface Props {
  model: DailyWorkbenchModel
  eyebrow?: string
}

export function DailyDecisionHero({ model, eyebrow = '今日先看' }: Props) {
  return (
    <section className={`daily-hero daily-tone-${model.statusTone}`} aria-label="今日纪律状态">
      <div className="daily-hero-main">
        <div className="state-label">{eyebrow}</div>
        <h2>{model.verdictText}</h2>
        <p>{model.trustSummary}；{model.riskSummary}。</p>
        <dl className="daily-hero-meta">
          <div>
            <dt>状态</dt>
            <dd>{model.statusLabel}</dd>
          </div>
          <div>
            <dt>更新时间</dt>
            <dd>{model.updatedAtText}</dd>
          </div>
        </dl>
      </div>

      <div className="daily-hero-side" aria-label="纪律边界">
        <div>
          <strong>禁止动作</strong>
          {model.prohibitedActions.length > 0 ? (
            <ul>
              {model.prohibitedActions.map((action) => <li key={action}>{action}</li>)}
            </ul>
          ) : (
            <p>暂无新增禁止动作。</p>
          )}
        </div>
        <div>
          <strong>可选人工动作</strong>
          {model.optionalActions.length > 0 ? (
            <ul>
              {model.optionalActions.map((action) => <li key={action}>{action}</li>)}
            </ul>
          ) : (
            <p>暂无可选动作；继续观察现有本地事实。</p>
          )}
        </div>
      </div>
    </section>
  )
}
