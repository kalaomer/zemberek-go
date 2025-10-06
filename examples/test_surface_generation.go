package main

import (
	"fmt"

	"github.com/kalaomer/zemberek-go/morphology/morphotactics"
)

func main() {
	// Test tokenization of ">dI"
	tokens := morphotactics.TokenizeSuffixTemplate(">dI")
	fmt.Println("Tokens for '>dI':")
	for i, tok := range tokens {
		fmt.Printf("  %d. Type: %v, Value: %c (%d)\n", i, tok.Type, tok.Value, tok.Value)
	}

	fmt.Println()

	// Test tokenization of ">Iyor"
	tokens2 := morphotactics.TokenizeSuffixTemplate(">Iyor")
	fmt.Println("Tokens for '>Iyor':")
	for i, tok := range tokens2 {
		fmt.Printf("  %d. Type: %v, Value: %c (%d)\n", i, tok.Type, tok.Value, tok.Value)
	}

	// Check BUFFER_LETTER constant value
	fmt.Printf("\nBUFFER_LETTER value: %d\n", morphotactics.BUFFER_LETTER)
}
