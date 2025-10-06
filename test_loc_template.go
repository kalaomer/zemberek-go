package main

import (
	"fmt"
	"github.com/kalaomer/zemberek-go/morphology"
)

func main() {
	morph := morphology.CreateWithDefaults()
	
	// Test different words with locative
	words := []string{
		"evde",      // ev + de
		"evda",      // wrong
		"masada",    // masa + da
		"masade",    // wrong
		"kitapta",   // kitap + ta (voiceless p)
		"kitapda",   // wrong
		"dosyada",   // dosya + da
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
