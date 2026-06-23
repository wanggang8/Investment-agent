package repository

import (
	"context"

	"investment-agent/internal/domain/model"
)

// RiskAlert 是 P35 本地风险预警与 SOP 状态事实。
type RiskAlert struct {
	AlertID               string
	RiskType              model.RiskType
	Severity              model.RiskSeverity
	SOPStatus             model.RiskSOPStatus
	Symbol                string
	TriggerSummary        string
	TriggerContextJSON    string
	ProhibitedActionsJSON string
	SuggestedActionsJSON  string
	RelatedDecisionID     string
	RelatedReportID       string
	RelatedNotificationID string
	RelatedAuditEventID   string
	LastTriggeredAt       string
	ResolvedAt            string
	ResolutionReason      string
	CreatedAt             string
	UpdatedAt             string
}

// RiskAlertFilter 是风险预警列表过滤条件。
type RiskAlertFilter struct {
	SOPStatuses []model.RiskSOPStatus
	Symbol      string
}

// RiskAlertRepository 定义风险预警事实持久化边界。
type RiskAlertRepository interface {
	UpsertRiskAlert(ctx context.Context, alert RiskAlert) error
	GetRiskAlert(ctx context.Context, alertID string) (RiskAlert, error)
	ListRiskAlerts(ctx context.Context, filter RiskAlertFilter) ([]RiskAlert, error)
	UpdateRiskAlertStatus(ctx context.Context, alertID string, status model.RiskSOPStatus, reason string, updatedAt string) error
}
