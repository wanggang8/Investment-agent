# Investment Agent Documentation

> Updated: 2026-06-23

This is the concise documentation map. Historical phase status and release caveats formerly kept here moved to [release/history.md](release/history.md).

## Start Here

| Document | Use it for |
| --- | --- |
| [product-overview.md](product-overview.md) | Product concepts, user workflows, and safety boundaries. |
| [quickstart.md](quickstart.md) | Docker Compose setup, local configuration, operations, upgrade, uninstall, and troubleshooting. |
| [deployment.md](deployment.md) | Release deployment package behavior and GitHub Actions release gates. |
| [release/README.md](release/README.md) | Current release materials, acceptance records, and caveats. |
| [release/history.md](release/history.md) | Long phase history moved from the former docs README. |

## Contract Truth Sources

| Level | Documents | Notes |
| --- | --- | --- |
| L1 contracts | [requirements.md](requirements.md), [data-model.md](data-model.md), [api.md](api.md), [workflow.md](workflow.md), [frontend-contract.md](frontend-contract.md) | Behavior, interfaces, state, and frontend contract truth sources. Change through OpenSpec deltas. |
| L2 architecture and plans | [architecture.md](architecture.md), [functional-spec.md](functional-spec.md), [development-plan.md](development-plan.md) | Architecture, feature breakdown, and implementation plan history. P95 owns detailed architecture corrections. |
| L3 UI and diagrams | [ui-design.md](ui-design.md), [ui-flow.md](ui-flow.md), [ui/prototype.md](ui/prototype.md), [diagrams/](diagrams/) | Product UI, page flows, and reusable diagram assets. |

## Operations And Validation

| Document | Use it for |
| --- | --- |
| [configuration.md](configuration.md) | Runtime config, local secrets, and safety defaults. |
| [ops-local-scheduler.md](ops-local-scheduler.md) | Explicit local scheduling notes; default operation is not automatic. |
| [testing-plan.md](testing-plan.md) | Test and acceptance strategy. |
| [project-acceptance-gate-matrix.md](project-acceptance-gate-matrix.md) | Release gate matrix. |
| [migration-plan.md](migration-plan.md) | Migration and recovery planning. |

## Governance

| Document | Use it for |
| --- | --- |
| [GOVERNANCE.md](GOVERNANCE.md) | Documentation levels, OpenSpec workflow, archive rules, and forbidden truth-source patterns. |
| [../openspec/project.md](../openspec/project.md) | OpenSpec project map and change governance. |
| [../openspec/PROGRESS.md](../openspec/PROGRESS.md) | Machine-readable current progress. |

## Safety Boundary

Unless a future approved change explicitly alters the contract, these docs must not claim broker connectivity, automatic trading, one-click trading, delegated order placement, external push, automatic confirmation, automatic rule application, return guarantees, paid/login/authorization-only sources, Level2 data, or high-frequency data.
