## ADDED Requirements

### Requirement: Full real product acceptance true pass

After P70, if the project claims full real product acceptance instead of the limited current-data release scope, it SHALL execute a strict P71 acceptance run that treats current-data gate failure, VecLite retrieval degradation, real UI failure, real LLM failure, safety failure, and package verification failure as blockers.

#### Scenario: Current data must pass without scope exclusion

- **GIVEN** P71 is evaluating full real product acceptance
- **WHEN** the current data-source quality strict gate is executed for `000300`
- **THEN** the command MUST return `policy=passed` and `gate=pass`
- **AND** P67 `resolved_with_scope_exclusion`, fixture-only regression, or waiver documentation MUST NOT be accepted as a P71 current-data pass.

#### Scenario: VecLite degradation blocks full acceptance

- **GIVEN** P71 is running real UI consultation or retrieval acceptance
- **WHEN** VecLite/RAG index health is missing, corrupted, incompatible, stale, empty when required, or the workflow reports `VECTOR_INDEX_UNAVAILABLE`
- **THEN** P71 SHALL be marked blocked for retrieval index readiness
- **AND** the acceptance record SHALL NOT describe retrieval-enhanced context as passed.

#### Scenario: Full UI acceptance uses real local operation

- **GIVEN** P71 runs the frontend acceptance suite
- **WHEN** primary product routes and key actions are checked
- **THEN** the browser MUST operate against a real local Go backend and Vite frontend
- **AND** frontend mocks, mocked network responses, or fixture-only current data MUST NOT be used as pass evidence for real product acceptance
- **AND** unexpected API failures, console errors, page errors, mobile/desktop overflow, missing generated decision detail, missing LLM material, or forbidden trading/automation affordances MUST block the P71 pass.

#### Scenario: Post-P70 package refresh follows strict acceptance

- **GIVEN** strict P71 current-data, VecLite, real UI, real LLM, safety, and redaction gates pass
- **WHEN** final distribution package evidence is generated
- **THEN** the package SHALL be built from the accepted post-P70 commit
- **AND** package verify and repeat acceptance SHALL pass
- **AND** the package manifest SHALL state that P69, P70, and P71 acceptance materials are included only if they are present in the packaged source commit.
