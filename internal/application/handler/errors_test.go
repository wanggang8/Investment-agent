package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"investment-agent/internal/pkg/apperr"
	"investment-agent/pkg/httputil"
)

func TestWriteHandlerErrorUsesAppErrorEnvelope(t *testing.T) {
	w := httptest.NewRecorder()
	err := apperr.New(apperr.CodeEvidenceNotFound, apperr.CategoryConflict, "未找到有效证据")

	WriteHandlerError(w, "req_test", err)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
	var body httputil.Envelope
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.RequestID != "req_test" {
		t.Fatalf("expected request_id, got %q", body.RequestID)
	}
	if body.Error == nil || body.Error.Code != string(apperr.CodeEvidenceNotFound) || body.Error.Message != "未找到有效证据" {
		t.Fatalf("unexpected error body: %+v", body.Error)
	}
}

func TestWriteHandlerErrorHidesUnknownDetails(t *testing.T) {
	w := httptest.NewRecorder()

	WriteHandlerError(w, "req_test", assertErr("sql: database is locked at /tmp/local.db"))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
	var body httputil.Envelope
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Error == nil || body.Error.Code != string(apperr.CodeInternalError) {
		t.Fatalf("unexpected error body: %+v", body.Error)
	}
	if body.Error.Detail != "" || body.Error.Message == "sql: database is locked at /tmp/local.db" {
		t.Fatalf("raw detail leaked: %+v", body.Error)
	}
}

type assertErr string

func (e assertErr) Error() string { return string(e) }
