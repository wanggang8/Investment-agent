## 1. OpenSpec 与范围

- [x] 1.1 确认 P42、P43、P44 与 P19-P24 历史追溯均已归档，当前无其他活跃 change。
- [x] 1.2 确认 P45 只做 P44 后路线图治理，不修改运行时代码、数据库 schema、HTTP API、Eino 工作流或前端页面。
- [x] 1.3 明确 P44 后任何功能实现都必须先创建独立 OpenSpec change。

## 2. P44 后候选队列

- [x] 2.1 在 `docs/development-plan.md` 中增加 P45 路线图治理章节。
- [x] 2.2 将后续候选方向拆分为本地知识治理、决策闭环解释、数据质量回归和运维发布体验。
- [x] 2.3 为每类候选方向说明建议阶段、change id、依赖、验收思路和适合度。
- [x] 2.4 推荐 P46 `p46-local-knowledge-import-governance` 作为下一功能阶段，但不得在 P45 中实现 P46。

## 3. 安全边界

- [x] 3.1 保留禁止券商 API、自动交易、一键交易、代下单、外部推送、自动确认、自动应用规则、自动修复承诺和收益承诺的边界。
- [x] 3.2 保留登录源、付费源、授权源、Level2、高频源默认不纳入的边界。
- [x] 3.3 保留 LLM 只生成分析材料、不写最终裁决的边界。

## 4. 计划复审

- [x] 4.1 计划完成后执行只读子 agent 复审，确认无 Critical / Important。
- [x] 4.2 复审通过后再执行文档同步任务。

## 5. 治理文档同步

- [x] 5.1 更新 `openspec/PROGRESS.md`，将当前活跃 change 指向 `p45-post-p44-roadmap-governance`。
- [x] 5.2 更新 `docs/GOVERNANCE.md`，写入当前活跃 change、P45 范围与执行边界。
- [x] 5.3 更新 `openspec/project.md` 的阶段映射，加入 P45 和 P46-P49 候选。
- [x] 5.4 更新 `docs/development-plan.md` 顶部阶段状态与尾部阶段建议，清理 P44 作为下一阶段的过期描述。
- [x] 5.5 在 OpenSpec delta 中记录 P44 后路线图治理规则。
- [x] 5.6 更新 `AGENTS.md`，清理 P44 作为下一阶段的过期协作指引。

## 6. 验收

- [x] 6.1 运行 `openspec validate p45-post-p44-roadmap-governance --strict`。
- [x] 6.2 运行 `openspec validate --all --strict`。
- [x] 6.3 运行 `git diff --check`。
- [x] 6.4 确认 `git status --short` 只包含预期治理文档、`AGENTS.md` 和 OpenSpec 变更。

## 7. 归档前复审

- [x] 7.1 执行完成后再次只读子 agent 复审，确认无 Critical / Important。
- [x] 7.2 复审通过后执行 archive，并将 P45 归档。
