package mcpserver

import (
	"testing"

	"github.com/hairglasses-studio/mcpkit/finops"
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

func TestSetupRalphReturnsSameModule(t *testing.T) {
	reg := registry.NewToolRegistry()
	mod := SetupRalph(reg, nil, RalphConfig{})

	if mod == nil {
		t.Fatal("SetupRalph returned nil module")
	}
	if mod.Name() != "ralph" {
		t.Errorf("Name() = %q, want %q", mod.Name(), "ralph")
	}
}

func TestSetSamplerWiresClient(t *testing.T) {
	reg := registry.NewToolRegistry()
	mod := SetupRalph(reg, nil, RalphConfig{})

	if mod.sampler != nil {
		t.Error("sampler should be nil initially")
	}

	// SetSampler should not panic with nil.
	mod.SetSampler(nil)
	if mod.sampler != nil {
		t.Error("sampler should still be nil after SetSampler(nil)")
	}
}

func TestRalphConfigProfileCapsIterations(t *testing.T) {
	profile := rdcycle.PersonalProfile()
	rcfg := RalphConfig{Profile: profile}

	if rcfg.Profile.MaxIterations <= 0 {
		t.Fatal("PersonalProfile should have positive MaxIterations")
	}
	if rcfg.Profile.MaxIterations > 200 {
		t.Errorf("PersonalProfile MaxIterations=%d, expected <= 200", rcfg.Profile.MaxIterations)
	}
}

func TestRalphConfigWorkAPIProfileHigherLimits(t *testing.T) {
	personal := rdcycle.PersonalProfile()
	workAPI := rdcycle.WorkAPIProfile()

	if workAPI.MaxIterations <= personal.MaxIterations {
		t.Errorf("WorkAPI MaxIterations (%d) should exceed Personal (%d)",
			workAPI.MaxIterations, personal.MaxIterations)
	}
	if workAPI.DollarBudget <= personal.DollarBudget {
		t.Errorf("WorkAPI DollarBudget ($%.2f) should exceed Personal ($%.2f)",
			workAPI.DollarBudget, personal.DollarBudget)
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

	// With no completed tasks, Selector infers "scan" is current → sonnet.
	model := selector(0, nil)
	if model != "claude-sonnet-4-6" {
		t.Errorf("model for scan phase = %q, want %q", model, "claude-sonnet-4-6")
	}

	// After scan is complete, next task is plan → uses default (opus).
	model = selector(1, []string{"scan"})
	if model != "claude-opus-4-6" {
		t.Errorf("model for plan phase = %q, want %q", model, "claude-opus-4-6")
	}

	// After scan+plan+implement complete, next is verify → sonnet.
	model = selector(3, []string{"scan", "plan", "implement"})
	if model != "claude-sonnet-4-6" {
		t.Errorf("model for verify phase = %q, want %q", model, "claude-sonnet-4-6")
	}
}

func TestRalphConfigWithCostPolicy(t *testing.T) {
	cp := finops.NewCostPolicy(finops.WithDollarBudget(10.0))
	rcfg := RalphConfig{
		Profile:    rdcycle.PersonalProfile(),
		CostPolicy: cp,
	}

	if rcfg.CostPolicy == nil {
		t.Fatal("CostPolicy should be set")
	}
	remaining := rcfg.CostPolicy.RemainingBudget()
	if remaining != 10.0 {
		t.Errorf("RemainingBudget = %.2f, want 10.00", remaining)
	}
}

func TestRalphModuleToolCount(t *testing.T) {
	reg := registry.NewToolRegistry()
	mod := SetupRalph(reg, nil, RalphConfig{})

	tools := mod.Tools()
	if len(tools) != 3 {
		t.Errorf("got %d tools, want 3", len(tools))
	}
}
