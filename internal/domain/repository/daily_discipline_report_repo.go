package repository

import "context"

// DailyDisciplineReport 持久化每日纪律检查报告的最新生成结果。
type DailyDisciplineReport struct {
	ReportID      string
	LocalDate     string
	Scope         string
	SymbolSetHash string
	SourceType    string
	SourceID      string
	DecisionID    string
	Status        string
	Summary       string
	FailureCode   string
	FailureReason string
	CreatedAt     string
	UpdatedAt     string
}

// DailyDisciplineReportListFilter 过滤每日纪律检查报告列表。
type DailyDisciplineReportListFilter struct {
	Status string
	Limit  int
}

// DailyDisciplineReportRepository 持久化每日纪律检查报告。
type DailyDisciplineReportRepository interface {
	UpsertDailyDisciplineReport(ctx context.Context, report DailyDisciplineReport) error
	GetDailyDisciplineReport(ctx context.Context, reportID string) (DailyDisciplineReport, error)
	GetDailyDisciplineReportByKey(ctx context.Context, localDate, scope, symbolSetHash string) (DailyDisciplineReport, error)
	ListDailyDisciplineReports(ctx context.Context, filter DailyDisciplineReportListFilter) ([]DailyDisciplineReport, error)
}
