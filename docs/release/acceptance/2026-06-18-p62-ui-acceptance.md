# P62 UI 验收与设计系统加固记录

日期：2026-06-18

范围：P62 `p62-design-system-accessibility-hardening`

## 结论

P62 范围内设计系统 primitives、键盘路径、可访问语义、390px/768px/1280px reflow 和真实本地 UI smoke 通过。

该结论只覆盖设计系统与可访问性加固，不代表 P63 全量真实 UI 回归或最终 release-ready 刷新。

## 实现摘要

- 新增 `web/src/components/ui/` primitives：Button、Field、StatusBadge、PageHeader、SummaryCard、DetailSection、ResponsiveTable、EmptyState、ErrorState。
- 统一状态 tone：success、warning、danger、degraded、unknown、readonly、blocked。
- 接入代表性页面和共享组件：Workbench、Positions、Data Quality、Risk Alerts、Rules、Audit、Notifications、Local Install、Local Knowledge、Settings。
- Audit 时间线折叠按钮补充 `aria-expanded` / `aria-controls` 和详情 region。
- Playwright smoke 新增 P62 三视口 reflow、移动菜单键盘路径、表单键盘输入、Audit 折叠区键盘展开和 Local Knowledge 校验按钮路径。

## 验收命令

| Gate | Result |
| --- | --- |
| `npm --prefix web test` | Pass: 47 files, 157 tests |
| `npm --prefix web run build` | Pass |
| `go test ./...` | Pass |
| `bash scripts/e2e-smoke.sh` | Pass: 3 Playwright tests |

## 浏览器证据

截图与浏览器结果位于 `docs/release/ui-audit-assets/2026-06-18-p62/`。

覆盖页面：

- `/workbench`
- `/positions`
- `/data-quality`
- `/risk-alerts`
- `/local-knowledge`

覆盖 viewport：

- `390 x 844`
- `768 x 900`
- `1280 x 900`

结果摘要：

- 截图数：15
- 页面级横向溢出失败：0
- 结果 JSON：`docs/release/ui-audit-assets/2026-06-18-p62/browser-results.json`

## 键盘路径

`bash scripts/e2e-smoke.sh` 已覆盖：

- 移动端导航按钮聚焦、Enter 打开菜单、链接聚焦并进入 Workbench。
- `/positions` 现金字段聚焦、键盘输入、Tab 到总资产字段。
- `/audit` P30 smoke 审计事件通过键盘展开引用，并验证 `aria-expanded=true`。
- `/local-knowledge` 校验预览按钮通过键盘触发，并展示批次与索引计划。

## 安全边界

P62 未新增后端 API、SQLite schema、Eino workflow、LLM 能力、数据源能力、券商接口、交易执行、外部推送、自动确认、自动规则应用、自动修复或发布状态刷新。

E2E safety scan 继续覆盖关键页面按钮、链接和部分 body 文案，未发现自动下单、一键交易、代下单、券商接口、外部推送、短信、邮件、Webhook、第三方推送、自动确认、自动修复、自动规则应用、收益承诺、完整密钥、密钥形态、SQL、私有路径或完整 prompt。
