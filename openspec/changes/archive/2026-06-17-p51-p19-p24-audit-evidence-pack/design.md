# P51 设计：P19-P24 审计证据包

## 设计目标

P51 不是补写历史，也不是重新实现 P19-P24。它把当前仓库中可验证的事实整理为审计证据包，让后续 P52 验收矩阵和 P53 发布材料有可引用的边界基础。

## 证据包结构

新增 `docs/p19-p24-audit-evidence-pack.md`，建议结构如下：

1. **总览**
   - 更新时间。
   - 适用范围。
   - 核心结论：P19-P24 能力已交付但缺逐阶段完整 archive 包；本文件不是历史 archive。
2. **P14-P18 旁证核查**
   - 说明 P14-P18 存在标准 archive 包。
   - 如果用户表达为“P19-P14”，P14-P18 不需要补证据包，只需保留存在性核查。
3. **P19-P24 分阶段矩阵**
   - 阶段目标。
   - 当前交付边界。
   - 文档证据。
   - 代码证据。
   - 测试证据。
   - 可重跑命令。
   - 不可声明事项。
   - 残余缺口。
4. **跨阶段证据**
   - P25-P29 对真实公开源的后续补强。
   - P30/P39 对 E2E 的后续补强。
   - P40/P44/P49 对本地运维、安装、升级的后续补强。
5. **发布前使用建议**
   - P51 可作为历史审计入口。
   - P52 必须把证据包转化为验收门禁矩阵。
   - P53 不应仅凭 P51 宣称发布就绪。

## 每阶段证据映射

### P19 公开 HTTP 数据桥接

- 文档证据：`docs/development-plan.md` P19 摘要、`docs/configuration.md` 数据源配置。
- 代码证据：`internal/application/workflow/data_sources.go` 的 `ConfiguredMarketDataSource`、`ConfiguredIntelligenceSource`。
- 测试证据：`internal/application/workflow/data_sources_test.go` 的 public HTTP、fixture、fallback、失败分类测试；`internal/application/handler/market_handler_test.go` 的市场刷新通知/降级测试。
- 不可声明：不等于接通所有真实公开源；真实 collector 由 P25-P29 补强。

### P20 A 股 ETF/基金证据 payload

- 文档证据：`docs/development-plan.md` P20 摘要、`docs/data-model.md` 证据事实表、`docs/api.md` 本地知识/证据写入要求。
- 代码证据：`internal/application/workflow/evidence_verification_graph.go`、`internal/application/workflow/evidence_ingestion.go`、`internal/infrastructure/persistence/sqlite/intelligence_repo_impl.go`。
- 测试证据：`internal/application/workflow/workflow_integration_test.go` 的 `source_verifications` 写入；`internal/infrastructure/persistence/sqlite/intelligence_repo_impl_test.go`。
- 不可声明：不等于完成真实公告源历史补采；真实公告 collector 由 P26/P29 补强。

### P21 应用内通知中心

- 文档证据：`docs/data-model.md` `notifications` 表、`docs/api.md` 通知接口、`docs/frontend-contract.md` 通知页契约。
- 代码证据：`internal/application/service/notification.go`、`internal/application/handler/notification_handler.go`、`internal/infrastructure/persistence/sqlite/notification_repo_impl.go`、`web/src/services/notification.ts`。
- 测试证据：`internal/application/handler/notification_handler_test.go`、`internal/infrastructure/persistence/sqlite/notification_repo_impl_test.go`、`web/src/services/notification.test.ts`。
- 不可声明：不提供邮件、短信、系统 Push、Webhook、WebSocket 或外部推送。

### P22 规则体系与提案增强

- 文档证据：`docs/development-plan.md` P22 摘要、`docs/data-model.md` 规则提案事务边界、`docs/api.md` 规则提案说明。
- 代码证据：`internal/application/workflow/evolution_proposal_graph.go`、`internal/domain/rule/gatekeeper_logic.go`、`internal/application/service/rule_proposal.go`、`web/src/components/rules/RuleProposalPanel.tsx`。
- 测试证据：`internal/application/workflow/workflow_integration_test.go` 的 proposal 状态/通知/审计测试、`web/src/components/rules/RuleProposalPanel.test.tsx`。
- 不可声明：规则不会自动应用，仍需用户确认和守门人审计。

### P23 复盘深化

- 文档证据：`docs/development-plan.md` P23 摘要、`docs/api.md` 复盘字段、`docs/frontend-contract.md` review summary。
- 代码证据：`internal/application/dto/review.go`、`internal/application/handler/review_handler.go`、`web/src/components/review/ReviewSummaryPanel.tsx`。
- 测试证据：`internal/application/handler/review_handler_test.go`、`web/src/pages/ReviewSummaryPage.test.tsx`。
- 不可声明：复盘只聚合本地事实，不做无来源结论。

### P24 本地运行硬化

- 文档证据：`docs/development-plan.md` P24 摘要、`docs/configuration.md` 配置校验和运维命令。
- 代码证据：`cmd/agent/main.go` 的 `--validate-config`、`--backup`、`--restore`、`--restore-confirm`、`backupSQLite`、`restoreSQLite`；`internal/infrastructure/config/config.go`。
- 测试证据：`cmd/agent/main_test.go` 的 validate-config、backup、restore、安全恢复测试；`internal/infrastructure/config/config_test.go`。
- 不可声明：不提供云同步、多用户权限、复杂安装器或自动修复。

## 验收策略

P51 为文档审计变更，验收以一致性和可追溯性为主：

- OpenSpec 当前 change 严格校验通过。
- OpenSpec 全量严格校验通过。
- `git diff --check` 通过。
- 子 agent 计划复审和执行后复审均无 Critical / Important。
- 文档只列仓库内可核验证据，不引用未提交截图或外部临时结果。

## 安全边界

证据包不得新增或暗示以下能力：券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复承诺、自动覆盖真实库、收益承诺、登录源、付费源、授权源、Level2 或高频源。
