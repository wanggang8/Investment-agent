import type { ReactNode } from 'react'
import { joinClassNames, type UITone } from './types'

type StatusBadgeProps = {
  tone: UITone
  children: ReactNode
  className?: string
}

export function StatusBadge({ tone, children, className }: StatusBadgeProps) {
  return (
    <span role="status" aria-label={textFromNode(children)} className={joinClassNames('ui-status-badge', `ui-status-badge-${tone}`, className)}>
      {children}
    </span>
  )
}

function textFromNode(node: ReactNode) {
  if (typeof node === 'string' || typeof node === 'number') return String(node)
  return undefined
}
