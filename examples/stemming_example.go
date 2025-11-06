//go:build demo
// +build demo

package main

import (
	"fmt"

	"github.com/kalaomer/zemberek-go/morphology"
)

func main() {
	fmt.Println("=== Zemberek-Go Stemming Example ===\n")

	// Create morphology instance
	morph := morphology.CreateWithDefaults()

	// Example 1: Simple text
	fmt.Println("Example 1: Simple Turkish text")
	text1 := "Kitapları okuyorum"
	tokens1 := morphology.StemTextWithPositions(text1, morph)

	fmt.Printf("Text: \"%s\"\n", text1)
	fmt.Println("Stems with positions:")
	for i, token := range tokens1 {
		fmt.Printf("  %d. Stem='%s', Original='%s', Bytes=[%d:%d]\n",
			i+1, token.Stem, token.Original, token.StartByte, token.EndByte)
	}
	fmt.Println()

	// Example 2: Sentence with punctuation
	fmt.Println("Example 2: Sentence with punctuation")
	text2 := "Bugün hava çok güzel, değil mi?"
	tokens2 := morphology.StemTextWithPositions(text2, morph)

	fmt.Printf("Text: \"%s\"\n", text2)
	fmt.Println("Stems with positions:")
	for i, token := range tokens2 {
		fmt.Printf("  %d. Stem='%s', Original='%s', Bytes=[%d:%d]\n",
			i+1, token.Stem, token.Original, token.StartByte, token.EndByte)
	}
	fmt.Println()

	// Example 3: UTF-8 byte offset validation
	fmt.Println("Example 3: UTF-8 byte offset validation")
	text3 := "Türkçe sözcükler"
	tokens3 := morphology.StemTextWithPositions(text3, morph)

	fmt.Printf("Text: \"%s\"\n", text3)
	fmt.Println("Extracting words using byte offsets:")
	for i, token := range tokens3 {
		extracted := text3[token.StartByte:token.EndByte]
		match := "✓"
		if extracted != token.Original {
			match = "✗"
		}
		fmt.Printf("  %d. %s text[%d:%d] = \"%s\" (expected \"%s\")\n",
			i+1, match, token.StartByte, token.EndByte, extracted, token.Original)
	}
	fmt.Println()

	// Example 4: Complex sentence
	fmt.Println("Example 4: Complex sentence with various forms")
	text4 := "Öğrencilerin kitaplarını okuyorlar ve yazıyorlar."
	tokens4 := morphology.StemTextWithPositions(text4, morph)

	fmt.Printf("Text: \"%s\"\n", text4)
	fmt.Println("Stems with positions:")
	for i, token := range tokens4 {
		fmt.Printf("  %d. Stem='%s', Original='%s', Bytes=[%d:%d]\n",
			i+1, token.Stem, token.Original, token.StartByte, token.EndByte)
	}
	fmt.Println()

	// Example 5: Without morphology (fallback mode)
	fmt.Println("Example 5: Without morphology (fallback mode)")
	text5 := "Test kelimeleri"
	tokens5 := morphology.StemTextWithPositions(text5, nil)

	fmt.Printf("Text: \"%s\"\n", text5)
	fmt.Println("Tokens (no stemming):")
	for i, token := range tokens5 {
		fmt.Printf("  %d. Stem='%s' (same as original), Bytes=[%d:%d]\n",
			i+1, token.Stem, token.StartByte, token.EndByte)
	}
	fmt.Println()

	// Example 6: FTS5 integration example
	fmt.Println("Example 6: How to use with SQLite FTS5")
	fmt.Println("In your C bridge code:")
	fmt.Println(`
  // Go function call:
  tokens := morphology.StemTextWithPositions(inputText, morphology)

  // For each token:
  for _, token := range tokens {
      // Call SQLite's xToken callback:
      rc := xToken(
          pCtx,
          0,                    // flags
          token.Stem,           // token text (stem)
          len(token.Stem),      // token length
          token.StartByte,      // start offset in original text
          token.EndByte,        // end offset in original text
      )
      if rc != SQLITE_OK {
          return rc
      }
  }
`)

	fmt.Println("\n=== Usage Summary ===")
	fmt.Println("Key Function:")
	fmt.Println("  morphology.StemTextWithPositions(text string, morph *TurkishMorphology) []StemToken")
	fmt.Println("\nReturns:")
	fmt.Println("  []StemToken with:")
	fmt.Println("    - Stem: stemmed form of the word")
	fmt.Println("    - Original: original word from text")
	fmt.Println("    - StartByte: UTF-8 byte offset start")
	fmt.Println("    - EndByte: UTF-8 byte offset end")
	fmt.Println("\nFeatures:")
	fmt.Println("  ✓ UTF-8 safe byte offset tracking")
	fmt.Println("  ✓ Morphological stemming via GetStem()")
	fmt.Println("  ✓ Handles Turkish characters correctly")
	fmt.Println("  ✓ Skips punctuation, returns only words")
	fmt.Println("  ✓ Fallback mode without morphology")
	fmt.Println("  ✓ Ready for SQLite FTS5 integration")
}
