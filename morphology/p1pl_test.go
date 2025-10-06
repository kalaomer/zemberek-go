package morphology

import (
	"fmt"
	"testing"
)

func TestP1plSuffix(t *testing.T) {
	morphology := CreateWithDefaults()

	words := []string{
		"evimiz",      // ev+ImIz
		"kutumuz",     // kutu+ImIz → kutumuz (u is back+rounded, I→u, I→u)
		"kutucuğumuz", // kutu+cuğ+ImIz
	}

	for _, word := range words {
		result := morphology.Analyze(word)
		fmt.Printf("'%s' → %d analyses\n", word, len(result.AnalysisResults))
		for _, a := range result.AnalysisResults {
			fmt.Printf("  %s\n", a.String())
		}
	}
}
