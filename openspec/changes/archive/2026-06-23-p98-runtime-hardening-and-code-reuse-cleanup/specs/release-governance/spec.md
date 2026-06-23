## ADDED Requirements

### Requirement: P98 SHALL harden release runtime mode and frontend redaction reuse

P98 SHALL add release-mode guardrails and shared frontend redaction without changing investment runtime capabilities.

#### Scenario: Release runtime rejects stub data

- **GIVEN** runtime mode is configured as `release`
- **AND** `data_sources.use_stub` is `true`
- **WHEN** configuration validation runs
- **THEN** validation SHALL fail with an actionable message
- **AND** release/Docker defaults SHALL keep `data_sources.use_stub=false`.

#### Scenario: Development fallback remains available

- **GIVEN** runtime mode is omitted or configured as `development`
- **WHEN** local example or test configuration enables stub data
- **THEN** validation SHALL continue to allow local stub data
- **AND** this SHALL NOT create a release claim for real provider operation.

#### Scenario: Frontend diagnostic redaction is shared

- **GIVEN** frontend pages or components display diagnostic or failure text
- **WHEN** the text contains key-shaped tokens, SQL fragments, prompt fragments, raw diagnostic payloads, stack traces, or local paths
- **THEN** the text SHALL be redacted through a shared utility
- **AND** current page/component tests SHALL continue to prove sensitive details are not displayed.
