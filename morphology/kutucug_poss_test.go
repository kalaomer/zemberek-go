package morphology

import (
	"fmt"
	"testing"
)

func TestKutucugPossession(t *testing.T) {
	morphology := CreateWithDefaults()

	words := []string{
		"kutucuğum",    // P1sg: Im
		"kutucuğun",    // P2sg: In
		"kutucuğu",     // P3sg: sI
		"kutucuğumuz",  // P1pl: ImIz
	}

	for _, word := range words {
		result := morphology.Analyze(word)
		fmt.Printf("'%s' → %d analyses\n", word, len(result.AnalysisResults))
		if len(result.AnalysisResults) > 0 {
			for i, analysis := range result.AnalysisResults {
				fmt.Printf("  %d: %s\n", i+1, analysis.String())
			}
		}
	}
}
