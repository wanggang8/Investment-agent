# P39 Design

## Context

当前前端已经具备 Dashboard、Portfolio、Evidence、Decision、Review、Rules、Risk Alerts、Daily Discipline Report、Settings 等主要页面。P39 的核心不是再扩业务能力，而是用浏览器级验收把这些能力串成完整用户旅程，并让降级、安全边界和跨页入口变得稳定可测。

## Goals

1. 覆盖空库用户从初始化到第一份每日纪律报告的路径。
2. 覆盖主动咨询、决策详情、用户线下确认记录、审计追踪和复盘入口。
3. 覆盖规则提案治理、守门人/最终确认边界和不自动应用规则。
4. 覆盖 P34/P35/P38 的 source health、risk alert、retrieval quality 降级展示。
5. 在桌面与窄屏下检查关键页面无 console error、基础可访问性可用、无自动交易入口。

## Non-Goals

- 不实现真实交易或券商集成。
- 不把 Playwright 测试变成端到端外部公网依赖。
- 不重写前端架构或引入大型 UI 框架。
- 不用 E2E 替代单元测试；Vitest 仍覆盖组件逻辑，Playwright 覆盖浏览器旅程。

## Approach

### 1. Stable E2E Fixture

建立临时 SQLite + 本地配置 + deterministic seed。fixture 应包含账户、持仓、市场快照、证据、风险预警、规则提案、daily report 和 retrieval quality 示例。所有 fixture 只读或本地写入，不接外部服务。

### 2. Journey Specs

按真实用户顺序组织 Playwright 测试：首次进入、初始化、刷新/生成、查看报告、咨询、确认、审计、复盘、规则治理。每段测试只断言目标页面状态和安全边界，避免 brittle 文案过度绑定。

### 3. Degraded Paths

构造缺账户/缺市场/索引缺失/证据不足/LLM 降级/能力圈外/规则提案待确认等状态，确认前端能解释原因并提供安全下一步。

### 4. Safety and UX Checks

每个关键页面运行 console error 捕获、no-auto-trading 文案/入口断言、基础可访问性检查和窄屏布局 smoke。测试不得读取 SQLite 或本地 VecLite 文件；只通过 API/页面状态观察。

## Risks

- E2E 易变慢或脆弱。P39 应优先稳定 fixture、少量关键路径、明确等待条件。
- 旅程过长会难定位失败。测试按用户阶段拆分，但共享 seed。
- 安全文案断言可能过宽。只断言禁止入口和关键边界，不要求每页重复同一句话。

## Verification

- `go test ./...`
- `npm --prefix web test -- --run`
- `npm --prefix web run build`
- Playwright E2E / browser smoke 命令
- `openspec validate p39-frontend-full-user-journey-e2e --strict`
- `openspec validate --all --strict`
