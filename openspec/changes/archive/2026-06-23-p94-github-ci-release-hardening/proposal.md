# P94 GitHub CI Release Hardening

## Why

Before publishing the repository to GitHub, the project needs CI/CD that behaves like the reference `wanggang8/sub2api` repository: every main/PR change must run backend and frontend quality gates, and version tags must create release artifacts through GitHub Actions.

P91 added the first GitHub workflow layer. P94 hardens it so the public repository does not accept code that fails Go vet, Go lint, frontend lint, tests, OpenSpec, release packaging, or security scanning.

## What Changes

- Harden `.github/workflows/ci.yml` for PR/main pushes with OpenSpec, `go vet`, `golangci-lint`, Go tests, frontend lint/test/build, P91/P92/P93 checks, package smoke, and whitespace checks.
- Harden `.github/workflows/release.yml` so tag pushes `v*` and manual dispatch run the same release preflight before building/uploading release packages.
- Add `.github/workflows/security-scan.yml` with `govulncheck`, frontend production dependency audit, and P93 code reality / secret scan.
- Make existing frontend lint and Go lint pass by cleaning rule scope, hook dependency warnings, unused Go helpers, and staticcheck findings.
- Document GitHub Actions trigger behavior in deployment docs.

## Out Of Scope

- Creating a Git tag or publishing a GitHub Release during P94.
- Physical second-machine validation.
- Broker integration, trading, one-click trading, delegated order placement, external push, automatic confirmation, automatic rule application, return guarantees, paid/login/authorization-only sources, Level2 data, or high-frequency data.
- Managed cloud deployment or hosted database provisioning.

## Acceptance

P94 is acceptable only if:

- Local equivalents of the CI gates pass: OpenSpec, `go vet`, `golangci-lint`, Go tests, frontend lint/test/build, P91/P92/P93 checks, package smoke, and dependency audit.
- GitHub Actions trigger on PR/main pushes and tag pushes as documented.
- Release artifacts are still built without runtime LLM secrets.
- Security scanning does not require broker, trading, paid/login, Level2, or high-frequency sources.
