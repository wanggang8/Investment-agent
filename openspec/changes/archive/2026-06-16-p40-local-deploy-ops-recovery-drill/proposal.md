# P40: 本地部署、运维与恢复演练

## Summary

补齐本地交付体验：让新环境能按文档完成依赖检查、配置初始化、启动诊断、备份恢复演练和浏览器 smoke，并把数据源健康、日志、临时文件与诊断输出纳入可审计边界。

## Why

P33-P39 已把主要产品路径和浏览器级 E2E 串通，但本地运行仍依赖用户手动理解 Go、Node、Playwright、SQLite、VecLite、配置文件、备份恢复和失败诊断。P40 要把这些运维动作产品化为可重复、可诊断、不会污染仓库的本地交付流程，降低新环境启动和恢复时的误操作风险。

## Scope

- 增加本地初始化或自检入口，检查 Go、Node、Playwright browser、SQLite path、VecLite path、配置文件和必要目录。
- 增加启动前诊断与可修复提示，失败时输出安全、脱敏、可审计的信息。
- 增加备份恢复 E2E smoke，确认恢复后的历史事实可由 API/前端读取，且不会无确认覆盖现有 DB。
- 增加或强化数据源健康面板，展示最近成功时间、失败分类、新鲜度和影响范围。
- 明确本地日志、临时文件、诊断文件、Playwright 输出和 gitignore 治理。

## Out of Scope

- 不接券商 API。
- 不自动交易、不一键交易、不代下单。
- 不外部推送、不接邮件/短信/Webhook/系统 Push。
- 不自动应用规则、不绕过守门人审计或用户最终确认。
- 不引入付费源、登录源、授权源、Level2 或高频源。
- 不承诺收益或确定性涨跌预测。

## Validation

- `go test ./...`
- `npm --prefix web test -- --run`
- `npm --prefix web run build`
- 本地部署/恢复 smoke 命令
- `bash scripts/e2e-smoke.sh`
- `openspec validate p40-local-deploy-ops-recovery-drill --strict`
- `openspec validate --all --strict`
