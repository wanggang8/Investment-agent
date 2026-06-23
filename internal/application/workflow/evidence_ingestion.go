package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

// PublicEvidenceError carries source-specific collector failure metadata for audit.
type PublicEvidenceError struct {
	SourceName string
	ErrorCode  string
	Count      int
	Err        error
}

func publicEvidenceErrorOf(err error) (PublicEvidenceError, bool) {
	var sourceErr PublicEvidenceError
	if errors.As(err, &sourceErr) {
		return sourceErr, true
	}
	return PublicEvidenceError{}, false
}

func (e PublicEvidenceError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.ErrorCode
}

func (e PublicEvidenceError) Unwrap() error { return e.Err }

// PublicEvidenceIngestionService 封装 P26 collector 输出到 intelligence_items/audit_events 的入库逻辑。
type PublicEvidenceIngestionService struct {
	Collector        PublicEvidenceCollector
	IntelligenceRepo repository.IntelligenceRepository
	AuditRepo        repository.AuditRepository
	GenerateAuditID  func() string
	RequestID        string
}

// IngestPublicEvidence 调用 collector 并将结果写入 intelligence_items、rag_chunks、source_verifications 和 audit_events。
func (s *PublicEvidenceIngestionService) IngestPublicEvidence(ctx context.Context, symbol string, start, end time.Time) error {
	if err := s.validate(); err != nil {
		return err
	}
	payloads, err := s.Collector.FetchPublicEvidence(ctx, symbol, start, end)
	if err != nil {
		return s.auditFailure(ctx, symbol, err)
	}
	payloads = applyPublicEvidenceRuntimePolicy(payloads, time.Now().UTC())

	saved := make([]repository.IntelligenceSummary, 0, len(payloads))
	for _, p := range payloads {
		summary, err := s.savePayload(ctx, p)
		if err != nil {
			return err
		}
		saved = append(saved, summary)
	}
	if err := s.saveVerifications(ctx, payloads, saved); err != nil {
		return err
	}
	if err := s.auditPartialFailures(ctx, symbol); err != nil {
		return err
	}

	return s.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{
		AuditEventID:  s.nextAuditID(),
		RequestID:     strings.TrimSpace(s.RequestID),
		Actor:         "system",
		Action:        "run_local_task",
		Status:        "success",
		InputRefType:  "symbol",
		InputRef:      symbol,
		OutputRefType: "public_evidence",
		OutputRef:     fmt.Sprintf("source=public_evidence count=%d", len(payloads)),
		CreatedAt:     time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *PublicEvidenceIngestionService) validate() error {
	if s.Collector == nil || s.IntelligenceRepo == nil || s.AuditRepo == nil {
		return apperr.New(apperr.CodeInternalError, apperr.CategoryInternal, "P26 证据采集依赖未配置")
	}
	return nil
}

func (s *PublicEvidenceIngestionService) savePayload(ctx context.Context, p PublicEvidencePayload) (repository.IntelligenceSummary, error) {
	createdAt := time.Now().UTC().Format(time.RFC3339)
	capturedAt := p.CapturedAt.Format(time.RFC3339)
	if p.CapturedAt.IsZero() {
		capturedAt = createdAt
	}

	intelligenceID := deterministicEvidenceID("intel", p)
	summaryID := deterministicEvidenceID("summary", p)
	chunkID := deterministicEvidenceID("chunk", p)
	verificationGroupID := deterministicEvidenceID("group", publicEvidenceEventPayload(p))
	contentHash := p.ContentHash

	item := repository.IntelligenceItem{
		IntelligenceID: intelligenceID,
		SourceName:     p.SourceName,
		SourceLevel:    string(p.SourceLevel),
		OriginalURL:    p.URL,
		PublishedAt:    p.PublishedAt,
		CapturedAt:     capturedAt,
		ContentHash:    contentHash,
		RawTitle:       p.Title,
		RawTextRef:     p.Text,
		CreatedAt:      createdAt,
	}
	if err := s.IntelligenceRepo.SaveIntelligenceItem(ctx, item); err != nil && !apperr.IsCode(err, apperr.CodeConflict) {
		return repository.IntelligenceSummary{}, err
	}

	summary := repository.IntelligenceSummary{
		SummaryID:           summaryID,
		IntelligenceID:      intelligenceID,
		Symbol:              p.Symbol,
		Entity:              p.Symbol,
		EventType:           p.SourceType,
		Summary:             p.Text,
		SourceLevel:         string(p.SourceLevel),
		EvidenceRole:        p.EvidenceRole,
		TimeWeight:          publicEvidencePayloadTimeWeight(p),
		RelevanceScore:      1,
		VerificationGroupID: verificationGroupID,
		CreatedAt:           createdAt,
	}
	chunk := repository.RAGChunk{
		ChunkID:      chunkID,
		SummaryID:    summaryID,
		Symbol:       p.Symbol,
		ChunkText:    p.Text,
		ChunkHash:    contentHash,
		IndexStatus:  "pending",
		MetadataJSON: publicEvidenceMetadataJSON(p, s.RequestID),
		CreatedAt:    createdAt,
	}
	if err := s.IntelligenceRepo.SaveIntelligenceSummary(ctx, summary, []repository.RAGChunk{chunk}); err != nil && !apperr.IsCode(err, apperr.CodeConflict) {
		return repository.IntelligenceSummary{}, err
	}
	return summary, nil
}

func (s *PublicEvidenceIngestionService) saveVerifications(ctx context.Context, payloads []PublicEvidencePayload, summaries []repository.IntelligenceSummary) error {
	groups := map[string][]int{}
	for i, p := range payloads {
		groups[publicEvidenceEventID(p)] = append(groups[publicEvidenceEventID(p)], i)
	}
	for eventID, indexes := range groups {
		sources := map[string]bool{}
		highGradeSources := map[string]bool{}
		highest := ""
		latest := ""
		evidenceIDs := make([]string, 0, len(indexes))
		sample := payloads[indexes[0]]
		for _, idx := range indexes {
			p := payloads[idx]
			sources[p.SourceName] = true
			if highGradeCount(string(p.SourceLevel)) > 0 {
				highGradeSources[p.SourceName] = true
			}
			if highest == "" || sourceLevelRank(string(p.SourceLevel)) > sourceLevelRank(highest) {
				highest = string(p.SourceLevel)
			}
			if latest == "" || p.PublishedAt > latest {
				latest = p.PublishedAt
			}
			evidenceIDs = append(evidenceIDs, summaries[idx].SummaryID)
		}
		sort.Strings(evidenceIDs)
		evidenceIDsJSON, err := json.Marshal(evidenceIDs)
		if err != nil {
			return apperr.Wrap(apperr.CodeInternalError, apperr.CategoryInternal, "证据 ID 序列化失败", err)
		}
		status := "failed"
		if publicEvidenceVerificationSatisfied(sample, len(sources), len(highGradeSources)) {
			status = "satisfied"
		}
		verification := repository.SourceVerification{
			VerificationID:                  "verification-" + stableHash(eventID),
			VerificationGroupID:             deterministicEvidenceID("group", publicEvidenceEventPayload(sample)),
			EventID:                         eventID,
			Symbol:                          sample.Symbol,
			EventType:                       sample.SourceType,
			EvidenceRole:                    sample.EvidenceRole,
			VerificationStatus:              status,
			IndependentSourceCount:          len(sources),
			HighGradeIndependentSourceCount: len(highGradeSources),
			HighestSourceLevel:              highest,
			LatestPublishedAt:               latest,
			EvidenceIDsJSON:                 string(evidenceIDsJSON),
			CreatedAt:                       time.Now().UTC().Format(time.RFC3339),
		}
		if keep, err := s.shouldKeepExistingVerification(ctx, verification); err != nil {
			return err
		} else if keep {
			continue
		}
		if err := s.IntelligenceRepo.SaveSourceVerification(ctx, verification); err != nil && !apperr.IsCode(err, apperr.CodeConflict) {
			return err
		}
	}
	return nil
}

func (s *PublicEvidenceIngestionService) shouldKeepExistingVerification(ctx context.Context, next repository.SourceVerification) (bool, error) {
	existing, err := s.IntelligenceRepo.GetLatestSourceVerificationByFilter(ctx, next.Symbol, next.EventID)
	if err != nil {
		if apperr.IsCode(err, apperr.CodeNotFound) {
			return false, nil
		}
		return false, err
	}
	if existing.VerificationStatus == "satisfied" && next.VerificationStatus != "satisfied" {
		return true, nil
	}
	if existing.IndependentSourceCount > next.IndependentSourceCount || existing.HighGradeIndependentSourceCount > next.HighGradeIndependentSourceCount {
		return true, nil
	}
	return false, nil
}

func applyPublicEvidenceRuntimePolicy(payloads []PublicEvidencePayload, now time.Time) []PublicEvidencePayload {
	result := make([]PublicEvidencePayload, 0, len(payloads))
	for _, p := range payloads {
		p.TimeWeight = publicEvidenceTimeWeight(p.PublishedAt, now)
		if p.TimeWeight <= 0.2 {
			p.EvidenceRole = string(model.EvidenceBackground)
		}
		result = append(result, p)
	}
	return result
}

func publicEvidencePayloadTimeWeight(p PublicEvidencePayload) float64 {
	if p.TimeWeight > 0 {
		return p.TimeWeight
	}
	return 1
}

func (s *PublicEvidenceIngestionService) auditPartialFailures(ctx context.Context, symbol string) error {
	reporter, ok := s.Collector.(interface{ PublicEvidenceFailures() []PublicEvidenceError })
	if !ok {
		return nil
	}
	for _, failure := range reporter.PublicEvidenceFailures() {
		errorCode := p52CollectorFailureCategory(failure.ErrorCode)
		if failure.SourceName != "" {
			errorCode = failure.SourceName + ":" + errorCode
		}
		if err := s.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{
			AuditEventID:  s.nextAuditID(),
			RequestID:     strings.TrimSpace(s.RequestID),
			Actor:         "system",
			Action:        "run_local_task",
			Status:        "degraded",
			ErrorCode:     errorCode,
			InputRefType:  "symbol",
			InputRef:      symbol,
			OutputRefType: "public_evidence",
			OutputRef:     fmt.Sprintf("source=%s count=%d", failure.SourceName, failure.Count),
			CreatedAt:     time.Now().UTC().Format(time.RFC3339),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s *PublicEvidenceIngestionService) auditFailure(ctx context.Context, symbol string, err error) error {
	errorCode := string(apperr.CodeInternalError)
	outputRef := "source=unknown count=0"
	var sourceErr PublicEvidenceError
	if errors.As(err, &sourceErr) {
		errorCode = p52CollectorFailureCategory(sourceErr.ErrorCode)
		if sourceErr.SourceName != "" {
			errorCode = sourceErr.SourceName + ":" + errorCode
			outputRef = fmt.Sprintf("source=%s count=%d", sourceErr.SourceName, sourceErr.Count)
		}
	} else if appErr, ok := apperr.AsAppError(err); ok {
		errorCode = string(appErr.Code)
	}
	auditErr := s.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{
		AuditEventID:  s.nextAuditID(),
		RequestID:     strings.TrimSpace(s.RequestID),
		Actor:         "system",
		Action:        "run_local_task",
		Status:        "failed",
		ErrorCode:     errorCode,
		InputRefType:  "symbol",
		InputRef:      symbol,
		OutputRefType: "public_evidence",
		OutputRef:     outputRef,
		CreatedAt:     time.Now().UTC().Format(time.RFC3339),
	})
	if auditErr != nil {
		return apperr.Wrap(apperr.CodeInternalError, apperr.CategoryInternal, "写入审计失败", auditErr)
	}
	return err
}

func (s *PublicEvidenceIngestionService) nextAuditID() string {
	if s.GenerateAuditID != nil {
		return s.GenerateAuditID()
	}
	return "audit-" + stableHash(time.Now().UTC().Format(time.RFC3339Nano))
}

func deterministicEvidenceID(prefix string, p PublicEvidencePayload) string {
	return prefix + "-" + stableHash(p.SourceName, p.SourceRecordKey, p.ContentHash)
}

func publicEvidenceEventPayload(p PublicEvidencePayload) PublicEvidencePayload {
	p.SourceName = "event"
	p.SourceRecordKey = publicEvidenceEventID(p)
	p.ContentHash = publicEvidenceEventID(p)
	return p
}

func publicEvidenceEventID(p PublicEvidencePayload) string {
	return stableHash(p.Symbol, p.SourceType, p.EvidenceRole, normalizeEvidenceTitle(p.Title), evidenceDateKey(p.PublishedAt))
}

func publicEvidenceTimeWeight(publishedAt string, now time.Time) float64 {
	published, err := time.Parse(time.RFC3339, strings.TrimSpace(publishedAt))
	if err != nil {
		return 1
	}
	age := now.Sub(published.UTC())
	if age < 0 || age <= 24*time.Hour {
		return 1
	}
	if age <= 7*24*time.Hour {
		return 0.8
	}
	if age <= 30*24*time.Hour {
		return 0.5
	}
	return 0.2
}

func evidenceDateKey(publishedAt string) string {
	if len(publishedAt) >= 10 {
		return publishedAt[:10]
	}
	return publishedAt
}

func normalizeEvidenceTitle(title string) string {
	return strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(title))), " ")
}

func publicEvidenceMetadataJSON(p PublicEvidencePayload, requestID string) string {
	metadata := map[string]any{
		"source_name":    p.SourceName,
		"url":            p.URL,
		"attachment_url": p.AttachmentURL,
		"source_type":    p.SourceType,
		"source_record":  p.SourceRecordKey,
		"raw":            p.Raw,
	}
	if strings.TrimSpace(requestID) != "" {
		metadata["request_id"] = strings.TrimSpace(requestID)
	}
	raw, err := json.Marshal(metadata)
	if err != nil {
		return "{}"
	}
	return string(raw)
}

func highGradeCount(level string) int {
	if level == "S" || level == "A" {
		return 1
	}
	return 0
}

func publicEvidenceVerificationSatisfied(sample PublicEvidencePayload, independentSources int, highGradeSources int) bool {
	if independentSources < 2 {
		return false
	}
	if sample.EvidenceRole != "formal" {
		return false
	}
	switch model.EventType(sample.SourceType) {
	case model.EventMajorPositive, model.EventMajorNegative, model.EventBuyLogicBreak:
		return highGradeSources >= 2
	default:
		return highGradeSources >= 1
	}
}

func sourceLevelRank(level string) int {
	switch level {
	case "S":
		return 4
	case "A":
		return 3
	case "B":
		return 2
	case "C":
		return 1
	default:
		return 0
	}
}
