package deepseek

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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
