#!/usr/bin/env python3
"""P90 capital-flow public source preverification."""

from __future__ import annotations

import json
import subprocess
import sys
import time
import urllib.parse
import urllib.request
from datetime import datetime, timezone
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
ARTIFACT_DIR = ROOT / "docs" / "release" / "ui-audit-assets" / "2026-06-22-p90-capital-flow-provider"
ARTIFACT = ARTIFACT_DIR / "p90-source-preverification.json"
TIMEOUT = 20


def fetch_json(url: str, headers: dict[str, str] | None = None) -> dict:
    request_headers = {
        "User-Agent": "Mozilla/5.0 InvestmentAgent-P90 readonly acceptance",
        "Accept": "application/json,text/plain,*/*",
        "Connection": "close",
        **(headers or {}),
    }
    last_error: Exception | None = None
    for attempt in range(3):
        req = urllib.request.Request(url, headers=request_headers)
        try:
            with urllib.request.urlopen(req, timeout=TIMEOUT) as resp:
                raw = resp.read().decode("utf-8", errors="replace")
        except Exception as err:
            last_error = err
            cmd = ["curl", "-L", "--retry", "2", "--retry-delay", "1", "--max-time", str(TIMEOUT), "-s", "-A", request_headers["User-Agent"]]
            for key, value in (headers or {}).items():
                cmd.extend(["-H", f"{key}: {value}"])
            cmd.append(url)
            try:
                raw = subprocess.check_output(cmd, text=True)
            except Exception as curl_err:
                last_error = curl_err
                raw = ""
        try:
            return json.loads(raw)
        except Exception as err:
            last_error = err
            if attempt < 2:
                time.sleep(1 + attempt)
    raise SystemExit(f"status=failed\nreason=source_json_parse_error:{last_error}")


def preverify_capital_flow(symbol: str) -> dict:
    secid = f"1.{symbol}" if symbol.startswith("6") else f"0.{symbol}"
    query = urllib.parse.urlencode({
        "secid": secid,
        "fields1": "f1,f2,f3",
        "fields2": "f51,f52,f53,f54,f55,f56,f62,f63",
        "ut": "b2884a393a59ad64002292a3e90d46a5",
    }, safe=",")
    url = f"https://emdatah5.eastmoney.com/dc/ZJLX/getDBHistoryData?{query}"
    page = f"https://emdatah5.eastmoney.com/dc/zjlx/stock?fc={urllib.parse.quote(secid)}"
    js = "https://emdatah5.eastmoney.com/dc/Content/js/zjlx/stock.min.js"
    payload = fetch_json(url, {"Referer": page})
    rows = (((payload.get("data") or {}).get("klines")) or [])
    if not rows:
        raise SystemExit("status=failed\nreason=capital_flow_no_rows")
    fields = str(rows[-1]).split(",")
    if len(fields) < 2:
        raise SystemExit("status=failed\nreason=capital_flow_fields")
    date = fields[0]
    raw_net_flow = float(fields[1])
    net_inflow = raw_net_flow if raw_net_flow > 0 else 0.0
    net_outflow = abs(raw_net_flow) if raw_net_flow < 0 else 0.0
    if not date:
        raise SystemExit("status=failed\nreason=capital_flow_required_fields")
    return {
        "category": "capital_flow",
        "requirement_ids": ["REQ-04-016", "REQ-05-003"],
        "provider_status": "verified_runtime_public",
        "real_pass_eligible": True,
        "authority": "Eastmoney H5 public capital-flow endpoint",
        "public_access_shape": "HTTPS GET used by public H5 page; browser/mobile-style GET may be required; no login, paid token, authorization, Level2, or broker session observed",
        "stable_request_or_page_evidence": url,
        "reference_page": page,
        "reference_js": js,
        "fields": {
            "date": date,
            "net_inflow": net_inflow,
            "net_outflow": net_outflow,
            "raw_net_flow": raw_net_flow,
        },
        "field_semantics": "H5 history exposes daily net capital flow as f52/index 1; positive maps to net_inflow, negative maps to net_outflow, raw value is preserved as raw_net_flow.",
        "update_frequency": "daily historical capital-flow series; P90 acceptance uses one low-frequency product refresh",
        "legal_access_limits": "public page/API observed without login or token; P90 treats it as low-frequency read-only public market data and not as Level2/high-frequency feed",
        "rate_limit_assumption": "single acceptance request; no polling or high-frequency refresh",
        "failure_behavior": "source_unavailable; do not synthesize values",
        "sqlite_target_path": "market_snapshots.market_metrics_json.metadata.p88_structured_fields.capital_flow",
        "source_name": "eastmoney_h5_zjlx_public",
    }


def validate(payload: dict) -> None:
    item = payload.get("category") or {}
    if item.get("category") != "capital_flow" or item.get("real_pass_eligible") is not True:
        raise SystemExit("status=failed\nreason=capital_flow_not_eligible")
    fields = item.get("fields") or {}
    for key in ("date", "net_inflow", "net_outflow", "raw_net_flow"):
        if key not in fields:
            raise SystemExit(f"status=failed\nreason=missing_field:{key}")
    forbidden = json.dumps({
        "provider_status": item.get("provider_status"),
        "source_name": item.get("source_name"),
        "stable_request_or_page_evidence": item.get("stable_request_or_page_evidence"),
    }, ensure_ascii=False).lower()
    for token in ("fixture", "stub", "accepted-local", "manual seed", "broker", "level2", "paid token"):
        if token in forbidden:
            raise SystemExit(f"status=failed\nreason=forbidden_evidence:{token}")


def main() -> None:
    ARTIFACT_DIR.mkdir(parents=True, exist_ok=True)
    payload = {
        "status": "passed",
        "generated_at": datetime.now(timezone.utc).isoformat(),
        "checked_at_unix": int(time.time()),
        "symbol": "600000",
        "category": preverify_capital_flow("600000"),
        "claim_boundary": "Only this verified runtime public source can upgrade P90 capital-flow rows; fixture/stub/accepted-local/manual seed evidence remains ineligible.",
    }
    validate(payload)
    ARTIFACT.write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    if "--check" in sys.argv:
        stored = json.loads(ARTIFACT.read_text(encoding="utf-8"))
        validate(stored)
    print(f"p90_source_preverification:status=passed:eligible=1:artifact={ARTIFACT}")


if __name__ == "__main__":
    main()
