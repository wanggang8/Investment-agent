import type { DailyWorkbenchModel } from '../../features/dashboard/dailyWorkbenchModel'
import { ReferenceHero } from '../reference'

interface Props {
  model: DailyWorkbenchModel
  eyebrow?: string
}

export function DailyDecisionHero({ model, eyebrow = '今日先看' }: Props) {
  return (
    <ReferenceHero
      iconLabel={eyebrow}
      title={model.verdictText}
      statusText={model.trustSummary}
      description={`${model.trustSummary}；${model.riskSummary}。`}
      stateTitle="当前纪律状态"
      stateValue={model.statusLabel}
      stateSummary={model.verdictText}
      stateDetail={model.updatedAtText}
      stateRegionLabel="今日纪律状态"
      prohibitedTitle="禁止动作"
      prohibitedActions={model.prohibitedActions}
      optionalActions={model.optionalActions}
    />
  )
}
