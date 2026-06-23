import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { DetailSection } from './DetailSection'
import { PageHeader } from './PageHeader'
import { SummaryCard } from './SummaryCard'

describe('structural UI primitives', () => {
  it('renders a page header with status metrics and next manual actions', () => {
    render(
      <PageHeader
        eyebrow="本地工作台"
        title="数据质量可观测"
        description="只读检查本地质量事实。"
        status={{ tone: 'warning', label: '需要人工检查' }}
        metrics={[{ label: '证据', value: '3 条' }]}
        actions={[{ label: '查看证据', href: '/evidence' }]}
      />,
    )

    expect(screen.getByRole('heading', { name: '数据质量可观测' })).toBeVisible()
    expect(screen.getByRole('status', { name: '需要人工检查' })).toBeVisible()
    expect(screen.getByText('证据')).toBeVisible()
    expect(screen.getByRole('link', { name: '查看证据' })).toHaveAttribute('href', '/evidence')
  })

  it('keeps summary card tone readable', () => {
    render(<SummaryCard title="数据源健康" value="降级" detail="需要人工复核" tone="degraded" />)

    expect(screen.getByRole('article', { name: '数据源健康' })).toHaveClass('ui-summary-card-degraded')
    expect(screen.getByText('需要人工复核')).toBeVisible()
  })

  it('toggles detail section through a button with aria-expanded', () => {
    render(
      <DetailSection title="展开引用" summary="查看本地审计引用">
        <p>audit_event_id: event_1</p>
      </DetailSection>,
    )

    const toggle = screen.getByRole('button', { name: /展开引用/ })
    expect(toggle).toHaveAttribute('aria-expanded', 'false')
    expect(screen.queryByText('audit_event_id: event_1')).not.toBeInTheDocument()

    fireEvent.click(toggle)

    expect(toggle).toHaveAttribute('aria-expanded', 'true')
    expect(screen.getByText('audit_event_id: event_1')).toBeVisible()
  })
})
