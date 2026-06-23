# P59 决策解释体验重构

## Why

P57 已把产品体验打磨拆分为 P58-P63，P58 已完成每日工作台重构。当前 `/consultation`、`/decisions/:decisionId`、`/evidence` 和 `/decision-loop` 虽然功能可用，但仍更像技术模块集合：用户需要自行把主动咨询输入、最终裁决、证据、LLM 分析、规则链、审计与闭环记录拼成完整决策故事。

P59 要把这些页面串成可理解的解释链路：先说明用户问了什么和系统产出了什么，再说明为什么、依据是什么、哪些动作被禁止、哪些动作只能人工复核，以及后续闭环如何只读追踪。该阶段不新增后端能力，不改变裁决、交易或规则边界，只重组前端信息架构、展示模型、交互、测试和验收证据。

## What Changes

- 重构主动咨询体验：
  - `/consultation` 明确输入假设、生成建议、查看解释路径和只读/人工复核边界。
  - 真实 LLM consultation 成功后可继续查看生成的决策详情，失败时展示可恢复错误和安全说明。
- 重构决策详情解释体验：
  - `/decisions/:decisionId` 首屏展示最终裁决、禁止动作、可选人工动作、数据可信度、关键原因和安全边界。
  - Evidence、LLM、rules、audit、expected return 和 confirmation 分层展示，长 trace 默认折叠或降为二级详情。
  - Nullable/missing DTO 必须保持安全空态，不崩溃、不误判为可执行。
- 重构证据与闭环导航：
  - `/evidence` 增强证据可信度说明、来源等级、验证状态和可到达的决策解释入口。
  - `/decision-loop` 用只读时间线串联建议、确认、线下记录、风险、复盘和审计，并保留缺口说明。
- 更新测试与真实 UI 验收：
  - Vitest 覆盖真实 LLM-like DTO、空态/降级态、链接可达、安全文案和禁止执行入口。
  - Playwright 或浏览器验收覆盖 `/consultation`、`/decisions/:decisionId`、`/evidence`、`/decision-loop` 桌面与 390px 移动端。
  - 使用本地真实后端和前端执行一次真实 LLM consultation 验收，不提交临时 key 或私有配置。

## Scope

- 前端 React/Vite/TypeScript：
  - `web/src/pages/DecisionDetailPage.tsx`
  - `web/src/components/decision/DecisionTrace.tsx`
  - `web/src/pages/EvidencePage.tsx`
  - `web/src/components/evidence/EvidenceTable.tsx`
  - `web/src/pages/DecisionLoopPage.tsx`
  - 相关 service/type 使用、测试、CSS 和 E2E smoke
  - 必要的轻量 view model / presentational components
- P59 UI 验收报告、截图资产、OpenSpec 和治理文档更新。

## Out of Scope

- 不修改 SQLite schema、HTTP API、Eino 工作流、Go 后端业务逻辑、真实数据 collector 或 LLM 裁决逻辑。
- 不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复承诺、自动覆盖真实库、收益承诺。
- 不新增登录源、付费源、授权源、Level2 或高频源。
- 不提前执行 P60 组合/风险/数据质量体验重构、P61 治理/运维产品化、P62 组件系统或 P63 全量发布刷新。
- 不把主动咨询改成开放式聊天产品；LLM 仍只生成分析材料，最终裁决和安全边界仍由现有规则链与后端结果决定。

## Validation

- 方案阶段：
  - `openspec validate p59-decision-explainability-experience --strict`
  - `openspec validate --all --strict`
  - `git diff --check`
  - 敏感信息扫描无 key、完整 prompt、私有 SQLite、raw vendor payload 或临时配置泄露。
  - 子 agent 方案复审无 Critical / Important，且复审覆盖 Product Design brief、信息架构、真实 LLM 验收和安全边界。
- 实现阶段：
  - `npm test -- --run src/pages/DecisionDetailPage.test.tsx src/pages/EvidencePage.test.tsx src/pages/DecisionLoopPage.test.tsx`
  - `npm test`
  - `npm run build`
  - `go test ./...`，若无后端修改仍执行以证明集成基线。
  - `E2E_BASE_URL=<local vite url> npm run test:e2e`
  - 真实启动本地后端和 Vite 前端，通过浏览器操作 `/consultation`、生成后的 `/decisions/:decisionId`、`/evidence`、`/decision-loop`。
  - 390px 移动端检查无页面级横向溢出，采集 P59 桌面/移动截图。
  - 扫描 UI 文案，确认无自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、收益承诺入口或暗示。
- 归档前：
  - 子 agent 执行后复审无 Critical / Important。
  - archive 后提交前复审无 Critical / Important。
