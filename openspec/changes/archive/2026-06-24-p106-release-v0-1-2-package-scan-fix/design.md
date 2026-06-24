# Design

## Root Cause

`scripts/local-release-package.sh` scans tracked release files for prompt-like payloads using a conservative pattern. The Data Quality page passed a caller-specific redaction label with the key `prompt` and a value long enough to match the scanner. The value was already redacted, but the source shape still resembled a prompt payload.

## Fix

Keep the scanner strict and change the UI redaction label to a short neutral token that does not resemble a prompt payload. This addresses the source of the packaging failure without weakening release safety.

## Version

Use `v0.1.2` because `v0.1.1` was already pushed and failed remote workflows. `v0.1.2` is a patch release for release-package scan compatibility and does not add runtime investment capability.
