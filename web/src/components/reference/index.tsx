import type { ReactNode } from 'react'
import { FileText } from 'lucide-react'
import { Link } from 'react-router-dom'
import type { ReferenceAction, ReferenceChecklistItem, ReferenceMetric, ReferenceProgressStep, ReferenceSnapshotItem } from './referenceTypes'

export type { ReferenceAction, ReferenceChecklistItem, ReferenceMetric, ReferenceProgressStep, ReferenceSnapshotItem, ReferenceTone } from './referenceTypes'

type ReferenceHeroProps = {
  iconLabel: string
  title: string
  statusText: string
  description: string
  stateTitle: string
  stateValue: ReactNode
  stateDetail?: ReactNode
  stateSummary?: ReactNode
  stateRegionLabel?: string
  prohibitedTitle: string
  prohibitedActions: string[]
  optionalTitle?: string
  optionalActions?: string[]
}

export function ReferenceHero({
  iconLabel,
  title,
  statusText,
  description,
  stateTitle,
  stateValue,
  stateDetail,
  stateSummary,
  stateRegionLabel,
  prohibitedTitle,
  prohibitedActions,
  optionalTitle = '可选人工动作',
  optionalActions = [],
}: ReferenceHeroProps) {
  return (
    <section className="reference-hero" aria-label="纪律报告概览">
      <div className="reference-hero-icon" aria-label={iconLabel}>
        <FileText size={28} strokeWidth={2.1} aria-hidden="true" />
        <span className="reference-sr-only">{iconLabel}</span>
      </div>
      <div className="reference-hero-main">
        <h2>{title}</h2>
        <p className="reference-hero-status">{statusText}</p>
        <p>{description}</p>
      </div>
      <div className="reference-hero-state" role="region" aria-label={stateRegionLabel ?? stateTitle}>
        <span>{stateTitle}</span>
        <strong>{stateValue}</strong>
        {stateSummary ? <p>{stateSummary}</p> : null}
        {stateDetail ? <p>{stateDetail}</p> : null}
      </div>
      <div className="reference-hero-prohibited">
        <span>{prohibitedTitle}</span>
        {prohibitedActions.length ? (
          <ul>
            {prohibitedActions.map((action) => <li key={action}>{action}</li>)}
          </ul>
        ) : (
          <p>暂无新增禁止动作。</p>
        )}
        {optionalActions.length ? (
          <>
            <span>{optionalTitle}</span>
            <ul>
              {optionalActions.map((action) => <li key={action}>{action}</li>)}
            </ul>
          </>
        ) : null}
      </div>
    </section>
  )
}

type PriorityActionQueueProps = {
  title: string
  actions: ReferenceAction[]
}

export function PriorityActionQueue({ title, actions }: PriorityActionQueueProps) {
  return (
    <section className="reference-action-queue" aria-label={title}>
      <div className="reference-panel-heading">
        <div>
          <h2>{title}</h2>
          <small>按优先级 · 系统只提示和记录，不会替你执行。</small>
        </div>
        <span>{actions.length} 项待处理</span>
      </div>
      <ol>
        {actions.map((action, index) => (
          <li key={`${action.label}:${action.href}`} className={`reference-action-row reference-action-${normalizePriority(action.priority)}`}>
            <span className="reference-action-index">{index + 1}</span>
            <div className="reference-action-body">
              <div>
                <strong>{action.label}</strong>
                <span className="reference-priority-chip">{priorityText(action.priority)}</span>
              </div>
              <p>{action.detail}</p>
            </div>
            {action.meta ? <small>{action.meta}</small> : null}
            <Link to={action.href}>{actionVerb(action.label)}</Link>
          </li>
        ))}
      </ol>
    </section>
  )
}

type StatusMetricGridProps = {
  title: string
  updatedAt?: string
  metrics: ReferenceMetric[]
}

export function StatusMetricGrid({ title, updatedAt, metrics }: StatusMetricGridProps) {
  return (
    <section className="reference-metric-grid" aria-label={title}>
      <div className="reference-panel-heading">
        <h2>{title}</h2>
        {updatedAt ? <small>更新于 {updatedAt}</small> : null}
      </div>
      <div className="reference-metric-cards">
        {metrics.map((metric) => (
          <article key={metric.label} className={`reference-metric-card reference-tone-${metric.tone ?? 'readonly'}`}>
            <div className="reference-metric-icon" aria-hidden="true">{metric.icon}</div>
            <span>{metric.label}</span>
            <strong>{metric.value}</strong>
            {metric.status ? <em>{metric.status}</em> : null}
            {metric.details?.length ? (
              <ul>
                {metric.details.map((detail, index) => <li key={index}>{detail}</li>)}
              </ul>
            ) : null}
          </article>
        ))}
      </div>
    </section>
  )
}

type SnapshotStripProps = {
  title: string
  updatedAt?: string
  items: ReferenceSnapshotItem[]
}

export function SnapshotStrip({ title, updatedAt, items }: SnapshotStripProps) {
  return (
    <section className="reference-snapshot-strip" aria-label={title}>
      <div className="reference-panel-heading">
        <h2>{title}</h2>
        {updatedAt ? <small>更新于 {updatedAt}</small> : null}
      </div>
      <dl>
        {items.map((item) => (
          <div key={item.label}>
            <dt>{item.label}</dt>
            <dd>{item.value}</dd>
            {item.status ? <span>{item.status}</span> : null}
          </div>
        ))}
      </dl>
    </section>
  )
}

type ProgressTrackerProps = {
  title: string
  actions?: ReactNode
  steps: ReferenceProgressStep[]
  children?: ReactNode
}

export function ProgressTracker({ title, actions, steps, children }: ProgressTrackerProps) {
  return (
    <section className="reference-progress-tracker" aria-label={title}>
      <div className="reference-panel-heading">
        <h2>{title}</h2>
        {actions}
      </div>
      <ol>
        {steps.map((step) => (
          <li key={step.label} className={`reference-progress-step reference-progress-${step.status}`}>
            <span aria-hidden="true" />
            <strong>{step.label}</strong>
            {step.detail ? <small>{step.detail}</small> : null}
          </li>
        ))}
      </ol>
      {children ? <div className="reference-progress-detail">{children}</div> : null}
    </section>
  )
}

type EvidenceChecklistProps = {
  title: string
  items: ReferenceChecklistItem[]
  action?: {
    label: string
    href: string
  }
}

export function EvidenceChecklist({ title, items, action }: EvidenceChecklistProps) {
  return (
    <section className="reference-checklist" aria-label={title}>
      <div className="reference-panel-heading">
        <h2>{title}</h2>
      </div>
      <ul>
        {items.map((item) => (
          <li key={item.label} className={`reference-checklist-${item.status}`}>
            <span aria-hidden="true" />
            <strong>{item.label}</strong>
            <em>{item.value}</em>
          </li>
        ))}
      </ul>
      {action ? <Link to={action.href}>{action.label}</Link> : null}
    </section>
  )
}

type LedgerSurfaceProps = {
  title: string
  children: ReactNode
  action?: ReactNode
}

export function LedgerSurface({ title, children, action }: LedgerSurfaceProps) {
  return (
    <section className="reference-ledger-surface" aria-label={title}>
      <div className="reference-panel-heading">
        <h2>{title}</h2>
        {action}
      </div>
      {children}
    </section>
  )
}

function normalizePriority(priority: ReferenceAction['priority']) {
  if (priority === 'blocking') return 'high'
  if (priority === 'review') return 'medium'
  if (priority === 'follow_up') return 'low'
  return priority
}

function priorityText(priority: ReferenceAction['priority']) {
  if (priority === 'high' || priority === 'blocking') return '高优先级'
  if (priority === 'medium' || priority === 'review') return '中优先级'
  return '低优先级'
}

function actionVerb(label: string) {
  if (label.includes('维护')) return '去维护'
  if (label.includes('处理')) return '去处理'
  if (label.includes('咨询') || label.includes('发起')) return '去咨询'
  if (label.includes('记录')) return '去记录'
  return '去查看'
}
