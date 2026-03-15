package mcpserver

import (
	"context"
	"strings"
	"testing"

	"github.com/hairglasses-studio/mcpkit/registry"
)

func TestGenerateMetadata(t *testing.T) {
	reg := registry.NewToolRegistry()
	reg.RegisterModule(&FontModule{})
	reg.RegisterModule(&ThemeModule{})
	reg.RegisterModule(&StatuslineModule{})

	meta := GenerateMetadata(reg)

	if meta.Name != "claudekit" {
		t.Errorf("expected name %q, got %q", "claudekit", meta.Name)
	}
	if !strings.Contains(meta.Description, "terminal customization") {
		t.Errorf("description should mention terminal customization, got %q", meta.Description)
	}

	toolCount := reg.ToolCount()
	if len(meta.Tools) != toolCount {
		t.Errorf("expected %d tools, got %d", toolCount, len(meta.Tools))
	}
	if toolCount == 0 {
		t.Error("expected at least one tool in the registry")
	}

	if len(meta.Transports) != 1 || meta.Transports[0].Type != "stdio" {
		t.Errorf("expected single stdio transport, got %v", meta.Transports)
	}
}

func TestPublishMissingToken(t *testing.T) {
	t.Setenv("CLAUDEKIT_REGISTRY_TOKEN", "")

	reg := registry.NewToolRegistry()
	err := Publish(context.Background(), reg)
	if err == nil {
		t.Fatal("expected error when CLAUDEKIT_REGISTRY_TOKEN is not set")
	}
	if !strings.Contains(err.Error(), "CLAUDEKIT_REGISTRY_TOKEN") {
		t.Errorf("error should mention CLAUDEKIT_REGISTRY_TOKEN, got: %v", err)
	}
}
