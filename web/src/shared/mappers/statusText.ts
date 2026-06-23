export const dashboardStateText: Record<string, string> = {
  first_use: '首次使用',
  normal: '正常',
  insufficient_data: '数据不足',
  frozen_watch: '冻结观察',
  high_risk: '高风险',
}

export const positionStateText: Record<string, string> = {
  normal: '正常',
  sell_only: '仅卖出',
  frozen_watch: '冻结观察',
}

export const marketStateText: Record<string, string> = {
  normal: '正常',
  warning: '预警',
  danger: '危险',
  cold: '偏冷',
  neutral: '中性',
  hot: '偏热',
  extreme: '极端',
  fresh: '新鲜',
  stale: '过期',
  missing: '缺失',
  unknown: '未知状态',
}

export const sourceHealthStatusText: Record<string, string> = {
  healthy: '健康',
  fresh: '新鲜',
  stale: '过期',
  failed: '失败',
  missing: '缺失',
  unknown: '未知状态',
  corrupted: '损坏',
  incompatible: '版本不兼容',
  degraded: '降级',
  no_data: '无数据',
  source_unavailable: '不可用',
  unavailable: '不可用',
  parse_error: '解析失败',
  'parse-error': '解析失败',
  disabled: '未启用',
  stubbed: '使用测试数据',
}

export const sourceCategoryText: Record<string, string> = {
  index_constituents: '指数样本',
  index_weights: '指数权重',
  index_valuation_files: '指数估值文件',
  constituent_financials: '成分财务',
  capital_flow: '资金流向',
  margin_financing: '融资融券',
  sentiment_proxy: '情绪替代指标',
}

export const workflowStatusText: Record<string, string> = {
  completed: '已完成',
  failed: '失败',
  degraded: '降级完成',
  running: '处理中',
}

export const verdictStatusText: Record<string, string> = {
  buy_allowed: '允许买入',
  hold: '持有',
  reduce: '减仓',
  sell_only: '仅卖出',
  frozen_watch: '冻结观察',
  rejected: '拒绝',
  insufficient_data: '数据不足',
  high_risk: '高风险',
}

export const evidenceRoleText: Record<string, string> = {
  formal: '正式证据',
  background: '背景材料',
}

export const verificationStatusText: Record<string, string> = {
  satisfied: '已满足',
  failed: '未通过',
  background_only: '仅背景材料',
}

export const retrievalQualityStatusText: Record<string, string> = {
  hit: '命中',
  miss: '未命中',
  degraded: '降级',
  empty: '无结果',
}

export const retrievalFallbackSourceText: Record<string, string> = {
  veclite: 'VecLite 索引',
  sqlite_summary: 'SQLite 摘要',
  none: '无 fallback',
}

export const sourceConsistencyStatusText: Record<string, string> = {
  checked: '已检查',
  not_checked: '未检查',
  mismatch: '不一致',
}

export const precisionStatusText: Record<string, string> = {
  available: '可展示区间',
  insufficient: '样本不足',
  unavailable: '不可用',
}

export const sellEvaluationStatusText: Record<string, string> = {
  triggered: '已触发人工复核',
  review_required: '需人工复核',
  not_triggered: '未触发',
  not_applicable: '不适用',
}

export const auditStatusText: Record<string, string> = {
  success: '成功',
  failed: '失败',
  degraded: '降级',
}

export const riskTypeText: Record<string, string> = {
  valuation_high: '估值高位',
  buy_thesis_broken: '买入逻辑破坏',
  liquidity_danger: '流动性危险',
  sentiment_extreme: '情绪极端',
  position_limit_breach: '仓位超限',
  insufficient_evidence: '证据不足',
  data_degraded: '数据降级',
}

export const riskSOPStatusText: Record<string, string> = {
  triggered: '已触发',
  active: '处理中',
  observing: '观察中',
  escalated: '已升级',
  resolved: '已解除',
  archived: '已归档',
}

export const riskSeverityText: Record<string, string> = {
  info: '提示',
  warning: '预警',
  critical: '严重',
}

export const auditActionText: Record<string, string> = {
  generate_decision: '生成决策',
  refresh_market_data: '刷新市场数据',
  retrieve_evidence: '检索证据',
  confirm_operation: '记录用户确认',
  mark_error: '标记错误案例',
  create_proposal: '创建规则提案',
  audit_rule_change: '审计规则变更',
  update_rule: '更新规则',
  update_settings: '更新系统设置',
  update_capability: '更新能力圈',
  rebuild_index: '重建索引',
  run_local_task: '运行本地任务',
  risk_alert: '风险预警',
}

export const ruleProposalStatusText: Record<string, string> = {
  draft: '草稿',
  pending_user_confirm: '待用户确认',
  under_gatekeeper_audit: '守门人审计中',
  pending_final_confirm: '待最终确认',
  rejected: '已拒绝',
  applied: '已应用',
}

export const auditActorText: Record<string, string> = {
  system: '系统',
  user: '用户',
  gatekeeper: '守门人',
}

export const systemStatusText: Record<string, string> = {
  ok: '可用',
  ready: '可用',
  available: '可用',
  configured: '已配置',
  degraded: '降级',
  failed: '失败',
  missing: '缺失',
  rebuilding: '重建中',
  unavailable: '不可用',
  disabled: '未启用',
  unknown: '未知状态',
}

export const opsStatusText: Record<string, string> = {
  success: '成功',
  degraded: '降级',
  failed: '失败',
  empty: '暂无数据',
  unknown: '未知状态',
  ok: '可用',
  ready: '可用',
  available: '可用',
  configured: '已配置',
  missing: '缺失',
  unavailable: '不可用',
  disabled: '未启用',
}

export const confidenceText: Record<string, string> = {
  high: '高',
  medium: '中',
  low: '低',
}

export const capabilityStatusText: Record<string, string> = {
  in_scope: '能力圈内',
  out_of_scope: '能力圈外',
  unknown: '未知状态',
}

export const severityText: Record<string, string> = {
  normal: '正常',
  warning: '预警',
  danger: '危险',
  high: '高风险',
  medium: '中风险',
  low: '低风险',
}

export const returnScenarioText: Record<string, string> = {
  upside: '乐观情景',
  base: '基准情景',
  downside: '悲观情景',
}

export function scenarioText(value?: string) {
  if (!value) return '暂无'
  return returnScenarioText[value] ?? '未知情景'
}

export function textOrRaw(map: Record<string, string>, value?: string, fallback = '未知状态') {
  if (!value) return '暂无'
  return map[value] ?? fallback
}
