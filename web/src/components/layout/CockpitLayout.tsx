import type { ReactNode } from 'react'

interface CockpitLayoutProps {
  systemPanel: ReactNode
  decisionPanel: ReactNode
  evidencePanel: ReactNode
}

export function CockpitLayout({ systemPanel, decisionPanel, evidencePanel }: CockpitLayoutProps) {
  return (
    <section className="cockpit-layout" aria-label="Agent 决策驾驶舱">
      <aside className="cockpit-column cockpit-column-left">{systemPanel}</aside>
      <main className="cockpit-column cockpit-column-main">{decisionPanel}</main>
      <aside className="cockpit-column cockpit-column-right">{evidencePanel}</aside>
    </section>
  )
}
