package main

import (
	"fmt"
	"github.com/kalaomer/zemberek-go/morphology"
)

func main() {
	morph := morphology.CreateWithDefaults()
	
	// Test with words we know exist
	words := []string{
		"dosyada",
		"hukukta",
		"davada",
		"kararına",
		"dairesinde",
	}
	
	for _, word := range words {
		wa := morph.Analyze(word)
		fmt.Printf("%-15s: ", word)
		if len(wa.AnalysisResults) > 0 {
			for i, ar := range wa.AnalysisResults {
				if i > 0 {
					fmt.Printf("                 ")
				}
				fmt.Printf("%d) %s\n", i+1, ar.FormatString())
			}
		} else {
			fmt.Println("❌ No analysis")
		}
	}
}
