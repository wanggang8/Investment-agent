#!/usr/bin/env python3
"""Validate that registered local HTTP routes are documented."""

from __future__ import annotations

import re
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
REGISTERED_ROUTE_FILES = [
    ROOT / "cmd/server/main.go",
    ROOT / "internal/application/handler/app.go",
]
DOCUMENTED_ROUTE_FILES = [
    ROOT / "docs/api.md",
    ROOT / "docs/frontend-contract.md",
]
ROUTE_RE = re.compile(r'HandleFunc\("([A-Z]+) ([^"]+)"')
DOC_ROUTE_RE = re.compile(r"`(GET|POST|PUT|PATCH|DELETE) (/api/v1[^`\s]*)`")


def normalize(route: str) -> str:
    method, path = route.split(" ", 1)
    return f"{method} {path.split('?', 1)[0]}"


def registered_routes() -> set[str]:
    routes: set[str] = set()
    for path in REGISTERED_ROUTE_FILES:
        text = path.read_text(encoding="utf-8")
        for method, route_path in ROUTE_RE.findall(text):
            routes.add(normalize(f"{method} {route_path}"))
    return routes


def documented_routes() -> set[str]:
    routes: set[str] = set()
    for path in DOCUMENTED_ROUTE_FILES:
        text = path.read_text(encoding="utf-8")
        for method, route_path in DOC_ROUTE_RE.findall(text):
            routes.add(normalize(f"{method} {route_path}"))
    return routes


def main() -> int:
    registered = registered_routes()
    documented = documented_routes()
    missing_docs = sorted(registered - documented)
    stale_docs = sorted(documented - registered)

    if missing_docs or stale_docs:
        print("api_route_contract_check:status=failed")
        if missing_docs:
            print("registered_routes_missing_docs:")
            for route in missing_docs:
                print(f"- {route}")
        if stale_docs:
            print("documented_routes_not_registered:")
            for route in stale_docs:
                print(f"- {route}")
        return 1

    print(f"api_route_contract_check:status=passed routes={len(registered)}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
