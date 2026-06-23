package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
	"investment-agent/internal/pkg/idgen"
)

const (
	DataSourceQualityModeFixture = "fixture"
	DataSourceQualityModeCurrent = "current"

	DataSourceQualityStatusPassed   = "passed"
	DataSourceQualityStatusDegraded = "degraded"
	DataSourceQualityStatusFailed   = "failed"

	DataSourceQualityPolicyPassed         = "passed"
	DataSourceQualityPolicyWaiverRequired = "waiver_required"
	DataSourceQualityPolicyBlocked        = "blocked"

	DataSourceQualityReleaseGatePass           = "pass"
	DataSourceQualityReleaseGateWaiverRequired = "waiver_required"
	DataSourceQualityReleaseGateBlock          = "block"

	dataSourceQualitySafetyNote         = "本地数据源质量回归只检查分类和脱敏摘要，不刷新数据、不修改规则、不改变账户事实。"
	dataSourceQualityPolicySafetyNote   = "当前数据质量策略只读取本地 source health 并给出发布影响，不刷新数据、不执行修复动作、不调用外部源、不触发交易。"
	dataQualityGateResolutionSafetyNote = "当前数据门禁处置只记录本地人工声明边界，不改变数据质量事实、不刷新数据、不触发交易。"

	DataQualityGateResolutionTypeWaiver         = "waiver"
	DataQualityGateResolutionTypeScopeExclusion = "scope_exclusion"
	DataQualityGateResolutionStatusActive       = "active"
	DataQualityGateResolutionStatusRetired      = "retired"

	DataQualityReleaseClaimPass                       = "pass"
	DataQualityReleaseClaimRequiresResolution         = "requires_resolution"
	DataQualityReleaseClaimResolvedWithWaiver         = "resolved_with_waiver"
	DataQualityReleaseClaimResolvedWithScopeExclusion = "resolved_with_scope_exclusion"
)

var sourceQualitySQLFromPattern = regexp.MustCompile(`(?i)\bFROM\s+[^\s,;，；。]+`)

// DataSourceQualityRegressionRequest bounds a local source-quality regression run.
type DataSourceQualityRegressionRequest struct {
	Mode   string
	Symbol string
}

type DataQualityGateResolutionCheckRequest struct {
	Symbol string
}

type DataQualityGateResolutionCreateRequest struct {
	RequestID      string
	Symbol         string
	ResolutionType string
	Scope          string
	Reason         string
	ReleaseImpact  string
	EvidenceRef    string
}

type DataQualityGateResolutionListRequest struct {
	Symbol string
	Status string
}

// DataSourceQualityService runs local data-source quality regression checks.
type DataSourceQualityService struct {
	repos repository.Repositories
	tx    repository.Transactor
	now   func() time.Time
}

// NewDataSourceQualityService creates a read-oriented data source quality service.
func NewDataSourceQualityService(repos repository.Repositories, tx ...repository.Transactor) *DataSourceQualityService {
	var transactor repository.Transactor
	if len(tx) > 0 {
		transactor = tx[0]
	}
	return &DataSourceQualityService{repos: repos, tx: transactor, now: func() time.Time { return time.Now().UTC() }}
}

// Run executes fixture or current source-health regression.
func (s *DataSourceQualityService) Run(ctx context.Context, req DataSourceQualityRegressionRequest) (dto.DataSourceQualityRegressionResponse, error) {
	mode := strings.TrimSpace(req.Mode)
	if mode == "" {
		mode = DataSourceQualityModeFixture
	}
	switch mode {
	case DataSourceQualityModeFixture:
		return s.regressionFromCases(mode, fixtureDataSourceQualityCases(), nil), nil
	case DataSourceQualityModeCurrent:
		items, missing, err := s.currentSourceHealth(ctx, strings.TrimSpace(req.Symbol))
		if err != nil {
			return dto.DataSourceQualityRegressionResponse{}, err
		}
		return s.regressionFromSourceHealth(mode, items, missing), nil
	default:
		return dto.DataSourceQualityRegressionResponse{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "unsupported data source quality regression mode")
	}
}

func (s *DataSourceQualityService) CheckGateResolution(ctx context.Context, req DataQualityGateResolutionCheckRequest) (dto.DataQualityGateResolutionCheck, error) {
	symbol := strings.TrimSpace(req.Symbol)
	regression, err := s.Run(ctx, DataSourceQualityRegressionRequest{Mode: DataSourceQualityModeCurrent, Symbol: symbol})
	if err != nil {
		return dto.DataQualityGateResolutionCheck{}, err
	}
	if symbol == "" {
		symbol = firstAffectedSymbol(regression)
	}
	fingerprint := dataQualityPolicyFingerprint(symbol, regression)
	out := dto.DataQualityGateResolutionCheck{
		Symbol:                symbol,
		PolicyFingerprint:     fingerprint,
		PolicySummary:         DataSourceQualityAuditOutputRef(regression),
		Policy:                regression.Policy,
		ReleaseClaimState:     DataQualityReleaseClaimRequiresResolution,
		CleanDataClaimAllowed: false,
		AllowedClaims:         dataQualityAllowedClaims(DataQualityReleaseClaimRequiresResolution),
		ProhibitedClaims:      dataQualityProhibitedClaims(false),
		SafetyNote:            dataQualityGateResolutionSafetyNote,
	}
	if regression.Policy.Verdict == DataSourceQualityPolicyPassed {
		out.ReleaseClaimState = DataQualityReleaseClaimPass
		out.CleanDataClaimAllowed = true
		out.AllowedClaims = dataQualityAllowedClaims(out.ReleaseClaimState)
		out.ProhibitedClaims = dataQualityProhibitedClaims(true)
		return out, nil
	}
	if s.repos.DataQualityGateResolutionRepo == nil {
		return out, nil
	}
	active, err := s.repos.DataQualityGateResolutionRepo.GetActiveDataQualityGateResolution(ctx, symbol, fingerprint)
	if err != nil {
		if apperr.IsCode(err, apperr.CodeNotFound) {
			return out, nil
		}
		return dto.DataQualityGateResolutionCheck{}, err
	}
	if !resolutionTypeAllowedForPolicy(active.ResolutionType, regression.Policy.Verdict) {
		return out, nil
	}
	record := dataQualityGateResolutionDTO(active)
	out.ActiveResolution = &record
	switch active.ResolutionType {
	case DataQualityGateResolutionTypeWaiver:
		out.ReleaseClaimState = DataQualityReleaseClaimResolvedWithWaiver
	case DataQualityGateResolutionTypeScopeExclusion:
		out.ReleaseClaimState = DataQualityReleaseClaimResolvedWithScopeExclusion
	}
	out.AllowedClaims = dataQualityAllowedClaims(out.ReleaseClaimState)
	out.ProhibitedClaims = dataQualityProhibitedClaims(false)
	return out, nil
}

func (s *DataSourceQualityService) CreateGateResolution(ctx context.Context, req DataQualityGateResolutionCreateRequest) (dto.DataQualityGateResolutionCheck, error) {
	if s.tx != nil {
		var out dto.DataQualityGateResolutionCheck
		err := s.tx.WithinTx(ctx, func(txCtx context.Context, txRepos repository.Repositories) error {
			txSvc := &DataSourceQualityService{repos: txRepos, now: s.now}
			next, err := txSvc.createGateResolutionNoTx(txCtx, req)
			if err != nil {
				return err
			}
			out = next
			return nil
		})
		return out, err
	}
	return s.createGateResolutionNoTx(ctx, req)
}

func (s *DataSourceQualityService) createGateResolutionNoTx(ctx context.Context, req DataQualityGateResolutionCreateRequest) (dto.DataQualityGateResolutionCheck, error) {
	check, err := s.CheckGateResolution(ctx, DataQualityGateResolutionCheckRequest{Symbol: req.Symbol})
	if err != nil {
		return dto.DataQualityGateResolutionCheck{}, err
	}
	if check.Policy.Verdict == DataSourceQualityPolicyPassed {
		return dto.DataQualityGateResolutionCheck{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "current data policy already passes")
	}
	if s.repos.DataQualityGateResolutionRepo == nil {
		return dto.DataQualityGateResolutionCheck{}, apperr.New(apperr.CodeInternalError, apperr.CategoryInternal, "data quality gate resolution repository missing")
	}
	resolutionType := strings.TrimSpace(req.ResolutionType)
	if !resolutionTypeAllowedForPolicy(resolutionType, check.Policy.Verdict) {
		return dto.DataQualityGateResolutionCheck{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "resolution type is not allowed for current policy")
	}
	scope := sanitizeDataSourceQualityText(req.Scope)
	reason := sanitizeDataSourceQualityText(req.Reason)
	releaseImpact := sanitizeDataSourceQualityText(req.ReleaseImpact)
	evidenceRef := sanitizeDataSourceQualityText(req.EvidenceRef)
	if strings.TrimSpace(scope) == "" || strings.TrimSpace(reason) == "" || strings.TrimSpace(releaseImpact) == "" {
		return dto.DataQualityGateResolutionCheck{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "resolution scope, reason and release impact are required")
	}
	if check.ActiveResolution != nil {
		if check.ActiveResolution.ResolutionType == resolutionType {
			if err := s.appendDataQualityGateResolutionAudit(ctx, req.RequestID, "create", dataQualityGateResolutionRepositoryRecord(*check.ActiveResolution), check.ReleaseClaimState); err != nil {
				return dto.DataQualityGateResolutionCheck{}, err
			}
			return check, nil
		}
		return dto.DataQualityGateResolutionCheck{}, apperr.New(apperr.CodeConflict, apperr.CategoryConflict, "active resolution already exists for current policy")
	}
	blockingJSON, err := json.Marshal(check.Policy.BlockingReasons)
	if err != nil {
		return dto.DataQualityGateResolutionCheck{}, err
	}
	waiverJSON, err := json.Marshal(check.Policy.WaiverReasons)
	if err != nil {
		return dto.DataQualityGateResolutionCheck{}, err
	}
	now := s.now().UTC().Format(time.RFC3339)
	record := repository.DataQualityGateResolution{
		ResolutionID:        idgen.NewGenerator().New("dqgr"),
		Symbol:              check.Symbol,
		PolicyFingerprint:   check.PolicyFingerprint,
		PolicyVerdict:       check.Policy.Verdict,
		ReleaseGate:         check.Policy.ReleaseGate,
		PolicySummary:       check.PolicySummary,
		ResolutionType:      resolutionType,
		Status:              DataQualityGateResolutionStatusActive,
		Scope:               scope,
		Reason:              reason,
		ReleaseImpact:       releaseImpact,
		EvidenceRef:         evidenceRef,
		BlockingReasonsJSON: string(blockingJSON),
		WaiverReasonsJSON:   string(waiverJSON),
		CreatedBy:           "local_user",
		CreatedAt:           now,
		SafetyNote:          dataQualityGateResolutionSafetyNote,
	}
	if err := s.repos.DataQualityGateResolutionRepo.CreateDataQualityGateResolution(ctx, record); err != nil {
		return dto.DataQualityGateResolutionCheck{}, err
	}
	out, err := s.CheckGateResolution(ctx, DataQualityGateResolutionCheckRequest{Symbol: check.Symbol})
	if err != nil {
		return dto.DataQualityGateResolutionCheck{}, err
	}
	if err := s.appendDataQualityGateResolutionAudit(ctx, req.RequestID, "create", record, out.ReleaseClaimState); err != nil {
		return dto.DataQualityGateResolutionCheck{}, err
	}
	return out, nil
}

func (s *DataSourceQualityService) RetireGateResolution(ctx context.Context, resolutionID string, requestIDs ...string) (dto.DataQualityGateResolutionCheck, error) {
	if s.tx != nil {
		var out dto.DataQualityGateResolutionCheck
		err := s.tx.WithinTx(ctx, func(txCtx context.Context, txRepos repository.Repositories) error {
			txSvc := &DataSourceQualityService{repos: txRepos, now: s.now}
			next, err := txSvc.retireGateResolutionNoTx(txCtx, resolutionID, requestIDs...)
			if err != nil {
				return err
			}
			out = next
			return nil
		})
		return out, err
	}
	return s.retireGateResolutionNoTx(ctx, resolutionID, requestIDs...)
}

func (s *DataSourceQualityService) retireGateResolutionNoTx(ctx context.Context, resolutionID string, requestIDs ...string) (dto.DataQualityGateResolutionCheck, error) {
	if s.repos.DataQualityGateResolutionRepo == nil {
		return dto.DataQualityGateResolutionCheck{}, apperr.New(apperr.CodeInternalError, apperr.CategoryInternal, "data quality gate resolution repository missing")
	}
	current, err := s.repos.DataQualityGateResolutionRepo.GetDataQualityGateResolution(ctx, strings.TrimSpace(resolutionID))
	if err != nil {
		return dto.DataQualityGateResolutionCheck{}, err
	}
	if err := s.repos.DataQualityGateResolutionRepo.RetireDataQualityGateResolution(ctx, current.ResolutionID, "local_user", s.now().UTC().Format(time.RFC3339)); err != nil {
		return dto.DataQualityGateResolutionCheck{}, err
	}
	out, err := s.CheckGateResolution(ctx, DataQualityGateResolutionCheckRequest{Symbol: current.Symbol})
	if err != nil {
		return dto.DataQualityGateResolutionCheck{}, err
	}
	requestID := ""
	if len(requestIDs) > 0 {
		requestID = requestIDs[0]
	}
	if err := s.appendDataQualityGateResolutionAudit(ctx, requestID, "retire", current, out.ReleaseClaimState); err != nil {
		return dto.DataQualityGateResolutionCheck{}, err
	}
	return out, nil
}

func (s *DataSourceQualityService) ListGateResolutions(ctx context.Context, req DataQualityGateResolutionListRequest) ([]dto.DataQualityGateResolutionRecord, error) {
	if s.repos.DataQualityGateResolutionRepo == nil {
		return nil, apperr.New(apperr.CodeInternalError, apperr.CategoryInternal, "data quality gate resolution repository missing")
	}
	items, err := s.repos.DataQualityGateResolutionRepo.ListDataQualityGateResolutions(ctx, repository.DataQualityGateResolutionFilter{
		Symbol: strings.TrimSpace(req.Symbol),
		Status: strings.TrimSpace(req.Status),
	})
	if err != nil {
		return nil, err
	}
	out := make([]dto.DataQualityGateResolutionRecord, 0, len(items))
	for _, item := range items {
		out = append(out, dataQualityGateResolutionDTO(item))
	}
	return out, nil
}

func (s *DataSourceQualityService) currentSourceHealth(ctx context.Context, symbol string) ([]dto.SourceHealthItem, []string, error) {
	if s.repos.MarketRepo == nil {
		return nil, []string{"p34_source_health"}, nil
	}
	var (
		market model.MarketSnapshot
		err    error
	)
	if symbol != "" {
		market, err = s.repos.MarketRepo.GetLatestMarketSnapshotBySymbol(ctx, symbol)
	} else {
		market, err = s.repos.MarketRepo.GetLatestMarketSnapshot(ctx)
	}
	if err != nil {
		if apperr.IsCode(err, apperr.CodeNotFound) {
			return nil, []string{"p34_source_health"}, nil
		}
		return nil, nil, err
	}
	items := SourceHealthFromMarketSnapshot(market)
	if len(items) == 0 {
		return nil, []string{"p34_source_health"}, nil
	}
	return items, nil, nil
}

func (s *DataSourceQualityService) regressionFromSourceHealth(mode string, items []dto.SourceHealthItem, missing []string) dto.DataSourceQualityRegressionResponse {
	if len(items) == 0 {
		return s.regressionFromCases(mode, []dto.DataSourceQualityCase{{
			CaseID:            "current_source_health_missing",
			DataCategory:      "p34_source_health",
			ExpectedFreshness: "fresh",
			ActualFreshness:   "missing",
			Status:            DataSourceQualityStatusDegraded,
			DiagnosticPreview: sanitizeDataSourceQualityText("未找到可评估的 P34 source health"),
		}}, missing)
	}
	cases := make([]dto.DataSourceQualityCase, 0, len(items))
	missingSet := map[string]struct{}{}
	for _, item := range items {
		freshness := strings.TrimSpace(item.Freshness)
		if freshness == "" {
			freshness = "missing"
		}
		status := DataSourceQualityStatusDegraded
		diagnostic := fmt.Sprintf("source health %s freshness=%s failure=%s", item.DataCategory, freshness, item.FailureCategory)
		if freshness == "fresh" || freshness == "stubbed" {
			status = DataSourceQualityStatusPassed
		}
		if !recognizedSourceFreshness(freshness) {
			status = DataSourceQualityStatusFailed
			diagnostic = "unrecognized freshness: " + freshness
		}
		if status != DataSourceQualityStatusPassed {
			missingSet[firstNonEmptySourceQuality(strings.TrimSpace(item.DataCategory), "p34_source_health")] = struct{}{}
		}
		cases = append(cases, dto.DataSourceQualityCase{
			CaseID:            firstNonEmptySourceQuality(strings.TrimSpace(item.DataCategory), "p34_source_health"),
			SourceName:        item.SourceName,
			SourceLevel:       item.SourceLevel,
			SourceType:        item.SourceType,
			DataCategory:      item.DataCategory,
			ExpectedFreshness: "fresh",
			ActualFreshness:   freshness,
			Status:            status,
			DataDate:          item.DataDate,
			FailureCategory:   item.FailureCategory,
			AffectedSymbols:   item.AffectedSymbols,
			DiagnosticPreview: sanitizeDataSourceQualityText(diagnostic),
		})
	}
	for _, item := range missing {
		if strings.TrimSpace(item) != "" {
			missingSet[strings.TrimSpace(item)] = struct{}{}
		}
	}
	missingCategories := sortedStringSet(missingSet)
	return s.regressionFromCases(mode, cases, missingCategories)
}

func (s *DataSourceQualityService) regressionFromCases(mode string, cases []dto.DataSourceQualityCase, missing []string) dto.DataSourceQualityRegressionResponse {
	status := DataSourceQualityStatusPassed
	degraded, failed := 0, 0
	for i := range cases {
		cases[i].DiagnosticPreview = sanitizeDataSourceQualityText(cases[i].DiagnosticPreview)
		switch cases[i].Status {
		case DataSourceQualityStatusFailed:
			failed++
			status = DataSourceQualityStatusFailed
		case DataSourceQualityStatusDegraded:
			degraded++
			if status != DataSourceQualityStatusFailed {
				status = DataSourceQualityStatusDegraded
			}
		}
	}
	if len(missing) > 0 && status == DataSourceQualityStatusPassed {
		status = DataSourceQualityStatusDegraded
	}
	return dto.DataSourceQualityRegressionResponse{
		Mode:              mode,
		Status:            status,
		GeneratedAt:       s.now().UTC().Format(time.RFC3339),
		Summary:           dataSourceQualitySummary(mode, status, len(cases), degraded, failed),
		Cases:             cases,
		MissingCategories: missing,
		Policy:            dataSourceQualityPolicy(mode, cases, degraded, failed),
		SafetyNote:        dataSourceQualitySafetyNote,
	}
}

func dataSourceQualityPolicy(mode string, cases []dto.DataSourceQualityCase, degraded int, failed int) dto.DataSourceQualityPolicy {
	policy := dto.DataSourceQualityPolicy{
		Verdict:       DataSourceQualityPolicyPassed,
		ReleaseGate:   DataSourceQualityReleaseGatePass,
		DegradedCount: degraded,
		FailedCount:   failed,
		NextActions: []string{
			"保留当前只读数据质量证据",
			"发布材料引用 policy verdict 和 release gate",
		},
		SafetyNote: dataSourceQualityPolicySafetyNote,
	}
	if mode != DataSourceQualityModeCurrent {
		return policy
	}
	if len(cases) == 0 {
		policy.BlockingReasons = append(policy.BlockingReasons, "未找到可评估的 source health facts")
	}
	for _, item := range cases {
		category := firstNonEmptySourceQuality(item.DataCategory, item.CaseID, "p34_source_health")
		freshness := firstNonEmptySourceQuality(item.ActualFreshness, "missing")
		failureCategory := strings.TrimSpace(item.FailureCategory)
		switch {
		case item.Status == DataSourceQualityStatusFailed:
			policy.BlockingReasons = append(policy.BlockingReasons, sanitizeDataSourceQualityText(fmt.Sprintf("%s failed freshness=%s", category, freshness)))
		case freshness == "missing":
			policy.BlockingReasons = append(policy.BlockingReasons, sanitizeDataSourceQualityText(fmt.Sprintf("%s missing source health", category)))
		case !recognizedSourceFreshness(freshness):
			policy.BlockingReasons = append(policy.BlockingReasons, sanitizeDataSourceQualityText(fmt.Sprintf("%s unrecognized freshness=%s", category, freshness)))
		case !recognizedSourceFailureCategory(failureCategory):
			policy.BlockingReasons = append(policy.BlockingReasons, sanitizeDataSourceQualityText(fmt.Sprintf("%s unrecognized failure_category=%s", category, failureCategory)))
		case item.Status == DataSourceQualityStatusDegraded && isCoreSourceQualityCategory(item):
			policy.BlockingReasons = append(policy.BlockingReasons, sanitizeDataSourceQualityText(fmt.Sprintf("%s core category degraded freshness=%s", category, freshness)))
		case item.Status == DataSourceQualityStatusDegraded:
			policy.WaiverReasons = append(policy.WaiverReasons, sanitizeDataSourceQualityText(fmt.Sprintf("%s optional category degraded freshness=%s", category, freshness)))
		}
	}
	policy.BlockingCount = len(policy.BlockingReasons)
	policy.WaiverCount = len(policy.WaiverReasons)
	if policy.BlockingCount > 0 {
		policy.Verdict = DataSourceQualityPolicyBlocked
		policy.ReleaseGate = DataSourceQualityReleaseGateBlock
		policy.NextActions = []string{
			"发布前处理 blocking source health 或在发布范围中明确排除当前本地数据健康声明",
			"重新运行 current data-source quality policy gate",
			"不得把当前数据源质量声明为 clean",
		}
		return policy
	}
	if policy.WaiverCount > 0 {
		policy.Verdict = DataSourceQualityPolicyWaiverRequired
		policy.ReleaseGate = DataSourceQualityReleaseGateWaiverRequired
		policy.NextActions = []string{
			"在发布材料中记录 waiver reason 和影响范围",
			"不得把 waiver_required 描述为 clean pass",
			"需要 clean claim 时先补充人工处理并重跑 current policy gate",
		}
	}
	return policy
}

// DataSourceQualityAuditOutputRef returns a compact sanitized audit output reference.
func DataSourceQualityAuditOutputRef(out dto.DataSourceQualityRegressionResponse) string {
	degraded, failed := 0, 0
	for _, item := range out.Cases {
		switch item.Status {
		case DataSourceQualityStatusDegraded:
			degraded++
		case DataSourceQualityStatusFailed:
			failed++
		}
	}
	return strings.Join([]string{
		"data_source_quality",
		"mode=" + firstNonEmptySourceQuality(out.Mode, DataSourceQualityModeFixture),
		"status=" + firstNonEmptySourceQuality(out.Status, DataSourceQualityStatusFailed),
		"policy=" + firstNonEmptySourceQuality(out.Policy.Verdict, DataSourceQualityPolicyBlocked),
		"gate=" + firstNonEmptySourceQuality(out.Policy.ReleaseGate, DataSourceQualityReleaseGateBlock),
		fmt.Sprintf("cases=%d", len(out.Cases)),
		fmt.Sprintf("degraded=%d", degraded),
		fmt.Sprintf("failed=%d", failed),
		"no_auto_trading",
	}, ":")
}

func (s *DataSourceQualityService) appendDataQualityGateResolutionAudit(ctx context.Context, requestID string, action string, record repository.DataQualityGateResolution, claimState string) error {
	if s.repos.AuditRepo == nil {
		return nil
	}
	output := strings.Join([]string{
		"data_quality_gate_resolution",
		"action=" + firstNonEmptySourceQuality(action, "unknown"),
		"symbol=" + sanitizeDataSourceQualityText(record.Symbol),
		"policy=" + sanitizeDataSourceQualityText(record.PolicyVerdict),
		"gate=" + sanitizeDataSourceQualityText(record.ReleaseGate),
		"resolution=" + sanitizeDataSourceQualityText(record.ResolutionType),
		"claim_state=" + sanitizeDataSourceQualityText(claimState),
		"clean_data_claim=false",
		"no_auto_trading",
	}, ":")
	return s.repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{
		AuditEventID: idgen.NewGenerator().New("audit"),
		RequestID:    strings.TrimSpace(requestID),
		WorkflowType: "data-source-quality-gate-resolution",
		NodeName:     "DataQualityGateResolutionService",
		Actor:        "user",
		Action:       "run_local_task",
		NodeAction:   "data_quality_gate_resolution_" + action,
		Status:       "success",
		InputRefType: "data_quality_gate_resolution",
		InputRef: strings.Join([]string{
			"data_quality_gate_resolution",
			"symbol=" + sanitizeDataSourceQualityText(record.Symbol),
			"fingerprint=" + sanitizeDataSourceQualityText(record.PolicyFingerprint),
			"resolution_id=" + sanitizeDataSourceQualityText(record.ResolutionID),
		}, ":"),
		OutputRefType: "data_quality_gate_resolution",
		OutputRef:     output,
		CreatedAt:     s.now().UTC().Format(time.RFC3339),
	})
}

func dataQualityPolicyFingerprint(symbol string, out dto.DataSourceQualityRegressionResponse) string {
	type fingerprintCase struct {
		CaseID          string `json:"case_id"`
		DataCategory    string `json:"data_category"`
		ActualFreshness string `json:"actual_freshness"`
		Status          string `json:"status"`
		FailureCategory string `json:"failure_category"`
	}
	cases := make([]fingerprintCase, 0, len(out.Cases))
	for _, item := range out.Cases {
		cases = append(cases, fingerprintCase{
			CaseID:          strings.TrimSpace(item.CaseID),
			DataCategory:    strings.TrimSpace(item.DataCategory),
			ActualFreshness: strings.TrimSpace(item.ActualFreshness),
			Status:          strings.TrimSpace(item.Status),
			FailureCategory: strings.TrimSpace(item.FailureCategory),
		})
	}
	sort.Slice(cases, func(i, j int) bool {
		left := cases[i].DataCategory + ":" + cases[i].CaseID
		right := cases[j].DataCategory + ":" + cases[j].CaseID
		return left < right
	})
	blocking := append([]string(nil), out.Policy.BlockingReasons...)
	waiver := append([]string(nil), out.Policy.WaiverReasons...)
	sort.Strings(blocking)
	sort.Strings(waiver)
	payload := struct {
		Symbol          string            `json:"symbol"`
		Verdict         string            `json:"verdict"`
		ReleaseGate     string            `json:"release_gate"`
		DegradedCount   int               `json:"degraded_count"`
		FailedCount     int               `json:"failed_count"`
		BlockingCount   int               `json:"blocking_count"`
		WaiverCount     int               `json:"waiver_count"`
		BlockingReasons []string          `json:"blocking_reasons"`
		WaiverReasons   []string          `json:"waiver_reasons"`
		Cases           []fingerprintCase `json:"cases"`
	}{
		Symbol:          strings.TrimSpace(symbol),
		Verdict:         out.Policy.Verdict,
		ReleaseGate:     out.Policy.ReleaseGate,
		DegradedCount:   out.Policy.DegradedCount,
		FailedCount:     out.Policy.FailedCount,
		BlockingCount:   out.Policy.BlockingCount,
		WaiverCount:     out.Policy.WaiverCount,
		BlockingReasons: blocking,
		WaiverReasons:   waiver,
		Cases:           cases,
	}
	data, _ := json.Marshal(payload)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func resolutionTypeAllowedForPolicy(resolutionType string, policyVerdict string) bool {
	switch strings.TrimSpace(policyVerdict) {
	case DataSourceQualityPolicyWaiverRequired:
		return resolutionType == DataQualityGateResolutionTypeWaiver || resolutionType == DataQualityGateResolutionTypeScopeExclusion
	case DataSourceQualityPolicyBlocked:
		return resolutionType == DataQualityGateResolutionTypeScopeExclusion
	default:
		return false
	}
}

func dataQualityGateResolutionDTO(item repository.DataQualityGateResolution) dto.DataQualityGateResolutionRecord {
	return dto.DataQualityGateResolutionRecord{
		ResolutionID:      item.ResolutionID,
		Symbol:            item.Symbol,
		PolicyFingerprint: item.PolicyFingerprint,
		PolicyVerdict:     item.PolicyVerdict,
		ReleaseGate:       item.ReleaseGate,
		PolicySummary:     item.PolicySummary,
		ResolutionType:    item.ResolutionType,
		Status:            item.Status,
		Scope:             item.Scope,
		Reason:            item.Reason,
		ReleaseImpact:     item.ReleaseImpact,
		EvidenceRef:       item.EvidenceRef,
		CreatedBy:         item.CreatedBy,
		RetiredBy:         item.RetiredBy,
		CreatedAt:         item.CreatedAt,
		RetiredAt:         item.RetiredAt,
		SafetyNote:        item.SafetyNote,
	}
}

func dataQualityGateResolutionRepositoryRecord(item dto.DataQualityGateResolutionRecord) repository.DataQualityGateResolution {
	return repository.DataQualityGateResolution{
		ResolutionID:      item.ResolutionID,
		Symbol:            item.Symbol,
		PolicyFingerprint: item.PolicyFingerprint,
		PolicyVerdict:     item.PolicyVerdict,
		ReleaseGate:       item.ReleaseGate,
		PolicySummary:     item.PolicySummary,
		ResolutionType:    item.ResolutionType,
		Status:            item.Status,
		Scope:             item.Scope,
		Reason:            item.Reason,
		ReleaseImpact:     item.ReleaseImpact,
		EvidenceRef:       item.EvidenceRef,
		CreatedBy:         item.CreatedBy,
		RetiredBy:         item.RetiredBy,
		CreatedAt:         item.CreatedAt,
		RetiredAt:         item.RetiredAt,
		SafetyNote:        item.SafetyNote,
	}
}

func dataQualityAllowedClaims(state string) []string {
	switch state {
	case DataQualityReleaseClaimPass:
		return []string{"可以声明当前本地数据门禁通过"}
	case DataQualityReleaseClaimResolvedWithWaiver:
		return []string{"可以声明已记录当前数据质量豁免"}
	case DataQualityReleaseClaimResolvedWithScopeExclusion:
		return []string{"可以声明当前本地数据健康已排除在 clean claim 外"}
	default:
		return []string{"需要先记录本地处置或让当前数据策略通过"}
	}
}

func dataQualityProhibitedClaims(cleanAllowed bool) []string {
	if cleanAllowed {
		return nil
	}
	return []string{
		"不得声明当前本地数据 clean",
		"不得声明 current data healthy",
		"不得把 resolution 描述为 policy passed",
	}
}

func firstAffectedSymbol(out dto.DataSourceQualityRegressionResponse) string {
	for _, item := range out.Cases {
		if len(item.AffectedSymbols) > 0 && strings.TrimSpace(item.AffectedSymbols[0]) != "" {
			return strings.TrimSpace(item.AffectedSymbols[0])
		}
	}
	return ""
}

func fixtureDataSourceQualityCases() []dto.DataSourceQualityCase {
	return []dto.DataSourceQualityCase{
		fixtureQualityCase("fresh_index", "csindex", "A", "index_basic", "index_constituents", "fresh", "000300", "2026-06-05", "公开指数样本 freshness=fresh"),
		fixtureQualityCase("no_data_window", "csindex", "A", "index_basic", "index_weights", "no_data", "000300", "2026-06-05", "公开指数样本窗口无记录 freshness=no_data"),
		fixtureQualityCase("source_unavailable", "csindex", "A", "index_basic", "index_valuation_files", "source_unavailable", "000300", "2026-06-05", "公开指数样本请求失败 freshness=source_unavailable"),
		fixtureQualityCase("parse_error", "eastmoney_fund", "B", "fund_nav", "fund_profile", "parse_error", "510300", "2026-06-05", "公开基金样本结构不兼容 freshness=parse_error"),
		fixtureQualityCase("stale", "sentiment_proxy_fixture", "C", "sentiment_proxy", "sentiment_proxy", "stale", "510300", "2026-06-01", "情绪替代样本 freshness=stale"),
		fixtureQualityCase("redaction", "local_fixture", "C", "diagnostic", "diagnostic_redaction", "fresh", "000300", "2026-06-05", redactionFixtureDiagnostic()),
	}
}

func fixtureQualityCase(caseID, sourceName, sourceLevel, sourceType, category, freshness, symbol, dataDate, diagnostic string) dto.DataSourceQualityCase {
	failure := ""
	if freshness != "fresh" && freshness != "stubbed" {
		failure = freshness
	}
	return dto.DataSourceQualityCase{
		CaseID:            caseID,
		SourceName:        sourceName,
		SourceLevel:       sourceLevel,
		SourceType:        sourceType,
		DataCategory:      category,
		ExpectedFreshness: freshness,
		ActualFreshness:   freshness,
		Status:            DataSourceQualityStatusPassed,
		DataDate:          dataDate,
		FailureCategory:   failure,
		AffectedSymbols:   []string{symbol},
		DiagnosticPreview: sanitizeDataSourceQualityText(diagnostic),
	}
}

func redactionFixtureDiagnostic() string {
	return strings.Join([]string{
		"s" + "k-123456789012",
		"s" + "k-proj-abc_def-123456",
		"/" + "Users/private/local.txt",
		"select    " + "*    from secret",
		"prompt" + ": full content",
		"raw " + "HTTP",
		"GET /secret HTTP/1.1",
		"HTTP/" + "1.1 500 Internal Server Error",
		"-----BEGIN RSA " + "PRIVATE KEY-----abc-----END RSA " + "PRIVATE KEY-----",
	}, " ")
}

func recognizedSourceFreshness(value string) bool {
	switch value {
	case "fresh", "stubbed", "no_data", "source_unavailable", "parse_error", "stale", "missing", "unknown":
		return true
	default:
		return false
	}
}

func recognizedSourceFailureCategory(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return true
	}
	return recognizedSourceFreshness(value)
}

func isCoreSourceQualityCategory(item dto.DataSourceQualityCase) bool {
	category := strings.TrimSpace(item.DataCategory)
	switch category {
	case "index_constituents", "index_weights", "index_valuation_files":
		return true
	}
	return strings.TrimSpace(item.SourceLevel) == "A" || strings.TrimSpace(item.SourceType) == "index_basic"
}

func dataSourceQualitySummary(mode, status string, total, degraded, failed int) string {
	return fmt.Sprintf("数据源质量回归 mode=%s status=%s cases=%d degraded=%d failed=%d", mode, status, total, degraded, failed)
}

func sanitizeDataSourceQualityText(text string) string {
	text = sanitizeDecisionLoopText(text)
	return sourceQualitySQLFromPattern.ReplaceAllString(text, "FROM [REDACTED_SOURCE]")
}

func firstNonEmptySourceQuality(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func sortedStringSet(items map[string]struct{}) []string {
	out := make([]string, 0, len(items))
	for item := range items {
		out = append(out, item)
	}
	sort.Strings(out)
	return out
}
