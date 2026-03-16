package mcpserver

import (
	"testing"

	"github.com/hairglasses-studio/mcpkit/ralph"
	"github.com/hairglasses-studio/mcpkit/registry"
)

func TestRalphStatusHandlerIdle(t *testing.T) {
	reg := registry.NewToolRegistry()
	mod := SetupRalph(reg, nil, RalphConfig{})
	tools := mod.Tools()

	td := findTool(tools, "ralph_status")
	if td == nil {
		t.Fatal("ralph_status tool not found")
	}

	result := callTool(t, td, map[string]interface{}{})

	var out ralphStatusOutput
	extractJSON(t, result, &out)

	if out.Status != ralph.StatusIdle {
		t.Errorf("Status = %q, want %q", out.Status, ralph.StatusIdle)
	}
	if out.Iteration != 0 {
		t.Errorf("Iteration = %d, want 0", out.Iteration)
	}
}

func TestRalphStopHandlerNoLoop(t *testing.T) {
	reg := registry.NewToolRegistry()
	mod := SetupRalph(reg, nil, RalphConfig{})
	tools := mod.Tools()

	td := findTool(tools, "ralph_stop")
	if td == nil {
		t.Fatal("ralph_stop tool not found")
	}

	result := callTool(t, td, map[string]interface{}{})

	var out ralphStopOutput
	extractJSON(t, result, &out)

	if out.Status != ralph.StatusIdle {
		t.Errorf("Status = %q, want %q", out.Status, ralph.StatusIdle)
	}
	if out.Message != "no loop is running" {
		t.Errorf("Message = %q, want %q", out.Message, "no loop is running")
	}
}

func TestRalphCurrentSampler(t *testing.T) {
	reg := registry.NewToolRegistry()
	mod := SetupRalph(reg, nil, RalphConfig{})

	// Initially nil.
	if mod.CurrentSampler() != nil {
		t.Error("expected CurrentSampler() to be nil initially")
	}

	// After SetSampler(nil), still nil.
	mod.SetSampler(nil)
	if mod.CurrentSampler() != nil {
		t.Error("expected CurrentSampler() to be nil after SetSampler(nil)")
	}
}

func TestRalphDescriptionNotEmpty(t *testing.T) {
	reg := registry.NewToolRegistry()
	mod := SetupRalph(reg, nil, RalphConfig{})

	if mod.Description() == "" {
		t.Error("Description() should not be empty")
	}
	if mod.Name() != "ralph" {
		t.Errorf("Name() = %q, want %q", mod.Name(), "ralph")
	}
}
