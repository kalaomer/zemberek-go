//go:build demo
// +build demo

package main

import (
	"fmt"

	"github.com/kalaomer/zemberek-go/morphology/lexicon"
)

func main() {
	items, err := lexicon.LoadBinaryLexicon()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Total dictionary items: %d\n", len(items))
}
