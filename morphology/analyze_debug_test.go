package morphology

import (
	"fmt"
	"testing"
)

func TestAnalyzeKutucuğDebug(t *testing.T) {
	morphology := CreateWithDefaults()

	// Simple test: can we analyze "kutu"?
	result1 := morphology.Analyze("kutu")
	fmt.Printf("'kutu' → %d analyses\n", len(result1.AnalysisResults))

	// Can we analyze "kutucuk"? (works)
	result2 := morphology.Analyze("kutucuk")
	fmt.Printf("'kutucuk' → %d analyses\n", len(result2.AnalysisResults))
	if len(result2.AnalysisResults) > 0 {
		fmt.Printf("  %s\n", result2.AnalysisResults[0].String())
	}

	// Can we analyze "kutucuğ"? (doesn't work)
	result3 := morphology.Analyze("kutucuğ")
	fmt.Printf("'kutucuğ' → %d analyses\n", len(result3.AnalysisResults))

	// Let's manually check if the analyzer sees the transitions
	// Get analyzer
	analyzer := morphology.Analyzer

	// Get stem candidates for "kutucuğ"
	candidates := analyzer.StemTransitions.GetPrefixMatches("kutucuğ", false)
	fmt.Printf("\nStem candidates for 'kutucuğ': %d\n", len(candidates))
	for _, c := range candidates {
		fmt.Printf("  '%s' -> %s\n", c.Surface, c.To.ID)
	}

	// The stem should be "kutu" with tail "cuğ"
	// Let's check if nom_ST has outgoing to dim_S with template ">cI!ğ"
}
