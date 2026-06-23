# 文档治理规范

> 版本：v1.0  
> 最后更新：2026-06-22
> 配套：`openspec/project.md`、`docs/development-plan.md`

## 1. 目标

在 AI 辅助开发时，避免规格文档无序增殖、双真源漂移。所有**契约级**变更必须经过 OpenSpec 变更包；实现任务以变更包内的 `tasks.md` 为准，不以 `docs/superpowers/plans/` 为准。

## 2. 文档分级

| 级别 | 路径 | 性质 | 修改方式 |
| --- | --- | --- | --- |
| **L1 契约真源** | `docs/requirements.md`、`docs/data-model.md`、`docs/api.md`、`docs/workflow.md`、`docs/frontend-contract.md` | 系统行为与接口的权威定义 | 仅通过 OpenSpec change 的 delta，经审阅后合并 |
| **L2 架构与计划** | `docs/architecture.md`、`docs/functional-spec.md`、`docs/development-plan.md` | 架构说明、功能拆分、阶段计划 | 契约变更时同步更新；计划勾选可在实现完成后更新 |
| **L3 体验与图示** | `docs/ui-*.md`、`docs/ui/`、`docs/diagrams/` | UI 与图示 | 随相关 change 更新，或独立 UI change |
| **L4 变更脚手架** | `openspec/changes/<name>/` | 提案、delta、设计、任务 | 活跃变更期间自由编辑；完成后 **archive** |
| **L5 归档** | `openspec/changes/archive/` | 历史变更审计 | 只读，不修改 |
| **禁止作真源** | `docs/superpowers/plans/`、会话临时 md、未归档的 `design.md` | 探索/一次性材料 | 不得作为验收或契约依据 |

## 3. OpenSpec 与 `docs/` 的关系

- **契约真源**：`docs/`（L1–L3），不是 `openspec/specs/` 的全量副本。
- **变更入口**：`openspec/changes/<change-id>/`。
- **归档时**：将 delta 合并进对应的 `docs/*.md`，更新文档头中的版本/日期，再将 change 移入 `archive/`。
- **`openspec/specs/`**：仅存放从 `docs/` 抽取的**行为摘要**（可选、按域增量），或保持为空；**不以两套全文并存**。

## 4. 标准工作流

标准工作流适用于从 OpenSpec change 开始的新增变更。P19–P24 属于已交付但无 archive 包的历史状态校准；如需追补 proposal、delta、tasks 或验收记录，应作为独立治理 change 处理，不回写或伪造历史归档。

```text
1. /opsx:propose <change-id>   或  openspec new change "<id>"
2. 编写 proposal.md、specs/（delta）、design.md、tasks.md
3. 人工审阅 delta 与范围（in scope / out of scope）
4. /opsx:apply                  按 tasks.md 实现（可配合 Superpowers executing-plans）
5. 按 change 范围完成验证：OpenSpec 校验、后端测试、前端测试/构建、本地任务
6. archive 前执行只读子 agent 复审，且无 Critical / Important 问题
7. /opsx:archive                合并 delta → docs/，change 进入 archive/
```

### Superpowers 使用边界

| 技能 | 允许 | 禁止 |
| --- | --- | --- |
| `executing-plans` / `subagent-driven-development` | 执行 `openspec/changes/*/tasks.md` | — |
| `verification-before-completion` | 按 change 内验收命令验证 | — |
| `writing-plans` | — | 在 `docs/superpowers/plans/` 新建与 change 重复的计划 |
| `brainstorming` | 探索阶段，结论须回写到 change 或 ADR | 直接修改 L1 契约文件 |

## 5. Delta 格式（契约变更）

变更包内 `specs/<domain>/spec.md` 使用：

```markdown
## ADDED Requirements
### Requirement: ...
#### Scenario: ...

## MODIFIED Requirements
### Requirement: ...
（完整替换后的条文）

## REMOVED Requirements
### Requirement: ...
```

合并到 `docs/` 时，由人工或 Agent 在 archive 步骤把上述条文写入对应章节，并在文件头更新 `最后更新` 日期。

## 6. 变更命名

- 格式：`<阶段>-<简短描述>`，kebab-case。  
- 示例：`p0-engineering-skeleton`、`p1-sqlite-migrations`、`feat-gatekeeper-audit`。

## 7. 当前活跃变更

当前活跃变更：无。

P95 已归档到 `openspec/changes/archive/2026-06-23-p95-architecture-api-engineering-hardening/`。P95 已完成架构/API/工程加固：新增 `scripts/go-packages.sh` 并将 CI/release Go gates 限定到项目后端包，P93 source inventory 改为 tracked + nonignored untracked release-relevant source files，新增 API route contract check，SQLite 改为每条连接的 PRAGMA hook，支持 `DEEPSEEK_API_KEY_FILE`，修复 release manifest backend test 命令，并更新架构/部署/API 文档。GitHub CI commit `054d2708440da925e2e4cc4ae65065ac002b8905` run `https://github.com/wanggang8/Investment-agent/actions/runs/28001447118` 已通过 full backend `Go tests`；Security Scan run `https://github.com/wanggang8/Investment-agent/actions/runs/28001447137` 已通过。P95 不新增投资运行时能力，不声称物理第二机器复验、券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、收益承诺、登录源、付费源、授权源、Level2 或高频源。

P96 已归档到 `openspec/changes/archive/2026-06-23-p96-public-docs-readme-productization/`。P96 已完成 public README/docs 产品化：新增 root `README.md`、`docs/product-overview.md`、`docs/quickstart.md`，将 `docs/README.md` 收敛为文档地图，并把长阶段历史拆到 `docs/release/history.md`；验收记录见 `docs/release/acceptance/2026-06-23-p96-public-docs-readme-productization.md`。P96 不修改运行时代码，不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、收益承诺、登录源、付费源、授权源、Level2 或高频源。

P94 已归档到 `openspec/changes/archive/2026-06-23-p94-github-ci-release-hardening/`。P94 已完成 GitHub CI/CD hardening：PR/main CI 覆盖 OpenSpec、`go vet`、bounded `golangci-lint`、Go tests、frontend lint/test/build、P91/P92/P93 checks、release package smoke 和 whitespace check；tag `v*` release workflow 在打包前运行 release preflight 并上传 package/manifest；新增独立 security scan workflow，覆盖 `govulncheck`、frontend production dependency audit 和 P93 code reality / secret scan。P94 同步清理 frontend lint warning、Go staticcheck/unused findings，并更新 `docs/deployment.md` 与 P94 acceptance record。P94 不创建 Git tag，不发布 GitHub Release，不声称物理第二机器复验，不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、收益承诺、登录源、付费源、授权源、Level2 或高频源。

P93 已归档到 `openspec/changes/archive/2026-06-22-p93-final-code-reality-design-audit/`。P93 已完成代码真实性与设计合理性最终复核，生成 `docs/release/acceptance/2026-06-22-p93-final-code-reality-design-audit.md` 和 `scripts/p93_code_reality_audit.py`：P92 仍是 341 行原始需求逐项台账，P93 对 P92 341 行做 row-level cross-check，并确认每行 source section 可解析到当前生产代码/验收证据 bundle，release-blocking findings 为 0。P93 发现并移除未引用的 `web/src/pages/PlaceholderPage.tsx`、8 个历史 `internal/application/workflow/nodes/*` wrapper 和旧 helper，清空本地 `configs/config.local.yaml` 中的 key 占位并关闭 stub 默认，新增 bounded `sk-...` secret 扫描，确认生产 route 不接 placeholder/demo 页面，Docker/`.env.example` 默认 `use_stub=false` 且无嵌入 API key。P93 验证通过 `go test ./...`、`go vet ./...`、`npm --prefix web test`、`npm --prefix web run build`、`openspec validate --all --strict`、P92/P93 audit check、独立 secret scan、`git diff --check` 和 subagent 复审。P93 不新增运行时能力，不声称物理第二机器复验、券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、自动迁移、自动恢复、真实库覆盖、收益承诺、登录源、付费源、授权源、Level2 或高频源。

P92 已归档到 `openspec/changes/archive/2026-06-22-p92-final-original-requirement-audit-ledger/`。P92 已完成原始需求逐项最终独立复核台账，生成最终 ledger 与 summary：从 P88 全量 341 行矩阵叠加 P89/P90 最终补丁证据，确认 341 行原始需求中 330 个 full-release-required rows 全部 `real_pass`，11 个附录/参考 rows 为 `reference_only`，full-release-required 非 `real_pass` 为 0。P92 台账逐行列出功能入口/UI surface、预期行为或数据影响、API/SQLite/readback/审计证据、验收命令/证据文件和安全边界。P92 不新增运行时能力，不声称物理第二机器复验、券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、自动迁移、自动恢复、真实库覆盖、收益承诺、登录源、付费源、授权源、Level2 或高频源。

P91 已归档到 `openspec/changes/archive/2026-06-22-p91-github-release-docker-deployment/`。P91 已完成 GitHub Release 与 Docker Compose 一键部署路径：Dockerfile、Compose、`.env.example`、本地部署配置、install/upgrade/uninstall/backup/status/doctor 脚本、GitHub CI/release workflow、部署文档、deployment checker、真实 Docker first-install/upgrade/status/uninstall 验收和最终包刷新均通过。安装脚本会自动区分首次安装与升级；升级前备份；卸载默认保留 SQLite、VecLite、logs、backups 与 `.env`；只有 `--purge` 加精确确认短语才删除本地数据；默认端口只绑定 `127.0.0.1`。P91 不新增投资业务能力，不声称物理第二机器复验、券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、自动迁移、自动恢复、真实库覆盖、收益承诺、登录源、付费源、授权源、Level2 或高频源。

P90 已归档到 `openspec/changes/archive/2026-06-22-p90-capital-flow-provider-closure/`。P90 已用真实公开 Eastmoney H5 capital-flow provider、产品 Settings UI market refresh、market snapshot API 与 SQLite readback 完成 P89 后剩余 2 个 full-release-required rows（`REQ-04-016`、`REQ-05-003`）升级为 `real_pass`；证据见 `docs/release/acceptance/2026-06-22-p90-capital-flow-provider-closure.md`、`docs/release/acceptance/2026-06-22-p90-capital-flow-provider-matrix.md` 与 `docs/release/ui-audit-assets/2026-06-22-p90-capital-flow-provider/final-validation.log`。P90 层结论为 `release_ready_full_original_requirement_real_pass_candidate_with_p90_capital_flow_closure`，且 P89-chain full-release-required rows 已无已知非 `real_pass` 剩余。P90 不刷新 P76 package，不声称物理第二机器复验、远程发布、Git tag、full package refresh、券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、自动迁移、自动恢复、真实库覆盖、收益承诺、登录源、付费源、授权源、Level2 或高频源。

P88 已归档到 `openspec/changes/archive/2026-06-22-p88-remaining-full-release-blockers-closure/`。P88 处理 P86 后剩余 27 个 full-release-required blockers，fresh real browser UI/API/SQLite/Go acceptance 覆盖 source-verified `sell_only`、single-source `frozen_watch`、historical expected-return probabilities、sample<5 degradation、quarterly rebalance 和 SOP addendum proposal。P88 closure 矩阵结论为 `release_ready_scoped_with_p88_remaining_blocker_progress`：27 行中 17 行升级为 `real_pass`、10 行保留 `partial`；全量 341 rows 中仍有 10 个 full-release-required rows 非 `real_pass`。P88 不声称 full original-requirement pass，不刷新 P76 package；`REQ-04-016`、`REQ-05-003`、`REQ-05-004`、`REQ-05-005`、`REQ-08-004`、`REQ-08-023`、`REQ-09-004`、`REQ-09-023`、`REQ-09-024`、`REQ-09-025` 仍需后续真实 provider 或真实 UI/API/SQLite 动态概率/假设跟踪验收。P88 不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、自动迁移、自动恢复、真实库覆盖、收益承诺、登录源、付费源、授权源、Level2 或高频源。

P86 已归档到 `openspec/changes/archive/2026-06-22-p86-core-goal-knowledge-safety-final-closure/`。P86 fresh integrated runner 复跑 P74/P81/P82/P83/P84/P85/P87 真实 UI/API/SQLite/Go evidence，P86 矩阵结论为 `release_ready_scoped_with_p86_final_integrated_progress`：341 rows 中 303 `real_pass`、11 `reference_only`、27 `partial`；仍有 27 个 full-release-required rows 非 `real_pass`。P86 不声称 full original-requirement pass，不刷新 P76 package；剩余 blocker 主要是 source-verified buy-logic/frozen-watch workflow transition、资金流/两融/成分股财务真实字段、public collector production preverification、历史回测/动态概率、季度再平衡、SOP 增补提案和 sell-only transition 证据。

P86 证据：fresh integrated real browser UI/API/SQLite/workflow evidence 使用本地 Go backend、Vite frontend 和临时 SQLite，复跑 `p86-core-goal-knowledge-safety-final-acceptance.sh`，覆盖内置大师知识与 LLM context、`159915` accepted-local 动态源、SOP/action/failure-state、治理追溯、组合/确认数据影响、预期收益读回、仓位状态/配置安全；`python3 scripts/p86_core_goal_knowledge_safety_final_closure.py --check` 通过，新增 110 行 `real_pass`，27 行因证据口径不足保留 `partial`。P86 不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、自动迁移、自动恢复、真实库覆盖、收益承诺、登录源、付费源、授权源、Level2 或高频源。

P87 已归档到 `openspec/changes/archive/2026-06-22-p87-portfolio-state-allocation-safety-closure/`：评估 32 行，5 行升级为 `real_pass`，27 行因源验证转移、季度再平衡、提案、source-readiness、审计或 release-safety 口径不足 deferred。P87 当前层结论为 `release_ready_scoped_with_p87_portfolio_state_allocation_progress`；不得扩展为 P86/full original-requirement pass。

P87 证据：fresh portfolio state/allocation safety real browser UI/API/SQLite/Go evidence 使用本地 Go backend、Vite frontend 和临时 SQLite，通过 `/positions` 写入/读回核心、卫星、现金/货币基金、`buy_date`、`normal`/`sell_only`/`frozen_watch`，验证 cash_ratio=8%、cash+money-fund bucket=9%、core=64%、satellite=27%、focused handler/rule tests、forbidden broker/order/push/auto-confirm absence 和只读 `--check` 校验。P87 当前层结论为 `release_ready_scoped_with_p87_portfolio_state_allocation_progress`：341 rows 中 193 `real_pass`、11 `reference_only`、137 `partial`；仍有 137 个 full-release-required rows 非 `real_pass`。P87 不刷新 P76 package，不声称 source-verified buy-logic-break/frozen-watch workflow transition、完整季度再平衡、提案确认/拒绝、月度归因、完整审计史、public collector production readiness、full release/upgrade preflight closure、full original-requirement pass、券商接口、自动交易、外部推送、自动确认或自动规则应用。

P85 已归档到 `openspec/changes/archive/2026-06-22-p85-expected-return-analysis-accuracy-closure/`。Fresh expected-return analysis-accuracy real browser UI/API/SQLite/Go evidence 使用本地 Go backend、Vite frontend 和临时 SQLite，通过 `/consultation` 目标收益率与上一轮基准情景中枢 UI 输入、完整样本/下行情景/样本不足三类咨询、决策详情 readback、SQLite 字段级 readback、focused workflow/handler tests、forbidden broker/order/push/auto-confirm absence 和只读 `--check` 校验。P85 评估 31 行，新增 15 行 `real_pass`，16 行因历史准确性/回测/概率口径不足 deferred；当前 P85 层结论为 `release_ready_scoped_with_p85_expected_return_progress`：341 rows 中 188 `real_pass`、11 `reference_only`、1 `scoped_pass`、141 `partial`；仍有 142 个 full-release-required rows 非 `real_pass`。本环境 `DEEPSEEK_API_KEY` 不存在，P85 不声称 fresh real LLM output，只声明确定性 workflow + 真实 UI/API/SQLite readback 证据。P85 不刷新 P76 package，不声称 full original-requirement pass，不新增未来收益准确性、未来市场方向准确性、真实历史回测模型、自动概率下调、纵向假设跟踪、券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、自动迁移、自动恢复、真实库覆盖、收益承诺、登录源、付费源、授权源、Level2 或高频源。

P84 已完成归档前证据与子 agent 审查。Fresh portfolio/confirmation real browser UI/API/SQLite/Go evidence 使用本地 Go backend、Vite frontend 和临时 SQLite，通过 `/positions` 本地账户校准、持仓编辑、批量导入、线下交易、修正审计、`/decisions/decision_p84_pending` 手动确认、`/decision-loop`、`/review`、`/audit`、`/workbench` 下游读回、SQLite 字段级 readback、forbidden broker/order/push table absence、auto-confirmation absence、focused handler tests 和只读 `--check` 校验。P84 评估 35 行，新增 3 行 `real_pass`，32 行因证据口径不足 deferred；当前 P84 层结论为 `release_ready_scoped_with_p84_portfolio_confirmation_progress`：341 rows 中 173 `real_pass`、11 `reference_only`、1 `scoped_pass`、156 `partial`；仍有 157 个 full-release-required rows 非 `real_pass`。P84 不升级宽泛月度归因/规则提案应用时间行，不刷新 P76 package，不声称 full original-requirement pass，不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、自动迁移、自动恢复、真实库覆盖、收益承诺、登录源、付费源、授权源、Level2 或高频源。

P83 已归档到 `openspec/changes/archive/2026-06-22-p83-governance-traceability-backfill/`。Fresh governance traceability real browser UI/API/SQLite/Go evidence 使用本地 Go backend、Vite frontend 和临时 SQLite，通过 `/review`、`/rules`、`/audit`、`/notifications`、`/daily-discipline/reports`、`/local-install` 真实浏览器 UI/readback、monthly/quarterly review API、notification mark-read、local-install redaction、SQLite 字段级 readback、forbidden table absence、focused Go tests 和只读 `--check` 校验。P83 评估 43 行，新增 10 行 `real_pass`，33 行因证据口径过宽 deferred 到 P86；当前 P83 层结论为 `release_ready_scoped_with_p83_governance_traceability_progress`：341 rows 中 170 `real_pass`、11 `reference_only`、3 `scoped_pass`、157 `partial`；仍有 160 个 full-release-required rows 非 `real_pass`。P83 不刷新 P76 package，不声称 full original-requirement pass，不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、自动迁移、自动恢复、真实库覆盖、收益承诺、登录源、付费源、授权源、Level2 或高频源。

P82 已归档到 `openspec/changes/archive/2026-06-22-p82-sop-action-ui-sqlite-closure/`。P82 fresh SOP/action real browser UI journey 使用本地 Go backend、Vite frontend 和临时 SQLite，通过 SOP A-F、failure states、mark-error、gatekeeper、显式最终确认应用本地规则版本、SQLite/readback、forbidden table absence 和 `--check` 只读校验。P82 评估 53 行，新增 44 行 `real_pass`，9 行因证据口径过宽被明确 deferred；当前 P82 层结论为 `release_ready_scoped_with_p82_sop_action_progress`：341 rows 中 160 `real_pass`、11 `reference_only`、3 `scoped_pass`、167 `partial`；仍有 170 个 full-release-required rows 非 `real_pass`。P82 不刷新 P76 package，不声称 full original-requirement pass，不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、自动迁移、自动恢复、真实库覆盖、收益承诺、登录源、付费源、授权源、Level2 或高频源。

P81 已归档到 `openspec/changes/archive/2026-06-22-p81-dynamic-source-field-coverage/`。P81 已完成动态源字段覆盖真实验收，fresh `159915` accepted-local 非 `510300` 真实浏览器 UI journey、Go collector/readiness 测试、SQLite/readback、source-health provenance、formal evidence、RAG indexing、真实 LLM-backed UI readback 和 forbidden table absence 均通过。P81 新增 59 行 `real_pass`；当前 P81 层结论为 `release_ready_scoped_with_p81_dynamic_source_progress`：341 rows 中 116 `real_pass`、11 `reference_only`、4 `scoped_pass`、210 `partial`；仍有 214 个 full-release-required rows 非 `real_pass`。P81 不刷新 P76 package，不声称 full original-requirement pass，不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、自动迁移、自动恢复、真实库覆盖、收益承诺、登录源、付费源、授权源、Level2 或高频源。

P80 已归档到 `openspec/changes/archive/2026-06-22-p80-review-audit-governance-real-use-closure/`。P80 已完成 P79 后复盘、审计、错误标注、规则提案确认与守门人审计真实使用闭环：fresh P75 SOP/failure-state 真实浏览器 UI journey 通过，字段级 SQLite/readback 摘要覆盖 `error_cases`、`operation_confirmations`、`rule_proposals`、`gatekeeper_audits`、`audit_events`、`risk_alerts` 与 forbidden table absence；当前 P80 层结论为 `release_ready_scoped_with_p80_review_audit_governance_progress`，341 rows 中 57 `real_pass`、11 `reference_only`、5 `scoped_pass`、268 `partial`，仍有 273 个 full-release-required rows 非 `real_pass`。P80 不升级月度/季度归因、完整提案生成、最终规则应用时间或 full original-requirement pass，不刷新 P76 package，不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、自动迁移、自动恢复、真实库覆盖、收益承诺、登录源、付费源、授权源、Level2 或高频源。

P79 已归档到 `openspec/changes/archive/2026-06-21-p79-real-use-data-impact-and-expected-return-closure/`。P79 已完成 P78 后真实使用数据影响与预期收益闭环：fresh P72 `510300` 真实用户 UI 场景和 P75 accepted-local 非 `510300` 真实 UI 场景均通过，字段级 SQLite/readback 门禁覆盖持仓、确认、交易前后状态、证据引用和安全负证据；ExpectedReturnNode 在 LLM `quality_failure` 时丢弃失败材料并使用 deterministic-local 安全兜底，普通 analyst unavailable 仍降级。当前 P79 层结论为 `release_ready_scoped_with_p79_real_use_data_impact_progress`：341 rows 中 43 `real_pass`、11 `reference_only`、22 `scoped_pass`、265 `partial`，仍有 287 个 full-release-required rows 非 `real_pass`。P79 不升级 `REQ-11-018`、`REQ-14-005`、`REQ-14-007` 等宽泛审计/月度归因行，不刷新 P76 package，不声称 full original-requirement pass，不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、自动迁移、自动恢复、真实库覆盖、收益承诺、登录源、付费源、授权源、Level2 或高频源。

P78 已归档到 `openspec/changes/archive/2026-06-21-p78-requirements-real-pass-batch-closure/`。P78 已完成 P77 后原始需求 `real_pass` 批次化收敛第一批：当前 P78 层结论为 `release_ready_scoped_with_p78_real_pass_batch_progress`，341 rows 中 20 `real_pass`、11 `reference_only`、22 `scoped_pass`、288 `partial`，仍有 310 个 full-release-required rows 非 `real_pass`。P78 不重写 P75/P77 历史矩阵，不刷新 P76 package，不声称 full original-requirement pass，不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、自动迁移、自动恢复、真实库覆盖、收益承诺、登录源、付费源、授权源、Level2 或高频源。

P77 已归档到 `openspec/changes/archive/2026-06-21-p77-requirements-real-pass-upgrade-gate/`。P77 已建立 P75 后原子需求 `real_pass` 升级门禁与第一批升级证据：生成新的 P77 evidence layer，不重写 P75 历史矩阵；当前结论为 `release_ready_scoped_with_p77_real_pass_progress`，341 rows 中 17 `real_pass`、11 `reference_only`、22 `scoped_pass`、291 `partial`，仍有 313 个 full-release-required rows 非 `real_pass`。P77 不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、自动迁移、自动恢复、真实库覆盖、收益承诺、登录源、付费源、授权源、Level2 或高频源；不得把 scoped/partial/deterministic-local、fixture/mock/stub、截图、route smoke、scope exclusion、waiver、临时 DB 或单标的证据冒充 full original-requirement pass。

P76 已归档到 `openspec/changes/archive/2026-06-21-p76-post-p75-final-package-refresh/`。P76 完成 P75 后最终本地分发包刷新：已从 clean source commit `8a317f25917b8ff18ec9b5049e6a6188206a22d3` 生成 `p76-post-p75-final` package，package verify 与 cross-machine-equivalent local repeat acceptance 均通过，包内确认包含 P72-P75 acceptance Markdown 与 OpenSpec archives。P76 不新增运行时能力，不扩大 P75 `release_ready_scoped_with_traceability_gaps` 结论，不声称物理第二机器复验、远程发布、Git tag、自动升级、自动迁移、自动恢复、自动修复、真实库覆盖、券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、未来 provider 可用性或收益承诺。

P75 已归档到 `openspec/changes/archive/2026-06-21-p75-requirements-traceability-and-real-use-closure/`。P75 已执行原始需求追踪与真实使用闭环审计：把 `docs/requirements.md` sections 1-19 拆成 341 个原子级 requirement rows，并逐项追踪实现、UI、数据、工作流、LLM、场景、SQLite 数据影响、审计、安全边界和发布声明证据。当前结论为 `release_ready_scoped_with_traceability_gaps`，不是 `release_ready_full_requirements_traceable`：当前矩阵为 291 `partial`、33 `scoped_pass`、17 `deterministic_local_evidence`、0 `real_pass`、0 `blocked`、0 `not_implemented`；expanded G9 forbidden-term scan 已完成分类人审，未发现 forbidden runtime affordance。P75 必须继续区分 `real_pass`、`scoped_pass`、`deterministic_local_evidence`、`partial`、`not_implemented` 和 `blocked`；不得把 P71-P74 scoped evidence、`510300` 单路径、fixture/mock/stub、scope exclusion、waiver、截图或 route smoke 冒充 full product pass。P75 不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、自动迁移、自动恢复、真实库覆盖、收益承诺、登录源、付费源、授权源、Level2 或高频源。

P74 已完成，待/已归档到 `openspec/changes/archive/2026-06-19-p74-built-in-knowledge-and-data-readiness/`。P74 建立内置知识与数据准备度闭环，覆盖 7 位大师经验、纪律、SOP、标的画像注册表、ETF/基金/指数数据依赖矩阵、采集数据准备度 API/UI、LLM 上下文引用和真实/降级场景验收；结论为 `release_ready_built_in_knowledge_data_readiness`。P74 不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、自动迁移、自动恢复、真实库覆盖、收益承诺、登录源、付费源、授权源、Level2 或高频源。

P71 已归档到 `openspec/changes/archive/2026-06-18-p71-real-product-acceptance-true-pass/`；当前状态：`release_ready_full_real_product_acceptance`。P71 已让当前本地 `000300` strict current-data gate 真正返回 `policy=passed` / `gate=pass`，并已加固 VecLite/RAG rebuild；`bash scripts/p71-real-product-acceptance.sh` 已在真实本地 backend、Vite frontend、真实 LLM provider 和 `use_stub=false` 配置下通过，真实 UI consultation 结果为 `workflow_status=completed`、3 份 LLM analyst report 全部 `parsed`/`passed`，retrieval 为 `fallback_source=veclite`、`index_health=healthy`、`index_freshness=fresh`。post-P70/P71 package refresh 已从 clean source commit `2c195a05cee3b6cdda031e86409d562bcc7ee379` 生成，verify 和 isolated repeat acceptance 均通过。不得把该结论扩展为未来 public-source/LLM provider 可用性承诺、物理第二机器复验、券商接口、自动交易、自动修复或收益承诺。

P72 已归档到 `openspec/changes/archive/2026-06-19-p72-real-user-fund-scenario-data-impact-acceptance/`。P72 已将真实用户基金/ETF `510300` 场景作为验收对象，完成 UI 操作、正式公开证据采集、真实 LLM 咨询、人工确认、SQLite 数据影响、审计事件、衍生页面回显、确定性计算准确性、真实 LLM 分析链路和安全边界验收；结论为 `release_ready_full_real_user_scenario_acceptance`。P72 不新增产品运行时能力，不承诺未来收益、未来市场方向或预测准确性，不接券商、不自动交易、不外部推送。

P67 已归档到 `openspec/changes/archive/2026-06-18-p67-current-data-gate-resolution-workflow/`；已处理 P66 current-data policy gate block 的本地人工处置、豁免或范围排除工作流，P67 当时 000300 为 `resolved_with_scope_exclusion`，但不得把 P66 `policy=blocked` 或 P67 scope exclusion 伪装为 clean pass。P71 的 current-data clean claim 来自 fresh strict `policy=passed` / `gate=pass` 证据。

P68 已归档到 `openspec/changes/archive/2026-06-18-p68-post-p67-release-readiness-governance/`；已刷新 P67 后 release-ready 表述、发布候选材料和打包复验策略，当前 release 状态为 `release_ready_limited_current_data_scope`；当时建议的 P69 clean tree 分发包刷新已完成。

P69 已归档到 `openspec/changes/archive/2026-06-18-p69-clean-tree-package-refresh/`；已从 clean detached worktree 生成 `p69-clean-tree` package，`source_status=clean`，verify 和 repeat acceptance 均通过；当前 package evidence 覆盖 P65-P68 后提交，但不声称包内包含 P69 文档。

P70 已归档到 `openspec/changes/archive/2026-06-18-p70-final-release-decision-and-risk-closure/`；已完成最终发布决策与风险收口，最终状态为 `release_ready_limited_current_data_scope`，确认 limited local release scope 无必需下一阶段；不新增运行时能力，不改变 P66/P67 current-data 边界，不声称 P69 package 包含 P69/P70 文档。

P73 已归档到 `openspec/changes/archive/2026-06-19-p73-product-effectiveness-ux-validation/`。P73 已完成产品目标效果与 UX 验收，`bash scripts/p73-product-effectiveness-ux-validation.sh` 已在真实本地 Go backend、Vite frontend、临时 SQLite 和真实浏览器 UI 操作下通过；browser results、截图和 effect replay summary 已写入 `docs/release/ui-audit-assets/2026-06-19-p73/`。验收覆盖纪律执行、证据充分性、可追溯性、复盘有效性、真实 UX 任务理解、背景材料阻断、人工确认不变更持仓、风险/readback/规则效果联动、390px reflow 和 unsafe input 阻断；当前结论为 `release_ready_product_effectiveness_ux_acceptance`。P73 不以未来收益率作为通过标准，不新增券商接口、自动交易或收益承诺。

下一执行队列：无必需下一阶段。物理第二机器复验仍未执行，另行刷新包含 P81-P90 材料的最终分发包也需单独立项。

最近完成：P54 `p54-release-handoff-and-repeatability-hardening` 已归档到 `openspec/changes/archive/2026-06-17-p54-release-handoff-and-repeatability-hardening/`；已新增发布交付说明和验收可重复性规则，承接 P53 `release_ready` 但不扩大声明。

P55 已归档到 `openspec/changes/archive/2026-06-17-p55-full-ui-acceptance-and-design-audit/`；范围为真实启动本地项目，通过浏览器操作 UI 验收主要功能，并使用 Product Design audit 审查前端设计是否需要优化；未修改运行时代码。P55 发现 full UI acceptance blocker：真实 LLM-backed consultation 生成的决策详情可因 nullable `final_verdict.optional_actions` 崩溃。

P56 已归档到 `openspec/changes/archive/2026-06-17-p56-ui-acceptance-blocker-fixes/`；已修复 P55-B1，完成任务分组导航、产品化 UI 基础层、`/positions` 与 `/data-quality` 移动端 reflow，并通过真实 LLM UI 验收和 Playwright E2E。当前 full UI acceptance 状态为 `p56_scope_pass`，但不等同新的最终 release_ready 声明。

P57 已归档到 `openspec/changes/archive/2026-06-17-p57-product-experience-polish-roadmap/`；已新增 `docs/product-experience-polish-roadmap.md`，固化产品设计、UI 设计和功能设计打磨路线图，拆分 P58-P63 后续阶段。发布状态刷新后移到产品体验打磨完成或明确豁免后执行。

P58 已归档到 `openspec/changes/archive/2026-06-17-p58-daily-workbench-redesign/`；已完成 Dashboard / Workbench 每日投资纪律 cockpit 重构、共享 daily workbench view model、首屏今日状态、下一步人工动作、信号摘要、桌面/390px 截图和真实本地 UI 验收。

P59 已归档到 `openspec/changes/archive/2026-06-17-p59-decision-explainability-experience/`；已重构 Consultation、Decision Detail、Evidence 和 Decision Loop 的解释链路，把主动咨询、最终裁决、证据、LLM 分析、规则链、审计与闭环记录串成可理解的决策故事，并完成真实 UI 验收、nullable DTO 加固、桌面/390px 截图和安全降级记录。

P60 已归档到 `openspec/changes/archive/2026-06-18-p60-portfolio-risk-data-quality-experience/`；已完成 Positions、Risk Alerts、Data Quality 的维护、处置和质量可观测体验重构，并通过真实本地 UI 验收、桌面/390px 截图、E2E smoke 和安全边界扫描。

P61 已归档到 `openspec/changes/archive/2026-06-18-p61-governance-ops-productization/`；已完成 Rules、Audit、Notifications、Daily Reports、Daily Auto Run、Local Install、Local Knowledge、Settings 的治理和运维页面产品化。

P62 已归档到 `openspec/changes/archive/2026-06-18-p62-design-system-accessibility-hardening/`；已完成设计系统 primitives、状态 tone、键盘路径、可访问语义、390px/768px/1280px reflow 和视觉回归门禁。P62 不新增后端 API、SQLite schema、Eino workflow、LLM 能力、券商接口、交易能力、外部推送、自动确认、自动规则应用、自动修复或发布状态刷新。

P63 已归档到 `openspec/changes/archive/2026-06-18-p63-full-ui-regression-release-refresh/`；范围为全量真实 UI 回归与发布状态刷新。已真实启动后端和前端，覆盖真实 LLM consultation、全路由浏览器操作、移动端/桌面端/错误态/降级态/安全边界扫描，并更新 `docs/release/acceptance/2026-06-18-p63-full-ui-regression.md`、`docs/release/release-candidate-2026-06-18.md`、`docs/release/release-handoff-2026-06-18.md`；当前 release 状态为 `release_ready`。

P64 已归档到 `openspec/changes/archive/2026-06-18-p64-release-packaging-version-tagging/`；已实现本地发布包脚本、sidecar manifest、archive checksum、verify 模式、输出目录限制、tracked/allowlisted source staging、敏感信息/禁止路径扫描和 P64 packaging handoff。P64 未新增后端业务能力、SQLite schema、HTTP API、Eino workflow、LLM 能力、券商接口、交易能力、外部推送、自动确认、自动规则应用、自动修复、自动升级、自动迁移、自动恢复或发布到远程渠道。

P65 已归档到 `openspec/changes/archive/2026-06-18-p65-cross-machine-release-repeat-acceptance/`；已使用 P64 package workflow 生成 P65 candidate archive，并在本地跨机器等价隔离环境中完成 package verify、安装、OpenSpec/Go/frontend/E2E smoke 复验和发布材料更新。P65 未声称物理第二机器已执行，仍不得新增远程发布、Git tag、自动升级、自动迁移、自动恢复、自动修复、券商接口、交易、外推或收益承诺。

P66 已归档到 `openspec/changes/archive/2026-06-18-p66-current-data-zero-degradation-policy/`；已把 P63/P65 release caveat 中的 current data-source degraded 状态转化为明确 policy verdict、阻断/豁免规则、可重复验收和发布声明边界。P66 当时本地库 strict current-data gate 结论为 `policy=blocked` / `gate=block`，后续 release-ready 声明必须通过 P66 gate 或显式记录豁免/范围排除；不得新增付费/授权/登录源、Level2、高频源、真实 provider 调用、自动刷新、自动修复、券商接口、自动交易或未来数据源可用性承诺。P71 已在当前验收窗口提供 fresh strict pass 证据。

后续新阶段必须先创建独立 OpenSpec change，并继续保留安全边界：不得新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复承诺、自动覆盖真实库、收益承诺、登录源、付费源、授权源、Level2 或高频源。
