import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { Button } from './Button'

describe('Button', () => {
  it('renders a named primary button with a stable type', () => {
    render(<Button type="button">保存本地校准</Button>)

    const button = screen.getByRole('button', { name: '保存本地校准' })
    expect(button).toHaveAttribute('type', 'button')
    expect(button).toHaveClass('ui-button-primary')
  })

  it('exposes working and disabled state without relying on color', () => {
    render(<Button isWorking workingLabel="正在记录">记录本地动作</Button>)

    const button = screen.getByRole('button', { name: '正在记录' })
    expect(button).toBeDisabled()
    expect(button).toHaveAttribute('aria-busy', 'true')
  })

  it('supports danger tone for destructive local records', () => {
    render(<Button variant="danger">移除当前持仓</Button>)

    expect(screen.getByRole('button', { name: '移除当前持仓' })).toHaveClass('ui-button-danger')
  })
})
