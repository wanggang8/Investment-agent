# Tasks: P50 Post-P49 治理与验收路线图

## 1. 范围与事实确认

- [x] 1.1 确认 P50 只做治理规划，不修改运行时代码。
- [x] 1.2 核对 P19-P24 已有 `p19-p24-historical-archive-traceability` 归档，但仍缺逐阶段完整 archive 包。
- [x] 1.3 核对 P14-P18 已有 archive，避免把“P19-P14”误处理成重复补档。

## 2. 路线图固化

- [x] 2.1 固化 P51 `p51-p19-p24-audit-evidence-pack` 为下一阶段。
- [x] 2.2 固化 P52 `p52-project-acceptance-gate-matrix` 为 P51 后续阶段。
- [x] 2.3 固化 P53 `p53-release-candidate-materials` 必须依赖 P51/P52。

## 3. 验收策略规划

- [x] 3.1 定义 P52 需要覆盖的单元测试、集成测试、E2E 测试、真实测试、冒烟测试、发布前门禁。
- [x] 3.2 定义真实源/真实 LLM 测试必须显式 opt-in，并记录失败分类与降级边界。
- [x] 3.3 定义安全边界验收：不得出现券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复承诺、自动覆盖真实库、收益承诺。

## 4. 文档同步

- [x] 4.1 更新 `openspec/PROGRESS.md`，将下一阶段指向 P51。
- [x] 4.2 更新 `openspec/project.md` 与 `docs/GOVERNANCE.md` 的当前变更/后续说明。
- [x] 4.3 更新 `docs/development-plan.md` 与 `AGENTS.md` 的阶段状态和下一步建议。

## 5. 复审与验证

- [x] 5.1 子 agent 复审 P50 计划，无 Critical / Important 后继续。
- [x] 5.2 执行 `openspec validate p50-post-p49-governance-validation-roadmap --strict`。
- [x] 5.3 执行 `openspec validate --all --strict`。
- [x] 5.4 执行 `git diff --check`。
- [x] 5.5 子 agent 复审最终变更，无 Critical / Important 后归档。

## 6. 归档

- [x] 6.1 执行 OpenSpec archive，将 delta 合并入对应 `docs/` 真源，并按需同步 `openspec/specs/` 摘要。
- [x] 6.2 更新 archive 后进度状态并提交。
