# P54 发布交付与可重复性加固

## Why

P53 已执行 P52 G0-G9 并生成 `release_ready` 材料，但验收记录中仍有三类需要面向交付明确说明的事项：

- G3/G4/G8 初次本地进程被 kill 后原命令重试通过。
- G5 current data-source quality regression 为 `degraded`，虽然不阻断发布，但限制当前本地数据快照质量声明。
- G6 真实公开源初次使用的临时配置缺少真实模式 market prerequisite，修正后通过。

P54 需要把这些事项整理成可交付、可复验、可审计的发布交付说明，避免用户把 `release_ready` 误解为无条件、未来可用性或收益/交易能力承诺。

## What Changes

- 新增 `docs/release/README.md`，作为发布材料索引。
- 新增 `docs/release/release-handoff-2026-06-17.md`，说明本次交付内容、验收状态、已知降级、复验入口、不可声明事项和下一阶段建议。
- 新增 `docs/release/acceptance-repeatability.md`，固化 P53 后的复验规则：
  - 重跑前提。
  - G0-G9 命令顺序。
  - 允许一次原命令重试的条件。
  - G5 current degraded 处理规则。
  - G6 临时配置前置条件。
  - G7 LLM key 脱敏规则。
  - 何时必须改为 `release_blocked`。
- 更新治理、进度、开发计划和文档地图。
- 增加 OpenSpec 行为摘要，约束发布交付材料必须引用 P53 实际验收记录，且不得扩大 P53 的 release_ready 声明。

## Scope

- 仅修改文档、OpenSpec change、OpenSpec specs 摘要和进度状态。
- 不重新执行 P53 全量验收，不改变 P53 结论。
- 不修改后端、前端、脚本、配置样例、SQLite schema、HTTP API、Eino workflow 或测试代码。

## Out of Scope

- 修复或优化导致 P53 初次重试的本机资源问题。
- 自动化新的验收 runner。
- 改写 P53 验收结果。
- 引入券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复承诺、自动覆盖真实库、收益承诺、登录源、付费源、授权源、Level2 或高频源。
- 提交完整 key、临时 SQLite、raw HTTP 响应、完整 prompt、原始 SQL 或私有路径。

## Validation

- `openspec validate p54-release-handoff-and-repeatability-hardening --strict`
- `openspec validate --all --strict`
- `git diff --check`
- 脱敏扫描 release/P54 提交材料。
- 子 agent 计划复审、执行后复审和提交前复审均无 Critical / Important。
