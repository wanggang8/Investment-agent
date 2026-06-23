import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { StatusBadge } from './StatusBadge'

describe('StatusBadge', () => {
  it('renders status text for degraded states instead of color-only state', () => {
    render(<StatusBadge tone="degraded">数据降级</StatusBadge>)

    const badge = screen.getByRole('status', { name: '数据降级' })
    expect(badge).toHaveTextContent('数据降级')
    expect(badge).toHaveClass('ui-status-badge-degraded')
  })

  it('supports readonly and blocked tones', () => {
    render(
      <div>
        <StatusBadge tone="readonly">只读追踪</StatusBadge>
        <StatusBadge tone="blocked">已阻断</StatusBadge>
      </div>,
    )

    expect(screen.getByRole('status', { name: '只读追踪' })).toHaveClass('ui-status-badge-readonly')
    expect(screen.getByRole('status', { name: '已阻断' })).toHaveClass('ui-status-badge-blocked')
  })
})
