package mcpserver

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hairglasses-studio/mcpkit/registry"
	"github.com/mark3labs/mcp-go/mcp"
)

// findTool locates a ToolDefinition by name in a slice of tools.
func findTool(tools []registry.ToolDefinition, name string) *registry.ToolDefinition {
	for _, td := range tools {
		if td.Tool.Name == name {
			return &td
		}
	}
	return nil
}

// callTool invokes a tool handler with the given arguments and returns the result.
func callTool(t *testing.T, td *registry.ToolDefinition, args map[string]interface{}) *mcp.CallToolResult {
	t.Helper()
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      td.Tool.Name,
			Arguments: args,
		},
	}
	result, err := td.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	return result
}

// callToolExpectErr invokes a tool handler and expects the result to have IsError set.
func callToolExpectErr(t *testing.T, td *registry.ToolDefinition, args map[string]interface{}) *mcp.CallToolResult {
	t.Helper()
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      td.Tool.Name,
			Arguments: args,
		},
	}
	result, err := td.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error (unexpected Go error): %v", err)
	}
	return result
}

// extractJSON extracts the text content from a CallToolResult and unmarshals it into dest.
func extractJSON(t *testing.T, result *mcp.CallToolResult, dest interface{}) {
	t.Helper()
	if len(result.Content) == 0 {
		t.Fatal("result has no content")
	}
	tc, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("first content is not TextContent: %T", result.Content[0])
	}
	if err := json.Unmarshal([]byte(tc.Text), dest); err != nil {
		t.Fatalf("failed to unmarshal result JSON: %v\nraw: %s", err, tc.Text)
	}
}
