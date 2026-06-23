# Tasks: p3-foundation-governance

> 目标：在 P4 HTTP API 前补齐横向基础能力，统一错误、ID、时间、事务、审计和测试策略。实现代码必须写中文注释，说明边界和迁移原因。

## 1. 统一错误体系

- [x] 1.1 创建 `internal/pkg/apperr` 包
- [x] 1.2 定义 `Code`、`Category`、`AppError`、`New`、`Wrap`、`IsCode`、`AsAppError`
- [x] 1.3 定义通用错误码：`INTERNAL_ERROR`、`BAD_REQUEST`、`NOT_FOUND`、`CONFLICT`、`INVALID_STATE`
- [x] 1.4 迁移工作流错误码为 `apperr.Code` 兼容常量
- [x] 1.5 将 `ErrInvalidConfirmationTransition` 等仓储业务错误映射为统一错误
- [x] 1.6 定义错误码到 HTTP status 的映射函数，供 P4 handler 使用
- [x] 1.7 定义错误码到审计 `error_code` 的兼容映射
- [x] 1.8 新增 `internal/pkg/apperr` 单元测试，覆盖包装、分类、HTTP 映射和 `errors.Is/As`

## 2. 统一 ID 与时间

- [x] 2.1 创建 `internal/pkg/clock` 包
- [x] 2.2 定义 `Clock`、`SystemClock`、`FixedClock`、`FormatRFC3339UTC`
- [x] 2.3 创建 `internal/pkg/idgen` 包
- [x] 2.4 定义 `DecisionID`、`EvidenceRefID`、`AuditEventID`、`TransactionID`、`RuleAppliedVersionID` 等生成函数
- [x] 2.5 迁移 workflow 与 repository 中关键 `fmt.Sprintf("xxx_%s")` ID 拼接到 `idgen`
- [x] 2.6 迁移 workflow 与 repository 中关键 `time.Now().UTC().Format(time.RFC3339)` 到 `clock`
- [x] 2.7 新增 ID 和时间单元测试，覆盖空输入、稳定输出、UTC/RFC3339、固定时间

## 3. 事务边界与仓储接口

- [x] 3.1 在 `DecisionRepository` 增加或调整组合方法，确保 `decision_records`、`evidence_refs`、DecisionRecordNode 审计事件同事务写入
- [x] 3.2 保持 `SaveOperationConfirmation` 用户确认事务边界，并改用统一错误和 ID/时间工具
- [x] 3.3 保持 `SaveEvidenceFacts` 证据事实事务边界，并补失败回滚测试
- [x] 3.4 保持 `SaveGatekeeperAuditAndUpdateProposalStatus` 守门人审计事务边界，并补失败回滚测试
- [x] 3.5 保持 `ApplyRuleVersion` 规则应用事务边界，确认新版本状态由仓储强制为 `active`
- [x] 3.6 为所有跨表事务补字段级成功断言和失败回滚断言

## 4. 审计事件契约

- [x] 4.1 整理审计 action、node_name、node_action、input/output ref 的代码常量或集中映射
- [x] 4.2 让 `AuditWriter` 使用统一 ID、时间和错误码映射
- [x] 4.3 补齐 Daily/Consultation 主链路审计字段测试
- [x] 4.4 补齐 Evidence/Market/Evolution/Gatekeeper 辅助 Graph 审计字段测试
- [x] 4.5 确认失败和降级审计均包含 `error_code`

## 5. 测试策略与契约验证

- [x] 5.1 补充 P0 health handler 或 httputil 基础错误响应测试，不实现 P4 业务 API
- [x] 5.2 补充 workflow 分支测试：行情缺失、规则版本缺失、能力圈 unknown、source verification failed
- [x] 5.3 补充 repository 错误分类测试：not_found、conflict、invalid_state、constraint
- [x] 5.4 补充前端 API client 对错误信封的基础解析测试或类型测试，不实现业务页面
- [x] 5.5 运行 `go test ./...`
- [x] 5.6 运行 `cd web && npm run build`

## 6. 文档与归档准备

- [x] 6.1 确认 `specs/foundation-governance/spec.md` 覆盖统一错误、ID、时间、事务、审计和测试策略
- [x] 6.2 确认 `specs/workflow/spec.md` 覆盖工作流与基础治理关系
- [x] 6.3 确认 `specs/data-model/spec.md` 覆盖事务、ID、时间和仓储错误分类
- [x] 6.4 确认 `specs/api/spec.md` 覆盖 P4 错误响应映射
- [x] 6.5 确认 `specs/frontend-contract/spec.md` 覆盖 P5 错误展示映射
- [x] 6.6 实现完成后再次多角度审查 P0-P3 与本变更是否匹配
- [x] 6.7 归档时将 delta 合并到对应 L1/L2 文档并更新 `openspec/PROGRESS.md`
