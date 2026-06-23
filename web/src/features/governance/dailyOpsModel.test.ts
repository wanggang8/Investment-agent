import { describe, expect, it } from 'vitest'
import { buildDailyOpsModel } from './dailyOpsModel'
import type { DailyAutoRunStatus } from '../../types/dailyAutoRun'
import type { DailyDisciplineReport } from '../../types/dailyDisciplineReport'

const report: DailyDisciplineReport = {
  report_id: 'report_1',
  local_date: '2026-06-18',
  scope: 'holdings',
  status: 'insufficient_data',
  summary: '缺少本地持仓，暂不能完成纪律复盘。',
  evidence: { evidence_count: 0, independent_source_count: 0, high_grade_independent_source_count: 0 },
  trend: { success_count: 2, degraded_count: 1, failed_count: 0, insufficient_data_count: 1 },
  missing_action: '补齐本地持仓。',
  missing_categories: ['holdings'],
  safety_note: '不会自动执行交易，需人工复核。',
}

const autoRun: DailyAutoRunStatus = {
  enabled: true,
  status: 'failed',
  run_time: '08:30',
  timezone: 'Asia/Shanghai',
  scope: 'holdings',
  failure_code: 'missing_prerequisites',
  failure_reason: '缺少本地持仓',
  latest_audit_link: '/audit?input_ref=auto',
  safety_note: '仅记录本地每日刷新和纪律评估结果，不会自动执行交易，需人工复核。',
}

describe('buildDailyOpsModel', () => {
  it('combines report and auto-run health into manual next steps', () => {
    const model = buildDailyOpsModel({ reports: [report], autoRun })

    expect(model.overallTone).toBe('warning')
    expect(model.overallLabel).toBe('每日纪律与自动运行需要检查')
    expect(model.metrics).toContainEqual(expect.objectContaining({ label: '最新报告', value: '数据不足' }))
    expect(model.metrics).toContainEqual(expect.objectContaining({ label: '自动运行', value: '失败' }))
    expect(model.nextActions.map((action) => action.label)).toEqual(expect.arrayContaining(['补齐本地持仓', '查看审计记录', '查看报告详情']))
    expect(JSON.stringify(model)).not.toMatch(/自动修复|自动确认|自动规则应用|覆盖真实库|后台交易|收益承诺/)
  })
})

