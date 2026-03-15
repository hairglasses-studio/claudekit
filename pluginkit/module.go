package pluginkit

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hairglasses-studio/mcpkit/registry"
)

// PluginModule adapts a PluginConfig into an mcpkit ToolModule.
type PluginModule struct {
	config  *PluginConfig
	handler *SubprocessHandler
}

// NewPluginModule creates a ToolModule from a plugin config.
func NewPluginModule(cfg *PluginConfig) *PluginModule {
	timeout := 30 * time.Second
	if cfg.Handler.Timeout != "" {
		if d, err := time.ParseDuration(cfg.Handler.Timeout); err == nil {
			timeout = d
		}
	}

	return &PluginModule{
		config: cfg,
		handler: &SubprocessHandler{
			Command: cfg.Handler.Command,
			Timeout: timeout,
		},
	}
}

func (m *PluginModule) Name() string        { return m.config.Name }
func (m *PluginModule) Description() string { return m.config.Description }

func (m *PluginModule) Tools() []registry.ToolDefinition {
	var tools []registry.ToolDefinition
	for _, td := range m.config.Tools {
		toolName := td.Name

		// Build the input schema from the plugin YAML definition.
		inputSchema := registry.ToolInputSchema{
			Type:       "object",
			Properties: make(map[string]any),
		}
		if td.InputSchema != nil {
			if props, ok := td.InputSchema["properties"].(map[string]interface{}); ok {
				for k, v := range props {
					inputSchema.Properties[k] = v
				}
			}
			if req, ok := td.InputSchema["required"].([]interface{}); ok {
				for _, r := range req {
					if s, ok := r.(string); ok {
						inputSchema.Required = append(inputSchema.Required, s)
					}
				}
			}
		}

		tools = append(tools, registry.ToolDefinition{
			Tool: registry.Tool{
				Name:        td.Name,
				Description: td.Description,
				InputSchema: inputSchema,
			},
			Handler: func(ctx context.Context, request registry.CallToolRequest) (*registry.CallToolResult, error) {
				// Marshal the request arguments to JSON for the subprocess.
				var input json.RawMessage
				if request.Params.Arguments != nil {
					argBytes, err := json.Marshal(request.Params.Arguments)
					if err != nil {
						return registry.MakeErrorResult("failed to marshal arguments: " + err.Error()), nil
					}
					input = argBytes
				} else {
					input = json.RawMessage(`{}`)
				}

				result, err := m.handler.Call(ctx, toolName, input)
				if err != nil {
					return registry.MakeErrorResult(err.Error()), nil
				}

				return registry.MakeTextResult(string(result)), nil
			},
		})
	}
	return tools
}
