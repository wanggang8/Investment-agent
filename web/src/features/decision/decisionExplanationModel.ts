import type { DecisionDetailResponse } from '../../types/decision'
import type { EvidenceItem } from '../../types/evidence'

export interface ExplanationLink {
  label: string
  href: string
}

export interface DecisionExplanationModel {
  storyTitle: string
  decisionContext: string[]
  keyReasons: string[]
  prohibitedActions: string[]
  optionalActions: string[]
  trustSummary: string[]
  explanationLinks: ExplanationLink[]
  safetyNotes: string[]
  missingDataWarnings: string[]
}

export function buildDecisionExplanationModel(decision: DecisionDetailResponse): DecisionExplanationModel {
  const evidenceItems = safeArray<EvidenceItem>(decision.evidence_chain)
  const analystReports = safeArray<DecisionDetailResponse['analyst_reports'][number]>(decision.analyst_reports)
  const triggeredRules = safeArray<DecisionDetailResponse['triggered_rules'][number]>(decision.triggered_rules)
  const formalEvidence = evidenceItems.filter((item) => (item.evidence_role ?? 'formal') !== 'background')
  const backgroundEvidence = evidenceItems.filter((item) => (item.evidence_role ?? 'formal') === 'background')
  const keyReasons = [
    ...analystReports.flatMap((report) => {
      const reasons = safeArray<string>(report.key_reasons)
      return reasons.length ? reasons.map(compactReason) : report.conclusion ? [`${report.agent_name}：${compactReason(report.conclusion)}`] : []
    }),
    ...triggeredRules.map((rule) => `${rule.rule_name}：${rule.description}`).filter(Boolean),
    ...formalEvidence.map((item) => `${item.source_name}：${item.summary}`).filter(Boolean),
  ].slice(0, 6)
  const missingDataWarnings: string[] = []

  if (formalEvidence.length === 0) {
    missingDataWarnings.push('缺少可展示的正式证据，最终结论需要人工复核。')
  }
  if (!analystReports.length) {
    missingDataWarnings.push('暂无可展示的 LLM 分析材料，页面只展示规则与本地事实。')
  }
  if (decision.retrieval_quality?.status === 'degraded' || decision.retrieval_quality?.status === 'empty') {
    missingDataWarnings.push('检索质量降级或无结果，证据需要重新核验后再复核。')
  }

  return {
    storyTitle: decision.final_verdict?.display_text || '暂无最终裁决',
    decisionContext: [
      `决策 ${decision.decision_id}`,
      decision.symbol ? `标的 ${decision.symbol}` : '标的 暂无',
      decision.question ? `问题：${decision.question}` : '问题：暂无',
      decision.generated_at ? `生成时间：${decision.generated_at}` : '生成时间：暂无',
    ],
    keyReasons: keyReasons.length ? keyReasons : ['暂无关键理由；请先补齐证据或重新执行咨询。'],
    prohibitedActions: safeTextList(decision.final_verdict?.prohibited_actions, '暂无；缺失字段不代表允许交易或自动执行'),
    optionalActions: safeTextList(decision.final_verdict?.optional_actions, '暂无；仅可人工复核'),
    trustSummary: [
      `正式证据 ${formalEvidence.length} 条，背景材料 ${backgroundEvidence.length} 条`,
      `最高信源等级 ${highestSourceLevel(evidenceItems)}`,
      `检索状态 ${decision.retrieval_quality?.status || '暂无'}，召回 ${decision.retrieval_quality?.top_k ?? 0} 条`,
      `LLM 材料 ${analystReports.length} 份；解析/质量通过 ${analystReports.filter(isAnalystReportUsable).length} 份`,
    ],
    explanationLinks: [
      { label: '查看证据', href: '/evidence' },
      { label: '查看决策闭环', href: '/decision-loop' },
      { label: '查看审计', href: '/audit' },
    ],
    safetyNotes: [
      '只展示本地分析和人工复核路径，不会自动交易、自动确认或自动应用规则。',
      'LLM 材料只作为分析输入，最终裁决仍以规则链与守门条件为准。',
      '缺失、降级或 nullable 字段不会被解释为允许交易、自动确认或自动应用规则。',
    ],
    missingDataWarnings,
  }
}

function safeArray<T>(value: unknown): T[] {
  return Array.isArray(value) ? value.filter((item): item is T => Boolean(item)) : []
}

function safeTextList(value: unknown, emptyText: string) {
  const items = safeArray<string>(value).filter((item) => item.trim().length > 0)
  return items.length ? items : [emptyText]
}

function highestSourceLevel(items: EvidenceItem[]) {
  const order = ['S', 'A', 'B', 'C']
  const levels = items.map((item) => item.source_level).filter(Boolean)
  return order.find((level) => levels.includes(level)) ?? '暂无'
}

function isAnalystReportUsable(report: DecisionDetailResponse['analyst_reports'][number]) {
  return report.quality_status === 'passed' || report.parse_status === 'parsed'
}

function compactReason(value: string) {
  const normalized = value.replace(/\s+/g, ' ').trim()
  if (normalized.length <= 100) {
    return normalized
  }
  return `${normalized.slice(0, 97)}...`
}
