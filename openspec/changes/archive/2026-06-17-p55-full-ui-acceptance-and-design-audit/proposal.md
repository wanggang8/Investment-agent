# P55 前端全功能真实验收与设计审查

## Why

P53 执行了项目级验收门禁，但主要依赖自动化 smoke、CLI 和命令级验证。用户明确要求：必须真实启动项目，通过前端 UI 操作来验收所有功能，并使用前端设计/UI 相关能力审查界面是否需要优化。

P55 需要补齐“真实 UI 操作验收 + Product Design 审查”证据，避免仅凭自动化 smoke 或发布材料判断前端体验已充分验收。

## What Changes

- 本地启动真实后端 server 与 Vite 前端。
- 使用浏览器逐页面/逐功能操作并截图留证。
- 覆盖当前路由入口：
  - `/`
  - `/workbench`
  - `/decision-loop`
  - `/data-quality`
  - `/positions`
  - `/consultation`
  - `/decisions/:decisionId`
  - `/evidence`
  - `/rules`
  - `/audit`
  - `/notifications`
  - `/risk-alerts`
  - `/risk-alerts/:alertId`
  - `/daily-auto-run`
  - `/daily-discipline/reports`
  - `/daily-discipline/reports/:reportId`
  - `/review`
  - `/local-install`
  - `/local-knowledge`
  - `/settings`
- 新增 `docs/release/ui-acceptance-2026-06-17.md`，记录真实 UI 操作验收结果、截图路径、问题和阻断项。
- 新增 `docs/release/ui-design-audit-2026-06-17.md`，使用 Product Design audit 框架记录 UX、视觉、一致性和可访问性风险。
- 新增 `docs/release/ui-audit-assets/2026-06-17-p55/` 保存本次审查截图。
- 更新治理、进度、开发计划和文档地图。
- 增加 OpenSpec 行为摘要，约束未来不得把命令级 smoke 等同于完整 UI 功能验收。

## Scope

- 真实启动本地 server 和前端。
- 使用前端 UI 操作验收页面可达性、关键控件、状态展示、导航、表单/按钮、错误/空态、只读安全边界；对包含多个功能区的页面按功能区分别记录。
- 做设计审查并给出优化建议。
- 仅新增或修改文档、截图资产、OpenSpec change、OpenSpec specs 摘要和进度状态。

## Out of Scope

- 在 P55 中直接修改前端代码或样式。
- 新增自动化 UI runner。
- 改写 P53 release_ready 结论。
- 新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复承诺、自动覆盖真实库、收益承诺、登录源、付费源、授权源、Level2 或高频源。
- 提交完整 key、临时 SQLite、raw HTTP 响应、完整 prompt、原始 SQL、私有路径或未脱敏供应商 payload。

## Validation

- `openspec validate p55-full-ui-acceptance-and-design-audit --strict`
- `openspec validate --all --strict`
- `git diff --check`
- 本地 server/web 启动并可访问。
- 浏览器逐页面验收截图和结果记录完整。
- Product Design audit 结果记录完整。
- 子 agent 计划复审、执行后复审和提交前复审均无 Critical / Important。
