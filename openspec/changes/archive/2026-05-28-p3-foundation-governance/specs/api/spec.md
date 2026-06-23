## ADDED Requirements

### Requirement: API errors must be derived from application errors

系统 SHALL 在 HTTP API 层只根据统一应用错误生成错误响应。

#### Scenario: Application error reaches handler

- **WHEN** handler 接收到应用错误
- **THEN** HTTP status 必须来自错误映射
- **AND** 响应体必须包含 `request_id`
- **AND** 响应体必须包含 `error.code`、`error.message`

#### Scenario: Unknown error reaches handler

- **WHEN** handler 接收到未知错误
- **THEN** HTTP status 必须为 500
- **AND** `error.code` 必须为 `INTERNAL_ERROR`
- **AND** 不得暴露底层错误详情

### Requirement: API error mapping must be stable

系统 SHALL 为 P4 API 提供稳定错误码到 HTTP 状态的映射表。

#### Scenario: Evidence is missing

- **WHEN** 应用错误码为 `EVIDENCE_NOT_FOUND`
- **THEN** HTTP status 必须为 409
- **AND** 前端必须能显示信息不足状态

#### Scenario: Source verification fails

- **WHEN** 应用错误码为 `SOURCE_VERIFICATION_FAILED`
- **THEN** HTTP status 必须为 409
- **AND** 前端必须能显示冻结观察状态

#### Scenario: Validation fails

- **WHEN** 输入参数或状态流转不合法
- **THEN** HTTP status 必须为 400
- **AND** 错误码必须来自统一错误体系
