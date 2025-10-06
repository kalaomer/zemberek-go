package lexicon

import (
	"fmt"
	"testing"
)

func TestDictionaryLoading(t *testing.T) {
	paths := GetDefaultDictionaryPaths()
	fmt.Printf("Dictionary paths: %v\n", paths)

	for _, path := range paths {
		fmt.Printf("\nTrying to load: %s\n", path)
		items, err := LoadMasterDictionary(path)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else {
			fmt.Printf("  Loaded: %d items\n", len(items))
			if len(items) > 0 {
				fmt.Printf("  First item: %s (%v)\n", items[0].Root, items[0].PrimaryPos)
			}
		}
	}

	fmt.Println("\nLoading all dictionaries...")
	items, err := LoadAllDefaultDictionaries()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	fmt.Printf("Total items loaded: %d\n", len(items))

	lex := NewRootLexicon(items)
	fmt.Printf("Lexicon size: %d\n", lex.Size())
}
