import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'
import { FinalVerdictCard } from './FinalVerdictCard'

describe('FinalVerdictCard', () => {
  afterEach(() => cleanup())

  it('renders nullable action arrays without crashing', () => {
    render(
      <FinalVerdictCard
        summary={{
          verdict: '暂停交易类建议',
          final_verdict_status: 'insufficient_data',
          prohibited_actions: null,
          optional_actions: null,
          action_required: false,
          confirmation_status: 'not_required',
        }}
      />,
    )

    expect(screen.getByRole('heading', { name: '暂停交易类建议' })).toBeInTheDocument()
    expect(screen.getByText('暂无新增禁止事项。')).toBeInTheDocument()
    expect(screen.getByText('暂无可选记录。')).toBeInTheDocument()
  })
})
