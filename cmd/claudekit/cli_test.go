package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hairglasses-studio/mcpkit/ralph"
)

// captureStdout runs fn and returns everything written to os.Stdout.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// --- 1. usage() ---

func TestUsageOutput(t *testing.T) {
	out := captureStdout(t, func() {
		usage()
	})
	if !strings.Contains(out, "claudekit") {
		t.Error("usage output should contain 'claudekit'")
	}
	if !strings.Contains(out, "fonts status") {
		t.Error("usage output should mention 'fonts status'")
	}
	if !strings.Contains(out, "theme apply") {
		t.Error("usage output should mention 'theme apply'")
	}
	if !strings.Contains(out, "ralph tail") {
		t.Error("usage output should mention 'ralph tail'")
	}
	if !strings.Contains(out, "skill list") {
		t.Error("usage output should mention 'skill list'")
	}
}

// --- 2. projectDir() ---

func TestProjectDir(t *testing.T) {
	// projectDir uses os.Getwd, so it should return a non-empty, non-"." path
	// in a normal environment.
	dir := projectDir()
	if dir == "" {
		t.Error("projectDir should not return empty string")
	}
	// In the test environment Getwd should succeed, so we shouldn't get "."
	if dir == "." {
		t.Error("projectDir returned '.' which indicates Getwd failure")
	}
}

// --- 3. newToolRegistry() ---

func TestNewToolRegistry(t *testing.T) {
	reg := newToolRegistry()
	if reg == nil {
		t.Fatal("newToolRegistry returned nil")
	}

	tools := reg.ListTools()
	if len(tools) == 0 {
		t.Error("registry should have tools registered")
	}

	modules := reg.ListModules()
	if len(modules) == 0 {
		t.Error("registry should have modules registered")
	}

	// Verify key modules are present
	moduleSet := make(map[string]bool)
	for _, m := range modules {
		moduleSet[m] = true
	}
	for _, expected := range []string{"fonts", "theme", "statusline", "env", "skills"} {
		if !moduleSet[expected] {
			t.Errorf("expected module %q to be registered", expected)
		}
	}
}

// --- 4. printProgress with empty log ---

func TestPrintProgressEmptyLog(t *testing.T) {
	p := ralph.Progress{
		Status:       ralph.StatusIdle,
		Iteration:    0,
		CompletedIDs: nil,
		Log:          nil,
	}

	out := captureStdout(t, func() {
		printProgress(p)
	})

	if !strings.Contains(out, "idle") {
		t.Errorf("expected 'idle' in output, got: %s", out)
	}
	if !strings.Contains(out, "Iteration: 0") {
		t.Errorf("expected 'Iteration: 0' in output, got: %s", out)
	}
	// Should NOT contain "Spec:" since SpecFile is empty
	if strings.Contains(out, "Spec:") {
		t.Errorf("should not print Spec line when SpecFile is empty, got: %s", out)
	}
	// Should NOT contain "Elapsed:" since StartedAt is zero
	if strings.Contains(out, "Elapsed:") {
		t.Errorf("should not print Elapsed line when StartedAt is zero, got: %s", out)
	}
}

// --- 5. printProgress with entries ---

func TestPrintProgressWithEntries(t *testing.T) {
	started := time.Now().Add(-5 * time.Minute)
	p := ralph.Progress{
		SpecFile:     "my-spec.json",
		Iteration:    3,
		CompletedIDs: []string{"scan", "plan", "verify"},
		Status:       ralph.StatusRunning,
		StartedAt:    started,
		Log: []ralph.IterationLog{
			{Iteration: 1, TaskID: "scan", Result: "scanned 10 files"},
			{Iteration: 2, TaskID: "plan", Result: "created plan"},
			{Iteration: 3, TaskID: "verify", Result: "all checks passed"},
		},
	}

	out := captureStdout(t, func() {
		printProgress(p)
	})

	if !strings.Contains(out, "running") {
		t.Errorf("expected 'running' in output, got: %s", out)
	}
	if !strings.Contains(out, "Iteration: 3") {
		t.Errorf("expected 'Iteration: 3' in output, got: %s", out)
	}
	if !strings.Contains(out, "Spec: my-spec.json") {
		t.Errorf("expected 'Spec: my-spec.json' in output, got: %s", out)
	}
	if !strings.Contains(out, "Elapsed:") {
		t.Errorf("expected 'Elapsed:' line in output, got: %s", out)
	}
	// CompletedIDs should appear
	if !strings.Contains(out, "scan") {
		t.Errorf("expected completed ID 'scan' in output, got: %s", out)
	}
}

// --- 6. truncate edge cases ---

func TestTruncateEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		n     int
		want  string
	}{
		{"empty string", "", 10, ""},
		{"shorter than max", "hi", 10, "hi"},
		{"exact length", "abcde", 5, "abcde"},
		{"one over", "abcdef", 5, "ab..."},
		{"max is 3 and string is longer", "abcd", 3, "..."},
		{"single char within limit", "x", 1, "x"},
		{"long string truncated", "hello world this is a long string", 15, "hello world ..."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncate(tt.input, tt.n)
			if got != tt.want {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.n, got, tt.want)
			}
		})
	}
}

// --- 7. runFonts routing ---

func TestRunFontsUnknownCommand(t *testing.T) {
	err := runFonts(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown fonts command")
	}
	if !strings.Contains(err.Error(), "unknown fonts command") {
		t.Errorf("unexpected error message: %v", err)
	}
	if !strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("error should mention the bad command name: %v", err)
	}
}

func TestRunFontsEmptyCommand(t *testing.T) {
	err := runFonts(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty fonts command")
	}
	if !strings.Contains(err.Error(), "unknown fonts command") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// --- 8. runTheme routing ---

func TestRunThemeUnknownCommand(t *testing.T) {
	err := runTheme("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown theme command")
	}
	if !strings.Contains(err.Error(), "unknown theme command") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRunThemeEmptyCommand(t *testing.T) {
	err := runTheme("")
	if err == nil {
		t.Fatal("expected error for empty theme command")
	}
	if !strings.Contains(err.Error(), "unknown theme command") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// --- 9. runEnv routing ---

func TestRunEnvUnknownCommand(t *testing.T) {
	err := runEnv(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown env command")
	}
	if !strings.Contains(err.Error(), "unknown env command") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRunEnvEmptyCommand(t *testing.T) {
	err := runEnv(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty env command")
	}
	if !strings.Contains(err.Error(), "unknown env command") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// --- 10. runSkill routing ---

func TestRunSkillUnknownCommand(t *testing.T) {
	err := runSkill("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown skill command")
	}
	if !strings.Contains(err.Error(), "unknown skill command") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRunSkillEmptyCommand(t *testing.T) {
	err := runSkill("")
	if err == nil {
		t.Fatal("expected error for empty skill command")
	}
	if !strings.Contains(err.Error(), "unknown skill command") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// --- Additional routing tests for more coverage ---

func TestRunMCPUnknownCommand(t *testing.T) {
	err := runMCP(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown mcp command")
	}
	if !strings.Contains(err.Error(), "unknown mcp command") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRunStatuslineUnknownCommand(t *testing.T) {
	err := runStatusline(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown statusline command")
	}
	if !strings.Contains(err.Error(), "unknown statusline command") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRunPluginUnknownCommand(t *testing.T) {
	err := runPlugin("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown plugin command")
	}
	if !strings.Contains(err.Error(), "unknown plugin command") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRunRalphUnknownCommand(t *testing.T) {
	err := runRalph(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown ralph command")
	}
	if !strings.Contains(err.Error(), "unknown ralph command") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// --- runRalph requires os.Args for subcommands ---

func TestRunRalphTailMissingFile(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"claudekit", "ralph", "tail"}

	err := runRalph(context.Background(), "tail")
	if err == nil {
		t.Fatal("expected error for missing file argument")
	}
	if !strings.Contains(err.Error(), "usage:") {
		t.Errorf("expected usage message, got: %v", err)
	}
}

func TestRunRalphStatusMissingFile(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"claudekit", "ralph", "status"}

	err := runRalph(context.Background(), "status")
	if err == nil {
		t.Fatal("expected error for missing file argument")
	}
	if !strings.Contains(err.Error(), "usage:") {
		t.Errorf("expected usage message, got: %v", err)
	}
}

func TestRunSkillInstallMissingName(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"claudekit", "skill", "install"}

	err := runSkill("install")
	if err == nil {
		t.Fatal("expected error for missing skill name")
	}
	if !strings.Contains(err.Error(), "usage:") {
		t.Errorf("expected usage message, got: %v", err)
	}
}

func TestRunSkillRemoveMissingName(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"claudekit", "skill", "remove"}

	err := runSkill("remove")
	if err == nil {
		t.Fatal("expected error for missing skill name")
	}
	if !strings.Contains(err.Error(), "usage:") {
		t.Errorf("expected usage message, got: %v", err)
	}
}

func TestRunPluginAddMissingPath(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"claudekit", "plugin", "add"}

	err := runPlugin("add")
	if err == nil {
		t.Fatal("expected error for missing plugin path")
	}
	if !strings.Contains(err.Error(), "usage:") {
		t.Errorf("expected usage message, got: %v", err)
	}
}

// --- ralphStatus with JSON output ---

func TestRalphStatusJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.progress.json")

	progress := ralph.Progress{
		SpecFile:     "spec.json",
		Iteration:    2,
		CompletedIDs: []string{"task1"},
		Status:       ralph.StatusRunning,
	}
	data, err := json.Marshal(progress)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"claudekit", "ralph", "status", path, "--json"}

	out := captureStdout(t, func() {
		err = ralphStatus(path)
	})
	if err != nil {
		t.Fatalf("ralphStatus: %v", err)
	}

	// Verify it's valid JSON
	var parsed map[string]interface{}
	if jsonErr := json.Unmarshal([]byte(out), &parsed); jsonErr != nil {
		t.Errorf("expected JSON output, got: %s (error: %v)", out, jsonErr)
	}
	if parsed["spec_file"] != "spec.json" {
		t.Errorf("expected spec_file=spec.json, got %v", parsed["spec_file"])
	}
}

func TestRalphStatusNonJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.progress.json")

	progress := ralph.Progress{
		SpecFile:     "myspec.json",
		Iteration:    7,
		CompletedIDs: []string{"a", "b"},
		Status:       ralph.StatusCompleted,
	}
	data, err := json.Marshal(progress)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"claudekit", "ralph", "status", path}

	out := captureStdout(t, func() {
		err = ralphStatus(path)
	})
	if err != nil {
		t.Fatalf("ralphStatus: %v", err)
	}

	if !strings.Contains(out, "completed") {
		t.Errorf("expected 'completed' in output, got: %s", out)
	}
	if !strings.Contains(out, "Iteration: 7") {
		t.Errorf("expected 'Iteration: 7' in output, got: %s", out)
	}
	if !strings.Contains(out, "Spec: myspec.json") {
		t.Errorf("expected 'Spec: myspec.json' in output, got: %s", out)
	}
}

// --- fontsPreview (pure function, no external deps) ---

func TestFontsPreview(t *testing.T) {
	out := captureStdout(t, func() {
		err := fontsPreview()
		if err != nil {
			t.Errorf("fontsPreview: %v", err)
		}
	})
	// fontsPreview just calls fontkit.Preview which returns a string
	// The output should contain some font-related content
	if out == "" {
		t.Error("fontsPreview should produce output")
	}
}

// --- themePreview (pure function, iterates all flavors) ---

func TestThemePreview(t *testing.T) {
	out := captureStdout(t, func() {
		err := themePreview()
		if err != nil {
			t.Errorf("themePreview: %v", err)
		}
	})
	// Should have all four flavor names
	for _, flavor := range []string{"Latte", "Frappé", "Macchiato", "Mocha"} {
		if !strings.Contains(out, flavor) {
			t.Errorf("themePreview should mention flavor %q", flavor)
		}
	}
}

// --- statuslinePreview (pure, uses sample JSON) ---

func TestStatuslinePreview(t *testing.T) {
	out := captureStdout(t, func() {
		err := statuslinePreview()
		if err != nil {
			t.Errorf("statuslinePreview: %v", err)
		}
	})
	if out == "" {
		t.Error("statuslinePreview should produce output")
	}
}

// --- mcpTools (calls newToolRegistry then prints) ---

func TestMcpTools(t *testing.T) {
	out := captureStdout(t, func() {
		err := mcpTools()
		if err != nil {
			t.Errorf("mcpTools: %v", err)
		}
	})
	if !strings.Contains(out, "registered") {
		t.Errorf("expected 'registered' in output, got: %s", out)
	}
	// Should mention at least some tool names
	if !strings.Contains(out, "font_status") && !strings.Contains(out, "theme_apply") {
		t.Errorf("expected known tool names in output, got: %s", out)
	}
}

// --- hasFlag additional tests ---

func TestHasFlagNoArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"claudekit"}

	if hasFlag("anything") {
		t.Error("hasFlag should return false when no flags present")
	}
}

func TestHasFlagPartialMatch(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"claudekit", "--json-output"}

	// --json should NOT match --json-output
	if hasFlag("json") {
		t.Error("hasFlag should not match partial flag names")
	}
}

// --- parseFlag additional edge cases ---

func TestParseFlagEqualsEmpty(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"claudekit", "--flavor="}

	got := parseFlag("flavor", "mocha")
	if got != "" {
		t.Errorf("parseFlag with --flavor= should return empty, got %q", got)
	}
}

func TestParseFlagValueAtEnd(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	// --interval at the end with no following value
	os.Args = []string{"claudekit", "--interval"}

	got := parseFlag("interval", "2")
	if got != "2" {
		t.Errorf("parseFlag should return fallback when flag is last arg, got %q", got)
	}
}
