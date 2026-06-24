# P111 Design: 高保真参考图视觉重构

## Reference Source Of Truth

P111 的视觉真源是用户选定的第二方案图：

`/Users/vick/.codex/generated_images/019ef8a7-f5c0-7442-95b9-e72bcdc89de1/ig_05724f56eb7089ab016a3b9109e1848191a87e68883d0c9826.png`

P111 不再把该图当作“方向感”，而是当作可检查的设计目标。实现必须逐项对照：布局、模块、密度、排版、颜色、边框、圆角、按钮、状态标签、图标、分割线、卡片比例和页面节奏。

## Fidelity Requirements

P111 页面完成标准分为四级：

| 等级 | 说明 | 归档允许 |
| --- | --- | --- |
| P0 mismatch | 影响页面骨架或产品识别，例如没有 top toolbar、大横幅、行动队列、指标矩阵，或仍是旧式卡片堆叠 | 禁止 |
| P1 mismatch | 影响首屏层级和可扫描性，例如模块顺序错误、密度明显不对、按钮/状态样式不一致 | 禁止 |
| P2 mismatch | 明显视觉不一致，例如色彩、边框、字号、间距、图标风格、圆角和 shadow 偏差大 | 禁止 |
| P3 mismatch | 小型 polish，例如个别字距、低风险文案长度、图标细节轻微差异 | 可记录后归档 |

每个页面的 mismatch ledger 必须包含：

- Reference evidence：参考图对应模块或视觉规则。
- Render evidence：当前页面截图和视口。
- Mismatch level：P0/P1/P2/P3/pass。
- Fix made：已修复的具体改动。
- Intentional deviation：若无法完全一致，说明原因和安全/契约依据。

## Page Coverage

### Tier 1: Core Cockpit

必须最高保真：

- `/` 今日纪律
- `/workbench` 决策工作台

目标：

- 使用参考图同款 shell：左侧导航、顶部工具栏、主内容 max-width 与 gutters。
- 首屏包含 report hero、current discipline state、prohibited actions。
- 主体包含 priority manual action queue、status overview metric cards、portfolio/fund snapshot、recent consultation/progress preview、evidence/rules checklist。
- `/` 偏每日纪律报告；`/workbench` 偏聚合决策工作台，但两者共享同一视觉骨架。

### Tier 2: Maintenance And Evidence

必须用参考图模块语言重构：

- `/positions`
- `/data-quality`
- `/risk-alerts`
- `/evidence`
- `/decision-loop`
- `/decisions/:decisionId`
- `/consultation`

目标：

- 维护类页面使用 status hero + action queue + metric grid + ledger/form surface。
- 证据/闭环类页面使用 evidence/rule checklist + progress tracker + ledger table/list。
- 决策详情/咨询页使用 report hero、process tracker、analysis panels、evidence checklist 和 manual confirmation panel。

### Tier 3: Governance And Ops

必须统一到参考图语言，不允许回到工程后台感：

- `/rules`
- `/review`
- `/audit`
- `/notifications`
- `/daily-discipline/reports`
- `/daily-auto-run`
- `/local-install`
- `/local-knowledge`
- `/settings`

目标：

- 顶部状态工具栏保持一致。
- 页面首屏必须有 summary hero 或 status strip。
- 列表、表格、通知、审计、规则提案统一为 ledger/progress/checklist variants。
- 运维页面只能展示本地检查、脱敏摘要和人工动作，不暗示自动修复。

## Reference Component System

### App Shell

- Sidebar width: 240px 左右，深蓝渐变或近似深色面。
- Sidebar brand：盾牌/纪律类图标 + Investment Agent + 本地投资纪律工作台。
- Nav group：图标 + label，分组标题较弱，active 背景为青绿色半透明块。
- Bottom status：版本/本地模式/离线优先。
- Topbar：页面标题 + 日期/本地时间；右侧 2-3 个 pill controls。

### Report Hero

- 大横向白色 panel。
- 左侧 circular icon well。
- 中间主标题 + warning/status line + explanatory copy。
- 右侧两列：当前纪律状态、禁止动作。
- 右侧列用竖分割线，禁止动作使用红色语义但不刺激交易。

### Priority Action Queue

- 标题行：下一步人工动作 + 按优先级 + info icon + count badge。
- Ordered rows：1/2/3 等数字方块，颜色按优先级。
- 每行包含 title、priority chip、detail、meta、right aligned outline action button。
- 低优先级行视觉弱化，但仍可读。

### Status Overview

- 2x2 或 4-card metric grid。
- 每个 metric card 有 icon well、label、large value、status chip、divider、bottom detail rows。
- 语义色来自 success/warning/danger/degraded/readonly，但必须有文字。

### Progress Tracker

- 水平步骤：输入假设、信息核查、LLM 分析材料、规则裁决、最终建议、人工确认。
- 每步有圆形 icon/status、连线、状态文字。
- 移动端改为纵向 timeline。

### Evidence And Rule Checklist

- 紧凑 checklist：icon、label、count/status、check indicator。
- 用于证据、规则、审计、source health、RAG/LLM 状态。

### Ledger Surface

- 轻边框白色 panel，密度接近参考图。
- 表格/列表要有清楚 header、分割线、行 hover/focus、局部横向滚动。

## Implementation Architecture

新增或重构共享前端组件，避免页面各写一套：

- `web/src/components/reference/ReferenceTopBar.tsx`
- `web/src/components/reference/ReferenceHero.tsx`
- `web/src/components/reference/PriorityActionQueue.tsx`
- `web/src/components/reference/StatusMetricGrid.tsx`
- `web/src/components/reference/SnapshotStrip.tsx`
- `web/src/components/reference/ProgressTracker.tsx`
- `web/src/components/reference/EvidenceChecklist.tsx`
- `web/src/components/reference/LedgerSurface.tsx`
- `web/src/components/reference/referenceTypes.ts`

CSS 优先集中在 `web/src/styles/global.css` 的 P111 section，或在现有全局样式中明确 `reference-*` class。不得引入新的 styling framework。

页面先接入共享组件，再做页面特定数据映射。若 API 缺字段，只能用现有 DTO 派生、显示“暂无/待检查”，不得伪造新能力。

## Testing Strategy

测试必须先于主要实现：

- `ReferenceComponents.test.tsx`：验证 hero、action queue、metric grid、progress tracker、checklist 的结构、状态文字和安全文案。
- `AppLayout.test.tsx`：验证 topbar、sidebar、原路由语义、安全边界和无交易入口。
- 核心页面 tests：Dashboard/Workbench/DataQuality/Risk/Evidence/DecisionLoop 需断言 P111 模块存在。
- Scan tests 或脚本：forbidden affordance、sensitive/redaction、reference screenshot QA JSON schema。

视觉 QA 不以测试替代，必须有真实浏览器截图和 mismatch ledger。

## Browser QA And Page Gates

P111 页面完成流程：

1. 实现页面。
2. 启动本地后端 + Vite 前端。
3. 捕获桌面截图，优先 1492 x 1068；如果 browser 限制则记录实际尺寸。
4. 用 `view_image` 打开参考图和渲染图。
5. 填写该页面 mismatch ledger。
6. 修复所有 P0/P1/P2。
7. 重新截图并更新 ledger。
8. 页面标记 pass 后进入下一页。

P111 不允许“所有页面最后统一再看”。每页都必须有自己的 pass 记录。

## Safety

所有页面继续只做本地查看、维护、记录、复核和导航。P111 不新增或暗示：

- 券商接口
- 自动交易
- 一键交易
- 代下单
- 外部推送
- 自动确认
- 自动规则应用
- 自动修复
- 真实库覆盖
- 收益承诺
- 登录源、付费源、授权源、Level2 或高频源
