import { Link } from 'react-router-dom'

export function ApiDiagnosticsPage() {
  return (
    <div>
      <h1 className="page-title">接口诊断</h1>
      <section className="daily-hero daily-tone-readonly" aria-label="接口诊断入口">
        <div className="daily-hero-main">
          <div className="state-label">本地诊断入口</div>
          <h2>查看本地接口与页面状态</h2>
          <p>这里提供安全的本地诊断导航。页面只展示脱敏状态和人工复验路径，不显示底层响应、私有路径或完整命令输出。</p>
          <div className="daily-signal-grid quality-signal-grid">
            <article className="ui-summary-card ui-summary-card-readonly">
              <div className="state-label">接口状态</div>
              <strong>请在设置页复验</strong>
            </article>
            <article className="ui-summary-card ui-summary-card-readonly">
              <div className="state-label">诊断材料</div>
              <strong>脱敏展示</strong>
            </article>
          </div>
        </div>
        <aside className="daily-hero-side" aria-label="接口诊断下一步">
          <strong>下一步本地复验</strong>
          <div className="link-row">
            <Link to="/settings">查看设置</Link>
            <Link to="/local-install">复验本地安装</Link>
            <Link to="/data-quality">查看数据质量</Link>
          </div>
        </aside>
      </section>
      <article className="cockpit-card">
        <div className="state-label">安全边界</div>
        <h2>不会触发后台动作</h2>
        <p>本页不刷新市场数据、不修复索引、不生成决策、不确认规则，也不连接任何交易通道。</p>
      </article>
    </div>
  )
}
