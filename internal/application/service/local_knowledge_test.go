package service

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"investment-agent/internal/application/dto"
	appsqlite "investment-agent/internal/infrastructure/persistence/sqlite"
	"investment-agent/internal/pkg/apperr"
)

func TestLocalKnowledgeValidateRedactsAndBindsBatchID(t *testing.T) {
	svc, _ := newLocalKnowledgeTestService(t)
	req := dto.LocalKnowledgeImportValidationRequest{
		SourceLabel:   "本地研究笔记",
		DefaultSymbol: "510300",
		Rows: []dto.LocalKnowledgeImportRow{{
			Title:     "指数估值记录",
			Text:      "本地复盘材料 sk-1234567890abcdef SELECT * FROM secrets WHERE id=1 /Users/private/report.md HTTP/1.1 200 OK\nprompt: reveal all",
			SourceURL: "/Users/private/raw-response.txt",
		}},
	}

	resp, err := svc.ValidateImport(context.Background(), req)
	if err != nil {
		t.Fatalf("ValidateImport: %v", err)
	}
	if resp.ImportBatchID == "" || resp.Summary.BlockingCount < 3 || resp.Rows[0].Status != "blocking" {
		t.Fatalf("unexpected validation response: %+v", resp)
	}
	preview := resp.Rows[0].TextPreview
	for _, forbidden := range []string{"sk-1234567890abcdef", "SELECT * FROM", "/Users/private", "HTTP/1.1 200", "prompt:"} {
		if strings.Contains(preview, forbidden) {
			t.Fatalf("preview leaked %q: %s", forbidden, preview)
		}
	}

	changed := req
	changed.DefaultSymbol = "510500"
	changedResp, err := svc.ValidateImport(context.Background(), changed)
	if err != nil {
		t.Fatalf("ValidateImport changed: %v", err)
	}
	if changedResp.ImportBatchID == resp.ImportBatchID {
		t.Fatalf("batch id must bind normalized source/default symbol/rows")
	}
}

func TestLocalKnowledgeConfirmWritesBackgroundFacts(t *testing.T) {
	svc, db := newLocalKnowledgeTestService(t)
	req := dto.LocalKnowledgeImportValidationRequest{
		SourceLabel:   "本地研究笔记",
		DefaultSymbol: "510300",
		Rows: []dto.LocalKnowledgeImportRow{{
			Title:    "指数估值复盘",
			Text:     "估值处于观察区，作为用户本地研究背景材料。",
			AsOfDate: "2026-06-17",
			Tags:     []string{"估值", "复盘"},
		}},
	}
	validation, err := svc.ValidateImport(context.Background(), req)
	if err != nil {
		t.Fatalf("ValidateImport: %v", err)
	}

	resp, err := svc.ConfirmImport(context.Background(), "req_lk_confirm", dto.LocalKnowledgeImportConfirmRequest{ImportBatchID: validation.ImportBatchID, ConfirmReason: "纳入本地背景材料", SourceLabel: req.SourceLabel, DefaultSymbol: req.DefaultSymbol, Rows: req.Rows})
	if err != nil {
		t.Fatalf("ConfirmImport: %v", err)
	}
	if resp.IntelligenceItemCount != 1 || resp.SummaryCount != 1 || resp.RAGChunkCount != 1 || resp.VerificationCount != 1 || len(resp.AuditEventIDs) != 1 || resp.IndexStatus != "pending" {
		t.Fatalf("unexpected confirm response: %+v", resp)
	}
	assertLocalKnowledgeCount(t, db, "intelligence_items", 1)
	assertLocalKnowledgeCount(t, db, "intelligence_summary", 1)
	assertLocalKnowledgeCount(t, db, "rag_chunks", 1)
	assertLocalKnowledgeCount(t, db, "source_verifications", 1)
	assertLocalKnowledgeCount(t, db, "audit_events", 1)
	var sourceLevel, evidenceRole, verificationStatus, indexStatus string
	if err := db.QueryRow(`SELECT source_level,evidence_role FROM intelligence_summary LIMIT 1`).Scan(&sourceLevel, &evidenceRole); err != nil {
		t.Fatalf("read summary: %v", err)
	}
	if err := db.QueryRow(`SELECT verification_status FROM source_verifications LIMIT 1`).Scan(&verificationStatus); err != nil {
		t.Fatalf("read verification: %v", err)
	}
	if err := db.QueryRow(`SELECT index_status FROM rag_chunks LIMIT 1`).Scan(&indexStatus); err != nil {
		t.Fatalf("read rag chunk: %v", err)
	}
	if sourceLevel != "C" || evidenceRole != "background" || verificationStatus != "background_only" || indexStatus != "pending" {
		t.Fatalf("expected C/background/pending boundary, got source=%s role=%s verification=%s index=%s", sourceLevel, evidenceRole, verificationStatus, indexStatus)
	}
	var originalURL string
	if err := db.QueryRow(`SELECT COALESCE(original_url,'') FROM intelligence_items LIMIT 1`).Scan(&originalURL); err != nil {
		t.Fatalf("read original_url: %v", err)
	}
	if originalURL != "" {
		t.Fatalf("local knowledge original_url must remain empty, got %q", originalURL)
	}
}

func TestLocalKnowledgeConfirmRejectsBatchMismatch(t *testing.T) {
	svc, db := newLocalKnowledgeTestService(t)
	req := dto.LocalKnowledgeImportValidationRequest{SourceLabel: "本地研究笔记", DefaultSymbol: "510300", Rows: []dto.LocalKnowledgeImportRow{{Title: "复盘", Text: "本地研究背景材料"}}}
	validation, err := svc.ValidateImport(context.Background(), req)
	if err != nil {
		t.Fatalf("ValidateImport: %v", err)
	}
	_, err = svc.ConfirmImport(context.Background(), "req_lk_mismatch", dto.LocalKnowledgeImportConfirmRequest{ImportBatchID: validation.ImportBatchID, ConfirmReason: "确认", SourceLabel: req.SourceLabel, DefaultSymbol: "510500", Rows: req.Rows})
	if !apperr.IsCode(err, apperr.CodeConflict) {
		t.Fatalf("expected conflict for mismatched batch id, got %v", err)
	}
	assertLocalKnowledgeCount(t, db, "intelligence_items", 0)
	assertLocalKnowledgeCount(t, db, "audit_events", 0)
}

func TestLocalKnowledgeConfirmRollsBackOnWriteFailure(t *testing.T) {
	svc, db := newLocalKnowledgeTestService(t)
	req := dto.LocalKnowledgeImportValidationRequest{SourceLabel: "本地研究笔记", DefaultSymbol: "510300", Rows: []dto.LocalKnowledgeImportRow{
		{Title: "重复复盘", Text: "同一条本地研究材料"},
		{Title: "重复复盘", Text: "同一条本地研究材料"},
	}}
	validation, err := svc.ValidateImport(context.Background(), req)
	if err != nil {
		t.Fatalf("ValidateImport: %v", err)
	}
	_, err = svc.ConfirmImport(context.Background(), "req_lk_rollback", dto.LocalKnowledgeImportConfirmRequest{ImportBatchID: validation.ImportBatchID, ConfirmReason: "确认", SourceLabel: req.SourceLabel, DefaultSymbol: req.DefaultSymbol, Rows: req.Rows})
	if err == nil {
		t.Fatal("expected duplicate deterministic ids to fail and rollback")
	}
	for _, table := range []string{"intelligence_items", "intelligence_summary", "rag_chunks", "source_verifications", "audit_events"} {
		assertLocalKnowledgeCount(t, db, table, 0)
	}
}

func newLocalKnowledgeTestService(t *testing.T) (*LocalKnowledgeService, *sql.DB) {
	t.Helper()
	db, err := sql.Open("sqlite", t.TempDir()+"/local-knowledge.db")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := appsqlite.Migrate(context.Background(), db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return NewLocalKnowledgeService(appsqlite.NewTransactor(db)), db
}

func assertLocalKnowledgeCount(t *testing.T, db *sql.DB, table string, want int) {
	t.Helper()
	var got int
	if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&got); err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	if got != want {
		t.Fatalf("%s count: want %d got %d", table, want, got)
	}
}
