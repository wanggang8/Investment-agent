## 1. OpenSpec 与范围

- [x] 1.1 确认 P46 已归档、当前无活跃 change，P47 为下一功能候选。
- [x] 1.2 确认 P47 聚焦只读决策闭环解释：建议、确认、线下记录、风险/审计/复盘线索和缺口说明。
- [x] 1.3 确认 P47 不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复承诺、收益承诺、登录/付费/授权/Level2/高频源。

## 2. 后端 DTO、Repository 与服务

- [x] 2.1 新增 `internal/application/dto/decision_loop.go`，定义 list/detail response、loop item、stage、manual action、link。
- [x] 2.2 扩展 `internal/domain/repository/decision_repo.go`，新增只读方法 `ListOperationConfirmations(ctx, decisionID string)` 与 `ListPositionTransactionsByConfirmation(ctx, confirmationID string)`。
- [x] 2.3 在 `internal/infrastructure/persistence/sqlite/decision_repo_impl.go` 实现上述只读查询，按 `created_at`/`occurred_at` 排序；不得新增 migration 或写路径。
- [x] 2.4 更新测试 stub 编译点，确保新增 repository 接口不破坏既有单测。
- [x] 2.5 新增 `internal/application/service/decision_loop.go`，从 decisions、confirmations、transactions、error cases、risk alerts、audit events 聚合 loop item。
- [x] 2.6 服务必须限制 `limit` 默认 10、最大 50，支持 `symbol` 过滤，并对 payload/note 做安全摘要，不返回原始 JSON。
- [x] 2.7 新增 `internal/application/service/decision_loop_test.go`，覆盖完整闭环、planned/watch/not_required、executed_manually 缺交易、risk/audit/error 链接、limit/symbol、安全摘要。

## 3. 后端 Handler 与 API

- [x] 3.1 新增 `internal/application/handler/decision_loop_handler.go`。
- [x] 3.2 在 `internal/application/handler/app.go` 初始化 `DecisionLoopSvc` 并注册：
  - `GET /api/v1/decision-loops`
  - `GET /api/v1/decision-loops/{decision_id}`
- [x] 3.3 新增 `internal/application/handler/decision_loop_handler_test.go`，覆盖列表、详情、not found、limit 参数、安全响应不泄露 payload/SQL/private path。

## 4. 前端页面与服务

- [x] 4.1 新增 `web/src/types/decisionLoop.ts` 与 `web/src/services/decisionLoop.ts`。
- [x] 4.2 新增 `web/src/pages/DecisionLoopPage.tsx`，展示闭环列表、阶段、缺口、人工记录和链接。
- [x] 4.3 新增 `web/src/pages/DecisionLoopPage.test.tsx`，覆盖成功展示、空态、错误态、缺口展示、无写入动作按钮。
- [x] 4.4 在 `web/src/App.tsx` 注册 `/decision-loop`，在 `web/src/app/AppLayout.tsx` 增加导航入口。
- [x] 4.5 在 `web/src/pages/WorkbenchPage.tsx` 与 `web/src/pages/ReviewSummaryPage.tsx` 增加只读导航链接到 `/decision-loop`，并更新对应测试。
- [x] 4.6 更新 `web/e2e/local-smoke.spec.ts`，覆盖 `/decision-loop` 可达与安全文本扫描。

## 5. 文档与契约

- [x] 5.1 更新 `docs/api.md`，新增 P47 decision loop list/detail API。
- [x] 5.2 更新 `docs/data-model.md`，说明 P47 只读复用现有事实表和新增只读查询边界。
- [x] 5.3 更新 `docs/frontend-contract.md`，新增 `/decision-loop` 页面契约。
- [x] 5.4 更新 `docs/development-plan.md`、`openspec/project.md`、`openspec/PROGRESS.md`、`docs/GOVERNANCE.md` 和 `AGENTS.md` 当前阶段状态。
- [x] 5.5 在 OpenSpec delta 中记录 P47 行为要求。

## 6. 执行前复审

- [x] 6.1 计划完成后执行只读子 agent 复审，确认无 Critical / Important。
- [x] 6.2 复审通过后再执行实现任务。

## 7. 验收

- [x] 7.1 运行 `go test ./...`。
- [x] 7.2 运行 `npm --prefix web test -- --run`。
- [x] 7.3 运行 `npm --prefix web run build`。
- [x] 7.4 运行 `bash scripts/e2e-smoke.sh`。
- [x] 7.5 运行 `openspec validate p47-decision-loop-explainability --strict`。
- [x] 7.6 运行 `openspec validate --all --strict`。
- [x] 7.7 运行 `git diff --check`。
- [x] 7.8 运行安全扫描：`rg -n 'sk-[A-Za-z0-9][A-Za-z0-9_-]{8,}|BEGIN (RSA|OPENSSH|PRIVATE) KEY|/Users/[^[:space:]，；。、]+|(?i:select[[:space:]]+\*[[:space:]]+from)|(?i:raw[[:space:]]+http)|(?i:prompt[[:space:]]*:)|完整[[:space:]]*prompt|HTTP/[0-9.]+[[:space:]]+[0-9]{3}|券商接口|自动交易|一键交易|代下单|外部推送|自动确认|自动应用规则|收益承诺|Level2|高频源' web/src/pages/DecisionLoopPage.tsx web/src/services/decisionLoop.ts web/src/types/decisionLoop.ts internal/application/dto/decision_loop.go internal/application/service/decision_loop.go internal/application/handler/decision_loop_handler.go docs/api.md docs/frontend-contract.md docs/data-model.md`，人工复核命中项，确认不存在未脱敏敏感内容或高风险操作入口；允许安全边界说明文本命中。

## 8. 归档前复审

- [x] 8.1 执行完成后再次只读子 agent 复审，确认无 Critical / Important。
- [x] 8.2 复审通过后执行 archive，并将 P47 归档。
