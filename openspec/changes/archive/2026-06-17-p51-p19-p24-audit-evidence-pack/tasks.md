# Tasks: P51 P19-P24 审计证据包

## 1. 范围与计划

- [x] 1.1 确认当前无活跃 change，P51 为下一阶段。
- [x] 1.2 创建 `p51-p19-p24-audit-evidence-pack` OpenSpec change。
- [x] 1.3 确认 P51 只做文档审计证据包，不修改运行时代码。
- [x] 1.4 子 agent 复审 P51 计划，无 Critical / Important 后继续。

## 2. 证据收集

- [x] 2.1 核对 P14-P18 archive 存在性，说明无需重复补档。
- [x] 2.2 核对 P19-P24 缺逐阶段完整 archive 包，引用既有历史追溯说明。
- [x] 2.3 为 P19-P24 收集文档证据、代码证据、测试证据和不可声明事项。
- [x] 2.4 明确 P25-P29、P30/P39、P40/P44/P49 对 P19-P24 的后续补强关系。

## 3. 审计证据包

- [x] 3.1 新增 `docs/p19-p24-audit-evidence-pack.md`。
- [x] 3.2 写入总览、P14-P18 旁证核查和 P19-P24 分阶段矩阵。
- [x] 3.3 写入跨阶段证据、发布前使用建议和安全边界。
- [x] 3.4 确认证据包不把缺 archive 状态写成已 archive，不伪造历史完成记录。

## 4. 文档与进度同步

- [x] 4.1 更新 `openspec/PROGRESS.md`，标记 P51 活跃并指向 P52。
- [x] 4.2 更新 `openspec/project.md`、`docs/GOVERNANCE.md`、`docs/development-plan.md`、`AGENTS.md`。
- [x] 4.3 按需更新 `docs/README.md` 文档地图，加入 P51 审计证据包入口。

## 5. 验证与复审

- [x] 5.1 执行 `openspec validate p51-p19-p24-audit-evidence-pack --strict`。
- [x] 5.2 执行 `openspec validate --all --strict`。
- [x] 5.3 执行 `git diff --check`。
- [x] 5.4 子 agent 复审最终变更，无 Critical / Important 后归档。

## 6. 归档

- [x] 6.1 执行 OpenSpec archive，将 delta 合并入对应 `docs/` 真源，并按需同步 `openspec/specs/` 摘要。
- [x] 6.2 更新 archive 后进度状态并提交。
