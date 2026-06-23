## MODIFIED Requirements

### Requirement: P7 verification preserves hardening boundaries
P7 与 P8 实现 SHALL 保留 P6 已确立的验收边界，并为真实数据、RAG/VecLite、DeepSeek 路径和前端体验测试补充可验证场景。

#### Scenario: P7 backend verification passes
- **WHEN** P7 审查修复完成
- **THEN** `go test ./internal/infrastructure/... ./internal/application/...` MUST pass

#### Scenario: P7 retrieval verification passes
- **WHEN** P7 检索降级修复完成
- **THEN** `go test ./internal/infrastructure/... ./internal/application/workflow/...` MUST pass

#### Scenario: P7 analyst verification passes
- **WHEN** P7 分析服务修复完成
- **THEN** `go test ./internal/application/workflow/... ./internal/infrastructure/...` MUST pass

#### Scenario: P8 frontend build passes
- **WHEN** P8.1 驾驶舱图表与关键交互实现完成
- **THEN** `cd web && npm run build` MUST pass

#### Scenario: P8 frontend tests pass
- **WHEN** P8.2 前端测试与契约校验实现完成
- **THEN** `cd web && npm run build && npm test` MUST pass

#### Scenario: P8 does not weaken safety boundaries
- **WHEN** 前端图表、交互或测试被启用
- **THEN** 前端 MUST NOT 提供自动交易入口
- **AND** 前端 MUST NOT 提供一键交易入口
- **AND** 用户确认区 MUST 继续表达为线下动作记录
