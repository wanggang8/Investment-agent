import type { DailySignal } from '../../features/dashboard/dailyWorkbenchModel'
import { SummaryCard, type UITone } from '../ui'

interface Props {
  signals: DailySignal[]
}

export function WorkbenchSignalGrid({ signals }: Props) {
  return (
    <section className="daily-signal-grid" aria-label="今日信号摘要">
      {signals.map((signal) => (
        <SummaryCard
          key={signal.label}
          title={signal.label}
          value={signal.value}
          detail={signal.detail}
          tone={signal.tone as UITone}
          action={signal.href ? { label: '查看', href: signal.href } : undefined}
        />
      ))}
    </section>
  )
}
