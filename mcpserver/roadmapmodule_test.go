package mcpserver

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hairglasses-studio/mcpkit/registry"
)

func TestNewRoadmapModuleTools(t *testing.T) {
	// Create a temp roadmap.json
	dir := t.TempDir()
	data := `{"title":"test","phases":[{"id":"p1","name":"Phase 1","status":"planned","items":[{"id":"i1","description":"Item 1","status":"planned"}]}]}`
	if err := os.WriteFile(filepath.Join(dir, "roadmap.json"), []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	mod := NewRoadmapModule(dir)
	if mod.Name() != "roadmap" {
		t.Errorf("Name() = %q, want %q", mod.Name(), "roadmap")
	}

	tools := mod.Tools()
	if len(tools) != 4 {
		t.Fatalf("got %d tools, want 4", len(tools))
	}

	want := map[string]bool{
		"roadmap_read":       false,
		"roadmap_update":     false,
		"roadmap_gaps":       false,
		"roadmap_next_phase": false,
	}
	for _, td := range tools {
		if _, ok := want[td.Tool.Name]; ok {
			want[td.Tool.Name] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("missing tool %q", name)
		}
	}
}

func TestSetupRoadmapRegistersTools(t *testing.T) {
	dir := t.TempDir()
	data := `{"title":"test","phases":[]}`
	if err := os.WriteFile(filepath.Join(dir, "roadmap.json"), []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	reg := registry.NewToolRegistry()
	SetupRoadmap(reg, dir)

	tools := reg.ListTools()
	if len(tools) < 4 {
		t.Errorf("expected at least 4 tools, got %d", len(tools))
	}
}
