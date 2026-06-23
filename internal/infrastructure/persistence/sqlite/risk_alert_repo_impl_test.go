package sqlite

import (
	"context"
	"testing"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

func TestRiskAlertRepositoryUpsertsActiveAlertAndListsByStatus(t *testing.T) {
	store, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	ctx := context.Background()
	if err := Migrate(ctx, store.DB); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	repo := NewRiskAlertRepository(store.DB)

	alert := repository.RiskAlert{
		AlertID:               "risk_1",
		RiskType:              model.RiskTypeValuationHigh,
		Severity:              model.RiskSeverityWarning,
		SOPStatus:             model.RiskSOPActive,
		Symbol:                "510300",
		TriggerSummary:        "PE 分位高于 80%",
		TriggerContextJSON:    `{"pe_percentile":88}`,
		ProhibitedActionsJSON: `["新增买入"]`,
		SuggestedActionsJSON:  `["人工复核分批止盈"]`,
		RelatedDecisionID:     "dec_1",
		RelatedReportID:       "daily_1",
		RelatedNotificationID: "notif_1",
		LastTriggeredAt:       "2026-06-15T09:30:00Z",
		CreatedAt:             "2026-06-15T09:30:00Z",
		UpdatedAt:             "2026-06-15T09:30:00Z",
	}
	if err := repo.UpsertRiskAlert(ctx, alert); err != nil {
		t.Fatalf("upsert risk alert: %v", err)
	}

	alert.AlertID = "risk_2"
	alert.TriggerSummary = "PE 分位仍高于 80%"
	alert.UpdatedAt = "2026-06-15T10:30:00Z"
	if err := repo.UpsertRiskAlert(ctx, alert); err != nil {
		t.Fatalf("upsert duplicate active risk alert: %v", err)
	}

	items, err := repo.ListRiskAlerts(ctx, repository.RiskAlertFilter{SOPStatuses: []model.RiskSOPStatus{model.RiskSOPActive}})
	if err != nil {
		t.Fatalf("list active risk alerts: %v", err)
	}
	if len(items) != 1 || items[0].AlertID != "risk_1" || items[0].TriggerSummary != "PE 分位仍高于 80%" {
		t.Fatalf("expected one updated active alert, got %+v", items)
	}
}

func TestRiskAlertRepositoryLifecycleTransition(t *testing.T) {
	store, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	ctx := context.Background()
	if err := Migrate(ctx, store.DB); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	repo := NewRiskAlertRepository(store.DB)

	alert := repository.RiskAlert{AlertID: "risk_1", RiskType: model.RiskTypeInsufficientEvidence, Severity: model.RiskSeverityCritical, SOPStatus: model.RiskSOPActive, Symbol: "510300", TriggerSummary: "缺少高等级独立证据", CreatedAt: "2026-06-15T09:30:00Z", UpdatedAt: "2026-06-15T09:30:00Z"}
	if err := repo.UpsertRiskAlert(ctx, alert); err != nil {
		t.Fatalf("upsert risk alert: %v", err)
	}
	if err := repo.UpdateRiskAlertStatus(ctx, "risk_1", model.RiskSOPResolved, "补齐证据后解除", "2026-06-15T11:00:00Z"); err != nil {
		t.Fatalf("update risk alert status: %v", err)
	}

	got, err := repo.GetRiskAlert(ctx, "risk_1")
	if err != nil {
		t.Fatalf("get risk alert: %v", err)
	}
	if got.SOPStatus != model.RiskSOPResolved || got.ResolutionReason != "补齐证据后解除" || got.ResolvedAt != "2026-06-15T11:00:00Z" {
		t.Fatalf("unexpected resolved alert: %+v", got)
	}
	if err := repo.UpdateRiskAlertStatus(ctx, "risk_missing", model.RiskSOPArchived, "不存在", "2026-06-15T12:00:00Z"); !apperr.IsCode(err, apperr.CodeNotFound) {
		t.Fatalf("expected not found for missing alert update, got %v", err)
	}
}
