# Tasks: P54 发布交付与可重复性加固

## 1. 范围与计划

- [x] 1.1 确认当前无活跃 change，P54 为下一阶段。
- [x] 1.2 创建 `p54-release-handoff-and-repeatability-hardening` OpenSpec change。
- [x] 1.3 确认 P54 只做发布交付和可重复性文档，不修改运行时代码、不改变 P53 release_ready 结论。
- [x] 1.4 子 agent 复审 P54 计划，无 Critical / Important 后继续。

## 2. 发布交付文档

- [x] 2.1 新增 `docs/release/README.md`，列出 P53/P54 发布材料入口。
- [x] 2.2 新增 `docs/release/release-handoff-2026-06-17.md`，说明交付状态、验收结果、已知降级、复验入口和安全边界。
- [x] 2.3 新增 `docs/release/acceptance-repeatability.md`，固化复验目录、命令顺序、重试规则、G5 degraded 规则、G6 配置前提、G7 脱敏规则和 `release_blocked` 条件。

## 3. 文档与进度同步

- [x] 3.1 更新 `openspec/PROGRESS.md`，标记 P54 活跃。
- [x] 3.2 更新 `openspec/project.md`、`docs/GOVERNANCE.md`、`docs/development-plan.md`、`AGENTS.md`。
- [x] 3.3 更新 `docs/README.md`，加入 P54 release handoff 和 repeatability 入口。

## 4. 验证与复审

- [x] 4.1 执行 `openspec validate p54-release-handoff-and-repeatability-hardening --strict`。
- [x] 4.2 执行 `openspec validate --all --strict`。
- [x] 4.3 执行 `git diff --check`。
- [x] 4.4 执行提交材料脱敏扫描。
- [x] 4.5 子 agent 复审最终变更，无 Critical / Important 后归档。

## 5. 归档

- [x] 5.1 执行 OpenSpec archive，将 delta 合并入对应 `docs/` 真源，并按需同步 `openspec/specs/` 摘要。
- [x] 5.2 归档后确认无活跃 change。
- [x] 5.3 提交前子 agent 复审无 Critical / Important。
- [x] 5.4 提交 P54。
