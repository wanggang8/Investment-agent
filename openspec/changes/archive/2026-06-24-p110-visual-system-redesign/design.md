# P110 Design: 视觉系统重设计

## Design Intent

P110 不推翻 P57-P63 的产品设计，而是在功能稳定后进行视觉成熟度升级。新 UI 应继续服务本地投资纪律，而不是刺激交易。视觉上从“后台台账”提升为“研究终端”：更强的信息层级、更好的阅读节奏、更少无意义边框、更稳定的状态语言。

## Visual Direction Options

P110 开始时先生成三版方向图，供用户选择：

1. **Research Terminal**：高密度研究终端。强调分栏、表格、状态轨道和证据链，适合重度使用。
2. **Calm Command Center**：冷静指挥台。强调首屏结论、人工动作和信任摘要，适合日常打开即判断。
3. **Ledger Pro**：高级审计账本。强调时间线、规则、审计和可追溯材料，适合复盘和治理感。

用户确认采用 **Calm Command Center** 作为主方向，并把 **Ledger Pro** 的审计账本气质用于 Evidence、Decision Loop、Audit 等证据/追踪页面。P110 代码实现以该组合方向为准：日常页面优先展示“今天能不能动、下一步人工做什么、哪些证据可信”，证据与闭环页面优先展示“来源、规则、时间线、readback 和审计轨迹”。

## Architecture

P110 保持现有架构：

- 路由仍由 `web/src/App.tsx` 和 `web/src/app/AppLayout.tsx` 组织。
- 页面仍位于 `web/src/pages/`，业务 view model 仍位于 `web/src/features/`。
- 通用 primitives 仍位于 `web/src/components/ui/`。
- 页面数据仍只能通过 `web/src/services/` 或 `web/src/shared/api/` 获取。
- 视觉 token、layout、状态和响应式规则优先集中在 `web/src/styles/global.css`，必要时小范围调整组件 class。

## Core Experience Changes

### Navigation

将当前深色侧栏从“后台目录”升级为更清晰的产品导航：分组仍保留，但 active、section label、品牌区和移动 topbar 要更有层级。导航不新增路由，不改变页面语义。

### Dashboard And Workbench

首屏继续遵循 P58：今日状态、下一步人工动作、信号摘要、详细驾驶舱。视觉上减少重复卡片边框，增强状态 hero、行动队列和信号组的主次关系，使用户在 5 秒内看清“能不能动”和“下一步人工做什么”。

### Decision And Evidence

Consultation、Decision Detail、Evidence、Decision Loop 保持解释链路：输入假设、信息核查、LLM 分析材料、规则裁决、最终建议、用户确认和审计 readback。视觉上把长文本、证据、规则和审计从平铺卡片改成更可扫描的分层阅读结构。

Evidence、Decision Loop 和后续审计类页面应使用 Ledger Pro 子方向：更强的台账 surface、清晰的行列信息、低刺激状态色和可追溯阅读节奏。它们仍然只展示证据与解释，不变成交易终端或自动执行入口。

### Portfolio, Risk, Data Quality

Positions、Risk Alerts、Data Quality 继续是维护/处置/质量面板。视觉上区分“事实录入”“处置队列”“质量阻断”，让用户能快速定位本地人工动作，而不是把所有内容看成同质卡片。

## Safety And Copy Boundaries

所有 UI copy 必须保持：

- 允许：查看、维护本地账户与持仓、刷新/检查数据、处理风险预警、查看解释、记录线下结果、人工确认。
- 禁止：自动交易、一键交易、代下单、券商连接、外部推送、自动确认、自动规则应用、自动修复、收益承诺。

视觉设计不得用红绿涨跌刺激、收益英雄数字、交易按钮样式或“立即操作”式引导制造交易冲动。

## Testing Strategy

P110 完成时必须覆盖：

- 前端组件/page tests 和 production build。
- 后端 Go tests，证明视觉改造未影响 Go 编译与共享契约。
- 真实浏览器桌面和移动截图。
- 390px、768px、1280px reflow。
- forbidden copy scan 和敏感信息 scan。
- OpenSpec strict validation 与 whitespace check。
