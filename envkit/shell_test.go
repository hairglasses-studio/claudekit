package envkit

import (
	"testing"
)

func TestDetectShell(t *testing.T) {
	info := DetectShell()

	// $SHELL should be set on macOS/Linux
	if info.Shell == "" {
		t.Error("Shell should not be empty")
	}
	if info.Shell == "unknown" {
		t.Log("SHELL env var not set, got 'unknown' — acceptable in CI")
	}

	// Manager should always have a value
	if info.Manager == "" {
		t.Error("Manager should not be empty (expected at least 'none')")
	}
}

func TestShellSummary(t *testing.T) {
	info := ShellInfo{
		Shell:   "zsh",
		Manager: "oh-my-zsh",
		RCFile:  "/home/test/.zshrc",
	}

	summary := info.ShellSummary()
	if summary == "" {
		t.Error("ShellSummary should not be empty")
	}

	// Verify all fields are included
	tests := []string{"zsh", "oh-my-zsh", "/home/test/.zshrc"}
	for _, want := range tests {
		if !contains(summary, want) {
			t.Errorf("summary should contain %q, got: %s", want, summary)
		}
	}
}

func TestShellSummaryNoManager(t *testing.T) {
	info := ShellInfo{
		Shell:   "bash",
		Manager: "none",
		RCFile:  "/home/test/.bashrc",
	}

	summary := info.ShellSummary()
	if contains(summary, "Plugin manager") {
		t.Error("summary should not mention plugin manager when it's 'none'")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
