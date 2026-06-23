# P94 Design

P94 keeps the CI/CD design simple and GitHub-native:

- `ci.yml` is the required PR/main quality gate. It validates OpenSpec, backend quality, frontend quality, deployment checks, release package smoke, and generated-diff hygiene.
- `release.yml` is tag/manual packaging. Tags matching `v*` create GitHub release artifacts after the same preflight used by CI.
- `security-scan.yml` is independent so dependency and vulnerability failures are clearly separated from build/test failures.

The Go lint gate uses a bounded `golangci-lint` linter set: `govet`, `ineffassign`, `staticcheck`, and `unused`. This makes dead code and obvious correctness issues blocking while avoiding style-only churn before the first public push.

The frontend lint gate keeps E2E dynamic response typing flexible but preserves source linting, React hooks, React refresh, TypeScript, and build checks.

No workflow receives runtime LLM credentials from the repository. Release packages continue to rely on `.env.example` and local runtime `.env` configuration.
