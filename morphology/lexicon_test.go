package morphology

import (
	"fmt"
	"testing"

	"github.com/kalaomer/zemberek-go/morphology/lexicon"
)

func TestLexiconLoading(t *testing.T) {
	lex, err := lexicon.LoadDefaultLexicon()
	if err != nil {
		t.Fatalf("Failed to load lexicon: %v", err)
	}

	fmt.Printf("Lexicon loaded: %d items\n", lex.Size())

	if lex.Size() == 0 {
		t.Error("Lexicon is empty!")
	}

	// Check for a common word
	allItems := lex.GetAllItems()
	found := false
	for _, item := range allItems {
		if item.Root == "kitap" {
			found = true
			fmt.Printf("Found 'kitap': %v\n", item.PrimaryPos)
			break
		}
	}

	if !found {
		t.Error("'kitap' not found in lexicon")
	}
}
