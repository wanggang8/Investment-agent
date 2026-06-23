# P44 Design

## Overview

P44 按“命令可见性优先”的思路补齐本地上手体验：通过一个单点脚本统一产生日志与摘要，再用前端页面将命令和配置模板收敛到一个清晰路径。核心目标是降低新环境误配与复现错误成本，不把诊断变成功能性写操作。

## Approach

1. **脚本化安装诊断与打包入口**
   - 新增 `scripts/local-install-diagnostics.sh`。
   - 步骤包括：
     - 使用 `INVESTMENT_AGENT_CONFIG` 运行 `cmd/agent --preflight --diagnostics`。
     - 默认导出 `preflight` 诊断文件到 `tmp/` 下。
     - 可选执行 `scripts/recovery-smoke.sh`。
     - 可选执行 `scripts/e2e-smoke.sh`。
   - 生成摘要 JSON（含各步骤状态、退出码、步骤命令和生成物路径），避免人工拼接。

2. **前端只读运维页面**
   - 新增 `web/src/pages/LocalInstallPage.tsx` 与测试。
   - 页面提供：
     - 安装步骤：配置路径、日志与临时目录约定、命令入口。
     - 配置向导：动态生成“启动用最小 config 草稿”。
     - 诊断摘要导入：支持上传脚本导出的摘要 JSON 并展示步骤状态（只读）。
     - 备份恢复与 smoke 命令：集中展示，强调不污染真实私有数据和不执行交易。

3. **路由/导航与验收覆盖**
   - 为页面添加路由 `/local-install`，在主导航中暴露。
   - 将关键 smoke 路径扩展为包含该页。

4. **文档与边界同步**
   - `docs/frontend-contract.md`：新增 P44 页面契约。
   - `docs/configuration.md`、`docs/ops-local-scheduler.md`：补充新脚本与打包输出说明。

## Risk & Guardrails

- 恢复 smoke 与 e2e 本身会消耗时间；通过命令行开关按需跳过，避免 CI 运行阻塞。
- 任何文件摘要必须明确 `tmp/`、命令与输出路径在 `.gitignore` 覆盖，防止持久化私密数据。
- 前端只读摘要展示，不能触发系统任务、不能上传或修改本地数据库。

