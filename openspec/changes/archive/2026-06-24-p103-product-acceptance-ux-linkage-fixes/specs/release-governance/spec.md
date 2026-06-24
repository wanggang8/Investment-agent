## ADDED Requirements

### Requirement: P103 product acceptance UX linkage fixes

The product SHALL address P102 non-blocking UX findings without expanding investment runtime capabilities.

#### Scenario: Portfolio empty state is onboarding-safe

- **GIVEN** no local portfolio snapshot exists
- **WHEN** the user opens the portfolio page
- **THEN** the page SHALL present first-use onboarding and local account calibration guidance instead of a generic system failure.

#### Scenario: Decision analysis remains auditable without overwhelming the page

- **GIVEN** a decision contains real LLM analyst reports
- **WHEN** the user opens the decision detail page
- **THEN** the page SHALL show the final verdict and safety boundary first
- **AND** the full analysis material SHALL remain available through explicit expansion.

#### Scenario: Decision loop deep link focuses the target

- **GIVEN** a decision-loop URL includes `decision_id`
- **WHEN** the linked decision is present in the loop response
- **THEN** the page SHALL focus that decision's loop record and keep trace links read-only.

#### Scenario: Release claims remain bounded

- **GIVEN** P103 fixes P102 UX findings
- **WHEN** release readiness is described
- **THEN** the project MAY claim the checked local product UX issues were fixed
- **AND** it SHALL NOT claim Docker installation, package distribution, GitHub Release, physical second-machine validation, broker connectivity, trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, paid/login/auth-only sources, Level2 data, high-frequency data, future provider availability, or investment returns.
