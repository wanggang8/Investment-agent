# Design: P113 布局高保真精修

## Approach

P113 不重新发明视觉方向，而是把 P111/P112 的 reference cockpit 语言做成更稳定的布局系统。优先修共享 CSS，因为当前多个页面共用 `daily-hero`、metric grid、compact action 和 ledger/card primitives；同一处约束不稳会在数据质量、设置、本地安装、本地知识和治理页一起表现为横向溢出或卡片压缩。

## Layout Rules

- 桌面 report hero 允许状态卡换行成两行，不强迫五个以上卡片挤在同一行。
- 390px 移动端禁止关键内容使用横向 scroll rail；metric/status 区使用 compact two-column grid 或 single-column content list。
- 卡片内部所有标题、数值、标签必须 `min-width: 0` 且允许正常换行；长英文、路径、ID 使用安全断行。
- 移动端可点击 action/link/button 目标不低于约 36px，主要操作接近 40px。
- 工程化 raw 内容默认折叠或移到次级层级，首屏优先展示用户状态、解释和下一步动作。
- 所有页面 first viewport 必须有明确的 status/report 结构和下一步线索，不能只露出大面积解释卡片。

## Page Focus

- `/data-quality`: 修复桌面 hero 压缩和移动横向溢出，让数据可信度、风险、信号、规则一致性和下一步处理形成稳定 report。
- `/settings`、`/local-install`、`/local-knowledge`: 修复多指标移动错位，降低运维读数堆叠感。
- `/decisions/:id`: 强化决策标题、裁决状态、关键指标和证据链层级，减少普通后台详情页观感。
- `/rules`、`/audit`、`/notifications`、`/daily-discipline/reports`、`/daily-auto-run`: 修复移动触控和过早暴露 raw/工程内容的问题。

## Validation Strategy

Rendered QA 是 P113 的核心验收：每个页面必须有桌面和 390px 移动截图，截图后检查横向溢出、重叠、裁切、按钮高度、console health 和参考图一致性。发现 P0/P1/P2 布局问题时，不进入归档，必须修复并重新截图。
