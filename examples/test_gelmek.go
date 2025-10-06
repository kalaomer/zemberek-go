package main

import (
	"fmt"

	"github.com/kalaomer/zemberek-go/morphology/lexicon"
)

func main() {
	items, _ := lexicon.LoadBinaryLexicon()

	for _, item := range items {
		if item.Lemma == "gelmek" {
			fmt.Printf("Found: Lemma=%s, Root=%s, POS=%v\n", item.Lemma, item.Root, item.PrimaryPos)
		}
	}
}
