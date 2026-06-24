# Investment Agent — OpenSpec 项目说明

## 文档真源

本项目的权威文档在 `docs/`，不在 `openspec/specs/` 的全量镜像。L1 契约真源定义系统行为与接口；L2–L3 文档提供架构、计划、体验和图示说明。

| 域 | 级别 | 真源文件 |
| --- | --- | --- |
| 产品需求 | L1 | `docs/requirements.md` |
| 数据模型 | L1 | `docs/data-model.md` |
| HTTP API | L1 | `docs/api.md` |
| Eino 工作流 | L1 | `docs/workflow.md` |
| 前端契约 | L1 | `docs/frontend-contract.md` |
| 功能拆分与验收 | L2 | `docs/functional-spec.md` |
| 架构与技术栈 | L2 | `docs/architecture.md` |
| 开发阶段计划 | L2 | `docs/development-plan.md` |
| UI | L3 | `docs/ui-design.md`、`docs/ui-flow.md` |

治理规则见 **`docs/GOVERNANCE.md`**。

## 变更原则

1. 实现前必须有 `openspec/changes/<id>/` 且包含已审阅的 `tasks.md`。
2. 修改 L1 契约时，change 内必须有 delta（`specs/**/*.md`）。
3. `archive` 时合并 delta 到 `docs/`，勿只更新 `openspec/specs/`。
4. 禁止将 `docs/superpowers/plans/` 当作真源。

## 阶段与 change 映射

| 开发计划 | 建议 change id | 状态 |
| --- | --- | --- |
| P0 工程骨架 | `p0-engineering-skeleton` | done |
| P1 数据底座 | `p1-data-foundation` | done |
| P2 领域规则 | `p2-domain-rules` | done |
| P3 工作流 | `p3-eino-workflows` | done |
| P4 HTTP API | `p4-http-api` | done |
| P5 前端 | `p5-frontend-cockpit` | done |
| P6 验收 | `p6-e2e-hardening` | done |
| P7 真实数据与分析底座 | `p7-real-data-integration` | done |
| P8 前端体验与测试 | `p8-frontend-experience-tests` | done |
| P9 复盘自动化与交付 | `p9-review-automation-delivery` | done |
| P10 产品完整性 | `p10-product-completeness` | done |
| P11 治理与阶段重置 | `p11-governance-and-phase-reset` | done |
| P12 真实数据最小可行源 | `p12-real-data-minimum-viable-sources` | done |
| P13 索引与检索硬化 | `p13-index-and-retrieval-hardening` | done |
| P14 守门人节点图 | `p14-gatekeeper-node-graph` | done |
| P15 证据质量增强 | `p15-evidence-quality-enrichment` | done |
| P16 前端运维复盘表面 | `p16-frontend-ops-review-surface` | done |
| P17 本地调度与运维文档 | `p17-local-scheduler-and-ops-docs` | done |
| P18 规则演化提案硬化 | `p18-evolution-proposal-hardening` | done |
| P19 公开 HTTP 数据桥接 | `p19-real-data-bridge` | delivered; archive not present |
| P20 A 股 ETF/基金证据源 | `p20-public-etf-fund-evidence` | delivered; archive not present |
| P21 应用内通知中心 | `p21-in-app-notification-center` | delivered; archive not present |
| P22 规则提案增强 | `p22-rule-proposal-enrichment` | delivered; archive not present |
| P23 复盘深化与可追溯性 | `p23-review-depth-and-traceability` | delivered; archive not present |
| P24 本地运行硬化 | `p24-local-runtime-hardening` | delivered; archive not present |
| P25 真实公开数据源调研验证 | `p25-real-public-data-sources` | done |
| P26 公告与证据源 collector | `p26-public-evidence-collectors` | done |
| P27 基金净值与 ETF 市场数据 collector | `p27-fund-etf-market-data-collectors` | done |
| P28 预期收益与动态卖出评估 | `p28-expected-return-dynamic-sell` | done |
| P29 公开证据真实采集 smoke | `p29-public-evidence-collector-smoke` | done |
| P30 真实环境 E2E / Playwright smoke 验收 | `p30-real-e2e-smoke` | done |
| P31 每日自动运行闭环 | `p31-daily-auto-run-loop` | done |
| P32 每日纪律报告产品化 | `p32-daily-discipline-report-productization` | done |
| P33–P40 路线图治理 | `p33-p40-roadmap-finalization` | done |
| P33 账户与持仓录入/校准体验 | `p33-account-position-onboarding` | done |
| P34 真实数据覆盖扩展 | `p34-real-data-coverage-expansion` | done |
| P35 风险预警与 SOP 编排 | `p35-risk-alert-sop-orchestration` | done |
| P36 规则进化效果验证 | `p36-rule-evolution-effect-validation` | done |
| P37 真实 LLM 使用与质量评估 | `p37-real-llm-quality-evaluation` | done |
| P38 RAG / VecLite 检索质量加固 | `p38-rag-veclite-retrieval-quality` | done |
| P39 前端完整用户旅程与全路径 E2E | `p39-frontend-full-user-journey-e2e` | done |
| P40 本地部署、运维与恢复演练 | `p40-local-deploy-ops-recovery-drill` | done |
| P41 后路线图治理 | `p41-post-p40-roadmap-governance` | done |
| P42 用户决策工作台 | `p42-user-decision-workbench` | done |
| P43 数据质量可观测 | `p43-data-quality-observability` | done |
| P44 本地安装诊断与打包 | `p44-local-install-diagnostics-packaging` | done |
| P45 P44 后路线图治理 | `p45-post-p44-roadmap-governance` | done |
| P46 本地知识库与数据导入治理 | `p46-local-knowledge-import-governance` | done |
| P47 组合复盘与决策闭环可解释性 | `p47-decision-loop-explainability` | done |
| P48 数据源覆盖与质量回归包 | `p48-data-source-quality-regression-pack` | done |
| P49 运维发布与本地升级体验 | `p49-local-release-upgrade-experience` | done |
| P50 P49 后治理与验收路线图 | `p50-post-p49-governance-validation-roadmap` | done |
| P51 P19-P24 审计证据包 | `p51-p19-p24-audit-evidence-pack` | done |
| P52 项目验收门禁矩阵 | `p52-project-acceptance-gate-matrix` | done |
| P53 验收执行与发布候选材料 | `p53-acceptance-execution-and-release-candidate-materials` | done |
| P54 发布交付与可重复性加固 | `p54-release-handoff-and-repeatability-hardening` | done |
| P55 前端全功能真实验收与设计审查 | `p55-full-ui-acceptance-and-design-audit` | done |
| P56 UI 验收阻断与产品化设计修复 | `p56-ui-acceptance-blocker-fixes` | done |
| P57 产品体验打磨总规划 | `p57-product-experience-polish-roadmap` | done |
| P58 今日工作台重构 | `p58-daily-workbench-redesign` | done |
| P59 决策解释体验重构 | `p59-decision-explainability-experience` | done |
| P60 组合、风险与数据质量体验重构 | `p60-portfolio-risk-data-quality-experience` | done |
| P61 治理和运维页面产品化 | `p61-governance-ops-productization` | done |
| P62 设计系统与可访问性验收 | `p62-design-system-accessibility-hardening` | done |
| P63 全量真实 UI 回归与发布状态刷新 | `p63-full-ui-regression-release-refresh` | done |
| P64 发布打包与版本标记 | `p64-release-packaging-version-tagging` | done |
| P65 跨机器发布包复验 | `p65-cross-machine-release-repeat-acceptance` | done |
| P66 当前数据零退化策略 | `p66-current-data-zero-degradation-policy` | done |
| P67 当前数据门禁处置工作流 | `p67-current-data-gate-resolution-workflow` | done |
| P68 P67 后发布状态治理 | `p68-post-p67-release-readiness-governance` | done |
| P69 Clean tree 最终分发包刷新 | `p69-clean-tree-package-refresh` | done |
| P70 最终发布决策与风险收口 | `p70-final-release-decision-and-risk-closure` | done |
| P71 真实产品验收真通过 | `p71-real-product-acceptance-true-pass` | done |
| P72 真实用户基金场景与数据影响验收 | `p72-real-user-fund-scenario-data-impact-acceptance` | done |
| P73 产品目标效果与 UX 验证 | `p73-product-effectiveness-ux-validation` | done |
| P74 内置知识与数据准备度 | `p74-built-in-knowledge-and-data-readiness` | done |
| P75 原始需求追踪与真实使用闭环 | `p75-requirements-traceability-and-real-use-closure` | done |
| P76 P75 后最终分发包刷新 | `p76-post-p75-final-package-refresh` | done |
| P77 原子需求真 pass 升级门禁 | `p77-requirements-real-pass-upgrade-gate` | done |
| P78 原始需求 real-pass 批次收敛 | `p78-requirements-real-pass-batch-closure` | done |
| P79 真实使用数据影响与预期收益闭环 | `p79-real-use-data-impact-and-expected-return-closure` | done |
| P80 复盘审计与规则治理真实使用闭环 | `p80-review-audit-governance-real-use-closure` | done |
| P81 动态源字段覆盖 | `p81-dynamic-source-field-coverage` | done |
| P82 SOP/action UI-to-SQLite 闭环 | `p82-sop-action-ui-sqlite-closure` | done |
| P83 治理追溯回填 | `p83-governance-traceability-backfill` | done |
| P84 组合与确认数据影响闭环 | `p84-portfolio-confirmation-data-impact-closure` | done |
| P85 预期收益与分析准确性闭环 | `p85-expected-return-analysis-accuracy-closure` | done |
| P87 组合状态、仓位纪律与安全闭环 | `p87-portfolio-state-allocation-safety-closure` | done |
| P86 核心目标、知识/RAG 与安全最终闭环 | `p86-core-goal-knowledge-safety-final-closure` | done |
| P88 剩余 full-release blockers 闭环 | `p88-remaining-full-release-blockers-closure` | done |
| P89 剩余真实 provider 与动态概率闭环 | `p89-real-provider-and-dynamic-probability-closure` | done |
| P90 capital-flow provider 闭环 | `p90-capital-flow-provider-closure` | done |
| P91 GitHub Release 与 Docker 部署 | `p91-github-release-docker-deployment` | done |
| P92 原始需求最终独立复核台账 | `p92-final-original-requirement-audit-ledger` | done |
| P93 最终代码真实性与设计审查 | `p93-final-code-reality-design-audit` | done |
| P94 GitHub CI/CD hardening | `p94-github-ci-release-hardening` | done |
| P95 架构/API/工程加固 | `p95-architecture-api-engineering-hardening` | done |
| P96 Public docs/README 产品化 | `p96-public-docs-readme-productization` | done |
| P97 默认本地配置文件修正 | `p97-default-local-config-file` | done |
| P98 运行时边界与前端复用加固 | `p98-runtime-hardening-and-code-reuse-cleanup` | done |
| P99 初始发布版本号 | `p99-initial-release-version` | done |
| P100 本地源码最终验收 | `p100-local-source-final-acceptance` | done |
| P101 本地配置路径统一 | `p101-unify-local-config-path` | done |
| P102 产品验收审计 | `p102-product-acceptance-audit` | done |
| P103 产品验收 UX 联动修复 | `p103-product-acceptance-ux-linkage-fixes` | done |
| P104 产品操作联动验收 | `p104-full-product-operation-linkage-acceptance` | done |
| P105 v0.1.1 发布版本 | `p105-release-v0-1-1` | done |
| P106 v0.1.2 package scan 修复 | `p106-release-v0-1-2-package-scan-fix` | done |
| P107 README 双语与中文功能介绍 | `p107-readme-bilingual-product-intro` | done |

当前机器可读进度以 `openspec/PROGRESS.md` 为准。P19–P24 为已交付能力的历史状态校准；仓库中未补建对应 `openspec/changes/archive/` 包，后续若要追补审计材料，应作为独立治理 change 处理。P33–P40 当前计划内功能队列已完成；P41 已完成 P40 后路线图治理；P42 已完成用户决策工作台；P43 已完成数据质量可观测；P44 已完成本地安装诊断与打包；P45 已完成 P44 后路线图治理；P46 已完成本地知识导入治理；P47 已完成决策闭环解释并归档；P48 已完成数据源质量回归包并归档；P49 已完成本地发布与升级体验并归档；P50 已完成 P49 后治理与验收路线图并归档；P51 已完成 P19-P24 审计证据包并归档；P52 已完成项目验收门禁矩阵并归档；P53 已完成验收执行与发布候选材料并归档；P54 已完成发布交付与可重复性加固并归档；P55 已完成前端全功能真实验收与设计审查并归档；P56 已完成 UI 验收阻断与产品化设计修复并归档；P57 已完成产品体验打磨总规划并归档；P58 已完成今日工作台重构并归档；P59 已完成决策解释体验重构并归档；P60 已完成组合、风险与数据质量体验重构；P61 已完成治理和运维页面产品化；P62 已完成设计系统与可访问性验收；P63-P107 均已完成并归档至对应 `openspec/changes/archive/`。P107 已补齐 README 双语入口与中文项目功能介绍；不新增运行时投资能力、发布包或 GitHub Release 声明。当前活跃 change 为无；下一建议阶段为无。

## 实现约束（摘要）

- Go 1.22+、Eino、SQLite、VecLite、DeepSeek；前端 React + Vite + TS。
- **禁止自动交易**；DeepSeek 不写最终裁决。
- 关键动作写 `audit_events`（见 `docs/data-model.md`）。

## 归档检查清单

以下清单适用于有 `openspec/changes/<id>/` 的标准变更。P19–P24 为已交付但无 archive 包的历史状态校准；如需追补审计材料，应新建独立治理 change，不应伪造历史 archive。

- [ ] delta 已合并到对应 `docs/*.md`
- [ ] 对应 `docs/*.md` 的版本或最后更新日期已同步
- [ ] `docs/development-plan.md` 相关任务已勾选
- [ ] 验收命令在 change 的 `tasks.md` 中已通过
- [ ] `docs/GOVERNANCE.md` 活跃变更表已更新
- [ ] `openspec/PROGRESS.md` 已推进阶段、状态和下一阶段入口
