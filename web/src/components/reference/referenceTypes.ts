import type { ReactNode } from 'react'

export type ReferenceTone = 'success' | 'warning' | 'danger' | 'degraded' | 'unknown' | 'readonly'

export type ReferenceActionPriority = 'high' | 'medium' | 'low' | 'blocking' | 'review' | 'follow_up'

export type ReferenceAction = {
  label: string
  detail: string
  href: string
  priority: ReferenceActionPriority
  meta?: string
}

export type ReferenceMetric = {
  label: string
  value: ReactNode
  status?: string
  tone?: ReferenceTone
  icon?: ReactNode
  details?: ReactNode[]
}

export type ReferenceSnapshotItem = {
  label: string
  value: ReactNode
  status?: string
}

export type ReferenceProgressStep = {
  label: string
  status: 'done' | 'active' | 'pending' | 'blocked'
  detail?: string
}

export type ReferenceChecklistItem = {
  label: string
  value: ReactNode
  status: 'done' | 'active' | 'pending' | 'blocked'
}
