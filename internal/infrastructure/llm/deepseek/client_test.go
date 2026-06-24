package deepseek

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"investment-agent/internal/domain/analyst"
	"investment-agent/internal/pkg/apperr"
)

func TestClientRequiresAPIKey(t *testing.T) {
	_, err := NewClient(Config{}, nil).Analyze(context.Background(), analyst.Request{AgentName: "value", Symbol: "510300"})
	if !apperr.IsCode(err, apperr.CodeAnalystUnavailable) {
		t.Fatalf("expected ANALYST_UNAVAILABLE, got %v", err)
	}
	if ErrorCategory(err) != "missing_key" {
		t.Fatalf("category = %q, want missing_key", ErrorCategory(err))
	}
}

func TestClientTreatsWhitespaceAPIKeyAsMissing(t *testing.T) {
	_, err := NewClient(Config{APIKey: "   "}, nil).Analyze(context.Background(), analyst.Request{AgentName: "value", Symbol: "510300"})
	if ErrorCategory(err) != "missing_key" {
		t.Fatalf("category = %q, want missing_key", ErrorCategory(err))
	}
}

func TestClientParsesAnalysisMaterial(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Fatalf("missing authorization header")
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Fatalf("accept header = %q, want application/json", r.Header.Get("Accept"))
		}
		if ua := r.Header.Get("User-Agent"); !strings.Contains(ua, "investment-agent") {
			t.Fatalf("user-agent header = %q, want investment-agent marker", ua)
		}
		var body struct {
			Model    string `json:"model"`
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if body.Model != "gpt-5.4-mini" {
			t.Fatalf("model = %q, want gpt-5.4-mini", body.Model)
		}
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"估值分析材料"}}]}`))
	}))
	defer server.Close()

	resp, err := NewClient(Config{APIKey: "test-key", BaseURL: server.URL, Model: "gpt-5.4-mini"}, server.Client()).Analyze(context.Background(), analyst.Request{AgentName: "value", Symbol: "510300", EvidenceSummary: "证据", RuleBoundary: "最终裁决由规则引擎负责"})
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if resp.Reports["value"] != "估值分析材料" {
		t.Fatalf("unexpected report: %+v", resp.Reports)
	}
	if resp.Metadata == nil || resp.Metadata["model"] != "gpt-5.4-mini" || resp.Metadata["parse_status"] != "parsed" || resp.Metadata["quality_status"] != "passed" || resp.Metadata["prompt_version"] == "" {
		t.Fatalf("missing metadata: %+v", resp.Metadata)
	}
}

func TestClientRetriesQualityFailureWithStricterBoundary(t *testing.T) {
	requests := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Messages []struct {
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if len(body.Messages) != 2 {
			t.Fatalf("unexpected messages: %+v", body.Messages)
		}
		requests = append(requests, body.Messages[1].Content)
		if len(requests) == 1 {
			_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"建议买入该标的。"}}]}`))
			return
		}
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"估值分析材料：仅描述估值、风险和证据缺口，不给出交易指令。"}}]}`))
	}))
	defer server.Close()

	resp, err := NewClient(Config{APIKey: "test-key", BaseURL: server.URL, Model: "gpt-5.4-mini"}, server.Client()).Analyze(context.Background(), analyst.Request{AgentName: "value", Symbol: "510300"})
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("expected safety retry, got %d requests", len(requests))
	}
	if !strings.Contains(requests[1], "安全重试") || !strings.Contains(requests[1], "不得使用“建议买入”") {
		t.Fatalf("retry prompt missing stricter boundary: %s", requests[1])
	}
	if resp.Reports["value"] == "" || resp.Metadata["retry"] != "quality_failed_safety_reprompt" {
		t.Fatalf("expected successful retry metadata, resp=%+v", resp)
	}
}

func TestClientRetriesTransportTimeoutOnce(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			time.Sleep(25 * time.Millisecond)
			return
		}
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"估值分析材料：仅描述风险和证据缺口。"}}]}`))
	}))
	defer server.Close()

	httpClient := server.Client()
	httpClient.Timeout = 5 * time.Millisecond
	resp, err := NewClient(Config{APIKey: "test-key", BaseURL: server.URL, Model: "gpt-5.4-mini"}, httpClient).Analyze(context.Background(), analyst.Request{AgentName: "value", Symbol: "510300"})
	if err != nil {
		t.Fatalf("Analyze after timeout retry: %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Fatalf("calls=%d, want 2", got)
	}
	if resp.Metadata["retry"] != "timeout_retry" || resp.Reports["value"] == "" {
		t.Fatalf("missing timeout retry metadata or report: %+v", resp)
	}
}

func TestBuildPromptIncludesKnowledgeReadinessContext(t *testing.T) {
	prompt := buildPrompt(analyst.Request{
		AgentName:               "value",
		Symbol:                  "510300",
		EvidenceSummary:         "正式证据满足",
		PositionContext:         "持仓上下文",
		RuleBoundary:            "最终裁决由规则引擎负责",
		KnowledgeContextSummary: "principles=master.graham.margin_of_safety; data_readiness=valuation_percentiles=degraded; boundary=背景知识不能满足正式证据",
	})

	for _, want := range []string{"知识与数据准备度：", "master.graham.margin_of_safety", "valuation_percentiles=degraded", "背景知识不能满足正式证据", "最终裁决由规则引擎负责"} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("expected prompt to contain %q, got %s", want, prompt)
		}
	}
	if strings.Contains(prompt, "sk-") || strings.Contains(prompt, "raw HTTP") {
		t.Fatalf("prompt must not include sensitive raw material: %s", prompt)
	}
}

func TestClientClassifiesHTTPEmptyParseAndTimeoutErrors(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		client  func(*httptest.Server) *http.Client
		want    string
	}{
		{name: "http error", handler: func(w http.ResponseWriter, r *http.Request) { http.Error(w, "bad", http.StatusBadGateway) }, want: "http_error"},
		{name: "empty response", handler: func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"choices":[{"message":{"content":""}}]}`))
		}, want: "empty_response"},
		{name: "parse error", handler: func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(`not-json`)) }, want: "parse_error"},
		{name: "timeout", handler: func(w http.ResponseWriter, r *http.Request) { time.Sleep(20 * time.Millisecond) }, client: func(s *httptest.Server) *http.Client { c := s.Client(); c.Timeout = time.Millisecond; return c }, want: "timeout"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(tc.handler)
			defer server.Close()
			httpClient := server.Client()
			if tc.client != nil {
				httpClient = tc.client(server)
			}
			_, err := NewClient(Config{APIKey: "test-key", BaseURL: server.URL, Model: "gpt-5.4-mini"}, httpClient).Analyze(context.Background(), analyst.Request{AgentName: "value", Symbol: "510300"})
			if !apperr.IsCode(err, apperr.CodeAnalystUnavailable) {
				t.Fatalf("expected analyst unavailable, got %v", err)
			}
			if ErrorCategory(err) != tc.want {
				t.Fatalf("category = %q, want %q", ErrorCategory(err), tc.want)
			}
		})
	}
}

func TestClientRejectsProhibitedLLMOutput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"最终裁决：立即全仓买入，保证收益翻倍。"}}]}`))
	}))
	defer server.Close()

	_, err := NewClient(Config{APIKey: "test-key", BaseURL: server.URL, Model: "gpt-5.4-mini"}, server.Client()).Analyze(context.Background(), analyst.Request{AgentName: "value", Symbol: "510300"})

	if !apperr.IsCode(err, apperr.CodeAnalystUnavailable) {
		t.Fatalf("expected analyst unavailable, got %v", err)
	}
	if ErrorCategory(err) != "quality_failed" {
		t.Fatalf("category = %q, want quality_failed", ErrorCategory(err))
	}
	metadata := ErrorMetadata(err)
	if metadata["model"] != "gpt-5.4-mini" || metadata["prompt_version"] == "" || metadata["parse_status"] != "parsed" || metadata["quality_status"] != "failed" || metadata["output_summary"] == "" {
		t.Fatalf("expected sanitized failed metadata, got %+v", metadata)
	}
}

func TestEvaluateQualityAllowsNormalAnalysisAndRejectsUnsafeClaims(t *testing.T) {
	if result := EvaluateQuality("基于现有证据，估值偏高，需要继续观察风险，不构成交易建议。"); result.Status != "passed" {
		t.Fatalf("normal analysis rejected: %+v", result)
	}
	if result := EvaluateQuality("仅供规则引擎使用，不包含最终裁决或交易指令。"); result.Status != "passed" {
		t.Fatalf("safe boundary statement rejected: %+v", result)
	}
	if result := EvaluateQuality("该分析材料不下单、不建仓，只供人工复核。"); result.Status != "passed" {
		t.Fatalf("safe no-trade statement rejected: %+v", result)
	}
	for _, text := range []string{"保证收益", "稳定获利", "明天必涨", "确定会上涨", "必然上涨", "立即买入", "建议买入", "建议卖出", "买入510300", "卖出该标的", "最终裁决：买入"} {
		if result := EvaluateQuality(text); result.Status != "failed" {
			t.Fatalf("unsafe output %q not rejected: %+v", text, result)
		}
	}
}

func TestLiveAnalyzeFromEnv(t *testing.T) {
	if os.Getenv("INVESTMENT_AGENT_LIVE_LLM_DIAG") != "1" {
		t.Skip("set INVESTMENT_AGENT_LIVE_LLM_DIAG=1 for a bounded live provider diagnostic")
	}
	resp, err := NewClient(Config{
		APIKey:         os.Getenv("DEEPSEEK_API_KEY"),
		BaseURL:        os.Getenv("DEEPSEEK_BASE_URL"),
		Model:          os.Getenv("DEEPSEEK_MODEL"),
		TimeoutSeconds: 60,
	}, nil).Analyze(context.Background(), analyst.Request{
		AgentName:       "value",
		Symbol:          "510300",
		EvidenceSummary: "P37 real LLM smoke: 本地最小样本，仅验证模型调用、解析、质量门禁和审计记录。",
		PositionContext: "本 smoke 不读取、不写入账户或持仓，不创建确认单或交易流水。",
		RuleBoundary:    "LLM 只生成分析材料，最终裁决由规则引擎负责；不得输出交易指令、收益承诺或最终裁决。",
	})
	if err != nil {
		t.Fatalf("category=%s metadata=%v err=%v", ErrorCategory(err), ErrorMetadata(err), err)
	}
	t.Logf("metadata=%v report_len=%d", resp.Metadata, len(resp.Reports["value"]))
}
