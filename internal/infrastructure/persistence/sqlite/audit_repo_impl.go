package sqlite

import (
	"context"
	"database/sql"

	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

// AuditRepository 是审计事件表的 SQLite 实现。
type AuditRepository struct{ db dbtx }

// NewAuditRepository 创建审计仓储实例。
func NewAuditRepository(db *sql.DB) *AuditRepository { return &AuditRepository{db: db} }

// AppendAuditEvent 追加写入审计事件。审计事件不做更新，便于保留完整历史。
func (r *AuditRepository) AppendAuditEvent(ctx context.Context, e repository.AuditEvent) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO audit_events (audit_event_id,request_id,decision_id,workflow_type,node_name,actor,action,node_action,proposal_id,confirmation_id,error_case_id,status,error_code,before_state,after_state,rule_version,snapshot_id,input_ref_type,input_ref,output_ref_type,output_ref,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, e.AuditEventID, nullString(e.RequestID), nullString(e.DecisionID), nullString(e.WorkflowType), nullString(e.NodeName), e.Actor, e.Action, nullString(e.NodeAction), nullString(e.ProposalID), nullString(e.ConfirmationID), nullString(e.ErrorCaseID), e.Status, nullString(e.ErrorCode), nullString(e.BeforeState), nullString(e.AfterState), nullString(e.RuleVersion), nullString(e.SnapshotID), nullString(e.InputRefType), nullString(e.InputRef), nullString(e.OutputRefType), nullString(e.OutputRef), e.CreatedAt)
	return apperr.FromRepositoryError(err)
}

// GetAuditEvent 按事件 ID 读取审计事件。
func (r *AuditRepository) GetAuditEvent(ctx context.Context, id string) (repository.AuditEvent, error) {
	var e repository.AuditEvent
	err := r.db.QueryRowContext(ctx, `SELECT audit_event_id,COALESCE(request_id,''),COALESCE(decision_id,''),COALESCE(workflow_type,''),COALESCE(node_name,''),actor,action,COALESCE(node_action,''),COALESCE(proposal_id,''),COALESCE(confirmation_id,''),COALESCE(error_case_id,''),status,COALESCE(error_code,''),COALESCE(before_state,''),COALESCE(after_state,''),COALESCE(rule_version,''),COALESCE(snapshot_id,''),COALESCE(input_ref_type,''),COALESCE(input_ref,''),COALESCE(output_ref_type,''),COALESCE(output_ref,''),created_at FROM audit_events WHERE audit_event_id=?`, id).Scan(&e.AuditEventID, &e.RequestID, &e.DecisionID, &e.WorkflowType, &e.NodeName, &e.Actor, &e.Action, &e.NodeAction, &e.ProposalID, &e.ConfirmationID, &e.ErrorCaseID, &e.Status, &e.ErrorCode, &e.BeforeState, &e.AfterState, &e.RuleVersion, &e.SnapshotID, &e.InputRefType, &e.InputRef, &e.OutputRefType, &e.OutputRef, &e.CreatedAt)
	return e, apperr.FromRepositoryError(err)
}

// ListAuditEvents reads the local audit timeline.
func (r *AuditRepository) ListAuditEvents(ctx context.Context) ([]repository.AuditEvent, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT audit_event_id,COALESCE(request_id,''),COALESCE(decision_id,''),COALESCE(workflow_type,''),COALESCE(node_name,''),actor,action,COALESCE(node_action,''),COALESCE(proposal_id,''),COALESCE(confirmation_id,''),COALESCE(error_case_id,''),status,COALESCE(error_code,''),COALESCE(before_state,''),COALESCE(after_state,''),COALESCE(rule_version,''),COALESCE(snapshot_id,''),COALESCE(input_ref_type,''),COALESCE(input_ref,''),COALESCE(output_ref_type,''),COALESCE(output_ref,''),created_at FROM audit_events ORDER BY created_at DESC`)
	if err != nil {
		return nil, apperr.FromRepositoryError(err)
	}
	defer rows.Close()
	var out []repository.AuditEvent
	for rows.Next() {
		var e repository.AuditEvent
		if err := rows.Scan(&e.AuditEventID, &e.RequestID, &e.DecisionID, &e.WorkflowType, &e.NodeName, &e.Actor, &e.Action, &e.NodeAction, &e.ProposalID, &e.ConfirmationID, &e.ErrorCaseID, &e.Status, &e.ErrorCode, &e.BeforeState, &e.AfterState, &e.RuleVersion, &e.SnapshotID, &e.InputRefType, &e.InputRef, &e.OutputRefType, &e.OutputRef, &e.CreatedAt); err != nil {
			return nil, apperr.FromRepositoryError(err)
		}
		out = append(out, e)
	}
	return out, apperr.FromRepositoryError(rows.Err())
}
