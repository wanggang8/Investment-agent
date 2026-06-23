# Proposal: P61 治理和运维页面产品化

## Summary

将 Rules、Audit、Notifications、Daily Reports、Daily Auto Run、Local Install、Local Knowledge、Settings 从工程工具式页面打磨为可扫描、可处置、边界清楚的治理和运维工作台。

## Motivation

P58-P60 已经把日常工作台、决策解释、组合/风险/数据质量页面调整为产品化 operational cockpit。剩余治理和运维页面仍存在若干体验问题：

- Rules 仍暴露较多 raw JSON 和规则 diff，缺少提案状态、样本、门禁和人工确认路径的优先级组织。
- Audit 主要是事件时间线，缺少首屏摘要、筛选概览和下一步定位。
- Notifications 仍像普通列表，未形成本地应用内收件箱和处理状态。
- Daily Reports / Daily Auto Run 尚未统一“本地纪律复盘”和“运行诊断”的入口层级。
- Local Install、Local Knowledge、Settings 的配置、诊断、脱敏和安全提示模式不够统一。

P61 的目标是降低工程感，让用户能回答：哪些治理事项需要我看、哪些运维状态阻断了日常判断、下一步应去哪里人工处理，以及哪些动作明确不会自动修复、自动确认、自动交易或外部推送。

## In Scope

- `/rules`：规则治理总览、提案卡片、样本/过拟合/守门人状态、人工确认边界、减少首屏 raw JSON。
- `/audit`：审计摘要、时间线分组、事件筛选和关联入口。
- `/notifications`：本地应用内通知收件箱、未读/严重程度/来源分类、处理状态、只读外推边界。
- `/daily-discipline/reports`：报告历史的信息架构，突出纪律复盘状态、证据覆盖、数据缺口和人工动作。
- `/daily-auto-run`：本地自动运行状态、最近执行、失败诊断、缺失前提、手动复验入口说明和安全边界。
- `/local-install`、`/local-knowledge`、`/settings`：统一诊断、配置、脱敏、草稿/预览/确认和安全提示模式。
- 新增前端 view model / mapper 测试，必要时调整页面组件测试和 Playwright smoke。
- 真实启动本地后端和前端，通过浏览器操作 P61 主要页面并采集桌面/390px 证据。

## Out of Scope

- 不新增券商接口、登录源、付费源、授权源、Level2 或高频源。
- 不新增自动交易、一键交易、代下单、外部推送、短信/邮件/第三方通知发送。
- 不新增自动确认、自动规则应用、自动修复、自动覆盖真实数据库或收益承诺。
- 不修改 SQLite schema，不新增 Eino workflow，不改变 LLM 只生成分析材料的边界。
- 不把 P61 结论表述为最终 release-ready；最终刷新留给 P63。

## Product Design Brief

P61 治理/运维页面应延续 P58-P60 的 operational cockpit 语言：首屏先给状态、原因、下一步人工动作和安全边界；详情再展示表格、时间线、配置草稿或诊断材料。页面风格应克制、密集、可扫描，适合本地投资纪律工具的重复使用，而不是营销页、聊天 demo 或工程调试台。交互级别为 full interactivity：继续使用现有后端 API/service、真实按钮状态和表单路径，不做静态 mock。

## Validation

- `npm test -- --run` 覆盖新增 view model 和 P61 页面测试。
- `npm --prefix web test`
- `npm --prefix web run build`
- `go test ./...`
- 启动真实本地后端和 Vite 前端，浏览器操作 P61 页面。
- 390px 和桌面 viewport 检查无页面级横向溢出。
- `bash scripts/e2e-smoke.sh`
- `openspec validate p61-governance-ops-productization --strict`
- `openspec validate --all --strict`
- `git diff --check`
- 敏感信息和 forbidden copy 扫描。
