package fontkit_test

import (
	"fmt"

	"github.com/hairglasses-studio/claudekit/fontkit"
)

func ExampleDefaultFont() {
	f := fontkit.DefaultFont
	fmt.Printf("Name: %s\n", f.Name)
	fmt.Printf("Nerd glyphs: %v\n", f.HasNerdGlyphs)
	// Output:
	// Name: Monaspice Ne Nerd Font
	// Nerd glyphs: true
}

func ExampleFallbackChain() {
	chain := fontkit.FallbackChain()
	for _, name := range chain {
		fmt.Println(name)
	}
	// Output:
	// MonaspiceNeNFM-Regular
	// MonaspaceNeon-Regular
	// Menlo
}

func ExampleAllFamilies() {
	families := fontkit.AllFamilies()
	fmt.Printf("Total families: %d\n", len(families))
	fmt.Printf("First: %s\n", families[0].Name)
	// Output:
	// Total families: 10
	// First: Monaspace Argon
}

func ExampleTerminal_SupportsConfig() {
	// Ghostty, iTerm2, and WezTerm support font configuration.
	fmt.Printf("ghostty: %v\n", fontkit.TerminalGhostty.SupportsConfig())
	fmt.Printf("iterm2: %v\n", fontkit.TerminalITerm2.SupportsConfig())
	fmt.Printf("unknown: %v\n", fontkit.TerminalUnknown.SupportsConfig())
	// Output:
	// ghostty: true
	// iterm2: true
	// unknown: false
}
