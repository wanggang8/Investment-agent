type EmptyStateProps = {
  title: string
  description: string
  action?: {
    label: string
    href: string
  }
}

export function EmptyState({ title, description, action }: EmptyStateProps) {
  return (
    <section className="ui-empty-state" aria-label={title}>
      <h2>{title}</h2>
      <p>{description}</p>
      {action ? <a href={action.href}>{action.label}</a> : null}
    </section>
  )
}
