package mcpserver

import (
	"testing"
)

func TestEnvModuleNameDescription(t *testing.T) {
	mod := &EnvModule{}

	if mod.Name() != "env" {
		t.Errorf("Name() = %q, want %q", mod.Name(), "env")
	}
	if mod.Description() == "" {
		t.Error("Description() should not be empty")
	}
}

func TestEnvModuleToolCount(t *testing.T) {
	mod := &EnvModule{}
	tools := mod.Tools()

	if len(tools) != 2 {
		t.Errorf("got %d tools, want 2", len(tools))
	}

	wantNames := map[string]bool{
		"env_status":   false,
		"env_snapshot": false,
	}
	for _, td := range tools {
		if _, ok := wantNames[td.Tool.Name]; ok {
			wantNames[td.Tool.Name] = true
		}
	}
	for name, found := range wantNames {
		if !found {
			t.Errorf("missing tool %q", name)
		}
	}
}
