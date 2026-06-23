# Design: P62 设计系统与可访问性验收

## Current State

P58-P61 已经形成了一套可用的 operational cockpit 体验：

- 首页和 workbench 先展示今日状态、数据可信度、禁止动作和下一步人工动作。
- 决策解释、组合维护、风险处置、数据质量、治理和运维页面已经开始复用 `daily-hero`、`daily-signal-grid`、`cockpit-card`、`table-wrap` 等样式。
- 390px reflow、真实浏览器验收和 forbidden copy scan 已经在多个阶段执行。

当前缺口在于：这些规则还分散在页面和 CSS class 中，按钮、状态、字段、详情区和表格的可访问语义不够集中；键盘路径和 768px 中间 viewport 还没有形成稳定门禁；后续 P63 全量回归前需要先降低 UI 漂移风险。

## Options Considered

### Option A: 全量页面重写为设计系统

优点是统一彻底，但会把 P62 扩成大范围 UI 迁移，容易引入产品行为回归，也会抢 P63 全量回归的范围。

### Option B: 只做可访问性检查，不抽组件

优点是变更小，但无法解决 P58-P61 模式分散的问题，后续页面继续复制局部 class 时仍会漂移。

### Option C: 轻量 primitives + 代表性接入 + 浏览器门禁

这是 P62 采用方案。先把最稳定的 UI 原语沉淀为组件和测试，在代表性页面接入证明可用，再通过键盘、reflow、截图和扫描门禁约束后续工作。

## Component Architecture

新增 `web/src/components/ui/`，以组合式、低侵入方式提供：

- `Button`：primary、secondary、ghost、danger、link-like 语义，支持 `type`、`disabled`、`aria-label`、loading/working copy。
- `Field`：label、hint、error、required、input/select/textarea 插槽，确保控件和说明有稳定关联。
- `StatusBadge`：success、warning、danger、degraded、unknown、readonly、blocked tone；必须包含文本，不只依赖颜色。
- `PageHeader`：页面标题、说明、状态 badge、主要指标和下一步人工动作区域。
- `SummaryCard`：指标/状态摘要卡，限制标题层级和固定布局，避免 hover 或动态内容造成 layout shift。
- `DetailSection`：可选折叠详情区，统一 `aria-expanded`、键盘操作和局部详情说明。
- `ResponsiveTable`：统一 caption/aria-label、局部横向滚动、移动端 stack/reflow 规则。
- `EmptyState`、`ErrorState`：安全空态和错误态，不渲染 raw stack、完整 key、私有路径或 raw vendor payload。

组件只处理展示和可访问语义，不读取 API、SQLite、VecLite、localStorage、sessionStorage、本地文件或临时配置。页面仍负责数据获取、service 调用和业务状态组合。

## Styling

P62 不引入新视觉主题。样式复用当前 `web/src/styles/global.css` 中的颜色、间距、shadow 和 cockpit tokens，并补充少量基础类：

- `.ui-button`、`.ui-field`、`.ui-status-badge`
- `.ui-page-header`、`.ui-summary-card`、`.ui-detail-section`
- `.ui-responsive-table`、`.ui-empty-state`、`.ui-error-state`
- 全局 `:focus-visible` 规则，确保键盘焦点在深浅背景上都可见

设计要求：

- 状态不得只靠颜色区分；必须有文字、图标或可访问 label。
- 390px、768px 和 1280px 都应保持信息层级稳定。
- 页面级横向滚动禁止；表格、JSON、日志类二维内容只能在明确局部容器滚动。
- 不使用营销式 hero、装饰性渐变球、过大的卡片套卡片或单一色相主题。

## Adoption Strategy

P62 按风险和复用价值接入，不追求一次性替换所有旧 markup：

1. 先为 primitives 写组件测试和样式。
2. 在共享状态/空态/错误态最多的页面接入：Data Quality、Risk Alerts、Positions、Rules/Audit/Notifications、Local Install/Local Knowledge/Settings。
3. 对 Dashboard/Workbench 保持已验证的信息架构，只抽取可复用的标题、状态和 summary card 片段。
4. 旧页面局部 class 可暂时保留，但新增或改动区域应优先使用 primitives。

## Accessibility And Keyboard Gates

P62 必须覆盖以下浏览器级行为：

- 主导航和移动菜单可通过键盘聚焦、打开、关闭并进入目标页面。
- 表单 label 与控件关联；错误提示能被文本识别。
- 折叠区使用 button 语义，并维护 `aria-expanded`。
- 关键按钮 disabled/working 状态可见，不仅依赖颜色。
- 表格有 caption 或可访问名称；移动端 reflow 不丢失列含义。
- 主要页面有可查询 landmark 和稳定标题。

## Browser Validation

真实 UI 验收必须启动本地 Go 后端和 Vite 前端，使用浏览器操作代表性页面，而不是只截图静态 HTML。验收记录至少包含：

- 390px、768px、1280px 截图或等价浏览器证据。
- `document.documentElement.scrollWidth` 和 `document.body.scrollWidth` 与 viewport 的比较结果。
- 键盘路径记录：主导航、移动菜单、表单、折叠区、关键按钮。
- forbidden copy scan 和敏感信息扫描结果。

## Risks

- 组件抽象过度导致大范围回归。控制方式：只抽稳定展示模式，页面业务逻辑保持在原 feature 中。
- 旧测试依赖页面文案。控制方式：先补组件测试，再最小同步页面测试和 E2E 断言。
- 中间 viewport 漏检。控制方式：P62 明确增加 768px 证据，不只覆盖 390px 和桌面。
- 可访问性只停留在属性层面。控制方式：通过浏览器键盘 smoke 验证真实可操作路径。
