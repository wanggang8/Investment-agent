import type { SourceHealthItem } from '../../types/market'
import type { LocalKnowledgeImportValidationResponse } from '../../types/localKnowledge'
import type { CapabilitySettings, SystemStatus } from '../../types/settings'
import { sourceHealthStatusText, systemStatusText, textOrRaw } from '../../shared/mappers'
import type { OpsAction, OpsMetric, OpsTone } from './modelTypes'

export interface LocalOpsModel {
  overallLabel: string
  overallTone: OpsTone
  metrics: OpsMetric[]
  nextActions: OpsAction[]
  safetyNotes: string[]
}

export function buildLocalOpsModel(input: {
  system?: SystemStatus
  capability?: CapabilitySettings
  sourceHealth?: SourceHealthItem[]
  validation?: LocalKnowledgeImportValidationResponse
}): LocalOpsModel {
  const sourceHealth = input.sourceHealth ?? []
  const degradedSources = sourceHealth.filter((item) => ['stale', 'failed', 'missing', 'unknown', 'degraded', 'no_data', 'source_unavailable', 'unavailable', 'parse_error', 'disabled'].includes(item.freshness) || (item.failure_category && item.failure_category !== 'none'))
  const systemNeedsCheck = ['degraded', 'failed', 'missing', 'unavailable', 'disabled', 'unknown'].some((status) => [input.system?.sqlite_status, input.system?.veclite_status, input.system?.deepseek_status].includes(status))
  const warningRows = input.validation?.summary.warning_count ?? 0
  const blockingRows = input.validation?.summary.blocking_count ?? 0
  const needsCheck = systemNeedsCheck || degradedSources.length > 0 || warningRows > 0 || blockingRows > 0

  return {
    overallLabel: needsCheck ? '本地配置与诊断需要检查' : '本地配置与诊断可用',
    overallTone: blockingRows ? 'danger' : needsCheck ? 'warning' : 'success',
    metrics: buildMetrics(input, degradedSources.length),
    nextActions: buildLocalActions(degradedSources.length, input.validation),
    safetyNotes: ['系统页只展示本地配置、诊断摘要、脱敏预览和人工复验路径；敏感配置、私有路径和底层诊断材料仅在二级详情中脱敏呈现。'],
  }
}

function buildMetrics(input: {
  system?: SystemStatus
  capability?: CapabilitySettings
  validation?: LocalKnowledgeImportValidationResponse
}, degradedSourceCount: number): OpsMetric[] {
  return [
    { label: 'SQLite', value: textOrRaw(systemStatusText, input.system?.sqlite_status), tone: input.system?.sqlite_status === 'ok' ? 'success' : 'warning' },
    { label: 'VecLite', value: textOrRaw(systemStatusText, input.system?.veclite_status), tone: input.system?.veclite_status === 'ok' ? 'success' : 'warning' },
    { label: 'DeepSeek', value: textOrRaw(systemStatusText, input.system?.deepseek_status), tone: input.system?.deepseek_status === 'configured' || input.system?.deepseek_status === 'ok' ? 'success' : 'warning' },
    { label: '能力圈标的', value: input.capability?.symbols?.join('、') || '暂无' },
    { label: '降级数据源', value: String(degradedSourceCount), tone: degradedSourceCount ? 'warning' : 'success' },
    { label: '知识预览', value: input.validation ? `${input.validation.summary.valid_count} 可写入 / ${input.validation.summary.warning_count} 需关注` : '暂无' },
    { label: '索引计划', value: input.validation ? `${input.validation.index_plan.rag_chunk_count} 片段 / ${textOrRaw(sourceHealthStatusText, input.validation.index_plan.index_status)}` : '暂无' },
  ]
}

function buildLocalActions(degradedSourceCount: number, validation?: LocalKnowledgeImportValidationResponse): OpsAction[] {
  const actions: OpsAction[] = []
  actions.push({
    label: '查看数据源健康',
    detail: degradedSourceCount ? '检查降级、失败或解析异常的数据源。' : '复核公开只读数据源的最新健康状态。',
    href: '/data-quality',
  })
  actions.push({ label: '复验本地安装', detail: '按本地诊断命令重新检查运行前提。', href: '/local-install' })
  if (validation) actions.push({ label: '复核知识导入', detail: '确认脱敏预览、阻断项和索引计划。', href: '/local-knowledge' })
  actions.push({ label: '查看设置', detail: '检查能力圈、系统状态和本地市场刷新入口。', href: '/settings' })
  return actions
}

export function localOpsMetricTitle(label: string) {
  const titles: Record<string, string> = {
    SQLite: '本地数据库',
    VecLite: '检索索引',
    DeepSeek: '分析模型',
  }
  return titles[label] || label
}
