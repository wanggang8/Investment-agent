# Investment Agent 文档地图

本文档是项目文档入口，用于降低后续规划和维护时的查找成本。当前仓库仍以 `docs/` 为契约真源，OpenSpec 负责变更提案、delta、任务和归档审计。

## 当前阶段状态

- P0–P18：工程骨架、规则、工作流、API、前端、测试、治理与本地运维基础已完成。
- P19–P24：可用 MVP 路径已完成，包括公开 HTTP 数据桥接、A 股 ETF/基金证据 payload 解析、应用内通知、规则提案增强、复盘深化、配置校验与本地备份恢复。仓库未补建对应逐阶段完整 `openspec/changes/archive/` 包；P51 已新增 `docs/p19-p24-audit-evidence-pack.md` 作为当前事实审计证据包。
- P25：已完成真实公开权威数据源验证，明确巨潮资讯、深交所、证监会进入 P26 首批公告/监管证据源；中证指数、东方财富基金进入 P27 候选范围。P19/P20 不应被解读为已经接通所有真实外部源。
- P26：已归档到 `openspec/changes/archive/2026-06-05-p26-public-evidence-collectors/`；包含首批只读公告/监管 collector、默认 90 天/分页边界、标准证据 payload、RAG 入库、source verification、审计、幂等去重和 `cmd/agent --task public-evidence-refresh` 手动触发入口。
- P27：已归档到 `openspec/changes/archive/2026-06-06-p27-fund-etf-market-data-collectors/`；包含东方财富基金 B 级净值/历史净值/资产配置/基金档案辅助源，以及已校准到当前公开 `index-basic-info` shape 的中证指数 A 级指数基础资料 collector；中证指数样本/权重/估值文件扩展接口仍按候选 metadata 低频读取。
- P28：已归档到 `openspec/changes/archive/2026-06-06-p28-expected-return-dynamic-sell/`；包含样本门槛、样本上下文、情景解释、非交易式动态卖出评估提示、历史 JSON 回放兼容和前端展示。
- P29：已归档到 `openspec/changes/archive/2026-06-06-p29-public-evidence-collector-smoke/`；包含公开证据真实采集 smoke、显式日期窗口、错误分类、全源 `no_data` 成功空刷新、CNInfo orgId 配置映射和 OpenSpec 严格校验通过证据。
- P30：已归档到 `openspec/changes/archive/2026-06-07-p30-real-e2e-smoke/`；包含本地真实环境 E2E / Playwright smoke 验收，使用临时 SQLite/配置和可控 seed 数据验证后端健康检查、前端决策详情 expected return、证据页、审计页和临时产物治理。
- P31-P50：已完成每日自动运行闭环、每日纪律报告产品化、P33-P40/P41/P45/P50 路线图治理、账户/持仓体验、真实数据覆盖、风险预警、规则效果验证、真实 LLM 质量评估、RAG 检索质量、全路径 E2E、本地部署恢复、本地知识导入、决策闭环解释、数据源质量回归和本地发布升级体验。
- P51：已归档到 `openspec/changes/archive/2026-06-17-p51-p19-p24-audit-evidence-pack/`；范围是整理 P19-P24 当前事实审计证据包，不伪造历史 archive，不修改运行时代码。
- P52：已归档到 `openspec/changes/archive/2026-06-17-p52-project-acceptance-gate-matrix/`；已新增 `docs/project-acceptance-gate-matrix.md`，范围是建立单元、集成、E2E、真实源、真实 LLM、冒烟和发布前门禁矩阵，不宣称验收已通过。
- P53：已归档到 `openspec/changes/archive/2026-06-17-p53-acceptance-execution-and-release-candidate-materials/`；已执行 P52 G0-G9 验收并生成 `docs/release/acceptance/2026-06-17-p53-acceptance-run.md`、`docs/release/release-candidate-2026-06-17.md`，当前结论为 `release_ready`，含 G5 current degraded 与重试记录。
- P54：已归档到 `openspec/changes/archive/2026-06-17-p54-release-handoff-and-repeatability-hardening/`；已新增 `docs/release/README.md`、`docs/release/release-handoff-2026-06-17.md`、`docs/release/acceptance-repeatability.md`，用于交付说明与验收可重复性规则。
- P55：已归档到 `openspec/changes/archive/2026-06-17-p55-full-ui-acceptance-and-design-audit/`；已真实启动项目并通过前端 UI 操作执行全功能验收和 Product Design audit。P55 当时 full UI acceptance 结论为 `blocked`，原因是 real LLM-backed consultation 生成的决策详情可因 nullable `final_verdict.optional_actions` 崩溃；详见 `docs/release/ui-acceptance-2026-06-17.md` 与 `docs/release/ui-design-audit-2026-06-17.md`。
- P56：已归档到 `openspec/changes/archive/2026-06-17-p56-ui-acceptance-blocker-fixes/`；已修复 P55-B1，完成任务分组导航、产品化 UI 基础层、`/positions` 与 `/data-quality` 移动端 reflow，并通过真实 UI 浏览器验收。当前 full UI acceptance 状态为 `p56_scope_pass`。
- P57：已归档到 `openspec/changes/archive/2026-06-17-p57-product-experience-polish-roadmap/`；已新增 `docs/product-experience-polish-roadmap.md`，固化产品设计、UI 设计和功能设计打磨路线图，拆分 P58-P63 后续阶段。发布状态刷新后移到产品体验打磨完成或明确豁免后执行。
- P58：已归档到 `openspec/changes/archive/2026-06-17-p58-daily-workbench-redesign/`；已完成 Dashboard / Workbench 每日投资纪律 cockpit 重构和真实本地 UI 验收。
- P59：已归档到 `openspec/changes/archive/2026-06-17-p59-decision-explainability-experience/`；已重构 Consultation、Decision Detail、Evidence 和 Decision Loop 的解释链路，完成真实 UI 验收、nullable DTO 加固、桌面/390px 截图和安全降级记录。
- P60：已归档到 `openspec/changes/archive/2026-06-18-p60-portfolio-risk-data-quality-experience/`；已完成 Positions、Risk Alerts、Data Quality 的维护、处置和质量可观测体验重构，并完成真实本地 UI 验收。
- P61：已归档到 `openspec/changes/archive/2026-06-18-p61-governance-ops-productization/`；已完成 Rules、Audit、Notifications、Daily Reports、Daily Auto Run、Local Install、Local Knowledge、Settings 的治理和运维页面产品化体验。
- P62：已归档到 `openspec/changes/archive/2026-06-18-p62-design-system-accessibility-hardening/`；已固化设计系统 primitives、状态体系、键盘路径、可访问语义、390px/768px/1280px reflow 和视觉回归门禁，验收记录见 `docs/release/acceptance/2026-06-18-p62-ui-acceptance.md`。
- P63：已归档到 `openspec/changes/archive/2026-06-18-p63-full-ui-regression-release-refresh/`；已执行 P52 G0-G9、全路由真实 UI 回归、真实 LLM-backed consultation UI journey，并刷新 `docs/release/acceptance/2026-06-18-p63-full-ui-regression.md`、`docs/release/release-candidate-2026-06-18.md`、`docs/release/release-handoff-2026-06-18.md`；当前 release 状态为 `release_ready`。
- P64：已归档到 `openspec/changes/archive/2026-06-18-p64-release-packaging-version-tagging/`；已实现本地发布包脚本、sidecar manifest、archive checksum、verify 模式、输出目录限制、tracked/allowlisted source staging、敏感信息/禁止路径扫描和 P64 packaging handoff，不新增运行时交易、券商、外推、自动确认、自动升级、自动迁移、自动恢复或自动修复能力。
- P65：已归档到 `openspec/changes/archive/2026-06-18-p65-cross-machine-release-repeat-acceptance/`；已使用 P64 package workflow 生成 P65 candidate archive，并在本地跨机器等价隔离环境中完成 package verify、安装、OpenSpec/Go/frontend/E2E smoke 复验和发布材料更新；不声称物理第二机器已执行，不提前承诺远程发布、Git tag、自动升级、自动迁移、自动恢复或自动修复。
- P66：已归档到 `openspec/changes/archive/2026-06-18-p66-current-data-zero-degradation-policy/`；已把 P63/P65 release caveat 中的 current data-source degraded 状态转化为明确 policy verdict、阻断/豁免规则、可重复验收和发布声明边界；P66 当时本地库 strict current-data gate 为 `policy=blocked` / `gate=block`，后续 release-ready 声明必须通过 P66 gate 或显式记录豁免/范围排除；不新增外部源、真实 provider 调用、自动刷新、自动修复、券商接口或交易能力。P71 已在当前验收窗口提供 fresh strict pass 证据。
- P67：已归档到 `openspec/changes/archive/2026-06-18-p67-current-data-gate-resolution-workflow/`；已把 P66 `policy=blocked` / `gate=block` 转化为本地人工处置、豁免或范围排除记录与 release claim state；P67 当时 000300 为 `resolved_with_scope_exclusion`，但不改变 P66 policy，不声明当前数据 clean，不新增外部源、真实 provider 调用、自动刷新、自动修复、券商接口或交易能力。P71 的 current-data clean claim 来自 fresh strict gate pass，不来自 P67 scope exclusion。
- P68：已归档到 `openspec/changes/archive/2026-06-18-p68-post-p67-release-readiness-governance/`；已复核 P67 scope exclusion 后的 release-ready 表述、发布候选材料和打包复验策略；当前 release 状态为 `release_ready_limited_current_data_scope`，当时建议的 P69 clean-tree package refresh 已完成；只做治理和材料边界，不新增运行时能力。
- P69：已归档到 `openspec/changes/archive/2026-06-18-p69-clean-tree-package-refresh/`；`p69-clean-tree` package / verify / repeat acceptance 已通过，source commit `cc0a64781e199a7745432b63bce26de4402042b5`、`source_status=clean`；P69 本身只记录包后验收，不声称生成的包包含 P69 文档。
- P70：已归档到 `openspec/changes/archive/2026-06-18-p70-final-release-decision-and-risk-closure/`；已完成最终发布决策与风险收口，最终状态为 `release_ready_limited_current_data_scope`，确认 limited local release scope 无必需下一阶段。
- P71：已归档到 `openspec/changes/archive/2026-06-18-p71-real-product-acceptance-true-pass/`；已完成当前本地 `000300` strict current-data 真 pass、VecLite acceptance hardening、真实 LLM-backed strict UI runner、post-P70/P71 package refresh、package verify 和 isolated repeat acceptance，结论为 `release_ready_full_real_product_acceptance`；详见 `docs/release/acceptance/2026-06-18-p71-real-product-acceptance.md`。
- P72：已归档到 `openspec/changes/archive/2026-06-19-p72-real-user-fund-scenario-data-impact-acceptance/`；已完成真实用户基金/ETF `510300` 场景验收，覆盖组合维护、正式公开证据采集、local knowledge/RAG、真实 LLM 咨询、人工确认、SQLite 数据影响、每日纪律/风险/通知/审计/规则/workbench 回显和安全边界，结论为 `release_ready_full_real_user_scenario_acceptance`；详见 `docs/release/acceptance/2026-06-18-p72-real-user-fund-scenario.md`。
- P73：已归档到 `openspec/changes/archive/2026-06-19-p73-product-effectiveness-ux-validation/`；产品目标效果与 UX 验证已通过，覆盖纪律执行、证据充分性、可追溯性、复盘有效性、真实 UX 任务理解、背景材料阻断、安全手动确认、风险/readback/规则效果联动、390px reflow 和 unsafe input；结论为 `release_ready_product_effectiveness_ux_acceptance`。P73 不承诺未来收益、未来市场方向、真实用户研究或物理第二机器复验；详见 `docs/release/acceptance/2026-06-19-p73-product-effectiveness-ux-validation.md`。
- P74：已归档到 `openspec/changes/archive/2026-06-19-p74-built-in-knowledge-and-data-readiness/`；内置知识与数据准备度验收已通过，覆盖大师经验、纪律/SOP、标的画像、数据依赖矩阵、readiness API/UI、LLM 上下文引用和真实/降级场景；结论为 `release_ready_built_in_knowledge_data_readiness`。
- P75：已归档到 `openspec/changes/archive/2026-06-21-p75-requirements-traceability-and-real-use-closure/`；原始需求原子级追踪与真实使用闭环审计已完成，结论为 `release_ready_scoped_with_traceability_gaps`，不是 full original-requirement pass；详见 `docs/release/acceptance/2026-06-20-p75-real-use-closure.md` 与 `docs/release/acceptance/2026-06-20-p75-requirements-traceability-matrix.md`。
- P76：已归档到 `openspec/changes/archive/2026-06-21-p76-post-p75-final-package-refresh/`；已从 clean source commit `8a317f25917b8ff18ec9b5049e6a6188206a22d3` 生成 post-P75 package，verify 和 isolated repeat acceptance 均通过，包内确认包含 P72-P75 acceptance Markdown 与 OpenSpec archives；详见 `docs/release/acceptance/2026-06-21-p76-post-p75-package-refresh.md`。
- P77：已归档到 `openspec/changes/archive/2026-06-21-p77-requirements-real-pass-upgrade-gate/`；已建立 P75 后原子需求 `real_pass` 升级门禁和第一批升级证据，当前 P77 层结论为 `release_ready_scoped_with_p77_real_pass_progress`：17 行 `real_pass`、11 行 `reference_only`，313 个 full-release-required rows 仍非 `real_pass`；详见 `docs/release/acceptance/2026-06-21-p77-real-pass-upgrade-acceptance.md` 与 `docs/release/acceptance/2026-06-21-p77-requirements-real-pass-upgrade-matrix.md`。
- P78：已归档到 `openspec/changes/archive/2026-06-21-p78-requirements-real-pass-batch-closure/`；已生成 P78 批次收敛层，当前 P78 层结论为 `release_ready_scoped_with_p78_real_pass_batch_progress`：20 行 `real_pass`、11 行 `reference_only`，310 个 full-release-required rows 仍非 `real_pass`；详见 `docs/release/acceptance/2026-06-21-p78-requirements-real-pass-batch-closure.md` 与 `docs/release/acceptance/2026-06-21-p78-requirements-real-pass-batch-matrix.md`。
- P79：已归档到 `openspec/changes/archive/2026-06-21-p79-real-use-data-impact-and-expected-return-closure/`；已生成 P79 真实使用数据影响层，当前 P79 层结论为 `release_ready_scoped_with_p79_real_use_data_impact_progress`：43 行 `real_pass`、11 行 `reference_only`，287 个 full-release-required rows 仍非 `real_pass`；P79 不升级宽泛审计/月度归因行，不刷新 P76 package，不声称 full original-requirement pass；详见 `docs/release/acceptance/2026-06-21-p79-real-use-data-impact-and-expected-return-closure.md` 与 `docs/release/acceptance/2026-06-21-p79-real-use-data-impact-and-expected-return-matrix.md`。
- P80：已归档到 `openspec/changes/archive/2026-06-22-p80-review-audit-governance-real-use-closure/`；已生成 P80 复盘审计与规则治理真实使用闭环层，当前 P80 层结论为 `release_ready_scoped_with_p80_review_audit_governance_progress`：57 行 `real_pass`、11 行 `reference_only`，273 个 full-release-required rows 仍非 `real_pass`；P80 不升级月度/季度归因、完整提案生成、最终规则应用时间，不刷新 P76 package，不声称 full original-requirement pass。
- P81：已归档到 `openspec/changes/archive/2026-06-22-p81-dynamic-source-field-coverage/`；已生成 P81 动态源字段覆盖层，当前 P81 层结论为 `release_ready_scoped_with_p81_dynamic_source_progress`：116 行 `real_pass`、11 行 `reference_only`，214 个 full-release-required rows 仍非 `real_pass`；详见 `docs/release/acceptance/2026-06-22-p81-dynamic-source-field-coverage.md` 与 `docs/release/acceptance/2026-06-22-p81-dynamic-source-field-coverage-matrix.md`。
- P82：已归档到 `openspec/changes/archive/2026-06-22-p82-sop-action-ui-sqlite-closure/`；已生成 P82 SOP/action UI-to-SQLite 闭环层，当前 P82 层结论为 `release_ready_scoped_with_p82_sop_action_progress`：160 行 `real_pass`、11 行 `reference_only`，170 个 full-release-required rows 仍非 `real_pass`；P82 评估 53 行，新增 44 行 `real_pass`，9 行 deferred 到 P83/P84/P86；详见 `docs/release/acceptance/2026-06-22-p82-sop-action-ui-sqlite-closure.md` 与 `docs/release/acceptance/2026-06-22-p82-sop-action-ui-sqlite-matrix.md`。
- P83：已归档到 `openspec/changes/archive/2026-06-22-p83-governance-traceability-backfill/`；已生成 P83 governance traceability evidence layer，当前 P83 层结论为 `release_ready_scoped_with_p83_governance_traceability_progress`：170 行 `real_pass`、11 行 `reference_only`，160 个 full-release-required rows 仍非 `real_pass`；P83 评估 43 行，新增 10 行 `real_pass`，33 行因证据口径过宽 deferred 到 P86；详见 `docs/release/acceptance/2026-06-22-p83-governance-traceability-backfill.md` 与 `docs/release/acceptance/2026-06-22-p83-governance-traceability-matrix.md`。
- P84：已归档到 `openspec/changes/archive/2026-06-22-p84-portfolio-confirmation-data-impact-closure/`；已生成 P84 portfolio/confirmation data-impact evidence layer，当前 P84 层结论为 `release_ready_scoped_with_p84_portfolio_confirmation_progress`：173 行 `real_pass`、11 行 `reference_only`，157 个 full-release-required rows 仍非 `real_pass`；P84 评估 35 行，新增 3 行 `real_pass`，32 行因证据口径不足 deferred 到 P87 或 P86；详见 `docs/release/acceptance/2026-06-22-p84-portfolio-confirmation-data-impact-closure.md` 与 `docs/release/acceptance/2026-06-22-p84-portfolio-confirmation-data-impact-matrix.md`。
- P85：已归档到 `openspec/changes/archive/2026-06-22-p85-expected-return-analysis-accuracy-closure/`；已生成 P85 expected-return analysis-accuracy evidence layer，当前 P85 层结论为 `release_ready_scoped_with_p85_expected_return_progress`：188 行 `real_pass`、11 行 `reference_only`，142 个 full-release-required rows 仍非 `real_pass`；P85 评估 31 行，新增 15 行 `real_pass`，16 行因历史准确性/回测/概率口径不足 deferred；本环境无 `DEEPSEEK_API_KEY`，不声称 fresh real LLM output；详见 `docs/release/acceptance/2026-06-22-p85-expected-return-analysis-accuracy-closure.md` 与 `docs/release/acceptance/2026-06-22-p85-expected-return-analysis-accuracy-matrix.md`。
- P87：已归档到 `openspec/changes/archive/2026-06-22-p87-portfolio-state-allocation-safety-closure/`；组合状态/仓位纪律/安全闭环已完成验收，32 行评估、5 行升级为 `real_pass`、27 行 deferred；P87 后仍有 137 个 full-release-required rows 非 `real_pass`。
- P86：已归档到 `openspec/changes/archive/2026-06-22-p86-core-goal-knowledge-safety-final-closure/`；fresh integrated runner 复跑 P74/P81/P82/P83/P84/P85/P87 真实 UI/API/SQLite/Go evidence，当前 P86 层结论为 `release_ready_scoped_with_p86_final_integrated_progress`：303 行 `real_pass`、11 行 `reference_only`、27 个 full-release-required rows 仍非 `real_pass`；详见 `docs/release/acceptance/2026-06-22-p86-core-goal-knowledge-safety-final-closure.md` 与 `docs/release/acceptance/2026-06-22-p86-core-goal-knowledge-safety-final-matrix.md`。

## 推荐阅读路径

| 目标 | 起点 | 继续阅读 |
| --- | --- | --- |
| 了解产品愿景 | `docs/requirements.md` | `docs/functional-spec.md` |
| 规划下一阶段 | `openspec/PROGRESS.md` | `docs/development-plan.md`、`docs/GOVERNANCE.md` |
| 审计 P19-P24 历史交付 | `docs/p19-p24-audit-evidence-pack.md` | `docs/p19-p24-historical-archive-traceability.md`、`openspec/PROGRESS.md` |
| 发布前验收门禁 | `docs/project-acceptance-gate-matrix.md` | `docs/testing-plan.md`、`docs/p19-p24-audit-evidence-pack.md` |
| 查看发布候选状态 | `docs/release/release-candidate-2026-06-18.md` | `docs/release/acceptance/2026-06-18-p72-real-user-fund-scenario.md`、`docs/release/acceptance/2026-06-18-p71-real-product-acceptance.md`、`docs/project-acceptance-gate-matrix.md` |
| 发布交付与复验 | `docs/release/README.md` | `docs/release/release-handoff-2026-06-18.md`、`docs/release/acceptance-repeatability.md` |
| 前端真实 UI 验收 | `docs/release/ui-acceptance-p56-2026-06-17.md` | `docs/release/ui-design-review-p56-2026-06-17.md`、`docs/release/ui-audit-assets/2026-06-17-p56/`、P55 历史记录：`docs/release/ui-acceptance-2026-06-17.md` |
| 产品体验打磨路线图 | `docs/product-experience-polish-roadmap.md` | `docs/frontend-contract.md`、`docs/development-plan.md` |
| 修改契约或行为 | `docs/GOVERNANCE.md` | `openspec/project.md`、相关 L1 契约文档 |
| 对接 HTTP API | `docs/api.md` | `docs/data-model.md`、`docs/frontend-contract.md` |
| 理解数据持久化 | `docs/data-model.md` | `docs/workflow.md` |
| 理解工作流 | `docs/workflow.md` | `docs/architecture.md` |
| 前端开发 | `docs/frontend-contract.md` | `docs/ui-design.md`、`docs/ui-flow.md`、`docs/ui/prototype.md` |
| 本地运行和配置 | `docs/configuration.md` | `configs/config.example.yaml` |
| 本地调度、备份恢复 | `docs/ops-local-scheduler.md` | `examples/scheduler/` |
| 测试与验收 | `docs/testing-plan.md` | `docs/development-plan.md`；有标准 OpenSpec change 时再看对应 `tasks.md` |
| 迁移和交付策略 | `docs/migration-plan.md` | `docs/architecture.md` |

## 文档分级

| 级别 | 文档 | 维护规则 |
| --- | --- | --- |
| L1 契约真源 | `requirements.md`、`data-model.md`、`api.md`、`workflow.md`、`frontend-contract.md` | 行为、接口、状态流转和字段变更必须通过 OpenSpec delta 后合并 |
| L2 架构与计划 | `architecture.md`、`functional-spec.md`、`development-plan.md` | 随契约和阶段变化同步；阶段状态不能滞后于实现 |
| L3 UI 与图示 | `ui-design.md`、`ui-flow.md`、`ui/`、`diagrams/` | 随相关 UI 或流程变更同步 |
| 运维与交付文档 | `configuration.md`、`ops-local-scheduler.md`、`testing-plan.md`、`migration-plan.md` | 随 CLI、配置、部署和验收方式变化同步 |
| OpenSpec 治理入口 | `openspec/project.md`、`openspec/PROGRESS.md`、`openspec/changes/` | `project.md` 与 `PROGRESS.md` 是当前治理入口；活跃 change 可编辑；archive 后作为历史审计，不回改历史归档 |

## 当前目录结构决策

本轮只新增文档地图并校准过期内容，不大规模移动文档路径。原因：

- `docs/GOVERNANCE.md` 和 `openspec/project.md` 已将多个文件列为契约真源，直接搬迁会造成链接和治理规则同时失效。
- `docs/diagrams/`、`docs/ui/`、OpenSpec archive 中存在大量历史引用，批量移动会增加维护风险。
- 当前最重要的问题是阶段状态、实现事实和契约内容漂移，而不是路径本身。

现阶段保留结构：

```text
docs/
  README.md
  requirements.md
  functional-spec.md
  development-plan.md
  architecture.md
  data-model.md
  api.md
  workflow.md
  frontend-contract.md
  configuration.md
  ops-local-scheduler.md
  testing-plan.md
  migration-plan.md
  ui-design.md
  ui-flow.md
  GOVERNANCE.md
  templates/
  ui/
  diagrams/
```

后续若文档继续增长，可单独发起文档架构变更，将内容渐进拆分为：

```text
docs/product/
docs/architecture/
docs/frontend/
docs/operations/
docs/governance/
docs/assets/
```

该拆分不应与功能开发混在同一 change 中执行。

## 维护规则

1. `docs/` 是契约真源；`openspec/specs/` 只保存从 `docs/` 抽取的行为摘要或保持为空，不与 `docs/` 维护两套全文。
2. 修改 L1 契约前必须创建 OpenSpec change；delta 写在 `openspec/changes/<id>/specs/**/*.md`，并在 archive 时合并回 `docs/`。
3. 实现完成后必须同步 `openspec/PROGRESS.md`、`AGENTS.md`、`docs/GOVERNANCE.md` 和 `docs/development-plan.md` 的阶段状态。
4. 新增 API 字段必须同时检查 `docs/api.md`、`docs/frontend-contract.md` 和对应 DTO/类型。
5. 新增表或迁移必须同步 `docs/data-model.md`。
6. 新增 CLI flag、配置项或运维流程必须同步 `docs/configuration.md` 和 `docs/ops-local-scheduler.md`。
7. 不在历史 archive 中回写内容；若历史描述过期，在当前文档中写明现状。
8. 文档不得承诺自动交易、券商接口、付费/登录数据源或外部通知渠道，除非后续需求明确重新立项。
