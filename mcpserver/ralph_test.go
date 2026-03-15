package mcpserver

import (
	"testing"

	"github.com/hairglasses-studio/mcpkit/rdcycle"
	"github.com/hairglasses-studio/mcpkit/registry"
)

func TestSetupRalphRegistersTools(t *testing.T) {
	reg := registry.NewToolRegistry()
	SetupRalph(reg, nil, RalphConfig{})

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

func TestRalphConfigProfileCapsIterations(t *testing.T) {
	profile := rdcycle.PersonalProfile()
	rcfg := RalphConfig{Profile: profile}

	// Profile MaxIterations should cap the effective value.
	if rcfg.Profile.MaxIterations <= 0 {
		t.Fatal("PersonalProfile should have positive MaxIterations")
	}

	// Verify the profile provides a meaningful cap.
	if rcfg.Profile.MaxIterations > 200 {
		t.Errorf("PersonalProfile MaxIterations=%d, expected <= 200", rcfg.Profile.MaxIterations)
	}
}

func TestRalphConfigModelTiers(t *testing.T) {
	tiers := rdcycle.ModelTierConfig{
		Default: "claude-opus-4-6",
		TaskOverrides: map[string]string{
			"scan":   "claude-sonnet-4-6",
			"verify": "claude-sonnet-4-6",
		},
	}

	selector := tiers.Selector()
	if selector == nil {
		t.Fatal("Selector() returned nil")
	}
}
