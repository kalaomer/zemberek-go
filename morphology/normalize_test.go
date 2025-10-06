package morphology

import (
	"fmt"
	"testing"
)

func TestNormalization(t *testing.T) {
	morphology := CreateWithDefaults()

	words := []string{"kutucuk", "kutucuğ", "kutucuğumuz"}

	for _, word := range words {
		normalized := morphology.NormalizeForAnalysis(word)
		fmt.Printf("'%s' → '%s' (same: %v)\n", word, normalized, word == normalized)

		if word != normalized {
			// Show byte representation
			fmt.Printf("  Original bytes: %v\n", []byte(word))
			fmt.Printf("  Normalized bytes: %v\n", []byte(normalized))
		}
	}
}
