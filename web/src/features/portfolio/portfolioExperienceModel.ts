import type { PageErrorState } from '../../shared/utils'
import { formatCurrency, formatPercent } from '../../shared/utils'
import type { PortfolioCurrentResponse } from '../../types/portfolio'

export type ExperienceTone = 'success' | 'warning' | 'danger' | 'unknown'
export type ActionPriority = 'blocking' | 'review' | 'follow_up'

export interface PortfolioExperienceAction {
  label: string
  detail: string
  priority: ActionPriority
  href?: string
}

export interface PortfolioExperienceMetric {
  label: string
  value: string
  tone?: ExperienceTone
}

export interface PortfolioMaintenanceMode {
  id: string
  label: string
  description: string
}

export interface PortfolioExperienceModel {
  statusLabel: string
  statusTone: ExperienceTone
  stageLabel: string
  stageDescription: string
  summaryMetrics: PortfolioExperienceMetric[]
  nextActions: PortfolioExperienceAction[]
  maintenanceModes: PortfolioMaintenanceMode[]
  safetyNotes: string[]
  warnings: string[]
}

export interface PortfolioExperienceInput {
  portfolio?: PortfolioCurrentResponse
  importReady?: boolean
  importBatchID?: string
  error?: PageErrorState
}

export function buildPortfolioExperienceModel(input: PortfolioExperienceInput): PortfolioExperienceModel {
  const snapshot = input.portfolio?.snapshot
  const positions = input.portfolio?.positions ?? []
  const warnings = input.error ? [input.error.message] : []
  const isFirstUse = !snapshot || snapshot.position_count === 0 || positions.length === 0
  const highRiskRatio = snapshot?.high_risk_ratio ?? 0
  const importPending = Boolean(input.importReady && input.importBatchID)

  let statusLabel: string
  let statusTone: ExperienceTone
  let stageLabel = '等待数据'
  let stageDescription = '暂时无法确认本地账户事实，请先查看错误提示或重新加载页面。'

  if (input.error) {
    statusLabel = '组合状态待检查'
    statusTone = 'warning'
  } else if (importPending) {
    statusLabel = '批量导入等待人工确认'
    statusTone = 'warning'
    stageLabel = '导入待确认'
    stageDescription = '导入校验已通过，仍需用户显式确认后才会写入本地账户事实。'
  } else if (isFirstUse) {
    statusLabel = '需要初始化本地账户'
    statusTone = 'warning'
    stageLabel = '首次初始化'
    stageDescription = '先录入现金、总资产和至少一条持仓，今日纪律和主动咨询才有本地账户上下文。'
  } else if (highRiskRatio > 0.3) {
    statusLabel = '高风险仓位需要人工复核'
    statusTone = 'danger'
    stageLabel = '高风险复核'
    stageDescription = '高风险比例处于较高水平，先查看风险预警并复核禁止动作。'
  } else {
    statusLabel = '组合事实可用于纪律评估'
    statusTone = 'success'
    stageLabel = '日常维护'
    stageDescription = '本地账户事实已具备基础上下文，可继续校准、补记线下交易或查看风险。'
  }

  return {
    statusLabel,
    statusTone,
    stageLabel,
    stageDescription,
    summaryMetrics: [
      { label: '总资产', value: formatCurrency(snapshot?.total_assets ?? 0) },
      { label: '现金占比', value: formatPercent(snapshot?.cash_ratio ?? 0) },
      { label: '持仓数量', value: String(snapshot?.position_count ?? positions.length) },
      { label: '高风险比例', value: formatPercent(highRiskRatio), tone: highRiskRatio > 0.3 ? 'danger' : 'success' },
    ],
    nextActions: buildPortfolioActions({ isFirstUse, highRiskRatio, importPending }),
    maintenanceModes: [
      { id: 'calibration', label: '初始化/校准', description: '录入或校准现金、总资产和持仓基础事实。' },
      { id: 'holding', label: '持仓维护', description: '编辑或移除当前本地持仓事实。' },
      { id: 'offline_transaction', label: '线下交易记录', description: '补记用户已经自行完成的线下买入、卖出或减仓。' },
      { id: 'batch_import', label: '批量导入', description: '先校验导入内容，再显式确认写入本地事实。' },
      { id: 'correction', label: '错误修正', description: '记录本地事实更正审计，不自动改写交易。' },
    ],
    safetyNotes: ['这里只记录本地事实、账户事实和审计事实，不连接券商、不自动交易、不代下单。'],
    warnings,
  }
}

function buildPortfolioActions(input: { isFirstUse: boolean; highRiskRatio: number; importPending: boolean }): PortfolioExperienceAction[] {
  if (input.importPending) {
    return [
      { label: '确认批量导入', detail: '校验通过后仍需人工确认，才会写入本地账户事实。', priority: 'blocking' },
      { label: '复核导入明细', detail: '确认行数、标的、数量、价格和买入理由。', priority: 'review' },
    ]
  }
  if (input.isFirstUse) {
    return [
      { label: '录入本地账户与持仓', detail: '先填写现金、总资产、标的、数量、价格和买入理由。', priority: 'blocking' },
      { label: '校验批量导入', detail: '如果已有表格内容，先做本地导入校验。', priority: 'follow_up' },
    ]
  }
  const actions: PortfolioExperienceAction[] = [
    { label: '校准本地账户事实', detail: '更新现金、总资产、现价和持仓说明。', priority: 'review' },
    { label: '补记线下交易', detail: '只记录用户已经自行完成的线下动作。', priority: 'follow_up' },
  ]
  if (input.highRiskRatio > 0.3) {
    actions.unshift({ label: '查看风险预警', detail: '先处理高风险仓位相关 SOP。', priority: 'blocking', href: '/risk-alerts' })
  }
  return actions
}
