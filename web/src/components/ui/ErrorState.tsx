import { Button } from './Button'

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
  return message
    .replace(/sk-[A-Za-z0-9_-]+/g, '已脱敏密钥')
    .replace(/\b(select|insert|update|delete|drop|alter|create)\b[\s\S]*/gi, 'SQL 已脱敏')
    .replace(/\bprompt\s*[:=][^\n\r]*/gi, 'prompt 已脱敏')
    .replace(/raw\s+(vendor|http|stack|payload)[^\n\r]*/gi, 'raw 诊断已脱敏')
    .replace(/stack trace/gi, '已脱敏错误摘要')
    .replace(/(?:[A-Za-z]:\\Users\\|[A-Za-z]:\\|\/Users\/|\/tmp\/|\/opt\/private\/|\/home\/)[^\s"']+/g, '已脱敏路径')
}
