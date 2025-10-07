//go:build demo
// +build demo

package main

import (
	"fmt"
	"github.com/kalaomer/zemberek-go/morphology/lexicon"
)

func main() {
	fmt.Println("=== Loading Master Dictionaries ===\n")

	// Load lexicon from master dictionaries
	lex, err := lexicon.LoadDefaultLexicon()
	if err != nil {
		fmt.Printf("Error loading lexicon: %v\n", err)
		return
	}

	fmt.Printf("✅ Lexicon loaded successfully!\n")
	fmt.Printf("   Total items: %d\n\n", lex.Size())

	// Test some words
	testWords := []string{
		"kitap",
		"gitmek",
		"güzel",
		"yarın",
		"okul",
		"bilgisayar",
		"telefon",
		"merhaba",
	}

	fmt.Println("=== Testing Word Lookup ===\n")
	for _, word := range testWords {
		items := lex.GetItems(word)
		if len(items) > 0 {
			fmt.Printf("✅ '%s' found - %d entry(ies)\n", word, len(items))
			for _, item := range items {
				fmt.Printf("   %s\n", item.String())
			}
		} else {
			fmt.Printf("❌ '%s' not found\n", word)
		}
	}

	fmt.Println("\n=== Sample Dictionary Items ===\n")
	allItems := lex.GetAllItems()
	if len(allItems) > 10 {
		for i := 0; i < 10; i++ {
			fmt.Printf("%d. %s\n", i+1, allItems[i].String())
		}
	}

	fmt.Println("\n=== Lexicon Statistics ===")
	fmt.Printf("Total items: %d\n", lex.Size())
	fmt.Printf("Unique lemmas: %d\n", len(lex.ItemMap))
}
