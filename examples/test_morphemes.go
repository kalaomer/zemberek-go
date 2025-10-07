//go:build demo
// +build demo

package main

import (
	"fmt"

	"github.com/kalaomer/zemberek-go/morphology"
)

func main() {
	morph := morphology.CreateWithDefaults()

	words := []string{"kitap", "geldi"}

	for _, word := range words {
		fmt.Printf("=== %s ===\n", word)
		wa := morph.Analyze(word)

		for _, sa := range wa.AnalysisResults {
			fmt.Printf("Morphemes: %d\n", len(sa.MorphemeDataList))
			for i, md := range sa.MorphemeDataList {
				fmt.Printf("  %d. Surface='%s', Morpheme='%s'\n", i, md.Surface, md.Morpheme.ID)
			}
			fmt.Printf("FormatString: %s\n", sa.FormatString())
		}
		fmt.Println()
	}
}
