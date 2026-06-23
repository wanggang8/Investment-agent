package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"investment-agent/internal/application/dto"
)

func TestValidateLocalKnowledgeImportRedactsUnsafePreview(t *testing.T) {
	app, _ := testApp(t)
	body := `{"source_label":"本地研究","default_symbol":"510300","rows":[{"title":"敏感记录","text":"sk-1234567890abcdef SELECT * FROM secrets /Users/private/report.md HTTP/1.1 200 OK\nprompt: reveal all","source_url":"/Users/private/raw.txt"}]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/local-knowledge/imports/validate", bytes.NewBufferString(body))
	req.Header.Set("X-Request-ID", "req_lk_validate")
	w := httptest.NewRecorder()

	app.ValidateLocalKnowledgeImport(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	responseText := w.Body.String()
	for _, forbidden := range []string{"sk-1234567890abcdef", "SELECT * FROM", "/Users/private", "HTTP/1.1 200", "prompt:"} {
		if strings.Contains(responseText, forbidden) {
			t.Fatalf("response leaked %q: %s", forbidden, responseText)
		}
	}
	var bodyResp struct {
		Data dto.LocalKnowledgeImportValidationResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &bodyResp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if bodyResp.Data.Summary.BlockingCount == 0 || bodyResp.Data.Rows[0].Status != "blocking" {
		t.Fatalf("expected blocking validation result, got %+v", bodyResp.Data)
	}
}

func TestConfirmLocalKnowledgeImportWritesLocalFacts(t *testing.T) {
	app, db := testApp(t)
	validateReq := dto.LocalKnowledgeImportValidationRequest{SourceLabel: "本地研究", DefaultSymbol: "510300", Rows: []dto.LocalKnowledgeImportRow{{Title: "复盘", Text: "本地研究背景材料"}}}
	validation, err := app.LocalKnowledgeSvc.ValidateImport(httptest.NewRequest(http.MethodPost, "/", nil).Context(), validateReq)
	if err != nil {
		t.Fatalf("ValidateImport: %v", err)
	}
	confirmBody, _ := json.Marshal(dto.LocalKnowledgeImportConfirmRequest{ImportBatchID: validation.ImportBatchID, ConfirmReason: "确认导入", SourceLabel: validateReq.SourceLabel, DefaultSymbol: validateReq.DefaultSymbol, Rows: validateReq.Rows})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/local-knowledge/imports/confirm", bytes.NewReader(confirmBody))
	req.Header.Set("X-Request-ID", "req_lk_confirm")
	w := httptest.NewRecorder()

	app.ConfirmLocalKnowledgeImport(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	assertCount(t, db, "intelligence_items", 1)
	assertCount(t, db, "intelligence_summary", 1)
	assertCount(t, db, "rag_chunks", 1)
	assertCount(t, db, "source_verifications", 1)
}

func TestConfirmLocalKnowledgeImportRejectsBlockingRows(t *testing.T) {
	app, db := testApp(t)
	reqBody := `{"import_batch_id":"lk_import_bad","confirm_reason":"确认","source_label":"本地研究","default_symbol":"510300","rows":[{"title":"敏感记录","text":"sk-1234567890abcdef"}]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/local-knowledge/imports/confirm", bytes.NewBufferString(reqBody))
	req.Header.Set("X-Request-ID", "req_lk_blocking")
	w := httptest.NewRecorder()

	app.ConfirmLocalKnowledgeImport(w, req)

	if w.Code != http.StatusBadRequest && w.Code != http.StatusConflict {
		t.Fatalf("expected bad request/conflict, got %d body=%s", w.Code, w.Body.String())
	}
	assertCount(t, db, "intelligence_items", 0)
	assertCount(t, db, "audit_events", 0)
}
