## ADDED Requirements

### Requirement: Frontend must consume stable error codes

系统 SHALL 让前端只依赖稳定错误码和展示状态，不直接解析底层错误文本。

#### Scenario: API returns evidence missing

- **WHEN** API 返回 `EVIDENCE_NOT_FOUND`
- **THEN** 前端必须展示信息不足状态
- **AND** 不得依赖后端原始错误字符串判断页面状态

#### Scenario: API returns source verification failed

- **WHEN** API 返回 `SOURCE_VERIFICATION_FAILED`
- **THEN** 前端必须展示冻结观察或证据验证失败状态

#### Scenario: API returns internal error

- **WHEN** API 返回 `INTERNAL_ERROR`
- **THEN** 前端必须展示通用失败状态
- **AND** 不得显示 SQL、文件路径或外部服务原始错误

### Requirement: Frontend contract tests must verify error-state mapping

系统 SHALL 在 P5 前端实现中覆盖错误码到 ViewModel 状态的映射测试。

#### Scenario: Error response is mapped

- **WHEN** 前端 API client 接收错误响应信封
- **THEN** ViewModel 必须映射为契约定义的页面状态
- **AND** 展示文案必须来自前端契约或前端文案表
