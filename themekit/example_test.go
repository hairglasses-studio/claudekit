package themekit_test

import (
	"fmt"

	"github.com/hairglasses-studio/claudekit/themekit"
)

func ExampleCatppuccin() {
	p := themekit.Catppuccin(themekit.Mocha)
	base := p.Get("base")
	fmt.Printf("Theme: %s\n", p.Name)
	fmt.Printf("Base color: #%s\n", base.Hex)
	// Output:
	// Theme: Catppuccin Mocha
	// Base color: #1e1e2e
}

func ExampleColor_ANSI() {
	p := themekit.Catppuccin(themekit.Mocha)
	red := p.Get("red")
	fmt.Printf("Red hex: #%s\n", red.Hex)
	fmt.Printf("Red RGB: (%d, %d, %d)\n", red.R, red.G, red.B)
	// Output:
	// Red hex: #f38ba8
	// Red RGB: (243, 139, 168)
}

func ExampleAllFlavors() {
	for _, f := range themekit.AllFlavors() {
		p := themekit.Catppuccin(f)
		fmt.Println(p.Name)
	}
	// Output:
	// Catppuccin Latte
	// Catppuccin Frappé
	// Catppuccin Macchiato
	// Catppuccin Mocha
}

func ExamplePalette_Get() {
	p := themekit.Catppuccin(themekit.Latte)
	text := p.Get("text")
	fmt.Printf("%s: #%s\n", text.Name, text.Hex)
	// Output:
	// text: #4c4f69
}
