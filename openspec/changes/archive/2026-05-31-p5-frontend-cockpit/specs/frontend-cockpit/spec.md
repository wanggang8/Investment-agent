## ADDED Requirements

### Requirement: Frontend DTO and API service contract
The P5 frontend SHALL define typed DTOs and API services for dashboard, portfolio, decision, evidence, rule, audit, settings, market, and review data using the P4 HTTP API response envelope.

#### Scenario: API response envelope is parsed consistently
- **WHEN** any P5 service receives an HTTP response
- **THEN** it SHALL parse `request_id`, `data`, `meta`, and `error` from the common response envelope
- **AND** callers SHALL receive typed success data or a stable frontend error state

#### Scenario: P5 DTOs cover plan domains
- **WHEN** P5 pages consume dashboard, portfolio, decision, evidence, rule, audit, settings, market, or review data
- **THEN** the frontend SHALL use DTO types defined for those domains
- **AND** the DTO field names SHALL remain compatible with `docs/frontend-contract.md` and `docs/api.md`

#### Scenario: Frontend data source boundary
- **WHEN** a P5 page or component needs application data
- **THEN** it SHALL use `web/src/services/` or shared API services
- **AND** it SHALL NOT access SQLite, VecLite, or local files directly

#### Scenario: Stable HTTP error handling
- **WHEN** the API returns 409, 500, or 503 responses
- **THEN** the frontend SHALL map stable error codes to insufficient-data, frozen-watch, data-source, or generic-failure display states
- **AND** it SHALL NOT show SQL, file paths, full API keys, or raw external service errors

### Requirement: Agent cockpit dashboard
The P5 frontend SHALL implement the Agent decision cockpit for today's discipline view using the three-column cockpit layout described by the UI documents.

#### Scenario: Cockpit first screen content
- **WHEN** the user opens the dashboard page
- **THEN** the first screen SHALL show discipline status, risk lines, today recommendation, account summary, and evidence summary

#### Scenario: Three-column cockpit layout
- **WHEN** the dashboard renders on desktop
- **THEN** it SHALL organize navigation and system status on the left, decision work area in the center, and evidence/rule context on the right

#### Scenario: Insufficient data state
- **WHEN** the dashboard state is `insufficient_data` or the mapped error state indicates insufficient data
- **THEN** the page SHALL show missing items and suspension reasons
- **AND** it SHALL not present trade-type advice as actionable

#### Scenario: Frozen watch state
- **WHEN** the dashboard state is `frozen_watch`
- **THEN** the page SHALL show the waiting conditions and evidence verification status

#### Scenario: Confirmation action boundary
- **WHEN** the confirmation area is displayed
- **THEN** it SHALL only allow recording `planned`, `executed_manually`, `watch`, or `marked_error`
- **AND** it SHALL not display automatic trading, one-click buy, one-click sell, or delegated execution controls

### Requirement: Decision detail and supporting pages
The P5 frontend SHALL implement decision detail, evidence, rules, audit, portfolio, settings, and review pages according to the frontend contract and UI flow.

#### Scenario: Decision detail shows full trace
- **WHEN** a user opens a decision detail
- **THEN** the page SHALL show final verdict, summary, triggered rules, account snapshot, evidence chain, agent opinions, arbitration chain, user confirmation area, and audit information in the order required by `docs/ui-flow.md`
- **AND** agent opinions SHALL be presented as analysis material rather than final verdicts

#### Scenario: Portfolio page uses portfolio API
- **WHEN** the portfolio page needs current holdings
- **THEN** it SHALL use `GET /api/v1/portfolio/current`
- **AND** it SHALL not directly access SQLite or VecLite

#### Scenario: Evidence page shows verification fields
- **WHEN** the evidence page renders evidence records
- **THEN** it SHALL show `source_level`, `evidence_role`, and `verification_status`

#### Scenario: Rules page shows final confirmation
- **WHEN** a rule proposal has status `pending_final_confirm`
- **THEN** the rules page SHALL show final confirm and reject actions
- **AND** it SHALL make clear that the rule is not active until final confirmation succeeds

#### Scenario: Settings page protects secrets and rules
- **WHEN** the settings page displays system status and configuration
- **THEN** it SHALL show capability settings, system status, market snapshot status, notification settings, and index status
- **AND** it SHALL not display complete secret values

#### Scenario: Review page shows audit summary
- **WHEN** the review page renders summary data
- **THEN** it SHALL show decision count, confirmation actions, error cases, rule proposals, and audit event summary

#### Scenario: Audit page preserves audit semantics
- **WHEN** audit events are displayed
- **THEN** the page SHALL distinguish `action`, `node_name`, and `node_action`
- **AND** it SHALL show `status`, `error_code`, input references, and output references where available

### Requirement: P5 verification
P5 implementation SHALL pass the development-plan verification commands.

#### Scenario: P5 verification commands
- **WHEN** each P5 implementation section is completed
- **THEN** the frontend build verification SHALL pass for P5.1, P5.2, and P5.3
- **AND** P5.0 baseline SHALL retain backend test and frontend build verification
