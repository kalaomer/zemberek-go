//go:build demo
// +build demo

package main

import (
	"fmt"

	"github.com/kalaomer/zemberek-go/morphology"
)

func main() {
	morph := morphology.CreateWithDefaults()

	// Check lexicon
	kitapItems := morph.Lexicon.GetItems("kitap")
	fmt.Printf("Lexicon items for 'kitap': %d\n", len(kitapItems))
	for _, item := range kitapItems {
		fmt.Printf("  - Lemma: %s, Root: %s, POS: %v\n", item.Lemma, item.Root, item.PrimaryPos)
	}
	fmt.Println()

	// Test analysis
	words := []string{"kitap", "kitaplar", "gidiyorum", "geldi"}
	for _, word := range words {
		fmt.Printf("=== Analysis for: %s ===\n", word)
		wa := morph.Analyze(word)
		fmt.Printf("Number of analyses: %d\n", len(wa.AnalysisResults))
		for i, sa := range wa.AnalysisResults {
			fmt.Printf("  %d. %s\n", i+1, sa.FormatString())
		}
		fmt.Println()
	}
}
