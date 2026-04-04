package mcpserver

import (
	"encoding/json"
	"net/http"
	"strings"

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
// Security model: local-only, read-only. Intended to run on localhost.
// CORS is restricted to localhost origins. No mutation endpoints are exposed.
//
// Endpoints:
//   - GET /tools   — list all registered tools as JSON
//   - GET /skills  — list available skills from the marketplace
//   - GET /health  — health check (200 OK)
func WebMCPHandler(reg *registry.ToolRegistry) http.Handler {
	mux := http.NewServeMux()

	// Restrict CORS to localhost origins.
	cors := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "http://127.0.0.1") {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
			next(w, r)
		}
	}

	mux.HandleFunc("GET /tools", cors(func(w http.ResponseWriter, r *http.Request) {
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
		if err := json.NewEncoder(w).Encode(tools); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))

	mux.HandleFunc("GET /skills", cors(func(w http.ResponseWriter, r *http.Request) {
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
		if err := json.NewEncoder(w).Encode(entries); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))

	mux.HandleFunc("GET /health", cors(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))

	return mux
}
