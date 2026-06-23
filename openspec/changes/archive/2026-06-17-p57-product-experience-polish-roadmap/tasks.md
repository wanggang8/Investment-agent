# Tasks: P57 产品体验打磨总规划

## 1. 方案准备

- [x] 1.1 阅读 `docs/GOVERNANCE.md` 与 `openspec/project.md`。
- [x] 1.2 确认当前无活跃 change，工作区干净。
- [x] 1.3 使用 Product Design get-context / P55-P56 audit / research 资料定义产品体验打磨 brief。
- [x] 1.4 创建 `p57-product-experience-polish-roadmap` OpenSpec change。

## 2. 路线图设计

- [x] 2.1 明确 P57 从发布状态刷新改为产品体验打磨总规划，发布刷新后移。
- [x] 2.2 定义产品北极星、核心用户问题和体验原则。
- [x] 2.3 梳理一级核心页面、二级解释页面、三级治理/运维页面。
- [x] 2.4 拆分 P58-P63 后续阶段、范围、出入边界和验收要求，并明确 Notifications 与 Daily Auto Run 纳入 P61 治理/运维产品化。
- [x] 2.5 明确每个阶段必须真实启动项目、浏览器操作 UI、使用 Product Design skill 审查并执行子 agent 复审。

## 3. OpenSpec 与文档 delta

- [x] 3.1 编写 `proposal.md`。
- [x] 3.2 编写 `design.md`。
- [x] 3.3 编写 `specs/frontend-experience-tests/spec.md` delta。
- [x] 3.4 archive 时将 P57 路线图合并到 `docs/frontend-contract.md`、`docs/development-plan.md` 和进度文档。

## 4. 校验与审查

- [x] 4.1 运行 `openspec validate p57-product-experience-polish-roadmap --strict`。
- [x] 4.2 运行 `openspec validate --all --strict`。
- [x] 4.3 运行 `git diff --check`。
- [x] 4.4 执行敏感信息扫描，确认无 key、完整 prompt、私有 SQLite、raw vendor payload、供应商原始响应或敏感路径泄露。
- [x] 4.5 子 agent 方案复审无 Critical / Important 后归档。

## 5. 归档与提交

- [x] 5.1 执行 OpenSpec archive，将 delta 合并到 `docs/` 真源并同步规格摘要。
- [x] 5.2 archive 后确认无活跃 change，更新 `openspec/PROGRESS.md`、`docs/GOVERNANCE.md`、`openspec/project.md`、`docs/development-plan.md`、`AGENTS.md`。
- [x] 5.3 提交前子 agent 复审无 Critical / Important。
- [x] 5.4 提交 P57。
