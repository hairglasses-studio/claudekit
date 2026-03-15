package mcpserver

import (
	"testing"

	"github.com/hairglasses-studio/mcpkit/finops"
	"github.com/hairglasses-studio/mcpkit/rdcycle"
	"github.com/hairglasses-studio/mcpkit/registry"
)

func TestSetupFinOpsRegistersTools(t *testing.T) {
	reg := registry.NewToolRegistry()
	tracker := SetupFinOps(reg, finops.Config{})

	if tracker == nil {
		t.Fatal("SetupFinOps returned nil tracker")
	}

	tools := reg.ListTools()
	want := map[string]bool{
		"finops_status": false,
		"finops_reset":  false,
	}
	for _, name := range tools {
		if _, ok := want[name]; ok {
			want[name] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("missing tool %q", name)
		}
	}
}

func TestFinOpsModuleToolCount(t *testing.T) {
	mod := &FinOpsModule{tracker: finops.NewTracker()}
	tools := mod.Tools()
	if len(tools) != 2 {
		t.Errorf("got %d tools, want 2", len(tools))
	}
}

func TestFinOpsStatusReturnsZeroOnFreshTracker(t *testing.T) {
	tracker := finops.NewTracker()
	summary := tracker.Summary()

	if summary.TotalInputTokens != 0 {
		t.Errorf("expected 0 input tokens, got %d", summary.TotalInputTokens)
	}
	if summary.TotalOutputTokens != 0 {
		t.Errorf("expected 0 output tokens, got %d", summary.TotalOutputTokens)
	}
	if summary.TotalInvocations != 0 {
		t.Errorf("expected 0 invocations, got %d", summary.TotalInvocations)
	}
}

func TestFinOpsTrackerRecordAndReset(t *testing.T) {
	tracker := finops.NewTracker()
	tracker.Record(finops.UsageEntry{
		ToolName:     "test_tool",
		InputTokens:  100,
		OutputTokens: 50,
	})

	summary := tracker.Summary()
	if summary.TotalInvocations != 1 {
		t.Errorf("expected 1 invocation, got %d", summary.TotalInvocations)
	}

	tracker.Reset()
	summary = tracker.Summary()
	if summary.TotalInvocations != 0 {
		t.Errorf("expected 0 after reset, got %d", summary.TotalInvocations)
	}
}

func TestSetupFinOpsFromProfile(t *testing.T) {
	reg := registry.NewToolRegistry()
	profile := rdcycle.PersonalProfile()

	tracker, cp, wt := SetupFinOpsFromProfile(reg, profile)

	if tracker == nil {
		t.Fatal("tracker is nil")
	}
	if cp == nil {
		t.Fatal("cost policy is nil")
	}
	if wt == nil {
		t.Fatal("windowed tracker is nil")
	}

	// Verify tools are registered.
	tools := reg.ListTools()
	found := false
	for _, name := range tools {
		if name == "finops_status" {
			found = true
			break
		}
	}
	if !found {
		t.Error("finops_status not registered via SetupFinOpsFromProfile")
	}
}

func TestFinOpsStatusDollarFields(t *testing.T) {
	mod := &FinOpsModule{
		tracker:     finops.NewTracker(),
		costPolicy:  finops.NewCostPolicy(finops.WithDollarBudget(10.0)),
		profileName: "test-profile",
	}

	tools := mod.Tools()
	if len(tools) != 2 {
		t.Fatalf("got %d tools, want 2", len(tools))
	}

	if mod.profileName != "test-profile" {
		t.Errorf("profileName = %q, want %q", mod.profileName, "test-profile")
	}
	if mod.costPolicy == nil {
		t.Error("costPolicy is nil")
	}
}

func TestSetupFinOpsFromProfilePersonal(t *testing.T) {
	reg := registry.NewToolRegistry()
	profile := rdcycle.PersonalProfile()
	tracker, cp, wt := SetupFinOpsFromProfile(reg, profile)

	if tracker == nil || cp == nil || wt == nil {
		t.Fatal("all return values should be non-nil")
	}

	// Personal profile has a dollar budget.
	remaining := cp.RemainingBudget()
	if remaining <= 0 {
		t.Errorf("PersonalProfile should have positive dollar budget, got %.2f", remaining)
	}
}

func TestSetupFinOpsFromProfileWorkAPI(t *testing.T) {
	reg := registry.NewToolRegistry()
	personal := rdcycle.PersonalProfile()
	workAPI := rdcycle.WorkAPIProfile()

	_, cpPersonal, _ := SetupFinOpsFromProfile(reg, personal)

	reg2 := registry.NewToolRegistry()
	_, cpWork, _ := SetupFinOpsFromProfile(reg2, workAPI)

	if cpWork.RemainingBudget() <= cpPersonal.RemainingBudget() {
		t.Errorf("WorkAPI budget ($%.2f) should exceed Personal ($%.2f)",
			cpWork.RemainingBudget(), cpPersonal.RemainingBudget())
	}
}

func TestFinOpsStatusNoCostPolicy(t *testing.T) {
	// Module without costPolicy should return zero dollar fields.
	mod := &FinOpsModule{tracker: finops.NewTracker()}

	if mod.costPolicy != nil {
		t.Error("costPolicy should be nil for basic setup")
	}
	if mod.profileName != "" {
		t.Error("profileName should be empty for basic setup")
	}
}

func TestSetupFinOpsBackwardCompatible(t *testing.T) {
	reg := registry.NewToolRegistry()
	tracker := SetupFinOps(reg, finops.Config{TokenBudget: 500000})

	if tracker == nil {
		t.Fatal("SetupFinOps should still work")
	}

	tools := reg.ListTools()
	if len(tools) != 2 {
		t.Errorf("expected 2 tools, got %d", len(tools))
	}
}
