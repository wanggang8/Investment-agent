## 1. OpenSpec 与范围

- [x] 1.1 确认 P44 聚焦本地安装诊断与打包体验：安装自检、配置向导、诊断导出、备份/恢复演练和 smoke 汇总。
- [x] 1.2 确认 P44 仅做本地运维体验增强，不接入券商、不交易、不中间介入规则生效。
- [x] 1.3 对齐 P43 合约约束：页面只读导航、无自动交易/自动规则入口、脱敏显示。

## 2. 安装诊断打包脚本

- [x] 2.1 新增 `scripts/local-install-diagnostics.sh`。
- [x] 2.2 脚本支持 `--config`、`--output-dir`、`--skip-recovery`、`--skip-e2e`。
- [x] 2.3 脚本必须运行 `go run ./cmd/agent --preflight --diagnostics <path>`，并导出摘要。
- [x] 2.4 脚本输出包含步骤名称、状态、退出码、命令、输出文件路径；默认不泄露 key。
- [x] 2.5 脚本成功时打印可复用的 smoke 摘要路径；失败项允许继续记录并返回失败摘要。

## 3. 前端本地安装与诊断页面

- [x] 3.1 新增 `web/src/pages/LocalInstallPage.tsx`，提供配置向导（仅读）和关键命令。
- [x] 3.2 页面支持上传脚本导出的 JSON 摘要并做只读渲染。
- [x] 3.3 页面禁止出现交易、自动确认、自动规则应用、外部推送等文案。
- [x] 3.4 添加页面测试（`web/src/pages/LocalInstallPage.test.tsx`）：
  - 命令与安全文案展示。
  - 配置向导输出与摘要导入。

## 4. 路由与入口

- [x] 4.1 在 `web/src/App.tsx` 注册 `/local-install` 路由。
- [x] 4.2 在 `web/src/app/AppLayout.tsx` 增加导航入口。
- [x] 4.3 更新 `web/e2e/local-smoke.spec.ts` 包含 `/local-install` 可达与关键禁止词扫描。

## 5. 文档与测试

- [x] 5.1 在 `docs/frontend-contract.md` 增加 P44 页约束。
- [x] 5.2 更新 `docs/configuration.md` 与 `docs/ops-local-scheduler.md` 的本地安装诊断节。
- [x] 5.3 运行并通过 `go test ./...`。
- [x] 5.4 运行并通过 `npm --prefix web test -- --run`。
- [x] 5.5 运行并通过 `npm --prefix web run build`。
- [x] 5.6 运行 `bash scripts/local-install-diagnostics.sh --skip-e2e`。
- [x] 5.7 运行 `bash scripts/e2e-smoke.sh`。
- [x] 5.8 运行 `openspec validate p44-local-install-diagnostics-packaging --strict`。
- [x] 5.9 运行 `openspec validate --all --strict`。

## 6. 归档前复审

- [x] 6.1 仅读子 agent 复审：确认无 Critical/Important。 
- [x] 6.2 复审通过后执行 archive。
