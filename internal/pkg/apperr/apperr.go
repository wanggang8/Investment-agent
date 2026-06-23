package apperr

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
)

type Code string

type Category string

const (
	CodeInternalError             Code = "INTERNAL_ERROR"
	CodeBadRequest                Code = "BAD_REQUEST"
	CodeNotFound                  Code = "NOT_FOUND"
	CodeConflict                  Code = "CONFLICT"
	CodeInvalidState              Code = "INVALID_STATE"
	CodeDataRequired              Code = "DATA_REQUIRED"
	CodeDataStale                 Code = "DATA_STALE"
	CodeDataSourceUnavailable     Code = "DATA_SOURCE_UNAVAILABLE"
	CodeMarketSnapshotWriteFailed Code = "MARKET_SNAPSHOT_WRITE_FAILED"
	CodeRuleVersionMissing        Code = "RULE_VERSION_MISSING"
	CodeEvidenceNotFound          Code = "EVIDENCE_NOT_FOUND"
	CodeSourceVerificationFailed  Code = "SOURCE_VERIFICATION_FAILED"
	CodeVectorIndexUnavailable    Code = "VECTOR_INDEX_UNAVAILABLE"
	CodeAnalystUnavailable        Code = "ANALYST_UNAVAILABLE"
	CodeDecisionRecordFailed      Code = "DECISION_RECORD_FAILED"
)

const (
	CategoryInternal     Category = "internal"
	CategoryBadRequest   Category = "bad_request"
	CategoryNotFound     Category = "not_found"
	CategoryConflict     Category = "conflict"
	CategoryInvalidState Category = "invalid_state"
)

type AppError struct {
	Code       Code
	Category   Category
	Message    string
	HTTPStatus int
	Cause      error
}

func New(code Code, category Category, message string) *AppError {
	return &AppError{Code: code, Category: category, Message: message, HTTPStatus: HTTPStatusForCode(code)}
}

func Wrap(code Code, category Category, message string, cause error) *AppError {
	return &AppError{Code: code, Category: category, Message: message, HTTPStatus: HTTPStatusForCode(code), Cause: cause}
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func IsCode(err error, code Code) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}

func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

func HTTPStatusForCode(code Code) int {
	switch code {
	case CodeBadRequest:
		return http.StatusBadRequest
	case CodeNotFound:
		return http.StatusNotFound
	case CodeDataSourceUnavailable:
		return http.StatusServiceUnavailable
	case CodeMarketSnapshotWriteFailed:
		return http.StatusInternalServerError
	case CodeConflict, CodeInvalidState, CodeDataRequired, CodeDataStale, CodeRuleVersionMissing, CodeEvidenceNotFound, CodeSourceVerificationFailed, CodeVectorIndexUnavailable, CodeAnalystUnavailable, CodeDecisionRecordFailed:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func AuditErrorCode(code Code) string {
	if code == "" {
		return string(CodeInternalError)
	}
	return string(code)
}

func FromRepositoryError(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := AsAppError(err); ok {
		return err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return Wrap(CodeNotFound, CategoryNotFound, "record not found", err)
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "constraint failed") || strings.Contains(msg, "unique constraint") {
		return Wrap(CodeConflict, CategoryConflict, "constraint conflict", err)
	}
	return Wrap(CodeInternalError, CategoryInternal, "repository error", err)
}
