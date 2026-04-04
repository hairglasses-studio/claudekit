package mcpserver

import (
	"context"
	"time"

	"github.com/hairglasses-studio/mcpkit/handler"
	"github.com/hairglasses-studio/mcpkit/memory"
	"github.com/hairglasses-studio/mcpkit/registry"
)

// MemoryModule exposes agent memory CRUD tools via MCP.
type MemoryModule struct {
	store memory.Store
}

func (m *MemoryModule) Name() string { return "memory" }
func (m *MemoryModule) Description() string {
	return "Agent memory — persist and retrieve data across tool invocations"
}

// SetupMemory creates an in-memory store and registers the memory module.
func SetupMemory(reg *registry.ToolRegistry) *MemoryModule {
	store := memory.NewInMemoryStore()
	mod := &MemoryModule{store: store}
	reg.RegisterModule(mod)
	return mod
}

// --- memory_get ---

type MemoryGetInput struct {
	Key string `json:"key" jsonschema:"required,description=Memory key to retrieve"`
}

type MemoryGetOutput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Found bool   `json:"found"`
	Tier  string `json:"tier,omitempty"`
}

// --- memory_set ---

type MemorySetInput struct {
	Key   string   `json:"key" jsonschema:"required,description=Memory key to store"`
	Value string   `json:"value" jsonschema:"required,description=Value to store"`
	Tier  string   `json:"tier,omitempty" jsonschema:"description=Memory tier: episodic|semantic|procedural,enum=episodic,enum=semantic,enum=procedural"`
	Tags  []string `json:"tags,omitempty" jsonschema:"description=Tags for categorization and filtering"`
}

type MemorySetOutput struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}

// --- memory_list ---

type MemoryListInput struct {
	Tier   string   `json:"tier,omitempty" jsonschema:"description=Filter by tier,enum=episodic,enum=semantic,enum=procedural"`
	Tags   []string `json:"tags,omitempty" jsonschema:"description=Filter by tags"`
	Prefix string   `json:"prefix,omitempty" jsonschema:"description=Filter by key prefix"`
	Limit  int      `json:"limit,omitempty" jsonschema:"description=Max results (default 100)"`
}

type MemoryListOutput struct {
	Entries []MemoryEntryInfo `json:"entries"`
	Count   int               `json:"count"`
}

type MemoryEntryInfo struct {
	Key       string   `json:"key"`
	Value     string   `json:"value"`
	Tier      string   `json:"tier,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	UpdatedAt string   `json:"updated_at,omitempty"`
}

// --- memory_search ---

type MemorySearchInput struct {
	Query string   `json:"query" jsonschema:"required,description=Text to search for in memory values"`
	Tier  string   `json:"tier,omitempty" jsonschema:"description=Filter by tier,enum=episodic,enum=semantic,enum=procedural"`
	Tags  []string `json:"tags,omitempty" jsonschema:"description=Filter by tags"`
	Limit int      `json:"limit,omitempty" jsonschema:"description=Max results (default 50)"`
}

type MemorySearchOutput struct {
	Results []MemoryEntryInfo `json:"results"`
	Count   int               `json:"count"`
}

func (m *MemoryModule) Tools() []registry.ToolDefinition {
	return []registry.ToolDefinition{
		handler.TypedHandler[MemoryGetInput, MemoryGetOutput](
			"memory_get",
			"Retrieve a value from agent memory by key.",
			func(ctx context.Context, input MemoryGetInput) (MemoryGetOutput, error) {
				entry, found, err := m.store.Get(ctx, input.Key)
				if err != nil {
					return MemoryGetOutput{}, err
				}
				if !found {
					return MemoryGetOutput{Key: input.Key, Found: false}, nil
				}
				return MemoryGetOutput{
					Key:   entry.Key,
					Value: entry.Value,
					Found: true,
					Tier:  string(entry.Tier),
				}, nil
			},
		),
		handler.TypedHandler[MemorySetInput, MemorySetOutput](
			"memory_set",
			"Store a value in agent memory. Supports tiers (episodic/semantic/procedural) and tags for organization.",
			func(ctx context.Context, input MemorySetInput) (MemorySetOutput, error) {
				entry := memory.MemoryEntry{
					Key:       input.Key,
					Value:     input.Value,
					Tier:      memory.Tier(input.Tier),
					Tags:      input.Tags,
					UpdatedAt: time.Now(),
				}
				if err := m.store.Set(ctx, entry); err != nil {
					return MemorySetOutput{}, err
				}
				return MemorySetOutput{
					Key:     input.Key,
					Message: "Stored successfully",
				}, nil
			},
		),
		handler.TypedHandler[MemoryListInput, MemoryListOutput](
			"memory_list",
			"List memory entries with optional filtering by tier, tags, or key prefix.",
			func(ctx context.Context, input MemoryListInput) (MemoryListOutput, error) {
				limit := input.Limit
				if limit <= 0 {
					limit = 100
				}
				entries, err := m.store.List(ctx, memory.ListOptions{
					Tier:   memory.Tier(input.Tier),
					Tags:   input.Tags,
					Prefix: input.Prefix,
					Limit:  limit,
				})
				if err != nil {
					return MemoryListOutput{}, err
				}

				out := MemoryListOutput{Count: len(entries)}
				for _, e := range entries {
					out.Entries = append(out.Entries, memoryEntryToInfo(e))
				}
				return out, nil
			},
		),
		handler.TypedHandler[MemorySearchInput, MemorySearchOutput](
			"memory_search",
			"Search memory entries by text query. Matches against values with optional tier and tag filters.",
			func(ctx context.Context, input MemorySearchInput) (MemorySearchOutput, error) {
				limit := input.Limit
				if limit <= 0 {
					limit = 50
				}
				entries, err := m.store.Search(ctx, input.Query, memory.SearchOptions{
					Tier:  memory.Tier(input.Tier),
					Tags:  input.Tags,
					Limit: limit,
				})
				if err != nil {
					return MemorySearchOutput{}, err
				}

				out := MemorySearchOutput{Count: len(entries)}
				for _, e := range entries {
					out.Results = append(out.Results, memoryEntryToInfo(e))
				}
				return out, nil
			},
		),
	}
}

func memoryEntryToInfo(e memory.MemoryEntry) MemoryEntryInfo {
	info := MemoryEntryInfo{
		Key:   e.Key,
		Value: e.Value,
		Tier:  string(e.Tier),
		Tags:  e.Tags,
	}
	if !e.UpdatedAt.IsZero() {
		info.UpdatedAt = e.UpdatedAt.Format(time.RFC3339)
	}
	return info
}
