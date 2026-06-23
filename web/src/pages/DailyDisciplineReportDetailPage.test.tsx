import { cleanup, render, screen, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { DailyDisciplineReportDetailPage } from './DailyDisciplineReportDetailPage'

vi.mock('../services/dailyDisciplineReport', () => ({
  getDailyDisciplineReport: vi.fn(),
}))

import { getDailyDisciplineReport } from '../services/dailyDisciplineReport'

describe('DailyDisciplineReportDetailPage', () => {
  afterEach(() => {
    cleanup()
    vi.clearAllMocks()
  })

  it('shows report summary, evidence, tracking links and safety note', async () => {
    vi.mocked(getDailyDisciplineReport).mockResolvedValue({ request_id: 'rid_detail', data: {
      report_id: 'daily_report:2026-06-08:holdings:v1',
      local_date: '2026-06-08',
      scope: 'holdings',
      status: 'success',
      summary: '今日纪律报告完成，维持观察。',
      decision_link: '/decisions/decision_8',
      audit_link: '/audit?input_ref=run_8',
      notification_link: '/notifications?source_id=run_8',
      auto_run_link: '/daily-auto-run?run_id=run_8',
      final_verdict: '维持观察',
      verdict_status: 'hold',
      missing_categories: ['holdings', 'market_snapshot'],
      p34_source_coverage: {
        summary: 'P34 可用扩展数据：index_constituents、sentiment_proxy',
        missing_categories: ['index_valuation_files'],
        source_health: [
          { source_name: 'csindex', source_level: 'A', source_type: 'index_basic', data_category: 'index_constituents', freshness: 'fresh', data_date: '2026-06-05', affected_symbols: ['000300'] },
          { source_name: 'csindex', source_level: 'A', source_type: 'index_basic', data_category: 'index_valuation_files', freshness: 'parse_error', data_date: '2026-06-05', affected_symbols: ['000300'] },
        ],
      },
      risk_alerts: [
        {
          alert_id: 'risk_detail_1',
          risk_type: 'data_degraded',
          severity: 'critical',
          sop_status: 'escalated',
          symbol: '510300',
          trigger_summary: '数据源新鲜度降级',
          prohibited_actions: ['交易类建议'],
          suggested_actions: ['检查数据源'],
          link: '/risk-alerts/risk_detail_1',
          safety_note: '风险预警只用于本地人工复核，不会自动交易。',
          created_at: '2026-06-08T02:00:00Z',
          updated_at: '2026-06-08T02:00:00Z',
        },
      ],
      evidence: { evidence_count: 6, independent_source_count: 3, high_grade_independent_source_count: 2 },
      trend: { success_count: 4, degraded_count: 1, failed_count: 0, insufficient_data_count: 2 },
      safety_note: '不会自动执行交易，所有动作需人工复核。',
    } })

    render(
      <MemoryRouter initialEntries={['/daily-discipline/reports/daily_report%3A2026-06-08%3Aholdings%3Av1']}>
        <Routes>
          <Route path="/daily-discipline/reports/:reportId" element={<DailyDisciplineReportDetailPage />} />
        </Routes>
      </MemoryRouter>,
    )

    await waitFor(() => expect(getDailyDisciplineReport).toHaveBeenCalledWith('daily_report:2026-06-08:holdings:v1'))
    expect(screen.getByRole('heading', { name: '每日纪律报告详情' })).toBeInTheDocument()
    expect(screen.getByText('今日纪律报告完成，维持观察。')).toBeInTheDocument()
    expect(screen.getAllByText('成功').length).toBeGreaterThan(0)
    expect(screen.getByText('2026-06-08')).toBeInTheDocument()
    expect(screen.getAllByText(/维持观察/).length).toBeGreaterThan(0)
    expect(screen.getByText('6')).toBeInTheDocument()
    expect(screen.getByText('3')).toBeInTheDocument()
    expect(screen.getAllByText('2').length).toBeGreaterThan(0)
    expect(screen.getByText('holdings')).toBeInTheDocument()
    expect(screen.getByText('market_snapshot')).toBeInTheDocument()
    expect(screen.getByText('P34 扩展数据覆盖')).toBeInTheDocument()
    expect(screen.getByText('P34 可用扩展数据：index_constituents、sentiment_proxy')).toBeInTheDocument()
    expect(screen.getByText('index_valuation_files')).toBeInTheDocument()
    expect(screen.getByText(/csindex · 指数样本 · 新鲜；数据日：2026-06-05；等级：A；影响标的：000300/)).toBeInTheDocument()
    expect(screen.getByText(/csindex · 指数估值文件 · 解析失败；数据日：2026-06-05；等级：A；影响标的：000300/)).toBeInTheDocument()
    expect(screen.getByText('风险预警')).toBeInTheDocument()
    expect(screen.getByText('数据源新鲜度降级')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '查看风险预警' })).toHaveAttribute('href', '/risk-alerts/risk_detail_1')
    expect(screen.getByText('不会自动执行交易，所有动作需人工复核。')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '查看决策详情' })).toHaveAttribute('href', '/decisions/decision_8')
    expect(screen.getByRole('link', { name: '查看审计详情' })).toHaveAttribute('href', '/audit?input_ref=run_8')
    expect(screen.getByRole('link', { name: '查看通知' })).toHaveAttribute('href', '/notifications?source_id=run_8')
    expect(screen.getByRole('link', { name: '查看每日自动运行' })).toHaveAttribute('href', '/daily-auto-run?run_id=run_8')
  })

  it('shows safe failure copy without leaking raw report failure reason', async () => {
    vi.mocked(getDailyDisciplineReport).mockResolvedValue({ request_id: 'rid_failed', data: {
      report_id: 'daily_report:2026-06-08:holdings:v1',
      local_date: '2026-06-08',
      scope: 'holdings',
      status: 'failed',
      summary: '每日纪律报告暂时无法生成。',
      audit_link: '/audit?input_ref=run_failed',
      failure_code: 'missing_prerequisites',
      failure_reason: 'sqlite /tmp/internal stack trace',
      evidence: { evidence_count: 0, independent_source_count: 0, high_grade_independent_source_count: 0 },
      trend: { success_count: 0, degraded_count: 0, failed_count: 1, insufficient_data_count: 0 },
      safety_note: '不会自动执行交易，所有动作需人工复核。',
    } })

    render(
      <MemoryRouter initialEntries={['/daily-discipline/reports/daily_report%3A2026-06-08%3Aholdings%3Av1']}>
        <Routes>
          <Route path="/daily-discipline/reports/:reportId" element={<DailyDisciplineReportDetailPage />} />
        </Routes>
      </MemoryRouter>,
    )

    await waitFor(() => expect(getDailyDisciplineReport).toHaveBeenCalledWith('daily_report:2026-06-08:holdings:v1'))
    expect(screen.queryByText(/sqlite \/tmp\/internal stack trace/)).not.toBeInTheDocument()
    expect(screen.getByText('缺少本地账户或持仓。')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '前往账户初始化' })).toHaveAttribute('href', '/positions')
    expect(screen.getByRole('link', { name: '查看审计详情' })).toHaveAttribute('href', '/audit?input_ref=run_failed')
  })

  it('shows safe load failure copy without leaking raw errors', async () => {
    vi.mocked(getDailyDisciplineReport).mockRejectedValue(new Error('internal stack trace'))

    render(
      <MemoryRouter initialEntries={['/daily-discipline/reports/daily_report%3A2026-06-08%3Aholdings%3Av1']}>
        <Routes>
          <Route path="/daily-discipline/reports/:reportId" element={<DailyDisciplineReportDetailPage />} />
        </Routes>
      </MemoryRouter>,
    )

    await waitFor(() => expect(screen.getByRole('alert')).toHaveTextContent('每日纪律报告详情加载失败'))
    expect(screen.queryByText(/internal stack trace/)).not.toBeInTheDocument()
  })
})
