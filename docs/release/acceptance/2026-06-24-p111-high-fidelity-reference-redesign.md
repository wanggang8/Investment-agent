# P111 High-Fidelity Reference Redesign Acceptance

日期：2026-06-24
Change：`p111-high-fidelity-reference-redesign`
视觉真源：`/Users/vick/.codex/generated_images/019ef8a7-f5c0-7442-95b9-e72bcdc89de1/ig_05724f56eb7089ab016a3b9109e1848191a87e68883d0c9826.png`

## 结论

P111 已按用户选定的第二套参考图完成全产品高保真视觉重构。验收不再只以“组件存在”作为通过标准，而是同时使用：

- 参考图与最新渲染截图同屏人工对照。
- Playwright Chromium 采集 18 个路由的桌面 1492×1068 与移动 390×844 截图。
- `visual-mismatch-ledger.md` 记录 P0/P1/P2/P3/pass 结果。
- 前端测试、构建、OpenSpec、Go 和安全边界扫描作为回归门禁。

当前已知 P0/P1/P2 高保真视觉 mismatch：0。

## 已发现并修正的主要差距

| 差距 | 等级 | 修正 |
| --- | --- | --- |
| 首页/工作台只做了 reference 外观，但快照条被拉成整行，导致最近咨询和证据快照被挤出首屏 | P1 | 改为参考图右栏结构：左侧人工动作队列，右侧状态总览 + 持仓与资金快照 |
| 侧栏顶部多出运行模式卡片，参考图只有导航组和底部本地状态 | P2 | 将运行边界保留为辅助语义，视觉上改为底部 `v0.1.0 / 本地模式 · 离线优先` 状态 |
| Hero 图标使用 CSS 绘制，不符合参考图图标系统 | P2 | 改为 `lucide-react` 图标，保留可访问标签 |
| 首屏密度偏松，metric cards 和 action rows 比参考图更高 | P2 | 压缩 reference hero、action row、metric card、snapshot strip 的 padding 和字号节奏 |
| `/data-quality` 首屏先出现说明和标的切换控件，抢占参考图式状态报告入口 | P1 | 调整为先展示数据质量状态报告，再展示标的筛选 ledger surface |
| 通用页面 hero 有旧式左彩条和灰色嵌套右卡 | P2 | 统一为白色报告面板 + 右侧细分隔行动区 |
| 桌面顶部误露移动导航按钮、移动端 topbar 文本对比不足 | P2 | 桌面隐藏 nav toggle，移动端 topbar 使用深色背景与白色标题 |

## 页面覆盖

覆盖路由：

`/`、`/workbench`、`/positions`、`/data-quality`、`/risk-alerts`、`/consultation`、`/decisions/decision_smoke_p30`、`/decision-loop`、`/evidence`、`/rules`、`/review`、`/audit`、`/notifications`、`/daily-discipline/reports`、`/daily-auto-run`、`/local-install`、`/local-knowledge`、`/settings`。

截图目录：

`docs/release/ui-audit-assets/2026-06-24-p111-high-fidelity-reference-redesign/`

关键证据：

- `desktop-dashboard.png`
- `desktop-workbench.png`
- `desktop-data-quality.png`
- `desktop-evidence.png`
- `desktop-settings.png`
- `mobile-dashboard.png`
- `desktop-contact-sheet.png`
- `mobile-contact-sheet.png`
- `visual-mismatch-ledger.md`
- `visual-mismatch-ledger.json`

说明：优先尝试 Codex in-app Browser 截图，但 CDP `Page.captureScreenshot` 超时；已记录 fallback，并使用 Playwright Chromium 进行确定性截图采集。

## 安全边界

P111 只修改前端视觉、布局、共享组件、样式、测试和文档。不新增后端 API、SQLite schema、Eino workflow、LLM 能力、数据源、投资规则、Docker/安装/发布包能力。

P111 不声称也不提供券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、真实库覆盖、收益承诺、登录源、付费源、授权源、Level2 或高频源。

## 验证命令

最终验证命令：

```bash
npm --prefix web test -- src/components/reference/ReferenceComponents.test.tsx src/app/AppLayout.test.tsx src/pages/WorkbenchPage.test.tsx src/features/dashboard/DashboardFeature.test.tsx
npm --prefix web test -- src/pages/DataQualityPage.test.tsx src/pages/DecisionDetailPage.test.tsx src/pages/EvidencePage.test.tsx src/pages/DecisionLoopPage.test.tsx src/pages/ReviewSummaryPage.test.tsx
P111_BASE_URL=http://127.0.0.1:14111 P111_OUTPUT_DIR=/Users/vick/Desktop/project/Investment-agent/docs/release/ui-audit-assets/2026-06-24-p111-high-fidelity-reference-redesign node web/scripts/p111_visual_reference_audit.mjs
```

最终门禁结果：

```bash
npm --prefix web test
# 52 files passed, 189 tests passed

npm --prefix web run build
# passed

go test ./...
# passed; macOS sqlite-vec cgo deprecation warnings only

go vet ./...
# passed; macOS sqlite-vec cgo deprecation warnings only

openspec validate p111-high-fidelity-reference-redesign --strict
# Change is valid

openspec validate --all --strict
# 35 passed, 0 failed

git diff --check
# passed

rg -n "sk-(proj-|live-|test-)?[A-Za-z0-9_-]{32,}" . --glob '!web/node_modules/**' --glob '!node_modules/**' --glob '!tmp/**' --glob '!docs/release/ui-audit-assets/**'
# no matches
```

Forbidden affordance scan 命中的内容均为文档边界、免责声明或测试断言；未发现新增自动交易、券商、代下单、外部推送、自动确认、自动规则应用、自动修复或收益承诺入口。
