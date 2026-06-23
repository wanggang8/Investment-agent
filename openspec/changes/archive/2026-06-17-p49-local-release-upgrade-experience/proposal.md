# P49: 本地发布与升级体验

## Summary

新增本地发布/升级检查入口，把版本检查、升级前备份提醒、迁移前预检和升级后 smoke 汇总收敛为可重复运行的本地报告。P49 不执行升级、不运行迁移、不覆盖真实 SQLite、不自动修复，也不新增交易或外部推送能力。

## Why

P40 已完成本地预检和恢复演练，P44 已完成安装诊断与打包入口。随着更多本地用户开始从旧版本升级，仍需要一个明确的“升级前/升级后 checklist”：确认当前版本、目标版本、备份是否已准备、迁移风险是否可见，以及升级后应跑哪些 smoke。现有命令分散，容易跳过备份或把恢复 smoke 误用于真实库。

## What Changes

- 新增本地 CLI 检查入口：
  - `go run ./cmd/agent --release-upgrade-check --target-version <version> --diagnostics ./tmp/release-upgrade.json`
  - 输出版本、备份提醒、迁移预检、升级后 smoke 命令和安全边界。
- 新增本地脚本包装：
  - `scripts/local-release-upgrade-check.sh`
  - 复用 CLI，默认写入 `tmp/local-release-upgrade/<timestamp>/release-upgrade.json` 和摘要。
- 扩展本地诊断打包脚本与文档：
  - 可选运行 release/upgrade check。
  - 在配置、运维和开发计划文档中记录升级流程、输出产物和禁止事项。
- 新增测试覆盖：
  - CLI 输出与 JSON 诊断不泄露密钥、完整私有路径、raw SQL、prompt 或 raw HTTP。
  - 检查入口只读，不创建 SQLite、VecLite、备份文件、不运行迁移、不写 audit。

## Scope

- 复用 P40 `--preflight`、`--backup`、`--recovery-smoke` 和 P44 本地诊断脚本的命令约定。
- 可新增 CLI DTO、报告构建函数、脚本、单测和文档。
- 报告可读取本地配置和迁移文件清单；不得打开运行时并触发 SQLite migration。

## Out of Scope

- 自动下载、自动升级、自动迁移、自动修复、自动覆盖真实库。
- 新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用。
- 新增登录源、付费源、授权源、Level2、高频源或收益承诺。
- 写入账户/持仓/决策事实、修改规则、创建通知、调用 LLM 或公网数据源。

## Validation

- `go test ./...`
- `bash scripts/local-release-upgrade-check.sh --target-version test-p49 --skip-preflight`
- `bash scripts/local-install-diagnostics.sh --skip-recovery --skip-e2e --include-release-upgrade --target-version test-p49`
- `openspec validate p49-local-release-upgrade-experience --strict`
- `openspec validate --all --strict`
- `git diff --check`
- P49 安全扫描（见 `tasks.md` 7.7）
