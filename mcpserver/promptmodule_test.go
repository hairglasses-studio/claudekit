package mcpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestPromptModuleName(t *testing.T) {
	mod := &PromptModule{APIKey: "test"}
	if mod.Name() != "prompt" {
		t.Errorf("Name() = %q, want %q", mod.Name(), "prompt")
	}
	if mod.Description() == "" {
		t.Error("Description() is empty")
	}
}

func TestPromptModuleToolCount(t *testing.T) {
	mod := &PromptModule{APIKey: "test"}
	tools := mod.Tools()
	if len(tools) != 1 {
		t.Errorf("got %d tools, want 1", len(tools))
	}
}

func newPromptMockServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(handler)
}

func promptOKHandler(improved string, system string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := promptAPIResponse{
			Prompt: promptAPIResponsePrompt{
				Messages: []promptMsg{
					{Role: "user", Content: improved},
				},
			},
			System: system,
		}
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func TestPromptImproveHandler(t *testing.T) {
	srv := newPromptMockServer(t, promptOKHandler("Improved: hello world", ""))
	defer srv.Close()

	mod := &PromptModule{APIKey: "test-key", BaseURL: srv.URL}
	tools := mod.Tools()
	td := findTool(tools, "prompt_improve")
	if td == nil {
		t.Fatal("prompt_improve tool not found")
	}

	result := callTool(t, td, map[string]interface{}{
		"prompt": "hello world",
	})

	var out PromptImproveOutput
	extractJSON(t, result, &out)

	if out.ImprovedPrompt != "Improved: hello world" {
		t.Errorf("improved_prompt = %q, want %q", out.ImprovedPrompt, "Improved: hello world")
	}
	if out.TargetModel != "claude-sonnet-4-6" {
		t.Errorf("target_model = %q, want %q", out.TargetModel, "claude-sonnet-4-6")
	}
}

func TestPromptImproveWithSystem(t *testing.T) {
	var receivedSystem string
	srv := newPromptMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		var req promptAPIRequest
		json.NewDecoder(r.Body).Decode(&req)
		receivedSystem = req.System

		resp := promptAPIResponse{
			Prompt: promptAPIResponsePrompt{
				Messages: []promptMsg{
					{Role: "user", Content: "improved prompt"},
				},
			},
			System: "improved system",
		}
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	defer srv.Close()

	mod := &PromptModule{APIKey: "test-key", BaseURL: srv.URL}
	tools := mod.Tools()
	td := findTool(tools, "prompt_improve")

	result := callTool(t, td, map[string]interface{}{
		"prompt": "hello",
		"system": "be helpful",
	})

	var out PromptImproveOutput
	extractJSON(t, result, &out)

	if receivedSystem != "be helpful" {
		t.Errorf("system not sent to API: got %q", receivedSystem)
	}
	if out.ImprovedSystem != "improved system" {
		t.Errorf("improved_system = %q, want %q", out.ImprovedSystem, "improved system")
	}
}

func TestPromptImproveWithFeedback(t *testing.T) {
	var receivedFeedback string
	srv := newPromptMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		var req promptAPIRequest
		json.NewDecoder(r.Body).Decode(&req)
		receivedFeedback = req.Feedback

		promptOKHandler("improved", "")(w, r)
	})
	defer srv.Close()

	mod := &PromptModule{APIKey: "test-key", BaseURL: srv.URL}
	tools := mod.Tools()
	td := findTool(tools, "prompt_improve")

	callTool(t, td, map[string]interface{}{
		"prompt":   "hello",
		"feedback": "add more examples",
	})

	if receivedFeedback != "add more examples" {
		t.Errorf("feedback not sent to API: got %q", receivedFeedback)
	}
}

func TestPromptImproveNoAPIKey(t *testing.T) {
	mod := &PromptModule{APIKey: ""}
	tools := mod.Tools()
	td := findTool(tools, "prompt_improve")

	result := callToolExpectErr(t, td, map[string]interface{}{
		"prompt": "hello",
	})

	if !result.IsError {
		t.Error("expected IsError to be true")
	}
}

func TestPromptImproveAPIError(t *testing.T) {
	srv := newPromptMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":{"type":"server_error","message":"internal"}}`))
	})
	defer srv.Close()

	mod := &PromptModule{APIKey: "test-key", BaseURL: srv.URL}
	tools := mod.Tools()
	td := findTool(tools, "prompt_improve")

	result := callToolExpectErr(t, td, map[string]interface{}{
		"prompt": "hello",
	})

	if !result.IsError {
		t.Error("expected IsError to be true")
	}
}

func TestPromptImproveRateLimit(t *testing.T) {
	var attempts atomic.Int32
	srv := newPromptMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		n := attempts.Add(1)
		if n <= 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error":{"type":"rate_limit","message":"slow down"}}`))
			return
		}
		promptOKHandler("retry succeeded", "")(w, r)
	})
	defer srv.Close()

	mod := &PromptModule{APIKey: "test-key", BaseURL: srv.URL}
	tools := mod.Tools()
	td := findTool(tools, "prompt_improve")

	result := callTool(t, td, map[string]interface{}{
		"prompt": "hello",
	})

	var out PromptImproveOutput
	extractJSON(t, result, &out)

	if out.ImprovedPrompt != "retry succeeded" {
		t.Errorf("improved_prompt = %q, want %q", out.ImprovedPrompt, "retry succeeded")
	}

	if got := attempts.Load(); got != 3 {
		t.Errorf("expected 3 attempts, got %d", got)
	}
}
