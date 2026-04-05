package envkit_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hairglasses-studio/claudekit/envkit"
)

func ExampleShellInfo_ShellSummary() {
	info := envkit.ShellInfo{
		Shell:   "zsh",
		Manager: "oh-my-zsh",
		RCFile:  "/home/user/.zshrc",
	}
	fmt.Println(info.ShellSummary())
	// Output:
	// Shell: zsh
	// Plugin manager: oh-my-zsh
	// RC file: /home/user/.zshrc
}

func ExampleSnapshotDir() {
	// Create a temporary directory with a managed config file.
	dir, err := os.MkdirTemp("", "envkit-example")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.RemoveAll(dir)

	// Write a test config file.
	configDir := filepath.Join(dir, ".config", "starship.toml")
	os.MkdirAll(filepath.Dir(configDir), 0o755)
	os.WriteFile(configDir, []byte("[character]\nsymbol = \">\""), 0o644)

	snap, err := envkit.SnapshotDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Captured files: %d\n", len(snap.Files))
	// Output:
	// Captured files: 1
}

func ExampleDefaultMiseTools() {
	tools := envkit.DefaultMiseTools()
	// Default tools include Go, Node, and Python.
	fmt.Printf("Go version: %s\n", tools["go"])
	fmt.Printf("Node version: %s\n", tools["node"])
	// Output:
	// Go version: latest
	// Node version: lts
}
