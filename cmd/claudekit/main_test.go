package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/hairglasses-studio/mcpkit/ralph"
)

func TestRalphStatusReadsProgressFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.progress.json")

	progress := ralph.Progress{
		SpecFile:     "test-spec.json",
		Iteration:    5,
		CompletedIDs: []string{"scan", "plan"},
		Status:       ralph.StatusRunning,
	}
	data, err := json.Marshal(progress)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	// Verify LoadProgress works with our file.
	loaded, err := ralph.LoadProgress(path)
	if err != nil {
		t.Fatalf("LoadProgress: %v", err)
	}
	if loaded.Iteration != 5 {
		t.Errorf("Iteration = %d, want 5", loaded.Iteration)
	}
	if loaded.Status != ralph.StatusRunning {
		t.Errorf("Status = %q, want %q", loaded.Status, ralph.StatusRunning)
	}
	if len(loaded.CompletedIDs) != 2 {
		t.Errorf("CompletedIDs = %v, want [scan plan]", loaded.CompletedIDs)
	}
}

func TestRalphStatusMissingFile(t *testing.T) {
	// LoadProgress returns empty progress for missing files (not an error).
	// ralphStatus should still succeed and print the empty state.
	dir := t.TempDir()
	err := ralphStatus(filepath.Join(dir, "nonexistent.json"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input string
		n     int
		want  string
	}{
		{"hello", 10, "hello"},
		{"hello world", 8, "hello..."},
		{"abc", 3, "abc"},
		{"abcdef", 5, "ab..."},
	}
	for _, tt := range tests {
		got := truncate(tt.input, tt.n)
		if got != tt.want {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.n, got, tt.want)
		}
	}
}

func TestHasFlag(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"claudekit", "ralph", "status", "file.json", "--json"}
	if !hasFlag("json") {
		t.Error("hasFlag should find --json")
	}
	if hasFlag("verbose") {
		t.Error("hasFlag should not find --verbose")
	}
}

func TestParseFlagWithValue(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"claudekit", "--interval", "5"}
	got := parseFlag("interval", "2")
	if got != "5" {
		t.Errorf("parseFlag = %q, want %q", got, "5")
	}
}

func TestParseFlagWithEquals(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"claudekit", "--interval=3"}
	got := parseFlag("interval", "2")
	if got != "3" {
		t.Errorf("parseFlag = %q, want %q", got, "3")
	}
}

func TestParseFlagFallback(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"claudekit"}
	got := parseFlag("interval", "2")
	if got != "2" {
		t.Errorf("parseFlag = %q, want %q", got, "2")
	}
}
