import type { PageErrorState } from '../../shared/utils'
import type { EvidenceItem, SourceVerification } from '../../types/evidence'
import type { MarketSnapshot, SourceHealthItem } from '../../types/market'
import type { ReviewSummary } from '../../types/review'
import type { SystemStatus } from '../../types/settings'
import type { DataQualityGateResolutionCheck, DataSourceQualityRegression } from '../../types/dataSourceQuality'

export type DataQualityTone = 'success' | 'warning' | 'danger' | 'unknown'

export interface DataQualitySignal {
  label: string
  value: string
  detail: string
  tone: DataQualityTone
}

export interface DataQualityAction {
  label: string
  detail: string
  href: string
}

export interface DataQualityExperienceModel {
  overallLabel: string
  overallTone: DataQualityTone
  qualitySignals: DataQualitySignal[]
  nextActions: DataQualityAction[]
  safetyNotes: string[]
  warnings: string[]
}

export interface DataQualityExperienceInput {
  system?: SystemStatus
  market?: MarketSnapshot
  sourceHealth?: SourceHealthItem[]
  evidenceItems?: EvidenceItem[]
  evidenceTotal?: number
  verification?: SourceVerification
  review?: ReviewSummary
  currentRegression?: DataSourceQualityRegression
  gateResolution?: DataQualityGateResolutionCheck
  errors?: PageErrorState[]
}

const degradedValues = new Set(['degraded', 'stale', 'missing', 'parse_error', 'source_unavailable', 'failed', 'unknown', 'disabled', 'no_data', 'unavailable', 'empty', 'quality_failed', 'insufficient', 'insufficient_data'])

export function buildDataQualityExperienceModel(input: DataQualityExperienceInput): DataQualityExperienceModel {
  const sourceHealth = input.sourceHealth ?? []
  const degradedSources = sourceHealth.filter((item) => isDegraded(item.freshness) || isDegraded(item.failure_category))
  const evidenceTotal = input.evidenceTotal ?? input.evidenceItems?.length ?? 0
  const degradedWorkflows = input.review?.degraded_workflows ?? []
  const missingEvidenceCount = input.review?.missing_evidence_count ?? 0
  const degradedCount = input.review?.degraded_count ?? 0
  const hasDiagnostics = Boolean(input.review?.ops_status?.explanation)
  const opsStatus = input.review?.ops_status
  const dataSourceNeedsInspection = isDegraded(opsStatus?.data_source_status)
  const indexNeedsInspection = isDegraded(opsStatus?.index_status)
  const reviewNeedsInspection = isDegraded(opsStatus?.review_status)
  const opsNeedsInspection = dataSourceNeedsInspection || indexNeedsInspection || reviewNeedsInspection
  const policySignal = buildPolicySignal(input.currentRegression, input.gateResolution)

  const sourceSignal: DataQualitySignal = degradedSources.length || isDegraded(input.market?.data_status) || dataSourceNeedsInspection
    ? { label: '数据源健康', value: `${Math.max(degradedSources.length, isDegraded(input.market?.data_status) || dataSourceNeedsInspection ? 1 : 0)} 项需检查`, detail: '存在过期、缺失或解析失败的数据源事实。', tone: 'warning' }
    : { label: '数据源健康', value: '当前可用', detail: '未发现数据源健康降级记录。', tone: sourceHealth.length ? 'success' : 'unknown' }

  const evidenceSignal: DataQualitySignal = evidenceTotal > 0 && input.verification?.verification_status === 'satisfied' && !isDegraded(input.system?.veclite_status) && !indexNeedsInspection
    ? { label: '证据与 RAG', value: `${evidenceTotal} 条证据`, detail: `${input.verification.independent_source_count} 个独立信源，VecLite ${input.system?.veclite_status ?? '暂无'}` , tone: 'success' }
    : { label: '证据与 RAG', value: evidenceTotal > 0 ? `${evidenceTotal} 条证据待检查` : '证据不足', detail: '证据、核验或索引状态不足，不能当作完整可信证据。', tone: 'warning' }

  const llmSignal: DataQualitySignal = isDegraded(input.system?.deepseek_status) || degradedCount > 0 || missingEvidenceCount > 0 || opsNeedsInspection
    ? { label: 'LLM 分析', value: `复盘降级 ${degradedCount} 条`, detail: `缺证据 ${missingEvidenceCount} 条；LLM 状态 ${input.system?.deepseek_status ?? '暂无'}` , tone: 'warning' }
    : { label: 'LLM 分析', value: '当前可用', detail: `LLM 状态 ${input.system?.deepseek_status ?? '暂无'}` , tone: input.system?.deepseek_status ? 'success' : 'unknown' }

  const impactSignal: DataQualitySignal = degradedWorkflows.length
    ? { label: '影响范围', value: `${degradedWorkflows.length} 个受影响工作流`, detail: '存在降级决策或复盘线索，需要查看对应本地页面。', tone: 'warning' }
    : { label: '影响范围', value: '暂无受影响工作流', detail: '当前复盘未返回受影响工作流。', tone: 'success' }

  const qualitySignals = policySignal ? [policySignal, sourceSignal, evidenceSignal, llmSignal, impactSignal] : [sourceSignal, evidenceSignal, llmSignal, impactSignal]
  const overallTone = resolveOverallTone(qualitySignals, input.errors)
  const policyReasons = sanitizePolicyReasons(input.currentRegression)

  return {
    overallLabel: resolveOverallLabel(overallTone, input.currentRegression, input.gateResolution),
    overallTone,
    qualitySignals,
    nextActions: buildDataQualityActions(degradedWorkflows[0]?.decision_id, Boolean(input.currentRegression)),
    safetyNotes: ['数据质量页只做只读检查和本地导航，不发起后台变更、规则确认、规则生效或资金动作。'],
    warnings: [
      ...(input.errors ?? []).map((error) => error.message),
      ...(hasDiagnostics ? ['存在已脱敏诊断摘要。'] : []),
      ...policyReasons,
    ],
  }
}

function buildPolicySignal(currentRegression?: DataSourceQualityRegression, gateResolution?: DataQualityGateResolutionCheck): DataQualitySignal | undefined {
  const policy = currentRegression?.policy
  if (!policy) return undefined
  switch (gateResolution?.release_claim_state) {
    case 'pass':
      return { label: '当前数据策略', value: '通过', detail: '当前数据源健康满足发布门禁。', tone: 'success' }
    case 'resolved_with_waiver':
      return { label: '当前数据策略', value: '已记录豁免', detail: '存在 waiver 记录；不得描述为 clean pass。', tone: 'warning' }
    case 'resolved_with_scope_exclusion':
      return { label: '当前数据策略', value: '已排除 clean claim', detail: 'current local data health 已排除在 clean claim 外。', tone: 'warning' }
  }
  switch (policy.release_gate) {
    case 'pass':
      return { label: '当前数据策略', value: '通过', detail: '当前数据源健康满足发布门禁。', tone: 'success' }
    case 'waiver_required':
      return { label: '当前数据策略', value: '需豁免记录', detail: `需记录 ${policy.waiver_count} 项 waiver；不得描述为 clean pass。`, tone: 'warning' }
    case 'block':
      return { label: '当前数据策略', value: '阻断', detail: `${policy.blocking_count} 项阻断；不得声明当前数据质量 clean。`, tone: 'danger' }
    default:
      return { label: '当前数据策略', value: '待检查', detail: 'policy verdict 未识别，需要复核。', tone: 'warning' }
  }
}

function resolveOverallLabel(overallTone: DataQualityTone, currentRegression?: DataSourceQualityRegression, gateResolution?: DataQualityGateResolutionCheck) {
  if (gateResolution?.release_claim_state === 'resolved_with_scope_exclusion') return '当前数据声明已限定范围'
  if (gateResolution?.release_claim_state === 'resolved_with_waiver') return '当前数据质量已记录豁免'
  if (currentRegression?.policy?.release_gate === 'block') return '当前数据质量阻断发布声明'
  return overallTone === 'success' ? '数据质量可用于日常判断' : '数据质量需要检查'
}

function sanitizePolicyReasons(currentRegression?: DataSourceQualityRegression) {
  const policy = currentRegression?.policy
  if (!policy) return []
  return [...(policy.blocking_reasons ?? []), ...(policy.waiver_reasons ?? [])]
    .map((item) => sanitizePolicyText(item))
    .filter(Boolean)
    .slice(0, 4)
}

function isDegraded(value?: string) {
  return value ? degradedValues.has(value) : false
}

function resolveOverallTone(signals: DataQualitySignal[], errors?: PageErrorState[]): DataQualityTone {
  if (errors?.length) return 'warning'
  if (signals.some((signal) => signal.tone === 'danger')) return 'danger'
  if (signals.some((signal) => signal.tone === 'warning' || signal.tone === 'unknown')) return 'warning'
  return 'success'
}

function buildDataQualityActions(impactedDecisionID?: string, hasPolicy?: boolean): DataQualityAction[] {
  const actions: DataQualityAction[] = [
    { label: '查看数据源设置', detail: '检查本地数据源启用状态和配置。', href: '/settings' },
    { label: '查看证据', detail: '核对证据来源等级、核验状态和检索索引。', href: '/evidence' },
    { label: '查看质量复盘', detail: '回到本地复盘摘要确认缺证据和降级工作流。', href: '/review' },
    { label: '查看风险预警', detail: '确认数据降级是否触发本地风险 SOP。', href: '/risk-alerts' },
  ]
  if (impactedDecisionID) {
    actions.splice(1, 0, { label: '查看受影响决策', detail: '从降级工作流回到具体决策解释。', href: `/decisions/${encodeURIComponent(impactedDecisionID)}` })
  }
  if (hasPolicy) {
    actions.unshift({ label: '查看当前数据策略', detail: '复核策略结论、发布门禁和人工处置原因。', href: '/data-quality' })
  }
  return actions
}

function sanitizePolicyText(value: string) {
  return value
    .replace(/sk-[A-Za-z0-9_-]+/g, '[REDACTED_KEY]')
    .replace(/\/Users\/[A-Za-z0-9._-]+[^\s，。；;]*/g, '[REDACTED_PATH]')
    .replace(/\bSELECT\s+\*\s+FROM\b/gi, 'SELECT [REDACTED]')
    .replace(/prompt\s*:/gi, 'prompt [REDACTED]')
    .replace(/raw\s+(vendor|provider|HTTP)[^\s，。；;]*/gi, 'raw [REDACTED]')
}
