# Tasks: P56 UI 验收阻断与产品化设计修复

## 1. 方案与审查

- [x] 1.1 确认当前无活跃 change，P56 为 P55 UI acceptance blocker 后续阶段。
- [x] 1.2 阅读 `docs/GOVERNANCE.md`、`openspec/project.md`、`openspec/PROGRESS.md` 和 P55 验收/设计审查材料。
- [x] 1.3 使用 Product Design get-context / audit / research 相关 skill 输出 P56 产品设计依据。
- [x] 1.4 创建 `p56-ui-acceptance-blocker-fixes` OpenSpec change。
- [x] 1.5 子 agent 复审 P56 方案，无 Critical / Important 后进入实现。

## 2. P55 阻断修复

- [x] 2.1 定位决策详情 nullable DTO 崩溃路径，确认 `DecisionTrace` 和相关 adapter 的字段假设。
- [x] 2.2 增加真实 LLM-like decision DTO fixture，覆盖 `optional_actions: null`、`prohibited_actions: null`、缺失数组字段和 unknown/degraded metadata。
- [x] 2.3 实现 null-safe DTO normalization 或组件边界安全处理。
- [x] 2.4 增加前端回归测试，确认真实 LLM-like DTO 可渲染且不白屏。

## 3. 产品化 UI 基础层

- [x] 3.1 重构 app shell：任务分组导航、当前路由状态、移动端菜单/抽屉或 compact nav。
- [x] 3.2 调整全局 CSS token、字体层级、按钮、卡片、状态 badge、表单和表格基础样式。
- [x] 3.3 建立或复用 Field、Button、StatusBadge、ResponsiveTable/DetailList 等轻量组件，避免重复默认控件。
- [x] 3.4 确认所有核心页面仍可从导航到达，不破坏现有路由矩阵。

## 4. 页面级修复

- [x] 4.1 优化 Dashboard / Workbench 第一屏信息层级，突出今日结论、风险、数据质量和下一步人工动作。
- [x] 4.2 优化 Consultation 表单与提交后路径，保持只读/人工复核边界。
- [x] 4.3 优化 Decision Detail：最终裁决优先、证据/LLM/规则/审计 section 分层、长 trace 更可扫描。
- [x] 4.4 修复 `/positions` 移动端横向溢出，表单在窄屏纵向堆叠，持仓表可读。
- [x] 4.5 修复 `/data-quality` 移动端横向溢出，长 source/status token 可换行或局部滚动。
- [x] 4.6 检查 Settings、Local Install、Local Knowledge 的表单一致性，做低风险统一样式修补。

## 5. 测试与自动化验证

- [x] 5.1 运行前端单元/组件测试。
- [x] 5.2 运行前端构建。
- [x] 5.3 按受影响范围运行后端测试；若无后端修改，记录不适用原因并至少执行 OpenSpec 与前端验证。
- [x] 5.4 执行 `openspec validate p56-ui-acceptance-blocker-fixes --strict`。
- [x] 5.5 执行 `openspec validate --all --strict`。
- [x] 5.6 执行 `git diff --check`。
- [x] 5.7 执行敏感信息扫描，确认无 key、完整 prompt、私有 SQLite、raw vendor payload 或敏感路径泄露。

## 6. 真实 UI 验收

- [x] 6.1 启动真实后端 server 与 Vite 前端，使用临时 SQLite / 临时配置。
- [x] 6.2 使用浏览器操作验证真实 LLM consultation 后的决策详情可打开且不崩溃。
- [x] 6.3 重跑 P55 核心路由矩阵：Dashboard、Workbench、Decision Loop、Data Quality、Positions、Consultation、Decision Detail、Evidence、Rules、Audit、Notifications、Risk Alerts、Daily Auto Run、Daily Reports、Review、Local Install、Local Knowledge、Settings。
- [x] 6.4 在桌面 1280x720 和移动 390x844 采集截图。
- [x] 6.5 验证 `/positions`、`/data-quality` 和核心导航在移动端页面本身无横向溢出。
- [x] 6.6 扫描 UI 文案，确认没有自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用或收益承诺入口。

## 7. 报告、复审与归档

- [x] 7.1 新增 P56 UI 验收报告与截图资产。
- [x] 7.2 新增或更新 P56 设计复审记录，说明 Product Design skill、调研资料和实际 UI 改造的对应关系。
- [x] 7.3 子 agent 执行后复审，无 Critical / Important 后归档。
- [x] 7.4 执行 OpenSpec archive，将 delta 合并到 `docs/` 真源并同步规格摘要。
- [x] 7.5 archive 后确认无活跃 change，更新 `openspec/PROGRESS.md`、`docs/GOVERNANCE.md`、`openspec/project.md`、`docs/development-plan.md`、`AGENTS.md`。
- [x] 7.6 提交前子 agent 复审无 Critical / Important。
- [x] 7.7 提交 P56。
