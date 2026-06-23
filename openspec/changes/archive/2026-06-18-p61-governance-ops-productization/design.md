# Design: P61 治理和运维页面产品化

## Current State

P61 页面多数已经具备基础能力和测试，但交互层级偏工程化：

- `/rules` 可读取当前规则、规则提案、用户确认和最终确认，但首屏仍以规则 JSON、阈值和提案详情为主。
- `/audit` 已有审计时间线和错误态，但缺少总览摘要、筛选状态和“为什么该看这条”的入口层。
- `/notifications` 有轮询、未读数、标记已读和风险预警链接，但缺少本地应用内收件箱结构和边界说明。
- `/daily-discipline/reports` 能展示历史报告卡片，但还没有把报告状态、证据覆盖、趋势和下一步人工动作组合成首屏。
- `/daily-auto-run` 能显示本地自动运行状态和诊断，但需要更清楚地区分配置、最近运行、失败原因、复验路径和安全边界。
- `/local-install`、`/local-knowledge`、`/settings` 已有表单、诊断、预览和本地刷新能力，但配置/诊断/脱敏说明模式不统一。

## Architecture

P61 继续采用 P58-P60 的前端 view model 模式：

- 在 `web/src/features/governance/` 下新增轻量 view model：
  - `rulesGovernanceModel.ts`
  - `auditOpsModel.ts`
  - `notificationInboxModel.ts`
  - `dailyOpsModel.ts`
  - `localOpsModel.ts`
- View model 只消费现有 service DTO 和页面已有本地 UI 状态，不读取 SQLite、VecLite、localStorage、sessionStorage、本地文件或临时配置。
- 页面组件负责调用现有 service、处理 loading/error/message，再把 DTO 传入 view model。
- 写入动作继续走现有 service：规则确认、通知标记已读、市场刷新、本地知识 validate/confirm。P61 只调整文案、布局、状态解释和安全边界，不新增后端写入。
- 样式复用 P58-P60 的 `.daily-hero`、`.cockpit-card`、`.cockpit-grid`、`.metric-grid`、`.action-row` 等 operational tokens；必要时新增少量 P61 专用 class。

## Page Design

### Rules

首屏展示“规则治理状态”：

- 当前规则版本、提案数量、待用户确认、待最终确认、门禁失败/风险提案数量。
- 下一步人工动作：查看待确认提案、复核门禁失败、查看审计、查看规则效果验证。
- 安全文案：最终确认仍是本地规则治理动作，不会交易、不外推、不绕过守门人。

提案卡片突出 title、status、reason、sample_count、overfit_risk、guardrail_decision、audit_summary 和关联链接。Raw JSON 只保留在折叠/次级区域或结构化摘要里，不作为首屏主内容。

### Audit

首屏展示“审计检查状态”：

- 事件总数、最近事件、关键来源、错误/风险/规则/通知相关计数。
- 时间线按 event_type 或 source_type 做可扫描分组，保留现有事件详情。
- 空态和错误态必须说明“当前没有匹配审计事件”或“只能展示安全错误摘要”，不显示 raw stack。

### Notifications

首屏展示“本地通知收件箱”：

- 未读数、严重程度分布、风险/数据源/每日运行/规则等来源分布。
- 明确“本地应用内通知”，不承诺短信、邮件、Webhook、第三方推送或外部通知。
- 标记已读/全部已读继续使用现有 service，只表达为本地处理状态。

### Daily Reports And Daily Auto Run

Daily Reports 首屏回答“最近纪律复盘是否可用”：

- 最新日期、状态、证据覆盖、趋势计数、缺口和下一步人工动作。
- 报告卡片保留详情入口，避免像数据 dump。

Daily Auto Run 首屏回答“本地自动运行是否健康”：

- enabled/status、计划时间、最近/下次执行、失败原因、缺失前提、关联决策/通知/审计。
- 文案必须说明默认关闭或显式启用；失败诊断只指导人工复验，不承诺自动修复。

### Local Install, Local Knowledge, Settings

三类系统页面统一为“配置与诊断”体验：

- Local Install：配置草稿、关键命令、上传摘要、失败步骤、下一步人工复验。
- Local Knowledge：导入草稿、脱敏预览、索引计划、确认理由和本地事实写入边界。
- Settings：能力圈、系统状态、数据源健康、市场刷新和错误摘要。

这些页面允许本地草稿、预览、上传本地 JSON 摘要或现有安全写入，但不得暗示自动覆盖真实库、自动修复、自动确认或自动规则应用。

## Testing Strategy

- View model 单测覆盖成功、空态、错误/降级、未知状态、安全文案和 forbidden copy。
- 页面测试覆盖首屏状态、下一步动作、现有写入 service 调用和安全边界。
- Playwright smoke 打开 P61 页面，检查关键标题、状态和 forbidden copy。
- 真实 UI 验收启动 Go 后端和 Vite 前端，浏览器实际访问并操作：
  - `/rules` 查看规则治理和提案。
  - `/audit` 查看审计摘要/时间线。
  - `/notifications` 查看本地收件箱并在安全条件下标记已读。
  - `/daily-discipline/reports` 查看报告历史。
  - `/daily-auto-run` 查看状态与关联入口。
  - `/local-install` 上传诊断摘要 fixture。
  - `/local-knowledge` 执行 validate；confirm 仅在测试数据安全时执行。
  - `/settings` 查看系统状态并在 fixture 环境下触发市场刷新。

## Risks

- P61 页面多，容易扩大范围。控制方式：只做信息架构、view model 和现有 service 的 UI 调整，不改后端契约。
- Existing tests 可能依赖旧文案。控制方式：先写 view model/page tests，再同步 E2E 文案。
- 真实本地库状态不可控。控制方式：浏览器验收记录实际可操作路径；无法执行的写入动作不伪造，通过 fixture/E2E 覆盖。

