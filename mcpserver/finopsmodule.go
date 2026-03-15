package mcpserver

import (
	"context"

	"github.com/hairglasses-studio/mcpkit/finops"
	"github.com/hairglasses-studio/mcpkit/handler"
	"github.com/hairglasses-studio/mcpkit/rdcycle"
	"github.com/hairglasses-studio/mcpkit/registry"
)

// FinOpsModule exposes token usage tracking and budget management tools.
type FinOpsModule struct {
	tracker     *finops.Tracker
	costPolicy  *finops.CostPolicy
	profileName string
}

func (m *FinOpsModule) Name() string        { return "finops" }
func (m *FinOpsModule) Description() string { return "Token usage tracking and budget enforcement" }

// SetupFinOps installs finops middleware on the registry and registers status/reset tools.
// Must be called before other module registrations so middleware wraps all handlers.
func SetupFinOps(reg *registry.ToolRegistry, cfg finops.Config) *finops.Tracker {
	tracker := finops.NewTracker(cfg)
	reg.SetMiddleware([]registry.Middleware{finops.Middleware(tracker)})
	reg.RegisterModule(&FinOpsModule{tracker: tracker})
	return tracker
}

// SetupFinOpsFromProfile installs profile-aware finops middleware with dollar-cost tracking.
// Returns the tracker, cost policy, and windowed tracker for use by other modules (e.g. ralph).
func SetupFinOpsFromProfile(reg *registry.ToolRegistry, profile rdcycle.BudgetProfile) (*finops.Tracker, *finops.CostPolicy, *finops.WindowedTracker) {
	tracker, cp, wt := rdcycle.BuildFinOpsStack(profile)
	reg.SetMiddleware([]registry.Middleware{finops.Middleware(tracker)})
	reg.RegisterModule(&FinOpsModule{tracker: tracker, costPolicy: cp, profileName: profile.Name})
	return tracker, cp, wt
}

// --- finops_status ---

type FinOpsStatusInput struct{}

type FinOpsStatusOutput struct {
	TotalInputTokens  int64            `json:"total_input_tokens"`
	TotalOutputTokens int64            `json:"total_output_tokens"`
	TotalInvocations  int64            `json:"total_invocations"`
	ByTool            map[string]int64 `json:"by_tool,omitempty"`
	DollarCost        float64          `json:"dollar_cost,omitempty"`
	DollarBudget      float64          `json:"dollar_budget,omitempty"`
	ProfileName       string           `json:"profile_name,omitempty"`
}

// --- finops_reset ---

type FinOpsResetInput struct{}

type FinOpsResetOutput struct {
	Message string `json:"message"`
}

func (m *FinOpsModule) Tools() []registry.ToolDefinition {
	return []registry.ToolDefinition{
		handler.TypedHandler[FinOpsStatusInput, FinOpsStatusOutput](
			"finops_status",
			"Get current token usage summary across all tool invocations. Shows total input/output tokens, invocation count, and per-tool breakdown.",
			func(_ context.Context, _ FinOpsStatusInput) (FinOpsStatusOutput, error) {
				summary := m.tracker.Summary()

				out := FinOpsStatusOutput{
					TotalInputTokens:  summary.TotalInputTokens,
					TotalOutputTokens: summary.TotalOutputTokens,
					TotalInvocations:  summary.TotalInvocations,
					ByTool:            summary.ByTool,
					ProfileName:       m.profileName,
				}
				if m.costPolicy != nil {
					out.DollarCost = m.costPolicy.TotalCost()
					out.DollarBudget = m.costPolicy.RemainingBudget() + m.costPolicy.TotalCost()
				}
				return out, nil
			},
		),
		handler.TypedHandler[FinOpsResetInput, FinOpsResetOutput](
			"finops_reset",
			"Reset all token usage counters to zero. Use to start a fresh tracking session.",
			func(_ context.Context, _ FinOpsResetInput) (FinOpsResetOutput, error) {
				m.tracker.Reset()
				return FinOpsResetOutput{Message: "Token usage counters reset"}, nil
			},
		),
	}
}
