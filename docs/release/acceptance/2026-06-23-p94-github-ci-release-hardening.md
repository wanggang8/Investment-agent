# P94 GitHub CI Release Hardening Acceptance

> Date: 2026-06-23
> Conclusion: `release_ready_for_github_public_push_with_ci_release_hardening`

## Scope

P94 hardens the GitHub-facing CI/CD layer before the first public push. It does not add investment runtime capabilities.

## Implemented Gates

- PR/main CI: OpenSpec, `go vet`, bounded `golangci-lint`, Go tests, frontend lint, frontend tests, frontend build, P91 deployment checks, P92/P93 audit checks, release package smoke, whitespace check.
- Tag/manual release: release preflight, package build, package verify, artifact upload, GitHub release attachment for `v*` tags.
- Security scan: `govulncheck`, frontend production dependency audit, P93 code reality / secret scan.

## Local Validation

The following local equivalents were executed during P94:

- `npm --prefix web run lint`: passed.
- `go vet ./...`: passed.
- `golangci-lint run --timeout=5m --enable-only=govet,ineffassign,staticcheck,unused ./...`: passed with `0 issues`.
- `openspec validate --all --strict`: passed.
- `go test ./...`: passed.
- `npm --prefix web test`: passed, 48 files / 176 tests.
- `npm --prefix web run build`: passed.
- `python3 scripts/p92_final_requirement_audit.py --check`: passed.
- `python3 scripts/p93_code_reality_audit.py --check`: passed.
- `npm --prefix web audit --omit=dev --audit-level=high`: passed, 0 vulnerabilities.
- `bash scripts/local-release-package.sh --release-label ci-smoke --output-dir tmp/ci-release-check` and verify: passed.

## Boundary

P94 does not create a Git tag, does not publish a GitHub release, does not claim physical second-machine validation, and does not add broker connectivity, automatic trading, one-click trading, delegated order placement, external push, automatic confirmation, automatic rule application, return guarantees, paid/login/authorization-only sources, Level2 data, or high-frequency data.
