## ADDED Requirements

### Requirement: P106 release package scan compatibility for v0.1.2

The repository SHALL keep release package prompt-payload scanning strict while avoiding source-level false positives in redacted UI labels.

#### Scenario: Package scanner passes redacted UI labels

- **GIVEN** the release package script scans tracked source files
- **WHEN** frontend redaction labels are inspected
- **THEN** caller-specific replacement labels SHALL NOT use long JSON-like `prompt: "..."` payload shapes that match prompt-payload forbidden-content rules
- **AND** the release package smoke and verify steps SHALL pass before `v0.1.2` is tagged.

#### Scenario: v0.1.2 patch release stays bounded

- **GIVEN** `v0.1.2` is described in release materials
- **WHEN** release readiness is communicated
- **THEN** the project MAY claim the release-package scan compatibility fix and current source version metadata synchronization
- **AND** it SHALL NOT claim Docker installation validation, physical second-machine validation, broker connectivity, trading, automatic confirmation, automatic rule application, future provider availability, or investment returns unless separately validated.
