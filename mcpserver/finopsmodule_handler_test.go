package mcpserver

import (
	"testing"

	"github.com/hairglasses-studio/mcpkit/finops"
)

func TestFinOpsStatusHandler(t *testing.T) {
	tracker := finops.NewTracker()
	mod := &FinOpsModule{tracker: tracker}
	tools := mod.Tools()

	td := findTool(tools, "finops_status")
	if td == nil {
		t.Fatal("finops_status tool not found")
	}

	result := callTool(t, td, map[string]interface{}{})

	var out FinOpsStatusOutput
	extractJSON(t, result, &out)

	if out.TotalInputTokens != 0 {
		t.Errorf("TotalInputTokens = %d, want 0", out.TotalInputTokens)
	}
	if out.TotalOutputTokens != 0 {
		t.Errorf("TotalOutputTokens = %d, want 0", out.TotalOutputTokens)
	}
	if out.TotalInvocations != 0 {
		t.Errorf("TotalInvocations = %d, want 0", out.TotalInvocations)
	}
}

func TestFinOpsStatusHandlerWithCostPolicy(t *testing.T) {
	tracker := finops.NewTracker()
	cp := finops.NewCostPolicy(finops.WithDollarBudget(25.0))
	mod := &FinOpsModule{
		tracker:     tracker,
		costPolicy:  cp,
		profileName: "test-budget",
	}
	tools := mod.Tools()

	td := findTool(tools, "finops_status")
	if td == nil {
		t.Fatal("finops_status tool not found")
	}

	result := callTool(t, td, map[string]interface{}{})

	var out FinOpsStatusOutput
	extractJSON(t, result, &out)

	if out.DollarBudget != 25.0 {
		t.Errorf("DollarBudget = %f, want 25.0", out.DollarBudget)
	}
	if out.DollarCost != 0.0 {
		t.Errorf("DollarCost = %f, want 0.0", out.DollarCost)
	}
	if out.ProfileName != "test-budget" {
		t.Errorf("ProfileName = %q, want %q", out.ProfileName, "test-budget")
	}
}

func TestFinOpsResetHandler(t *testing.T) {
	tracker := finops.NewTracker()

	// Record some usage first.
	tracker.Record(finops.UsageEntry{
		ToolName:     "test_tool",
		InputTokens:  500,
		OutputTokens: 200,
	})
	tracker.Record(finops.UsageEntry{
		ToolName:     "another_tool",
		InputTokens:  300,
		OutputTokens: 100,
	})

	// Verify there is usage before reset.
	summary := tracker.Summary()
	if summary.TotalInvocations != 2 {
		t.Fatalf("pre-reset: expected 2 invocations, got %d", summary.TotalInvocations)
	}

	mod := &FinOpsModule{tracker: tracker}
	tools := mod.Tools()

	resetTD := findTool(tools, "finops_reset")
	if resetTD == nil {
		t.Fatal("finops_reset tool not found")
	}

	result := callTool(t, resetTD, map[string]interface{}{})

	var out FinOpsResetOutput
	extractJSON(t, result, &out)

	if out.Message != "Token usage counters reset" {
		t.Errorf("Message = %q, want %q", out.Message, "Token usage counters reset")
	}

	// Verify counters are now zero.
	summary = tracker.Summary()
	if summary.TotalInvocations != 0 {
		t.Errorf("post-reset: TotalInvocations = %d, want 0", summary.TotalInvocations)
	}
	if summary.TotalInputTokens != 0 {
		t.Errorf("post-reset: TotalInputTokens = %d, want 0", summary.TotalInputTokens)
	}
	if summary.TotalOutputTokens != 0 {
		t.Errorf("post-reset: TotalOutputTokens = %d, want 0", summary.TotalOutputTokens)
	}
}
