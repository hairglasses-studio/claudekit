package mcpserver

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/hairglasses-studio/mcpkit/handler"
	"github.com/hairglasses-studio/mcpkit/prompts"
	"github.com/hairglasses-studio/mcpkit/registry"
	"github.com/hairglasses-studio/mcpkit/resources"
	"github.com/mark3labs/mcp-go/mcp"
)

type ContractToolModule struct {
	ToolRegistry     *registry.ToolRegistry
	ResourceRegistry *resources.ResourceRegistry
	PromptRegistry   *prompts.PromptRegistry
	Version          string
}

type ContractResourceModule struct {
	ToolRegistry   *registry.ToolRegistry
	PromptRegistry *prompts.PromptRegistry
	Version        string
}

type ContractPromptModule struct{}

type claudekitToolSchemaInput struct {
	Name string `json:"name" jsonschema:"required,description=Exact tool name to inspect"`
}

type claudekitToolSchemaOutput struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Category     string `json:"category"`
	IsWrite      bool   `json:"is_write"`
	InputSchema  any    `json:"input_schema,omitempty"`
	OutputSchema any    `json:"output_schema,omitempty"`
}

func (m *ContractToolModule) Name() string { return "contract" }
func (m *ContractToolModule) Description() string {
	return "Discovery and server health tools for claudekit-mcp"
}

func (m *ContractToolModule) Tools() []registry.ToolDefinition {
	toolSchema := handler.TypedHandler[claudekitToolSchemaInput, claudekitToolSchemaOutput](
		"claudekit_tool_schema",
		"Inspect one claudekit tool descriptor including schemas and write-safety hints.",
		func(_ context.Context, input claudekitToolSchemaInput) (claudekitToolSchemaOutput, error) {
			td, ok := m.ToolRegistry.GetTool(input.Name)
			if !ok {
				return claudekitToolSchemaOutput{}, fmt.Errorf("tool not found: %s", input.Name)
			}
			annotated := registry.ApplyToolMetadata(td, "", false)
			return claudekitToolSchemaOutput{
				Name:         td.Tool.Name,
				Description:  td.Tool.Description,
				Category:     td.Category,
				IsWrite:      td.IsWrite,
				InputSchema:  td.Tool.InputSchema,
				OutputSchema: annotated.Tool.OutputSchema,
			}, nil
		},
	)
	toolSchema.Category = "discovery"
	toolSchema.SearchTerms = []string{"tool schema", "tool descriptor", "input schema", "output schema"}

	toolStats := handler.TypedHandler[struct{}, map[string]any](
		"claudekit_tool_stats",
		"Show claudekit tool counts by category and resource/prompt coverage.",
		func(_ context.Context, _ struct{}) (map[string]any, error) {
			stats := m.ToolRegistry.GetToolStats()
			resourceCount := 0
			promptCount := 0
			if m.ResourceRegistry != nil {
				resourceCount = m.ResourceRegistry.ResourceCount() + m.ResourceRegistry.TemplateCount()
			}
			if m.PromptRegistry != nil {
				promptCount = m.PromptRegistry.PromptCount()
			}
			return map[string]any{
				"tool_count":      stats.TotalTools,
				"module_count":    stats.ModuleCount,
				"resource_count":  resourceCount,
				"prompt_count":    promptCount,
				"by_category":     stats.ByCategory,
				"by_runtime_group": stats.ByRuntimeGroup,
				"write_tools":     stats.WriteToolsCount,
				"read_only_tools": stats.ReadOnlyCount,
			}, nil
		},
	)
	toolStats.Category = "discovery"
	toolStats.SearchTerms = []string{"tool stats", "catalog stats", "tool counts"}

	serverHealth := handler.TypedHandler[struct{}, map[string]any](
		"claudekit_server_health",
		"Show the claudekit MCP contract shape, version, runtime, and discovery coverage.",
		func(_ context.Context, _ struct{}) (map[string]any, error) {
			stats := m.ToolRegistry.GetToolStats()
			resourceCount := 0
			promptCount := 0
			if m.ResourceRegistry != nil {
				resourceCount = m.ResourceRegistry.ResourceCount() + m.ResourceRegistry.TemplateCount()
			}
			if m.PromptRegistry != nil {
				promptCount = m.PromptRegistry.PromptCount()
			}
			return map[string]any{
				"server":         "claudekit",
				"version":        m.Version,
				"status":         "ok",
				"go_version":     runtime.Version(),
				"tool_count":     stats.TotalTools,
				"resource_count": resourceCount,
				"prompt_count":   promptCount,
				"discovery_tools": []string{
					"claudekit_tool_search",
					"claudekit_tool_catalog",
					"claudekit_tool_schema",
					"claudekit_tool_stats",
					"claudekit_server_health",
				},
			}, nil
		},
	)
	serverHealth.Category = "discovery"
	serverHealth.SearchTerms = []string{"server health", "contract", "status", "overview"}

	return []registry.ToolDefinition{toolSchema, toolStats, serverHealth}
}

func (m *ContractResourceModule) Name() string { return "server_context" }
func (m *ContractResourceModule) Description() string {
	return "Reusable overview and workflow guides for claudekit-mcp"
}

func (m *ContractResourceModule) Resources() []resources.ResourceDefinition {
	return []resources.ResourceDefinition{
		{
			Resource: mcp.NewResource(
				"claudekit://server/overview",
				"claudekit Overview",
				mcp.WithResourceDescription("Server card covering discovery-first usage and the highest-value claudekit workflows"),
				mcp.WithMIMEType("text/markdown"),
			),
			Handler: func(_ context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
				toolCount := 0
				promptCount := 0
				if m.ToolRegistry != nil {
					toolCount = m.ToolRegistry.ToolCount()
				}
				if m.PromptRegistry != nil {
					promptCount = m.PromptRegistry.PromptCount()
				}
				return []mcp.ResourceContents{
					mcp.TextResourceContents{
						URI:      "claudekit://server/overview",
						MIMEType: "text/markdown",
						Text: strings.Join([]string{
							"# claudekit",
							"",
							fmt.Sprintf("- Version: `%s`", m.Version),
							fmt.Sprintf("- Registered tools: `%d`", toolCount),
							fmt.Sprintf("- Registered prompt workflows: `%d`", promptCount),
							"",
							"1. Start with `claudekit_tool_search` or `claudekit_tool_catalog`.",
							"2. Use `env_status`, `font_status`, and `theme_list` for read-first environment inspection.",
							"3. Use `workflow_list` and `workflow_run` when the task matches a packaged workflow rather than a one-off tool call.",
							"4. Use the roadmap and rdcycle tools only after discovery identifies the smallest useful surface.",
						}, "\n"),
					},
				}, nil
			},
			Category: "overview",
			Tags:     []string{"claudekit", "overview", "workflow"},
		},
		{
			Resource: mcp.NewResource(
				"claudekit://workflows/start-here",
				"claudekit Start Here",
				mcp.WithResourceDescription("Compact workflow for discovery-first claudekit usage"),
				mcp.WithMIMEType("text/markdown"),
			),
			Handler: func(_ context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
				return []mcp.ResourceContents{
					mcp.TextResourceContents{
						URI:      "claudekit://workflows/start-here",
						MIMEType: "text/markdown",
						Text: strings.Join([]string{
							"1. Use `claudekit_server_health` to confirm the contract surface.",
							"2. Use `claudekit_tool_search` or `claudekit_tool_catalog` before loading a broad workflow.",
							"3. Use `env_status`, `theme_list`, `font_status`, or `workflow_list` for the first read path.",
							"4. Only then move to `workflow_run`, roadmap mutation, or rdcycle write tools if the read path justifies it.",
						}, "\n"),
					},
				}, nil
			},
			Category: "workflow",
			Tags:     []string{"workflow", "discovery", "triage"},
		},
	}
}

func (m *ContractResourceModule) Templates() []resources.TemplateDefinition { return nil }

func (m *ContractPromptModule) Name() string { return "server_prompts" }
func (m *ContractPromptModule) Description() string {
	return "Prompt workflows for discovery-first claudekit triage"
}

func (m *ContractPromptModule) Prompts() []prompts.PromptDefinition {
	return []prompts.PromptDefinition{
		{
			Prompt: mcp.NewPrompt(
				"claudekit_start_triage",
				mcp.WithPromptDescription("Use claudekit in a discovery-first sequence before choosing a deeper workflow"),
			),
			Handler: func(_ context.Context, _ mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
				return mcp.NewGetPromptResult("Start claudekit triage", []mcp.PromptMessage{
					mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(
						"Triage this claudekit task in a discovery-first way. Start with `claudekit_tool_search` or `claudekit_tool_catalog`, then `claudekit_server_health`, then use the smallest read-first tools such as `env_status`, `font_status`, `theme_list`, or `workflow_list` before taking any write action.",
					)),
				}), nil
			},
			Category: "workflow",
			Tags:     []string{"triage", "workflow", "discovery"},
		},
	}
}
