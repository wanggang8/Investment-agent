# Design: P112 参考图高保真细节修复

## Design Target

P112 不创建新风格。它继续以 P111 锁定的第二方案参考图为唯一视觉真源，把 P111 后仍存在的“普通后台感”和“二级页不严谨”修成统一的 report cockpit 产品体验。

参考图的关键规则：

- 页面首屏应先呈现状态报告，而不是说明文字或大表单。
- report hero 应紧凑、有 icon well、有状态/警示/动作分区，并且下方内容露出。
- 二级页面应复用同一套 report + next action + metric/checklist/ledger 模块语言。
- 侧栏应低噪声、分组少而清楚，active item 有明确高亮但不过度装饰。
- 卡片边框、圆角、阴影、文字大小和图标 tone 都应克制、清楚、可扫描。
- 移动端首屏必须尽快露出下一步动作或关键状态，不能只显示一整屏 hero。

## Known P111 Gaps To Fix

| Area | Current Gap | P112 Target |
| --- | --- | --- |
| Secondary pages | 多个页面仍是旧管理后台结构套新卡片 | 统一为 reference report composition |
| Hero height | `/data-quality`、`/settings`、`/local-install` 等过高 | 首屏紧凑，下方内容可见 |
| Risk/ops pages | 主状态区下移，前置提示层太多 | 主状态报告进入首屏 |
| Sidebar | 分组碎、入口密度偏噪 | 更接近参考图的核心/证据/系统秩序 |
| Dashboard details | 状态语义、progress、metric、checklist 精细度不足 | 更接近参考图状态色、图标、列表节奏 |
| Mobile | hero 占屏过大 | 390px 首屏露出 action/status continuation |

## Implementation Approach

1. First tighten shared components and CSS tokens, because most page drift comes from reusable component geometry rather than isolated content.
2. Then fix core cockpit pages (`/` and `/workbench`) to reset the visual standard.
3. Then fix secondary pages by family: portfolio/data/risk, decision/evidence, governance, ops/settings.
4. After each family, capture screenshots and update the mismatch ledger.
5. At the end, dispatch sub agents for desktop and mobile visual review. Any material finding loops back into implementation.

## Acceptance Bar

P112 can only pass when:

- All 18 desktop routes have fresh screenshots.
- Key mobile routes have fresh screenshots and no horizontal overflow.
- P112 mismatch ledger has no open P0/P1/P2.
- Sub agent desktop review has no Critical/Important/P0/P1/P2.
- Sub agent mobile review has no Critical/Important/P0/P1/P2.
- Automated tests/build/OpenSpec/Go/whitespace gates pass.

## Safety Boundary

This is visual/frontend work. It must not add or imply broker connectivity, automated trading, one-click trading, order placement, external push, auto-confirmation, auto rule application, auto repair, real database overwrite, return promises, paid/auth/login data sources, Level2, high-frequency source, Docker release, GitHub Release, release package refresh, or physical second-machine validation.
