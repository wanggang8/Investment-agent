import { useId, useState, type ReactNode } from 'react'

type DetailSectionProps = {
  title: string
  summary?: ReactNode
  defaultOpen?: boolean
  children: ReactNode
}

export function DetailSection({ title, summary, defaultOpen = false, children }: DetailSectionProps) {
  const [open, setOpen] = useState(defaultOpen)
  const regionID = useId()
  return (
    <section className="ui-detail-section">
      <button type="button" className="ui-detail-toggle" aria-expanded={open} aria-controls={regionID} onClick={() => setOpen((current) => !current)}>
        {title}
      </button>
      {summary ? <p className="ui-detail-summary">{summary}</p> : null}
      {open ? (
        <div id={regionID} className="ui-detail-body">
          {children}
        </div>
      ) : null}
    </section>
  )
}
