import type { ReactNode } from 'react'
import { StatusBadge } from './StatusBadge'
import type { UITone } from './types'

type HeaderMetric = {
  label: string
  value: ReactNode
}

type HeaderAction = {
  label: string
  href: string
}

type PageHeaderProps = {
  eyebrow?: string
  title: string
  description?: ReactNode
  status?: {
    tone: UITone
    label: string
  }
  metrics?: HeaderMetric[]
  actions?: HeaderAction[]
}

export function PageHeader({ eyebrow, title, description, status, metrics = [], actions = [] }: PageHeaderProps) {
  return (
    <header className="ui-page-header">
      <div className="ui-page-header-main">
        {eyebrow ? <div className="state-label">{eyebrow}</div> : null}
        <h1>{title}</h1>
        {description ? <p>{description}</p> : null}
      </div>
      <div className="ui-page-header-side">
        {status ? <StatusBadge tone={status.tone}>{status.label}</StatusBadge> : null}
        {metrics.length ? (
          <dl className="ui-page-header-metrics">
            {metrics.map((metric) => (
              <div key={metric.label}>
                <dt>{metric.label}</dt>
                <dd>{metric.value}</dd>
              </div>
            ))}
          </dl>
        ) : null}
        {actions.length ? (
          <div className="ui-page-header-actions">
            {actions.map((action) => (
              <a key={action.href} href={action.href}>{action.label}</a>
            ))}
          </div>
        ) : null}
      </div>
    </header>
  )
}
