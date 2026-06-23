import { cleanup, render, screen, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { DailyDisciplineReportsPage } from './DailyDisciplineReportsPage'

vi.mock('../services/dailyDisciplineReport', () => ({
  listDailyDisciplineReports: vi.fn(),
}))

import { listDailyDisciplineReports } from '../services/dailyDisciplineReport'

describe('DailyDisciplineReportsPage', () => {
  afterEach(() => {
    cleanup()
    vi.clearAllMocks()
  })

  it('shows report history with status labels and encoded detail links', async () => {
    vi.mocked(listDailyDisciplineReports).mockResolvedValue({ request_id: 'rid_list', data: { reports: [
      {
        report_id: 'daily/report:2026-06-08:ETF 300',
        local_date: '2026-06-08',
        scope: 'holdings',
        status: 'insufficient_data',
        summary: '证据不足，等待补齐持仓数据。',
        evidence: { evidence_count: 0, independent_source_count: 0, high_grade_independent_source_count: 0 },
        trend: { success_count: 1, degraded_count: 0, failed_count: 0, insufficient_data_count: 1 },
        safety_note: '不会自动执行交易，需人工复核。',
      },
    ] } })

    render(<MemoryRouter><DailyDisciplineReportsPage /></MemoryRouter>)

    await waitFor(() => expect(screen.getByRole('heading', { name: '每日纪律报告历史' })).toBeInTheDocument())
    expect(screen.getByText('每日纪律复盘状态')).toBeInTheDocument()
    expect(screen.getByText('每日纪律与自动运行需要检查')).toBeInTheDocument()
    expect(screen.getByText('最新报告')).toBeInTheDocument()
    expect(screen.getByText('补齐本地持仓')).toBeInTheDocument()
    expect(screen.getByText('查看报告详情')).toBeInTheDocument()
    expect(screen.getAllByText('数据不足').length).toBeGreaterThan(0)
    expect(screen.getAllByText('2026-06-08').length).toBeGreaterThan(0)
    expect(screen.getByText('证据不足，等待补齐持仓数据。')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '前往账户初始化' })).toHaveAttribute('href', '/positions')
    expect(screen.getByRole('link', { name: '查看报告' })).toHaveAttribute('href', '/daily-discipline/reports/daily%2Freport%3A2026-06-08%3AETF%20300')
  })

  it('shows empty state when no reports exist', async () => {
    vi.mocked(listDailyDisciplineReports).mockResolvedValue({ request_id: 'rid_empty', data: { reports: [] } })

    render(<MemoryRouter><DailyDisciplineReportsPage /></MemoryRouter>)

    await waitFor(() => expect(screen.getByText('每日纪律复盘状态')).toBeInTheDocument())
    await waitFor(() => expect(screen.getByText('暂无每日纪律报告')).toBeInTheDocument())
  })

  it('shows safe load failure copy without leaking raw errors', async () => {
    vi.mocked(listDailyDisciplineReports).mockRejectedValue(new Error('internal stack trace'))

    render(<MemoryRouter><DailyDisciplineReportsPage /></MemoryRouter>)

    await waitFor(() => expect(screen.getByRole('alert')).toHaveTextContent('每日纪律报告加载失败'))
    expect(screen.queryByText(/internal stack trace/)).not.toBeInTheDocument()
  })
})
