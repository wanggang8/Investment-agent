# P94 Tasks

## 1. Reference Inspection

- [x] Inspect `wanggang8/sub2api` GitHub workflows for CI, release, and security scan structure.
- [x] Compare the current Investment Agent workflow coverage against the reference pattern.

## 2. CI Hardening

- [x] Add `go vet` to PR/main CI.
- [x] Add `golangci-lint` to PR/main CI.
- [x] Add frontend lint to PR/main CI.
- [x] Keep OpenSpec, Go tests, frontend tests/build, deployment checks, release package smoke, and whitespace checks in PR/main CI.
- [x] Add P92/P93 audit checks to PR/main CI.

## 3. Release And Security Automation

- [x] Harden tag/manual release workflow with release preflight checks before packaging.
- [x] Add independent security scan workflow with `govulncheck`, frontend production dependency audit, and P93 code reality / secret scan.

## 4. Lint Closure

- [x] Make frontend lint pass without suppressing product-critical checks globally.
- [x] Make Go lint pass by removing unused helpers and staticcheck findings.

## 5. Documentation And Acceptance

- [x] Document GitHub Actions trigger behavior in deployment docs.
- [x] Generate P94 acceptance record.
- [x] Run local equivalents of the CI/release/security gates.
- [x] Archive P94 after validation passes.
