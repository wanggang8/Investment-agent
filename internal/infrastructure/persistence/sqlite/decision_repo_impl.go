package sqlite

import (
	"context"
	"database/sql"

	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

// DecisionRepository 是决策记录、证据引用和用户确认表的 SQLite 实现。
type DecisionRepository struct{ db dbtx }

// NewDecisionRepository 创建决策仓储实例。
func NewDecisionRepository(db *sql.DB) *DecisionRepository { return &DecisionRepository{db: db} }

// SaveDecisionRecord 在同一事务中保存决策记录与证据引用。
func (r *DecisionRepository) SaveDecisionRecord(ctx context.Context, d repository.DecisionRecord, refs []repository.EvidenceRef) error {
	err := withTx(ctx, r.db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO decision_records (decision_id,request_id,workflow_type,symbol,question,workflow_status,record_type,dashboard_state,capability_status,capability_reason,source_verification_status,risk_reason_code,media_heat_summary_json,user_emotion_tags_json,triggered_rules_json,errors_json,final_verdict_status,final_verdict_text,prohibited_actions_json,optional_actions_json,confirmation_status,portfolio_snapshot_id,market_snapshot_id,rule_version,analyst_reports_json,expected_return_scenarios_json,arbitration_chain_json,context_snapshot_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, d.DecisionID, d.RequestID, d.WorkflowType, nullString(d.Symbol), nullString(d.Question), d.WorkflowStatus, d.RecordType, d.DashboardState, nullString(d.CapabilityStatus), nullString(d.CapabilityReason), nullString(d.SourceVerificationStatus), nullString(d.RiskReasonCode), nullString(d.MediaHeatSummaryJSON), nullString(d.UserEmotionTagsJSON), nullString(d.TriggeredRulesJSON), nullString(d.ErrorsJSON), d.FinalVerdictStatus, d.FinalVerdictText, nullString(d.ProhibitedActionsJSON), nullString(d.OptionalActionsJSON), d.ConfirmationStatus, nullString(d.PortfolioSnapshotID), nullString(d.MarketSnapshotID), d.RuleVersion, nullString(d.AnalystReportsJSON), nullString(d.ExpectedReturnScenariosJSON), nullString(d.ArbitrationChainJSON), nullString(d.ContextSnapshotJSON), d.CreatedAt)
		if err != nil {
			return err
		}
		for _, e := range refs {
			_, err = tx.ExecContext(ctx, `INSERT INTO evidence_refs (evidence_ref_id,evidence_id,decision_id,summary_id,source_name,source_level,evidence_role,published_at,captured_at,original_url,summary,content_hash,time_weight,relevance_score,independent_source_count,high_grade_independent_source_count,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, e.EvidenceRefID, e.EvidenceID, e.DecisionID, e.SummaryID, e.SourceName, e.SourceLevel, e.EvidenceRole, nullString(e.PublishedAt), nullString(e.CapturedAt), nullString(e.OriginalURL), e.Summary, nullString(e.ContentHash), e.TimeWeight, e.RelevanceScore, e.IndependentSourceCount, e.HighGradeIndependentSourceCount, e.CreatedAt)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return apperr.FromRepositoryError(err)
}

// GetDecisionRecord 读取决策记录及其证据引用快照。
func (r *DecisionRepository) GetDecisionRecord(ctx context.Context, id string) (repository.DecisionRecord, []repository.EvidenceRef, error) {
	var d repository.DecisionRecord
	err := r.db.QueryRowContext(ctx, `SELECT decision_id,request_id,workflow_type,COALESCE(symbol,''),COALESCE(question,''),workflow_status,record_type,dashboard_state,COALESCE(capability_status,''),COALESCE(capability_reason,''),COALESCE(source_verification_status,''),COALESCE(risk_reason_code,''),COALESCE(media_heat_summary_json,''),COALESCE(user_emotion_tags_json,''),COALESCE(triggered_rules_json,''),COALESCE(errors_json,''),final_verdict_status,final_verdict_text,COALESCE(prohibited_actions_json,''),COALESCE(optional_actions_json,''),confirmation_status,COALESCE(portfolio_snapshot_id,''),COALESCE(market_snapshot_id,''),rule_version,COALESCE(analyst_reports_json,''),COALESCE(expected_return_scenarios_json,''),COALESCE(arbitration_chain_json,''),COALESCE(context_snapshot_json,''),created_at FROM decision_records WHERE decision_id=?`, id).Scan(&d.DecisionID, &d.RequestID, &d.WorkflowType, &d.Symbol, &d.Question, &d.WorkflowStatus, &d.RecordType, &d.DashboardState, &d.CapabilityStatus, &d.CapabilityReason, &d.SourceVerificationStatus, &d.RiskReasonCode, &d.MediaHeatSummaryJSON, &d.UserEmotionTagsJSON, &d.TriggeredRulesJSON, &d.ErrorsJSON, &d.FinalVerdictStatus, &d.FinalVerdictText, &d.ProhibitedActionsJSON, &d.OptionalActionsJSON, &d.ConfirmationStatus, &d.PortfolioSnapshotID, &d.MarketSnapshotID, &d.RuleVersion, &d.AnalystReportsJSON, &d.ExpectedReturnScenariosJSON, &d.ArbitrationChainJSON, &d.ContextSnapshotJSON, &d.CreatedAt)
	if err != nil {
		return d, nil, apperr.FromRepositoryError(err)
	}
	rows, err := r.db.QueryContext(ctx, `SELECT evidence_ref_id,evidence_id,decision_id,summary_id,source_name,source_level,evidence_role,COALESCE(published_at,''),COALESCE(captured_at,''),COALESCE(original_url,''),summary,COALESCE(content_hash,''),time_weight,relevance_score,COALESCE(independent_source_count,0),COALESCE(high_grade_independent_source_count,0),created_at FROM evidence_refs WHERE decision_id=? ORDER BY evidence_ref_id`, id)
	if err != nil {
		return d, nil, apperr.FromRepositoryError(err)
	}
	defer rows.Close()
	var refs []repository.EvidenceRef
	for rows.Next() {
		var e repository.EvidenceRef
		if err := rows.Scan(&e.EvidenceRefID, &e.EvidenceID, &e.DecisionID, &e.SummaryID, &e.SourceName, &e.SourceLevel, &e.EvidenceRole, &e.PublishedAt, &e.CapturedAt, &e.OriginalURL, &e.Summary, &e.ContentHash, &e.TimeWeight, &e.RelevanceScore, &e.IndependentSourceCount, &e.HighGradeIndependentSourceCount, &e.CreatedAt); err != nil {
			return d, nil, apperr.FromRepositoryError(err)
		}
		refs = append(refs, e)
	}
	return d, refs, apperr.FromRepositoryError(rows.Err())
}

// ListDecisionRecords reads decision records for list pages and review aggregation.
func (r *DecisionRepository) ListDecisionRecords(ctx context.Context) ([]repository.DecisionRecord, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT decision_id,COALESCE(symbol,''),workflow_status,COALESCE(source_verification_status,''),final_verdict_status,COALESCE(triggered_rules_json,'[]'),confirmation_status,created_at FROM decision_records ORDER BY created_at DESC`)
	if err != nil {
		return nil, apperr.FromRepositoryError(err)
	}
	defer rows.Close()
	var out []repository.DecisionRecord
	for rows.Next() {
		var d repository.DecisionRecord
		if err := rows.Scan(&d.DecisionID, &d.Symbol, &d.WorkflowStatus, &d.SourceVerificationStatus, &d.FinalVerdictStatus, &d.TriggeredRulesJSON, &d.ConfirmationStatus, &d.CreatedAt); err != nil {
			return nil, apperr.FromRepositoryError(err)
		}
		out = append(out, d)
	}
	return out, apperr.FromRepositoryError(rows.Err())
}

// ListErrorCases reads marked error cases for review aggregation.
func (r *DecisionRepository) ListErrorCases(ctx context.Context) ([]repository.ErrorCase, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT error_case_id,decision_id,confirmation_id,COALESCE(actual_outcome,''),COALESCE(root_cause_tag,''),COALESCE(lesson_learned,''),created_at FROM error_cases ORDER BY created_at DESC`)
	if err != nil {
		return nil, apperr.FromRepositoryError(err)
	}
	defer rows.Close()
	out := []repository.ErrorCase{}
	for rows.Next() {
		var item repository.ErrorCase
		if err := rows.Scan(&item.ErrorCaseID, &item.DecisionID, &item.ConfirmationID, &item.ActualOutcome, &item.RootCauseTag, &item.LessonLearned, &item.CreatedAt); err != nil {
			return nil, apperr.FromRepositoryError(err)
		}
		out = append(out, item)
	}
	return out, apperr.FromRepositoryError(rows.Err())
}

// CountErrorCases counts marked error cases.
func (r *DecisionRepository) CountErrorCases(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM error_cases`).Scan(&count)
	return count, apperr.FromRepositoryError(err)
}

// GetDecisionConfirmationState reads fields needed before confirmation writes.
func (r *DecisionRepository) GetDecisionConfirmationState(ctx context.Context, decisionID string) (string, string, error) {
	var recordType, currentStatus string
	err := r.db.QueryRowContext(ctx, `SELECT record_type,confirmation_status FROM decision_records WHERE decision_id=?`, decisionID).Scan(&recordType, &currentStatus)
	return recordType, currentStatus, apperr.FromRepositoryError(err)
}

// SaveOperationConfirmation 保存用户对建议的线下处理结果。
func (r *DecisionRepository) SaveOperationConfirmation(ctx context.Context, c repository.OperationConfirmation) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO operation_confirmations (confirmation_id,decision_id,confirmation_type,operation_type,symbol,quantity,price,fees,executed_at,error_case_id,payload_json,note,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`, c.ConfirmationID, c.DecisionID, c.ConfirmationType, nullString(c.OperationType), nullString(c.Symbol), c.Quantity, c.Price, c.Fees, nullString(c.ExecutedAt), nullString(c.ErrorCaseID), nullString(c.PayloadJSON), nullString(c.Note), c.CreatedAt)
	return apperr.FromRepositoryError(err)
}

// ListOperationConfirmations reads local user confirmation facts for one decision.
func (r *DecisionRepository) ListOperationConfirmations(ctx context.Context, decisionID string) ([]repository.OperationConfirmation, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT confirmation_id,decision_id,confirmation_type,COALESCE(operation_type,''),COALESCE(symbol,''),COALESCE(quantity,0),COALESCE(price,0),COALESCE(fees,0),COALESCE(executed_at,''),COALESCE(error_case_id,''),'',COALESCE(note,''),created_at FROM operation_confirmations WHERE decision_id=? ORDER BY created_at ASC, confirmation_id ASC`, decisionID)
	if err != nil {
		return nil, apperr.FromRepositoryError(err)
	}
	defer rows.Close()
	var out []repository.OperationConfirmation
	for rows.Next() {
		var c repository.OperationConfirmation
		if err := rows.Scan(&c.ConfirmationID, &c.DecisionID, &c.ConfirmationType, &c.OperationType, &c.Symbol, &c.Quantity, &c.Price, &c.Fees, &c.ExecutedAt, &c.ErrorCaseID, &c.PayloadJSON, &c.Note, &c.CreatedAt); err != nil {
			return nil, apperr.FromRepositoryError(err)
		}
		out = append(out, c)
	}
	return out, apperr.FromRepositoryError(rows.Err())
}

// UpdateDecisionConfirmationStatus updates a decision confirmation state.
func (r *DecisionRepository) UpdateDecisionConfirmationStatus(ctx context.Context, decisionID, status string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE decision_records SET confirmation_status=? WHERE decision_id=?`, status, decisionID)
	return apperr.FromRepositoryError(err)
}

// UpdateDecisionConfirmationStatusIfCurrent updates a confirmation state only when the stored state still matches.
func (r *DecisionRepository) UpdateDecisionConfirmationStatusIfCurrent(ctx context.Context, decisionID, expectedStatus, nextStatus string) (bool, error) {
	result, err := r.db.ExecContext(ctx, `UPDATE decision_records SET confirmation_status=? WHERE decision_id=? AND confirmation_status=?`, nextStatus, decisionID, expectedStatus)
	if err != nil {
		return false, apperr.FromRepositoryError(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, apperr.FromRepositoryError(err)
	}
	return rows > 0, nil
}

// SavePositionTransaction records a manually executed position operation.
func (r *DecisionRepository) SavePositionTransaction(ctx context.Context, tx repository.PositionTransaction) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO position_transactions (transaction_id,confirmation_id,symbol,operation_type,quantity,price,fees,occurred_at,before_position_json,after_position_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?)`, tx.TransactionID, tx.ConfirmationID, tx.Symbol, tx.OperationType, tx.Quantity, tx.Price, tx.Fees, tx.OccurredAt, nullString(tx.BeforePositionJSON), nullString(tx.AfterPositionJSON), tx.CreatedAt)
	return apperr.FromRepositoryError(err)
}

// ListPositionTransactionsByConfirmation reads local position transaction facts for one confirmation.
func (r *DecisionRepository) ListPositionTransactionsByConfirmation(ctx context.Context, confirmationID string) ([]repository.PositionTransaction, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT transaction_id,confirmation_id,symbol,operation_type,quantity,price,COALESCE(fees,0),occurred_at,'','',created_at FROM position_transactions WHERE confirmation_id=? ORDER BY occurred_at ASC, transaction_id ASC`, confirmationID)
	if err != nil {
		return nil, apperr.FromRepositoryError(err)
	}
	defer rows.Close()
	var out []repository.PositionTransaction
	for rows.Next() {
		var tx repository.PositionTransaction
		if err := rows.Scan(&tx.TransactionID, &tx.ConfirmationID, &tx.Symbol, &tx.OperationType, &tx.Quantity, &tx.Price, &tx.Fees, &tx.OccurredAt, &tx.BeforePositionJSON, &tx.AfterPositionJSON, &tx.CreatedAt); err != nil {
			return nil, apperr.FromRepositoryError(err)
		}
		out = append(out, tx)
	}
	return out, apperr.FromRepositoryError(rows.Err())
}

// SaveErrorCase records a user-marked error case.
func (r *DecisionRepository) SaveErrorCase(ctx context.Context, e repository.ErrorCase) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO error_cases (error_case_id,decision_id,confirmation_id,actual_outcome,root_cause_tag,lesson_learned,created_at) VALUES (?,?,?,?,?,?,?)`, e.ErrorCaseID, e.DecisionID, e.ConfirmationID, nullString(e.ActualOutcome), nullString(e.RootCauseTag), nullString(e.LessonLearned), e.CreatedAt)
	return apperr.FromRepositoryError(err)
}

// GetOperationConfirmation 读取一条用户确认记录。
func (r *DecisionRepository) GetOperationConfirmation(ctx context.Context, id string) (repository.OperationConfirmation, error) {
	var c repository.OperationConfirmation
	err := r.db.QueryRowContext(ctx, `SELECT confirmation_id,decision_id,confirmation_type,COALESCE(operation_type,''),COALESCE(symbol,''),COALESCE(quantity,0),COALESCE(price,0),COALESCE(fees,0),COALESCE(executed_at,''),COALESCE(error_case_id,''),COALESCE(payload_json,''),COALESCE(note,''),created_at FROM operation_confirmations WHERE confirmation_id=?`, id).Scan(&c.ConfirmationID, &c.DecisionID, &c.ConfirmationType, &c.OperationType, &c.Symbol, &c.Quantity, &c.Price, &c.Fees, &c.ExecutedAt, &c.ErrorCaseID, &c.PayloadJSON, &c.Note, &c.CreatedAt)
	return c, apperr.FromRepositoryError(err)
}
