## ADDED Requirements

### Requirement: P86 core goal knowledge safety final closure

After P81-P85 and P87, P86 SHALL reconcile the remaining core-goal, source/data transition, knowledge/LLM/RAG, expected-return, review/audit, implementation, release-safety, and unclassified rows into a final row-level matrix and SHALL only claim full original-requirement pass if every full-release-required row is resolved by valid evidence or explicitly reclassified with documented rationale.

#### Scenario: P86 row inventory completes the post-P87 remainder

- **GIVEN** P86 starts from the P87 evidence matrix
- **WHEN** the P86 plan is reviewed
- **THEN** P86 SHALL cover exactly the 137 remaining full-release-required non-`real_pass` rows from P87
- **AND** no P87 remaining row SHALL be omitted from the final P86 inventory.

#### Scenario: Integrated real user acceptance

- **GIVEN** P86 evaluates the product goal
- **WHEN** end-to-end acceptance runs
- **THEN** it SHALL use real local UI operation, API responses, workflow metadata, read-only SQLite evidence, and deterministic checks where applicable
- **AND** it SHALL cover setup, portfolio/account state, data readiness, knowledge/RAG, consultation, expected return, risk/SOP, manual confirmation, review, audit, release governance, and safety.

#### Scenario: Row upgrade requires direct evidence

- **GIVEN** P86 generates the final matrix
- **WHEN** a row is upgraded to `real_pass`
- **THEN** the matrix SHALL cite direct row-level evidence from P86 or cumulative P81-P87 artifacts
- **AND** it SHALL NOT rely only on seeded SQLite rows, route smoke, screenshots, fixture/mock/stub data, or broad narrative.

#### Scenario: Full-pass claim is evidence gated

- **GIVEN** P86 writes final release materials
- **WHEN** any full-release-required row remains partial, blocked, scoped-only, unsupported, or unverified
- **THEN** P86 SHALL NOT claim full original-requirement pass
- **AND** it SHALL list the exact remaining rows and blockers.

#### Scenario: Forbidden capabilities remain out of product scope

- **GIVEN** P86 passes integrated acceptance
- **WHEN** final claims are written
- **THEN** they SHALL NOT introduce or imply broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, real DB overwrite, paid/login/authorized source, Level2 source, high-frequency source, future provider availability, or return promises.
