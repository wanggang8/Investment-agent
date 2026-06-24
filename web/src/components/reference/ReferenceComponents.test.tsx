import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { describe, expect, it } from 'vitest'
import {
  EvidenceChecklist,
  LedgerSurface,
  PriorityActionQueue,
  ProgressTracker,
  ReferenceHero,
  SnapshotStrip,
  StatusMetricGrid,
} from './index'

describe('P111 reference components', () => {
  it('renders the report hero with discipline state and prohibited actions', () => {
    render(
      <ReferenceHero
        iconLabel="纪律报告"
        title="今日尚未生成每日纪律报告"
        statusText="证据不足或尚未完成核验"
        description="当前不具备生成每日纪律报告的条件，请先完成数据核查。"
        stateTitle="当前纪律状态"
        stateValue="冻结观察"
        stateDetail="证据不足或尚未完成核验"
        prohibitedTitle="禁止动作"
        prohibitedActions={['暂停交易类建议']}
      />,
    )

    expect(screen.getByRole('region', { name: '纪律报告概览' })).toHaveClass('reference-hero')
    expect(screen.getByText('今日尚未生成每日纪律报告')).toBeInTheDocument()
    expect(screen.getByText('冻结观察')).toBeInTheDocument()
    expect(screen.getByText('暂停交易类建议')).toBeInTheDocument()
  })

  it('renders numbered priority actions with right aligned navigation', () => {
    render(
      <MemoryRouter>
        <PriorityActionQueue
          title="下一步人工动作"
          actions={[
            { label: '维护本地账户与持仓', detail: '更新现金、仓位、成本。', href: '/positions', priority: 'high', meta: '上次维护：2天前' },
            { label: '查看数据质量', detail: '处理数据阻断项。', href: '/data-quality', priority: 'medium', meta: '阻断项：2项' },
          ]}
        />
      </MemoryRouter>,
    )

    const region = screen.getByRole('region', { name: '下一步人工动作' })
    expect(region).toHaveClass('reference-action-queue')
    expect(screen.getByText('1')).toBeInTheDocument()
    expect(screen.getByText('高优先级')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '去维护' })).toHaveAttribute('href', '/positions')
  })

  it('renders metric grid, snapshot strip, progress tracker, checklist, and ledger surface', () => {
    render(
      <MemoryRouter>
        <StatusMetricGrid
          title="状态总览"
          metrics={[
            { label: '数据可信度', value: '62%', status: '需处理', tone: 'warning', details: ['阻断项 2', '告警项 3'] },
            { label: '风险状态', value: '中高', status: '需关注', tone: 'danger', details: ['高风险 1', '中风险 2'] },
          ]}
        />
        <SnapshotStrip
          title="持仓与资金快照"
          items={[
            { label: '总资产（估）', value: '¥1,028,450.28' },
            { label: '仓位水平', value: '78%', status: '中等' },
          ]}
        />
        <ProgressTracker
          title="最近咨询 · 解释预览"
          steps={[
            { label: '输入假设', status: 'done' },
            { label: '规则裁决', status: 'active', detail: '进行中' },
            { label: '等待人工确认', status: 'pending' },
          ]}
        />
        <EvidenceChecklist
          title="证据与规则快照"
          items={[
            { label: '信息核查来源', value: '7/7 完成', status: 'done' },
            { label: '关键规则通过率', value: '83%', status: 'active' },
          ]}
        />
        <LedgerSurface title="证据列表">
          <p>P30SmokeSource</p>
        </LedgerSurface>
      </MemoryRouter>,
    )

    expect(screen.getByRole('region', { name: '状态总览' })).toHaveClass('reference-metric-grid')
    expect(screen.getByRole('region', { name: '持仓与资金快照' })).toHaveClass('reference-snapshot-strip')
    expect(screen.getByRole('region', { name: '最近咨询 · 解释预览' })).toHaveClass('reference-progress-tracker')
    expect(screen.getByRole('region', { name: '证据与规则快照' })).toHaveClass('reference-checklist')
    expect(screen.getByRole('region', { name: '证据列表' })).toHaveClass('reference-ledger-surface')
  })
})
