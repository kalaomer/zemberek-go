package main

import (
	"fmt"
	"github.com/kalaomer/zemberek-go/morphology"
)

func main() {
	morph := morphology.CreateWithDefaults()
	
	words := []string{"dosya", "dosyası", "dosyada", "dosyadaki"}
	
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
