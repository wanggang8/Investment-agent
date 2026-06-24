import { useState } from 'react'
import { NavLink, Outlet } from 'react-router-dom'

const navGroups = [
  {
    label: '今日',
    items: [
      { to: '/', label: '今日纪律', end: true },
      { to: '/workbench', label: '决策工作台' },
    ],
  },
  {
    label: '决策',
    items: [
      { to: '/consultation', label: '决策咨询' },
      { to: '/decision-loop', label: '决策闭环' },
    ],
  },
  {
    label: '组合',
    items: [
      { to: '/positions', label: '持仓' },
      { to: '/risk-alerts', label: '风险预警' },
    ],
  },
  {
    label: '证据',
    items: [
      { to: '/data-quality', label: '数据质量' },
      { to: '/evidence', label: '情报与证据' },
      { to: '/local-knowledge', label: '本地知识' },
    ],
  },
  {
    label: '治理',
    items: [
      { to: '/rules', label: '规则与纪律' },
      { to: '/audit', label: '复盘与审计' },
      { to: '/review', label: '复盘摘要' },
      { to: '/notifications', label: '通知中心' },
    ],
  },
  {
    label: '系统',
    items: [
      { to: '/settings', label: '设置' },
      { to: '/local-install', label: '本地安装' },
      { to: '/daily-auto-run', label: '每日自动运行' },
      { to: '/daily-discipline/reports', label: '纪律报告' },
    ],
  },
]

// AppLayout 提供全局导航与页面内容区域。
export function AppLayout() {
  const [navOpen, setNavOpen] = useState(false)

  return (
    <div className={navOpen ? 'app-shell nav-open' : 'app-shell'}>
      <header className="app-topbar">
        <div>
          <strong>Investment Agent</strong>
          <span>本地投资纪律工作台</span>
        </div>
        <div className="app-topbar-status" aria-label="运行状态">
          <span>本地运行</span>
          <span>安全边界开启</span>
        </div>
        <button type="button" className="nav-toggle" aria-expanded={navOpen} aria-controls="primary-navigation" onClick={() => setNavOpen((open) => !open)}>
          导航
        </button>
      </header>
      <nav id="primary-navigation" className="app-nav" aria-label="主导航">
        <div className="app-brand">
          <strong>Investment Agent</strong>
          <span>本地投资纪律工作台</span>
        </div>
        <div className="nav-mode-panel" aria-label="本地运行边界">
          <div>
            <strong>本地模式</strong>
            <span>离线优先</span>
          </div>
          <div>
            <strong>只读导航</strong>
            <span>系统只提示和记录，不会自动执行。</span>
          </div>
        </div>
        {navGroups.map((group) => (
          <section key={group.label} className="nav-group" aria-label={group.label}>
            <div className="nav-group-label">{group.label}</div>
            {group.items.map((item) => (
              <NavLink
                key={item.to}
                to={item.to}
                end={item.end}
                className={({ isActive }) => (isActive ? 'active' : undefined)}
                onClick={() => setNavOpen(false)}
              >
                {item.label}
              </NavLink>
            ))}
          </section>
        ))}
      </nav>
      <main className="app-main command-center-shell">
        <Outlet />
      </main>
    </div>
  )
}
