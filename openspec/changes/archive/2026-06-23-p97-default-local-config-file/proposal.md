# P97 Default Local Config File

## Why

Local runtime startup currently falls back directly to `configs/config.example.yaml` when no explicit config path is supplied. That makes the example file behave like the default runtime file, which is confusing and encourages editing a tracked template for local secrets.

## What Changes

- Make `configs/config.yaml` the preferred default local runtime config file.
- Keep `configs/config.example.yaml` as a committed template and fallback for fresh checkout compatibility.
- Ignore `configs/config.yaml` in Git so real local keys and machine-specific paths are not committed.
- Update configuration documentation and startup guidance.

## Out Of Scope

- No investment runtime capability changes.
- No broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, Level2/high-frequency/paid-login sources, or return guarantees.
