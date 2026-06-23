package service

import (
	"context"
	"encoding/json"
	"strings"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
	"investment-agent/internal/pkg/clock"
	"investment-agent/internal/pkg/idgen"
)

// SourceHealthRiskInput 是 source health 进入风险预警的最小上下文。
type SourceHealthRiskInput struct {
	SourceName   string `json:"source_name,omitempty"`
	Freshness    string `json:"freshness,omitempty"`
	DataCategory string `json:"data_category,omitempty"`
	DataDate     string `json:"data_date,omitempty"`
}

// RiskAlertTriggerInput 是风险预警触发输入。
type RiskAlertTriggerInput struct {
	RequestID             string
	DecisionID            string
	ReportID              string
	RiskType              model.RiskType
	Severity              model.RiskSeverity
	Symbol                string
	TriggerSummary        string
	TriggerContextJSON    string
	ProhibitedActionsJSON string
	SuggestedActionsJSON  string
}

// RiskAlertService 编排 P35 本地风险预警、通知和审计。
type RiskAlertService struct {
	tx  repository.Transactor
	clk clock.Clock
	ids idgen.Generator
}

func NewRiskAlertService(tx repository.Transactor) *RiskAlertService {
	return &RiskAlertService{tx: tx, clk: clock.SystemClock{}, ids: idgen.NewGenerator()}
}

func (s *RiskAlertService) BuildRiskAlertTriggers(decision repository.DecisionRecord, market model.MarketSnapshot, sourceHealth []SourceHealthRiskInput) []RiskAlertTriggerInput {
	symbol := strings.TrimSpace(decision.Symbol)
	if symbol == "" {
		symbol = strings.TrimSpace(market.Symbol)
	}
	var out []RiskAlertTriggerInput
	base := func(riskType model.RiskType, severity model.RiskSeverity, summary string, context map[string]any) RiskAlertTriggerInput {
		return RiskAlertTriggerInput{RequestID: decision.RequestID, DecisionID: decision.DecisionID, RiskType: riskType, Severity: severity, Symbol: symbol, TriggerSummary: summary, TriggerContextJSON: marshalRiskContext(context), ProhibitedActionsJSON: decision.ProhibitedActionsJSON, SuggestedActionsJSON: decision.OptionalActionsJSON}
	}
	if decision.FinalVerdictStatus == string(model.VerdictSellOnly) || decision.RiskReasonCode == string(model.RiskTypeBuyThesisBroken) {
		out = append(out, base(model.RiskTypeBuyThesisBroken, model.RiskSeverityCritical, "买入逻辑破坏，进入只卖不买复核", map[string]any{"final_verdict_status": decision.FinalVerdictStatus, "risk_reason_code": decision.RiskReasonCode}))
	}
	if market.PEPercentile >= 80 || market.PBPercentile >= 80 {
		out = append(out, base(model.RiskTypeValuationHigh, model.RiskSeverityWarning, "估值分位处于高位，暂停新增买入并复核止盈计划", map[string]any{"pe_percentile": market.PEPercentile, "pb_percentile": market.PBPercentile}))
	}
	if market.LiquidityState == model.LiquidityDanger {
		out = append(out, base(model.RiskTypeLiquidityDanger, model.RiskSeverityCritical, "流动性处于 danger，禁止市价式大额操作", map[string]any{"liquidity_state": market.LiquidityState}))
	}
	if market.SentimentState == model.SentimentExtreme {
		out = append(out, base(model.RiskTypeSentimentExtreme, model.RiskSeverityWarning, "市场情绪极端，进入冷静复核", map[string]any{"sentiment_state": market.SentimentState}))
	}
	if decision.RiskReasonCode == string(model.RiskTypePositionLimitBreach) {
		out = append(out, base(model.RiskTypePositionLimitBreach, model.RiskSeverityWarning, "仓位超过纪律上限，禁止新增买入", map[string]any{"risk_reason_code": decision.RiskReasonCode}))
	}
	if decision.SourceVerificationStatus == string(model.VerificationFailed) || decision.FinalVerdictStatus == string(model.VerdictInsufficientData) {
		out = append(out, base(model.RiskTypeInsufficientEvidence, model.RiskSeverityCritical, "证据不足，暂停交易类建议", map[string]any{"source_verification_status": decision.SourceVerificationStatus, "final_verdict_status": decision.FinalVerdictStatus}))
	}
	if degraded := degradedSourceHealth(sourceHealth); len(degraded) > 0 {
		out = append(out, base(model.RiskTypeDataDegraded, model.RiskSeverityWarning, "数据源新鲜度或可用性降级", map[string]any{"source_health": degraded}))
	}
	return out
}

func (s *RiskAlertService) TriggerRiskAlert(ctx context.Context, input RiskAlertTriggerInput) (repository.RiskAlert, error) {
	if err := validateRiskAlertTriggerInput(input); err != nil {
		return repository.RiskAlert{}, err
	}
	now := s.clk.NowRFC3339()
	alert := repository.RiskAlert{
		AlertID:               s.ids.New("risk"),
		RiskType:              input.RiskType,
		Severity:              input.Severity,
		SOPStatus:             model.RiskSOPActive,
		Symbol:                strings.TrimSpace(input.Symbol),
		TriggerSummary:        strings.TrimSpace(input.TriggerSummary),
		TriggerContextJSON:    input.TriggerContextJSON,
		ProhibitedActionsJSON: input.ProhibitedActionsJSON,
		SuggestedActionsJSON:  input.SuggestedActionsJSON,
		RelatedDecisionID:     input.DecisionID,
		RelatedReportID:       input.ReportID,
		RelatedNotificationID: s.ids.New("notif"),
		RelatedAuditEventID:   s.ids.New("audit"),
		LastTriggeredAt:       now,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
	if err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		if repos.RiskAlertRepo == nil || repos.NotificationRepo == nil || repos.AuditRepo == nil {
			return apperr.New(apperr.CodeInternalError, apperr.CategoryInternal, "risk alert repositories not configured")
		}
		if existing := findExistingActiveRisk(ctx, repos.RiskAlertRepo, alert.RiskType, alert.Symbol); existing.AlertID != "" {
			alert.AlertID = existing.AlertID
			if existing.RelatedNotificationID != "" {
				alert.RelatedNotificationID = existing.RelatedNotificationID
			}
		}
		if err := repos.RiskAlertRepo.UpsertRiskAlert(ctx, alert); err != nil {
			return err
		}
		if err := repos.NotificationRepo.SaveNotification(ctx, repository.Notification{NotificationID: alert.RelatedNotificationID, Type: "risk_alert", Severity: string(alert.Severity), Title: riskAlertTitle(alert.RiskType), Message: alert.TriggerSummary, SourceType: "risk_alert", SourceID: alert.AlertID, CreatedAt: now}); err != nil {
			return err
		}
		return repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: alert.RelatedAuditEventID, RequestID: input.RequestID, DecisionID: input.DecisionID, WorkflowType: "risk_alert_sop", NodeName: "RiskAlertService", Actor: string(model.AuditActorSystem), Action: string(model.AuditActionRiskAlert), NodeAction: "trigger_risk_alert", Status: string(model.AuditStatusSuccess), InputRefType: "decision", InputRef: input.DecisionID, OutputRefType: "risk_alert", OutputRef: alert.AlertID, CreatedAt: now})
	}); err != nil {
		return repository.RiskAlert{}, err
	}
	return alert, nil
}

func (s *RiskAlertService) UpdateRiskAlertLifecycle(ctx context.Context, alertID string, status model.RiskSOPStatus, reason string) (repository.RiskAlert, error) {
	if strings.TrimSpace(alertID) == "" || !status.Valid() {
		return repository.RiskAlert{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "invalid risk alert lifecycle update")
	}
	now := s.clk.NowRFC3339()
	auditID := s.ids.New("audit")
	var updated repository.RiskAlert
	if err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		if repos.RiskAlertRepo == nil || repos.AuditRepo == nil {
			return apperr.New(apperr.CodeInternalError, apperr.CategoryInternal, "risk alert repositories not configured")
		}
		before, err := repos.RiskAlertRepo.GetRiskAlert(ctx, alertID)
		if err != nil {
			return err
		}
		if !validRiskAlertLifecycleTransition(before.SOPStatus, status) {
			return apperr.New(apperr.CodeInvalidState, apperr.CategoryInvalidState, "invalid risk alert lifecycle transition")
		}
		if err := repos.RiskAlertRepo.UpdateRiskAlertStatus(ctx, alertID, status, reason, now); err != nil {
			return err
		}
		updated, err = repos.RiskAlertRepo.GetRiskAlert(ctx, alertID)
		if err != nil {
			return err
		}
		return repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: auditID, WorkflowType: "risk_alert_sop", NodeName: "RiskAlertService", Actor: string(model.AuditActorUser), Action: string(model.AuditActionRiskAlert), NodeAction: "update_risk_alert_lifecycle", Status: string(model.AuditStatusSuccess), BeforeState: string(before.SOPStatus), AfterState: string(status), OutputRefType: "risk_alert", OutputRef: alertID, CreatedAt: now})
	}); err != nil {
		return repository.RiskAlert{}, err
	}
	return updated, nil
}

func validateRiskAlertTriggerInput(input RiskAlertTriggerInput) error {
	if !input.RiskType.Valid() || !input.Severity.Valid() || strings.TrimSpace(input.Symbol) == "" || strings.TrimSpace(input.TriggerSummary) == "" {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "invalid risk alert trigger input")
	}
	return nil
}

func validRiskAlertLifecycleTransition(from model.RiskSOPStatus, to model.RiskSOPStatus) bool {
	if from == to {
		return true
	}
	if from.IsTerminal() {
		return false
	}
	switch from {
	case model.RiskSOPTriggered, model.RiskSOPActive, model.RiskSOPObserving:
		return to == model.RiskSOPObserving || to == model.RiskSOPEscalated || to == model.RiskSOPResolved || to == model.RiskSOPArchived
	case model.RiskSOPEscalated:
		return to == model.RiskSOPResolved || to == model.RiskSOPArchived
	default:
		return false
	}
}

func degradedSourceHealth(items []SourceHealthRiskInput) []SourceHealthRiskInput {
	var degraded []SourceHealthRiskInput
	for _, item := range items {
		switch item.Freshness {
		case "fresh", "":
			continue
		default:
			degraded = append(degraded, item)
		}
	}
	return degraded
}

func SourceHealthRiskInputsFromExpectedReturnJSON(raw string) []SourceHealthRiskInput {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var payload struct {
		SourceHealth []SourceHealthRiskInput `json:"source_health"`
	}
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return nil
	}
	return payload.SourceHealth
}

func findExistingActiveRisk(ctx context.Context, repo repository.RiskAlertRepository, riskType model.RiskType, symbol string) repository.RiskAlert {
	items, err := repo.ListRiskAlerts(ctx, repository.RiskAlertFilter{SOPStatuses: []model.RiskSOPStatus{model.RiskSOPTriggered, model.RiskSOPActive, model.RiskSOPObserving, model.RiskSOPEscalated}, Symbol: symbol})
	if err != nil {
		return repository.RiskAlert{}
	}
	for _, item := range items {
		if item.RiskType == riskType {
			return item
		}
	}
	return repository.RiskAlert{}
}

func marshalRiskContext(value map[string]any) string {
	if len(value) == 0 {
		return ""
	}
	b, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(b)
}

func riskAlertTitle(riskType model.RiskType) string {
	switch riskType {
	case model.RiskTypeValuationHigh:
		return "估值高位风险"
	case model.RiskTypeBuyThesisBroken:
		return "买入逻辑破坏风险"
	case model.RiskTypeLiquidityDanger:
		return "流动性风险"
	case model.RiskTypeSentimentExtreme:
		return "情绪极端风险"
	case model.RiskTypePositionLimitBreach:
		return "仓位超限风险"
	case model.RiskTypeInsufficientEvidence:
		return "证据不足风险"
	case model.RiskTypeDataDegraded:
		return "数据降级风险"
	default:
		return "风险预警"
	}
}
