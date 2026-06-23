package workflow

import (
	"testing"
	"time"
)

func TestEvidenceInRangeUsesExclusiveEndBoundary(t *testing.T) {
	start := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 6, 6, 0, 0, 0, 0, time.UTC)
	if !evidenceInRange("2024-06-06 23:59:59", start, end) {
		t.Fatal("expected end date to include the whole day")
	}
	if evidenceInRange("2024-06-07 00:00:00", start, end) {
		t.Fatal("expected next day midnight to be excluded")
	}
}

func TestEvidenceInRangeRejectsUnparseableDates(t *testing.T) {
	start := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 6, 6, 0, 0, 0, 0, time.UTC)
	if evidenceInRange("not-a-date", start, end) {
		t.Fatal("expected unparseable date to be rejected")
	}
	if !evidenceInRange("2024年06月05日", start, end) {
		t.Fatal("expected Chinese date format to be accepted")
	}
}
