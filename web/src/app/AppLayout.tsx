import { useState } from 'react'
import {
  Bell,
  BookOpen,
  BriefcaseBusiness,
  CalendarClock,
  ClipboardCheck,
  ClipboardList,
  Database,
  DatabaseZap,
  FileSearch,
  FileText,
  GitBranch,
  Home,
  MessageCircleQuestion,
  PackageCheck,
  RefreshCw,
  Scale,
  Settings,
  ShieldCheck,
  Siren,
  WalletCards,
  Waves,
} from 'lucide-react'
import { NavLink, Outlet, useLocation } from 'react-router-dom'

const navGroups = [
  {
    label: '今日',
    items: [
      { to: '/', label: '今日纪律', end: true, icon: Home },
      { to: '/workbench', label: '决策工作台', icon: BriefcaseBusiness },
    ],
  },
  {
    label: '决策',
    items: [
      { to: '/consultation', label: '决策咨询', icon: MessageCircleQuestion },
      { to: '/decision-loop', label: '决策闭环', icon: GitBranch },
    ],
  },
  {
    label: '组合',
    items: [
      { to: '/positions', label: '持仓', icon: WalletCards },
      { to: '/risk-alerts', label: '风险预警', icon: Siren },
    ],
  },
  {
    label: '证据',
    items: [
      { to: '/data-quality', label: '数据质量', icon: DatabaseZap },
      { to: '/evidence', label: '情报与证据', icon: FileSearch },
      { to: '/local-knowledge', label: '本地知识', icon: BookOpen },
    ],
  },
  {
    label: '治理',
    items: [
      { to: '/rules', label: '规则与纪律', icon: Scale },
      { to: '/audit', label: '复盘与审计', icon: ClipboardList },
      { to: '/review', label: '复盘摘要', icon: ClipboardCheck },
      { to: '/notifications', label: '通知中心', icon: Bell },
    ],
  },
  {
    label: '系统',
    items: [
      { to: '/settings', label: '设置', icon: Settings },
      { to: '/local-install', label: '本地安装', icon: PackageCheck },
      { to: '/daily-auto-run', label: '每日自动运行', icon: CalendarClock },
      { to: '/daily-discipline/reports', label: '纪律报告', icon: FileText },
    ],
  },
]

const routeTitles: Record<string, string> = {
  '/': '今日纪律',
  '/workbench': '用户决策工作台',
  '/consultation': '决策咨询',
  '/decision-loop': '决策闭环',
  '/positions': '持仓与账户',
  '/risk-alerts': '风险预警',
  '/data-quality': '数据质量',
  '/evidence': '情报与证据',
  '/local-knowledge': '本地知识',
  '/rules': '规则与纪律',
  '/audit': '复盘与审计',
  '/review': '复盘摘要',
  '/notifications': '通知中心',
  '/settings': '设置',
  '/local-install': '本地安装',
  '/daily-auto-run': '每日自动运行',
  '/daily-discipline/reports': '纪律报告',
}

// AppLayout 提供全局导航与页面内容区域。
export function AppLayout() {
  const [navOpen, setNavOpen] = useState(false)
  const location = useLocation()
  const pageTitle = routeTitle(location.pathname)
  const now = new Date()
  const today = formatTopbarDate(now)
  const dataCutoff = formatCutoff(now)

  return (
    <div className={navOpen ? 'app-shell nav-open' : 'app-shell'}>
      <header className="app-topbar reference-topbar">
        <div className="reference-topbar-title">
          <strong>{pageTitle}</strong>
          <span>{today}（本地时间）</span>
        </div>
        <div className="app-topbar-status reference-topbar-actions" role="group" aria-label="页面状态工具栏">
          <span><Waves size={15} aria-hidden="true" />本地模式</span>
          <span><Database size={15} aria-hidden="true" />数据截至 {dataCutoff}</span>
          <button type="button"><RefreshCw size={15} aria-hidden="true" />刷新摘要</button>
        </div>
        <button type="button" className="nav-toggle" aria-expanded={navOpen} aria-controls="primary-navigation" onClick={() => setNavOpen((open) => !open)}>
          导航
        </button>
      </header>
      <nav id="primary-navigation" className="app-nav reference-sidebar" aria-label="主导航">
        <div className="app-brand">
          <span className="reference-brand-mark" aria-hidden="true"><ShieldCheck size={25} strokeWidth={2.2} /></span>
          <strong>Investment Agent</strong>
          <span>本地投资纪律工作台</span>
        </div>
        <div className="nav-mode-panel" aria-label="本地运行边界">
          <div>
            <strong>运行模式</strong>
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
                <item.icon size={17} strokeWidth={2.1} aria-hidden="true" />
                {item.label}
              </NavLink>
            ))}
          </section>
        ))}
        <div className="reference-sidebar-footer" aria-label="本地状态">
          <span>v0.1.0</span>
          <strong>本地模式 · 离线优先</strong>
          <small>系统只提示和记录，不会替你执行。</small>
        </div>
      </nav>
      <main className="app-main command-center-shell">
        <Outlet />
        <p className="reference-safety-footer">
          本系统为你的本地投资纪律研究与记录工具，仅提供信息、建议与解释，不构成任何交易指令或投资建议。最终决策与执行由你线下、人工、独立完成。
        </p>
      </main>
    </div>
  )
}

function routeTitle(pathname: string) {
  if (pathname.startsWith('/decisions/')) return '决策详情'
  if (pathname.startsWith('/daily-discipline/reports/')) return '纪律报告详情'
  return routeTitles[pathname] ?? 'Investment Agent'
}

function formatTopbarDate(date: Date) {
  const dateText = new Intl.DateTimeFormat('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' }).format(date)
  const weekday = new Intl.DateTimeFormat('zh-CN', { weekday: 'long' }).format(date)
  return `${dateText.replace(/\//g, '年').replace(/年(\d{2})$/, '月$1日')} ${weekday}`
}

function formatCutoff(date: Date) {
  const parts = new Intl.DateTimeFormat('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  }).formatToParts(date)
  const value = (type: string) => parts.find((part) => part.type === type)?.value ?? ''
  return `${value('month')}-${value('day')} ${value('hour')}:${value('minute')}`
}
