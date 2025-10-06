package morphology

import (
	"fmt"
	"testing"
)

func TestMorphologicalAnalysis_BasicWord(t *testing.T) {
	// Create morphology instance
	morphology := CreateWithDefaults()

	// Test word from Java example
	word := "kutucuğumuz"
	result := morphology.Analyze(word)

	fmt.Printf("Word: %s\n", word)
	fmt.Printf("Analysis count: %d\n", len(result.AnalysisResults))

	for i, analysis := range result.AnalysisResults {
		fmt.Printf("\nAnalysis %d:\n", i+1)
		fmt.Printf("  Stem: %s\n", analysis.GetStem())
		fmt.Printf("  String: %s\n", analysis.String())
	}

	// At least one analysis should be found
	if len(result.AnalysisResults) == 0 {
		t.Errorf("Expected at least one analysis for '%s', got 0", word)
	}
}

func TestMorphologicalAnalysis_SimpleNoun(t *testing.T) {
	morphology := CreateWithDefaults()

	word := "kitap"
	result := morphology.Analyze(word)

	fmt.Printf("\nWord: %s\n", word)
	fmt.Printf("Analysis count: %d\n", len(result.AnalysisResults))

	for i, analysis := range result.AnalysisResults {
		fmt.Printf("\nAnalysis %d:\n", i+1)
		fmt.Printf("  Stem: %s\n", analysis.GetStem())
	}

	if len(result.AnalysisResults) == 0 {
		t.Errorf("Expected at least one analysis for '%s', got 0", word)
	}
}

func TestMorphologicalAnalysis_PluralNoun(t *testing.T) {
	morphology := CreateWithDefaults()

	word := "kitaplar"
	result := morphology.Analyze(word)

	fmt.Printf("\nWord: %s\n", word)
	fmt.Printf("Analysis count: %d\n", len(result.AnalysisResults))

	for i, analysis := range result.AnalysisResults {
		fmt.Printf("\nAnalysis %d:\n", i+1)
		fmt.Printf("  Stem: %s\n", analysis.GetStem())
		fmt.Printf("  String: %s\n", analysis.String())
	}

	if len(result.AnalysisResults) == 0 {
		t.Errorf("Expected at least one analysis for '%s', got 0", word)
	}
}
