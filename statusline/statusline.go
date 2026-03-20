// Package statusline provides a Claude Code statusline script and installer.
package statusline

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	"github.com/hairglasses-studio/claudekit/themekit"
)

// Style controls statusline verbosity.
type Style string

const (
	StyleFull    Style = "full"
	StyleCompact Style = "compact"
	StyleMinimal Style = "minimal"
)

// SessionData is the JSON piped to the statusline via stdin.
type SessionData struct {
	Model    string  `json:"model"`
	CWD      string  `json:"working_directory"`
	Session  string  `json:"session_id"`
	CostUSD  float64 `json:"cost_usd"`
	Tokens   int     `json:"total_tokens"`
	MaxToks  int     `json:"max_tokens"`
	Duration int     `json:"duration_ms"`
}

// FontTier describes which icon set to use based on available fonts.
type FontTier int

const (
	TierNerdFont FontTier = iota // Full Nerd Font icons
	TierUnicode                  // Unicode symbols
	TierASCII                    // Pure ASCII
)

// DetectFontTier checks which icon tier is available.
// In a real terminal, this would test glyph rendering.
// For now, we check if Monaspice is the active font by env hint.
func DetectFontTier() FontTier {
	if os.Getenv("CLAUDEKIT_FONT_TIER") == "nerd" {
		return TierNerdFont
	}
	if os.Getenv("CLAUDEKIT_FONT_TIER") == "unicode" {
		return TierUnicode
	}
	// Default: try nerd font icons (they degrade gracefully in most terminals)
	return TierNerdFont
}

type icons struct {
	Model  string
	Folder string
	Git    string
	Cost   string
	Time   string
	Sep    string
}

func iconsForTier(tier FontTier) icons {
	switch tier {
	case TierNerdFont:
		return icons{Model: "\uf10b", Folder: "\uf07b", Git: "\ue725", Cost: "\uf155", Time: "\uf017", Sep: " │ "}
	case TierUnicode:
		return icons{Model: "⬡", Folder: "●", Git: "◆", Cost: "$", Time: "◷", Sep: " │ "}
	default:
		return icons{Model: "[M]", Folder: ">", Git: "*", Cost: "$", Time: "T", Sep: " | "}
	}
}

// Render generates the statusline output from session JSON on the reader.
func Render(r io.Reader) string {
	var data SessionData
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return fmt.Sprintf("claudekit: statusline: %v", err)
	}

	tier := DetectFontTier()
	ic := iconsForTier(tier)

	// Model display name
	model := shortModel(data.Model)

	// CWD slug
	cwd := cwdSlug(data.CWD)

	// Context progress bar
	progress := contextBar(data.Tokens, data.MaxToks)

	// Cost
	cost := fmt.Sprintf("%.2f", data.CostUSD)

	// Duration
	dur := formatDuration(data.Duration)

	// Apply Catppuccin theme colors if configured
	theme := detectTheme()
	reset := "\033[0m"

	modelColor := ""
	cwdColor := ""
	costColor := ""
	timeColor := ""
	sepColor := ""

	if theme != nil {
		modelColor = theme.Get("mauve").ANSI()
		cwdColor = theme.Get("blue").ANSI()
		costColor = theme.Get("green").ANSI()
		timeColor = theme.Get("subtext0").ANSI()
		sepColor = theme.Get("overlay0").ANSI()
	} else {
		reset = ""
	}

	sep := sepColor + ic.Sep + reset

	// Line 1: Model | CWD | Context | Cost
	line1 := fmt.Sprintf("%s%s %s%s%s%s%s %s%s%s%s%s%s %s%s",
		modelColor, ic.Model, model, reset,
		sep,
		cwdColor, ic.Folder, cwd, reset,
		sep, progress,
		sep, costColor, ic.Cost, cost+reset)

	// Line 2: Duration
	line2 := fmt.Sprintf("%s%s %s%s", timeColor, ic.Time, dur, reset)

	return line1 + "\n" + line2 + "\n"
}

// detectTheme returns a Catppuccin palette if CLAUDEKIT_THEME is set, or nil.
func detectTheme() *themekit.Palette {
	t := os.Getenv("CLAUDEKIT_THEME")
	if t == "" {
		return nil
	}
	var flavor themekit.Flavor
	switch t {
	case "latte":
		flavor = themekit.Latte
	case "frappe":
		flavor = themekit.Frappe
	case "macchiato":
		flavor = themekit.Macchiato
	case "mocha":
		flavor = themekit.Mocha
	default:
		flavor = themekit.Mocha
	}
	p := themekit.Catppuccin(flavor)
	return &p
}

func shortModel(model string) string {
	switch {
	case strings.Contains(model, "opus"):
		return "Opus"
	case strings.Contains(model, "sonnet"):
		return "Sonnet"
	case strings.Contains(model, "haiku"):
		return "Haiku"
	default:
		parts := strings.Split(model, "-")
		if len(parts) > 1 {
			return parts[len(parts)-1]
		}
		return model
	}
}

func cwdSlug(path string) string {
	if path == "" {
		return "~"
	}
	home, _ := os.UserHomeDir()
	if home != "" {
		path = strings.Replace(path, home, "~", 1)
	}
	parts := strings.Split(path, "/")
	if len(parts) <= 3 {
		return path
	}
	// Abbreviate: ~/h-s/claudekit
	var slug []string
	for i, p := range parts {
		if i == len(parts)-1 {
			slug = append(slug, p)
		} else if p == "~" {
			slug = append(slug, p)
		} else if len(p) > 2 {
			// Take first char + first char after hyphen if present
			abbr := string(p[0])
			for j := 1; j < len(p); j++ {
				if p[j] == '-' && j+1 < len(p) {
					abbr += "-" + string(p[j+1])
					break
				}
			}
			slug = append(slug, abbr)
		} else {
			slug = append(slug, p)
		}
	}
	return strings.Join(slug, "/")
}

func contextBar(tokens, maxTokens int) string {
	if maxTokens <= 0 {
		return "░░░░░░░░░░"
	}
	ratio := float64(tokens) / float64(maxTokens)
	filled := int(math.Round(ratio * 10))
	if filled > 10 {
		filled = 10
	}
	return strings.Repeat("▓", filled) + strings.Repeat("░", 10-filled)
}

func formatDuration(ms int) string {
	secs := ms / 1000
	if secs < 60 {
		return fmt.Sprintf("%ds", secs)
	}
	mins := secs / 60
	secs = secs % 60
	return fmt.Sprintf("%dm%ds", mins, secs)
}

// ScriptContent returns the shell script that Claude Code will execute.
func ScriptContent() string {
	return `#!/bin/bash
# claudekit statusline — Claude Code status bar
# Reads session JSON from stdin, outputs styled status text

# Self-detect binary path for Go-rendered output
CLAUDEKIT_BIN="$(command -v claudekit 2>/dev/null)"

if [ -n "$CLAUDEKIT_BIN" ]; then
    # Use Go binary for rich rendering
    echo "$(cat)" | "$CLAUDEKIT_BIN" statusline render
else
    # Inline fallback: minimal status from JSON
    INPUT="$(cat)"
    MODEL=$(echo "$INPUT" | grep -o '"model":"[^"]*"' | head -1 | cut -d'"' -f4)
    COST=$(echo "$INPUT" | grep -o '"cost_usd":[0-9.]*' | head -1 | cut -d: -f2)
    echo "${MODEL:-claude} | \$${COST:-0.00}"
fi
`
}
