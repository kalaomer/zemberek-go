package analysis

import (
	"fmt"
	"strings"
)

// WordAnalysis represents analysis results for a word
type WordAnalysis struct {
	Input            string
	AnalysisResults  []*SingleAnalysis
	NormalizedInput  string
}

// EmptyInputResult is a singleton for empty input
var EmptyInputResult = &WordAnalysis{
	Input:           "",
	AnalysisResults: make([]*SingleAnalysis, 0),
}

// NewWordAnalysis creates a new WordAnalysis
func NewWordAnalysis(input string, analysisResults []*SingleAnalysis, normalizedInput string) *WordAnalysis {
	normalized := input
	if normalizedInput != "" {
		normalized = normalizedInput
	}

	return &WordAnalysis{
		Input:            input,
		AnalysisResults:  analysisResults,
		NormalizedInput:  normalized,
	}
}

// IsCorrect returns true if the analysis is correct (not unknown)
func (wa *WordAnalysis) IsCorrect() bool {
	return len(wa.AnalysisResults) > 0 && !wa.AnalysisResults[0].IsUnknown()
}

// String returns string representation
func (wa *WordAnalysis) String() string {
	var analyses []string
	for _, a := range wa.AnalysisResults {
		analyses = append(analyses, a.String())
	}

	return fmt.Sprintf("WordAnalysis{input='%s', normalizedInput='%s', analysisResults=%s}",
		wa.Input, wa.NormalizedInput, strings.Join(analyses, " "))
}

// GetAnalysisCount returns the number of analyses
func (wa *WordAnalysis) GetAnalysisCount() int {
	return len(wa.AnalysisResults)
}
