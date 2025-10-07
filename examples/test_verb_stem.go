//go:build demo
// +build demo

package main

import (
	"fmt"

	"github.com/kalaomer/zemberek-go/core/turkish"
	"github.com/kalaomer/zemberek-go/morphology"
)

func main() {
	morph := morphology.CreateWithDefaults()

	gelmek := morph.Lexicon.GetItems("gelmek")[0]
	fmt.Printf("Item: %s, Root: %s, POS: %v\n", gelmek.Lemma, gelmek.Root, gelmek.PrimaryPos)

	// Check phonetic attributes
	phoneticAttrs := make(map[turkish.PhoneticAttribute]bool)
	turkish.Instance.CalculatePhoneticAttributes(gelmek.Root, phoneticAttrs, gelmek.Attributes)

	fmt.Println("\nPhonetic attributes:")
	for attr, val := range phoneticAttrs {
		if val {
			fmt.Printf("  %v\n", attr)
		}
	}

	// Get root state
	rootState := morph.Morphotactics.GetRootState(gelmek, phoneticAttrs)
	fmt.Printf("\nRoot state: %s\n", rootState.ID)

	// Get transitions from VerbRoot
	transitions := morph.Morphotactics.GetStemTransitions().GetTransitions(gelmek, phoneticAttrs)
	fmt.Printf("\nStem transitions count: %d\n", len(transitions))
	for i, trans := range transitions {
		if i < 5 {
			fmt.Printf("  %d. %s -> %s\n", i, trans.GetSurface(), trans.GetToState().ID)
		}
	}
}
