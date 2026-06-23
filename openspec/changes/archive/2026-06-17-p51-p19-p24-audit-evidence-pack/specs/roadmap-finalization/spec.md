# roadmap-finalization Specification Delta

## ADDED Requirements

### Requirement: Maintain a current-state P19-P24 audit evidence pack

The project SHALL provide a current-state audit evidence pack for P19-P24 before project acceptance gates or release candidate materials are prepared.

#### Scenario: P19-P24 audit evidence is reviewed

- **GIVEN** P19-P24 are marked delivered but do not have per-phase complete archive packages
- **WHEN** audit evidence is requested for these phases
- **THEN** the project SHALL provide a document that lists each phase's delivery boundary, archive status, documentation evidence, code evidence, test evidence, rerunnable commands, residual gaps, and claims that must not be made
- **AND** the document SHALL state that it is not a replacement for missing historical archive packages.

### Requirement: Preserve historical integrity in P19-P24 audit evidence

The project SHALL NOT fabricate P19-P24 historical archive packages, completion timestamps, or unverified implementation claims when creating audit evidence.

#### Scenario: Missing archive packages are documented

- **WHEN** a P19-P24 phase lacks a complete archive package
- **THEN** the audit evidence SHALL use missing or not present archive language
- **AND** it SHALL cite current repository evidence instead of backfilled historical tasks or fabricated archive paths.

### Requirement: Keep P19-P24 audit evidence separate from acceptance gates

The P19-P24 audit evidence pack SHALL remain separate from the project acceptance gate matrix.

#### Scenario: Acceptance planning follows the audit evidence pack

- **GIVEN** the P19-P24 audit evidence pack has been completed
- **WHEN** the next stage is selected
- **THEN** the next recommended stage SHALL be P52 `p52-project-acceptance-gate-matrix`
- **AND** P52 SHALL convert the evidence into release-blocking and non-blocking validation gates instead of treating P51 alone as release readiness.
