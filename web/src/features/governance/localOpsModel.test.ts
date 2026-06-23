import { describe, expect, it } from 'vitest'
import { buildLocalOpsModel } from './localOpsModel'
import type { SourceHealthItem } from '../../types/market'
import type { SystemStatus, CapabilitySettings } from '../../types/settings'
import type { LocalKnowledgeImportValidationResponse } from '../../types/localKnowledge'

const system: SystemStatus = {
  sqlite_status: 'ok',
  veclite_status: 'degraded',
  deepseek_status: 'failed',
  data_sources: ['stub', 'csindex'],
  log_level: 'info',
  sqlite_path: '/Users/private/investment-agent.db',
}

const capability: CapabilitySettings = {
  asset_types: ['ETF'],
  symbols: ['510300'],
  excluded_symbols: [],
  strategy_scope: ['定投'],
}

const sourceHealth: SourceHealthItem[] = [
  { source_name: 'csindex', source_level: 'A', source_type: 'index_basic', data_category: 'index_valuation_files', freshness: 'parse_error', failure_category: 'parse_error', affected_symbols: ['510300'] },
]

const validation: LocalKnowledgeImportValidationResponse = {
  import_batch_id: 'lk_1',
  summary: { total_count: 2, valid_count: 1, warning_count: 1, blocking_count: 0 },
  rows: [],
  index_plan: { rag_chunk_count: 2, index_status: 'pending' },
  safety_note: '仅写入本地背景材料，后续需由用户在相关页面人工复核。',
}

describe('buildLocalOpsModel', () => {
  it('summarizes configuration, diagnostics, knowledge preview, and redacts sensitive material', () => {
    const model = buildLocalOpsModel({ system, capability, sourceHealth, validation })

    expect(model.overallTone).toBe('warning')
    expect(model.overallLabel).toBe('本地配置与诊断需要检查')
    expect(model.metrics).toContainEqual(expect.objectContaining({ label: 'VecLite', value: '降级' }))
    expect(model.metrics).toContainEqual(expect.objectContaining({ label: 'DeepSeek', value: '失败' }))
    expect(model.metrics).toContainEqual(expect.objectContaining({ label: '知识预览', value: '1 可写入 / 1 需关注' }))
    expect(model.nextActions.map((action) => action.label)).toEqual(expect.arrayContaining(['查看数据源健康', '复验本地安装', '复核知识导入']))

    const text = JSON.stringify(model)
    expect(text).not.toMatch(/\/Users\/private|sqlite_path|sk-|SELECT \* FROM|raw vendor|完整 prompt/)
    expect(text).not.toMatch(/自动修复|自动确认|自动规则应用|覆盖真实库|外部推送|自动交易/)
  })
})

