# P39: 前端完整用户旅程与全路径 E2E

## Why

P33-P38 已补齐账户/持仓录入、真实公开数据覆盖、风险预警、规则进化效果验证、真实 LLM 质量和 RAG / VecLite 检索质量。当前能力已经分散在多个页面和本地任务里，但仍缺少一个从空库初始化到每日使用、主动咨询、用户确认、复盘和规则治理的浏览器级完整验收。P39 需要把产品路径串起来，确保关键页面在真实浏览器中可达、可理解、无 console error，并继续守住不自动交易和不自动应用规则边界。

## What Changes

- 增加前端完整用户旅程 E2E：空库 onboarding、配置/账户初始化、市场/证据刷新、每日纪律、主动咨询、确认记录、风险预警、复盘和规则治理。
- 扩展 Playwright 验收覆盖降级路径：缺账户、缺市场、VecLite/RAG 降级、LLM 降级、证据不足、能力圈外、规则提案待确认。
- 增加前端可用性检查：窄屏、基础可访问性、关键页面无 console error / unhandled error。
- 保持 Vitest 单测与 Playwright E2E 分层，避免测试互相收集或污染本地 SQLite。
- 将 P34/P35/P38 的源健康、风险预警和检索质量展示纳入跨页验收。

## Out of Scope

- 不新增券商 API、不自动交易、不外部推送。
- 不自动应用规则、不绕过守门人审计或用户最终确认。
- 不扩大数据源范围，不引入登录、付费、授权、Level2 或高频源。
- 不把 E2E fixture 当作真实投资建议或收益承诺。

## Impact

- 前端：页面状态、跨页入口、错误/空态、窄屏和基础 a11y。
- 测试：Playwright E2E fixture、browser smoke、console/a11y/no-auto-trading 断言。
- 后端/CLI：仅在必要时补充本地 smoke seed 或只读 fixture，不改变交易边界。
- 文档：`docs/frontend-contract.md`、`docs/development-plan.md`、`docs/configuration.md`、`docs/ops-local-scheduler.md` 的验收与运行说明。
