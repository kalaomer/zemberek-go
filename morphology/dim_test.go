package morphology

import (
	"fmt"
	"testing"
)

func TestDiminutiveSuffix(t *testing.T) {
	morphology := CreateWithDefaults()

	tests := []struct {
		word string
	}{
		{"kutu"},
		{"kutucuk"},
		{"kutucuğ"},
		{"kutucuğumuz"},
	}

	for _, tt := range tests {
		result := morphology.Analyze(tt.word)
		fmt.Printf("\nWord: %s\n", tt.word)
		fmt.Printf("Analysis count: %d\n", len(result.AnalysisResults))
		for i, analysis := range result.AnalysisResults {
			fmt.Printf("  [%d] Stem: %s, String: %s\n", i+1, analysis.GetStem(), analysis.String())
		}
	}
}
