## 1. P32 归档后治理状态

- [x] 1.1 确认 `p32-daily-discipline-report-productization` 已进入 `openspec/changes/archive/2026-06-11-p32-daily-discipline-report-productization/`。
- [x] 1.2 确认 `openspec/specs/daily-discipline-report/spec.md` 已记录 P32 行为摘要。
- [x] 1.3 更新 `openspec/PROGRESS.md`：P32 标记为 done，current phase 进入 P33 roadmap preparation，next change 指向 P33。
- [x] 1.4 更新 `docs/GOVERNANCE.md` 和 `AGENTS.md`，移除 P31/P32 旧活跃描述，写入当前路线图治理 change。

## 2. P33-P40 路线图固化

- [x] 2.1 更新 `docs/development-plan.md`，将 P32 标记为已归档，并补充 P33–P40 为当前剩余计划内功能队列。
- [x] 2.2 为 P33 账户与持仓录入/校准体验补充任务组：账户初始化、持仓维护、线下交易流水、一致性校验、错误修正和首次使用引导。
- [x] 2.3 为 P34 真实数据覆盖扩展补充任务组：指数样本/权重/估值、财务/资金/情绪类数据、新鲜度/失败分类、工作流接入和健康状态。
- [x] 2.4 为 P35 风险预警与 SOP 编排补充任务组：风险中心、SOP 状态流转、通知/审计/报告写入、解除/升级和前端展示。
- [x] 2.5 为 P36 规则进化效果验证补充任务组：错误聚类、规则前后对比、样本代表性、过拟合检查、历史回放和应用后追踪。
- [x] 2.6 为 P37 真实 LLM 使用与质量评估补充任务组：真实配置 smoke、prompt 版本、错误分类、输出质量 fixture 和调用审计。
- [x] 2.7 为 P38 RAG / VecLite 检索质量加固补充任务组：测试集、混合检索或重排、引用一致性、新鲜度/重建和降级展示。
- [x] 2.8 为 P39 前端完整用户旅程与全路径 E2E 补充任务组：新手引导、数据源配置、账户初始化、跨页体验、Playwright E2E、窄屏/可访问性/console error 检查。
- [x] 2.9 为 P40 本地部署、运维与恢复演练补充任务组：初始化检查、启动前自检、备份恢复 E2E、数据源健康面板、日志/临时文件治理。

## 3. 边界与遗漏判断

- [x] 3.1 在 `docs/development-plan.md` 中明确 P33–P40 是当前计划内剩余功能队列。
- [x] 3.2 在 `docs/development-plan.md` 中明确 P19–P24 历史 archive 追溯属于独立治理候选，不属于 P33–P40 功能实现。
- [x] 3.3 在 `docs/development-plan.md` 中明确 P40 后新增产品愿景必须另建路线图 change。
- [x] 3.4 保留并复核禁止自动交易、外部推送、登录源、付费源、授权源、Level2、收益承诺和确定性涨跌预测边界。

## 4. 验收

- [x] 4.1 运行 `openspec validate p33-p40-roadmap-finalization --strict`。
- [x] 4.2 运行 `openspec validate --all --strict`。
- [x] 4.3 运行 `git status --short`，确认只包含预期文档、OpenSpec 归档和路线图变更。
- [x] 4.4 本 change 不改运行时代码；无需运行 Go、前端或 E2E 测试，除非后续修改了实现文件。
