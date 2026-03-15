package themekit

import (
	"encoding/json"
	"os"
	"testing"
)

func TestExportITerm2(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	path, err := ExportITerm2(Catppuccin(Mocha))
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var profile map[string]any
	if err := json.Unmarshal(data, &profile); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	profiles := profile["Profiles"].([]any)
	p := profiles[0].(map[string]any)

	if got := p["Name"].(string); got != "Claudekit Catppuccin Mocha" {
		t.Errorf("Name = %q", got)
	}

	// Verify background color components
	bg := p["Background Color"].(map[string]any)
	if r := bg["Red Component"].(float64); r < 0.11 || r > 0.13 {
		t.Errorf("Background Red = %f, expected ~0.118 (30/255)", r)
	}
}
