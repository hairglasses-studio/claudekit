package mcpserver

import (
	"testing"

	"github.com/hairglasses-studio/mcpkit/memory"
)

func TestMemorySetHandler(t *testing.T) {
	mod := &MemoryModule{store: memory.NewInMemoryStore()}
	tools := mod.Tools()

	td := findTool(tools, "memory_set")
	if td == nil {
		t.Fatal("memory_set tool not found")
	}

	result := callTool(t, td, map[string]interface{}{
		"key":   "greeting",
		"value": "hello world",
		"tier":  "semantic",
		"tags":  []interface{}{"test", "demo"},
	})

	var out MemorySetOutput
	extractJSON(t, result, &out)

	if out.Key != "greeting" {
		t.Errorf("Key = %q, want %q", out.Key, "greeting")
	}
	if out.Message != "Stored successfully" {
		t.Errorf("Message = %q, want %q", out.Message, "Stored successfully")
	}
}

func TestMemoryGetHandler(t *testing.T) {
	mod := &MemoryModule{store: memory.NewInMemoryStore()}
	tools := mod.Tools()

	// First set a value.
	setTD := findTool(tools, "memory_set")
	if setTD == nil {
		t.Fatal("memory_set tool not found")
	}
	callTool(t, setTD, map[string]interface{}{
		"key":   "color",
		"value": "blue",
		"tier":  "episodic",
	})

	// Now get it.
	getTD := findTool(tools, "memory_get")
	if getTD == nil {
		t.Fatal("memory_get tool not found")
	}
	result := callTool(t, getTD, map[string]interface{}{
		"key": "color",
	})

	var out MemoryGetOutput
	extractJSON(t, result, &out)

	if !out.Found {
		t.Error("expected Found=true")
	}
	if out.Key != "color" {
		t.Errorf("Key = %q, want %q", out.Key, "color")
	}
	if out.Value != "blue" {
		t.Errorf("Value = %q, want %q", out.Value, "blue")
	}
	if out.Tier != "episodic" {
		t.Errorf("Tier = %q, want %q", out.Tier, "episodic")
	}
}

func TestMemoryGetHandlerNotFound(t *testing.T) {
	mod := &MemoryModule{store: memory.NewInMemoryStore()}
	tools := mod.Tools()

	td := findTool(tools, "memory_get")
	if td == nil {
		t.Fatal("memory_get tool not found")
	}

	result := callTool(t, td, map[string]interface{}{
		"key": "nonexistent",
	})

	var out MemoryGetOutput
	extractJSON(t, result, &out)

	if out.Found {
		t.Error("expected Found=false for nonexistent key")
	}
	if out.Key != "nonexistent" {
		t.Errorf("Key = %q, want %q", out.Key, "nonexistent")
	}
}

func TestMemoryListHandler(t *testing.T) {
	mod := &MemoryModule{store: memory.NewInMemoryStore()}
	tools := mod.Tools()

	// Set multiple entries.
	setTD := findTool(tools, "memory_set")
	if setTD == nil {
		t.Fatal("memory_set tool not found")
	}

	for _, kv := range []struct{ k, v string }{
		{"alpha", "first"},
		{"beta", "second"},
		{"gamma", "third"},
	} {
		callTool(t, setTD, map[string]interface{}{
			"key":   kv.k,
			"value": kv.v,
		})
	}

	// List all entries.
	listTD := findTool(tools, "memory_list")
	if listTD == nil {
		t.Fatal("memory_list tool not found")
	}

	result := callTool(t, listTD, map[string]interface{}{})

	var out MemoryListOutput
	extractJSON(t, result, &out)

	if out.Count != 3 {
		t.Errorf("Count = %d, want 3", out.Count)
	}
	if len(out.Entries) != 3 {
		t.Errorf("len(Entries) = %d, want 3", len(out.Entries))
	}
}

func TestMemorySearchHandler(t *testing.T) {
	mod := &MemoryModule{store: memory.NewInMemoryStore()}
	tools := mod.Tools()

	// Set entries with distinct values.
	setTD := findTool(tools, "memory_set")
	if setTD == nil {
		t.Fatal("memory_set tool not found")
	}

	callTool(t, setTD, map[string]interface{}{
		"key":   "fruit",
		"value": "apple pie is delicious",
	})
	callTool(t, setTD, map[string]interface{}{
		"key":   "drink",
		"value": "orange juice is refreshing",
	})
	callTool(t, setTD, map[string]interface{}{
		"key":   "dessert",
		"value": "apple crumble is great",
	})

	// Search for "apple".
	searchTD := findTool(tools, "memory_search")
	if searchTD == nil {
		t.Fatal("memory_search tool not found")
	}

	result := callTool(t, searchTD, map[string]interface{}{
		"query": "apple",
	})

	var out MemorySearchOutput
	extractJSON(t, result, &out)

	if out.Count != 2 {
		t.Errorf("Count = %d, want 2", out.Count)
	}
}
