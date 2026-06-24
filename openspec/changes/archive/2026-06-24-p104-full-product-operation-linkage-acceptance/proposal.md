# P104 Full Product Operation Linkage Acceptance

## Why

P102/P103 proved the primary local product journey with a real model and fixed the most visible non-blocking UX linkage issues. The remaining acceptance risk is repeatability: a human can inspect the product, but there is not yet a single fresh gate that maps broad product operations to API responses, SQLite side effects, downstream readback, audit trails, and forbidden automation absence.

## What

- Add a P104 product operation/linkage acceptance matrix that enumerates the local product surfaces and the evidence expected for each operation class.
- Add a repeatable local-source runner that starts an isolated backend with a temporary SQLite database, performs representative write/read operations through HTTP APIs, checks SQLite side effects, and records safety negative evidence.
- Add a P104 acceptance record summarizing the fresh run and its boundary.
- Update governance/progress materials and archive the change when validation passes.

## Scope

In scope:

- Local source acceptance only.
- Product operation logic, associated data, downstream readback, audit traceability, and safety boundary verification.
- Temporary SQLite database and local backend process started by the P104 runner.

Out of scope:

- Docker, installer, GitHub Release, package refresh, physical second-machine validation, or remote deployment.
- New investment runtime capabilities.
- Broker interfaces, automatic trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, or return guarantees.
