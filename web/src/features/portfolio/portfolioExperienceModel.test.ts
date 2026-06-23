import { describe, expect, it } from 'vitest'
import { buildPortfolioExperienceModel } from './portfolioExperienceModel'
import type { PortfolioCurrentResponse } from '../../types/portfolio'

const emptyPortfolio: PortfolioCurrentResponse = {
  snapshot: {
    snapshot_id: 'snap_empty',
    snapshot_time: '2026-06-18T09:30:00Z',
    cash: 100,
    total_assets: 100,
    cash_ratio: 1,
    high_risk_ratio: 0,
    position_count: 0,
  },
  positions: [],
}

const investedPortfolio: PortfolioCurrentResponse = {
  snapshot: {
    snapshot_id: 'snap_position',
    snapshot_time: '2026-06-18T09:30:00Z',
    cash: 32000,
    total_assets: 100000,
    cash_ratio: 0.32,
    high_risk_ratio: 0.12,
    position_count: 2,
  },
  positions: [
    {
      position_id: 'pos_510300',
      symbol: '510300',
      name: '沪深300ETF',
      quantity: 100,
      cost_price: 4,
      current_price: 4.2,
      market_value: 420,
      unrealized_profit_ratio: 0.05,
      position_state: 'normal',
      buy_reason: '长期配置',
      asset_tag: 'core',
    },
  ],
}

describe('buildPortfolioExperienceModel', () => {
  it('guides first-use portfolio initialization before maintenance actions', () => {
    const model = buildPortfolioExperienceModel({ portfolio: emptyPortfolio })

    expect(model.statusLabel).toBe('需要初始化本地账户')
    expect(model.statusTone).toBe('warning')
    expect(model.stageLabel).toBe('首次初始化')
    expect(model.nextActions.map((action) => action.label)).toContain('录入本地账户与持仓')
    expect(model.summaryMetrics).toContainEqual(expect.objectContaining({ label: '持仓数量', value: '0' }))
    expect(model.safetyNotes.join(' ')).toContain('不连接券商')
  })

  it('summarizes a maintainable portfolio with local-only actions', () => {
    const model = buildPortfolioExperienceModel({ portfolio: investedPortfolio })

    expect(model.statusLabel).toBe('组合事实可用于纪律评估')
    expect(model.statusTone).toBe('success')
    expect(model.stageLabel).toBe('日常维护')
    expect(model.summaryMetrics).toEqual(expect.arrayContaining([
      expect.objectContaining({ label: '总资产', value: '¥100,000.00' }),
      expect.objectContaining({ label: '现金占比', value: '32.00%' }),
      expect.objectContaining({ label: '高风险比例', value: '12.00%' }),
    ]))
    expect(model.nextActions.map((action) => action.label)).toEqual(expect.arrayContaining(['校准本地账户事实', '补记线下交易']))
    expect(model.maintenanceModes.map((mode) => mode.label)).toEqual(['初始化/校准', '持仓维护', '线下交易记录', '批量导入', '错误修正'])
  })

  it('does not describe high risk ratio as normal success', () => {
    const model = buildPortfolioExperienceModel({
      portfolio: {
        ...investedPortfolio,
        snapshot: { ...investedPortfolio.snapshot, high_risk_ratio: 0.46 },
      },
    })

    expect(model.statusTone).toBe('danger')
    expect(model.statusLabel).toBe('高风险仓位需要人工复核')
    expect(model.nextActions.map((action) => action.label)).toContain('查看风险预警')
  })

  it('surfaces import confirmation as a dedicated local fact stage', () => {
    const model = buildPortfolioExperienceModel({ portfolio: investedPortfolio, importReady: true, importBatchID: 'import_1' })

    expect(model.stageLabel).toBe('导入待确认')
    expect(model.nextActions[0]).toMatchObject({ label: '确认批量导入', priority: 'blocking' })
    expect(model.safetyNotes.join(' ')).toContain('本地事实')
  })
})
