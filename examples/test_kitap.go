//go:build demo
// +build demo

package main

import (
	"fmt"

	"github.com/kalaomer/zemberek-go/morphology"
)

func main() {
	morph := morphology.CreateWithDefaults()

	words := []string{"kitap", "kitaplar", "gidiyorum", "geldi"}

	for _, word := range words {
		fmt.Printf("Word: %s\n", word)
		results := morph.Analyze(word)

		for _, analysis := range results.AnalysisResults {
			if analysis.Item.SecondaryPos.GetStringForm() != "ProperNoun" {
				fmt.Printf("Morpheme count: %d\n", len(analysis.MorphemeDataList))
				fmt.Println("Morphemes:")
				for i, md := range analysis.MorphemeDataList {
					fmt.Printf("  %d. surface='%s', morpheme='%s'\n", i, md.Surface, md.Morpheme.ID)
				}
				fmt.Printf("formatLong: %s\n", analysis.FormatString())
				break
			}
		}
		fmt.Println()
	}
}
