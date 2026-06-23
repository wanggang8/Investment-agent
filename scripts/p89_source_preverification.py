#!/usr/bin/env python3
"""P89 real public structured-source preverification.

This script performs low-frequency read-only checks against public pages/APIs.
It writes field-level evidence only when required values are present. Fixture,
stub, accepted-local, and manually seeded evidence are never eligible here.
"""

from __future__ import annotations

import json
import math
import subprocess
import sys
import time
import urllib.parse
import urllib.request
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
ARTIFACT_DIR = ROOT / "docs" / "release" / "ui-audit-assets" / "2026-06-22-p89-real-provider-dynamic-probability"
ARTIFACT = ARTIFACT_DIR / "p89-source-preverification.json"
TIMEOUT = 20


def fetch_json(url: str, headers: dict[str, str] | None = None) -> dict:
    req = urllib.request.Request(url, headers={
        "User-Agent": "Mozilla/5.0 InvestmentAgent-P89 readonly acceptance",
        "Accept": "application/json,text/plain,*/*",
        "Connection": "close",
        **(headers or {}),
    })
    try:
        with urllib.request.urlopen(req, timeout=TIMEOUT) as resp:
            raw = resp.read().decode("utf-8", errors="replace")
    except Exception:
        cmd = ["curl", "-L", "--max-time", str(TIMEOUT), "-s"]
        for key, value in (headers or {}).items():
            cmd.extend(["-H", f"{key}: {value}"])
        cmd.append(url)
        raw = subprocess.check_output(cmd, text=True)
    if raw.startswith("(") and raw.endswith(")"):
        raw = raw[1:-1]
    return json.loads(raw)


def capital_flow(symbol: str) -> dict:
    secid = f"1.{symbol}" if symbol.startswith("6") else f"0.{symbol}"
    query = urllib.parse.urlencode({
        "lmt": "1",
        "klt": "101",
        "fields1": "f1,f2,f3,f7",
        "fields2": "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f62,f63",
        "secid": secid,
    }, safe=",")
    url = f"https://push2.eastmoney.com/api/qt/stock/fflow/daykline/get?{query}"
    payload = fetch_json(url)
    rows = (((payload.get("data") or {}).get("klines")) or [])
    if not rows:
        raise SystemExit("status=failed\nreason=capital_flow_no_rows")
    fields = str(rows[0]).split(",")
    if len(fields) < 4:
        raise SystemExit("status=failed\nreason=capital_flow_fields")
    date = fields[0]
    net_inflow = float(fields[1])
    net_outflow = abs(float(fields[3]))
    return {
        "category": "capital_flow",
        "requirement_ids": ["REQ-05-003"],
        "provider_status": "verified_runtime_public",
        "real_pass_eligible": True,
        "authority": "Eastmoney public market data endpoint",
        "public_access_shape": "HTTPS GET, no login observed, no paid token observed, daily kline request",
        "stable_request_or_page_evidence": url,
        "fields": {"date": date, "net_inflow": net_inflow, "net_outflow": net_outflow},
        "update_frequency": "daily historical fund-flow kline; P89 acceptance uses one low-frequency request",
        "legal_access_limits": "public endpoint observed without login or token; not treated as Level2 or high-frequency; production use remains low-frequency/manual unless provider terms are separately reviewed",
        "rate_limit_assumption": "single acceptance request; no polling or high-frequency refresh",
        "failure_behavior": "source_unavailable; do not synthesize values",
        "sqlite_target_path": "market_snapshots.market_metrics_json.metadata.p88_structured_fields.capital_flow",
        "source_name": "eastmoney_push2_public",
    }


def margin_financing() -> dict:
    query = urllib.parse.urlencode({
        "isPagination": "true",
        "sqlId": "COMMON_SSE_SJ_GPSJ_GPHYSJ_MX_L",
        "pageHelp.pageSize": "2",
        "pageHelp.pageNo": "1",
    })
    url = f"https://query.sse.com.cn/marketdata/tradedata/queryMargin.do?{query}"
    payload = fetch_json(url, {"Referer": "https://www.sse.com.cn/market/othersdata/margin/detail/"})
    rows = (((payload.get("pageHelp") or {}).get("data")) or [])
    if len(rows) < 2:
        raise SystemExit("status=failed\nreason=margin_financing_rows")
    latest, previous = rows[0], rows[1]
    latest_balance = float(latest["rzye"])
    previous_balance = float(previous["rzye"])
    change_rate = 0.0 if previous_balance == 0 else (latest_balance - previous_balance) / previous_balance
    date = str(latest["opDate"])
    date = f"{date[:4]}-{date[4:6]}-{date[6:]}" if len(date) == 8 else date
    return {
        "category": "margin_financing",
        "requirement_ids": ["REQ-05-004"],
        "provider_status": "verified_runtime_public",
        "real_pass_eligible": True,
        "authority": "Shanghai Stock Exchange public margin data endpoint",
        "public_access_shape": "HTTPS GET with SSE page referer, no login observed, no paid token observed",
        "stable_request_or_page_evidence": url,
        "reference_page": "https://www.sse.com.cn/market/othersdata/margin/detail/",
        "fields": {"date": date, "margin_balance": latest_balance, "balance_change_rate": change_rate},
        "update_frequency": "daily exchange margin series",
        "legal_access_limits": "official public SSE page/API observed without login; P89 uses low-frequency read-only acceptance only",
        "rate_limit_assumption": "single acceptance request; no polling or high-frequency refresh",
        "failure_behavior": "source_unavailable; do not synthesize values",
        "sqlite_target_path": "market_snapshots.margin_balance, market_snapshots.margin_balance_change, and market_snapshots.market_metrics_json.metadata.p88_structured_fields.margin_financing",
        "source_name": "sse_query_margin_public",
        "symbol_scope": "market-level official margin series; not a broker/order/trading source",
    }


def constituent_financial(symbol: str) -> dict:
    query = urllib.parse.urlencode({
        "sortColumns": "NOTICE_DATE",
        "sortTypes": "-1",
        "pageSize": "1",
        "pageNumber": "1",
        "reportName": "RPT_LICO_FN_CPD",
        "columns": "ALL",
        "filter": f'(SECURITY_CODE="{symbol}")',
    })
    url = f"https://datacenter-web.eastmoney.com/api/data/v1/get?{query}"
    payload = fetch_json(url)
    rows = (((payload.get("result") or {}).get("data")) or [])
    if not rows:
        raise SystemExit("status=failed\nreason=constituent_financial_no_rows")
    row = rows[0]
    disclosure_date = str(row.get("NOTICE_DATE") or "").replace(" 00:00:00", "")
    fields = {
        "revenue": float(row["TOTAL_OPERATE_INCOME"]),
        "net_profit": float(row["PARENT_NETPROFIT"]),
        "growth": float(row["SJLTZ"]),
        "disclosure_date": disclosure_date,
    }
    if not disclosure_date or fields["revenue"] == 0 or fields["net_profit"] == 0:
        raise SystemExit("status=failed\nreason=constituent_financial_required_fields")
    return {
        "category": "constituent_financial",
        "requirement_ids": ["REQ-05-005"],
        "provider_status": "verified_runtime_public",
        "real_pass_eligible": True,
        "authority": "Eastmoney public datacenter financial report endpoint",
        "public_access_shape": "HTTPS GET, no login observed, no paid token observed, quarterly/annual report table",
        "stable_request_or_page_evidence": url,
        "fields": fields,
        "update_frequency": "quarterly/annual disclosure cadence",
        "legal_access_limits": "public endpoint observed without login or token; P89 treats it as public read-only report data and not as an authorization source",
        "rate_limit_assumption": "single acceptance request; no polling or high-frequency refresh",
        "failure_behavior": "source_unavailable; do not synthesize values",
        "sqlite_target_path": "market_snapshots.market_metrics_json.metadata.p88_structured_fields.constituent_financial",
        "source_name": "eastmoney_datacenter_public",
    }


def validate(payload: dict) -> None:
    categories = payload.get("categories") or []
    if len(categories) != 3:
        raise SystemExit("status=failed\nreason=category_count")
    forbidden = ("accepted", "fixture", "stub", "manual")
    for item in categories:
        if not item.get("real_pass_eligible"):
            if not item.get("blocker"):
                raise SystemExit(f"status=failed\nreason=blocked_without_reason:{item.get('category')}")
            continue
        source_name = str(item.get("source_name") or "").lower()
        if any(token in source_name for token in forbidden):
            raise SystemExit(f"status=failed\nreason=forbidden_source:{item.get('category')}:{source_name}")
        fields = item.get("fields") or {}
        if item["category"] == "capital_flow" and not all(key in fields for key in ("date", "net_inflow", "net_outflow")):
            raise SystemExit("status=failed\nreason=capital_fields")
        if item["category"] == "margin_financing" and not all(key in fields for key in ("date", "margin_balance", "balance_change_rate")):
            raise SystemExit("status=failed\nreason=margin_fields")
        if item["category"] == "constituent_financial" and not all(key in fields for key in ("revenue", "net_profit", "growth", "disclosure_date")):
            raise SystemExit("status=failed\nreason=financial_fields")


def build() -> dict:
    symbol = "600000"
    builders = [
        ("capital_flow", ["REQ-05-003"], capital_flow),
        ("margin_financing", ["REQ-05-004"], lambda _symbol: margin_financing()),
        ("constituent_financial", ["REQ-05-005"], constituent_financial),
    ]
    categories = []
    for category, requirement_ids, builder in builders:
        try:
            categories.append(builder(symbol))
        except Exception as exc:
            categories.append({
                "category": category,
                "requirement_ids": requirement_ids,
                "provider_status": "blocked",
                "real_pass_eligible": False,
                "required_fields": {
                    "capital_flow": ["date", "net_inflow", "net_outflow"],
                    "margin_financing": ["date", "margin_balance", "balance_change_rate"],
                    "constituent_financial": ["revenue", "net_profit", "growth", "disclosure_date"],
                }[category],
                "failure_behavior": "source_unavailable; do not synthesize values",
                "sqlite_target_path": {
                    "capital_flow": "market_snapshots.market_metrics_json.metadata.p88_structured_fields.capital_flow",
                    "margin_financing": "market_snapshots.margin_balance, market_snapshots.margin_balance_change, and market_snapshots.market_metrics_json.metadata.p88_structured_fields.margin_financing",
                    "constituent_financial": "market_snapshots.market_metrics_json.metadata.p88_structured_fields.constituent_financial",
                }[category],
                "blocker": f"P89 live public provider verification failed for {category}: {exc}",
            })
    return {
        "change": "p89-real-provider-and-dynamic-probability-closure",
        "artifact": str(ARTIFACT.relative_to(ROOT)),
        "captured_at_epoch": int(time.time()),
        "symbol": symbol,
        "policy": "Only verified live public no-login/no-paid/no-auth/no-Level2/no-high-frequency provider evidence is eligible; fixture/stub/accepted-local/manual seed evidence is excluded.",
        "categories": categories,
        "summary": {
            "category_count": len(categories),
            "real_pass_eligible_count": sum(1 for item in categories if item["real_pass_eligible"]),
            "blocked_count": sum(1 for item in categories if not item["real_pass_eligible"]),
        },
    }


def main() -> None:
    if "--check" in sys.argv:
        payload = json.loads(ARTIFACT.read_text(encoding="utf-8"))
    else:
        ARTIFACT_DIR.mkdir(parents=True, exist_ok=True)
        payload = build()
        ARTIFACT.write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    validate(payload)
    print(f"p89_source_preverification:status=passed:eligible={payload['summary']['real_pass_eligible_count']}:blocked={payload['summary']['blocked_count']}:artifact={ARTIFACT}")


if __name__ == "__main__":
    main()
