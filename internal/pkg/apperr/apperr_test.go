package apperr

import (
	"database/sql"
	"errors"
	"net/http"
	"testing"
)

func TestAppErrorWrapAndMatchCode(t *testing.T) {
	cause := sql.ErrNoRows
	err := Wrap(CodeNotFound, CategoryNotFound, "record not found", cause)
	if !IsCode(err, CodeNotFound) {
		t.Fatalf("expected not found code")
	}
	if !errors.Is(err, cause) {
		t.Fatalf("expected wrapped cause")
	}
	appErr, ok := AsAppError(err)
	if !ok || appErr.Category != CategoryNotFound || appErr.HTTPStatus != http.StatusNotFound {
		t.Fatalf("unexpected app error: %+v", appErr)
	}
}

func TestHTTPStatusForCode(t *testing.T) {
	cases := map[Code]int{
		CodeBadRequest:               http.StatusBadRequest,
		CodeNotFound:                 http.StatusNotFound,
		CodeConflict:                 http.StatusConflict,
		CodeInvalidState:             http.StatusConflict,
		CodeEvidenceNotFound:         http.StatusConflict,
		CodeSourceVerificationFailed: http.StatusConflict,
		CodeInternalError:            http.StatusInternalServerError,
	}
	for code, want := range cases {
		if got := HTTPStatusForCode(code); got != want {
			t.Fatalf("%s status=%d want=%d", code, got, want)
		}
	}
}

func TestFromRepositoryError(t *testing.T) {
	if err := FromRepositoryError(nil); err != nil {
		t.Fatalf("nil err mapped to %v", err)
	}
	if err := FromRepositoryError(sql.ErrNoRows); !IsCode(err, CodeNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
	if err := FromRepositoryError(errors.New("constraint failed: UNIQUE constraint failed")); !IsCode(err, CodeConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
	original := New(CodeInvalidState, CategoryInvalidState, "invalid state")
	if err := FromRepositoryError(original); err != original {
		t.Fatalf("expected existing app error unchanged")
	}
}
