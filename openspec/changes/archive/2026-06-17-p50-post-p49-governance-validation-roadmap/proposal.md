# P50 Post-P49 治理与验收路线图

## Why

P49 已完成本地发布与升级体验，但进入后续发布材料前仍有两个治理风险需要先收敛：

- P19-P24 已标记为能力交付，但仓库没有逐阶段完整 `openspec/changes/archive/` 包；当前仅有 `p19-p24-historical-archive-traceability` 说明入口，不能替代 proposal / delta / tasks / 验收证据级追溯。
- 项目已经积累单元测试、集成测试、E2E、真实源 smoke、本地安装诊断、发布升级检查等多类验证入口，但缺少一份发布前验收矩阵来说明哪些门禁必须通过、哪些真实测试需要显式配置、哪些降级结果可以接受。

本阶段先做治理规划，明确后续顺序：先补 P19-P24 审计证据包，再补全项目验收门禁矩阵，最后再进入发布材料或 RC 阶段。

## What Changes

- 固化 P50 之后的候选队列：
  - P51 `p51-p19-p24-audit-evidence-pack`：补 P19-P24 历史交付证据矩阵，不伪造历史 archive。
  - P52 `p52-project-acceptance-gate-matrix`：建立全项目验收矩阵，覆盖单元、集成、E2E、真实源、真实 LLM、冒烟、安装诊断、发布升级检查。
  - P53 `p53-release-candidate-materials`：仅在 P51/P52 通过后整理本地发布材料、RC 检查清单与交付说明。
- 明确 P51 的输出应是审计证据包，而不是补建伪历史归档。
- 明确 P52 的输出应是可执行验收门禁与结果记录策略，而不是泛泛测试说明。
- 更新进度、治理和开发计划文档，使下一阶段自动指向 P51。

## Scope

- 仅修改 OpenSpec、治理文档、进度文档、开发计划。
- 不修改后端、前端、数据库 schema、工作流、配置默认值、脚本行为或测试代码。
- 不新增真实外部源、登录源、付费源、授权源、Level2、高频源。

## Out of Scope

- 直接发布、打包或生成 RC。
- 重写 P19-P24 历史 archive。
- 实现新的验收自动化命令。
- 接入券商 API、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、自动覆盖真实库、收益承诺。

## Validation

- `openspec validate p50-post-p49-governance-validation-roadmap --strict`
- `openspec validate --all --strict`
- `git diff --check`
- 只读复审：确认 P50 没有把 P19-P24 缺失 archive 误写成已归档，也没有扩大交易、外推或自动化边界。
