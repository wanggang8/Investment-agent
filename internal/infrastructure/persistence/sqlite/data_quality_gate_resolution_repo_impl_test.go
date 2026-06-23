package sqlite

import (
	"context"
	"testing"

	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

func TestDataQualityGateResolutionRepositoryCreateListAndRetire(t *testing.T) {
	store, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	ctx := context.Background()
	if err := Migrate(ctx, store.DB); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	repo := NewDataQualityGateResolutionRepository(store.DB)

	first := repository.DataQualityGateResolution{
		ResolutionID:        "dqgr_1",
		Symbol:              "000300",
		PolicyFingerprint:   "fp_current_block",
		PolicyVerdict:       "blocked",
		ReleaseGate:         "block",
		PolicySummary:       "data_source_quality:mode=current:policy=blocked:gate=block",
		ResolutionType:      "scope_exclusion",
		Status:              "active",
		Scope:               "release excludes current local data health",
		Reason:              "core source health stale in local DB",
		ReleaseImpact:       "do not claim current data clean",
		EvidenceRef:         "docs/release/acceptance/p66",
		BlockingReasonsJSON: `["index_valuation_files core category degraded freshness=stale"]`,
		WaiverReasonsJSON:   `[]`,
		CreatedBy:           "local_user",
		CreatedAt:           "2026-06-18T10:00:00Z",
		SafetyNote:          "只记录本地人工处置，不改变数据质量事实。",
	}
	if err := repo.CreateDataQualityGateResolution(ctx, first); err != nil {
		t.Fatalf("create first resolution: %v", err)
	}

	active, err := repo.GetActiveDataQualityGateResolution(ctx, "000300", "fp_current_block")
	if err != nil {
		t.Fatalf("get active resolution: %v", err)
	}
	if active.ResolutionID != "dqgr_1" || active.ResolutionType != "scope_exclusion" {
		t.Fatalf("unexpected active resolution: %+v", active)
	}

	conflict := first
	conflict.ResolutionID = "dqgr_2"
	conflict.ResolutionType = "waiver"
	if err := repo.CreateDataQualityGateResolution(ctx, conflict); err == nil {
		t.Fatal("expected unique active resolution conflict")
	}

	second := first
	second.ResolutionID = "dqgr_3"
	second.PolicyFingerprint = "fp_optional_waiver"
	second.PolicyVerdict = "waiver_required"
	second.ReleaseGate = "waiver_required"
	second.ResolutionType = "waiver"
	second.CreatedAt = "2026-06-18T11:00:00Z"
	if err := repo.CreateDataQualityGateResolution(ctx, second); err != nil {
		t.Fatalf("create second resolution: %v", err)
	}

	items, err := repo.ListDataQualityGateResolutions(ctx, repository.DataQualityGateResolutionFilter{Symbol: "000300"})
	if err != nil {
		t.Fatalf("list resolutions: %v", err)
	}
	if len(items) != 2 || items[0].ResolutionID != "dqgr_3" || items[1].ResolutionID != "dqgr_1" {
		t.Fatalf("expected newest first list, got %+v", items)
	}

	if err := repo.RetireDataQualityGateResolution(ctx, "dqgr_1", "local_user", "2026-06-18T12:00:00Z"); err != nil {
		t.Fatalf("retire resolution: %v", err)
	}
	retired, err := repo.GetDataQualityGateResolution(ctx, "dqgr_1")
	if err != nil {
		t.Fatalf("get retired resolution: %v", err)
	}
	if retired.Status != "retired" || retired.RetiredBy != "local_user" || retired.RetiredAt == "" {
		t.Fatalf("unexpected retired resolution: %+v", retired)
	}
	if _, err := repo.GetActiveDataQualityGateResolution(ctx, "000300", "fp_current_block"); !apperr.IsCode(err, apperr.CodeNotFound) {
		t.Fatalf("expected retired resolution ignored for active lookup, got %v", err)
	}
}
