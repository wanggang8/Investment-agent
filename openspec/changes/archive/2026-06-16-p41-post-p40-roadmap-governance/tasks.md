## 1. OpenSpec 与范围

- [x] 1.1 确认 P40 已归档，当前无其他活跃 change。
- [x] 1.2 确认本 change 只做 P40 后路线图治理，不修改运行时代码、数据库 schema、HTTP API 或前端页面。
- [x] 1.3 明确 P40 后新增功能必须先创建独立 OpenSpec change，不得直接实现。

## 2. P41+ 候选队列

- [x] 2.1 在 `docs/development-plan.md` 中增加 P41+ 路线图治理章节。
- [x] 2.2 将后续候选方向拆分为产品能力增强、数据质量增强、运维体验增强和历史审计追溯。
- [x] 2.3 为每类候选方向说明依赖、验收思路和是否适合下一阶段优先处理。
- [x] 2.4 明确 P19-P24 历史 archive 追溯属于独立治理候选，不属于功能增强阶段。

## 3. 安全边界

- [x] 3.1 保留禁止券商 API、自动交易、外部推送、自动应用规则和收益承诺的边界。
- [x] 3.2 保留登录源、付费源、授权源、Level2、高频源默认不纳入的边界。
- [x] 3.3 保留 LLM 只生成分析材料、不写最终裁决的边界。

## 4. 治理文档同步

- [x] 4.1 更新 `openspec/PROGRESS.md`，将当前活跃 change 指向 `p41-post-p40-roadmap-governance`。
- [x] 4.2 更新 `docs/GOVERNANCE.md` 和 `AGENTS.md`，写入当前活跃 change 与执行边界。
- [x] 4.3 更新 `openspec/project.md` 的阶段映射，加入 P41 后路线图治理。
- [x] 4.4 在 OpenSpec delta 中记录 P40 后路线图治理规则。

## 5. 验收

- [x] 5.1 运行 `openspec validate p41-post-p40-roadmap-governance --strict`。
- [x] 5.2 运行 `openspec validate --all --strict`。
- [x] 5.3 运行 `git diff --check`。
- [x] 5.4 确认 `git status --short` 只包含预期文档和 OpenSpec 变更。
