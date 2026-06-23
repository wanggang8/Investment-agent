# Tasks: P52 项目验收门禁矩阵

## 1. 范围与计划

- [x] 1.1 确认当前无活跃 change，P52 为下一阶段。
- [x] 1.2 创建 `p52-project-acceptance-gate-matrix` OpenSpec change。
- [x] 1.3 确认 P52 只定义验收门禁矩阵，不运行完整验收、不修改运行时代码。
- [x] 1.4 子 agent 复审 P52 计划，无 Critical / Important 后继续。

## 2. 验收入口梳理

- [x] 2.1 梳理 OpenSpec、Go、前端、E2E、smoke、安装诊断、发布升级、安全扫描入口。
- [x] 2.2 梳理真实公开源和真实 LLM opt-in 验收入口及失败分类。
- [x] 2.3 梳理 P51 证据包如何输入 P52 门禁矩阵。

## 3. 门禁矩阵文档

- [x] 3.1 新增 `docs/project-acceptance-gate-matrix.md`。
- [x] 3.2 定义 G0-G9 门禁表，包含命令、前置条件、通过标准、允许降级、产物、是否阻断发布。
- [x] 3.3 定义真实源/真实 LLM opt-in 策略和失败分类。
- [x] 3.4 定义验收记录格式和 P53 使用要求。
- [x] 3.5 明确 P52 不宣称验收已通过，只定义门禁。

## 4. 文档与进度同步

- [x] 4.1 更新 `openspec/PROGRESS.md`，标记 P52 活跃并指向 P53。
- [x] 4.2 更新 `openspec/project.md`、`docs/GOVERNANCE.md`、`docs/development-plan.md`、`AGENTS.md`。
- [x] 4.3 更新 `docs/README.md` 文档地图，加入项目验收门禁矩阵入口。

## 5. 验证与复审

- [x] 5.1 执行 `openspec validate p52-project-acceptance-gate-matrix --strict`。
- [x] 5.2 执行 `openspec validate --all --strict`。
- [x] 5.3 执行 `git diff --check`。
- [x] 5.4 子 agent 复审最终变更，无 Critical / Important 后归档。

## 6. 归档

- [x] 6.1 执行 OpenSpec archive，将 delta 合并入对应 `docs/` 真源，并按需同步 `openspec/specs/` 摘要。
- [x] 6.2 更新 archive 后进度状态并提交。
