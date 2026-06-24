import type { DailySignal } from '../../features/dashboard/dailyWorkbenchModel'
import { StatusMetricGrid, type ReferenceMetric } from '../reference'
import { Activity, AlertTriangle, Scale, ShieldCheck } from 'lucide-react'

interface Props {
  signals: DailySignal[]
}

export function WorkbenchSignalGrid({ signals }: Props) {
  const metrics: ReferenceMetric[] = signals.map((signal) => ({
    label: signal.label,
    value: compactValue(signal),
    status: signal.tone === 'success' ? '通过' : signal.tone === 'danger' ? '需关注' : '需处理',
    tone: signal.tone,
    icon: metricIcon(signal.label),
    details: [signal.detail],
  }))

  return (
    <section className="reference-signal-overview" aria-label="今日信号摘要">
      <span className="reference-sr-only">{signals.map((signal) => signal.value).join(' ')}</span>
      <StatusMetricGrid title="状态总览" metrics={metrics} />
    </section>
  )
}

function compactValue(signal: DailySignal) {
  if (signal.label.includes('可信度')) {
    const evidence = signal.value.match(/证据\s*(\d+)\s*条/)
    return evidence ? `${evidence[1]} 条` : signal.value
  }
  if (signal.label.includes('规则')) {
    const rule = signal.value.match(/待确认规则\s*(\d+)/)
    const review = signal.value.match(/复盘\s*(\d+)/)
    if (rule || review) return `${rule?.[1] ?? 0} · ${review?.[1] ?? 0}`
  }
  if (signal.label.includes('风险处置')) {
    const count = signal.value.match(/(\d+)\s*条/)
    return count ? `${count[1]} 条` : signal.value
  }
  return signal.value
}

function metricIcon(label: string) {
  const props = { size: 18, strokeWidth: 2.2 }
  if (label.includes('风险')) return <AlertTriangle {...props} />
  if (label.includes('规则')) return <Scale {...props} />
  if (label.includes('可信')) return <ShieldCheck {...props} />
  return <Activity {...props} />
}
