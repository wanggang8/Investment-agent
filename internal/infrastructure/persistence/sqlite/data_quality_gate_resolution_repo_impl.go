package sqlite

import (
	"context"
	"database/sql"
	"strings"

	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

// DataQualityGateResolutionRepository is the SQLite implementation for current-data gate resolutions.
type DataQualityGateResolutionRepository struct{ db dbtx }

// NewDataQualityGateResolutionRepository creates a resolution repository.
func NewDataQualityGateResolutionRepository(db *sql.DB) *DataQualityGateResolutionRepository {
	return &DataQualityGateResolutionRepository{db: db}
}

func (r *DataQualityGateResolutionRepository) CreateDataQualityGateResolution(ctx context.Context, resolution repository.DataQualityGateResolution) error {
	if err := validateDataQualityGateResolution(resolution); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx, `INSERT INTO data_quality_gate_resolutions (resolution_id,symbol,policy_fingerprint,policy_verdict,release_gate,policy_summary,resolution_type,status,scope,reason,release_impact,evidence_ref,blocking_reasons_json,waiver_reasons_json,created_by,retired_by,created_at,retired_at,safety_note) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		resolution.ResolutionID,
		strings.TrimSpace(resolution.Symbol),
		strings.TrimSpace(resolution.PolicyFingerprint),
		strings.TrimSpace(resolution.PolicyVerdict),
		strings.TrimSpace(resolution.ReleaseGate),
		resolution.PolicySummary,
		strings.TrimSpace(resolution.ResolutionType),
		strings.TrimSpace(resolution.Status),
		resolution.Scope,
		resolution.Reason,
		resolution.ReleaseImpact,
		nullString(resolution.EvidenceRef),
		nullString(resolution.BlockingReasonsJSON),
		nullString(resolution.WaiverReasonsJSON),
		firstNonEmptySQLite(resolution.CreatedBy, "local_user"),
		nullString(resolution.RetiredBy),
		resolution.CreatedAt,
		nullString(resolution.RetiredAt),
		resolution.SafetyNote,
	)
	return apperr.FromRepositoryError(err)
}

func (r *DataQualityGateResolutionRepository) GetDataQualityGateResolution(ctx context.Context, resolutionID string) (repository.DataQualityGateResolution, error) {
	row := r.db.QueryRowContext(ctx, dataQualityGateResolutionSelectSQL()+` WHERE resolution_id=?`, strings.TrimSpace(resolutionID))
	return scanDataQualityGateResolution(row)
}

func (r *DataQualityGateResolutionRepository) GetActiveDataQualityGateResolution(ctx context.Context, symbol, policyFingerprint string) (repository.DataQualityGateResolution, error) {
	row := r.db.QueryRowContext(ctx, dataQualityGateResolutionSelectSQL()+` WHERE symbol=? AND policy_fingerprint=? AND status='active'`, strings.TrimSpace(symbol), strings.TrimSpace(policyFingerprint))
	return scanDataQualityGateResolution(row)
}

func (r *DataQualityGateResolutionRepository) ListDataQualityGateResolutions(ctx context.Context, filter repository.DataQualityGateResolutionFilter) ([]repository.DataQualityGateResolution, error) {
	query := dataQualityGateResolutionSelectSQL()
	var args []any
	var clauses []string
	if strings.TrimSpace(filter.Symbol) != "" {
		clauses = append(clauses, "symbol=?")
		args = append(args, strings.TrimSpace(filter.Symbol))
	}
	if strings.TrimSpace(filter.Status) != "" {
		clauses = append(clauses, "status=?")
		args = append(args, strings.TrimSpace(filter.Status))
	}
	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}
	query += " ORDER BY created_at DESC, resolution_id DESC"
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperr.FromRepositoryError(err)
	}
	defer rows.Close()
	var items []repository.DataQualityGateResolution
	for rows.Next() {
		item, err := scanDataQualityGateResolution(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, apperr.FromRepositoryError(rows.Err())
}

func (r *DataQualityGateResolutionRepository) RetireDataQualityGateResolution(ctx context.Context, resolutionID, retiredBy, retiredAt string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE data_quality_gate_resolutions SET status='retired', retired_by=?, retired_at=? WHERE resolution_id=? AND status='active'`, firstNonEmptySQLite(retiredBy, "local_user"), retiredAt, strings.TrimSpace(resolutionID))
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

func validateDataQualityGateResolution(resolution repository.DataQualityGateResolution) error {
	if strings.TrimSpace(resolution.ResolutionID) == "" ||
		strings.TrimSpace(resolution.Symbol) == "" ||
		strings.TrimSpace(resolution.PolicyFingerprint) == "" ||
		strings.TrimSpace(resolution.PolicyVerdict) == "" ||
		strings.TrimSpace(resolution.ReleaseGate) == "" ||
		strings.TrimSpace(resolution.ResolutionType) == "" ||
		strings.TrimSpace(resolution.Status) == "" ||
		strings.TrimSpace(resolution.Scope) == "" ||
		strings.TrimSpace(resolution.Reason) == "" ||
		strings.TrimSpace(resolution.ReleaseImpact) == "" ||
		strings.TrimSpace(resolution.CreatedAt) == "" ||
		strings.TrimSpace(resolution.SafetyNote) == "" {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "data quality gate resolution requires policy, resolution and audit fields")
	}
	return nil
}

func dataQualityGateResolutionSelectSQL() string {
	return `SELECT resolution_id,symbol,policy_fingerprint,policy_verdict,release_gate,policy_summary,resolution_type,status,scope,reason,release_impact,COALESCE(evidence_ref,''),COALESCE(blocking_reasons_json,''),COALESCE(waiver_reasons_json,''),created_by,COALESCE(retired_by,''),created_at,COALESCE(retired_at,''),safety_note FROM data_quality_gate_resolutions`
}

type dataQualityGateResolutionScanner interface {
	Scan(dest ...any) error
}

func scanDataQualityGateResolution(scanner dataQualityGateResolutionScanner) (repository.DataQualityGateResolution, error) {
	var item repository.DataQualityGateResolution
	if err := scanner.Scan(&item.ResolutionID, &item.Symbol, &item.PolicyFingerprint, &item.PolicyVerdict, &item.ReleaseGate, &item.PolicySummary, &item.ResolutionType, &item.Status, &item.Scope, &item.Reason, &item.ReleaseImpact, &item.EvidenceRef, &item.BlockingReasonsJSON, &item.WaiverReasonsJSON, &item.CreatedBy, &item.RetiredBy, &item.CreatedAt, &item.RetiredAt, &item.SafetyNote); err != nil {
		return repository.DataQualityGateResolution{}, apperr.FromRepositoryError(err)
	}
	return item, nil
}

func firstNonEmptySQLite(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
