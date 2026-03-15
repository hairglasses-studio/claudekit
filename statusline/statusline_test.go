package statusline

import (
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	input := `{"model":"claude-opus-4-6","working_directory":"/Users/dev/project","cost_usd":0.42,"total_tokens":15000,"max_tokens":200000,"duration_ms":120000}`
	output := Render(strings.NewReader(input))

	if !strings.Contains(output, "Opus") {
		t.Error("should contain model name Opus")
	}
	if !strings.Contains(output, "0.42") {
		t.Error("should contain cost")
	}
	if !strings.Contains(output, "2m0s") {
		t.Error("should contain duration")
	}
}

func TestShortModel(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"claude-opus-4-6", "Opus"},
		{"claude-sonnet-4-6", "Sonnet"},
		{"claude-haiku-4-5-20251001", "Haiku"},
		{"gpt-4", "4"},
	}
	for _, tt := range tests {
		if got := shortModel(tt.input); got != tt.want {
			t.Errorf("shortModel(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestCwdSlug(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", "~"},
		{"/usr/local", "/usr/local"},
	}
	for _, tt := range tests {
		got := cwdSlug(tt.input)
		if tt.input == "" && got != "~" {
			t.Errorf("cwdSlug(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestContextBar(t *testing.T) {
	tests := []struct {
		tokens, max int
		filled      int
	}{
		{0, 100, 0},
		{50, 100, 5},
		{100, 100, 10},
		{0, 0, 0}, // zero max
	}
	for _, tt := range tests {
		bar := contextBar(tt.tokens, tt.max)
		got := strings.Count(bar, "▓")
		if got != tt.filled {
			t.Errorf("contextBar(%d, %d) filled=%d, want %d: %q", tt.tokens, tt.max, got, tt.filled, bar)
		}
	}
}

func TestInstall(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	result, err := Install(InstallOpts{Style: StyleFull})
	if err != nil {
		t.Fatal(err)
	}

	if result.ScriptPath == "" {
		t.Error("ScriptPath should not be empty")
	}
	if result.SettingsPath == "" {
		t.Error("SettingsPath should not be empty")
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		ms   int
		want string
	}{
		{5000, "5s"},
		{120000, "2m0s"},
		{90000, "1m30s"},
	}
	for _, tt := range tests {
		if got := formatDuration(tt.ms); got != tt.want {
			t.Errorf("formatDuration(%d) = %q, want %q", tt.ms, got, tt.want)
		}
	}
}
