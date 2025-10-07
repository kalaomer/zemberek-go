//go:build demo
// +build demo

package main

import (
	"fmt"

	"github.com/kalaomer/zemberek-go/morphology"
)

func main() {
	morph := morphology.CreateWithDefaults()

	// Check lexicon for gitmek
	gitmekItems := morph.Lexicon.GetItems("gitmek")
	fmt.Printf("Lexicon items for 'gitmek': %d\n", len(gitmekItems))
	for _, item := range gitmekItems {
		fmt.Printf("  - Lemma: %s, Root: %s, POS: %v, Attrs: %v\n", item.Lemma, item.Root, item.PrimaryPos, item.Attributes)
	}
	fmt.Println()

	// Test analysis
	wa := morph.Analyze("gidiyorum")
	fmt.Printf("Analysis for 'gidiyorum': %d results\n", len(wa.AnalysisResults))
	for i, sa := range wa.AnalysisResults {
		fmt.Printf("  %d. %s\n", i+1, sa.FormatString())
	}
}
