package mcpserver

import (
	"testing"

	"github.com/hairglasses-studio/mcpkit/registry"
)

func TestSetupRalphRegistersTools(t *testing.T) {
	reg := registry.NewToolRegistry()
	SetupRalph(reg, nil)

	tools := reg.ListTools()
	want := map[string]bool{
		"ralph_start":  false,
		"ralph_stop":   false,
		"ralph_status": false,
	}

	for _, name := range tools {
		if _, ok := want[name]; ok {
			want[name] = true
		}
	}

	for name, found := range want {
		if !found {
			t.Errorf("expected tool %q to be registered", name)
		}
	}
}
