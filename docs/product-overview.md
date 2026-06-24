# Product Overview

> Updated: 2026-06-23

Investment Agent is a local-first investment discipline product. It helps a user keep portfolio facts, evidence, analysis, decisions, reviews, and governance records in one auditable workspace. It is designed for manual decision-making: the app can organize and analyze information, but the user remains responsible for any real-world action.

## Core Workflows

### 1. Daily Discipline

The daily cockpit surfaces current status, pending manual actions, data-quality signals, risk alerts, notifications, and review work. It is meant to answer: what needs attention today, what evidence changed, and what action is still waiting for a person.

### 2. Portfolio Maintenance

The portfolio pages record local account calibration, holdings, cash and money-fund buckets, buy dates, position states, manual offline trades, and corrections. These records feed downstream risk checks, decision context, audit events, and review pages.

### 3. Evidence And Data Readiness

The app keeps supported public evidence, source-health records, market snapshots, local knowledge, and retrieval indexes close to the decision workflow. Readiness and data-quality views make degraded or missing evidence visible instead of silently treating it as complete.

### 4. Consultation And Decision Explanation

Consultation combines deterministic local rules, portfolio facts, evidence references, expected-return analysis, retrieval context, and LLM analyst reports when a provider key is configured. The output is an explanation and recommendation context, not an order instruction. Decision detail pages preserve assumptions, analyst results, final verdict structure, evidence links, rule traces, and audit references.

### 5. Manual Confirmation

Important actions require explicit user confirmation. Confirmation records the user's choice and the data impact in local SQLite. The system does not submit trades, create broker orders, or auto-confirm decisions.

### 6. Review, Audit, And Rule Governance

Review and governance pages support monthly or quarterly review, error-case marking, rule proposals, gatekeeper checks, notification state, audit trails, and release traceability. Rule proposals remain governed records; the product does not automatically apply new rules without explicit handling.

## Local Data Model In Plain Terms

Local state is persisted in SQLite and local VecLite/sqlite-vec index files. The app stores decisions, portfolio facts, evidence metadata, source-health state, audit events, notifications, rule records, review records, and retrieval material needed for the local workflow. SQLite remains the authoritative fact store; the vector index can be rebuilt from SQLite `rag_chunks`. The L1 data contract remains [data-model.md](data-model.md).

## Safety Boundary

Investment Agent does not provide:

- broker connectivity;
- automatic trading;
- one-click trading;
- delegated order placement;
- external push delivery;
- automatic decision confirmation;
- automatic rule application;
- return guarantees;
- paid, login-only, or authorization-only data-source claims;
- Level2 or high-frequency data.

The product may show analysis, risks, gates, assumptions, and local records. Real-world investment decisions remain outside the software boundary.

## Where To Go Next

- Run locally: [quickstart.md](quickstart.md)
- Understand architecture: [architecture.md](architecture.md)
- Read product requirements: [requirements.md](requirements.md)
- Inspect release evidence: [release/README.md](release/README.md)
