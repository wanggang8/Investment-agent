package workflow

import (
	"net/url"
	"time"
)

const publicEvidenceDefaultBackfillDays = 90
const publicEvidencePageSize = 50

func evidenceDateRange(start, end time.Time) (time.Time, time.Time) {
	if end.IsZero() {
		end = time.Now().UTC()
	}
	if start.IsZero() {
		start = end.AddDate(0, 0, -publicEvidenceDefaultBackfillDays)
	}
	return start, end
}

func evidenceInRange(publishedAt string, start, end time.Time) bool {
	if publishedAt == "" {
		return false
	}
	parsed, err := parseEvidenceTime(publishedAt)
	if err != nil {
		return false
	}
	if !start.IsZero() && parsed.Before(start) {
		return false
	}
	if !end.IsZero() && !parsed.Before(evidenceEndExclusive(end)) {
		return false
	}
	return true
}

func evidenceEndExclusive(end time.Time) time.Time {
	if end.Hour() == 0 && end.Minute() == 0 && end.Second() == 0 && end.Nanosecond() == 0 {
		return end.AddDate(0, 0, 1)
	}
	return end
}

func parseEvidenceTime(value string) (time.Time, error) {
	layouts := []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02", "2006年01月02日"}
	var lastErr error
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed, nil
		}
		lastErr = err
	}
	return time.Time{}, lastErr
}

func resolveEvidenceURL(base, ref string) string {
	if ref == "" {
		return ""
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return ref
	}
	refURL, err := url.Parse(ref)
	if err != nil {
		return ref
	}
	return baseURL.ResolveReference(refURL).String()
}

func totalPages(total, pageSize int) int {
	if total <= 0 || pageSize <= 0 {
		return 1
	}
	pages := total / pageSize
	if total%pageSize != 0 {
		pages++
	}
	if pages < 1 {
		return 1
	}
	return pages
}
