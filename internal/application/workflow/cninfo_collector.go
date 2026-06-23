package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/pkg/apperr"
)

// CninfoCollector 实现巨潮资讯公告采集（P26 首批）。
type CninfoCollector struct {
	HTTPClient    *http.Client
	BaseURL       string
	OrgIDBySymbol map[string]string
}

type cninfoResponse struct {
	TotalAnnouncement int                  `json:"totalAnnouncement"`
	TotalRecordNum    int                  `json:"totalRecordNum"`
	Announcements     []cninfoAnnouncement `json:"announcements"`
	HasMore           bool                 `json:"hasMore"`
	TotalPages        int                  `json:"totalpages"`
}

type cninfoAnnouncement struct {
	SecCode           string `json:"secCode"`
	SecName           string `json:"secName"`
	OrgID             string `json:"orgId"`
	AnnouncementID    string `json:"announcementId"`
	AnnouncementTitle string `json:"announcementTitle"`
	AnnouncementTime  int64  `json:"announcementTime"`
	AdjunctURL        string `json:"adjunctUrl"`
	AdjunctSize       int    `json:"adjunctSize"`
	AdjunctType       string `json:"adjunctType"`
	AnnouncementType  string `json:"announcementType"`
}

func (c *CninfoCollector) FetchPublicEvidence(ctx context.Context, symbol string, start, end time.Time) ([]PublicEvidencePayload, error) {
	if c.BaseURL == "" {
		c.BaseURL = "https://www.cninfo.com.cn"
	}
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{Timeout: 30 * time.Second}
	}
	start, end = evidenceDateRange(start, end)
	cleanSymbol := cninfoCleanSymbol(symbol)

	payloads := []PublicEvidencePayload{}
	for page := 1; ; page++ {
		apiResp, err := c.fetchPage(ctx, symbol, start, end, page)
		if err != nil {
			return nil, err
		}
		payloads = append(payloads, c.mapAnnouncements(cleanSymbol, start, end, apiResp.Announcements)...)
		if page >= apiResp.TotalPages || !apiResp.HasMore {
			break
		}
	}
	if len(payloads) == 0 {
		return nil, PublicEvidenceError{SourceName: "cninfo", ErrorCode: "no_data", Count: 0, Err: apperr.New(apperr.CodeNotFound, apperr.CategoryNotFound, "未找到公告")}
	}
	return NormalizePublicEvidenceItems(payloads)
}

func (c *CninfoCollector) fetchPage(ctx context.Context, symbol string, start, end time.Time, page int) (cninfoResponse, error) {
	params := url.Values{}
	params.Set("pageNum", strconv.Itoa(page))
	params.Set("pageSize", strconv.Itoa(publicEvidencePageSize))
	params.Set("column", "szse_main")
	params.Set("tabName", "fulltext")
	params.Set("stock", c.cninfoStockParam(symbol))
	params.Set("searchkey", "")
	params.Set("category", "")
	params.Set("seDate", fmt.Sprintf("%s~%s", start.Format("2006-01-02"), end.Format("2006-01-02")))
	params.Set("sortName", "announcementTime")
	params.Set("sortType", "desc")
	params.Set("isHLtitle", "true")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/new/hisAnnouncement/query", strings.NewReader(params.Encode()))
	if err != nil {
		return cninfoResponse{}, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "构建请求失败", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; investment-agent/1.0)")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return cninfoResponse{}, PublicEvidenceError{SourceName: "cninfo", ErrorCode: "source_unavailable", Count: 0, Err: apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "请求失败", err)}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return cninfoResponse{}, PublicEvidenceError{SourceName: "cninfo", ErrorCode: "source_unavailable", Count: 0, Err: apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)))}
	}
	var apiResp cninfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return cninfoResponse{}, PublicEvidenceError{SourceName: "cninfo", ErrorCode: "parse_error", Count: 0, Err: apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "解析响应失败", err)}
	}
	if apiResp.TotalPages <= 0 {
		apiResp.TotalPages = 1
	}
	return apiResp, nil
}

func (c *CninfoCollector) cninfoStockParam(symbol string) string {
	symbol = strings.TrimSpace(symbol)
	if strings.Contains(symbol, ",") {
		return symbol
	}
	if orgID := strings.TrimSpace(c.OrgIDBySymbol[symbol]); orgID != "" {
		return symbol + "," + orgID
	}
	switch symbol {
	case "510300":
		return "510300,9900000091"
	case "000001":
		return "000001,gssz0000001"
	}
	return symbol
}

func cninfoCleanSymbol(symbol string) string {
	symbol = strings.TrimSpace(symbol)
	if before, _, ok := strings.Cut(symbol, ","); ok {
		return strings.TrimSpace(before)
	}
	return symbol
}

func (c *CninfoCollector) mapAnnouncements(symbol string, start, end time.Time, announcements []cninfoAnnouncement) []PublicEvidencePayload {
	capturedAt := time.Now().UTC()
	payloads := make([]PublicEvidencePayload, 0, len(announcements))
	for _, ann := range announcements {
		publishedAt := time.Unix(ann.AnnouncementTime/1000, 0).Format(time.RFC3339)
		if !evidenceInRange(publishedAt, start, end) {
			continue
		}
		payloads = append(payloads, PublicEvidencePayload{
			SourceName:     "cninfo",
			SourceLevel:    model.SourceLevelA,
			SourceType:     "public_disclosure",
			EvidenceRole:   "formal",
			Symbol:         symbol,
			SourceRecordID: ann.AnnouncementID,
			Title:          ann.AnnouncementTitle,
			Text:           ann.AnnouncementTitle,
			URL:            fmt.Sprintf("%s/new/disclosure/detail?plate=&orgId=%s&stockCode=%s&announcementId=%s", c.BaseURL, ann.OrgID, ann.SecCode, ann.AnnouncementID),
			AttachmentURL:  resolveEvidenceURL(c.BaseURL, ann.AdjunctURL),
			PublishedAt:    publishedAt,
			CapturedAt:     capturedAt,
			Raw:            map[string]any{"sec_code": ann.SecCode, "sec_name": ann.SecName, "org_id": ann.OrgID, "adjunct_size": ann.AdjunctSize, "adjunct_type": ann.AdjunctType, "announcement_type": ann.AnnouncementType},
		})
	}
	return payloads
}
