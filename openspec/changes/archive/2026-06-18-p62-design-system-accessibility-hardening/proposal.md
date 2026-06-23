# Proposal: P62 设计系统与可访问性验收

## Summary

将 P58-P61 反复出现的 operational cockpit UI 规则固化为轻量设计系统基础层，并补齐键盘路径、可访问语义、390px/768px/1280px reflow 与视觉证据门禁。

## Why

P58-P61 已经完成今日工作台、决策解释、组合/风险/数据质量、治理/运维页面的产品化改造，但这些页面仍主要依赖局部 class 和页面内模式复用。继续推进 P63 全量真实 UI 回归前，需要先把重复的按钮、字段、状态、首屏摘要、详情区、表格和空/错误态沉淀为可测试、可复用、可验收的基础层。

P62 的目标不是重新设计所有页面，而是把已经验证有效的 UI 语言变成稳定组件和浏览器门禁，降低后续维护时的漂移风险，确保主导航、移动菜单、表单、折叠区和关键按钮都能通过键盘和窄屏使用。

## What Changes

- 新增轻量前端 UI primitives：Button、Field、StatusBadge、PageHeader、SummaryCard、DetailSection、ResponsiveTable、EmptyState、ErrorState。
- 统一状态 tone：success、warning、danger、degraded、unknown、readonly、blocked，并要求状态必须有可读文本和可访问名称。
- 将代表性页面逐步接入 primitives，优先覆盖 P58-P61 共享模式和高频关键路径，避免大面积重写。
- 增加组件测试、键盘 smoke、reflow 检查和 390px/768px/1280px 截图证据。
- 新增 P62 UI 验收记录，说明实际启动本地项目、浏览器操作路径、截图、overflow 检查和安全扫描结果。

## In Scope

- `web/src/components/ui/` 下建立轻量设计系统组件和导出入口。
- 复用现有 CSS tokens，补充 focus-visible、按钮尺寸、字段错误态、状态徽标、summary card、详情区和响应式表格样式。
- 代表性接入页面：Dashboard/Workbench、Positions、Risk Alerts、Data Quality、Rules/Audit/Notifications、Local Install/Local Knowledge/Settings 中的共享控件与状态展示。
- Playwright 覆盖主导航、移动菜单、表单控件、折叠/展开区域和关键按钮的键盘路径。
- 浏览器验收覆盖 390px、768px、1280px 三档 viewport，确认无页面级横向溢出；二维表格仅允许在局部容器滚动。
- 安全文案和敏感信息扫描继续覆盖自动交易、外部推送、自动确认、自动修复、完整 key、raw payload、私有路径等禁止内容。

## Out of Scope

- 不新增后端 API、SQLite schema、Eino workflow、LLM 能力或数据源能力。
- 不新增券商接口、登录源、付费源、授权源、Level2、高频源、自动交易、一键交易、代下单或收益承诺。
- 不新增外部推送、短信、邮件、Webhook、自动确认、自动规则应用、自动修复或覆盖真实库。
- 不把 P62 结论表述为最终 release-ready；最终全量真实 UI 回归和发布口径刷新留给 P63。
- 不把所有页面一次性重写为新组件；P62 以基础组件、代表性接入和门禁证明为主。

## Product Design Brief

P62 延续 P57-P61 的产品定位：本地投资纪律工作台，而不是券商交易终端、AI chat demo、营销页或工程调试台。视觉风格保持克制、密集、可扫描，支持重复工作流。设计系统应服务于“今天能不能动、为什么、需要什么人工动作、数据和规则是否可信”，并让高风险、降级、未知、只读、blocked 状态在视觉和文案上都不可误读为普通成功。

交互级别为 full interactivity：使用真实本地后端、真实 Vite 前端、现有 service/API 和浏览器操作验收，不做静态 mock。

## Validation

- `npm --prefix web test`
- `npm --prefix web run build`
- `go test ./...`
- 启动真实本地后端和 Vite 前端，通过浏览器操作代表性页面和关键控件。
- Playwright keyboard smoke 覆盖主导航、移动菜单、表单、折叠区和关键按钮。
- 390px、768px、1280px viewport 截图和 reflow 检查。
- `bash scripts/e2e-smoke.sh`
- `openspec validate p62-design-system-accessibility-hardening --strict`
- `openspec validate --all --strict`
- `git diff --check`
- 敏感信息和 forbidden copy 扫描。
