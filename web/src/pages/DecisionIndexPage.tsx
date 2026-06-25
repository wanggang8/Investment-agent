import { Link } from 'react-router-dom'

export function DecisionIndexPage() {
  return (
    <div>
      <h1 className="page-title">决策详情</h1>
      <section className="daily-hero daily-tone-readonly" aria-label="决策详情入口">
        <div className="daily-hero-main">
          <div className="state-label">决策详情入口</div>
          <h2>先选择一条本地决策记录</h2>
          <p>这里用于打开已经生成的本地决策详情。尚未选中决策时，不展示空白页，也不会自动生成咨询或交易动作。</p>
          <div className="daily-signal-grid quality-signal-grid">
            <article className="ui-summary-card ui-summary-card-readonly">
              <div className="state-label">当前决策</div>
              <strong>未选择</strong>
            </article>
            <article className="ui-summary-card ui-summary-card-readonly">
              <div className="state-label">下一步</div>
              <strong>人工选择</strong>
            </article>
          </div>
        </div>
        <aside className="daily-hero-side" aria-label="决策详情下一步">
          <strong>下一步人工动作</strong>
          <p>从决策闭环、复盘摘要或主动咨询中进入具体决策。</p>
          <div className="link-row">
            <Link to="/decision-loop">查看决策闭环</Link>
            <Link to="/review">查看复盘摘要</Link>
            <Link to="/consultation">发起主动咨询</Link>
          </div>
        </aside>
      </section>
      <article className="cockpit-card">
        <div className="state-label">安全边界</div>
        <h2>只读入口</h2>
        <p>本页只提供本地导航，不创建确认、不记录交易、不改变风险状态，也不应用规则。</p>
      </article>
    </div>
  )
}
