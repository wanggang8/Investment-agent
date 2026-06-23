package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

// DailyDisciplineReportRepository 是每日纪律检查报告的 SQLite 实现。
type DailyDisciplineReportRepository struct{ db dbtx }

// NewDailyDisciplineReportRepository 创建每日纪律检查报告仓储。
func NewDailyDisciplineReportRepository(db *sql.DB) *DailyDisciplineReportRepository {
	return &DailyDisciplineReportRepository{db: db}
}

func (r *DailyDisciplineReportRepository) UpsertDailyDisciplineReport(ctx context.Context, report repository.DailyDisciplineReport) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO daily_discipline_reports (report_id,local_date,scope,symbol_set_hash,source_type,source_id,decision_id,status,summary,failure_code,failure_reason,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?) ON CONFLICT(local_date,scope,symbol_set_hash) DO UPDATE SET source_type=excluded.source_type,source_id=excluded.source_id,decision_id=excluded.decision_id,status=excluded.status,summary=excluded.summary,failure_code=excluded.failure_code,failure_reason=excluded.failure_reason,updated_at=excluded.updated_at`, report.ReportID, report.LocalDate, report.Scope, report.SymbolSetHash, report.SourceType, nullString(report.SourceID), nullString(report.DecisionID), report.Status, nullString(report.Summary), nullString(report.FailureCode), nullString(report.FailureReason), report.CreatedAt, report.UpdatedAt)
	return apperr.FromRepositoryError(err)
}

func (r *DailyDisciplineReportRepository) GetDailyDisciplineReport(ctx context.Context, reportID string) (repository.DailyDisciplineReport, error) {
	return r.getOne(ctx, `SELECT report_id,local_date,scope,symbol_set_hash,source_type,COALESCE(source_id,''),COALESCE(decision_id,''),status,COALESCE(summary,''),COALESCE(failure_code,''),COALESCE(failure_reason,''),created_at,updated_at FROM daily_discipline_reports WHERE report_id=?`, reportID)
}

func (r *DailyDisciplineReportRepository) GetDailyDisciplineReportByKey(ctx context.Context, localDate, scope, symbolSetHash string) (repository.DailyDisciplineReport, error) {
	return r.getOne(ctx, `SELECT report_id,local_date,scope,symbol_set_hash,source_type,COALESCE(source_id,''),COALESCE(decision_id,''),status,COALESCE(summary,''),COALESCE(failure_code,''),COALESCE(failure_reason,''),created_at,updated_at FROM daily_discipline_reports WHERE local_date=? AND scope=? AND symbol_set_hash=?`, localDate, scope, symbolSetHash)
}

func (r *DailyDisciplineReportRepository) ListDailyDisciplineReports(ctx context.Context, filter repository.DailyDisciplineReportListFilter) ([]repository.DailyDisciplineReport, error) {
	query := `SELECT report_id,local_date,scope,symbol_set_hash,source_type,COALESCE(source_id,''),COALESCE(decision_id,''),status,COALESCE(summary,''),COALESCE(failure_code,''),COALESCE(failure_reason,''),created_at,updated_at FROM daily_discipline_reports`
	args := []any{}
	if filter.Status != "" {
		query += ` WHERE status=?`
		args = append(args, filter.Status)
	}
	query += ` ORDER BY local_date DESC, updated_at DESC`
	if filter.Limit > 0 {
		query += ` LIMIT ?`
		args = append(args, filter.Limit)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperr.FromRepositoryError(err)
	}
	defer rows.Close()

	reports := []repository.DailyDisciplineReport{}
	for rows.Next() {
		var report repository.DailyDisciplineReport
		if err := scanDailyDisciplineReport(rows, &report); err != nil {
			return nil, apperr.FromRepositoryError(err)
		}
		reports = append(reports, report)
	}
	if err := rows.Err(); err != nil {
		return nil, apperr.FromRepositoryError(err)
	}
	return reports, nil
}

func (r *DailyDisciplineReportRepository) getOne(ctx context.Context, query string, args ...any) (repository.DailyDisciplineReport, error) {
	var report repository.DailyDisciplineReport
	err := scanDailyDisciplineReport(r.db.QueryRowContext(ctx, query, args...), &report)
	return report, apperr.FromRepositoryError(err)
}

type dailyDisciplineReportScanner interface {
	Scan(dest ...any) error
}

func scanDailyDisciplineReport(scanner dailyDisciplineReportScanner, report *repository.DailyDisciplineReport) error {
	if report == nil {
		return fmt.Errorf("daily discipline report destination is nil")
	}
	return scanner.Scan(&report.ReportID, &report.LocalDate, &report.Scope, &report.SymbolSetHash, &report.SourceType, &report.SourceID, &report.DecisionID, &report.Status, &report.Summary, &report.FailureCode, &report.FailureReason, &report.CreatedAt, &report.UpdatedAt)
}
