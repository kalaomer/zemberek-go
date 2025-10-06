package main

import (
	"fmt"

	"github.com/kalaomer/zemberek-go/morphology"
)

func main() {
	tm := morphology.CreateWithDefaults()

	testWords := []string{"kitap", "kitaplar", "gidiyorum", "geldi"}

	for _, word := range testWords {
		wa := tm.Analyze(word)
		if len(wa.AnalysisResults) > 0 {
			result := wa.AnalysisResults[0]

			// Print morpheme count
			fmt.Printf("Word: %s\n", word)
			fmt.Printf("Morpheme count: %d\n", len(result.MorphemeDataList))
			fmt.Printf("Morphemes:\n")
			for i, md := range result.MorphemeDataList {
				fmt.Printf("  %d. surface='%s', morpheme='%s'\n", i, md.Surface, md.Morpheme.ID)
			}
			fmt.Printf("formatLong: %s\n", result.FormatString())
			fmt.Println()
		}
	}
}
