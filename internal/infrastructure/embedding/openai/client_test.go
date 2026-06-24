package openai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientEmbedsTextWithOpenAICompatibleRequest(t *testing.T) {
	var gotPath string
	var gotAuth string
	var gotModel string
	var gotInput any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		var body struct {
			Model string `json:"model"`
			Input any    `json:"input"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		gotModel = body.Model
		gotInput = body.Input
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"embedding":[0.25,-0.5,0.75]}]}`))
	}))
	defer server.Close()
	client := NewClient(Config{APIKey: "embed-key", BaseURL: server.URL + "/v1", Model: "text-embedding-3-small"}, server.Client())

	vec, err := client.Embed(context.Background(), "监管处罚 风险")
	if err != nil {
		t.Fatalf("Embed: %v", err)
	}

	if gotPath != "/v1/embeddings" {
		t.Fatalf("expected /v1/embeddings path, got %s", gotPath)
	}
	if gotAuth != "Bearer embed-key" {
		t.Fatalf("expected bearer auth, got %q", gotAuth)
	}
	if gotModel != "text-embedding-3-small" || gotInput != "监管处罚 风险" {
		t.Fatalf("unexpected request model/input: model=%q input=%v", gotModel, gotInput)
	}
	if len(vec) != 3 || vec[0] != 0.25 || vec[1] != -0.5 || vec[2] != 0.75 {
		t.Fatalf("unexpected embedding: %+v", vec)
	}
}

func TestClientRejectsEmptyEmbeddingResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[]}`))
	}))
	defer server.Close()
	client := NewClient(Config{BaseURL: server.URL, Model: "text-embedding-3-small"}, server.Client())

	if _, err := client.Embed(context.Background(), "text"); err == nil {
		t.Fatal("expected empty embedding response to fail")
	}
}
