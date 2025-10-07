package morphology

import (
	"fmt"
	"testing"
)

func TestParseKutucuğ(t *testing.T) {
	morphology := CreateWithDefaults()

	// TODO: Remove skip when DIM voicing is supported in Go port
	t.Skip("dim suffix voicing support is not yet implemented")

	// Test if we can analyze the generated surface
	// First get what surface kutu+dim would generate
	word1 := "kutucuk" // Template >cI~k (LAST_VOICED k)
	word2 := "kutucuğ" // Template >cI!ğ (LAST_NOT_VOICED ğ)

	result1 := morphology.Analyze(word1)
	result2 := morphology.Analyze(word2)

	fmt.Printf("'%s' analysis: %d\n", word1, len(result1.AnalysisResults))
	for _, a := range result1.AnalysisResults {
		fmt.Printf("  %s\n", a.String())
	}

	fmt.Printf("\n'%s' analysis: %d\n", word2, len(result2.AnalysisResults))
	for _, a := range result2.AnalysisResults {
		fmt.Printf("  %s\n", a.String())
	}

	// Both should work
	if len(result1.AnalysisResults) == 0 {
		t.Error("kutucuk should have analysis")
	}
	if len(result2.AnalysisResults) == 0 {
		t.Error("kutucuğ should have analysis")
	}
}
