package sqlite

import (
	"context"
	"database/sql"

	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

// RuleRepository 是规则版本、规则提案和守门人审计表的 SQLite 实现。
type RuleRepository struct{ db dbtx }

// NewRuleRepository 创建规则仓储实例。
func NewRuleRepository(db *sql.DB) *RuleRepository { return &RuleRepository{db: db} }

// SaveRuleVersion 保存一份正式规则版本快照。
func (r *RuleRepository) SaveRuleVersion(ctx context.Context, v repository.RuleVersion) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO rule_versions (rule_version,status,rules_json,effective_at,created_from_proposal_id,created_at) VALUES (?,?,?,?,?,?)`, v.RuleVersion, v.Status, v.RulesJSON, v.EffectiveAt, nullString(v.CreatedFromProposalID), v.CreatedAt)
	return apperr.FromRepositoryError(err)
}

// GetRuleVersion 按版本号读取规则快照。
func (r *RuleRepository) GetRuleVersion(ctx context.Context, id string) (repository.RuleVersion, error) {
	var v repository.RuleVersion
	err := r.db.QueryRowContext(ctx, `SELECT rule_version,status,rules_json,effective_at,COALESCE(created_from_proposal_id,''),created_at FROM rule_versions WHERE rule_version=?`, id).Scan(&v.RuleVersion, &v.Status, &v.RulesJSON, &v.EffectiveAt, &v.CreatedFromProposalID, &v.CreatedAt)
	return v, apperr.FromRepositoryError(err)
}

// GetActiveRuleVersion reads the active rule version.
func (r *RuleRepository) GetActiveRuleVersion(ctx context.Context) (repository.RuleVersion, error) {
	var v repository.RuleVersion
	err := r.db.QueryRowContext(ctx, `SELECT rule_version,status,rules_json,effective_at,COALESCE(created_from_proposal_id,''),created_at FROM rule_versions WHERE status='active' LIMIT 1`).Scan(&v.RuleVersion, &v.Status, &v.RulesJSON, &v.EffectiveAt, &v.CreatedFromProposalID, &v.CreatedAt)
	return v, apperr.FromRepositoryError(err)
}

// SaveRuleProposal 保存规则提案，提案本身不会直接改变 active 规则版本。
func (r *RuleRepository) SaveRuleProposal(ctx context.Context, p repository.RuleProposal) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO rule_proposals (proposal_id,proposal_type,status,source_error_case_id,title,proposal_version,before_rule_json,after_rule_json,reason,impact_scope_json,risk_notes_json,sample_count,final_confirmed_at,final_confirmed_note,applied_rule_version,related_error_cases_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, p.ProposalID, p.ProposalType, p.Status, nullString(p.SourceErrorCaseID), p.Title, p.ProposalVersion, nullString(p.BeforeRuleJSON), nullString(p.AfterRuleJSON), nullString(p.Reason), nullString(p.ImpactScopeJSON), nullString(p.RiskNotesJSON), p.SampleCount, nullString(p.FinalConfirmedAt), nullString(p.FinalConfirmedNote), nullString(p.AppliedRuleVersion), nullString(p.RelatedErrorCasesJSON), p.CreatedAt)
	return apperr.FromRepositoryError(err)
}

// UpdateRuleProposalStatus 更新规则提案状态。
func (r *RuleRepository) UpdateRuleProposalStatus(ctx context.Context, proposalID string, status string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE rule_proposals SET status=? WHERE proposal_id=?`, status, proposalID)
	if err != nil {
		return apperr.FromRepositoryError(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return apperr.FromRepositoryError(err)
	}
	if rows == 0 {
		return apperr.FromRepositoryError(sql.ErrNoRows)
	}
	return nil
}

// ApplyRuleProposal records final confirmation metadata.
func (r *RuleRepository) ApplyRuleProposal(ctx context.Context, proposalID, status, finalConfirmedAt, finalConfirmedNote, appliedRuleVersion string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE rule_proposals SET status=?,final_confirmed_at=?,final_confirmed_note=?,applied_rule_version=? WHERE proposal_id=?`, status, finalConfirmedAt, nullString(finalConfirmedNote), appliedRuleVersion, proposalID)
	if err != nil {
		return apperr.FromRepositoryError(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return apperr.FromRepositoryError(err)
	}
	if rows == 0 {
		return apperr.FromRepositoryError(sql.ErrNoRows)
	}
	return nil
}

// ArchiveActiveRuleVersions archives currently active rule versions.
func (r *RuleRepository) ArchiveActiveRuleVersions(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `UPDATE rule_versions SET status='archived' WHERE status='active'`)
	return apperr.FromRepositoryError(err)
}

// GetRuleProposal 按提案 ID 读取规则提案。
func (r *RuleRepository) GetRuleProposal(ctx context.Context, id string) (repository.RuleProposal, error) {
	var p repository.RuleProposal
	err := r.db.QueryRowContext(ctx, `SELECT proposal_id,proposal_type,status,COALESCE(source_error_case_id,''),title,proposal_version,COALESCE(before_rule_json,''),COALESCE(after_rule_json,''),COALESCE(reason,''),COALESCE(impact_scope_json,''),COALESCE(risk_notes_json,''),sample_count,COALESCE(final_confirmed_at,''),COALESCE(final_confirmed_note,''),COALESCE(applied_rule_version,''),COALESCE(related_error_cases_json,''),created_at FROM rule_proposals WHERE proposal_id=?`, id).Scan(&p.ProposalID, &p.ProposalType, &p.Status, &p.SourceErrorCaseID, &p.Title, &p.ProposalVersion, &p.BeforeRuleJSON, &p.AfterRuleJSON, &p.Reason, &p.ImpactScopeJSON, &p.RiskNotesJSON, &p.SampleCount, &p.FinalConfirmedAt, &p.FinalConfirmedNote, &p.AppliedRuleVersion, &p.RelatedErrorCasesJSON, &p.CreatedAt)
	return p, apperr.FromRepositoryError(err)
}

// ListRuleProposals reads rule proposals with their latest gatekeeper audit summary.
func (r *RuleRepository) ListRuleProposals(ctx context.Context) ([]repository.RuleProposalWithAudit, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT rp.proposal_id,rp.proposal_type,rp.status,COALESCE(rp.source_error_case_id,''),rp.title,rp.proposal_version,COALESCE(rp.before_rule_json,'{}'),COALESCE(rp.after_rule_json,'{}'),COALESCE(rp.reason,''),COALESCE(rp.impact_scope_json,''),COALESCE(rp.risk_notes_json,''),rp.sample_count,COALESCE(rp.final_confirmed_at,''),COALESCE(rp.final_confirmed_note,''),COALESCE(rp.applied_rule_version,''),COALESCE(rp.related_error_cases_json,''),rp.created_at,COALESCE(ga.audit_result,''),COALESCE(ga.audit_reason,'') FROM rule_proposals rp LEFT JOIN gatekeeper_audits ga ON ga.gatekeeper_audit_id=(SELECT gatekeeper_audit_id FROM gatekeeper_audits WHERE proposal_id=rp.proposal_id ORDER BY created_at DESC LIMIT 1) ORDER BY rp.created_at DESC`)
	if err != nil {
		return nil, apperr.FromRepositoryError(err)
	}
	defer rows.Close()
	var out []repository.RuleProposalWithAudit
	for rows.Next() {
		var p repository.RuleProposalWithAudit
		if err := rows.Scan(&p.ProposalID, &p.ProposalType, &p.Status, &p.SourceErrorCaseID, &p.Title, &p.ProposalVersion, &p.BeforeRuleJSON, &p.AfterRuleJSON, &p.Reason, &p.ImpactScopeJSON, &p.RiskNotesJSON, &p.SampleCount, &p.FinalConfirmedAt, &p.FinalConfirmedNote, &p.AppliedRuleVersion, &p.RelatedErrorCasesJSON, &p.CreatedAt, &p.AuditResult, &p.AuditReason); err != nil {
			return nil, apperr.FromRepositoryError(err)
		}
		out = append(out, p)
	}
	return out, apperr.FromRepositoryError(rows.Err())
}

// SaveGatekeeperAudit 保存守门人审计结论。
func (r *RuleRepository) SaveGatekeeperAudit(ctx context.Context, a repository.GatekeeperAudit) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO gatekeeper_audits (gatekeeper_audit_id,proposal_id,audit_result,audit_reason,required_changes,violates_fundamental_rule,has_rule_conflict,backtest_metrics_json,allow_apply,audited_rule_version,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?)`, a.GatekeeperAuditID, a.ProposalID, a.AuditResult, nullString(a.AuditReason), nullString(a.RequiredChanges), boolInt(a.ViolatesFundamentalRule), boolInt(a.HasRuleConflict), nullString(a.BacktestMetricsJSON), boolInt(a.AllowApply), a.AuditedRuleVersion, a.CreatedAt)
	return apperr.FromRepositoryError(err)
}

// GetGatekeeperAudit 按审计 ID 读取守门人审计结论。
func (r *RuleRepository) GetGatekeeperAudit(ctx context.Context, id string) (repository.GatekeeperAudit, error) {
	var a repository.GatekeeperAudit
	var v, h, ap int
	err := r.db.QueryRowContext(ctx, `SELECT gatekeeper_audit_id,proposal_id,audit_result,COALESCE(audit_reason,''),COALESCE(required_changes,''),violates_fundamental_rule,has_rule_conflict,COALESCE(backtest_metrics_json,''),allow_apply,audited_rule_version,created_at FROM gatekeeper_audits WHERE gatekeeper_audit_id=?`, id).Scan(&a.GatekeeperAuditID, &a.ProposalID, &a.AuditResult, &a.AuditReason, &a.RequiredChanges, &v, &h, &a.BacktestMetricsJSON, &ap, &a.AuditedRuleVersion, &a.CreatedAt)
	a.ViolatesFundamentalRule = v == 1
	a.HasRuleConflict = h == 1
	a.AllowApply = ap == 1
	return a, apperr.FromRepositoryError(err)
}

// GetLatestGatekeeperAuditByProposal 读取指定提案最新守门人审计结论。
func (r *RuleRepository) GetLatestGatekeeperAuditByProposal(ctx context.Context, proposalID string) (repository.GatekeeperAudit, error) {
	var a repository.GatekeeperAudit
	var v, h, ap int
	err := r.db.QueryRowContext(ctx, `SELECT gatekeeper_audit_id,proposal_id,audit_result,COALESCE(audit_reason,''),COALESCE(required_changes,''),violates_fundamental_rule,has_rule_conflict,COALESCE(backtest_metrics_json,''),allow_apply,audited_rule_version,created_at FROM gatekeeper_audits WHERE proposal_id=? ORDER BY created_at DESC LIMIT 1`, proposalID).Scan(&a.GatekeeperAuditID, &a.ProposalID, &a.AuditResult, &a.AuditReason, &a.RequiredChanges, &v, &h, &a.BacktestMetricsJSON, &ap, &a.AuditedRuleVersion, &a.CreatedAt)
	a.ViolatesFundamentalRule = v == 1
	a.HasRuleConflict = h == 1
	a.AllowApply = ap == 1
	return a, apperr.FromRepositoryError(err)
}

// boolInt 将 Go bool 转为 SQLite 使用的 0/1。
func boolInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
