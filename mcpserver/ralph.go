package mcpserver

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/hairglasses-studio/mcpkit/finops"
	"github.com/hairglasses-studio/mcpkit/handler"
	"github.com/hairglasses-studio/mcpkit/ralph"
	"github.com/hairglasses-studio/mcpkit/rdcycle"
	"github.com/hairglasses-studio/mcpkit/registry"
	"github.com/hairglasses-studio/mcpkit/sampling"
)

// RalphConfig holds budget-aware configuration for the ralph module.
type RalphConfig struct {
	Profile     rdcycle.BudgetProfile
	ModelTiers  rdcycle.ModelTierConfig
	CostTracker *finops.Tracker
	CostPolicy  *finops.CostPolicy
}

// ralphModule is a config-injecting ralph module that applies budget profiles,
// model selection, and cost tracking to each loop run.
// RalphModule is exported so callers can wire a sampler after registration.
type RalphModule = ralphModule

type ralphModule struct {
	mu       sync.Mutex
	loop     *ralph.Loop
	cancel   context.CancelFunc
	registry *registry.ToolRegistry
	sampler  sampling.SamplingClient
	rcfg     RalphConfig
}

func (m *ralphModule) Name() string { return "ralph" }
func (m *ralphModule) Description() string {
	return "Autonomous loop runner with budget-aware configuration"
}

// SetupRalph registers the ralph autonomous loop module with budget-aware config.
// Returns the module so the caller can wire a sampler later via SetSampler.
func SetupRalph(reg *registry.ToolRegistry, sampler sampling.SamplingClient, rcfg RalphConfig) *ralphModule {
	m := &ralphModule{
		registry: reg,
		sampler:  sampler,
		rcfg:     rcfg,
	}
	reg.RegisterModule(m)
	return m
}

// SetSampler sets the sampling client after module registration.
// This is needed because the MCP server (which provides sampling) is created
// after modules are registered.
func (m *ralphModule) SetSampler(s sampling.SamplingClient) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sampler = s
}

// CurrentSampler returns the currently configured sampling client.
func (m *ralphModule) CurrentSampler() sampling.SamplingClient {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sampler
}

type ralphStartInput struct {
	SpecFile      string `json:"spec_file" jsonschema:"required,description=Path to the task specification JSON file"`
	MaxIterations int    `json:"max_iterations,omitempty" jsonschema:"description=Maximum loop iterations (capped by budget profile)"`
}

type ralphStartOutput struct {
	Status   ralph.Status `json:"status"`
	SpecFile string       `json:"spec_file"`
	Message  string       `json:"message"`
	Profile  string       `json:"profile,omitempty"`
}

type ralphStopInput struct{}

type ralphStopOutput struct {
	Status  ralph.Status `json:"status"`
	Message string       `json:"message"`
}

type ralphStatusInput struct{}

type ralphStatusOutput struct {
	Status       ralph.Status `json:"status"`
	Iteration    int          `json:"iteration"`
	CompletedIDs []string     `json:"completed_ids"`
	SpecFile     string       `json:"spec_file"`
}

func (m *ralphModule) Tools() []registry.ToolDefinition {
	return []registry.ToolDefinition{
		handler.TypedHandler[ralphStartInput, ralphStartOutput](
			"ralph_start",
			"Start an autonomous loop that iteratively executes tasks from a spec file. The loop reads the spec, calls tools via LLM decisions, and tracks progress to disk. Iteration count is capped by the active budget profile.",
			func(ctx context.Context, input ralphStartInput) (ralphStartOutput, error) {
				m.mu.Lock()
				defer m.mu.Unlock()

				if m.loop != nil {
					status := m.loop.Status()
					if status.Status == ralph.StatusRunning {
						return ralphStartOutput{
							Status:   ralph.StatusRunning,
							SpecFile: status.SpecFile,
							Message:  "loop is already running",
						}, nil
					}
					// Cancel previous loop's context to avoid goroutine leak.
					if m.cancel != nil {
						m.cancel()
						m.cancel = nil
					}
				}

				maxIter := input.MaxIterations
				if m.rcfg.Profile.MaxIterations > 0 {
					if maxIter <= 0 || maxIter > m.rcfg.Profile.MaxIterations {
						maxIter = m.rcfg.Profile.MaxIterations
					}
				}

				// Derive project root: prefer cwd, fall back to spec file's directory.
				projectRoot, _ := os.Getwd()
				if projectRoot == "" {
					projectRoot = filepath.Dir(input.SpecFile)
				}

				config := ralph.Config{
					SpecFile:      input.SpecFile,
					ProjectRoot:   projectRoot,
					MaxIterations: maxIter,
					MaxTokens:     m.rcfg.Profile.MaxTokensPerReq,
					ToolRegistry:  m.registry,
					Sampler:       m.sampler,
					CostTracker:   m.rcfg.CostTracker,
				}

				if m.rcfg.ModelTiers.Default != "" {
					config.ModelSelector = m.rcfg.ModelTiers.Selector()
				}

				// Cost logging and auto-stop on budget breach.
				costPolicy := m.rcfg.CostPolicy
				config.Hooks.OnCostUpdate = func(iteration int, summary finops.UsageSummary) {
					log.Printf("[ralph] iteration %d: %d input tokens, %d output tokens",
						iteration, summary.TotalInputTokens, summary.TotalOutputTokens)
					if costPolicy != nil && costPolicy.RemainingBudget() <= 0 {
						log.Printf("[ralph] cost budget exceeded ($%.2f used) — stopping loop", costPolicy.TotalCost())
						m.mu.Lock()
						if m.loop != nil {
							m.loop.Stop()
						}
						m.mu.Unlock()
					}
				}

				loop, err := ralph.NewLoop(config)
				if err != nil {
					return ralphStartOutput{}, fmt.Errorf("failed to create loop: %w", err)
				}
				m.loop = loop
				ctx, cancel := context.WithCancel(context.Background())
				m.cancel = cancel

				go loop.Run(ctx)

				return ralphStartOutput{
					Status:   ralph.StatusRunning,
					SpecFile: input.SpecFile,
					Message:  "loop started",
					Profile:  m.rcfg.Profile.Name,
				}, nil
			},
		),
		handler.TypedHandler[ralphStopInput, ralphStopOutput](
			"ralph_stop",
			"Stop the currently running autonomous loop. The loop will finish its current iteration and then stop gracefully.",
			func(ctx context.Context, input ralphStopInput) (ralphStopOutput, error) {
				m.mu.Lock()
				defer m.mu.Unlock()

				if m.loop == nil {
					return ralphStopOutput{
						Status:  ralph.StatusIdle,
						Message: "no loop is running",
					}, nil
				}

				m.loop.Stop()
				if m.cancel != nil {
					m.cancel()
				}
				status := m.loop.Status()
				return ralphStopOutput{
					Status:  status.Status,
					Message: "stop signal sent",
				}, nil
			},
		),
		handler.TypedHandler[ralphStatusInput, ralphStatusOutput](
			"ralph_status",
			"Get the current status of the autonomous loop, including iteration count, completed task IDs, and overall status.",
			func(ctx context.Context, input ralphStatusInput) (ralphStatusOutput, error) {
				m.mu.Lock()
				defer m.mu.Unlock()

				if m.loop == nil {
					return ralphStatusOutput{Status: ralph.StatusIdle}, nil
				}

				p := m.loop.Status()
				return ralphStatusOutput{
					Status:       p.Status,
					Iteration:    p.Iteration,
					CompletedIDs: p.CompletedIDs,
					SpecFile:     p.SpecFile,
				}, nil
			},
		),
	}
}
