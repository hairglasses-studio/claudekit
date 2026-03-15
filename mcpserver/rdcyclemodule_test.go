package mcpserver

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hairglasses-studio/mcpkit/registry"
)

func TestNewRDCycleModuleTools(t *testing.T) {
	dir := t.TempDir()
	// rdcycle needs a roadmap.json to exist
	data := `{"title":"test","phases":[{"id":"p1","name":"Phase 1","status":"planned","items":[]}]}`
	if err := os.WriteFile(filepath.Join(dir, "roadmap.json"), []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	mod := NewRDCycleModule(dir)
	if mod.Name() != "rdcycle" {
		t.Errorf("Name() = %q, want %q", mod.Name(), "rdcycle")
	}

	tools := mod.Tools()
	if len(tools) != 9 {
		t.Fatalf("got %d tools, want 9", len(tools))
	}

	want := map[string]bool{
		"rdcycle_scan":      false,
		"rdcycle_plan":      false,
		"rdcycle_verify":    false,
		"rdcycle_artifacts": false,
		"rdcycle_commit":    false,
		"rdcycle_report":    false,
		"rdcycle_schedule":  false,
		"rdcycle_notes":     false,
		"rdcycle_improve":   false,
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

func TestSetupRDCycleRegistersTools(t *testing.T) {
	dir := t.TempDir()
	data := `{"title":"test","phases":[]}`
	if err := os.WriteFile(filepath.Join(dir, "roadmap.json"), []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	reg := registry.NewToolRegistry()
	SetupRDCycle(reg, dir)

	tools := reg.ListTools()
	if len(tools) < 9 {
		t.Errorf("expected at least 9 tools, got %d", len(tools))
	}
}
