# P93 Final Code Reality And Design Audit

## Why

P92 proves final requirement acceptance from archived evidence, but the user requested a fresh code-facing audit: re-read original requirements, inspect implementation paths, and determine whether the product is truly implemented rather than hardcoded, demo-like, dead-code-backed, or poorly designed.

## What Changes

- Add a code reality audit script and final audit report.
- Cross-check original requirement sections against production Go, React, configuration, scripts, tests, and release artifacts.
- Classify suspicious terms such as `stub`, `mock`, `demo`, `placeholder`, hardcoded values, and temporary code by context.
- Check whether visible product routes are real product pages rather than placeholder/demo pages.
- Validate release-relevant hardening with Go, frontend, OpenSpec, static scans, and package checks.

## Out Of Scope

- No new product runtime features.
- No new investment behavior, providers, LLM behavior, UI routes, API endpoints, SQLite schema, or workflow nodes unless the audit finds a release-blocking defect.
- No physical second-machine validation.
- No broker integration, trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, return guarantees, paid/login/auth-only sources, Level2 data, or high-frequency data.

## Acceptance

P93 is acceptable only if:

- A final code audit report lists implementation evidence by original requirement section.
- The report classifies demo/mock/stub/hardcoded/dead-code risks and identifies whether each is acceptable, test-only, config-only, or release-blocking.
- The audit verifies that production UI routes do not rely on placeholder/demo pages.
- The audit verifies that release config defaults use real data paths (`use_stub=false`) and no secrets are embedded.
- The audit runs `go test ./...`, frontend tests/build, OpenSpec strict validation, the P92 checker, the P93 checker, and `git diff --check`.
- Any blocker found is either fixed in P93 or explicitly recorded as blocking rather than hidden.
