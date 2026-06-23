# P60 UI 验收与设计审查记录

日期：2026-06-18

## 范围

- Change：`p60-portfolio-risk-data-quality-experience`
- 页面：`/positions`、`/risk-alerts`、`/data-quality`
- 环境：真实本地后端配置 + Vite 前端；配置文件含临时测试 key，未在本文档输出。
- 浏览器：Codex in-app Browser；截图为视口截图。

## 真实 UI 操作

### `/positions`

- 打开组合与持仓维护页。
- 执行一次本地账户校准：
  - 现金：`88`
  - 总资产：`116`
  - 标的：`510300 / 沪深300ETF`
  - 数量：`8`
  - 成本价：`3.2`
  - 现价：`3.5`
  - 买入理由：`P60 UI 验收本地初始化`
- 页面显示：`账户校准已保存为本地事实；不会连接交易接口。`
- 首屏状态更新为：`组合事实可用于纪律评估`。
- 额外发现：第一次使用不一致的 `总资产 != 现金 + 持仓市值` 数据时，后端拒绝写入；这是有效校验，不是前端崩溃。

### `/risk-alerts`

- 打开风险预警中心。
- 验证总览、`待看队列`、`处理中队列`、`需复盘队列`、`已记录队列` 可见。
- 当前真实库返回 `0 条` 风险预警，因此无 eligible SOP 生命周期按钮可操作；验收记录为只读队列和空态通过。
- 验证无横向溢出、无嵌套 cockpit card、无禁用操作入口。
- 执行后复审要求 SOP action 文案显式表达本地记录语义；已改为 `记录继续观察`、`记录升级复核`、`记录本地解除预警`。

### `/data-quality`

- 打开数据质量可观测页。
- 验证 `数据质量总览` 与四类信号：
  - `数据源健康信号`
  - `证据与 RAG信号`
  - `LLM 分析信号`
  - `影响范围信号`
- 点击 `查看数据源设置` 并成功导航到 `/settings`，再返回 `/data-quality`。
- 验证只读边界、脱敏边界和禁用入口扫描。
- 发现并修复：数据质量 view model 的负向安全文案仍包含 `自动修复` / `自动确认` / `自动应用规则`，会触发 forbidden copy scan；已改为 `不发起后台变更、规则确认、规则生效或资金动作`。
- 执行后复审要求首屏展示完整本地检查动作，并覆盖 `ops_status` 中 `empty`、`failed`、`quality_failed` 等状态；已补模型和页面测试。

## 390px 移动端

浏览器设置为 `390 x 844` 后验证：

| 页面 | body.scrollWidth | documentElement.scrollWidth | viewport | 结果 |
| --- | ---: | ---: | ---: | --- |
| `/positions` | 390 | 390 | 390 | 通过 |
| `/risk-alerts` | 390 | 390 | 390 | 通过 |
| `/data-quality` | 390 | 390 | 390 | 通过 |

## 产品与 UI 设计审查

- 三页都复用 P58/P59 的 operational cockpit 模式：状态先行、下一步动作明确、详情下沉。
- 首屏不再只是表单或列表，均先给出当前状态、维护/处置/质量信号和人工下一步。
- 桌面端保持工作台密度，没有营销式 hero；移动端改为单列，导航折叠后主要内容没有横向溢出。
- 卡片半径、边框、状态色和 spacing 与现有 `daily-hero` / `cockpit-card` token 一致。
- 未发现页面级文本重叠、按钮文字溢出、卡片套卡片或第一视口不可识别产品状态的问题。

## 截图资产

- `docs/release/ui-audit-assets/2026-06-18-p60/positions-desktop.jpg`
- `docs/release/ui-audit-assets/2026-06-18-p60/positions-mobile.jpg`
- `docs/release/ui-audit-assets/2026-06-18-p60/risk-alerts-desktop.jpg`
- `docs/release/ui-audit-assets/2026-06-18-p60/risk-alerts-mobile.jpg`
- `docs/release/ui-audit-assets/2026-06-18-p60/data-quality-desktop.jpg`
- `docs/release/ui-audit-assets/2026-06-18-p60/data-quality-mobile.jpg`

## 结论

P60 范围内真实 UI 验收通过。当前结论只覆盖组合维护、风险处置队列和数据质量可观测三页，不代表 P61 之后的最终 release-ready 状态。
