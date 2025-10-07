//go:build demo
// +build demo

package main

import (
	"fmt"

	"github.com/kalaomer/zemberek-go/core/turkish"
	"github.com/kalaomer/zemberek-go/morphology/lexicon"
	"github.com/kalaomer/zemberek-go/morphology/morphotactics"
)

func main() {
	fmt.Println("=== Zemberek-Go Morphology Test ===\n")

	// Create a simple lexicon with test words
	items := []*lexicon.DictionaryItem{
		lexicon.NewDictionaryItem("kalem", "kalem", turkish.Noun, turkish.NonePos, nil, "", 0),
		lexicon.NewDictionaryItem("kitap", "kitap", turkish.Noun, turkish.NonePos, nil, "", 0),
		lexicon.NewDictionaryItem("ev", "ev", turkish.Noun, turkish.NonePos, nil, "", 0),
	}

	// Create lexicon
	lex := lexicon.NewRootLexicon(items)
	fmt.Printf("Lexicon created with %d items\n", lex.Size())

	// Create morphotactics
	tm := morphotactics.NewTurkishMorphotactics(lex)
	fmt.Printf("Morphotactics initialized\n")

	// Get stem transitions
	stemTrans := tm.GetStemTransitions()
	fmt.Printf("Stem transitions manager ready\n\n")

	// Test words
	testWords := []string{"kalem", "kitap", "ev", "kalemin", "kitaplar", "evler"}

	for _, word := range testWords {
		fmt.Printf("Testing: %s\n", word)
		matches := stemTrans.GetPrefixMatches(word, false)
		fmt.Printf("  Found %d stem matches:\n", len(matches))
		for _, match := range matches {
			fmt.Printf("    - Surface: %s, Item: %s, State: %s\n",
				match.Surface,
				match.Item.Lemma,
				match.To.ID)
		}
		fmt.Println()
	}

	// Test morpheme definitions
	fmt.Println("=== Morpheme Definitions ===")
	fmt.Printf("Noun: %s (POS: %d)\n", morphotactics.Noun.ID, morphotactics.Noun.Pos)
	fmt.Printf("A3sg: %s\n", morphotactics.A3sg.ID)
	fmt.Printf("A3pl: %s\n", morphotactics.A3pl.ID)
	fmt.Printf("Pnon: %s\n", morphotactics.Pnon.ID)
	fmt.Printf("Nom: %s\n", morphotactics.Nom.ID)
	fmt.Printf("Dat: %s\n", morphotactics.Dat.ID)
	fmt.Printf("Acc: %s\n", morphotactics.Acc.ID)

	// Test state connections
	fmt.Println("\n=== State Connections ===")
	fmt.Printf("NounS outgoing transitions: %d\n", len(tm.NounS.Outgoing))
	fmt.Printf("A3sgS outgoing transitions: %d\n", len(tm.A3sgS.Outgoing))
	fmt.Printf("A3plS outgoing transitions: %d\n", len(tm.A3plS.Outgoing))
	fmt.Printf("PnonS outgoing transitions: %d\n", len(tm.PnonS.Outgoing))

	// Show some transition details
	fmt.Println("\n=== Sample Transitions from NounS ===")
	for i, trans := range tm.NounS.Outgoing {
		if i >= 2 {
			break
		}
		fmt.Printf("%d. %s\n", i+1, trans)
	}
}
