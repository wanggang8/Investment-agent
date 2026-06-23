package workflow

import "strings"

func p52CollectorFailureCategory(value string) string {
	switch strings.TrimSpace(value) {
	case "no_data":
		return "no_data"
	case "parse_error", "parse_failure":
		return "parse_failure"
	case "rate_limit":
		return "rate_limit"
	case "authentication_or_key", "missing_key":
		return "authentication_or_key"
	case "source_schema_change":
		return "source_schema_change"
	case "source_unavailable", "timeout", "unavailable", "http_error", "network":
		return "network"
	default:
		if strings.TrimSpace(value) == "" {
			return "network"
		}
		return strings.TrimSpace(value)
	}
}

func p52AnalystFailureCategory(value string) string {
	switch strings.TrimSpace(value) {
	case "missing_key", "authentication_or_key":
		return "authentication_or_key"
	case "quality_failed", "quality_failure":
		return "quality_failure"
	case "parse_error", "empty_response", "parse_failure":
		return "parse_failure"
	case "redaction_failed", "redaction_failure":
		return "redaction_failure"
	case "timeout", "unavailable", "http_error", "model_unavailable":
		return "model_unavailable"
	default:
		if strings.TrimSpace(value) == "" {
			return "model_unavailable"
		}
		return strings.TrimSpace(value)
	}
}
