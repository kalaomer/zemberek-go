package main

import (
	"fmt"
	"github.com/kalaomer/zemberek-go/morphology"
)

func main() {
	morph := morphology.CreateWithDefaults()
	
	// Test simple locative first
	fmt.Println("Testing: dosyada")
	wa1 := morph.Analyze("dosyada")
	if len(wa1.AnalysisResults) > 0 {
		fmt.Printf("✅ dosyada: %s\n", wa1.AnalysisResults[0].FormatString())
	} else {
		fmt.Println("❌ dosyada: No analysis")
	}
	
	// Test with -ki
	fmt.Println("\nTesting: dosyadaki")
	wa2 := morph.Analyze("dosyadaki")
	if len(wa2.AnalysisResults) > 0 {
		fmt.Printf("✅ dosyadaki: %s\n", wa2.AnalysisResults[0].FormatString())
	} else {
		fmt.Println("❌ dosyadaki: No analysis")
	}
}
