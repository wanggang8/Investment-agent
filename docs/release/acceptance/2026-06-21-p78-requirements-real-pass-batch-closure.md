# P78 Requirements Real-Pass Batch Closure Acceptance

> Date: 2026-06-21
> Change: `p78-requirements-real-pass-batch-closure`
> Conclusion: `release_ready_scoped_with_p78_real_pass_batch_progress`

## Summary

- Source matrix: `docs/release/acceptance/2026-06-21-p77-requirements-real-pass-upgrade-matrix.md`
- P78 matrix: `docs/release/acceptance/2026-06-21-p78-requirements-real-pass-batch-matrix.md`
- Summary JSON: `docs/release/ui-audit-assets/2026-06-21-p78/real-pass-batch-summary.json`
- Full-release-required rows: 330
- Full-release-required `real_pass` rows after P78: 20
- Remaining full-release-required non-`real_pass` rows: 310
- Newly upgraded by P78: 3

## P78 Batch A Upgrades

- `REQ-09-002` `9.2`: 只输出概率估算，不承诺结果。
- `REQ-09-017` `9.4`: 免责声明。
- `REQ-09-026` `9.7`: 历史类似样本数少于 20 个时，不输出精确概率，只输出“样本不足，仅作参考”。

## Fresh Evidence

- Expected-return Go tests: `docs/release/ui-audit-assets/2026-06-21-p78/expected-return-go-tests.log` and `docs/release/ui-audit-assets/2026-06-21-p78/expected-return-go-tests.json`
- Accepted-local non-`510300` real UI journey: `docs/release/ui-audit-assets/2026-06-21-p78-non-510300`
- Expected-return SQLite readback: `docs/release/ui-audit-assets/2026-06-21-p78/expected-return-ui-readback.json`

Commands:

```bash
go test -v ./internal/application/workflow -run 'TestBuildExpectedReturnIncludesSampleContextForAllPrecisionStates|TestBuildExpectedReturnProducesAdvisorySellEvaluation|TestBuildExpectedReturnDoesNotTriggerTargetWithoutConfiguredTarget|TestBuildExpectedReturnUsesScenarioBoundsForSellTriggers|TestBuildExpectedReturnCoversAllSellEvaluationTriggers|TestExpectedReturnNodeUsesWorkflowPricesForSellEvaluation|TestExpectedReturnNodeUsesMatchingSymbolPosition|TestExpectedReturnNodeUsesWorkflowDynamicSellInputs|TestExpectedReturnNodeIncludesP34SupportingDataContext|TestExpectedReturnSampleCountFromWorkflowDataUsesMarketHistory|TestExpectedReturnSampleCountFromWorkflowDataDoesNotInventSamples|TestBuildExpectedReturnExplainsMissingPriceContext' -count=1 && go test -v ./internal/domain/rule -run TestExpectedReturnDoesNotOverrideVerdict -count=1
P75_ARTIFACT_DIR=docs/release/ui-audit-assets/2026-06-21-p78-non-510300 bash scripts/p75-non-510300-real-ui-journey.sh
python3 scripts/p78_requirements_real_pass_batch_closure.py --check
```

## Boundaries

- P78 does not rewrite P75 or P77 historical matrices.
- P78 does not refresh the P76 package; a separate package refresh is required before claiming distribution archives include P78 materials.
- P78 does not claim full original-requirement pass while any full-release-required row remains non-`real_pass`.
- P78 does not add broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, real database overwrite, provider availability promises, or investment return promises.
