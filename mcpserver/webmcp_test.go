package mcpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hairglasses-studio/mcpkit/registry"
)

func TestWebMCPToolsEndpoint(t *testing.T) {
	reg := registry.NewToolRegistry()
	reg.RegisterModule(&FontModule{})
	reg.RegisterModule(&ThemeModule{})
	reg.RegisterModule(&StatuslineModule{})

	handler := WebMCPHandler(reg)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/tools")
	if err != nil {
		t.Fatalf("GET /tools: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var tools []toolInfo
	if err := json.NewDecoder(resp.Body).Decode(&tools); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(tools) == 0 {
		t.Fatal("expected at least one tool")
	}

	// Check that known tools are present.
	names := make(map[string]bool)
	for _, tool := range tools {
		names[tool.Name] = true
	}

	for _, expected := range []string{"font_status", "theme_apply"} {
		if !names[expected] {
			t.Errorf("expected tool %q in response, got tools: %v", expected, names)
		}
	}
}

func TestWebMCPHealthEndpoint(t *testing.T) {
	reg := registry.NewToolRegistry()
	handler := WebMCPHandler(reg)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if result["status"] != "ok" {
		t.Errorf("expected status ok, got %q", result["status"])
	}
}
