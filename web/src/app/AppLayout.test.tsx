import { render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { describe, expect, it } from 'vitest'
import { AppLayout } from './AppLayout'

describe('AppLayout', () => {
  it('renders the P110 command-center navigation with local-only safety context', () => {
    render(
      <MemoryRouter initialEntries={['/workbench']}>
        <Routes>
          <Route element={<AppLayout />}>
            <Route path="/workbench" element={<div>工作台内容</div>} />
          </Route>
        </Routes>
      </MemoryRouter>,
    )

    expect(screen.getAllByText('Investment Agent').length).toBeGreaterThanOrEqual(1)
    expect(screen.getAllByText('本地投资纪律工作台').length).toBeGreaterThanOrEqual(1)
    expect(screen.getByText('本地模式')).toBeInTheDocument()
    expect(screen.getByText('离线优先')).toBeInTheDocument()
    expect(screen.getByText('只读导航')).toBeInTheDocument()
    expect(screen.getByText('系统只提示和记录，不会自动执行。')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '今日纪律' })).toHaveAttribute('href', '/')
    expect(screen.getByRole('link', { name: '决策工作台' })).toHaveAttribute('href', '/workbench')
    expect(screen.getByRole('link', { name: '情报与证据' })).toHaveAttribute('href', '/evidence')
    expect(screen.queryByRole('link', { name: /交易|下单|券商/ })).not.toBeInTheDocument()
  })
})
