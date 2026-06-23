# P52 项目验收门禁矩阵

## Why

P51 已补齐 P19-P24 当前事实审计证据包，但项目进入发布候选材料前仍缺一份统一的验收门禁矩阵。当前仓库已有单元测试、集成测试、前端测试、E2E smoke、真实公开源 smoke、真实 LLM smoke、本地安装诊断、备份恢复、发布升级检查和安全边界验证入口；P52 需要把这些入口整理成可执行、可记录、可判断是否阻断发布的矩阵。

## What Changes

- 新增 `docs/project-acceptance-gate-matrix.md`，定义发布前验收门禁：
  - 单元测试。
  - 集成测试。
  - 前端组件/页面测试与构建。
  - E2E / Playwright smoke。
  - fixture / current 数据质量和本地 smoke。
  - 真实公开源测试（显式 opt-in）。
  - 真实 LLM 测试（显式 opt-in）。
  - 本地安装诊断、备份恢复和发布升级检查。
  - 安全边界和脱敏检查。
- 每个门禁记录命令、前置条件、通过标准、允许降级、产物位置、是否阻断发布。
- 明确真实源/真实 LLM 测试失败必须分类，不得把网络、限流、凭证、模型不可用、解析失败混成单一失败。
- 明确 P52 只定义验收矩阵，不运行完整验收，不生成发布候选材料；P53 才整理发布材料。
- 更新治理、进度和开发计划，标记 P52 活跃，下一阶段指向 P53 `p53-release-candidate-materials`。

## Scope

- 仅修改文档、OpenSpec change、OpenSpec specs 摘要和进度状态。
- 不修改后端、前端、数据库 schema、迁移、API 行为、Eino 工作流、脚本行为或测试代码。
- 不执行真实外部源或真实 LLM 调用；仅定义显式 opt-in 验收策略。

## Out of Scope

- 新增测试自动化、修复测试、调整脚本行为。
- 实际运行完整发布验收。
- 生成 RC、发布说明或交付包。
- 接入券商 API、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复承诺、自动覆盖真实库、收益承诺、登录源、付费源、授权源、Level2 或高频源。

## Validation

- `openspec validate p52-project-acceptance-gate-matrix --strict`
- `openspec validate --all --strict`
- `git diff --check`
- 只读复审：确认 P52 只定义验收门禁矩阵，不宣称验收已通过，不扩大安全边界。
