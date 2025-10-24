package main

// #include <stdlib.h>
import "C"
import (
	"fmt"
	"strings"
	"sync"

	"github.com/kalaomer/zemberek-go/morphology"
	"github.com/kalaomer/zemberek-go/morphology/lexicon"
	"github.com/kalaomer/zemberek-go/normalization"
)

var (
	morphologyInstance *morphology.TurkishMorphology
	normalizerInstance *normalization.TurkishSentenceNormalizer
	initOnce           sync.Once
	initError          error
)

// Initialize initializes the morphology and normalizer instances
func initialize() {
	initOnce.Do(func() {
		// Initialize morphology
		items, err := lexicon.LoadBinaryLexicon()
		if err != nil {
			initError = fmt.Errorf("failed to load lexicon: %v", err)
			return
		}

		lex := lexicon.NewRootLexicon(items)
		morphologyInstance = morphology.NewBuilder(lex).Build()

		// Initialize normalizer with empty stem words for now
		normalizerInstance, err = normalization.NewTurkishSentenceNormalizer([]string{}, "resources/normalization")
		if err != nil {
			// Normalizer is optional, continue without it
			normalizerInstance = nil
		}
	})
}

// NormalizeTurkish normalizes informal Turkish text to formal Turkish
//
//export NormalizeTurkish
func NormalizeTurkish(text *C.char) *C.char {
	initialize()
	if initError != nil {
		return C.CString(fmt.Sprintf("Error: %v", initError))
	}

	input := C.GoString(text)
	if input == "" {
		return C.CString("")
	}

	if normalizerInstance == nil {
		return C.CString(input) // Return as-is if normalizer not available
	}

	normalized := normalizerInstance.Normalize(input)
	return C.CString(normalized)
}

// AnalyzeTurkish performs morphological analysis on Turkish text
//
//export AnalyzeTurkish
func AnalyzeTurkish(word *C.char) *C.char {
	initialize()
	if initError != nil {
		return C.CString(fmt.Sprintf("Error: %v", initError))
	}

	input := C.GoString(word)
	if input == "" {
		return C.CString("")
	}

	wa := morphologyInstance.Analyze(input)
	if len(wa.AnalysisResults) == 0 {
		return C.CString(fmt.Sprintf("No analysis found for: %s", input))
	}

	// Format results
	var results []string
	for _, sa := range wa.AnalysisResults {
		results = append(results, sa.FormatString())
	}

	return C.CString(strings.Join(results, " | "))
}

// StemTurkish extracts stems from Turkish words
//
//export StemTurkish
func StemTurkish(word *C.char) *C.char {
	initialize()
	if initError != nil {
		return C.CString(fmt.Sprintf("Error: %v", initError))
	}

	input := C.GoString(word)
	if input == "" {
		return C.CString("")
	}

	wa := morphologyInstance.Analyze(input)
	if len(wa.AnalysisResults) == 0 {
		return C.CString(input) // Return original if no analysis
	}

	// Get stem from first analysis
	stem := wa.AnalysisResults[0].GetStem()
	return C.CString(stem)
}

// HasTurkishAnalysis checks if a Turkish word has morphological analysis
//
//export HasTurkishAnalysis
func HasTurkishAnalysis(word *C.char) C.int {
	initialize()
	if initError != nil {
		return 0
	}

	input := C.GoString(word)
	if input == "" {
		return 0
	}

	hasAnalysis := morphologyInstance.HasAnalysis(input)
	if hasAnalysis {
		return 1
	}
	return 0
}

func main() {
	// Required for CGO build
}
