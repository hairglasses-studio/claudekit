package mcpserver

import (
	"context"
	"testing"

	"github.com/hairglasses-studio/mcpkit/memory"
	"github.com/hairglasses-studio/mcpkit/registry"
)

func TestSetupMemoryRegistersTools(t *testing.T) {
	reg := registry.NewToolRegistry()
	mod := SetupMemory(reg)

	if mod == nil {
		t.Fatal("SetupMemory returned nil")
	}

	tools := reg.ListTools()
	want := map[string]bool{
		"memory_get":    false,
		"memory_set":    false,
		"memory_list":   false,
		"memory_search": false,
	}
	for _, name := range tools {
		if _, ok := want[name]; ok {
			want[name] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("missing tool %q", name)
		}
	}
}

func TestMemoryModuleToolCount(t *testing.T) {
	mod := &MemoryModule{store: memory.NewInMemoryStore()}
	tools := mod.Tools()
	if len(tools) != 4 {
		t.Errorf("got %d tools, want 4", len(tools))
	}
}

func TestMemorySetGetRoundtrip(t *testing.T) {
	store := memory.NewInMemoryStore()
	ctx := context.Background()

	err := store.Set(ctx, memory.MemoryEntry{
		Key:   "test-key",
		Value: "test-value",
		Tier:  memory.TierSemantic,
		Tags:  []string{"test"},
	})
	if err != nil {
		t.Fatal(err)
	}

	entry, found, err := store.Get(ctx, "test-key")
	if err != nil {
		t.Fatal(err)
	}
	if !found {
		t.Fatal("expected entry to be found")
	}
	if entry.Value != "test-value" {
		t.Errorf("got value %q, want %q", entry.Value, "test-value")
	}
	if entry.Tier != memory.TierSemantic {
		t.Errorf("got tier %q, want %q", entry.Tier, memory.TierSemantic)
	}
}

func TestMemoryGetNonexistent(t *testing.T) {
	store := memory.NewInMemoryStore()
	ctx := context.Background()

	_, found, err := store.Get(ctx, "nonexistent")
	if err != nil {
		t.Fatal(err)
	}
	if found {
		t.Error("expected not found for nonexistent key")
	}
}

func TestMemoryListWithFilter(t *testing.T) {
	store := memory.NewInMemoryStore()
	ctx := context.Background()

	_ = store.Set(ctx, memory.MemoryEntry{Key: "a", Value: "alpha", Tier: memory.TierSemantic, Tags: []string{"greek"}})
	_ = store.Set(ctx, memory.MemoryEntry{Key: "b", Value: "beta", Tier: memory.TierEpisodic, Tags: []string{"greek"}})
	_ = store.Set(ctx, memory.MemoryEntry{Key: "c", Value: "gamma", Tier: memory.TierSemantic})

	entries, err := store.List(ctx, memory.ListOptions{Tier: memory.TierSemantic})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 semantic entries, got %d", len(entries))
	}
}

func TestMemorySearchText(t *testing.T) {
	store := memory.NewInMemoryStore()
	ctx := context.Background()

	_ = store.Set(ctx, memory.MemoryEntry{Key: "greeting", Value: "hello world"})
	_ = store.Set(ctx, memory.MemoryEntry{Key: "farewell", Value: "goodbye world"})

	results, err := store.Search(ctx, "hello", memory.SearchOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if len(results) > 0 && results[0].Key != "greeting" {
		t.Errorf("expected key %q, got %q", "greeting", results[0].Key)
	}
}
