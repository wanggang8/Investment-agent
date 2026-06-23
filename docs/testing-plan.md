# P6 验收计划

> 对齐范围：`docs/development-plan.md` P6.1 与 `docs/functional-spec.md` A01-A17。  
> 验收目标：确认端到端行为、降级状态、审计记录、配置文档和禁止自动交易边界满足既有契约。

## 1. 验收方式

- 后端以 API、Repository、工作流和审计事件的可观测结果为准。
- 前端以页面字段、状态文案、操作入口和构建结果为准。
- 外部数据源、DeepSeek、VecLite 异常场景可使用本地 stub、fixture、空数据或依赖注入方式验证。
- 验收期间不得写入真实密钥，不得调用真实交易接口。

## 2. A01-A17 端到端验收断言

### A01 首次使用

- 前置条件：本地 SQLite 已完成 migration，但没有账户快照、持仓和决策记录。
- 操作：请求 `GET /api/v1/dashboard/today`，或打开前端今日纪律页。
- 期望：接口返回 `DATA_REQUIRED` 或 `dashboard_state=first_use`；页面展示初始化引导。
- 数据检查：`decision_records` 不新增记录。
- 关联任务：P6.1 A01。

### A02 正常每日纪律

- 前置条件：存在账户快照、持仓、市场快照、规则版本和满足正式裁决的证据。
- 操作：触发每日纪律工作流或请求今日纪律 API。
- 期望：生成正式建议，并返回证据摘要、触发规则和最终裁决。
- 数据检查：`decision_records` 增加 1 条，`evidence_refs` 至少 1 条；`audit_events` 至少包含 `StateSnapshotNode`、`EvidenceRetrievalNode`、`RuleArbitrationNode`、`DecisionRecordNode`。
- 关联任务：P6.1 A02。

### A03 证据不足

- 前置条件：目标标的缺少可用正式证据。
- 操作：触发每日纪律或主动咨询。
- 期望：返回 `EVIDENCE_NOT_FOUND`，`final_verdict.status=insufficient_data`。
- 前端检查：不展示可执行的交易类建议，展示缺失项和暂停原因。
- 关联任务：P6.1 A03。

### A04 VecLite 不可用

- 前置条件一：VecLite 索引不可用，但 SQLite 中存在足够 `intelligence_summary` 与 `rag_chunks`。
- 操作一：触发证据检索或证据页刷新。
- 期望一：`workflow_status=degraded`，页面展示降级说明和可用摘要。
- 前置条件二：VecLite 不可用，且 SQLite 摘要不足。
- 操作二：触发同一流程。
- 期望二：页面状态为 `insufficient_data`。
- 关联任务：P6.1 A04。

### A05 能力圈外

- 前置条件：咨询标的不在能力圈配置内。
- 操作：发起主动咨询。
- 期望：`final_verdict.status=rejected`。
- 工作流检查：不得调用 `ValueAnalystNode` 与 `TrendRiskOfficerNode`。
- 关联任务：P6.1 A05。

### A06 用户记录计划

- 前置条件：存在可确认的决策记录。
- 操作：提交 `confirmation_type=planned`。
- 期望：写入 `operation_confirmations` 与 `audit_events`。
- 数据检查：不写 `position_transactions`，不新增账户快照。
- 关联任务：P6.1 A06。

### A07 用户记录已手动执行

- 前置条件：存在可确认的决策记录，并提供线下成交信息。
- 操作：提交 `confirmation_type=executed_manually`。
- 期望：同一事务写入 `operation_confirmations`、`position_transactions`、`positions`、`portfolio_snapshots`、`position_snapshots`、`audit_events`。
- 页面检查：文案表达为记录线下动作，不出现自动交易含义。
- 关联任务：P6.1 A07。

### A08 已手动执行失败

- 前置条件：构造交易流水或快照写入失败场景。
- 操作：提交 `confirmation_type=executed_manually`。
- 期望：接口返回失败错误。
- 数据检查：事务回滚，不留下部分确认记录、交易流水或快照。
- 关联任务：P6.1 A08。

### A09 用户标记错误

- 前置条件：存在用户认为有误的决策记录。
- 操作：提交 `confirmation_type=marked_error`。
- 期望：同一事务写入 `operation_confirmations`、`error_cases`、`audit_events`。
- 响应检查：返回 `error_case_id`。
- 关联任务：P6.1 A09。

### A10 C 级信源

- 前置条件：证据来源等级为 C。
- 操作：刷新证据并触发裁决。
- 期望：C 级证据只能以 `evidence_role=background` 返回。
- 数据检查：C 级证据不得出现在正式裁决的 `formal` 证据引用中。
- 关联任务：P6.1 A10。

### A11 LLM 不可用

- 前置条件：DeepSeek 调用失败或配置为空。
- 操作：触发包含分析节点的工作流。
- 期望：返回 `ANALYST_UNAVAILABLE`，`workflow_status=degraded`。
- 裁决检查：最终裁决来自规则引擎，LLM 输出不覆盖最终裁决。
- 关联任务：P6.1 A11。

### A12 守门人审计通过

- 前置条件：存在 `under_gatekeeper_audit` 状态且 `sample_count>=3` 的规则提案。
- 操作：触发守门人审计通过路径。
- 期望：提案状态变为 `pending_final_confirm`，不写 `rule_versions`。
- 负向检查：`sample_count<3` 的提案不得进入守门人审计，接口返回 `BAD_REQUEST`。
- 关联任务：P6.1 A12。

### A13 规则最终确认

- 前置条件：存在 `pending_final_confirm` 状态且 `sample_count>=3` 的规则提案。
- 操作：提交最终确认 `confirm=true`。
- 期望：创建新 active `rule_versions`，旧 active 归档，提案状态为 `applied`。
- 负向检查：`sample_count<3` 的提案最终确认返回 `BAD_REQUEST`，不得写 `rule_versions`。
- 关联任务：P6.1 A13。

### A14 审计事件

- 前置条件：系统已有工作流、确认、规则或市场刷新审计事件。
- 操作：请求审计页数据并打开前端审计页。
- 期望：接口字段同时包含 `action`、`node_name`、`node_action`。
- 前端检查：分别展示业务动作、节点名称和节点动作，并展示 `status`、`error_code`、输入引用、输出引用。
- 关联任务：P6.1 A14。

### A15 禁止自动交易

- 前置条件：后端路由和前端页面已构建。
- 操作：审查 API 路由、前端确认区和页面文案。
- 期望：API 列表中不存在买入、卖出、撤单、改单接口。
- 前端检查：确认区不出现一键交易、自动下单、委托执行等入口或文案；用户操作只记录线下动作。
- 关联任务：P6.1 A15。

### A16 市场数据刷新

- 前置条件：市场刷新 API 可用，支持成功、部分失败、全部失败和写入失败测试数据。
- 操作一：请求 `POST /api/v1/market/refresh`，所有标的成功。
- 期望一：新增 `market_snapshots`，`audit_events.status=success`。
- 操作二：构造部分标的失败。
- 期望二：返回 200，写入成功标的，返回 `failed_symbols`，`audit_events.status=degraded`。
- 操作三：构造全部标的失败。
- 期望三：返回 `DATA_SOURCE_UNAVAILABLE` 或 `DATA_STALE`。
- 操作四：构造快照写入失败。
- 期望四：返回 `MARKET_SNAPSHOT_WRITE_FAILED`，不留下部分写入。
- 关联任务：P6.1 A16。

### A17 预期收益评估

- 前置条件：存在不同样本数量的预期收益输入。
- 操作：触发每日纪律或主动咨询，并查看决策详情页。
- 期望：预期收益只作为情景概率展示，不覆盖最终规则裁决，不承诺收益。
- 状态检查：`available` 必须包含 upside/base/downside 且可返回概率；`insufficient` 不得返回精确概率且必须写样本不足说明；`unavailable` 必须返回空 `scenarios` 并写定性原因。
- 关联任务：P6.1 A17。

## 3. 阶段验收命令

```bash
go test ./...
go run ./cmd/agent --help
go test ./internal/application/workflow/... ./internal/application/handler/...
./scripts/e2e-smoke.sh
E2E_SERVER_PORT=18081 E2E_WEB_PORT=14174 ./scripts/e2e-smoke.sh
npm --prefix web test -- --run
npm --prefix web run build
openspec validate --all --strict
```

P32 每日纪律报告产品化定向验收命令：

```bash
go test ./cmd/smoke-seed ./internal/infrastructure/persistence/sqlite
npm --prefix web test -- --run
npm --prefix web run build
E2E_SERVER_PORT=18081 E2E_WEB_PORT=14174 ./scripts/e2e-smoke.sh
openspec validate p32-daily-discipline-report-productization --strict
```

P32 smoke 断言覆盖今日纪律报告、历史报告列表、报告详情、报告摘要和“不会自动执行交易”安全文案；E2E seed 使用本地临时 SQLite，不写真实账户、密钥或交易接口。
说明：`npm test` 纳入阶段验收，用于覆盖前端 API client、页面空态/错误态、确认失败、规则最终确认和禁止自动交易入口。`cmd/agent` 验收只确认本地任务入口和安全边界，不代表开启后台调度或自动交易。P31 每日自动运行验收还应覆盖默认关闭、启用后按 `run_time`/`timezone` 等待触发、缺持仓失败诊断、执行副作用前写 `running` 幂等状态、重复运行复用、重试、超时、应用内通知、审计摘要和 `/daily-auto-run` 前端状态页。

## 4. 验收结果记录

- P32 Task 5 计划/执行记录（2026-06-08）：先扩展 `web/e2e/local-smoke.spec.ts`，在未写入 P32 report seed 时运行 `E2E_SERVER_PORT=18081 E2E_WEB_PORT=14174 ./scripts/e2e-smoke.sh`，确认失败于缺少 `P32 smoke 今日纪律报告已生成`；随后实现 `cmd/smoke-seed` 写入 report index。已运行：`go test ./cmd/smoke-seed ./internal/infrastructure/persistence/sqlite` 通过；`npm --prefix web test -- --run` 通过（25 files / 78 tests）；`npm --prefix web run build` 通过；`E2E_SERVER_PORT=18082 E2E_WEB_PORT=14175 ./scripts/e2e-smoke.sh` 通过（Playwright 1 passed）。原指定 `E2E_SERVER_PORT=18081 E2E_WEB_PORT=14174 ./scripts/e2e-smoke.sh` 因 18081 已有 stale server listener / bind 冲突未完成；未声称该端口组合通过。`openspec validate p32-daily-discipline-report-productization --strict` 已通过。

- `go test ./...`：通过（2026-06-07，P31 复审修复后最终验收），输出 `Go test: 398 passed in 25 packages`。
- `E2E_SERVER_PORT=18081 E2E_WEB_PORT=14174 ./scripts/e2e-smoke.sh`：通过（2026-06-07，P31 每日自动运行闭环验收），Playwright `1 passed`，覆盖健康检查、决策详情 expected return/动态卖出提示、证据页、审计页和每日自动运行失败状态页。
- `go test ./cmd/smoke-seed ./cmd/server ./internal/application/workflow ./internal/application/handler`：通过（2026-06-07，P31 每日自动运行闭环验收），输出 `Go test: 272 passed in 4 packages`；复审修复后相关后端聚焦测试输出 `Go test: 318 passed in 5 packages`。
- `npm --prefix web test -- --run`：通过（2026-06-07，P31 每日自动运行闭环验收），Vitest `22 passed (22)`、`68 passed (68)`。
- `npm --prefix web run build`：通过（2026-06-07，P31 每日自动运行闭环验收），输出 `✓ built in 86ms`。
- `git status --short`：通过（2026-06-07，P31 6.x 验收），未出现 scheduler、E2E 或 smoke 临时产物。
- `./scripts/e2e-smoke.sh`：通过（2026-06-07，`p30-real-e2e-smoke` 验收），Playwright `1 passed`，覆盖健康检查、决策详情 expected return/动态卖出提示、证据页和审计页。
- `go test ./...`：通过（2026-06-07，`p30-real-e2e-smoke` 验收），输出 `Go test: 380 passed in 25 packages`。
- `npm --prefix web test -- --run`：通过（2026-06-07，`p30-real-e2e-smoke` 验收），Vitest `20 passed (20)`、`65 passed (65)`。
- `npm --prefix web run build`：通过（2026-06-07，`p30-real-e2e-smoke` 验收），输出 `✓ built in 93ms`。
- `openspec validate --all --strict`：通过（2026-06-07，`p30-real-e2e-smoke` 验收），输出 `18 passed, 0 failed`。
- `git status --short`：通过（2026-06-07，`p30-real-e2e-smoke` 验收），未出现 `.playwright-mcp/`、临时 SQLite、日志、trace 或 screenshot 等 smoke 产物。
