# Proposal: P115 真实用户场景全链路验收

## Why

P104 已证明一批产品操作可以通过 HTTP API 写入、SQLite 读回、下游页面联动和安全负证据验证。但 P104 更像“操作联动核心门禁”，覆盖的是代表性操作集合，不是尽可能贴近真实用户全天候使用的完整场景验收。

用户当前希望在 P114 UI 重构后，进一步确认产品不是只“看起来能用”，而是在真实使用场景中，页面操作、后端 API、SQLite 数据、审计事件、跨页面状态、错误降级和安全边界都能闭环。P115 因此需要把验收对象从单页面/单接口扩展为多阶段用户旅程。

## What

- 建立 P115 真实用户场景验收矩阵，覆盖首次使用、日常纪律、主动咨询、组合维护、线下交易、风险处置、数据质量、规则治理、证据/RAG、本地知识、日报/自动运行、通知、复盘、设置与运维、异常降级、移动端等场景。
- 新增或扩展 repeatable runner，使用隔离临时 SQLite 和本地 backend/frontend，按真实用户旅程执行操作。
- 每个场景必须记录：
  - 用户入口与操作步骤。
  - API 请求/响应或浏览器交互证据。
  - SQLite 字段级 readback。
  - 跨页面或跨 endpoint 联动读回。
  - 审计/通知/风险/确认等副作用。
  - 禁止能力负证据。
- 对已有 P104 runner 能覆盖的场景复用其能力；对 P104 未覆盖的场景补充新的 API/browser/SQLite 检查。
- 生成 P115 acceptance record 和证据包，明确哪些场景 fresh pass、哪些因外部 provider/LLM/key/network 等条件记录为 degraded 或 scoped。

## Out Of Scope

- 不新增投资规则、交易能力、自动执行能力或收益承诺。
- 不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认或自动规则应用。
- 不处理 Docker、安装器、发布包刷新、GitHub Release 或物理第二机器复验。
- 不把 stub/fixture/deterministic-local 验收冒充真实外部 provider 可用性。
- 不把单一标的、单一浏览器、单一 happy path 冒充全量现实世界通过。

## Acceptance

- `openspec validate p115-real-user-scenario-acceptance --strict` 通过。
- P115 场景矩阵覆盖所有可见产品路由和主要用户操作类别。
- P115 runner 至少覆盖组合、持仓、交易记录、导入、咨询、确认、决策错误标注、决策闭环、证据、数据质量、风险、规则、通知、日报、复盘、本地知识、设置、settings 禁止规则/SOP 直接修改、browser 交互补齐、安全边界。
- 每个场景有 API/browser、SQLite、跨页面或下游 readback 证据。
- 失败/降级场景必须验证不会生成无依据建议、不会自动交易、不会自动确认、不会自动应用规则。
- 验收通过后运行 `go test ./...`、`npm --prefix web test -- --run`、`npm --prefix web run build`、`openspec validate --all --strict`、P92/P93 或其 P114 后等价真实性检查、`git diff --check`。
