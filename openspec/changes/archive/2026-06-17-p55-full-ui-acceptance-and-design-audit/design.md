# P55 设计：前端全功能真实验收与设计审查

## 设计目标

P55 是真实 UI 验收和设计审查阶段，不是功能开发阶段。它通过真实启动项目、浏览器操作和截图证据确认前端主要功能是否可用，并用 Product Design audit 框架评估 UI 是否需要优化。

## 验收运行方式

使用临时运行目录：

- `tmp/ui-acceptance/p55-2026-06-17/`

建议启动步骤：

1. 生成临时 SQLite 和配置。
2. 使用 `go run ./cmd/smoke-seed` 写入 smoke 数据。
3. 启动 `go run ./cmd/server --config <tmp config>`。
4. 启动 `npm --prefix web run dev -- --host 127.0.0.1 --port <port>`，设置 `VITE_API_PROXY_TARGET` 指向后端。
5. 用 in-app Browser 访问前端地址。

如果已有脚本 `scripts/e2e-smoke.sh` 提供可复用启动方式，可参考其配置和 seed，但 P55 必须进行人工式浏览器操作和截图，不只运行脚本。

## 功能覆盖矩阵

验收至少覆盖：

- 全局导航和布局。
- Dashboard 今日状态。
- Workbench 用户决策工作台。
- Decision Loop 决策闭环解释。
- Data Quality 数据质量面板。
- Portfolio 持仓/账户录入、线下交易、导入/确认/修正可见控件。
- Consultation / Decision detail 主动咨询和决策详情。
- Evidence 证据页。
- Rules 规则提案/效果追踪。
- Audit 审计时间线。
- Notifications 通知。
- Risk Alerts 风险预警列表与详情。
- Daily Auto Run 每日自动运行状态。
- Daily Discipline 报告列表和详情。
- Review 复盘摘要。
- Local Install 本地安装诊断。
- Local Knowledge 本地知识导入治理。
- Settings 设置页。

每个页面记录：

- URL。
- 截图路径。
- 可达性。
- 关键 UI 元素是否可见。
- 关键操作或状态切换；包含多个功能区的页面必须按主要功能区分别记录，不只做页面可达性检查。
- 是否出现 console error、空白页、重叠、文本溢出、明显错误文案。
- 安全边界：不得出现自动交易、一键交易、代下单、外部推送、收益承诺等入口。

## Product Design audit

使用 Product Design `audit` 框架。审查维度：

- 信息架构和导航。
- 任务入口和发现性。
- 页面层级与扫描效率。
- 表格/卡片密度。
- 空态/错误态/降级态说明。
- 可访问性风险：对比度、焦点、标签、目标尺寸、键盘路径、响应式 reflow。
- 交易安全边界的可见性。

产物：

- `docs/release/ui-acceptance-2026-06-17.md`
- `docs/release/ui-design-audit-2026-06-17.md`
- `docs/release/ui-audit-assets/2026-06-17-p55/*.png`

## 结果处理

- 如果页面不可达、关键控件无法操作、出现运行时错误或出现越界能力入口，记录为 `blocked`。
- 如果页面可用但存在明显体验或设计问题，记录为 `needs_optimization`。
- 如果仅有轻微视觉/文案/密度问题，记录为 `minor`.
- P55 不直接修复问题；后续若需要优化，创建独立 UI improvement change。

## 安全边界

P55 不提交临时数据库、完整 key、raw API 响应、完整 prompt、私有路径或未脱敏日志。截图不得包含完整 key。
