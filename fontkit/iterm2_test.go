package fontkit

import (
	"encoding/json"
	"os"
	"testing"
)

func TestConfigureITerm2(t *testing.T) {
	// Use a temp dir to avoid polluting real iTerm2 config
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	path, err := ConfigureITerm2(ITerm2Opts{FontSize: 14})
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
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profiles))
	}

	p := profiles[0].(map[string]any)
	if got := p["Normal Font"].(string); got != "MonaspiceNeNFM-Regular 14" {
		t.Errorf("Normal Font = %q, want MonaspiceNeNFM-Regular 14", got)
	}
	if got := p["Non Ascii Font"].(string); got != "MenloRegular 14" {
		t.Errorf("Non Ascii Font = %q, want MenloRegular 14", got)
	}
	if got := p["Use Non-ASCII Font"].(bool); !got {
		t.Error("Use Non-ASCII Font should be true")
	}
}
