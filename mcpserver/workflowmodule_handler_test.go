package mcpserver

import (
	"testing"
)

func TestWorkflowListHandler(t *testing.T) {
	mod := NewWorkflowModule()
	tools := mod.Tools()

	td := findTool(tools, "workflow_list")
	if td == nil {
		t.Fatal("workflow_list tool not found")
	}

	result := callTool(t, td, map[string]interface{}{})

	var out WorkflowListOutput
	extractJSON(t, result, &out)

	if out.Count != 1 {
		t.Errorf("Count = %d, want 1", out.Count)
	}

	foundFullSetup := false
	for _, wf := range out.Workflows {
		if wf.Name == "full-setup" {
			foundFullSetup = true
			if wf.Description == "" {
				t.Error("full-setup description should not be empty")
			}
		}
	}
	if !foundFullSetup {
		t.Error("expected full-setup workflow in list")
	}
}

func TestWorkflowRunHandler(t *testing.T) {
	mod := NewWorkflowModule()
	tools := mod.Tools()

	td := findTool(tools, "workflow_run")
	if td == nil {
		t.Fatal("workflow_run tool not found")
	}

	result := callTool(t, td, map[string]interface{}{
		"name": "full-setup",
	})

	var out WorkflowRunOutput
	extractJSON(t, result, &out)

	if out.Name != "full-setup" {
		t.Errorf("Name = %q, want %q", out.Name, "full-setup")
	}
	if out.Status != "completed" {
		t.Errorf("Status = %q, want %q", out.Status, "completed")
	}
	if out.Steps < 5 {
		t.Errorf("Steps = %d, want >= 5", out.Steps)
	}
}

func TestWorkflowRunHandlerNotFound(t *testing.T) {
	mod := NewWorkflowModule()
	tools := mod.Tools()

	td := findTool(tools, "workflow_run")
	if td == nil {
		t.Fatal("workflow_run tool not found")
	}

	result := callToolExpectErr(t, td, map[string]interface{}{
		"name": "nonexistent-workflow",
	})

	if !result.IsError {
		t.Error("expected IsError=true for nonexistent workflow")
	}
}
