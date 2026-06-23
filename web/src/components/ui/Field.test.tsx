import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { Field } from './Field'

describe('Field', () => {
  it('associates label hint and error with the control', () => {
    render(
      <Field id="symbol" label="标的代码" hint="只记录本地事实" error="请填写标的代码" required>
        <input />
      </Field>,
    )

    const input = screen.getByLabelText('标的代码')
    expect(input).toHaveAttribute('id', 'symbol')
    expect(input).toHaveAttribute('required')
    expect(input).toHaveAttribute('aria-invalid', 'true')
    expect(input.getAttribute('aria-describedby')).toContain('symbol-hint')
    expect(input.getAttribute('aria-describedby')).toContain('symbol-error')
    expect(screen.getByText('只记录本地事实')).toHaveAttribute('id', 'symbol-hint')
    expect(screen.getByRole('alert')).toHaveTextContent('请填写标的代码')
  })

  it('keeps textarea controls accessible through the same label pattern', () => {
    render(
      <Field id="reason" label="确认理由">
        <textarea />
      </Field>,
    )

    expect(screen.getByLabelText('确认理由').tagName).toBe('TEXTAREA')
  })
})
