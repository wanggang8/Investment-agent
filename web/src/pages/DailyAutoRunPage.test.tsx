import { cleanup, render, screen, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { DailyAutoRunPage } from './DailyAutoRunPage'

vi.mock('../services/dailyAutoRun', () => ({
  getDailyAutoRunStatus: vi.fn(),
}))

import { getDailyAutoRunStatus } from '../services/dailyAutoRun'

describe('DailyAutoRunPage', () => {
  afterEach(() => {
    cleanup()
    vi.clearAllMocks()
  })

  it('shows failed status, missing prerequisites and tracking links', async () => {
    vi.mocked(getDailyAutoRunStatus).mockResolvedValue({ request_id: 'rid_auto', data: {
      enabled: true,
      run_time: '08:30',
      timezone: 'Asia/Shanghai',
      scope: 'holdings',
      status: 'failed',
      last_run_at: '2026-06-07T00:30:00Z',
      next_run_at: '2026-06-08T00:30:00Z',
      failure_code: 'missing_prerequisites',
      failure_reason: '缺少本地持仓',
      latest_decision_link: '/decisions?request_id=auto_run_1',
      latest_notification_link: '/notifications?source_id=key_1',
      latest_audit_link: '/audit?input_ref=key_1',
      missing_action: '请先录入本地账户、组合和当前持仓，再等待下一次自动运行或手动触发。',
      safety_note: '仅记录本地每日刷新和纪律评估结果，不会自动执行交易，需人工复核。',
    } })

    render(<DailyAutoRunPage />)

    await waitFor(() => expect(screen.getAllByText('失败').length).toBeGreaterThan(0))
    expect(screen.getByText('每日自动运行健康')).toBeInTheDocument()
    expect(screen.getByText('每日纪律与自动运行需要检查')).toBeInTheDocument()
    expect(screen.getByText('补齐本地持仓')).toBeInTheDocument()
    expect(screen.getByText('查看审计记录')).toBeInTheDocument()
    expect(screen.getByText('缺少本地账户或持仓。')).toBeInTheDocument()
    expect(screen.getByText('请先录入本地账户、组合和当前持仓，再等待下一次自动运行或手动触发。')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '查看最新每日决策' })).toHaveAttribute('href', '/decisions?request_id=auto_run_1')
    expect(screen.getByRole('link', { name: '查看通知' })).toHaveAttribute('href', '/notifications?source_id=key_1')
    expect(screen.getByRole('link', { name: '查看审计详情' })).toHaveAttribute('href', '/audit?input_ref=key_1')
    expect(screen.getByText('仅记录本地每日刷新和纪律评估结果，不会自动执行交易，需人工复核。')).toBeInTheDocument()
  })

  it('maps diagnostic failure reasons to safe copy', async () => {
    vi.mocked(getDailyAutoRunStatus).mockResolvedValue({ request_id: 'rid_auto_raw', data: {
      enabled: true,
      status: 'failed',
      failure_code: 'market_refresh_failed',
      failure_reason: 'sqlite /tmp/internal stack trace',
      safety_note: '仅记录本地每日刷新和纪律评估结果，不会自动执行交易，需人工复核。',
    } })

    render(<DailyAutoRunPage />)

    await waitFor(() => expect(screen.getByText('每日自动运行暂时无法完成，请查看审计记录。')).toBeInTheDocument())
    expect(screen.queryByText(/sqlite \/tmp\/internal stack trace/)).not.toBeInTheDocument()
  })

  it('shows disabled state without implying automatic trading', async () => {
    vi.mocked(getDailyAutoRunStatus).mockResolvedValue({ request_id: 'rid_disabled', data: {
      enabled: false,
      status: 'disabled',
      safety_note: '仅记录本地每日刷新和纪律评估结果，不会自动执行交易，需人工复核。',
    } })

    render(<DailyAutoRunPage />)

    await waitFor(() => expect(screen.getAllByText('关闭').length).toBeGreaterThan(0))
    expect(screen.getByText('每日自动运行健康')).toBeInTheDocument()
    expect(screen.getByText('每日自动运行未启用。')).toBeInTheDocument()
    expect(screen.getByText('仅记录本地每日刷新和纪律评估结果，不会自动执行交易，需人工复核。')).toBeInTheDocument()
  })

  it('shows safe load failure copy without leaking raw errors', async () => {
    vi.mocked(getDailyAutoRunStatus).mockRejectedValue(new Error('internal stack trace'))

    render(<DailyAutoRunPage />)

    await waitFor(() => expect(screen.getByRole('alert')).toHaveTextContent('每日自动运行状态加载失败'))
    expect(screen.queryByText(/internal stack trace/)).not.toBeInTheDocument()
  })
})
