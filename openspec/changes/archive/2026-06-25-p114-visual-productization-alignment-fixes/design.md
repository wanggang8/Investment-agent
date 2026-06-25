# Design: P114 视觉产品化与对齐残留修复

## Design Source

P114 延续 P111 第二方案参考图作为视觉真源：深色侧栏、白色顶栏、report hero、优先级行动队列、状态卡片、资金快照、证据/规则 ledger 和克制的 8px 卡片半径。

## Audit Lenses

P114 不以 `issueCount=0` 作为唯一通过条件。审查必须包含人工视觉判断：

1. 表单与按钮：输入框、select、textarea、按钮、helper text、错误提示和动作区必须视觉对齐；桌面端动作按钮不应漂在不同高度，移动端必须整宽或合理分组。
2. 卡片同层级：同一 grid/row 内同层级卡片应具备一致的标题区、内边距、底部动作位置和视觉重量；内容长短不应造成“高低参差像拼贴”。
3. 产品化首层：首屏优先展示状态、解释、下一步人工动作和安全边界；raw JSON、命令、路径、diagnostic、内部枚举、长技术串默认折叠或进入二级详情。
4. 参考图一致性：页面应呈现报告式、纪律式、ledger 式产品界面，而不是普通后台表单堆叠。
5. 后端归属：只有当前端无法从现有 DTO 推导产品化摘要时，才进入 `backend-summary-needed`。

## Implementation Strategy

- 先运行真实渲染审查，形成 P114 finding ledger。
- 优先改共享 CSS / primitives：form action row、card grid equalization、secondary details、compact product summary。
- 再处理页面级结构：Positions、Consultation、Settings、Local Install、Local Knowledge、Rules、Evidence、Decision Detail、Data Quality、Daily Auto Run 等含表单/动作/诊断页面。
- 对产品化内容暴露先用前端映射和折叠解决；若 API 缺摘要字段，再补最小后端只读字段。

## Validation

- Browser 优先截图和 DOM 检查。
- 每个修复页面至少覆盖桌面和 390px 移动。
- 保存 P114 evidence folder 和 acceptance record。
- 运行 `npm --prefix web test -- --run`、`npm --prefix web run build`、`openspec validate --all --strict`、`git diff --check`。
