# P28 预期收益与动态卖出评估增强任务

## 1. 前置确认

- [x] 复核现有预期收益实现
  - 确认 `BuildExpectedReturn`、`ExpectedReturnNode`、decision DTO、持久化 JSON 和前端契约当前字段。
  - 确认预期收益只作为分析材料，不覆盖 `final_verdict`。

- [x] 复核安全边界
  - 不接券商交易 API。
  - 不自动下单，不提供一键买卖。
  - 不预测确定涨跌，不承诺收益。
  - 不引入付费、登录、授权、Level2 或外部通知渠道。

## 2. Payload 与契约设计

- [x] 扩展预期收益 payload
  - 补充 `sample_window`、`screening_condition`、`scenario_trigger`、`sell_evaluation`、`reassessment_trigger`。
  - 保留 `sample_count`、`precision_status`、`reason`、`scenarios`、`disclaimer`。
  - 明确 schema：`sell_evaluation.status/triggers/prompts/actions/non_trading_disclaimer`；`reassessment_trigger.reason/boundary/current_value`；scenario 内可包含 `trigger` 说明。
  - 确保 decision record JSON 可复现历史详情；旧 JSON 缺字段时按空/不适用处理。

- [x] 设计样本门槛策略
  - `<5`：`unavailable`，不返回收益区间。
  - `5–19`：`insufficient`，可返回区间但不返回精确概率。
  - `>=20`：`available`，可返回概率。
  - 所有输出必须包含样本数量、样本窗口和筛选条件，缺失时写明原因。

- [x] 设计动态卖出评估策略
  - 当前净值进入乐观情景下沿：提示启动移动止盈评估。
  - 当前净值突破基准情景上沿：提示分批止盈评估。
  - 当前净值跌破悲观情景下沿：提示重新核验买入逻辑。
  - 基准情景中枢下移超过 15%：仅当本地存在可复现的上一轮基准中枢输入时提示重新评估买入逻辑并考虑减仓；当前 consultation workflow 尚无历史中枢字段时不得伪造触发。
  - 达到用户自定义目标收益：提示查看或记录人工止盈计划；目标收益仅来自本地持仓/设置中的明确字段，缺失时不触发。
  - 所有提示只进入可选动作/分析材料，不更新账户，不创建交易、确认、通知或外部推送。

## 3. 后端实现

- [x] 扩展 workflow 预期收益类型
  - 更新 `ExpectedReturnOutput`、`ExpectedReturnScenario` 或新增辅助结构。
  - 保持向后兼容现有 decision detail 读取。

- [x] 实现样本门槛与情景解释
  - 按样本数输出 precision status、reason、sample window、screening condition。
  - consultation response 从当前本地持仓、最新市场快照和已保存公开市场元数据派生可解释样本数；缺少历史样本时不得伪造成完整历史回测。
  - 不足样本时不得返回精确概率。

- [x] 实现动态卖出评估
  - 基于当前价格/净值、按标的匹配的持仓成本或本地可复现基准价格、情景上下沿和持仓收益状态生成 sell evaluation。
  - 明确价格基准优先级：优先匹配标的的持仓成本/确认记录可复现成本，其次市场快照当前价格；缺价格或缺基准时返回不适用原因。
  - 当前内置情景使用可复现 `return_rate` 边界；未来若接入字符串区间上下沿，无法解析时不得触发卖出评估。
  - 不改变最终裁决优先级；规则裁决仍由 `RuleArbitrationNode` 负责。

- [x] 持久化与 API DTO 对齐
  - `decision_records.expected_return_scenarios_json` 保存新增字段。
  - `GET /api/v1/decisions/{id}` 和 consultation response 返回新增字段。

## 4. 测试与验收

- [x] 增加预期收益单元测试
  - 覆盖 `<5`、`5–19`、`>=20` 三类样本门槛。
  - 覆盖概率是否允许展示；`insufficient/unavailable` 不泄露精确概率。
  - 覆盖 sample window、screening condition 和 scenario trigger。

- [x] 增加动态卖出评估测试
  - 当前净值进入乐观情景下沿。
  - 当前净值突破基准情景上沿。
  - 当前净值跌破悲观情景下沿。
  - 基准情景中枢下移超过 15%。
  - 达到用户目标收益；目标收益缺失时不触发。
  - 缺当前价格、缺基准价格或情景边界无法解析时返回不适用原因。
  - 验证只生成提示/可选动作，不更新账户、不自动交易、不创建 confirmation/transaction/notification。

- [x] 增加持久化/API 测试
  - 决策记录保存新增预期收益字段。
  - 决策详情读取历史 JSON 时能返回新增字段。
  - 旧 JSON 缺字段时保持可读。
  - handler/API 层断言 `<5` scenarios 为空，`5–19` probability absent/null，`>=20` 可展示 probability，且 disclaimer 始终存在。

- [x] 增加前端类型与渲染测试
  - 更新 `web/src/types/decision.ts`。
  - 更新 expected return 渲染组件，展示 sample window、screening condition、scenario trigger、sell evaluation 和 reassessment trigger。
  - 更新决策详情相关测试，覆盖新增 P28 字段展示和旧字段兼容。

- [x] 运行后端测试
  - `go test ./...`

- [x] 运行前端验证
  - 按 `web/package.json` 实际脚本运行相关测试或构建命令。

- [x] 运行 OpenSpec 校验
  - `openspec validate --all --strict`

## 5. 文档同步

- [x] 更新 `docs/development-plan.md`
  - 新增 P28 阶段状态、范围和验收目标。

- [x] 更新 `docs/api.md` 与 `docs/frontend-contract.md`
  - 补充 expected_return_scenarios 新字段。

- [x] 更新 `docs/workflow.md`、`docs/data-model.md`
  - 说明 ExpectedReturnNode 新增输入/输出、mutates、audit 字段。
  - 说明 `decision_records.expected_return_scenarios_json` 的新增 JSON shape、历史回放和旧 JSON 兼容策略。

- [x] 更新 `openspec/PROGRESS.md`、`docs/GOVERNANCE.md`、`AGENTS.md`
  - 归档后标记 P28 done，并清空当前 active change。
