## ADDED Requirements

### Requirement: P73 Product Effectiveness And UX Validation

P73 SHALL validate whether the product supports its stated investment-discipline assistant goal, not merely whether features execute.

#### Scenario: Product goal metrics are recorded

- **GIVEN** P73 evaluates product effectiveness
- **WHEN** acceptance materials are written
- **THEN** they SHALL record discipline adherence, evidence sufficiency, traceability, review usefulness, and UX comprehension checks
- **AND** they SHALL NOT use future investment returns as the required pass metric.

#### Scenario: Real UX task validation is required

- **GIVEN** a user operates the local product
- **WHEN** P73 runs browser acceptance
- **THEN** it SHALL cover first-use or missing-prerequisite guidance, daily discipline, portfolio maintenance, data quality/evidence review, consultation, decision detail, manual confirmation, risk/notification/audit/rules/review readback, and invalid or unsafe input
- **AND** page errors, unexpected API failures, console errors, forbidden affordances, or critical UX confusion SHALL block pass.

#### Scenario: Effect replay validates discipline behavior

- **GIVEN** local facts, evidence records, decisions, confirmations, risk alerts, rule-effect facts, and audit events exist
- **WHEN** P73 runs effect replay checks
- **THEN** C-level background-only material SHALL NOT satisfy formal evidence
- **AND** insufficient evidence SHALL result in safe blocking, gap qualification, or non-trade records
- **AND** manual confirmation SHALL be the only accepted path that mutates local portfolio facts
- **AND** rule proposal/effect validation SHALL expose sample, overfit, gate, or tracking state when available.

#### Scenario: UX audit findings are dispositioned

- **GIVEN** P73 captures representative UI screenshots and task results
- **WHEN** the UX audit is written
- **THEN** findings SHALL be classified as critical, major, minor, or accepted gap
- **AND** critical findings SHALL block product-effectiveness pass until fixed and rerun.

#### Scenario: P73 claims remain bounded

- **GIVEN** P73 passes
- **WHEN** release materials state the result
- **THEN** they MAY claim product-effectiveness and UX validation for the accepted local scope
- **AND** they SHALL NOT claim investment return improvement, future market prediction, future public-source or model-provider availability, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, or automatic real DB overwrite.
