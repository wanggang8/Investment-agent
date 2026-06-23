# P19-P24 审计证据包

> 更新时间：2026-06-17
>
> 作用：为 P19-P24 可用 MVP 阶段提供当前仓库事实证据矩阵。本文不是 P19-P24 的历史 OpenSpec archive，也不补写历史完成时间。

## 1. 总览结论

P19-P24 当前状态为：能力已交付并纳入当前文档真源，但缺少逐阶段完整 `openspec/changes/archive/` 包。仓库已有 `docs/p19-p24-historical-archive-traceability.md` 作为历史追溯说明，P51 在其基础上补充“可核验证据包”。

本文件可用于后续审计、P52 验收门禁矩阵和 P53 发布候选材料的引用基础，但不得单独用于声明“项目已发布就绪”。项目级验收门禁仍必须由 P52 明确定义。

## 2. P14-P18 旁证核查

用户若将范围写作“P19-P14”，按当前仓库事实应解释为 P14-P19 逆序核查。P14-P18 已有标准 archive，P51 不重复补档；P19 之后才是缺逐阶段 archive 的重点。

| 阶段 | archive 状态 |
| --- | --- |
| P14 | `openspec/changes/archive/2026-06-02-p14-gatekeeper-node-graph/` 已存在 |
| P15 | `openspec/changes/archive/2026-06-03-p15-evidence-quality-enrichment/` 已存在 |
| P16 | `openspec/changes/archive/2026-06-03-p16-frontend-ops-review-surface/` 已存在 |
| P17 | `openspec/changes/archive/2026-06-03-p17-local-scheduler-and-ops-docs/` 已存在 |
| P18 | `openspec/changes/archive/2026-06-03-p18-evolution-proposal-hardening/` 已存在 |

## 3. P19-P24 分阶段证据矩阵

### P19 公开 HTTP 数据桥接

| 项目 | 证据 |
| --- | --- |
| 当前交付边界 | 可配置公开 HTTP JSON 行情与情报来源，保留 fixture/stub fallback；支持稳定失败和降级。 |
| archive 状态 | 缺少逐阶段完整 archive 包；不得写成已 archive。 |
| 文档证据 | `docs/development-plan.md` P19 摘要；`docs/configuration.md` 数据源配置说明。 |
| 代码证据 | `internal/application/workflow/data_sources.go`：`ConfiguredMarketDataSource`、`ConfiguredIntelligenceSource`；`internal/infrastructure/wiring/workflow.go`：按配置接入工作流依赖。 |
| 测试证据 | `internal/application/workflow/data_sources_test.go` 覆盖 public HTTP、fixture、fallback、payload shape、中文公开信源等级和失败分类；`internal/application/handler/market_handler_test.go` 覆盖市场刷新失败通知与去重。 |
| 可重跑命令 | `go test ./internal/application/workflow -run 'TestConfigured(Market|Intelligence)Source'`；`go test ./internal/application/handler -run TestRefreshMarket` |
| 不可声明事项 | P19 不是“已接通所有真实外部公开源”；真实 collector 验证和接入由 P25-P29 补强。 |
| 残余缺口 | 缺 P19 原始 proposal/tasks/archive；只能以当前代码、测试和文档作为审计证据。 |

### P20 A 股 ETF/基金证据 payload 解析

| 项目 | 证据 |
| --- | --- |
| 当前交付边界 | 支持公开证据 payload 解析、URL/标题发布时间去重、多源验证、SQLite 事实写入；高等级独立信源不足时保持信息不足或降级。 |
| archive 状态 | 缺少逐阶段完整 archive 包；不得写成已 archive。 |
| 文档证据 | `docs/development-plan.md` P20 摘要；`docs/data-model.md` 证据事实表、`source_verifications`、RAG chunk 和审计事务边界；`docs/api.md` 证据写入要求。 |
| 代码证据 | `internal/application/workflow/evidence_verification_graph.go`；`internal/application/workflow/evidence_ingestion.go`；`internal/infrastructure/persistence/sqlite/intelligence_repo_impl.go`。 |
| 测试证据 | `internal/application/workflow/workflow_integration_test.go` 覆盖 `source_verifications`、`intelligence_items`、`intelligence_summary`、`rag_chunks` 写入；`internal/infrastructure/persistence/sqlite/intelligence_repo_impl_test.go` 覆盖持久化。 |
| 可重跑命令 | `go test ./internal/application/workflow -run 'Test.*(Evidence|SourceVerification|WorkflowIntegration)'`；`go test ./internal/infrastructure/persistence/sqlite -run Test.*Intelligence` |
| 不可声明事项 | P20 不是“已完成真实公告源历史补采”；P26/P29 才补强真实公告 collector 和真实采集 smoke。 |
| 残余缺口 | 缺 P20 原始 proposal/tasks/archive；真实源稳定性验收不属于 P20 原始交付范围。 |

### P21 应用内通知中心

| 项目 | 证据 |
| --- | --- |
| 当前交付边界 | 本地 `notifications` 表、Repository、Service、API、前端服务和通知页契约；支持列表、单条已读、全部已读、未读去重刷新。 |
| archive 状态 | 缺少逐阶段完整 archive 包；不得写成已 archive。 |
| 文档证据 | `docs/data-model.md` `notifications` 表；`docs/api.md` 通知接口；`docs/frontend-contract.md` 通知页契约。 |
| 代码证据 | `internal/application/service/notification.go`；`internal/application/handler/notification_handler.go`；`internal/infrastructure/persistence/sqlite/notification_repo_impl.go`；`web/src/services/notification.ts`。 |
| 测试证据 | `internal/application/handler/notification_handler_test.go`；`internal/infrastructure/persistence/sqlite/notification_repo_impl_test.go`；`web/src/services/notification.test.ts`。 |
| 可重跑命令 | `go test ./internal/application/handler -run TestNotificationHandlers`；`go test ./internal/infrastructure/persistence/sqlite -run TestNotificationRepository`；`npm --prefix web test -- --run web/src/services/notification.test.ts` |
| 不可声明事项 | 不提供邮件、短信、系统 Push、Webhook、WebSocket 或外部推送；通知不执行交易，不自动应用规则。 |
| 残余缺口 | 缺 P21 原始 proposal/tasks/archive；只能按当前本地通知实现和测试核验。 |

### P22 规则体系与提案增强

| 项目 | 证据 |
| --- | --- |
| 当前交付边界 | 规则提案包含 before/after payload、source facts、impact scope、risk notes、样本、原因和重复抑制；保留用户确认、守门人审计和最终确认流程。 |
| archive 状态 | 缺少逐阶段完整 archive 包；不得写成已 archive。 |
| 文档证据 | `docs/development-plan.md` P22 摘要；`docs/data-model.md` `rule_proposals` 与事务边界；`docs/api.md` 规则提案说明。 |
| 代码证据 | `internal/application/workflow/evolution_proposal_graph.go`；`internal/domain/rule/gatekeeper_logic.go`；`internal/application/service/rule_proposal.go`；`web/src/components/rules/RuleProposalPanel.tsx`。 |
| 测试证据 | `internal/application/workflow/workflow_integration_test.go` 覆盖 proposal 状态、通知、审计和守门人路径；`web/src/components/rules/RuleProposalPanel.test.tsx` 覆盖前端展示边界。 |
| 可重跑命令 | `go test ./internal/application/workflow -run 'Test.*(Evolution|Gatekeeper|RuleProposal)'`；`npm --prefix web test -- --run web/src/components/rules/RuleProposalPanel.test.tsx` |
| 不可声明事项 | 规则提案不会自动应用，不绕过用户确认，不触发交易。 |
| 残余缺口 | 缺 P22 原始 proposal/tasks/archive；规则效果回放和过拟合检查由 P36 补强。 |

### P23 复盘深化

| 项目 | 证据 |
| --- | --- |
| 当前交付边界 | 复盘 DTO 和页面展示可追溯归因摘要、错误标签、缺证据主题、提案结果、降级 workflow、ops status 和 tracking links；空窗口返回 empty/unknown。 |
| archive 状态 | 缺少逐阶段完整 archive 包；不得写成已 archive。 |
| 文档证据 | `docs/development-plan.md` P23 摘要；`docs/api.md` 复盘字段；`docs/frontend-contract.md` review summary 契约。 |
| 代码证据 | `internal/application/dto/review.go`；`internal/application/handler/review_handler.go`；`web/src/components/review/ReviewSummaryPanel.tsx`；`web/src/pages/ReviewSummaryPage.tsx`。 |
| 测试证据 | `internal/application/handler/review_handler_test.go`；`web/src/pages/ReviewSummaryPage.test.tsx`。 |
| 可重跑命令 | `go test ./internal/application/handler -run TestReview`；`npm --prefix web test -- --run web/src/pages/ReviewSummaryPage.test.tsx` |
| 不可声明事项 | 复盘只聚合本地事实，不做无来源推断，不生成交易动作。 |
| 残余缺口 | 缺 P23 原始 proposal/tasks/archive；更完整浏览器全路径由 P30/P39 补强。 |

### P24 本地运行硬化

| 项目 | 证据 |
| --- | --- |
| 当前交付边界 | 本地配置校验、server 启动前校验、SQLite 备份、安全恢复、恢复确认和 CLI smoke；恢复默认拒绝覆盖现有 DB。 |
| archive 状态 | 缺少逐阶段完整 archive 包；不得写成已 archive。 |
| 文档证据 | `docs/development-plan.md` P24 摘要；`docs/configuration.md` 配置校验和运维命令。 |
| 代码证据 | `cmd/agent/main.go`：`--validate-config`、`--backup`、`--restore`、`--restore-confirm`、`backupSQLite`、`restoreSQLite`；`internal/infrastructure/config/config.go`。 |
| 测试证据 | `cmd/agent/main_test.go` 覆盖 validate-config、backup、restore、安全恢复；`internal/infrastructure/config/config_test.go` 覆盖配置解析和校验。 |
| 可重跑命令 | `go test ./cmd/agent -run 'Test.*(Validate|Backup|Restore|Recovery)'`；`go test ./internal/infrastructure/config/...` |
| 不可声明事项 | 不提供云同步、多用户权限、复杂安装器、自动修复或自动覆盖真实库。 |
| 残余缺口 | 缺 P24 原始 proposal/tasks/archive；P40/P44/P49 已对本地部署、安装诊断和升级检查做后续补强。 |

## 4. 跨阶段补强关系

| 补强阶段 | 与 P19-P24 的关系 |
| --- | --- |
| P25 | 验证真实公开源，明确 P19/P20 只是基础 HTTP bridge 和 payload parser，不等于已接通所有真实外部源。 |
| P26/P29 | 补强公告/监管证据 collector、真实采集 smoke、错误分类和临时 SQLite 入库验收。 |
| P27/P34/P48 | 补强基金/ETF/指数市场数据 collector、数据源健康、质量回归和脱敏摘要。 |
| P30/P39 | 补强真实环境 E2E、Playwright smoke、空库到日报和前端全路径验收。 |
| P35/P36/P38 | 补强风险预警、规则效果验证、检索质量和守门人相关链路。 |
| P40/P44/P49 | 补强本地部署、恢复演练、安装诊断、发布升级检查和脱敏诊断汇总。 |

## 5. 发布前使用建议

1. P51 可以作为 P19-P24 历史审计入口，说明“当前有哪些事实证据可查”。
2. P51 不等于项目验收通过；P52 必须定义单元、集成、E2E、真实源、真实 LLM、冒烟、安装诊断、发布升级检查和安全边界门禁。
3. P53 发布候选材料不得只引用 P51；必须同时引用 P52 的门禁结果。
4. 若审计方要求逐阶段原始 proposal/tasks/archive，只能说明当前缺失，不能补写伪历史包。

## 6. 安全边界

P51 不新增运行时能力，也不改变任何安全边界。不得把本文解释为支持以下能力：券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复承诺、自动覆盖真实库、收益承诺、登录源、付费源、授权源、Level2 或高频源。
