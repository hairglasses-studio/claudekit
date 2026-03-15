package mcpserver

import (
	"context"
	"testing"

	"github.com/hairglasses-studio/mcpkit/registry"
	"github.com/hairglasses-studio/mcpkit/workflow"
)

func TestNewWorkflowModuleTools(t *testing.T) {
	mod := NewWorkflowModule()
	if mod.Name() != "workflow" {
		t.Errorf("Name() = %q, want %q", mod.Name(), "workflow")
	}

	tools := mod.Tools()
	if len(tools) != 2 {
		t.Errorf("got %d tools, want 2", len(tools))
	}

	want := map[string]bool{
		"workflow_list": false,
		"workflow_run":  false,
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

func TestWorkflowModuleRegistersWithRegistry(t *testing.T) {
	reg := registry.NewToolRegistry()
	reg.RegisterModule(NewWorkflowModule())

	tools := reg.ListTools()
	if len(tools) < 2 {
		t.Errorf("expected at least 2 tools, got %d", len(tools))
	}
}

func TestFullSetupGraphValidates(t *testing.T) {
	g := buildFullSetupGraph()
	if err := g.Validate(); err != nil {
		t.Errorf("full-setup graph validation failed: %v", err)
	}
}

func TestFullSetupGraphRuns(t *testing.T) {
	g := buildFullSetupGraph()
	engine, err := workflow.NewEngine(g)
	if err != nil {
		t.Fatal(err)
	}

	result, err := engine.Run(context.Background(), "test-run", workflow.NewState())
	if err != nil {
		t.Fatal(err)
	}

	if result.Status != workflow.RunStatusCompleted {
		t.Errorf("status = %q, want %q", result.Status, workflow.RunStatusCompleted)
	}
	if result.Steps < 5 {
		t.Errorf("expected at least 5 steps, got %d", result.Steps)
	}

	// Verify state was set by each node
	if done, ok := workflow.Get[bool](result.FinalState, "env_snapshot_done"); !ok || !done {
		t.Error("expected env_snapshot_done to be true in final state")
	}
}
