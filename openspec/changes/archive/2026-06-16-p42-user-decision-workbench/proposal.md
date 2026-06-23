## Why

P41 已将 P40 后首个产品能力增强阶段推荐为 P42 `p42-user-decision-workbench`。当前系统已经具备今日纪律、账户/持仓、风险预警、规则治理、复盘摘要、决策咨询和浏览器 E2E 基础，但这些入口分散在多个页面。用户每天打开系统时仍需要在多个页面之间拼接“今天要看什么、缺什么、能做什么、不能做什么”。

本变更用于建立一个用户决策工作台，把现有本地事实、每日纪律、风险、规则、复盘和咨询入口组合成单一阅读与导航面板，同时继续保持只读展示、人工复核和不自动交易边界。

## What Changes

- 新增 P42 用户决策工作台页面或等价入口，聚合今日纪律、组合状态、风险预警、规则治理、复盘摘要和主动咨询入口。
- 工作台只使用现有 API/service DTO，不直接读取 SQLite、VecLite 或本地文件。
- 工作台提供跨页导航和检查清单，不自动触发交易、外部推送、自动确认或自动应用规则。
- 前端测试和浏览器 smoke 覆盖工作台的空库、数据完整、降级和安全文案路径。
- 更新 `docs/development-plan.md`、`docs/frontend-contract.md`、`docs/GOVERNANCE.md`、`AGENTS.md`、`openspec/PROGRESS.md` 和 `openspec/project.md` 的 P42 活跃状态。

## Capabilities

### New Capabilities

- `user-decision-workbench`: 覆盖 P42 用户决策工作台、聚合入口、跨页导航、安全边界和前端验收。

### Modified Capabilities

- `frontend-ops-review-surface`: 增加工作台对风险、规则、复盘和诊断状态的安全聚合要求。
- `frontend-experience-tests`: 增加工作台浏览器与组件测试要求。

## Impact

- 影响文档：`docs/development-plan.md`、`docs/frontend-contract.md`、`docs/GOVERNANCE.md`、`AGENTS.md`、`openspec/PROGRESS.md`、`openspec/project.md`。
- 影响前端：新增或扩展工作台页面、导航、组件测试和 E2E smoke。
- 预期不需要新增数据库 migration 或后端写接口；如发现现有 DTO 不足，必须在本 change 中明确只读 DTO 扩展，不得引入交易执行能力。
