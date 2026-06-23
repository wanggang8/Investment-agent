# P53 验收执行与发布候选材料

## Why

P52 已定义项目验收门禁矩阵，但尚未执行 G0-G9，也未生成实际验收记录。直接整理发布材料会把“验收标准已定义”误当成“验收已通过”。P53 需要按 P52 矩阵执行门禁、记录真实结果，再基于结果生成发布候选材料。

本阶段目标是让项目进入可审计的发布前验收状态：所有通过项、降级项、跳过项、阻断项都必须有命令、产物、原因和发布影响。

## What Changes

- 执行 `docs/project-acceptance-gate-matrix.md` 的 G0-G9 门禁。
- 使用本地临时目录保存命令输出和 smoke 产物：`tmp/acceptance/p53-*`。
- 新增验收记录：`docs/release/acceptance/2026-06-17-p53-acceptance-run.md`。
- 新增发布候选材料：`docs/release/release-candidate-2026-06-17.md`。
- 若任何阻断门禁失败，发布候选材料必须写明 `release_blocked`，不得声明 release ready。
- 更新治理、进度、开发计划和文档地图。
- 增加 OpenSpec 行为摘要，约束 P53 必须引用实际门禁结果或 waiver，而不是引用 P52 作为通过证据。

## Scope

- 执行已有测试、smoke、真实公开源 opt-in、真实 LLM opt-in、安装诊断、发布升级检查和安全扫描。
- 仅新增或修改验收/发布文档、OpenSpec change、OpenSpec specs 摘要和进度状态。
- 可使用 `tmp/acceptance/` 下的临时配置、临时 SQLite 和日志产物；这些产物不提交。
- 真实 LLM 使用本地配置文件 `configs/config.local.yaml` 中的测试 key/base_url/model，不通过环境变量注入。

## Out of Scope

- 修复 P53 验收发现的运行时代码问题；若有阻断问题，记录为后续修复阶段。
- 新增验收 runner 脚本或改变既有测试命令语义。
- 新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复承诺、自动覆盖真实库、收益承诺、登录源、付费源、授权源、Level2 或高频源。
- 把真实公开源或真实 LLM 的一次通过结果写成未来可用性、收益或交易能力承诺。
- 提交包含完整 API key、私有路径、raw HTTP 响应、完整 prompt 或原始 SQL 的验收产物。

## Validation

- `openspec validate p53-acceptance-execution-and-release-candidate-materials --strict`
- `openspec validate --all --strict`
- `git diff --check`
- 按 P52 G0-G9 执行门禁并在验收记录中写明结果。
- 子 agent 计划复审、执行后复审和提交前复审均无 Critical / Important。
