package main

import (
	"fmt"
	"github.com/kalaomer/zemberek-go/morphology"
)

func main() {
	morph := morphology.CreateWithDefaults()
	
	// Test with "hukuk" which we know works
	words := []string{
		"hukuk",       // base
		"hukukta",     // +Loc should be: hukuk + ta (voiceless k → t)
		"hukuktan",    // +Abl should be: hukuk + tan
		"hukukun",     // +Gen
	}
	
	for _, word := range words {
		wa := morph.Analyze(word)
		fmt.Printf("%-15s: ", word)
		if len(wa.AnalysisResults) > 0 {
			fmt.Printf("✅ %s\n", wa.AnalysisResults[0].FormatString())
		} else {
			fmt.Println("❌ No analysis")
		}
	}
}
