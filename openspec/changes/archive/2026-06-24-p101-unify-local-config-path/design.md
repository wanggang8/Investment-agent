# Design

## Config Path Rule

The maintained local runtime config path is:

```text
configs/config.yaml
```

This file is already ignored by Git and is the default used by `config.Load("")` in the Go runtime. P101 aligns older shell acceptance scripts with that default.

## Script Behavior

Each script keeps its explicit override variable but changes the default fallback:

```bash
LOCAL_CONFIG="${P71_LOCAL_CONFIG:-$ROOT_DIR/configs/config.yaml}"
```

This pattern allows a one-off private config path when needed while making the ordinary path consistent.

## Documentation

Current docs and acceptance records created or updated after P101 should refer to `configs/config.yaml` as the local ignored config. Historical archives remain unchanged when they describe historical runs.

## OpenAI-Compatible Request Compatibility

The local analyst client still calls the configured API root with:

```text
POST <base_url>/chat/completions
```

and still sends the same OpenAI-compatible `model` and `messages` body. P101 adds only request headers expected by the working sibling `ai-agent`/SDK path:

- `Accept: application/json`
- `User-Agent: investment-agent/llm-openai-compatible`

P101 also retries one transport timeout before reporting failure. It does not retry unsafe or malformed model output, does not loosen the local quality gate, and does not convert LLM output into a trading decision.

## Timeout

The default/example `deepseek.timeout_seconds` is 60 seconds because the configured OpenAI-compatible gateway can take longer than the old 15-second local default to return response headers.

## Validation

P101 validates three layers:

1. Static path check: maintained scripts no longer default to `configs/config.local.yaml`.
2. LLM request compatibility: focused tests cover the JSON accept/user-agent headers and one-time timeout retry.
3. Machine gates: OpenSpec, Go/Frontend as appropriate, P92/P93, and whitespace checks.
4. Real LLM reruns: P71/P72/P86 should use the configured `configs/config.yaml` by default. If the external provider rejects the model, rate limits, times out after bounded retry, or returns an incompatible response, P101 must classify that honestly rather than hiding it.
