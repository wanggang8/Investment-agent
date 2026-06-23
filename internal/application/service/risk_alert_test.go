package service

import (
	"context"
	"testing"
	"time"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/infrastructure/persistence/sqlite"
	"investment-agent/internal/pkg/apperr"
	"investment-agent/internal/pkg/clock"
	"investment-agent/internal/pkg/idgen"
)

func TestRiskAlertServiceBuildsTriggersFromDecisionAndMarket(t *testing.T) {
	sourceHealth := SourceHealthRiskInputsFromExpectedReturnJSON(`{"source_health":[{"source_name":"csindex_extended","freshness":"stale","data_category":"index_valuation_files","data_date":"2026-06-05"}]}`)
	if len(sourceHealth) != 1 || sourceHealth[0].Freshness != "stale" {
		t.Fatalf("expected source health inputs from expected return json, got %+v", sourceHealth)
	}
	svc := NewRiskAlertService(transactorStub{})
	inputs := svc.BuildRiskAlertTriggers(repository.DecisionRecord{
		RequestID:                "req_1",
		DecisionID:               "dec_1",
		Symbol:                   "510300",
		SourceVerificationStatus: string(model.VerificationFailed),
		FinalVerdictStatus:       string(model.VerdictSellOnly),
		RiskReasonCode:           "buy_thesis_broken",
		ProhibitedActionsJSON:    `["新增买入","加仓"]`,
		OptionalActionsJSON:      `["人工复核卖出计划"]`,
	}, model.MarketSnapshot{Symbol: "510300", PEPercentile: 88, LiquidityState: model.LiquidityDanger, SentimentState: model.SentimentExtreme}, sourceHealth)

	seen := map[model.RiskType]RiskAlertTriggerInput{}
	for _, input := range inputs {
		seen[input.RiskType] = input
	}
	for _, riskType := range []model.RiskType{model.RiskTypeBuyThesisBroken, model.RiskTypeValuationHigh, model.RiskTypeLiquidityDanger, model.RiskTypeSentimentExtreme, model.RiskTypeInsufficientEvidence, model.RiskTypeDataDegraded} {
		if _, ok := seen[riskType]; !ok {
			t.Fatalf("missing risk type %s in %+v", riskType, inputs)
		}
	}
	positionInputs := svc.BuildRiskAlertTriggers(repository.DecisionRecord{DecisionID: "dec_2", Symbol: "510300", RiskReasonCode: string(model.RiskTypePositionLimitBreach)}, model.MarketSnapshot{}, nil)
	if len(positionInputs) != 1 || positionInputs[0].RiskType != model.RiskTypePositionLimitBreach {
		t.Fatalf("expected position limit risk, got %+v", positionInputs)
	}
	if seen[model.RiskTypeBuyThesisBroken].Severity != model.RiskSeverityCritical || seen[model.RiskTypeDataDegraded].TriggerContextJSON == "" {
		t.Fatalf("unexpected generated triggers: %+v", seen)
	}
}

func TestRiskAlertServiceUsesValuationHighRiskBoundaryAtEightyPercent(t *testing.T) {
	svc := NewRiskAlertService(transactorStub{})
	inputs := svc.BuildRiskAlertTriggers(repository.DecisionRecord{
		RequestID:  "req_valuation_boundary",
		DecisionID: "dec_valuation_boundary",
		Symbol:     "510300",
	}, model.MarketSnapshot{Symbol: "510300", PEPercentile: 80, PBPercentile: 30}, nil)

	for _, input := range inputs {
		if input.RiskType == model.RiskTypeValuationHigh {
			return
		}
	}
	t.Fatalf("expected valuation high risk alert at 80%% boundary, got %+v", inputs)
}

func TestRiskAlertServiceTriggersRiskWithNotificationAndAudit(t *testing.T) {
	store, err := sqlite.Open(":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	ctx := context.Background()
	if err := sqlite.Migrate(ctx, store.DB); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	svc := NewRiskAlertService(sqlite.NewTransactor(store.DB))
	svc.clk = clock.FixedClock{Time: time.Date(2026, 6, 15, 9, 30, 0, 0, time.UTC)}
	svc.ids = idgen.NewFixedGenerator(map[string][]string{"risk": {"risk_1"}, "notif": {"notif_1"}, "audit": {"audit_1"}})

	alert, err := svc.TriggerRiskAlert(ctx, RiskAlertTriggerInput{
		RequestID:             "req_1",
		DecisionID:            "dec_1",
		ReportID:              "daily_1",
		RiskType:              model.RiskTypeValuationHigh,
		Severity:              model.RiskSeverityWarning,
		Symbol:                "510300",
		TriggerSummary:        "PE 分位高于 80%",
		TriggerContextJSON:    `{"pe_percentile":88}`,
		ProhibitedActionsJSON: `["新增买入"]`,
		SuggestedActionsJSON:  `["人工复核分批止盈"]`,
	})
	if err != nil {
		t.Fatalf("trigger risk alert: %v", err)
	}
	if alert.AlertID != "risk_1" || alert.SOPStatus != model.RiskSOPActive || alert.RelatedNotificationID != "notif_1" || alert.RelatedAuditEventID != "audit_1" {
		t.Fatalf("unexpected alert: %+v", alert)
	}

	repos := repository.Repositories{RiskAlertRepo: sqlite.NewRiskAlertRepository(store.DB), NotificationRepo: sqlite.NewNotificationRepository(store.DB), AuditRepo: sqlite.NewAuditRepository(store.DB)}
	items, err := repos.RiskAlertRepo.ListRiskAlerts(ctx, repository.RiskAlertFilter{SOPStatuses: []model.RiskSOPStatus{model.RiskSOPActive}})
	if err != nil || len(items) != 1 {
		t.Fatalf("list risk alerts: items=%+v err=%v", items, err)
	}
	notifications, err := repos.NotificationRepo.ListNotifications(ctx)
	if err != nil || len(notifications) != 1 || notifications[0].Type != "risk_alert" || notifications[0].SourceID != "risk_1" {
		t.Fatalf("unexpected notifications: %+v err=%v", notifications, err)
	}

	svc.ids = idgen.NewFixedGenerator(map[string][]string{"risk": {"risk_2"}, "notif": {"notif_2"}, "audit": {"audit_2"}})
	if _, err := svc.TriggerRiskAlert(ctx, RiskAlertTriggerInput{RequestID: "req_2", DecisionID: "dec_2", ReportID: "daily_2", RiskType: model.RiskTypeValuationHigh, Severity: model.RiskSeverityCritical, Symbol: "510300", TriggerSummary: "PE 分位仍高于 80%"}); err != nil {
		t.Fatalf("trigger duplicate risk alert: %v", err)
	}
	notifications, err = repos.NotificationRepo.ListNotifications(ctx)
	if err != nil || len(notifications) != 1 || notifications[0].SourceID != "risk_1" || notifications[0].Severity != string(model.RiskSeverityCritical) {
		t.Fatalf("expected deduplicated risk notification, got %+v err=%v", notifications, err)
	}
	updatedAlert, err := repos.RiskAlertRepo.GetRiskAlert(ctx, "risk_1")
	if err != nil || updatedAlert.RelatedNotificationID != "notif_1" {
		t.Fatalf("expected duplicate trigger to preserve existing notification id, got %+v err=%v", updatedAlert, err)
	}
	audits, err := repos.AuditRepo.ListAuditEvents(ctx)
	if err != nil || len(audits) == 0 || audits[0].Action != string(model.AuditActionRiskAlert) || audits[0].OutputRef != "risk_1" {
		t.Fatalf("unexpected audits: %+v err=%v", audits, err)
	}
}

func TestRiskAlertServiceUpdatesLifecycleWithAuditOnly(t *testing.T) {
	store, err := sqlite.Open(":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()
	ctx := context.Background()
	if err := sqlite.Migrate(ctx, store.DB); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	repo := sqlite.NewRiskAlertRepository(store.DB)
	if err := repo.UpsertRiskAlert(ctx, repository.RiskAlert{AlertID: "risk_1", RiskType: model.RiskTypeDataDegraded, Severity: model.RiskSeverityWarning, SOPStatus: model.RiskSOPActive, Symbol: "510300", TriggerSummary: "source health stale", CreatedAt: "2026-06-15T09:30:00Z", UpdatedAt: "2026-06-15T09:30:00Z"}); err != nil {
		t.Fatalf("seed risk alert: %v", err)
	}

	svc := NewRiskAlertService(sqlite.NewTransactor(store.DB))
	svc.clk = clock.FixedClock{Time: time.Date(2026, 6, 15, 10, 30, 0, 0, time.UTC)}
	svc.ids = idgen.NewFixedGenerator(map[string][]string{"audit": {"audit_2"}})
	updated, err := svc.UpdateRiskAlertLifecycle(ctx, "risk_1", model.RiskSOPResolved, "数据源恢复")
	if err != nil {
		t.Fatalf("update lifecycle: %v", err)
	}
	if updated.SOPStatus != model.RiskSOPResolved || updated.ResolutionReason != "数据源恢复" || updated.ResolvedAt != "2026-06-15T10:30:00Z" {
		t.Fatalf("unexpected updated alert: %+v", updated)
	}

	if err := repo.UpsertRiskAlert(ctx, repository.RiskAlert{AlertID: "risk_2", RiskType: model.RiskTypeLiquidityDanger, Severity: model.RiskSeverityCritical, SOPStatus: model.RiskSOPActive, Symbol: "159915", TriggerSummary: "流动性 danger", CreatedAt: "2026-06-15T09:30:00Z", UpdatedAt: "2026-06-15T09:30:00Z"}); err != nil {
		t.Fatalf("seed second risk alert: %v", err)
	}
	svc.ids = idgen.NewFixedGenerator(map[string][]string{"audit": {"audit_3", "audit_4"}})
	if updated, err = svc.UpdateRiskAlertLifecycle(ctx, "risk_2", model.RiskSOPEscalated, "连续触发"); err != nil || updated.SOPStatus != model.RiskSOPEscalated {
		t.Fatalf("expected escalated alert, got %+v err=%v", updated, err)
	}
	if updated, err = svc.UpdateRiskAlertLifecycle(ctx, "risk_2", model.RiskSOPArchived, "用户归档"); err != nil || updated.SOPStatus != model.RiskSOPArchived || updated.ResolvedAt == "" {
		t.Fatalf("expected archived alert, got %+v err=%v", updated, err)
	}
	if _, err := svc.UpdateRiskAlertLifecycle(ctx, "risk_2", model.RiskSOPActive, "非法恢复"); !apperr.IsCode(err, apperr.CodeInvalidState) {
		t.Fatalf("expected invalid state for terminal lifecycle transition, got %v", err)
	}

	notifications, err := sqlite.NewNotificationRepository(store.DB).ListNotifications(ctx)
	if err != nil || len(notifications) != 0 {
		t.Fatalf("lifecycle update should not create notification: %+v err=%v", notifications, err)
	}
}
