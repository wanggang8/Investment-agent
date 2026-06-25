# Proposal: P114 视觉产品化与对齐残留修复

## Why

P113 已修复机械布局问题，但用户复核后仍指出明显的产品级 UI 问题：表单与按钮错位、同层级卡片高度/节奏不一致、部分页面仍暴露工程化或诊断型内容，导致界面不像参考图中严谨、整齐、可交付的产品。

P113 的 `issueCount=0` 主要来自 DOM overflow、offscreen、touch target 等自动指标，不能充分覆盖视觉基线、等高关系、表单动作区、产品化信息架构和文案层级。P114 需要补上人工视觉审查和针对性修复。

## What

- 新建 P114 全产品 residual UI audit，覆盖桌面与 390px 移动视口。
- 明确审查以下问题类型：
  - 表单控件、按钮、辅助文案和动作区不在同一视觉基线。
  - 同层级卡片高度、标题区、内边距、按钮位置不一致。
  - 页面首层仍展示 raw JSON、命令、路径、diagnostic、枚举字段或过强工程语气。
  - 页面使用了通用后台式表单/卡片，而不是 P111 参考图的 report / ledger / action queue 风格。
  - 移动端表单按钮堆叠、宽度、间距、触控和文案换行不够精致。
- 对每个问题判断修复归属：
  - `frontend-mapping`: 前端通过组件、布局、文案映射和折叠即可解决。
  - `backend-summary-needed`: API 缺少产品化摘要字段，前端只能展示 raw/诊断内容，需后端补只读 summary。
  - `intentional-technical-secondary`: 作为二级详情保留，但不得出现在首层。
- 按审查结果修复前端；仅在证据显示 API 数据结构不足时，才提出或实现最小后端只读字段补充。
- 重新截图并生成 P114 mismatch ledger，确认无 P0/P1/P2 产品化视觉问题后归档。

## Out Of Scope

- 不重做全新视觉方向；仍以 P111 第二方案参考图为视觉真源。
- 不新增投资规则、交易能力、自动执行能力或收益承诺。
- 不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认或自动规则应用。
- 不处理 Docker、安装、发布包刷新、GitHub Release 或物理第二机器复验。

## Acceptance

- P114 change 通过 `openspec validate p114-visual-productization-alignment-fixes --strict`。
- 完成全产品桌面与 390px 移动截图审查，mismatch ledger 覆盖表单/按钮、卡片等高、产品化内容、移动端和参考图一致性。
- 所有 P0/P1/P2 UI residual findings 均修复并复验。
- 明确记录是否需要后端修改；如不需要，说明具体原因；如需要，列出 API/DTO 字段和最小后端变更。
- 通过前端测试、构建、OpenSpec 全量校验和 whitespace gate。
