import { Button } from './Button'
import { redactSensitiveText } from '../../shared/utils'

type ErrorStateProps = {
  title: string
  message: string
  retryLabel?: string
  onRetry?: () => void
}

export function ErrorState({ title, message, retryLabel, onRetry }: ErrorStateProps) {
  return (
    <section className="ui-error-state" role="alert" aria-label={title}>
      <h2>{title}</h2>
      <p>{sanitizeErrorMessage(message)}</p>
      {onRetry && retryLabel ? <Button variant="secondary" onClick={onRetry}>{retryLabel}</Button> : null}
    </section>
  )
}

function sanitizeErrorMessage(message: string) {
  return redactSensitiveText(message)
}
