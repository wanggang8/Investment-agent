package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/pkg/apperr"
)

// SzseCollector 实现深交所公告采集（P26 首批）。
type SzseCollector struct {
	HTTPClient *http.Client
	BaseURL    string
}

type szseResponse struct {
	CompanyCount  int            `json:"companyCount"`
	AnnounceCount int            `json:"announceCount"`
	DisclosureTip string         `json:"disclosureTip"`
	RecordCount   int            `json:"recordCount"`
	Data          []szseDataItem `json:"data"`
}

type szseDataItem struct {
	SecCode    string             `json:"secCode"`
	SecName    string             `json:"secName"`
	AnnounList []szseAnnouncement `json:"announList"`
}

type szseAnnouncement struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	AttachPath      string `json:"attachPath"`
	AttachFormat    string `json:"attachFormat"`
	AttachSize      int    `json:"attachSize"`
	AnnID           string `json:"annId"`
	BigCategoryID   string `json:"bigCategoryId"`
	BigCategoryName string `json:"bigCategoryName"`
	PublishTime     string `json:"publishTime"`
	ImportantRatio  string `json:"importantRatio"`
}

func (s *SzseCollector) FetchPublicEvidence(ctx context.Context, symbol string, start, end time.Time) ([]PublicEvidencePayload, error) {
	if s.BaseURL == "" {
		s.BaseURL = "https://www.szse.cn"
	}
	if s.HTTPClient == nil {
		s.HTTPClient = &http.Client{Timeout: 30 * time.Second}
	}
	start, end = evidenceDateRange(start, end)

	payloads := []PublicEvidencePayload{}
	pages := 1
	for page := 1; page <= pages; page++ {
		apiResp, err := s.fetchPage(ctx, symbol, page)
		if err != nil {
			return nil, err
		}
		pages = totalPages(apiResp.RecordCount, publicEvidencePageSize)
		payloads = append(payloads, s.mapAnnouncements(symbol, start, end, apiResp.Data)...)
	}
	if len(payloads) == 0 {
		return nil, PublicEvidenceError{SourceName: "szse", ErrorCode: "no_data", Count: 0, Err: apperr.New(apperr.CodeNotFound, apperr.CategoryNotFound, "未找到公告")}
	}
	return NormalizePublicEvidenceItems(payloads)
}

func (s *SzseCollector) fetchPage(ctx context.Context, symbol string, page int) (szseResponse, error) {
	params := url.Values{}
	params.Set("random", fmt.Sprintf("%.0f", float64(time.Now().UnixNano())/1e6))
	params.Set("pageSize", strconv.Itoa(publicEvidencePageSize))
	params.Set("pageNum", strconv.Itoa(page))
	params.Set("plateCode", "")
	params.Set("annType", "")
	params.Set("keyword", symbol)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.BaseURL+"/api/disc/announcement/searchQuery?"+params.Encode(), nil)
	if err != nil {
		return szseResponse{}, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "构建请求失败", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; investment-agent/1.0)")
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return szseResponse{}, PublicEvidenceError{SourceName: "szse", ErrorCode: "source_unavailable", Count: 0, Err: apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "请求失败", err)}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return szseResponse{}, PublicEvidenceError{SourceName: "szse", ErrorCode: "source_unavailable", Count: 0, Err: apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)))}
	}
	var apiResp szseResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return szseResponse{}, PublicEvidenceError{SourceName: "szse", ErrorCode: "parse_error", Count: 0, Err: apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "解析响应失败", err)}
	}
	return apiResp, nil
}

func (s *SzseCollector) mapAnnouncements(symbol string, start, end time.Time, data []szseDataItem) []PublicEvidencePayload {
	capturedAt := time.Now().UTC()
	payloads := []PublicEvidencePayload{}
	for _, dataItem := range data {
		for _, ann := range dataItem.AnnounList {
			if !evidenceInRange(ann.PublishTime, start, end) {
				continue
			}
			payloads = append(payloads, PublicEvidencePayload{
				SourceName:     "szse",
				SourceLevel:    model.SourceLevelA,
				SourceType:     "public_disclosure",
				EvidenceRole:   "formal",
				Symbol:         symbol,
				SourceRecordID: ann.AnnID,
				Title:          ann.Title,
				Text:           ann.Title,
				URL:            fmt.Sprintf("%s/disclosure/listed/bulletinDetail/index.html?announcementId=%s", s.BaseURL, ann.ID),
				AttachmentURL:  resolveEvidenceURL(s.BaseURL, ann.AttachPath),
				PublishedAt:    ann.PublishTime,
				CapturedAt:     capturedAt,
				Raw:            map[string]any{"sec_code": dataItem.SecCode, "sec_name": dataItem.SecName, "id": ann.ID, "attach_format": ann.AttachFormat, "attach_size": ann.AttachSize, "big_category_id": ann.BigCategoryID, "big_category_name": ann.BigCategoryName, "important_ratio": ann.ImportantRatio},
			})
		}
	}
	return payloads
}
