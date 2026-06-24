# Investment Agent 开发计划

> 文档版本：v1.1  
> 最后更新：2026-06-24
> 适用范围：后端、前端、数据存储、Eino 工作流、审计与验收。  
> 配套文档：`docs/requirements.md`、`docs/architecture.md`、`docs/data-model.md`、`docs/api.md`、`docs/workflow.md`、`docs/frontend-contract.md`、`docs/ui-design.md`、`docs/ui-flow.md`。  
> 文档治理：`docs/GOVERNANCE.md`；变更流程：OpenSpec（`openspec/project.md`、`AGENTS.md`）。

## 1. 目标

本计划用于把现有需求、架构、数据模型、API、工作流和 UI 契约转换为可分派、可验收的开发任务。

开发完成后，系统应具备以下能力：

1. 使用 SQLite 保存账户、持仓、行情、情报摘要、决策记录、用户确认、规则版本和审计事件。
2. 使用 VecLite 作为可重建的辅助检索索引。
3. 使用 Eino 组织每日纪律、主动咨询、证据验证、规则提案和守门人审计工作流。
4. 使用 DeepSeek 生成分析材料，但最终裁决由领域规则完成。
5. 使用 React + Vite + TypeScript 提供 Agent 决策驾驶舱。
6. 严格禁止自动交易。所有账户变化只来自用户记录的线下动作。
7. 每次正式建议、用户确认、规则提案、守门人审计和规则应用都生成 `audit_events`。

### 1.1 阶段定位（与 `docs/requirements.md` 的关系）

| 层级 | 阶段 | 状态 | 说明 |
| --- | --- | --- | --- |
| 全量愿景 | `docs/requirements.md` v3.0 | 真源 | 描述完整产品能力，不按 MVP 裁剪 |
| 工程骨架 | P0–P18 | **已完成** | 数据模型、规则引擎、工作流、API、前端、测试与治理基础已完成 |
| 可用 MVP | P19–P24 | **已完成** | 公开 HTTP 数据桥接、A 股 ETF/基金证据 payload 解析、应用内通知、规则提案增强、复盘深化、配置校验与本地备份恢复已完成；真实外部公开源 collector 尚未完成 |
| 真实公开源验证 | P25 | **已完成** | 已输出首轮验证结论：巨潮、深交所、证监会可作为 P26 首批公告/监管证据源；中证指数、东方财富基金可作为 P27 首批指数/基金数据源；上交所、AMAC 查询和新浪财经需按验证结论分级处理 |
| 公告与证据源 collector | P26 | **已完成** | 已实现巨潮资讯、深交所、证监会首批只读 collector、默认 90 天/分页边界、标准证据 payload、去重幂等、RAG 入库、source verification、审计和 `cmd/agent --task public-evidence-refresh` 手动触发入口；AMAC 暂缓为二线背景候选 |
| 基金净值与 ETF 市场数据 collector | P27 | **已完成** | 已接入东方财富基金真实净值/历史净值/资产配置/基金档案基础 metadata collector，并将中证指数基础信息 collector 校准到当前公开 `index-basic-info` shape；默认关闭，不接交易、登录、付费、授权、Level2 或实时估算净值正式化；中证指数样本/权重/估值文件扩展接口仍按候选 metadata 低频读取 |
| 预期收益与动态卖出评估 | P28 | **已完成** | 已增强预期收益样本上下文、精度状态、情景触发、动态卖出评估和重新评估触发；所有输出仅作为人工复核材料，不承诺收益、不覆盖规则裁决、不执行交易 |
| 公开证据真实采集 smoke | P29 | **已完成** | 已修复 P26 公开证据 collector 的真实接口参数、no-data/source-unavailable/parse-error 诊断和临时 SQLite 入库 smoke 验收；CNInfo 真实公告 smoke 已通过，SZSE/CSRC 按 no_data/source_unavailable 分类降级，不阻塞可用 A 级公告源 |
| 真实环境 E2E / Playwright smoke | P30 | **已完成** | 已实现本地 server + Vite + Playwright smoke，使用临时 SQLite/配置和可控 seed 数据验证健康检查、决策详情 expected return、证据页与审计页只读路径，并治理浏览器临时产物 |
| 每日自动运行闭环 | P31 | **已完成** | 已实现默认关闭但可显式启用的本地每日自动运行，串联市场/证据刷新、每日纪律、应用内通知、审计、幂等和前端状态展示 |
| 每日纪律报告产品化 | P32 | **已完成** | P32 已归档到 `openspec/changes/archive/2026-06-11-p32-daily-discipline-report-productization/`；已将 daily workflow 与 P31 auto-run 结果统一为今日纪律报告、历史报告、详情回看和 smoke 可验证入口，并修复复审 findings |
| 账户与持仓录入/校准体验 | P33 | **已完成** | 已完成账户/持仓初始化、本地事实录入、校准、导入和 Dashboard/每日纪律缺前提引导 |
| 真实数据覆盖扩展 | P34 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-15-p34-real-data-coverage-expansion/`；已扩展真实公开数据覆盖、数据源健康、失败分类、工作流输入上下文和前端状态展示 |
| 风险预警与 SOP 编排 | P35 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-16-p35-risk-alert-sop-orchestration/`；已建立本地风险预警事实、SOP 状态流转、通知/审计/每日纪律报告联动和风险预警中心；不接券商、不自动交易、不外部推送 |
| 规则进化效果验证 | P36 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-16-p36-rule-evolution-effect-validation/`；已建立规则提案来源解释、样本代表性、过拟合检查、历史回放、守门人门禁接入、应用后追踪和前端展示；不自动应用规则、不自动交易 |
| 真实 LLM 使用与质量评估 | P37 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-16-p37-real-llm-quality-evaluation/`；已建立真实 LLM 配置 smoke、模型/超时配置、错误分类、prompt 与输出摘要、质量门禁、工作流/API metadata 和脱敏审计；LLM 只生成分析材料，不写最终裁决 |
| RAG / VecLite 检索质量加固 | P38 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-16-p38-rag-veclite-retrieval-quality/`；已建立本地 retrieval quality fixture、quality-aware ranking、source verification / RAG metadata 一致性校验、index freshness、审计/API/前端展示和 `retrieval-quality-smoke`；不绕过规则裁决、不自动交易 |
| 前端完整用户旅程与全路径 E2E | P39 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-16-p39-frontend-full-user-journey-e2e/`；已补齐空库 onboarding、配置/账户初始化、每日纪律、主动咨询、确认记录、风险预警、复盘、规则治理、降级路径和浏览器级 console/a11y/窄屏验收；不接券商、不自动交易、不外部推送、不自动应用规则、不收益承诺 |
| 本地部署、运维与恢复演练 | P40 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-16-p40-local-deploy-ops-recovery-drill/`；已建立本地预检、启动诊断、备份恢复 smoke、数据源健康面板、日志/临时文件治理和恢复演练文档；不接券商、不自动交易、不外部推送、不自动应用规则、不收益承诺 |
| P40 后路线图治理 | P41 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-16-p41-post-p40-roadmap-governance/`；已固化 P42-P44 与历史审计追溯候选队列 |
| 用户决策工作台 | P42 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-16-p42-user-decision-workbench/`；已新增只读日常工作台入口 |
| 数据质量可观测 | P43 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-16-p43-data-quality-observability/`；已新增只读数据质量面板 |
| 本地安装诊断与打包 | P44 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-16-p44-local-install-diagnostics-packaging/`；已新增本地安装、诊断打包和 smoke 汇总入口 |
| P44 后路线图治理 | P45 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-17-p45-post-p44-roadmap-governance/`；已固化 P46-P49 候选队列、依赖、验收与安全边界，未修改运行时代码 |
| 本地知识导入治理 | P46 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-17-p46-local-knowledge-import-governance/`；已实现 validate、脱敏预览、显式确认、C/background 背景事实写入和 pending 索引计划 |
| 决策闭环解释 | P47 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-17-p47-decision-loop-explainability/`；已实现建议、确认、线下记录、风险/审计/复盘线索和缺口说明的只读解释链 |
| 数据源质量回归包 | P48 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-17-p48-data-source-quality-regression-pack/`；已提供 fixture/current 回归、freshness/失败分类验证、脱敏摘要、只读 API 和 CLI 脱敏审计 |
| 本地发布与升级体验 | P49 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-17-p49-local-release-upgrade-experience/`；已新增本地发布/升级检查、脱敏诊断 JSON、升级前备份提醒、迁移预检和升级后 smoke 汇总 |
| P49 后治理与验收路线图 | P50 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-17-p50-post-p49-governance-validation-roadmap/`；已固化 P51 补 P19-P24 审计证据包、P52 建立项目验收门禁矩阵、P53 执行验收并整理发布候选材料 |
| P19-P24 审计证据包 | P51 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-17-p51-p19-p24-audit-evidence-pack/`；已新增 `docs/p19-p24-audit-evidence-pack.md`，补当前事实证据矩阵，不伪造历史 archive，不修改运行时代码；下一阶段建议 P52 验收门禁矩阵 |
| 项目验收门禁矩阵 | P52 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-17-p52-project-acceptance-gate-matrix/`；已新增 `docs/project-acceptance-gate-matrix.md`，定义单元、集成、E2E、真实源、真实 LLM、冒烟、安装诊断、发布升级和安全边界门禁，不宣称验收已通过 |
| 验收执行与发布候选材料 | P53 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-17-p53-acceptance-execution-and-release-candidate-materials/`；已按 P52 G0-G9 执行实际验收，生成验收记录和发布候选材料；当前结论为 `release_ready`，含 G5 current degraded 与重试记录 |
| 发布交付与可重复性加固 | P54 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-17-p54-release-handoff-and-repeatability-hardening/`；已新增 release handoff、复验规则、重试处理、G5 current degraded 解释和安全边界文档 |
| 前端全功能真实验收与设计审查 | P55 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-17-p55-full-ui-acceptance-and-design-audit/`；已真实启动项目，通过前端 UI 操作验收主要功能，并使用 Product Design audit 审查 UI 优化项；full UI acceptance 结论为 `blocked`，原因是真实 LLM 决策详情 nullable `final_verdict.optional_actions` 前端崩溃 |
| UI 验收阻断与产品化设计修复 | P56 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-17-p56-ui-acceptance-blocker-fixes/`；已修复 P55-B1，完成任务分组导航、产品化 UI 基础层、`/positions` 与 `/data-quality` 移动端 reflow，并通过真实 LLM UI 验收和 Playwright E2E；当前 full UI acceptance 状态为 `p56_scope_pass` |
| 产品体验打磨总规划 | P57 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-17-p57-product-experience-polish-roadmap/`；已新增 `docs/product-experience-polish-roadmap.md`，固化产品设计、UI 设计和功能设计打磨路线图，拆分 P58-P63 后续阶段；原发布状态刷新后移 |
| 今日工作台重构 | P58 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-17-p58-daily-workbench-redesign/`；已实现 Dashboard/Workbench 每日投资纪律 cockpit，并完成范围内单测、构建、Go 测试、真实 smoke/E2E、桌面/移动 UI 验收和 safety scan |
| 决策解释体验重构 | P59 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-17-p59-decision-explainability-experience/`；已重构 Consultation、Decision Detail、Evidence、Decision Loop 的解释链路，完成真实 UI 验收、nullable DTO 加固、桌面/移动截图和安全降级记录 |
| 组合、风险与数据质量体验重构 | P60 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-18-p60-portfolio-risk-data-quality-experience/`；完成 Positions、Risk Alerts、Data Quality 的维护、处置和质量可观测体验重构、真实本地 UI 验收、桌面/移动截图和安全边界扫描 |
| 治理和运维页面产品化 | P61 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-18-p61-governance-ops-productization/`；已完成 Rules、Audit、Notifications、Daily Reports、Daily Auto Run、Local Install、Local Knowledge、Settings 的产品化重构、真实本地 UI 验收、桌面/移动截图和安全边界扫描 |
| 设计系统与可访问性验收 | P62 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-18-p62-design-system-accessibility-hardening/`；已固化组件基础层、状态体系、键盘路径、WCAG reflow 和视觉回归，验收记录见 `docs/release/acceptance/2026-06-18-p62-ui-acceptance.md` |
| 全量真实 UI 回归与发布状态刷新 | P63 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-18-p63-full-ui-regression-release-refresh/`；已执行 P52 G0-G9、20 个主要路由 × 3 视口真实 UI 回归、真实 LLM-backed consultation UI journey，并刷新 `docs/release/acceptance/2026-06-18-p63-full-ui-regression.md`、`docs/release/release-candidate-2026-06-18.md`、`docs/release/release-handoff-2026-06-18.md`，当前 release 状态为 `release_ready` |
| 发布打包与版本标记 | P64 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-18-p64-release-packaging-version-tagging/`；已实现本地发布包脚本、sidecar manifest、archive checksum、verify 模式、输出目录限制、tracked/allowlisted source staging、敏感信息/禁止路径扫描和 P64 packaging handoff；不新增交易、券商、外推、自动确认、自动升级、自动迁移、自动恢复或自动修复能力 |
| 跨机器发布包复验 | P65 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-18-p65-cross-machine-release-repeat-acceptance/`；已使用 P64 package workflow 生成 P65 candidate archive，并在本地跨机器等价隔离环境中完成 package verify、安装、OpenSpec/Go/frontend/E2E smoke 复验和发布材料更新；未声称物理第二机器已执行 |
| 当前数据零退化策略 | P66 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-18-p66-current-data-zero-degradation-policy/`；已把 P63/P65 release caveat 中的 current data-source degraded 状态转化为明确 policy verdict、阻断/豁免规则、可重复验收和发布声明边界；P66 当时本地库 strict current-data gate 为 `policy=blocked` / `gate=block`，后续 release-ready 声明必须通过 P66 gate 或显式记录豁免/范围排除；P71 已在当前验收窗口提供 fresh strict pass 证据 |
| 当前数据门禁处置工作流 | P67 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-18-p67-current-data-gate-resolution-workflow/`；围绕 P66 `policy=blocked` / `gate=block` 提供本地人工处置、豁免或范围排除记录与 release claim state、API/UI/CLI、真实 UI 操作验收和 release acceptance；P67 当时 000300 为 `resolved_with_scope_exclusion`，但不得把该 scope exclusion 声明为当前数据 clean；P71 clean claim 来自 fresh strict gate pass |
| P67 后发布状态治理 | P68 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-18-p68-post-p67-release-readiness-governance/`；已刷新 P68 发布状态治理材料，当前 release 状态为 `release_ready_limited_current_data_scope`，当时建议的 P69 clean-tree package refresh 已完成；不直接新增交易、外部源、自动刷新、自动修复或 provider 能力 |
| Clean tree 最终分发包刷新 | P69 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-18-p69-clean-tree-package-refresh/`；`p69-clean-tree` package / verify / repeat acceptance 已通过，source commit `cc0a64781e199a7745432b63bce26de4402042b5`、`source_status=clean`；P69 本身只记录包后验收，不声称生成的包包含 P69 文档；不得新增远程发布、Git tag、自动升级、自动迁移、自动恢复、自动修复、券商接口、交易、外推或收益承诺 |
| 最终发布决策与风险收口 | P70 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-18-p70-final-release-decision-and-risk-closure/`；最终状态为 `release_ready_limited_current_data_scope`，确认 limited local release scope 无必需下一阶段；可选后续仅包括物理第二机器复验、P66 当前数据真 pass、post-P70 package refresh 或 VecLite acceptance hardening |
| 真实产品验收真通过 | P71 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-18-p71-real-product-acceptance-true-pass/`；当前本地 `000300` strict current-data gate 已真 pass，VecLite rebuild / retrieval 已达到 `hit` + `fallback=veclite` + `index=healthy` + `index_freshness=fresh`，真实本地 UI / LLM strict runner 已通过，post-P70/P71 package refresh、verify/repeat acceptance 已通过，结论为 `release_ready_full_real_product_acceptance` |
| 真实用户基金场景与数据影响验收 | P72 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-19-p72-real-user-fund-scenario-data-impact-acceptance/`；真实基金/ETF `510300` 场景验收已完成，覆盖 UI 操作、正式公开证据采集、真实 LLM 咨询、人工确认、SQLite 数据影响、审计事件、衍生页面回显、确定性计算准确性、真实 LLM 分析链路和安全边界；结论为 `release_ready_full_real_user_scenario_acceptance`；不新增券商接口、自动交易、收益承诺或未来预测准确性声明 |
| 产品目标效果与 UX 验证 | P73 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-19-p73-product-effectiveness-ux-validation/`；真实浏览器 UX 任务、产品效果 replay checker、runner、截图、browser results、effect replay summary 和 UX audit 已通过；结论为 `release_ready_product_effectiveness_ux_acceptance`；不新增券商接口、自动交易、收益承诺、未来市场方向预测、真实用户研究或物理第二机器复验 |
| 内置知识与数据准备度 | P74 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-19-p74-built-in-knowledge-and-data-readiness/`；已新增 7 位大师经验内置知识注册表、readiness API、数据依赖矩阵、LLM 上下文摘要、`/data-quality` 完整展示、决策详情脱敏回显、API 场景矩阵和真实浏览器 UI 验收；结论为 `release_ready_built_in_knowledge_data_readiness`；不新增券商接口、自动交易、收益承诺、未来预测、登录/付费/授权/Level2/高频源或自动规则应用 |
| 原始需求追踪与真实使用闭环 | P75 | **已完成，存在 scoped gaps** | 已归档到 `openspec/changes/archive/2026-06-21-p75-requirements-traceability-and-real-use-closure/`；已生成 341 个原子需求追踪 rows、P75 acceptance record 和 traceability summary；结论为 `release_ready_scoped_with_traceability_gaps`，不是 full original-requirement pass；expanded G9 分类人审未发现 forbidden runtime affordance；主要剩余缺口包括动态非 `510300` 真实外部查询、全字段基金/指数/benchmark join、缺失数据传播、确定性阈值覆盖、SOP A-F 真实 UI/data-impact、完整 UI action-to-SQLite-to-readback 验收 |
| P75 后最终分发包刷新 | P76 | **已完成** | 已归档到 `openspec/changes/archive/2026-06-21-p76-post-p75-final-package-refresh/`；已从 clean source commit `8a317f25917b8ff18ec9b5049e6a6188206a22d3` 生成 `p76-post-p75-final` package，verify 和 isolated repeat acceptance 均通过，包内确认包含 P72-P75 acceptance Markdown 与 OpenSpec archives；P76 不扩大 P75 scoped 结论，不声称物理第二机器复验、远程发布、Git tag、自动升级/迁移/恢复/修复、券商接口、交易、外推、自动确认、自动规则应用、未来 provider 可用性或收益承诺 |
| 原子需求真 pass 升级门禁 | P77 | **已完成，存在 scoped gaps** | 已归档到 `openspec/changes/archive/2026-06-21-p77-requirements-real-pass-upgrade-gate/`；已建立 P75 后 `real_pass` 升级门禁、P77 矩阵和第一批新鲜证据；当前结论 `release_ready_scoped_with_p77_real_pass_progress`，341 rows 中 17 `real_pass`、11 `reference_only`、22 `scoped_pass`、291 `partial`，仍有 313 个 full-release-required rows 非 `real_pass`；禁止重写 P75 历史矩阵或把 scoped/partial/deterministic-local、fixture/mock/stub、截图、route smoke、scope exclusion、waiver、临时 DB、单标的证据冒充 full original-requirement pass |
| 原始需求 real-pass 批次收敛 | P78 | **已完成，存在 scoped gaps** | 已归档到 `openspec/changes/archive/2026-06-21-p78-requirements-real-pass-batch-closure/`；已把 P77 后 full-release-required 非 `real_pass` rows 分类成可执行批次，并用新鲜 Go 测试、真实 UI 操作和 SQLite readback 推进第一批预期收益/降级/免责声明 rows；当前 20 行 `real_pass`，仍有 310 个 full-release-required rows 非 `real_pass`；不刷新 P76 package，不声称 full original-requirement pass |
| 真实使用数据影响与预期收益闭环 | P79 | **已完成，存在 scoped gaps** | 已归档到 `openspec/changes/archive/2026-06-21-p79-real-use-data-impact-and-expected-return-closure/`；fresh P72 `510300` 真实用户 UI 场景和 P75 accepted-local 非 `510300` 真实 UI 场景均通过，字段级 SQLite/readback 门禁和 ExpectedReturnNode LLM `quality_failure` 安全兜底已完成；当前 43 行 `real_pass`，仍有 287 个 full-release-required rows 非 `real_pass`；不升级宽泛审计/月度归因行，不刷新 P76 package，不声称 full original-requirement pass |
| 复盘审计与规则治理真实使用闭环 | P80 | **已完成，存在 scoped gaps** | 已归档到 `openspec/changes/archive/2026-06-22-p80-review-audit-governance-real-use-closure/`；fresh P75 SOP/failure real UI journey 和字段级 SQLite/readback 已通过，当前 57 行 `real_pass`，仍有 273 个 full-release-required rows 非 `real_pass`；不刷新 P76 package，不声称 full original-requirement pass |
| 动态源字段覆盖 | P81 | **已完成并归档** | 已归档到 `openspec/changes/archive/2026-06-22-p81-dynamic-source-field-coverage/`；fresh `159915` accepted-local 非 `510300` 真实浏览器 UI journey、Go collector/readiness 测试、SQLite/readback、source-health provenance、formal evidence、RAG indexing 和真实 LLM-backed UI readback 均通过；当前 116 行 `real_pass`，仍有 214 个 full-release-required rows 非 `real_pass`；不刷新 P76 package，不声称 full original-requirement pass |
| SOP/action UI-to-SQLite 闭环 | P82 | **已完成并归档** | 已归档到 `openspec/changes/archive/2026-06-22-p82-sop-action-ui-sqlite-closure/`；fresh P82 real browser UI journey、显式最终确认应用本地规则版本、SQLite/readback、forbidden table absence 和只读 checker 均通过；P82 评估 53 行，新增 44 行 `real_pass`，9 行因证据口径过宽 deferred；当前 160 行 `real_pass`，仍有 170 个 full-release-required rows 非 `real_pass`；不刷新 P76 package，不声称 full original-requirement pass |
| 治理追溯回填 | P83 | **已完成并归档** | 已归档到 `openspec/changes/archive/2026-06-22-p83-governance-traceability-backfill/`；Fresh P83 real browser UI/API/SQLite/Go evidence 已通过，覆盖 review/governance/release traceability、monthly/quarterly review、notification mark-read、local-install redaction、SQLite 字段级 readback、focused Go tests 和 forbidden table absence；P83 评估 43 行，新增 10 行 `real_pass`，33 行因证据口径过宽 deferred 到 P86；当前 170 行 `real_pass`，仍有 160 个 full-release-required rows 非 `real_pass`；不刷新 P76 package，不声称 full original-requirement pass |
| 组合与确认数据影响闭环 | P84 | **已完成并归档** | 已归档到 `openspec/changes/archive/2026-06-22-p84-portfolio-confirmation-data-impact-closure/`；Fresh P84 real browser UI/API/SQLite/Go evidence 已通过，覆盖 `/positions` 本地账户校准、持仓编辑、批量导入、线下交易、修正审计、手动确认、下游读回、SQLite 字段级 readback、focused handler tests 和 forbidden table/auto-confirm absence；P84 评估 35 行，新增 3 行 `real_pass`，32 行因证据口径不足 deferred；当前 173 行 `real_pass`，仍有 157 个 full-release-required rows 非 `real_pass`；子 agent re-review 结论 `can_archive` |
| 预期收益与分析准确性闭环 | P85 | **已完成并归档** | 已归档到 `openspec/changes/archive/2026-06-22-p85-expected-return-analysis-accuracy-closure/`；Fresh P85 real browser UI/API/SQLite/Go evidence 已通过，覆盖 `/consultation` 目标收益率与上一轮基准情景中枢 UI 输入、完整样本/下行情景/样本不足三类咨询、决策详情 readback、SQLite 字段级 readback、focused workflow/handler tests 和 forbidden table/auto-confirm absence；31 行评估、15 行升级为 `real_pass`、16 行因历史准确性/回测/概率口径不足 deferred；当前 188 行 `real_pass`，仍有 142 个 full-release-required rows 非 `real_pass`；本环境无 `DEEPSEEK_API_KEY`，不声称 fresh real LLM output |
| 组合状态与仓位安全闭环 | P87 | **已完成并归档** | 已归档到 `openspec/changes/archive/2026-06-22-p87-portfolio-state-allocation-safety-closure/`；Fresh P87 real browser UI/API/SQLite/Go evidence 已通过，评估 32 行、5 行升级为 `real_pass`、27 行 deferred；当前剩余 137 个 full-release-required rows 非 `real_pass`，P87 不刷新 P76 package、不声称 full original-requirement pass |
| 剩余 137 行最终执行队列 | P86 | **已完成并归档** | 已归档到 `openspec/changes/archive/2026-06-22-p86-core-goal-knowledge-safety-final-closure/`；P86 fresh integrated runner 和 closure 已通过，复跑 P74/P81/P82/P83/P84/P85/P87 真实 UI/API/SQLite/Go evidence，新增 110 行 `real_pass`；当前结论为 `release_ready_scoped_with_p86_final_integrated_progress`，341 rows 中 303 `real_pass`、11 `reference_only`、27 `partial`；仍有 27 个 full-release-required rows 非 `real_pass`，不声称 full original-requirement pass |
| 剩余 27 行 blockers 闭环 | P88 | **已完成并归档** | 已归档到 `openspec/changes/archive/2026-06-22-p88-remaining-full-release-blockers-closure/`；P88 fresh real browser UI/API/SQLite/Go acceptance 已通过，覆盖 source-verified `sell_only`、single-source `frozen_watch`、historical expected-return probabilities、sample<5 degradation、quarterly rebalance 和 SOP addendum proposal；27 行中 17 行升级为 `real_pass`，10 行保留 `partial`；仍有 10 个 full-release-required rows 非 `real_pass`，不声称 full original-requirement pass |
| 剩余 10 行真实 provider 与动态概率闭环 | P89 | **待立项** | 建议 change：`p89-real-provider-and-dynamic-probability-closure`；处理 P88 后剩余 `REQ-04-016`、`REQ-05-003`、`REQ-05-004`、`REQ-05-005`、`REQ-08-004`、`REQ-08-023`、`REQ-09-004`、`REQ-09-023`、`REQ-09-024`、`REQ-09-025`，重点是真实 no-login/no-paid/no-Level2/no-high-frequency provider 入库/readback 与真实 UI/API/SQLite 动态概率/假设跟踪验收 |

进度真源：`openspec/PROGRESS.md`。P19–P88 已交付并归档；当前无活跃 change，下一建议阶段为 P89。P75/P77/P78/P79/P80/P81/P82/P83/P84/P85/P87/P86/P88 形成 scoped progress 结论时，不能作为全量原始需求真实通过声明。

## 2. 开发原则

| 原则 | 要求 |
| --- | --- |
| 契约优先 | 数据表、DTO、枚举、状态流转以现有文档为准。 |
| 规则优先 | DeepSeek 只提供分析材料，不能写最终裁决。 |
| 证据优先 | 没有有效证据时展示信息不足状态，暂停交易类建议。 |
| 本地优先 | 默认本地 SQLite、VecLite、本地 Web 控制台，不暴露公网服务。 |
| 审计优先 | 关键动作必须写 `audit_events`。 |
| 用户确认优先 | 系统建议不会改变账户状态，只有 `executed_manually` 会更新本地账户。 |

## 3. 阶段总览

| 阶段 | 目标 | 主要产物 | 验收方式 |
| --- | --- | --- | --- |
| P0 工程骨架 | 建立可启动的 Go 后端与 React 前端 | 目录、启动命令、配置读取、健康检查 | 后端与前端可本地启动 |
| P1 数据底座 | 建立 SQLite 表、Repository、种子数据 | migration、Repository、seed | 数据写读测试通过 |
| P2 领域规则 | 实现核心枚举、状态机、规则裁决 | domain model、rule engine | 单元测试覆盖主要裁决 |
| P3 工作流 | 实现 Eino Graph 与审计事件 | Daily / Consultation / Evidence / Evolution / Gatekeeper | 工作流测试通过 |
| P3-foundation 基础治理 | 在 P4 前统一错误、ID、时间、事务、审计和测试策略 | `apperr`、`idgen`、`clock`、事务仓储、审计契约 | Go 测试和前端构建通过 |
| P4 HTTP API | 实现 `docs/api.md` 中核心 API | handler、DTO、错误码 | API 契约测试通过 |
| P5 前端驾驶舱 | 实现核心页面与状态展示 | Dashboard、Decision Detail、Evidence、Rules、Audit | 页面字段与前端契约一致 |
| P6 验收加固 | 完成端到端场景、错误状态和文档校验 | 测试报告、验收清单 | 关键场景全部通过 |
| P7 真实数据与分析底座 | 接入真实行情、情报、RAG/VecLite 与 DeepSeek 分析材料 | 数据源适配、RAG 索引、分析师节点、降级与审计 | 数据刷新、检索和分析降级测试通过 |
| P8 前端体验与测试 | 增强驾驶舱图表、交互体验、空态/错误态和前端测试 | 图表组件、页面交互、前端测试套件 | 前端测试与构建通过，无自动交易入口 |
| P9 复盘自动化与交付 | 实现本地任务入口、月度/季度复盘、规则有效性评估和交付说明 | `cmd/agent`、周期复盘、交付文档 | 本地任务可手动触发，复盘与审计可追踪 |
| P10 产品完成度 | 完成 P0-P9 后的产品级补齐与交付校验 | 节点级工作流、前端操作入口、数据源配置、索引降级 | 产品完成度规格归档，核心测试通过 |
| P11 治理与阶段重置 | 清理 P10 后治理状态，建立 P11-P18 执行门槛 | OpenSpec 状态、进度表、治理说明 | 无非预期活跃 change，归档前复审规则明确 |
| P12 最小真实只读数据源 | 接入最小可用行情与情报 provider | 只读 provider、降级、审计 | stub 继续可用，非 stub 可写入 SQLite facts |
| P13 索引与检索加固 | 稳定本地索引健康、重建与降级 | JSON 文件索引健康、重建统计、检索降级 | 索引缺失/损坏/不兼容可解释与恢复 |
| P14 Gatekeeper 节点图 | 将 GatekeeperAuditGraph 对齐节点级 Eino Graph | 节点注册、审计、样本/冲突/规则检查 | 节点名与文档一致，拒绝条件可审计 |
| P15 证据质量增强 | 提升证据质量字段与 formal/background 边界 | 时效权重、独立信源计数、材料分层 | 少于两个 A/S 独立信源不得 satisfied |
| P16 前端运维与复盘视图 | 展示数据源、索引、复盘与跨页追踪状态 | 状态面板、review summary、追踪入口 | 空态/失败态/降级态/成功态测试覆盖 |
| P17 本地调度与运维文档 | 提供默认关闭的本地周期任务配置 | launchd/cron 示例、`cmd/agent` 帮助、审计说明 | 示例无交易能力，配置不绕过人工确认 |
| P18 复盘提案链路 | 补齐 EvolutionProposalGraph 与规则提案链路 | 复盘产物、ProposalDraft/Record、审计追踪 | 规则提案只进入审查队列，不自动应用 |
| P19 公开 HTTP 数据桥接 | 接入可配置公开 HTTP JSON 行情与情报来源，保留 fixture/stub fallback | `ConfiguredMarketDataSource`、`ConfiguredIntelligenceSource`、config 扩展、数据源解析与降级测试 | `cmd/agent --task market-refresh` 与 evidence refresh 可从自备 HTTP/fixture/stub 写入本地事实；真实外部公开源 collector 待 P25+ 验证和实现 |
| P20 A 股 ETF/基金证据源 | 覆盖 A 股 ETF/基金公告、披露、净值/行情相关公开证据 payload 形态 | 公开证据 payload 解析、中文信源等级映射、去重、多源验证 | 高等级独立信源 payload 可驱动 `source_verifications`，不足时明确降级；真实公告源采集器待 P25+ 验证和实现 |
| P21 应用内通知中心 | 在本地应用内展示工作流和运行状态通知 | `notifications` 表、Repository/Service/API、前端通知页、轮询与已读状态 | 用户可查看未读通知并标记已读；不发送邮件、短信、系统 Push、Webhook 或 WebSocket |
| P22 规则体系与提案增强 | 扩展规则提案覆盖面与提案可解释性 | before/after payload、source facts、impact scope、risk notes、重复提案抑制 | 提案仍需用户确认、守门人审计、最终确认，不自动应用 |
| P23 复盘深化 | 输出可追溯复盘摘要和运行状态 | attribution summaries、错误标签、缺证据主题、提案结果、降级 workflow、tracking links | 空数据返回 unknown/empty，不生成无来源结论 |
| P24 本地运行硬化 | 提供本地配置校验、SQLite 备份与安全恢复 | `--validate-config`、`--backup`、`--restore`、`--restore-confirm`、server startup validation | 备份使用一致性快照；恢复拒绝无确认或覆盖现有 DB |
| P25 真实公开数据源调研验证 | 验证公开权威源是否可安全稳定接入 | 数据源卡片、接口/页面请求验证、合规边界、字段标准化、补采与刷新策略 | 首轮验证结论已输出；可先接、需二次验证、暂不接的源都有明确结论；P26/P27 范围清楚 |
| P26 公告与证据源 collector | 接入已通过 P25 验证的公告和监管证据源 | 巨潮、深交所、证监会首批 collector；AMAC 行业统计/自律栏目暂缓为二线背景源；标准证据 payload、RAG 入库、source verification、审计和幂等去重 | httptest/fixture 验证首批公开 JSON 响应、失败降级、重复执行不重复写入事实；collector 默认只读低频，不触发交易或外部推送 |
| P27 基金净值与 ETF 市场数据 collector | 接入已通过 P25 验证的基金净值、ETF 与指数日频数据源 | 东方财富基金净值、历史净值、资产配置和基金档案基础 metadata collector；中证指数基础信息 collector 当前公开字段校准；样本/权重/估值文件扩展接口按候选 metadata 低频读取；默认关闭的 `market_collectors` 配置；`market_metrics_json` source metadata 写入 | 东方财富基金 ETF/fund symbol 可获取最近公开净值并写入 market snapshot；中证指数基础信息可按当前公开 `index-basic-info` shape 写入 metadata；真实源失败时不伪造 stub 行情；缺失估值分位时不伪造并进入信息不足 |
| P28 预期收益与动态卖出评估增强 | 增强 expected return 样本上下文、情景触发和动态卖出复核材料 | `expected_return_scenarios` 新增 sample window、screening condition、scenario trigger、sell evaluation、reassessment trigger；后端 DTO/持久化/API 与前端展示同步 | `<5`、`5–19`、`>=20` 样本门槛可解释；卖出评估只提示人工复核，不更新账户、不创建交易、不覆盖最终规则裁决 |
| P29 公开证据真实采集 smoke | 修复 P26 公开证据 collector 的真实接口可运行性和入库验收 | CNInfo 当前接口参数修复、no_data/source_unavailable/parse_error 诊断、临时 SQLite smoke、真实公开公告入库文档 | `public-evidence-refresh` 在显式真实源配置下写入 `intelligence_items`、`intelligence_summary`、`rag_chunks`、`source_verifications` 和 `audit_events`；无公告和源不可用可区分 |

## 4. 任务分解

### P0：工程骨架

#### P0.1 初始化 Go 后端

参考文档：`docs/architecture.md`。

创建或确认以下目录：

```text
cmd/server/main.go
cmd/agent/main.go
internal/domain/
internal/application/
internal/infrastructure/
pkg/logger/
pkg/httputil/
configs/config.example.yaml
```

任务：

- [x] 建立 Go module。
- [x] 增加 `cmd/server/main.go`，提供 HTTP 服务启动入口。
- [x] 增加 `/api/v1/health`。
- [x] 增加配置读取：服务端口、SQLite 路径、VecLite 路径、DeepSeek 配置、日志级别。
- [x] 增加统一 logger。
- [x] 编写启动说明到 `docs/configuration.md`。

验收：

```bash
go test ./...
go run ./cmd/server
curl http://localhost:8080/api/v1/health
```

期望：

```json
{"status":"ok"}
```

#### P0.2 初始化前端工程

参考文档：`docs/ui-design.md`、`docs/frontend-contract.md`。

创建或确认以下目录：

```text
web/
web/src/app/
web/src/pages/
web/src/components/
web/src/services/
web/src/types/
web/src/styles/
```

任务：

- [x] 使用 React + Vite + TypeScript 建立前端工程。
- [x] 增加基础路由：今日纪律、持仓、决策咨询、情报与证据、规则与纪律、复盘与审计、设置。
- [x] 增加 API client，统一处理 `request_id`、`data`、`error`。
- [x] 增加全局状态样式：正常、信息不足、冻结观察、高危。

验收：

```bash
cd web
npm install
npm run build
npm run dev
```

期望：前端可打开，并展示空驾驶舱骨架。

### P1：数据底座

#### P1.1 SQLite migration

参考文档：`docs/data-model.md`。

创建文件：

```text
internal/infrastructure/persistence/sqlite/migration/001_init.sql
internal/infrastructure/persistence/sqlite/migration/002_seed_rules.sql
internal/infrastructure/persistence/sqlite/migrate.go
```

必须包含的核心表：

- `portfolio_snapshots`
- `positions`
- `position_snapshots`
- `operation_confirmations`
- `position_transactions`
- `market_snapshots`
- `rule_versions`
- `decision_records`
- `intelligence_items`
- `intelligence_summary`
- `rag_chunks`
- `evidence_refs`
- `source_verifications`
- `capability_configs`
- `user_settings`
- `audit_events`
- `error_cases`
- `rule_proposals`
- `gatekeeper_audits`

任务：

- [x] 按 `docs/data-model.md` 创建字段和索引。
- [x] 对枚举字段使用 CHECK 约束。
- [x] 对历史表采用追加式写法，不提供物理删除路径。
- [x] 写入默认 `rule_versions`，版本号为 `v3.0`。
- [x] 写入默认信源等级配置。

验收：

```bash
go test ./internal/infrastructure/persistence/sqlite/...
```

期望：migration 可创建空库，重复启动不会破坏已有数据。

#### P1.2 Repository 层

参考文档：`docs/data-model.md`、`docs/workflow.md`。

创建文件：

```text
internal/domain/repository/portfolio_repo.go
internal/domain/repository/decision_repo.go
internal/domain/repository/intelligence_repo.go
internal/domain/repository/rule_repo.go
internal/domain/repository/audit_repo.go
internal/infrastructure/persistence/sqlite/portfolio_repo_impl.go
internal/infrastructure/persistence/sqlite/decision_repo_impl.go
internal/infrastructure/persistence/sqlite/intelligence_repo_impl.go
internal/infrastructure/persistence/sqlite/rule_repo_impl.go
internal/infrastructure/persistence/sqlite/audit_repo_impl.go
```

任务：

- [x] 实现账户快照写读。
- [x] 实现持仓当前态与持仓快照写读。
- [x] 实现决策记录与证据引用写读。
- [x] 实现用户确认写入。
- [x] 实现情报摘要、RAG 文本块、多源验证写读。
- [x] 实现规则版本、规则提案、守门人审计写读。
- [x] 实现审计事件追加写入。

验收：

```bash
go test ./internal/domain/repository/... ./internal/infrastructure/persistence/sqlite/...
```

期望：每个 Repository 都有写入、读取、事务失败回滚测试。

### P2：领域规则

#### P2.1 核心模型与枚举

参考文档：`docs/api.md` 第 4 节、`docs/data-model.md` 第 3 节。

创建文件：

```text
internal/domain/model/enums.go
internal/domain/model/portfolio.go
internal/domain/model/market.go
internal/domain/model/evidence.go
internal/domain/model/decision.go
internal/domain/model/rule.go
internal/domain/model/audit.go
```

任务：

- [x] 定义 dashboard_state、workflow_status、position_state、verification_status。
- [x] 定义 confirmation_status、confirmation_type、final_verdict.status。
- [x] 定义 rule_proposal.status、audit_result、audit.action、audit.status。
- [x] 定义 WorkflowContext 对应领域结构。
- [x] 编写枚举合法性测试。

验收：

```bash
go test ./internal/domain/model/...
```

#### P2.2 规则裁决引擎

参考文档：`docs/requirements.md` 第 2、6 节，`docs/workflow.md` 第 5 节。

创建文件：

```text
internal/domain/rule/rules_engine.go
internal/domain/rule/source_policy.go
internal/domain/rule/capability_policy.go
internal/domain/rule/risk_policy.go
internal/domain/rule/expectation_engine.go
internal/domain/rule/gatekeeper_logic.go
```

任务：

- [x] 实现能力圈外时返回 `rejected`，并拒绝交易类分析。
- [x] 实现证据不足时返回 `insufficient_data`。
- [x] 实现普通正式证据允许 S/A/B 级来源，C 级只能作为 `background`。
- [x] 实现重大利好、重大利空、买入逻辑破坏必须至少 2 个 A 或 S 级独立信源，不满足时返回 `frozen_watch`。
- [x] 实现买入逻辑破坏时返回 `sell_only`。
- [x] 实现情绪极端时暂停主动交易建议。
- [x] 实现 PE/PB 分位区间规则：高危区、观察区、舒适区、低估区分别对应禁止买入、仅持有、按计划定投、分批补仓。
- [x] 实现移动止盈规则：浮盈 20% 时 `optional_actions` 可包含卖出 30% 与启动移动止盈，浮盈 30% 时 `optional_actions` 可包含再卖出 30%，剩余仓位按 10% 回撤触发减仓或卖出评估。
- [x] 实现 R-5 现金冗余规则：现金仓位低于 5% 时限制新增买入，现金仓位 5%-10% 为正常冗余区间。
- [x] 实现核心-卫星仓位规则：核心资产目标 60%-70%，卫星资产目标 20%-30%，偏离 ±15% 或卫星超上限时提示再平衡。
- [x] 实现预期收益评估：输出上行情景、基准情景、下行情景及置信度，只作为分析材料，不覆盖最终裁决。
- [x] 实现 C 级信源只能作为 `background` 材料。
- [x] 实现规则提案完整状态机：生成、送审、放弃、守门人通过、守门人拒绝、需要用户复核、最终确认、最终拒绝、已拒绝/已应用后再次操作。

验收：

```bash
go test ./internal/domain/rule/...
```

期望：能力圈、证据不足、多源验证、普通证据 S/A/B、重大事件 2 个 A 或 S、C 级 background、PE/PB 分位区间、移动止盈、现金冗余、核心-卫星仓位、规则提案完整状态机均有测试。

规则测试断言至少覆盖：

| 场景 | 输入条件 | 期望裁决 / 输出 |
| --- | --- | --- |
| 高危估值 | `pe_percentile>80` 或 `pb_percentile>80` | `prohibited_actions` 包含新增买入，`triggered_rules` 包含估值高危 |
| 观察区 | `pe_percentile` 或 `pb_percentile` 位于 50%-80% | `final_verdict.status=hold`，只允许继续观察 |
| 舒适区 | PE/PB 分位 30%-50%，买入逻辑完好，证据满足 | 可进入按计划定投建议，但仍受仓位与现金规则约束 |
| 低估区 | PE/PB 分位 <30%，买入逻辑完好，情绪非极端 | 可进入分批配置建议，`optional_actions` 不得包含立即满仓 |
| PE/PB 冲突 | PE 低估但 PB 高危，或任一核心估值指标高危 | 高风险指标优先，禁止新增买入 |
| 浮盈 20% | 持仓浮盈达到 20%，此前未记录同阶段止盈建议 | `optional_actions` 包含卖出 30% 与启动移动止盈 |
| 浮盈 30% | 持仓浮盈达到 30%，20% 阶段已处理 | `optional_actions` 包含再卖出 30%，不得重复生成 20% 阶段建议 |
| 回撤 10% | 已启动移动止盈，当前价较阶段高点回撤 10% | `final_verdict.status=reduce` 或 `sell_only`，触发移动止盈规则 |
| 现金不足 | `cash_ratio<0.05` | `prohibited_actions` 包含新增买入，触发 R-5 |
| 现金正常 | `0.05<=cash_ratio<=0.10` | 不因现金规则单独禁止交易 |
| 卫星超上限 | `satellite_ratio>0.30` 或偏离目标超过 15% | 优先提示再平衡，止盈资金优先回归核心资产 |
| 核心低于目标 | `core_ratio<0.60` 且有可配置资金 | `optional_actions` 优先包含提高核心资产占比 |
| 放弃送审 | `pending_user_confirm` + `confirm=false` | 提案状态变为 `rejected`，不写 `gatekeeper_audits`、不写 `rule_versions` |
| 样本不足提案 | `sample_count<3` | 只能生成 `draft` 或停留在 `pending_user_confirm`；除非后续 EvolutionGraph 或受控内部任务生成满足样本条件的新提案版本，否则不得进入 `pending_final_confirm` 或 `applied` |
| 守门人通过 | `under_gatekeeper_audit` + `approved` | 提案状态变为 `pending_final_confirm`，不写 `rule_versions` |
| 守门人拒绝 | `under_gatekeeper_audit` + `rejected` | 提案状态变为 `rejected`，不写 `rule_versions` |
| 需要用户复核 | `under_gatekeeper_audit` + `needs_user_review` | 提案回到 `pending_user_confirm`，写 `gatekeeper_audits` 与 `audit_events`；用户可放弃或重新送审，如需修改只能由后续 EvolutionGraph 或受控内部任务生成新提案版本 |
| 最终确认 | `pending_final_confirm` + `confirm=true` + `sample_count>=3` | 提案状态变为 `applied`，创建新 active `rule_versions`，旧 active 归档 |
| 样本不足最终确认 | `pending_final_confirm` + `confirm=true` + `sample_count<3` | 返回 `BAD_REQUEST`，不写 `rule_versions`，提案状态不变 |
| 最终拒绝 | `pending_final_confirm` + `confirm=false` | 提案状态变为 `rejected`，不写 `rule_versions` |
| 终态重复操作 | `rejected` 或 `applied` 后再次确认 | 返回 `BAD_REQUEST`，不改变提案和规则版本 |

### P3：Eino 工作流

#### P3.1 WorkflowContext 与节点框架

参考文档：`docs/workflow.md`。

创建文件：

```text
internal/application/workflow/context.go
internal/application/workflow/node.go
internal/application/workflow/audit_writer.go
internal/application/workflow/errors.go
```

任务：

- [x] 定义 WorkflowContext。
- [x] 定义节点输入输出约定。
- [x] 每个节点必须返回状态、错误码、审计事件片段。
- [x] 工作流节点审计必须填写 `action`、`node_name`、`node_action`、`status`、`input_ref_type/input_ref`；产生输出时填写 `output_ref_type/output_ref`。
- [x] `status=failed` 时必须填写 `error_code`；降级有明确原因时填写 `error_code`。
- [x] 失败节点必须写入 `audit_events`。

验收：

```bash
go test ./internal/application/workflow/...
```

#### P3.2 DailyDisciplineGraph 与 ConsultationGraph

参考文档：`docs/workflow.md` 第 2、3、5、6 节。

当前实现文件：

```text
internal/application/workflow/eino_graph.go
internal/application/workflow/steps.go
internal/application/workflow/dependencies.go
internal/application/workflow/expected_return.go
```

任务：

- [x] 每日纪律读取账户、持仓、市场、证据和规则版本。
- [x] 主动咨询必须包含能力圈检查。
- [x] DeepSeek 节点只写分析报告。
- [x] ExpectedReturnNode 按 `sample_count` 映射 `precision_status`：`>=20 available`、`5~19 insufficient`、`<5 unavailable`。
- [x] `expected_return_scenarios` DTO 符合三种状态约束：available 可返回概率，insufficient 不返回精确概率且写样本不足说明，unavailable 返回空 `scenarios` 且写定性原因。
- [x] RuleArbitrationNode 生成最终裁决。
- [x] DecisionRecordNode 写 `decision_records`、`evidence_refs`、`audit_events`，并把预期收益情景保存到 `decision_records.expected_return_scenarios_json`。

验收：

```bash
go test ./internal/application/workflow/... -run 'Daily|Consultation'
```

期望：正常、信息不足、能力圈外、LLM 不可用、预期收益样本不足五类场景均通过。

#### P3.3 Evidence、Evolution、Gatekeeper 工作流

参考文档：`docs/workflow.md`、`docs/data-model.md` 第 5、6 节。

创建文件：

```text
internal/application/workflow/evidence_verification_graph.go
internal/application/workflow/market_refresh_graph.go
internal/application/workflow/evolution_proposal_graph.go
internal/application/workflow/gatekeeper_audit_graph.go
```

任务：

- [x] 情报写入 `intelligence_items`、`intelligence_summary`、`rag_chunks`。
- [x] VecLite 索引从 `rag_chunks` 构建，可由 SQLite 重建。
- [x] 多源验证写 `source_verifications`。
- [x] 市场刷新实现为独立 MarketRefreshGraph：读取数据源、标准化市场状态、写入 `market_snapshots` 与 `audit_events`。
- [x] 错误案例生成规则提案，但不改正式规则。
- [x] 守门人审计只生成 `gatekeeper_audits`。
- [x] 审计通过后状态为 `pending_final_confirm`，不写正式规则。
- [x] 用户最终确认后才写 `rule_versions`。

验收：

```bash
go test ./internal/application/workflow/... -run 'Evidence|Market|Evolution|Gatekeeper'
```

### P3-foundation：基础治理

参考文档：`docs/architecture.md`、`docs/data-model.md`、`docs/workflow.md`、`docs/api.md`、`docs/frontend-contract.md`。

目标：在 P4 HTTP API 前统一横向基础能力，避免 API、前端和后续集成建立在分散错误、ID、时间、事务和审计规则之上。

任务：

- [x] 建立 `internal/pkg/apperr`，统一错误码、分类、HTTP 映射和审计错误码映射。
- [x] 建立 `internal/pkg/idgen` 与 `internal/pkg/clock`，统一关键实体 ID 与 UTC/RFC3339 时间。
- [x] 调整跨表仓储事务边界，覆盖决策记录、证据事实、用户确认、守门人审计和规则应用。
- [x] 集中工作流审计 action、node_name、node_action 与 input/output ref 映射。
- [x] 补充基础包、仓储事务、工作流分支、HTTP 错误信封和前端错误状态映射验证。

验收：

```bash
go test ./...
cd web && npm run build
```

期望：所有测试与构建通过，且 P4/P5 可直接复用统一错误信封和前端错误状态映射。

### P4：HTTP API

#### P4.1 API DTO 与错误处理

参考文档：`docs/api.md`、`docs/frontend-contract.md`。

创建文件：

```text
internal/application/dto/common.go
internal/application/dto/dashboard.go
internal/application/dto/decision.go
internal/application/dto/evidence.go
internal/application/dto/rule.go
internal/application/dto/audit.go
internal/application/handler/errors.go
pkg/httputil/response.go
```

任务：

- [x] 所有响应包含 `request_id`。
- [x] 错误响应符合 `docs/api.md` 第 2、3 节。
- [x] DTO 字段名与 `docs/frontend-contract.md` 一致。
- [x] `EVIDENCE_NOT_FOUND` 返回 409，并让前端显示信息不足状态。

验收：

```bash
go test ./internal/application/dto/... ./internal/application/handler/... ./pkg/httputil/...
```

#### P4.2 核心 API

创建文件：

```text
internal/application/handler/dashboard_handler.go
internal/application/handler/decision_handler.go
internal/application/handler/portfolio_handler.go
internal/application/handler/evidence_handler.go
internal/application/handler/rule_handler.go
internal/application/handler/audit_handler.go
internal/application/handler/settings_handler.go
internal/application/handler/market_handler.go
internal/application/handler/review_handler.go
```

任务：

- [x] `dashboard_handler.go`：`GET /api/v1/dashboard/today`。
- [x] `decision_handler.go`：`POST /api/v1/decisions/consult`、`GET /api/v1/decisions/{decision_id}`、`GET /api/v1/decisions`、`POST /api/v1/decisions/{decision_id}/confirmations`。
- [x] `portfolio_handler.go`：`POST /api/v1/portfolio/init`、`GET /api/v1/portfolio/current`、`POST /api/v1/portfolio/adjustments`。
- [x] `evidence_handler.go`：`POST /api/v1/evidence/refresh`、`GET /api/v1/evidence`、`GET /api/v1/evidence/verification`、`POST /api/v1/evidence/rebuild-index`。
- [x] `market_handler.go`：`POST /api/v1/market/refresh`、`GET /api/v1/market/snapshots/latest`。
- [x] `rule_handler.go`：`GET /api/v1/rules/current`、`GET /api/v1/rule-proposals`、`POST /api/v1/rule-proposals/{proposal_id}/confirm`、`POST /api/v1/rule-proposals/{proposal_id}/final-confirm`。
- [x] `settings_handler.go`：`GET /api/v1/settings/system`、`PUT /api/v1/settings`、`GET /api/v1/settings/capability`、`PUT /api/v1/settings/capability`。
- [x] `audit_handler.go`：`GET /api/v1/audit-events`。
- [x] `review_handler.go`：`GET /api/v1/review/summary`。

验收断言：

- [x] `executed_manually` 成功后同时写入 6 类记录：`operation_confirmations`、`position_transactions`、`positions`、`portfolio_snapshots`、`position_snapshots`、`audit_events`。
- [x] `marked_error` 成功后同一事务写入 `operation_confirmations`、`error_cases`、`audit_events`，响应返回 `error_case_id`。
- [x] `planned` 与 `watch` 不写 `position_transactions`，不新增账户快照。
- [x] `planned` 与 `watch` 可互相转换，并可升级为 `executed_manually` 或 `marked_error`；每次成功转换都创建新的 `operation_confirmations`。
- [x] `executed_manually` 与 `marked_error` 是确认终态，再次确认返回 `BAD_REQUEST`，不得重复写账户快照、交易流水或错误案例。
- [x] `record_type!=formal_trade_advice` 或 `confirmation_status=not_required` 时，确认接口返回 `BAD_REQUEST`。
- [x] 守门人审计通过只进入 `pending_final_confirm`，不写 `rule_versions`。
- [x] 最终确认应用规则后创建新 active `rule_versions`，旧 active 归档。
- [x] `POST /api/v1/evidence/refresh` 同步完成情报采集、摘要、`source_verifications` 写入和索引更新；索引失败不得回滚 SQLite 事实数据。
- [x] 市场状态枚举统一使用 `liquidity_state=normal/warning/danger`、`sentiment_state=cold/neutral/hot/extreme`。
- [x] `PUT /api/v1/settings` 只能保存通知、页面偏好、普通数据源；`PUT /api/v1/settings/capability` 只保存能力圈；规则阈值、裁决优先级和 SOP 必须生成 `rule_proposals`。
- [x] `POST /api/v1/market/refresh` 覆盖全部成功、部分成功、全部失败、快照写入失败四类响应；部分成功时返回 200，并在 `failed_symbols` 中写明失败标的与原因。
- [x] `sample_count<3` 的规则提案调用送审接口返回 `BAD_REQUEST`，不得写 `gatekeeper_audits` 或进入 `under_gatekeeper_audit`。
- [x] `sample_count<3` 的规则提案即使状态异常进入 `pending_final_confirm`，最终确认接口也必须返回 `BAD_REQUEST`，不得写 `rule_versions`。

验收：

```bash
go test ./internal/application/handler/...
```

期望：API 契约测试覆盖成功响应、错误响应、状态流转、事务写入。

### P5：前端驾驶舱

#### P5.0 架构治理准备

参考文档：`docs/architecture.md`、`docs/data-model.md`、`docs/workflow.md`、`docs/frontend-contract.md`。

任务：

- [x] 应用层和工作流包只依赖仓储接口与事务协调接口，不依赖 SQLite 具体实现。
- [x] HTTP handler 保持请求解析、用例调用和响应写入职责，不直接访问数据库或管理 SQL 事务。
- [x] 跨表业务事实和对应 `audit_events` 使用统一事务协调路径。
- [x] 关键业务 ID、时间和契约枚举复用共享实现，并支持确定性测试。
- [x] 前端按 feature 和 shared 目录准备 P5 页面扩展结构。

验收：

```bash
go test ./...
cd web && npm run build
```

#### P5.1 类型与 API client

参考文档：`docs/frontend-contract.md`、`docs/api.md`。

创建文件：

```text
web/src/types/api.ts
web/src/types/dashboard.ts
web/src/types/portfolio.ts
web/src/types/decision.ts
web/src/types/evidence.ts
web/src/types/rule.ts
web/src/types/audit.ts
web/src/types/settings.ts
web/src/types/market.ts
web/src/types/review.ts
web/src/services/client.ts
web/src/services/dashboard.ts
web/src/services/portfolio.ts
web/src/services/decision.ts
web/src/services/evidence.ts
web/src/services/rule.ts
web/src/services/audit.ts
web/src/services/settings.ts
web/src/services/market.ts
web/src/services/review.ts
```

任务：

- [x] 定义通用响应类型。
- [x] 定义驾驶舱、持仓、决策、证据、市场、规则、审计、设置、复盘 DTO。
- [x] 统一处理 409、500、503 错误。
- [x] 前端不访问 SQLite、VecLite、本地文件。

验收：

```bash
cd web
npm run build
```

#### P5.2 Agent 决策驾驶舱

参考文档：`docs/ui-design.md`、`docs/ui-flow.md`。

创建文件：

```text
web/src/pages/DashboardPage.tsx
web/src/components/layout/CockpitLayout.tsx
web/src/components/dashboard/DisciplineStatus.tsx
web/src/components/dashboard/PortfolioSummary.tsx
web/src/components/dashboard/TriggeredRules.tsx
web/src/components/dashboard/EvidenceSummary.tsx
web/src/components/dashboard/FinalVerdictCard.tsx
web/src/components/dashboard/UserConfirmationPanel.tsx
```

任务：

- [x] 实现三栏驾驶舱布局。
- [x] 首屏展示纪律状态、风险红线、今日建议、账户摘要、证据摘要。
- [x] 信息不足状态展示缺失项和暂停原因。
- [x] 冻结观察状态展示等待条件。
- [x] 用户确认区只允许记录计划、已手动执行、待观察、标记错误。
- [x] 页面不得出现自动交易入口。

验收：

```bash
cd web
npm run build
```

#### P5.3 决策详情、证据、规则与审计页面

创建文件：

```text
web/src/pages/DecisionDetailPage.tsx
web/src/pages/EvidencePage.tsx
web/src/pages/RulesPage.tsx
web/src/pages/AuditPage.tsx
web/src/pages/PortfolioPage.tsx
web/src/pages/SettingsPage.tsx
web/src/pages/ReviewSummaryPage.tsx
web/src/components/decision/DecisionTrace.tsx
web/src/components/evidence/EvidenceTable.tsx
web/src/components/rules/RuleProposalPanel.tsx
web/src/components/audit/AuditEventTimeline.tsx
web/src/components/portfolio/PortfolioTable.tsx
web/src/components/settings/CapabilitySettingsPanel.tsx
web/src/components/review/ReviewSummaryPanel.tsx
```

任务：

- [x] 决策详情按 `ui-flow.md` 第 6 节展示。
- [x] 持仓页使用 `GET /api/v1/portfolio/current`，不直接访问 SQLite 或 VecLite。
- [x] 证据页展示 source_level、evidence_role、verification_status。
- [x] 规则提案页展示 `pending_final_confirm` 状态和最终确认动作。
- [x] 设置页展示能力圈配置、系统状态、市场快照状态、通知配置和索引状态，不展示完整密钥。
- [x] 复盘页展示建议数量、确认动作、错误案例、规则提案和审计事件汇总。
- [x] 审计页区分 `action`、`node_name` 与 `node_action`，并按 `status`、`error_code`、输入引用、输出引用展示审计详情。

验收：

```bash
cd web
npm run build
```

### P6：验收加固

#### P6.1 端到端场景

创建文件：

```text
docs/testing-plan.md
```

必须覆盖 `docs/functional-spec.md` 的 A01-A17 可测试验收断言：

- [x] A01 首次使用：无账户数据，展示引导，且不创建 `decision_records`。
- [x] A02 正常每日纪律：生成建议、证据、审计事件。
- [x] A03 证据不足：返回 `EVIDENCE_NOT_FOUND`，暂停交易类建议。
- [x] A04 VecLite 不可用：SQLite 摘要充足时降级展示；不足时信息不足。
- [x] A05 能力圈外：拒绝交易类分析。
- [x] A06 用户记录计划：写 `operation_confirmations`，不更新账户。
- [x] A07 用户记录已手动执行：写 `operation_confirmations`、`position_transactions`、`positions`、`portfolio_snapshots`、`position_snapshots`、`audit_events`。
- [x] A08 已手动执行失败：事务回滚，不留下部分确认记录。
- [x] A09 用户标记错误：写 `operation_confirmations`、`error_cases`、`audit_events`，返回 `error_case_id`。
- [x] A10 C 级信源：只能作为 `background`，不得进入正式裁决证据。
- [x] A11 LLM 不可用：规则引擎降级裁决，`workflow_status=degraded`。
- [x] A12 守门人审计通过：进入 `pending_final_confirm`，不写 `rule_versions`；`sample_count<3` 的提案不得进入守门人审计，接口返回 `BAD_REQUEST`。
- [x] A13 规则最终确认：创建新 active `rule_versions`，旧 active 归档；`sample_count<3` 的提案不得最终确认，接口返回 `BAD_REQUEST` 且不写 `rule_versions`。
- [x] A14 审计事件：前端区分 `action`、`node_name` 和 `node_action`。
- [x] A15 禁止自动交易：不存在交易执行接口，前端无一键交易入口。
- [x] A16 市场数据刷新：全部成功新增快照且审计成功；部分成功返回 200、写 `failed_symbols` 且审计降级；全部失败返回 `DATA_SOURCE_UNAVAILABLE` 或 `DATA_STALE`；写入失败返回 `MARKET_SNAPSHOT_WRITE_FAILED` 且市场快照事务回滚。
- [x] A17 预期收益评估：只展示情景概率，不覆盖最终规则裁决，不承诺收益；`available` 必含 upside/base/downside 且可返回概率；`insufficient` 不返回精确概率，`probability=null` 且写样本不足说明；`unavailable` 返回 `scenarios=[]` 且写定性原因。

验收：

```bash
go test ./...
cd web && npm run build && npm test
```

#### P6.2 配置与启动文档

创建文件：

```text
docs/configuration.md
docs/migration-plan.md
```

必须包含：

- [x] SQLite 数据文件路径。
- [x] VecLite 索引文件路径。
- [x] DeepSeek API Key 环境变量。
- [x] 数据源开关。
- [x] 日志级别。
- [x] migration 执行方式。
- [x] seed 数据说明。
- [x] 本地启动命令。

### P7：真实数据与分析底座

P7 目标是把 P0-P6 的骨架接入可替换的真实数据与分析基础能力。外部依赖缺失时必须降级为信息不足或定性说明，不得伪造证据。

#### P7.1 真实行情与情报数据源

参考文档：`docs/requirements.md`、`docs/architecture.md`、`docs/data-model.md`、`docs/workflow.md`。

任务：

- [x] 增加行情数据源适配层，支持按配置启用真实数据源或本地 stub。
- [x] 增加情报数据源适配层，支持新闻、公告或手工导入数据进入 `intelligence_items`。
- [x] 为市场快照刷新写入 `market_snapshots` 与 `audit_events`。
- [x] 对部分失败、全部失败、数据过期分别返回既有错误或降级状态。
- [x] 在配置文档中说明数据源开关、凭证环境变量和本地 stub 用法，不写真实密钥。

验收：

```bash
go test ./internal/infrastructure/... ./internal/application/...
```

期望：真实数据源可替换为 stub；失败场景可审计；不会因外部依赖缺失而阻塞本地开发。

#### P7.2 RAG/VecLite 检索与索引

参考文档：`docs/requirements.md`、`docs/architecture.md`、`docs/data-model.md`。

任务：

- [x] 实现 VecLite 索引读写适配，索引路径来自配置。
- [x] 将 `rag_chunks` 与 `intelligence_summary` 纳入检索构建流程。
- [x] 支持从 SQLite 文本块重建 VecLite 索引。
- [x] VecLite 不可用时按既有约定降级到 SQLite 摘要或信息不足。
- [x] 记录检索输入、命中证据和降级原因到审计事件或可追踪上下文。

验收：

```bash
go test ./internal/infrastructure/... ./internal/application/workflow/...
```

期望：检索可用、可重建、可降级；C 级信源仍不得作为正式裁决依据。

#### P7.3 DeepSeek 分析师材料

参考文档：`docs/requirements.md`、`docs/architecture.md`、`docs/workflow.md`。

任务：

- [x] 增加 DeepSeek 客户端封装，API Key 从环境变量读取。
- [x] 将价值分析、趋势风险和预期收益节点从占位实现改为可调用分析服务。
- [x] 明确 prompt 输入只包含允许使用的证据、持仓上下文和规则边界。
- [x] 解析 DeepSeek 输出为 `analyst_reports` 或等价结构，不写最终裁决。
- [x] LLM 不可用、超时或输出不可解析时，工作流进入降级状态，并由规则引擎继续生成最终裁决。
- [x] 对非显然的 prompt 约束、降级和审计逻辑添加中文注释。

验收：

```bash
go test ./internal/application/workflow/... ./internal/infrastructure/...
```

期望：DeepSeek 仅提供分析材料；最终裁决仍由规则引擎负责；LLM 故障不产生自动交易动作。

### P8：前端体验与测试

P8 目标是在既有 API 与前端契约基础上增强可读性和可验证性，不改变产品边界。

#### P8.1 驾驶舱图表与关键交互

参考文档：`docs/ui-design.md`、`docs/ui-flow.md`、`docs/frontend-contract.md`。

任务：

- [x] 在今日纪律、持仓、复盘页面增加图表组件，用于展示仓位、风险、证据覆盖和复盘摘要。
- [x] 图表数据只来自 API DTO，不直接读取 SQLite、VecLite 或本地文件。
- [x] 增强证据、决策链、审计时间线的筛选和展开交互。
- [x] 对信息不足、数据过期、LLM 降级、VecLite 不可用等状态提供明确空态和错误态。
- [x] 保持用户确认区为线下动作记录，不增加一键交易或自动下单入口。

验收：

```bash
cd web
npm run build
```

期望：核心页面可展示图表和关键状态；构建通过；无自动交易入口。

#### P8.2 前端测试与契约校验

参考文档：`docs/frontend-contract.md`、`docs/testing-plan.md`。

任务：

- [x] 建立前端测试脚本，覆盖 API client、关键状态渲染和用户确认流程。
- [x] 覆盖信息不足、冻结观察、高危、降级、错误响应等 UI 状态。
- [x] 覆盖规则提案最终确认流程，确认 `pending_final_confirm` 可见且不会自动应用规则。
- [x] 覆盖禁止自动交易入口的页面断言。
- [x] 将前端测试命令写入本阶段验收。

验收：

```bash
cd web
npm run build
npm test
```

期望：前端测试和构建通过；页面字段继续符合 `docs/frontend-contract.md`。

### P9：复盘自动化与交付

P9 目标是把每日任务、周期复盘、规则有效性评估和本地交付流程串联起来，形成可长期使用的本地系统。

#### P9.1 `cmd/agent` 本地任务入口

参考文档：`docs/requirements.md`、`docs/architecture.md`、`docs/workflow.md`。

任务：

- [x] 实现 `cmd/agent`，提供每日纪律、市场刷新、情报索引、复盘任务的本地触发入口。
- [x] 支持手动触发和本地调度配置，默认不启用任何自动交易能力。
- [x] 每次任务执行写入 `audit_events`，记录输入摘要、状态和错误码。
- [x] 任务失败时返回可读错误，并保留已有数据一致性。
- [x] 对调度、安全边界和审计写入添加必要中文注释。

验收：

```bash
go test ./...
go run ./cmd/agent --help
```

期望：本地任务入口可启动；任务行为可追踪；不会执行交易。

#### P9.2 月度/季度复盘与规则有效性评估

参考文档：`docs/requirements.md`、`docs/functional-spec.md`、`docs/workflow.md`。

任务：

- [x] 完善月度复盘，汇总确认动作、错误案例、规则提案和审计事件。
- [x] 完善季度复盘，评估规则命中、误判、缺证据和降级情况。
- [x] 将规则有效性评估结果写入规则提案或复盘摘要，仍需守门人审计和用户最终确认。
- [x] 前端复盘页展示周期摘要、规则建议和追踪入口。
- [x] 保持规则提案不会自动应用。

验收：

```bash
go test ./internal/application/workflow/... ./internal/application/handler/...
cd web && npm run build
```

期望：月度/季度复盘可生成摘要；规则变更仍经审计和用户确认。

#### P9.3 本地交付与运维说明

参考文档：`docs/configuration.md`、`docs/migration-plan.md`、`docs/testing-plan.md`。

任务：

- [x] 补充本地启动、初始化、数据备份、索引重建和恢复说明。
- [x] 补充常见故障处理：数据源不可用、VecLite 索引损坏、DeepSeek 缺配置、SQLite 写入失败。
- [x] 补充 P7-P9 后的完整验收命令。
- [x] 确认文档不包含真实密钥、账号、token 或个人敏感信息。

验收：

```bash
go test ./...
cd web && npm run build
```

期望：本地交付说明完整；关键故障有处理路径；安全边界清晰。

### P19–P24：可用 MVP 交付摘要

P19–P24 已完成并提交。以下内容作为交付摘要和下一轮规划基线保留；不要再按旧任务清单重复规划。

#### P19 公开 HTTP 数据桥接

已交付：

- [x] 行情数据源支持可配置公开 HTTP endpoint、fixture 与 stub fallback。
- [x] 市场 payload 支持标准字段、`data`/`market`/`quote` 包装、对象或数组、`close_price`/`close`/`nav`/`net_value` 价格字段。
- [x] 情报数据源支持公开 HTTP endpoint、`items`/`data`/顶层数组、公开字段名 `source`/`content`/`summary`/`original_url`。
- [x] A 股 ETF/基金公开信源支持中文来源等级推断，交易所、基金公司公告、巨潮资讯、上交所、深交所等映射为高等级来源。
- [x] HTTP 失败、非 2xx、解析失败、缺价格、过期数据返回稳定错误，并保留 fixture/stub 离线能力。

当前范围：公开 HTTP 数据源优先；不引入 Python/AKShare 代理作为当前实现方案；不接券商交易 API；不引入付费、登录或浏览器爬虫数据源。

#### P20 A 股 ETF/基金证据源增强

已交付：

- [x] 公开证据 payload 解析、去重和多源验证路径可写入 SQLite 事实。
- [x] URL 优先去重；无 URL 时按来源、标题和发布时间去重，保留不同来源的同题同日独立证据。
- [x] 证据刷新时即使向量索引失败，也保留 SQLite 摘要、文本块和审计事实。
- [x] 高等级独立信源不足时保持信息不足或降级状态，不伪造结论。

当前范围：优先支持 A 股 ETF、基金公告、披露、净值/行情相关公开证据 payload 形态；真实公告源采集器、历史补采和 smoke 验证待 P25+ 验证后实现；不扩展到个股、港股、美股或登录源。

#### P21 应用内通知中心

已交付：

- [x] 新增 `notifications` 本地表、repository、service、handler 和 DTO。
- [x] 新增 `GET /api/v1/notifications`、`POST /api/v1/notifications/{notification_id}/read`、`POST /api/v1/notifications/read-all`。
- [x] 前端新增通知中心页面，展示未读数、通知列表、单条已读、全部已读，并通过本地 API 轮询刷新。
- [x] 市场数据源失败、部分失败、证据索引失败、手动重建索引失败、规则提案待确认、复盘降级或缺证据等场景写入应用内通知。
- [x] 同一 `type/source_type/source_id` 的未读通知去重并刷新内容，避免重复告警刷屏。

当前范围：只做应用内通知；不提供邮件、短信、系统 Push、Bark、Webhook 或 WebSocket。通知动作不执行交易，不自动应用规则。

#### P22 规则体系与提案增强

已交付：

- [x] 规则提案包含 before/after payload、source facts、impact scope、risk notes、样本和原因。
- [x] 提案生成保留用户确认、守门人审计和最终确认流程。
- [x] 相同规则/事实窗口避免重复提案。
- [x] 前端展示增强后的提案信息，并保留最终确认要求。

当前范围：规则提案更可解释，但不会自动应用。

#### P23 复盘深化

已交付：

- [x] 复盘 DTO 增加 `attribution_summaries`、`recurring_error_tags`、`missing_evidence_themes`、`rule_proposal_outcomes`、`degraded_workflows`、`ops_status`、`tracking_links` 等可追溯字段。
- [x] 复盘结论来自本地决策、确认、行情、证据、规则提案、错误案例和审计事件。
- [x] 空窗口返回 `empty`/`unknown` 状态，不生成无来源结论。
- [x] 复盘降级或缺证据时写入应用内通知。

当前范围：复盘只聚合本地事实，不做无来源推断；情绪日志、基准对比等更深功能可作为下一轮独立需求重新评估。

#### P24 本地运行硬化

已交付：

- [x] 配置校验覆盖 SQLite path、vector path、数据源 endpoint、通知设置等本地运行前提。
- [x] `cmd/server` 启动前执行配置校验。
- [x] `cmd/agent` 支持 `--validate-config`、`--backup`、`--restore`、`--restore-confirm`。
- [x] SQLite 备份使用一致性快照；缺失源 DB 时失败且不创建空库。
- [x] restore 默认拒绝执行；无确认或目标 DB 已存在时拒绝覆盖；恢复使用不可预测临时文件。

当前范围：只做本地运行硬化、配置校验、备份恢复和 CLI smoke；不做云同步、多用户权限或复杂安装器。

#### P25–P27：真实公开数据源接入状态

P25 已完成真实公开源验证，P26 已完成首批公告与监管证据 collector 实现并归档；P27 已完成首批基金净值与 ETF/指数市场数据 collector 实现并归档。

已完成状态：

- P25：验证巨潮资讯、上交所、深交所、证监会、基金业协会、东方财富基金、新浪财经、中证指数等公开源的访问方式、字段、频率和限制。
- P25：明确 P19/P20 只是基础 HTTP bridge 和 payload parser，不等同于已接通所有真实外部源。
- P26：实现巨潮资讯、深交所、证监会首批只读公告/监管 collector，包含默认 90 天/分页边界、RAG 入库、source verification、审计、幂等去重和 `cmd/agent --task public-evidence-refresh` 手动触发入口。
- P27：实现东方财富基金首批只读市场数据 collector 与中证指数基础信息 collector，包含基金净值/历史净值/资产配置/基金档案基础 metadata、中证指数当前公开 `index-basic-info` shape、market-refresh 写入、配置、幂等和测试；中证指数样本/权重/估值文件扩展接口仍按候选 metadata 低频读取。

仍作为后续独立候选的方向：

- 更完整的真实环境 E2E / Playwright 覆盖。
- 成分股财务结构化入库和更细数据源扩展。
- 情绪日志、基准对比和更深复盘归因。
- 季度再平衡任务和更完整 SOP A–F 场景图。
- 更精细流动性、复合高危区和大师智慧映射。
- 若未来需要外部通知渠道，必须单独立项并重新确认安全边界。

## 5. 依赖关系

```text
P0–P18 工程骨架、治理与本地运维基础
  ↓
P19–P24 可用 MVP 路径
  ├─ 公开 HTTP 数据桥接与 A 股 ETF/基金证据源
  ├─ 应用内通知中心
  ├─ 规则提案增强
  ├─ 复盘深化
  └─ 本地配置校验、备份与恢复
  ↓
P25 真实公开源调研验证（已完成）
  ├─ 确认公开权威源访问方式、字段、频率和合规边界
  ├─ 修正 P19/P20 基础能力与真实 collector 的边界
  └─ 输出 P26/P27 可实现范围
  ↓
P26 公告与证据源 collector（已完成并归档）
  ├─ 巨潮资讯、深交所、证监会首批 collector
  ├─ 证据入库、RAG chunks、source verification、审计与幂等
  └─ cmd/agent 手动触发入口
  ↓
P27 基金净值与 ETF 市场数据 collector（已完成并归档）
  ├─ 东方财富基金净值、历史净值、资产配置、基金档案基础 metadata
  ├─ 中证指数基础信息（当前公开 `index-basic-info` shape 已校准），样本、权重、估值文件扩展接口按候选 metadata 低频读取
  └─ 默认关闭配置、market-refresh 写入、幂等与测试
```

规划建议：

- P19–P24 已作为当前可用 MVP 基线，不再作为待开发阶段重复拆分。
- P25 已完成真实公开源调研验证；P26 已完成首批公告/监管证据 collector 实现并归档；P27 已完成东方财富基金首批市场数据 collector 与中证指数候选 collector 实现并归档。
- 下一阶段应从剩余候选增强中选择独立主题，先创建新的 OpenSpec change，再进入实现。
- 外部通知、券商接口、付费/登录数据源、自动交易等仍不在当前边界内；若要引入，必须单独立项并重新审查安全边界。
- 长文档拆分或目录迁移应作为独立文档架构 change，不要与功能开发混在一起。

## 6. P31–P40 后续开发完成路线图

以下路线图用于把 P30 后剩余工作落到可规划的阶段任务。每一项都应先创建独立 OpenSpec change，再进入实现；不得在没有 change 的情况下直接修改 L1 契约或扩大产品边界。

### P31：每日自动运行闭环（已实现，待归档）

目标：让系统在用户显式启用后，能够按本地配置自动执行每日信息刷新与纪律检查，并生成可审计的每日结果。

已完成：

- 增加默认关闭的本地 scheduler 配置，支持运行时间、时区、启停开关、失败重试、最大运行时长和最大标的数；启用后按 `run_time` / `timezone` 计算下一次触发，不在 server 启动时立即运行。
- 串联市场刷新、公开证据刷新、每日纪律工作流、应用内通知和审计记录；运行入口随本地 server 生命周期启动，但仅在 `daily_auto_run.enabled: true` 后工作。
- 每日运行 scope 以本地账户/组合当前持仓为准；缺持仓时写入 `missing_prerequisites` 失败状态、应用内通知和审计诊断，不生成正式交易建议。
- 以本地日期、scope、持仓 symbol set 和任务版本生成幂等 key；执行副作用前先写入 `running` 状态，重复运行复用既有结果并写 `daily_auto_run_reuse` 审计，避免重复通知刷屏和重复生成每日结果。
- 自动运行步骤支持有限重试和超时保护，失败时写入结构化错误码、可读原因和 `status=...;step=...;safety=no_auto_trading` 诊断摘要。
- 前端新增“每日自动运行”页面，展示 enabled、计划/运行/成功/降级/失败状态、last/next run、缺失项说明、通知/审计/决策跳转和非自动交易安全边界。

验收方向：

- 显式启用 scheduler 后，本地服务可自动触发每日刷新与纪律评估，并写入状态、审计和应用内通知。
- 关闭 scheduler 后不会自动运行，也不会因 server 启动写入每日自动运行结果。
- 缺数据、源失败或超时时不伪造建议，前端展示明确降级或失败原因。
- 不新增自动交易、外部推送、Webhook、券商接口或收益承诺。

### P32：每日纪律报告产品化（已完成）

目标：把 Daily workflow 从可触发任务提升为用户每天打开即可使用的核心页面和历史报告体系。

已立项目标：

- 新增每日纪律报告索引模型、持久化仓储、迁移和 wiring，用于统一保存自动运行、手动运行、成功、降级和缺前提状态。
- 新增今日报告、历史列表和详情 API/DTO，聚合关联 decision/evidence/audit/notification 链接，避免前端直接读取本地存储。
- 产品化今日纪律页、历史报告列表页和详情页，展示报告状态、摘要、缺前提、关联材料和“不会自动执行交易”的安全边界。
- 将 P32 smoke seed 写入报告索引，并扩展 Playwright smoke 覆盖今日报告、历史列表和详情回看。

已实现：

- 完善今日纪律页，集中展示组合状态、持仓状态、估值、风险红线、触发规则、证据摘要和最终裁决。
- 保存每日纪律报告历史，并支持查看最近 N 次运行结果。
- 增加“今日是否已完成纪律检查”的状态展示。
- 增加红线、估值、证据覆盖和降级原因的趋势展示。
- 将每日自动运行结果、手动运行结果和失败状态统一到同一报告模型。
- 复审 findings 已修复并重新验证；P32 已归档到 `openspec/changes/archive/2026-06-11-p32-daily-discipline-report-productization/`。

验收方向：

- 空库展示初始化引导，不生成正式建议。
- 数据完整时生成每日纪律报告并写入 `decision_records`、`evidence_refs` 和 `audit_events`。
- 同一交易日重复运行保持幂等或有明确版本关系。
- 前端能查看今日报告和历史报告。

### P33–P40：当前剩余计划内功能队列

P33–P40 是当前路线图中剩余的计划内功能阶段。除 P19–P24 历史 archive 追溯、P40 后新路线图这类治理事项外，后续功能开发应从 P33–P40 中选择，并为每个阶段单独创建 OpenSpec change。

推荐执行顺序：

1. **第一组：P33 → P39**。先完成账户与持仓初始化，再补齐浏览器级完整用户旅程，确保空库用户能走到第一份每日纪律报告。
2. **第二组：P34 → P35 → P38**。先扩展真实公开数据覆盖，再建立风险预警与 SOP，随后加固 RAG / VecLite 检索质量。
3. **第三组：P36 → P37 → P40**。完善规则进化效果验证、真实 LLM 质量评估和本地交付运维演练。

边界：P33–P40 仍不得接券商交易 API、自动交易、外部推送、登录源、付费源、授权源、Level2 或高频源；不得承诺收益或预测确定涨跌；LLM 不得覆盖最终规则裁决。P19–P24 如需 archive 级追溯，应单独发起治理 change，不作为 P33–P40 功能任务。P40 后若继续扩展产品愿景，应新建下一轮路线图 change。

### P33：账户与持仓录入/校准体验

目标：让真实用户能从空库完成账户初始化、持仓录入、线下交易记录和数据校准。

需要完成：

- 新增账户初始化向导，录入现金、总资产、持仓、成本、买入原因、资产标签和风险偏好基础信息。
- 增加持仓新增、编辑、校验和删除前置确认；历史事实仍按追加式审计记录处理。
- 支持线下交易流水录入后的持仓、现金和快照一致性校验。
- 支持批量导入或表格化录入持仓与历史交易。
- 增加录入错误修正流程，避免静默覆盖历史数据。
- 前端展示首次使用到可运行每日纪律的完整引导状态。
- 当前 change id：`p33-account-position-onboarding`，已归档到 `openspec/changes/archive/2026-06-12-p33-account-position-onboarding/`，完成后端、HTTP API、Portfolio 前端入口、Dashboard/每日纪律缺前提引导、文档和验收。
- 依赖：P32 每日纪律报告索引已完成；P33 输出供 P39 全路径 E2E 使用。

验收方向：

- 空库用户能通过 UI 完成账户初始化。
- 录入持仓后 Dashboard 进入可分析状态或明确提示仍缺哪些数据。
- `executed_manually` 只记录用户线下动作，并保持账户、持仓、交易流水和审计一致。
- 错误输入被校验拦截，不能写入不一致快照。
- P33 验收命令：`go test ./...`、`npm --prefix web test -- --run`、`npm --prefix web run build`、P33 定向 smoke、`openspec validate p33-account-position-onboarding --strict`、`openspec validate --all --strict`、`git status --short`。

### P34：真实数据覆盖扩展

目标：扩展日常纪律和风险判断所需的真实公开数据覆盖面，同时保留只读、低频、可降级边界。

需要完成：

- 继续验证并接入中证指数样本、权重、估值文件等公开数据。
- 增加成分股财务、资金流向、融资融券或可替代情绪指标的数据模型和 collector。
- 为每类数据定义新鲜度、缺失、过期和失败分类。
- 将新增数据接入每日纪律、风险预警和 expected return 输入上下文。
- 增加数据源健康状态和最近成功/失败记录。
- 当前 change id：`p34-real-data-coverage-expansion`，已归档到 `openspec/changes/archive/2026-06-15-p34-real-data-coverage-expansion/`，完成真实公开数据覆盖扩展、结构化 source health、失败分类、工作流输入和前端状态展示。
- 依赖：P26/P27/P29 公开源 collector 基线；输出供 P35 风险预警和 P38 检索质量使用。

验收方向：

- 指定 ETF/基金/指数能刷新完整日频指标或明确说明缺失原因。
- 估值、财务、资金或情绪数据缺失时不伪造结果。
- 数据源失败能分类为 `no_data`、`source_unavailable`、`parse_error` 或等价状态。
- 新数据能被工作流读取并进入审计上下文。

- P34 验收命令：`go test ./...`、`npm --prefix web test -- --run`、`npm --prefix web run build`、P34 fixture/真实公开源 smoke、`openspec validate p34-real-data-coverage-expansion --strict`、`openspec validate --all --strict`、`git status --short`。

### P35：风险预警与 SOP 编排

目标：把规则裁决、通知和前端展示组织成可持续追踪的风险预警体系。

需要完成：

- 建立风险预警中心，覆盖估值高位、买入逻辑破坏、流动性不足、情绪极端、仓位超限和证据不足。
- 将冷静机制、冻结观察、只卖不买、分批止盈提醒、重新评估提示等 SOP 明确为状态流转。
- 将风险预警写入通知、审计和每日纪律报告。
- 支持风险状态解除、持续观察和升级。
- 前端展示每个风险的触发证据、当前状态、建议人工动作和禁止动作。
- 建议 change id：`p35-risk-alert-sop-orchestration`。
- 依赖：P34 的数据新鲜度与失败分类；输出供 P39 浏览器级风险路径验收使用。

验收方向：

- 构造不同风险输入时，系统进入对应 SOP 状态。
- 风险预警可在前端、通知中心和审计页追踪。
- 风险解除后状态可恢复或归档。
- 不产生自动交易动作或确定性涨跌预测。

### P36：规则进化效果验证

目标：让规则提案不只可生成和审批，还能证明来源、样本、风险和应用后效果。

需要完成：

- 增强错误案例聚类和规则提案来源解释。
- 增加规则修改前后对比、影响范围、风险说明和样本代表性评估。
- 增加过拟合检查、最小样本门槛和历史回放验证。
- 规则应用后追踪命中率、误判率、缺证据率和降级情况。
- 前端展示规则有效性趋势和提案应用后的追踪结果。
- 建议 change id：`p36-rule-evolution-effect-validation`。
- 依赖：P22/P23/P28/P35 的提案、复盘、动态卖出和风险状态基础；仍需守门人审计与用户最终确认。

验收方向：

- 样本不足或过拟合风险高的提案被守门人拒绝。
- 通过审计的提案仍需用户最终确认才会应用。
- 已应用规则能在后续复盘中看到效果追踪。
- 规则提案不会自动修改根本规则。

### P37：真实 LLM 使用与质量评估

目标：验证 DeepSeek/LLM 在真实配置下能稳定生成分析材料，并确保不会影响最终裁决边界。

需要完成：

- 增加真实 DeepSeek 配置 smoke：`cmd/agent --task llm-smoke --symbol <symbol>` 显式读取本地私有配置；常规工作流默认不需要真实 key 即可降级通过。
- 管理 prompt 版本、输入摘要、输出摘要和解析状态。
- 增加 LLM 超时、空响应、格式错误、不可用和质量不足的错误分类。
- 增加输出质量评估 fixture，检查是否越权预测、承诺收益或覆盖最终裁决。
- 审计 LLM 调用输入，避免发送不必要的敏感信息。
- 建议 change id：`p37-real-llm-quality-evaluation`。
- 依赖：P7/P28 的 analyst reports 与 expected return 基线；LLM 仍只生成分析材料。

验收方向：

- 有真实 key 时可通过 `llm-smoke` 生成 analyst report，并记录 prompt/model/parse/quality 脱敏审计摘要。
- 无 key 或调用失败时工作流降级，规则裁决仍可运行。
- LLM 输出不能写最终 verdict，不能触发交易动作。
- prompt 版本和调用结果可审计。

### P38：RAG / VecLite 检索质量加固

目标：从“能检索、能降级”提升到“召回质量可评估、证据引用可验证”。

需要完成：

- 建立检索质量测试集，覆盖公告、监管文件、基金信息和背景材料。
- 增加混合检索、重排序或等价质量增强策略。
- 校验证据引用与 `source_verifications`、信源等级、时效权重一致。
- 增加索引新鲜度、损坏、重建和版本兼容检查。
- 前端或审计中展示检索降级原因和召回摘要。
- 建议 change id：`p38-rag-veclite-retrieval-quality`。
- 依赖：P13/P15/P26/P27/P34 的索引、证据质量和公开数据基线。

验收方向：

- 给定查询能召回预期证据。
- VecLite 不可用时降级到 SQLite 摘要或信息不足。
- C 级或过期证据不会成为 formal 裁决依据。
- 索引重建后结果稳定且可审计。

### P39：前端完整用户旅程与全路径 E2E

目标：补齐真实用户从初始化到每日使用、咨询、确认、复盘和规则治理的浏览器级验收。

需要完成：

- 增加新手引导、数据源配置向导和账户初始化流程。
- 完善主动咨询、用户确认、错误标注、规则提案、复盘历史和审计追踪的跨页体验。
- 扩展 Playwright E2E，覆盖每日纪律、用户确认、规则治理、降级路径、市场刷新、expected return 和禁止自动交易。
- 增加窄屏、基础可访问性和 console error 检查。
- 保持 Vitest 单测与 Playwright E2E 分层，避免互相收集。
- 已归档 change id：`p39-frontend-full-user-journey-e2e`。
- 依赖：P33 账户初始化；可吸收 P34/P35/P38 已完成路径作为更完整的浏览器验收。

验收方向：

- 空库用户能从初始化走到第一份每日纪律报告。
- 浏览器级 E2E 覆盖 A01–A17 的关键路径。
- 所有关键页面无 console error 和未处理异常。
- 页面不出现自动下单、一键交易或收益承诺文案。

### P40：本地部署、运维与恢复演练

目标：把当前本地运行能力整理成稳定、可重复、可诊断的交付体验。

已完成：

- 增加 `cmd/agent --preflight --diagnostics`，验证 Go、Node、npm、Playwright browser、SQLite path、VecLite path、配置文件、数据源和 LLM 配置。
- 增加本地启动前自检和可修复提示，诊断不输出密钥原文。
- 增加 `scripts/recovery-smoke.sh`，用临时 SQLite / VecLite / 配置路径验证备份恢复后本地事实可读。
- 增加数据源健康和本地运行就绪面板，展示最近成功时间、失败分类、数据新鲜度、影响范围和安全文案。
- 明确本地日志、临时文件、诊断文件和 gitignore 治理策略。
- change id：`p40-local-deploy-ops-recovery-drill`，已归档到 `openspec/changes/archive/2026-06-16-p40-local-deploy-ops-recovery-drill/`。
- 依赖：P33–P39 主要用户路径和运行能力已稳定；P40 后新增产品愿景需另建路线图 change。

验收方向：

- 新环境按文档可完成初始化、启动、测试和 smoke。
- 备份恢复不会无确认覆盖现有 DB。
- 诊断文件不会污染 git 工作树。
- 配置错误能给出明确修复建议。

### P41-roadmap：P40 后下一轮路线图治理

目标：在 P33-P40 当前计划内功能队列完成后，重新固化下一轮候选方向、依赖、验收边界和禁止事项，避免无活跃 change 时直接实现新功能。

需要完成：

- 建立 P40 后候选队列，并把候选方向拆分为产品能力增强、数据质量增强、运维体验增强和历史审计追溯。
- 为每类候选方向说明依赖关系、验收思路、是否适合下一阶段优先处理，以及需要单独创建的 OpenSpec change。
- 明确 P19-P24 历史 archive 追溯属于独立治理候选，不伪造历史归档，也不混入新产品功能阶段。
- 明确保留既有安全边界：不接券商 API、不自动交易、不外部推送、不自动应用规则、不承诺收益、不预测确定涨跌；登录源、付费源、授权源、Level2 和高频源默认不纳入；LLM 不写最终裁决。
- 建议 change id：`p41-post-p40-roadmap-governance`。
- 依赖：P40 已归档；当前无其他活跃 change。

历史候选方向（P45 归档时固化，P49 已在当前阶段进入实现）：

| 候选方向 | 建议后续 change | 依赖 | 验收思路 | 下一阶段适合度 |
| --- | --- | --- | --- | --- |
| 产品能力增强 | `p42-user-decision-workbench` | P33 账户/持仓、P32 每日纪律报告、P35 风险预警、P36 规则治理、P39 全路径 E2E | 明确每日使用工作台、主动咨询、组合复盘、规则治理入口的用户旅程；每个入口只展示本地事实、分析材料、审计链路和人工确认状态；浏览器 E2E 覆盖空库、数据完整、降级和窄屏路径 | **高**：最贴近日常使用价值，适合作为 P41 后首个产品功能阶段 |
| 数据质量增强 | `p43-data-quality-observability` | P34 真实数据覆盖、P38 RAG / VecLite 检索质量、P37 真实 LLM 质量评估、P40 预检/诊断 | 明确 source health、证据新鲜度、RAG 命中率、LLM smoke/质量门禁的统一质量视图；验收覆盖 no_data、source_unavailable、parse_error、stale、missing、unknown 和脱敏审计 | **中高**：适合在产品工作台后推进，也可在真实数据风险较高时提前 |
| 运维体验增强 | `p44-local-install-diagnostics-packaging` | P40 本地预检、恢复 smoke、P31 本地调度、P39/P40 smoke 脚本 | 明确本地安装检查、配置向导、诊断导出、备份/恢复演练和 smoke 汇总的操作体验；验收覆盖临时目录、日志脱敏、失败保留排查材料和不污染真实私有数据库 | **中**：适合交付给更多本地用户前推进，不应替代核心投资工作流 |
| 历史审计追溯 | `p19-p24-historical-archive-traceability` | P19-P24 已交付但无完整 archive 包，P41 已明确追溯边界 | 明确只整理历史证据、验收记录、文档一致性和缺口说明；不得伪造历史 change，不得回写虚假任务完成时间；验收以文档一致性、OpenSpec 校验和审计说明为主 | **条件触发**：如需要审计级追溯则优先，否则不阻塞 P42/P43/P44 |

默认推荐顺序：先做 P42 产品能力增强，再按数据质量风险选择 P43，随后推进 P44；P19-P24 历史审计追溯按审计需求独立排期。任何候选一旦进入实现，都必须先创建独立 OpenSpec change，并在 proposal 中重新确认 out of scope 和安全边界。

验收方向：

- `openspec/PROGRESS.md`、`docs/GOVERNANCE.md`、`AGENTS.md` 和 `openspec/project.md` 均指向 P41 活跃 change。
- 路线图明确后续候选方向、依赖、出入范围和验收方式。
- 每个后续功能都要求独立 OpenSpec change 后再实现。
- 本 change 不修改运行时代码。

### P42：用户决策工作台

目标：把今日纪律、组合风险、规则治理、复盘摘要和主动咨询入口聚合为单一日常工作台，让用户每天打开系统后能先看到“今天要看什么、缺什么、能做什么、不能做什么”。

需要完成：

- 新增 `/workbench` 页面或等价首屏入口，展示“今日先看”“组合与风险”“规则与复盘”“主动咨询入口”四类信息。
- 工作台只复用现有 API/service DTO，不直接读取 SQLite、VecLite、本地日志或配置文件。
- 工作台提供到每日纪律报告、持仓、风险预警、规则提案、复盘摘要、审计和决策咨询的导航入口。
- 工作台在空库、数据缺失、source health 降级、LLM/RAG/VecLite 不可用时展示安全状态和可检查下一步。
- 工作台不得展示自动交易、一键交易、代下单、自动外推、自动确认或自动应用规则入口。
- 建议 change id：`p42-user-decision-workbench`。
- 依赖：P32、P33、P35、P36、P39、P41 已完成。

验收方向：

- `npm --prefix web test -- --run`、`npm --prefix web run build`、`bash scripts/e2e-smoke.sh` 通过。
- `/workbench` 在本地 fixture 下可达，核心区域可见，窄屏可用。
- 工作台只导航、过滤或展示本地事实，不触发交易、外推、自动确认或自动规则应用。
- OpenSpec 校验通过，archive 前只读复审无 Critical / Important。

### P43：数据质量可观测

目标：把 source health、证据新鲜度、RAG/VecLite 检索质量、LLM 质量门禁和本地诊断状态聚合为统一只读页面，让用户能快速判断“哪里降级、影响什么、下一步检查哪里”。

需要完成：

- 新增 `/data-quality` 页面或等价入口，展示“数据源健康”“证据与检索”“LLM 质量”“影响范围与下一步”四类信息。
- 页面只复用现有 API/service DTO；如必须扩展后端，只允许新增只读聚合 DTO/API，不新增数据库 migration。
- 页面提供到设置、证据、复盘、审计、风险预警、决策详情和工作台的导航入口。
- 页面在 source_unavailable、parse_error、stale、missing、unknown、LLM/RAG/VecLite 不可用时展示安全状态和可检查下一步。
- 页面不得展示完整 key、完整 prompt、私有本地路径、SQL、供应商原始错误或账户敏感明细。
- 页面不得展示自动刷新修复、外部推送、自动确认、自动应用规则、自动交易、一键交易、代下单或收益承诺入口。
- 建议 change id：`p43-data-quality-observability`。
- 依赖：P34、P37、P38、P40、P42 已完成。

验收方向：

- `npm --prefix web test -- --run`、`npm --prefix web run build`、`bash scripts/e2e-smoke.sh` 通过。
- `/data-quality` 在本地 fixture 下可达，核心区域可见，窄屏可用。
- 页面只导航、过滤或展示本地质量事实，不触发刷新、修复、外推、自动确认、自动规则应用或交易。
- OpenSpec 校验通过，执行前与 archive 前只读复审无 Critical / Important。

### P44：本地安装诊断与打包

目标：把本地预检、配置向导、诊断导出、备份/恢复演练和 smoke 汇总整理成统一、可重复、可审计的本地交付体验。

已完成：

- 新增 `scripts/local-install-diagnostics.sh`，集中运行预检、恢复 smoke 和可选 e2e smoke，并输出 `install-summary.json`。
- 新增 `/local-install` 页面，展示配置草稿、关键命令、诊断摘要导入和安全边界。
- 更新前端路由、导航和 smoke 覆盖，确保页面可达且不出现高风险入口。
- 更新本地配置与运维文档，明确日志、临时目录、摘要路径和脱敏边界。
- change id：`p44-local-install-diagnostics-packaging`，已归档到 `openspec/changes/archive/2026-06-16-p44-local-install-diagnostics-packaging/`。

验收方向：

- `go test ./...`、`npm --prefix web test -- --run`、`npm --prefix web run build` 和 `openspec validate --all --strict` 通过。
- 页面和脚本只展示本地诊断事实，不接券商、不交易、不外部推送、不自动确认、不自动应用规则、不承诺收益。

### P45-roadmap：P44 后路线图治理

目标：在 P42、P43、P44 与 P19-P24 历史追溯均已完成归档后，重新固化下一轮候选方向、依赖、验收边界和禁止事项，避免 `next_change_id` 为空时直接实现新功能。

需要完成：

- 建立 P44 后候选队列，并拆分为本地知识治理、决策闭环解释、数据质量回归和运维发布体验。
- 为每类候选方向说明建议阶段、change id、依赖、验收方式和适合度。
- 推荐 P46 `p46-local-knowledge-import-governance` 作为下一功能候选，但 P45 不实现 P46。
- 明确保留既有安全边界：不接券商 API、不自动交易、不一键交易、不代下单、不外部推送、不自动确认、不自动应用规则、不自动修复承诺、不承诺收益、不预测确定涨跌；登录源、付费源、授权源、Level2 和高频源默认不纳入；LLM 不写最终裁决。
- 建议 change id：`p45-post-p44-roadmap-governance`。
- 依赖：P44 已归档；当前无其他活跃 change。

候选方向：

| 候选方向 | 建议阶段 | 建议后续 change | 依赖 | 验收思路 | 下一阶段适合度 |
| --- | --- | --- | --- | --- | --- |
| 本地知识治理 | P46 | `p46-local-knowledge-import-governance` | P33 账户/持仓、P38 RAG/VecLite、P43 数据质量、P44 本地诊断 | 本地文件/笔记/CSV 导入前检查、脱敏预览、索引重建计划和失败可回滚；不得导入私钥、完整 key 或原始 SQL | **高**：承接检索质量与本地诊断，能安全扩大用户自有研究材料上下文 |
| 决策闭环解释 | P47 | `p47-decision-loop-explainability` | P32 每日纪律、P35 风险预警、P36 规则进化、P42 工作台 | 将建议、人工确认、线下执行记录和结果复盘串成可解释闭环；继续保持人工记录，不触发交易 | **中高**：适合在知识导入后增强日常复盘价值 |
| 数据质量回归 | P48 | `p48-data-source-quality-regression-pack` | P34 数据覆盖、P37 LLM 质量、P38 检索质量、P43 质量面板 | 形成可重复运行的真实/fixture 数据质量回归包，覆盖 source health、freshness、parse_error、no_data 和脱敏摘要 | **中高**：适合在数据源变化或真实源风险升高时优先 |
| 运维发布体验 | P49 | `p49-local-release-upgrade-experience` | P40 恢复演练、P44 诊断打包 | 本地版本检查、升级前备份提醒、迁移前预检和升级后 smoke 汇总；不自动修复、不覆盖真实库 | **中**：适合更多本地用户开始升级时推进 |

P45 归档时的默认推荐顺序：先做 P46，再按使用反馈选择 P47 或 P48；P49 适合在本地发布/升级需求更明确后推进。P49 当前已创建独立 OpenSpec change，并在 proposal 中重新确认 out of scope 和安全边界。

验收方向：

- `openspec validate p45-post-p44-roadmap-governance --strict` 通过。
- `openspec validate --all --strict` 通过。
- `git diff --check` 通过。
- 执行前与 archive 前只读复审无 Critical / Important。
- 本 change 不修改运行时代码。

### P46：本地知识导入治理

目标：把用户自有研究记录以“校验 -> 脱敏预览 -> 显式确认 -> 本地背景事实”的方式纳入检索上下文，同时保持 C/background、可审计、可回滚和不触发交易的边界。

需要完成：

- 新增 `POST /api/v1/local-knowledge/imports/validate`，返回 import batch、逐行脱敏预览、风险、计数和索引计划。
- 新增 `POST /api/v1/local-knowledge/imports/confirm`，重新校验 rows，重算 `import_batch_id`，批次不匹配或存在 blocking 风险时拒绝写入。
- 确认写入复用 `intelligence_items`、`intelligence_summary`、`rag_chunks`、`source_verifications` 和 `audit_events`，默认 `source_level=C`、`evidence_role=background`、`verification_status=background_only`、`index_status=pending`。
- 新增 `/local-knowledge` 页面，展示导入草稿、脱敏预览、索引计划、确认结果和安全边界。
- 更新 API、数据模型、前端契约、smoke 和 OpenSpec 进度。
- change id：`p46-local-knowledge-import-governance`。

验收方向：

- `go test ./...`、`npm --prefix web test -- --run`、`npm --prefix web run build`、`bash scripts/e2e-smoke.sh` 通过。
- `openspec validate p46-local-knowledge-import-governance --strict`、`openspec validate --all --strict`、`git diff --check` 通过。
- 安全扫描覆盖完整 key、私有路径、原始 SQL、私钥、完整 prompt 和禁止能力文案；命中项必须人工复核。
- 执行前与 archive 前只读复审无 Critical / Important。

### P47：决策闭环解释

目标：新增只读决策闭环解释能力，把建议、用户确认、本地线下记录、风险线索、复盘线索和审计线索串成可阅读的本地解释链。

已实现范围：

- 新增 `GET /api/v1/decision-loops` 与 `GET /api/v1/decision-loops/{decision_id}`。
- 新增 `DecisionLoopService`，从 `decision_records`、`operation_confirmations`、`position_transactions`、`error_cases`、`risk_alerts` 和 `audit_events` 聚合 `DecisionLoopItem`。
- 新增只读 repository 查询 `ListOperationConfirmations` 与 `ListPositionTransactionsByConfirmation`；未新增 migration。
- 新增 `/decision-loop` 页面、主导航入口、工作台入口、复盘摘要入口和 Playwright smoke 覆盖。
- 响应不返回 raw payload、私有路径、完整 key、原始 SQL、完整 prompt 或供应商原始响应；备注仅以脱敏 `note_preview` 展示。

安全边界：

- 不新增券商接口、交易执行、一键式交易、代下单、外部推送、自动确认、自动规则应用、自动修复承诺、收益承诺、登录源、付费源、授权源、Level2 或高频源。
- P47 只读聚合不得写入 `decision_records`、`operation_confirmations`、`position_transactions`、`risk_alerts`、`audit_events`、`notifications` 或规则版本。

验收方向：

- `go test ./...`、`npm --prefix web test -- --run`、`npm --prefix web run build`、`bash scripts/e2e-smoke.sh` 通过。
- `openspec validate p47-decision-loop-explainability --strict`、`openspec validate --all --strict`、`git diff --check` 通过。
- 执行完成后必须进行只读子 agent 复审，无 Critical / Important 后方可归档。

### P48：数据源质量回归包

目标：提供本地可重复运行的数据源质量回归包，验证 source health/freshness、`no_data`、`source_unavailable`、`parse_error`、`stale` 分类和脱敏摘要，帮助区分真实源波动与本地解析/展示逻辑退化。

已实现范围：

- 新增 `GET /api/v1/data-source-quality/regression`，支持 `fixture` 和 `current` 两种模式。
- 新增 `DataSourceQualityService` 与 DTO，默认 fixture 离线回归覆盖 fresh、no_data、source_unavailable、parse_error、stale 和敏感诊断脱敏。
- 将 P34 source health 从 market snapshot 提取逻辑移到 service 层，供 `/market/source-health` 与 P48 回归复用。
- 新增 `cmd/agent --task data-source-quality-regression --source fixture|current --symbol 000300`，写入一条脱敏本地任务审计摘要。
- 不新增 migration、collector、调度器、外部推送、自动修复或交易能力。

安全边界：

- `fixture` 模式不得访问公网；`current` 模式只读现有 `market_snapshots.market_metrics_json.metadata.p34_source_health`。
- 响应和审计摘要不得包含完整 key、私有路径、原始 SQL、完整 prompt、raw HTTP、private key 或供应商原始响应。
- P48 不修改账户、持仓、确认记录、风险 SOP、规则版本或通知。

验收方向：

- `go test ./...`、`npm --prefix web test -- --run`、`npm --prefix web run build`、`bash scripts/e2e-smoke.sh` 通过。
- `openspec validate p48-data-source-quality-regression-pack --strict`、`openspec validate --all --strict`、`git diff --check` 通过。
- 安全扫描覆盖 P48 code/docs，命中项必须人工复核；执行完成后必须进行只读子 agent 复审，无 Critical / Important 后方可归档。

### P49：本地发布与升级体验

已完成：

- 新增本地发布/升级检查，覆盖版本检查、升级前备份提醒、迁移前预检和升级后 smoke 汇总。
- 诊断 JSON 与脚本摘要保持脱敏，不输出完整 key、私有路径、原始 SQL 或敏感配置。
- 继续禁止自动升级、自动迁移、自动修复、覆盖真实库、交易、外推和收益承诺。
- change id：`p49-local-release-upgrade-experience`，已归档到 `openspec/changes/archive/2026-06-17-p49-local-release-upgrade-experience/`。

验收方向：

- `go test ./...`、`npm --prefix web test -- --run`、`npm --prefix web run build`、`bash scripts/local-release-upgrade-check.sh --target-version vNEXT --output-dir <tmp>/release-upgrade` 通过。
- `openspec validate p49-local-release-upgrade-experience --strict`、`openspec validate --all --strict`、`git diff --check` 通过。
- 执行完成后只读子 agent 复审无 Critical / Important。

### P50-roadmap：P49 后治理与验收路线图

目标：在 P49 完成后，先收敛 P19-P24 历史审计与全项目验收门禁，再进入发布候选材料。

已完成：

- 明确 P19-P24 已有 `p19-p24-historical-archive-traceability` 审计说明，但仍缺逐阶段完整 archive 包；后续只能补当前事实证据包，不伪造历史 archive。
- 固化 P51 `p51-p19-p24-audit-evidence-pack` 为下一阶段，输出交付边界、事实源、代码/测试/文档证据、可重跑命令、缺口和不可声明事项。
- 固化 P52 `p52-project-acceptance-gate-matrix` 为 P51 后续阶段，覆盖单元测试、集成测试、E2E 测试、真实源测试、真实 LLM 测试、冒烟测试、安装诊断、发布升级检查和安全边界门禁。
- 固化 P53 `p53-acceptance-execution-and-release-candidate-materials` 必须依赖 P51/P52 完成后再进入。
- change id：`p50-post-p49-governance-validation-roadmap`，已归档到 `openspec/changes/archive/2026-06-17-p50-post-p49-governance-validation-roadmap/`。

验收方向：

- `openspec validate p50-post-p49-governance-validation-roadmap --strict` 通过。
- `openspec validate --all --strict` 通过。
- `git diff --check` 通过。
- 执行前与 archive 前只读复审无 Critical / Important。
- 本 change 不修改运行时代码。

### P51：P19-P24 审计证据包

目标：为 P19-P24 可用 MVP 阶段建立当前事实审计证据包，作为 P52 验收门禁矩阵和 P53 发布候选材料的前置引用基础。

已完成：

- 新增 `docs/p19-p24-audit-evidence-pack.md`。
- 明确 P14-P18 已有标准 archive，P19-P24 缺逐阶段完整 archive 包。
- 按 P19-P24 分阶段列出交付边界、archive 状态、文档证据、代码证据、测试证据、可重跑命令、不可声明事项和残余缺口。
- 说明 P25-P29、P30/P39、P40/P44/P49 对 P19-P24 的后续补强关系。
- 保留安全边界：不接券商、不自动交易、不外部推送、不自动确认、不自动应用规则、不自动修复、不覆盖真实库、不承诺收益，不新增登录源、付费源、授权源、Level2 或高频源。
- change id：`p51-p19-p24-audit-evidence-pack`，已归档到 `openspec/changes/archive/2026-06-17-p51-p19-p24-audit-evidence-pack/`。

验收方向：

- `openspec validate p51-p19-p24-audit-evidence-pack --strict` 通过。
- `openspec validate --all --strict` 通过。
- `git diff --check` 通过。
- 执行前与 archive 前只读复审无 Critical / Important。
- 本 change 不修改运行时代码。

### P52：项目验收门禁矩阵

目标：定义发布候选材料前必须使用的项目级验收门禁矩阵，明确哪些测试和 smoke 阻断发布、哪些真实测试可降级、如何记录验收结果。

已完成：

- 新增 `docs/project-acceptance-gate-matrix.md`。
- 定义 G0-G9 门禁，覆盖治理、Go、集成、前端、E2E、fixture/current smoke、真实公开源 opt-in、真实 LLM opt-in、本地安装/升级和安全脱敏检查。
- 每个门禁写明命令、前置条件、通过标准、允许降级、产物位置和是否阻断发布。
- 定义真实源/真实 LLM 失败分类：网络、限流、凭证、source schema 变化、no_data、parse failure、模型不可用、质量失败、脱敏失败等。
- 定义 P53 发布候选材料必须引用实际验收结果或 waiver，不得把 P52 文档本身当作验收通过。
- change id：`p52-project-acceptance-gate-matrix`，已归档到 `openspec/changes/archive/2026-06-17-p52-project-acceptance-gate-matrix/`。

验收方向：

- `openspec validate p52-project-acceptance-gate-matrix --strict` 通过。
- `openspec validate --all --strict` 通过。
- `git diff --check` 通过。
- 执行前与 archive 前只读复审无 Critical / Important。
- 本 change 不修改运行时代码。

### P53：验收执行与发布候选材料

目标：按 P52 G0-G9 执行真实验收，并基于实际结果生成验收记录和发布候选材料。

已完成：

- 新增 `docs/release/acceptance/2026-06-17-p53-acceptance-run.md`。
- 新增 `docs/release/release-candidate-2026-06-17.md`。
- 执行 G0-G9，覆盖 OpenSpec、Go、前端、E2E、fixture/current smoke、真实公开源、真实 LLM、本地安装/升级和安全脱敏。
- 当前 release status 为 `release_ready`。
- 已记录 G5 current data-source quality degraded：`cases=1:degraded=1:failed=0`，不阻断发布，但限制当前本地数据快照质量声明。
- 已记录 G3/G4/G8 初次本地进程被 kill 后原命令重试通过。
- change id：`p53-acceptance-execution-and-release-candidate-materials`，已归档到 `openspec/changes/archive/2026-06-17-p53-acceptance-execution-and-release-candidate-materials/`。

验收方向：

- `openspec validate p53-acceptance-execution-and-release-candidate-materials --strict` 通过。
- `openspec validate --all --strict` 通过。
- `git diff --check` 通过。
- 执行完成后只读子 agent 复审无 Critical / Important。
- 本 change 不修改运行时代码。

### P54：发布交付与可重复性加固

目标：把 P53 `release_ready` 转化为可交付、可复验、可审计的发布说明，明确重试、降级和安全边界。

已完成：

- 新增 `docs/release/README.md`。
- 新增 `docs/release/release-handoff-2026-06-17.md`。
- 新增 `docs/release/acceptance-repeatability.md`。
- 固化重试规则：仅本地资源 kill、端口启动竞态、浏览器/Vite 启动竞态或短时资源压力允许一次原命令 retry；第二次失败按 blocked 处理，除非显式 waiver。
- 固化 G5 current degraded 规则：fixture pass + current degraded + failed=0 可非阻断，但限制当前本地 DB 健康声明；fixture failed 或 current failed>0 默认阻断。
- 固化 G6 真实源配置前提和 G7 LLM 脱敏规则。
- 明确 P54 不重新执行验收、不改变 P53 `release_ready` 结论、不修改运行时代码。
- change id：`p54-release-handoff-and-repeatability-hardening`，已归档到 `openspec/changes/archive/2026-06-17-p54-release-handoff-and-repeatability-hardening/`。

验收方向：

- `openspec validate p54-release-handoff-and-repeatability-hardening --strict` 通过。
- `openspec validate --all --strict` 通过。
- `git diff --check` 通过。
- 脱敏扫描 release/P54 提交材料通过。
- 执行完成后只读子 agent 复审无 Critical / Important。
- 本 change 不修改运行时代码。

## 7. 开发验收总清单

### 后端

- [x] `go test ./...` 通过。
- [x] migration 可创建完整 SQLite 表。
- [x] Repository 事务测试通过。
- [x] 领域规则单测覆盖主要裁决。
- [x] Eino 工作流测试覆盖正常、降级、失败状态。
- [x] HTTP API 契约测试覆盖主要响应。
- [x] 所有关键动作写入 `audit_events`。
- [x] 真实行情与情报数据源支持配置、stub 和失败降级。
- [x] RAG/VecLite 检索可用、可重建，并能降级到 SQLite 摘要或信息不足。
- [x] DeepSeek 仅生成分析材料，不写最终裁决。
- [x] `cmd/agent` 本地任务可手动触发并记录审计事件。
- [x] 公开 HTTP 行情与情报数据源支持配置、fixture/stub 和失败降级。
- [x] 应用内通知中心已覆盖数据源、索引、规则提案和复盘降级场景。
- [x] `cmd/agent` 支持本地任务、配置校验、SQLite 备份和安全恢复。

### 前端

- [x] `npm run build` 通过.
- [x] 页面字段与 `docs/frontend-contract.md` 一致.
- [x] 驾驶舱三栏结构符合 `docs/ui-design.md`.
- [x] 信息不足、冻结观察、高危状态有明确展示.
- [x] 用户确认区没有自动交易入口.
- [x] 规则提案最终确认流程可见.
- [x] 审计页区分 `action`、`node_name` 与 `node_action`，并按 `status`、`error_code`、输入引用、输出引用展示审计详情.
- [x] 驾驶舱图表使用 API DTO，不直接访问本地存储或索引。
- [x] 前端测试覆盖关键状态、错误响应、确认流程和禁止自动交易入口。
- [x] 通知中心页面展示未读数、通知列表、单条已读和全部已读。
- [x] 复盘页展示 attribution、错误标签、缺证据主题、提案结果、降级 workflow 和追踪入口。
后续独立候选增强（不计入当前 P0–P24 验收；P25 已进入真实公开源调研验证）：

- 更完整真实浏览器 E2E。
- 情绪日志、基准对比、季度再平衡或更深 SOP 场景。

### 产品边界

- [x] DeepSeek 不生成最终裁决。
- [x] 系统不主动推荐具体标的。
- [x] 系统不承诺收益。
- [x] 系统不自动下单。
- [x] C 级信源不作为正式裁决依据。
- [x] 规则提案不自动应用。
- [x] 月度/季度复盘只生成建议和评估，不绕过守门人审计与用户最终确认。

### 端到端验收基线

- [x] 无账户时展示引导，不生成正式建议。
- [x] 每日纪律、主动咨询、证据刷新、市场刷新、规则提案、用户确认和复盘路径都有后端或前端自动化测试覆盖。
- [x] 证据不足、VecLite 不可用、能力圈外、LLM 降级、C 级信源、规则提案不自动应用等边界有回归测试或契约约束。
- [x] 市场刷新覆盖全部成功、部分成功、全部失败和写入失败等响应。
- [x] 页面和 API 不提供自动交易入口。
后续独立候选增强（不计入当前 P0–P24 验收；P25 已进入真实公开源调研验证）：浏览器级 Playwright 全路径仍可作为产品化验收增强单独规划。

## 8. 风险与处理

| 风险 | 表现 | 处理 |
| --- | --- | --- |
| 数据模型字段多 | migration 与 DTO 容易不一致 | 以 `docs/data-model.md` 为准，API 测试校验字段 |
| 工作流复杂 | 节点职责混杂 | 每个节点只读写 WorkflowContext 的指定字段 |
| LLM 输出不可控 | 分析内容偏离规则 | 只把 LLM 输出作为 analyst_reports，最终裁决由规则引擎生成 |
| VecLite 不可用 | 检索失败 | 从 SQLite 的 `rag_chunks` 和 `intelligence_summary` 重建或展示信息不足 |
| 前端误导交易 | 用户以为系统可代为执行 | 文案统一为记录线下动作，不提供自动交易按钮 |
| 规则误应用 | 守门人审计通过后立即生效 | 必须经过 `pending_final_confirm` 与用户最终确认 |

## 9. 建议提交节奏

每个阶段建议独立提交；P19–P110 已作为当前本地源码交付基线完成并归档。P111 当前作为高保真参考图视觉重构活跃 change，范围只覆盖前端视觉层级、reference shell、共享 reference components、全覆盖路由截图对照和响应式验收；不得借 P111 声称新增投资运行时能力、发布包刷新、物理第二机器复验、券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、真实库覆盖或收益承诺。

历史提交主题示例：

1. `chore: initialize project skeleton`
2. `feat: add sqlite schema and repositories`
3. `feat: add domain rule engine`
4. `feat: add eino workflows`
5. `feat: add http api handlers`
6. `feat: add cockpit frontend`
7. `test: add end-to-end acceptance coverage`
8. `feat: integrate real data rag and analyst services`
9. `feat: enhance frontend experience and tests`
10. `feat: add review automation and local delivery`
11. `feat: add public http data bridge`
12. `feat: add in-app notification center`
13. `feat: enrich rule proposals and review traceability`
14. `feat: harden local config backup and restore`

提交前固定执行：

```bash
go test ./...
cd web && npm run build && npm test
```

如果当前阶段不包含前端工程，只执行 Go 测试。
