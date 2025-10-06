package main

import (
	"fmt"

	"github.com/kalaomer/zemberek-go/morphology/morphotactics"
)

func main() {
	// Test tokenization of "Iyor"
	tokens := morphotactics.TokenizeSuffixTemplate("Iyor")
	fmt.Println("Tokens for 'Iyor':")
	for i, tok := range tokens {
		fmt.Printf("  %d. Type: %v, Value: %c (%d)\n", i, tok.Type, tok.Value, tok.Value)
	}
}
