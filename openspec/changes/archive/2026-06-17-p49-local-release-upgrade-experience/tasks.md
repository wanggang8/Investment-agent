## 1. OpenSpec 与范围

- [x] 1.1 确认 P48 已归档，P49 为当前活跃 change。
- [x] 1.2 确认 P49 聚焦本地版本检查、升级前备份提醒、迁移前预检和升级后 smoke 汇总。
- [x] 1.3 确认 P49 不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复承诺、自动覆盖真实库、收益承诺、登录/付费/授权/Level2/高频源。

## 2. CLI 发布升级检查

- [x] 2.1 在 `cmd/agent/main.go` 增加 `--release-upgrade-check` 与 `--target-version`。
- [x] 2.2 新增 release/upgrade report DTO 与构建函数，输出版本、备份提醒、迁移预检、smoke 计划和安全边界。
- [x] 2.3 复用 `--diagnostics` 写 JSON 报告；报告不得包含密钥、完整私有路径、raw SQL、prompt、raw HTTP 或供应商原始响应。
- [x] 2.4 确保该入口只读：不得调用 `openRuntime`，不得触发 SQLite migration，不创建备份/恢复文件，不写 audit。
- [x] 2.5 更新 `cmd/agent/main_test.go`，覆盖成功输出、target 缺失 warning、诊断 JSON、脱敏和只读边界。

## 3. 本地脚本

- [x] 3.1 新增 `scripts/local-release-upgrade-check.sh`，支持 `--config`、`--target-version`、`--output-dir`、`--skip-preflight`。
- [x] 3.2 扩展 `scripts/local-install-diagnostics.sh`，增加可选 `--include-release-upgrade` 与 `--target-version`，默认不运行 release/upgrade check。
- [x] 3.3 脚本摘要只记录步骤状态、退出码、相对命令模板和生成物位置；不得输出完整密钥或原始响应。

## 4. 文档与契约

- [x] 4.1 新增 OpenSpec delta 记录 P49 行为要求。
- [x] 4.2 更新 `docs/configuration.md` 和 `docs/ops-local-scheduler.md`，说明升级检查命令、备份提醒、迁移预检和升级后 smoke。
- [x] 4.3 更新 `docs/development-plan.md`、`openspec/project.md`、`openspec/PROGRESS.md`、`docs/GOVERNANCE.md` 和 `AGENTS.md` 当前阶段状态。

## 5. 执行前复审

- [x] 5.1 计划完成后执行只读子 agent 复审，确认无 Critical / Important。
- [x] 5.2 复审通过后再执行实现任务。

## 6. 执行后复审

- [x] 6.1 执行完成后再次只读子 agent 复审，确认无 Critical / Important。
- [x] 6.2 复审通过后执行 archive，并将 P49 归档。

## 7. 验收

- [x] 7.1 运行 `go test ./...`。
- [x] 7.2 运行 `bash scripts/local-release-upgrade-check.sh --target-version test-p49 --skip-preflight`。
- [x] 7.3 运行 `bash scripts/local-install-diagnostics.sh --skip-recovery --skip-e2e --include-release-upgrade --target-version test-p49`。
- [x] 7.4 运行 `openspec validate p49-local-release-upgrade-experience --strict`。
- [x] 7.5 运行 `openspec validate --all --strict`。
- [x] 7.6 运行 `git diff --check`。
- [x] 7.7 运行安全扫描：`rg -n 'sk-[A-Za-z0-9][A-Za-z0-9_-]{8,}|BEGIN (RSA|OPENSSH|PRIVATE) KEY|/Users/[^[:space:]，；。、]+|(?i:select[[:space:]]+\*[[:space:]]+from)|(?i:raw[[:space:]]+http)|(?i:prompt[[:space:]]*:)|完整[[:space:]]*prompt|HTTP/[0-9.]+[[:space:]]+[0-9]{3}|券商接口|自动交易|一键交易|代下单|外部推送|自动确认|自动应用规则|自动修复|自动覆盖|收益承诺|登录源|付费源|授权源|Level2|高频源' cmd/agent scripts docs openspec/changes/p49-local-release-upgrade-experience`，人工复核命中项，确认不存在未脱敏敏感内容或高风险操作入口；允许安全边界说明文本命中。
