# Tasks: p6-e2e-hardening

> 对齐 `docs/development-plan.md` P6：验收加固。实现代码必须写必要中文注释，说明复杂业务边界、降级原因、事务语义、审计语义和禁止自动交易约束；不得写只复述代码的注释。

## 1. 端到端场景（P6.1）

创建或确认以下文件：

```text
docs/testing-plan.md
```

必须覆盖 `docs/functional-spec.md` 的 A01-A17 可测试验收断言：

- [x] 1.1 创建或确认 `docs/testing-plan.md`
- [x] 1.2 覆盖 A01 首次使用：无账户数据，展示引导，且不创建 `decision_records`
- [x] 1.3 覆盖 A02 正常每日纪律：生成建议、证据、审计事件
- [x] 1.4 覆盖 A03 证据不足：返回 `EVIDENCE_NOT_FOUND`，暂停交易类建议
- [x] 1.5 覆盖 A04 VecLite 不可用：SQLite 摘要充足时降级展示；不足时信息不足
- [x] 1.6 覆盖 A05 能力圈外：拒绝交易类分析
- [x] 1.7 覆盖 A06 用户记录计划：写 `operation_confirmations`，不更新账户
- [x] 1.8 覆盖 A07 用户记录已手动执行：写 `operation_confirmations`、`position_transactions`、`positions`、`portfolio_snapshots`、`position_snapshots`、`audit_events`
- [x] 1.9 覆盖 A08 已手动执行失败：事务回滚，不留下部分确认记录
- [x] 1.10 覆盖 A09 用户标记错误：写 `operation_confirmations`、`error_cases`、`audit_events`，返回 `error_case_id`
- [x] 1.11 覆盖 A10 C 级信源：只能作为 `background`，不得进入正式裁决证据
- [x] 1.12 覆盖 A11 LLM 不可用：规则引擎降级裁决，`workflow_status=degraded`
- [x] 1.13 覆盖 A12 守门人审计通过：进入 `pending_final_confirm`，不写 `rule_versions`；`sample_count<3` 的提案不得进入守门人审计，接口返回 `BAD_REQUEST`
- [x] 1.14 覆盖 A13 规则最终确认：创建新 active `rule_versions`，旧 active 归档；`sample_count<3` 的提案不得最终确认，接口返回 `BAD_REQUEST` 且不写 `rule_versions`
- [x] 1.15 覆盖 A14 审计事件：前端区分 `action`、`node_name` 和 `node_action`
- [x] 1.16 覆盖 A15 禁止自动交易：不存在交易执行接口，前端无一键交易入口
- [x] 1.17 覆盖 A16 市场数据刷新：全部成功新增快照且审计成功；部分成功返回 200、写 `failed_symbols` 且审计降级；全部失败返回 `DATA_SOURCE_UNAVAILABLE` 或 `DATA_STALE`；写入失败返回 `MARKET_SNAPSHOT_WRITE_FAILED` 且市场快照事务回滚
- [x] 1.18 覆盖 A17 预期收益评估：只展示情景概率，不覆盖最终规则裁决，不承诺收益；`available` 必含 upside/base/downside 且可返回概率；`insufficient` 不返回精确概率，`probability=null` 且写样本不足说明；`unavailable` 返回 `scenarios=[]` 且写定性原因
- [x] 1.19 若 A01-A17 暴露实现缺口，只按既有契约修正相关后端、前端或测试代码，并补充必要中文注释

## 2. P6.1 验收命令

验收：

```bash
go test ./...
cd web && npm run build
```

- [x] 2.1 执行 `go test ./...`
- [x] 2.2 执行 `cd web && npm run build`
- [x] 2.3 将验收结果记录到 `docs/testing-plan.md` 或本阶段可追踪验收记录中

## 3. 配置与启动文档（P6.2）

创建或确认以下文件：

```text
docs/configuration.md
docs/migration-plan.md
```

必须包含：

- [x] 3.1 创建或确认 `docs/configuration.md`
- [x] 3.2 创建或确认 `docs/migration-plan.md`
- [x] 3.3 文档覆盖 SQLite 数据文件路径
- [x] 3.4 文档覆盖 VecLite 索引文件路径
- [x] 3.5 文档覆盖 DeepSeek API Key 环境变量
- [x] 3.6 文档覆盖数据源开关
- [x] 3.7 文档覆盖日志级别
- [x] 3.8 文档覆盖 migration 执行方式
- [x] 3.9 文档覆盖 seed 数据说明
- [x] 3.10 文档覆盖本地启动命令
- [x] 3.11 确认文档不包含真实密钥、私有路径或计划外需求
