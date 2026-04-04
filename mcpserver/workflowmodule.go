package mcpserver

import (
	"context"
	"fmt"

	"github.com/hairglasses-studio/mcpkit/handler"
	"github.com/hairglasses-studio/mcpkit/registry"
	"github.com/hairglasses-studio/mcpkit/workflow"
)

// WorkflowModule exposes predefined workflow graphs for multi-step terminal setup.
type WorkflowModule struct {
	graphs map[string]*workflowEntry
}

type workflowEntry struct {
	graph       *workflow.Graph
	description string
}

func (m *WorkflowModule) Name() string { return "workflow" }
func (m *WorkflowModule) Description() string {
	return "Predefined workflow graphs for multi-step operations"
}

// NewWorkflowModule creates a workflow module with predefined graphs.
func NewWorkflowModule() *WorkflowModule {
	m := &WorkflowModule{graphs: make(map[string]*workflowEntry)}
	m.graphs["full-setup"] = &workflowEntry{
		graph:       buildFullSetupGraph(),
		description: "Complete terminal setup: detect → font install + theme apply → statusline → env snapshot",
	}
	return m
}

// buildFullSetupGraph creates the full-setup workflow:
// detect → font_install, theme_apply (parallel) → statusline_install → env_snapshot → END
func buildFullSetupGraph() *workflow.Graph {
	g := workflow.NewGraph()

	// Each node is a placeholder that sets state to track progress.
	// Actual tool calls are made by the MCP client when running the workflow.
	g.AddNode("detect", func(ctx context.Context, state workflow.State) (workflow.State, error) {
		return workflow.Set(state, "detect_done", true), nil
	})
	g.AddNode("font_install", func(ctx context.Context, state workflow.State) (workflow.State, error) {
		return workflow.Set(state, "font_install_done", true), nil
	})
	g.AddNode("theme_apply", func(ctx context.Context, state workflow.State) (workflow.State, error) {
		return workflow.Set(state, "theme_apply_done", true), nil
	})
	g.AddNode("statusline_install", func(ctx context.Context, state workflow.State) (workflow.State, error) {
		return workflow.Set(state, "statusline_install_done", true), nil
	})
	g.AddNode("env_snapshot", func(ctx context.Context, state workflow.State) (workflow.State, error) {
		return workflow.Set(state, "env_snapshot_done", true), nil
	})

	g.SetStart("detect")
	// detect → font_install and theme_apply (sequential edges; parallelism is a client concern)
	g.AddEdge("detect", "font_install")
	g.AddEdge("font_install", "theme_apply")
	g.AddEdge("theme_apply", "statusline_install")
	g.AddEdge("statusline_install", "env_snapshot")
	g.AddEdge("env_snapshot", workflow.EndNode)

	return g
}

// --- workflow_list ---

type WorkflowListInput struct{}

type WorkflowListOutput struct {
	Workflows []WorkflowInfo `json:"workflows"`
	Count     int            `json:"count"`
}

type WorkflowInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// --- workflow_run ---

type WorkflowRunInput struct {
	Name   string            `json:"name" jsonschema:"required,description=Workflow name to execute (e.g. full-setup)"`
	Params map[string]string `json:"params,omitempty" jsonschema:"description=Optional parameters to pass to the workflow"`
}

type WorkflowRunOutput struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Steps    int    `json:"steps"`
	Duration string `json:"duration,omitempty"`
	Error    string `json:"error,omitempty"`
}

func (m *WorkflowModule) Tools() []registry.ToolDefinition {
	return []registry.ToolDefinition{
		handler.TypedHandler[WorkflowListInput, WorkflowListOutput](
			"workflow_list",
			"List available predefined workflows with their descriptions.",
			func(_ context.Context, _ WorkflowListInput) (WorkflowListOutput, error) {
				var out WorkflowListOutput
				for name, entry := range m.graphs {
					out.Workflows = append(out.Workflows, WorkflowInfo{
						Name:        name,
						Description: entry.description,
					})
				}
				out.Count = len(out.Workflows)
				return out, nil
			},
		),
		handler.TypedHandler[WorkflowRunInput, WorkflowRunOutput](
			"workflow_run",
			"Execute a predefined workflow by name. Runs all steps in the graph sequentially.",
			func(ctx context.Context, input WorkflowRunInput) (WorkflowRunOutput, error) {
				entry, ok := m.graphs[input.Name]
				if !ok {
					return WorkflowRunOutput{}, fmt.Errorf("workflow %q not found", input.Name)
				}

				if err := entry.graph.Validate(); err != nil {
					return WorkflowRunOutput{}, fmt.Errorf("workflow %q invalid: %w", input.Name, err)
				}

				engine, err := workflow.NewEngine(entry.graph)
				if err != nil {
					return WorkflowRunOutput{}, fmt.Errorf("creating engine: %w", err)
				}

				initial := workflow.NewState()
				for k, v := range input.Params {
					initial = workflow.Set(initial, k, v)
				}

				result, err := engine.Run(ctx, input.Name, initial)
				if err != nil {
					return WorkflowRunOutput{
						Name:   input.Name,
						Status: "failed",
						Error:  err.Error(),
					}, nil
				}

				return WorkflowRunOutput{
					Name:     input.Name,
					Status:   string(result.Status),
					Steps:    result.Steps,
					Duration: result.Duration.String(),
				}, nil
			},
		),
	}
}
