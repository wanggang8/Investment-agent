# P8 复审修复记录

> Change：`p8-full-review-fixes`  
> 日期：2026-06-01

## 本轮复审发现

- 后端：守门人审计样本数、流动性裁决、市场刷新指标。
- 数据：手动买入资金校验、证据稳定哈希、A/S 独立信源数量映射。
- 前端：规则版本缺失状态、Dashboard 确认入口、Portfolio 空态与买入理由。
- 测试与治理：能力圈排除集成测试、fetch mock 清理、验收记录可追溯性。

## 已处理

- 修复守门人审计使用真实 `proposal.status` 与 `proposal.sample_count` 推进状态。
- 将 `liquidity_state` 纳入规则裁决，危险流动性禁止新增买入和大额市价操作，且不再保留买入类可选动作。
- 市场刷新生成 `market_metrics_json`，保留 `close_price` 与 `turnover_rate`。
- 手动买入在写入任何确认事实前校验现金充足。
- 证据刷新使用稳定内容哈希与 chunk 哈希，规则输入保留高等级独立信源数量。
- 补能力圈 `excluded_symbols_json` 主动咨询集成测试。
- 前端将 `RULE_VERSION_MISSING` 映射为 `high_risk`，Dashboard 改为详情页确认入口，Portfolio 增加空态与买入理由列。
- `client.test.ts` 改为 `afterEach` 统一清理全局 fetch mock。

## 定向验证

- `go test ./internal/domain/rule ./internal/application/service ./internal/application/handler`：通过，`Go test: 92 passed in 3 packages`。
- `go test ./internal/application/workflow ./internal/application/service ./internal/infrastructure/persistence/sqlite`：通过，`Go test: 81 passed in 3 packages`。
- `go test ./internal/application/handler ./internal/application/workflow ./internal/application/service ./internal/domain/rule ./internal/infrastructure/persistence/sqlite`：通过，`Go test: 157 passed in 5 packages`。
- `go test ./internal/domain/rule ./internal/application/workflow ./internal/application/service ./internal/infrastructure/persistence/sqlite`：通过，`Go test: 98 passed in 4 packages`。
- `npm test -- --run src/services/client.test.ts src/features/dashboard/DashboardFeature.test.tsx src/pages/PortfolioPage.test.tsx`：通过，3 个测试文件、8 个测试。

## 最新复审处理记录

- 2026-06-01 17:29 复审结论：无 Critical；Important 集中在前端确认/咨询入口、费用字段、市场结构列、低样本风险说明、确认状态语义、架构边界与测试证据。
- 已处理：Dashboard 详情链接使用 `decision_id`；`/consultation` 接入主动咨询表单；错误原因标签改为契约枚举；审计时间线展示时间、工作流、状态变化、规则版本、快照。
- 已处理：手动执行贯通 `fees`，买入校验含费用，卖出现金扣费用；`watch` 后拒绝转手动执行。
- 已处理：市场刷新写入 `close_price`、`turnover_rate` 结构列；低样本规则提案写入 `risk_notes_json`；预期收益样本数持久化回放；DeepSeek 客户端改依赖 domain analyst 接口，移除 infrastructure 到 application/workflow 的反向依赖。

## 最新验证

- `go test ./internal/application/handler ./internal/infrastructure/persistence/sqlite ./internal/application/workflow ./internal/infrastructure/llm/deepseek`：通过，`Go test: 126 passed in 4 packages`。
- `npm test -- --run src/features/dashboard/DashboardFeature.test.tsx src/pages/DecisionDetailPage.test.tsx src/components/audit/AuditEventTimeline.test.tsx`：通过，3 个测试文件、8 个测试。
- `go test ./...`：通过，`Go test: 177 passed in 24 packages`。
- `cd web && npm run build && npm test`：通过，`✓ built in 97ms`；Vitest `17 passed (17)`、`44 passed (44)`。

## 第三轮复审修复记录

- 已处理：主动咨询 `scenario` 与契约示例统一为 `hold_review` / `buy_review` / `sell_review` / `rebalance_review`，前后端测试覆盖。
- 已处理：手动买入费用计入现金扣减、交易流水和持仓成本价；市场刷新在快照写入失败时生成独立失败审计。
- 已处理：规则提案页展示来源误判案例、提案理由、守门人结果、审计摘要、影响范围、风险提示与变更前后规则。
- 已处理：补齐审计字段、`root_cause_tag` 枚举提交、`RULE_VERSION_MISSING`、market 结构列、hash 内容差异等测试证据。
- 已处理：OpenSpec delta 同步更新 `p0-p8-review-fixes`，tasks 8.1-8.6 全部完成。

## 第三轮修复验证

- `go test ./...`：通过，`Go test: 178 passed in 24 packages`。
- `cd web && npm run build && npm test`：通过，`✓ built in 89ms`；Vitest `17 passed (17)`、`46 passed (46)`。

## 第四轮复审修复记录

- 已处理：市场快照写入失败只保留 graph 层单条失败审计，handler 不再追加重复失败事件。
- 已处理：`SourceVerificationDTO` 返回 `high_grade_independent_source_count`，并补接口测试。
- 已处理：DeepSeek 测试改依赖 domain analyst 请求类型，避免 infrastructure 测试继续引用 application/workflow。
- 已处理：确认错误案例字段持久化断言、卖出费用现金断言、非法 consult scenario 拒绝断言、测试命名修正。
- 已处理：前端审计事件兼容 `event_id`，手动执行表单支持费用输入，规则提案审计结果中文展示，规则内容优先展示 `content`。
- 已处理：OpenSpec delta 补 scenario、hash、高等级信源计数、失败审计、费用字段说明；`docs/testing-plan.md` 同步最新全量验收结果。

## 第四轮修复验证

- `go test ./...`：通过，`Go test: 181 passed in 24 packages`。
- `cd web && npm run build && npm test`：通过，`✓ built in 85ms`；Vitest `17 passed (17)`、`47 passed (47)`。

## 第五轮复审修复记录

- 已处理：tasks 增加第 9 节，OpenSpec delta 补前端 `event_id` 兼容、手动执行费用输入、规则内容展示和审计 DTO 字段要求。
- 已处理：补审计列表 API 全字段测试、市场刷新失败审计 SQLite 持久化断言、source verification 高等级信源数仓储读写断言。
- 已处理：主动咨询 `scenario` 收敛为 domain 枚举，测试覆盖全部合法契约值和非法值。
- 已处理：前端 Evidence 类型增加 `high_grade_independent_source_count`，`getEvidenceVerification()` 返回类型改为单个 verification DTO。
- 已处理：前端补强 `event_id` 稳定展开、费用非负有限数字校验、规则提案审计结果与 content 展示、高等级信源数展示测试。
- 已处理：手动执行费用后端断言改为容差比较，减少浮点精度导致的测试脆弱性。

## 第五轮修复验证

- `go test ./internal/application/handler ./internal/application/workflow ./internal/infrastructure/persistence/sqlite`：通过，`Go test: 135 passed in 3 packages`。
- `npm test -- --run src/components/dashboard/UserConfirmationPanel.test.tsx src/components/audit/AuditEventTimeline.test.tsx src/components/rules/RuleProposalPanel.test.tsx src/components/evidence/EvidenceTable.test.tsx`：通过，4 个测试文件、15 个测试。
- `go test ./...`：通过，`Go test: 188 passed in 24 packages`。
- `cd web && npm run build && npm test`：通过，`✓ built in 87ms`；Vitest `17 passed (17)`、`49 passed (49)`。

## 末轮只读复审修复记录

- 已处理：`EvidenceDTO`、证据列表和决策证据链补 `high_grade_independent_source_count` 返回链路，并补接口与仓储测试。
- 已处理：证据刷新写入 `latest_published_at` 与 `evidence_ids_json`，增强 source verification 可追溯性。
- 已处理：`operation_confirmations` 结构化保存 `fees`，确认记录、交易流水、现金和成本价链路一致。
- 已处理：consult scenario 合法值测试改为断言成功链路；前端补 evidence service 真实 fetch 契约测试。
- 已处理：市场和证据浮点断言改为容差比较，补 DeepSeek infrastructure 不依赖 application 的 import 边界测试。

## 末轮修复验证

- `go test ./internal/application/handler ./internal/application/workflow ./internal/infrastructure/persistence/sqlite`：通过，`Go test: 137 passed in 3 packages`。
- `npm test -- --run src/services/evidence.test.ts src/components/evidence/EvidenceTable.test.tsx`：通过，2 个测试文件、3 个测试。
- `go test ./...`：通过，`Go test: 190 passed in 24 packages`。
- `cd web && npm run build && npm test`：通过，`✓ built in 86ms`；Vitest `18 passed (18)`、`51 passed (51)`。
