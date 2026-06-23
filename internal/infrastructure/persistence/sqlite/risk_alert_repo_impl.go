package sqlite

import (
	"context"
	"database/sql"
	"strings"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

// RiskAlertRepository 是 P35 风险预警事实的 SQLite 实现。
type RiskAlertRepository struct{ db dbtx }

// NewRiskAlertRepository 创建风险预警仓储实例。
func NewRiskAlertRepository(db *sql.DB) *RiskAlertRepository {
	return &RiskAlertRepository{db: db}
}

func (r *RiskAlertRepository) UpsertRiskAlert(ctx context.Context, alert repository.RiskAlert) error {
	if err := validateRiskAlert(alert); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx, `INSERT INTO risk_alerts (alert_id,risk_type,severity,sop_status,symbol,trigger_summary,trigger_context_json,prohibited_actions_json,suggested_actions_json,related_decision_id,related_report_id,related_notification_id,related_audit_event_id,last_triggered_at,resolved_at,resolution_reason,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) ON CONFLICT(risk_type, symbol) WHERE sop_status IN ('triggered','active','observing','escalated') DO UPDATE SET severity=excluded.severity,sop_status=excluded.sop_status,trigger_summary=excluded.trigger_summary,trigger_context_json=excluded.trigger_context_json,prohibited_actions_json=excluded.prohibited_actions_json,suggested_actions_json=excluded.suggested_actions_json,related_decision_id=excluded.related_decision_id,related_report_id=excluded.related_report_id,related_notification_id=excluded.related_notification_id,related_audit_event_id=excluded.related_audit_event_id,last_triggered_at=excluded.last_triggered_at,updated_at=excluded.updated_at`, alert.AlertID, string(alert.RiskType), string(alert.Severity), string(alert.SOPStatus), alert.Symbol, alert.TriggerSummary, nullString(alert.TriggerContextJSON), nullString(alert.ProhibitedActionsJSON), nullString(alert.SuggestedActionsJSON), nullString(alert.RelatedDecisionID), nullString(alert.RelatedReportID), nullString(alert.RelatedNotificationID), nullString(alert.RelatedAuditEventID), nullString(alert.LastTriggeredAt), nullString(alert.ResolvedAt), nullString(alert.ResolutionReason), alert.CreatedAt, alert.UpdatedAt)
	return apperr.FromRepositoryError(err)
}

func (r *RiskAlertRepository) GetRiskAlert(ctx context.Context, alertID string) (repository.RiskAlert, error) {
	row := r.db.QueryRowContext(ctx, riskAlertSelectSQL()+` WHERE alert_id=?`, alertID)
	return scanRiskAlert(row)
}

func (r *RiskAlertRepository) ListRiskAlerts(ctx context.Context, filter repository.RiskAlertFilter) ([]repository.RiskAlert, error) {
	query := riskAlertSelectSQL()
	var args []any
	var clauses []string
	if len(filter.SOPStatuses) > 0 {
		placeholders := make([]string, 0, len(filter.SOPStatuses))
		for _, status := range filter.SOPStatuses {
			if !status.Valid() {
				return nil, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "invalid risk alert status")
			}
			placeholders = append(placeholders, "?")
			args = append(args, string(status))
		}
		clauses = append(clauses, "sop_status IN ("+strings.Join(placeholders, ",")+")")
	}
	if strings.TrimSpace(filter.Symbol) != "" {
		clauses = append(clauses, "symbol=?")
		args = append(args, strings.TrimSpace(filter.Symbol))
	}
	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}
	query += " ORDER BY updated_at DESC, alert_id DESC"
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperr.FromRepositoryError(err)
	}
	defer rows.Close()
	var alerts []repository.RiskAlert
	for rows.Next() {
		alert, err := scanRiskAlert(rows)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}
	return alerts, apperr.FromRepositoryError(rows.Err())
}

func (r *RiskAlertRepository) UpdateRiskAlertStatus(ctx context.Context, alertID string, status model.RiskSOPStatus, reason string, updatedAt string) error {
	if !status.Valid() {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "invalid risk alert status")
	}
	resolvedAt := ""
	if status.IsTerminal() {
		resolvedAt = updatedAt
	}
	result, err := r.db.ExecContext(ctx, `UPDATE risk_alerts SET sop_status=?,resolution_reason=?,resolved_at=?,updated_at=? WHERE alert_id=?`, string(status), nullString(reason), nullString(resolvedAt), updatedAt, alertID)
	if err != nil {
		return apperr.FromRepositoryError(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return apperr.FromRepositoryError(err)
	}
	if rows == 0 {
		return apperr.New(apperr.CodeNotFound, apperr.CategoryNotFound, "record not found")
	}
	return nil
}

func validateRiskAlert(alert repository.RiskAlert) error {
	if strings.TrimSpace(alert.AlertID) == "" || strings.TrimSpace(alert.Symbol) == "" || strings.TrimSpace(alert.TriggerSummary) == "" {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "risk alert requires id, symbol and trigger summary")
	}
	if !alert.RiskType.Valid() || !alert.Severity.Valid() || !alert.SOPStatus.Valid() {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "invalid risk alert enum")
	}
	return nil
}

func riskAlertSelectSQL() string {
	return `SELECT alert_id,risk_type,severity,sop_status,symbol,trigger_summary,COALESCE(trigger_context_json,''),COALESCE(prohibited_actions_json,''),COALESCE(suggested_actions_json,''),COALESCE(related_decision_id,''),COALESCE(related_report_id,''),COALESCE(related_notification_id,''),COALESCE(related_audit_event_id,''),COALESCE(last_triggered_at,''),COALESCE(resolved_at,''),COALESCE(resolution_reason,''),created_at,updated_at FROM risk_alerts`
}

type riskAlertScanner interface {
	Scan(dest ...any) error
}

func scanRiskAlert(scanner riskAlertScanner) (repository.RiskAlert, error) {
	var alert repository.RiskAlert
	var riskType, severity, sopStatus string
	if err := scanner.Scan(&alert.AlertID, &riskType, &severity, &sopStatus, &alert.Symbol, &alert.TriggerSummary, &alert.TriggerContextJSON, &alert.ProhibitedActionsJSON, &alert.SuggestedActionsJSON, &alert.RelatedDecisionID, &alert.RelatedReportID, &alert.RelatedNotificationID, &alert.RelatedAuditEventID, &alert.LastTriggeredAt, &alert.ResolvedAt, &alert.ResolutionReason, &alert.CreatedAt, &alert.UpdatedAt); err != nil {
		return repository.RiskAlert{}, apperr.FromRepositoryError(err)
	}
	alert.RiskType = model.RiskType(riskType)
	alert.Severity = model.RiskSeverity(severity)
	alert.SOPStatus = model.RiskSOPStatus(sopStatus)
	return alert, nil
}
