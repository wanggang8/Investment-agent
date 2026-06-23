import type { DisplayState } from '../../types/api'
import { getStatusNoticeCopy } from './statusNoticeCopy'

interface Props {
  state: DisplayState
  safeMessage?: string
  code?: string
}

export function StatusNotice({ state, safeMessage, code }: Props) {
  const copy = getStatusNoticeCopy(state, code)
  return (
    <article className={`cockpit-card state-notice state-notice-${state}`}>
      <div className="state-label">状态提示</div>
      <h2>{copy.title}</h2>
      <p>{safeMessage || copy.message}</p>
    </article>
  )
}
