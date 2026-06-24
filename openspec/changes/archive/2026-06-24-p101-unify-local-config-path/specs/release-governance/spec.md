## ADDED Requirements

### Requirement: P101 unified local config path

The project SHALL use `configs/config.yaml` as the default ignored local config path for current local runtime and current local-source acceptance scripts.

#### Scenario: Historical acceptance scripts use the runtime default

- **GIVEN** a user configures LLM and local runtime settings in `configs/config.yaml`
- **WHEN** current local-source acceptance scripts are run without explicit override variables
- **THEN** they SHALL read `configs/config.yaml` by default
- **AND** they SHALL NOT require a separate `configs/config.local.yaml` file.

#### Scenario: Explicit overrides remain available

- **GIVEN** an operator needs a one-off private config file
- **WHEN** a script-specific variable such as `P71_LOCAL_CONFIG`, `P72_LOCAL_CONFIG`, `P75_LOCAL_CONFIG`, or `P63_LOCAL_CONFIG` is provided
- **THEN** that script SHALL use the explicit path
- **AND** this override SHALL NOT change the default documented local config path.

### Requirement: P101 OpenAI-compatible local LLM request compatibility

The local analyst LLM client SHALL remain compatible with OpenAI Chat Completions gateways that expect JSON accept headers, stable user-agent identification, and longer bounded response times.

#### Scenario: Compatible headers are sent

- **GIVEN** a configured OpenAI-compatible LLM gateway
- **WHEN** the analyst client sends a chat completion request
- **THEN** it SHALL send `Accept: application/json`
- **AND** it SHALL send a stable `User-Agent`
- **AND** it SHALL continue using the configured `<base_url>/chat/completions` path.

#### Scenario: Transport timeout is retried once

- **GIVEN** the first LLM request times out before receiving response headers
- **WHEN** the retry succeeds
- **THEN** the analyst client SHALL return parsed analysis material
- **AND** it SHALL mark metadata with a bounded timeout retry
- **AND** it SHALL NOT loosen the local parser or quality gate.

#### Scenario: Default timeout allows slower compatible gateways

- **GIVEN** a local config omits `deepseek.timeout_seconds`
- **WHEN** defaults are applied
- **THEN** the configured timeout SHALL be 60 seconds.

#### Scenario: Release claims stay bounded

- **GIVEN** P101 changes script defaults and LLM request compatibility
- **WHEN** release readiness is described
- **THEN** the project MAY claim local config path consistency for source-runtime validation
- **AND** it MAY claim OpenAI-compatible LLM request compatibility for Chat Completions style gateways
- **AND** it SHALL NOT claim Docker installation, package distribution, GitHub Release, physical second-machine validation, broker connectivity, trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, paid/login/auth-only sources, Level2 data, high-frequency data, future provider availability, or investment returns.
