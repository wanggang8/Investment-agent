# Proposal: P63 全量真实 UI 回归与发布状态刷新

## Summary

在 P58-P62 产品体验打磨完成后，真实启动本地后端与前端，执行全路由 UI 回归、真实 LLM consultation、错误/降级/移动端/桌面端/安全边界扫描，并刷新发布候选与交付口径。

## Why

P53 曾基于当时的 G0-G9 验收给出 `release_ready`，但 P55 后续真实 UI 验收发现过真实 LLM 决策详情 nullable DTO 崩溃；P56-P62 又连续改造了导航、核心工作台、决策解释、组合/风险/数据质量、治理运维页面和设计系统可访问性。

因此，P63 不能只沿用 P53 的发布状态，也不能只引用 P58-P62 的局部验收。需要在产品体验打磨完成后重新跑一次全量真实 UI 回归和发布门禁，明确当前 commit 是否可以恢复或刷新 `release_ready`，并把任何 degraded / blocked / waiver 记录进新的发布材料。

## What Changes

- 新增 P63 全量 UI 验收记录，覆盖真实后端、真实 Vite 前端、全主要路由浏览器操作、三视口 reflow、console/page error、安全文案和敏感信息扫描。
- 执行真实 LLM-backed consultation；成功则记录模型、质量门禁和新决策详情可打开，失败则按 P52 分类记录为外部依赖降级或阻断对应声明。
- 重新执行 P52 G0-G9 项目验收门禁，并记录每个 gate 的结果、命令、产物、降级原因和发布影响。
- 新增刷新后的 release candidate 文档和 handoff 文档，明确当前状态为 `release_ready`、`release_degraded` 或 `blocked`。
- 更新 P57 产品体验路线图、开发计划、文档地图和进度文档，说明 P63 结论及其边界。

## In Scope

- 使用现有本地后端、前端、CLI、测试和脚本执行验收；必要时新增/调整 Playwright 验收脚本或小型辅助脚本，只服务于回归证据采集。
- 浏览器操作覆盖核心路径：Dashboard、Workbench、Consultation、Decision Detail、Evidence、Decision Loop、Positions、Data Quality、Risk Alerts、Risk Alert Detail、Rules、Audit、Notifications、Daily Reports、Daily Report Detail、Daily Auto Run、Review、Local Install、Local Knowledge、Settings。
- 390px、768px、1280px viewport 下检查页面级横向溢出、首屏状态可读性、关键按钮/链接可达性和可访问名称。
- 记录 console error、page error、HTTP 失败、外部依赖失败、真实 LLM 质量/解析状态和 UI 安全降级表现。
- 刷新 release candidate / handoff 材料，并清楚说明 P53 与 P63 的关系。

## Out of Scope

- 不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、自动覆盖真实库、收益承诺、登录源、付费源、授权源、Level2 或高频源。
- 不新增业务功能、后端 API、SQLite schema、Eino workflow 或新的 LLM 能力；若验收发现阻断问题，只允许在既有前端/后端行为边界内做最小缺陷修复，仍不得新增上述能力。
- 不把外部公开源或模型供应商的一次通过表述为未来可用性承诺。
- 不覆盖或迁移真实用户数据库；真实验收使用临时 SQLite 或明确只读/临时配置。
- 不做新的视觉改版；P63 以回归验收和发布材料刷新为主，发现产品/UI 问题时只修复阻断级缺陷，非阻断项记录为后续 backlog。

## Product Design Brief

产品：Investment Agent 本地投资纪律工作台。

视觉/设计来源：现有代码、P57 路线图、P58-P62 产品化 UI 与设计系统 primitives；没有外部 Figma 或新视觉目标。

交互等级：full interactivity。必须真实启动本地后端和 Vite 前端，通过浏览器操作 UI 验收；真实 LLM 使用用户提供的测试配置或现有配置，失败需按外部依赖分类，不伪装成功。

## Validation

- `openspec validate p63-full-ui-regression-release-refresh --strict`
- `openspec validate --all --strict`
- `git diff --check`
- `npm --prefix web test`
- `npm --prefix web run build`
- `go test ./...`
- `bash scripts/e2e-smoke.sh`
- P52 G0-G9 全门禁命令
- 真实后端 + Vite 前端全路由浏览器验收
- 真实 LLM consultation UI journey
- 390px / 768px / 1280px 截图或等价浏览器证据
- 安全文案和敏感信息扫描
