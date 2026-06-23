import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { RuleProposalPanel } from './RuleProposalPanel'

import type { RuleProposal } from '../../types/rule'

const proposal: RuleProposal = {
  proposal_id: 'proposal_1',
  proposal_type: 'threshold',
  status: 'pending_final_confirm',
  title: '调整观察区阈值',
  proposal_version: 'v1',
  sample_count: 6,
  created_at: '2026-05-31T00:00:00Z',
}

describe('RuleProposalPanel', () => {
  afterEach(() => cleanup())

  it('shows final confirmation without auto applying rules', () => {
    const onFinalConfirm = vi.fn()
    render(<RuleProposalPanel proposals={[proposal]} onFinalConfirm={onFinalConfirm} />)

    expect(screen.getByText(/待最终确认/)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '确认应用到正式规则' })).toBeInTheDocument()
    expect(screen.queryByText(/自动应用|立即生效/)).not.toBeInTheDocument()
    expect(screen.getByText(/正式规则不会自动生效/)).toBeInTheDocument()
  })

  it('shows unknown proposal status safely', () => {
    render(<RuleProposalPanel proposals={[{ ...proposal, status: 'backend_new_state' as RuleProposal['status'] }]} />)

    expect(screen.getByText(/状态：未知状态/)).toBeInTheDocument()
    expect(screen.queryByText(/backend_new_state/)).not.toBeInTheDocument()
  })

  it('renders audit result mapping and rule content fields', () => {
    render(<RuleProposalPanel proposals={[{ ...proposal, audit_result: 'needs_user_review', before_rule: { content: '旧规则文本', id: 'r1' }, after_rule: { content: '新规则文本', id: 'r1' } }]} />)

    expect(screen.getByText('守门人结果：需要用户复核')).toBeInTheDocument()
    expect(screen.queryByText('needs_user_review')).not.toBeInTheDocument()
    expect(screen.getByLabelText('变更前规则')).toHaveTextContent('旧规则文本')
    expect(screen.getByLabelText('变更后规则')).toHaveTextContent('新规则文本')
    expect(screen.getByLabelText('变更前规则')).not.toHaveTextContent('"id"')
  })
})
