# Tasks: P110 视觉系统重设计

## 1. Scope And Governance

- [x] 1.1 阅读 `docs/GOVERNANCE.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 1.2 阅读 `docs/frontend-contract.md`、`docs/ui-design.md`、`docs/product-experience-polish-roadmap.md`。
- [x] 1.3 检查 P102/P104 最新验收截图与证据，确认当前 UI 可用但审美可升级。
- [x] 1.4 创建 `p110-visual-system-redesign` OpenSpec change。
- [x] 1.5 更新进度文件，将 P110 标记为当前活跃 change。
- [x] 1.6 运行 `openspec validate p110-visual-system-redesign --strict`。

## 2. Visual Direction Selection

- [x] 2.1 生成三版独立视觉方向图：Research Terminal、Calm Command Center、Ledger Pro。
- [x] 2.2 向用户展示三版方向并等待选择。
- [x] 2.3 用户选择按推荐方向执行，采用 Calm Command Center 为主、Ledger Pro 用于证据/闭环页；无需再生成修订方向图。
- [x] 2.4 将最终选定方向记录到 P110 design/acceptance 材料中。

## 3. Tests First

- [x] 3.1 为视觉 token/status 语义保留添加或更新前端测试。
- [x] 3.2 为 Dashboard/Workbench 核心层级和安全 CTA 添加或更新测试。
- [x] 3.3 为关键页面 reflow、长文本和局部滚动添加或更新测试。
- [x] 3.4 为 forbidden copy/redaction 边界添加或更新扫描脚本或测试。

## 4. Implementation

- [x] 4.1 更新 `web/src/styles/global.css` 的视觉 token、排版、间距、surface、border、状态 tone 和响应式规则。
- [x] 4.2 优化 `AppLayout` 导航、品牌区、active 状态和移动 topbar。
- [x] 4.3 优化 Dashboard/Workbench 首屏 hero、人工动作队列、信号摘要和详细驾驶舱。
- [x] 4.4 优化 Consultation、Decision Detail、Evidence、Decision Loop 的解释链路阅读层级。
- [x] 4.5 优化 Positions、Risk Alerts、Data Quality 的事实维护、处置队列和质量面板视觉。
- [x] 4.6 保持所有页面只使用现有 service/API DTO，不新增后端契约。

## 5. Validation

- [x] 5.1 运行 `npm --prefix web test -- --run`。
- [x] 5.2 运行 `npm --prefix web run build`。
- [x] 5.3 运行 `go test ./...`。
- [x] 5.4 启动真实本地后端和 Vite 前端。
- [x] 5.5 采集核心路由 390px、768px、1280px 截图或等价 browser evidence。
- [x] 5.6 检查页面级横向溢出，仅允许局部二维容器滚动。
- [x] 5.7 执行 forbidden copy scan。
- [x] 5.8 执行敏感信息/redaction scan。
- [x] 5.9 运行 `openspec validate p110-visual-system-redesign --strict` 与 `openspec validate --all --strict`。
- [x] 5.10 运行 `git diff --check`。

## 6. Documentation, Review, Archive

- [x] 6.1 新增 P110 UI 验收记录和截图/视觉证据资产。
- [x] 6.2 更新 `docs/ui-design.md`、`docs/frontend-contract.md`、`docs/product-experience-polish-roadmap.md`、`docs/development-plan.md` 中与 P110 相关的视觉系统记录。
- [x] 6.3 子 agent 或等价复审无 Critical / Important 后归档；本环境工具限制未在用户未显式要求时 spawn subagent，已执行本地等价 diff/门禁复审。
- [x] 6.4 执行 OpenSpec archive，合并 delta 并更新进度。
- [x] 6.5 提交 P110。
