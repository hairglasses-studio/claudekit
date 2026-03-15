package statusline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// InstallOpts configures statusline installation.
type InstallOpts struct {
	Style Style // Statusline style (default: full)
}

// InstallResult describes what was written.
type InstallResult struct {
	ScriptPath   string
	SettingsPath string
	BackedUp     string // Path to backup if existing config was present
}

// Install writes the statusline script and updates Claude Code settings.
func Install(opts InstallOpts) (*InstallResult, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	claudeDir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(claudeDir, 0o755); err != nil {
		return nil, fmt.Errorf("create ~/.claude: %w", err)
	}

	result := &InstallResult{}

	// Write script
	scriptPath := filepath.Join(claudeDir, "statusline.sh")
	if err := os.WriteFile(scriptPath, []byte(ScriptContent()), 0o755); err != nil {
		return nil, fmt.Errorf("write statusline script: %w", err)
	}
	result.ScriptPath = scriptPath

	// Update settings.json
	settingsPath := filepath.Join(claudeDir, "settings.json")
	result.SettingsPath = settingsPath

	settings := make(map[string]any)

	// Read existing settings
	if data, err := os.ReadFile(settingsPath); err == nil {
		if err := json.Unmarshal(data, &settings); err == nil {
			// Back up if there's an existing statusLine config
			if _, ok := settings["statusLine"]; ok {
				backupPath := settingsPath + ".claudekit-backup"
				if err := os.WriteFile(backupPath, data, 0o644); err != nil {
					return nil, fmt.Errorf("backup settings: %w", err)
				}
				result.BackedUp = backupPath
			}
		}
	}

	// Set statusline config
	settings["statusLine"] = map[string]any{
		"type":    "command",
		"command": scriptPath,
		"padding": 1,
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(settingsPath, data, 0o644); err != nil {
		return nil, fmt.Errorf("write settings: %w", err)
	}

	return result, nil
}
