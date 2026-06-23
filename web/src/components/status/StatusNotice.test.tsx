import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { StatusNotice } from './StatusNotice'
import { getStatusNoticeCopy } from './statusNoticeCopy'

describe('StatusNotice', () => {
  it('shows distinct messages for data stale analyst and vector degradation codes', () => {
    expect(getStatusNoticeCopy('insufficient_data', 'DATA_STALE').title).toBe('数据已过期')
    expect(getStatusNoticeCopy('insufficient_data', 'ANALYST_UNAVAILABLE').title).toBe('分析服务降级')
    expect(getStatusNoticeCopy('insufficient_data', 'VECTOR_INDEX_UNAVAILABLE').title).toBe('索引不可用')
  })

  it('renders high risk as a dedicated state', () => {
    render(<StatusNotice state="high_risk" />)

    expect(screen.getByText('高危状态')).toBeInTheDocument()
    expect(screen.getByText(/禁止新增买入/)).toBeInTheDocument()
  })

  it('renders safe message without exposing detail-shaped content', () => {
    render(<StatusNotice state="generic_failure" safeMessage="系统暂时无法处理请求，请稍后重试。" />)

    expect(screen.getByText('系统暂时无法处理请求，请稍后重试。')).toBeInTheDocument()
    expect(screen.queryByText(/SQL|stack|\/tmp/)).not.toBeInTheDocument()
  })
})
