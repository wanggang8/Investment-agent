# Proposal: P110 视觉系统重设计

## Summary

在 P102-P104 产品验收和操作联动门禁通过后，针对当前前端“功能可信但审美偏后台台账”的问题，发起一轮视觉系统重设计。P110 目标是把 Investment Agent 从可用的本地投资纪律工作台，提升为更成熟、更高级、更可扫描的个人投资纪律研究终端。

## Why

当前 UI 已经通过多轮 AI/真实浏览器验收，核心路径、API/SQLite/readback、安全边界和人工确认流程均稳定。最新 P102 截图显示，页面结构清楚、语义可信，但视觉语言仍偏早期后台：深色侧栏 + 白卡片 + 多处边框台账，缺少更强的信息节奏、密度控制、状态层级和高品质产品质感。

在功能验收稳定后进行独立视觉系统升级，可以降低重设计对业务逻辑的干扰；同时保留 P57-P63 已固化的产品定位：冷静、可审计、风险前置、人工决策、不刺激交易。

## What Changes

- 先生成三版独立视觉方向稿，围绕 Dashboard/Workbench 核心工作台建立视觉目标。
- 选定方向后，重塑前端视觉 tokens、导航、页面头部、状态摘要、行动队列、数据面板、表单、表格和关键解释区的审美一致性。
- 优先覆盖核心页面：Dashboard、Workbench、Consultation、Decision Detail、Positions、Data Quality、Risk Alerts、Evidence 和 Decision Loop。
- 复用现有 React/Vite/TS 架构、service/API DTO、UI primitives 和安全文案边界。
- 增加视觉回归截图、reflow、forbidden copy scan 和前端 test/build 验收记录。

## In Scope

- 视觉系统方向探索：三版桌面 mock，目标尺寸 1440 x 1024。
- `web/src/styles/global.css` 与前端 UI primitives 的视觉 token、排版、间距、边框、状态 tone 和响应式样式优化。
- `AppLayout` 导航与核心页面视觉层级优化，但不改变路由语义。
- Dashboard/Workbench 的首屏状态、下一步人工动作、信号摘要和详细驾驶舱视觉升级。
- Consultation、Decision Detail、Decision Loop、Evidence 的解释链路密度和阅读层级优化。
- Positions、Risk Alerts、Data Quality 的维护、处置和质量面板审美升级。
- 390px、768px、1280px viewport 的截图与无页面级横向溢出检查。
- 安全文案扫描：不出现交易自动化、收益承诺、外部推送或敏感信息暴露。

## Out Of Scope

- 不新增后端 API、SQLite schema、Eino workflow、LLM 能力、数据源能力或检索能力。
- 不改变投资规则、最终裁决逻辑、分析 prompt、确认流程、审计语义或发布版本号。
- 不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、真实库覆盖或收益承诺。
- 不新增登录源、付费源、授权源、Level2、高频源或外部商业数据依赖。
- 不处理 Docker、安装/升级/卸载、GitHub Release、分发包刷新或物理第二机器复验。
- 不把视觉改造结论扩大为新的业务功能验收或投资效果承诺。

## Product Design Brief

Investment Agent 应呈现为“冷静的投资纪律研究终端”：比普通后台更精致，比券商行情盘更克制，比 AI chat demo 更结构化。核心用户每天需要回答四个问题：今天能不能动、为什么、需要做什么人工动作、数据和规则是否可信。

视觉方向必须保持高密度但可扫描，使用清晰排版、稳定布局、克制色彩和明确状态语义。风险、禁止动作、数据不足、冻结观察和人工确认边界的视觉优先级必须高于收益想象。所有行动入口均为本地查看、维护、记录或人工复核，不得暗示系统可自动执行交易或自动应用规则。

交互级别为 full interactivity：最终实现必须基于现有真实前端、真实 service/API 和浏览器验收，不做静态 mock 交付。

## Validation

- `npm --prefix web test -- --run`
- `npm --prefix web run build`
- `go test ./...`
- 启动真实本地后端和 Vite 前端，采集核心路由桌面/移动截图。
- 390px、768px、1280px reflow 检查，无页面级横向溢出。
- Forbidden copy scan：券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、收益承诺等不得出现为能力入口。
- 敏感信息扫描：不渲染完整 key、私有路径、SQL、完整 prompt、raw vendor payload、本地数据库路径或 raw stack。
- `openspec validate p110-visual-system-redesign --strict`
- `openspec validate --all --strict`
- `git diff --check`
