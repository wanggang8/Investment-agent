import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { cleanup } from '@testing-library/react'
import { APIClientError } from '../services/client'
import { PortfolioPage } from './PortfolioPage'

vi.mock('../services/portfolio', () => ({
  adjustPortfolio: vi.fn(),
  confirmPortfolioImport: vi.fn(),
  correctPortfolioFact: vi.fn(),
  editHolding: vi.fn(),
  getPortfolioCurrent: vi.fn(),
  recordOfflineTransaction: vi.fn(),
  removeHolding: vi.fn(),
  reviewQuarterlyRebalance: vi.fn(),
  validatePortfolioImport: vi.fn(),
}))

import { adjustPortfolio, confirmPortfolioImport, correctPortfolioFact, editHolding, getPortfolioCurrent, recordOfflineTransaction, reviewQuarterlyRebalance, validatePortfolioImport } from '../services/portfolio'

describe('PortfolioPage', () => {
  afterEach(() => {
    cleanup()
    vi.clearAllMocks()
  })

  it('shows API error message instead of a silent empty page', async () => {
    vi.mocked(getPortfolioCurrent).mockRejectedValue(new APIClientError({ requestId: 'rid', code: 'DATA_STALE', message: '本地数据已过期，请刷新后再查看。', displayState: 'insufficient_data' }))

    render(<PortfolioPage />)

    await waitFor(() => expect(screen.getByText('本地数据已过期，请刷新后再查看。')).toBeInTheDocument())
  })

  it('maps missing portfolio snapshot to first-use onboarding instead of generic failure', async () => {
    vi.mocked(getPortfolioCurrent).mockRejectedValue(new APIClientError({ requestId: 'rid_missing', code: 'NOT_FOUND', message: '系统暂时无法处理请求，请稍后重试。', displayState: 'generic_failure' }))

    render(<PortfolioPage />)

    await waitFor(() => expect(screen.getByText('需要初始化本地账户')).toBeInTheDocument())
    expect(screen.getByText('首次初始化')).toBeInTheDocument()
    expect(screen.getByText('录入本地账户与持仓')).toBeInTheDocument()
    expect(screen.queryByText('读取失败')).not.toBeInTheDocument()
    expect(screen.queryByText('系统暂时无法处理请求，请稍后重试。')).not.toBeInTheDocument()
  })

  it('shows empty success state for portfolio page', async () => {
    vi.mocked(getPortfolioCurrent).mockResolvedValue({ request_id: 'rid', data: { snapshot: { snapshot_id: 'snap_empty', snapshot_time: '2026-01-01T00:00:00Z', cash: 100, total_assets: 100, cash_ratio: 1, high_risk_ratio: 0, position_count: 0 }, positions: [] } })

    render(<PortfolioPage />)

    await waitFor(() => expect(screen.getByRole('heading', { name: '组合与持仓维护' })).toBeInTheDocument())
    expect(screen.getByText('组合维护状态')).toBeInTheDocument()
    expect(screen.getByText('需要初始化本地账户')).toBeInTheDocument()
    expect(screen.getByText('首次初始化')).toBeInTheDocument()
    expect(screen.getByText('录入本地账户与持仓')).toBeInTheDocument()
    await waitFor(() => expect(screen.getByText('快照：snap_empty')).toBeInTheDocument())
    expect(screen.getByText('现金：¥100.00')).toBeInTheDocument()
    expect(screen.getByText('现金占比：100.00%')).toBeInTheDocument()
    expect(screen.getByText('暂无持仓结构数据。')).toBeInTheDocument()
    expect(screen.getByText('当前持仓')).toBeInTheDocument()
    expect(screen.getByText('暂无持仓记录。')).toBeInTheDocument()
    expect(screen.queryByText('本地数据已过期，请刷新后再查看。')).not.toBeInTheDocument()
  })

  it('shows buy reason for portfolio positions', async () => {
    vi.mocked(getPortfolioCurrent).mockResolvedValue({
      request_id: 'rid_position',
      data: {
        snapshot: { snapshot_id: 'snap_position', snapshot_time: '2026-01-01T00:00:00Z', cash: 100, total_assets: 130, cash_ratio: 100 / 130, high_risk_ratio: 0, position_count: 1 },
        positions: [{ position_id: 'pos_1', symbol: '510300', name: '沪深300', quantity: 10, cost_price: 2, current_price: 3, market_value: 30, unrealized_profit_ratio: 0.5, position_state: 'normal', buy_reason: '低估分批配置' }],
      },
    })

    render(<PortfolioPage />)

    await waitFor(() => expect(screen.getByText('低估分批配置')).toBeInTheDocument())
    expect(screen.getByText('组合事实可用于纪律评估')).toBeInTheDocument()
    expect(screen.getByText('日常维护')).toBeInTheDocument()
    expect(screen.getByText('校准本地账户事实')).toBeInTheDocument()
    expect(screen.getByText('补记线下交易')).toBeInTheDocument()
    expect(screen.getAllByText('买入理由').length).toBeGreaterThan(0)
  })

  it('prioritizes risk review when high risk ratio is elevated', async () => {
    vi.mocked(getPortfolioCurrent).mockResolvedValue({
      request_id: 'rid_high_risk',
      data: {
        snapshot: { snapshot_id: 'snap_high_risk', snapshot_time: '2026-01-01T00:00:00Z', cash: 100, total_assets: 130, cash_ratio: 100 / 130, high_risk_ratio: 0.46, position_count: 1 },
        positions: [{ position_id: 'pos_1', symbol: '510300', name: '沪深300', quantity: 10, cost_price: 2, current_price: 3, market_value: 30, unrealized_profit_ratio: 0.5, position_state: 'normal', buy_reason: '低估分批配置' }],
      },
    })

    render(<PortfolioPage />)

    await waitFor(() => expect(screen.getByText('高风险仓位需要人工复核')).toBeInTheDocument())
    expect(screen.getByText('高风险复核')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '查看风险预警' })).toHaveAttribute('href', '/risk-alerts')
  })

  it('submits local portfolio calibration without trading endpoint', async () => {
    vi.mocked(getPortfolioCurrent)
      .mockResolvedValueOnce({ request_id: 'rid', data: { snapshot: { snapshot_id: 'snap_empty', snapshot_time: '2026-01-01T00:00:00Z', cash: 100, total_assets: 100, cash_ratio: 1, high_risk_ratio: 0, position_count: 0 }, positions: [] } })
      .mockResolvedValue({ request_id: 'rid_new', data: { snapshot: { snapshot_id: 'snap_new', snapshot_time: '2026-01-02T00:00:00Z', cash: 200, total_assets: 500, cash_ratio: 0.4, high_risk_ratio: 0, position_count: 1 }, positions: [{ position_id: 'pos_new', symbol: '510300', name: '沪深300', quantity: 10, cost_price: 2, current_price: 3, market_value: 30, unrealized_profit_ratio: 0.5, position_state: 'normal' }] } })
    vi.mocked(adjustPortfolio).mockResolvedValue({ request_id: 'adj', data: { snapshot_id: 'snap_new', position_count: 1, audit_event_ids: ['audit_1'] } })

    render(<PortfolioPage />)

    await waitFor(() => expect(screen.getByLabelText('现金')).toBeInTheDocument())
    expect(screen.getByLabelText('现金').getAttribute('aria-describedby')).toContain('portfolio-cash-hint')
    fireEvent.change(screen.getByLabelText('现金'), { target: { value: '200' } })
    fireEvent.change(screen.getByLabelText('总资产'), { target: { value: '500' } })
    fireEvent.change(screen.getByLabelText('标的代码'), { target: { value: '510300' } })
    fireEvent.change(screen.getByLabelText('标的名称'), { target: { value: '沪深300' } })
    fireEvent.change(screen.getByLabelText('数量'), { target: { value: '10' } })
    fireEvent.change(screen.getByLabelText('成本价'), { target: { value: '2' } })
    fireEvent.change(screen.getByLabelText('现价'), { target: { value: '3' } })
    fireEvent.change(screen.getByLabelText('买入日期'), { target: { value: '2026-01-05' } })
    fireEvent.change(screen.getByLabelText('纪律状态'), { target: { value: 'sell_only' } })
    fireEvent.change(screen.getByLabelText('买入理由'), { target: { value: '低估配置' } })
    fireEvent.click(screen.getByRole('button', { name: '保存本地校准' }))

    await waitFor(() => expect(adjustPortfolio).toHaveBeenCalledWith(expect.objectContaining({ cash: 200, total_assets: 500, adjust_reason: '用户本地账户校准' })))
    expect(adjustPortfolio).toHaveBeenCalledWith(expect.objectContaining({ positions: [expect.objectContaining({ buy_date: '2026-01-05', position_state: 'sell_only' })] }))
    expect(screen.getByText('账户校准已保存为本地事实；不会连接交易接口。')).toBeInTheDocument()
    await waitFor(() => expect(screen.getByText('快照：snap_new')).toBeInTheDocument())
    expect(screen.getByText('现金：¥200.00')).toBeInTheDocument()
    expect(screen.getByText('现金占比：40.00%')).toBeInTheDocument()
    expect(screen.getByText('50.00%')).toBeInTheDocument()
    expect(getPortfolioCurrent).toHaveBeenCalledTimes(2)
    expect(screen.getAllByText(/交易接口/).length).toBeGreaterThan(0)
  })

  it('shows onboarding fields and submits local account facts with buy reason', async () => {
    vi.mocked(getPortfolioCurrent).mockResolvedValue({ request_id: 'rid', data: { snapshot: { snapshot_id: 'snap_empty', snapshot_time: '2026-01-01T00:00:00Z', cash: 100, total_assets: 100, cash_ratio: 1, high_risk_ratio: 0, position_count: 0 }, positions: [] } })
    vi.mocked(adjustPortfolio).mockResolvedValue({ request_id: 'adj', data: { snapshot_id: 'snap_new', position_count: 1, audit_event_ids: ['audit_1'] } })

    render(<PortfolioPage />)

    await waitFor(() => expect(screen.getByText('首次使用引导')).toBeInTheDocument())
    fireEvent.change(screen.getByLabelText('标的代码'), { target: { value: '510300' } })
    fireEvent.change(screen.getByLabelText('标的名称'), { target: { value: '沪深300ETF' } })
    fireEvent.change(screen.getByLabelText('买入理由'), { target: { value: '低估配置' } })
    fireEvent.change(screen.getByLabelText('现金'), { target: { value: '70' } })
    fireEvent.change(screen.getByLabelText('总资产'), { target: { value: '100' } })
    fireEvent.change(screen.getByLabelText('数量'), { target: { value: '10' } })
    fireEvent.change(screen.getByLabelText('成本价'), { target: { value: '2' } })
    fireEvent.change(screen.getByLabelText('现价'), { target: { value: '3' } })
    fireEvent.change(screen.getByLabelText('买入日期'), { target: { value: '2026-01-06' } })
    fireEvent.change(screen.getByLabelText('纪律状态'), { target: { value: 'frozen_watch' } })
    fireEvent.click(screen.getByRole('button', { name: '保存本地校准' }))

    await waitFor(() => expect(adjustPortfolio).toHaveBeenCalledWith(expect.objectContaining({ positions: [expect.objectContaining({ buy_reason: '低估配置', buy_date: '2026-01-06', position_state: 'frozen_watch' })] })))
  })

  it('does not submit local facts from empty default placeholders', async () => {
    vi.mocked(getPortfolioCurrent).mockResolvedValue({ request_id: 'rid', data: { snapshot: { snapshot_id: 'snap_empty', snapshot_time: '2026-01-01T00:00:00Z', cash: 100, total_assets: 100, cash_ratio: 1, high_risk_ratio: 0, position_count: 0 }, positions: [] } })

    render(<PortfolioPage />)

    await waitFor(() => expect(screen.getByText('首次使用引导')).toBeInTheDocument())
    fireEvent.click(screen.getByRole('button', { name: '保存本地校准' }))
    fireEvent.click(screen.getByRole('button', { name: '记录线下交易' }))
    fireEvent.click(screen.getByRole('button', { name: '校验批量导入' }))

    expect(adjustPortfolio).not.toHaveBeenCalled()
    expect(recordOfflineTransaction).not.toHaveBeenCalled()
    expect(validatePortfolioImport).not.toHaveBeenCalled()
    expect(screen.getByText('请先填写账户和持仓必填信息。')).toBeInTheDocument()
  })

  it('keeps import confirmation disabled until validation succeeds', async () => {
    vi.mocked(getPortfolioCurrent).mockResolvedValue({ request_id: 'rid_position', data: { snapshot: { snapshot_id: 'snap_position', snapshot_time: '2026-01-01T00:00:00Z', cash: 100, total_assets: 130, cash_ratio: 100 / 130, high_risk_ratio: 0, position_count: 1 }, positions: [{ position_id: 'pos_1', symbol: '510300', name: '沪深300ETF', quantity: 10, cost_price: 2, current_price: 3, market_value: 30, unrealized_profit_ratio: 0.5, position_state: 'normal', buy_reason: '低估配置' }] } })
    vi.mocked(validatePortfolioImport).mockResolvedValue({ request_id: 'validate', data: { import_batch_id: 'batch_bad', summary: { row_count: 1, valid_count: 0, invalid_count: 1 }, rows: [{ row_number: 1, valid: false, errors: ['symbol 不能为空'] }] } })

    render(<PortfolioPage />)

    const confirmButton = await screen.findByRole('button', { name: '确认批量导入' })
    expect(confirmButton).toBeDisabled()
    fireEvent.click(screen.getByRole('button', { name: '校验批量导入' }))
    await waitFor(() => expect(validatePortfolioImport).toHaveBeenCalled())
    expect(confirmButton).toBeDisabled()
    fireEvent.click(confirmButton)
    expect(confirmPortfolioImport).not.toHaveBeenCalled()
  })

  it('submits holding edit, offline transaction, import validation, import confirm, and correction actions', async () => {
    vi.mocked(getPortfolioCurrent).mockResolvedValue({
      request_id: 'rid_position',
      data: {
        snapshot: { snapshot_id: 'snap_position', snapshot_time: '2026-01-01T00:00:00Z', cash: 100, total_assets: 130, cash_ratio: 100 / 130, high_risk_ratio: 0, position_count: 1 },
        positions: [{ position_id: 'pos_1', symbol: '510300', name: '沪深300ETF', quantity: 10, cost_price: 2, current_price: 3, market_value: 30, unrealized_profit_ratio: 0.5, position_state: 'sell_only', buy_date: '2026-01-05', buy_reason: '低估配置' }],
      },
    })
    vi.mocked(editHolding).mockResolvedValue({ request_id: 'edit', data: { snapshot_id: 'snap_edit', position_id: 'pos_1', audit_event_ids: ['audit_edit'], safety_statement: '不连接券商、不自动交易。' } })
    vi.mocked(recordOfflineTransaction).mockResolvedValue({ request_id: 'tx', data: { snapshot_id: 'snap_tx', transaction_id: 'tx_1', audit_event_ids: ['audit_tx'], safety_statement: '不连接券商、不自动交易。' } })
    vi.mocked(validatePortfolioImport).mockResolvedValue({ request_id: 'validate', data: { import_batch_id: 'batch_1', summary: { row_count: 1, valid_count: 1, invalid_count: 0 }, rows: [{ row_number: 1, valid: true }] } })
    vi.mocked(confirmPortfolioImport).mockResolvedValue({ request_id: 'confirm', data: { snapshot_id: 'snap_import', import_batch_id: 'batch_1', audit_event_ids: ['audit_import'], safety_statement: '不连接券商、不自动交易。' } })
    vi.mocked(correctPortfolioFact).mockResolvedValue({ request_id: 'corr', data: { correction_id: 'corr_1', audit_event_ids: ['audit_corr'], safety_statement: '不连接券商、不自动交易。' } })

    render(<PortfolioPage />)

    await waitFor(() => expect(screen.getByText('线下交易记录')).toBeInTheDocument())
    expect(screen.getByText('用于补记已由用户自行完成的线下买入、卖出或减仓。')).toBeInTheDocument()
    expect(screen.getByRole('option', { name: '卖出' })).toBeInTheDocument()
    expect(screen.getByRole('option', { name: '减仓' })).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: '保存持仓编辑' }))
    fireEvent.change(screen.getByLabelText('线下交易类型'), { target: { value: 'sell' } })
    fireEvent.click(screen.getByRole('button', { name: '记录线下交易' }))
    fireEvent.click(screen.getByRole('button', { name: '校验批量导入' }))
    await waitFor(() => expect(validatePortfolioImport).toHaveBeenCalled())
    fireEvent.click(screen.getByRole('button', { name: '确认批量导入' }))
    fireEvent.click(screen.getByRole('button', { name: '记录修正审计' }))

    expect(screen.getByLabelText('纪律状态')).toHaveValue('sell_only')
    await waitFor(() => expect(editHolding).toHaveBeenCalledWith(expect.objectContaining({ position_id: 'pos_1', reason: '用户本地持仓编辑', position: expect.objectContaining({ position_state: 'sell_only', buy_date: '2026-01-05' }) })))
    expect(recordOfflineTransaction).toHaveBeenCalledWith(expect.objectContaining({ operation_type: 'sell', symbol: '510300', fees: 0, executed_at: expect.any(String), note: '用户补记线下交易' }))
    expect(confirmPortfolioImport).toHaveBeenCalledWith(expect.objectContaining({ import_batch_id: 'batch_1', confirm_reason: '确认导入本地账户事实' }))
    expect(correctPortfolioFact).toHaveBeenCalledWith(expect.objectContaining({ target_type: 'position', target_id: 'pos_1', before_json: expect.any(String), after_json: expect.any(String), correction_reason: '用户更正本地事实' }))
  })

  it('runs quarterly rebalance review through the portfolio UI without trading calls', async () => {
    vi.mocked(getPortfolioCurrent).mockResolvedValue({
      request_id: 'rid_rebalance',
      data: {
        snapshot: { snapshot_id: 'snap_rebalance', snapshot_time: '2026-06-22T00:00:00Z', cash: 100, total_assets: 1000, cash_ratio: 0.1, high_risk_ratio: 0, position_count: 2 },
        positions: [
          { position_id: 'pos_core', symbol: '510300', name: '沪深300ETF', quantity: 100, cost_price: 3, current_price: 6, market_value: 600, unrealized_profit_ratio: 1, position_state: 'normal', buy_reason: '核心配置', asset_tag: 'core' },
          { position_id: 'pos_sat', symbol: '159915', name: '创业板ETF', quantity: 100, cost_price: 2, current_price: 3, market_value: 300, unrealized_profit_ratio: 0.5, position_state: 'normal', buy_reason: '卫星配置', asset_tag: 'satellite' },
        ],
      },
    })
    vi.mocked(reviewQuarterlyRebalance).mockResolvedValue({
      request_id: 'rebalance',
      data: {
        review_id: 'rebalance_1',
        review_date: '2026-06-22',
        total_assets: 1000,
        drift_threshold: 0.15,
        audit_event_ids: ['audit_rebalance'],
        safety_statement: '季度再平衡仅生成人工计划金额，不连接券商、不自动交易、不创建订单。',
        items: [
          { bucket: 'core', target_ratio: 0.5, actual_ratio: 0.6, drift_ratio: 0.1, target_value: 500, actual_value: 600, recommendation: 'hold', manual_amount: 0 },
          { bucket: 'satellite', target_ratio: 0.2, actual_ratio: 0.3, drift_ratio: 0.1, target_value: 200, actual_value: 300, recommendation: 'hold', manual_amount: 0 },
          { bucket: 'cash', target_ratio: 0.3, actual_ratio: 0.1, drift_ratio: -0.2, target_value: 300, actual_value: 100, recommendation: 'manual_raise_cash', manual_amount: 200 },
        ],
      },
    })

    render(<PortfolioPage />)

    const button = await screen.findByRole('button', { name: '运行季度再平衡复核' })
    fireEvent.click(button)

    await waitFor(() => expect(reviewQuarterlyRebalance).toHaveBeenCalledWith(expect.objectContaining({ target_core_ratio: 0.5, target_satellite_ratio: 0.2, target_cash_ratio: 0.3, drift_threshold: 0.15 })))
    expect(screen.getByText('季度再平衡复核已生成，仅作为人工计划。')).toBeInTheDocument()
    expect(screen.getByText(/现金：目标 30.00%，实际 10.00%/)).toBeInTheDocument()
    expect(screen.getByText('季度再平衡仅生成人工计划金额，不连接券商、不自动交易、不创建订单。')).toBeInTheDocument()
    expect(recordOfflineTransaction).not.toHaveBeenCalled()
  })
})
