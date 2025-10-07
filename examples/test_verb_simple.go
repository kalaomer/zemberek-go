//go:build demo
// +build demo

package main

import (
	"fmt"

	"github.com/kalaomer/zemberek-go/morphology"
)

func main() {
	morph := morphology.CreateWithDefaults()

	// Test gelmek
	gelmek := morph.Lexicon.GetItems("gelmek")[0]
	fmt.Printf("gelmek: Root=%s, POS=%v\n", gelmek.Root, gelmek.PrimaryPos)

	// Test gitmek
	gitmek := morph.Lexicon.GetItems("gitmek")[0]
	fmt.Printf("gitmek: Root=%s, POS=%v\n", gitmek.Root, gitmek.PrimaryPos)

	// Try analyzing manually
	fmt.Println("\n=== Manual analysis test ===")
	wa1 := morph.Analyze("gel")
	fmt.Printf("'gel': %d results\n", len(wa1.AnalysisResults))

	wa2 := morph.Analyze("geldi")
	fmt.Printf("'geldi': %d results\n", len(wa2.AnalysisResults))

	wa3 := morph.Analyze("git")
	fmt.Printf("'git': %d results\n", len(wa3.AnalysisResults))

	wa4 := morph.Analyze("gidiyor")
	fmt.Printf("'gidiyor': %d results\n", len(wa4.AnalysisResults))
}
