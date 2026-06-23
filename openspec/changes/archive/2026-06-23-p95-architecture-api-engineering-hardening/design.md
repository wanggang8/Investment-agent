# P95 Design

P95 is an engineering hardening stage, not a product capability stage.

The Go package issue is handled by adding a small script that derives backend packages from `go list ./...` and rejects any package under ignored frontend dependency trees. CI and release workflows should call that script instead of raw `go test ./...` where package discovery matters.

P93 should become independent of local ignored runtime state. Its production-file inventory should use Git tracked files plus nonignored untracked files so new source files are scanned during local development while ignored `tmp`, VecLite, or generated frontend output cannot change the report.

The API route contract should stay lightweight: parse `internal/application/handler/app.go` and `cmd/server/main.go`, normalize documented routes from `docs/api.md` and `docs/frontend-contract.md`, and fail when an implemented route is undocumented or when docs reference a nonexistent route. Query strings in examples are ignored for route identity.

SQLite hardening stays local and conservative. `Open` should enable foreign key enforcement, configure a bounded busy timeout, and attempt WAL mode for file-backed databases. `:memory:` tests should remain compatible.

Docker secrets should support `DEEPSEEK_API_KEY_FILE` so Compose secrets can be used without removing the simpler `.env` path. The runtime config loader remains responsible for the final key value.

P95 may update `docs/architecture.md`, `docs/deployment.md`, `docs/api.md`, `docs/frontend-contract.md`, scripts, CI workflows, and focused tests. It must not rewrite root README or the public documentation information architecture; that belongs to P96.
