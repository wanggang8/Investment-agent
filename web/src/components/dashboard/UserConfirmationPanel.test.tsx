import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import type { ConfirmationRequest } from '../../types/decision'
import { UserConfirmationPanel } from './UserConfirmationPanel'

describe('UserConfirmationPanel', () => {
  afterEach(() => cleanup())

  it('only renders offline action records and never auto trading entries', () => {
    render(<UserConfirmationPanel confirmationStatus="pending" availableActions={['planned', 'executed_manually', 'auto_trade', 'one_click_trade']} />)

    expect(screen.getByRole('button', { name: '记录计划' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '已手动执行' })).toBeInTheDocument()
    expect(screen.queryByText(/自动交易|自动执行交易|一键交易|代下单/)).not.toBeInTheDocument()
    expect(screen.getByText(/只记录你的线下动作/)).toBeInTheDocument()
    expect(screen.getByText(/不会替你买入或卖出/)).toBeInTheDocument()
  })

  it('shows unknown confirmation status safely', () => {
    render(<UserConfirmationPanel confirmationStatus="backend_new_state" availableActions={[]} />)

    expect(screen.getByText(/当前确认状态：未知状态/)).toBeInTheDocument()
    expect(screen.queryByText(/backend_new_state/)).not.toBeInTheDocument()
  })

  it('validates manual execution fields before submission', () => {
    const onSubmit = vi.fn()
    render(<UserConfirmationPanel decisionId="decision_1" confirmationStatus="pending" availableActions={['executed_manually']} onSubmit={onSubmit} />)

    fireEvent.click(screen.getByRole('button', { name: '已手动执行' }))
    fireEvent.click(screen.getByRole('button', { name: '提交确认' }))

    expect(screen.getByText(/请完整填写标的、动作、数量、价格和执行时间/)).toBeInTheDocument()
    expect(onSubmit).not.toHaveBeenCalled()
  })

  it('rejects future manual execution time before submission', () => {
    const onSubmit = vi.fn()
    render(<UserConfirmationPanel decisionId="decision_1" confirmationStatus="pending" availableActions={['executed_manually']} onSubmit={onSubmit} />)

    fireEvent.click(screen.getByRole('button', { name: '已手动执行' }))
    fireEvent.change(screen.getByLabelText('标的代码'), { target: { value: '510300' } })
    fireEvent.change(screen.getByLabelText('线下动作'), { target: { value: 'sell' } })
    fireEvent.change(screen.getByLabelText('数量'), { target: { value: '4' } })
    fireEvent.change(screen.getByLabelText('价格'), { target: { value: '3.5' } })
    fireEvent.change(screen.getByLabelText('执行时间'), { target: { value: '2999-05-29T03:00' } })
    fireEvent.click(screen.getByRole('button', { name: '提交确认' }))

    expect(screen.getByText(/执行时间不能晚于当前时间/)).toBeInTheDocument()
    expect(onSubmit).not.toHaveBeenCalled()
  })

  it('submits manual execution payload without auto trading copy', () => {
    const onSubmit = vi.fn()
    render(<UserConfirmationPanel decisionId="decision_1" confirmationStatus="pending" availableActions={['executed_manually']} onSubmit={onSubmit} />)

    fireEvent.click(screen.getByRole('button', { name: '已手动执行' }))
    fireEvent.change(screen.getByLabelText('标的代码'), { target: { value: '510300' } })
    fireEvent.change(screen.getByLabelText('线下动作'), { target: { value: 'sell' } })
    fireEvent.change(screen.getByLabelText('数量'), { target: { value: '4' } })
    fireEvent.change(screen.getByLabelText('价格'), { target: { value: '3.5' } })
    fireEvent.change(screen.getByLabelText('费用'), { target: { value: '1.2' } })
    fireEvent.change(screen.getByLabelText('执行时间'), { target: { value: '2026-05-29T03:00' } })
    fireEvent.click(screen.getByRole('button', { name: '提交确认' }))

    expect(onSubmit).toHaveBeenCalledWith('decision_1', expect.objectContaining<Partial<ConfirmationRequest>>({
      confirmation_type: 'executed_manually',
      operation_type: 'sell',
      symbol: '510300',
      quantity: 4,
      price: 3.5,
      fees: 1.2,
      executed_at: '2026-05-29T03:00:00Z',
    }))
    expect(screen.queryByText(/自动交易|自动执行交易|一键交易|代下单/)).not.toBeInTheDocument()
  })

  it('rejects negative manual execution fees', () => {
    const onSubmit = vi.fn()
    render(<UserConfirmationPanel decisionId="decision_1" confirmationStatus="pending" availableActions={['executed_manually']} onSubmit={onSubmit} />)

    fireEvent.click(screen.getByRole('button', { name: '已手动执行' }))
    fireEvent.change(screen.getByLabelText('标的代码'), { target: { value: '510300' } })
    fireEvent.change(screen.getByLabelText('线下动作'), { target: { value: 'sell' } })
    fireEvent.change(screen.getByLabelText('数量'), { target: { value: '4' } })
    fireEvent.change(screen.getByLabelText('价格'), { target: { value: '3.5' } })
    fireEvent.change(screen.getByLabelText('费用'), { target: { value: '-1' } })
    fireEvent.change(screen.getByLabelText('执行时间'), { target: { value: '2026-05-29T03:00' } })
    fireEvent.click(screen.getByRole('button', { name: '提交确认' }))

    expect(screen.getByText(/费用必须是非负数字/)).toBeInTheDocument()
    expect(onSubmit).not.toHaveBeenCalled()
  })

  it('submits marked error with enum root cause tag', () => {
    const onSubmit = vi.fn()
    render(<UserConfirmationPanel decisionId="decision_1" confirmationStatus="pending" availableActions={['marked_error']} onSubmit={onSubmit} />)

    fireEvent.click(screen.getByRole('button', { name: '标记错误' }))
    fireEvent.change(screen.getByLabelText('实际结果'), { target: { value: '实际未达预期' } })
    fireEvent.change(screen.getByLabelText('原因标签'), { target: { value: 'rule_threshold_issue' } })
    fireEvent.change(screen.getByLabelText('复盘记录'), { target: { value: '阈值需要再评估' } })
    fireEvent.click(screen.getByRole('button', { name: '提交确认' }))

    expect(onSubmit).toHaveBeenCalledWith('decision_1', expect.objectContaining<Partial<ConfirmationRequest>>({
      confirmation_type: 'marked_error',
      actual_outcome: '实际未达预期',
      root_cause_tag: 'rule_threshold_issue',
      lesson_learned: '阈值需要再评估',
    }))
  })
})
