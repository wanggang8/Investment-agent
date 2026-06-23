# Delta for API（合并目标：`docs/api.md`）

## ADDED Requirements

### Requirement: Health Check

系统 SHALL 提供无需认证的健康检查端点，供部署与本地开发探测服务存活。

#### Scenario: 服务正常

- **WHEN** 客户端 `GET /api/v1/health`
- **THEN** 响应 HTTP 200
- **AND** 响应体为 JSON：`{"status":"ok"}`

## MODIFIED Requirements

（无）

## REMOVED Requirements

（无）
