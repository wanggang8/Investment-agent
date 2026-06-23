# Design: P56 UI 验收阻断与产品化设计修复

## Product Design Skill Usage

P56 采用 Product Design plugin 的以下方法作为设计输入，但契约真源仍遵守 `docs/GOVERNANCE.md` 与 OpenSpec change：

- `product-design:get-context`：确认产品 brief、用户、目标、视觉方向和交互需求。
- `product-design:audit`：沿用 P55 审查框架，按 UX、视觉、一致性、可访问性和验收风险分类。
- `product-design:research`：补充 dashboard、表单、移动 reflow、金融信任/可解释性的外部参考。

Superpowers `brainstorming` 用于方案探索；不在 `docs/superpowers/plans/` 新建规格真源，结论写回本 change。

## Product Brief

- 产品：本地优先的 Investment Agent，面向个人投资纪律、证据检索、风险预警、人工确认和复盘。
- 核心用户：需要每天快速判断“现在是否可行动、为什么、风险在哪里、下一步只做哪些人工动作”的个人投资者/操作者。
- 体验目标：像一个本地投资纪律工作台，而不是券商交易终端、AI 聊天 demo 或工程调试后台。
- 安全姿态：所有关键页面继续明确“只读/人工复核/不自动交易/不自动应用规则”边界，但文案应融入状态与动作区，不堆砌成噪声。
- 视觉方向：安静、密集、可信、可扫描；使用中性底色和语义状态色，避免默认浏览器控件和过度 marketing hero。

## Research Inputs

本阶段调研结论用于设计判断，不直接替代本项目契约。可追溯来源：

- Nielsen Norman Group dashboard 指南：dashboard 应帮助用户快速获得可行动的 at-a-glance 信息，利用预注意处理减少扫描成本。来源：`https://www.nngroup.com/articles/dashboards-preattentive/`
- Nielsen Norman Group form design 指南：表单应减少认知负担，明确标签、格式、必填/可选和错误反馈。来源：`https://www.nngroup.com/articles/web-form-design/`、`https://www.nngroup.com/articles/4-principles-reduce-cognitive-load/`
- W3C WCAG 2.1 Reflow：在 320 CSS px 宽度或 400% 缩放下，除真正二维内容外，不应要求同时水平和垂直滚动。来源：`https://www.w3.org/WAI/WCAG21/Understanding/reflow.html`
- Material Design data table 指南：表格适合行列扫描，但在窄屏应避免让用户为了读内容横向滚动完整页面。来源：`https://m2.material.io/components/data-tables`
- Explainable robo-advisor 相关研究：金融决策系统需要解释、透明度和信任线索，尤其在算法建议或市场波动场景下。来源：`https://ceur-ws.org/Vol-3222/paper6.pdf`、`https://www.aodr.org/xml//46040/46040.pdf`
- Robo-advisor UX 行业材料：金融产品应把用户流程、信任、透明度和 plain language 放在核心位置，避免让用户误解为自动代为交易。来源：`https://www.zymr.com/blog/how-to-build-robo-advisors-platform`

## Current Problems From P55

### Blocking Defect

- 真实 LLM consultation 写入 `decision_62160bd3494023dd` 后，打开决策详情时 React 崩溃。
- 原因是 `DecisionTrace.tsx` 直接对 `final_verdict.optional_actions` 调用 `.join()`，但真实 DTO 中该字段为 `null`。
- 该缺陷说明前端没有在 DTO 边界统一处理 nullable / unknown shape。

### Product Experience Problems

- 信息架构：17 个导航项平铺在侧栏，缺少“今日/决策/组合/证据/治理/系统”等任务分组。
- 移动端：固定侧栏在 390px 视口下占用过多宽度，`/positions` 和 `/data-quality` 出现横向溢出。
- 视觉系统：`web/src/index.css` 仍有 demo 风格根样式、过大的 `h1`、紫色 accent 和居中容器；`global.css` 卡片圆角/阴影偏散，控件样式不统一。
- 表单体验：持仓、知识导入、本地安装等页面使用默认输入控件或紧凑 inline label，不利于金融数据录入和纠错。
- 决策详情：内容更像 debug trace，最终裁决、证据、LLM 分析、规则边界、人工确认之间的阅读顺序不够产品化。

## Proposed Information Architecture

P56 不删除现有路由，但重组导航呈现：

| Group | Primary routes | Intent |
| --- | --- | --- |
| 今日 | `/`, `/workbench` | 快速判断今天能不能行动、优先看什么 |
| 决策 | `/consultation`, `/decisions/:id`, `/decision-loop` | 发起咨询、阅读裁决、追踪人工确认 |
| 组合 | `/positions`, `/risk-alerts` | 维护本地账户事实，处理风险 SOP |
| 证据 | `/data-quality`, `/evidence`, `/local-knowledge` | 判断数据质量、证据覆盖和本地知识导入状态 |
| 治理 | `/rules`, `/audit`, `/review`, `/notifications` | 审计、复盘、规则提案和通知 |
| 系统 | `/settings`, `/local-install`, `/daily-auto-run`, `/daily-discipline/reports` | 本地配置、安装诊断、自动运行和报告 |

桌面端应保留可扫描的侧栏或分组 rail；移动端应改为顶部栏 + 可展开菜单或分组抽屉，避免侧栏常驻挤压内容。

## UI System Direction

- Typography：页面标题控制在工作台尺度，避免 hero 级 `h1`；卡片内标题更紧凑。
- Color：中性背景 + 语义状态色，区分正常、观察、风险、降级、失败、成功；避免一套单色紫/蓝主题。
- Components：
  - AppShell / Navigation：支持分组、当前路由、高优先级入口、移动菜单。
  - Button：primary/secondary/ghost/danger，明确 disabled/loading/focus。
  - Field：label、hint、error、input/select/textarea 统一宽度与间距。
  - DataPanel / Metric / StatusBadge：用于 dashboard、data-quality 和 workbench。
  - ResponsiveTable / DetailList：桌面表格，移动端转为键值列表或分组卡片。
  - Empty/Error/Degraded State：统一可读，不把未知状态显示成成功。
- Layout：页面区块全宽、内部内容约束；避免卡片套卡片；固定格式元素使用稳定尺寸，防止 hover/文本导致布局跳动。

## Page-Level Plan

### App Shell

- 按任务分组重排导航。
- 移动端隐藏常驻侧栏，提供可达的菜单按钮和分组导航。
- 保持所有现有路由可达，避免破坏 P39/P55 路由矩阵。

### Dashboard and Workbench

- 第一屏聚焦“今日结论”：纪律状态、组合风险、数据质量、最近决策、下一步人工动作。
- 降低重复安全文案噪声，把禁止自动交易边界绑定到相关操作区。
- 保持 dashboard at-a-glance 可扫描，不做营销式 landing page。

### Decision Detail

- 首屏先呈现最终裁决、安全边界、必要动作和禁止动作。
- 将 LLM 分析、证据、规则命中、审计 trace 分成清晰 section；长 trace 使用展开/折叠或更紧凑的列表。
- 所有数组/对象字段在前端 adapter 边界做 null-safe normalization。

### Consultation

- 表单使用统一 field 组件，明确 symbol、场景、问题和提交状态。
- 提交后给出下一步路径：查看决策详情、记录人工确认、查看证据。

### Positions

- 桌面端保留高密度录入；移动端表单 label/input 纵向堆叠。
- 持仓表在移动端转为持仓卡片或定义列表，避免整体页面横向滚动。
- 成功、错误和校验信息靠近对应输入区。

### Data Quality

- 将 source health、Evidence/RAG、LLM quality、影响工作流分成清晰 summary + drill-down。
- 长 token、source id、request id、诊断文本必须换行或截断，不扩大页面宽度。
- 移动端改为纵向状态卡/键值列表。

## Functional Hardening

- 引入前端 DTO normalization utility，至少覆盖 decision detail 当前使用的 nullable list 字段。
- 为真实 LLM-like DTO 建立 fixture，包含：
  - `optional_actions: null`
  - `prohibited_actions: null`
  - 缺失数组字段
  - unknown status / degraded metadata
- 为核心页面增加异常边界或局部 fallback，避免单个 section shape 异常导致整页白屏。

## Validation Strategy

- Unit/component：
  - nullable decision DTO 不崩溃。
  - DecisionTrace 展示空列表、未知状态、降级状态。
  - Navigation 分组和移动菜单可访问。
  - Positions/DataQuality 移动布局关键元素存在，且禁用交易类文案不存在。
- Build/static：
  - TypeScript/Vite build 通过。
  - `git diff --check` 与敏感信息扫描通过。
- Browser/manual：
  - 启动真实后端与前端。
  - 桌面 1280x720 和移动 390x844 操作核心页面。
  - 重跑 P55 blocker 路径：真实 LLM consultation 后打开新决策详情不得崩溃。
  - 重测 `/positions` 和 `/data-quality`，`document.body.scrollWidth <= window.innerWidth` 或仅局部二维表格容器滚动，页面本身不得横向溢出。
  - 截图保存到 P56 release 资产目录。

## Review Requirements

方案完成后、实现完成后、归档前均需要子 agent 复审。复审必须覆盖：

- 是否遵守 `docs/GOVERNANCE.md` 与 OpenSpec 工作流。
- 是否引用 Product Design skill 和研究依据，而不是只做主观 UI 判断。
- 是否先修 P55-B1，再处理产品化 UI 优化。
- 是否保持安全边界：无交易、无外推、无自动确认、无自动规则应用、无收益承诺。
- 是否有可执行验收：单元测试、构建、真实浏览器操作、移动端截图、脱敏扫描。
