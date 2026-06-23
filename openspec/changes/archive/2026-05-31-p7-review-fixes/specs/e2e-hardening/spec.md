## MODIFIED Requirements

### Requirement: P7 verification preserves hardening boundaries
P7 实现 SHALL 保留 P6 已确立的验收边界，并为真实数据、RAG/VecLite 与 DeepSeek 路径补充可验证场景。

#### Scenario: P7 backend verification passes
- **WHEN** P7 审查修复完成
- **THEN** `go test ./internal/infrastructure/... ./internal/application/...` MUST pass

#### Scenario: P7 retrieval verification passes
- **WHEN** P7 检索降级修复完成
- **THEN** `go test ./internal/infrastructure/... ./internal/application/workflow/...` MUST pass

#### Scenario: P7 analyst verification passes
- **WHEN** P7 分析服务修复完成
- **THEN** `go test ./internal/application/workflow/... ./internal/infrastructure/...` MUST pass

#### Scenario: P7 does not weaken safety boundaries
- **WHEN** 真实数据、RAG/VecLite 或 DeepSeek 路径被启用
- **THEN** DeepSeek MUST NOT 生成最终裁决
- **AND** 系统 MUST NOT 自动下单
- **AND** C 级信源 MUST NOT 作为正式裁决依据
