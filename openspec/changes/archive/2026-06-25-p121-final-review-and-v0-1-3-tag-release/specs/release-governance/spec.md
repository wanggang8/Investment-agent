## ADDED Requirements

### Requirement: P121 final review gates v0.1.3 tag publication

The repository SHALL publish `v0.1.3` only after a fresh P121 final review validates the current source tree and release-facing materials.

#### Scenario: Fresh release review passes before tagging

- **GIVEN** P114-P120 have been archived and the current source version is being advanced to `v0.1.3`
- **WHEN** P121 release verification is executed
- **THEN** OpenSpec strict validation, Go tests, Go vet, frontend tests, frontend build, P92 audit check, P121 final release review, whitespace checks, and local release package smoke/verify SHALL pass before the tag is created.

#### Scenario: P93 stale boundary stays explicit

- **GIVEN** P93 was completed before the P114-P120 post-redesign acceptance layer
- **WHEN** `v0.1.3` release notes describe final review status
- **THEN** they SHALL NOT claim fresh P93 pass after P114-P120
- **AND** they SHALL record that P121 is the fresh release-governance review for the current tree.

#### Scenario: v0.1.3 release claims stay bounded

- **GIVEN** `v0.1.3` is described in release materials and tag content
- **WHEN** release readiness is communicated
- **THEN** the project MAY claim the scoped P114-P120 product/UI/real-use/control acceptance layer and P121 release gates
- **AND** it SHALL NOT claim new Docker installation validation, upgrade validation, physical second-machine validation, broker connectivity, trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, future provider availability, prediction accuracy, or investment returns unless separately validated.

#### Scenario: Release package does not expose local workstation paths

- **GIVEN** historical design and acceptance materials may mention local generated-image or workspace paths in source
- **WHEN** the local release package script copies text files into the archive
- **THEN** local absolute paths SHALL be replaced with stable placeholders before forbidden-content scanning and package inclusion
- **AND** the package smoke and verify gates SHALL fail if any forbidden local path remains in included text files.
