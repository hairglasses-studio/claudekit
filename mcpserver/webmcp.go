package mcpserver

import (
	"encoding/json"
	"net/http"

	"github.com/hairglasses-studio/mcpkit/registry"
)

// toolInfo is a JSON-serializable summary of a registered tool.
type toolInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category,omitempty"`
}

// WebMCPHandler returns an http.Handler that serves the claudekit MCP tools
// over HTTP using a basic REST pattern. This is a minimal bridge for
// browser-based clients to discover available tools.
//
// Endpoints:
//   - GET /tools   — list all registered tools as JSON
//   - GET /health  — health check (200 OK)
func WebMCPHandler(reg *registry.ToolRegistry) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tools", func(w http.ResponseWriter, r *http.Request) {
		defs := reg.GetAllToolDefinitions()
		tools := make([]toolInfo, 0, len(defs))
		for _, td := range defs {
			tools = append(tools, toolInfo{
				Name:        td.Tool.Name,
				Description: td.Tool.Description,
				Category:    td.Category,
			})
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tools)
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	return mux
}
