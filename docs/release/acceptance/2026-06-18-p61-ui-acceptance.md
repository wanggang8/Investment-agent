# P61 UI 验收与设计审查记录

日期：2026-06-18

## 范围

- Change：`p61-governance-ops-productization`
- 页面：`/rules`、`/audit`、`/notifications`、`/daily-discipline/reports`、`/daily-auto-run`、`/local-install`、`/local-knowledge`、`/settings`
- 环境：真实本地后端临时 SQLite/VecLite + Vite 前端；使用 smoke fixture 与真实 HTTP API，不使用静态 mock。
- 浏览器：Playwright Chromium；截图为桌面与 390px 移动端 full-page 截图。

## 自动化验收

- `npm test`：42 个文件、144 个测试通过。
- `npm run build`：TypeScript 与 Vite build 通过。
- `go test ./...`：通过。
- `bash scripts/e2e-smoke.sh`：2 个 Playwright 测试通过。

P61 已扩展 Playwright smoke，覆盖治理/运维页面路由、关键状态、下一步动作、禁用文案扫描和 390px 横向溢出检查。

## 真实 UI 操作

### `/rules`

- 打开规则与纪律页。
- 验证 `规则治理状态`、当前规则、提案计数和 `查看审计记录` 下一步入口。
- 当前手动验收临时库为 `暂无规则提案` 空态；有提案、守门人结果和最终确认边界已由 `bash scripts/e2e-smoke.sh` 的 P39 journey 覆盖。
- 验证规则 raw detail 展示会替换高风险负向词，不把 `自动应用规则` 原文渲染到 UI。

### `/audit`

- 打开复盘与审计页。
- 验证 `审计检查状态`、审计分类摘要、审计时间线和 `审计下一步` 可见。
- 页面以可扫描时间线展示事件，而不是首屏 raw dump。

### `/notifications`

- 打开通知中心。
- 验证 `本地通知收件箱`、未读/严重/预警/总数指标、来源分类和本地处理入口。
- 当前手动验收库无单条未读可处理按钮；页面保持本地 inbox 边界。

### `/daily-discipline/reports`

- 打开每日纪律报告历史页。
- 验证 `每日纪律复盘状态`、最新报告/证据覆盖/自动运行/执行范围指标和人工复盘入口。
- 当前手动验收临时库为 `暂无每日纪律报告` 空态；报告详情导航已由 `bash scripts/e2e-smoke.sh` 覆盖。

### `/daily-auto-run`

- 打开每日自动运行页。
- 验证 `每日自动运行健康`、运行状态、执行范围和安全边界。
- 当前真实状态为默认关闭；页面没有把关闭或未知状态显示为成功。

### `/local-install`

- 打开本地安装与诊断页。
- 上传脱敏 `install-summary-fixture.json`。
- 验证页面显示 `失败步骤：1 个`、关键命令、启动草稿和本地复验入口。

### `/local-knowledge`

- 打开本地知识导入页。
- 执行 `校验预览`，验证 `知识预览`、脱敏预览、索引计划和安全说明。
- 使用安全样本执行 `写入本地事实`，验证写入结果摘要可见。

### `/settings`

- 打开设置页。
- 验证 `本地配置与诊断状态`、能力圈、系统状态、数据源健康和市场刷新入口。
- 执行 `刷新市场数据`，真实 API 返回成功，页面显示本地行情事实和审计记录更新完成。

## 390px 移动端

所有 P61 页面均在 `390 x 844` 视口下验证 `body.scrollWidth` 与 `documentElement.scrollWidth` 不超过 viewport。

| 页面 | body.scrollWidth | documentElement.scrollWidth | viewport | 结果 |
| --- | ---: | ---: | ---: | --- |
| `/rules` | 390 | 390 | 390 | 通过 |
| `/audit` | 390 | 390 | 390 | 通过 |
| `/notifications` | 390 | 390 | 390 | 通过 |
| `/daily-discipline/reports` | 390 | 390 | 390 | 通过 |
| `/daily-auto-run` | 390 | 390 | 390 | 通过 |
| `/local-install` | 390 | 390 | 390 | 通过 |
| `/local-knowledge` | 390 | 390 | 390 | 通过 |
| `/settings` | 390 | 390 | 390 | 通过 |

## 安全与脱敏扫描

- P61 手动浏览器验收 `failedCount=0`。
- `consoleErrors=[]`。
- P61 页面 body scan 未发现自动下单、一键交易、代下单、券商接口、外部推送、短信、邮件、Webhook、第三方推送、自动确认、自动修复、自动规则应用、收益承诺、完整密钥、`sk-`、SQL、私有路径或完整 prompt。
- 本阶段不新增后端 API、SQLite schema、Eino workflow、券商接口、交易执行或外部通知发送能力。

## 产品与 UI 设计审查

- 八个治理/运维页面统一为 operational cockpit 模式：首屏状态、指标、下一步人工动作，然后再进入列表、表格或表单详情。
- 关键任务路径从工程工具感转向工作台体验：规则治理、审计检查、通知处理、每日复盘、本地诊断、知识导入和设置维护都能回答“现在是什么状态、下一步去哪处理”。
- 桌面端保持高密度但可扫描；移动端单列 reflow，无页面级横向溢出。
- 高风险负向文案从产品界面中收敛为更安全的表达，例如规则生效需手动确认、站外通知、人工复验。

## 截图资产

- `docs/release/ui-audit-assets/2026-06-18-p61/browser-results.json`
- `docs/release/ui-audit-assets/2026-06-18-p61/_rules-desktop.png`
- `docs/release/ui-audit-assets/2026-06-18-p61/_rules-390.png`
- `docs/release/ui-audit-assets/2026-06-18-p61/_audit-desktop.png`
- `docs/release/ui-audit-assets/2026-06-18-p61/_audit-390.png`
- `docs/release/ui-audit-assets/2026-06-18-p61/_notifications-desktop.png`
- `docs/release/ui-audit-assets/2026-06-18-p61/_notifications-390.png`
- `docs/release/ui-audit-assets/2026-06-18-p61/_daily-discipline_reports-desktop.png`
- `docs/release/ui-audit-assets/2026-06-18-p61/_daily-discipline_reports-390.png`
- `docs/release/ui-audit-assets/2026-06-18-p61/_daily-auto-run-desktop.png`
- `docs/release/ui-audit-assets/2026-06-18-p61/_daily-auto-run-390.png`
- `docs/release/ui-audit-assets/2026-06-18-p61/_local-install-desktop.png`
- `docs/release/ui-audit-assets/2026-06-18-p61/_local-install-390.png`
- `docs/release/ui-audit-assets/2026-06-18-p61/_local-knowledge-desktop.png`
- `docs/release/ui-audit-assets/2026-06-18-p61/_local-knowledge-390.png`
- `docs/release/ui-audit-assets/2026-06-18-p61/_settings-desktop.png`
- `docs/release/ui-audit-assets/2026-06-18-p61/_settings-390.png`

## 结论

P61 范围内真实 UI 验收通过。当前结论只覆盖治理和运维页面产品化，不代表 P62/P63 或最终 release-ready 状态。
