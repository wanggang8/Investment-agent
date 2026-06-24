# P107 README Bilingual Product Intro Acceptance

> Date: 2026-06-24  
> Change: `p107-readme-bilingual-product-intro`  
> Scope: public README documentation only

## Result

Status: `passed`

P107 adds an explicit language switch to the root README and adds `README.zh-CN.md` as a Simplified Chinese public README with product positioning, feature modules, product flow, safety boundary, architecture summary, quickstart, documentation links, CI/release notes, and governance notes.

## Files Changed

- `README.md`
- `README.zh-CN.md`
- `docs/release/acceptance/2026-06-24-p107-readme-bilingual-product-intro.md`
- `openspec/changes/p107-readme-bilingual-product-intro/`

## Claim Boundary

P107 does not modify runtime code, CI workflows, API contracts, SQLite schema, Eino workflow, frontend routes, Docker behavior, installer behavior, or release package artifacts.

P107 does not claim broker connectivity, automatic trading, one-click trading, delegated orders, external push, automatic confirmation, automatic rule application, return guarantees, paid/login/authorized data sources, Level2 data, high-frequency data, Docker validation, installer validation, physical second-machine validation, or a new GitHub Release.

## Validation

- `openspec validate p107-readme-bilingual-product-intro --strict` passed.
- README local-link validation for `README.md` and `README.zh-CN.md` passed.
- `openspec validate --all --strict` passed.
- `git diff --check` passed.
