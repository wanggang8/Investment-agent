# Tasks: 补充 P6 后续总计划

## 1. 更新开发总计划

- [x] 1.1 在 `docs/development-plan.md` 阶段总览中追加 P7、P8、P9。
- [x] 1.2 在 `docs/development-plan.md` 任务分解中追加 P7：真实数据、RAG/VecLite、DeepSeek 分析师、降级和审计任务。
- [x] 1.3 在 `docs/development-plan.md` 任务分解中追加 P8：图表展示、交互体验、空态/错误态、前端测试和构建验收任务。
- [x] 1.4 在 `docs/development-plan.md` 任务分解中追加 P9：`cmd/agent` 本地任务、每日任务入口、月度/季度复盘、规则有效性评估和交付说明任务。
- [x] 1.5 在 `docs/development-plan.md` 依赖关系或阶段说明中补充 P7 → P8 → P9 的先后关系。
- [x] 1.6 在 `docs/development-plan.md` 开发验收总清单中追加 P7–P9 相关验收项。
- [x] 1.7 在 `docs/development-plan.md` 建议提交节奏中追加 P7、P8、P9 提交建议。

## 2. 更新 OpenSpec 进度与映射

- [x] 2.1 更新 `openspec/PROGRESS.md`，将 `next_phase` 设置为 `P7`。
- [x] 2.2 更新 `openspec/PROGRESS.md`，将 `next_change_id` 设置为 `p7-real-data-integration`。
- [x] 2.3 更新 `openspec/PROGRESS.md` 阶段状态表，追加 P7、P8、P9，状态为 `pending`。
- [x] 2.4 更新 `openspec/project.md` 阶段映射，追加 P7 `p7-real-data-integration`。
- [x] 2.5 更新 `openspec/project.md` 阶段映射，追加 P8 `p8-frontend-experience-tests`。
- [x] 2.6 更新 `openspec/project.md` 阶段映射，追加 P9 `p9-review-automation-delivery`。

## 3. 一致性检查

- [x] 3.1 确认新增 P7–P9 不修改 P0–P6 已完成记录。
- [x] 3.2 确认本 change 只修改规划与治理文档，不改实现代码。
- [x] 3.3 确认 P7–P9 没有加入 requirements 与 architecture 之外的新产品能力。
- [x] 3.4 确认计划文本保留安全边界：禁止自动交易、DeepSeek 不写最终裁决、规则提案需审计和用户确认。
- [x] 3.5 执行 `openspec status --change extend-post-p6-roadmap`，确认 artifacts 完整。
