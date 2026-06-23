package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"regexp"
	"sort"
	"strings"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
	"investment-agent/internal/pkg/clock"
	"investment-agent/internal/pkg/idgen"
)

const localKnowledgeSafetyNote = "本地知识导入仅写入本地背景材料，不接券商、不交易、不外部推送、不自动确认、不自动应用规则。"

var (
	localKnowledgeSecretPattern     = regexp.MustCompile(`sk-[A-Za-z0-9]{12,}`)
	localKnowledgePrivateKeyPattern = regexp.MustCompile(`(?is)BEGIN (RSA|OPENSSH|PRIVATE) KEY.*END (RSA|OPENSSH|PRIVATE) KEY`)
	localKnowledgePrivatePath       = regexp.MustCompile(`/Users/[^\s,;:"']+`)
	localKnowledgeSelectStar        = regexp.MustCompile(`(?is)\bSELECT\s+\*\s+FROM\b[^;\n]*`)
	localKnowledgeRawHTTP           = regexp.MustCompile(`(?i)HTTP/[0-9.]+\s+[0-9]{3}|raw HTTP`)
	localKnowledgePrompt            = regexp.MustCompile(`(?i)prompt:|完整\s*prompt`)
)

type LocalKnowledgeService struct {
	tx  repository.Transactor
	clk clock.Clock
	ids idgen.Generator
}

func NewLocalKnowledgeService(tx repository.Transactor) *LocalKnowledgeService {
	return &LocalKnowledgeService{tx: tx, clk: clock.SystemClock{}, ids: idgen.NewGenerator()}
}

func (s *LocalKnowledgeService) ValidateImport(_ context.Context, req dto.LocalKnowledgeImportValidationRequest) (dto.LocalKnowledgeImportValidationResponse, error) {
	return validateLocalKnowledgeImport(req)
}

func (s *LocalKnowledgeService) ConfirmImport(ctx context.Context, requestID string, req dto.LocalKnowledgeImportConfirmRequest) (dto.LocalKnowledgeImportConfirmResponse, error) {
	if strings.TrimSpace(req.ConfirmReason) == "" {
		return dto.LocalKnowledgeImportConfirmResponse{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "confirm_reason 不能为空")
	}
	validation, err := validateLocalKnowledgeImport(dto.LocalKnowledgeImportValidationRequest{SourceLabel: req.SourceLabel, DefaultSymbol: req.DefaultSymbol, Rows: req.Rows})
	if err != nil {
		return dto.LocalKnowledgeImportConfirmResponse{}, err
	}
	if validation.ImportBatchID != strings.TrimSpace(req.ImportBatchID) {
		return dto.LocalKnowledgeImportConfirmResponse{}, apperr.New(apperr.CodeConflict, apperr.CategoryConflict, "import_batch_id 与当前导入内容不匹配，请重新校验")
	}
	if validation.Summary.BlockingCount > 0 {
		return dto.LocalKnowledgeImportConfirmResponse{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "存在阻塞风险，不能确认导入")
	}
	now := s.clk.NowRFC3339()
	auditID := s.ids.New("audit")
	sourceLabel := normalizeLocalKnowledgeSource(req.SourceLabel)
	err = s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		for i, row := range req.Rows {
			result := validation.Rows[i]
			ids := localKnowledgeIDs(validation.ImportBatchID, result.ContentHash)
			symbol := firstNonEmptyLocal(strings.TrimSpace(row.Symbol), strings.TrimSpace(req.DefaultSymbol))
			tagsJSON := localKnowledgeTagsJSON(row.Tags)
			if err := repos.IntelligenceRepo.SaveIntelligenceItem(ctx, repository.IntelligenceItem{IntelligenceID: ids.IntelligenceID, SourceName: sourceLabel, SourceLevel: "C", OriginalURL: "", PublishedAt: strings.TrimSpace(row.AsOfDate), CapturedAt: now, ContentHash: result.ContentHash, RawTitle: result.TitlePreview, RawTextRef: "local_knowledge_import:" + validation.ImportBatchID, CreatedAt: now}); err != nil {
				return err
			}
			if err := repos.IntelligenceRepo.SaveIntelligenceSummary(ctx, repository.IntelligenceSummary{SummaryID: ids.SummaryID, IntelligenceID: ids.IntelligenceID, Symbol: symbol, Entity: firstNonEmptyLocal(symbol, sourceLabel), EventType: "local_knowledge", ImpactDirection: "neutral", Summary: result.TextPreview, SourceLevel: "C", EvidenceRole: string(model.EvidenceBackground), TimeWeight: 0.3, RelevanceScore: 0.5, VerificationGroupID: ids.VerificationGroupID, CreatedAt: now}, []repository.RAGChunk{{
				ChunkID:      ids.ChunkID,
				SummaryID:    ids.SummaryID,
				Symbol:       symbol,
				ChunkText:    result.TextPreview,
				ChunkHash:    result.ContentHash,
				IndexStatus:  "pending",
				MetadataJSON: localKnowledgeMetadataJSON(validation.ImportBatchID, sourceLabel, tagsJSON),
				CreatedAt:    now,
			}}); err != nil {
				return err
			}
			evidenceIDsJSON, _ := json.Marshal([]string{ids.SummaryID})
			if err := repos.IntelligenceRepo.SaveSourceVerification(ctx, repository.SourceVerification{VerificationID: ids.VerificationID, VerificationGroupID: ids.VerificationGroupID, EventID: ids.EventID, Symbol: symbol, EventType: "local_knowledge", EvidenceRole: string(model.EvidenceBackground), VerificationStatus: string(model.VerificationBackgroundOnly), IndependentSourceCount: 1, HighGradeIndependentSourceCount: 0, HighestSourceLevel: "C", LatestPublishedAt: strings.TrimSpace(row.AsOfDate), EvidenceIDsJSON: string(evidenceIDsJSON), CreatedAt: now}); err != nil {
				return err
			}
		}
		inputRef, _ := json.Marshal(map[string]any{"import_batch_id": validation.ImportBatchID, "row_count": len(req.Rows)})
		return repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: auditID, RequestID: requestID, Actor: string(model.AuditActorUser), Action: string(model.AuditActionRunLocalTask), Status: string(model.AuditStatusSuccess), InputRefType: "local_knowledge_import", InputRef: string(inputRef), CreatedAt: now})
	})
	if err != nil {
		return dto.LocalKnowledgeImportConfirmResponse{}, err
	}
	return dto.LocalKnowledgeImportConfirmResponse{ImportBatchID: validation.ImportBatchID, IntelligenceItemCount: len(req.Rows), SummaryCount: len(req.Rows), RAGChunkCount: len(req.Rows), VerificationCount: len(req.Rows), IndexStatus: "pending", AuditEventIDs: []string{auditID}, SafetyNote: localKnowledgeSafetyNote}, nil
}

func validateLocalKnowledgeImport(req dto.LocalKnowledgeImportValidationRequest) (dto.LocalKnowledgeImportValidationResponse, error) {
	sourceLabel := normalizeLocalKnowledgeSource(req.SourceLabel)
	if sourceLabel == "" {
		return dto.LocalKnowledgeImportValidationResponse{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "source_label 不能为空")
	}
	if len(req.Rows) == 0 {
		return dto.LocalKnowledgeImportValidationResponse{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "rows 不能为空")
	}
	rows := make([]dto.LocalKnowledgeImportRowResult, 0, len(req.Rows))
	validCount := 0
	blockingCount := 0
	warningCount := 0
	totalChunks := 0
	hashParts := []string{sourceLabel, strings.TrimSpace(req.DefaultSymbol)}
	for i, row := range req.Rows {
		result := validateLocalKnowledgeRow(i+1, req.SourceLabel, req.DefaultSymbol, row)
		rows = append(rows, result)
		hashParts = append(hashParts, result.ContentHash)
		totalChunks += result.EstimatedChunk
		if result.Status == "valid" || result.Status == "warning" {
			validCount++
		}
		for _, risk := range result.Risks {
			if risk.Severity == "blocking" {
				blockingCount++
			}
			if risk.Severity == "warning" {
				warningCount++
			}
		}
	}
	return dto.LocalKnowledgeImportValidationResponse{
		ImportBatchID: "lk_import_" + shortHash(strings.Join(hashParts, "\x1f"), 20),
		Summary:       dto.LocalKnowledgeImportValidationSummary{TotalCount: len(req.Rows), ValidCount: validCount, BlockingCount: blockingCount, WarningCount: warningCount},
		Rows:          rows,
		IndexPlan:     dto.LocalKnowledgeImportIndexPlan{RAGChunkCount: totalChunks, IndexStatus: "pending"},
		SafetyNote:    localKnowledgeSafetyNote,
	}, nil
}

func validateLocalKnowledgeRow(rowNumber int, sourceLabel string, defaultSymbol string, row dto.LocalKnowledgeImportRow) dto.LocalKnowledgeImportRowResult {
	title := strings.TrimSpace(row.Title)
	text := strings.TrimSpace(row.Text)
	risks := []dto.LocalKnowledgeImportRisk{}
	if title == "" {
		risks = append(risks, dto.LocalKnowledgeImportRisk{Code: "missing_title", Severity: "blocking", Message: "标题不能为空"})
	}
	if text == "" {
		risks = append(risks, dto.LocalKnowledgeImportRisk{Code: "missing_text", Severity: "blocking", Message: "正文不能为空"})
	}
	combined := strings.Join([]string{sourceLabel, title, text, row.SourceURL, strings.Join(row.Tags, "\n")}, "\n")
	risks = append(risks, localKnowledgeContentRisks(combined)...)
	hasBlocking := false
	hasWarning := false
	for _, risk := range risks {
		if risk.Severity == "blocking" {
			hasBlocking = true
		}
		if risk.Severity == "warning" {
			hasWarning = true
		}
	}
	status := "valid"
	if hasBlocking {
		status = "blocking"
	} else if hasWarning {
		status = "warning"
	}
	symbol := firstNonEmptyLocal(strings.TrimSpace(row.Symbol), strings.TrimSpace(defaultSymbol))
	redactedTitle := localKnowledgeRedact(title)
	redactedText := localKnowledgeRedact(text)
	return dto.LocalKnowledgeImportRowResult{RowNumber: rowNumber, Status: status, Symbol: symbol, TitlePreview: truncateLocalKnowledge(redactedTitle, 80), TextPreview: truncateLocalKnowledge(redactedText, 260), ContentHash: shortHash(strings.Join([]string{symbol, title, text, strings.TrimSpace(row.SourceURL), strings.TrimSpace(row.AsOfDate), strings.Join(sortedTags(row.Tags), ",")}, "\x1f"), 32), EstimatedChunk: estimatedLocalKnowledgeChunks(text), Risks: risks}
}

func localKnowledgeContentRisks(value string) []dto.LocalKnowledgeImportRisk {
	risks := []dto.LocalKnowledgeImportRisk{}
	if localKnowledgeSecretPattern.MatchString(value) || localKnowledgePrivateKeyPattern.MatchString(value) {
		risks = append(risks, dto.LocalKnowledgeImportRisk{Code: "suspected_secret", Severity: "blocking", Message: "疑似密钥或私钥材料"})
	}
	if localKnowledgeSelectStar.MatchString(value) {
		risks = append(risks, dto.LocalKnowledgeImportRisk{Code: "raw_sql", Severity: "blocking", Message: "疑似原始 SQL"})
	}
	if localKnowledgePrivatePath.MatchString(value) {
		risks = append(risks, dto.LocalKnowledgeImportRisk{Code: "private_path", Severity: "blocking", Message: "疑似本地私有路径，预览已脱敏"})
	}
	if localKnowledgeRawHTTP.MatchString(value) {
		risks = append(risks, dto.LocalKnowledgeImportRisk{Code: "raw_http", Severity: "blocking", Message: "疑似原始 HTTP 响应，需改为摘要"})
	}
	if localKnowledgePrompt.MatchString(value) {
		risks = append(risks, dto.LocalKnowledgeImportRisk{Code: "full_prompt", Severity: "blocking", Message: "疑似完整 prompt，需改为摘要"})
	}
	return risks
}

func localKnowledgeRedact(value string) string {
	out := localKnowledgePrivateKeyPattern.ReplaceAllString(value, "[REDACTED_PRIVATE_KEY]")
	if localKnowledgeRawHTTP.MatchString(out) {
		out = "[REDACTED_HTTP]"
	}
	if localKnowledgePrompt.MatchString(out) {
		out = "[REDACTED_PROMPT]"
	}
	out = localKnowledgeSecretPattern.ReplaceAllString(out, "[REDACTED_KEY]")
	out = localKnowledgeSelectStar.ReplaceAllString(out, "[REDACTED_SQL]")
	out = localKnowledgePrivatePath.ReplaceAllString(out, "/Users/[REDACTED]")
	return out
}

func shortHash(value string, length int) string {
	sum := sha256.Sum256([]byte(value))
	out := hex.EncodeToString(sum[:])
	if length > 0 && length < len(out) {
		return out[:length]
	}
	return out
}

type localKnowledgeIDSet struct {
	IntelligenceID      string
	SummaryID           string
	ChunkID             string
	VerificationID      string
	VerificationGroupID string
	EventID             string
}

func localKnowledgeIDs(batchID string, contentHash string) localKnowledgeIDSet {
	seed := shortHash(batchID+":"+contentHash, 20)
	return localKnowledgeIDSet{IntelligenceID: "intel_lk_" + seed, SummaryID: "sum_lk_" + seed, ChunkID: "chunk_lk_" + seed, VerificationID: "sv_lk_" + seed, VerificationGroupID: "vg_lk_" + seed, EventID: "event_lk_" + seed}
}

func normalizeLocalKnowledgeSource(value string) string {
	return strings.TrimSpace(value)
}

func truncateLocalKnowledge(value string, max int) string {
	value = strings.TrimSpace(value)
	if len([]rune(value)) <= max {
		return value
	}
	runes := []rune(value)
	return string(runes[:max])
}

func estimatedLocalKnowledgeChunks(text string) int {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0
	}
	runeCount := len([]rune(text))
	chunks := runeCount / 500
	if runeCount%500 != 0 {
		chunks++
	}
	if chunks < 1 {
		return 1
	}
	return chunks
}

func firstNonEmptyLocal(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func sortedTags(tags []string) []string {
	out := make([]string, 0, len(tags))
	for _, tag := range tags {
		if strings.TrimSpace(tag) != "" {
			out = append(out, strings.TrimSpace(tag))
		}
	}
	sort.Strings(out)
	return out
}

func localKnowledgeTagsJSON(tags []string) string {
	data, _ := json.Marshal(sortedTags(tags))
	return string(data)
}

func localKnowledgeMetadataJSON(batchID, sourceLabel, tagsJSON string) string {
	var tags []string
	_ = json.Unmarshal([]byte(tagsJSON), &tags)
	data, _ := json.Marshal(map[string]any{"source_type": "local_knowledge_import", "import_batch_id": batchID, "source_label": sourceLabel, "tags": tags})
	return string(data)
}
