# Tasks: P55 前端全功能真实验收与设计审查

## 1. 范围与计划

- [x] 1.1 确认当前无活跃 change，P55 为用户要求的真实 UI 验收阶段。
- [x] 1.2 创建 `p55-full-ui-acceptance-and-design-audit` OpenSpec change。
- [x] 1.3 确认 P55 真实启动项目、通过浏览器操作 UI 验收全部主要页面，并做 Product Design audit；不修改运行时代码。
- [x] 1.4 子 agent 复审 P55 计划，无 Critical / Important 后继续。

## 2. 环境启动

- [x] 2.1 创建 `tmp/ui-acceptance/p55-2026-06-17/`。
- [x] 2.2 生成临时配置和 SQLite。
- [x] 2.3 运行 smoke seed，确保前端页面有可验收数据。
- [x] 2.4 启动后端 server。
- [x] 2.5 启动 Vite 前端 dev server。
- [x] 2.6 使用浏览器打开本地前端并确认健康检查通过。

## 3. UI 全功能验收

- [x] 3.1 验收全局导航、布局、响应式基础和安全边界。
- [x] 3.2 验收 Dashboard `/`。
- [x] 3.3 验收 Workbench `/workbench`。
- [x] 3.4 验收 Decision Loop `/decision-loop`。
- [x] 3.5 验收 Data Quality `/data-quality`。
- [x] 3.6 验收 Portfolio `/positions`。
- [x] 3.7 验收 Consultation `/consultation` 和 Decision Detail `/decisions/:decisionId`；发现真实 LLM 决策详情 nullable `optional_actions` 前端崩溃 blocker。
- [x] 3.8 验收 Evidence `/evidence`。
- [x] 3.9 验收 Rules `/rules`。
- [x] 3.10 验收 Audit `/audit`。
- [x] 3.11 验收 Notifications `/notifications`。
- [x] 3.12 验收 Risk Alerts `/risk-alerts` 和详情路径。
- [x] 3.13 验收 Daily Auto Run `/daily-auto-run`。
- [x] 3.14 验收 Daily Discipline 列表和详情。
- [x] 3.15 验收 Review `/review`。
- [x] 3.16 验收 Local Install `/local-install`。
- [x] 3.17 验收 Local Knowledge `/local-knowledge`。
- [x] 3.18 验收 Settings `/settings`。

## 4. Product Design 审查

- [x] 4.1 保存每个主要页面截图到 `docs/release/ui-audit-assets/2026-06-17-p55/`。
- [x] 4.2 使用 Product Design audit 框架记录 UX、视觉、一致性和可访问性风险。
- [x] 4.3 区分 `blocked`、`needs_optimization`、`minor` 和 `pass`。
- [x] 4.4 确认设计审查不宣称完整 WCAG 合规，只记录截图和浏览器操作可见证据。

## 5. 报告与同步

- [x] 5.1 新增 `docs/release/ui-acceptance-2026-06-17.md`。
- [x] 5.2 新增 `docs/release/ui-design-audit-2026-06-17.md`。
- [x] 5.3 更新 `openspec/PROGRESS.md`，标记 P55 活跃。
- [x] 5.4 更新 `openspec/project.md`、`docs/GOVERNANCE.md`、`docs/development-plan.md`、`AGENTS.md`。
- [x] 5.5 更新 `docs/README.md`，加入 UI 验收和设计审查入口。

## 6. 验证与复审

- [x] 6.1 执行 `openspec validate p55-full-ui-acceptance-and-design-audit --strict`。
- [x] 6.2 执行 `openspec validate --all --strict`。
- [x] 6.3 执行 `git diff --check`。
- [x] 6.4 执行提交材料脱敏扫描。
- [x] 6.5 子 agent 复审最终变更，无 Critical / Important 后归档。

## 7. 归档

- [x] 7.1 执行 OpenSpec archive，将 delta 合并入对应 `docs/` 真源，并按需同步 `openspec/specs/` 摘要。
- [x] 7.2 归档后确认无活跃 change。
- [x] 7.3 提交前子 agent 复审无 Critical / Important。
- [x] 7.4 提交 P55。
