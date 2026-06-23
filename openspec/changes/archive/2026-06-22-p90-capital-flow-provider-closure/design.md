# P90 Design

## Provider Selection

The blocked P89 path used Eastmoney `push2` / `push2his` endpoints. P90 uses the Eastmoney H5 capital-flow endpoint observed from the public mobile capital-flow page:

- Page evidence: `https://emdatah5.eastmoney.com/dc/zjlx/stock?fc=1.600000&fn=浦发银行`
- JS evidence: `https://emdatah5.eastmoney.com/dc/Content/js/zjlx/stock.min.js`
- Runtime endpoint: `https://emdatah5.eastmoney.com/dc/ZJLX/getDBHistoryData`
- Request parameters: `secid=1.600000`, `fields1=f1,f2,f3`, `fields2=f51,f52,f53,f54,f55,f56,f62,f63`

The H5 table renders historical rows as:

- `f51` / index 0: date
- `f52` / index 1: daily net capital flow
- `f62` / index 6: close price
- `f63` / index 7: percent change

P90 maps daily net capital flow directionally:

- `net_inflow = max(f52, 0)`
- `net_outflow = max(-f52, 0)`
- `raw_net_flow = f52`

This preserves the public field semantics and does not invent gross inflow/outflow fields that the H5 history table does not publish.

## Product Proof Path

P90 acceptance must prove the product path:

1. Seed only local portfolio/capability facts needed to make the Settings UI refresh button visible.
2. Start the real local server with `use_stub=false` and `market_collectors.sources=[p89_structured_public]`.
3. Use Playwright to open `/settings` and click `刷新市场数据`.
4. Read `/api/v1/market/snapshots/latest?symbol=600000`.
5. Verify UI and API expose `capital_flow.date`, `capital_flow.net_inflow`, `capital_flow.net_outflow`, and `capital_flow.raw_net_flow`.
6. Verify SQLite readback from the runtime market snapshot ID, not a seeded snapshot.

## Safety

P90 remains read-only. It does not add broker/order/execution tables, trading UI, external push, automatic confirmation, automatic rule application, or return promises.

## Evidence Outputs

- `docs/release/ui-audit-assets/2026-06-22-p90-capital-flow-provider/p90-inventory.json`
- `docs/release/ui-audit-assets/2026-06-22-p90-capital-flow-provider/p90-source-preverification.json`
- `docs/release/ui-audit-assets/2026-06-22-p90-capital-flow-provider/browser-results.json`
- `docs/release/ui-audit-assets/2026-06-22-p90-capital-flow-provider/db-readback-check.log`
- `docs/release/ui-audit-assets/2026-06-22-p90-capital-flow-provider/p90-acceptance-summary.json`
- `docs/release/ui-audit-assets/2026-06-22-p90-capital-flow-provider/final-validation.log`
- `docs/release/acceptance/2026-06-22-p90-capital-flow-provider-closure.md`
- `docs/release/acceptance/2026-06-22-p90-capital-flow-provider-matrix.md`
