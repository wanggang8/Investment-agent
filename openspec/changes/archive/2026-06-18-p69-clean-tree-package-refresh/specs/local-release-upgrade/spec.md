## ADDED Requirements

### Requirement: P69 SHALL refresh final package evidence from a clean tree

P69 SHALL regenerate package evidence from a clean committed source tree after P65-P68 so final package claims are not based on stale dirty candidate archives.

#### Scenario: Clean package is generated from the P68 source commit

- **GIVEN** P68 has been committed and the main repository has a clean HEAD
- **WHEN** P69 generates final package evidence
- **THEN** the package SHALL be generated from a clean detached worktree or equivalent clean checkout at the P68 source commit
- **AND** the package manifest SHALL record `source_status=clean`
- **AND** the package manifest SHALL record the source commit used for the package
- **AND** P69 release materials SHALL NOT claim the generated archive includes P69 documentation unless a later package refresh is performed after P69 commit.

#### Scenario: Clean package is verified and repeated

- **WHEN** the P69 package archive is generated
- **THEN** package verification SHALL confirm archive checksum consistency, required entries, forbidden path exclusions, and manifest safety boundaries
- **AND** repeat acceptance SHALL run from an extracted package workspace rather than from the active repository checkout
- **AND** repeat acceptance SHALL cover OpenSpec validation, Go tests, frontend dependency installation, frontend tests, frontend build, and local E2E smoke.

#### Scenario: Package evidence is handed off

- **WHEN** P69 updates release materials
- **THEN** the materials SHALL include package identity, source commit, `source_status`, checksum, archive entry count, verify result, repeat command matrix, known caveats, and Not Claimed boundaries
- **AND** the materials SHALL preserve P68 `release_ready_limited_current_data_scope`
- **AND** the materials SHALL NOT claim P66 current-data policy passed, physical second-machine execution, remote publishing, Git tag creation, installer signing, automatic upgrade, automatic migration, automatic restore, automatic repair, real database overwrite, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, provider availability, or investment returns.
