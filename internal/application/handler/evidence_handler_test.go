package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/application/service"
	"investment-agent/internal/application/workflow"
	"investment-agent/internal/domain/repository"
	"investment-agent/pkg/httputil"
)

func TestListEvidencePreservesSourceMetadata(t *testing.T) {
	app, db := testApp(t)
	_, err := db.Exec(`INSERT INTO intelligence_items (intelligence_id,source_name,source_level,original_url,published_at,captured_at,content_hash,raw_title,raw_text_ref,created_at) VALUES (?,?,?,?,?,?,?,?,?,?)`, "intel_meta", "交易所公告", "A", "https://example.com/a", "2026-01-02T00:00:00Z", "2026-01-02T01:00:00Z", "hash_meta", "标题", "raw", "2026-01-02T01:00:00Z")
	if err != nil {
		t.Fatalf("seed item: %v", err)
	}
	_, err = db.Exec(`INSERT INTO intelligence_summary (summary_id,intelligence_id,symbol,entity,source_level,evidence_role,event_type,summary,time_weight,relevance_score,verification_group_id,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`, "sum_meta", "intel_meta", "510300", "旧实体", "A", "formal", "major_positive", "摘要", 0.7, 0.8, "vg_meta", "2026-01-02T01:00:00Z")
	if err != nil {
		t.Fatalf("seed summary: %v", err)
	}

	_, err = db.Exec(`INSERT INTO source_verifications (verification_id,verification_group_id,event_id,symbol,event_type,evidence_role,verification_status,independent_source_count,high_grade_independent_source_count,highest_source_level,evidence_ids_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`, "ver_meta", "vg_meta", "event_1", "510300", "normal", "formal", "satisfied", 3, 2, "A", `["sum_meta"]`, "2026-01-02T01:00:00Z")
	if err != nil {
		t.Fatalf("seed verification: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/evidence", nil)
	req.Header.Set("X-Request-ID", "req_ev_list")
	w := httptest.NewRecorder()

	app.ListEvidence(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data struct {
			Items []dto.EvidenceDTO `json:"items"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Data.Items) != 1 {
		t.Fatalf("expected one evidence, got %+v", body.Data.Items)
	}
	item := body.Data.Items[0]
	if item.SourceName != "交易所公告" || item.OriginalURL != "https://example.com/a" || item.PublishedAt != "2026-01-02T00:00:00Z" || item.CapturedAt != "2026-01-02T01:00:00Z" || item.ContentHash != "hash_meta" || !floatClose(item.TimeWeight, 0.7) || !floatClose(item.RelevanceScore, 0.8) || item.HighGradeIndependentSourceCount != 2 {
		t.Fatalf("metadata not preserved: %+v", item)
	}
}

func TestGetEvidenceVerificationReturnsHighGradeIndependentSourceCount(t *testing.T) {
	app, db := testApp(t)
	_, err := db.Exec(`INSERT INTO source_verifications (verification_id,verification_group_id,event_id,symbol,event_type,evidence_role,verification_status,independent_source_count,high_grade_independent_source_count,highest_source_level,evidence_ids_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`, "ver_high_grade", "vg_high_grade", "event_1", "510300", "normal", "formal", "satisfied", 3, 2, "S", `["ev_1"]`, "2026-01-02T01:00:00Z")
	if err != nil {
		t.Fatalf("seed verification: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/evidence/verification", nil)
	req.Header.Set("X-Request-ID", "req_ev_verification")
	w := httptest.NewRecorder()

	app.GetEvidenceVerification(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data dto.SourceVerificationDTO `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.IndependentSourceCount != 3 || body.Data.HighGradeIndependentSourceCount != 2 {
		t.Fatalf("expected source counts, got %+v", body.Data)
	}
}

func TestRefreshEvidenceWritesVerificationAuditFields(t *testing.T) {
	app, db := testApp(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/evidence/refresh", bytes.NewBufferString(`{"symbol":"510300","refresh_scope":"all"}`))
	req.Header.Set("X-Request-ID", "req_ev_refresh_verification")
	w := httptest.NewRecorder()

	app.RefreshEvidence(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var latestPublishedAt, evidenceIDsJSON string
	if err := db.QueryRow(`SELECT COALESCE(latest_published_at,''),COALESCE(evidence_ids_json,'') FROM source_verifications ORDER BY created_at DESC LIMIT 1`).Scan(&latestPublishedAt, &evidenceIDsJSON); err != nil {
		t.Fatalf("read verification: %v", err)
	}
	if latestPublishedAt == "" {
		t.Fatalf("expected latest_published_at to be stored")
	}
	var evidenceIDs []string
	if err := json.Unmarshal([]byte(evidenceIDsJSON), &evidenceIDs); err != nil {
		t.Fatalf("decode evidence ids: %v", err)
	}
	if len(evidenceIDs) == 0 || evidenceIDs[0] == "" {
		t.Fatalf("expected verification evidence ids, got %q", evidenceIDsJSON)
	}
}

func TestRefreshEvidenceIndexFailureKeepsSQLiteFacts(t *testing.T) {
	app, db := testApp(t)
	app.Deps.VectorIndexWriter = failingVectorIndexWriter{}
	query := `SELECT index_status FROM rag_chunks ORDER BY created_at DESC LIMIT 1`
	if _, err := db.Exec(`INSERT INTO intelligence_items (intelligence_id,source_name,source_level,captured_at,content_hash,created_at) VALUES (?,?,?,?,?,?)`, "intel_pending", "交易所", "A", "2026-01-02T00:00:00Z", "hash_pending", "2026-01-02T00:00:00Z"); err != nil {
		t.Fatalf("seed item: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO intelligence_summary (summary_id,intelligence_id,symbol,source_level,evidence_role,event_type,summary,created_at) VALUES (?,?,?,?,?,?,?,?)`, "sum_pending", "intel_pending", "510300", "A", "formal", "normal", "待索引摘要", "2026-01-02T00:00:00Z"); err != nil {
		t.Fatalf("seed summary: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO rag_chunks (chunk_id,summary_id,chunk_text,chunk_hash,index_status,created_at) VALUES (?,?,?,?,?,?)`, "chunk_pending", "sum_pending", "待索引文本", "chunk_hash_pending", "pending", "2026-01-02T00:00:00Z"); err != nil {
		t.Fatalf("seed chunk: %v", err)
	}
	var preStatus string
	if err := db.QueryRow(query).Scan(&preStatus); err != nil {
		t.Fatalf("read seeded rag chunk status: %v", err)
	}
	if preStatus != "pending" {
		t.Fatalf("new chunks should start pending before vector indexing, got %s", preStatus)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/evidence/refresh", bytes.NewBufferString(`{"symbol":"510300","refresh_scope":"all"}`))
	req.Header.Set("X-Request-ID", "req_ev_index")
	w := httptest.NewRecorder()

	app.RefreshEvidence(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	assertCount(t, db, "intelligence_items", 2)
	assertCount(t, db, "intelligence_summary", 2)
	assertCount(t, db, "source_verifications", 1)
	var indexStatus string
	if err := db.QueryRow(`SELECT index_status FROM rag_chunks WHERE chunk_id <> 'chunk_pending' ORDER BY created_at DESC LIMIT 1`).Scan(&indexStatus); err != nil {
		t.Fatalf("read rag chunk status: %v", err)
	}
	if indexStatus != "failed" {
		t.Fatalf("expected failed index status, got %s", indexStatus)
	}
	var body struct {
		httputil.Envelope
		Data dto.EvidenceRefreshResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.IndexStatus != "failed" || body.Data.FailedReason == "" {
		t.Fatalf("expected index failure fields, got %+v", body.Data)
	}
	if body.Data.IntelligenceItemCount != 1 || body.Data.SummaryCount != 1 || body.Data.VerificationCount != 1 || body.Data.RAGChunkCount != 0 || len(body.Data.AuditEventIDs) != 6 {
		t.Fatalf("expected contract count fields, got %+v", body.Data)
	}
	var notificationType, sourceType, sourceID string
	if err := db.QueryRow(`SELECT type,COALESCE(source_type,''),COALESCE(source_id,'') FROM notifications WHERE read_at IS NULL ORDER BY created_at DESC LIMIT 1`).Scan(&notificationType, &sourceType, &sourceID); err != nil {
		t.Fatalf("read vector index notification: %v", err)
	}
	if notificationType != "vector_index_failure" || sourceType != "evidence_refresh" || sourceID != "510300" {
		t.Fatalf("expected vector index failure notification, got type=%s source=%s/%s", notificationType, sourceType, sourceID)
	}
}

type handlerTestIntelligenceSource struct {
	items []workflow.IntelligenceSourceItem
}

func (s handlerTestIntelligenceSource) FetchIntelligence(context.Context, string) ([]workflow.IntelligenceSourceItem, error) {
	return s.items, nil
}

func TestRefreshEvidenceReportsAllWrittenSummaries(t *testing.T) {
	app, _ := testApp(t)
	app.Deps.IntelligenceSource = handlerTestIntelligenceSource{items: []workflow.IntelligenceSourceItem{
		{SourceName: "交易所公告", SourceLevel: "A", Title: "公告一", Text: "正文一", PublishedAt: "2026-06-01T00:00:00Z"},
		{SourceName: "基金公司公告", SourceLevel: "A", Title: "公告二", Text: "正文二", PublishedAt: "2026-06-01T01:00:00Z"},
	}}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/evidence/refresh", bytes.NewBufferString(`{"symbol":"510300","refresh_scope":"all"}`))
	req.Header.Set("X-Request-ID", "req_ev_multi_summary")
	w := httptest.NewRecorder()

	app.RefreshEvidence(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data dto.EvidenceRefreshResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.IntelligenceItemCount != 2 || body.Data.SummaryCount != 2 {
		t.Fatalf("expected all summary counts, got %+v", body.Data)
	}
}

func TestRebuildEvidenceIndexWritesConfiguredVectorIndex(t *testing.T) {
	app, db := testApp(t)
	path := filepath.Join(t.TempDir(), "index.json")
	app.VectorIndex = service.NewFileVectorIndex(path)
	if _, err := db.Exec(`INSERT INTO intelligence_items (intelligence_id,source_name,source_level,captured_at,content_hash,created_at) VALUES (?,?,?,?,?,?)`, "intel_rebuild", "交易所", "A", "2026-01-02T00:00:00Z", "hash_rebuild", "2026-01-02T00:00:00Z"); err != nil {
		t.Fatalf("seed item: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO intelligence_summary (summary_id,intelligence_id,symbol,source_level,evidence_role,event_type,summary,created_at) VALUES (?,?,?,?,?,?,?,?)`, "sum_rebuild", "intel_rebuild", "510300", "A", "formal", "normal", "重建摘要", "2026-01-02T00:00:00Z"); err != nil {
		t.Fatalf("seed summary: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO rag_chunks (chunk_id,summary_id,chunk_text,chunk_hash,index_status,created_at) VALUES (?,?,?,?,?,?)`, "chunk_rebuild", "sum_rebuild", "重建文本", "chunk_hash", "pending", "2026-01-02T00:00:00Z"); err != nil {
		t.Fatalf("seed chunk: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/evidence/rebuild-index", nil)
	req.Header.Set("X-Request-ID", "req_rebuild_index")
	w := httptest.NewRecorder()
	app.RebuildEvidenceIndex(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	chunks, err := service.NewFileVectorIndex(path).Search(req.Context(), "510300")
	if err != nil {
		t.Fatalf("search rebuilt index: %v", err)
	}
	if len(chunks) != 1 || chunks[0].ChunkID != "chunk_rebuild" {
		t.Fatalf("expected rebuilt chunk, got %+v", chunks)
	}
	var body struct {
		Data dto.RebuildIndexResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.IndexedCount != 1 || body.Data.SkippedCount != 0 || body.Data.IndexHealth.Status != service.VectorIndexHealthHealthy || body.Data.LastRebuildAt == "" {
		t.Fatalf("expected rebuild stats and health, got %+v", body.Data)
	}
}

type failingVectorIndex struct{}

func (f failingVectorIndex) Upsert(context.Context, repository.RAGChunk) error { return errors.New("index unavailable") }
func (f failingVectorIndex) Search(context.Context, string) ([]repository.RAGChunk, error) { return nil, nil }

func TestRebuildEvidenceIndexFailureWritesNotification(t *testing.T) {
	app, db := testApp(t)
	app.VectorIndex = failingVectorIndex{}
	if _, err := db.Exec(`INSERT INTO intelligence_items (intelligence_id,source_name,source_level,captured_at,content_hash,created_at) VALUES (?,?,?,?,?,?)`, "intel_rebuild_fail", "交易所", "A", "2026-01-02T00:00:00Z", "hash_rebuild_fail", "2026-01-02T00:00:00Z"); err != nil {
		t.Fatalf("seed item: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO intelligence_summary (summary_id,intelligence_id,symbol,source_level,evidence_role,event_type,summary,created_at) VALUES (?,?,?,?,?,?,?,?)`, "sum_rebuild_fail", "intel_rebuild_fail", "510300", "A", "formal", "normal", "重建失败摘要", "2026-01-02T00:00:00Z"); err != nil {
		t.Fatalf("seed summary: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO rag_chunks (chunk_id,summary_id,chunk_text,chunk_hash,index_status,created_at) VALUES (?,?,?,?,?,?)`, "chunk_rebuild_fail", "sum_rebuild_fail", "重建失败文本", "chunk_hash_fail", "pending", "2026-01-02T00:00:00Z"); err != nil {
		t.Fatalf("seed chunk: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/evidence/rebuild-index", nil)
	req.Header.Set("X-Request-ID", "req_rebuild_index_fail")
	w := httptest.NewRecorder()
	app.RebuildEvidenceIndex(w, req)

	if w.Code < 400 {
		t.Fatalf("expected failure response, got %d body=%s", w.Code, w.Body.String())
	}
	var notificationType, sourceType, sourceID string
	if err := db.QueryRow(`SELECT type,COALESCE(source_type,''),COALESCE(source_id,'') FROM notifications WHERE read_at IS NULL ORDER BY created_at DESC LIMIT 1`).Scan(&notificationType, &sourceType, &sourceID); err != nil {
		t.Fatalf("read rebuild notification: %v", err)
	}
	if notificationType != "vector_index_failure" || sourceType != "evidence_rebuild" || sourceID != "local_vector_index" {
		t.Fatalf("expected rebuild vector index failure notification, got type=%s source=%s/%s", notificationType, sourceType, sourceID)
	}
}

func TestRebuildEvidenceIndexReportsMissingHealthWhenNoChunks(t *testing.T) {
	app, _ := testApp(t)
	path := filepath.Join(t.TempDir(), "missing-index.json")
	app.VectorIndex = service.NewFileVectorIndex(path)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/evidence/rebuild-index", nil)
	req.Header.Set("X-Request-ID", "req_rebuild_empty_index")
	w := httptest.NewRecorder()
	app.RebuildEvidenceIndex(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data dto.RebuildIndexResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.IndexedCount != 0 || body.Data.IndexHealth.Status != service.VectorIndexHealthHealthy || body.Data.IndexHealth.Path != path {
		t.Fatalf("expected visible index health for empty rebuild, got %+v", body.Data)
	}
}

func TestRefreshEvidenceWritesConfiguredVectorIndex(t *testing.T) {
	app, _ := testApp(t)
	path := filepath.Join(t.TempDir(), "refresh-index.json")
	index := service.NewFileVectorIndex(path)
	app.Deps.VectorIndexWriter = index

	req := httptest.NewRequest(http.MethodPost, "/api/v1/evidence/refresh", bytes.NewBufferString(`{"symbol":"510300","refresh_scope":"all"}`))
	req.Header.Set("X-Request-ID", "req_ev_refresh_index")
	w := httptest.NewRecorder()
	app.RefreshEvidence(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	chunks, err := service.NewFileVectorIndex(path).Search(req.Context(), "510300")
	if err != nil {
		t.Fatalf("search refreshed index: %v", err)
	}
	if len(chunks) == 0 {
		t.Fatalf("expected refreshed evidence chunks in vector index")
	}
}

func TestRefreshEvidenceReportsVectorIndexFailure(t *testing.T) {
	app, db := testApp(t)
	app.Deps.VectorIndexWriter = failingVectorIndexWriter{}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/evidence/refresh", bytes.NewBufferString(`{"symbol":"510300","refresh_scope":"all"}`))
	req.Header.Set("X-Request-ID", "req_ev_refresh_index_failed")
	w := httptest.NewRecorder()
	app.RefreshEvidence(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	assertCount(t, db, "intelligence_items", 1)
	var body struct {
		Data dto.EvidenceRefreshResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.IndexStatus != "failed" || body.Data.FailedReason == "" || body.Data.RAGChunkCount != 0 {
		t.Fatalf("expected failed index response, got %+v", body.Data)
	}
}

func TestRefreshEvidenceUsesActualSourcesForVerification(t *testing.T) {
	app, db := testApp(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/evidence/refresh", bytes.NewBufferString(`{"symbol":"510300","refresh_scope":"all","sources":["official","exchange"]}`))
	req.Header.Set("X-Request-ID", "req_ev_actual_one_source")
	w := httptest.NewRecorder()
	app.RefreshEvidence(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var status string
	var independentCount int
	if err := db.QueryRow(`SELECT verification_status,independent_source_count FROM source_verifications ORDER BY created_at DESC LIMIT 1`).Scan(&status, &independentCount); err != nil {
		t.Fatalf("read verification: %v", err)
	}
	if status != "failed" || independentCount != 1 {
		t.Fatalf("expected one actual source to fail verification, status=%s count=%d", status, independentCount)
	}
}

func TestGetEvidenceVerificationFiltersBySymbolAndEventID(t *testing.T) {
	app, db := testApp(t)
	if _, err := db.Exec(`INSERT INTO source_verifications (verification_id,verification_group_id,event_id,symbol,event_type,evidence_role,verification_status,independent_source_count,high_grade_independent_source_count,highest_source_level,evidence_ids_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`, "ver_a", "vg_a", "event_a", "510300", "normal", "formal", "satisfied", 2, 2, "A", `["ev_a"]`, "2026-01-02T01:00:00Z"); err != nil {
		t.Fatalf("seed first verification: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO source_verifications (verification_id,verification_group_id,event_id,symbol,event_type,evidence_role,verification_status,independent_source_count,high_grade_independent_source_count,highest_source_level,evidence_ids_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`, "ver_b", "vg_b", "event_b", "159915", "normal", "formal", "background_only", 2, 1, "A", `["ev_b"]`, "2026-01-03T01:00:00Z"); err != nil {
		t.Fatalf("seed second verification: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/evidence/verification?symbol=510300&event_id=event_a", nil)
	req.Header.Set("X-Request-ID", "req_ev_verification_filter")
	w := httptest.NewRecorder()
	app.GetEvidenceVerification(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data dto.SourceVerificationDTO `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.VerificationID != "ver_a" || body.Data.VerificationStatus != "satisfied" {
		t.Fatalf("expected filtered verification, got %+v", body.Data)
	}
}

func TestGetEvidenceVerificationReturnsNotFoundForMismatchedFilter(t *testing.T) {
	app, db := testApp(t)
	if _, err := db.Exec(`INSERT INTO source_verifications (verification_id,verification_group_id,event_id,symbol,event_type,evidence_role,verification_status,independent_source_count,high_grade_independent_source_count,highest_source_level,evidence_ids_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`, "ver_mismatch", "vg_mismatch", "event_a", "510300", "normal", "formal", "satisfied", 2, 2, "A", `["ev_a"]`, "2026-01-02T01:00:00Z"); err != nil {
		t.Fatalf("seed verification: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/evidence/verification?symbol=159915&event_id=event_a", nil)
	req.Header.Set("X-Request-ID", "req_ev_verification_mismatch")
	w := httptest.NewRecorder()
	app.GetEvidenceVerification(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d body=%s", w.Code, w.Body.String())
	}
}

type failingVectorIndexWriter struct{}

func (failingVectorIndexWriter) Upsert(context.Context, repository.RAGChunk) error {
	return errors.New("index write failed")
}
