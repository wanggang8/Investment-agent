import type { ReactNode } from 'react'
import { joinClassNames, type UITone } from './types'

type SummaryCardProps = {
  title: string
  value: ReactNode
  detail?: ReactNode
  tone?: UITone
  action?: {
    label: string
    href: string
  }
}

export function SummaryCard({ title, value, detail, tone = 'readonly', action }: SummaryCardProps) {
  return (
    <article aria-label={title} className={joinClassNames('ui-summary-card', `ui-summary-card-${tone}`)}>
      <div className="state-label">{title}</div>
      <strong>{value}</strong>
      {detail ? <p>{detail}</p> : null}
      {action ? <a href={action.href}>{action.label}</a> : null}
    </article>
  )
}
