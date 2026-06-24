import type { DailyAction } from '../../features/dashboard/dailyWorkbenchModel'
import { PriorityActionQueue } from '../reference'

interface Props {
  actions: DailyAction[]
  title?: string
}

export function ManualActionQueue({ actions, title = '下一步人工动作' }: Props) {
  return <PriorityActionQueue title={title} actions={actions} />
}
