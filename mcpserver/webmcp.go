package mcpserver

import (
	"encoding/json"
	"net/http"

	"github.com/hairglasses-studio/claudekit/skillkit"
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
//   - GET /skills  — list available skills from the marketplace
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

	mux.HandleFunc("GET /skills", func(w http.ResponseWriter, r *http.Request) {
		index := skillkit.BuiltinIndex()
		type skillEntry struct {
			Name        string   `json:"name"`
			Title       string   `json:"title"`
			Description string   `json:"description"`
			Tools       []string `json:"tools"`
		}
		entries := make([]skillEntry, 0, len(index))
		for _, e := range index {
			entries = append(entries, skillEntry{
				Name:        e.Name,
				Title:       e.Title,
				Description: e.Description,
				Tools:       e.Tools,
			})
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entries)
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	return mux
}
