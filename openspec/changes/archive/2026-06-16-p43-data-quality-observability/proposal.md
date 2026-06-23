# P43 Data Quality Observability Proposal

## Why

P42 已把日常入口聚合为用户决策工作台，但数据质量状态仍分散在设置页、复盘页、证据/检索 metadata、LLM 质量记录和运维预检文档中。用户在遇到 source health 降级、证据过期、RAG/VecLite 命中异常、LLM smoke 失败或质量门禁失败时，仍需要跨多个页面拼接“问题在哪里、影响哪些工作流、下一步检查什么”。

P43 用于建立数据质量可观测面，把现有 source health、证据新鲜度、RAG/VecLite 检索质量、LLM 质量门禁和本地诊断状态聚合成统一只读视图，并继续保持本地、安全、脱敏和不自动修复边界。

## What Changes

- 新增 P43 数据质量可观测页面或等价入口，展示数据源健康、证据新鲜度、RAG/VecLite 状态、LLM 质量状态和影响范围。
- 优先复用现有 API/service DTO；如确需后端扩展，只允许新增只读聚合 DTO，不新增持久化事实或数据库 migration。
- 页面提供到设置页、证据页、复盘摘要、审计、决策详情、风险预警和工作台的导航，不执行刷新、修复、外推、规则应用或交易。
- 前端测试和浏览器 smoke 覆盖成功、空库、source_unavailable、parse_error、stale、missing、unknown、LLM/RAG 降级和脱敏安全扫描。
- 更新 `docs/frontend-contract.md`、`docs/development-plan.md`、`docs/GOVERNANCE.md`、`AGENTS.md`、`openspec/PROGRESS.md` 和 `openspec/project.md` 的 P43 活跃状态。

## Out Of Scope

- 不接券商 API、不自动交易、不外部推送、不自动确认、不自动应用规则。
- 不承诺自动恢复数据源，不绕过 source verification、规则裁决、守门人审计或用户最终确认。
- 不展示完整 API key、prompt 全文、私有本地路径、SQL 错误、供应商原始错误或账户敏感明细。
- 不把 LLM 输出、RAG 命中或单一数据源状态提升为最终投资裁决。
- 不新增登录源、付费源、授权源、Level2、高频源或验证码绕过。

## Spec Deltas

- `data-quality-observability`: 新增 P43 数据质量可观测、脱敏、影响范围和只读导航要求。
- `frontend-ops-review-surface`: 增加数据质量可观测页面对 source/RAG/LLM/diagnostic 状态的前端要求。
- `frontend-experience-tests`: 增加 P43 组件与浏览器测试要求。

## Validation Plan

- `openspec validate p43-data-quality-observability --strict`
- `openspec validate --all --strict`
- `npm --prefix web test -- --run`
- `npm --prefix web run build`
- `bash scripts/e2e-smoke.sh`
- 如修改后端：`go test ./...`
- `git diff --check`
