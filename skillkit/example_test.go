package skillkit_test

import (
	"fmt"

	"github.com/hairglasses-studio/claudekit/skillkit"
)

func ExampleBuiltinIndex() {
	index := skillkit.BuiltinIndex()
	for _, entry := range index {
		fmt.Printf("%s: %s\n", entry.Name, entry.Title)
	}
	// Output:
	// font-setup: Font Setup Skill
	// theme-setup: Theme Setup Skill
	// env-setup: Environment Setup Skill
}

func ExampleFindInIndex() {
	entry := skillkit.FindInIndex("font-setup")
	if entry != nil {
		fmt.Printf("Found: %s\n", entry.Title)
		fmt.Printf("Tools: %v\n", entry.Tools)
	}
	// Output:
	// Found: Font Setup Skill
	// Tools: [font_status font_install font_configure]
}

func ExampleFindInIndex_notFound() {
	entry := skillkit.FindInIndex("nonexistent")
	fmt.Printf("Found: %v\n", entry != nil)
	// Output:
	// Found: false
}
