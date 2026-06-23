## 1. OpenSpec 与范围

- [x] 1.1 确认 P42 已归档，当前无其他活跃 change。
- [x] 1.2 确认 P43 只聚合现有质量状态、DTO 和安全导航；不得新增券商 API、自动交易、外部推送、自动确认或自动规则应用。
- [x] 1.3 确认 P43 默认不新增数据库 migration；如必须扩展后端，只允许只读聚合 DTO/API。

## 2. 契约与文档

- [x] 2.1 更新 `docs/frontend-contract.md`，定义数据质量可观测页面区域、数据来源、状态、脱敏和安全边界。
- [x] 2.2 更新 `docs/development-plan.md`，加入 P43 已立项目标、任务和验收命令。
- [x] 2.3 更新 `docs/GOVERNANCE.md`、`AGENTS.md`、`openspec/PROGRESS.md` 和 `openspec/project.md` 的 P43 active 状态。

## 3. 数据质量可观测面

- [x] 3.1 新增 `/data-quality` 页面或等价入口，展示“数据源健康”“证据与检索”“LLM 质量”“影响范围与下一步”四类信息。
- [x] 3.2 页面复用现有 services/API DTO；如新增聚合 DTO，必须只读且不得写数据库。
- [x] 3.3 页面提供到设置、证据、复盘、审计、风险预警、决策详情和工作台的导航入口。
- [x] 3.4 页面在空库、source_unavailable、parse_error、stale、missing、unknown、LLM/RAG/VecLite 不可用时显示明确安全状态。
- [x] 3.5 页面不得展示自动刷新修复、外部推送、自动确认、自动应用规则、自动交易、一键交易、代下单或收益承诺入口。
- [x] 3.6 页面不得展示完整 key、完整 prompt、私有本地路径、SQL、供应商原始错误或账户敏感明细。

## 4. 测试与 E2E

- [x] 4.1 增加数据质量页面或组件 Vitest，覆盖成功、空库、降级、错误、unknown 和脱敏安全文案。
- [x] 4.2 扩展 Playwright smoke，覆盖 `/data-quality` 可达、核心区域可见、窄屏可用和禁止入口扫描。
- [x] 4.3 运行 `npm --prefix web test -- --run`。
- [x] 4.4 运行 `npm --prefix web run build`。
- [x] 4.5 运行 `bash scripts/e2e-smoke.sh`。

## 5. 验收与归档

- [x] 5.1 如修改后端，运行 `go test ./...`；如未修改后端，在任务记录中说明原因。（P43 仅修改前端与文档，未修改后端，故无需运行 `go test ./...`。）
- [x] 5.2 运行 `openspec validate p43-data-quality-observability --strict`。
- [x] 5.3 运行 `openspec validate --all --strict`。
- [x] 5.4 运行 `git diff --check`。
- [x] 5.5 执行实现前只读审查，且无 Critical / Important 问题后再开始实现。
- [x] 5.6 执行 archive 前只读复审，且无 Critical / Important 问题。
