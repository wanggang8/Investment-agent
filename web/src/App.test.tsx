import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'
import App from './App'

function renderAt(path: string) {
  window.history.pushState({}, '', path)
  return render(<App />)
}

describe('App productized fallback routes', () => {
  afterEach(() => {
    cleanup()
  })

  it('renders the decisions index instead of a blank screen', () => {
    renderAt('/decisions')

    expect(screen.getByRole('heading', { name: '决策详情' })).toBeInTheDocument()
    expect(screen.getByText('先选择一条本地决策记录')).toBeInTheDocument()
    expect(screen.getByText('只读入口')).toBeInTheDocument()
  })

  it('renders productized API diagnostics instead of a raw 404', () => {
    renderAt('/api-diagnostics')

    expect(screen.getByRole('heading', { name: '接口诊断' })).toBeInTheDocument()
    expect(screen.getByText('查看本地接口与页面状态')).toBeInTheDocument()
    expect(screen.getByText('不会触发后台动作')).toBeInTheDocument()
  })
})
