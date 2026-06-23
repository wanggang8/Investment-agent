## 1. OpenSpec 与范围

- [x] 1.1 确认 P41 已归档，当前无其他活跃 change。
- [x] 1.2 确认 P42 只聚合现有本地事实、DTO 和安全导航；不得新增券商 API、自动交易、外部推送、自动确认或自动规则应用。
- [x] 1.3 确认 P42 默认不新增数据库 migration；如必须扩展 DTO，只允许只读聚合字段。

## 2. 契约与文档

- [x] 2.1 更新 `docs/frontend-contract.md`，定义用户决策工作台页面区域、数据来源、状态和安全边界。
- [x] 2.2 更新 `docs/development-plan.md`，加入 P42 已立项目标、任务和验收命令。
- [x] 2.3 更新 `docs/GOVERNANCE.md`、`AGENTS.md`、`openspec/PROGRESS.md` 和 `openspec/project.md` 的 P42 active 状态。

## 3. 前端工作台

- [x] 3.1 新增 `/workbench` 页面或等价首屏入口，展示“今日先看”“组合与风险”“规则与复盘”“主动咨询入口”四类信息。
- [x] 3.2 工作台复用现有 services/API DTO，不直接读取 SQLite、VecLite 或本地文件。
- [x] 3.3 工作台提供到每日纪律报告、持仓、风险预警、规则提案、复盘摘要、审计和决策咨询的导航入口。
- [x] 3.4 工作台在空库、数据缺失、source health 降级、LLM/RAG/VecLite 不可用时显示明确安全状态。
- [x] 3.5 工作台不得展示自动交易、一键交易、代下单、自动外推、自动确认或自动应用规则入口。

## 4. 测试与 E2E

- [x] 4.1 增加工作台页面或组件 Vitest，覆盖成功、空库、降级、错误和安全文案。
- [x] 4.2 扩展 Playwright smoke，覆盖 `/workbench` 可达、核心区域可见、窄屏可用和禁止入口扫描。
- [x] 4.3 运行 `npm --prefix web test -- --run`。
- [x] 4.4 运行 `npm --prefix web run build`。
- [x] 4.5 运行 `bash scripts/e2e-smoke.sh`。

## 5. 验收与归档

- [x] 5.1 如修改后端，运行 `go test ./...`；如未修改后端，在任务记录中说明原因：P42 未修改后端代码、migration 或 Go DTO。
- [x] 5.2 运行 `openspec validate p42-user-decision-workbench --strict`。
- [x] 5.3 运行 `openspec validate --all --strict`。
- [x] 5.4 运行 `git diff --check`。
- [x] 5.5 运行 archive 前只读复审，且无 Critical / Important 问题。
