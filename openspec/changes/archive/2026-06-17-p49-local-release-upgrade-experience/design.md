# P49 Design

## Overview

P49 将本地升级流程做成“检查报告”，不是升级执行器。实现重点是把 P40/P44 已有的预检、备份、恢复 smoke 和本地诊断产物串成稳定 checklist，让用户在升级前知道必须先备份，在升级后知道应运行哪些 smoke，同时避免任何隐式写库、迁移、外部访问或自动修复。

## Approach

1. **CLI 报告入口**
   - 在 `cmd/agent` 增加 `--release-upgrade-check` 和 `--target-version`。
   - 可复用 `--diagnostics` 写 JSON 报告。
   - 报告包含 `generated_at`、`current_version`、`target_version`、`status`、`checks`、`backup_reminder`、`pre_upgrade_commands`、`post_upgrade_smoke_commands`、`safety_note`。
   - 该入口只加载配置与本地文件状态；不得调用 `openRuntime`，避免触发 SQLite migration。

2. **检查项**
   - `version_check`：展示当前本地版本占位和目标版本，目标为空时 warning。
   - `backup_reminder`：检查配置中 SQLite 路径是否存在；存在则提示先运行 `--backup`，不存在则提示首次安装或路径需确认。
   - `migration_precheck`：只读检查 `internal/infrastructure/persistence/sqlite/migrations` 中迁移文件是否存在且可读，记录迁移文件数量和最后文件名。
   - `smoke_plan`：列出升级后建议运行的 `--preflight`、`scripts/recovery-smoke.sh`、`scripts/e2e-smoke.sh`、`local-install-diagnostics.sh`。
   - `artifact_boundary`：说明报告和脚本输出位于 `tmp/`，不应提交私密本地产物。

3. **脚本包装与诊断集成**
   - 新增 `scripts/local-release-upgrade-check.sh`，统一创建输出目录并调用 CLI。
   - `scripts/local-install-diagnostics.sh` 增加可选 `--include-release-upgrade` 与 `--target-version`，默认不改变现有行为。

4. **脱敏与安全边界**
   - 报告不得输出 API key、私有绝对路径、原始 SQL、完整 prompt、raw HTTP 或供应商原始响应。
   - 路径字段使用占位符、basename 或相对命令模板，不输出用户 home 绝对路径。
   - 所有命令文本保持手动执行语义，不自动确认、自动修复或自动恢复。

## Risks

- 用户可能把检查报告理解成升级已完成。缓解：命令输出和 JSON 明确标注只读检查、未执行升级、未运行迁移。
- 迁移预检可能无法证明迁移一定成功。缓解：只声称迁移文件可读和升级后 smoke 计划，不承诺自动修复。
- 绝对路径泄漏。缓解：测试覆盖 JSON 与 stdout，不直接输出配置文件或 SQLite 绝对路径。

## Verification

- 单测覆盖 CLI 成功、目标版本缺失 warning、诊断 JSON、脱敏、只读边界。
- 脚本测试以 `--skip-preflight` 快速验证输出摘要。
- 严格 OpenSpec 校验、全量 Go 测试、git whitespace 校验和安全扫描。
