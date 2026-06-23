## ADDED Requirements

### Requirement: P76 SHALL refresh package evidence after P75

P76 SHALL regenerate final local package evidence after P75 so package freshness claims are not based on the earlier P71 archive.

#### Scenario: Package source is the clean post-P75 package commit

- **GIVEN** P75 has been committed and archived
- **WHEN** P76 generates package evidence
- **THEN** the package SHALL be generated from a clean source commit that includes committed P72-P75 evidence and any P76 acceptance-harness correction required to make repeat acceptance deterministic
- **AND** the package manifest SHALL record `source_status=clean`
- **AND** the package manifest SHALL record the package source commit
- **AND** release materials SHALL state whether P72-P75 acceptance Markdown and OpenSpec archives are included in the packaged source.

#### Scenario: Package verify and repeat acceptance pass

- **WHEN** the P76 package archive is generated
- **THEN** package verification SHALL confirm archive checksum consistency, required entries, forbidden path exclusions, and manifest safety boundaries
- **AND** repeat acceptance SHALL run from an extracted package workspace rather than from the active repository checkout
- **AND** repeat acceptance SHALL cover OpenSpec validation, Go tests, frontend dependency installation, frontend tests, frontend build, and local E2E smoke.

#### Scenario: Package handoff remains bounded

- **WHEN** P76 updates release materials
- **THEN** the materials SHALL include package identity, source commit, source status, checksum, archive entry count, verify result, repeat result, known caveats, and Not Claimed boundaries
- **AND** the materials SHALL NOT claim that the archive includes P76 package-after-the-fact evidence or `docs/release/ui-audit-assets/`
- **AND** the materials SHALL preserve P75 `release_ready_scoped_with_traceability_gaps`
- **AND** the materials SHALL NOT claim physical second-machine execution, remote publishing, Git tag creation, installer signing, automatic upgrade, automatic migration, automatic restore, automatic repair, real database overwrite, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, future provider availability, or investment returns.
