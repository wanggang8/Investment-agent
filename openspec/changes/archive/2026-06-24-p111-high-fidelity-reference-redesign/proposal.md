# Proposal: P111 高保真参考图视觉重构

## Summary

P111 在 P110 轻量视觉系统升级之后，发起一次以已选第二方案图为唯一视觉真源的高保真产品 UI 重构。P111 不是继续调色或局部换肤，而是把 Investment Agent 的核心页面重构为接近参考图的 Calm Command Center：左侧高密度导航、顶部状态工具栏、大型纪律状态横幅、优先级人工动作队列、状态指标矩阵、资金/持仓快照、最近咨询进度、证据与规则快照，并把同一套模块语言延展到维护、证据、决策、风险、治理和运维页面。

参考图：

`/Users/vick/.codex/generated_images/019ef8a7-f5c0-7442-95b9-e72bcdc89de1/ig_05724f56eb7089ab016a3b9109e1848191a87e68883d0c9826.png`

## Why

P110 已完成全局 tokens、导航状态和少量共享 surface，但实际截图与用户选定的第二方案差距明显。差距不是某一个页面未完成，而是 P110 没有做到参考图级别的信息架构、组件结构和视觉密度复刻。

P111 的必要性是把“方向接近”升级为“可对照验收”：每个页面完成后必须与参考图的模块语言、密度、层级、状态色、边框、按钮、排版和响应式行为做截图对比。未达标的页面不得归档。

## What Changes

- 建立 P111 reference design system：top status toolbar、command sidebar、report hero、priority action queue、status metric cards、portfolio snapshot strip、progress tracker、evidence/rule checklist、ledger table surface、ops panel variants。
- 重构 `/` 今日纪律与 `/workbench` 为参考图级别的核心 cockpit，而不是 P110 的纵向卡片堆叠。
- 将同一套模块扩展到 `/data-quality`、`/risk-alerts`、`/evidence`、`/decision-loop`、`/positions`、`/consultation`、`/decisions/:id`、`/rules`、`/review`、`/audit`、`/notifications`、`/daily-discipline/reports`、`/daily-auto-run`、`/local-install`、`/local-knowledge`、`/settings`。
- 增加逐页视觉 QA：每个页面在桌面完成后必须生成截图、填写 reference mismatch ledger、修复 P0/P1/P2 mismatch 后再进入下一页。
- 保留现有 React/Vite/TypeScript 架构、service/API DTO、路由语义、人工确认流程和安全边界。

## In Scope

- 前端视觉与布局重构：AppLayout、Dashboard/Workbench、dashboard shared components、核心页面样式和必要页面结构调整。
- 新增或重构只读展示组件：reference shell、status toolbar、priority list、metric grid、snapshot strip、process tracker、evidence checklist、ledger surface。
- 页面级参考图对比门禁：桌面参考尺寸优先使用 1492 x 1068 或等价比例；同时保留 390px、768px、1280px reflow。
- 逐页 QA 资产：reference image、rendered screenshots、mismatch ledger、pass/fail matrix、console/reflow JSON。
- 测试覆盖：视觉结构 class/landmark、导航语义、安全边界、核心页面模块存在、forbidden affordance scan、redaction scan。

## Out Of Scope

- 不新增后端 API、SQLite schema、Eino workflow、LLM 能力、RAG/VecLite 能力、真实数据源或投资规则。
- 不改变最终裁决逻辑、规则提案逻辑、确认流程、审计语义或 release/package/version 行为。
- 不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、真实库覆盖或收益承诺。
- 不新增登录源、付费源、授权源、Level2、高频源或外部商业数据依赖。
- 不把 P111 视觉验收扩大为新的投资效果、收益准确性或物理第二机器验收声明。

## Product Design Brief

P111 的设计目标是“把用户选定的第二方案真正做进产品”。Investment Agent 必须看起来像一个成熟的本地投资纪律指挥台：冷静、密集、可扫描、有清楚的人工动作优先级和证据/规则可信度，而不是一个普通后台或卡片堆叠页面。

所有页面都应从参考图抽取同一套视觉语言：

- 左侧导航有清晰分组、图标、active rail、底部本地模式状态。
- 顶部工具栏承载页面标题、日期、本地模式、数据截至和刷新类动作。
- 首屏模块不是营销 hero，而是纪律报告/状态横幅。
- 行动队列使用编号、优先级、右侧动作按钮、细分线和进度状态。
- 指标区使用小图标、强数值、状态标签和底部细项。
- 证据/规则/进度使用 checklist、timeline 或 ledger，而不是大段文本堆叠。

P111 每页必须有对比证据。若某页因数据结构与参考图不同无法一比一复刻，必须在 mismatch ledger 中写明“功能必要偏差”，并用同一套视觉语言完成等价设计。

## Validation

- `openspec validate p111-high-fidelity-reference-redesign --strict`
- `npm --prefix web test -- --run`
- `npm --prefix web run build`
- `go test ./...`
- `go vet ./...`
- 启动真实本地后端和 Vite 前端。
- 每个 P111 覆盖页面采集桌面截图并填写 mismatch ledger；P0/P1/P2 mismatch 必须修复后才能标记页面完成。
- 核心路由采集 390px、768px、1280px reflow 截图/JSON，确认无页面级横向溢出。
- Forbidden affordance scan：不得新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、收益承诺。
- Sensitive/redaction scan：不得暴露完整 key、私有路径、SQL、完整 prompt、raw vendor payload、本地数据库路径或 raw stack。
- `openspec validate --all --strict`
- `git diff --check`
