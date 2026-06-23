# P97 Design

`config.Load("")` should resolve default paths in this order:

1. `INVESTMENT_AGENT_CONFIG`, when set.
2. `configs/config.yaml`, when it exists.
3. `configs/config.example.yaml`, as a fresh-checkout fallback.

This preserves existing Docker behavior because Docker already sets `INVESTMENT_AGENT_CONFIG=/app/configs/config.docker.yaml`.

`configs/config.yaml` must be ignored by Git. Users create it from `configs/config.example.yaml` or from their own local template and may put local LLM keys in it if they choose. The example config remains tracked and key-free.
