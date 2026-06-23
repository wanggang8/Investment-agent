# P44: 本地安装诊断与打包体验

## Summary

在 P40 的本地运维基础上，补齐“新环境快速上手 + 诊断打包 + 备份恢复演练 + smoke 汇总”体验。该阶段不新增交易能力、不改动交易语义，仅提供可重复、可审计、可脱敏、以本地文件为边界的运维与打包入口。

## Why

当前用户已能完成真实任务调试与恢复演练，但仍缺少一条“从配置向导到诊断摘要导出”的统一体验。团队在新机或新环境时仍需手工查阅文档、逐条执行命令并拼装结果，容易遗漏、误读和污染本地目录。

## What Changes

- 增加本地安装/诊断打包脚本：集中执行依赖预检、备份恢复演练、e2e smoke（可选），并输出不含密钥的总结文件。
- 增加前端“本地安装与诊断”只读页面：提供配置向导文案、关键命令、诊断汇总导入入口与安全边界说明。
- 扩展导航与前端 smoke 覆盖新页面，维持安全扫描边界（不出现交易/自动应用规则等语义）。
- 更新 `docs/configuration.md`、`docs/ops-local-scheduler.md`、`docs/frontend-contract.md` 的本地运维体验描述。

## Scope

- 保持现有范围一致，仅覆盖本地安装与运维诊断体验，不扩展交易或外部推送能力。

## Out of Scope

- 接入券商 API。
- 自动交易、一键交易、代下单。
- 外部推送（邮件/短信/Webhook/System Push）。
- 自动交易规则确认、自动规则应用。
- 承诺收益或确定性收益判断。
- 连接付费/授权/登录源。

## Validation

- `go test ./...`
- `npm --prefix web test -- --run`
- `npm --prefix web run build`
- `bash scripts/local-install-diagnostics.sh --skip-e2e`（本地预检/恢复 smoke 总结）
- `bash scripts/e2e-smoke.sh`
- `openspec validate p44-local-install-diagnostics-packaging --strict`
- `openspec validate --all --strict`
