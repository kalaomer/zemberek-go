package morphology

import (
	"fmt"
	"testing"

	"github.com/kalaomer/zemberek-go/morphology/lexicon"
)

func TestDebugStemTransitions(t *testing.T) {
	morphology := CreateWithDefaults()

	word := "kitap"

	// Check lexicon size
	fmt.Printf("Lexicon size: %d\n", morphology.Lexicon.Size())

	// Check if word is in lexicon
	allItems := morphology.Lexicon.GetAllItems()
	foundItems := make([]*lexicon.DictionaryItem, 0)
	for _, item := range allItems {
		if item.Root == word {
			foundItems = append(foundItems, item)
		}
	}
	fmt.Printf("Items matching '%s': %d\n", word, len(foundItems))
	for _, item := range foundItems {
		fmt.Printf("  - %s (%v)\n", item.Root, item.PrimaryPos)
	}

	// Check stem transitions
	stemTrans := morphology.Morphotactics.GetStemTransitions()
	candidates := stemTrans.GetPrefixMatches(word, false)
	fmt.Printf("\nStem candidates for '%s': %d\n", word, len(candidates))
	for _, cand := range candidates {
		fmt.Printf("  - Surface: %s, Item: %s\n", cand.Surface, cand.Item.Root)
	}

	// Try analysis
	result := morphology.Analyze(word)
	fmt.Printf("\nAnalysis results: %d\n", len(result.AnalysisResults))
}
