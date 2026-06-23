# Delta for Domain Rules（合并目标：`docs/workflow.md`、`docs/data-model.md`）

## ADDED Requirements

### Requirement: Domain enums must validate contract values

系统 SHALL 在领域层定义并验证 API、数据模型与工作流中使用的核心枚举。

#### Scenario: Valid enum value

- **GIVEN** 一个来自契约文档的合法枚举值
- **WHEN** 调用领域枚举校验
- **THEN** 校验通过

#### Scenario: Invalid enum value

- **GIVEN** 一个契约外枚举值
- **WHEN** 调用领域枚举校验
- **THEN** 校验失败

### Requirement: Rule engine must prioritize safety and discipline

系统 SHALL 由领域规则引擎给出最终规则裁决，大模型输出不得覆盖最终裁决。

#### Scenario: Out of capability scope

- **GIVEN** 标的不在能力圈内
- **WHEN** 执行规则裁决
- **THEN** `final_verdict.status` 为 `rejected`
- **AND** 拒绝交易类分析

#### Scenario: Insufficient evidence

- **GIVEN** 有效证据不足
- **WHEN** 执行规则裁决
- **THEN** `final_verdict.status` 为 `insufficient_data`

#### Scenario: Major event lacks independent high-grade sources

- **GIVEN** 重大利好、重大利空或买入逻辑破坏事件未满足至少 2 个 A/S 独立信源
- **WHEN** 执行规则裁决
- **THEN** `final_verdict.status` 为 `frozen_watch`

#### Scenario: Buy logic broken

- **GIVEN** 买入逻辑破坏且证据满足正式裁决要求
- **WHEN** 执行规则裁决
- **THEN** `final_verdict.status` 为 `sell_only`

### Requirement: Valuation, cash, portfolio and take-profit rules must be deterministic

系统 SHALL 以确定性领域规则处理 PE/PB 分位、现金冗余、核心-卫星仓位和移动止盈。

#### Scenario: Valuation bands

- **WHEN** PE/PB 分位落在高危、观察、舒适、低估区间
- **THEN** 输出对应禁止买入、持有观察、按计划定投或分批配置建议

#### Scenario: Cash redundancy

- **WHEN** 现金比例低于 5%
- **THEN** 限制新增买入
- **AND** 现金比例 5%-10% 不因现金规则单独禁止交易

#### Scenario: Take-profit stages

- **WHEN** 浮盈达到 20%、30% 或已启动移动止盈后回撤 10%
- **THEN** 输出对应分批止盈或减仓/卖出评估动作

### Requirement: Rule proposal state machine must guard rule evolution

系统 SHALL 在领域层实现规则提案状态机，确保规则应用必须经过用户确认与守门人审计。

#### Scenario: Final confirm with enough samples

- **GIVEN** 提案状态为 `pending_final_confirm` 且 `sample_count>=3`
- **WHEN** 用户最终确认
- **THEN** 提案状态变为 `applied`
- **AND** 允许创建新 active 规则版本

#### Scenario: Sample count insufficient

- **GIVEN** 提案样本数小于 3
- **WHEN** 尝试进入最终应用路径
- **THEN** 返回错误且不允许写入正式规则版本

## MODIFIED Requirements

（无）

## REMOVED Requirements

（无）
