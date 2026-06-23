# Design: P63 全量真实 UI 回归与发布状态刷新

## Current State

P58-P62 已完成产品体验打磨的主要实施阶段：

- P58 将 Dashboard / Workbench 打磨为每日投资纪律 cockpit。
- P59 串联 Consultation、Decision Detail、Evidence、Decision Loop 的解释链路。
- P60 重构 Positions、Risk Alerts、Data Quality 的维护和处置体验。
- P61 产品化 Rules、Audit、Notifications、Daily Reports、Daily Auto Run、Local Install、Local Knowledge、Settings。
- P62 固化前端 UI primitives、状态 tone、键盘路径、可访问语义和三视口 reflow 门禁。

P53 的 `release_ready` 仍是历史验收结论；它不自动覆盖 P58-P62 后的新代码与 UI。P63 的职责是重新执行真实验收，并给出当前 commit 的发布状态。

## Approach

P63 采用“验收优先，阻断修复最小化”的方式：

1. 先建立验收计划和 change，明确全路由、全门禁、安全边界和发布材料。
2. 先由子 agent 复审计划，确认没有范围越界或验收遗漏。
3. 执行自动化门禁：OpenSpec、diff check、Go、Vitest、build、E2E smoke、P52 G0-G9。
4. 真实启动后端和前端，使用浏览器操作全主要路由，采集桌面/平板/移动端证据。
5. 执行真实 LLM consultation UI journey，并打开新决策详情；若外部服务失败，按 P52 分类记录 release impact。
6. 只修复阻断级运行时/UI 缺陷；非阻断产品建议记录到验收文档或后续 backlog。
7. 生成 P63 acceptance、release candidate、handoff 和必要截图/JSON 资产。
8. 执行后子 agent 复审；无 Critical / Important 后 archive；archive 后提交前再复审。

## Browser Coverage

P63 浏览器验收至少覆盖以下路由：

- `/`
- `/workbench`
- `/consultation`
- `/decisions/:decisionId` 或真实 consultation 生成的新决策详情
- `/evidence`
- `/decision-loop`
- `/positions`
- `/data-quality`
- `/risk-alerts`
- `/risk-alerts/:alertId`
- `/rules`
- `/audit`
- `/notifications`
- `/daily-auto-run`
- `/daily-discipline/reports`
- `/daily-discipline/reports/:reportId`
- `/review`
- `/local-install`
- `/local-knowledge`
- `/settings`

每个路由记录：

- 页面标题或主要 landmark 是否可见。
- 关键状态、下一步人工动作或空/错误/降级态是否可读。
- 关键按钮/链接是否可点击或可通过键盘到达。
- 390px、768px、1280px 是否无页面级横向溢出。
- console error、page error、关键 HTTP 失败是否存在。
- 是否出现禁止能力文案或敏感信息泄露。

## Real LLM Journey

真实 LLM journey 使用现有配置或用户提供的临时测试配置：

1. 打开 `/consultation`。
2. 输入代表性咨询内容，触发真实 consultation。
3. 等待 workflow 返回。
4. 打开新生成的 decision detail。
5. 检查最终裁决、LLM 材料、证据、规则、审计和 decision loop 链接。
6. 记录模型、base URL 是否脱敏、错误分类、质量门禁和发布影响。

若真实 LLM 返回 503、限流、额度、认证、模型不可用或网络失败，P63 不伪装为通过；按 P52 G7 规则记录为 degraded / blocked / waiver，并限制 release claim。

## Release Status Model

P63 release candidate 使用以下三态：

| 状态 | 含义 |
| --- | --- |
| `release_ready` | G0-G9 默认阻断门禁通过；真实 UI 全路由无阻断；真实 LLM/公开源 opt-in 通过或有不阻断 waiver；安全扫描通过 |
| `release_degraded` | 核心本地能力和安全边界通过，但存在外部依赖、current 数据质量或非阻断 UI 降级，需明确限制发布声明 |
| `blocked` | 默认阻断门禁失败、真实 UI 核心路径崩溃、安全/脱敏问题、自动交易等禁止能力入口、或无法解释的验收失败 |

## Artifacts

P63 预期新增或更新：

- `docs/release/acceptance/2026-06-18-p63-full-ui-regression.md`
- `docs/release/release-candidate-2026-06-18.md`
- `docs/release/release-handoff-2026-06-18.md`
- `docs/release/ui-audit-assets/2026-06-18-p63/browser-results.json`
- `docs/release/ui-audit-assets/2026-06-18-p63/*.png`
- 必要时新增 `web/e2e/` 中的 P63 全路由验收用例或辅助脚本

## Risks

- 真实 LLM 或公开源外部依赖不稳定。控制方式：分类记录，不把外部失败误写为产品通过。
- 全路由 UI 回归发现阻断级缺陷。控制方式：只做最小修复，修复后重新跑相关验收。
- 截图/日志资产膨胀。控制方式：只保存代表性证据和 JSON 摘要，不提交 Playwright trace、临时 DB、完整日志或敏感输出。
- 发布口径过度乐观。控制方式：release candidate 必须引用实际 gate 结果，并明确 Not Claimed。
