# P101 Unify Local Config Path

## Why

The local program now defaults to `configs/config.yaml`, but several historical real-UI/real-LLM acceptance scripts still default to `configs/config.local.yaml`. This creates an avoidable split-brain setup: a user can configure the local program correctly while the acceptance scripts still report a missing key.

## What Changes

- Change historical acceptance script defaults from `configs/config.local.yaml` to `configs/config.yaml`.
- Preserve explicit environment-variable overrides such as `P71_LOCAL_CONFIG`, `P72_LOCAL_CONFIG`, `P75_LOCAL_CONFIG`, and `P63_LOCAL_CONFIG`.
- Align the hand-written LLM client with OpenAI-compatible gateway expectations by sending JSON accept/user-agent headers and retrying one transport timeout.
- Raise the default/example LLM timeout to 60 seconds for slower compatible model gateways.
- Update current documentation and acceptance notes so the maintained local config path is `configs/config.yaml`.
- Re-run validation and real LLM-dependent local-source acceptance after the user-provided OpenAI-compatible config is present.

## Scope Boundaries

- Does not commit or print API keys.
- Does not change the LLM endpoint path, body schema, response parser, or analyst quality gate.
- Does not add providers, API endpoints, SQLite schema, frontend routes, or investment runtime capability.
- Does not validate Docker, install/upgrade/uninstall, GitHub Release, package refresh, or physical second-machine runs.
- Does not claim broker connectivity, automatic trading, one-click trading, delegated orders, external push, automatic confirmation, automatic rule application, paid/login/auth/Level2/HFT sources, or return guarantees.
