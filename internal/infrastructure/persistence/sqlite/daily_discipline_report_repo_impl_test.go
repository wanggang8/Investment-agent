package sqlite

import (
	"context"
	"testing"

	"investment-agent/internal/domain/repository"
)

func TestDailyDisciplineReportRepositoryUpsertAndGet(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	repo := NewDailyDisciplineReportRepository(db)

	report := repository.DailyDisciplineReport{
		ReportID:      "report_1",
		LocalDate:     "2026-06-08",
		Scope:         "holdings",
		SymbolSetHash: "hash_1",
		SourceType:    "auto_run",
		SourceID:      "run_1",
		DecisionID:    "decision_1",
		Status:        "running",
		Summary:       "生成中",
		CreatedAt:     testTime,
		UpdatedAt:     testTime,
	}
	if err := repo.UpsertDailyDisciplineReport(ctx, report); err != nil {
		t.Fatal(err)
	}

	report.Status = "success"
	report.Summary = "纪律报告已生成"
	report.UpdatedAt = "2026-06-08T00:31:00Z"
	if err := repo.UpsertDailyDisciplineReport(ctx, report); err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetDailyDisciplineReport(ctx, report.ReportID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ReportID != "report_1" || got.Status != "success" || got.Summary != "纪律报告已生成" || got.UpdatedAt != "2026-06-08T00:31:00Z" {
		t.Fatalf("unexpected report: %+v", got)
	}

	byKey, err := repo.GetDailyDisciplineReportByKey(ctx, report.LocalDate, report.Scope, report.SymbolSetHash)
	if err != nil {
		t.Fatal(err)
	}
	if byKey.ReportID != report.ReportID {
		t.Fatalf("expected report by key %s, got %+v", report.ReportID, byKey)
	}
}

func TestDailyDisciplineReportRepositoryListLatestNewestFirst(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	repo := NewDailyDisciplineReportRepository(db)

	reports := []repository.DailyDisciplineReport{
		{ReportID: "report_old", LocalDate: "2026-06-06", Scope: "holdings", SymbolSetHash: "hash_old", SourceType: "manual", SourceID: "manual_1", Status: "success", Summary: "old", CreatedAt: testTime, UpdatedAt: "2026-06-06T00:10:00Z"},
		{ReportID: "report_new_a", LocalDate: "2026-06-08", Scope: "holdings", SymbolSetHash: "hash_new_a", SourceType: "auto_run", SourceID: "run_1", Status: "failed", Summary: "new a", FailureCode: "DATA_SOURCE_UNAVAILABLE", FailureReason: "source down", CreatedAt: testTime, UpdatedAt: "2026-06-08T00:10:00Z"},
		{ReportID: "report_new_b", LocalDate: "2026-06-08", Scope: "holdings", SymbolSetHash: "hash_new_b", SourceType: "auto_run", SourceID: "run_2", Status: "success", Summary: "new b", CreatedAt: testTime, UpdatedAt: "2026-06-08T00:20:00Z"},
	}
	for _, report := range reports {
		if err := repo.UpsertDailyDisciplineReport(ctx, report); err != nil {
			t.Fatal(err)
		}
	}

	got, err := repo.ListDailyDisciplineReports(ctx, repository.DailyDisciplineReportListFilter{Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].ReportID != "report_new_b" || got[1].ReportID != "report_new_a" {
		t.Fatalf("expected newest first with limit, got %+v", got)
	}

	failed, err := repo.ListDailyDisciplineReports(ctx, repository.DailyDisciplineReportListFilter{Status: "failed"})
	if err != nil {
		t.Fatal(err)
	}
	if len(failed) != 1 || failed[0].ReportID != "report_new_a" || failed[0].FailureReason != "source down" {
		t.Fatalf("expected failed report filter, got %+v", failed)
	}
}

func TestDailyDisciplineReportRepositoryIdempotentKeyKeepsStableReportID(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	repo := NewDailyDisciplineReportRepository(db)

	first := repository.DailyDisciplineReport{ReportID: "report_first", LocalDate: "2026-06-08", Scope: "holdings", SymbolSetHash: "hash_same", SourceType: "auto_run", SourceID: "run_1", Status: "running", Summary: "first", CreatedAt: testTime, UpdatedAt: testTime}
	if err := repo.UpsertDailyDisciplineReport(ctx, first); err != nil {
		t.Fatal(err)
	}
	second := first
	second.ReportID = "report_second"
	second.SourceID = "run_2"
	second.Status = "success"
	second.Summary = "replacement"
	second.UpdatedAt = "2026-06-08T00:40:00Z"
	if err := repo.UpsertDailyDisciplineReport(ctx, second); err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetDailyDisciplineReportByKey(ctx, first.LocalDate, first.Scope, first.SymbolSetHash)
	if err != nil {
		t.Fatal(err)
	}
	if got.ReportID != first.ReportID || got.CreatedAt != first.CreatedAt {
		t.Fatalf("expected stable report identity from first upsert, got %+v", got)
	}
	if got.Status != "success" || got.Summary != "replacement" || got.SourceID != "run_2" || got.UpdatedAt != "2026-06-08T00:40:00Z" {
		t.Fatalf("expected mutable report fields from second upsert, got %+v", got)
	}
}
