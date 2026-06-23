package httputil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteSuccess(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteSuccess(rec, "req_test", map[string]string{"status": "ok"})

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}

	var body Envelope
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.RequestID != "req_test" {
		t.Errorf("request_id = %q", body.RequestID)
	}
}
