# E2E Hardening Specification

## Purpose
Document end-to-end acceptance and verification hardening for the Investment Agent workflow.
## Requirements
### Requirement: End-to-end acceptance plan
The P6 change SHALL provide an end-to-end acceptance plan that covers all A01-A17 testable acceptance assertions from `docs/functional-spec.md` without adding requirements beyond `docs/development-plan.md`.

#### Scenario: A01-A17 assertions are covered
- **WHEN** `docs/testing-plan.md` is reviewed
- **THEN** it SHALL include acceptance coverage for A01 through A17
- **AND** each assertion SHALL state expected observable outcomes aligned with `docs/development-plan.md` P6.1

#### Scenario: No automatic trading acceptance path
- **WHEN** the A15 acceptance assertion is verified
- **THEN** the plan SHALL confirm that no trade execution API or one-click trading frontend entry exists
- **AND** user actions SHALL remain limited to recording offline actions

#### Scenario: Evidence and degradation paths are represented
- **WHEN** A03, A04, A10, A11, A16, and A17 are verified
- **THEN** the plan SHALL cover evidence insufficiency, VecLite degradation, C-level source handling, LLM degradation, market refresh outcomes, and expected-return display states

### Requirement: Configuration and startup documentation
The P6 change SHALL provide configuration and startup documentation for local operation, migration, and seed data as listed in `docs/development-plan.md` P6.2.

#### Scenario: Runtime configuration is documented
- **WHEN** `docs/configuration.md` is reviewed
- **THEN** it SHALL document SQLite data file path, VecLite index file path, DeepSeek API Key environment variable, data source switches, log level, and local startup commands

#### Scenario: Migration and seed are documented
- **WHEN** `docs/migration-plan.md` is reviewed
- **THEN** it SHALL document migration execution and seed data behavior
- **AND** it SHALL avoid real secrets or environment-specific private values

### Requirement: P6 verification commands
The P6 implementation SHALL preserve the verification commands defined in `docs/development-plan.md` for the acceptance hardening phase.

#### Scenario: Backend and frontend verification pass
- **WHEN** P6 implementation is completed
- **THEN** `go test ./...` SHALL pass
- **AND** `cd web && npm run build` SHALL pass

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

#### Scenario: P7 does not weaken safety boundaries
- **WHEN** 真实数据、RAG/VecLite 或 DeepSeek 路径被启用
- **THEN** DeepSeek MUST NOT 生成最终裁决
- **AND** 系统 MUST NOT 自动下单
- **AND** C 级信源 MUST NOT 作为正式裁决依据

### Requirement: P8 review hardening coverage
The system SHALL include regression tests and documented verification gates for decision confirmation, rule proposal flow, account snapshot consistency, and frontend behavior.

#### Scenario: Confirmation transaction is atomic
- **WHEN** a manual execution confirmation fails during any dependent write
- **THEN** tests SHALL verify that no partial state remains.

#### Scenario: Consult-to-confirm is executable when contract allows it
- **WHEN** a decision detail advertises offline confirmation actions
- **THEN** tests SHALL verify the confirmation endpoint accepts the corresponding request.

#### Scenario: Full validation includes frontend tests
- **WHEN** the project is prepared for review or archive
- **THEN** the documented validation SHALL include `go test ./...`, `cd web && npm run build`, and `cd web && npm test`.

### Requirement: P9 delivery acceptance is verifiable
The system SHALL include P9 acceptance coverage for local task entrypoints, periodic review behavior, frontend build, and local delivery documentation.

#### Scenario: Local agent command is accepted
- **WHEN** P9 acceptance is executed
- **THEN** `go test ./...` and `go run ./cmd/agent --help` complete successfully.

#### Scenario: Review and handler coverage is accepted
- **WHEN** P9 review automation acceptance is executed
- **THEN** `go test ./internal/application/workflow/... ./internal/application/handler/...` completes successfully.

#### Scenario: Frontend delivery build is accepted
- **WHEN** P9 frontend acceptance is executed
- **THEN** `cd web && npm run build` completes successfully.

#### Scenario: Full P7-P9 delivery commands are documented
- **WHEN** local delivery documentation is reviewed
- **THEN** it lists the complete P7-P9 validation commands required by the development plan.

### Requirement: P39 Browser Full User Journey Acceptance
The system SHALL provide browser-level E2E coverage for a complete local user journey from first use through daily report review, active consultation, offline confirmation, audit inspection, periodic review, rule governance, and risk alert inspection.

#### Scenario: Empty local state reaches first daily report safely
- **WHEN** the E2E fixture starts from an empty or missing-prerequisite local state
- **THEN** the browser journey SHALL expose safe onboarding or prerequisite guidance for configuration, account, and position setup
- **AND** the journey SHALL reach a first daily discipline report using fixed-ID local fixture data
- **AND** it SHALL NOT require public network access, real secrets, broker credentials, or manual database inspection

#### Scenario: Consultation to confirmation remains a local record flow
- **WHEN** the browser journey performs an active consultation and opens the resulting decision detail
- **THEN** the page SHALL expose decision trace, evidence, retrieval quality, and audit references where available
- **AND** any confirmation action SHALL record only an offline user fact
- **AND** the journey SHALL NOT expose automatic trading, broker order placement, one-click order placement, or portfolio mutation without user-recorded offline confirmation

#### Scenario: Review and rule governance are inspectable but not automatic
- **WHEN** the browser journey opens periodic review and rule governance surfaces
- **THEN** it SHALL show review summaries, rule proposal status, gatekeeper or final confirmation boundaries, and tracking entrypoints where available
- **AND** pending proposals SHALL remain visible as review/governance facts
- **AND** the journey SHALL NOT automatically apply rules or bypass gatekeeper audit and final user confirmation

#### Scenario: Existing P34 P35 P38 statuses are included
- **WHEN** the browser journey traverses dashboard, evidence, decision, review, rules, and risk alert pages
- **THEN** it SHALL include source health, risk alert/SOP, and retrieval quality states from the existing API/service DTOs
- **AND** degraded, empty, missing, or unknown states SHALL be visible as safe non-success states

### Requirement: P40 Local Recovery Smoke Acceptance
The system SHALL include local smoke coverage for deployment readiness and backup recovery that verifies restored data through supported APIs or browser-visible behavior.

#### Scenario: Recovery smoke verifies readable restored facts
- **WHEN** the P40 recovery smoke restores a backup into temporary local paths
- **THEN** it SHALL run migrations or compatibility checks as needed
- **AND** it SHALL verify at least one restored decision, audit event, position, report, or equivalent historical fact through API or browser-visible behavior
- **AND** it SHALL NOT require manual SQLite inspection

#### Scenario: Recovery smoke preserves safety boundaries
- **WHEN** recovery smoke, startup smoke, or browser smoke executes
- **THEN** it SHALL use temporary local paths and non-secret fixture data by default
- **AND** it SHALL NOT connect to broker APIs, place orders, send external pushes, automatically confirm user actions, or automatically apply rules
