# Proposal: P112 参考图高保真细节修复

## Summary

P112 在 P111 已完成全产品参考图重构后，针对用户复核指出的“仍不够严谨、合理、细微处不像参考图”问题，执行一次逐页高保真细节修复。P112 不改变视觉方向，不新增运行时能力，而是把 P111 后仍存在的可见差距修到可交付：二级页面首屏节奏、report hero 高度、侧栏分组密度、状态卡语义、progress/checklist 图标、ledger 组件密度、移动端首屏效率和页面级参考图一致性。

参考图仍为 P111 锁定的第二方案：

`/Users/vick/.codex/generated_images/019ef8a7-f5c0-7442-95b9-e72bcdc89de1/ig_05724f56eb7089ab016a3b9109e1848191a87e68883d0c9826.png`

## Why

P111 的结果达到全局风格统一，但用户复核后指出与参考图仍有明显细节差距。重新抓取全 18 个桌面路由和关键移动路由后，确认 P111 存在以下可见问题：

- 部分二级页面仍像旧管理后台套新样式，而不是参考图的 report cockpit 信息架构。
- 多个页面首屏 hero 过高或主状态区下移，例如 `/data-quality`、`/settings`、`/local-install`、`/local-knowledge`、`/risk-alerts`、`/notifications`、`/daily-discipline/reports`、`/daily-auto-run`。
- 侧栏分组比参考图更碎、更拥挤。
- 首页/工作台的状态语义、metric tone、progress tracker、evidence checklist 和行动队列细节仍不够精细。
- 移动端首屏效率不足，hero 占用过大。

P112 的目标是把这些“会被设计审查指出的问题”修掉，并在完成后使用子 agent 进行逐页参考图对比审查。若子 agent 发现 Critical / Important / P0 / P1 / P2 视觉问题，必须继续修复后再次审查。

## What Changes

- 收敛全局 reference tokens：hero 高度、surface radius、边框/阴影、type scale、sidebar 分组、metric tone、list/ledger density。
- 重修 `AppLayout` 侧栏分组和密度，使其更接近参考图的核心工作区/系统与证据式秩序。
- 将二级页面统一为 reference report composition：紧凑 report hero + 右侧 next action/status block + 下方 ledger/checklist/action queue，而不是大面积旧式卡片堆叠。
- 压缩 `/data-quality`、`/positions`、`/rules`、`/audit`、`/settings`、`/local-install`、`/local-knowledge` 等页面的首屏占用。
- 修复 `/risk-alerts`、`/notifications`、`/daily-discipline/reports`、`/daily-auto-run` 等页面主状态区过晚出现的问题。
- 优化首页/工作台的状态卡、行动队列、progress tracker、evidence checklist 和 snapshot strip 细节。
- 优化 390px 移动端，让首屏更快露出下一步动作和证据/状态。
- 生成 fresh screenshot evidence、mismatch ledger、page pass matrix，并进行子 agent 逐页对比审查。

## In Scope

- 前端 React/Vite/TypeScript 页面和共享 UI 组件。
- CSS/tokens/布局/组件结构调整。
- UI 测试、截图验收、视觉审查文档和 acceptance record。
- P112 所需的本地临时后端、Vite、Browser 截图和移动/桌面 reflow 验证。

## Out Of Scope

- 不新增后端 API、SQLite schema、Eino workflow、LLM 能力、RAG/VecLite 能力、真实数据源或投资规则。
- 不改变最终裁决逻辑、规则提案逻辑、确认流程、审计语义、release/package/version 行为。
- 不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、真实库覆盖或收益承诺。
- 不新增登录源、付费源、授权源、Level2、高频源或外部商业数据依赖。
- 不把 P112 视觉验收扩大为新的投资效果、收益准确性、Docker、GitHub Release、发布包或物理第二机器验收声明。

## Validation

- `openspec validate p112-reference-fidelity-detail-pass --strict`
- `npm --prefix web test -- --run`
- `npm --prefix web run build`
- `go test ./...`
- `go vet ./...`
- 启动真实本地后端与 Vite 前端。
- 使用 Browser 采集全 18 个桌面路由截图与关键移动路由截图。
- 使用 `view_image` 对参考图和最新渲染截图做人工对比。
- 生成 P112 visual mismatch ledger；P0/P1/P2 mismatch 必须修复后重新截图。
- 使用子 agent 对所有页面进行逐页参考图对比审查；若存在 Critical / Important / P0 / P1 / P2，继续修复并再次审查。
- Forbidden affordance scan：不得新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、收益承诺。
- Sensitive/redaction scan：不得暴露完整 key、私有路径、SQL、完整 prompt、raw vendor payload、本地数据库路径或 raw stack。
- `openspec validate --all --strict`
- `git diff --check`
