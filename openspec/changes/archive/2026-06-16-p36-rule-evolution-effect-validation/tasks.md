## 1. OpenSpec 与范围

- [x] 1.1 确认 P36 只覆盖规则进化效果验证、提案来源解释、样本代表性、过拟合检查、历史回放、应用后追踪和前端展示。
- [x] 1.2 确认 P36 不自动应用规则、不自动交易、不接券商 API、不外部推送、不引入登录/付费/授权/Level2/高频源、不承诺收益、不预测确定涨跌。
- [x] 1.3 对齐 P22/P23 规则提案与复盘、P28 动态卖出评估、P35 风险预警、`docs/api.md`、`docs/data-model.md`、`docs/workflow.md`、`docs/frontend-contract.md` 的现有契约。

## 2. 数据模型与领域状态

- [x] 2.1 定义规则效果验证事实模型：validation ID、proposal ID、rule version、validation status、sample count/window、representativeness status、overfit risk、replay result、guardrail decision、metrics JSON、risk notes、related IDs、timestamps。
- [x] 2.2 定义应用后追踪事实模型或等价存储：tracking ID、applied rule version、proposal ID、period、hit/misjudgment/missing-evidence/degraded/risk-alert metrics、trend direction、related links、timestamps。
- [x] 2.3 增加 SQLite migration、domain model、repository interface 和 sqlite repository tests。
- [x] 2.4 定义并校验状态枚举：not_evaluated、insufficient、passed、failed、needs_more_samples、needs_user_review；overfit risk low/medium/high；trend improved/flat/worsened/unknown。
- [x] 2.5 确认验证与追踪写入不会更新 positions、portfolio_snapshots、operation_confirmations、position_transactions、broker state、orders 或 external notifications。

## 3. 效果验证服务

- [x] 3.1 实现 RuleEffectValidationService，聚合 rule proposals、error cases、decision records、confirmations、review facts、risk alerts 和 audit events。
- [x] 3.2 实现来源解释：输出关联错误案例、复盘周期、决策、确认、风险预警和审计线索，缺失事实显式标记。
- [x] 3.3 实现样本代表性与门槛判断：样本不足或来源过窄时输出 insufficient / needs_more_samples。
- [x] 3.4 实现过拟合检查：识别单样本调参、窄场景命中、冲突结果和风险恶化信号。
- [x] 3.5 实现历史回放验证：对比候选规则与基线规则在本地历史事实上的命中率、误判率、缺证据率、降级率和风险预警影响。
- [x] 3.6 增加服务 tests，覆盖样本不足、过拟合高、回放不利、验证通过、缺事实降级和非交易边界。

## 4. 守门人、规则提案与复盘接入

- [x] 4.1 将效果验证结果接入规则提案 DTO 和查询服务，展示 validation status、sample summary、overfit risk、replay result、guardrail decision 和 validation link。
- [x] 4.2 将效果验证结果接入守门人审计：样本不足、代表性不足、过拟合高或回放不利时拒绝或返回 needs_user_review。
- [x] 4.3 确认验证通过仍不能自动应用规则，必须保持 pending_final_confirm 与用户最终确认状态机。
- [x] 4.4 将应用后追踪接入 review summary，展示规则命中、误判、缺证据、降级和风险预警趋势。
- [x] 4.5 增加 workflow/service tests，覆盖提案验证、审计门禁、应用后追踪和复盘输出。

## 5. HTTP API 与前端

- [x] 5.1 新增或扩展规则效果验证 API：提案验证详情、触发/刷新验证、应用后追踪查询；统一响应信封、错误码和只读安全文案。
- [x] 5.2 更新 app routing/handler wiring，并增加 handler tests。
- [x] 5.3 新增前端 rule effect validation types、services、status mappers。
- [x] 5.4 更新规则提案详情页，展示来源解释、样本代表性、过拟合风险、历史回放结果、门禁结论和相关追踪链接。
- [x] 5.5 更新复盘页，展示规则应用后效果趋势和风险提示。
- [x] 5.6 增加前端 tests，覆盖 not_evaluated、insufficient、passed、failed、high overfit、worsened tracking、空状态和禁止自动应用文案。

## 6. 文档与验收

- [x] 6.1 在 P36 delta 中记录待归档合并到 `docs/api.md` 的 rule effect validation API/DTO/错误分类和事务边界。
- [x] 6.2 在 P36 delta 中记录待归档合并到 `docs/data-model.md` 的验证/追踪模型、状态枚举、索引和非交易约束。
- [x] 6.3 在 P36 delta 中记录待归档合并到 `docs/workflow.md` 与 `docs/frontend-contract.md` 的验证编排、复盘接入和前端展示。
- [x] 6.4 更新 `docs/development-plan.md`、`openspec/PROGRESS.md`、`AGENTS.md`、`docs/GOVERNANCE.md` 的 P36 active 状态。
- [x] 6.5 运行 `go test ./...`。
- [x] 6.6 运行 `npm --prefix web test -- --run`。
- [x] 6.7 运行 `npm --prefix web run build`。
- [x] 6.8 运行 P36 规则效果场景 smoke，覆盖 insufficient / high_overfit / replay_failed / passed / tracking_worsened。
- [x] 6.9 运行 `openspec validate p36-rule-evolution-effect-validation --strict`。
- [x] 6.10 运行 `openspec validate --all --strict`。
- [x] 6.11 运行 `git status --short`，确认只包含预期修改且无临时产物。