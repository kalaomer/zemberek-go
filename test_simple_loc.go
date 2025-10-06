package main

import (
	"fmt"
	"github.com/kalaomer/zemberek-go/morphology"
)

func main() {
	morph := morphology.CreateWithDefaults()
	
	// Test simplest form
	fmt.Println("Testing dosya + variants:")
	tests := []string{"dosya", "dosyam", "dosyan", "dosyası", "dosyamız", "dosyalarım"}
	
	for _, word := range tests {
		wa := morph.Analyze(word)
		fmt.Printf("  %-15s: ", word)
		if len(wa.AnalysisResults) > 0 {
			fmt.Printf("✅ %s\n", wa.AnalysisResults[0].FormatString())
		} else {
			fmt.Println("❌")
		}
	}
	
	fmt.Println("\nNow with +Loc:")
	tests2 := []string{"dosyada", "dosyamda", "dosyanda", "dosyasında"}
	
	for _, word := range tests2 {
		wa := morph.Analyze(word)
		fmt.Printf("  %-15s: ", word)
		if len(wa.AnalysisResults) > 0 {
			fmt.Printf("✅ %s\n", wa.AnalysisResults[0].FormatString())
		} else {
			fmt.Println("❌")
		}
	}
}
