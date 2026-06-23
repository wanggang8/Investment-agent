# P60 组合、风险与数据质量体验重构

## Why

P57 将产品体验打磨拆成 P58-P63；P58 已完成每日工作台，P59 已完成决策解释链路。当前 `/positions`、`/risk-alerts` 和 `/data-quality` 已具备基础能力，但用户仍需要自己判断：

- 账户与持仓现在处于初始化、校准、维护、导入还是纠错状态。
- 风险预警哪些待看、哪些处理中、哪些已记录，下一步应做什么人工动作。
- 数据质量中的 source health、证据/RAG、LLM 和受影响工作流之间如何互相影响。

P60 要把这三页打磨成日常维护与处置页面：用户先看到当前状态、风险/质量影响、下一步人工动作，再进入具体本地事实记录或只读诊断。该阶段不新增后端能力，不改变 API、SQLite schema、Eino 工作流、交易边界或规则裁决，只重组前端信息架构、展示模型、交互、测试和验收证据。

## What Changes

- 重构组合与持仓维护体验：
  - `/positions` 首屏展示组合状态、现金/持仓/高风险比例、维护阶段和下一步人工动作。
  - 明确区分初始化、校准、持仓编辑、线下交易补记、批量导入、错误修正。
  - 写入动作继续只调用现有本地账户事实 API，并保留“不连接券商、不自动交易”的显式确认边界。
- 重构风险预警处置体验：
  - `/risk-alerts` 从列表改为处置队列，按待看、处理中、需复盘、已记录等用户语义分组。
  - 风险详情展示触发原因、禁止动作、建议人工动作、关联决策/报告/通知/审计和 SOP 当前状态。
  - 生命周期按钮继续只调用现有本地 SOP API；不得出现交易、推送或自动确认入口。
- 重构数据质量体验：
  - `/data-quality` 先展示质量总览、source health、证据/RAG、LLM、影响范围和下一步本地导航。
  - 把降级、缺证据、索引/LLM 问题和受影响决策串成可扫描质量面板。
  - 继续只读展示，不触发自动刷新、自动修复、自动确认或规则应用。
- 更新测试与真实 UI 验收：
  - Vitest 覆盖 view model、三页空态/降级态/安全文案/链接/写入边界。
  - Playwright 或浏览器验收覆盖 `/positions`、`/risk-alerts`、`/data-quality` 桌面与 390px 移动端。
  - 使用本地真实后端和 Vite 前端操作关键维护/处置路径，采集截图和安全扫描证据。

## Scope

- 前端 React/Vite/TypeScript：
  - `web/src/pages/PortfolioPage.tsx`
  - `web/src/pages/RiskAlertPage.tsx`
  - `web/src/pages/DataQualityPage.tsx`
  - 相关测试、CSS、E2E smoke
  - 必要的轻量 view model / presentational components
- P60 UI 验收报告、截图资产、OpenSpec 和治理文档更新。

## Out of Scope

- 不修改 SQLite schema、HTTP API、Eino 工作流、Go 后端业务逻辑、真实数据 collector、LLM 裁决逻辑或规则引擎。
- 不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复承诺、自动覆盖真实库、收益承诺。
- 不新增登录源、付费源、授权源、Level2 或高频源。
- 不把 `/positions` 改造成券商账户、资金流水同步或自动仓位管理产品。
- 不把 `/risk-alerts` 改造成通知外推、自动处置或交易执行系统。
- 不把 `/data-quality` 改造成一键修复、自动刷新真实源或自动覆盖本地库的运维工具。
- 不提前执行 P61 治理/运维产品化、P62 组件系统或 P63 全量发布刷新。

## Validation

- 方案阶段：
  - `openspec validate p60-portfolio-risk-data-quality-experience --strict`
  - `openspec validate --all --strict`
  - `git diff --check`
  - 敏感信息扫描无 key、完整 prompt、私有 SQLite、raw vendor payload、私有路径或临时配置泄露。
  - 子 agent 方案复审无 Critical / Important，且复审覆盖 Product Design brief、P57-P59 alignment、信息架构、真实 UI 验收和安全边界。
- 实现阶段：
  - `npm test -- --run src/pages/PortfolioPage.test.tsx src/pages/RiskAlertPage.test.tsx src/pages/DataQualityPage.test.tsx`
  - `npm test`
  - `npm run build`
  - `go test ./...`，若无后端修改仍执行以证明集成基线。
  - `bash scripts/e2e-smoke.sh` 或 `E2E_BASE_URL=<local vite url> npm run test:e2e`
  - 真实启动本地后端和 Vite 前端，通过浏览器操作 `/positions`、`/risk-alerts`、`/data-quality`。
  - 390px 移动端检查无页面级横向溢出，采集 P60 桌面/移动截图。
  - 扫描 UI 文案，确认无自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、收益承诺入口或暗示。
- 归档前：
  - 子 agent 执行后复审无 Critical / Important。
  - archive 后提交前复审无 Critical / Important。
