# P110 Visual System Redesign Acceptance

日期：2026-06-24

## 1. 范围

P110 将前端视觉系统从可用后台台账质感升级为冷静的投资纪律研究终端。用户确认采用 **Calm Command Center** 作为主方向，并在 Evidence、Decision Loop 等证据/闭环页面吸收 **Ledger Pro** 的台账阅读感。

本阶段覆盖 AppLayout、Workbench、Data Quality、Risk Alerts、Evidence 和 Decision Loop 的视觉层级、共享 tokens、surface、状态面板、人工动作队列、证据/闭环台账 surface 和 390px/768px/1280px reflow。P110 不新增后端 API、SQLite schema、Eino workflow、LLM 能力、投资规则、数据源、Docker/发布包/物理第二机器验收或任何交易/外推/自动确认/自动规则能力。

## 2. 视觉证据

方向图：

- Research Terminal：`/Users/vick/.codex/generated_images/019ef8a7-f5c0-7442-95b9-e72bcdc89de1/ig_05724f56eb7089ab016a3b90a0f16c819199cbba4340c4ec24.png`
- Calm Command Center：`/Users/vick/.codex/generated_images/019ef8a7-f5c0-7442-95b9-e72bcdc89de1/ig_05724f56eb7089ab016a3b9109e1848191a87e68883d0c9826.png`
- Ledger Pro：`/Users/vick/.codex/generated_images/019ef8a7-f5c0-7442-95b9-e72bcdc89de1/ig_05724f56eb7089ab016a3b9189b78c81919758896fe18aafa3.png`

Browser QA assets：

- 桌面：`docs/release/ui-audit-assets/2026-06-24-p110-visual-system-redesign/01-workbench-desktop.png` 至 `05-decision-loop-desktop.png`
- 390px：`06-workbench-mobile.png` 至 `10-decision-loop-mobile.png`
- 768px：`11-workbench-tablet.png` 至 `14-decision-loop-tablet.png`
- 1280px：`15-workbench-1280.png` 至 `18-decision-loop-1280.png`
- 自动检查：`desktop-browser-qa.json`、`responsive-browser-qa.json`

桌面、390px、768px、1280px 自动检查均显示核心路由无页面级横向溢出。`responsive-browser-qa.json` 中 console 记录的 404/409 为既有本地 API 状态响应，在 P100 已按 classified API response 处理，不是 P110 视觉回归。

## 3. 实现摘要

- `web/src/styles/global.css` 新增 command center tokens、shared shell、manual action panel、trust signal strip、ledger surface 和移动端 reflow。
- `web/src/app/AppLayout.tsx` 增加本地模式/只读导航安全状态，并保留原导航名称与路由语义。
- Workbench 人工动作队列和信号摘要接入共享视觉 class。
- Evidence 与 Decision Loop 接入 Ledger Pro 台账 surface。
- 新增 `web/src/styles/visualSystem.test.ts` 与 `web/src/app/AppLayout.test.tsx`，覆盖 tokens、共享 class、响应式 guardrail、导航语义和安全边界。

## 4. 验收命令

已通过：

- `npm --prefix web test -- --run`
- `npm --prefix web run build`
- `go test ./...`
- `go vet ./...`
- `openspec validate p110-visual-system-redesign --strict`
- `openspec validate --all --strict`
- forbidden copy scan
- sensitive/redaction scan
- `git diff --check`

补充说明：

- `go test ./...` 与 `go vet ./...` 仅出现 sqlite-vec 依赖在 macOS SDK 上的 deprecated warning，命令退出码为 0。
- Production UI forbidden affordance scan 只命中 `PortfolioPage` 的“买入理由”字段，且同一行说明“不生成收益承诺”，不属于交易入口。
- Production source sensitive scan 使用 bounded `sk-...`、private path、private key 和 raw diagnostic patterns，结果无命中。
- OpenSpec archive 已将 change 归档到 `openspec/changes/archive/2026-06-24-p110-visual-system-redesign/`，并将 frontend experience delta 合并到 `openspec/specs/frontend-experience-tests/spec.md`。
